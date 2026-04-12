package tui

import (
	"fmt"

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
	theme   *Theme
	width   int
	height  int
}

// ApplyTheme applies the given theme to the vaultTreeModel.
func (m *vaultTreeModel) ApplyTheme(t *Theme) {
	m.theme = t
}

// Compile-time assertion: vaultTreeModel satisfies childModel.
var _ childModel = &vaultTreeModel{}

// newVaultTreeModel creates a new vault tree stub.
func newVaultTreeModel(mgr *vault.Manager, actions *ActionManager, msgs *MessageManager, theme *Theme) *vaultTreeModel {
	return &vaultTreeModel{mgr: mgr, actions: actions, msgs: msgs, theme: theme}
}

// Update processes messages for the vault tree.
func (m *vaultTreeModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// View renders a placeholder for the vault tree panel.
func (m *vaultTreeModel) View() string {
	if m.width == 0 || m.height == 0 {
		panic(fmt.Sprintf("vaultTreeModel.View() called without SetSize: width=%d height=%d", m.width, m.height))
	}
	return lipgloss.NewStyle().Foreground(m.theme.SemanticInfo).
		Render("[vault tree - Phase 7]")
}

// SetSize stores the allocated panel dimensions.
func (m *vaultTreeModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
