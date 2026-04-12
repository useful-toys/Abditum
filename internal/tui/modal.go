package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// modalModel is a reusable overlay dialog that sits on top of the main frame.
// It implements modalView - rootModel overlays this using lipgloss.Place() in View().
//
// Keyboard: j/k or arrows to move selection, Enter to confirm, ESC to close (onSelect(-1)).
// Rendering: returns only the rendered box; rootModel positions via lipgloss.Place.
type modalModel struct {
	title         string
	body          string
	options       []string
	selectedIndex int
	onSelect      func(int) tea.Cmd
}

// newModal creates a new modal overlay.
func newModal(title, body string, options []string, onSelect func(int) tea.Cmd) *modalModel {
	return &modalModel{
		title:    title,
		body:     body,
		options:  options,
		onSelect: onSelect,
	}
}

// Update handles keyboard input for the modal.
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
	return nil
}

// View renders the modal box. Returns only the box - rootModel positions it.
// MUST only be called after SetSize has been called by rootModel.
func (m *modalModel) View() string {
	boxW := 50

	var content strings.Builder
	if m.body != "" {
		content.WriteString(m.body + "\n\n")
	}

	selectedStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11"))
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

	for i, opt := range m.options {
		prefix := "  "
		if i == m.selectedIndex {
			content.WriteString(selectedStyle.Render(fmt.Sprintf("\u25b6 %s", opt)) + "\n")
		} else {
			content.WriteString(normalStyle.Render(prefix+opt) + "\n")
		}
	}
	if len(m.options) == 0 {
		content.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("245")).
			Render("  Press Enter or ESC to close") + "\n")
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Width(boxW).
		Render(content.String())
}

// Shortcuts returns nil - basic dialog modals show no command bar shortcuts.
func (m *modalModel) Shortcuts() []Shortcut { return nil }

// SetAvailableSize stores terminal dimensions for layout calculations.
// For modalModel, this is called by rootModel before View() per the SetAvailableSize-before-View contract,
// but modalModel uses fixed width and doesn't need the dimensions.
// Other modalView implementations may use these dimensions to constrain their layout.
func (m *modalModel) SetAvailableSize(maxWidth, maxHeight int) {}

// Compile-time assertion: modalModel must satisfy the modalView interface.
var _ modalView = &modalModel{}
