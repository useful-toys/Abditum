package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/vault"
)

// vaultTreeModel is the left-panel child model active during workAreaVault.
// Stub in Phase 5 - real tree in Phase 7.
type vaultTreeModel struct {
	mgr     *vault.Manager
	actions *ActionManager
	msgs    *MessageManager
}

// Compile-time assertion: vaultTreeModel satisfies childModel.
var _ childModel = &vaultTreeModel{}

// newVaultTreeModel creates a new vault tree stub.
func newVaultTreeModel(mgr *vault.Manager, actions *ActionManager, msgs *MessageManager) *vaultTreeModel {
	return &vaultTreeModel{mgr: mgr, actions: actions, msgs: msgs}
}

// Update processes messages for the vault tree.
func (m *vaultTreeModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// View renders a placeholder for the vault tree panel.
func (m *vaultTreeModel) View(width, height int, theme *Theme) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Semantic.Info)).
		Render("[vault tree - Phase 7]")
}
