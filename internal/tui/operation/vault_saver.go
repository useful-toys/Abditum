package operation

// vaultSaver é a interface mínima que as operations precisam do vault.Manager.
// Usar uma interface aqui (em vez do tipo concreto) facilita os testes.
type vaultSaver interface {
	IsModified() bool
	Salvar(forcarSobrescrita bool) error
}
