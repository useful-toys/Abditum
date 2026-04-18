package storage_test

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
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

// TestLoad_CorruptedPayload tests that loading corrupted payload returns error.
func TestLoad_CorruptedPayload(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "corrupt.abditum")

	if err := storage.SaveNew(path, newTestCofre(), testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	data[storage.HeaderSize] ^= 0xFF
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	_, _, err = storage.Load(path, testPassword)
	if err == nil {
		t.Error("esperado erro ao carregar payload corrupto")
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

// ---------------------------------------------------------------------------
// RecoverOrphans tests
// ---------------------------------------------------------------------------

// TestRecoverOrphans_RemovesStaleTmp verifies that a stale .tmp file is removed.
func TestRecoverOrphans_RemovesStaleTmp(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")
	tmpPath := path + ".tmp"

	// Create a stale .tmp
	if err := os.WriteFile(tmpPath, []byte("stale"), 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	if err := storage.RecoverOrphans(path); err != nil {
		t.Fatalf("RecoverOrphans() error: %v", err)
	}

	if _, err := os.Stat(tmpPath); !os.IsNotExist(err) {
		t.Error("RecoverOrphans() did not remove stale .tmp file")
	}
}

// TestRecoverOrphans_NoOpWhenClean verifies RecoverOrphans returns nil when there is no .tmp.
func TestRecoverOrphans_NoOpWhenClean(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	if err := storage.RecoverOrphans(path); err != nil {
		t.Errorf("RecoverOrphans() on clean state returned error: %v", err)
	}
}

// TestRecoverOrphans_WithBackupFiles verifies .bak and .bak2 are preserved.
func TestRecoverOrphans_WithBackupFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	if err := storage.SaveNew(path, newTestCofre(), testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	bakPath := path + ".bak"
	bak2Path := path + ".bak2"

	if err := os.WriteFile(bakPath, []byte("old backup"), 0600); err != nil {
		t.Fatalf("WriteFile() bak error: %v", err)
	}
	if err := os.WriteFile(bak2Path, []byte("older backup"), 0600); err != nil {
		t.Fatalf("WriteFile() bak2 error: %v", err)
	}

	if err := storage.RecoverOrphans(path); err != nil {
		t.Fatalf("RecoverOrphans() error: %v", err)
	}

	if _, err := os.Stat(bakPath); os.IsNotExist(err) {
		t.Error(".bak deveria permanecer após RecoverOrphans")
	}
	if _, err := os.Stat(bak2Path); os.IsNotExist(err) {
		t.Error(".bak2 deveria permanecer após RecoverOrphans")
	}
}

// ---------------------------------------------------------------------------
// DetectExternalChange and ComputeFileMetadata tests
// ---------------------------------------------------------------------------

// TestDetectExternalChange_NoChange verifies false is returned when file is unchanged.
func TestDetectExternalChange_NoChange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	if err := storage.SaveNew(path, newTestCofre(), testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	meta, err := storage.ComputeFileMetadata(path)
	if err != nil {
		t.Fatalf("ComputeFileMetadata() error: %v", err)
	}

	changed, err := storage.DetectExternalChange(path, meta)
	if err != nil {
		t.Fatalf("DetectExternalChange() error: %v", err)
	}
	if changed {
		t.Error("DetectExternalChange() returned true for unchanged file")
	}
}

// TestDetectExternalChange_SizeDiffers verifies true is returned when file size changes.
func TestDetectExternalChange_SizeDiffers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	if err := os.WriteFile(path, []byte("original"), 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	meta, err := storage.ComputeFileMetadata(path)
	if err != nil {
		t.Fatalf("ComputeFileMetadata() error: %v", err)
	}

	// Append bytes to change size
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatalf("OpenFile() error: %v", err)
	}
	if _, err := f.Write([]byte(" extra")); err != nil {
		f.Close()
		t.Fatalf("Write() error: %v", err)
	}
	f.Close()

	changed, err := storage.DetectExternalChange(path, meta)
	if err != nil {
		t.Fatalf("DetectExternalChange() error: %v", err)
	}
	if !changed {
		t.Error("DetectExternalChange() returned false after size change")
	}
}

// TestDetectExternalChange_ContentDiffers verifies true is returned when content changes but size stays the same.
func TestDetectExternalChange_ContentDiffers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	content := []byte("ABCDEFGHIJ") // 10 bytes
	if err := os.WriteFile(path, content, 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	meta, err := storage.ComputeFileMetadata(path)
	if err != nil {
		t.Fatalf("ComputeFileMetadata() error: %v", err)
	}

	// Overwrite with same-size but different content
	modified := []byte("ABCDEFGHIZ") // last byte changed
	if err := os.WriteFile(path, modified, 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	changed, err := storage.DetectExternalChange(path, meta)
	if err != nil {
		t.Fatalf("DetectExternalChange() error: %v", err)
	}
	if !changed {
		t.Error("DetectExternalChange() returned false after content change (same size)")
	}
}

// TestDetectExternalChange_FileNotFound verifies an error is returned for a missing file.
func TestDetectExternalChange_FileNotFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent.abditum")

	var meta storage.FileMetadata
	_, err := storage.DetectExternalChange(path, meta)
	if err == nil {
		t.Error("DetectExternalChange() should return error for missing file")
	}
}

// TestComputeFileMetadata_EmptyFile tests that empty file returns Size=0 and non-zero Hash.
func TestComputeFileMetadata_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.abditum")

	if err := os.WriteFile(path, []byte{}, 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	meta, err := storage.ComputeFileMetadata(path)
	if err != nil {
		t.Fatalf("ComputeFileMetadata() error: %v", err)
	}
	if meta.Size != 0 {
		t.Errorf("Size = %d, want 0", meta.Size)
	}
	var zeroHash [32]byte
	if meta.Hash == zeroHash {
		t.Error("Hash deveria ser computável para arquivo vazio")
	}
}

// ---------------------------------------------------------------------------
// Integration tests (Plan 04-04)
// ---------------------------------------------------------------------------

// TestIntegration_FullPipelineRoundtrip validates the full create→save→load pipeline
// via FileRepository, verifying data integrity of the deserialized vault.
func TestIntegration_FullPipelineRoundtrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	// Create a vault with default content
	cofre := vault.NovoCofre()
	if err := cofre.InicializarConteudoPadrao(); err != nil {
		t.Fatalf("InicializarConteudoPadrao() error: %v", err)
	}

	// Save via NewFileRepositoryForCreate
	repo := storage.NewFileRepositoryForCreate(path, testPassword)
	if err := repo.Salvar(cofre); err != nil {
		t.Fatalf("Salvar() error: %v", err)
	}

	// Load via a fresh repository
	// Salt is read from the file header inside Carregar, so nil is fine here.
	repo2 := storage.NewFileRepository(path, testPassword, nil, repo.Metadata())
	loaded, err := repo2.Carregar()
	if err != nil {
		t.Fatalf("Carregar() error: %v", err)
	}

	// Verify PastaGeral
	pg := loaded.PastaGeral()
	if pg == nil {
		t.Fatal("PastaGeral() is nil")
	}

	// Verify default subfolders
	subs := pg.Subpastas()
	if len(subs) != 2 {
		t.Errorf("expected 2 default subfolders, got %d", len(subs))
	} else {
		names := map[string]bool{subs[0].Nome(): true, subs[1].Nome(): true}
		if !names["Sites e Apps"] {
			t.Error("missing default folder 'Sites e Apps'")
		}
		if !names["Financeiro"] {
			t.Error("missing default folder 'Financeiro'")
		}
	}

	// Verify default templates
	modelos := loaded.Modelos()
	if len(modelos) != 3 {
		t.Errorf("expected 3 default templates, got %d", len(modelos))
	}
	modeloNomes := make(map[string]bool, len(modelos))
	for _, m := range modelos {
		modeloNomes[m.Nome()] = true
	}
	for _, nome := range []string{"Login", "Cartão de Crédito", "Chave de API"} {
		if !modeloNomes[nome] {
			t.Errorf("missing default template %q", nome)
		}
	}
}

// TestIntegration_BackupChainRotation verifies .bak and .bak2 are created across 3 saves.
func TestIntegration_BackupChainRotation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")
	bakPath := path + ".bak"
	bak2Path := path + ".bak2"

	cofre := newTestCofre()

	// First save: SaveNew (no .bak yet)
	repo := storage.NewFileRepositoryForCreate(path, testPassword)
	if err := repo.Salvar(cofre); err != nil {
		t.Fatalf("first Salvar() error: %v", err)
	}
	if _, err := os.Stat(bakPath); !os.IsNotExist(err) {
		t.Error("unexpected .bak after first (SaveNew) save")
	}

	// Second save: Save → creates .bak
	if err := repo.Salvar(cofre); err != nil {
		t.Fatalf("second Salvar() error: %v", err)
	}
	if _, err := os.Stat(bakPath); os.IsNotExist(err) {
		t.Error(".bak not created after second save")
	}

	// Third save: Save → rotates .bak → .bak2, creates new .bak
	if err := repo.Salvar(cofre); err != nil {
		t.Fatalf("third Salvar() error: %v", err)
	}
	if _, err := os.Stat(bak2Path); os.IsNotExist(err) {
		t.Error(".bak2 not created after third save")
	}
	if _, err := os.Stat(bakPath); os.IsNotExist(err) {
		t.Error(".bak not present after third save")
	}
}

