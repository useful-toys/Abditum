package modal_test

import (
	"bytes"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/tui/testdata"
)

// makeKeyMsg creates a KeyPressMsg from a Key.
func makeKeyMsg(key tea.Key) tea.KeyPressMsg {
	return tea.KeyPressMsg(key)
}

// TestPasswordField_InitiallyEmpty tests that a new PasswordField starts empty.
func TestPasswordField_InitiallyEmpty(t *testing.T) {
	field := modal.NewPasswordField("Senha")
	if field.Len() != 0 {
		t.Errorf("NewPasswordField().Len() = %d, want 0", field.Len())
	}
	value := field.Value()
	if len(value) != 0 {
		t.Errorf("NewPasswordField().Value() returned %d bytes, want 0", len(value))
	}
}

// TestPasswordField_HandleKey_Rune tests appending runes to the password field.
func TestPasswordField_HandleKey_Rune(t *testing.T) {
	field := modal.NewPasswordField("Senha")

	// Send a rune key
	msg := makeKeyMsg(tea.Key{Text: "a", Code: 'a'})
	consumed := field.HandleKey(msg)
	if !consumed {
		t.Error("HandleKey(rune 'a') should return true (consumed)")
	}
	if field.Len() != 1 {
		t.Errorf("after appending 'a', Len() = %d, want 1", field.Len())
	}

	// Append another rune
	msg = makeKeyMsg(tea.Key{Text: "b", Code: 'b'})
	consumed = field.HandleKey(msg)
	if !consumed {
		t.Error("HandleKey(rune 'b') should return true (consumed)")
	}
	if field.Len() != 2 {
		t.Errorf("after appending 'b', Len() = %d, want 2", field.Len())
	}

	// Verify the value
	value := field.Value()
	if !bytes.Equal(value, []byte("ab")) {
		t.Errorf("Value() = %q, want \"ab\"", string(value))
	}
}

// TestPasswordField_HandleKey_Backspace tests removing characters via Backspace.
func TestPasswordField_HandleKey_Backspace(t *testing.T) {
	field := modal.NewPasswordField("Senha")

	// Append some characters
	for _, r := range []rune("password") {
		msg := makeKeyMsg(tea.Key{Text: string(r), Code: r})
		field.HandleKey(msg)
	}
	if field.Len() != 8 {
		t.Errorf("after appending 8 chars, Len() = %d, want 8", field.Len())
	}

	// Backspace once
	msg := makeKeyMsg(tea.Key{Code: tea.KeyBackspace})
	consumed := field.HandleKey(msg)
	if !consumed {
		t.Error("HandleKey(Backspace) should return true (consumed)")
	}
	if field.Len() != 7 {
		t.Errorf("after Backspace, Len() = %d, want 7", field.Len())
	}

	// Verify the value
	value := field.Value()
	if !bytes.Equal(value, []byte("passwor")) {
		t.Errorf("Value() = %q, want \"passwor\"", string(value))
	}
}

// TestPasswordField_HandleKey_BackspaceEmpty tests Backspace on empty field.
func TestPasswordField_HandleKey_BackspaceEmpty(t *testing.T) {
	field := modal.NewPasswordField("Senha")

	msg := makeKeyMsg(tea.Key{Code: tea.KeyBackspace})
	consumed := field.HandleKey(msg)
	if !consumed {
		t.Error("HandleKey(Backspace) on empty field should return true (consumed)")
	}
	if field.Len() != 0 {
		t.Errorf("Backspace on empty field should keep Len() = 0, got %d", field.Len())
	}
}

// TestPasswordField_HandleKey_ArrowNotConsumed tests that arrow keys are not consumed.
func TestPasswordField_HandleKey_ArrowNotConsumed(t *testing.T) {
	field := modal.NewPasswordField("Senha")

	// Arrow Left should not be consumed
	msg := makeKeyMsg(tea.Key{Code: tea.KeyLeft})
	consumed := field.HandleKey(msg)
	if consumed {
		t.Error("HandleKey(Left arrow) should return false (not consumed)")
	}

	// Arrow Right should not be consumed
	msg = makeKeyMsg(tea.Key{Code: tea.KeyRight})
	consumed = field.HandleKey(msg)
	if consumed {
		t.Error("HandleKey(Right arrow) should return false (not consumed)")
	}

	// Delete should not be consumed
	msg = makeKeyMsg(tea.Key{Code: tea.KeyDelete})
	consumed = field.HandleKey(msg)
	if consumed {
		t.Error("HandleKey(Delete) should return false (not consumed)")
	}

	// Home should not be consumed
	msg = makeKeyMsg(tea.Key{Code: tea.KeyHome})
	consumed = field.HandleKey(msg)
	if consumed {
		t.Error("HandleKey(Home) should return false (not consumed)")
	}

	// End should not be consumed
	msg = makeKeyMsg(tea.Key{Code: tea.KeyEnd})
	consumed = field.HandleKey(msg)
	if consumed {
		t.Error("HandleKey(End) should return false (not consumed)")
	}
}

