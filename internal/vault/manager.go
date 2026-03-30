package vault

import (
	"github.com/useful-toys/abditum/internal/crypto"
)

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

// Lock securely locks the vault by wiping sensitive data from memory.
// Per CRYPTO-04: password and all sensitive field values are overwritten with zeros.
// After locking, vault reference is cleared and IsLocked() returns true.
func (m *Manager) Lock() {
	if m.bloqueado {
		return // Already locked
	}

	// Wipe master password
	if m.senha != nil {
		crypto.Wipe(m.senha)
		m.senha = nil
	}

	// Wipe all sensitive field values recursively
	if m.cofre != nil {
		m.limparCamposSensiveis(m.cofre.pastaGeral)
	}

	// Clear vault reference
	m.cofre = nil
	m.bloqueado = true
}

// limparCamposSensiveis recursively wipes sensitive field values in all secrets.
func (m *Manager) limparCamposSensiveis(pasta *Pasta) {
	if pasta == nil {
		return
	}

	// Wipe sensitive fields in all secrets
	for _, segredo := range pasta.segredos {
		for i := range segredo.campos {
			if segredo.campos[i].tipo == TipoCampoSensivel {
				crypto.Wipe(segredo.campos[i].valor)
				segredo.campos[i].valor = nil
			}
		}
		// Wipe observation (always common type but still sensitive content)
		crypto.Wipe(segredo.observacao.valor)
		segredo.observacao.valor = nil
	}

	// Recurse into subfolders
	for _, subpasta := range pasta.subpastas {
		m.limparCamposSensiveis(subpasta)
	}
}
