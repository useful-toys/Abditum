package vault

// RepositorioCofre defines the storage interface for vault persistence.
// Implementation provided in Phase 4 (internal/storage package).
// Manager receives repository via dependency injection.
type RepositorioCofre interface {
	// Salvar persists the vault to storage.
	// Phase 4 implementation handles encryption, atomic writes, and backup chain.
	Salvar(cofre *Cofre) error

	// Carregar loads a vault from storage.
	// Phase 4 implementation handles decryption and validation.
	// Returns error if file doesn't exist, wrong password, or corrupted.
	Carregar() (*Cofre, error)

	// DetectarAlteracaoExterna verifica se o arquivo de cofre foi modificado
	// por processo externo desde o último Salvar ou Carregar.
	// Retorna false sem erro se não houver baseline (cofre recém-criado).
	DetectarAlteracaoExterna() (bool, error)
}
