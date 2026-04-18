package main

import (
	"fmt"
	"strings"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
)

func main() {
	theme := design.TokyoNight
	opts := []modal.ModalOption{
		{Keys: []design.Key{design.Keys.Enter}, Label: "Executar", Action: func() { return }},
		{Keys: []design.Key{design.Keys.Esc}, Label: "Cancelar", Action: func() { return }},
	}
	
	body := "\nDeseja executar a opera��o fake?\nIsso simular� 5 segundos de trabalho.\n"
	
	// Simulate calculateBodyWidth
	paddingH := 2 * design.DialogPaddingH
	
	// T�tulo
	titleWidth := lipgloss.Width("Opera��o Fake")
	titleWidth += 8
	
	// A��es
	actionWidth := 3
	for _, opt := range opts {
		_, keyWidth := design.RenderDialogAction(opt.Keys[0].Label, opt.Label, theme.Border.Focused, theme)
		actionWidth += keyWidth + 4 + 3
	}
	
	// Corpo
	bodyLines := strings.Split(body, "\n")
	maxBodyWidth := 0
	for _, line := range bodyLines {
		w := lipgloss.Width(line) + paddingH
		if w > maxBodyWidth {
			maxBodyWidth = w
		}
	}
	
	width := titleWidth
	if actionWidth > width { width = actionWidth }
	if maxBodyWidth > width { width = maxBodyWidth }
	
	fmt.Printf("titleWidth=%d actionWidth=%d maxBodyWidth=%d => width=%d\n", titleWidth, actionWidth, maxBodyWidth, width)
	fmt.Printf("For maxWidth=200: innerWidth=%d\n", 200-2)
}