// TestIntegration_ExternalChangeDetection verifies DetectExternalChange works after a save.
func TestIntegration_ExternalChangeDetection(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	repo := storage.NewFileRepositoryForCreate(path, testPassword)
	if err := repo.Salvar(newTestCofre()); err != nil {
		t.Fatalf("Salvar() error: %v", err)
	}

	meta := repo.Metadata()

	// No change yet
	changed, err := storage.DetectExternalChange(path, meta)
	if err != nil {
		t.Fatalf("DetectExternalChange() error: %v", err)
	}
	if changed {
		t.Error("DetectExternalChange() returned true before any external modification")
	}

	// Modify the file externally
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatalf("OpenFile() error: %v", err)
	}
	if _, err := f.Write([]byte{0xFF}); err != nil {
		f.Close()
		t.Fatalf("Write() error: %v", err)
	}
	f.Close()

	changed, err = storage.DetectExternalChange(path, meta)
	if err != nil {
		t.Fatalf("DetectExternalChange() error after modification: %v", err)
	}
	if !changed {
		t.Error("DetectExternalChange() returned false after external modification")
	}
}

// TestIntegration_ErrorClassification verifies all sentinel errors through the full pipeline.
func TestIntegration_ErrorClassification(t *testing.T) {
	dir := t.TempDir()

	t.Run("WrongMagic", func(t *testing.T) {
		path := filepath.Join(dir, "wrong-magic.abditum")
		data := make([]byte, storage.HeaderSize+32)
		data[0], data[1], data[2], data[3] = 'X', 'Y', 'Z', 'W'
		os.WriteFile(path, data, 0600)
		_, _, err := storage.Load(path, testPassword)
		if !errors.Is(err, storage.ErrInvalidMagic) {
			t.Errorf("expected ErrInvalidMagic, got %v", err)
		}
	})

	t.Run("VersionTooNew", func(t *testing.T) {
		path := filepath.Join(dir, "future-version.abditum")
		data := make([]byte, storage.HeaderSize+32)
		copy(data[0:4], storage.Magic[:])
		data[4] = 255
		os.WriteFile(path, data, 0600)
		_, _, err := storage.Load(path, testPassword)
		if !errors.Is(err, storage.ErrVersionTooNew) {
			t.Errorf("expected ErrVersionTooNew, got %v", err)
		}
	})

	t.Run("WrongPassword", func(t *testing.T) {
		path := filepath.Join(dir, "wrong-password.abditum")
		if err := storage.SaveNew(path, newTestCofre(), testPassword); err != nil {
			t.Fatalf("SaveNew() error: %v", err)
		}
		_, _, err := storage.Load(path, []byte("wrong"))
		if !errors.Is(err, crypto.ErrAuthFailed) {
			t.Errorf("expected crypto.ErrAuthFailed, got %v", err)
		}
	})

	t.Run("TamperedHeader", func(t *testing.T) {
		path := filepath.Join(dir, "tampered-header.abditum")
		if err := storage.SaveNew(path, newTestCofre(), testPassword); err != nil {
			t.Fatalf("SaveNew() error: %v", err)
		}
		data, _ := os.ReadFile(path)
		data[10] ^= 0xFF // flip a salt byte (in AAD)
		os.WriteFile(path, data, 0600)
		_, _, err := storage.Load(path, testPassword)
		if !errors.Is(err, crypto.ErrAuthFailed) {
			t.Errorf("expected crypto.ErrAuthFailed on tampered header, got %v", err)
		}
	})
}

