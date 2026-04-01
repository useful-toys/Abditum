package tui

import tea "charm.land/bubbletea/v2"

// dialogs provides factory functions for pre-defined modal dialogs (D-18).
// Each factory returns a tea.Cmd that emits a pushModalMsg, which rootModel
// intercepts and pushes onto the modal stack. Dialogs are stateless — they
// produce a tea.Cmd and nothing else. No shared object needs to be passed at
// construction time.
//
// Usage by children: return dialogs.Message(...) or dialogs.Confirm(...) as
// the Cmd from Update(). No direct access to rootModel's modal stack.

// Message creates an informational dialog that is dismissed via ESC or Enter.
// title is shown in the modal header; text is the body content.
func Message(title, text string) tea.Cmd {
	return func() tea.Msg {
		m := &modalModel{
			content: title + "\n\n" + text,
		}
		return pushModalMsg{modal: m}
	}
}

// Confirm creates a yes/no confirmation dialog.
// question is displayed to the user; onYes fires when the user confirms,
// onNo fires when the user cancels or presses ESC.
//
// Phase 5 stub: the modal content is rendered as plain text; interactive
// key handling (yes/no selection) is implemented in later phases when
// concrete confirm modals are wired into flows.
func Confirm(question string, onYes, onNo tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		m := &confirmModalModel{
			modalModel: modalModel{
				content: question + "\n\n[Y] Yes   [N] No / ESC",
			},
			onYes: onYes,
			onNo:  onNo,
		}
		return pushModalMsg{modal: &m.modalModel}
	}
}

// confirmModalModel is the internal implementation of a confirmation dialog.
// It embeds modalModel and overrides Update to handle yes/no key events.
type confirmModalModel struct {
	modalModel
	onYes tea.Cmd
	onNo  tea.Cmd
}

// Update handles key events for yes/no selection.
// "y" or "enter" fires onYes; "n", "q", or "esc" fires onNo.
func (m *confirmModalModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "y", "enter":
			return m.onYes
		case "n", "q", "esc":
			return m.onNo
		}
	}
	return nil
}
