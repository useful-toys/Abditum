package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

// TestRootModelInit verifies rootModel starts in the correct initial state.
func TestRootModelInit(t *testing.T) {
	m := newRootModel(nil, "")

	if m == nil {
		t.Fatal("newRootModel returned nil")
	}
	if m.area != workAreaPreVault {
		t.Errorf("expected workAreaPreVault, got %d", m.area)
	}
	if m.preVault == nil {
		t.Error("preVault should be non-nil after construction")
	}
	if len(m.modals) != 0 {
		t.Errorf("expected 0 modals, got %d", len(m.modals))
	}
	if cmd := m.Init(); cmd != nil {
		t.Error("Init() must return nil — tick must not start before workAreaVault")
	}
}

// TestModalStack_PushPop verifies the modal stack grows and shrinks correctly.
func TestModalStack_PushPop(t *testing.T) {
	m := newRootModel(nil, "")

	modal1 := newModal("title1", "body1", nil, nil)
	modal2 := newModal("title2", "body2", nil, nil)

	// Push first modal
	m.Update(pushModalMsg{modal: modal1})
	if len(m.modals) != 1 {
		t.Errorf("after push 1: expected 1 modal, got %d", len(m.modals))
	}

	// Push second modal
	m.Update(pushModalMsg{modal: modal2})
	if len(m.modals) != 2 {
		t.Errorf("after push 2: expected 2 modals, got %d", len(m.modals))
	}

	// Pop one
	m.Update(popModalMsg{})
	if len(m.modals) != 1 {
		t.Errorf("after pop 1: expected 1 modal, got %d", len(m.modals))
	}

	// Pop all
	m.Update(popModalMsg{})
	if len(m.modals) != 0 {
		t.Errorf("after pop 2: expected 0 modals, got %d", len(m.modals))
	}

	// Extra pop on empty stack must not panic
	m.Update(popModalMsg{})
	if len(m.modals) != 0 {
		t.Errorf("after extra pop: expected 0 modals, got %d", len(m.modals))
	}
}

// TestLiveModels_TypedNilSafety verifies that nil concrete pointer fields
// do NOT appear in liveModels() as typed-nil interface values.
func TestLiveModels_TypedNilSafety(t *testing.T) {
	m := newRootModel(nil, "")

	// Nil out preVault (only active child at this point)
	m.preVault = nil
	live := m.liveModels()
	if len(live) != 0 {
		t.Errorf("expected 0 live models after nil'ing preVault, got %d", len(live))
	}

	// Restore preVault
	m.preVault = newPreVaultModel(m.actions)
	live = m.liveModels()
	if len(live) != 1 {
		t.Errorf("expected 1 live model after restoring preVault, got %d", len(live))
	}

	// Verify none of the returned interfaces are nil (typed-nil check)
	for i, child := range live {
		if child == nil {
			t.Errorf("liveModels()[%d] is nil interface — typed-nil trap!", i)
		}
	}
}

// TestDispatchPriority_CtrlQ verifies ctrl+Q is intercepted before any child.
func TestDispatchPriority_CtrlQ(t *testing.T) {
	m := newRootModel(nil, "") // nil mgr → IsModified() not called (mgr nil check in dispatchKey)

	// Simulate ctrl+q key press
	_, cmd := m.Update(makeKeyPress("ctrl+q"))
	if cmd == nil {
		t.Error("ctrl+q should return a non-nil Cmd (tea.Quit or confirm dialog)")
	}
}

// TestWindowSizeMsg_PropagatesToChildren verifies SetSize is called on live children.
func TestWindowSizeMsg_PropagatesToChildren(t *testing.T) {
	m := newRootModel(nil, "")
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	if m.width != 80 || m.height != 24 {
		t.Errorf("rootModel size not updated: got %dx%d", m.width, m.height)
	}
	// preVault should have received the size (full width is passed via liveModels + SetSize)
	if m.preVault != nil && m.preVault.width != 80 {
		// Note: height passed to child via SetSize in liveModels broadcast,
		// but renderFrame also calls SetSize(m.width, workH) on preVault.
		// In WindowSizeMsg handling, SetSize(msg.Width, msg.Height) is called via liveModels.
		t.Errorf("preVault width not updated: got %d", m.preVault.width)
	}
}

// makeKeyPress creates a tea.KeyPressMsg that String() returns the given key string.
// Supports common key combinations used in rootModel dispatch.
func makeKeyPress(key string) tea.KeyPressMsg {
	switch key {
	case "ctrl+q":
		// tea.KeyPressMsg is type alias for tea.Key
		// Code = 'q', Mod = ModCtrl → String() returns "ctrl+q"
		return tea.KeyPressMsg{Code: 'q', Mod: tea.ModCtrl}
	case "ctrl+c":
		return tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}
	case "enter":
		return tea.KeyPressMsg{Code: tea.KeyEnter}
	case "esc":
		return tea.KeyPressMsg{Code: tea.KeyEsc}
	case "?":
		return tea.KeyPressMsg{Code: '?', Text: "?"}
	default:
		// For single character keys, set both Code and Text
		if len(key) == 1 {
			return tea.KeyPressMsg{Code: rune(key[0]), Text: key}
		}
		return tea.KeyPressMsg{}
	}
}
