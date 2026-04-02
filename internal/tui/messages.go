package tui

import "time"

// MsgKind classifies the semantic type of a status bar message.
type MsgKind int

const (
	MsgInfo  MsgKind = iota // operation completed successfully
	MsgWarn                  // attention - imminent lock, external conflict
	MsgError                 // failure - save failed, corruption
	MsgBusy                  // operation in progress - spinner animation (no TTL)
	MsgHint                  // contextual explanation - field description
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
