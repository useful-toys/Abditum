package operation

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
)

type fakeState int

const (
	stateExecuting fakeState = iota
	stateDone
)

type FakeOperationMsg struct {
	state fakeState
}

type FakeOperation struct {
	notifier tui.MessageController
}

func NewFakeOperation(notifier tui.MessageController) *FakeOperation {
	return &FakeOperation{notifier: notifier}
}

func (f *FakeOperation) Init() tea.Cmd {
	return tui.OpenModal(f.buildConfirmModal())
}

func (f *FakeOperation) Update(msg tea.Msg) tea.Cmd {
	m, ok := msg.(FakeOperationMsg)
	if !ok {
		return nil
	}

	switch m.state {
	case stateExecuting:
		f.notifier.SetBusy("Executando operação fake...")
		return fakeWorkCmd()
	case stateDone:
		f.notifier.Clear()
		return tui.OpenModal(f.buildResultModal())
	}
	return nil
}

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
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return FakeOperationMsg{state: stateExecuting}
					})
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

func fakeWorkCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(5 * time.Second)
		return FakeOperationMsg{state: stateDone}
	}
}
