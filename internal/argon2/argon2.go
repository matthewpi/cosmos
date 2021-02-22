//
// Copyright (c) 2021 Matthew Penner
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

// Package argon2 provides an easy way to use the argon2 key derivation function.
package argon2 // import "github.com/matthewpi/cosmos/internal/argon2"

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/blake2b"
)

// Argon2 Settings
var (
	Memory      uint32 = 65536
	Iterations  uint32 = 3
	Parallelism uint8  = 2
	SaltLength  uint32 = 16
	KeyLength   uint32 = 32
)

var (
	// ErrIncompatibleVersion is an incompatible version error.
	ErrIncompatibleVersion = errors.New("argon2: incompatible version")

	// ErrInvalidHash is an invalid hash error.
	ErrInvalidHash = errors.New("argon2: invalid hash")

	// ErrFailedVerify is a failed verification error.
	ErrFailedVerify = errors.New("argon2: failed verify")
)

// Hash hashes the input using the argon2id algorithm.
func Hash(input []byte) (string, error) {
	// Generate a random salt
	salt, err := generateRandomBytes(SaltLength)
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	// Get the argon2 id key.
	key, err := idKey(input, salt)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	b.WriteString("$argon2id$v=")
	b.WriteString(strconv.FormatInt(argon2.Version, 10))
	b.WriteString("$m=")
	b.WriteString(strconv.FormatUint(uint64(Memory), 10))
	b.WriteString(",t=")
	b.WriteString(strconv.FormatUint(uint64(Iterations), 10))
	b.WriteString(",p=")
	b.WriteString(strconv.FormatUint(uint64(Parallelism), 10))
	b.WriteString("$")
	b.WriteString(encodeBase64(salt))
	b.WriteString("$")
	b.WriteString(encodeBase64(key))

	// Verify that input against the hash, if this fails than we are fucked
	if err := Verify(input, b.String()); err != nil {
		return "", errors.Wrap(err, "failed to verify input against hash")
	}

	return b.String(), nil
}

// Verify verifies the input against a hash.
func Verify(input []byte, encodedHash string) error {
	// Decode the hash
	salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return errors.Wrap(err, "failed to decode hash")
	}

	// Get the argon2 id key
	comparisonHash, err := idKey(input, salt)
	if err != nil {
		return err
	}

	// Compare the two hashes
	if subtle.ConstantTimeCompare(hash, comparisonHash) != 1 {
		return ErrFailedVerify
	}

	return nil
}

// idKey gets the argon2 id key.
func idKey(input, salt []byte) ([]byte, error) {
	inputHash, err := blake2(input)
	if err != nil {
		return nil, err
	}

	return argon2.IDKey(inputHash, salt, Iterations, Memory, Parallelism, KeyLength), nil
}

// decodeHash decodes the argon2 hash from a string.
func decodeHash(encodedHash string) ([]byte, []byte, error) {
	values := strings.Split(encodedHash, "$")
	if len(values) != 6 {
		return nil, nil, ErrInvalidHash
	}

	var version int
	if _, err := fmt.Sscanf(values[2], "v=%d", &version); err != nil {
		return nil, nil, err
	}

	if version != argon2.Version {
		return nil, nil, ErrIncompatibleVersion
	}

	salt, err := decodeBase64(values[4])
	if err != nil {
		return nil, nil, err
	}

	hash, err := decodeBase64(values[5])
	if err != nil {
		return nil, nil, err
	}

	return salt, hash, nil
}

// blake2 runs the BLAKE2b_384 algorithm on an input.
func blake2(input []byte) ([]byte, error) {
	h, err := blake2b.New384(nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create blake2 hash")
	}

	_, err = h.Write(input)
	if err != nil {
		return nil, err
	}

	sum := h.Sum(nil)

	dst := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(dst, sum)

	return dst, nil
}

// encodeBase64 encodes a string using base64.
func encodeBase64(src []byte) string {
	return base64.RawStdEncoding.EncodeToString(src)
}

// decodeBase64 decodes a base64 string.
func decodeBase64(src string) ([]byte, error) {
	return base64.RawStdEncoding.DecodeString(src)
}

// generateRandomBytes generates crypto-secure random bytes.
func generateRandomBytes(n uint32) ([]byte, error) {
	bytes := make([]byte, n)

	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
