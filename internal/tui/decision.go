package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ─────────────────────────────────────────────────────────────────────────────
// Domain types
// ─────────────────────────────────────────────────────────────────────────────

// Severity controls the visual treatment of a DecisionDialog (border, title,
// action token colors). Maps to the severity column in the design-system table.
type Severity int

const (
	SeverityNeutral     Severity = iota // border.focused token, no symbol
	SeverityInformative                 // semantic.info token, ℹ symbol
	SeverityAlert                       // semantic.warning token, ⚠ symbol
	SeverityError                       // semantic.error token, ✕ symbol
	SeverityDestructive                 // semantic.warning border, semantic.error default key
)

// Intention controls the action bar layout.
type Intention int

const (
	IntentionAcknowledge Intention = iota // single action: Enter OK (right-aligned)
	IntentionConfirm                      // two or three actions: default + [middle] + Esc
)

// DecisionAction is a single selectable action in the action bar.
type DecisionAction struct {
	Key     string // display key string, e.g. "Enter", "N", "Esc"
	Label   string // display label, e.g. "Excluir", "Cancelar"
	Cmd     tea.Cmd
	Default bool // true → bold + severity default-key color
	Cancel  bool // true → triggered by Esc key at any focus position
}

// ─────────────────────────────────────────────────────────────────────────────
// DecisionDialog
// ─────────────────────────────────────────────────────────────────────────────

// DecisionDialog is a decision overlay (D-02) implementing modalView.
// It renders the canonical wireframe from tui-specification-novo.md §Diálogos de Decisão:
//   - Rounded border + title with optional severity symbol
//   - Body message in text.primary
//   - Action bar at the bottom of the box (not outside)
//
// Keyboard: Enter triggers the default action; Esc triggers the cancel action;
// explicit key bindings (e.g. "M", "A") trigger their matched action directly.
// Shortcuts() returns nil — command bar shows only F1 Ajuda while dialog is open.
type DecisionDialog struct {
	title     string
	body      string
	severity  Severity
	intention Intention
	actions   []DecisionAction // ordered: default first, cancel last
	width     int
}

// Compile-time assertion: DecisionDialog must satisfy modalView.
var _ modalView = &DecisionDialog{}

// ─────────────────────────────────────────────────────────────────────────────
// Constructor
// ─────────────────────────────────────────────────────────────────────────────

