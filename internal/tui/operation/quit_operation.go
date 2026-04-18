package operation

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/vault"
)

// QuitOperation implementa o fluxo de saída (ctrl+Q) do gerenciador.
//
// Três fluxos possíveis:
//   - Fluxo 3: sem cofre aberto → confirmação simples → sair
//   - Fluxo 4: cofre aberto mas sem alterações → confirmação simples → sair
//   - Fluxo 5: cofre com alterações → guardCofreAlterado → sair
type QuitOperation struct {
	notifier tui.MessageController
	manager  vaultSaver
	guard    *guardCofreAlterado // não-nil apenas no Fluxo 5
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
		// Fluxo 5: cofre com alterações não salvas — delegar ao guard
		q.guard = novoGuardCofreAlterado(
			q.notifier,
			q.manager,
			func() tea.Cmd { return tea.Quit },
			func() tea.Cmd { return tui.OperationCompleted() },
		)
		return q.guard.Init()
	}
	// Fluxo 3 (sem cofre) ou Fluxo 4 (cofre inalterado): confirmação simples
	return tui.OpenModal(q.buildConfirmModal())
}

// Update trata as mensagens internas da QuitOperation.
// Quando o Fluxo 5 está ativo, todas as mensagens são delegadas ao guard.
func (q *QuitOperation) Update(msg tea.Msg) tea.Cmd {
	if q.guard != nil {
		return q.guard.Update(msg)
	}
	// Fluxos 3/4: nenhuma mensagem interna — tudo tratado pelo modal de confirmação
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
