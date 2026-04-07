package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// welcomeModel renders the welcome background (ASCII art logo + version + action hints).
// It is active during workAreaWelcome and has no sub-states.
// Open/create vault flows are orchestrated via the modal stack, not this model.
type welcomeModel struct {
	actions *ActionManager
	theme   *Theme
	version string // Application version to display below logo
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
func newWelcomeModel(actions *ActionManager, theme *Theme, version string) *welcomeModel {
	return &welcomeModel{actions: actions, theme: theme, version: version}
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
// Version is displayed below the logo in text.secondary color.
func (m *welcomeModel) View() string {
	// 43 = width of AsciiArt (const in ascii.go) — each line is exactly 43 characters.
	// Background must match SurfaceBase to prevent default-terminal bg bleed from
	// the Width() padding fill characters leaking through the work area style.
	logoBlock := lipgloss.NewStyle().Width(43).Background(m.theme.SurfaceBase).Render(RenderLogo(m.theme))

	// Format version with semantic.secondary color (from theme)
	// Per spec: version token = text.secondary
	versionStyle := lipgloss.NewStyle().Foreground(m.theme.TextSecondary)
	versionLine := versionStyle.Render(m.version)

	content := lipgloss.JoinVertical(lipgloss.Center, logoBlock, "", versionLine)

	if m.width == 0 || m.height == 0 {
		// Terminal dimensions not yet set (edge case during init)
		// Return uncentered content as fallback
		return content
	}
	// WithWhitespaceStyle ensures the padding space around the centered content
	// uses SurfaceBase background instead of the terminal default, preventing
	// a spurious #000000 background on the first character of the welcome area.
	bgStyle := lipgloss.NewStyle().Background(m.theme.SurfaceBase)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content,
		lipgloss.WithWhitespaceStyle(bgStyle))
}

// SetSize stores the allocated terminal dimensions for layout.
func (m *welcomeModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}
