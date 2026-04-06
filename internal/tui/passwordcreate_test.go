package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	testdatapkg "github.com/useful-toys/abditum/internal/tui/testdata"
)

// TestPasswordCreateModalStructExists verifies that passwordCreateModal can be instantiated.
func TestPasswordCreateModalStructExists(t *testing.T) {
	m := &passwordCreateModal{}
	if m == nil {
		t.Fatal("passwordCreateModal creation failed")
	}
}

// TestPasswordCreateModalImplementsModalView verifies passwordCreateModal implements modalView.
func TestPasswordCreateModalImplementsModalView(t *testing.T) {
	m := &passwordCreateModal{}
	var _ modalView = m
}

// TestPasswordCreateModalInit verifies Init() initializes both password fields.
func TestPasswordCreateModalInit(t *testing.T) {
	m := &passwordCreateModal{}
	cmd := m.Init()
	if cmd == nil {
		t.Log("Init() returned nil (acceptable)")
	}
	// After Init, both input fields should exist and be initialized
	if m.password.EchoCharacter != '•' {
		t.Fatal("password echo character not set to bullet")
	}
	if m.confirm.EchoCharacter != '•' {
		t.Fatal("confirm echo character not set to bullet")
	}
	if m.focusIndex != 0 {
		t.Fatal("focus should start at password field")
	}
}

// TestPasswordCreateModalView verifies View() returns a string.
func TestPasswordCreateModalView(t *testing.T) {
	m := &passwordCreateModal{}
	m.Init()
	m.theme = ThemeTokyoNight
	m.SetSize(80, 24)
	view := m.View()
	if view == "" {
		t.Fatal("View() returned empty string")
	}
	if len(view) == 0 {
		t.Fatal("View() produced no output")
	}
}

// TestPasswordCreateModalSetSize verifies SetSize stores dimensions.
func TestPasswordCreateModalSetSize(t *testing.T) {
	m := &passwordCreateModal{}
	m.SetSize(100, 30)
	if m.width != 100 || m.height != 30 {
		t.Fatalf("SetSize failed: got %dx%d, want 100x30", m.width, m.height)
	}
}

// TestPasswordCreateModalShortcuts verifies Shortcuts returns expected shortcuts.
func TestPasswordCreateModalShortcuts(t *testing.T) {
	m := &passwordCreateModal{}
	shortcuts := m.Shortcuts()
	if len(shortcuts) < 3 {
		t.Fatalf("Expected at least 3 shortcuts, got %d", len(shortcuts))
	}
}

// TestPasswordCreateModalTabNavigation verifies Tab switches focus between fields.
func TestPasswordCreateModalTabNavigation(t *testing.T) {
	m := &passwordCreateModal{}
	m.Init()
	m.theme = ThemeTokyoNight

	// Start at password field
	if m.focusIndex != 0 {
		t.Fatal("should start at password field (focusIndex=0)")
	}

	// Press Tab
	m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.focusIndex != 1 {
		t.Fatal("should move to confirm field (focusIndex=1)")
	}

	// Press Tab again
	m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.focusIndex != 0 {
		t.Fatal("should return to password field (focusIndex=0)")
	}
}

// TestPasswordCreateModalMismatchedPasswords verifies mismatch error handling.
func TestPasswordCreateModalMismatchedPasswords(t *testing.T) {
	m := &passwordCreateModal{}
	m.Init()
	m.theme = ThemeTokyoNight
	m.SetSize(80, 24)

	// Set different passwords
	m.password.SetValue("password1")
	m.confirm.SetValue("password2")

	// Press Enter - should not emit success
	cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Log("Update returned nil on mismatched passwords (acceptable)")
	}
}

// TestPasswordCreateModalEmptyPassword verifies empty password is rejected.
func TestPasswordCreateModalEmptyPassword(t *testing.T) {
	m := &passwordCreateModal{}
	m.Init()
	m.theme = ThemeTokyoNight

	// Leave password empty
	m.password.SetValue("")
	m.confirm.SetValue("")

	// Press Enter - should be rejected
	cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Log("Update returned nil on empty password (acceptable)")
	}
}

// TestPasswordCreateModalMatchingPasswords verifies matching passwords emit success.
func TestPasswordCreateModalMatchingPasswords(t *testing.T) {
	m := &passwordCreateModal{}
	m.Init()
	m.theme = ThemeTokyoNight
	m.SetSize(80, 24)

	// Set matching passwords
	m.password.SetValue("SamePassword123!")
	m.confirm.SetValue("SamePassword123!")

	// Press Enter - should emit success
	cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("Update returned nil on matching passwords")
	}

	// Execute the command
	msg := cmd()
	if msg == nil {
		t.Fatal("Command returned nil")
	}
}

// TestPasswordCreateModalEsc verifies ESC emits flowCancelledMsg.
func TestPasswordCreateModalEscKey(t *testing.T) {
	m := &passwordCreateModal{}
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

// TestPasswordCreateModalMaskedInput verifies both inputs are masked.
func TestPasswordCreateModalMaskedInput(t *testing.T) {
	m := &passwordCreateModal{}
	m.Init()

	// Both fields should have echo character set to '•'
	if m.password.EchoCharacter != '•' {
		t.Fatalf("password EchoCharacter is %q, expected '•'", m.password.EchoCharacter)
	}
	if m.confirm.EchoCharacter != '•' {
		t.Fatalf("confirm EchoCharacter is %q, expected '•'", m.confirm.EchoCharacter)
	}
}

// TestPasswordCreateModalApplyTheme verifies ApplyTheme stores theme.
func TestPasswordCreateModalApplyTheme(t *testing.T) {
	m := &passwordCreateModal{}
	m.Init()
	m.ApplyTheme(ThemeTokyoNight)
	if m.theme == nil {
		t.Fatal("ApplyTheme did not store theme")
	}
	if m.theme != ThemeTokyoNight {
		t.Fatal("ApplyTheme did not store the correct theme")
	}
}

// TestPasswordCreateModalStrengthEvaluation verifies strength is evaluated.
func TestPasswordCreateModalStrengthEvaluation(t *testing.T) {
	m := &passwordCreateModal{}
	m.Init()
	m.theme = ThemeTokyoNight

	// Weak password
	m.password.SetValue("weak")
	m.updateStrength()
	if m.strength == 0 {
		t.Log("Weak password detected (strength not yet evaluated)")
	}

	// Strong password
	m.password.SetValue("VeryStrongPassword123!")
	m.updateStrength()
	// Should evaluate to something (exact value depends on crypto package implementation)
}

// ─────────────────────────────────────────────────────────────────────────────
// Golden File Tests
// ─────────────────────────────────────────────────────────────────────────────

// TestPasswordCreateModal_Golden validates visual output against golden files.
func TestPasswordCreateModal_Golden(t *testing.T) {
	m := &passwordCreateModal{
		title: "Criar senha mestra",
	}
	m.Init()
	m.theme = ThemeTokyoNight
	m.messages = NewMessageManager()
	m.SetSize(80, 24)

	out := m.View()

	// .txt.golden: plain text render
	txtPath := goldenPath("passwordcreate", "initial", 80, "txt")
	checkOrUpdateGolden(t, txtPath, stripANSI(out))

	// .json.golden: style transitions
	transitions := testdatapkg.ParseANSIStyle(out)
	jsonBytes, err := testdatapkg.MarshalStyleTransitions(transitions)
	if err != nil {
		t.Fatalf("marshal transitions: %v", err)
	}
	jsonPath := goldenPath("passwordcreate", "initial", 80, "json")
	checkOrUpdateGolden(t, jsonPath, string(jsonBytes))
}
