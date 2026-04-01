package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/vault"
)

// templateListModel is the left-panel child model active during workAreaTemplates.
// It renders the list of secret templates.
// Stub in Phase 5 — real implementation in Phase 8.
type templateListModel struct {
	mgr     *vault.Manager
	actions *ActionManager
	msgs    *MessageManager
	width   int
	height  int
}

// Compile-time assertion: templateListModel satisfies childModel.
var _ childModel = &templateListModel{}

// newTemplateListModel creates a new template list stub.
func newTemplateListModel(mgr *vault.Manager, actions *ActionManager, msgs *MessageManager) *templateListModel {
	return &templateListModel{mgr: mgr, actions: actions, msgs: msgs}
}

// Update processes messages for the template list panel.
func (m *templateListModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// View renders a placeholder for the template list panel.
// Phase 8 will replace this with the real template list.
func (m *templateListModel) View() string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("245")).
		Render("[template list — Phase 8]")
}

// SetSize stores the allocated panel dimensions.
func (m *templateListModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Context returns the current navigation context for flow dispatch.
func (m *templateListModel) Context() FlowContext {
	return FlowContext{}
}

// ChildFlows returns nil — no child-specific flows in stub.
func (m *templateListModel) ChildFlows() []flowDescriptor {
	return nil
}
