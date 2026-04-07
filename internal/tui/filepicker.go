package tui

import (
	"os"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// FilePickerMode controls which picker behavior is active.
type FilePickerMode int

const (
	FilePickerOpen FilePickerMode = iota // open an existing file
	FilePickerSave                       // save / name a new file
)

// treeNode represents a directory entry in the lazy recursive tree (D-01).
type treeNode struct {
	path       string
	name       string
	depth      int
	expanded   bool
	loaded     bool
	children   []*treeNode
	hasSubdirs bool // true if dir has at least one subdirectory
}

// visibleNode is a flattened tree entry visible in the Estrutura panel.
type visibleNode struct {
	node *treeNode
}

// filePickerModal is the spec-compliant two-panel file picker (D-00).
// It satisfies the modalView interface defined in modal.go.
type filePickerModal struct {
	// Construction-time fields (set by factory before Init)
	mode     FilePickerMode
	ext      string // required: e.g. ".abditum" (D-13)
	title    string
	messages *MessageManager
	theme    *Theme

	// Tree state (D-01)
	root         *treeNode
	visibleNodes []visibleNode
	treeCursor   int
	treeScroll   int

	// Current directory and file list (D-13, D-15)
	currentPath string
	files       []string      // names without ext, sorted case-insensitive
	fileInfos   []os.FileInfo // parallel to files — size+mtime metadata
	fileCursor  int           // -1 when no files in currentPath (D-15)
	fileScroll  int

	// Focus (D-06)
	focusPanel int // 0=tree, 1=files, 2=campo nome (Save mode only)

	// Save mode filename field (D-12)
	nameField     textinput.Model
	suggestedName string // pre-fill if non-empty at Init (D-14)

	// Test injection (D-07)
	timeFmt func(time.Time) string // if nil, defaults to Local "02/01/06 15:04"

	// Dimensions (set by rootModel via SetSize)
	width  int
	height int
}

// compile-time assertion: filePickerModal implements modalView.
var _ modalView = &filePickerModal{}

// SetSize stores the terminal dimensions for use in View().
func (m *filePickerModal) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Shortcuts returns the keyboard shortcuts shown in the command bar (D-18).
func (m *filePickerModal) Shortcuts() []Shortcut {
	return []Shortcut{
		{Key: "Tab", Label: "Painel"},
		{Key: "F1", Label: "Ajuda"},
	}
}

// Init initialises file-picker state. Called by the factory before push (D-14).
// Implemented in Plan 02.
func (m *filePickerModal) Init() tea.Cmd {
	return nil // implemented in Plan 02
}

// Update handles keyboard and async messages.
// Implemented in Plan 02.
func (m *filePickerModal) Update(msg tea.Msg) tea.Cmd {
	return nil // implemented in Plan 02
}

// View renders the two-panel modal box.
// Implemented in Plan 03.
func (m *filePickerModal) View() string {
	return "" // implemented in Plan 03
}
