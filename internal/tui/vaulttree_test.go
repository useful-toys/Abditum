package tui

import (
	"testing"
)

// TestVaultTreeModel_ViewPanicsWithoutSetSize verifies that calling View()
// without first calling SetSize() results in a panic — the rootModel contract.
func TestVaultTreeModel_ViewPanicsWithoutSetSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("vaultTreeModel.View() should panic without SetSize")
		}
	}()
	vm := newVaultTreeModel(nil, nil, nil, ThemeTokyoNight)
	vm.View() // Should panic here
}
