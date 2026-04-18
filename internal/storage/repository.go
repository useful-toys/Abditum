package storage

import (
	"fmt"
	"os"

	"github.com/useful-toys/abditum/internal/vault"
)

// Compile-time verification that FileRepository implements vault.RepositorioCofre.
var _ vault.RepositorioCofre = (*FileRepository)(nil)

// FileRepository implements vault.RepositorioCofre using the .abditum binary file format.
//
// FileRepository holds the password, current salt, vault path, and last-known
// FileMetadata as internal state. The salt is preserved across saves so the
// password-derived key is stable (salt only changes on explicit password change).
//
// External change detection is available via:
//   - repo.Metadata() — returns last known FileMetadata snapshot
//   - storage.DetectExternalChange(repo.Path(), repo.Metadata()) — compares to current file state
//
// Lifecycle:
//   - New vault: NewFileRepositoryForCreate → Salvar (uses SaveNew on first call)
//   - Existing vault: NewFileRepository (after Load) or load via Carregar
type FileRepository struct {
	path     string
	password []byte
	salt     []byte
	isNew    bool         // true until first Salvar succeeds
	metadata FileMetadata // last known FileMetadata (set after each Salvar/Carregar)
}

// NewFileRepository creates a FileRepository for an existing vault file.
//
// Call after a successful Load when you already have the password and salt
// (salt is at bytes [SaltOffset:SaltOffset+SaltSize] in the file header).
// The metadata parameter should be the FileMetadata returned by Load.
//
// Parameters:
//   - path: Absolute path to the existing vault file.
//   - password: The user's master password (UTF-8 bytes). Stored by reference — caller must not wipe until done.
//   - salt: The 32-byte Argon2id salt from the file header.
//   - metadata: FileMetadata recorded at Load time (for change detection).
func NewFileRepository(path string, password, salt []byte, metadata FileMetadata) *FileRepository {
	return &FileRepository{
		path:     path,
		password: password,
		salt:     salt,
		isNew:    false,
		metadata: metadata,
	}
}

// NewFileRepositoryForCreate creates a FileRepository for a vault that has not yet been saved.
//
// The first call to Salvar will use SaveNew (direct write without .tmp protocol).
// Subsequent calls use Save (atomic .tmp + rename protocol with backup chain).
//
// Parameters:
//   - path: Absolute path where the vault file will be created.
//   - password: The user's master password (UTF-8 bytes). Stored by reference — caller must not wipe until done.
func NewFileRepositoryForCreate(path string, password []byte) *FileRepository {
	return &FileRepository{
		path:     path,
		password: password,
		salt:     nil,
		isNew:    true,
	}
}

// Salvar implements vault.RepositorioCofre.
//
// First save (isNew == true): uses SaveNew — creates the vault file directly at the
// configured path with a fresh salt and nonce. Transitions isNew to false.
//
// Subsequent saves: uses Save with the stored salt — generates a fresh nonce but
// reuses the salt (stable key across saves), writes via atomic .tmp + rename protocol
// with .bak/.bak2 backup chain.
//
// After a successful save, the internal metadata snapshot is updated.
func (r *FileRepository) Salvar(cofre *vault.Cofre) error {
	if r.isNew {
		if err := SaveNew(r.path, cofre, r.password); err != nil {
			return err
		}
		r.isNew = false

		// Snapshot metadata and extract salt from the written file
		meta, err := ComputeFileMetadata(r.path)
		if err != nil {
			return fmt.Errorf("FileRepository.Salvar: compute metadata after SaveNew: %w", err)
		}
		r.metadata = meta

		salt, err := readSaltFromFile(r.path)
		if err != nil {
			return fmt.Errorf("FileRepository.Salvar: read salt after SaveNew: %w", err)
		}
		r.salt = salt
		return nil
	}

	if err := Save(r.path, cofre, r.password, r.salt); err != nil {
		return err
	}

	meta, err := ComputeFileMetadata(r.path)
	if err != nil {
		return fmt.Errorf("FileRepository.Salvar: compute metadata after Save: %w", err)
	}
	r.metadata = meta
	return nil
}

// Carregar implements vault.RepositorioCofre.
//
// Loads, decrypts, and deserializes the vault from the configured path.
// After a successful load, the internal metadata snapshot and salt are updated.
//
// Returns:
//   - *vault.Cofre: The decrypted vault.
//   - error: ErrInvalidMagic, ErrVersionTooNew, crypto.ErrAuthFailed, or ErrCorrupted.
func (r *FileRepository) Carregar() (*vault.Cofre, error) {
	cofre, meta, err := Load(r.path, r.password)
	if err != nil {
		return nil, err
	}
	r.metadata = meta

	salt, err := readSaltFromFile(r.path)
	if err != nil {
		return nil, fmt.Errorf("FileRepository.Carregar: read salt: %w", err)
	}
	r.salt = salt

	return cofre, nil
}

// Metadata returns the FileMetadata snapshot recorded after the last Salvar or Carregar.
//
// Use with DetectExternalChange to check whether the file was modified externally:
//
//	changed, err := storage.DetectExternalChange(repo.Path(), repo.Metadata())
func (r *FileRepository) Metadata() FileMetadata {
	return r.metadata
}

// Path returns the vault file path configured for this repository.
func (r *FileRepository) Path() string {
	return r.path
}

// DetectarAlteracaoExterna verifica se o arquivo de cofre foi modificado por processo
// externo desde o último Salvar ou Carregar.
// Retorna false sem erro se o metadata ainda não foi capturado (cofre recém-criado).
func (r *FileRepository) DetectarAlteracaoExterna() (bool, error) {
	// Metadata zero significa que nenhum Salvar ou Carregar ocorreu ainda.
	// Não há baseline para comparar — considerar sem alteração externa.
	if r.metadata == (FileMetadata{}) {
		return false, nil
	}
	return DetectExternalChange(r.path, r.metadata)
}

// UpdatePassword replaces the stored password.
//
// Call when the user changes their master password. The next Salvar call
// will use the new password. The caller is responsible for wiping the old
// password slice before discarding it.
func (r *FileRepository) UpdatePassword(password []byte) {
	r.password = password
}

// readSaltFromFile reads the 32-byte Argon2id salt from the vault file header.
//
// The salt occupies bytes [SaltOffset : SaltOffset+SaltSize] of the header.
// Returns a copy of the salt bytes (does not retain a reference to the mmap/buffer).
func readSaltFromFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(data) < SaltOffset+SaltSize {
		return nil, ErrInvalidMagic
	}
	salt := make([]byte, SaltSize)
	copy(salt, data[SaltOffset:SaltOffset+SaltSize])
	return salt, nil
}
