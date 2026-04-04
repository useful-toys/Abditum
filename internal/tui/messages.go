package tui

import (
	"strings"
	"time"

	"charm.land/lipgloss/v2"
)

// TickMsg is sent every second by models that need to drive MessageManager.Tick().
// Define locally in each model's Init() via: tea.Tick(time.Second, func(t time.Time) tea.Msg { return TickMsg(t) })
type TickMsg time.Time

// MsgKind classifies the semantic type of a status bar message.
type MsgKind int

const (
	MsgSuccess MsgKind = iota // operação concluída — #9ece6a, TTL 3s default
	MsgInfo                   // informação neutra — #7dcfff, TTL 3s default
	MsgWarn                   // atenção — #e0af68, permanente, clearOnInput
	MsgError                  // falha — #f7768e + bold, TTL 5s default
	MsgBusy                   // spinner — #7aa2f7, permanente
	MsgHint                   // dica contextual — #565f89 + italic, permanente
)

// DisplayMessage is the read-only view of the current message for rootModel.View().
type DisplayMessage struct {
	Text  string
	Kind  MsgKind
	Frame int // animation frame index for MsgBusy (0-3), incremented by Tick()
}

// activeMessage is the internal state of MessageManager.
type activeMessage struct {
	text         string
	kind         MsgKind
	startedAt    time.Time
	frame        int
	ttlTicks     int  // 0 = permanent; decremented by Tick()
	clearOnInput bool // true = cleared by next HandleInput() call
}

// MessageManager is the centralized, write-only (for children) message bar service.
type MessageManager struct {
	current *activeMessage
}

// NewMessageManager creates an empty MessageManager.
func NewMessageManager() *MessageManager {
	return &MessageManager{}
}

// Show sets the current message. Last-write-wins.
// ttlSeconds == 0 means permanent (no auto-expiry).
// MsgBusy ignores ttlSeconds - it persists until replaced or cleared.
func (m *MessageManager) Show(kind MsgKind, text string, ttlSeconds int, clearOnInput bool) {
	ttl := ttlSeconds
	if kind == MsgBusy {
		ttl = 0
	}
	m.current = &activeMessage{
		text:         text,
		kind:         kind,
		startedAt:    time.Now(),
		ttlTicks:     ttl,
		clearOnInput: clearOnInput,
	}
}

// Clear removes the current message immediately.
func (m *MessageManager) Clear() {
	m.current = nil
}

// Current returns the displayable message for rootModel.View(), or nil if none.
func (m *MessageManager) Current() *DisplayMessage {
	if m.current == nil {
		return nil
	}
	return &DisplayMessage{
		Text:  m.current.text,
		Kind:  m.current.kind,
		Frame: m.current.frame,
	}
}

// Tick is called by rootModel every tickMsg. Decrements TTL and expires messages.
func (m *MessageManager) Tick() {
	if m.current == nil {
		return
	}
	if m.current.kind == MsgBusy {
		m.current.frame = (m.current.frame + 1) % 4
		return
	}
	if m.current.ttlTicks > 0 {
		m.current.ttlTicks--
		if m.current.ttlTicks == 0 {
			m.current = nil
		}
	}
}

// HandleInput is called by rootModel on every KeyPressMsg.
// Clears messages that were shown with clearOnInput == true.
func (m *MessageManager) HandleInput() {
	if m.current != nil && m.current.clearOnInput {
		m.current = nil
	}
}

// RenderMessageBar renders the message bar (1 line) with full-width continuous border.
// When msg == nil, renders a plain border line. Exported for PoC and rootModel use.
// Anatomy when msg != nil: ── <symbol> <text> ─...─  (fills to width)
func RenderMessageBar(msg *DisplayMessage, width int) string {
	borderStyle := StyleBorder()
	borderChar := SymBorder

	if msg == nil || width <= 0 {
		if width <= 0 {
			return ""
		}
		return borderStyle.Render(strings.Repeat(borderChar, width))
	}

	var symbol string
	var symStyle lipgloss.Style
	switch msg.Kind {
	case MsgSuccess:
		symbol = SymSuccess
		symStyle = StyleSymbol(MsgSuccess)
	case MsgInfo:
		symbol = SymInfo
		symStyle = StyleSymbol(MsgInfo)
	case MsgWarn:
		symbol = SymWarn
		symStyle = StyleSymbol(MsgWarn)
	case MsgError:
		symbol = SymError
		symStyle = StyleSymbol(MsgError)
	case MsgBusy:
		symbol = SpinnerFrame(msg.Frame)
		symStyle = StyleSymbol(MsgBusy)
	default: // MsgHint
		symbol = SymHint
		symStyle = StyleSymbol(MsgHint)
	}

	// Prefix: "── " (2 border chars + space = 3 visible chars)
	prefix := borderStyle.Render("──") + " "
	// Suffix start: " ─" (1 space + 1 border = 2 visible chars)
	suffixStart := " " + borderStyle.Render(borderChar)

	// Calculate available width for symbol + text
	prefixW := lipgloss.Width(prefix)
	suffixW := lipgloss.Width(suffixStart)
	symbolRendered := symStyle.Render(symbol + " ")
	symbolW := lipgloss.Width(symbolRendered)
	availableTextW := width - prefixW - suffixW

	var content string
	if availableTextW <= symbolW {
		// Not enough room for any text — show symbol only
		content = symbolRendered
	} else {
		// Truncate text to fit available width
		maxTextW := availableTextW - symbolW
		text := msg.Text
		truncated := false
		for lipgloss.Width(symStyle.Render(text)) > maxTextW && len(text) > 0 {
			text = text[:len(text)-1]
			truncated = true
		}
		if truncated && len(text) > 0 {
			text = text + SymEllipsis
		}
		content = symbolRendered + symStyle.Render(text)
	}

	// Calculate fill: width minus all visible chars
	visibleSoFar := prefixW + lipgloss.Width(content) + suffixW
	fillLen := width - visibleSoFar
	if fillLen < 0 {
		fillLen = 0
	}
	fill := borderStyle.Render(strings.Repeat(borderChar, fillLen))

	return prefix + content + suffixStart + fill
}
