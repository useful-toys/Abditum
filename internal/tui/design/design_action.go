package design

import "charm.land/lipgloss/v2"

// RenderedAction encapsula o texto ANSI estilizado de uma Action e sua largura em colunas.
// Usar Width (em vez de len) porque sequências ANSI não contribuem para a largura visual.
type RenderedAction struct {
	Text  string // texto estilizado com sequências ANSI
	Width int    // largura visual em colunas (lipgloss.Width — nunca len)
}

// RenderAction renderiza uma action: tecla (Accent.Primary + bold) + espaço + rótulo (Text.Primary).
// key é o rótulo da tecla (ex: "⌃S", "F1"); label é o texto descritivo (ex: "Salvar", "Ajuda").
// Retorna RenderedAction com texto ANSI e largura em colunas.
func RenderAction(key, label string, theme *Theme) RenderedAction {
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Accent.Primary)).
		Bold(true)
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Text.Primary))

	rendered := keyStyle.Render(key) + " " + labelStyle.Render(label)
	return RenderedAction{
		Text:  rendered,
		Width: lipgloss.Width(rendered),
	}
}

// ActionSeparator retorna o separador " · " (espaço + SymHeaderSep + espaço) estilizado
// com Text.Secondary. Sempre tem Width == 3.
func ActionSeparator(theme *Theme) RenderedAction {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Secondary))
	rendered := style.Render(" " + SymHeaderSep + " ")
	return RenderedAction{
		Text:  rendered,
		Width: lipgloss.Width(rendered),
	}
}
