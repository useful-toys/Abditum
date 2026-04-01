package tui

import tea "charm.land/bubbletea/v2"

// templateDetailModel is the right-panel child model active during workAreaTemplates.
// It renders the focused template's field definitions. Stub in Phase 5 — real
// implementation in Phase 8.
type templateDetailModel struct {
	width  int
	height int
}

// newTemplateDetailModel creates a new template detail stub.
func newTemplateDetailModel() *templateDetailModel {
	return &templateDetailModel{}
}

// Update processes messages for the template detail panel.
func (m *templateDetailModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// View renders a placeholder for the template detail panel.
func (m *templateDetailModel) View() string {
	return "[template detail — Phase 8]"
}

// SetSize stores the allocated panel dimensions.
func (m *templateDetailModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Context returns the current navigation context for flow dispatch.
func (m *templateDetailModel) Context() FlowContext {
	return FlowContext{}
}

// ChildFlows returns nil.
func (m *templateDetailModel) ChildFlows() []flowDescriptor {
	return nil
}
