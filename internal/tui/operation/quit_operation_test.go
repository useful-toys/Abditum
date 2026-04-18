package operation

import (
	"errors"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/vault"
)

// stubManager simula vault.Manager para os testes da QuitOperation.
type stubManager struct {
	isModified                    bool
	salvarErr                     error
	salvarForcarSobrescritaCalled bool
}

func (s *stubManager) IsModified() bool { return s.isModified }
func (s *stubManager) Salvar(forcarSobrescrita bool) error {
	if forcarSobrescrita {
		s.salvarForcarSobrescritaCalled = true
	}
	return s.salvarErr
}

// --- Init ---

func TestQuitOperation_Init_SemCofre_AbreModalConfirmacao(t *testing.T) {
	op := newQuitOperationFromSaver(&stubNotifier{}, nil)
	msg := execCmd(op.Init())
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Init sem cofre: esperado OpenModalMsg, obteve %T", msg)
	}
}

func TestQuitOperation_Init_CofreInalterado_AbreModalConfirmacao(t *testing.T) {
	op := newQuitOperationFromSaver(&stubNotifier{}, &stubManager{isModified: false})
	msg := execCmd(op.Init())
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Init cofre inalterado: esperado OpenModalMsg, obteve %T", msg)
	}
}

func TestQuitOperation_Init_CofreAlterado_AbreModalDecisao(t *testing.T) {
	op := newQuitOperationFromSaver(&stubNotifier{}, &stubManager{isModified: true})
	msg := execCmd(op.Init())
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Init cofre alterado: esperado OpenModalMsg, obteve %T", msg)
	}
}

// --- Update ---

func TestQuitOperation_Update_IgnoraMensagemDesconhecida(t *testing.T) {
	op := newQuitOperationFromSaver(&stubNotifier{}, nil)
	type outraMsg struct{}
	if cmd := op.Update(outraMsg{}); cmd != nil {
		t.Error("Update(outraMsg): esperado nil")
	}
}

func TestQuitOperation_Update_Saving_Sucesso_EmiteQuit(t *testing.T) {
	n := &stubNotifier{}
	m := &stubManager{}
	op := newQuitOperationFromSaver(n, m)

	cmd := op.Update(quitMsg{state: quitStateSaving})
	if n.lastMethod != "SetBusy" {
		t.Errorf("Update(saving): esperado SetBusy, obteve %q", n.lastMethod)
	}
	// cmd é a goroutine de salvamento — executar para obter o resultado
	resultMsg := execCmd(cmd)
	// o resultado é quitSaveResultMsg, que Update deve processar
	resultCmd := op.Update(resultMsg)
	msg := execCmd(resultCmd)
	if n.lastMethod != "SetSuccess" {
		t.Errorf("após salvar OK: esperado SetSuccess, obteve %q", n.lastMethod)
	}
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Errorf("após salvar OK: esperado tea.QuitMsg, obteve %T", msg)
	}
}

func TestQuitOperation_Update_Saving_ErroExterno_AbreModalConflito(t *testing.T) {
	n := &stubNotifier{}
	m := &stubManager{salvarErr: vault.ErrModifiedExternally}
	op := newQuitOperationFromSaver(n, m)

	cmd := op.Update(quitMsg{state: quitStateSaving})
	resultMsg := execCmd(cmd)
	resultCmd := op.Update(resultMsg)
	msg := execCmd(resultCmd)
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("após ErrModifiedExternally: esperado OpenModalMsg, obteve %T", msg)
	}
	if n.lastMethod != "Clear" {
		t.Errorf("após ErrModifiedExternally: esperado Clear no notifier, obteve %q", n.lastMethod)
	}
}

func TestQuitOperation_Update_Saving_ErroGenerico_NaoSai(t *testing.T) {
	n := &stubNotifier{}
	m := &stubManager{salvarErr: errors.New("disco cheio")}
	op := newQuitOperationFromSaver(n, m)

	cmd := op.Update(quitMsg{state: quitStateSaving})
	resultMsg := execCmd(cmd)
	resultCmd := op.Update(resultMsg)
	msg := execCmd(resultCmd)
	if n.lastMethod != "SetError" {
		t.Errorf("após erro genérico: esperado SetError, obteve %q", n.lastMethod)
	}
	if _, ok := msg.(tui.OperationCompletedMsg); !ok {
		t.Errorf("após erro genérico: esperado OperationCompletedMsg, obteve %T", msg)
	}
}

func TestQuitOperation_Update_SavingForced_Sucesso_EmiteQuit(t *testing.T) {
	n := &stubNotifier{}
	m := &stubManager{}
	op := newQuitOperationFromSaver(n, m)

	cmd := op.Update(quitMsg{state: quitStateSavingForced})
	resultMsg := execCmd(cmd)
	resultCmd := op.Update(resultMsg)
	msg := execCmd(resultCmd)
	if !m.salvarForcarSobrescritaCalled {
		t.Error("SavingForced: Salvar(true) não foi chamado")
	}
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Errorf("SavingForced OK: esperado tea.QuitMsg, obteve %T", msg)
	}
}

func TestQuitOperation_Update_SavingForced_Erro_NaoSai(t *testing.T) {
	n := &stubNotifier{}
	m := &stubManager{salvarErr: errors.New("falha")}
	op := newQuitOperationFromSaver(n, m)

	cmd := op.Update(quitMsg{state: quitStateSavingForced})
	resultMsg := execCmd(cmd)
	resultCmd := op.Update(resultMsg)
	msg := execCmd(resultCmd)
	if n.lastMethod != "SetError" {
		t.Errorf("SavingForced erro: esperado SetError, obteve %q", n.lastMethod)
	}
	if _, ok := msg.(tui.OperationCompletedMsg); !ok {
		t.Errorf("SavingForced erro: esperado OperationCompletedMsg, obteve %T", msg)
	}
}
