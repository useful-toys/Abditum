package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	p := tea.NewProgram(newPocModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		os.Exit(1)
	}
}
