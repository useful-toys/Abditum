package modal_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/tui/testdata"
)

// stubMessageController implementa tui.MessageController com métodos no-op.
// Reutilizado em todos os testes de modal de senha.
type stubMessageController struct{}

func (s *stubMessageController) SetBusy(string)      {}
func (s *stubMessageController) SetSuccess(string)   {}
func (s *stubMessageController) SetError(string)     {}
func (s *stubMessageController) SetWarning(string)   {}
func (s *stubMessageController) SetInfo(string)      {}
func (s *stubMessageController) SetHintField(string) {}
func (s *stubMessageController) SetHintUsage(string) {}
func (s *stubMessageController) Clear()              {}

var _ tui.MessageController = (*stubMessageController)(nil)

func TestPasswordEntryModal_Empty(t *testing.T) {
	mc := &stubMessageController{}
	m := modal.NewPasswordEntryModal(mc,
		func(_ []byte) tea.Cmd { return nil },
		func() tea.Cmd { return nil },
	)
	testdata.TestRenderManaged(t, "password_entry", "empty", []string{"50x6"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestPasswordEntryModal_WithContent(t *testing.T) {
	mc := &stubMessageController{}
	m := modal.NewPasswordEntryModal(mc,
		func(_ []byte) tea.Cmd { return nil },
		func() tea.Cmd { return nil },
	)
	for _, r := range "12345678" {
		msg := makeKeyMsg(tea.Key{Text: string(r), Code: r})
		m.HandleKey(msg)
	}
	testdata.TestRenderManaged(t, "password_entry", "with_content", []string{"50x6"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestPasswordEntryModal_NotifyWrongPassword(t *testing.T) {
	mc := &stubMessageController{}
	m := modal.NewPasswordEntryModal(mc,
		func(_ []byte) tea.Cmd { return nil },
		func() tea.Cmd { return nil },
	)
	for _, r := range "12345678" {
		msg := makeKeyMsg(tea.Key{Text: string(r), Code: r})
		m.HandleKey(msg)
	}
	m.NotifyWrongPassword()
	if m.Len() != 0 {
		t.Errorf("NotifyWrongPassword: expected field cleared, got Len()=%d", m.Len())
	}
}

func TestPasswordEntryModal_Cursor_Empty(t *testing.T) {
	mc := &stubMessageController{}
	m := modal.NewPasswordEntryModal(mc,
		func(_ []byte) tea.Cmd { return nil },
		func() tea.Cmd { return nil },
	)
	c := m.Cursor(1, 0)
	if c == nil {
		t.Fatal("Cursor: expected non-nil cursor")
	}
	// topY=1, leftX=0, borda=1, linhaDoFieldNoBody=2, paddingH=2, len=0
	// Y = 1 + 1 + 2 = 4
	// X = 0 + 1 + 2 + 0 = 3
	if c.Position.Y != 4 {
		t.Errorf("Cursor.Y: expected 4, got %d", c.Position.Y)
	}
	if c.Position.X != 3 {
		t.Errorf("Cursor.X: expected 3, got %d", c.Position.X)
	}
}

func TestPasswordEntryModal_Cursor_WithContent(t *testing.T) {
	mc := &stubMessageController{}
	m := modal.NewPasswordEntryModal(mc,
		func(_ []byte) tea.Cmd { return nil },
		func() tea.Cmd { return nil },
	)
	for _, r := range "12345" {
		msg := makeKeyMsg(tea.Key{Text: string(r), Code: r})
		m.HandleKey(msg)
	}
	c := m.Cursor(1, 0)
	if c == nil {
		t.Fatal("Cursor: expected non-nil cursor")
	}
	// Y = 1 + 1 + 2 = 4
	// X = 0 + 1 + 2 + 5 = 8
	if c.Position.Y != 4 {
		t.Errorf("Cursor.Y: expected 4, got %d", c.Position.Y)
	}
	if c.Position.X != 8 {
		t.Errorf("Cursor.X: expected 8, got %d", c.Position.X)
	}
}
