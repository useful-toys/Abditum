package tui

import tea "charm.land/bubbletea/v2"

// preVaultModel is the welcome screen child model, active during workAreaPreVault.
// It renders the ASCII art logo and initial action hints as a static background.
// It does NOT manage open/create vault sub-states — those are modal flows (D-09).
type preVaultModel struct {
	width  int
	height int
}

// newPreVaultModel creates a new pre-vault welcome screen model.
func newPreVaultModel() *preVaultModel {
	return &preVaultModel{}
}

// Update processes messages for the pre-vault screen.
// Most messages are ignored — the screen is essentially static.
func (m *preVaultModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// View renders the ASCII art logo and action hints.
// Phase 5 placeholder: displays the logo only.
func (m *preVaultModel) View() string {
	return RenderLogo()
}

// SetSize stores the allocated terminal dimensions for layout.
func (m *preVaultModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Context returns an empty FlowContext.
// The pre-vault screen has no focused vault items.
func (m *preVaultModel) Context() FlowContext {
	return FlowContext{}
}

// ChildFlows returns nil.
// The pre-vault screen has no child-specific flows — all flows (open/create vault)
// are globally registered in FlowRegistry.
func (m *preVaultModel) ChildFlows() []flowDescriptor {
	return nil
}
