package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// welcomeModel renders the welcome background (ASCII art logo + action hints).
// It is active during workAreaWelcome and has no sub-states.
// Open/create vault flows are orchestrated via the modal stack, not this model.
type welcomeModel struct {
	actions *ActionManager
	theme   *Theme
	width   int
	height  int
}

// ApplyTheme applies the given theme to the welcomeModel.
func (m *welcomeModel) ApplyTheme(t *Theme) {
	m.theme = t
}

// Compile-time assertion: welcomeModel satisfies childModel.
var _ childModel = &welcomeModel{}

// newWelcomeModel creates a new welcome screen model.
func newWelcomeModel(actions *ActionManager, theme *Theme) *welcomeModel {
	return &welcomeModel{actions: actions, theme: theme}
}

// Update processes messages for the welcome screen.
// Phase 5.1: welcomeModel is display-only. No input handling until Phase 6.
func (m *welcomeModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// View renders the ASCII art logo centered on screen.
// Per spec (tui-specification-novo.md § Boas-vindas), the logo and version
// are centered horizontally and vertically via lipgloss.Place().
// Logo width is hardcoded to 43 columns matching the ASCII art width.
func (m *welcomeModel) View() string {
	// 43 = width of AsciiArt (const in ascii.go) — each line is exactly 43 characters
	logoBlock := lipgloss.NewStyle().Width(43).Render(RenderLogo(m.theme))

	content := lipgloss.JoinVertical(lipgloss.Center, logoBlock)

	if m.width == 0 || m.height == 0 {
		// Terminal dimensions not yet set (edge case during init)
		// Return uncentered content as fallback
		return content
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

// SetSize stores the allocated terminal dimensions for layout.
func (m *welcomeModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}
