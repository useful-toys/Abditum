package modal

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/crypto"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// PasswordField is a secure text input field that manages password buffers
// using []byte and supports memory wiping.
//
// The field displays a label and an input area with bullet masking for each
// character. It only supports printable rune input and backspace; all other
// keys are ignored.
//
// The internal buffer is never converted to a string — callers receive copies
// via Value() and are responsible for wiping sensitive data.
type PasswordField struct {
	label string // Label displayed above the field (e.g., "Senha", "Nova senha")
	value []byte // Password buffer — never converted to string internally
}

// NewPasswordField creates a new PasswordField with the given label.
// The field starts empty and ready to accept input.
func NewPasswordField(label string) *PasswordField {
	return &PasswordField{
		label: label,
		value: []byte{},
	}
}

// Value returns a copy of the password buffer.
// The caller is responsible for wiping the returned slice when done.
func (f *PasswordField) Value() []byte {
	if f.value == nil {
		return []byte{}
	}
	// Return a copy, not a reference
	result := make([]byte, len(f.value))
	// Use copy builtin to copy bytes
	n := 0
	for i := range f.value {
		result[i] = f.value[i]
		n = i + 1
	}
	return result[:n]
}

// Len returns the number of bytes in the password buffer.
func (f *PasswordField) Len() int {
	return len(f.value)
}

// Clear sets the password buffer to nil (empty).
func (f *PasswordField) Clear() {
	f.value = nil
}

// Wipe securely clears the password buffer by overwriting it with zeros,
// then clearing it. This should be called before discarding a PasswordField
// to ensure sensitive data doesn't remain in memory.
func (f *PasswordField) Wipe() {
	if f.value != nil {
		crypto.Wipe(f.value)
	}
	f.Clear()
}

// HandleKey processes a keyboard message.
//
// Behavior:
// - Printable characters (Key.Text non-empty): appends the text to the buffer and returns true (consumed)
// - tea.KeyBackspace: removes the last byte from the buffer and returns true (consumed)
// - Any other key: returns false (not consumed)
func (f *PasswordField) HandleKey(msg tea.KeyPressMsg) bool {
	key := msg.Key()

	// Check if it's a printable character
	if key.Text != "" {
		// Append the text (UTF-8 encoded) to the buffer
		f.value = append(f.value, []byte(key.Text)...)
		return true
	}

	// Check for backspace
	if key.Code == tea.KeyBackspace {
		// Remove the last byte
		if len(f.value) > 0 {
			f.value = f.value[:len(f.value)-1]
		}
		return true
	}

	// Other keys (arrows, delete, home, end, etc.) are not consumed
	return false
}

// Render renders the password field as two lines:
// Line 1: label with appropriate styling
// Line 2: bullet-masked input area
//
// Parameters:
// - innerWidth: the width available for rendering (typically modalWidth - 2 - 2*padding)
// - focused: whether the field is currently focused
// - theme: the theme to apply
//
// Returns a string with two lines separated by \n. Each line uses ANSI styling
// from the theme. The input area shows one bullet (•) per byte in the password,
// padded with spaces to reach innerWidth.
func (f *PasswordField) Render(innerWidth int, focused bool, theme *design.Theme) string {
	// Line 1: Label with appropriate styling
	labelStyle := lipgloss.NewStyle()
	if focused {
		labelStyle = labelStyle.
			Foreground(lipgloss.Color(theme.Accent.Primary)).
			Bold(true)
	} else {
		labelStyle = labelStyle.
			Foreground(lipgloss.Color(theme.Text.Secondary))
	}
	labelLine := labelStyle.Render(f.label)

	// Line 2: Input area with bullet masking and padding
	// Create bullets: one per byte in the password
	numBullets := len(f.value)
	bullets := strings.Repeat("•", numBullets)

	// Calculate padding to fill innerWidth
	padding := innerWidth - numBullets
	if padding < 0 {
		// If content exceeds innerWidth, truncate the bullets display
		// (keep the full buffer internally, but truncate display)
		bullets = bullets[:innerWidth]
		padding = 0
	}

	// Apply background color to the entire input line
	inputStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(theme.Surface.Input))
	inputLine := inputStyle.Render(bullets + strings.Repeat(" ", padding))

	// Return both lines separated by \n
	return labelLine + "\n" + inputLine
}
