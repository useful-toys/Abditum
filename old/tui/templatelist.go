package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/vault"
)

// templateListModel is the left-panel child model active during workAreaTemplates.
// Stub in Phase 5 - real implementation in Phase 8.
type templateListModel struct {
	mgr     *vault.Manager
	actions *ActionManager
	msgs    *MessageManager
}

// Compile-time assertion: templateListModel satisfies childModel.
var _ childModel = &templateListModel{}

// newTemplateListModel creates a new template list stub.
func newTemplateListModel(mgr *vault.Manager, actions *ActionManager, msgs *MessageManager) *templateListModel {
	return &templateListModel{mgr: mgr, actions: actions, msgs: msgs}
}

// Update processes messages for the template list.
func (m *templateListModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// View renders a placeholder for the template list panel.
func (m *templateListModel) View(width, height int, theme *Theme) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Semantic.Info)).
		Render("[template list - Phase 8]")
}
