package storage

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/useful-toys/abditum/internal/crypto"
	"github.com/useful-toys/abditum/internal/vault"
)

// SaveNew encrypts cofre with password and writes a new .abditum file at destPath.
//
// A fresh salt and nonce are generated for each call. The function writes the
// file with mode 0600 (owner read/write only).
//
// The file is written directly to destPath (not via atomic rename) because the
// file does not yet exist. Use Save for subsequent writes to an existing vault.
//
// Parameters:
//   - destPath: Absolute path for the new vault file.
//   - cofre: The vault to encrypt and persist.
//   - password: The user's master password (UTF-8 bytes).
//
// Returns an error if serialization, key derivation, encryption, or file I/O fails.
func SaveNew(destPath string, cofre *vault.Cofre, password []byte) error {
	jsonBytes, err := vault.SerializarCofre(cofre)
	if err != nil {
		return fmt.Errorf("storage.SaveNew: serialize: %w", err)
	}
	defer crypto.Wipe(jsonBytes)

	salt, err := crypto.GenerateSalt()
	if err != nil {
		return fmt.Errorf("storage.SaveNew: generate salt: %w", err)
	}

	// Build the 49-byte header with a pre-generated nonce so the full header
	// (including the nonce at bytes 37-48) can serve as GCM AAD.
	header := make([]byte, HeaderSize)
	copy(header[0:MagicSize], Magic[:])
	header[MagicSize] = CurrentFormatVersion
	copy(header[SaltOffset:SaltOffset+SaltSize], salt)

	nonce := header[NonceOffset : NonceOffset+NonceSize]
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("storage.SaveNew: generate nonce: %w", err)
	}

	profile, err := ProfileForVersion(CurrentFormatVersion)
	if err != nil {
		return fmt.Errorf("storage.SaveNew: profile: %w", err)
	}

	key, err := crypto.DeriveKey(password, salt, profile.ToArgonParams())
	if err != nil {
		return fmt.Errorf("storage.SaveNew: derive key: %w", err)
	}
	defer crypto.Wipe(key)

	// SealWithAAD uses the caller-provided nonce and the full 49-byte header as AAD.
	ciphertext, err := crypto.SealWithAAD(key, jsonBytes, nonce, header)
	if err != nil {
		return fmt.Errorf("storage.SaveNew: encrypt: %w", err)
	}

	fileData := append(header, ciphertext...)
	if err := os.WriteFile(destPath, fileData, 0600); err != nil {
		return fmt.Errorf("storage.SaveNew: write file: %w", err)
	}

	return nil
}

// Save encrypts cofre with password and atomically replaces the existing vault
// at vaultPath, preserving the provided salt so the password-derived key is stable.
//
// Atomic rotation protocol (all errors leave the original vault intact):
//  1. Write new content to vaultPath + ".tmp"
//  2. If vaultPath + ".bak" exists, rename it to vaultPath + ".bak2" (overwrites)
//  3. Rename vaultPath to vaultPath + ".bak"
//  4. atomicRename(vaultPath + ".tmp", vaultPath)
//  5. On success: remove vaultPath + ".bak2" (best-effort)
//  6. On any failure after step 1: remove vaultPath + ".tmp" (best-effort)
//
// Parameters:
//   - vaultPath: Absolute path to the existing vault file.
//   - cofre: The vault to encrypt and persist.
//   - password: The user's master password (UTF-8 bytes).
//   - salt: The 32-byte Argon2id salt from the existing file header.
func Save(vaultPath string, cofre *vault.Cofre, password, salt []byte) error {
	jsonBytes, err := vault.SerializarCofre(cofre)
	if err != nil {
		return fmt.Errorf("storage.Save: serialize: %w", err)
	}
	defer crypto.Wipe(jsonBytes)

	// Build header using the provided salt (preserves the derived key).
	header := make([]byte, HeaderSize)
	copy(header[0:MagicSize], Magic[:])
	header[MagicSize] = CurrentFormatVersion
	copy(header[SaltOffset:SaltOffset+SaltSize], salt)

	nonce := header[NonceOffset : NonceOffset+NonceSize]
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("storage.Save: generate nonce: %w", err)
	}

	profile, err := ProfileForVersion(CurrentFormatVersion)
	if err != nil {
		return fmt.Errorf("storage.Save: profile: %w", err)
	}

	key, err := crypto.DeriveKey(password, salt, profile.ToArgonParams())
	if err != nil {
		return fmt.Errorf("storage.Save: derive key: %w", err)
	}
	defer crypto.Wipe(key)

	ciphertext, err := crypto.SealWithAAD(key, jsonBytes, nonce, header)
	if err != nil {
		return fmt.Errorf("storage.Save: encrypt: %w", err)
	}

	fileData := append(header, ciphertext...)

	tmpPath := vaultPath + ".tmp"
	bakPath := vaultPath + ".bak"
	bak2Path := vaultPath + ".bak2"

	// Step 1: write to .tmp
	if err := os.WriteFile(tmpPath, fileData, 0600); err != nil {
		return fmt.Errorf("storage.Save: write tmp: %w", err)
	}

	// From here on, clean up .tmp on any failure.
	cleanup := func() { os.Remove(tmpPath) } //nolint:errcheck

	// Step 2: rotate existing .bak to .bak2
	if _, err := os.Stat(bakPath); err == nil {
		if err := os.Rename(bakPath, bak2Path); err != nil {
			cleanup()
			return fmt.Errorf("storage.Save: rotate bak to bak2: %w", err)
		}
	}

	// Step 3: rename current vault to .bak
	if err := os.Rename(vaultPath, bakPath); err != nil {
		cleanup()
		return fmt.Errorf("storage.Save: rename vault to bak: %w", err)
	}

	// Step 4: atomic rename .tmp to vault
	if err := atomicRename(tmpPath, vaultPath); err != nil {
		// Attempt to restore original
		os.Rename(bakPath, vaultPath) //nolint:errcheck
		cleanup()
		return fmt.Errorf("storage.Save: atomic rename: %w", err)
	}

	return nil
}

