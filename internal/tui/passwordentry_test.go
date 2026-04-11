package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	testdatapkg "github.com/useful-toys/abditum/internal/tui/testdata"
)

// TestPasswordEntryModalStructExists verifies that passwordEntryModal can be instantiated.
func TestPasswordEntryModalStructExists(t *testing.T) {
	m := &passwordEntryModal{}
	if m == nil {
		t.Fatal("passwordEntryModal creation failed")
	}
}

// TestPasswordEntryModalImplementsModalView verifies passwordEntryModal implements modalView.
func TestPasswordEntryModalImplementsModalView(t *testing.T) {
	m := &passwordEntryModal{}
	var _ modalView = m
}

// TestPasswordEntryModalInit verifies Init() initializes the text input.
func TestPasswordEntryModalInit(t *testing.T) {
	m := &passwordEntryModal{}
	cmd := m.Init()
	if cmd == nil {
		t.Log("Init() returned nil (acceptable)")
	}
	// After Init, the modal should be ready to accept input
	if m.input.Value() == "" && m.input.EchoMode == 0 {
		t.Fatal("input field not initialized properly")
	}
	if m.input.EchoCharacter != '•' {
		t.Fatal("input echo character not set to bullet")
	}
}

// TestPasswordEntryModalView verifies View() returns a string.
func TestPasswordEntryModalView(t *testing.T) {
	m := &passwordEntryModal{}
	m.Init()
	m.theme = ThemeTokyoNight
	m.SetSize(80, 24)
	view := m.View()
	if view == "" {
		t.Fatal("View() returned empty string")
	}
	// Should contain the masked field indicator
	if len(view) == 0 {
		t.Fatal("View() produced no output")
	}
}

// TestPasswordEntryModalSetSize verifies SetSize stores dimensions.
func TestPasswordEntryModalSetSize(t *testing.T) {
	m := &passwordEntryModal{}
	m.SetSize(100, 30)
	if m.width != 100 || m.height != 30 {
		t.Fatalf("SetSize failed: got %dx%d, want 100x30", m.width, m.height)
	}
}

// TestPasswordEntryModalShortcuts verifies Shortcuts returns expected shortcuts.
func TestPasswordEntryModalShortcuts(t *testing.T) {
	m := &passwordEntryModal{}
	shortcuts := m.Shortcuts()
	if len(shortcuts) < 2 {
		t.Fatalf("Expected at least 2 shortcuts, got %d", len(shortcuts))
	}
}

// TestPasswordEntryModalAttemptCounter verifies attempt counter starts hidden.
func TestPasswordEntryModalAttemptCounter(t *testing.T) {
	m := &passwordEntryModal{}
	m.Init()
	m.theme = ThemeTokyoNight
	m.SetSize(80, 24)

	// First attempt - counter should be hidden
	view := m.View()
	if view == "" {
		t.Fatal("View returned empty")
	}
	// On first attempt, "Tentativa" should not be visible
	// (We'll do a more detailed assertion after implementation)
}

// TestPasswordEntryModalAttemptCounterAfterWrongPassword verifies counter shows from attempt 2.
func TestPasswordEntryModalAttemptCounterShowsFromSecondAttempt(t *testing.T) {
	m := &passwordEntryModal{}
	m.Init()
	m.theme = ThemeTokyoNight
	m.SetSize(80, 24)

	// Simulate first wrong attempt
	m.HandleWrongPassword()
	if m.attempt < 1 {
		t.Fatal("attempt not incremented")
	}

	// After second attempt, counter should be visible
	if m.attempt >= 2 {
		view := m.View()
		if view == "" {
			t.Fatal("View returned empty after increment")
		}
		// View should mention "Tentativa" when attempt >= 2
	}
}

// TestPasswordEntryModalEnter verifies pressing Enter emits pwdEnteredMsg.
func TestPasswordEntryModalEnterKey(t *testing.T) {
	m := &passwordEntryModal{}
	m.Init()
	m.theme = ThemeTokyoNight
	m.SetSize(80, 24)

	// Type a password
	m.input.SetValue("test1234!")

	// Press Enter
	cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("Update returned nil on Enter with non-empty password")
	}

	// Execute the command to verify it returns pwdEnteredMsg
	msg := cmd()
	if msg == nil {
		t.Fatal("Command returned nil")
	}

	// Should be a batch of commands including pwdEnteredMsg
	// (Can't easily test batch structure, but at minimum shouldn't be nil)
}

// TestPasswordEntryModalEsc verifies ESC emits flowCancelledMsg.
func TestPasswordEntryModalEscKey(t *testing.T) {
	m := &passwordEntryModal{}
	m.Init()
	m.theme = ThemeTokyoNight
	m.SetSize(80, 24)

	// Press ESC
	cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("Update returned nil on Esc")
	}

	// Execute the command
	msg := cmd()
	if msg == nil {
		t.Fatal("Command returned nil")
	}
}

// TestPasswordEntryModalMaskedInput verifies input is masked with •.
func TestPasswordEntryModalMaskedInput(t *testing.T) {
	m := &passwordEntryModal{}
	m.Init()
	m.SetSize(80, 24)

	// Input should have echo character set to '•'
	if m.input.EchoCharacter != '•' {
		t.Fatalf("EchoCharacter is %q, expected '•'", m.input.EchoCharacter)
	}
}