// NewDecisionDialog creates a new DecisionDialog.
// actions must be ordered: default first, optional intermediate, cancel last.
// At least one action is required. For IntentionAcknowledge pass a single action
// with Default: true (the factory helpers below enforce this).
func NewDecisionDialog(
	severity Severity,
	intention Intention,
	title, body string,
	actions []DecisionAction,
) *DecisionDialog {
	return &DecisionDialog{
		title:     title,
		body:      body,
		severity:  severity,
		intention: intention,
		actions:   actions,
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Factory helpers (match the 5 × 2 = 10 combinations from the spec)
// ─────────────────────────────────────────────────────────────────────────────

// Acknowledge creates an IntentionAcknowledge dialog (single Enter OK action).
func Acknowledge(severity Severity, title, body string, onOK tea.Cmd) tea.Cmd {
	d := NewDecisionDialog(severity, IntentionAcknowledge, title, body, []DecisionAction{
		{Key: "Enter", Label: "OK", Cmd: onOK, Default: true},
	})
	return func() tea.Msg { return pushModalMsg{modal: d} }
}

// Decision creates an IntentionConfirm dialog with a default action and cancel.
// middleActions may be empty (binary confirm) or contain one additional action.
func Decision(severity Severity, title, body string, defaultAction DecisionAction, middleActions []DecisionAction, cancelAction DecisionAction) tea.Cmd {
	actions := []DecisionAction{defaultAction}
	actions = append(actions, middleActions...)
	cancelAction.Cancel = true
	actions = append(actions, cancelAction)
	d := NewDecisionDialog(severity, IntentionConfirm, title, body, actions)
	return func() tea.Msg { return pushModalMsg{modal: d} }
}

// ─────────────────────────────────────────────────────────────────────────────
// modalView implementation
// ─────────────────────────────────────────────────────────────────────────────

// Shortcuts returns nil — while a DecisionDialog is active the command bar shows
// only the global F1 Ajuda shortcut (injected by rootModel via ActionManager).
func (d *DecisionDialog) Shortcuts() []Shortcut { return nil }

// Update handles keyboard input. Returns only tea.Cmd (modalView contract).
func (d *DecisionDialog) Update(msg tea.Msg) tea.Cmd {
	kp, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return nil
	}

	key := kp.String()

	switch key {
	case "enter":
		return d.triggerFocused()

	case "esc":
		return d.triggerCancel()
	}

	// Check if any action has a matching explicit key binding.
	for _, a := range d.actions {
		if a.Cancel {
			continue // Esc already handled above
		}
		if strings.EqualFold(key, strings.ToLower(a.Key)) && a.Key != "Enter" {
			return d.triggerAction(a)
		}
	}

	return nil
}

// View renders the dialog box. rootModel positions it via lipgloss.Place.
// The box is built line-by-line to avoid lipgloss v2's missing BorderTitle method.
//
// Structure:
//
//	╭── ⚠  Excluir segredo ─────────────────────────────────────────╮
//	│                                                                 │
//	│  Gmail será excluído permanentemente. Esta ação não pode ser    │
//	│  desfeita.                                                      │
//	│                                                                 │
//	╰── Enter Excluir ─────────────────────────────────── Esc Cancelar ──╯
func (d *DecisionDialog) View(maxWidth, maxHeight int, theme *Theme) string {
	if maxWidth == 0 {
		panic(fmt.Sprintf("DecisionDialog.View() called without maxWidth: maxWidth=%d", maxWidth))
	}
	// Store width for use by boxWidth()
	d.width = maxWidth
	borderColor := d.borderColor(theme)
	boxW := d.boxWidth()

	borderSt := lipgloss.NewStyle().Foreground(lipgloss.Color(borderColor))
	titleSt := lipgloss.NewStyle().Foreground(lipgloss.Color(borderColor)).Bold(true)
	bodySt := lipgloss.NewStyle().Foreground(lipgloss.Color("#a9b1d6"))

	// ── Top border ───────────────────────────────────────────────────────────
	// "╭── {symbol}  {title} ────...────╮"
	// leftAnchor = "╭── " (4 runes), rightAnchor = " ──╮" (4 runes)
	symbol := d.symbol()
	titleText := d.title
	if symbol != "" {
		titleText = symbol + "  " + titleText
	}
	titleW := lipgloss.Width(titleText) // measure plain (no ANSI yet)
	const leftAnchorW = 4               // "╭── "
	const rightAnchorW = 4              // " ──╮"
	fillW := boxW - leftAnchorW - titleW - rightAnchorW
	if fillW < 1 {
		fillW = 1
	}
	topLine := borderSt.Render("╭──") + " " +
		titleSt.Render(titleText) + " " +
		borderSt.Render(strings.Repeat("─", fillW)+"──╮")

	// ── Empty padding line ───────────────────────────────────────────────────
	// innerW = columns between the two │ chars (exclusive)
	innerW := boxW - 2 // boxW includes the 2 border columns
	emptyPad := borderSt.Render("│") + strings.Repeat(" ", innerW) + borderSt.Render("│")

	// ── Body lines (word-wrapped) ────────────────────────────────────────────
	// innerW = boxW-2; padding each side = 2 chars; body content = innerW-4
	maxBodyW := innerW - 4
	if maxBodyW < 1 {
		maxBodyW = 1
	}
	wrappedBody := wrapBody(d.body, maxBodyW)

	var bodyLines []string
	for _, line := range wrappedBody {
		// Single render: apply both color and exact width together (mirrors help.go pattern).
		content := bodySt.Width(maxBodyW).Render(line)
		bodyLines = append(bodyLines, borderSt.Render("│")+"  "+content+"  "+borderSt.Render("│"))
	}

	// ── Action bar (bottom border) ───────────────────────────────────────────
	actionBar := d.renderActionBar(boxW, theme)

	// ── Assembly ─────────────────────────────────────────────────────────────
	lines := []string{topLine, emptyPad}
	lines = append(lines, bodyLines...)
	lines = append(lines, emptyPad, actionBar)

	return strings.Join(lines, "\n")
}

// wrapBody word-wraps text to maxWidth runes per line.
// It respects existing "\n" in the input and splits on spaces.
// Uses only stdlib — no external dependency.
func wrapBody(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		maxWidth = 1
	}
	var result []string

	// First split on explicit newlines.
	paragraphs := strings.Split(text, "\n")
	for _, para := range paragraphs {
		words := strings.Fields(para)
		if len(words) == 0 {
			result = append(result, "")
			continue
		}
		var line strings.Builder
		lineW := 0
		for i, word := range words {
			wordW := len([]rune(word))
			if i == 0 {
				line.WriteString(word)
				lineW = wordW
			} else if lineW+1+wordW <= maxWidth {
				line.WriteByte(' ')
				line.WriteString(word)
				lineW += 1 + wordW
			} else {
				result = append(result, line.String())
				line.Reset()
				line.WriteString(word)
				lineW = wordW
			}
		}
		if line.Len() > 0 {
			result = append(result, line.String())
		}
	}
	return result
}

