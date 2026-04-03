package tui

import (
	"fmt"
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// helpModal is the global keyboard shortcut reference overlay.
// It is pushed onto the modal stack when the user presses "?".
// Dismissed via ESC or "?".
type helpModal struct {
	actions *ActionManager
}

// Compile-time assertion: helpModal satisfies modalView.
var _ modalView = &helpModal{}

// newHelpModal creates a new help overlay modal.
func newHelpModal(actions *ActionManager) *helpModal {
	return &helpModal{actions: actions}
}

// Update handles keyboard input for the help modal.
func (m *helpModal) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc", "f1":
			return func() tea.Msg { return popModalMsg{} }
		}
	}
	return nil
}

// View renders the keyboard shortcut reference box.
// Returns only the box - rootModel positions it via lipgloss.Place.
func (m *helpModal) View() string {
	boxW := 60

	allActions := m.actions.All()
	content := m.buildContent(allActions)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Width(boxW).
		Render(content)
}

// Shortcuts returns the dismiss shortcut for the command bar.
func (m *helpModal) Shortcuts() []Shortcut {
	return []Shortcut{{Key: "esc", Label: "Close"}}
}

// buildContent formats the action list into a readable shortcut reference.
func (m *helpModal) buildContent(actions []Action) string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62"))
	keyStyle   := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	sepStyle   := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	groupStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))

	var b strings.Builder
	b.WriteString(titleStyle.Render("  Keyboard Shortcuts") + "\n")
	b.WriteString(sepStyle.Render(strings.Repeat("-", 50)) + "\n")

	if len(actions) == 0 {
		b.WriteString(sepStyle.Render("  No actions registered") + "\n")
	} else {
		// Collect unique groups in sorted order
		groupSeen := make(map[int]bool)
		var groupOrder []int
		grouped := make(map[int][]Action)
		for _, act := range actions {
			if !groupSeen[act.Group] {
				groupSeen[act.Group] = true
				groupOrder = append(groupOrder, act.Group)
			}
			grouped[act.Group] = append(grouped[act.Group], act)
		}
		sort.Ints(groupOrder)

		for _, grp := range groupOrder {
			label := m.actions.GroupLabel(grp)
			if label != "" {
				b.WriteString("\n" + groupStyle.Render(label) + "\n")
			}
			for _, act := range grouped[grp] {
				if len(act.Keys) == 0 {
					continue
				}
				b.WriteString(fmt.Sprintf("  %s  %s\n",
					keyStyle.Render(act.Keys[0]), act.Description))
			}
		}
	}

	b.WriteString("\n" + sepStyle.Render("  Esc or F1  close this help") + "\n")
	return b.String()
}