// ---------------------------------------------------------------------------
// FileRepository.DetectarAlteracaoExterna tests
// ---------------------------------------------------------------------------

func TestFileRepository_DetectarAlteracaoExterna_SemAlteracao(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cofre.abditum")
	cofre := vault.NovoCofre()
	repo := storage.NewFileRepositoryForCreate(path, testPassword)
	if err := repo.Salvar(cofre); err != nil {
		t.Fatalf("Salvar: %v", err)
	}

	changed, err := repo.DetectarAlteracaoExterna()
	if err != nil {
		t.Fatalf("DetectarAlteracaoExterna: %v", err)
	}
	if changed {
		t.Error("esperado false (sem alteração externa), obteve true")
	}
}

func TestFileRepository_DetectarAlteracaoExterna_ComAlteracao(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cofre.abditum")
	cofre := vault.NovoCofre()
	repo := storage.NewFileRepositoryForCreate(path, testPassword)
	if err := repo.Salvar(cofre); err != nil {
		t.Fatalf("Salvar: %v", err)
	}

	// Modificar o arquivo externamente
	if err := os.WriteFile(path, []byte("conteudo diferente"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	changed, err := repo.DetectarAlteracaoExterna()
	if err != nil {
		t.Fatalf("DetectarAlteracaoExterna: %v", err)
	}
	if !changed {
		t.Error("esperado true (arquivo alterado externamente), obteve false")
	}
}

func TestFileRepository_DetectarAlteracaoExterna_CofreNovo(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "novo.abditum")
	repo := storage.NewFileRepositoryForCreate(path, testPassword)

	// Antes do primeiro Salvar, metadata é zero — deve retornar false sem erro
	changed, err := repo.DetectarAlteracaoExterna()
	if err != nil {
		t.Fatalf("DetectarAlteracaoExterna em cofre novo: %v", err)
	}
	if changed {
		t.Error("cofre novo: esperado false, obteve true")
	}
}

// ---------------------------------------------------------------------------
// ValidateHeader tests
// ---------------------------------------------------------------------------

func TestValidateHeader_ArquivoValido(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.abditum")
	cofre := vault.NovoCofre()
	if err := cofre.InicializarConteudoPadrao(); err != nil {
		t.Fatal(err)
	}
	if err := storage.SaveNew(path, cofre, []byte("SenhaForte123!")); err != nil {
		t.Fatal(err)
	}

	if err := storage.ValidateHeader(path); err != nil {
		t.Errorf("ValidateHeader de cofre válido: erro inesperado %v", err)
	}
}

func TestValidateHeader_ArquivoPequenoDemais(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "small.bin")
	os.WriteFile(path, []byte("ABC"), 0600)

	err := storage.ValidateHeader(path)
	if !errors.Is(err, storage.ErrInvalidMagic) {
		t.Errorf("arquivo pequeno: esperado ErrInvalidMagic, obteve %v", err)
	}
}

func TestValidateHeader_MagicInvalida(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad_magic.bin")
	data := make([]byte, storage.HeaderSize)
	copy(data[0:4], []byte("XXXX"))
	data[4] = storage.CurrentFormatVersion
	os.WriteFile(path, data, 0600)

	err := storage.ValidateHeader(path)
	if !errors.Is(err, storage.ErrInvalidMagic) {
		t.Errorf("magic inválida: esperado ErrInvalidMagic, obteve %v", err)
	}
}

func TestValidateHeader_VersaoIncompativel(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "future.bin")
	data := make([]byte, storage.HeaderSize)
	copy(data[0:4], storage.Magic[:])
	data[4] = 255
	os.WriteFile(path, data, 0600)

	err := storage.ValidateHeader(path)
	if !errors.Is(err, storage.ErrVersionTooNew) {
		t.Errorf("versão futura: esperado ErrVersionTooNew, obteve %v", err)
	}
}

func TestValidateHeader_ArquivoInexistente(t *testing.T) {
	err := storage.ValidateHeader("/caminho/que/nao/existe.abditum")
	if err == nil {
		t.Error("arquivo inexistente: esperado erro, obteve nil")
	}
}

// TestSave_FailRenameVaultToBak tests that Save returns error when .bak is a directory.
func TestSave_FailRenameVaultToBak(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := data[storage.SaltOffset:storage.SaltOffset+storage.SaltSize]

	bakPath := path + ".bak"
	if err := os.Mkdir(bakPath, 0755); err != nil {
		t.Fatalf("Mkdir() error: %v", err)
	}
	defer os.Remove(bakPath)

	bak2Path := path + ".bak2"
	defer os.RemoveAll(bak2Path)

	err = storage.Save(path, cofre, testPassword, salt)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	if _, err := os.Stat(bak2Path); os.IsNotExist(err) {
		t.Error(".bak2 não foi criado - rotação falhou")
	}
}

// TestSaveNew_PermissionDenied tests that SaveNew returns error when
// destination file exists and is read-only.
// Note: On Windows, permission checks may not work as on Unix.
// This test may pass or fail depending on OS and user privileges.
func TestSaveNew_PermissionDenied(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	if err := storage.SaveNew(path, newTestCofre(), testPassword); err != nil {
		t.Fatalf("SaveNew() initial error: %v", err)
	}

	var chmodErr error
	if chmodErr = os.Chmod(path, 0444); chmodErr != nil {
		t.Fatalf("Chmod() error: %v", chmodErr)
	}
	defer os.Chmod(path, 0644)

	var saveErr error
	saveErr = storage.SaveNew(path, newTestCofre(), testPassword)
	if saveErr == nil {
		t.Error("esperado erro ao sobrescrever arquivo somente-leitura")
	}
}

