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
	version string // Application version to display below logo
}

// Compile-time assertion: welcomeModel satisfies childModel.
var _ childModel = &welcomeModel{}

// newWelcomeModel creates a new welcome screen model.
func newWelcomeModel(actions *ActionManager, version string) *welcomeModel {
	return &welcomeModel{actions: actions, version: version}
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
func (m *welcomeModel) View(width, height int, theme *Theme) string {
	// 43 = width of AsciiArt (const in ascii.go) — each line is exactly 43 characters.
	// No background is set here: the root workAreaStyle already applies SurfaceBase
	// to the entire work area. Setting background here would emit redundant SGR codes
	// that may conflict with the terminal's own background rendering.
	logoBlock := lipgloss.NewStyle().Width(43).Render(RenderLogo(theme))

	// Format version with semantic.secondary color (from theme)
	// Per spec: version token = text.secondary
	versionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Secondary))
	versionLine := versionStyle.Render(m.version)

	content := lipgloss.JoinVertical(lipgloss.Center, logoBlock, "", versionLine)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}
