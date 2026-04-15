package design

import (
	"regexp"
	"testing"

	"charm.land/lipgloss/v2"
)

var ansiEscapeRe = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func TestRenderAction_WidthMatchesLipglossWidth(t *testing.T) {
	theme := TokyoNight
	ra := RenderAction("⌃S", "Salvar", theme)

	measured := lipgloss.Width(ra.Text)
	if measured != ra.Width {
		t.Errorf("RenderAction: Width = %d, lipgloss.Width(Text) = %d", ra.Width, measured)
	}
}

func TestRenderAction_ContainsKeyAndLabel(t *testing.T) {
	theme := TokyoNight
	ra := RenderAction("F1", "Ajuda", theme)

	clean := ansiEscapeRe.ReplaceAllString(ra.Text, "")
	if len(clean) == 0 {
		t.Fatal("RenderAction: texto limpo está vazio")
	}
	if !containsSubstring(clean, "F1") {
		t.Errorf("RenderAction: texto limpo %q não contém tecla 'F1'", clean)
	}
	if !containsSubstring(clean, "Ajuda") {
		t.Errorf("RenderAction: texto limpo %q não contém rótulo 'Ajuda'", clean)
	}
}

func TestActionSeparator_Width(t *testing.T) {
	theme := TokyoNight
	sep := ActionSeparator(theme)

	if sep.Width != 3 {
		t.Errorf("ActionSeparator: Width = %d, want 3", sep.Width)
	}
	if lipgloss.Width(sep.Text) != 3 {
		t.Errorf("ActionSeparator: lipgloss.Width = %d, want 3", lipgloss.Width(sep.Text))
	}
}

func TestRenderAction_WidthPositive(t *testing.T) {
	theme := TokyoNight
	ra := RenderAction("⌃Q", "Sair", theme)
	if ra.Width <= 0 {
		t.Errorf("RenderAction: Width = %d, deve ser > 0", ra.Width)
	}
}

// containsSubstring verifica se s contém sub.
func containsSubstring(s, sub string) bool {
	if len(sub) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
