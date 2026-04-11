package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"fmt"
	"github.com/useful-toys/abditum/internal/vault"
)

// settingsModel fills the full work area during workAreaSettings.
// Stub in Phase 5 - real implementation in Phase 9.
type settingsModel struct {
	mgr     *vault.Manager
	actions *ActionManager
	msgs    *MessageManager
	theme   *Theme
	width   int
	height  int
}

// ApplyTheme applies the given theme to the settingsModel.
func (m *settingsModel) ApplyTheme(t *Theme) {
	m.theme = t
}

// Compile-time assertion: settingsModel satisfies childModel.
var _ childModel = &settingsModel{}

// newSettingsModel creates a new settings stub.
func newSettingsModel(mgr *vault.Manager, actions *ActionManager, msgs *MessageManager, theme *Theme) *settingsModel {
	return &settingsModel{mgr: mgr, actions: actions, msgs: msgs, theme: theme}
}

// Update processes messages for the settings screen.
func (m *settingsModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// View renders a placeholder for the settings screen.
func (m *settingsModel) View() string {
	if m.width == 0 || m.height == 0 {
		panic(fmt.Sprintf("settingsModel.View() called without SetSize: width=%d height=%d", m.width, m.height))
	}
	return lipgloss.NewStyle().Foreground(m.theme.SemanticInfo).
		Render("[settings - Phase 9]")
}

// SetSize stores the allocated screen dimensions.
func (m *settingsModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}
