package modal

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// Intent classifies the semantic intention of a modal option.
type Intent int

const (
	IntentConfirm Intent = iota
	IntentCancel
	IntentOther
)

// ModalOption represents an action available to the user within a modal.
type ModalOption struct {
	// Keys lists the keys that activate this option.
	// Keys[0].Label is displayed in the dialog footer.
	// Other Keys are functional aliases (ex: Enter as alias for "S Overwrite").
	Keys []design.Key
	// Label is the text displayed to the user describing the action.
	Label string
	// Intent classifies the semantic intention of this action.
	Intent Intent
	// Action is the function executed when the option is chosen.
	Action func() tea.Cmd
}