// TestPasswordField_Clear tests the Clear method.
func TestPasswordField_Clear(t *testing.T) {
	field := modal.NewPasswordField("Senha")

	// Add some data
	for _, r := range []rune("secret") {
		msg := makeKeyMsg(tea.Key{Text: string(r), Code: r})
		field.HandleKey(msg)
	}
	if field.Len() != 6 {
		t.Errorf("after adding 6 chars, Len() = %d, want 6", field.Len())
	}

	// Clear
	field.Clear()
	if field.Len() != 0 {
		t.Errorf("after Clear(), Len() = %d, want 0", field.Len())
	}

	value := field.Value()
	if len(value) != 0 {
		t.Errorf("after Clear(), Value() returned %d bytes, want 0", len(value))
	}
}

// TestPasswordField_Value_ReturnsCopy tests that Value() returns a copy, not a reference.
func TestPasswordField_Value_ReturnsCopy(t *testing.T) {
	field := modal.NewPasswordField("Senha")

	// Add some data
	for _, r := range []rune("secret") {
		msg := makeKeyMsg(tea.Key{Text: string(r), Code: r})
		field.HandleKey(msg)
	}

	// Get the value
	value1 := field.Value()
	if !bytes.Equal(value1, []byte("secret")) {
		t.Errorf("Value() = %q, want \"secret\"", string(value1))
	}

	// Modify the returned copy
	if len(value1) > 0 {
		value1[0] = 'X'
	}

	// Get the value again — it should still be "secret", not "Xecret"
	value2 := field.Value()
	if !bytes.Equal(value2, []byte("secret")) {
		t.Errorf("after modifying returned copy, Value() = %q, want \"secret\"", string(value2))
	}
}

// TestPasswordField_Wipe tests the Wipe method (clears and zeros memory).
func TestPasswordField_Wipe(t *testing.T) {
	field := modal.NewPasswordField("Senha")

	// Add some data
	for _, r := range []rune("topsecret") {
		msg := makeKeyMsg(tea.Key{Text: string(r), Code: r})
		field.HandleKey(msg)
	}
	if field.Len() != 9 {
		t.Errorf("after adding 9 chars, Len() = %d, want 9", field.Len())
	}

	// Wipe
	field.Wipe()
	if field.Len() != 0 {
		t.Errorf("after Wipe(), Len() = %d, want 0", field.Len())
	}

	value := field.Value()
	if len(value) != 0 {
		t.Errorf("after Wipe(), Value() returned %d bytes, want 0", len(value))
	}
}

// TestPasswordField_EmptyFocused tests rendering of an empty, focused password field.
func TestPasswordField_EmptyFocused(t *testing.T) {
	testdata.TestRenderManaged(t, "password_field", "empty_focused", []string{"44x2"},
		func(w, h int, theme *design.Theme) string {
			field := modal.NewPasswordField("Senha")
			return field.Render(w-2, true, theme)
		})
}

// TestPasswordField_EmptyBlurred tests rendering of an empty, blurred password field.
func TestPasswordField_EmptyBlurred(t *testing.T) {
	testdata.TestRenderManaged(t, "password_field", "empty_blurred", []string{"44x2"},
		func(w, h int, theme *design.Theme) string {
			field := modal.NewPasswordField("Senha")
			return field.Render(w-2, false, theme)
		})
}

// TestPasswordField_ContentFocused tests rendering of a focused field with content (8 chars).
func TestPasswordField_ContentFocused(t *testing.T) {
	testdata.TestRenderManaged(t, "password_field", "content_focused", []string{"44x2"},
		func(w, h int, theme *design.Theme) string {
			field := modal.NewPasswordField("Senha")
			// Add 8 characters
			for _, r := range []rune("password") {
				msg := makeKeyMsg(tea.Key{Text: string(r), Code: r})
				field.HandleKey(msg)
			}
			return field.Render(w-2, true, theme)
		})
}

// TestPasswordField_ContentBlurred tests rendering of a blurred field with content (8 chars).
func TestPasswordField_ContentBlurred(t *testing.T) {
	testdata.TestRenderManaged(t, "password_field", "content_blurred", []string{"44x2"},
		func(w, h int, theme *design.Theme) string {
			field := modal.NewPasswordField("Senha")
			// Add 8 characters
			for _, r := range []rune("password") {
				msg := makeKeyMsg(tea.Key{Text: string(r), Code: r})
				field.HandleKey(msg)
			}
			return field.Render(w-2, false, theme)
		})
}
