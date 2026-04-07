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

// currentNode returns the treeNode under treeCursor, or nil if visibleNodes is empty.
func (m *filePickerModal) currentNode() *treeNode {
	if m.treeCursor < 0 || m.treeCursor >= len(m.visibleNodes) {
		return nil
	}
	return m.visibleNodes[m.treeCursor].node
}

// parentCursor returns the visibleNodes index of the parent of node, or -1 if not found (D-16).
func (m *filePickerModal) parentCursor(node *treeNode) int {
	if node == nil || node.depth == 0 {
		return -1
	}
	targetDepth := node.depth - 1
	targetPath := filepath.Dir(node.path)
	for i, vn := range m.visibleNodes {
		if vn.node.depth == targetDepth && vn.node.path == targetPath {
			return i
		}
	}
	return -1
}

// visibleTreeHeight returns the number of rows available for tree/file rows inside the modal.
func (m *filePickerModal) visibleTreeHeight() int {
	modalH := m.height * 8 / 10
	// subtract: top border(1) + Caminho header(1) + panel separator(1) + bottom border(1) = 4
	h := modalH - 4
	if m.mode == FilePickerSave {
		h -= 3 // campo nome section
	}
	if h < 1 {
		h = 1
	}
	return h
}

func (m *filePickerModal) visibleFilesHeight() int {
	return m.visibleTreeHeight()
}

// adjustTreeScroll keeps treeCursor visible in the tree panel viewport.
func (m *filePickerModal) adjustTreeScroll() {
	h := m.visibleTreeHeight()
	if m.treeCursor < m.treeScroll {
		m.treeScroll = m.treeCursor
	}
	if m.treeCursor >= m.treeScroll+h {
		m.treeScroll = m.treeCursor - h + 1
	}
}

// adjustFileScroll keeps fileCursor visible in the files panel viewport.
func (m *filePickerModal) adjustFileScroll() {
	h := m.visibleFilesHeight()
	if m.fileCursor < m.fileScroll {
		m.fileScroll = m.fileCursor
	}
	if m.fileCursor >= m.fileScroll+h {
		m.fileScroll = m.fileCursor - h + 1
	}
}

// updateTree handles keyboard input for the tree panel (D-10, D-16, D-22).
func (m *filePickerModal) updateTree(msg tea.KeyPressMsg) tea.Cmd {
	node := m.currentNode()
	h := m.visibleTreeHeight()

	switch msg.Code {
	case tea.KeyDown:
		if m.treeCursor < len(m.visibleNodes)-1 {
			m.treeCursor++
			if n := m.currentNode(); n != nil {
				m.currentPath = n.path
			}
			m.loadFilesForCursor()
		}
	case tea.KeyUp:
		if m.treeCursor > 0 {
			m.treeCursor--
			if n := m.currentNode(); n != nil {
				m.currentPath = n.path
			}
			m.loadFilesForCursor()
		}
	case tea.KeyHome:
		m.treeCursor = 0
		if n := m.currentNode(); n != nil {
			m.currentPath = n.path
		}
		m.loadFilesForCursor()
	case tea.KeyEnd:
		if len(m.visibleNodes) > 0 {
			m.treeCursor = len(m.visibleNodes) - 1
			if n := m.currentNode(); n != nil {
				m.currentPath = n.path
			}
			m.loadFilesForCursor()
		}
	case tea.KeyPgDown:
		m.treeCursor += h
		if m.treeCursor >= len(m.visibleNodes) {
			m.treeCursor = len(m.visibleNodes) - 1
		}
		if n := m.currentNode(); n != nil {
			m.currentPath = n.path
		}
		m.loadFilesForCursor()
	case tea.KeyPgUp:
		m.treeCursor -= h
		if m.treeCursor < 0 {
			m.treeCursor = 0
		}
		if n := m.currentNode(); n != nil {
			m.currentPath = n.path
		}
		m.loadFilesForCursor()
	case tea.KeyRight:
		if node == nil || !node.hasSubdirs {
			return nil // ▷ node or no subdirs — no-op
		}
		if node.expanded {
			return nil // already expanded — no-op
		}
		if err := m.expandNode(node); err != nil {
			// D-22: permission error — show dir basename only
			if m.messages != nil {
				m.messages.Show(MsgError, "✕ Sem permissão para acessar "+filepath.Base(node.path), 5, true)
			}
			return nil
		}
		m.visibleNodes = nil
		m.buildVisibleNodes(m.root, &m.visibleNodes)
		// Re-find treeCursor pointing to the same node
		for i, vn := range m.visibleNodes {
			if vn.node == node {
				m.treeCursor = i
				break
			}
		}
	case tea.KeyLeft:
		if node == nil || node.depth == 0 {
			return nil // root — no-op
		}
		if node.expanded {
			// Collapse the expanded node
			node.expanded = false
			m.visibleNodes = nil
			m.buildVisibleNodes(m.root, &m.visibleNodes)
			for i, vn := range m.visibleNodes {
				if vn.node == node {
					m.treeCursor = i
					break
				}
			}
		} else {
			// D-16: already collapsed → navigate to parent
			if pi := m.parentCursor(node); pi >= 0 {
				m.treeCursor = pi
				if n := m.currentNode(); n != nil {
					m.currentPath = n.path
				}
				m.loadFilesForCursor()
			}
		}
	case tea.KeyEnter:
		if len(m.files) > 0 {
			m.focusPanel = 1
			return m.emitHint()
		}
		// no files → no-op
	case tea.KeyTab:
		m.focusPanel = 1
		return m.emitHint()
	}
	m.adjustTreeScroll()
	return m.emitHint()
}

