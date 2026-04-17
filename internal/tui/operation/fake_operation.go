package operation

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
)

// fakeState representa o estado interno da FakeOperation.
type fakeState int

const (
	// stateAwaitingConfirmation é o estado inicial — aguarda resposta do modal de confirmação.
	stateAwaitingConfirmation fakeState = iota
	// stateExecuting indica que o trabalho assíncrono está em andamento.
	stateExecuting
)

// fakeConfirmedMsg é emitida pela action closure do modal de confirmação quando o
// usuário confirma. Tipo não exportado — invisível ao RootModel.
type fakeConfirmedMsg struct{}

// fakeWorkDoneMsg é emitida pela goroutine de trabalho fake após 5 segundos.
// Tipo não exportado — invisível ao RootModel.
type fakeWorkDoneMsg struct{}

// FakeOperation valida o padrão Operation de ponta a ponta sem lógica de negócio real.
// Fluxo: modal de confirmação → 5s de trabalho fake → modal de resultado.
type FakeOperation struct {
	state    fakeState
	notifier tui.MessageController
}

// NewFakeOperation cria uma FakeOperation.
// notifier é usado para reportar progresso na barra de mensagem durante o trabalho fake.
func NewFakeOperation(notifier tui.MessageController) *FakeOperation {
	return &FakeOperation{notifier: notifier}
}

// Init abre o modal de confirmação. Implementa tui.Operation.
func (f *FakeOperation) Init() tea.Cmd {
	return tui.OpenModal(f.buildConfirmModal())
}

// Update processa mensagens e avança a máquina de estados. Implementa tui.Operation.
func (f *FakeOperation) Update(msg tea.Msg) tea.Cmd {
	switch msg.(type) {
	case fakeConfirmedMsg:
		if f.state != stateAwaitingConfirmation {
			return nil
		}
		f.state = stateExecuting
		f.notifier.SetBusy("Executando operação fake...")
		return fakeWorkCmd()

	case fakeWorkDoneMsg:
		f.notifier.Clear()
		return tui.OpenModal(f.buildResultModal())
	}
	return nil
}

// buildConfirmModal cria o modal de confirmação inicial.
func (f *FakeOperation) buildConfirmModal() *modal.ConfirmModal {
	return modal.NewConfirmModal(
		"Operação Fake",
		"Deseja executar a operação fake?\nIsso simulará 5 segundos de trabalho.",
		[]modal.ModalOption{
			{
				Keys:   []design.Key{design.Keys.Enter},
				Label:  "Executar",
				Intent: modal.IntentConfirm,
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg { return fakeConfirmedMsg{} })
				},
			},
			{
				Keys:   []design.Key{design.Keys.Esc},
				Label:  "Cancelar",
				Intent: modal.IntentCancel,
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), tui.OperationCompleted())
				},
			},
		},
	)
}

// buildResultModal cria o modal de resultado após o trabalho fake.
func (f *FakeOperation) buildResultModal() *modal.ConfirmModal {
	return modal.NewConfirmModal(
		"Operação Fake",
		"Operação concluída com sucesso!",
		[]modal.ModalOption{
			{
				Keys:   []design.Key{design.Keys.Enter},
				Label:  "OK",
				Intent: modal.IntentConfirm,
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), tui.OperationCompleted())
				},
			},
		},
	)
}

// fakeWorkCmd retorna um Cmd que aguarda 5 segundos em goroutine e emite fakeWorkDoneMsg.
func fakeWorkCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(5 * time.Second)
		return fakeWorkDoneMsg{}
	}
}