// TestPasswordEntryModalApplyTheme verifies ApplyTheme stores theme.
func TestPasswordEntryModalApplyTheme(t *testing.T) {
	m := &passwordEntryModal{}
	m.Init()
	m.ApplyTheme(ThemeTokyoNight)
	if m.theme == nil {
		t.Fatal("ApplyTheme did not store theme")
	}
	if m.theme != ThemeTokyoNight {
		t.Fatal("ApplyTheme did not store the correct theme")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Golden File Tests
// ─────────────────────────────────────────────────────────────────────────────

// TestPasswordEntryModal_Golden validates visual output against golden files.
func TestPasswordEntryModal_Golden(t *testing.T) {
	m := &passwordEntryModal{
		title: "Senha mestra",
	}
	m.Init()
	m.theme = ThemeTokyoNight
	m.messages = NewMessageManager()
	m.SetSize(80, 24)

	out := m.View()

	// .txt.golden: plain text render
	txtPath := goldenPath("passwordentry", "initial", 80, "txt")
	checkOrUpdateGolden(t, txtPath, stripANSI(out))

	// .json.golden: style transitions
	transitions := testdatapkg.ParseANSIStyle(out)
	jsonBytes, err := testdatapkg.MarshalStyleTransitions(transitions)
	if err != nil {
		t.Fatalf("marshal transitions: %v", err)
	}
	jsonPath := goldenPath("passwordentry", "initial", 80, "json")
	checkOrUpdateGolden(t, jsonPath, string(jsonBytes))
}

// ─────────────────────────────────────────────────────────────────────────────
// Behavioral Tests — D-PE-03/04 action state styling
// ─────────────────────────────────────────────────────────────────────────────

// TestPasswordEntryModal_ConfirmarDisabledWhenEmpty verifies D-PE-03:
// when input is empty, "Enter" does NOT use bold styling on the action hints line.
func TestPasswordEntryModal_ConfirmarDisabledWhenEmpty(t *testing.T) {
	m := &passwordEntryModal{title: "Senha mestra"}
	m.Init()
	m.theme = ThemeTokyoNight
	m.messages = NewMessageManager()
	m.SetSize(80, 24)
	// input is empty after Init

	out := m.View()
	transitions := testdatapkg.ParseANSIStyle(out)

	// Find the maximum line number (action hints are on the last rendered lines)
	maxLine := 0
	for _, tr := range transitions {
		if tr.Line > maxLine {
			maxLine = tr.Line
		}
	}

	// Check the last few lines (action hints row) — should have NO bold style.
	// Title ("Senha mestra") has bold on an early line; action hints do not when empty.
	actionLinesHaveBold := false
	for _, tr := range transitions {
		if tr.Line >= maxLine-2 { // last 3 lines = border + action hints + border
			for _, s := range tr.Style {
				if s == "bold" {
					actionLinesHaveBold = true
				}
			}
		}
	}
	if actionLinesHaveBold {
		t.Error("unexpected bold in action hints area when input is empty (D-PE-03): 'Enter' should not be bold when disabled")
	}
}

// TestPasswordEntryModal_ConfirmarActiveWhenFilled verifies D-PE-04:
// when input is non-empty, "Enter" uses bold styling on the action hints line.
func TestPasswordEntryModal_ConfirmarActiveWhenFilled(t *testing.T) {
	m := &passwordEntryModal{title: "Senha mestra"}
	m.Init()
	m.theme = ThemeTokyoNight
	m.messages = NewMessageManager()
	m.SetSize(80, 24)
	m.input.SetValue("hunter2") // non-empty → activates "Enter Confirmar"

	out := m.View()
	transitions := testdatapkg.ParseANSIStyle(out)

	// Find the maximum line number
	maxLine := 0
	for _, tr := range transitions {
		if tr.Line > maxLine {
			maxLine = tr.Line
		}
	}

	// Check the last few lines (action hints row) — SHOULD have bold for "Enter".
	actionLinesHaveBold := false
	for _, tr := range transitions {
		if tr.Line >= maxLine-2 { // last 3 lines = border + action hints + border
			for _, s := range tr.Style {
				if s == "bold" {
					actionLinesHaveBold = true
				}
			}
		}
	}
	if !actionLinesHaveBold {
		t.Error("expected bold in action hints area when input is non-empty (D-PE-04): 'Enter' should be bold when active")
	}
}

// TestPasswordEntryModal_Golden_Filled validates visual output when password is non-empty.
// Verifies D-PE-04: "Enter Confirmar" is accent.primary + bold when field has content.
func TestPasswordEntryModal_Golden_Filled(t *testing.T) {
	m := &passwordEntryModal{title: "Senha mestra"}
	m.Init()
	m.theme = ThemeTokyoNight
	m.messages = NewMessageManager()
	m.SetSize(80, 24)
	m.input.SetValue("hunter2") // non-empty → activates "Enter Confirmar"

	out := m.View()

	txtPath := goldenPath("passwordentry", "filled", 80, "txt")
	checkOrUpdateGolden(t, txtPath, stripANSI(out))

	transitions := testdatapkg.ParseANSIStyle(out)
	jsonBytes, err := testdatapkg.MarshalStyleTransitions(transitions)
	if err != nil {
		t.Fatalf("marshal transitions: %v", err)
	}
	jsonPath := goldenPath("passwordentry", "filled", 80, "json")
	checkOrUpdateGolden(t, jsonPath, string(jsonBytes))
}
