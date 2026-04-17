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

	// Calcular largura baseada no corpo + padding
	bodyWidth := f.calculateBodyWidth(body, theme)
	// innerWidth = min(maxWidth - 2, max(bodyWidth, 20))
	// Ou seja: usar a largura calculada, mas nunca menor que 20 nem maior que maxWidth - 2
	innerWidth := bodyWidth
	if innerWidth < 20 {
		innerWidth = 20
	}
	maxAvailable := maxWidth - 2
	if innerWidth > maxAvailable {
		innerWidth = maxAvailable
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

	// Conteúdo com background — preencher até innerWidth para alinhar a borda direita
	contentVisualWidth := lipgloss.Width(content)
	availableContent := innerWidth - design.DialogPaddingH*2
	fill := availableContent - contentVisualWidth
	if fill < 0 {
		fill = 0
	}
	lineContent := bgStyle.Render(paddingH + content + strings.Repeat(" ", fill) + paddingH)

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

// calculateBodyWidth calcula a largura baseada no conteúdo do corpo.
func (f DialogFrame) calculateBodyWidth(body string, theme *design.Theme) int {
	paddingH := 2 * design.DialogPaddingH

	// Largura do título na borda superior
	titleWidth := lipgloss.Width(f.Title)
	if f.Symbol != "" {
		titleWidth += lipgloss.Width(f.Symbol) + 2
	}
	titleWidth += 8 // " título " + cantos

	// Largura das ações
	actionWidth := 3
	for _, opt := range f.Options {
		if len(opt.Keys) == 0 {
			continue
		}
		_, keyWidth := design.RenderDialogAction(opt.Keys[0].Label, opt.Label, f.BorderColor, theme)
		actionWidth += keyWidth + 4 + 3
	}

	// Largura do corpo (cada linha)
	bodyLines := strings.Split(body, "\n")
	maxBodyWidth := 0
	for _, line := range bodyLines {
		w := lipgloss.Width(line) + paddingH
		if w > maxBodyWidth {
			maxBodyWidth = w
		}
	}

	// Usar a maior largura needed
	width := titleWidth
	if actionWidth > width {
		width = actionWidth
	}
	if maxBodyWidth > width {
		width = maxBodyWidth
	}

	return width // innerWidth (sem bordas)
}

// renderBottomBorder gera a linha de rodapé conforme a spec:
//
// Cada ação é envolvida por espaços: " ação ".
// O preenchimento ─ fica entre os espaços externos de cada bloco de ação.
//
// Estrutura da linha completa (incluindo cantos):
//
//	╰─[─…─] ação1 [─…─] ação2 [─…─]─╯
//
// Posicionamento:
//
//	1 ação  → ╰──────────────────── ação ─╯  (alinhada à direita)
//	2 ações → ╰─ ação1 ──────────── ação2 ─╯
//	3 ações → ╰─ ação1 ─── ação2 ── ação3 ─╯
func (f DialogFrame) renderBottomBorder(innerWidth int, borderStyle lipgloss.Style, theme *design.Theme) string {
	// Cantos: ╰ e ╯ (1 coluna cada). O traço inicial e final estão no fill.
	cornerBL := borderStyle.Render(design.SymCornerBL)
	cornerBR := borderStyle.Render(design.SymCornerBR)

	dash := func(n int) string {
		if n <= 0 {
			return ""
		}
		return borderStyle.Render(strings.Repeat(design.SymBorderH, n))
	}

	// Renderizar as ações (RenderDialogAction retorna "key label" sem espaços externos)
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
		return cornerBL + dash(innerWidth) + cornerBR
	}

	// Anatomia de cada ação no rodapé: " ação " (1 espaço antes + 1 espaço depois)
	// O preenchimento ─ vai entre os blocos de ação (já incluindo os espaços internos de cada ação).
	//
	// Para N ações, innerWidth é ocupado assim:
	//   fill[0] + " " + ação[0] + " " + fill[1] + " " + ação[1] + " " + … + fill[N] = innerWidth
	//
	// fill[0] >= 1 (pelo menos "─" após o canto esquerdo)
	// fill[N] >= 1 (pelo menos "─" antes do canto direito)
	// fill[i > 0 e < N] >= 1 (pelo menos "─" entre ações)
	//
	// Total de espaços fixos: 2 * len(rendered) (1 antes e 1 depois de cada ação)
	// Total mínimo de fill: len(rendered) + 1 (um fill entre cada par + esquerda + direita)
	// Espaço disponível para fill: innerWidth - 2*len(ações) - soma(widths)

	n := len(rendered)
	totalActionWidth := 0
	for _, r := range rendered {
		totalActionWidth += r.width
	}
	// Espaços: 2 por ação (antes e depois)
	fixedSpaces := 2 * n
	// Fills: n+1 slots (antes da 1ª ação, entre ações, depois da última)
	totalFill := innerWidth - totalActionWidth - fixedSpaces
	if totalFill < n+1 {
		totalFill = n + 1 // mínimo 1 traço por slot
	}

	// Distribuir o fill: slot 0 recebe o restante para alinhar à direita quando n==1,
	// ou para 2+ ações: slot 0 e slots intermediários recebem 1 traço mínimo,
	// o excedente vai para o último slot antes do canto direito.
	//
	// Spec para 2 ações: ação1 à esquerda, ação2 à direita → fill grande no meio.
	// Spec para 1 ação: ação à direita → fill grande à esquerda.

	fills := make([]int, n+1)
	// Mínimo 1 traço em cada slot
	for i := range fills {
		fills[i] = 1
	}
	remaining := totalFill - (n + 1)

	if n == 1 {
		// 1 ação à direita: excedente todo no slot 0 (esquerda)
		fills[0] += remaining
	} else {
		// 2+ ações: excedente no slot n-1 (entre última e penúltima — empurra a última à direita)
		// Isso coloca ação1 à esquerda e ação2 à direita.
		fills[n-1] += remaining
	}

	var sb strings.Builder
	sb.WriteString(cornerBL)
	for i, r := range rendered {
		sb.WriteString(dash(fills[i]))
		sb.WriteString(" ")
		sb.WriteString(r.text)
		sb.WriteString(" ")
	}
	sb.WriteString(dash(fills[n]))
	sb.WriteString(cornerBR)
	return sb.String()
}
