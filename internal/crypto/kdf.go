package crypto

import (
	"io"

	"golang.org/x/crypto/argon2"
)

// GenerateSalt generates a cryptographically secure 32-byte salt using
// the operating system's random number generator (crypto/rand).
//
// The salt should be stored alongside the encrypted data (in the plaintext
// header of the vault file) and passed to DeriveKey() when opening the vault.
// A new salt should be generated only when:
//   - Creating a new vault
//   - Changing the master password
//
// Returns:
//   - []byte: A 32-byte salt
//   - error: ErrInsufficientEntropy if crypto/rand cannot provide enough random bytes
//
// Example:
//
//	salt, err := crypto.GenerateSalt()
//	if err != nil {
//	    return fmt.Errorf("failed to generate salt: %w", err)
//	}
//	// Store salt in vault file header
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 32)

	// Use io.ReadFull to ensure we get exactly 32 bytes.
	// rand.Read() can return fewer bytes if the entropy pool is low.
	if _, err := io.ReadFull(entropyReader, salt); err != nil {
		return nil, ErrInsufficientEntropy
	}

	return salt, nil
}

// DeriveKey derives a cryptographic key from a password and salt using Argon2id.
//
// Argon2id is a memory-hard key derivation function that provides resistance
// against GPU-based brute-force attacks. It is the current best practice for
// password-based key derivation (RFC 9106).
//
// The derived key should be used with AES-256-GCM encryption and must be
// zeroed after use by the caller using ZeroBytes().
//
// Parameters:
//   - password: The user's master password as a []byte. Must not be nil or empty.
//   - salt: A 32-byte salt generated with GenerateSalt(). Must not be nil or empty.
//   - params: Argon2id parameters (time, memory, threads, keyLen).
//
// Returns:
//   - []byte: The derived key (length specified in params.KeyLen)
//   - error: ErrInvalidParams if password or salt is nil/empty
//
// CRITICAL: The caller is responsible for zeroing the returned key after use:
//
//	key, err := crypto.DeriveKey(password, salt, params)
//	if err != nil {
//	    return err
//	}
//	defer crypto.ZeroBytes(key)
//
// Note: Argon2id memory parameter is in KiB. For 256 MiB, use 262144 (not 256).
func DeriveKey(password, salt []byte, params ArgonParams) ([]byte, error) {
	// Validate inputs
	if len(password) == 0 {
		return nil, ErrInvalidParams
	}

	if len(salt) == 0 {
		return nil, ErrInvalidParams
	}

	// Call Argon2id key derivation.
	// The argon2.IDKey function signature is:
	//   func IDKey(password, salt []byte, time, memory uint32, threads uint8, keyLen uint32) []byte
	//
	// This is the RFC 9106 reference implementation.
	key := argon2.IDKey(password, salt, params.Time, params.Memory, params.Threads, params.KeyLen)

	// Caller's responsibility to zero key after use with ZeroBytes()
	return key, nil
}
