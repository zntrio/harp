// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package cmdutil

import (
	"fmt"
	"os"
	"syscall"

	"github.com/awnumar/memguard"

	"zntr.io/harp/v2/pkg/sdk/security"

	"golang.org/x/term"
)

// ReadSecret reads password from Stdin and returns a lockedbuffer.
func ReadSecret(prompt string, confirmation bool) (*memguard.LockedBuffer, error) {
	var (
		err             error
		password        []byte
		passwordConfirm []byte
	)
	defer memguard.WipeBytes(password)
	defer memguard.WipeBytes(passwordConfirm)

	// Ask to password
	fmt.Fprintf(os.Stdout, "%s: ", prompt)
	//nolint:unconvert // stdin doesn't share same type on each platform
	password, err = term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, fmt.Errorf("unable to read secret")
	}

	fmt.Fprint(os.Stdout, "\n")

	// Check if confirmation is requested
	if !confirmation {
		// Return locked buffer
		return memguard.NewBufferFromBytes(password), nil
	}

	fmt.Fprintf(os.Stdout, "%s (confirmation): ", prompt)
	//nolint:unconvert // stdin doesn't share same type on each platform
	passwordConfirm, err = term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, fmt.Errorf("unable to read secret confirmation")
	}

	fmt.Fprint(os.Stdout, "\n")

	// Compare if equal
	if !security.SecureCompare(password, passwordConfirm) {
		return nil, fmt.Errorf("passphrase doesn't match")
	}

	// Return locked buffer
	return memguard.NewBufferFromBytes(password), nil
}