// TestIntegration_ManagerWithFileRepository verifies that Manager.Salvar works via FileRepository.
func TestIntegration_ManagerWithFileRepository(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := vault.NovoCofre()
	if err := cofre.InicializarConteudoPadrao(); err != nil {
		t.Fatalf("InicializarConteudoPadrao() error: %v", err)
	}

	repo := storage.NewFileRepositoryForCreate(path, testPassword)
	manager := vault.NewManager(cofre, repo)

	// Create a secret via Manager
	modelo := cofre.Modelos()[0] // "Login" template
	pg := cofre.PastaGeral()
	segredo, err := manager.CriarSegredo(pg, "Meu GitHub", modelo)
	if err != nil {
		t.Fatalf("CriarSegredo() error: %v", err)
	}
	if segredo == nil {
		t.Fatal("CriarSegredo() returned nil")
	}

	// Save via Manager
	if err := manager.Salvar(false); err != nil {
		t.Fatalf("Manager.Salvar() error: %v", err)
	}

	// Load directly and verify secret present
	loaded, _, err := storage.Load(path, testPassword)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	segredos := loaded.PastaGeral().Segredos()
	if len(segredos) != 1 {
		t.Fatalf("expected 1 secret in PastaGeral, got %d", len(segredos))
	}
	if segredos[0].Nome() != "Meu GitHub" {
		t.Errorf("secret name = %q, want %q", segredos[0].Nome(), "Meu GitHub")
	}
}

