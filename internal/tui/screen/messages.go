package screen

import "github.com/useful-toys/abditum/internal/tui/design"

// WorkAreaChangedMsg is emitted by HeaderView when the user clicks on a tab.
// Defined in package screen (not in package tui) to avoid import cycle.
// RootModel in package tui does: case screen.WorkAreaChangedMsg:
type WorkAreaChangedMsg struct {
	Area design.WorkArea
}