// ─────────────────────────────────────────────────────────────────────────────
// Internal helpers
// ─────────────────────────────────────────────────────────────────────────────

// borderColor returns the hex color for the box border and title.
func (d *DecisionDialog) borderColor(theme *Theme) string {
	switch d.severity {
	case SeverityInformative:
		return "#7dcfff"
	case SeverityAlert:
		return "#e0af68"
	case SeverityError:
		return "#f7768e"
	case SeverityDestructive:
		return "#e0af68" // destructive: border = semantic.warning
	default: // SeverityNeutral
		return "#414868"
	}
}

// defaultKeyColor returns the color for the default action key.
func (d *DecisionDialog) defaultKeyColor(theme *Theme) string {
	switch d.severity {
	case SeverityDestructive:
		return "#f7768e" // destructive: default key = semantic.error
	default:
		return "#7aa2f7"
	}
}

// symbol returns the prefix symbol for the title, or "" for Neutral.
func (d *DecisionDialog) symbol() string {
	switch d.severity {
	case SeverityInformative:
		return SymInfo
	case SeverityAlert, SeverityDestructive:
		return SymWarning
	case SeverityError:
		return SymError
	default:
		return ""
	}
}

// boxWidth returns the dialog width. Sized to fit content, capped at 80% of
// terminal width (or 80 columns absolute max). Never smaller than 40 columns.
func (d *DecisionDialog) boxWidth() int {
	maxW := 80 // absolute max
	if d.width > 0 {
		pct := d.width * 80 / 100
		if pct < maxW {
			maxW = pct
		}
	}
	if maxW < 40 {
		maxW = 40
	}

	// Minimum width needed to render the top border without truncation:
	// "╭── {symbol}  {title} ──╮" = leftAnchorW + titleW + rightAnchorW + 1 fill
	symbol := d.symbol()
	titleText := d.title
	if symbol != "" {
		titleText = symbol + "  " + titleText
	}
	titleW := lipgloss.Width(titleText)
	const leftAnchorW = 4  // "╭── "
	const rightAnchorW = 4 // " ──╮"
	minForTitle := leftAnchorW + titleW + rightAnchorW + 1

	// Minimum width needed to render the action bar without truncation.
	actionMinW := d.actionBarMinWidth()

	// Minimum width needed for a readable body (at least 20 chars of content).
	const minBodyInner = 20
	minForBody := minBodyInner + 4 + 2 // padding (4) + borders (2)

	minW := minForTitle
	if actionMinW > minW {
		minW = actionMinW
	}
	if minForBody > minW {
		minW = minForBody
	}
	if minW < 40 {
		minW = 40
	}

	// Use the minimum content width, but never exceed maxW.
	if minW > maxW {
		return minW
	}
	return minW
}

// actionBarMinWidth returns the minimum box width needed to render the action
// bar with at least 1 fill dash between left and right portions.
func (d *DecisionDialog) actionBarMinWidth() int {
	// Acknowledgement: "╰── seg ──╯" = 2(corners) + 2(lead) + 1(sp) + actionW + 1(sp) + 2(trail)
	// Confirmation:    "╰── left ──[fill]── right ──╯"
	// We compute plain-text widths of left + right + 2 corners + 1 fill min.

	var cancelW int
	var leftW int

	leftW += 4 // "╰── "
	for i, a := range d.actions {
		plain := a.Key + " " + a.Label
		if a.Cancel {
			cancelW = 1 + len([]rune(plain)) + 4 // " plain ──╯"
		} else {
			leftW += len([]rune(plain))
			if i < len(d.actions)-1 && !d.actions[i+1].Cancel {
				leftW += 4 // " ── "
			} else {
				leftW += 1 // trailing space
			}
		}
	}

	return leftW + 1 + cancelW // +1 for at least one fill dash
}

