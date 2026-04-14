package tui

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/crypto"
)

// passwordCreateModal is a modal for creating a password when creating a vault.
// Features:
// - Two masked textinput fields (Nova senha and Confirmação)
// - Tab navigation between fields
// - Real-time password strength meter using crypto.EvaluatePasswordStrength
// - Emits pwdCreatedMsg on matching passwords, flowCancelledMsg on ESC
//
// Visual layout follows tui-specification-novo.md §PasswordCreate:
//
//	╭── Definir senha mestra ───────────────────╮
//	│                                            │
//	│  Nova senha                                │
//	│  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
//	│                                            │
//	│  Confirmação                               │
//	│  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
//	│                                            │
//	│  Força: [label]   ← only when pwd non-empty│
//	│                                            │
//	╰── Enter Confirmar ──────────── Esc Cancelar ──╯
type passwordCreateModal struct {
	password   textinput.Model
	confirm    textinput.Model
	focusIndex int // 0 = Nova senha, 1 = Confirmação
	title      string
	width      int
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
	m.showPasswordHint()
	return nil
}

// showPasswordHint displays the hint for the Nova senha field.
func (m *passwordCreateModal) showPasswordHint() {
	m.messages.Show(MessageHint, SymBullet+" A senha mestra protege todo o cofre — use 12+ caracteres", 0, false)
}

// showConfirmHint displays the hint for the Confirmação field.
func (m *passwordCreateModal) showConfirmHint() {
	m.messages.Show(MessageHint, SymBullet+" Redigite a senha para confirmar", 0, false)
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
			// Switch focus between Nova senha and Confirmação
			if m.focusIndex == 0 {
				m.password.Blur()
				m.confirm.Focus()
				m.focusIndex = 1
				// D-PC-05: validate on Tab abandonment of Nova senha field
				confirmVal := m.confirm.Value()
				passwordVal := m.password.Value()
				if len(confirmVal) > 0 && confirmVal != passwordVal {
					m.messages.Show(MessageError, SymError+" As senhas não conferem — digite novamente", 5, false)
				} else {
					m.showConfirmHint()
				}
			} else {
				m.confirm.Blur()
				m.password.Focus()
				m.focusIndex = 0
				// D-PC-05: validate on Tab abandonment of Confirmação field
				confirmVal := m.confirm.Value()
				passwordVal := m.password.Value()
				if len(confirmVal) > 0 && confirmVal != passwordVal {
					m.messages.Show(MessageError, SymError+" As senhas não conferem — digite novamente", 5, false)
				} else {
					m.showPasswordHint()
				}
			}
			return nil
		case tea.KeyEnter:
			passwordVal := m.password.Value()
			confirmVal := m.confirm.Value()

			// D-PC-04: Block Enter when either field is empty or passwords diverge
			if len(passwordVal) == 0 {
				m.messages.Show(MessageError, SymError+" Digite uma senha", 3, false)
				return nil
			}
			if len(confirmVal) == 0 {
				m.messages.Show(MessageError, SymError+" Confirme a senha", 3, false)
				return nil
			}
			if passwordVal != confirmVal {
				// D-PC-07: exact error text
				m.messages.Show(MessageError, SymError+" As senhas não conferem — digite novamente", 5, false)
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
		// D-PC-05: Real-time equality validation on each keypress in Confirmação field
		confirmVal := m.confirm.Value()
		passwordVal := m.password.Value()
		if len(confirmVal) > 0 && confirmVal != passwordVal {
			m.messages.Show(MessageError, SymError+" As senhas não conferem — digite novamente", 5, false)
		} else if len(confirmVal) == 0 {
			m.showConfirmHint()
		} else {
			// passwords match — clear error, restore hint
			m.showConfirmHint()
		}
	}
	return cmd
}

