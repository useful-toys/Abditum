package tui

import (
	"testing"
)

// TestTemplateListModel_ViewPanicsWithoutSetSize verifies that calling View()
// without first calling SetSize() results in a panic — the rootModel contract.
func TestTemplateListModel_ViewPanicsWithoutSetSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("templateListModel.View() should panic without SetSize")
		}
	}()
	tlm := newTemplateListModel(nil, nil, nil, ThemeTokyoNight)
	tlm.View() // Should panic here
}
