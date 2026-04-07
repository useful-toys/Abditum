package tui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// AsciiArt is the 5-line wordmark for Abditum, ported from the reference project.
// Each line is rendered with a different color from the violet→cyan gradient palette.
const AsciiArt = `    ___    __        ___ __                  
   /   |  / /_  ____/ (_) /___  ______ ___   
  / /| | / __ \/ __  / / __/ / / / __ ` + "`" + `__ \  
 / ___ |/ /_/ / /_/ / / /_/ /_/ / / / / / / 
/_/  |_/_.___/\__,_/_/\__/\__,_/_/ /_/ /_/ `

// RenderLogo renders the Abditum wordmark with a violet→cyan gradient.
// The gradient must contain exactly 5 colors matching Design System § Gradiente do logo.
// Each of the 5 ASCII art lines gets its corresponding color.
// Returns the colored multi-line string ready for display.
func RenderLogo(t *Theme) string {
	colors := t.LogoGradient
	// Validate gradient has exactly 5 colors per design system spec
	if len(colors) != 5 {
		panic("LogoGradient must contain exactly 5 colors (one per ASCII art line)")
	}

	lines := strings.Split(AsciiArt, "\n")
	rendered := make([]string, 0, len(lines))

	for i, line := range lines {
		if i >= len(colors) {
			break // Should never happen due to validation above
		}
		style := lipgloss.NewStyle().Foreground(colors[i])
		rendered = append(rendered, style.Render(line))
	}

	return strings.Join(rendered, "\n")
}
