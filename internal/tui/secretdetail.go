package tui

import tea "charm.land/bubbletea/v2"

// secretDetailModel is the right-panel child model active during workAreaVault.
// It renders the focused secret's fields and values. Stub in Phase 5 — real
// implementation in Phase 8.
type secretDetailModel struct {
	width  int
	height int
}

// newSecretDetailModel creates a new secret detail stub.
func newSecretDetailModel() *secretDetailModel {
	return &secretDetailModel{}
}

// Update processes messages for the secret detail panel.
func (m *secretDetailModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// View renders a placeholder for the secret detail panel.
func (m *secretDetailModel) View() string {
	return "[secret detail — Phase 8]"
}

// SetSize stores the allocated panel dimensions.
func (m *secretDetailModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Context returns the current navigation context for flow dispatch.
// Stub: no selection state yet.
func (m *secretDetailModel) Context() FlowContext {
	return FlowContext{}
}

// ChildFlows returns nil.
func (m *secretDetailModel) ChildFlows() []flowDescriptor {
	return nil
}
