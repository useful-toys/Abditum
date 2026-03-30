package storage_test

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/useful-toys/abditum/internal/crypto"
	"github.com/useful-toys/abditum/internal/storage"
	"github.com/useful-toys/abditum/internal/vault"
)

// fastParams overrides Argon2id params for fast testing.
// These are NOT secure for production use -- only for test speed.
var testPassword = []byte("senha-de-teste-123")

// newTestCofre creates a minimal Cofre for storage tests.
func newTestCofre() *vault.Cofre {
	cofre := vault.NovoCofre()
	return cofre
}

// TestSaveNew_RoundTrip tests that SaveNew + Load returns equivalent Cofre.
func TestSaveNew_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()

	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	// File must exist
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("vault file not created: %v", err)
	}
	if info.Size() < storage.HeaderSize {
		t.Errorf("file too small: %d bytes, want >= %d", info.Size(), storage.HeaderSize)
	}

	loaded, meta, err := storage.Load(path, testPassword)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if loaded == nil {
		t.Fatal("Load() returned nil Cofre")
	}

	// FileMetadata: size must match file size
	if meta.Size != info.Size() {
		t.Errorf("FileMetadata.Size = %d, want %d", meta.Size, info.Size())
	}
	// Hash must be non-zero
	var zeroHash [32]byte
	if meta.Hash == zeroHash {
		t.Error("FileMetadata.Hash is all zeros")
	}
}

// TestSaveNew_HeaderMagic tests that SaveNew writes the ABDT magic bytes.
func TestSaveNew_HeaderMagic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	if err := storage.SaveNew(path, newTestCofre(), testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile() error: %v", err)
	}

	expected := storage.Magic
	if !bytes.Equal(data[0:4], expected[:]) {
		t.Errorf("magic bytes = %v, want %v", data[0:4], expected)
	}
	if data[4] != storage.CurrentFormatVersion {
		t.Errorf("version byte = %d, want %d", data[4], storage.CurrentFormatVersion)
	}
}

// TestLoad_WrongMagic tests that Load returns ErrInvalidMagic for wrong magic.
func TestLoad_WrongMagic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.abditum")

	// Write file with wrong magic
	data := make([]byte, storage.HeaderSize+32)
	data[0] = 'X'
	data[1] = 'Y'
	data[2] = 'Z'
	data[3] = 'W'
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	_, _, err := storage.Load(path, testPassword)
	if !errors.Is(err, storage.ErrInvalidMagic) {
		t.Errorf("Load(wrong magic) = %v, want ErrInvalidMagic", err)
	}
}

// TestLoad_FileTooShort tests that Load returns ErrInvalidMagic for truncated files.
func TestLoad_FileTooShort(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "short.abditum")

	// Write file shorter than header
	if err := os.WriteFile(path, []byte{0x41, 0x42, 0x44, 0x54}, 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	_, _, err := storage.Load(path, testPassword)
	if !errors.Is(err, storage.ErrInvalidMagic) {
		t.Errorf("Load(too short) = %v, want ErrInvalidMagic", err)
	}
}

// TestLoad_VersionTooNew tests that Load returns ErrVersionTooNew for future versions.
func TestLoad_VersionTooNew(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "future.abditum")

	// Write file with valid magic but version 255
	data := make([]byte, storage.HeaderSize+32)
	copy(data[0:4], storage.Magic[:])
	data[4] = 255 // far-future version
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	_, _, err := storage.Load(path, testPassword)
	if !errors.Is(err, storage.ErrVersionTooNew) {
		t.Errorf("Load(version=255) = %v, want ErrVersionTooNew", err)
	}
}

