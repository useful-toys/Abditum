package modal

import tea "charm.land/bubbletea/v2"

type Intent int

const (
	IntentConfirm Intent = iota
	IntentCancel
	IntentOther
)

type ModalOption struct {
	Keys   []string
	Label  string
	Intent Intent
	Action func() tea.Cmd
}
