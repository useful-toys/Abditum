package tui

import tea "charm.land/bubbletea/v2"

// childModel - D-01: 3 methods only.
// closures that capture child state directly.
type childModel interface {
	Update(tea.Msg) tea.Cmd
	View() string
	SetSize(w, h int)
}

// Shortcut is a key+label pair displayed in the command bar while a modal is active.
type Shortcut struct {
	Key   string
	Label string
}

// modalView - D-02: Separate from childModel. No SetSize - modals auto-size by content.
// rootModel.View() positions the modal via lipgloss.Place.
type modalView interface {
	Update(tea.Msg) tea.Cmd
	View() string
	Shortcuts() []Shortcut
}

// modalResult - D-03: Marker interface for messages carrying sensitive data.
// rootModel routes these EXCLUSIVELY to activeFlow.
type modalResult interface {
	isModalResult()
}

type passwordEntryResult struct {
	Password  []byte
	Cancelled bool
}

func (passwordEntryResult) isModalResult() {}

type passwordCreateResult struct {
	Password  []byte
	Cancelled bool
}

func (passwordCreateResult) isModalResult() {}

type filePickerResult struct {
	Path      string
	Cancelled bool
}

func (filePickerResult) isModalResult() {}

// flowHandler - D-04: Init() is called by rootModel immediately after setting activeFlow.
type flowHandler interface {
	Init() tea.Cmd
	Update(tea.Msg) tea.Cmd
}

// startFlowMsg - D-08: clears orphan modals, sets activeFlow, calls activeFlow.Init().
type startFlowMsg struct{ flow flowHandler }

// endFlowMsg - D-08: signals flow completion. rootModel sets activeFlow = nil.
type endFlowMsg struct{}

// endFlow returns a Cmd that emits endFlowMsg. Flows call this on any exit path.
func endFlow() tea.Cmd {
	return func() tea.Msg { return endFlowMsg{} }
}

// pushModalMsg - emitted by dialog factories; rootModel appends modal to its stack.
type pushModalMsg struct{ modal modalView }

// popModalMsg - emitted by modals to request removal from the stack.
type popModalMsg struct{}
