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
// Each of the 5 lines gets a distinct lipgloss foreground color.
// Returns the colored multi-line string ready for display.
func RenderLogo(t *Theme) string {
	colors := t.LogoGradient

	lines := strings.Split(AsciiArt, "\n")
	var renderedLogo strings.Builder

	for i, line := range lines {
		if i >= len(colors) {
			i = len(colors) - 1
		}
		style := lipgloss.NewStyle().Foreground(colors[i])
		renderedLogo.WriteString(style.Render(line) + "\n")
	}

	return renderedLogo.String()
}