// View renders the modal following the spec wireframe.
func (m *passwordCreateModal) View(maxWidth, maxHeight int, theme *Theme) string {
	m.width = maxWidth
	const fixedWidth = 50
	boxW := fixedWidth
	if m.width > 0 && m.width < fixedWidth {
		boxW = m.width
	}

	borderSt := lipgloss.NewStyle().Foreground(lipgloss.Color("#7aa2f7"))
	titleSt := lipgloss.NewStyle().Foreground(lipgloss.Color("#7aa2f7")).Bold(true)
	activeLabelSt := lipgloss.NewStyle().Foreground(lipgloss.Color("#7aa2f7")).Bold(true)
	inactiveLabelSt := lipgloss.NewStyle().Foreground(lipgloss.Color("#565f89"))
	secondarySt := lipgloss.NewStyle().Foreground(lipgloss.Color("#565f89"))
	fieldBgSt := lipgloss.NewStyle().Background(lipgloss.Color("#1e1f2e"))

	dashes := func(n int) string {
		if n < 0 {
			n = 0
		}
		return borderSt.Render(strings.Repeat("─", n))
	}

	innerW := boxW - 2

	// ── Top border ───────────────────────────────────────────────────────────
	titleText := m.title
	if titleText == "" {
		titleText = "Definir senha mestra"
	}
	titleW := lipgloss.Width(titleText)
	const leftAnchorW = 4
	const rightAnchorW = 4
	fillW := boxW - leftAnchorW - titleW - rightAnchorW
	if fillW < 1 {
		fillW = 1
	}
	topLine := borderSt.Render("╭──") + " " +
		titleSt.Render(titleText) + " " +
		borderSt.Render(strings.Repeat("─", fillW)+"──╮")

	// ── Body helpers ─────────────────────────────────────────────────────────
	contentW := innerW - 4 // 2 spaces padding each side
	if contentW < 1 {
		contentW = 1
	}

	pad := func(s string) string {
		plain := lipgloss.Width(s)
		trailing := contentW - plain
		if trailing < 0 {
			trailing = 0
		}
		return borderSt.Render("│") + "  " + s + strings.Repeat(" ", trailing) + "  " + borderSt.Render("│")
	}
	emptyLine := pad("")

	fieldLine := func() string {
		return borderSt.Render("│") + "  " + fieldBgSt.Width(contentW).Render("") + "  " + borderSt.Render("│")
	}

	// Nova senha label (active = accent.primary+bold, inactive = text.secondary)
	var pwdLabelSt lipgloss.Style
	if m.focusIndex == 0 {
		pwdLabelSt = activeLabelSt
	} else {
		pwdLabelSt = inactiveLabelSt
	}
	pwdLabelLine := pad(pwdLabelSt.Render("Nova senha"))

	// Confirmação label
	var confirmLabelSt lipgloss.Style
	if m.focusIndex == 1 {
		confirmLabelSt = activeLabelSt
	} else {
		confirmLabelSt = inactiveLabelSt
	}
	confirmLabelLine := pad(confirmLabelSt.Render("Confirmação"))

	var lines []string
	lines = append(lines, topLine, emptyLine)
	lines = append(lines, pwdLabelLine)
	lines = append(lines, fieldLine())
	lines = append(lines, emptyLine)
	lines = append(lines, confirmLabelLine)
	lines = append(lines, fieldLine())
	lines = append(lines, emptyLine)

	// D-PC-03: Strength meter only visible when Nova senha is non-empty
	if len(m.password.Value()) > 0 {
		var strengthLabel string
		var strengthColor string
		if m.strength == crypto.StrengthStrong {
			strengthColor = "#9ece6a"
			strengthLabel = "Forte"
		} else {
			strengthColor = "#e0af68"
			strengthLabel = "Fraca"
		}
		strengthSt := lipgloss.NewStyle().Foreground(lipgloss.Color(strengthColor))
		strengthLine := pad(secondarySt.Render("Força: ") + strengthSt.Render(strengthLabel))
		lines = append(lines, strengthLine, emptyLine)
	}

	// ── Action bar (bottom border) ────────────────────────────────────────────
	// D-PC-04: Enter Confirmar disabled when either empty or mismatch
	passwordVal := m.password.Value()
	confirmVal := m.confirm.Value()
	isActive := len(passwordVal) > 0 && len(confirmVal) > 0 && passwordVal == confirmVal

	var enterSt lipgloss.Style
	if isActive {
		enterSt = lipgloss.NewStyle().Foreground(lipgloss.Color("#7aa2f7")).Bold(true)
	} else {
		enterSt = lipgloss.NewStyle().Foreground(lipgloss.Color("#3b4261"))
	}

	enterToken := enterSt.Render("Enter") + secondarySt.Render(" Confirmar")
	enterPlain := "Enter Confirmar"
	leftPart := borderSt.Render("╰") + dashes(2) + " " + enterToken + " "
	leftPlainW := 4 + len([]rune(enterPlain)) + 1

	escToken := borderSt.Render("Esc") + secondarySt.Render(" Cancelar")
	rightPlain := " Esc Cancelar ──╯"
	rightPart := " " + escToken + " " + dashes(2) + borderSt.Render("╯")

	rightPlainW := len([]rune(rightPlain))
	actionFillW := boxW - leftPlainW - rightPlainW
	if actionFillW < 1 {
		actionFillW = 1
	}
	actionBar := leftPart + dashes(actionFillW) + rightPart

	lines = append(lines, actionBar)
	return strings.Join(lines, "\n")
}

// Shortcuts returns keyboard hints for the command bar.
func (m *passwordCreateModal) Shortcuts() []Shortcut {
	return []Shortcut{
		{Key: "Tab", Label: "Navegar"},
		{Key: "Enter", Label: "Criar"},
		{Key: "Esc", Label: "Cancelar"},
	}
}
