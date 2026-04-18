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

func (s *stubModal) Render(_ int, _ int, _ *design.Theme) string  { return "" }
func (s *stubModal) HandleKey(_ tea.KeyMsg) tea.Cmd               { return nil }
func (s *stubModal) HandleMouse(_ tea.MouseMsg) tea.Cmd           { return nil }
func (s *stubModal) Cursor(_, _ int) *tea.Cursor                  { return nil }

func TestRootModel_VaultOpenedMsg_NilManager_IsIgnored(t *testing.T) {
	r := NewRootModel()

	_, _ = r.Update(VaultOpenedMsg{Manager: nil})

	if r.workArea != design.WorkAreaWelcome {
		t.Errorf("VaultOpenedMsg(nil): workArea não deveria mudar, obteve %v", r.workArea)
	}
}

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

func TestRootModel_TickMsg_ReachesMessageLineView(t *testing.T) {
	r := NewRootModel()
	r.messageLineView.SetBusy("test")
	initialFrame := r.messageLineView.CurrentSpinnerFrame()

	_, _ = r.Update(TickMsg{})

	if r.messageLineView.CurrentSpinnerFrame() == initialFrame {
		t.Error("TickMsg: SpinnerFrame não avançou — TickMsg não chegou ao messageLineView")
	}
}

func TestRootModel_TickMsg_WithActiveOperation_SpinnerStillAnimates(t *testing.T) {
	r := NewRootModel()
	r.messageLineView.SetBusy("test")
	initialFrame := r.messageLineView.CurrentSpinnerFrame()

	// Criar uma operação stub
	op := &stubOperation{}
	r.activeOperation = op

	_, _ = r.Update(TickMsg{})

	// TickMsg deve chegar à operação (se ela processar TickMsg)
	// MAS TAMBÉM deve continuar atualizando o spinner!
	if r.messageLineView.CurrentSpinnerFrame() == initialFrame {
		t.Error("TickMsg com operação ativa: SpinnerFrame não avançou — spinner deveria continuar animando")
	}
}

// TestRootModel_TickMsg_ReturnsCmd verifica que Update(TickMsg) retorna um Cmd não-nil.
// tea.Every() dispara apenas uma vez — sem re-agendar, o spinner para após o primeiro tick.
func TestRootModel_TickMsg_ReturnsCmd(t *testing.T) {
	r := NewRootModel()

	_, cmd := r.Update(TickMsg{})

	if cmd == nil {
		t.Error("Update(TickMsg): retornou nil — ticker não será re-agendado e spinner para")
	}
}

// TestRootModel_TickMsg_SpinnerCycles verifica que múltiplos ticks avançam o frame
// na sequência correta: 0→1→2→3→0, confirmando que o re-agendamento funciona.
func TestRootModel_TickMsg_SpinnerCycles(t *testing.T) {
	r := NewRootModel()
	r.messageLineView.SetBusy("test")

	for tick := 1; tick <= 8; tick++ {
		_, _ = r.Update(TickMsg{})
		want := tick % 4
		got := r.messageLineView.CurrentSpinnerFrame()
		if got != want {
			t.Errorf("após %d ticks: SpinnerFrame = %d, want %d", tick, got, want)
		}
	}
}
