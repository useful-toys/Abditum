package modal_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/tui/testdata"
)

// newCreateModal creates a PasswordCreateModal for testing with stub MessageController.
func newCreateModal() *modal.PasswordCreateModal {
	mc := &stubMessageController{}
	return modal.NewPasswordCreateModal(mc,
		func(_ []byte) tea.Cmd { return nil },
		func() tea.Cmd { return nil },
	)
}

// TestPasswordCreateModal_Initial tests the initial render with both fields empty.
// Expected output: 50x9 grid, no strength meter.
func TestPasswordCreateModal_Initial(t *testing.T) {
	m := newCreateModal()
	testdata.TestRenderManaged(t, "password_create", "initial", []string{"50x9"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

// TestPasswordCreateModal_WithMeter tests render with Nova senha having 8 characters.
// Expected output: 50x11 grid (includes strength meter).
func TestPasswordCreateModal_WithMeter(t *testing.T) {
	m := newCreateModal()
	for _, r := range "12345678" {
		msg := makeKeyMsg(tea.Key{Text: string(r), Code: r})
		m.Update(msg)
	}
	testdata.TestRenderManaged(t, "password_create", "with_meter", []string{"50x11"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

// TestPasswordCreateModal_Confirmed tests render with both fields filled with matching strong password.
// Expected output: 50x11 grid (includes strength meter).
func TestPasswordCreateModal_Confirmed(t *testing.T) {
	m := newCreateModal()
	// Type strong password in Nova senha
	for _, r := range "MyP@ssw0rd123" {
		msg := makeKeyMsg(tea.Key{Text: string(r), Code: r})
		m.Update(msg)
	}
	// Switch to Confirmação
	m.Update(makeKeyMsg(tea.Key{Code: tea.KeyTab}))
	// Type same password
	for _, r := range "MyP@ssw0rd123" {
		msg := makeKeyMsg(tea.Key{Text: string(r), Code: r})
		m.Update(msg)
	}
	testdata.TestRenderManaged(t, "password_create", "confirmed", []string{"50x11"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

// TestPasswordCreateModal_Mismatch tests render with both fields filled but different.
// Expected output: 50x11 grid (includes strength meter).
func TestPasswordCreateModal_Mismatch(t *testing.T) {
	m := newCreateModal()
	// Type password in Nova senha
	for _, r := range "MyP@ssw0rd123" {
		msg := makeKeyMsg(tea.Key{Text: string(r), Code: r})
		m.Update(msg)
	}
	// Switch to Confirmação
	m.Update(makeKeyMsg(tea.Key{Code: tea.KeyTab}))
	// Type different password
	for _, r := range "Different1234" {
		msg := makeKeyMsg(tea.Key{Text: string(r), Code: r})
		m.Update(msg)
	}
	testdata.TestRenderManaged(t, "password_create", "mismatch", []string{"50x11"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

// TestPasswordCreateModal_TabSwitchesFocus tests that Tab key toggles focus between fields.
func TestPasswordCreateModal_TabSwitchesFocus(t *testing.T) {
	m := newCreateModal()
	if !m.FocusedOnNew() {
		t.Errorf("Initial focus: expected fieldNew, got fieldConfirm")
	}
	// Press Tab to switch to fieldConfirm
	m.Update(makeKeyMsg(tea.Key{Code: tea.KeyTab}))
	if m.FocusedOnNew() {
		t.Errorf("After first Tab: expected fieldConfirm, got fieldNew")
	}
	// Press Tab again to switch back to fieldNew
	m.Update(makeKeyMsg(tea.Key{Code: tea.KeyTab}))
	if !m.FocusedOnNew() {
		t.Errorf("After second Tab: expected fieldNew, got fieldConfirm")
	}
}

// TestPasswordCreateModal_Cursor_FocusedNew tests cursor position when focused on Nova senha.
// Expected: Y = topY + 1 + 2 = 4, X = leftX + 1 + 2 + Len()
func TestPasswordCreateModal_Cursor_FocusedNew(t *testing.T) {
	m := newCreateModal()
	for _, r := range "12345" {
		msg := makeKeyMsg(tea.Key{Text: string(r), Code: r})
		m.Update(msg)
	}
	c := m.Cursor(1, 0)
	if c == nil {
		t.Fatal("Cursor: expected non-nil cursor")
	}
	// Y = 1 + 1 + 2 = 4, X = 0 + 1 + 2 + 5 = 8
	if c.Position.Y != 4 {
		t.Errorf("Cursor.Y: expected 4, got %d", c.Position.Y)
	}
	if c.Position.X != 8 {
		t.Errorf("Cursor.X: expected 8, got %d", c.Position.X)
	}
}

// TestPasswordCreateModal_Cursor_FocusedConfirm tests cursor position when focused on Confirmação.
// Expected: Y = topY + 1 + 6 = 8, X = leftX + 1 + 2 + Len()
func TestPasswordCreateModal_Cursor_FocusedConfirm(t *testing.T) {
	m := newCreateModal()
	// Type in Nova senha
	for _, r := range "12345" {
		msg := makeKeyMsg(tea.Key{Text: string(r), Code: r})
		m.Update(msg)
	}
	// Switch to Confirmação
	m.Update(makeKeyMsg(tea.Key{Code: tea.KeyTab}))
	// Type in Confirmação
	for _, r := range "abcde" {
		msg := makeKeyMsg(tea.Key{Text: string(r), Code: r})
		m.Update(msg)
	}
	c := m.Cursor(1, 0)
	if c == nil {
		t.Fatal("Cursor: expected non-nil cursor")
	}
	// Y = 1 + 1 + 6 = 8, X = 0 + 1 + 2 + 5 = 8
	if c.Position.Y != 8 {
		t.Errorf("Cursor.Y: expected 8, got %d", c.Position.Y)
	}
	if c.Position.X != 8 {
		t.Errorf("Cursor.X: expected 8, got %d", c.Position.X)
	}
}