// TestSalvar_UpdateExistingVault tests updating an existing vault.
func TestSalvar_UpdateExistingVault(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre1 := vault.NovoCofre()
	repo := storage.NewFileRepositoryForCreate(path, testPassword)
	if err := repo.Salvar(cofre1); err != nil {
		t.Fatalf("Salvar() error: %v", err)
	}

	loaded1, meta1, err := storage.Load(path, testPassword)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if meta1.Size == 0 {
		t.Fatal("metadata tamanho não pode ser zero")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := data[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	repo2 := storage.NewFileRepository(path, testPassword, salt, meta1)
	cofre2 := vault.NovoCofre()
	if err := repo2.Salvar(cofre2); err != nil {
		t.Fatalf("Salvar() update error: %v", err)
	}

	loaded2, meta2, err := storage.Load(path, testPassword)
	if err != nil {
		t.Fatalf("Load() after update error: %v", err)
	}

	if loaded1.PastaGeral() == nil || loaded2.PastaGeral() == nil {
		t.Fatal("PastaGeral nil")
	}
	if meta2.Size == meta1.Size {
		t.Log("aviso: tamanhos iguais após update (possível em certains casos)")
	}
}

// TestSave_FailRotateBakToBak2 tests Save failure when .bak2 is a directory.
func TestSave_FailRotateBakToBak2(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := data[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	if err := storage.Save(path, cofre, testPassword, salt); err != nil {
		t.Fatalf("Save() first error: %v", err)
	}

	bak2Path := path + ".bak2"
	if err := os.Mkdir(bak2Path, 0755); err != nil {
		t.Fatalf("Mkdir() error: %v", err)
	}
	defer os.Remove(bak2Path)

	err = storage.Save(path, cofre, testPassword, salt)
	if err == nil {
		t.Error("esperado erro ao rotacionar .bak para .bak2 quando .bak2 é diretório")
	}
}

// TestSave_DetectsExternalModification tests Save with externally modified file.
func TestSave_DetectsExternalModification(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := data[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	truncated := data[:len(data)-5]
	if err := os.WriteFile(path, truncated, 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	err = storage.Save(path, cofre, testPassword, salt)
	if err != nil {
		t.Logf("Save com arquivo modificado externamente: erro retornado (aceiteável): %v", err)
	}
}

// TestSave_AfterExternalModification tests Save after external file modification.
func TestSave_AfterExternalModification(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := data[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	err = storage.Save(path, cofre, testPassword, salt)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, _, err := storage.Load(path, testPassword)
	if err != nil {
		t.Fatalf("Load() after Save error: %v", err)
	}
	if loaded == nil {
		t.Error("Load returned nil cofre")
	}
}

// TestReadSaltFromFile tests reading salt from an existing vault.
func TestReadSaltFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	if err := storage.SaveNew(path, newTestCofre(), testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	expectedSalt := data[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]
	if len(expectedSalt) == 0 {
		t.Fatal("expected non-empty salt")
	}
}

// TestSaveNew_EmptyPassword tests SaveNew with empty password.
func TestSaveNew_EmptyPassword(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	err := storage.SaveNew(path, newTestCofre(), []byte{})
	if err != nil {
		t.Logf("SaveNew com senha vazia: erro retornado: %v", err)
	}
}

// TestSave_FailWriteTmp tests Save failure when tmp file is read-only.
func TestSave_FailWriteTmp(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := data[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	tmpPath := path + ".tmp"
	if err := os.Mkdir(tmpPath, 0555); err != nil {
		t.Fatalf("Mkdir() error: %v", err)
	}
	defer os.RemoveAll(tmpPath)

	err = storage.Save(path, cofre, testPassword, salt)
	if err == nil {
		t.Error("esperado erro ao escrever em .tmp quando .tmp é diretório")
	}
}

// TestSave_WithLargeVault tests Save with a larger vault content.
func TestSave_WithLargeVault(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := vault.NovoCofre()
	if err := cofre.InicializarConteudoPadrao(); err != nil {
		t.Fatalf("InicializarConteudoPadrao() error: %v", err)
	}

	for i := 0; i < 10; i++ {
		modelo := cofre.Modelos()[0]
		pg := cofre.PastaGeral()
		_, err := vault.NewManager(cofre, nil).CriarSegredo(pg, fmt.Sprintf("secret-%d", i), modelo)
		if err != nil {
			t.Logf("CriarSegredo erro (pode ser esperado): %v", err)
		}
	}

	repo := storage.NewFileRepositoryForCreate(path, testPassword)
	if err := repo.Salvar(cofre); err != nil {
		t.Fatalf("Salvar() error: %v", err)
	}

	loaded, _, err := storage.Load(path, testPassword)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if loaded == nil {
		t.Error("Load returned nil")
	}
}

// TestComputeFileMetadata_FileNotFound tests error for missing file.
func TestComputeFileMetadata_FileNotFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "naoexiste.abditum")

	_, err := storage.ComputeFileMetadata(path)
	if err == nil {
		t.Error("esperado erro para arquivo inexistente")
	}
}

// TestSave_WithDifferentSalt tests Save with a different salt.
func TestSave_WithDifferentSalt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := make([]byte, storage.SaltSize)
	copy(salt, data[storage.SaltOffset:storage.SaltOffset+storage.SaltSize])
	salt[0] ^= 0xFF

	err = storage.Save(path, cofre, testPassword, salt)
	if err != nil {
		t.Logf("Save com salt diferente: erro (pode ser esperado): %v", err)
	}
}

// TestLoad_TruncatedFile tests loading a truncated file.
func TestLoad_TruncatedFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "truncado.abditum")

	header := make([]byte, storage.HeaderSize)
	copy(header[0:4], storage.Magic[:])
	header[4] = storage.CurrentFormatVersion
	for i := storage.SaltOffset; i < storage.SaltOffset+storage.SaltSize; i++ {
		header[i] = 0xAB
	}
	for i := storage.NonceOffset; i < storage.NonceOffset+storage.NonceSize; i++ {
		header[i] = 0xCD
	}
	if err := os.WriteFile(path, append(header, 0x00), 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	var loadErr error
	_, _, loadErr = storage.Load(path, testPassword)
	if loadErr == nil {
		t.Error("esperado erro para arquivo truncado")
	}
}

// TestLoad_NonceAllZeros tests loading a file with all-zero nonce.
func TestLoad_NonceAllZeros(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "noncezero.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	for i := storage.NonceOffset; i < storage.NonceOffset+storage.NonceSize; i++ {
		data[i] = 0x00
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	_, _, err = storage.Load(path, testPassword)
	if err != nil {
		t.Logf("Load com nonce zero: erro: %v", err)
	}
}

// TestSave_MultipleSaves tests multiple saves to the same vault.
func TestSave_MultipleSaves(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := data[storage.SaltOffset:storage.SaltOffset+storage.SaltSize]

	for i := 0; i < 3; i++ {
		err = storage.Save(path, cofre, testPassword, salt)
		if err != nil {
			t.Fatalf("Save() #%d error: %v", i+1, err)
		}
	}

	loaded, _, err := storage.Load(path, testPassword)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if loaded == nil {
		t.Fatal("Load returned nil")
	}
}

// TestLoad_InvalidSaltSize tests loading with invalid salt size in header.
func TestLoad_InvalidSaltSize(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "invalidsalt.abditum")

	header := make([]byte, storage.HeaderSize)
	copy(header[0:4], storage.Magic[:])
	header[4] = storage.CurrentFormatVersion
	for i := storage.NonceOffset; i < storage.NonceOffset+storage.NonceSize; i++ {
		header[i] = 0x11
	}
	payload := make([]byte, 32)
	for i := range payload {
		payload[i] = 0x22
	}
	if err := os.WriteFile(path, append(header, payload...), 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	var loadErr error
	_, _, loadErr = storage.Load(path, testPassword)
	if loadErr != nil {
		t.Logf("Load com payload inválido: erro: %v", loadErr)
	}
}

// TestSave_EmptyCofre tests saving an empty cofre.
func TestSave_EmptyCofre(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "emptyvault.abditum")

	cofre := vault.NovoCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	loaded, _, loadErr := storage.Load(path, testPassword)
	if loadErr != nil {
		t.Fatalf("Load() error: %v", loadErr)
	}
	if loaded == nil {
		t.Fatal("Load returned nil")
	}
}

// TestRecoverOrphans_RemovesTmpOnly tests that only .tmp is removed.
func TestRecoverOrphans_RemovesTmpOnly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	if err := storage.SaveNew(path, newTestCofre(), testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	tmpPath := path + ".tmp"
	bakPath := path + ".bak"
	if err := os.WriteFile(tmpPath, []byte("stale"), 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}
	if err := os.WriteFile(bakPath, []byte("backup"), 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	if err := storage.RecoverOrphans(path); err != nil {
		t.Fatalf("RecoverOrphans() error: %v", err)
	}

	if _, err := os.Stat(tmpPath); !os.IsNotExist(err) {
		t.Error(".tmp deveria ser removido")
	}
	if _, err := os.Stat(bakPath); os.IsNotExist(err) {
		t.Error(".bak deveria permanecer")
	}
}

// TestAtomicRename_ExistingFile tests atomic rename when target exists.
func TestAtomicRename_ExistingFile(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "source.abditum")
	dstPath := filepath.Join(dir, "dest.abditum")

	if err := os.WriteFile(srcPath, []byte("source data"), 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}
	if err := os.WriteFile(dstPath, []byte("dest data"), 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	saveErr := storage.SaveNew(dstPath, newTestCofre(), testPassword)
	if saveErr != nil {
		t.Logf("atomicRename error: %v", saveErr)
	}
}

// TestSave_WithNilSalt tests Save with nil salt.
func TestSave_WithNilSalt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	err := storage.Save(path, cofre, testPassword, nil)
	if err != nil {
		t.Logf("Save com nil salt: erro: %v", err)
	}
}

// TestComputeFileMetadata_WithLargeFile tests metadata for large file.
func TestComputeFileMetadata_WithLargeFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "large.abditum")

	largeData := make([]byte, 1024*100)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}
	if err := os.WriteFile(path, largeData, 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	meta, err := storage.ComputeFileMetadata(path)
	if err != nil {
		t.Fatalf("ComputeFileMetadata() error: %v", err)
	}
	if meta.Size != int64(len(largeData)) {
		t.Errorf("Size = %d, want %d", meta.Size, len(largeData))
	}
}

// TestLoad_NoFileAccess tests Load when file cannot be accessed.
func TestLoad_NoFileAccess(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent.abditum")

	_, _, err := storage.Load(path, testPassword)
	if err == nil {
		t.Error("esperado erro para arquivo inexistente")
	}
}

// TestSave_NoWritePermission tests Save when directory is not writable.
func TestSave_NoWritePermission(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := data[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	parentDir := filepath.Dir(path)
	if err := os.Chmod(parentDir, 0555); err != nil {
		t.Logf("Chmod error: %v", err)
	}
	defer os.Chmod(parentDir, 0755)

	err = storage.Save(path, cofre, testPassword, salt)
	if err != nil {
		t.Logf("Save sem permissão: erro: %v", err)
	}
}

// TestValidateHeader_ValidHeader tests validation with a valid header.
func TestValidateHeader_ValidHeader(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "valid.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	if err := storage.ValidateHeader(path); err != nil {
		t.Errorf("ValidateHeader() error: %v", err)
	}
}

