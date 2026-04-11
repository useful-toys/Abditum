package tui

import (
	"testing"
)

// TestSecretDetailModel_ViewPanicsWithoutSetSize verifies that calling View()
// without first calling SetSize() results in a panic — the rootModel contract.
func TestSecretDetailModel_ViewPanicsWithoutSetSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("secretDetailModel.View() should panic without SetSize")
		}
	}()
	sdm := newSecretDetailModel(nil, nil, nil, ThemeTokyoNight)
	sdm.View() // Should panic here
}
