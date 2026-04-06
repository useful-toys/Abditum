package tui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// headerModel is a stateless model that renders the application header.
// It receives all necessary data at render time.
type headerModel struct{}

// Render composes the two-line application header based on current state.
// Design decisions D-01 in 06-UI-SPEC.md define the layout and styling.
func (h *headerModel) Render(width int, vaultName string, isDirty bool, activeArea workArea, theme *Theme) string {
	if width == 0 {
		return ""
	}

	// Styles
	appNameStyle := lipgloss.NewStyle().Foreground(theme.AccentPrimary).Bold(true)
	separatorStyle := lipgloss.NewStyle().Foreground(theme.SurfaceRaised)
	vaultNameStyle := lipgloss.NewStyle().Foreground(theme.TextPrimary).Bold(true)
	dirtyStyle := lipgloss.NewStyle().Foreground(theme.SemanticWarning).Bold(true)
	tabStyle := lipgloss.NewStyle().Foreground(theme.TextSecondary)
	activeTabStyle := lipgloss.NewStyle().Foreground(theme.AccentPrimary).Bold(true)

	// Line 1 content
	var line1 string
	if vaultName == "" {
		// State 1 — No vault (welcome)
		line1 = appNameStyle.Render("  Abditum")
	} else {
		// State 2 — Vault open
		// Truncate vaultName if it's too long
		maxVaultNameWidth := width - lipgloss.Width("  Abditum ·   ") - lipgloss.Width(" •") - lipgloss.Width("  Cofre · Modelos · Config  ")
		if maxVaultNameWidth < 0 {
			maxVaultNameWidth = 0
		}
		displayVaultName := vaultName
		if lipgloss.Width(vaultName) > maxVaultNameWidth {
			runes := []rune(vaultName)
			if len(runes) > maxVaultNameWidth {
				displayVaultName = string(runes[:maxVaultNameWidth-1]) + "…"
			}
		}

		dirtyIndicator := ""
		if isDirty {
			dirtyIndicator = dirtyStyle.Render(" •")
		}

		// Tabs
		var tabs []string
		tabs = append(tabs, renderTab("Cofre", workAreaVault, activeArea, tabStyle, activeTabStyle))
		tabs = append(tabs, renderTab("Modelos", workAreaTemplates, activeArea, tabStyle, activeTabStyle))
		tabs = append(tabs, renderTab("Config", workAreaSettings, activeArea, tabStyle, activeTabStyle))
		allTabs := strings.Join(tabs, separatorStyle.Render(" · "))

		line1Left := lipgloss.JoinHorizontal(lipgloss.Top,
			appNameStyle.Render("  Abditum"),
			separatorStyle.Render(" · "),
			vaultNameStyle.Render(displayVaultName),
			dirtyIndicator,
		)

		// Calculate space for tabs on the right
		remainingWidth := width - lipgloss.Width(line1Left)
		if remainingWidth < 0 {
			remainingWidth = 0
		}
		line1Right := lipgloss.NewStyle().Width(remainingWidth).Align(lipgloss.Right).Render(allTabs)
		line1 = lipgloss.JoinHorizontal(lipgloss.Top, line1Left, line1Right)
	}

	// Line 2: Separator
	line2 := separatorStyle.Render(strings.Repeat("─", width))

	return line1 + "\n" + line2
}

func renderTab(label string, area, activeArea workArea, defaultStyle, activeStyle lipgloss.Style) string {
	if area == activeArea {
		return activeStyle.Render(label)
	}
	return defaultStyle.Render(label)
}