// TestLoad_WrongPassword tests that Load returns ErrAuthFailed for wrong password.
func TestLoad_WrongPassword(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	if err := storage.SaveNew(path, newTestCofre(), testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	_, _, err := storage.Load(path, []byte("senha-errada"))
	if !errors.Is(err, crypto.ErrAuthFailed) {
		t.Errorf("Load(wrong password) = %v, want crypto.ErrAuthFailed", err)
	}
}

// TestLoad_TamperedHeader tests that tampering the header returns ErrAuthFailed.
func TestLoad_TamperedHeader(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	if err := storage.SaveNew(path, newTestCofre(), testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	// Read file and tamper a salt byte
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	data[10] ^= 0xFF // Flip a salt byte (in header, affects AAD)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	_, _, err = storage.Load(path, testPassword)
	if !errors.Is(err, crypto.ErrAuthFailed) {
		t.Errorf("Load(tampered header) = %v, want crypto.ErrAuthFailed", err)
	}
}

// TestLoad_TamperedPayload tests that tampering the payload returns ErrAuthFailed.
func TestLoad_TamperedPayload(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	if err := storage.SaveNew(path, newTestCofre(), testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	// Flip a byte in the payload area
	if len(data) > storage.HeaderSize {
		data[storage.HeaderSize] ^= 0xFF
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	_, _, err = storage.Load(path, testPassword)
	if !errors.Is(err, crypto.ErrAuthFailed) {
		t.Errorf("Load(tampered payload) = %v, want crypto.ErrAuthFailed", err)
	}
}

// TestSave_CreatesBackup tests that Save creates a .bak file.
func TestSave_CreatesBackup(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")
	bakPath := path + ".bak"

	// First: create the initial vault
	cofre := newTestCofre()
	salt, err := crypto.GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() error: %v", err)
	}
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	// Load to get salt (we need it for Save)
	// Actually Save requires salt -- let's use SaveNew then Save
	// Save needs salt -- test using SaveNew which uses its own salt.
	// For Save, we need the salt from the file.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt = data[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	// Now Save (atomic overwrite)
	if err := storage.Save(path, cofre, testPassword, salt); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// .bak should exist
	if _, err := os.Stat(bakPath); os.IsNotExist(err) {
		t.Error("Save() did not create .bak backup")
	}
}

// TestSave_CreatesBak2WhenBakExists tests .bak2 creation when .bak exists.
func TestSave_CreatesBak2WhenBakExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")
	bakPath := path + ".bak"
	bak2Path := path + ".bak2"

	// Create initial file
	if err := storage.SaveNew(path, newTestCofre(), testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := data[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	cofre := newTestCofre()

	// First Save: creates .bak
	if err := storage.Save(path, cofre, testPassword, salt); err != nil {
		t.Fatalf("Save() first call error: %v", err)
	}
	if _, err := os.Stat(bakPath); os.IsNotExist(err) {
		t.Fatal(".bak not created after first Save")
	}

	// Read the new file to get its salt
	data2, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() after first Save error: %v", err)
	}
	salt2 := data2[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	// Second Save: .bak should become .bak2, new .bak created
	if err := storage.Save(path, cofre, testPassword, salt2); err != nil {
		t.Fatalf("Save() second call error: %v", err)
	}

	if _, err := os.Stat(bak2Path); os.IsNotExist(err) {
		t.Error("Save() did not create .bak2 when .bak already existed")
	}
}

// TestSave_NoTmpAfterSuccess tests that .tmp does not remain after successful Save.
func TestSave_NoTmpAfterSuccess(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")
	tmpPath := path + ".tmp"

	if err := storage.SaveNew(path, newTestCofre(), testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := data[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	if err := storage.Save(path, newTestCofre(), testPassword, salt); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	if _, err := os.Stat(tmpPath); !os.IsNotExist(err) {
		t.Error(".tmp file still exists after successful Save")
	}
}

// TestLoad_FileMetadata tests that Load returns correct FileMetadata.
func TestLoad_FileMetadata(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	if err := storage.SaveNew(path, newTestCofre(), testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("os.Stat() error: %v", err)
	}

	_, meta, err := storage.Load(path, testPassword)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if meta.Size != info.Size() {
		t.Errorf("FileMetadata.Size = %d, want %d", meta.Size, info.Size())
	}

	// Compute expected hash
	fileData, err2 := os.ReadFile(path)
	if err2 != nil {
		t.Fatalf("ReadFile() error: %v", err2)
	}
	expectedHash := sha256.Sum256(fileData)
	if meta.Hash != expectedHash {
		t.Errorf("FileMetadata.Hash mismatch: got %x, want %x", meta.Hash, expectedHash)
	}
	// Verify hash is non-zero (implicit from above, but explicit for clarity)
	var zeroHash [32]byte
	if meta.Hash == zeroHash {
		t.Error("FileMetadata.Hash should not be all zeros")
	}
}
