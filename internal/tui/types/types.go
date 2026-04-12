package types

import "io/fs"
import tea "charm.land/bubbletea/v2"
import "time"

type DirectoryLoadedMsg struct {
	Path  string
	Dirs  []fs.DirEntry
	Files []fs.DirEntry
}

// WorkArea represents the main content area of the application. (D-02)
type WorkArea int

const (
	WorkAreaWelcome WorkArea = iota
	WorkAreaVault
	WorkAreaTemplates
	WorkAreaSettings
)

// FilePickerMode defines the mode of the file picker.
type FilePickerMode int

const (
	FilePickerModeOpen FilePickerMode = iota
	FilePickerModeSave
)

// PwdEnteredMsg is sent when a password has been entered by the user.
// The Password field MUST be zeroed after use.
type PwdEnteredMsg struct {
	ID       int // ID of the password modal
	Password []byte
}

func (msg PwdEnteredMsg) Zero() {
	ZeroBytes(msg.Password)
}

// PwdCreatedMsg is sent when a password has been created by the user.
// The Password field MUST be zeroed after use.
type PwdCreatedMsg struct {
	ID       int // ID of the password modal
	Password []byte
}

func (msg PwdCreatedMsg) Zero() {
	ZeroBytes(msg.Password)
}

// ErrorMsg is sent when an error occurs during an async operation.
type ErrorMsg struct {
	ID  int // ID of the component that generated the error
	Err error
}

// ZeroBytes overwrites a byte slice with zeroes.
func ZeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// ModalResult is an empty interface for any message that can be returned by a modal.
type ModalResult interface{ tea.Msg }

// TickMsg is sent by the main tick command. (D-07)
type TickMsg time.Time

// ToggleThemeMsg is sent to toggle between light and dark themes. (D-12)
type ToggleThemeMsg struct{}

// Action represents a keyboard shortcut and its associated behavior.
type Action struct {
	Keys        []string
	Label       string
	Description string
	Group       int
	Scope       Scope
	Priority    int
	HideFromBar bool
	Enabled     func() bool
	Handler     func() tea.Cmd
}

// Scope defines where an action is active.
type Scope int

const (
	ScopeGlobal Scope = iota // Always active, highest priority (e.g., F1 for help)
	ScopeModal               // Active when a modal is open
	ScopeLocal               // Active when current work area is focused
)
