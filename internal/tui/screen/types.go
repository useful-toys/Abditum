package screen

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// Action is a temporary placeholder to avoid circular imports.
// The real Action type is defined in the tui package as a struct.
type Action = interface{}

// ChildView matches the interface from tui package to avoid importing it.
// Defines the contract for renderable components on the main screen.
type ChildView interface {
	// Render returns the string representation of the component for display.
	// height and width define available dimensions. theme is passed by pointer
	// to avoid unnecessary copying — design.Theme is about 400 bytes.
	Render(height, width int, theme *design.Theme) string

	// HandleKey processes keyboard events and returns a command or nil.
	HandleKey(msg tea.KeyMsg) tea.Cmd

	// HandleEvent processes custom events for the component.
	HandleEvent(event any)

	// HandleTeaMsg processes Bubble Tea framework messages.
	HandleTeaMsg(msg tea.Msg) tea.Cmd

	// Update is called to update the component's state.
	Update(msg tea.Msg) tea.Cmd

	// Actions returns the actions available in this view.
	// May return nil if the view has no actions of its own.
	Actions() []interface{}
}
