package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/storage"
)

// childModel represents a UI component managed by rootModel.
//
// CRITICAL CONTRACT: rootModel guarantees that SetSize(width, height) is ALWAYS called
// immediately before View(). Implementations must panic if View() is called
// with width or height equal to zero — this indicates a bug in rootModel.
//
// Layout and rendering MUST assume width > 0 and height > 0.
// SetSize receives EXACT allocated dimensions that the component MUST occupy.
//
// D-01: 3 methods only. Closures that capture child state directly.
type childModel interface {
	Update(tea.Msg) tea.Cmd
	// View renders the component. MUST only be called after SetSize.
	// If width or height are zero, panic with a descriptive message.
	View() string
	// SetSize stores exact allocated dimensions for layout.
	// Component MUST occupy exactly this space.
	SetSize(width, height int)
	ApplyTheme(*Theme)
}

// Shortcut is a key+label pair displayed in the command bar while a modal is active.
type Shortcut struct {
	Key   string
	Label string
}

// modalView represents an overlay modal dialog.
//
// CRITICAL CONTRACT: rootModel guarantees that SetAvailableSize(maxWidth, maxHeight) is ALWAYS called
// immediately before View(). Implementations must panic if View() is called
// with maxWidth or maxHeight equal to zero — this indicates a bug in rootModel.
//
// Modals auto-size by content but need terminal dimensions to calculate their
// own limits (e.g., 80% of terminal width). rootModel calls SetAvailableSize on every
// resize and before View().
// SetAvailableSize receives MAXIMUM available dimensions; modal MAY use less space (e.g., center and pad).
//
// D-02: Separate from childModel.
type modalView interface {
	Update(tea.Msg) tea.Cmd
	// View renders the modal. MUST only be called after SetAvailableSize.
	// If maxWidth or maxHeight are zero, panic with a descriptive message.
	View() string
	// Shortcuts returns the command bar shortcuts active while this modal is displayed.
	Shortcuts() []Shortcut
	// SetAvailableSize stores maximum available dimensions for layout calculations.
	// Modal MAY use less space (e.g., center and pad).
	SetAvailableSize(maxWidth, maxHeight int)
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
