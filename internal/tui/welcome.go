package tui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// logoLines contains the ASCII art logo for Abditum (5 lines for gradient).
var logoLines = [5]string{
	`    _    _         _ _ _                    `,
	`   / \  | |__   __| (_) |_ _   _ _ __ ___  `,
	`  / _ \ | '_ \ / _` + "`" + ` | | __| | | | '_ ` + "`" + ` _ \ `,
	` / ___ \| |_) | (_| | | |_| |_| | | | | | |`,
	`/_/   \_\_.__/ \__,_|_|\__|\__,_|_| |_| |_|`,
}

// renderWelcome renders the welcome screen work area.
// height: height of the work area (not the full terminal).
func renderWelcome(st styles, width, height int, version string) string {
	lines := make([]string, 0, height)

	// logo (5) + blank + subtitle + blank + hint = 9 content lines
	contentHeight := len(logoLines) + 1 + 1 + 1 + 1
	topPad := (height - contentHeight) / 2
	if topPad < 0 {
		topPad = 0
	}

	for i := 0; i < topPad; i++ {
		lines = append(lines, "")
	}

	// render logo with gradient colors
	colors := st.pal.logoGrad
	for i, line := range logoLines {
		colored := lipgloss.NewStyle().Foreground(lipgloss.Color(colors[i])).Render(line)
		lines = append(lines, centerText(colored, width, len([]rune(line))))
	}

	lines = append(lines, "")

	subtitle := "v" + version + "  ·  Gerenciador de Segredos"
	lines = append(lines, centerText(st.TextSecondary.Render(subtitle), width, len([]rune(subtitle))))

	lines = append(lines, "")

	hint := "F5 Novo cofre  ·  F6 Abrir cofre existente"
	lines = append(lines, centerText(st.TextDisabled.Render(hint), width, len([]rune(hint))))

	for len(lines) < height {
		lines = append(lines, "")
	}

	return strings.Join(lines[:height], "\n")
}

// centerText centers a pre-rendered string (with ANSI escapes) in a field of `width` columns.
// visLen is the visual (rune) length of the string without ANSI escapes.
func centerText(rendered string, width, visLen int) string {
	if visLen >= width {
		return rendered
	}
	pad := (width - visLen) / 2
	return strings.Repeat(" ", pad) + rendered
}
