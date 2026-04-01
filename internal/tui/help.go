package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// helpModal is the global keyboard shortcut reference overlay.
// It is pushed onto the modal stack when the user presses "?".
// It reads ActionManager.All() to show all currently registered actions,
// grouped by their Group field.
// Dismissed via ESC or "?" — rootModel pops it from the stack on popModalMsg.
type helpModal struct {
	actions *ActionManager
	width   int
	height  int
}

// Compile-time assertion: helpModal satisfies childModel.
var _ childModel = &helpModal{}

// newHelpModal creates a new help overlay modal.
func newHelpModal(actions *ActionManager) *helpModal {
	return &helpModal{actions: actions}
}

// Update handles keyboard input for the help modal.
// ESC or "?" dismisses the overlay by emitting popModalMsg.
func (m *helpModal) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc", "?":
			return func() tea.Msg { return popModalMsg{} }
		}
	}
	return nil
}

// View renders the keyboard shortcut reference box, centered on the terminal.
// It calls ActionManager.All() to list all currently registered actions.
func (m *helpModal) View() string {
	boxW := 60
	if m.width > 0 && m.width-4 < boxW {
		boxW = m.width - 4
	}

	allActions := m.actions.All()
	content := m.buildContent(allActions)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Width(boxW).
		Render(content)

	if m.width == 0 || m.height == 0 {
		return box
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

// buildContent formats the action list into a readable shortcut reference.
// Actions are grouped by their Group field; groups are displayed in the order
// they first appear.
func (m *helpModal) buildContent(actions []Action) string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62"))
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	sepStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	var b strings.Builder
	b.WriteString(titleStyle.Render("⌨  Keyboard Shortcuts") + "\n")
	b.WriteString(sepStyle.Render(strings.Repeat("─", 50)) + "\n")

	if len(actions) == 0 {
		b.WriteString(sepStyle.Render("  No actions registered") + "\n")
	} else {
		currentGroup := ""
		groupStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
		for _, act := range actions {
			if act.Group != currentGroup {
				currentGroup = act.Group
				if currentGroup != "" {
					b.WriteString("\n" + groupStyle.Render(currentGroup) + "\n")
				}
			}
			b.WriteString(fmt.Sprintf("  %s  %s\n",
				keyStyle.Render(act.Key), act.Description))
		}
	}

	b.WriteString("\n" + sepStyle.Render("  Esc or ?  close this help") + "\n")
	return b.String()
}

// SetSize stores terminal dimensions for centered placement.
func (m *helpModal) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Context returns an empty FlowContext — help modal has no navigation state.
func (m *helpModal) Context() FlowContext {
	return FlowContext{}
}

// ChildFlows returns nil — help modal has no child-specific flows.
func (m *helpModal) ChildFlows() []flowDescriptor {
	return nil
}
