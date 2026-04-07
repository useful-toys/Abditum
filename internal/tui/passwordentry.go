package tui

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/tokens"
)

// passwordEntryModal is a modal for entering a password when opening a vault.
// Features:
// - Single masked textinput field (fixed 8 • characters)
// - Attempt counter visible from attempt 2 onward
// - Emits pwdEnteredMsg on Enter, flowCancelledMsg on ESC
type passwordEntryModal struct {
	input       textinput.Model
	title       string
	attempt     int
	maxAttempts int
	width       int
	height      int
	theme       *Theme
	messages    *MessageManager
}

// Init initializes the password entry modal.
func (m *passwordEntryModal) Init() tea.Cmd {
	m.input = textinput.New()
	m.input.Placeholder = ""
	m.input.SetValue("")
	m.input.EchoMode = textinput.EchoPassword
	m.input.EchoCharacter = '•'
	m.input.Focus()
	m.attempt = 1
	m.maxAttempts = 5
	if m.messages == nil {
		m.messages = NewMessageManager()
	}
	m.showHint()
	return nil
}

// HandleWrongPassword increments the attempt counter and shows error message.
func (m *passwordEntryModal) HandleWrongPassword() {
	m.attempt++
	m.input.SetValue("")
	m.messages.Show(MsgError, tokens.SymError+" Senha incorreta", 5, false)
}

// showHint displays the initial hint message.
func (m *passwordEntryModal) showHint() {
	m.messages.Show(MsgHint, tokens.SymHint+" Digite a senha para desbloquear o cofre", 0, false)
}

// Update handles keyboard input.
func (m *passwordEntryModal) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Code {
		case tea.KeyEsc:
			// Cancel flow
			return tea.Batch(
				func() tea.Msg { return popModalMsg{} },
				func() tea.Msg { return flowCancelledMsg{} },
			)
		case tea.KeyEnter:
			password := []byte(m.input.Value())
			if len(password) == 0 {
				m.messages.Show(MsgError, tokens.SymError+" Digite uma senha", 3, false)
				return nil
			}
			// Return the entered password to the flow
			return tea.Batch(
				func() tea.Msg { return popModalMsg{} },
				func() tea.Msg { return pwdEnteredMsg{Password: password} },
			)
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return cmd
}

// View renders the modal.
func (m *passwordEntryModal) View() string {
	// Nil-safe theme fallback
	theme := m.theme
	if theme == nil {
		theme = ThemeTokyoNight
	}

	width := 50
	if m.width > 0 {
		width = m.width
		if width > 70 {
			width = 70
		}
	}

	var content strings.Builder

	// Title
	if m.title != "" {
		content.WriteString(lipgloss.NewStyle().
			Foreground(theme.TextPrimary).
			Bold(true).
			Render(m.title))
		content.WriteString("\n\n")
	}

	// Masked input field (always show 8 dots)
	maskedDisplay := strings.Repeat("•", 8)
	inputStyle := lipgloss.NewStyle().
		Foreground(theme.TextPrimary)
	content.WriteString(inputStyle.Render(maskedDisplay))
	content.WriteString("\n")

	// Attempt counter (visible from attempt 2 onward)
	if m.attempt >= 2 {
		counterStyle := lipgloss.NewStyle().
			Foreground(theme.TextSecondary)
		counter := counterStyle.Render("\nTentativa " + string(rune('0'+m.attempt)) + " de 5")
		content.WriteString(counter)
		content.WriteString("\n")
	}

	// Action hints
	hintsStyle := lipgloss.NewStyle().
		Foreground(theme.TextSecondary)
	hints := hintsStyle.Render("\n\nEnter Confirmar    Esc Cancelar")
	content.WriteString(hints)

	// Render with border
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.SemanticOff).
		Padding(1, 2).
		Width(width).
		Render(content.String())
}

// SetSize stores the modal dimensions.
func (m *passwordEntryModal) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Shortcuts returns keyboard hints for the command bar.
func (m *passwordEntryModal) Shortcuts() []Shortcut {
	return []Shortcut{
		{Key: "Enter", Label: "Confirmar"},
		{Key: "Esc", Label: "Cancelar"},
	}
}

// ApplyTheme applies a theme to the modal.
func (m *passwordEntryModal) ApplyTheme(t *Theme) {
	m.theme = t
}

// flowCancelledMsg signals that a flow was cancelled.
type flowCancelledMsg struct{}
