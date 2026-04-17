package operation

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
)

// stubNotifier implementa tui.MessageController para os testes da FakeOperation.
// Registra qual método foi chamado por último e com qual texto.
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

// execCmd executa um tea.Cmd e retorna a mensagem produzida.
// Retorna nil se cmd for nil.
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

func TestFakeOperation_Update_Confirmed_SetsBusyAndStartsWork(t *testing.T) {
	n := &stubNotifier{}
	op := NewFakeOperation(n)

	cmd := op.Update(fakeConfirmedMsg{})

	if op.state != stateExecuting {
		t.Errorf("Update(fakeConfirmedMsg): state esperado stateExecuting, obteve %v", op.state)
	}
	if n.lastMethod != "SetBusy" {
		t.Errorf("Update(fakeConfirmedMsg): esperado SetBusy, notifier recebeu %q", n.lastMethod)
	}
	if cmd == nil {
		t.Error("Update(fakeConfirmedMsg): esperado cmd não-nil (fakeWorkCmd)")
	}
}

func TestFakeOperation_Update_WorkDone_ClearsNotifierAndOpensResultModal(t *testing.T) {
	n := &stubNotifier{}
	op := NewFakeOperation(n)
	op.state = stateExecuting // avança estado manualmente para o correto

	cmd := op.Update(fakeWorkDoneMsg{})

	if n.lastMethod != "Clear" {
		t.Errorf("Update(fakeWorkDoneMsg): esperado Clear no notifier, obteve %q", n.lastMethod)
	}
	msg := execCmd(cmd)
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Update(fakeWorkDoneMsg): esperado OpenModalMsg, obteve %T", msg)
	}
}

func TestFakeOperation_Update_ConfirmedInWrongState_Ignored(t *testing.T) {
	op := NewFakeOperation(&stubNotifier{})
	op.state = stateExecuting // já está executando

	cmd := op.Update(fakeConfirmedMsg{})

	if cmd != nil {
		t.Error("Update(fakeConfirmedMsg) em stateExecuting: deveria ser ignorado (retornar nil)")
	}
}
