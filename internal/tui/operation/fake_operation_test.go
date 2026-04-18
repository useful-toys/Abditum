package operation

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
)

type stubNotifier struct {
	lastMethod string
	lastText   string
}

func (s *stubNotifier) SetBusy(text string)      { s.lastMethod = "SetBusy"; s.lastText = text }
func (s *stubNotifier) SetSuccess(text string)   { s.lastMethod = "SetSuccess"; s.lastText = text }
func (s *stubNotifier) SetError(text string)     { s.lastMethod = "SetError"; s.lastText = text }
func (s *stubNotifier) SetWarning(text string)   { s.lastMethod = "SetWarning"; s.lastText = text }
func (s *stubNotifier) SetInfo(text string)      { s.lastMethod = "SetInfo"; s.lastText = text }
func (s *stubNotifier) SetHintField(text string) { s.lastMethod = "SetHintField"; s.lastText = text }
func (s *stubNotifier) SetHintUsage(text string) { s.lastMethod = "SetHintUsage"; s.lastText = text }
func (s *stubNotifier) Clear()                   { s.lastMethod = "Clear"; s.lastText = "" }

func execCmd(cmd tea.Cmd) tea.Msg {
	if cmd == nil {
		return nil
	}
	return cmd()
}

func TestFakeOperation_Init_EmitsOpenModalMsg(t *testing.T) {
	op := NewFakeOperation(&stubNotifier{})

	cmd := op.Init()

	msg := execCmd(cmd)
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Init: esperado OpenModalMsg, obteve %T", msg)
	}
}

func TestFakeOperation_Update_IgnoresUnknownMsg(t *testing.T) {
	op := NewFakeOperation(&stubNotifier{})

	type randomMsg struct{}
	cmd := op.Update(randomMsg{})

	if cmd != nil {
		t.Error("Update(randomMsg): deveria retornar nil para mensagem desconhecida")
	}
}

func TestFakeOperation_Update_Executing_SetsBusyAndStartsWork(t *testing.T) {
	n := &stubNotifier{}
	op := NewFakeOperation(n)

	cmd := op.Update(fakeOperationMsg{state: stateExecuting})

	if n.lastMethod != "SetBusy" {
		t.Errorf("Update(stateExecuting): esperado SetBusy, notifier recebeu %q", n.lastMethod)
	}
	if cmd == nil {
		t.Error("Update(stateExecuting): esperado cmd não-nil (fakeWorkCmd)")
	}
}

func TestFakeOperation_Update_Done_ClearsNotifierAndOpensResultModal(t *testing.T) {
	n := &stubNotifier{}
	op := NewFakeOperation(n)

	cmd := op.Update(fakeOperationMsg{state: stateDone})

	if n.lastMethod != "Clear" {
		t.Errorf("Update(stateDone): esperado Clear no notifier, obteve %q", n.lastMethod)
	}
	msg := execCmd(cmd)
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Update(stateDone): esperado OpenModalMsg, obteve %T", msg)
	}
}