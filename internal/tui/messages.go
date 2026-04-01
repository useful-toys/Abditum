package tui

// MessageManager is the centralized API for the message bar zone.
// Children write to it (Set); rootModel.View() reads from it (Current).
// MessageManager is a shared mutable object: one instance lives in rootModel
// and is passed by pointer to every child at construction.
//
// Write-only from children's perspective — children must NOT read from it.
// rootModel may also write directly for global hints (e.g., "vault locked").
type MessageManager struct {
	current string
}

// NewMessageManager creates a new, empty MessageManager.
func NewMessageManager() *MessageManager {
	return &MessageManager{}
}

// Set replaces the current message. Called synchronously from Update().
// An empty string clears the message bar (no message shown).
func (m *MessageManager) Set(msg string) {
	m.current = msg
}

// Current returns the current message for rendering in the message bar.
// Called from rootModel.View() only.
func (m *MessageManager) Current() string {
	return m.current
}

// Clear resets the message bar to empty.
func (m *MessageManager) Clear() {
	m.current = ""
}
