package vault

// Manager orchestrates vault operations and maintains session state.
// All vault mutations go through Manager methods.
// Manager knows WHAT operations exist (high-level workflows),
// entities know HOW to execute (validation + mutation logic).
type Manager struct {
	cofre       *Cofre
	repositorio RepositorioCofre
	senha       []byte
	caminho     string
	bloqueado   bool
}

// NewManager creates a new Manager for the given vault and repository.
// The vault is initially unlocked.
func NewManager(cofre *Cofre, repositorio RepositorioCofre) *Manager {
	return &Manager{
		cofre:       cofre,
		repositorio: repositorio,
		senha:       nil,
		caminho:     "",
		bloqueado:   false,
	}
}

// Vault returns a pointer to the managed vault.
// Per D-08: returns live pointer, safety via package encapsulation.
// TUI cannot mutate private fields even with this pointer.
func (m *Manager) Vault() *Cofre {
	if m.bloqueado {
		return nil
	}
	return m.cofre
}

// IsLocked returns true if the vault is currently locked.
func (m *Manager) IsLocked() bool {
	return m.bloqueado
}

// IsModified returns true if the vault has unsaved changes.
func (m *Manager) IsModified() bool {
	if m.cofre == nil {
		return false
	}
	return m.cofre.modificado
}