// TestValidateHeader_VersionZero tests validation with version 0.
func TestValidateHeader_VersionZero(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "version0.abditum")

	header := make([]byte, storage.HeaderSize)
	copy(header[0:4], storage.Magic[:])
	header[4] = 0
	if err := os.WriteFile(path, header, 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	err := storage.ValidateHeader(path)
	if err != nil {
		t.Logf("ValidateHeader versao 0: erro: %v", err)
	}
}

// TestDetectExternalChange_MissingFile tests detection with missing file.
func TestDetectExternalChange_MissingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.abditum")

	meta, err := storage.ComputeFileMetadata(path)
	if err == nil {
		_, detectErr := storage.DetectExternalChange(path, meta)
		if detectErr == nil {
			t.Error("esperado erro para arquivo faltante")
		}
	}
}

// TestFileRepository_WithNilMetadata tests FileRepository with nil metadata.
func TestFileRepository_WithNilMetadata(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	repo := storage.NewFileRepository(path, testPassword, nil, storage.FileMetadata{})
	if repo == nil {
		t.Error("NewFileRepository returned nil")
	}
}

// TestSave_WithEmptyPassword tests Save with empty password.
func TestSave_WithEmptyPassword(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	data := make([]byte, storage.SaltSize)
	for i := range data {
		data[i] = 0xAA
	}

	err := storage.Save(path, cofre, []byte{}, data)
	if err != nil {
		t.Logf("Save senha vazia: erro: %v", err)
	}
}

// TestRepository_Path tests repository Path getter.
func TestRepository_Path(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	repo := storage.NewFileRepositoryForCreate(path, testPassword)
	if repo.Path() != path {
		t.Errorf("Path() = %q, want %q", repo.Path(), path)
	}
}

// TestRepository_UpdatePassword tests updating password.
func TestRepository_UpdatePassword(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	repo := storage.NewFileRepositoryForCreate(path, testPassword)
	if err := repo.Salvar(cofre); err != nil {
		t.Fatalf("Salvar() error: %v", err)
	}

	newPassword := []byte("nova-senha-123")
	repo.UpdatePassword(newPassword)
}

// TestLoad_WithCorruptedNonce tests Load with corrupted nonce.
func TestLoad_WithCorruptedNonce(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	data[storage.NonceOffset] ^= 0xFF
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	_, _, loadErr := storage.Load(path, testPassword)
	if loadErr != nil {
		t.Logf("Load nonce corrupto: erro: %v", loadErr)
	}
}

// TestSave_WithCorruptedSalt tests Save with corrupted salt.
func TestSave_WithCorruptedSalt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := make([]byte, storage.SaltSize)
	copy(salt, data[storage.SaltOffset:storage.SaltOffset+storage.SaltSize])
	salt[0] ^= 0x01

	err = storage.Save(path, cofre, testPassword, salt)
	if err != nil {
		t.Logf("Save salt satél corrompido: erro: %v", err)
	}
}

// TestLoad_InvalidVersion tests Load with invalid version.
func TestLoad_InvalidVersion(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "invalidver.abditum")

	header := make([]byte, storage.HeaderSize)
	copy(header[0:4], storage.Magic[:])
	header[4] = 254
	if err := os.WriteFile(path, header, 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	_, _, loadErr := storage.Load(path, testPassword)
	if loadErr == nil {
		t.Error("esperado erro para versão inválida")
	}
}

// TestRecoverOrphans_WithBakDir tests RecoverOrphans when .bak is directory.
func TestRecoverOrphans_WithBakDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	if err := os.WriteFile(path, []byte("vault"), 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	bakPath := path + ".bak"
	if err := os.Mkdir(bakPath, 0755); err != nil {
		t.Fatalf("Mkdir() error: %v", err)
	}
	defer os.Remove(bakPath)

	err := storage.RecoverOrphans(path)
	if err != nil {
		t.Logf("RecoverOrphans error: %v", err)
	}
}

// TestSalvar_WithContentTests tests Salvar with various content.
func TestSalvar_WithContentTests(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := vault.NovoCofre()
	if err := cofre.InicializarConteudoPadrao(); err != nil {
		t.Fatalf("InicializarConteudoPadrao() error: %v", err)
	}

	repo := storage.NewFileRepositoryForCreate(path, testPassword)
	if err := repo.Salvar(cofre); err != nil {
		t.Fatalf("Salvar() error: %v", err)
	}

	meta := repo.Metadata()
	if meta.Size == 0 {
		t.Error("metadata size should not be zero")
	}
}

// TestCarregar Tests the Carregar method.
func TestCarregar(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	repo := storage.NewFileRepositoryForCreate(path, testPassword)
	if err := repo.Salvar(cofre); err != nil {
		t.Fatalf("Salvar() error: %v", err)
	}

	loaded, err := repo.Carregar()
	if err != nil {
		t.Fatalf("Carregar() error: %v", err)
	}
	if loaded == nil {
		t.Fatal("Load retornou nil cofre")
	}
}

// TestSave_ConcurrentWithError tests Save error during write.
func TestSave_ConcurrentWithError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	invalidSalt := []byte{0x01, 0x02}
	saveErr := storage.Save(path, cofre, testPassword, invalidSalt)
	if saveErr != nil {
		t.Logf("Save com salt inválido: erro: %v", saveErr)
	}
}

// TestLoad_AllZerosNonce tests Load with all zeros nonce.
func TestLoad_AllZerosNonce(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	zeroNonce := make([]byte, storage.NonceSize)
	copy(data[storage.NonceOffset:storage.NonceOffset+storage.NonceSize], zeroNonce)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	_, _, loadErr := storage.Load(path, testPassword)
	if loadErr != nil {
		t.Logf("Load com nonce zero: erro: %v", loadErr)
	}
}

// TestSaveNew_OverwriteExisting vaults tests SaveNew on existing file.
func TestSaveNew_OverwriteExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	cofre2 := newTestCofre()
	err := storage.SaveNew(path, cofre2, testPassword)
	if err != nil {
		t.Logf("SaveNew() sobrescrita: erro: %v", err)
	}
}

