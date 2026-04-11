package tui

import (
	"testing"
)

// TestSettingsModel_ViewPanicsWithoutSetSize verifies that calling View()
// without first calling SetSize() results in a panic — the rootModel contract.
func TestSettingsModel_ViewPanicsWithoutSetSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("settingsModel.View() should panic without SetSize")
		}
	}()
	sm := newSettingsModel(nil, nil, nil, ThemeTokyoNight)
	sm.View() // Should panic here
}
