package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// Encrypt encrypts plaintext using AES-256-GCM with a unique nonce.
//
// This function generates a fresh 12-byte nonce using crypto/rand for every
// call, ensuring that the same plaintext encrypted twice with the same key
// produces distinct ciphertexts. Nonce uniqueness is CRITICAL for GCM security.
//
// The output format is:
//   - Nonce (12 bytes)
//   - Ciphertext (len(plaintext) bytes)
//   - Authentication tag (16 bytes, included by GCM)
//
// Total overhead: 28 bytes (12 + 16)
//
// Parameters:
//   - key: A 32-byte AES-256 key (must be exactly 32 bytes)
//   - plaintext: The data to encrypt (any length, including empty)
//
// Returns:
//   - []byte: The encrypted data (nonce + ciphertext + tag)
//   - error: ErrInvalidParams if key is not 32 bytes
//     ErrInsufficientEntropy if nonce generation fails
//
// Example:
//
//	ciphertext, err := crypto.Encrypt(key, plaintext)
//	if err != nil {
//	    return fmt.Errorf("encryption failed: %w", err)
//	}
//	// Write ciphertext to vault file
//
// CRITICAL: Never reuse a key+nonce pair. This function generates a fresh
// nonce for every call to prevent nonce reuse.
func Encrypt(key, plaintext []byte) ([]byte, error) {
	// Validate key length - AES-256 requires exactly 32 bytes
	if len(key) != 32 {
		return nil, ErrInvalidParams
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		// This should never happen with a valid 32-byte key
		return nil, err
	}

	// Create GCM mode wrapper
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		// This should never happen with a valid AES block
		return nil, err
	}

	// Generate fresh nonce IMMEDIATELY before Seal.
	// CRITICAL: nonce must be unique for every encryption with the same key.
	// Using io.ReadFull ensures we get exactly 12 bytes (GCM standard nonce size).
	nonce := make([]byte, gcm.NonceSize()) // 12 bytes for GCM
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, ErrInsufficientEntropy
	}

	// Encrypt and authenticate.
	// gcm.Seal appends the ciphertext and tag to the nonce slice.
	// The nil argument means "no additional authenticated data".
	// Output layout: nonce (12) || ciphertext (len(plaintext)) || tag (16)
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

// Decrypt decrypts ciphertext using AES-256-GCM and verifies the authentication tag.
//
// This function extracts the nonce from the ciphertext (first 12 bytes),
// then decrypts and verifies the remaining data. If the authentication tag
// is invalid (wrong key OR corrupted data OR tampered ciphertext), it returns
// ErrAuthFailed.
//
// Parameters:
//   - key: A 32-byte AES-256 key (must be exactly 32 bytes)
//   - ciphertext: Data encrypted with Encrypt() (nonce + ciphertext + tag)
//
// Returns:
//   - []byte: The decrypted plaintext
//   - error: ErrInvalidParams if key is not 32 bytes
//     ErrAuthFailed if ciphertext is too short, tag verification fails,
//     or the key is incorrect
//
// Example:
//
//	plaintext, err := crypto.Decrypt(key, ciphertext)
//	if err != nil {
//	    if errors.Is(err, crypto.ErrAuthFailed) {
//	        // Wrong password or corrupted vault - allow retry
//	        return promptForPasswordAgain()
//	    }
//	    return fmt.Errorf("decryption failed: %w", err)
//	}
//
// CRITICAL: Returns a single ErrAuthFailed for both "wrong key" and "corrupted data"
// to prevent timing attacks that could distinguish between these cases.
func Decrypt(key, ciphertext []byte) ([]byte, error) {
	// Validate key length
	if len(key) != 32 {
		return nil, ErrInvalidParams
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM mode wrapper
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Check minimum ciphertext length (must contain at least the nonce)
	nonceSize := gcm.NonceSize() // 12 bytes for GCM
	if len(ciphertext) < nonceSize {
		return nil, ErrAuthFailed
	}

	// Split nonce and ciphertext+tag
	// Layout: nonce (12) || ciphertext (variable) || tag (16)
	nonce := ciphertext[:nonceSize]
	ciphertextAndTag := ciphertext[nonceSize:]

	// Decrypt and verify authentication tag.
	// If the tag is invalid (wrong key OR corrupted data), gcm.Open returns an error.
	plaintext, err := gcm.Open(nil, nonce, ciphertextAndTag, nil)
	if err != nil {
		// Return single sentinel error for both wrong key and corruption.
		// This prevents timing attacks that could distinguish between the two cases.
		return nil, ErrAuthFailed
	}

	return plaintext, nil
}

// EncryptWithAAD encrypts plaintext using AES-256-GCM with additional authenticated data (AAD).
//
// Unlike Encrypt, this function returns the nonce and ciphertext separately. This is
// required for the .abditum file format, where the nonce is written to the binary
// header (bytes 37-48) and the ciphertext is written as the payload after the header.
// The AAD is the full 49-byte file header -- any header byte tampering causes authentication failure.
//
// Parameters:
//   - key: A 32-byte AES-256 key (must be exactly 32 bytes)
//   - plaintext: The data to encrypt (any length, including empty)
//   - aad: Additional authenticated data (the file header bytes)
//
// Returns:
//   - nonce: 12-byte GCM nonce (written to header bytes 37-48)
//   - ciphertext: Encrypted data + 16-byte GCM tag (written as file payload)
//   - error: ErrInvalidParams if key is not 32 bytes
//     ErrInsufficientEntropy if nonce generation fails
//
// CRITICAL: Never reuse a key+nonce pair. This function generates a fresh
// nonce for every call via crypto/rand.
func EncryptWithAAD(key, plaintext, aad []byte) (nonce []byte, ciphertext []byte, err error) {
	if len(key) != 32 {
		return nil, nil, ErrInvalidParams
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce = make([]byte, gcm.NonceSize()) // 12 bytes
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, ErrInsufficientEntropy
	}

	// Seal with AAD. Pass nil as dst so ciphertext does NOT include the nonce prefix.
	// The storage layer writes nonce to the header and ciphertext separately.
	sealed := gcm.Seal(nil, nonce, plaintext, aad)

	return nonce, sealed, nil
}

// DecryptWithAAD decrypts ciphertext using AES-256-GCM with additional authenticated data (AAD).
//
// The nonce is provided explicitly (read from the file header bytes 37-48).
// The AAD is the full 49-byte file header -- any header byte tampering causes authentication failure.
//
// Parameters:
//   - key: A 32-byte AES-256 key (must be exactly 32 bytes)
//   - ciphertext: Encrypted data + 16-byte GCM tag (file payload after header)
//   - nonce: 12-byte GCM nonce (from file header bytes 37-48)
//   - aad: Additional authenticated data (the full file header bytes 0-48)
//
// Returns:
//   - []byte: The decrypted plaintext
//   - error: ErrInvalidParams if key is not 32 bytes
//     ErrAuthFailed if tag verification fails (wrong key, tampered AAD, or corrupted ciphertext)
func DecryptWithAAD(key, ciphertext, nonce, aad []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, ErrInvalidParams
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(nonce) != gcm.NonceSize() {
		return nil, ErrAuthFailed
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, aad)
	if err != nil {
		return nil, ErrAuthFailed
	}

	return plaintext, nil
}
