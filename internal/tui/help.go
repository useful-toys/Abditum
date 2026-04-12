package tui

import (
	"fmt"
	"math"
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// formatKeyForHelp converts a raw bubbletea key string to the display format
// defined by the Abditum design system (tui-specification-novo.md §573-596).
//
// Rules:
//   - Simple ctrl (ctrl+single-letter): "ctrl+q" → "⌃Q"
//   - Complex multi-modifier: "ctrl+alt+shift+q" → "Ctrl+Alt+Shift+Q"
//   - Function keys: "f1" → "F1", "shift+f6" → "Shift+F6"
//   - Named special keys: "esc" → "Esc", "del" → "Delete", etc.
//   - Arrow keys: "up" → "↑", "down" → "↓", "left" → "←", "right" → "→"
//   - Fallback: strings.ToUpper(raw)
func formatKeyForHelp(raw string) string {
	// Named specials (must check before generic prefix handling)
	switch raw {
	case "esc":
		return "Esc"
	case "enter":
		return "Enter"
	case "space":
		return "Space"
	case "backspace":
		return "Backspace"
	case "tab":
		return "Tab"
	case "del":
		return "Delete"
	case "insert":
		return "Insert"
	case "pgup":
		return "PgUp"
	case "pgdown":
		return "PgDn"
	case "home":
		return "Home"
	case "end":
		return "End"
	case "up":
		return "↑"
	case "down":
		return "↓"
	case "left":
		return "←"
	case "right":
		return "→"
	}
	// Function keys: f1–f12 (bare function key with no modifier prefix)
	if len(raw) >= 2 && raw[0] == 'f' {
		suffix := raw[1:]
		allDigits := true
		for _, c := range suffix {
			if c < '0' || c > '9' {
				allDigits = false
				break
			}
		}
		if allDigits && len(suffix) > 0 {
			return "F" + suffix
		}
	}
	// Ctrl combinations
	if strings.HasPrefix(raw, "ctrl+") {
		rest := raw[5:] // after "ctrl+"
		// Simple ctrl: exactly one letter after "ctrl+" (e.g. "ctrl+q" → rest == "q")
		if len(rest) == 1 && rest[0] >= 'a' && rest[0] <= 'z' {
			return "⌃" + strings.ToUpper(rest)
		}
		// Complex multi-modifier (e.g. "ctrl+alt+shift+q"):
		// Title-case each segment separated by "+"
		parts := strings.Split(raw, "+")
		for i, p := range parts {
			if len(p) > 0 {
				parts[i] = strings.ToUpper(p[:1]) + p[1:]
			}
		}
		return strings.Join(parts, "+")
	}
	// shift+fN (e.g. "shift+f6" → "Shift+F6")
	if strings.HasPrefix(raw, "shift+") {
		rest := raw[6:]
		return "Shift+" + formatKeyForHelp(rest)
	}
	// Fallback: uppercase the entire string
	return strings.ToUpper(raw)
}

// helpModal is the global keyboard shortcut reference overlay.
// Pushed onto the modal stack when the user presses F1.
// Dismissed via ESC or F1.
type helpModal struct {
	actions    []Action         // all registered actions for the help overlay
	groupLabel func(int) string // resolves a group int to a display label
	width      int              // terminal width for dynamic sizing
	height     int              // terminal height for dynamic sizing
	scroll     int              // current scroll offset
}

// Compile-time assertion: helpModal satisfies modalView.
var _ modalView = &helpModal{}

// newHelpModal creates a new help overlay modal.
// It accepts a pre-computed slice of all actions and a groupLabel function
// so that golden tests can instantiate it with plain []Action fixtures,
// without needing an ActionManager.
func newHelpModal(actions []Action, groupLabel func(int) string) *helpModal {
	return &helpModal{actions: actions, groupLabel: groupLabel}
}

// SetAvailableSize sets the maximum available dimensions for dynamic modal sizing.
func (m *helpModal) SetAvailableSize(maxWidth, maxHeight int) {
	m.width = maxWidth
	m.height = maxHeight
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
	if m.width == 0 || m.height == 0 {
		panic(fmt.Sprintf("helpModal.View() called without SetSize: width=%d height=%d", m.width, m.height))
	}
	// Dynamic sizing per DS: max 60 cols or 70% of terminal
	maxW := 60
	if m.width > 0 {
		pctW := int(float64(m.width) * 0.7)
		if pctW < maxW {
			maxW = pctW
		}
	}
	boxW := maxW

	allActions := m.actions
	lines := m.buildContentLines(allActions)
	totalLines := len(lines)

	// Dialog layout: top border(1) + content(innerH) + bottom border(1)
	// Content area: top padding(1) + action lines(usableH) + bottom padding(1)
	// Total dialog = usableH + 4 lines — must fit in terminal.
	maxUsable := m.height - 4
	if maxUsable > 20 {
		maxUsable = 20 // cap for large terminals
	}
	if maxUsable < 3 {
		maxUsable = 3 // minimum usable action lines
	}
	usableH := maxUsable
	innerH := usableH + 2 // content area includes padding lines

	// Clamp visible window to available content
	start := m.scroll
	if start > totalLines-usableH {
		start = totalLines - usableH
	}
	if start < 0 {
		start = 0
	}
	end := start + usableH
	if end > totalLines {
		end = totalLines
	}
	visibleLines := lines[start:end]

	hasAbove := start > 0
	hasBelow := end < totalLines

	return m.renderDialog(visibleLines, boxW, innerH, hasAbove, hasBelow, totalLines, start, usableH)
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
		rightBorder := borderStyle.Render("│")
		if indicator != "" {
			indicatorStyled := actionStyle.Render(indicator)
			contentLines = append(contentLines, left+"  "+lipgloss.NewStyle().Width(innerW).Render(line)+"  "+indicatorStyled)
		} else {
			contentLines = append(contentLines, left+"  "+lipgloss.NewStyle().Width(innerW).Render(line)+"  "+rightBorder)
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
	return len(m.buildContentLines(m.actions))
}

// contentHeight returns the visible content height (excluding borders).
// Must match the usableH calculation in View() exactly.
func (m *helpModal) contentHeight() int {
	maxUsable := m.height - 4
	if maxUsable > 20 {
		maxUsable = 20
	}
	if maxUsable < 3 {
		maxUsable = 3
	}
	return maxUsable
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
		// Group 0: no header — actions listed without section label
		if grp == 0 {
			for _, act := range grouped[grp] {
				if len(act.Keys) == 0 {
					continue
				}
				lines = append(lines, "  "+keyStyle.Render(formatKeyForHelp(act.Keys[0]))+"  "+act.Description)
			}
			continue
		}
		label := m.groupLabel(grp)
		if label != "" {
			lines = append(lines, groupStyle.Render(label))
		}
		for _, act := range grouped[grp] {
			if len(act.Keys) == 0 {
				continue
			}
			lines = append(lines, "  "+keyStyle.Render(formatKeyForHelp(act.Keys[0]))+"  "+act.Description)
		}
	}

	return lines
}

// Shortcuts returns an empty slice — help modal has no command bar actions.
// Dismiss is handled via ESC/F1 in Update() and indicated in the bottom border.
func (m *helpModal) Shortcuts() []Shortcut {
	return nil
}
