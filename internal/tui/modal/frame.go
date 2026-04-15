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

	innerWidth := maxWidth - 2 // subtrair as duas bordas verticais

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

// renderBottomBorder gera a linha de rodapé: ╰─ [ação1] ── [ação2] ── [ação3] ─╯
// Posicionamento:
//
//	1 ação  → alinhada à direita
//	2 ações → 1ª à esquerda, 2ª à direita
//	3 ações → 1ª à esquerda, 2ª ao centro, 3ª à direita
func (f DialogFrame) renderBottomBorder(innerWidth int, borderStyle lipgloss.Style, theme *design.Theme) string {
	bl := borderStyle.Render(design.SymCornerBL)
	br := borderStyle.Render(design.SymCornerBR)
	dash := borderStyle.Render(design.SymBorderH)

	// Renderizar as ações
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
		// Adicionar espaços " " em torno da ação (1 espaço antes e depois)
		rendered = append(rendered, renderedOpt{text: " " + text + " ", width: w + 2})
	}

	if len(rendered) == 0 {
		// Sem ações: linha ─ completa
		return bl + borderStyle.Render(strings.Repeat(design.SymBorderH, innerWidth)) + br
	}

	// Montar a linha de rodapé conforme o número de ações
	line := make([]byte, 0, innerWidth*4)
	writeDashes := func(count int) {
		for i := 0; i < count; i++ {
			line = append(line, []byte(dash)...)
		}
	}
	writeAction := func(r renderedOpt) {
		line = append(line, []byte(r.text)...)
	}

	switch len(rendered) {
	case 1:
		// Ação única à direita
		totalActionWidth := rendered[0].width
		fill := innerWidth - totalActionWidth
		if fill < 0 {
			fill = 0
		}
		writeDashes(fill)
		writeAction(rendered[0])
	case 2:
		// 1ª à esquerda, 2ª à direita
		gap := innerWidth - rendered[0].width - rendered[1].width
		if gap < 1 {
			gap = 1
		}
		writeAction(rendered[0])
		writeDashes(gap)
		writeAction(rendered[1])
	case 3:
		// 1ª à esquerda, 2ª ao centro, 3ª à direita
		remaining := innerWidth - rendered[0].width - rendered[1].width - rendered[2].width
		if remaining < 2 {
			remaining = 2
		}
		gapLeft := remaining / 2
		gapRight := remaining - gapLeft
		writeAction(rendered[0])
		writeDashes(gapLeft)
		writeAction(rendered[1])
		writeDashes(gapRight)
		writeAction(rendered[2])
	default:
		// Mais de 3: renderizar apenas as 3 primeiras (DS: máximo 3)
		remaining := innerWidth - rendered[0].width - rendered[1].width - rendered[2].width
		if remaining < 2 {
			remaining = 2
		}
		gapLeft := remaining / 2
		gapRight := remaining - gapLeft
		writeAction(rendered[0])
		writeDashes(gapLeft)
		writeAction(rendered[1])
		writeDashes(gapRight)
		writeAction(rendered[2])
	}

	return bl + string(line) + br
}
