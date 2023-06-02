// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package container

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/awnumar/memguard"

	"zntr.io/harp/v2/pkg/container"
	"zntr.io/harp/v2/pkg/container/seal"
	sealv1 "zntr.io/harp/v2/pkg/container/seal/v1"
	sealv2 "zntr.io/harp/v2/pkg/container/seal/v2"
	"zntr.io/harp/v2/pkg/sdk/types"
	"zntr.io/harp/v2/pkg/tasks"
)

// SealTask implements secret container sealing task.
type SealTask struct {
	ContainerReader          tasks.ReaderProvider
	SealedContainerWriter    tasks.WriterProvider
	OutputWriter             tasks.WriterProvider
	PeerPublicKeys           []string
	DCKDMasterKey            string
	DCKDTarget               string
	JSONOutput               bool
	DisableContainerIdentity bool
	SealVersion              uint
	PreSharedKey             *memguard.LockedBuffer
}

// Run the task.
//
//nolint:funlen,gocyclo // to refactor
func (t *SealTask) Run(ctx context.Context) error {
	// Check arguments
	if types.IsNil(t.ContainerReader) {
		return errors.New("unable to run task with a nil containerReader provider")
	}
	if types.IsNil(t.SealedContainerWriter) {
		return errors.New("unable to run task with a nil sealedContainerWriter provider")
	}
	if types.IsNil(t.OutputWriter) {
		return errors.New("unable to run task with a nil outputWriter provider")
	}
	if len(t.PeerPublicKeys) == 0 {
		return errors.New("at least one public key must be provided for recovery")
	}

	// Create input reader
	reader, err := t.ContainerReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to open input reader: %w", err)
	}

	// Load input container
	in, err := container.Load(reader)
	if err != nil {
		return fmt.Errorf("unable to read input container: %w", err)
	}

	var containerKey string
	if !t.DisableContainerIdentity {
		opts := []seal.GenerateOption{}

		// Check container sealing master key usage
		if t.DCKDMasterKey != "" {
			// Process target
			if t.DCKDTarget == "" {
				return errors.New("target flag (string) is mandatory for key derivation")
			}

			// Decode master key
			masterKeyRaw, errDecode := base64.RawURLEncoding.DecodeString(t.DCKDMasterKey)
			if errDecode != nil {
				return fmt.Errorf("unable to decode master key: %w", errDecode)
			}

			// Enable deterministic container key generation
			opts = append(opts, seal.WithDeterministicKey(memguard.NewBufferFromBytes(masterKeyRaw), t.DCKDTarget))
		}

		// Initialize seal strategy
		var ss seal.Strategy
		switch t.SealVersion {
		case 1:
			ss = sealv1.New()
		case 2:
			ss = sealv2.New()
		default:
			ss = sealv1.New()
		}

		// Generate container key
		containerPublicKey, containerSecretKey, errGenerate := ss.GenerateKey(opts...)
		if errGenerate != nil {
			return fmt.Errorf("unable to generate container key: %w", errGenerate)
		}

		// Append to identities
		t.PeerPublicKeys = append(t.PeerPublicKeys, containerPublicKey)

		// Assign container key
		containerKey = containerSecretKey
	}

	// Seal options
	sopts := []container.Option{
		container.WithPeerPublicKeys(t.PeerPublicKeys),
	}

	// Process pre-shared key
	if t.PreSharedKey != nil {
		// Try to decode preshared key
		psk, errDecode := base64.RawURLEncoding.DecodeString(t.PreSharedKey.String())
		if errDecode != nil {
			return fmt.Errorf("unable to decode pre-shared key: %w", errDecode)
		}
		sopts = append(sopts, container.WithPreSharedKey(memguard.NewBufferFromBytes(psk)))
		t.PreSharedKey.Destroy()
	}

	// Seal the container
	sealedContainer, err := container.Seal(rand.Reader, in, sopts...)
	if err != nil {
		return fmt.Errorf("unable to seal container: %w", err)
	}

	// Open output file
	writer, err := t.SealedContainerWriter(ctx)
	if err != nil {
		return fmt.Errorf("unable to create output bundle: %w", err)
	}

	// Dump to writer
	if err = container.Dump(writer, sealedContainer); err != nil {
		return fmt.Errorf("unable to write sealed container: %w", err)
	}

	if !t.DisableContainerIdentity {
		// Get output writer
		outputWriter, err := t.OutputWriter(ctx)
		if err != nil {
			return fmt.Errorf("unable to retrieve output writer: %w", err)
		}

		// Display as json
		if t.JSONOutput {
			if err := json.NewEncoder(outputWriter).Encode(map[string]interface{}{
				"container_key": containerKey,
			}); err != nil {
				return fmt.Errorf("unable to display as json: %w", err)
			}
		} else {
			// Display container key
			if _, err := fmt.Fprintf(outputWriter, "Container key : %s\n", containerKey); err != nil {
				return fmt.Errorf("unable to display result: %w", err)
			}
		}
	}

	// No error
	return nil
}
