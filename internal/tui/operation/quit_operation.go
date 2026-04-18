package operation

import (
	"errors"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/vault"
)

// vaultSaver é a interface mínima que a QuitOperation precisa do vault.Manager.
// Usar uma interface aqui (em vez do tipo concreto) facilita os testes.
type vaultSaver interface {
	IsModified() bool
	Salvar(forcarSobrescrita bool) error
}

// quitState representa em qual etapa do fluxo de saída estamos.
type quitState int

const (
	quitStateSaving       quitState = iota // salvar normalmente
	quitStateSavingForced                  // salvar forçando sobrescrita de modificação externa
)

// quitMsg é a mensagem interna que dispara a ação de salvar.
type quitMsg struct {
	state quitState
}

// quitSaveResultMsg carrega o resultado da tentativa de salvar.
type quitSaveResultMsg struct {
	err     error
	forcado bool
}

// QuitOperation implementa o fluxo de saída (ctrl+Q) do gerenciador.
//
// Três fluxos possíveis:
//   - Fluxo 3: sem cofre aberto → confirmação simples → sair
//   - Fluxo 4: cofre aberto mas sem alterações → confirmação simples → sair
//   - Fluxo 5: cofre com alterações → modal de decisão → salvar/descartar/voltar
type QuitOperation struct {
	notifier tui.MessageController
	manager  vaultSaver
}

// NewQuitOperation cria uma QuitOperation com o vault.Manager concreto.
// manager pode ser nil quando nenhum cofre está aberto (Fluxo 3).
func NewQuitOperation(notifier tui.MessageController, manager *vault.Manager) *QuitOperation {
	var saver vaultSaver
	if manager != nil {
		saver = manager
	}
	return &QuitOperation{notifier: notifier, manager: saver}
}

// newQuitOperationFromSaver é usada internamente e nos testes (mesmo pacote).
// Permite injetar qualquer implementação de vaultSaver sem depender do tipo concreto.
func newQuitOperationFromSaver(notifier tui.MessageController, saver vaultSaver) *QuitOperation {
	return &QuitOperation{notifier: notifier, manager: saver}
}

// Init inicia o fluxo de saída exibindo o modal adequado conforme o estado do cofre.
func (q *QuitOperation) Init() tea.Cmd {
	if q.manager != nil && q.manager.IsModified() {
		// Fluxo 5: cofre com alterações não salvas
		return tui.OpenModal(q.buildModifiedModal())
	}
	// Fluxo 3 (sem cofre) ou Fluxo 4 (cofre inalterado): confirmação simples
	return tui.OpenModal(q.buildConfirmModal())
}

// Update trata as mensagens internas de salvamento da QuitOperation.
func (q *QuitOperation) Update(msg tea.Msg) tea.Cmd {
	switch m := msg.(type) {
	case quitMsg:
		switch m.state {
		case quitStateSaving:
			q.notifier.SetBusy("Salvando...")
			return func() tea.Msg {
				return quitSaveResultMsg{err: q.manager.Salvar(false), forcado: false}
			}
		case quitStateSavingForced:
			q.notifier.SetBusy("Salvando...")
			return func() tea.Msg {
				return quitSaveResultMsg{err: q.manager.Salvar(true), forcado: true}
			}
		}
	case quitSaveResultMsg:
		if m.err == nil {
			q.notifier.SetSuccess("Cofre salvo.")
			return tea.Quit
		}
		// Arquivo modificado externamente: perguntar ao usuário se deseja sobrescrever
		if !m.forcado && errors.Is(m.err, vault.ErrModifiedExternally) {
			q.notifier.Clear()
			return tui.OpenModal(q.buildConflictModal())
		}
		// Erro genérico ou erro mesmo após forçar: não sair, exibir erro
		q.notifier.SetError(m.err.Error())
		return tui.OperationCompleted()
	}
	return nil
}

// buildConfirmModal cria o modal de confirmação simples para os Fluxos 3 e 4.
func (q *QuitOperation) buildConfirmModal() *modal.ConfirmModal {
	return modal.NewConfirmModal(
		"Sair",
		"Deseja encerrar a aplicação?",
		[]modal.ModalOption{
			{
				Keys:  []design.Key{design.Keys.Enter},
				Label: "Confirmar",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), tea.Quit)
				},
			},
			{
				Keys:  []design.Key{design.Keys.Esc},
				Label: "Voltar",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), tui.OperationCompleted())
				},
			},
		},
	)
}

// buildModifiedModal cria o modal de decisão para o Fluxo 5 (cofre com alterações).
func (q *QuitOperation) buildModifiedModal() *modal.ConfirmModal {
	return modal.NewConfirmModal(
		"Sair",
		"Há alterações não salvas. O que deseja fazer?",
		[]modal.ModalOption{
			{
				Keys:  []design.Key{design.Keys.Enter},
				Label: "Salvar e sair",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return quitMsg{state: quitStateSaving}
					})
				},
			},
			{
				Keys:  []design.Key{design.Letter('d')},
				Label: "Descartar e sair",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), tea.Quit)
				},
			},
			{
				Keys:  []design.Key{design.Keys.Esc},
				Label: "Voltar",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), tui.OperationCompleted())
				},
			},
		},
	)
}

// buildConflictModal cria o modal de resolução de conflito quando o arquivo foi
// modificado externamente (por outro processo) enquanto o cofre estava aberto.
func (q *QuitOperation) buildConflictModal() *modal.ConfirmModal {
	return modal.NewConfirmModal(
		"Conflito",
		"O arquivo foi modificado externamente. Deseja sobrescrever?",
		[]modal.ModalOption{
			{
				Keys:  []design.Key{design.Keys.Enter},
				Label: "Sobrescrever e sair",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return quitMsg{state: quitStateSavingForced}
					})
				},
			},
			{
				Keys:  []design.Key{design.Keys.Esc},
				Label: "Voltar",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), tui.OperationCompleted())
				},
			},
		},
	)
}
