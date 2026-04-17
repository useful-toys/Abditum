package modal

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// DialogFrame define a aparência visual de um diálogo: borda superior com título,
// bordas laterais com indicadores de scroll opcionais, e borda de rodapé com ações.
// Não tem estado próprio — mesmos argumentos = mesmo output (idempotente).
type DialogFrame struct {
	// Title é o texto do cabeçalho.
	Title string
	// TitleColor é a cor do título — ex: theme.Text.Primary. Nunca hardcoded.
	TitleColor string
	// Symbol é o símbolo de severidade ou "" para omitir — ex: design.SymWarning.
	Symbol string
	// SymbolColor é a cor do símbolo — ex: design.SeverityDestructive.BorderColor(theme).
	SymbolColor string
	// BorderColor é a cor de toda a borda — ex: theme.Border.Focused.
	BorderColor string
	// Options lista as ações do rodapé (máximo 3 conforme DS).
	// A 1ª opção usa DefaultKeyColor; as demais usam BorderColor.
	Options []ModalOption
	// DefaultKeyColor é a cor da tecla da 1ª opção (ação principal).
	DefaultKeyColor string
	// Scroll é o estado de scroll para exibir indicadores na borda lateral direita.
	// nil = sem scroll.
	Scroll *ScrollState
}

// Render monta a string completa do diálogo a partir do corpo fornecido.
//
// body é uma string com linhas separadas por \n. Cada linha já deve estar renderizada
// com ANSI e ter largura visual de (maxWidth - 2 - 2*DialogPaddingH) colunas.
// O frame não reaplica padding horizontal — apenas adiciona as bordas laterais.
//
// Algoritmo:
//  1. Borda superior: ╭── [símbolo  ]título ───╮
//  2. Para cada linha do body:
//     │ [padding] linha [padding] │  com background surface.raised
//     Se Scroll != nil, substitui │ direito por indicador de scroll.
//  3. Borda de rodapé com ações posicionadas.
func (f DialogFrame) Render(body string, maxWidth int, theme *design.Theme) string {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(f.BorderColor))
	bgStyle := lipgloss.NewStyle().Background(lipgloss.Color(theme.Surface.Raised))

	// Calcular largura mínima necessária
	minWidth := f.calculateMinWidth(body, theme)
	innerWidth := maxWidth - 2 // subtrair as duas bordas verticais
	if innerWidth > minWidth {
		innerWidth = minWidth
	}

	// --- Borda superior ---
	topLine := f.renderTopBorder(innerWidth, borderStyle, theme)

	// --- Linhas do corpo ---
	lines := strings.Split(body, "\n")
	var bodyLines []string
	for i, line := range lines {
		bodyLines = append(bodyLines, f.renderBodyLine(i+1, len(lines), line, innerWidth, borderStyle, bgStyle, theme))
	}

	// --- Borda de rodapé ---
	bottomLine := f.renderBottomBorder(innerWidth, borderStyle, theme)

	var sb strings.Builder
	sb.WriteString(topLine)
	sb.WriteRune('\n')
	for _, l := range bodyLines {
		sb.WriteString(l)
		sb.WriteRune('\n')
	}
	sb.WriteString(bottomLine)
	return sb.String()
}

// renderTopBorder gera a linha: ╭── [símbolo  título] ───╮
func (f DialogFrame) renderTopBorder(innerWidth int, borderStyle lipgloss.Style, theme *design.Theme) string {
	titleText, titleWidth := design.RenderDialogTitle(f.Title, f.Symbol, f.SymbolColor, theme)

	// " título " — 1 espaço antes, 1 espaço depois
	titleSegmentWidth := titleWidth + 2
	// preenchimento ─ à esquerda (mínimo 1) e à direita
	fillLeft := 1
	fillRight := innerWidth - fillLeft - titleSegmentWidth
	if fillRight < 1 {
		fillRight = 1
	}

	tl := borderStyle.Render(design.SymCornerTL)
	tr := borderStyle.Render(design.SymCornerTR)
	dashL := borderStyle.Render(strings.Repeat(design.SymBorderH, fillLeft))
	dashR := borderStyle.Render(strings.Repeat(design.SymBorderH, fillRight))
	space := " "

	return tl + dashL + space + titleText + space + dashR + tr
}

// renderBodyLine gera uma linha do corpo com bordas e indicador de scroll.
// lineNum é 1-based dentro das linhas visíveis (para determinar posição das setas).
func (f DialogFrame) renderBodyLine(lineNum, totalLines int, content string, innerWidth int, borderStyle lipgloss.Style, bgStyle lipgloss.Style, theme *design.Theme) string {
	lBorder := borderStyle.Render(design.SymBorderV)
	paddingH := strings.Repeat(" ", design.DialogPaddingH)

	// Conteúdo com background
	lineContent := bgStyle.Render(paddingH + content + paddingH)

	// Borda direita: seta ou thumb ou │ normal
	var rBorder string
	if f.Scroll != nil {
		isFirstLine := lineNum == 1
		isLastLine := lineNum == totalLines
		thumbLine := f.Scroll.ThumbLine()

		switch {
		case isFirstLine && f.Scroll.CanScrollUp():
			arrow, _ := design.RenderScrollArrow(true, theme)
			rBorder = arrow
		case isLastLine && f.Scroll.CanScrollDown():
			arrow, _ := design.RenderScrollArrow(false, theme)
			rBorder = arrow
		case thumbLine != -1 && lineNum == thumbLine:
			thumb, _ := design.RenderScrollThumb(theme)
			rBorder = thumb
		default:
			rBorder = borderStyle.Render(design.SymBorderV)
		}
	} else {
		rBorder = borderStyle.Render(design.SymBorderV)
	}

	return lBorder + lineContent + rBorder
}

