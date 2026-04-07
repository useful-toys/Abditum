package tui

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
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
// Sequence:
//  1. Set focusPanel = 0 (tree, D-06)
//  2. Init timeFmt injection (D-07)
//  3. Init nameField for Save mode (D-12)
//  4. Resolve currentPath with fallback chain (D-14)
//  5. Build tree from "/" down to currentPath (D-01)
//  6. Flatten into visibleNodes
//  7. Position treeCursor on currentPath node
//  8. Load initial file list (D-15)
//  9. Emit initial hint (D-03)
func (m *filePickerModal) Init() tea.Cmd {
	// 1. Set initial focus to tree (D-06)
	m.focusPanel = 0

	// 2. Init timeFmt (D-07 injection point)
	if m.timeFmt == nil {
		m.timeFmt = func(t time.Time) string {
			return t.Local().Format("02/01/06 15:04")
		}
	}

	// 3. Init nameField for Save mode (D-12)
	m.nameField = textinput.New()
	m.nameField.Placeholder = ""
	m.nameField.Blur()
	if m.mode == FilePickerSave && m.suggestedName != "" {
		m.nameField.SetValue(m.suggestedName)
	}

	// 4. Resolve currentPath with fallback chain (D-14)
	var warnCmd tea.Cmd
	cwd, err := os.Getwd()
	if err == nil {
		if _, err2 := os.ReadDir(cwd); err2 == nil {
			m.currentPath = cwd
		} else {
			err = err2
		}
	}
	if err != nil {
		home, herr := os.UserHomeDir()
		if herr != nil {
			home = "/"
		}
		m.currentPath = home
		warnCmd = func() tea.Msg {
			if m.messages != nil {
				m.messages.Show(MsgWarn, "⚠ Diretório atual inacessível — navegando para home", 0, true)
			}
			return nil
		}
	}

	// 5. Build tree from "/" down to currentPath (D-01)
	m.root = m.buildTreeChain(m.currentPath)

	// 6. Flatten into visibleNodes
	m.visibleNodes = nil
	m.buildVisibleNodes(m.root, &m.visibleNodes)

	// 7. Set treeCursor to currentPath node
	for i, vn := range m.visibleNodes {
		if vn.node.path == m.currentPath {
			m.treeCursor = i
			break
		}
	}

	// 8. Load initial file list
	m.loadFilesForCursor()

	// 9. Emit initial hint
	hintCmd := m.emitHint()

	return tea.Batch(warnCmd, hintCmd)
}

// buildTreeChain creates a treeNode chain from root "/" down to targetPath,
// with each ancestor expanded. Returns the root node.
// Only the path-chain child is expanded at each level; siblings start collapsed.
func (m *filePickerModal) buildTreeChain(targetPath string) *treeNode {
	// Split targetPath into components: e.g. "/home/user/docs" → ["", "home", "user", "docs"]
	parts := strings.Split(filepath.ToSlash(targetPath), "/")
	// Build chain top-down starting at root "/"
	root := &treeNode{path: "/", name: "/", depth: 0, expanded: true, loaded: true}
	current := root
	currentBuiltPath := "/"
	for i := 1; i < len(parts); i++ {
		if parts[i] == "" {
			continue
		}
		childPath := filepath.Join(currentBuiltPath, parts[i])
		// Load children of current node to determine hasSubdirs
		entries, _ := os.ReadDir(current.path)
		hasSubdirs := false
		for _, e := range entries {
			if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
				hasSubdirs = true
				break
			}
		}
		current.hasSubdirs = hasSubdirs
		// Create the path-chain child: expand intermediates; leaf starts collapsed
		child := &treeNode{
			path:     childPath,
			name:     parts[i],
			depth:    i,
			expanded: i < len(parts)-1, // intermediate nodes are expanded; leaf is not
			loaded:   false,
		}
		current.children = []*treeNode{child}
		current.loaded = true
		current = child
		currentBuiltPath = childPath
	}
	return root
}

// buildVisibleNodes appends visible nodes to out via DFS.
// Expanded nodes have their children traversed; collapsed nodes are leaves in the flat list.
func (m *filePickerModal) buildVisibleNodes(node *treeNode, out *[]visibleNode) {
	if node == nil {
		return
	}
	*out = append(*out, visibleNode{node: node})
	if node.expanded {
		// Load children lazily if not yet loaded
		if !node.loaded {
			_ = m.expandNode(node) // errors handled inside expandNode
		}
		for _, child := range node.children {
			m.buildVisibleNodes(child, out)
		}
	}
}

