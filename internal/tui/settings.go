package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/vault"
)

// settingsModel fills the full work area during workAreaSettings.
// Stub in Phase 5 - real implementation in Phase 9.
type settingsModel struct {
	mgr     *vault.Manager
	actions *ActionManager
	msgs    *MessageManager
	width   int
	height  int
}

// Compile-time assertion: settingsModel satisfies childModel.
var _ childModel = &settingsModel{}

// newSettingsModel creates a new settings stub.
func newSettingsModel(mgr *vault.Manager, actions *ActionManager, msgs *MessageManager) *settingsModel {
	return &settingsModel{mgr: mgr, actions: actions, msgs: msgs}
}

// Update processes messages for the settings screen.
func (m *settingsModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// View renders a placeholder for the settings screen.
func (m *settingsModel) View() string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("245")).
		Render("[settings - Phase 9]")
}

// SetSize stores the allocated screen dimensions.
func (m *settingsModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}