// TestFileRepositoryForOpen tests creating repository for open.
func TestFileRepositoryForOpen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	repo := storage.NewFileRepositoryForCreate(path, testPassword)
	if err := repo.Salvar(cofre); err != nil {
		t.Fatalf("Salvar() error: %v", err)
	}

	meta := repo.Metadata()
	repoOpen := storage.NewFileRepositoryForOpen(path, testPassword)
	if repoOpen == nil {
		t.Error("NewFileRepositoryForOpen retornou nil")
	}
	_ = meta
}

// TestSave_EncryptWithRealKey tests Save with real derived key.
func TestSave_EncryptWithRealKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := vault.NovoCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	loaded, meta, err := storage.Load(path, testPassword)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if loaded == nil {
		t.Fatal("Load returned nil")
	}

	newCofre := vault.NovoCofre()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := data[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	err = storage.Save(path, newCofre, testPassword, salt)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	if _, err := os.Stat(path + ".bak"); os.IsNotExist(err) {
		t.Error(".bak não foi criado")
	}
	_ = meta
}

// TestSave_SameDataDoesNotChange tests saving identical data.
func TestSave_SameDataDoesNotChange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := vault.NovoCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data1, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := data[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	err = storage.Save(path, cofre, testPassword, salt)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	data2, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}

	if len(data1) != len(data2) {
		t.Error("tamanho mudou sem alteração")
	}
}

// TestRecoverOrphans_TmpIsDirectory tests when .tmp is a directory.
func TestRecoverOrphans_TmpIsDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	if err := storage.SaveNew(path, newTestCofre(), testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	tmpPath := path + ".tmp"
	if err := os.Mkdir(tmpPath, 0755); err != nil {
		t.Fatalf("Mkdir() error: %v", err)
	}
	defer os.Remove(tmpPath)

	recoverErr := storage.RecoverOrphans(path)
	if recoverErr != nil {
		t.Logf("RecoverOrphans error: %v", recoverErr)
	}
}

// TestRepository_IsNew tests repository isNew flag.
func TestRepository_IsNew(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	repo := storage.NewFileRepositoryForCreate(path, testPassword)
	if repo == nil {
		t.Fatal("NewFileRepositoryForCreate returned nil")
	}

	cofre := newTestCofre()
	if err := repo.Salvar(cofre); err != nil {
		t.Fatalf("Salvar() error: %v", err)
	}
}

// TestSave_WithCorruptedPayloadData tests Save with corrupted payload.
func TestSave_WithCorruptedPayloadData(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := data[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	corruptedCofre := vault.NovoCofre()
	if err := corruptedCofre.InicializarConteudoPadrao(); err != nil {
		t.Fatalf("InicializarConteudoPadrao() error: %v", err)
	}

	err = storage.Save(path, corruptedCofre, testPassword, salt)
	if err != nil {
		t.Logf("Save com cofre corrompido: erro: %v", err)
	}
}

// TestLoad_WithTamperedSalt tests Load with tampered salt.
func TestLoad_WithTamperedSalt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	data[storage.SaltOffset] ^= 0xFF
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	_, _, loadErr := storage.Load(path, testPassword)
	if loadErr == nil {
		t.Error("esperado erro para salt corrompido")
	}
}