// expandNode loads node.children from disk and sets node.expanded=true.
// Returns non-nil error if os.ReadDir fails (e.g. permission denied, D-22).
// On error: node remains collapsed (expanded=false), loaded=true (won't retry).
func (m *filePickerModal) expandNode(node *treeNode) error {
	entries, err := os.ReadDir(node.path)
	if err != nil {
		node.loaded = true
		node.expanded = false
		return err
	}
	node.loaded = true
	node.expanded = true
	node.children = node.children[:0] // reset children slice
	hasSubdirs := false
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue // skip hidden directories (D-13)
		}
		if e.IsDir() {
			hasSubdirs = true
			childPath := filepath.Join(node.path, e.Name())
			node.children = append(node.children, &treeNode{
				path:  childPath,
				name:  e.Name(),
				depth: node.depth + 1,
			})
		}
	}
	node.hasSubdirs = hasSubdirs
	// Sort children alphabetically case-insensitive (D-11)
	sort.Slice(node.children, func(i, j int) bool {
		return strings.ToLower(node.children[i].name) < strings.ToLower(node.children[j].name)
	})
	return nil
}

// loadFilesForCursor reads m.currentPath and populates m.files and m.fileInfos.
// Applies extension filter (m.ext) and hidden-file exclusion (D-13).
// Sets fileCursor per D-15: 0 if files exist, -1 if empty.
func (m *filePickerModal) loadFilesForCursor() {
	m.files = nil
	m.fileInfos = nil
	m.fileScroll = 0

	entries, err := os.ReadDir(m.currentPath)
	if err != nil {
		m.fileCursor = -1
		return
	}

	type fileEntry struct {
		name string
		info os.FileInfo
	}
	var found []fileEntry
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasPrefix(e.Name(), ".") {
			continue // skip hidden files (D-13)
		}
		if !strings.HasSuffix(e.Name(), m.ext) {
			continue // only show files matching extension (D-13)
		}
		info, _ := e.Info()
		// Store name without extension for display (D-13)
		found = append(found, fileEntry{
			name: strings.TrimSuffix(e.Name(), m.ext),
			info: info,
		})
	}
	// Sort alphabetically case-insensitive (D-11)
	sort.Slice(found, func(i, j int) bool {
		return strings.ToLower(found[i].name) < strings.ToLower(found[j].name)
	})

	for _, f := range found {
		m.files = append(m.files, f.name)
		m.fileInfos = append(m.fileInfos, f.info)
	}

	// Auto-selection D-15: 0 if files exist, -1 if empty
	if len(m.files) > 0 {
		m.fileCursor = 0
	} else {
		m.fileCursor = -1
	}
}

// emitHint returns a Cmd that shows the appropriate hint message for the current
// focus/state combination (D-03). Returns nil if m.messages is nil.
func (m *filePickerModal) emitHint() tea.Cmd {
	if m.messages == nil {
		return nil
	}
	var text string
	switch m.focusPanel {
	case 0: // tree panel
		if m.mode == FilePickerOpen {
			text = "• Navegue pelas pastas e selecione um cofre"
		} else {
			text = "• Navegue pelas pastas e escolha onde salvar"
		}
	case 1: // files panel
		if m.mode == FilePickerOpen {
			if len(m.files) > 0 {
				text = "• Selecione o cofre para abrir"
			} else {
				text = "• Nenhum cofre neste diretório — navegue para outra pasta"
			}
		} else {
			text = "• Arquivos existentes neste diretório"
		}
	case 2: // campo nome (Save mode only)
		if m.nameField.Value() == "" {
			text = "• Digite o nome do arquivo — " + m.ext + " será adicionado automaticamente"
		} else {
			text = "• Confirme para salvar o cofre"
		}
	}
	m.messages.Show(MsgHint, text, 0, true)
	return nil
}

// Update handles keyboard and async messages.
// Dispatches to panel-specific handlers based on m.focusPanel.
// Esc always cancels regardless of panel (D-19).
// Implemented in Task 2.
func (m *filePickerModal) Update(msg tea.Msg) tea.Cmd {
	return nil // implemented in Task 2
}

// View renders the two-panel modal box.
// Implemented in Plan 03.
func (m *filePickerModal) View() string {
	return "" // implemented in Plan 03
}
