package tui

import tea "charm.land/bubbletea/v2"

// dialogs provides stateless factory functions for pre-defined modal dialogs.
// Each factory returns a tea.Cmd that emits pushModalMsg, which rootModel
// intercepts to push the modal onto its stack.
//
// Children return these Cmds from their Update() — they never access the
// modal stack directly.

// NewMessage creates an informational dialog dismissed via Enter or ESC.
// title: modal border title. text: message body.
func NewMessage(title, text string) tea.Cmd {
	return func() tea.Msg {
		m := newModal(title, text, nil, nil)
		return pushModalMsg{modal: m}
	}
}

// NewConfirm creates a yes/no confirmation dialog.
// question: shown as the modal body.
// onYes: Cmd returned when user selects "Yes" (index 0).
// onNo: Cmd returned when user selects "No" or presses ESC (index 1 or -1).
func NewConfirm(question string, onYes, onNo tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		m := newModal("Confirm", question, []string{"Yes", "No"}, func(idx int) tea.Cmd {
			if idx == 0 {
				return onYes
			}
			return onNo
		})
		return pushModalMsg{modal: m}
	}
}
