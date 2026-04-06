package tui

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
)

// --- Test stubs ---

// stubModal implements modalView for tests.
type stubModal struct {
	updateCalls   int
	setSizeCalled bool
	received      tea.Msg
}

func (s *stubModal) Update(msg tea.Msg) tea.Cmd {
	s.updateCalls++
	s.received = msg
	return nil
}
func (s *stubModal) View() string          { return "stub" }
func (s *stubModal) Shortcuts() []Shortcut { return nil }
func (s *stubModal) SetSize(w, h int)      {}

// stubFlow implements flowHandler for tests.
type stubFlow struct {
	initCalled   bool
	updateCalled bool
	received     tea.Msg
}

func (f *stubFlow) Init() tea.Cmd {
	f.initCalled = true
	return nil
}
func (f *stubFlow) Update(msg tea.Msg) tea.Cmd {
	f.updateCalled = true
	f.received = msg
	return nil
}

// stubResult implements modalResult for routing tests.
type stubResult struct{}

func (stubResult) isModalResult() {}

// --- Tests ---

// TestRootModelInit verifies rootModel starts in the correct initial state (D-11).
func TestRootModelInit(t *testing.T) {
	m := NewRootModel()
	if m == nil {
		t.Fatal("NewRootModel returned nil")
	}
	if m.area != workAreaWelcome {
		t.Errorf("expected workAreaWelcome, got %d", m.area)
	}
	if m.welcome == nil {
		t.Error("welcomeModel should be non-nil after construction")
	}
	if len(m.modals) != 0 {
		t.Errorf("expected 0 modals, got %d", len(m.modals))
	}
	// Init() now starts global tick (D-10) — returns non-nil cmd
	if cmd := m.Init(); cmd == nil {
		t.Error("Init() must return a tick cmd in PoC mode")
	}
}

// TestModalStack_PushPop verifies modal stack grows/shrinks correctly.
func TestModalStack_PushPop(t *testing.T) {
	m := NewRootModel()

	modal1 := &stubModal{}
	modal2 := &stubModal{}

	m.Update(pushModalMsg{modal: modal1})
	if len(m.modals) != 1 {
		t.Errorf("after push 1: expected 1 modal, got %d", len(m.modals))
	}

	m.Update(pushModalMsg{modal: modal2})
	if len(m.modals) != 2 {
		t.Errorf("after push 2: expected 2 modals, got %d", len(m.modals))
	}

	m.Update(popModalMsg{})
	if len(m.modals) != 1 {
		t.Errorf("after pop 1: expected 1 modal, got %d", len(m.modals))
	}

	m.Update(popModalMsg{})
	if len(m.modals) != 0 {
		t.Errorf("after pop 2: expected 0 modals, got %d", len(m.modals))
	}

	// Extra pop on empty stack must not panic.
	m.Update(popModalMsg{})
	if len(m.modals) != 0 {
		t.Errorf("after extra pop: expected 0 modals, got %d", len(m.modals))
	}
}

// TestLiveWorkChildren_NilSafety verifies that nil concrete pointer fields do not appear
// as typed-nil interface values in liveWorkChildren() (Go typed-nil trap prevention).
func TestLiveWorkChildren_NilSafety(t *testing.T) {
	m := NewRootModel()

	// Nil out the only active child.
	m.welcome = nil
	live := m.liveWorkChildren()
	if len(live) != 0 {
		t.Errorf("expected 0 live children after nil'ing welcome, got %d", len(live))
	}

	// Restore welcome.
	m.welcome = newWelcomeModel(m.actions)
	live = m.liveWorkChildren()
	if len(live) != 1 {
		t.Errorf("expected 1 live child after restoring welcome, got %d", len(live))
	}

	// None of the returned interfaces must be nil (typed-nil trap).
	for i, child := range live {
		if child == nil {
			t.Errorf("liveWorkChildren()[%d] is nil interface - typed-nil trap!", i)
		}
	}
}

// TestStartEndFlow verifies startFlowMsg sets activeFlow and calls Init(),
// and endFlowMsg clears activeFlow (D-08).
func TestStartEndFlow(t *testing.T) {
	m := NewRootModel()
	flow := &stubFlow{}

	if m.activeFlow != nil {
		t.Fatal("expected activeFlow == nil initially")
	}

	m.Update(startFlowMsg{flow: flow})

	if m.activeFlow == nil {
		t.Fatal("activeFlow should be set after startFlowMsg")
	}
	if !flow.initCalled {
		t.Error("Init() should be called immediately after startFlowMsg")
	}

	m.Update(endFlowMsg{})

	if m.activeFlow != nil {
		t.Error("activeFlow should be nil after endFlowMsg")
	}
}

// TestModalResultRouting verifies that modalResult messages route exclusively
// to activeFlow and are silently dropped when no flow is active (D-03).
func TestModalResultRouting(t *testing.T) {
	m := NewRootModel()
	flow := &stubFlow{}
	result := stubResult{}

	// With no active flow, modalResult should be silently dropped (no panic).
	m.Update(result)

	// Set active flow.
	m.Update(startFlowMsg{flow: flow})
	flow.initCalled = false // reset tracking

	// Send modalResult - should reach flow.Update.
	m.Update(result)
	if !flow.updateCalled {
		t.Error("activeFlow.Update should be called when modalResult is dispatched to it")
	}
	if flow.received != result {
		t.Error("activeFlow.Update received wrong message")
	}
}

