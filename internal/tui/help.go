package tui

import (
	"math"
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// helpModal is the global keyboard shortcut reference overlay.
// Pushed onto the modal stack when the user presses F1.
// Dismissed via ESC or F1.
type helpModal struct {
	actions *ActionManager
	width   int // terminal width for dynamic sizing
	height  int // terminal height for dynamic sizing
	scroll  int // current scroll offset
}

// Compile-time assertion: helpModal satisfies modalView.
var _ modalView = &helpModal{}

// newHelpModal creates a new help overlay modal.
func newHelpModal(actions *ActionManager) *helpModal {
	return &helpModal{actions: actions}
}

// SetSize sets the terminal dimensions for dynamic modal sizing.
func (m *helpModal) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Update handles keyboard input for the help modal.
func (m *helpModal) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc", "f1":
			return func() tea.Msg { return popModalMsg{} }
		case "up":
			if m.scroll > 0 {
				m.scroll--
			}
		case "down":
			m.scroll++
		case "pgup":
			m.scroll = max(0, m.scroll-(m.contentHeight()-1))
		case "pgdown":
			m.scroll += m.contentHeight() - 1
		case "home":
			m.scroll = 0
		case "end":
			m.scroll = m.totalLines() - m.contentHeight()
		}
		// Clamp scroll
		maxScroll := max(0, m.totalLines()-m.contentHeight())
		if m.scroll > maxScroll {
			m.scroll = maxScroll
		}
	}
	return nil
}

// View renders the help modal with title in top border and action in bottom border.
// Follows DS dialog anatomy (§436-458): title embedded in top border, action bar in bottom border.
func (m *helpModal) View() string {
	// Dynamic sizing per DS: max 70 cols or 80% of terminal
	maxW := 70
	if m.width > 0 {
		pctW := int(float64(m.width) * 0.8)
		if pctW < maxW {
			maxW = pctW
		}
	}
	boxW := maxW

	allActions := m.actions.All()
	lines := m.buildContentLines(allActions)
	totalLines := len(lines)

	// Max height: 80% of terminal, minus top border + bottom border
	maxH := m.height
	if m.height > 0 {
		pctH := int(float64(m.height) * 0.8)
		if pctH < maxH {
			maxH = pctH
		}
	}
	contentH := maxH - 2 // minus top border line + bottom border line
	if contentH < 5 {
		contentH = 5
	}

	// Dialog grows to fit content, capped at contentH
	dialogH := min(totalLines+2, contentH) // +2 for top/bottom borders
	if dialogH < 7 {
		dialogH = 7 // minimum: title + 4 content + bottom border
	}
	innerH := dialogH - 2 // content area only

	// Apply scroll window
	start := m.scroll
	if start > totalLines-innerH {
		start = totalLines - innerH
	}
	if start < 0 {
		start = 0
	}
	end := start + innerH
	if end > totalLines {
		end = totalLines
	}
	visibleLines := lines[start:end]

	hasAbove := start > 0
	hasBelow := end < totalLines

	return m.renderDialog(visibleLines, boxW, innerH, hasAbove, hasBelow, totalLines, start, innerH)
}

// renderDialog builds the full dialog with title in top border and action in bottom border.
func (m *helpModal) renderDialog(lines []string, boxW, innerH int, hasAbove, hasBelow bool, totalLines, start, viewH int) string {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorBorderDefault))
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ColorTextPrimary))
	actionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTextSecondary))
	innerW := boxW - 6 // 2 borders + 4 padding (2 each side)

	// Top border: ╭── Title ──────────────╮
	titleText := "Ajuda — Atalhos e Ações"
	titleRendered := titleStyle.Render(titleText)
	titleW := lipgloss.Width(titleText)
	topFill := boxW - 6 - titleW
	if topFill < 1 {
		topFill = 1
	}
	topBorder := borderStyle.Render("╭── ") + titleRendered + borderStyle.Render(" "+strings.Repeat(SymBorder, topFill)+"╮")

	// Content lines with 2-col horizontal padding + 1-line vertical padding (DS §240-241)
	var contentLines []string

	// Top padding line (DS §241)
	contentLines = append(contentLines, borderStyle.Render("│")+strings.Repeat(" ", innerW+4)+borderStyle.Render("│"))

	for i, line := range lines {
		var indicator string
		if hasAbove && i == 0 {
			indicator = "↑"
		} else if hasBelow && i == len(lines)-1 {
			indicator = "↓"
		} else if totalLines > viewH && viewH > 1 {
			thumbLine := int(math.Round(float64(start) / float64(totalLines-viewH) * float64(viewH-1)))
			if i == thumbLine {
				indicator = "■"
			}
		}

		left := borderStyle.Render("│")
		if indicator != "" {
			right := actionStyle.Render(indicator)
			contentLines = append(contentLines, left+"  "+lipgloss.NewStyle().Width(innerW-2).Render(line)+right)
		} else {
			right := borderStyle.Render("│")
			contentLines = append(contentLines, left+"  "+lipgloss.NewStyle().Width(innerW).Render(line)+"  "+right)
		}
	}

	// Bottom padding line (DS §241)
	contentLines = append(contentLines, borderStyle.Render("│")+strings.Repeat(" ", innerW+4)+borderStyle.Render("│"))

	// Pad content to fill inner height
	emptyLine := borderStyle.Render("│") + strings.Repeat(" ", innerW+4) + borderStyle.Render("│")
	for len(contentLines) < innerH {
		contentLines = append(contentLines, emptyLine)
	}

	// Bottom border: ╰──────────────── Esc Fechar ─╯
	actionText := "Esc Fechar"
	actionRendered := borderStyle.Render(actionText)
	actionW := lipgloss.Width(actionText)
	bottomFill := boxW - 5 - actionW // ╰(1) + space(1) + action + space(1) + ─╯(2) = 5
	if bottomFill < 1 {
		bottomFill = 1
	}
	bottomBorder := borderStyle.Render("╰") + borderStyle.Render(strings.Repeat(SymBorder, bottomFill)) + " " + actionRendered + " " + borderStyle.Render("─╯")

	return topBorder + "\n" + strings.Join(contentLines, "\n") + "\n" + bottomBorder
}

// totalLines returns the total number of content lines.
func (m *helpModal) totalLines() int {
	allActions := m.actions.All()
	return len(m.buildContentLines(allActions))
}

// contentHeight returns the visible content height (excluding borders).
func (m *helpModal) contentHeight() int {
	maxH := m.height
	if m.height > 0 {
		pctH := int(float64(m.height) * 0.8)
		if pctH < maxH {
			maxH = pctH
		}
	}
	return maxH - 2 // minus top border + bottom border
}

// buildContentLines formats actions into display lines.
func (m *helpModal) buildContentLines(actions []Action) []string {
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorAccentPrimary))
	groupStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ColorTextSecondary))

	var lines []string

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
			lines = append(lines, groupStyle.Render(label))
		}
		for _, act := range grouped[grp] {
			if len(act.Keys) == 0 {
				continue
			}
			lines = append(lines, "  "+keyStyle.Render(act.Keys[0])+"  "+act.Description)
		}
	}

	return lines
}

// Shortcuts returns the dismiss shortcut for the command bar.
func (m *helpModal) Shortcuts() []Shortcut {
	return []Shortcut{{Key: "esc", Label: "Fechar ajuda"}}
}
