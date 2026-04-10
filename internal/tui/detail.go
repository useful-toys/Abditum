package tui

import (
	"strings"

	"github.com/useful-toys/abditum/internal/vault"
)

// revealState tracks reveal status of sensitive fields in the detail panel.
type revealState int

const (
	revealMasked revealState = iota
	revealHint               // first 3 chars + "••"
	revealFull
)

// detailModel tracks the detail panel state.
type detailModel struct {
	secret    *vault.Segredo
	scrollTop int
	reveal    revealState // applies to first sensitive field
	focused   bool
}

// renderDetail renders the detail panel for the currently selected secret.
// width: panel width. height: available lines.
func renderDetail(st styles, d detailModel, width, height int) string {
	if d.secret == nil {
		return renderDetailEmpty(st, width, height)
	}
	s := d.secret

	lines := []string{}

	// Line 1: name + breadcrumb
	name := s.Nome()
	star := ""
	if s.Favorito() {
		star = " " + st.AccentSecondary.Render("★")
	}
	breadcrumb := buildBreadcrumb(s.Pasta())
	nameStr := st.TextPrimary.Bold(true).Render(name) + star
	nameVisLen := len([]rune(name)) + len([]rune(star))
	bcMaxLen := width - nameVisLen - 2
	if bcMaxLen < 0 {
		bcMaxLen = 0
	}
	bcStr := st.TextSecondary.Render(truncateLeft(breadcrumb, bcMaxLen))
	bcVisLen := len([]rune(truncateLeft(breadcrumb, bcMaxLen)))

	gap := max(1, width-nameVisLen-bcVisLen)
	headerLine := nameStr + strings.Repeat(" ", gap) + bcStr
	lines = append(lines, headerLine)

	// Line 2: separator
	// reserve last col for scroll indicator
	lines = append(lines, st.BorderDefault.Render(strings.Repeat("─", width-1))+" ")

	// Fields
	campos := s.Campos()
	firstSensitiveIdx := -1
	for i, c := range campos {
		if c.Tipo() == vault.TipoCampoSensivel {
			firstSensitiveIdx = i
			break
		}
	}

	for i, campo := range campos {
		if i > 0 {
			lines = append(lines, "") // blank between fields
		}
		// Rótulo
		labelStr := st.TextSecondary.Render(campo.Nome())
		lines = append(lines, labelStr)

		// Valor
		var valStr string
		if campo.Tipo() == vault.TipoCampoSensivel {
			if i == firstSensitiveIdx {
				switch d.reveal {
				case revealMasked:
					valStr = st.TextSecondary.Render("••••••••")
				case revealHint:
					full := campo.ValorComoString()
					hint := buildRevealHint(full)
					valStr = st.TextPrimary.Render(hint)
				case revealFull:
					valStr = st.TextPrimary.Render(campo.ValorComoString())
				}
			} else {
				valStr = st.TextSecondary.Render("••••••••")
			}
		} else {
			v := campo.ValorComoString()
			if v == "" {
				valStr = ""
			} else {
				valStr = st.TextPrimary.Render(v)
			}
		}
		lines = append(lines, valStr)
	}

	// Observação
	obs := s.Observacao()
	if obs != "" {
		lines = append(lines, st.BorderDefault.Render(strings.Repeat("╌", width-1))+" ")
		// wrap observation
		for _, obsLine := range wrapText(obs, width-1) {
			lines = append(lines, st.TextPrimary.Render(obsLine))
		}
	}

	// Render with scroll
	return renderScrollable(st, lines, d.scrollTop, height, width)
}

// renderDetailEmpty renders placeholder when no secret is selected.
func renderDetailEmpty(st styles, width, height int) string {
	lines := make([]string, height)
	msg := "Cofre vazio"
	row := height / 2
	if row < 0 {
		row = 0
	}
	for i := range lines {
		if i == row {
			lines[i] = centerText(st.TextSecondary.Italic(true).Render(msg), width, len([]rune(msg)))
		} else {
			lines[i] = ""
		}
	}
	return strings.Join(lines, "\n")
}

// buildBreadcrumb builds "Pasta › Subpasta › ..." for the given folder.
func buildBreadcrumb(p *vault.Pasta) string {
	if p == nil {
		return ""
	}
	var parts []string
	cur := p
	for cur != nil {
		parts = append([]string{cur.Nome()}, parts...)
		cur = cur.Pai()
	}
	return strings.Join(parts, " › ")
}

// buildRevealHint shows the first 3 chars of a value followed by ••.
func buildRevealHint(val string) string {
	runes := []rune(val)
	if len(runes) <= 3 {
		return val
	}
	return string(runes[:3]) + "••"
}

// wrapText wraps text at width columns (simple word-wrap).
func wrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}
	words := strings.Fields(text)
	var lines []string
	var current strings.Builder
	lineLen := 0
	for _, w := range words {
		wLen := len([]rune(w))
		if lineLen > 0 && lineLen+1+wLen > width {
			lines = append(lines, current.String())
			current.Reset()
			lineLen = 0
		}
		if lineLen > 0 {
			current.WriteRune(' ')
			lineLen++
		}
		current.WriteString(w)
		lineLen += wLen
	}
	if current.Len() > 0 {
		lines = append(lines, current.String())
	}
	return lines
}

// renderScrollable renders lines into height rows with scroll indicators on the right.
func renderScrollable(st styles, content []string, scrollTop, height, width int) string {
	totalLines := len(content)

	// clamp scrollTop
	if scrollTop > totalLines-height {
		scrollTop = totalLines - height
	}
	if scrollTop < 0 {
		scrollTop = 0
	}

	result := make([]string, height)
	for row := 0; row < height; row++ {
		lineIdx := scrollTop + row
		var lineContent string
		if lineIdx < totalLines {
			lineContent = content[lineIdx]
		}

		// Pad to width-1 (reserve last col for scroll)
		visLen := visibleLen(lineContent)
		if visLen < width-1 {
			lineContent += strings.Repeat(" ", width-1-visLen)
		}

		// Scroll indicator
		scrollIndicator := " "
		if totalLines > height {
			if row == 0 && scrollTop > 0 {
				scrollIndicator = st.TextSecondary.Render("↑")
			} else if row == height-1 && scrollTop+height < totalLines {
				scrollIndicator = st.TextSecondary.Render("↓")
			} else {
				thumb := thumbPosition(scrollTop, totalLines-height, height)
				if row == thumb {
					scrollIndicator = st.TextSecondary.Render("■")
				}
			}
		}

		result[row] = lineContent + scrollIndicator
	}

	return strings.Join(result, "\n")
}
