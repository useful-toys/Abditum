package tui

import (
	"strings"
	"unicode/utf8"
)

// dialogSeverity classifies a dialog's visual treatment.
type dialogSeverity int

const (
	severityNeutral dialogSeverity = iota
	severityInfo
	severityAlert
	severityDestructive
	severityError
)

// dialogAction describes a button in a dialog.
type dialogAction struct {
	Key     string // display key (e.g. "Enter", "S", "Esc")
	Label   string // display label (e.g. "Salvar", "Cancelar")
	IsESC   bool   // marks the cancel/back action
	IsEnter bool   // marks the default action
}

// dialog represents a decision dialog (confirmation or acknowledgment).
type dialog struct {
	title    string
	body     string
	severity dialogSeverity
	actions  []dialogAction // 1–3 actions; last is typically Esc/cancel
}

// renderDialog renders a dialog centered in (width × height) space.
// Returns a string of `height` lines.
func renderDialog(st styles, d dialog, width, height int) string {
	sev := d.severity
	borderColor := severityBorderColor(st, sev)
	defaultKeyColor := severityDefaultKeyColor(st, sev)
	symbol := severitySymbol(sev)

	// --- Compute dialog width ---
	// Title width: ⚠  Title (with symbol if non-neutral)
	titleFull := ""
	if symbol != "" {
		titleFull = symbol + "  " + d.title
	} else {
		titleFull = d.title
	}

	// -- Compute actions bar width
	actionsBarContent := buildActionsBar(d.actions, st, defaultKeyColor, borderColor)
	actionsBarVisLen := visibleLen(actionsBarContent)

	// Preferred width: max of title width + 4, content width + 4, actions bar + 4
	bodyLines := wrapText(d.body, width*80/100-4)
	contentWidth := 0
	for _, l := range bodyLines {
		w := utf8.RuneCountInString(l)
		if w > contentWidth {
			contentWidth = w
		}
	}

	preferred := max(max(
		utf8.RuneCountInString(titleFull)+4,
		contentWidth+4),
		actionsBarVisLen+4)

	maxDialogW := width * 95 / 100
	if maxDialogW < 20 {
		maxDialogW = 20
	}
	dialogW := min(preferred, maxDialogW)
	if dialogW < 20 {
		dialogW = 20
	}

	// Rewrap body at final dialog inner width
	innerW := dialogW - 2                  // excluding border chars
	bodyLines = wrapText(d.body, innerW-2) // 2 inner padding

	// -- Build dialog lines --
	var dLines []string

	// Top border: ╭── ⚠ Title ────────╮
	topBorderLine := buildTopBorder(titleFull, dialogW, borderColor, st)
	dLines = append(dLines, topBorderLine)

	// padding top
	dLines = append(dLines, borderColor.Render("│")+strings.Repeat(" ", innerW)+borderColor.Render("│"))

	// Body lines
	for _, bl := range bodyLines {
		visLen := utf8.RuneCountInString(bl)
		pad := innerW - 2 - visLen
		if pad < 0 {
			pad = 0
		}
		dLines = append(dLines,
			borderColor.Render("│")+" "+
				st.TextPrimary.Render(bl)+
				strings.Repeat(" ", pad+1)+
				borderColor.Render("│"))
	}

	// padding bottom
	dLines = append(dLines, borderColor.Render("│")+strings.Repeat(" ", innerW)+borderColor.Render("│"))

	// Bottom border with actions: ╰── Enter Salvar ─── Esc Cancelar ──╯
	botBorderLine := buildBottomBorder(d.actions, dialogW, borderColor, defaultKeyColor, st)
	dLines = append(dLines, botBorderLine)

	// Center dialog vertically
	dialogH := len(dLines)
	topPad := (height - dialogH) / 2
	if topPad < 0 {
		topPad = 0
	}
	leftPad := (width - dialogW) / 2
	if leftPad < 0 {
		leftPad = 0
	}
	padStr := strings.Repeat(" ", leftPad)

	result := make([]string, height)
	for row := 0; row < height; row++ {
		dlgRow := row - topPad
		if dlgRow < 0 || dlgRow >= len(dLines) {
			result[row] = ""
		} else {
			result[row] = padStr + dLines[dlgRow]
		}
	}
	return strings.Join(result, "\n")
}

func buildTopBorder(title string, dialogW int, borderColor, st styles) string {
	// ╭── <symbol> Title ────────────────╮
	inner := dialogW - 2 // between ╭ and ╮
	prefix := "── " + title + " "
	prefixLen := utf8.RuneCountInString(prefix)
	remaining := inner - prefixLen
	if remaining < 0 {
		remaining = 0
	}
	return borderColor.Render("╭" + prefix + strings.Repeat("─", remaining) + "╮")
}

