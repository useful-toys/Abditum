package tui

import tea "charm.land/bubbletea/v2"

// templateListModel is the left-panel child model active during workAreaTemplates.
// It renders the list of secret templates. Stub in Phase 5 — real implementation
// in Phase 8.
type templateListModel struct {
	width  int
	height int
}

// newTemplateListModel creates a new template list stub.
func newTemplateListModel() *templateListModel {
	return &templateListModel{}
}

// Update processes messages for the template list panel.
func (m *templateListModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// View renders a placeholder for the template list panel.
func (m *templateListModel) View() string {
	return "[template list — Phase 8]"
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

// ChildFlows returns nil.
func (m *templateListModel) ChildFlows() []flowDescriptor {
	return nil
}
