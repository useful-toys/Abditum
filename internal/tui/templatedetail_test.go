package tui

import (
	"testing"
)

// TestTemplateDetailModel_ViewPanicsWithoutSetSize verifies that calling View()
// without first calling SetSize() results in a panic — the rootModel contract.
func TestTemplateDetailModel_ViewPanicsWithoutSetSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("templateDetailModel.View() should panic without SetSize")
		}
	}()
	tdm := newTemplateDetailModel(nil, nil, nil, ThemeTokyoNight)
	tdm.View() // Should panic here
}