// TestWindowSizeMsg_NoModalSetSize verifies that modals do NOT receive SetSize
// when tea.WindowSizeMsg is processed (modals are position-unaware per D-02).
func TestWindowSizeMsg_NoModalSetSize(t *testing.T) {
	m := NewRootModel()
	modal := &stubModal{}

	// Push a modal onto the stack.
	m.Update(pushModalMsg{modal: modal})
	if len(m.modals) != 1 {
		t.Fatal("expected modal on stack")
	}

	// Send window size message.
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// The modal must NOT have received the WindowSizeMsg.
	// WindowSizeMsg is handled by rootModel directly via liveWorkChildren().
	if modal.updateCalls > 0 {
		t.Errorf("modal.Update() was called %d time(s) on WindowSizeMsg - modals must not receive SetSize", modal.updateCalls)
	}
	if m.width != 80 || m.height != 24 {
		t.Errorf("rootModel size not updated: got %dx%d", m.width, m.height)
	}
	_ = time.Now() // keep time import
}

// makeKeyPress creates a tea.KeyPressMsg for testing.
func makeKeyPress(key string) tea.KeyPressMsg {
	switch key {
	case "ctrl+q":
		return tea.KeyPressMsg{Code: 'q', Mod: tea.ModCtrl}
	case "f1":
		return tea.KeyPressMsg{Code: tea.KeyF1}
	default:
		if len(key) == 1 {
			return tea.KeyPressMsg{Code: rune(key[0])}
		}
		return tea.KeyPressMsg{}
	}
}

// TestStartFlow_ClearsOrphanModals verifies that startFlowMsg resets the modal stack,
// preventing orphan modals from a previous flow leaking into the new one (D-08).
func TestStartFlow_ClearsOrphanModals(t *testing.T) {
	m := NewRootModel()
	m.Update(pushModalMsg{modal: &stubModal{}})
	m.Update(pushModalMsg{modal: &stubModal{}})
	if len(m.modals) != 2 {
		t.Fatal("precondition: expected 2 modals on stack")
	}

	flow := &stubFlow{}
	m.Update(startFlowMsg{flow: flow})

	if len(m.modals) != 0 {
		t.Errorf("startFlowMsg must clear orphan modals: expected 0, got %d", len(m.modals))
	}
}

// TestBroadcast_ReachesModals verifies that domain messages sent through the
// broadcast path are forwarded to active modals in addition to work-area children.
func TestBroadcast_ReachesModals(t *testing.T) {
	m := NewRootModel()
	modal := &stubModal{}
	m.Update(pushModalMsg{modal: modal})

	m.Update(vaultSavedMsg{})

	if modal.updateCalls == 0 {
		t.Error("modals must receive broadcast domain messages")
	}
}

// TestD09_ActionManagerBeforeModal verifies that a ScopeGlobal action registered
// on ActionManager fires BEFORE the topmost modal receives the key (D-09 priority #1 > #2).
func TestD09_ActionManagerBeforeModal(t *testing.T) {
	m := NewRootModel()
	modal := &stubModal{}
	m.Update(pushModalMsg{modal: modal})

	// "f1" is registered by newRootModel as ScopeGlobal.
	// ActionManager must consume it; the modal must never see the key.
	m.Update(makeKeyPress("f1"))

	if modal.received != nil {
		t.Error("D-09: ActionManager must consume ScopeGlobal key before topmost modal receives it")
	}
}

// TestKeyPress_CallsHandleInput verifies that messages.HandleInput is called for
// every KeyPressMsg, even when no registered action matches the key.
func TestKeyPress_CallsHandleInput(t *testing.T) {
	m := NewRootModel()
	m.messages.Show(MsgHint, "dismiss me", 0, true) // clearOnInput=true

	m.Update(makeKeyPress("z")) // unknown key, no action registered

	if m.messages.Current() != nil {
		t.Error("messages.HandleInput() must be called on every KeyPressMsg, including unhandled keys")
	}
}

// TestDecisionDialog_ModalStackIntegration verifies that:
// 1. An Acknowledge cmd pushes a DecisionDialog onto the modal stack.
// 2. View() renders without panic with a modal present.
// 3. Enter triggers popModalMsg, clearing the modal stack.
func TestDecisionDialog_ModalStackIntegration(t *testing.T) {
	m := NewRootModel()
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Push DecisionDialog via its Cmd — uses Acknowledge (ℹ Cofre criado)
	cmd := Acknowledge(SeverityInformative, "Cofre criado", "O cofre foi criado com sucesso em ~/documentos/pessoal.abditum.", nil)
	msg := cmd()
	m.Update(msg)

	if len(m.modals) != 1 {
		t.Fatalf("expected 1 modal after push, got %d", len(m.modals))
	}

	// Verify View() doesn't panic with modal on stack
	_ = m.View()

	// Send Enter — should trigger popModalMsg from DecisionDialog.
	// DecisionDialog.Update returns a cmd that returns popModalMsg when executed.
	// We must execute that cmd and feed the result back to rootModel.Update.
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24}) // ensure sized
	_, popCmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if popCmd != nil {
		popMsg := popCmd()
		m.Update(popMsg)
	}

	if len(m.modals) != 0 {
		t.Errorf("expected modal stack empty after Enter, got %d modal(s)", len(m.modals))
	}
}
