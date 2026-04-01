package tui

import "charm.land/lipgloss/v2"

// MessageManager is the centralized API for setting the message bar content.
// It is a shared mutable object instantiated in main.go and passed to rootModel,
// which in turn passes it to every child at construction time.
//
// Guiding analogy: just as ActionManager is the API for defining available actions,
// MessageManager is the API for setting what message is shown in the message bar (D-17).
//
// Usage:
//   - Writing: any child calls Set(text) during its own Update() — synchronous mutation
//   - Reading: rootModel.View() calls Current() when composing the message bar zone
//   - Children MUST NOT read from MessageManager — it is write-only from their perspective
//
// Since Bubble Tea re-renders after every Update(), the view is always fresh.
// No notification or broadcast mechanism is needed.
type MessageManager struct {
	current  string
	severity MessageSeverity
}

// MessageSeverity classifies the urgency of a message for display styling.
type MessageSeverity int

const (
	// MessageInfo is a neutral informational message.
	MessageInfo MessageSeverity = iota
	// MessageSuccess indicates a successful operation.
	MessageSuccess
	// MessageWarning indicates a warning that requires attention.
	MessageWarning
	// MessageError indicates an error or failed operation.
	MessageError
)

// Set updates the message bar content with the default Info severity.
func (m *MessageManager) Set(text string) {
	m.current = text
	m.severity = MessageInfo
}

// SetWithSeverity updates the message bar content and its display severity.
func (m *MessageManager) SetWithSeverity(text string, severity MessageSeverity) {
	m.current = text
	m.severity = severity
}

// Clear empties the message bar.
func (m *MessageManager) Clear() {
	m.current = ""
	m.severity = MessageInfo
}

// Current returns the current message text.
// Called by rootModel.View() when composing the message bar zone.
func (m *MessageManager) Current() string {
	return m.current
}

// Severity returns the current message severity for display styling.
func (m *MessageManager) Severity() MessageSeverity {
	return m.severity
}

// NewMessageManager creates a new MessageManager with no active message.
func NewMessageManager() *MessageManager {
	return &MessageManager{}
}

// messageBarStyle is the default lipgloss style for the message bar zone.
// Phase 5 uses a minimal placeholder style; refined in later phases.
var messageBarStyle = lipgloss.NewStyle().
	Padding(0, 1)
