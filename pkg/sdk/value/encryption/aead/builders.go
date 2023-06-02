// SPDX-FileCopyrightText: 2019 Elasticsearch B.V.
// SPDX-FileCopyrightText: 2019-2023 Thibault NORMAND <me@zenithar.org>
//
// SPDX-License-Identifier: Apache-2.0 AND MIT

package aead

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"strings"

	miscreant "github.com/miscreant/miscreant.go"

	"zntr.io/harp/v2/build/fips"
	"zntr.io/harp/v2/pkg/sdk/value"
	"zntr.io/harp/v2/pkg/sdk/value/encryption"

	"golang.org/x/crypto/chacha20poly1305"
)

var (
	aesgcmPrefix     = "aes-gcm"
	aespmacsivPrefix = "aes-pmac-siv"
	aessivPrefix     = "aes-siv"
	chachaPrefix     = "chacha"
	xchachaPrefix    = "xchacha"
)

func init() {
	encryption.Register(aesgcmPrefix, AESGCM)

	if !fips.Enabled() {
		encryption.Register(aespmacsivPrefix, AESPMACSIV)
		encryption.Register(aessivPrefix, AESSIV)
		encryption.Register(chachaPrefix, Chacha20Poly1305)
		encryption.Register(xchachaPrefix, XChacha20Poly1305)
	}
}

// AESGCM returns an AES-GCM value transformer instance.
func AESGCM(key string) (value.Transformer, error) {
	// Remove the prefix
	key = strings.TrimPrefix(key, "aes-gcm:")

	// Decode key
	k, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("aes: unable to decode transformer key: %w", err)
	}

	// Check key length
	switch len(k) {
	case 16, 24, 32:
	default:
		return nil, fmt.Errorf("aes: invalid key length, use 16 bytes (AES128) or 24 bytes (AES192) or 32 bytes (AES256)")
	}

	// Create AES block cipher
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, fmt.Errorf("aes: unable to initialize block cipher: %w", err)
	}

	// Initialize AEAD cipher chain
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("aes: unable to initialize aead chain : %w", err)
	}

	// Return transformer
	return &aeadTransformer{
		aead: aead,
	}, nil
}

// AESSIV returns an AES-SIV/AES-CMAC-SIV value transformer instance.
func AESSIV(key string) (value.Transformer, error) {
	// Remove the prefix
	key = strings.TrimPrefix(key, "aes-siv:")

	// Decode key
	k, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("aes: unable to decode transformer key: %w", err)
	}
	if l := len(k); l != 64 {
		return nil, fmt.Errorf("aes: invalid secret key length (%d)", l)
	}

	// Initialize AEAD
	aead, err := miscreant.NewAEAD("AES-SIV", k, 16)
	if err != nil {
		return nil, fmt.Errorf("aes: unable to initialize aes-pmac-siv: %w", err)
	}

	// Return transformer
	return &aeadTransformer{
		aead: aead,
	}, nil
}

// AESPMACSIV returns an AES-PMAC-SIV value transformer instance.
func AESPMACSIV(key string) (value.Transformer, error) {
	// Remove the prefix
	key = strings.TrimPrefix(key, "aes-pmac-siv:")

	// Decode key
	k, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("aes: unable to decode transformer key: %w", err)
	}
	if l := len(k); l != 64 {
		return nil, fmt.Errorf("aes: invalid secret key length (%d)", l)
	}

	// Initialize AEAD
	aead, err := miscreant.NewAEAD("AES-PMAC-SIV", k, 16)
	if err != nil {
		return nil, fmt.Errorf("aes: unable to initialize aes-pmac-siv: %w", err)
	}

	// Return transformer
	return &aeadTransformer{
		aead: aead,
	}, nil
}

// Chacha20Poly1305 returns an ChaCha20Poly1305 value transformer instance.
func Chacha20Poly1305(key string) (value.Transformer, error) {
	// Remove the prefix
	key = strings.TrimPrefix(key, "chacha:")

	// Decode key
	k, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("chacha: unable to decode transformer key: %w", err)
	}
	if l := len(k); l != keyLength {
		return nil, fmt.Errorf("chacha: invalid secret key length (%d)", l)
	}

	// Create Chacha20-Poly1305 aead cipher
	aead, err := chacha20poly1305.New(k)
	if err != nil {
		return nil, fmt.Errorf("chacha: unable to initialize chacha cipher: %w", err)
	}

	// Return transformer
	return &aeadTransformer{
		aead: aead,
	}, nil
}

// XChacha20Poly1305 returns an XChaCha20Poly1305 value transformer instance.
func XChacha20Poly1305(key string) (value.Transformer, error) {
	// Remove the prefix
	key = strings.TrimPrefix(key, "xchacha:")

	// Decode key
	k, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("xchacha: unable to decode transformer key: %w", err)
	}
	if l := len(k); l != keyLength {
		return nil, fmt.Errorf("xchacha: invalid secret key length (%d)", l)
	}

	// Create Chacha20-Poly1305 aead cipher
	aead, err := chacha20poly1305.NewX(k)
	if err != nil {
		return nil, fmt.Errorf("xchacha: unable to initialize chacha cipher: %w", err)
	}

	// Return transformer
	return &aeadTransformer{
		aead: aead,
	}, nil
}
