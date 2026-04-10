package tui

import (
	"strings"
	"time"
	"unicode/utf8"
)

// msgType represents the type of a bar message.
type msgType int

const (
	msgNone msgType = iota
	msgSuccess
	msgInfo
	msgWarning
	msgError
	msgBusy
	msgHint
)

// barMessage holds the current message bar state.
type barMessage struct {
	kind    msgType
	text    string
	expires time.Time // zero = never expires
}

var noMessage = barMessage{}

// spinnerFrames for the busy spinner.
var spinnerFrames = []string{"◐", "◓", "◑", "◒"}

// renderMsgBar renders the 1-line message bar.
// width: terminal width.
func renderMsgBar(st styles, width int, msg barMessage, spinFrame int) string {
	if msg.kind == msgNone {
		return st.BorderDefault.Render(strings.Repeat("─", width))
	}

	var sym, text string
	var msgStyle = st.TextPrimary

	switch msg.kind {
	case msgSuccess:
		sym = "✓"
		msgStyle = st.MsgSuccess
	case msgInfo:
		sym = "ℹ"
		msgStyle = st.MsgInfo
	case msgWarning:
		sym = "⚠"
		msgStyle = st.MsgWarning
	case msgError:
		sym = "✕"
		msgStyle = st.MsgError
	case msgBusy:
		sym = spinnerFrames[spinFrame%len(spinnerFrames)]
		msgStyle = st.MsgSpinner
	case msgHint:
		sym = "•"
		msgStyle = st.MsgHint
	}

	// Format: "── <sym> <text> ────"
	// padding: 2 cols before symbol
	prefix := "── "
	prefixLen := 3
	symLen := 1
	spaceLen := 1 // between sym and text

	textMaxLen := width - prefixLen - symLen - spaceLen - 2 // 2 trailing ─
	if textMaxLen < 0 {
		textMaxLen = 0
	}
	truncated := truncateMsgText(msg.text, textMaxLen)
	contentLen := prefixLen + symLen + spaceLen + utf8.RuneCountInString(truncated)
	trailing := max(0, width-contentLen)

	rendered := st.BorderDefault.Render(prefix) +
		msgStyle.Render(sym+" "+truncated) +
		st.BorderDefault.Render(strings.Repeat("─", trailing))
	return rendered
}

// truncateMsgText truncates msg text to n runes with "…" if needed.
func truncateMsgText(s string, n int) string {
	if n <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	if n <= 1 {
		return "…"
	}
	return string(runes[:n-1]) + "…"
}

// Action represents an item in the command bar.
type Action struct {
	Key         string // display key (e.g. "⌃S", "F1", "Del")
	Label       string // display label (e.g. "Salvar")
	Priority    int    // higher = more left in bar
	Group       int    // for Help grouping
	Enabled     bool
	HideFromBar bool
}

// renderCmdBar renders the 1-line command bar.
// actions: all actions for current context, sorted by priority descending.
func renderCmdBar(st styles, width int, actions []Action) string {
	// anchor: "F1 Ajuda" always at the right
	anchor := st.ActionKey.Render("F1") + " " + st.ActionLabel.Render("Ajuda")
	anchorLen := len("F1") + 1 + len("Ajuda") // 8

	// "  " prefix + anchor
	// available for left actions: width - 2 - anchorLen - gap
	sortedActions := filterVisibleActions(actions)

	// Build left side iteratively, respecting available space
	sep := st.ActionSep.Render(" · ")
	sepLen := 3 // " · "
	leftParts := []string{}
	usedLen := 0

	for _, a := range sortedActions {
		part := st.ActionKey.Render(a.Key) + " " + st.ActionLabel.Render(a.Label)
		partLen := utf8.RuneCountInString(a.Key) + 1 + utf8.RuneCountInString(a.Label)

		addLen := partLen
		if len(leftParts) > 0 {
			addLen += sepLen
		}

		// Check if adding this action still leaves room for anchor
		if 2+usedLen+addLen+2+anchorLen > width && len(leftParts) > 0 {
			break
		}
		leftParts = append(leftParts, part)
		usedLen += addLen
	}

	left := "  " + strings.Join(leftParts, sep)
	leftVisLen := 2 + usedLen

	// padding between left and anchor
	padLen := max(1, width-leftVisLen-anchorLen)
	return left + strings.Repeat(" ", padLen) + anchor
}

// filterVisibleActions returns enabled, non-hidden actions, sorted by priority descending.
func filterVisibleActions(actions []Action) []Action {
	var result []Action
	for _, a := range actions {
		if a.Enabled && !a.HideFromBar {
			result = append(result, a)
		}
	}
	// sort by priority descending (simple insertion sort, small N)
	for i := 1; i < len(result); i++ {
		for j := i; j > 0 && result[j].Priority > result[j-1].Priority; j-- {
			result[j], result[j-1] = result[j-1], result[j]
		}
	}
	return result
}