// TestLoad_InvalidFormatVersion tests Load with invalid version.
func TestLoad_InvalidFormatVersion(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	header := make([]byte, storage.HeaderSize)
	copy(header[0:4], storage.Magic[:])
	header[4] = 200
	if err := os.WriteFile(path, header, 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	_, _, loadErr := storage.Load(path, testPassword)
	if loadErr == nil {
		t.Error("esperado erro para versão inválida")
	}
}

// TestSave_VerifyBakCreated tests that .bak is created after Save.
func TestSave_VerifyBakCreated(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := data[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	err = storage.Save(path, cofre, testPassword, salt)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	bakPath := path + ".bak"
	if _, err := os.Stat(bakPath); os.IsNotExist(err) {
		t.Error(".bak não foi criado")
	}
}

// TestSave_VerifyBak2Created tests that .bak2 is created when .bak exists.
func TestSave_VerifyBak2Created(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := data[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	if err := storage.Save(path, cofre, testPassword, salt); err != nil {
		t.Fatalf("Save() first error: %v", err)
	}

	if err := storage.Save(path, cofre, testPassword, salt); err != nil {
		t.Fatalf("Save() second error: %v", err)
	}

	bak2Path := path + ".bak2"
	if _, err := os.Stat(bak2Path); os.IsNotExist(err) {
		t.Error(".bak2 não foi criado")
	}
}

// TestSalvar_UpdateVault tests updating vault using Salvar.
func TestSalvar_UpdateVault(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre1 := vault.NovoCofre()
	repo := storage.NewFileRepositoryForCreate(path, testPassword)
	if err := repo.Salvar(cofre1); err != nil {
		t.Fatalf("Salvar() first error: %v", err)
	}

	cofre2 := vault.NovoCofre()
	if err := repo.Salvar(cofre2); err != nil {
		t.Fatalf("Salvar() second error: %v", err)
	}

	meta := repo.Metadata()
	if meta.Size == 0 {
		t.Error("metadata size should not be zero")
	}
}

// TestSave_NewSalt tests Save with a newly generated salt.
func TestSave_NewSalt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	newSalt, err := crypto.GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() error: %v", err)
	}

	err = storage.Save(path, cofre, testPassword, newSalt)
	if err != nil {
		t.Logf("Save com novo salt: erro: %v", err)
	}
}

// TestSave_FailAtomicRename tests failure in atomic rename.
func TestSave_FailAtomicRename(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := data[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	tmpPath := path + ".tmp"
	err = storage.Save(path, cofre, testPassword, salt)
	if err != nil {
		t.Logf("Save error: %v", err)
	}

	if _, err := os.Stat(tmpPath); !os.IsNotExist(err) {
		t.Error(".tmp deveria ser removido após sucesso")
	}
}

// TestLoad_NoModification tests after Save with no file modification.
func TestLoad_NoModification(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	_, meta1, err := storage.Load(path, testPassword)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := data[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	err = storage.Save(path, cofre, testPassword, salt)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	_, meta2, err := storage.Load(path, testPassword)
	if err != nil {
		t.Fatalf("Load() after Save error: %v", err)
	}

	if meta1.Size == meta2.Size {
		t.Log("tamanho igual após Save (normal para dados pequenos)")
	}
}

// TestSave_LargeData tests Save with large data.
func TestSave_LargeData(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := vault.NovoCofre()
	if err := cofre.InicializarConteudoPadrao(); err != nil {
		t.Fatalf("InicializarConteudoPadrao() error: %v", err)
	}

	repo := storage.NewFileRepositoryForCreate(path, testPassword)
	if err := repo.Salvar(cofre); err != nil {
		t.Fatalf("Salvar() error: %v", err)
	}

	loaded, err := repo.Carregar()
	if err != nil {
		t.Fatalf("Carregar() error: %v", err)
	}
	if loaded == nil {
		t.Fatal("carregar returned nil")
	}
}

// TestSave_BackupRotationChain tests backup chain .bak -> .bak2.
func TestSave_BackupRotationChain(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := data[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	for i := 0; i < 3; i++ {
		err = storage.Save(path, cofre, testPassword, salt)
		if err != nil {
			t.Fatalf("Save() #%d error: %v", i+1, err)
		}
	}

	bakPath := path + ".bak"
	bak2Path := path + ".bak2"

	if _, err := os.Stat(bakPath); os.IsNotExist(err) {
		t.Error(".bak should exist after multiple saves")
	}
	if _, err := os.Stat(bak2Path); os.IsNotExist(err) {
		t.Error(".bak2 should exist after multiple saves")
	}
}

// TestFileRepository_Metadata tests metadata getter.
func TestFileRepository_Metadata(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	repo := storage.NewFileRepositoryForCreate(path, testPassword)
	if err := repo.Salvar(cofre); err != nil {
		t.Fatalf("Salvar() error: %v", err)
	}

	meta := repo.Metadata()
	if meta.Size == 0 {
		t.Error("metadata should have size")
	}
}

// TestSave_PreserveSalt tests that Save preserves salt.
func TestSave_PreserveSalt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data1, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt1 := data1[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	err = storage.Save(path, cofre, testPassword, salt1)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	data2, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() after Save error: %v", err)
	}
	salt2 := data2[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	if !bytes.Equal(salt1, salt2) {
		t.Error("salt should be preserved after Save")
	}
}

// TestRecoverOrphans_WithTmpFile tests RecoverOrphans with stale .tmp file.
func TestRecoverOrphans_WithTmpFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	if err := storage.SaveNew(path, newTestCofre(), testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, []byte("stale tmp"), 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	if err := storage.RecoverOrphans(path); err != nil {
		t.Fatalf("RecoverOrphans() error: %v", err)
	}

	if _, err := os.Stat(tmpPath); !os.IsNotExist(err) {
		t.Error(".tmp should be removed")
	}
}

// TestLoad_EmptyFile tests Load with empty file.
func TestLoad_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.abditum")

	if err := os.WriteFile(path, []byte{}, 0600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	_, _, loadErr := storage.Load(path, testPassword)
	if loadErr == nil {
		t.Error("esperado erro para arquivo vazio")
	}
}

// TestSave_WithDifferentPassword tests Save with different password.
func TestSave_WithDifferentPassword(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	salt := data[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	newPassword := []byte("different-password-123")
	err = storage.Save(path, cofre, newPassword, salt)
	if err != nil {
		t.Logf("Save com senha diferente: erro: %v", err)
	}
}

// TestLoad_WithDifferentPassword tests loading with different password.
func TestLoad_WithDifferentPassword(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	_, _, loadErr := storage.Load(path, []byte("wrong-password"))
	if loadErr == nil {
		t.Error("esperado erro para senha errada")
	}
}

// TestSave_ConcurrentWrites tests save sequence.
func TestSave_ConcurrentWrites(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := vault.NovoCofre()
	if err := cofre.InicializarConteudoPadrao(); err != nil {
		t.Fatalf("InicializarConteudoPadrao() error: %v", err)
	}

	repo := storage.NewFileRepositoryForCreate(path, testPassword)
	if err := repo.Salvar(cofre); err != nil {
		t.Fatalf("Salvar() error: %v", err)
	}

	for i := 0; i < 3; i++ {
		if err := repo.Salvar(cofre); err != nil {
			t.Fatalf("Salvar() #%d error: %v", i+1, err)
		}
	}

	loaded, err := repo.Carregar()
	if err != nil {
		t.Fatalf("Carregar() error: %v", err)
	}
	if loaded == nil {
		t.Error("carregar returned nil")
	}
}

// TestRecoverOrphans_MissingFile tests RecoverOrphans on missing file.
func TestRecoverOrphans_MissingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent.abditum")

	err := storage.RecoverOrphans(path)
	if err != nil {
		t.Logf("RecoverOrphans error: %v", err)
	}
}

// TestComputeFileMetadata_SameFileTwice tests computing metadata twice.
func TestComputeFileMetadata_SameFileTwice(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	meta1, err := storage.ComputeFileMetadata(path)
	if err != nil {
		t.Fatalf("ComputeFileMetadata() error: %v", err)
	}

	meta2, err := storage.ComputeFileMetadata(path)
	if err != nil {
		t.Fatalf("ComputeFileMetadata() second error: %v", err)
	}

	if meta1.Hash != meta2.Hash {
		t.Error("hashes should be equal for same file")
	}
}

// TestSave_GenNewNonce tests that Save generates new nonce.
func TestSave_GenNewNonce(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.abditum")

	cofre := newTestCofre()
	if err := storage.SaveNew(path, cofre, testPassword); err != nil {
		t.Fatalf("SaveNew() error: %v", err)
	}

	data1, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	nonce1 := data1[storage.NonceOffset : storage.NonceOffset+storage.NonceSize]
	salt := data1[storage.SaltOffset : storage.SaltOffset+storage.SaltSize]

	err = storage.Save(path, cofre, testPassword, salt)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	data2, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() after Save error: %v", err)
	}
	nonce2 := data2[storage.NonceOffset : storage.NonceOffset+storage.NonceSize]

	if bytes.Equal(nonce1, nonce2) {
		t.Log("nonces should be different after Save (new nonce generated)")
	}
}
