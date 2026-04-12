package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/tokens"
)

// passwordEntryModal is a modal for entering a password when opening a vault.
// Features:
// - Single masked textinput field (fixed 8 • characters displayed)
// - Attempt counter visible from attempt 2 onward
// - Emits pwdEnteredMsg on Enter, flowCancelledMsg on ESC
//
// Visual layout follows tui-specification-novo.md §PasswordEntry:
//
//	╭── Senha mestra ────────────────────────────╮
//	│                                            │
//	│  Senha                                     │
//	│  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
//	│                                            │
//	╰── Enter Confirmar ──────────── Esc Cancelar ──╯
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

// View renders the modal following the spec wireframe:
//
//	╭── Senha mestra ────────────────────────────╮
//	│                                            │
//	│  Senha                                     │
//	│  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
//	│                                            │
//	│  Tentativa N de 5                          │  ← attempt ≥ 2 only
//	╰── Enter Confirmar ──────────── Esc Cancelar ──╯
func (m *passwordEntryModal) View() string {
	if m.width == 0 || m.height == 0 {
		panic(fmt.Sprintf("passwordEntryModal.View() called without SetSize: width=%d height=%d", m.width, m.height))
	}
	const fixedWidth = 50
	boxW := fixedWidth
	if m.width > 0 && m.width < fixedWidth {
		boxW = m.width
	}

	borderSt := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorBorderFocused))
	titleSt := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorBorderFocused)).Bold(true)
	labelSt := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorAccentPrimary)).Bold(true)
	secondarySt := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTextSecondary))
	fieldBgSt := lipgloss.NewStyle().Background(lipgloss.Color(ColorSurfaceInput))

	dashes := func(n int) string {
		if n < 0 {
			n = 0
		}
		return borderSt.Render(strings.Repeat("─", n))
	}

	innerW := boxW - 2 // columns between │ chars

	// ── Top border ───────────────────────────────────────────────────────────
	// ╭── {title} ────...────╮
	titleText := m.title
	if titleText == "" {
		titleText = "Senha mestra"
	}
	titleW := lipgloss.Width(titleText)
	const leftAnchorW = 4  // "╭── "
	const rightAnchorW = 4 // " ──╮"
	fillW := boxW - leftAnchorW - titleW - rightAnchorW
	if fillW < 1 {
		fillW = 1
	}
	topLine := borderSt.Render("╭──") + " " +
		titleSt.Render(titleText) + " " +
		borderSt.Render(strings.Repeat("─", fillW)+"──╮")

	// ── Body lines ───────────────────────────────────────────────────────────
	pad := func(s string) string {
		// Render a content line: │  {s padded to innerW-2}  │
		// innerW-4 = usable content width (2 spaces padding each side)
		contentW := innerW - 4
		if contentW < 0 {
			contentW = 0
		}
		plain := lipgloss.Width(s)
		trailing := contentW - plain
		if trailing < 0 {
			trailing = 0
		}
		return borderSt.Render("│") + "  " + s + strings.Repeat(" ", trailing) + "  " + borderSt.Render("│")
	}
	emptyLine := pad("")

	// Label "Senha" — always accent.primary + bold (single-field dialog)
	labelLine := pad(labelSt.Render("Senha"))

	// Input field: surface.input background, fixed 8 • mask
	maskedPassword := strings.Repeat("•", 8)
	fieldW := innerW - 4 // matches content width
	if fieldW < 1 {
		fieldW = 1
	}
	fieldContent := fieldBgSt.Width(fieldW).Render(maskedPassword)
	fieldLine := borderSt.Render("│") + "  " + fieldContent + "  " + borderSt.Render("│")

	// Attempt counter (attempt ≥ 2)
	var lines []string
	lines = append(lines, topLine, emptyLine, labelLine, fieldLine, emptyLine)
	if m.attempt >= 2 {
		counterText := secondarySt.Render(fmt.Sprintf("Tentativa %d de %d", m.attempt, m.maxAttempts))
		lines = append(lines, pad(counterText))
	}

	// ── Action bar (bottom border) ────────────────────────────────────────────
	// ╰── Enter Confirmar ──────────── Esc Cancelar ──╯
	isEmpty := len(m.input.Value()) == 0
	var enterSt lipgloss.Style
	if isEmpty {
		enterSt = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTextDisabled))
	} else {
		enterSt = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorAccentPrimary)).Bold(true)
	}

	// Left: "╰── Enter Confirmar "
	enterToken := enterSt.Render("Enter") + secondarySt.Render(" Confirmar")
	enterPlain := "Enter Confirmar"
	leftPart := borderSt.Render("╰") + dashes(2) + " " + enterToken + " "
	leftPlainW := 4 + len([]rune(enterPlain)) + 1 // "╰── " + token + " "

	// Right: " Esc Cancelar ──╯"
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

// SetAvailableSize stores the maximum available modal dimensions.
func (m *passwordEntryModal) SetAvailableSize(maxWidth, maxHeight int) {
	m.width = maxWidth
	m.height = maxHeight
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
