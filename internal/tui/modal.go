package tui

import tea "charm.land/bubbletea/v2"

// modalModel represents a modal dialog overlaid on the main frame.
// Modals are displayed as overlays above the base frame using lipgloss.Place().
// They implement childModel so they receive domain message broadcasts and participate
// in liveModels() iteration alongside base child models.
//
// Modal stack mechanics (D-10):
//   - Push: rootModel appends to modals slice on pushModalMsg
//   - Pop: topmost modal is removed on dismiss (ESC or selection)
//   - Input routing: only the topmost modal receives keyboard/mouse events
//   - Domain messages: all modals in the stack receive domain messages via liveModels()
type modalModel struct {
	width  int
	height int
	// content is the rendered content string for this modal.
	// Concrete implementations (password entry, confirmation, help) embed modalModel
	// and provide their own View() via the outer struct.
	content string
}

// Update processes a message and returns a command.
// Stub implementation — concrete modal types embed modalModel and override Update.
func (m *modalModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// View returns the modal content as a string.
// rootModel composites this string into the frame using lipgloss.Place().
func (m *modalModel) View() string {
	return m.content
}

// SetSize stores the allocated terminal dimensions for this modal.
// rootModel calls SetSize before rendering so the modal can adapt its layout.
func (m *modalModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Context returns an empty FlowContext.
// Modals do not expose navigation state — they are transient UI overlays.
func (m *modalModel) Context() FlowContext {
	return FlowContext{}
}

// ChildFlows returns nil.
// Modals do not own flows — they are created by flows or Cmd factories.
func (m *modalModel) ChildFlows() []flowDescriptor {
	return nil
}