// updateFiles handles keyboard input for the files panel (D-06, D-10).
func (m *filePickerModal) updateFiles(msg tea.KeyPressMsg) tea.Cmd {
	h := m.visibleFilesHeight()
	switch msg.Code {
	case tea.KeyDown:
		if len(m.files) > 0 && m.fileCursor < len(m.files)-1 {
			m.fileCursor++
			m.adjustFileScroll()
		}
	case tea.KeyUp:
		if m.fileCursor > 0 {
			m.fileCursor--
			m.adjustFileScroll()
		}
	case tea.KeyHome:
		if len(m.files) > 0 {
			m.fileCursor = 0
			m.fileScroll = 0
		}
	case tea.KeyEnd:
		if len(m.files) > 0 {
			m.fileCursor = len(m.files) - 1
			m.adjustFileScroll()
		}
	case tea.KeyPgDown:
		if len(m.files) > 0 {
			m.fileCursor += h
			if m.fileCursor >= len(m.files) {
				m.fileCursor = len(m.files) - 1
			}
			m.adjustFileScroll()
		}
	case tea.KeyPgUp:
		m.fileCursor -= h
		if m.fileCursor < 0 {
			m.fileCursor = 0
		}
		m.adjustFileScroll()
	case tea.KeyEnter:
		if m.fileCursor < 0 || m.fileCursor >= len(m.files) {
			return nil
		}
		if m.mode == FilePickerOpen {
			// D-06: Open — confirm selection immediately
			fullPath := filepath.Join(m.currentPath, m.files[m.fileCursor]+m.ext)
			return tea.Batch(
				func() tea.Msg { return filePickerResult{Path: fullPath} },
				func() tea.Msg { return popModalMsg{} },
			)
		}
		// Save mode: copy selected filename to campo nome, move focus
		m.nameField.SetValue(m.files[m.fileCursor])
		m.focusPanel = 2
		m.nameField.Focus()
		return m.emitHint()
	case tea.KeyTab:
		if m.mode == FilePickerOpen {
			m.focusPanel = 0
		} else {
			m.focusPanel = 2
			m.nameField.Focus()
		}
		return m.emitHint()
	}
	return nil
}

// updateField handles keyboard input for the campo nome text field (D-12).
func (m *filePickerModal) updateField(msg tea.KeyPressMsg) tea.Cmd {
	switch msg.Code {
	case tea.KeyEnter:
		name := m.nameField.Value()
		if name == "" {
			return nil // D-12: no-op when empty
		}
		if !strings.HasSuffix(name, m.ext) {
			name += m.ext
		}
		fullPath := filepath.Join(m.currentPath, name)
		return tea.Batch(
			func() tea.Msg { return filePickerResult{Path: fullPath} },
			func() tea.Msg { return popModalMsg{} },
		)
	case tea.KeyTab:
		m.focusPanel = 0
		m.nameField.Blur()
		return m.emitHint()
	default:
		// Block invalid filesystem chars silently (D-12): / \ : * ? " < > |
		if msg.Text != "" {
			for _, r := range msg.Text {
				if strings.ContainsRune(`/\:*?"<>|`, r) {
					return nil // silently block
				}
			}
		}
		var cmd tea.Cmd
		m.nameField, cmd = m.nameField.Update(msg)
		return tea.Batch(cmd, m.emitHint())
	}
}

// Update handles keyboard and async messages.
// Dispatches to panel-specific handlers based on m.focusPanel.
// Esc always cancels regardless of panel (D-19).
func (m *filePickerModal) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// Global: Esc always cancels (D-19)
		if msg.Code == tea.KeyEsc {
			return tea.Batch(
				func() tea.Msg { return filePickerResult{Cancelled: true, Path: ""} },
				func() tea.Msg { return popModalMsg{} },
			)
		}
		// Dispatch to panel-specific handler
		switch m.focusPanel {
		case 0:
			return m.updateTree(msg)
		case 1:
			return m.updateFiles(msg)
		case 2:
			return m.updateField(msg)
		}
	}
	return nil
}

// View renders the two-panel modal box.
// Implemented in Plan 03.
func (m *filePickerModal) View() string {
	return "" // implemented in Plan 03
}
