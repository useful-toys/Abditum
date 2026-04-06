package tui

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/crypto"
	"github.com/useful-toys/abditum/internal/tui/tokens"
)

// passwordCreateModal is a modal for creating a password when creating a vault.
// Features:
// - Two masked textinput fields (password and confirmation)
// - Tab navigation between fields
// - Real-time password strength meter using crypto.EvaluatePasswordStrength
// - Emits pwdCreatedMsg on matching passwords, flowCancelledMsg on ESC
type passwordCreateModal struct {
	password   textinput.Model
	confirm    textinput.Model
	focusIndex int // 0 = password, 1 = confirm
	title      string
	width      int
	height     int
	theme      *Theme
	messages   *MessageManager
	strength   crypto.StrengthLevel
}

// Init initializes the password creation modal.
func (m *passwordCreateModal) Init() tea.Cmd {
	m.password = textinput.New()
	m.password.Placeholder = ""
	m.password.SetValue("")
	m.password.EchoMode = textinput.EchoPassword
	m.password.EchoCharacter = '•'
	m.password.Focus()

	m.confirm = textinput.New()
	m.confirm.Placeholder = ""
	m.confirm.SetValue("")
	m.confirm.EchoMode = textinput.EchoPassword
	m.confirm.EchoCharacter = '•'
	m.confirm.Blur()

	m.focusIndex = 0
	if m.messages == nil {
		m.messages = NewMessageManager()
	}
	m.showHint()
	return nil
}

// showHint displays the initial hint message.
func (m *passwordCreateModal) showHint() {
	m.messages.Show(MsgHint, tokens.SymHint+" Crie uma senha forte para seu cofre", 0, false)
}

// updateStrength evaluates password strength based on current input.
func (m *passwordCreateModal) updateStrength() {
	pwd := []byte(m.password.Value())
	m.strength = crypto.EvaluatePasswordStrength(pwd)
}

// Update handles keyboard input.
func (m *passwordCreateModal) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Code {
		case tea.KeyEsc:
			// Cancel flow
			return tea.Batch(
				func() tea.Msg { return popModalMsg{} },
				func() tea.Msg { return flowCancelledMsg{} },
			)
		case tea.KeyTab:
			// Switch focus between password and confirm
			if m.focusIndex == 0 {
				m.password.Blur()
				m.confirm.Focus()
				m.focusIndex = 1
			} else {
				m.confirm.Blur()
				m.password.Focus()
				m.focusIndex = 0
			}
			return nil
		case tea.KeyEnter:
			passwordVal := m.password.Value()
			confirmVal := m.confirm.Value()

			// Validate inputs
			if len(passwordVal) == 0 {
				m.messages.Show(MsgError, tokens.SymError+" Digite uma senha", 3, false)
				return nil
			}
			if len(confirmVal) == 0 {
				m.messages.Show(MsgError, tokens.SymError+" Confirme a senha", 3, false)
				return nil
			}
			if passwordVal != confirmVal {
				m.messages.Show(MsgError, tokens.SymError+" As senhas nao conferem", 3, false)
				m.confirm.SetValue("")
				m.password.Focus()
				m.focusIndex = 0
				return nil
			}

			// Return the created password to the flow
			return tea.Batch(
				func() tea.Msg { return popModalMsg{} },
				func() tea.Msg { return pwdCreatedMsg{Password: []byte(passwordVal)} },
			)
		}
	}

	var cmd tea.Cmd
	if m.focusIndex == 0 {
		m.password, cmd = m.password.Update(msg)
		m.updateStrength() // Update strength as user types
	} else {
		m.confirm, cmd = m.confirm.Update(msg)
	}
	return cmd
}

// View renders the modal.
func (m *passwordCreateModal) View() string {
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
			Foreground(m.theme.TextPrimary).
			Bold(true).
			Render(m.title))
		content.WriteString("\n\n")
	}

	// Password field label and input
	passwordLabel := "Senha:"
	if m.focusIndex == 0 {
		passwordLabel = lipgloss.NewStyle().
			Foreground(m.theme.AccentPrimary).
			Render(passwordLabel + " ✓")
	}
	content.WriteString(passwordLabel)
	content.WriteString("\n")

	// Masked password display
	maskedPassword := strings.Repeat("•", len(m.password.Value()))
	if len(maskedPassword) == 0 {
		maskedPassword = "______"
	}
	passwordStyle := lipgloss.NewStyle().
		Foreground(m.theme.TextPrimary)
	content.WriteString(passwordStyle.Render(maskedPassword))
	content.WriteString("\n\n")

	// Strength meter
	strengthColor := m.theme.SemanticOff
	strengthLabel := "Fraca"
	if m.strength == crypto.StrengthStrong {
		strengthColor = m.theme.SemanticSuccess
		strengthLabel = "Forte"
	}
	strengthMeter := lipgloss.NewStyle().
		Foreground(strengthColor).
		Render("[" + strengthLabel + "]")
	content.WriteString(strengthMeter)
	content.WriteString("\n\n")

	// Confirm field label and input
	confirmLabel := "Confirmar:"
	if m.focusIndex == 1 {
		confirmLabel = lipgloss.NewStyle().
			Foreground(m.theme.AccentPrimary).
			Render(confirmLabel + " ✓")
	}
	content.WriteString(confirmLabel)
	content.WriteString("\n")

	// Masked confirm display
	maskedConfirm := strings.Repeat("•", len(m.confirm.Value()))
	if len(maskedConfirm) == 0 {
		maskedConfirm = "______"
	}
	confirmStyle := lipgloss.NewStyle().
		Foreground(m.theme.TextPrimary)
	content.WriteString(confirmStyle.Render(maskedConfirm))
	content.WriteString("\n")

	// Action hints
	hintsStyle := lipgloss.NewStyle().
		Foreground(m.theme.TextSecondary)
	hints := hintsStyle.Render("\n\nTab Navegar    Enter Criar    Esc Cancelar")
	content.WriteString(hints)

	// Render with border
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.SemanticOff).
		Padding(1, 2).
		Width(width).
		Render(content.String())
}

// SetSize stores the modal dimensions.
func (m *passwordCreateModal) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Shortcuts returns keyboard hints for the command bar.
func (m *passwordCreateModal) Shortcuts() []Shortcut {
	return []Shortcut{
		{Key: "Tab", Label: "Navegar"},
		{Key: "Enter", Label: "Criar"},
		{Key: "Esc", Label: "Cancelar"},
	}
}

// ApplyTheme applies a theme to the modal.
func (m *passwordCreateModal) ApplyTheme(t *Theme) {
	m.theme = t
}
