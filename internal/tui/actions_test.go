package tui

import (
	"encoding/json"
	"fmt"
	"testing"

	tea "charm.land/bubbletea/v2"
	testdatapkg "github.com/useful-toys/abditum/internal/tui/testdata"
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

// TestActionManager_Dispatch_MultiKey verifies that an action with multiple keys
// dispatches its handler for any key in the Keys slice, not just Keys[0].
func TestActionManager_Dispatch_MultiKey(t *testing.T) {
	am := NewActionManager()
	called := false
	am.Register("owner", Action{
		Keys:    []string{"f10", "f11"},
		Scope:   ScopeGlobal,
		Enabled: func() bool { return true },
		Handler: func() tea.Cmd { called = true; return nil },
	})

	am.Dispatch("f10", false)
	if !called {
		t.Error("Dispatch(f10) must trigger the handler")
	}

	called = false
	am.Dispatch("f11", false)
	if !called {
		t.Error("Dispatch(f11) must trigger the handler — both f10 and f11 must work")
	}
}

// TestActionManager_Dispatch_MultiKey_InFlowOrModal verifies multi-key dispatch
// works correctly when inFlowOrModal=true (ScopeGlobal actions).
func TestActionManager_Dispatch_MultiKey_InFlowOrModal(t *testing.T) {
	am := NewActionManager()
	called := false
	am.Register("owner", Action{
		Keys:    []string{"f10", "f11"},
		Scope:   ScopeGlobal,
		Enabled: func() bool { return true },
		Handler: func() tea.Cmd { called = true; return nil },
	})

	am.Dispatch("f10", true)
	if !called {
		t.Error("Dispatch(f10, inFlow=true) must trigger the handler")
	}

	called = false
	am.Dispatch("f11", true)
	if !called {
		t.Error("Dispatch(f11, inFlow=true) must trigger the handler — both keys must work when inFlowOrModal=true")
	}
}

// TestActionManager_Visible_HidesFromBar verifies that HideFromBar actions are
// excluded from Visible() but still present in All().
func TestActionManager_Visible_HidesFromBar(t *testing.T) {
	am := NewActionManager()
	am.Register("owner",
		Action{Keys: []string{"a"}, Scope: ScopeGlobal, Enabled: func() bool { return true }, HideFromBar: false, Handler: func() tea.Cmd { return nil }},
		Action{Keys: []string{"b"}, Scope: ScopeGlobal, Enabled: func() bool { return true }, HideFromBar: true, Handler: func() tea.Cmd { return nil }},
	)

	visible := am.Visible()
	if len(visible) != 1 {
		t.Errorf("Visible() expected 1 action (HideFromBar excluded), got %d", len(visible))
	}
	if len(visible) > 0 && visible[0].Keys[0] != "a" {
		t.Errorf("expected 'a' to be visible, got %q", visible[0].Keys[0])
	}

	all := am.All()
	if len(all) != 2 {
		t.Errorf("All() expected 2 actions (including HideFromBar), got %d", len(all))
	}
}

// TestActionManager_Visible_DisabledExcluded verifies that Enabled()==false actions
// are excluded from Visible() but still present in All().
func TestActionManager_Visible_DisabledExcluded(t *testing.T) {
	am := NewActionManager()
	am.Register("owner",
		Action{Keys: []string{"a"}, Scope: ScopeGlobal, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		Action{Keys: []string{"b"}, Scope: ScopeGlobal, Enabled: func() bool { return false }, Handler: func() tea.Cmd { return nil }},
	)

	visible := am.Visible()
	if len(visible) != 1 {
		t.Errorf("Visible() expected 1 action (disabled excluded), got %d", len(visible))
	}
	if len(visible) > 0 && visible[0].Keys[0] != "a" {
		t.Errorf("expected 'a' to be visible, got %q", visible[0].Keys[0])
	}

	all := am.All()
	if len(all) != 2 {
		t.Errorf("All() expected 2 actions (including disabled), got %d", len(all))
	}
}

// TestRenderCommandBar_MultiKeyShowsFirst verifies that RenderCommandBar displays
// only Keys[0] for multi-key actions, not all keys.
func TestRenderCommandBar_MultiKeyShowsFirst(t *testing.T) {
	am := NewActionManager()
	am.Register("owner", Action{
		Keys:     []string{"f10", "f11"},
		Label:    "Truncar",
		Scope:    ScopeGlobal,
		Priority: 90,
		Enabled:  func() bool { return true },
		Handler:  func() tea.Cmd { return nil },
	})

	result := RenderCommandBar(am.Visible(), 80)
	if !containsKey(result, "f10") {
		t.Error("RenderCommandBar must show Keys[0] (f10)")
	}
	if containsKey(result, "f11") {
		t.Error("RenderCommandBar must NOT show Keys[1] (f11) — only Keys[0] is displayed")
	}
}

// TestRenderCommandBar_TruncatesLowestPriority verifies that when the terminal
// is too narrow for all actions, the lowest-priority actions are dropped first
// while the F1 anchor is preserved.
func TestRenderCommandBar_TruncatesLowestPriority(t *testing.T) {
	am := NewActionManager()
	am.Register("owner",
		Action{Keys: []string{"f10"}, Label: "High", Scope: ScopeGlobal, Priority: 90, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		Action{Keys: []string{"f9"}, Label: "Mid", Scope: ScopeGlobal, Priority: 50, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		Action{Keys: []string{"f2"}, Label: "Low", Scope: ScopeGlobal, Priority: 10, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		Action{Keys: []string{"f1"}, Label: "Help", Scope: ScopeGlobal, Priority: 0, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
	)

	// 30 cols: body ("  f10 High · f9 Mid · f2 Low") + anchor ("f1 Help") exceeds width.
	// Truncation removes lowest priority first: f2 goes, f9 stays.
	result := RenderCommandBar(am.Visible(), 30)

	if !containsKey(result, "f10") {
		t.Error("highest priority action (f10) must be kept")
	}
	if containsKey(result, "f2") {
		t.Error("lowest priority action (f2) must be truncated first")
	}
	if !containsKey(result, "f1") {
		t.Error("F1 anchor must always be preserved")
	}
}

// TestRenderCommandBar_NarrowTerminal_F1Only verifies that at extremely narrow
// widths (10 cols), only the F1 anchor remains — all body actions are truncated.
func TestRenderCommandBar_NarrowTerminal_F1Only(t *testing.T) {
	am := NewActionManager()
	am.Register("owner",
		Action{Keys: []string{"f10"}, Label: "High", Scope: ScopeGlobal, Priority: 90, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		Action{Keys: []string{"f9"}, Label: "Mid", Scope: ScopeGlobal, Priority: 50, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		Action{Keys: []string{"f2"}, Label: "Low", Scope: ScopeGlobal, Priority: 10, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		Action{Keys: []string{"f1"}, Label: "Help", Scope: ScopeGlobal, Priority: 0, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
	)

	// 10 cols: body doesn't fit at all, only anchor survives.
	result := RenderCommandBar(am.Visible(), 10)

	if !containsKey(result, "f1") {
		t.Error("at 10 cols, F1 anchor must still be visible")
	}
	if containsKey(result, "f10") || containsKey(result, "f9") || containsKey(result, "f2") {
		t.Error("at 10 cols, all body actions must be truncated — only F1 should remain")
	}
}

// containsKey checks if a rendered command bar string contains a key token.
// Works despite ANSI escape codes from lipgloss styling.
func containsKey(bar, key string) bool {
	return len(bar) > 0 && (bar[0] == '\x1b' || bar[0] == ' ') && containsSubstring(bar, key)
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestRenderCommandBar_Golden validates the visual output of RenderCommandBar
// against golden files for all 5 scenarios × 2 widths (= 10 sub-tests).
//
// Each sub-test produces two golden files:
//   - .txt.golden — raw ANSI output, validated byte-for-byte
//   - .json.golden — style transitions from ParseANSIStyle, validated byte-for-byte
//
// First run (no golden files present) auto-generates the baselines.
// Subsequent runs compare against the recorded baselines.
// Run with -update to intentionally regenerate all baselines.
func TestRenderCommandBar_Golden(t *testing.T) {
	type testCase struct {
		variant string
		actions []Action
	}

	// Scenario 3: typical — F4 Salvar(90), F2 Novo(80), F3 Abrir(70), F5 Deletar(50), F1 Ajuda(0)
	// (Priority desc order: Salvar, Novo, Abrir, Deletar; F1 is right-anchored)
	typicalActions := []Action{
		{Keys: []string{"f2"}, Label: "Novo", Priority: 80, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		{Keys: []string{"f3"}, Label: "Abrir", Priority: 70, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		{Keys: []string{"f4"}, Label: "Salvar", Priority: 90, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		{Keys: []string{"f5"}, Label: "Deletar", Priority: 50, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		{Keys: []string{"f1"}, Label: "Ajuda", Priority: 0, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
	}

	// Scenario 4: unsorted — same actions but passed in reverse priority order.
	// RenderCommandBar must sort them → golden output identical to "typical".
	unsortedActions := []Action{
		{Keys: []string{"f1"}, Label: "Ajuda", Priority: 0, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		{Keys: []string{"f5"}, Label: "Deletar", Priority: 50, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		{Keys: []string{"f4"}, Label: "Salvar", Priority: 90, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		{Keys: []string{"f3"}, Label: "Abrir", Priority: 70, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		{Keys: []string{"f2"}, Label: "Novo", Priority: 80, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
	}

	// Scenario 5: many — 7 body actions + F1 anchor
	manyActions := []Action{
		{Keys: []string{"f2"}, Label: "Novo", Priority: 70, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		{Keys: []string{"f3"}, Label: "Abrir", Priority: 60, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		{Keys: []string{"f4"}, Label: "Salvar", Priority: 50, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		{Keys: []string{"f5"}, Label: "Fechar", Priority: 40, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		{Keys: []string{"f6"}, Label: "Exportar", Priority: 30, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		{Keys: []string{"f7"}, Label: "Importar", Priority: 20, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		{Keys: []string{"f8"}, Label: "Config", Priority: 10, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		{Keys: []string{"f1"}, Label: "Ajuda", Priority: 0, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
	}

	cases := []testCase{
		// Scenario 1: no-actions — empty slice
		{"no-actions", []Action{}},
		// Scenario 2: f1-only — single F1 anchor
		{"f1-only", []Action{
			{Keys: []string{"f1"}, Label: "Ajuda", Priority: 0, Enabled: func() bool { return true }, Handler: func() tea.Cmd { return nil }},
		}},
		// Scenario 3: typical — 5 actions, mixed priorities
		{"typical", typicalActions},
		// Scenario 4: unsorted — same as typical but reverse order (tests sorting)
		{"unsorted", unsortedActions},
		// Scenario 5: many — 8 actions, truncates at narrow widths
		{"many", manyActions},
	}
	widths := []int{30, 60}

	for _, tc := range cases {
		for _, w := range widths {
			tc := tc
			w := w
			name := fmt.Sprintf("%s-%d", tc.variant, w)
			t.Run(name, func(t *testing.T) {
				out := RenderCommandBar(tc.actions, w)

				// .txt.golden: raw ANSI output
				txtPath := goldenPath("commandbar", tc.variant, w, "txt")
				checkOrUpdateGolden(t, txtPath, out)

				// .json.golden: style transitions
				transitions := testdatapkg.ParseANSIStyle(out)
				jsonBytes, err := json.MarshalIndent(transitions, "", "  ")
				if err != nil {
					t.Fatalf("marshal transitions: %v", err)
				}
				jsonPath := goldenPath("commandbar", tc.variant, w, "json")
				checkOrUpdateGolden(t, jsonPath, string(jsonBytes))
			})
		}
	}
}
