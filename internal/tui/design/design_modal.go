package design

import "charm.land/lipgloss/v2"

// DialogPaddingH é o padding horizontal interno dos diálogos (colunas em cada lado).
const DialogPaddingH = 2

// DialogPaddingV é o padding vertical interno dos diálogos (linhas acima e abaixo do corpo).
// O HelpModal usa 0 linhas de padding vertical — não inclui linhas em branco no corpo.
const DialogPaddingV = 1

// Severity representa a severidade visual de um diálogo de Notificação ou Confirmação.
// Diálogos de Ajuda e Funcionais não usam severidade.
type Severity int

const (
	// SeverityNeutral não tem símbolo; usa border.focused; tecla default em accent.primary.
	SeverityNeutral Severity = iota
	// SeverityInformative tem símbolo ℹ e usa semantic.info como cor de borda.
	SeverityInformative
	// SeverityAlert tem símbolo ⚠ e usa semantic.warning como cor de borda.
	SeverityAlert
	// SeverityDestructive tem símbolo ⚠, usa semantic.warning como cor de borda,
	// mas a tecla default usa semantic.error (ação destrutiva irrecuperável).
	SeverityDestructive
	// SeverityError tem símbolo ✕ e usa semantic.error como cor de borda.
	SeverityError
)

// Symbol retorna o símbolo Unicode da severidade.
// Retorna "" para SeverityNeutral.
// SeverityAlert e SeverityDestructive retornam ambos SymWarning — a distinção
// visual está na cor da tecla default (DefaultKeyColor), não no símbolo.
func (s Severity) Symbol() string {
	switch s {
	case SeverityInformative:
		return SymInfo
	case SeverityAlert, SeverityDestructive:
		return SymWarning
	case SeverityError:
		return SymError
	default: // SeverityNeutral
		return ""
	}
}

// BorderColor retorna a cor de borda para a severidade a partir do tema.
func (s Severity) BorderColor(theme *Theme) string {
	switch s {
	case SeverityInformative:
		return theme.Semantic.Info
	case SeverityAlert, SeverityDestructive:
		return theme.Semantic.Warning
	case SeverityError:
		return theme.Semantic.Error
	default: // SeverityNeutral
		return theme.Border.Focused
	}
}

// DefaultKeyColor retorna a cor da tecla da ação default (primeira opção) para a severidade.
// Para SeverityDestructive, a tecla default usa semantic.error para enfatizar o risco.
// Todas as demais severidades usam accent.primary.
func (s Severity) DefaultKeyColor(theme *Theme) string {
	if s == SeverityDestructive {
		return theme.Semantic.Error
	}
	return theme.Accent.Primary
}

// RenderDialogTitle renderiza o bloco título do diálogo.
// Se symbol != "", inclui "symbol  title" (símbolo + 2 espaços + título).
// Se symbol == "", inclui apenas "title".
// Cores: symbol em symbolColor, título em theme.Text.Primary + bold.
// Retorna texto ANSI e largura visual em colunas.
func RenderDialogTitle(title, symbol, symbolColor string, theme *Theme) (string, int) {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Bold(true)

	var rendered string
	if symbol != "" {
		symStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(symbolColor))
		rendered = symStyle.Render(symbol) + "  " + titleStyle.Render(title)
	} else {
		rendered = titleStyle.Render(title)
	}
	return rendered, lipgloss.Width(rendered)
}

// RenderDialogAction renderiza uma ação do rodapé: "key label".
// key é o Label da tecla (Keys[0].Label da ModalOption — ex: "Enter", "S", "Esc").
// key é renderizada em keyColor; label é renderizada em theme.Text.Primary.
// Retorna texto ANSI e largura visual em colunas.
func RenderDialogAction(key, label, keyColor string, theme *Theme) (string, int) {
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(keyColor)).
		Bold(true)
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Text.Primary))

	rendered := keyStyle.Render(key) + " " + labelStyle.Render(label)
	return rendered, lipgloss.Width(rendered)
}

// RenderScrollArrow renderiza ↑ (up=true) ou ↓ (up=false) em theme.Text.Secondary.
// Sempre retorna width == 1.
func RenderScrollArrow(up bool, theme *Theme) (string, int) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Secondary))
	sym := SymScrollDown
	if up {
		sym = SymScrollUp
	}
	return style.Render(sym), 1
}

// RenderScrollThumb renderiza ■ em theme.Text.Secondary.
// Sempre retorna width == 1.
func RenderScrollThumb(theme *Theme) (string, int) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Secondary))
	return style.Render(SymScrollThumb), 1
}