func buildBottomBorder(actions []dialogAction, dialogW int, borderColor, defaultKeyColor styles, st styles) string {
	// ╰── Enter Salvar ──── Esc Cancelar ──╯
	inner := dialogW - 2

	if len(actions) == 0 {
		return borderColor.Render("╰" + strings.Repeat("─", inner) + "╯")
	}

	// Build action segments
	parts := make([]struct{ str, vis string }, len(actions))
	for i, a := range actions {
		var keyStyle styles
		if a.IsEnter {
			keyStyle = defaultKeyColor
		} else {
			keyStyle = borderColor
		}
		parts[i].str = keyStyle.Bold(a.IsEnter).Render(a.Key) + " " + keyStyle.Bold(a.IsEnter).Render(a.Label)
		parts[i].vis = a.Key + " " + a.Label
	}

	switch len(actions) {
	case 1:
		// right-aligned
		dash := inner - utf8.RuneCountInString(parts[0].vis) - 3
		if dash < 1 {
			dash = 1
		}
		return borderColor.Render("╰") +
			borderColor.Render(strings.Repeat("─", dash)+(" ")) +
			parts[0].str +
			borderColor.Render(" ──╯")
	case 2:
		// default left, cancel right
		middle := inner - utf8.RuneCountInString(parts[0].vis) - utf8.RuneCountInString(parts[1].vis) - 6
		if middle < 1 {
			middle = 1
		}
		return borderColor.Render("╰── ") +
			parts[0].str +
			borderColor.Render(strings.Repeat("─", middle)+" ") +
			parts[1].str +
			borderColor.Render(" ──╯")
	default:
		// 3: default left, middle, cancel right
		// Simple equal-space distribution
		totalVis := 0
		for _, p := range parts {
			totalVis += utf8.RuneCountInString(p.vis)
		}
		gaps := inner - totalVis - 4 // 4 for ── at start and ──╯ at end
		gapSize := gaps / (len(parts) - 1)
		if gapSize < 1 {
			gapSize = 1
		}
		var sb strings.Builder
		sb.WriteString(borderColor.Render("╰── "))
		for i, p := range parts {
			sb.WriteString(p.str)
			if i < len(parts)-1 {
				sb.WriteString(borderColor.Render(" " + strings.Repeat("─", gapSize) + " "))
			}
		}
		sb.WriteString(borderColor.Render(" ──╯"))
		return sb.String()
	}
}

// buildActionsBar returns a preview string for width calculation.
func buildActionsBar(actions []dialogAction, st, defaultKeyColor, borderColor styles) string {
	var parts []string
	for _, a := range actions {
		parts = append(parts, a.Key+" "+a.Label)
	}
	return "── " + strings.Join(parts, " ─── ") + " ──"
}

func severityBorderColor(st styles, sev dialogSeverity) styles {
	_ = st
	switch sev {
	case severityDestructive, severityAlert:
		return st.SemanticWarning
	case severityError:
		return st.SemanticError
	case severityInfo:
		return st.SemanticInfo
	default:
		return st.BorderFocused
	}
}

func severityDefaultKeyColor(st styles, sev dialogSeverity) styles {
	switch sev {
	case severityDestructive:
		return st.SemanticError
	default:
		return st.AccentPrimary
	}
}

func severitySymbol(sev dialogSeverity) string {
	switch sev {
	case severityDestructive, severityAlert:
		return "⚠"
	case severityError:
		return "✕"
	case severityInfo:
		return "ℹ"
	default:
		return ""
	}
}

// overlayDialog overlays a dialog on top of a background rendering.
// background: full-screen string (height lines)
// dlgLines: dialog string (height lines, same dimensions as background)
// Returns merged output.
func overlayDialog(background, dlgStr string, width, height int) string {
	bgLines := strings.Split(background, "\n")
	dlgLines := strings.Split(dlgStr, "\n")

	result := make([]string, height)
	for i := 0; i < height; i++ {
		bg := ""
		if i < len(bgLines) {
			bg = bgLines[i]
		}
		dl := ""
		if i < len(dlgLines) {
			dl = dlgLines[i]
		}

		// If dialog line is non-empty, use it; else use background
		if strings.TrimSpace(stripANSI(dl)) != "" {
			result[i] = dl
		} else {
			result[i] = bg
		}
	}
	return strings.Join(result, "\n")
}
