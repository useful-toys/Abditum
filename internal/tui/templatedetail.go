package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/vault"
)

// templateDetailModel is the right-panel child model active during workAreaTemplates.
// Stub in Phase 5 - real implementation in Phase 8.
type templateDetailModel struct {
	mgr     *vault.Manager
	actions *ActionManager
	msgs    *MessageManager
	width   int
	height  int
}

// Compile-time assertion: templateDetailModel satisfies childModel.
var _ childModel = &templateDetailModel{}

// newTemplateDetailModel creates a new template detail stub.
func newTemplateDetailModel(mgr *vault.Manager, actions *ActionManager, msgs *MessageManager) *templateDetailModel {
	return &templateDetailModel{mgr: mgr, actions: actions, msgs: msgs}
}

// Update processes messages for the template detail panel.
func (m *templateDetailModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// View renders a placeholder for the template detail panel.
func (m *templateDetailModel) View() string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("245")).
		Render("[template detail - Phase 8]")
}

// SetSize stores the allocated panel dimensions.
func (m *templateDetailModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}
