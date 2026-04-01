package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// preVaultModel renders the welcome background (ASCII art logo + action hints).
// It is active during workAreaPreVault and has no sub-states.
// Open/create vault flows are orchestrated via the modal stack, not this model.
type preVaultModel struct {
	actions *ActionManager
	width   int
	height  int
}

// Compile-time assertion: preVaultModel satisfies childModel.
var _ childModel = &preVaultModel{}

// newPreVaultModel creates a new pre-vault welcome screen model.
func newPreVaultModel(actions *ActionManager) *preVaultModel {
	return &preVaultModel{actions: actions}
}

// Update processes messages for the pre-vault screen.
// Phase 5: preVaultModel is display-only. No input handling until Phase 6.
func (m *preVaultModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// View renders the ASCII art logo centered on screen with a "No vault open" hint.
func (m *preVaultModel) View() string {
	logoBlock := lipgloss.NewStyle().Width(43).Render(RenderLogo())
	hint := lipgloss.NewStyle().Foreground(lipgloss.Color("245")).
		Render("No vault open")
	content := logoBlock + "\n\n" + hint

	if m.width == 0 || m.height == 0 {
		return content
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
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
// All flows (open/create vault) are globally registered in FlowRegistry.
func (m *preVaultModel) ChildFlows() []flowDescriptor {
	return nil
}

// renderHints renders a list of hint lines using the muted style.
// Used by preVaultModel.View and any other model that shows key hints.
func renderHints(items []string) string {
	var b strings.Builder
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	for _, item := range items {
		b.WriteString(style.Render(item) + "\n")
	}
	return b.String()
}