// Load reads a .abditum vault file, verifies its integrity, and decrypts it.
//
// Load also returns FileMetadata (size + SHA-256 hash of the raw file bytes),
// which the caller can store to detect external modifications via DetectChange.
//
// Returns:
//   - *vault.Cofre: The decrypted vault.
//   - FileMetadata: Size and SHA-256 of the file as read.
//   - error: ErrInvalidMagic if the file is too short or has wrong magic bytes.
//     ErrVersionTooNew if the format version is not supported.
//     crypto.ErrAuthFailed if password is wrong or file is tampered.
//     ErrCorrupted if decryption succeeds but content is structurally invalid.
func Load(vaultPath string, password []byte) (*vault.Cofre, FileMetadata, error) {
	data, err := os.ReadFile(vaultPath)
	if err != nil {
		return nil, FileMetadata{}, fmt.Errorf("storage.Load: read file: %w", err)
	}

	var meta FileMetadata
	meta.Size = int64(len(data))
	meta.Hash = sha256.Sum256(data)

	// Validate minimum length before checking magic.
	if len(data) < HeaderSize {
		return nil, FileMetadata{}, ErrInvalidMagic
	}

	// Validate magic bytes.
	if data[0] != Magic[0] || data[1] != Magic[1] || data[2] != Magic[2] || data[3] != Magic[3] {
		return nil, FileMetadata{}, ErrInvalidMagic
	}

	version := data[MagicSize]

	// Validate version (ProfileForVersion returns ErrVersionTooNew for unknown versions).
	profile, err := ProfileForVersion(version)
	if err != nil {
		return nil, FileMetadata{}, err
	}

	salt := data[SaltOffset : SaltOffset+SaltSize]
	nonce := data[NonceOffset : NonceOffset+NonceSize]
	header := data[0:HeaderSize]
	ciphertext := data[HeaderSize:]

	key, err := crypto.DeriveKey(password, salt, profile.ToArgonParams())
	if err != nil {
		return nil, FileMetadata{}, fmt.Errorf("storage.Load: derive key: %w", err)
	}
	defer crypto.Wipe(key)

	jsonBytes, err := crypto.DecryptWithAAD(key, ciphertext, nonce, header)
	if err != nil {
		// Propagate ErrAuthFailed directly so callers can errors.Is check it.
		return nil, FileMetadata{}, err
	}
	defer crypto.Wipe(jsonBytes)

	cofre, err := vault.DeserializarCofre(jsonBytes, version)
	if err != nil {
		return nil, FileMetadata{}, fmt.Errorf("%w: %v", ErrCorrupted, err)
	}

	return cofre, meta, nil
}
