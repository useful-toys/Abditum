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
	if err := manager.Salvar(); err != nil {
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
