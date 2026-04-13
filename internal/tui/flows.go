package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/storage"
)

// childModel represents a UI component managed by rootModel.
// It must be displayed as the primary or secondary panel in a work area.
// childModel is managed by rootModel (via liveWorkChildren and activeChild).
// Each childModel implementation is a full tea.Model that handles its own state and updates.
type childModel interface {
	Update(tea.Msg) tea.Cmd
	View(width, height int, theme *Theme) string
}

// Shortcut is a key+label pair displayed in the command bar while a modal is active.
type Shortcut struct {
	Key   string
	Label string
}

// modalView represents an overlay modal dialog.
// Modals are managed by rootModel via the modal stack (m.modals).
// Unlike childModels, modals can be pushed/popped dynamically during a flow.
// Each modalView implementation must handle its own state and update logic.
type modalView interface {
	Update(tea.Msg) tea.Cmd
	View(maxWidth, maxHeight int, theme *Theme) string
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

// vaultOpenedMsg - emitted by flows when a vault is successfully opened/created.
// Carries the path and file metadata needed for external change detection.
type vaultOpenedMsg struct {
	Path     string               // path to the opened vault file
	Metadata storage.FileMetadata // file metadata snapshot taken after load/create
}

// overwriteConfirmedMsg - emitted by the overwrite Decision dialog in the create
// vault flow when the user chooses "Sobrescrever". Routed to activeFlow via modalResult.
type overwriteConfirmedMsg struct{}

func (overwriteConfirmedMsg) isModalResult() {}

// overwriteCancelledMsg - emitted by the overwrite Decision dialog in the create
// vault flow when the user chooses "Voltar". Returns to file picker.
type overwriteCancelledMsg struct{}

func (overwriteCancelledMsg) isModalResult() {}

// weakPwdProceedMsg - emitted by the weak password Decision dialog when the user
// chooses "Prosseguir" (proceed despite weak password). Carries the password.
type weakPwdProceedMsg struct {
	Password []byte
}

func (weakPwdProceedMsg) isModalResult() {}

// weakPwdReviseMsg - emitted by the weak password Decision dialog when the user
// chooses "Revisar" (return to password creation modal).
type weakPwdReviseMsg struct{}

func (weakPwdReviseMsg) isModalResult() {}
