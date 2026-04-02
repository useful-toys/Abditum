package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

// TestActionManager_DispatchScopeLocal verifies that ScopeLocal actions are
// skipped when inFlowOrModal=true and fire when inFlowOrModal=false.
func TestActionManager_DispatchScopeLocal(t *testing.T) {
	am := NewActionManager()
	called := false
	am.Register("owner", Action{
		Keys:    []string{"q"},
		Scope:   ScopeLocal,
		Enabled: func() bool { return true },
		Handler: func() tea.Cmd { called = true; return nil },
	})

	// inFlowOrModal=true → ScopeLocal must NOT fire.
	cmd := am.Dispatch("q", true)
	if cmd != nil || called {
		t.Error("ScopeLocal action must not dispatch when inFlowOrModal=true")
	}

	// inFlowOrModal=false → ScopeLocal must fire.
	am.Dispatch("q", false)
	if !called {
		t.Error("ScopeLocal action must dispatch when inFlowOrModal=false")
	}
}

// TestActionManager_DispatchScopeGlobal verifies that ScopeGlobal actions fire
// regardless of whether a flow or modal is active.
func TestActionManager_DispatchScopeGlobal(t *testing.T) {
	am := NewActionManager()
	called := false
	am.Register("owner", Action{
		Keys:    []string{"?"},
		Scope:   ScopeGlobal,
		Enabled: func() bool { return true },
		Handler: func() tea.Cmd { called = true; return nil },
	})

	am.Dispatch("?", true) // in flow/modal
	if !called {
		t.Error("ScopeGlobal action must dispatch even when inFlowOrModal=true")
	}
}

// TestActionManager_DispatchEnabled verifies that an action whose Enabled()
// returns false is never dispatched.
func TestActionManager_DispatchEnabled(t *testing.T) {
	am := NewActionManager()
	called := false
	am.Register("owner", Action{
		Keys:    []string{"x"},
		Scope:   ScopeGlobal,
		Enabled: func() bool { return false },
		Handler: func() tea.Cmd { called = true; return nil },
	})

	am.Dispatch("x", false)
	if called {
		t.Error("action with Enabled()=false must not dispatch")
	}
}

// TestActionManager_ActiveOwnerPriority verifies that the active owner's actions
// are checked before other registered owners.
func TestActionManager_ActiveOwnerPriority(t *testing.T) {
	am := NewActionManager()
	winner := ""
	am.Register("A", Action{
		Keys:    []string{"enter"},
		Scope:   ScopeGlobal,
		Enabled: func() bool { return true },
		Handler: func() tea.Cmd { winner = "A"; return nil },
	})
	am.Register("B", Action{
		Keys:    []string{"enter"},
		Scope:   ScopeGlobal,
		Enabled: func() bool { return true },
		Handler: func() tea.Cmd { winner = "B"; return nil },
	})

	am.SetActiveOwner("B")
	am.Dispatch("enter", false)
	if winner != "B" {
		t.Errorf("expected active owner 'B' to win dispatch, got %q", winner)
	}
}

// TestActionManager_ClearOwned removes all actions for the given owner.
func TestActionManager_ClearOwned(t *testing.T) {
	am := NewActionManager()
	called := false
	owner := "owner"
	am.Register(owner, Action{
		Keys:    []string{"x"},
		Scope:   ScopeGlobal,
		Enabled: func() bool { return true },
		Handler: func() tea.Cmd { called = true; return nil },
	})

	am.ClearOwned(owner)
	am.Dispatch("x", false)
	if called {
		t.Error("action must not dispatch after ClearOwned")
	}
	if len(am.Visible()) != 0 {
		t.Error("Visible() must return empty after ClearOwned removes the only owner")
	}
}

// TestActionManager_Visible returns only actions where Enabled() is true.
func TestActionManager_Visible(t *testing.T) {
	am := NewActionManager()
	am.Register("owner",
		Action{Keys: []string{"a"}, Scope: ScopeGlobal, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		Action{Keys: []string{"b"}, Scope: ScopeGlobal, Enabled: func() bool { return false }, Handler: func() tea.Cmd { return nil }},
	)

	visible := am.Visible()
	if len(visible) != 1 {
		t.Errorf("expected 1 visible action, got %d", len(visible))
	}
	if len(visible) > 0 && visible[0].Keys[0] != "a" {
		t.Errorf("expected 'a' to be visible, got %q", visible[0].Keys[0])
	}
}

// TestActionManager_All returns all actions regardless of Enabled state.
func TestActionManager_All(t *testing.T) {
	am := NewActionManager()
	am.Register("owner",
		Action{Keys: []string{"a"}, Scope: ScopeGlobal, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		Action{Keys: []string{"b"}, Scope: ScopeGlobal, Enabled: func() bool { return false }, Handler: func() tea.Cmd { return nil }},
	)

	all := am.All()
	if len(all) != 2 {
		t.Errorf("All() expected 2 actions (including disabled), got %d", len(all))
	}
}
