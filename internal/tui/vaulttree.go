package tui

import tea "charm.land/bubbletea/v2"

// vaultTreeModel is the left-panel child model active during workAreaVault.
// It renders the hierarchical folder/secret tree. Stub in Phase 5 — real
// implementation in Phase 7.
type vaultTreeModel struct {
	width  int
	height int
}

// newVaultTreeModel creates a new vault tree stub.
func newVaultTreeModel() *vaultTreeModel {
	return &vaultTreeModel{}
}

// Update processes messages for the vault tree.
func (m *vaultTreeModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// View renders a placeholder for the vault tree panel.
func (m *vaultTreeModel) View() string {
	return "[vault tree — Phase 7]"
}

// SetSize stores the allocated panel dimensions.
func (m *vaultTreeModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Context returns the current navigation context for flow dispatch.
// Stub: no selection state yet.
func (m *vaultTreeModel) Context() FlowContext {
	return FlowContext{}
}

// ChildFlows returns nil.
func (m *vaultTreeModel) ChildFlows() []flowDescriptor {
	return nil
}
