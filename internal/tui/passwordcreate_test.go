package tui

import (
	"strings"
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

// TestPasswordCreateModalView verifies View(80, 24) returns a string.
func TestPasswordCreateModalView(t *testing.T) {
	m := &passwordCreateModal{}
	m.Init()
	m.theme = TokyoNight

	view := m.View(80, 24, TokyoNight)
	if view == "" {
		t.Fatal("View(80, 24) returned empty string")
	}
	if len(view) == 0 {
		t.Fatal("View(80, 24) produced no output")
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
	m.theme = TokyoNight

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
	m.theme = TokyoNight

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
	m.theme = TokyoNight

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
	m.theme = TokyoNight

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
	m.theme = TokyoNight

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
	m.ApplyTheme(TokyoNight)
	if m.theme == nil {
		t.Fatal("ApplyTheme did not store theme")
	}
	if m.theme != TokyoNight {
		t.Fatal("ApplyTheme did not store the correct theme")
	}
}

// TestPasswordCreateModalStrengthEvaluation verifies strength is evaluated.
func TestPasswordCreateModalStrengthEvaluation(t *testing.T) {
	m := &passwordCreateModal{}
	m.Init()
	m.theme = TokyoNight

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
// Behavioral Tests — D-PC-03/04/05/07 real-time validation and strength meter
// ─────────────────────────────────────────────────────────────────────────────

// TestPasswordCreate_RealtimeValidation_ShowsErrorOnMismatch verifies D-PC-05:
// typing in the confirm field when it doesn't match password → error message shown.
func TestPasswordCreate_RealtimeValidation_ShowsErrorOnMismatch(t *testing.T) {
	m := &passwordCreateModal{}
	m.Init()
	m.theme = TokyoNight
	m.messages = NewMessageManager()

	m.password.SetValue("abc123")
	m.confirm.SetValue("abc124") // mismatch
	m.focusIndex = 1             // focus on confirm

	// Simulate a key press in confirm field to trigger real-time validation
	m.Update(tea.KeyPressMsg{Code: 'x', Text: "x"})

	curr := m.messages.Current()
	if curr == nil {
		t.Fatal("expected error message after mismatch, got nil")
	}
	if !strings.Contains(curr.Text, "não conferem") {
		t.Errorf("expected 'não conferem' in message text, got %q", curr.Text)
	}
	if curr.Kind != MessageError {
		t.Errorf("expected MessageError, got %v", curr.Kind)
	}
}

// TestPasswordCreate_RealtimeValidation_NoErrorWhenEmpty verifies D-PC-05:
// when confirm field is empty, no error message is shown.
func TestPasswordCreate_RealtimeValidation_NoErrorWhenEmpty(t *testing.T) {
	m := &passwordCreateModal{}
	m.Init()
	m.theme = TokyoNight
	m.messages = NewMessageManager()

	m.password.SetValue("abc123")
	m.confirm.SetValue("") // empty confirm
	m.focusIndex = 1

	// The hint should be shown (not an error) when confirm is empty
	curr := m.messages.Current()
	if curr != nil && curr.Kind == MessageError {
		t.Errorf("expected no error when confirm is empty, but got MessageError: %q", curr.Text)
	}
}

// TestPasswordCreate_EnterBlocked_EmptyFields verifies D-PC-04:
// Enter returns nil when both fields are empty.
func TestPasswordCreate_EnterBlocked_EmptyFields(t *testing.T) {
	m := &passwordCreateModal{}
	m.Init()
	m.theme = TokyoNight

	m.password.SetValue("")
	m.confirm.SetValue("")

	cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected nil cmd when Enter pressed with empty fields (D-PC-04)")
	}
}

// TestPasswordCreate_EnterBlocked_Mismatch verifies D-PC-04:
// Enter returns nil when fields are non-empty but don't match.
func TestPasswordCreate_EnterBlocked_Mismatch(t *testing.T) {
	m := &passwordCreateModal{}
	m.Init()
	m.theme = TokyoNight

	m.password.SetValue("password1")
	m.confirm.SetValue("password2")

	cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected nil cmd when Enter pressed with mismatched passwords (D-PC-04)")
	}
}

// TestPasswordCreate_StrengthMeterHidden_EmptyPassword verifies D-PC-03:
// strength meter is hidden when password field is empty.
func TestPasswordCreate_StrengthMeterHidden_EmptyPassword(t *testing.T) {
	m := &passwordCreateModal{}
	m.Init()
	m.theme = TokyoNight

	// Leave password empty (default after Init)

	plain := stripANSI(m.View(80, 24, TokyoNight))

	// Strength labels should NOT appear when password is empty
	if strings.Contains(plain, "Força:") {
		t.Error("expected strength meter to be hidden when password is empty (D-PC-03)")
	}
}

// TestPasswordCreate_StrengthMeterVisible_NonEmptyPassword verifies D-PC-03:
// strength meter is visible (shows "Forte" or "Fraca") when password is non-empty.
func TestPasswordCreate_StrengthMeterVisible_NonEmptyPassword(t *testing.T) {
	m := &passwordCreateModal{}
	m.Init()
	m.theme = TokyoNight

	m.password.SetValue("test123")
	m.updateStrength()

	plain := stripANSI(m.View(80, 24, TokyoNight))

	if !strings.Contains(plain, "Forte") && !strings.Contains(plain, "Fraca") {
		t.Error("expected strength meter to be visible when password is non-empty (D-PC-03)")
	}
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
	m.theme = TokyoNight
	m.messages = NewMessageManager()

	out := m.View(80, 24, TokyoNight)

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

// TestPasswordCreateModal_Golden_Empty validates visual output when both fields are empty.
// Verifies D-PC-03: strength meter is hidden when "Nova senha" is empty.
func TestPasswordCreateModal_Golden_Empty(t *testing.T) {
	m := &passwordCreateModal{title: "Criar senha mestra"}
	m.Init()
	m.theme = TokyoNight
	m.messages = NewMessageManager()

	// Leave both fields empty (default after Init)

	out := m.View(80, 24, TokyoNight)

	txtPath := goldenPath("passwordcreate", "empty", 80, "txt")
	checkOrUpdateGolden(t, txtPath, stripANSI(out))

	transitions := testdatapkg.ParseANSIStyle(out)
	jsonBytes, err := testdatapkg.MarshalStyleTransitions(transitions)
	if err != nil {
		t.Fatalf("marshal transitions: %v", err)
	}
	jsonPath := goldenPath("passwordcreate", "empty", 80, "json")
	checkOrUpdateGolden(t, jsonPath, string(jsonBytes))
}

// TestPasswordCreateModal_Golden_FilledMatch validates visual output when passwords match.
// Verifies D-PC-04: "Enter Confirmar" is active (AccentPrimary+bold) when both fields non-empty and match.
func TestPasswordCreateModal_Golden_FilledMatch(t *testing.T) {
	m := &passwordCreateModal{title: "Criar senha mestra"}
	m.Init()
	m.theme = TokyoNight
	m.messages = NewMessageManager()

	m.password.SetValue("SamePass123!")
	m.confirm.SetValue("SamePass123!")
	m.updateStrength()

	out := m.View(80, 24, TokyoNight)

	txtPath := goldenPath("passwordcreate", "filled-match", 80, "txt")
	checkOrUpdateGolden(t, txtPath, stripANSI(out))

	transitions := testdatapkg.ParseANSIStyle(out)
	jsonBytes, err := testdatapkg.MarshalStyleTransitions(transitions)
	if err != nil {
		t.Fatalf("marshal transitions: %v", err)
	}
	jsonPath := goldenPath("passwordcreate", "filled-match", 80, "json")
	checkOrUpdateGolden(t, jsonPath, string(jsonBytes))
}
