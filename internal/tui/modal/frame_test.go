package modal_test

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/tui/testdata"
)

// renderFrame is a helper that constructs a DialogFrame and renders it.
func renderFrame(
	title, symbol, symbolColor, borderColor, defaultKeyColor string,
	options []modal.ModalOption,
	scroll *modal.ScrollState,
	body string,
) testdata.RenderFn {
	return func(w, h int, theme *design.Theme) string {
		f := modal.DialogFrame{
			Title:           title,
			TitleColor:      theme.Text.Primary,
			Symbol:          symbol,
			SymbolColor:     symbolColor,
			BorderColor:     borderColor,
			Options:         options,
			DefaultKeyColor: defaultKeyColor,
			Scroll:          scroll,
		}
		return f.Render(body, w, theme)
	}
}

func twoOptions(theme *design.Theme) []modal.ModalOption {
	return []modal.ModalOption{
		{Keys: []design.Key{design.Keys.Enter}, Label: "Confirmar", Intent: modal.IntentConfirm, Action: func() tea.Cmd { return nil }},
		{Keys: []design.Key{design.Keys.Esc}, Label: "Cancelar", Intent: modal.IntentCancel, Action: func() tea.Cmd { return nil }},
	}
}

func TestDialogFrame_NoScroll(t *testing.T) {
	theme := design.TokyoNight
	opts := twoOptions(theme)
	body := "Linha 1\nLinha 2\nLinha 3"
	testdata.TestRenderManaged(t, "frame", "no_scroll", []string{"60x10"},
		renderFrame("Título do Diálogo", "", "", theme.Border.Focused, theme.Accent.Primary, opts, nil, body))
}

func TestDialogFrame_WithScrollTop(t *testing.T) {
	theme := design.TokyoNight
	opts := twoOptions(theme)
	scroll := &modal.ScrollState{Offset: 0, Total: 30, Viewport: 8}
	// Body has 8 visible lines (viewport size)
	var lines []string
	for i := 1; i <= 8; i++ {
		lines = append(lines, fmt.Sprintf("Linha %d de 30", i))
	}
	body := strings.Join(lines, "\n")
	testdata.TestRenderManaged(t, "frame", "scroll_top", []string{"60x10"},
		renderFrame("Diálogo com Scroll", "", "", theme.Border.Focused, theme.Accent.Primary, opts, scroll, body))
}

func TestDialogFrame_WithScrollMiddle(t *testing.T) {
	theme := design.TokyoNight
	opts := twoOptions(theme)
	scroll := &modal.ScrollState{Offset: 11, Total: 30, Viewport: 8}
	var lines []string
	for i := 12; i <= 19; i++ {
		lines = append(lines, fmt.Sprintf("Linha %d de 30", i))
	}
	body := strings.Join(lines, "\n")
	testdata.TestRenderManaged(t, "frame", "scroll_middle", []string{"60x10"},
		renderFrame("Diálogo com Scroll", "", "", theme.Border.Focused, theme.Accent.Primary, opts, scroll, body))
}

func TestDialogFrame_WithScrollBottom(t *testing.T) {
	theme := design.TokyoNight
	opts := twoOptions(theme)
	scroll := &modal.ScrollState{Offset: 22, Total: 30, Viewport: 8}
	var lines []string
	for i := 23; i <= 30; i++ {
		lines = append(lines, fmt.Sprintf("Linha %d de 30", i))
	}
	body := strings.Join(lines, "\n")
	testdata.TestRenderManaged(t, "frame", "scroll_bottom", []string{"60x10"},
		renderFrame("Diálogo com Scroll", "", "", theme.Border.Focused, theme.Accent.Primary, opts, scroll, body))
}

func TestDialogFrame_SeverityDestructive(t *testing.T) {
	theme := design.TokyoNight
	sev := design.SeverityDestructive
	opts := []modal.ModalOption{
		{Keys: []design.Key{design.Keys.Enter}, Label: "Excluir", Intent: modal.IntentConfirm, Action: func() tea.Cmd { return nil }},
		{Keys: []design.Key{design.Keys.Esc}, Label: "Cancelar", Intent: modal.IntentCancel, Action: func() tea.Cmd { return nil }},
	}
	body := "Esta ação não pode ser desfeita."
	testdata.TestRenderManaged(t, "frame", "severity_destructive", []string{"60x8"},
		renderFrame("Excluir item", sev.Symbol(), sev.BorderColor(theme), sev.BorderColor(theme), sev.DefaultKeyColor(theme), opts, nil, body))
}

func TestDialogFrame_SeverityError(t *testing.T) {
	theme := design.TokyoNight
	sev := design.SeverityError
	opts := []modal.ModalOption{
		{Keys: []design.Key{design.Keys.Enter}, Label: "OK", Intent: modal.IntentConfirm, Action: func() tea.Cmd { return nil }},
	}
	body := "Ocorreu um erro inesperado."
	testdata.TestRenderManaged(t, "frame", "severity_error", []string{"60x6"},
		renderFrame("Erro", sev.Symbol(), sev.BorderColor(theme), sev.BorderColor(theme), sev.DefaultKeyColor(theme), opts, nil, body))
}