// renderActionBar builds the bottom border line with embedded action tokens.
// Canonical layout (from spec):
//   - Acknowledge:        ╰──────────────────────── Enter OK ──╯
//   - Confirm 2-action:   ╰── Enter Excluir ─────── Esc Cancelar ──╯
//   - Confirm 3-action:   ╰── Enter Salvar ── N Descartar ── Esc Voltar ──╯
func (d *DecisionDialog) renderActionBar(boxW int, theme *Theme) string {
	borderColor := d.borderColor(theme)
	defaultKeyColor := d.defaultKeyColor(theme)

	borderFg := lipgloss.Color(borderColor)
	defaultFg := lipgloss.Color(defaultKeyColor)
	otherFg := lipgloss.Color(borderColor)

	borderSt := lipgloss.NewStyle().Foreground(borderFg)
	// dashes renders N dashes as a single ANSI span (avoids N×span overhead).
	dashes := func(n int) string { return borderSt.Render(strings.Repeat("─", n)) }

	// Build styled action segments.
	// Each segment: " Key Label " wrapped in "── ... ──"
	type segment struct {
		text   string // plain for width calculation
		styled string // ANSI-styled for rendering
	}

	renderToken := func(key, label string, isDefault bool) segment {
		var kg, lg lipgloss.Style
		if isDefault {
			kg = lipgloss.NewStyle().Foreground(defaultFg).Bold(true)
			lg = lipgloss.NewStyle().Foreground(defaultFg).Bold(true)
		} else {
			kg = lipgloss.NewStyle().Foreground(otherFg)
			lg = lipgloss.NewStyle().Foreground(otherFg)
		}
		plain := key + " " + label
		styled := kg.Render(key) + " " + lg.Render(label)
		return segment{text: plain, styled: styled}
	}

	var segs []segment
	var cancelSeg *segment

	for _, a := range d.actions {
		isDefault := a.Default
		s := renderToken(a.Key, a.Label, isDefault)
		if a.Cancel {
			cancelSeg = &s
		} else {
			segs = append(segs, s)
		}
	}

	// Acknowledgement: only one action, right-aligned.
	if d.intention == IntentionAcknowledge && len(segs) == 1 {
		// ╰────── Enter OK ──╯
		actionPart := " " + segs[0].styled + " " + dashes(2)
		actionPlain := " " + segs[0].text + " ──"

		// total width of actionPart (printable)
		actionW := len([]rune(actionPlain))
		// remaining width for leading dashes: boxW - 2 (corners) - actionW
		leadW := boxW - 2 - actionW
		if leadW < 2 {
			leadW = 2
		}
		lead := borderSt.Render("╰") + dashes(leadW)
		tail := borderSt.Render("╯")
		return lead + actionPart + tail
	}

	// Confirmation: left-aligned default, right-aligned cancel.
	// ╰── Enter Excluir ─────── Esc Cancelar ──╯
	// ╰── Enter Salvar ── N Descartar ── Esc Voltar ──╯

	// Build left portion: "╰── seg0 ──"  [── seg1 ──] ...
	var leftPlain, leftStyled strings.Builder
	leftPlain.WriteString("╰── ")
	leftStyled.WriteString(borderSt.Render("╰") + dashes(2) + " ")

	for i, s := range segs {
		leftPlain.WriteString(s.text)
		leftStyled.WriteString(s.styled)
		if i < len(segs)-1 {
			leftPlain.WriteString(" ── ")
			leftStyled.WriteString(" " + dashes(2) + " ")
		} else {
			// trailing space before the fill dashes
			leftPlain.WriteString(" ")
			leftStyled.WriteString(" ")
		}
	}

	// Build right portion: " Esc Cancelar ──╯"
	var rightPlain, rightStyled string
	if cancelSeg != nil {
		rightPlain = " " + cancelSeg.text + " ──╯"
		rightStyled = " " + cancelSeg.styled + " " + dashes(2) + borderSt.Render("╯")
	} else {
		rightPlain = " ──╯"
		rightStyled = " " + dashes(2) + borderSt.Render("╯")
	}

	leftPW := len([]rune(leftPlain.String()))
	rightPW := len([]rune(rightPlain))
	fillW := boxW - leftPW - rightPW
	if fillW < 1 {
		fillW = 1
	}
	fill := dashes(fillW)

	return leftStyled.String() + fill + rightStyled
}

// triggerFocused executes the action with Default:true and pops the modal.
// This is always triggered by Enter, regardless of any navigation state.
func (d *DecisionDialog) triggerFocused() tea.Cmd {
	for _, a := range d.actions {
		if a.Default {
			return d.triggerAction(a)
		}
	}
	// No default action defined — just pop.
	return func() tea.Msg { return popModalMsg{} }
}

// triggerAction executes a single action's Cmd (if any) and pops the modal.
func (d *DecisionDialog) triggerAction(a DecisionAction) tea.Cmd {
	pop := func() tea.Msg { return popModalMsg{} }
	if a.Cmd != nil {
		return tea.Batch(pop, a.Cmd)
	}
	return pop
}

// triggerCancel executes the cancel action and pops the modal.
func (d *DecisionDialog) triggerCancel() tea.Cmd {
	for _, a := range d.actions {
		if a.Cancel {
			pop := func() tea.Msg { return popModalMsg{} }
			if a.Cmd != nil {
				return tea.Batch(pop, a.Cmd)
			}
			return pop
		}
	}
	// No cancel action defined (e.g. Acknowledge) — just pop.
	return func() tea.Msg { return popModalMsg{} }
}
