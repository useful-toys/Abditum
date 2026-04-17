package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/vault"
)

// Operation encapsula uma operação multi-etapa com estado próprio.
// Cada operação é uma mini-máquina de estados que processa mensagens
// e avança seu fluxo de forma autônoma via comandos Tea.
type Operation interface {
	// Init retorna o primeiro comando da operação — normalmente abre um modal
	// para coletar input inicial. Análogo a tea.Model.Init().
	Init() tea.Cmd
	// Update processa mensagens Tea e avança o estado interno.
	// Retorna um comando ou nil. A operação ignora mensagens que não reconhece.
	Update(msg tea.Msg) tea.Cmd
}

// StartOperationMsg inicia uma operação (e encerra a atual, se houver).
// Emitida por ações no setup.go ou por operações encadeando outra.
type StartOperationMsg struct{ Op Operation }

// OperationCompletedMsg sinaliza conclusão da operação ativa sem continuação.
// Root limpa activeOperation ao receber esta mensagem.
type OperationCompletedMsg struct{}

// VaultOpenedMsg sinaliza que um cofre foi aberto ou criado com sucesso.
// Emitida pela operação; tratada pelo root para configurar vaultManager.
type VaultOpenedMsg struct{ Manager *vault.Manager }

// SecretExportedMsg sinaliza exportação bem-sucedida de um segredo.
type SecretExportedMsg struct{}

// StartOperation cria um Cmd que emite StartOperationMsg.
func StartOperation(op Operation) tea.Cmd {
	return func() tea.Msg { return StartOperationMsg{Op: op} }
}

// OperationCompleted cria um Cmd que emite OperationCompletedMsg.
func OperationCompleted() tea.Cmd {
	return func() tea.Msg { return OperationCompletedMsg{} }
}
