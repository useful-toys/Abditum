package design

import (
	"testing"

	"charm.land/lipgloss/v2"
)

// --- Severity ---

func TestSeverity_Symbol(t *testing.T) {
	tests := []struct {
		sev  Severity
		want string
	}{
		{SeverityNeutral, ""},
		{SeverityInformative, SymInfo},
		{SeverityAlert, SymWarning},
		{SeverityDestructive, SymWarning},
		{SeverityError, SymError},
	}
	for _, tt := range tests {
		if got := tt.sev.Symbol(); got != tt.want {
			t.Errorf("Severity(%d).Symbol() = %q, want %q", tt.sev, got, tt.want)
		}
	}
}

func TestSeverity_BorderColor_NotEmpty(t *testing.T) {
	theme := TokyoNight
	sevs := []Severity{SeverityNeutral, SeverityInformative, SeverityAlert, SeverityDestructive, SeverityError}
	for _, s := range sevs {
		if c := s.BorderColor(theme); c == "" {
			t.Errorf("Severity(%d).BorderColor(TokyoNight) = empty string", s)
		}
	}
}

func TestSeverity_DefaultKeyColor_NotEmpty(t *testing.T) {
	theme := TokyoNight
	sevs := []Severity{SeverityNeutral, SeverityInformative, SeverityAlert, SeverityDestructive, SeverityError}
	for _, s := range sevs {
		if c := s.DefaultKeyColor(theme); c == "" {
			t.Errorf("Severity(%d).DefaultKeyColor(TokyoNight) = empty string", s)
		}
	}
}

func TestSeverityDestructive_DefaultKeyColor_IsSemanticError(t *testing.T) {
	theme := TokyoNight
	got := SeverityDestructive.DefaultKeyColor(theme)
	want := theme.Semantic.Error
	if got != want {
		t.Errorf("SeverityDestructive.DefaultKeyColor = %q, want Semantic.Error = %q", got, want)
	}
}

// --- RenderDialogTitle ---

func TestRenderDialogTitle_WithoutSymbol(t *testing.T) {
	theme := TokyoNight
	text, width := RenderDialogTitle("Título", "", "", theme)
	if width <= 0 {
		t.Errorf("RenderDialogTitle width = %d, want > 0", width)
	}
	if lipgloss.Width(text) != width {
		t.Errorf("RenderDialogTitle: returned width %d != lipgloss.Width %d", width, lipgloss.Width(text))
	}
	clean := ansiEscapeRe.ReplaceAllString(text, "")
	if !containsSubstring(clean, "Título") {
		t.Errorf("RenderDialogTitle: clean text %q does not contain title", clean)
	}
}

func TestRenderDialogTitle_WithSymbol(t *testing.T) {
	theme := TokyoNight
	text, width := RenderDialogTitle("Aviso", SymWarning, theme.Semantic.Warning, theme)
	clean := ansiEscapeRe.ReplaceAllString(text, "")
	if !containsSubstring(clean, SymWarning) {
		t.Errorf("RenderDialogTitle with symbol: clean text %q does not contain symbol", clean)
	}
	if !containsSubstring(clean, "Aviso") {
		t.Errorf("RenderDialogTitle with symbol: clean text %q does not contain title", clean)
	}
	if lipgloss.Width(text) != width {
		t.Errorf("RenderDialogTitle: returned width %d != lipgloss.Width %d", width, lipgloss.Width(text))
	}
}

// --- RenderDialogAction ---

func TestRenderDialogAction_ContainsKeyAndLabel(t *testing.T) {
	theme := TokyoNight
	text, width := RenderDialogAction("Enter", "Confirmar", theme.Accent.Primary, theme)
	clean := ansiEscapeRe.ReplaceAllString(text, "")
	if !containsSubstring(clean, "Enter") {
		t.Errorf("RenderDialogAction: clean %q does not contain key", clean)
	}
	if !containsSubstring(clean, "Confirmar") {
		t.Errorf("RenderDialogAction: clean %q does not contain label", clean)
	}
	if lipgloss.Width(text) != width {
		t.Errorf("RenderDialogAction: returned width %d != lipgloss.Width %d", width, lipgloss.Width(text))
	}
}

// --- RenderScrollArrow ---

func TestRenderScrollArrow_Width(t *testing.T) {
	theme := TokyoNight
	upText, upW := RenderScrollArrow(true, theme)
	if upW != 1 {
		t.Errorf("RenderScrollArrow(up) width = %d, want 1", upW)
	}
	if lipgloss.Width(upText) != upW {
		t.Errorf("RenderScrollArrow(up) lipgloss.Width mismatch")
	}
	downText, downW := RenderScrollArrow(false, theme)
	if downW != 1 {
		t.Errorf("RenderScrollArrow(down) width = %d, want 1", downW)
	}
	if lipgloss.Width(downText) != downW {
		t.Errorf("RenderScrollArrow(down) lipgloss.Width mismatch")
	}
}

func TestRenderScrollArrow_CorrectSymbol(t *testing.T) {
	theme := TokyoNight
	upText, _ := RenderScrollArrow(true, theme)
	downText, _ := RenderScrollArrow(false, theme)
	cleanUp := ansiEscapeRe.ReplaceAllString(upText, "")
	cleanDown := ansiEscapeRe.ReplaceAllString(downText, "")
	if cleanUp != SymScrollUp {
		t.Errorf("RenderScrollArrow(true) = %q, want %q", cleanUp, SymScrollUp)
	}
	if cleanDown != SymScrollDown {
		t.Errorf("RenderScrollArrow(false) = %q, want %q", cleanDown, SymScrollDown)
	}
}

// --- RenderScrollThumb ---

func TestRenderScrollThumb_Width(t *testing.T) {
	theme := TokyoNight
	text, w := RenderScrollThumb(theme)
	if w != 1 {
		t.Errorf("RenderScrollThumb width = %d, want 1", w)
	}
	if lipgloss.Width(text) != w {
		t.Errorf("RenderScrollThumb lipgloss.Width mismatch")
	}
	clean := ansiEscapeRe.ReplaceAllString(text, "")
	if clean != SymScrollThumb {
		t.Errorf("RenderScrollThumb = %q, want %q", clean, SymScrollThumb)
	}
}
