package main

import (
	"charm.land/lipgloss/v2"
	"fmt"
	"github.com/useful-toys/abditum/internal/tui/design"
	"strings"
)

func main() {
	theme := design.TokyoNight

	// Título
	titleWidth := lipgloss.Width("Operação Fake")
	fmt.Printf("Title: %q width=%d\n", "Operação Fake", titleWidth)
	titleWidth += 8
	fmt.Printf("Title with spacing: %d\n", titleWidth)

	// Ações
	actionWidth := 3
	_, w1 := design.RenderDialogAction("Enter", "Executar", theme.Border.Focused, theme)
	_, w2 := design.RenderDialogAction("Esc", "Cancelar", theme.Border.Focused, theme)
	fmt.Printf("Ação 1: %d, Ação 2: %d\n", w1, w2)
	actionWidth += w1 + 4 + 3
	actionWidth += w2 + 4 + 3
	fmt.Printf("Total actionWidth: %d\n", actionWidth)

	// Corpo
	body := "\nDeseja executar a operação fake?\nIsso simulará 5 segundos de trabalho.\n"
	paddingH := 2 * design.DialogPaddingH

	bodyLines := strings.Split(body, "\n")
	maxBodyWidth := 0
	for _, line := range bodyLines {
		w := lipgloss.Width(line) + paddingH
		if w > maxBodyWidth {
			maxBodyWidth = w
			fmt.Printf("Line %q: width=%d (visual=%d + padding=%d)\n", line, w, lipgloss.Width(line), paddingH)
		}
	}
	fmt.Printf("Max bodyWidth: %d\n", maxBodyWidth)

	width := titleWidth
	if actionWidth > width {
		width = actionWidth
	}
	if maxBodyWidth > width {
		width = maxBodyWidth
	}

	fmt.Printf("\nFinal bodyWidth returned: %d\n", width)
	fmt.Printf("For terminal width 200, innerWidth would be: %d\n", 200-2)
	fmt.Printf("Using max(bodyWidth=%d, innerWidth=%d) = %d\n", width, 200-2, max(width, 200-2))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
