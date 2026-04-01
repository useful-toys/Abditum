package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// modalModel is a reusable overlay dialog that sits on top of the main frame.
// It implements childModel so it can receive domain message broadcasts via liveModels().
//
// Keyboard: j/k or arrows to move selection, Enter to confirm, ESC to close (onSelect(-1)).
// Rendering: centered in terminal using lipgloss.Place().
//
// rootModel owns the modal stack; modals are pushed via pushModalMsg and
// popped when onSelect fires or ESC is pressed.
type modalModel struct {
	title         string
	body          string
	options       []string
	selectedIndex int
	onSelect      func(int) tea.Cmd
	width         int
	height        int
}

// newModal creates a new modal overlay.
// title: shown in the border title.
// body: text shown above the option list.
// options: selectable choices (empty list = informational modal, dismiss with Enter/ESC).
// onSelect: callback receiving the selected index, or -1 on ESC.
func newModal(title, body string, options []string, onSelect func(int) tea.Cmd) *modalModel {
	return &modalModel{
		title:    title,
		body:     body,
		options:  options,
		onSelect: onSelect,
	}
}

// Update handles keyboard input for the modal.
// Returns the tea.Cmd from onSelect when the user confirms or dismisses.
// Returns a popModalMsg when the modal should be removed from the stack.
func (m *modalModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "j", "down":
			if m.selectedIndex < len(m.options)-1 {
				m.selectedIndex++
			}
		case "k", "up":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "enter":
			if m.onSelect != nil {
				cmd := m.onSelect(m.selectedIndex)
				return tea.Batch(cmd, func() tea.Msg { return popModalMsg{} })
			}
			return func() tea.Msg { return popModalMsg{} }
		case "esc":
			if m.onSelect != nil {
				cmd := m.onSelect(-1)
				return tea.Batch(cmd, func() tea.Msg { return popModalMsg{} })
			}
			return func() tea.Msg { return popModalMsg{} }
		}
	}
	// Domain messages (tick, vault events) are received but ignored by modals.
	return nil
}

// View renders the modal box. Returns a string (not tea.View) — rootModel
// overlays this using lipgloss.Place() in its View() compositor.
func (m *modalModel) View() string {
	boxW := 50
	if m.width > 0 && m.width-8 < boxW {
		boxW = m.width - 8
	}
	if boxW < 20 {
		boxW = 20
	}

	var content strings.Builder
	if m.body != "" {
		content.WriteString(m.body + "\n\n")
	}

	selectedStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11"))
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

	for i, opt := range m.options {
		prefix := "  "
		if i == m.selectedIndex {
			content.WriteString(selectedStyle.Render(fmt.Sprintf("▶ %s", opt)) + "\n")
		} else {
			content.WriteString(normalStyle.Render(prefix+opt) + "\n")
		}
	}
	if len(m.options) == 0 {
		content.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("245")).
			Render("  Press Enter or ESC to close") + "\n")
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Width(boxW).
		Render(content.String())

	if m.width == 0 || m.height == 0 {
		return box
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

// SetSize stores terminal dimensions for centered placement.
func (m *modalModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Context returns an empty FlowContext — modals don't expose navigation state.
func (m *modalModel) Context() FlowContext {
	return FlowContext{}
}

// ChildFlows returns nil — modals do not register child-specific flows.
func (m *modalModel) ChildFlows() []flowDescriptor {
	return nil
}

// popModalMsg is an internal message emitted by a modal when it is done.
// rootModel handles this by popping the topmost modal from the stack.
type popModalMsg struct{}

// Compile-time assertion: modalModel must satisfy the childModel interface.
var _ childModel = &modalModel{}
