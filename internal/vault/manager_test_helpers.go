package vault

// NewManagerForTest cria um Manager com caminho explícito para uso em testes.
// Permite golden tests do HeaderView que precisam de um vault com nome de arquivo real.
func NewManagerForTest(cofre *Cofre, caminho string) *Manager {
	return &Manager{
		cofre:       cofre,
		repositorio: nil,
		senha:       nil,
		caminho:     caminho,
		bloqueado:   false,
	}
}
