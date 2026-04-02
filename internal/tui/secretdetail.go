package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/vault"
)

// secretDetailModel is the right-panel child model active during workAreaVault.
// Stub in Phase 5 - real implementation in Phase 8.
type secretDetailModel struct {
	mgr     *vault.Manager
	actions *ActionManager
	msgs    *MessageManager
	width   int
	height  int
}

// Compile-time assertion: secretDetailModel satisfies childModel.
var _ childModel = &secretDetailModel{}

// newSecretDetailModel creates a new secret detail stub.
func newSecretDetailModel(mgr *vault.Manager, actions *ActionManager, msgs *MessageManager) *secretDetailModel {
	return &secretDetailModel{mgr: mgr, actions: actions, msgs: msgs}
}

// Update processes messages for the secret detail panel.
func (m *secretDetailModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// View renders a placeholder for the secret detail panel.
func (m *secretDetailModel) View() string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("245")).
		Render("[secret detail - Phase 8]")
}

// SetSize stores the allocated panel dimensions.
func (m *secretDetailModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}
