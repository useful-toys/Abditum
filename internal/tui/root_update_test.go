package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/vault"
)

// stubOperation é uma Operation mínima para uso nos testes do RootModel.
type stubOperation struct {
	initCalled    bool
	updateCalled  bool
	lastMsg       tea.Msg
	updateReturns tea.Cmd
}

func (s *stubOperation) Init() tea.Cmd {
	s.initCalled = true
	return nil
}

func (s *stubOperation) Update(msg tea.Msg) tea.Cmd {
	s.updateCalled = true
	s.lastMsg = msg
	return s.updateReturns
}

func TestRootModel_StartOperationMsg_SetsActiveOperationAndCallsInit(t *testing.T) {
	r := NewRootModel()
	op := &stubOperation{}

	_, _ = r.Update(StartOperationMsg{Op: op})

	if r.activeOperation != op {
		t.Error("StartOperationMsg: activeOperation não foi setado")
	}
	if !op.initCalled {
		t.Error("StartOperationMsg: Init() não foi chamado")
	}
}

func TestRootModel_StartOperationMsg_ReplacesExistingOperation(t *testing.T) {
	r := NewRootModel()
	op1 := &stubOperation{}
	op2 := &stubOperation{}

	_, _ = r.Update(StartOperationMsg{Op: op1})
	_, _ = r.Update(StartOperationMsg{Op: op2})

	if r.activeOperation != op2 {
		t.Error("StartOperationMsg: deve substituir a operação anterior sem encerrar")
	}
}

func TestRootModel_OperationCompletedMsg_ClearsActiveOperation(t *testing.T) {
	r := NewRootModel()
	op := &stubOperation{}
	r.activeOperation = op

	_, _ = r.Update(OperationCompletedMsg{})

	if r.activeOperation != nil {
		t.Error("OperationCompletedMsg: activeOperation deveria ser nil")
	}
}

type unknownMsg struct{}

func TestRootModel_UnknownMsg_RoutedToActiveOperation(t *testing.T) {
	r := NewRootModel()
	op := &stubOperation{}
	r.activeOperation = op

	_, _ = r.Update(unknownMsg{})

	if !op.updateCalled {
		t.Error("mensagem desconhecida: op.Update não foi chamado")
	}
	if _, ok := op.lastMsg.(unknownMsg); !ok {
		t.Errorf("mensagem desconhecida: op.Update recebeu %T, esperado unknownMsg", op.lastMsg)
	}
}

func TestRootModel_UnknownMsg_RoutedToOperation_EvenWhenModalActive(t *testing.T) {
	r := NewRootModel()
	op := &stubOperation{}
	r.activeOperation = op
	r.modals = append(r.modals, &stubModal{})

	_, _ = r.Update(unknownMsg{})

	if !op.updateCalled {
		t.Errorf("mensagem desconhecida com modal ativo: op.Update não foi chamado")
	}
}

// stubModal implementa ModalView para os testes do RootModel.
type stubModal struct{}

func (s *stubModal) Render(_ int, _ int, _ *design.Theme) string { return "" }
func (s *stubModal) HandleKey(_ tea.KeyMsg) tea.Cmd              { return nil }
func (s *stubModal) Update(_ tea.Msg) tea.Cmd                    { return nil }

func TestRootModel_VaultOpenedMsg_SetsManagerAndWorkArea(t *testing.T) {
	r := NewRootModel()
	mgr := &vault.Manager{}

	_, _ = r.Update(VaultOpenedMsg{Manager: mgr})

	if r.vaultManager != mgr {
		t.Error("VaultOpenedMsg: vaultManager não foi setado")
	}
	if r.workArea != design.WorkAreaVault {
		t.Errorf("VaultOpenedMsg: workArea = %v, esperado WorkAreaVault", r.workArea)
	}
}