// calculateMinWidth calcula a largura mínima necessária para o diálogo.
// Largura mínima = mínimo entre 20 colunas e o necessário para título + ações.
func (f DialogFrame) calculateMinWidth(body string, theme *design.Theme) int {
	minWidth := 20 // mínimo absoluto

	// Largura do título na borda superior
	titleWidth := lipgloss.Width(f.Title)
	if f.Symbol != "" {
		titleWidth += lipgloss.Width(f.Symbol) + 2 // símbolo + 2 espaços
	}
	// " título " = +2, " ╭──" = 4, " ╮" = 2
	titleWidth += 8

	// Largura das ações
	actionWidth := 4 // cantos + preenchimento mínimo: "╰─ " + " ─╯"
	for _, opt := range f.Options {
		if len(opt.Keys) == 0 {
			continue
		}
		_, keyWidth := design.RenderDialogAction(opt.Keys[0].Label, opt.Label, f.BorderColor, theme)
		actionWidth += keyWidth + 4 // ação + 2 espaços + 1 dash mínimo
	}

	// Largura do corpo (cada linha)
	bodyLines := strings.Split(body, "\n")
	bodyWidth := 0
	paddingH := 2 * design.DialogPaddingH
	for _, line := range bodyLines {
		w := lipgloss.Width(line) + paddingH + 2 // +2 para bordas │
		if w > bodyWidth {
			bodyWidth = w
		}
	}

	// Largura mínima = máximo entre título, ações e corpo (com mínimo de 20)
	width := titleWidth
	if actionWidth > width {
		width = actionWidth
	}
	if bodyWidth > width {
		width = bodyWidth
	}
	if width < minWidth {
		width = minWidth
	}

	return width + 2 // +2 para bordas verticais
}

// renderBottomBorder gera a linha de rodapé:
// Formato: "╰─ " + Ação1 + " ─ " + Ação2 + " ─ " + Ação3 + " ─╯"
//
// Posicionamento:
//
//	1 ação  → "╰─ " + Ação + " ─╯" (alinhada à direita)
//	2 ações → "╰─ " + Ação1 + " ─ " + Ação2 + " ─╯"
//	3 ações → "╰─ " + Ação1 + " ─ " + Ação2 + " ─ " + Ação3 + " ─╯"
func (f DialogFrame) renderBottomBorder(innerWidth int, borderStyle lipgloss.Style, theme *design.Theme) string {
	// Cantos inferiores conforme especificação: ╰─  e  ─╯
	bl := borderStyle.Render(design.SymCornerBL) + borderStyle.Render(design.SymBorderH) + " "
	br := " " + borderStyle.Render(design.SymBorderH) + borderStyle.Render(design.SymCornerBR)

	// Renderizar as ações (cada ação já inclui espaços internos via design.RenderDialogAction)
	type renderedOpt struct {
		text  string
		width int
	}
	var rendered []renderedOpt
	for i, opt := range f.Options {
		if len(opt.Keys) == 0 {
			continue
		}
		keyColor := f.BorderColor
		if i == 0 {
			keyColor = f.DefaultKeyColor
		}
		text, w := design.RenderDialogAction(opt.Keys[0].Label, opt.Label, keyColor, theme)
		rendered = append(rendered, renderedOpt{text: text, width: w})
	}

	if len(rendered) == 0 {
		// Sem ações: linha ─ completa
		innerLine := borderStyle.Render(strings.Repeat(design.SymBorderH, innerWidth))
		return bl + innerLine + br
	}

	// Montar a linha de rodapé
	// Formato: ╰─  Ação1  ─  Ação2  ─╯
	var sb strings.Builder
	sb.WriteString(bl)

	// Calcular largura total das ações
	totalActionsWidth := 0
	for _, r := range rendered {
		totalActionsWidth += r.width
	}

	// Separador entre ações: " ─ "
	separator := " " + borderStyle.Render(design.SymBorderH) + " "
	sepWidth := lipgloss.Width(separator)

	// Calcular espaços entre ações
	actionCount := len(rendered)
	if actionCount == 1 {
		// Ação única à direita: preencher até a ação
		fillCount := innerWidth - rendered[0].width - 2 // -2 por causa do "╰─ " e " ─╯"
		if fillCount < 1 {
			fillCount = 1
		}
		sb.WriteString(strings.Repeat(design.SymBorderH, fillCount))
		sb.WriteString(" ")
		sb.WriteString(rendered[0].text)
	} else {
		// Múltiplas ações
		totalSepWidth := (actionCount - 1) * sepWidth
		fillBetween := innerWidth - totalActionsWidth - totalSepWidth - 2 // -2 por causa do "╰─ " e " ─╯"
		if fillBetween < 1 {
			fillBetween = 1
		}

		// Distribuir espaços entre ações
		spacePerGap := fillBetween / (actionCount - 1)
		spaceRemainder := fillBetween % (actionCount - 1)

		for i, r := range rendered {
			if i > 0 {
				gap := spacePerGap
				if spaceRemainder > 0 {
					gap++
					spaceRemainder--
				}
				sb.WriteString(strings.Repeat(design.SymBorderH, gap))
				sb.WriteString(separator)
			}
			sb.WriteString(r.text)
		}
	}

	sb.WriteString(br)
	return sb.String()
}
