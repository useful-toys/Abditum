package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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

	// Dimensions calculated in View()
	viewportHeight int // Height of content viewport (set in View, used in Update)
}

// compile-time assertion: filePickerModal implements modalView.
var _ modalView = &filePickerModal{}

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
	// Do NOT call adjustTreeScroll() here: m.height is 0 at Init time (SetSize
	// has not been called yet). A scroll computed from height=0 would be wrong
	// and would hide parent nodes above the cursor. SetSize resets treeScroll=0
	// and then adjusts correctly once real dimensions are known.

	// 8. Load initial file list
	m.loadFilesForCursor()

	// 9. Emit initial hint
	hintCmd := m.emitHint()

	return tea.Batch(warnCmd, hintCmd)
}

// buildTreeChain creates a treeNode chain from the filesystem root down to targetPath,
// with each ancestor expanded. Returns the root node.
// Only the path-chain child is expanded at each level; siblings start collapsed.
// Works on both Unix ("/") and Windows ("C:\").
func (m *filePickerModal) buildTreeChain(targetPath string) *treeNode {
	// Determine the filesystem root: "C:\" on Windows, "/" on Unix.
	vol := filepath.VolumeName(targetPath)       // "C:" on Windows, "" on Unix
	rootPath := vol + string(filepath.Separator) // "C:\" or "/"

	root := &treeNode{path: rootPath, name: rootPath, depth: 0, expanded: true, loaded: true}

	// Split the path below the root into components.
	// filepath.ToSlash normalises separators; splitting on "/" works cross-platform.
	// e.g. "C:/g/Abditum" → ["C:", "g", "Abditum"] — skip the volume part.
	// e.g. "/home/user"   → ["", "home", "user"]    — skip the empty first element.
	parts := strings.Split(filepath.ToSlash(targetPath), "/")
	current := root
	currentBuiltPath := rootPath
	depth := 1
	for _, part := range parts {
		if part == "" || part == vol {
			continue // skip empty segment (Unix root split) and Windows drive letter
		}
		childPath := filepath.Join(currentBuiltPath, part)
		// Determine hasSubdirs for the current node
		entries, _ := os.ReadDir(current.path)
		hasSubdirs := false
		for _, e := range entries {
			if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
				hasSubdirs = true
				break
			}
		}
		current.hasSubdirs = hasSubdirs
		// Create the path-chain child: intermediates are expanded, the leaf is not.
		// Mark intermediates as loaded=true so buildVisibleNodes does not call expandNode
		// on them and destroy the single-child chain we just built (D-01).
		isLeaf := childPath == targetPath
		child := &treeNode{
			path:     childPath,
			name:     part,
			depth:    depth,
			expanded: !isLeaf,
			loaded:   !isLeaf, // intermediates: children already set; leaf: lazy-load on expand
		}
		current.children = []*treeNode{child}
		current.loaded = true
		current = child
		currentBuiltPath = childPath
		depth++
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
			text = "Navegue pelas pastas e selecione um cofre"
		} else {
			text = "Navegue pelas pastas e escolha onde salvar"
		}
	case 1: // files panel
		if m.mode == FilePickerOpen {
			if len(m.files) > 0 {
				text = "Selecione o cofre para abrir"
			} else {
				text = "Nenhum cofre neste diretório — navegue para outra pasta"
			}
		} else {
			text = "Arquivos existentes neste diretório"
		}
	case 2: // campo nome (Save mode only)
		if m.nameField.Value() == "" {
			text = "Digite o nome do arquivo — " + m.ext + " será adicionado automaticamente"
		} else {
			text = "Confirme para salvar o cofre"
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
	// viewportHeight is calculated and stored in View(), used here in Update()
	return m.viewportHeight
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
		// Only move to the files panel if there are files to interact with.
		// If the panel is empty, skip it: in Open mode stay in tree; in Save mode go to campo nome.
		if len(m.files) > 0 {
			m.focusPanel = 1
		} else if m.mode == FilePickerSave {
			m.focusPanel = 2
			m.nameField.Focus()
		}
		// else: Open mode + empty dir → no-op, stay in tree
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

// padRight pads s to exactly width visual columns (ANSI-aware via lipgloss.Width).
func padRight(s string, width int) string {
	w := lipgloss.Width(s)
	if w >= width {
		return s
	}
	return s + strings.Repeat(" ", width-w)
}

// renderTopBorder draws the top rounded border with the modal title centered (D-08, D-09).
func (m *filePickerModal) renderTopBorder(modalW int, theme *Theme) string {
	borderSt := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorBorderFocused))
	titleSt := lipgloss.NewStyle().Foreground(theme.TextPrimary).Bold(true)

	title := " " + m.title + " "
	titleRendered := titleSt.Render(title)
	prefixDashes := "── "
	suffixPrefix := " "
	used := lipgloss.Width(prefixDashes) + lipgloss.Width(title) + lipgloss.Width(suffixPrefix)
	remaining := modalW - 2 - used
	if remaining < 1 {
		remaining = 1
	}
	dashes := strings.Repeat("─", remaining)
	return borderSt.Render("╭") + borderSt.Render(prefixDashes) + titleRendered +
		borderSt.Render(suffixPrefix+dashes+"╮")
}

// renderCaminhoHeader draws the "/path/to/dir" row with 1-space lateral padding (D-20).
func (m *filePickerModal) renderCaminhoHeader(innerW int, theme *Theme) string {
	borderSt := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorBorderFocused))
	valueSt := lipgloss.NewStyle().Foreground(theme.TextPrimary)

	// 1 space of padding on each side; path uses the remaining space
	padding := 1
	avail := innerW - padding*2
	if avail < 3 {
		avail = 3
	}
	path := m.currentPath
	// Truncate path from left with … if too wide
	runes := []rune(path)
	for len(runes) > avail {
		runes = runes[1:]
	}
	if string(runes) != path {
		path = "…" + string(runes)
	}
	content := strings.Repeat(" ", padding) + valueSt.Render(path)
	content = padRight(content, innerW)
	return borderSt.Render("│") + content + borderSt.Render("│")
}

// renderPanelSeparator draws the ├── Estrutura ──┬── Arquivos ──┤ line (D-08, D-09).
func (m *filePickerModal) renderPanelSeparator(innerW, treeW, filesW int, theme *Theme) string {
	sepSt := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorBorderFocused))
	headerSt := lipgloss.NewStyle().Foreground(theme.TextSecondary).Bold(true)
	borderSt := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorBorderFocused))

	treeLabel := " Estrutura "
	treeLabelW := lipgloss.Width(treeLabel)
	treeDashes := treeW - treeLabelW
	if treeDashes < 0 {
		treeDashes = 0
	}
	treeLeft := strings.Repeat("─", treeDashes/2)
	treeRight := strings.Repeat("─", treeDashes-treeDashes/2)

	filesLabel := " Arquivos "
	filesLabelW := lipgloss.Width(filesLabel)
	filesDashes := filesW - filesLabelW
	if filesDashes < 0 {
		filesDashes = 0
	}
	filesLeft := strings.Repeat("─", filesDashes/2)
	filesRight := strings.Repeat("─", filesDashes-filesDashes/2)

	return borderSt.Render("├") +
		sepSt.Render(treeLeft) + headerSt.Render(treeLabel) + sepSt.Render(treeRight) +
		sepSt.Render("┬") +
		sepSt.Render(filesLeft) + headerSt.Render(filesLabel) + sepSt.Render(filesRight) +
		borderSt.Render("┤")
}

// renderTreeSepChar returns the character for the tree│files separator column at row i.
// Replaces │ with ↑/■/↓ scroll indicators when tree content overflows (D-08).
func renderTreeSepChar(scroll, total, visibleH, row int) string {
	sepSt := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorBorderFocused))
	indSt := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTextSecondary))
	if total <= visibleH {
		return sepSt.Render("│")
	}
	if row == 0 && scroll > 0 {
		return indSt.Render("↑")
	}
	if row == visibleH-1 && scroll+visibleH < total {
		return indSt.Render("↓")
	}
	thumbPos := 0
	if total-visibleH > 0 {
		thumbPos = (scroll * (visibleH - 2)) / (total - visibleH)
		if thumbPos < 0 {
			thumbPos = 0
		}
		if thumbPos > visibleH-3 {
			thumbPos = visibleH - 3
		}
	}
	if row == thumbPos+1 {
		return indSt.Render("■")
	}
	return sepSt.Render("│")
}

// renderFileSepChar returns the right modal border char at row i.
// Replaces │ with ↑/■/↓ scroll indicators when files content overflows (D-08).
func renderFileSepChar(scroll, total, visibleH, row int) string {
	borderSt := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorBorderFocused))
	indSt := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTextSecondary))
	if total <= visibleH {
		return borderSt.Render("│")
	}
	if row == 0 && scroll > 0 {
		return indSt.Render("↑")
	}
	if row == visibleH-1 && scroll+visibleH < total {
		return indSt.Render("↓")
	}
	thumbPos := 0
	if total-visibleH > 0 {
		thumbPos = (scroll * (visibleH - 2)) / (total - visibleH)
		if thumbPos < 0 {
			thumbPos = 0
		}
		if thumbPos > visibleH-3 {
			thumbPos = visibleH - 3
		}
	}
	if row == thumbPos+1 {
		return indSt.Render("■")
	}
	return borderSt.Render("│")
}

// renderTreeContent returns visibleH lines for the tree (Estrutura) panel (D-01, D-09, D-11).
func (m *filePickerModal) renderTreeContent(treeW, visibleH int, theme *Theme) []string {
	// Highlight the selected dir only when the tree panel has focus (D-09).
	// When another panel is focused the tree shows a passive cursor (bold + accent, no bg).
	var selectedSt lipgloss.Style
	if m.focusPanel == 0 {
		selectedBg := lipgloss.Color("#3d59a1") // special.highlight
		selectedSt = lipgloss.NewStyle().Background(selectedBg).Foreground(theme.AccentPrimary).Bold(true)
	} else {
		selectedSt = lipgloss.NewStyle().Foreground(theme.AccentPrimary).Bold(true)
	}
	normalSt := lipgloss.NewStyle().Foreground(theme.TextPrimary)
	indicatorSt := lipgloss.NewStyle().Foreground(theme.AccentSecondary)

	var lines []string
	end := m.treeScroll + visibleH
	if end > len(m.visibleNodes) {
		end = len(m.visibleNodes)
	}
	for i := m.treeScroll; i < end; i++ {
		node := m.visibleNodes[i].node
		indent := strings.Repeat("  ", node.depth)

		var indicator string
		switch {
		case node.depth == 0:
			indicator = "  " // root — no indicator
		case !node.hasSubdirs:
			indicator = "▷ "
		case node.expanded:
			indicator = "▼ "
		default:
			indicator = "▶ "
		}

		nameText := indent + node.name
		indicatorW := lipgloss.Width(indicator)
		maxNameW := treeW - indicatorW
		if maxNameW < 1 {
			maxNameW = 1
		}
		if lipgloss.Width(nameText) > maxNameW {
			nameText = nameText[:maxNameW-1] + "…"
		}

		indRendered := indicatorSt.Render(indicator)
		var line string
		if i == m.treeCursor {
			// Fill remaining space so background highlight covers full panel width (D-09)
			rowText := padRight(nameText, treeW-indicatorW)
			line = indRendered + selectedSt.Render(rowText)
		} else {
			nameRendered := normalSt.Render(nameText)
			line = indRendered + nameRendered
		}
		lines = append(lines, line)
	}
	for len(lines) < visibleH {
		lines = append(lines, "")
	}
	return lines
}

// renderFilesContent returns visibleH lines for the files (Arquivos) panel (D-09, D-11, D-15).
func (m *filePickerModal) renderFilesContent(filesW, visibleH int, theme *Theme) []string {
	normalSt := lipgloss.NewStyle().Foreground(theme.TextPrimary)
	bulletSt := lipgloss.NewStyle().Foreground(theme.TextSecondary)
	metaSt := lipgloss.NewStyle().Foreground(theme.TextSecondary)
	disabledSt := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTextSecondary))
	// Highlight the selected file only when the files panel has focus (D-09).
	// When the tree panel is focused the files panel shows a passive cursor (bold only, no bg).
	var selectedSt lipgloss.Style
	if m.focusPanel == 1 {
		selectedBg := lipgloss.Color("#3d59a1") // special.highlight
		selectedSt = lipgloss.NewStyle().Background(selectedBg).Foreground(theme.TextPrimary).Bold(true)
	} else {
		selectedSt = lipgloss.NewStyle().Foreground(theme.TextPrimary).Bold(true)
	}

	if len(m.files) == 0 {
		msg := "Nenhum cofre neste diretório"
		msgW := lipgloss.Width(msg)
		pad := (filesW - msgW) / 2
		if pad < 0 {
			pad = 0
		}
		emptyLine := strings.Repeat(" ", pad) + disabledSt.Render(msg)
		var lines []string
		emptyRow := visibleH / 2
		for i := 0; i < visibleH; i++ {
			if i == emptyRow {
				lines = append(lines, padRight(emptyLine, filesW))
			} else {
				lines = append(lines, "")
			}
		}
		return lines
	}

	// Column layout: ● name<nameW>  size<8>  date<14>
	const dateW = 14
	const sizeW = 8
	const colSep = 2
	nameW := filesW - 1 - 1 - sizeW - colSep - dateW - colSep // 1=bullet 1=space
	if nameW < 4 {
		nameW = 4
	}

	var lines []string
	end := m.fileScroll + visibleH
	if end > len(m.files) {
		end = len(m.files)
	}
	for i := m.fileScroll; i < end; i++ {
		name := m.files[i]
		if lipgloss.Width(name) > nameW {
			name = name[:nameW-1] + "…"
		}

		var sizePart, datePart string
		if i < len(m.fileInfos) && m.fileInfos[i] != nil {
			sizePart = fmt.Sprintf("%*s", sizeW, formatFileSize(m.fileInfos[i].Size()))
			datePart = m.timeFmt(m.fileInfos[i].ModTime())
		}

		namePadded := padRight(name, nameW)
		if i == m.fileCursor {
			rowText := " " + namePadded + strings.Repeat(" ", colSep) + sizePart +
				strings.Repeat(" ", colSep) + datePart
			rowText = padRight(rowText, filesW-1)
			line := bulletSt.Render("●") + selectedSt.Render(rowText)
			lines = append(lines, padRight(line, filesW))
		} else {
			bullet := bulletSt.Render("●")
			nameR := normalSt.Render(" " + namePadded)
			sizeR := metaSt.Render(strings.Repeat(" ", colSep) + sizePart)
			dateR := metaSt.Render(strings.Repeat(" ", colSep) + datePart)
			lines = append(lines, padRight(bullet+nameR+sizeR+dateR, filesW))
		}
	}
	for len(lines) < visibleH {
		lines = append(lines, "")
	}
	return lines
}

// renderFieldSeparator draws the ├────┴────┤ separator above the Save mode campo nome (D-08).
func (m *filePickerModal) renderFieldSeparator(innerW, treeW int) string {
	borderSt := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorBorderFocused))
	sepSt := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorBorderFocused))

	leftDashes := strings.Repeat("─", treeW)
	// rightLen: innerW total = treeW(─) + 1(┴) + rightLen(─)
	// so rightLen = innerW - treeW - 1
	// full row: 1(├) + treeW + 1(┴) + rightLen + 1(┤) = innerW + 2 = modalW ✓
	rightLen := innerW - treeW - 1
	if rightLen < 0 {
		rightLen = 0
	}
	rightDashes := strings.Repeat("─", rightLen)
	return borderSt.Render("├") + sepSt.Render(leftDashes) + sepSt.Render("┴") +
		sepSt.Render(rightDashes) + borderSt.Render("┤")
}

// renderFieldRow draws the filename textinput row for Save mode (D-12, D-09).
func (m *filePickerModal) renderFieldRow(innerW int, theme *Theme) string {
	borderSt := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorBorderFocused))
	fieldBgSt := lipgloss.NewStyle().Background(lipgloss.Color(ColorSurfaceInput))

	isFocused := m.focusPanel == 2
	var labelSt lipgloss.Style
	if isFocused {
		labelSt = lipgloss.NewStyle().Foreground(theme.AccentPrimary).Bold(true)
	} else {
		labelSt = lipgloss.NewStyle().Foreground(theme.TextSecondary)
	}

	prefix := " " // 1 space padding at left edge
	label := labelSt.Render("Arquivo:") + " "
	labelW := lipgloss.Width("Arquivo: ")
	fieldW := innerW - len(prefix) - labelW
	if fieldW < 4 {
		fieldW = 4
	}

	val := m.nameField.Value()
	var fieldContent string
	if isFocused {
		cursor := lipgloss.NewStyle().Foreground(theme.TextPrimary).Render("▌")
		fieldContent = val + cursor
	} else {
		fieldContent = val
	}
	if lipgloss.Width(fieldContent) > fieldW {
		// keep rightmost part
		runes := []rune(fieldContent)
		for lipgloss.Width(string(runes)) > fieldW && len(runes) > 1 {
			runes = runes[1:]
		}
		fieldContent = "…" + string(runes)
	}
	fieldRendered := fieldBgSt.Render(padRight(fieldContent, fieldW))

	content := prefix + label + fieldRendered
	content = padRight(content, innerW)
	return borderSt.Render("│") + content + borderSt.Render("│")
}

// renderBottomBorder draws the bottom rounded border with action text state (D-08, D-09).
func (m *filePickerModal) renderBottomBorder(innerW, treeW int, theme *Theme) string {
	borderSt := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorBorderFocused))
	sepSt := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorBorderFocused))
	activeSt := lipgloss.NewStyle().Foreground(theme.AccentPrimary).Bold(true)
	disabledSt := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTextDisabled))
	cancelSt := lipgloss.NewStyle().Foreground(theme.TextPrimary)

	var actionActive bool
	if m.mode == FilePickerOpen {
		actionActive = len(m.files) > 0
	} else {
		actionActive = m.nameField.Value() != ""
	}

	var actionLabel string
	if m.mode == FilePickerOpen {
		actionLabel = "Enter Abrir"
	} else {
		actionLabel = "Enter Salvar"
	}
	var actionRendered string
	if actionActive {
		actionRendered = " " + activeSt.Render(actionLabel) + " "
	} else {
		actionRendered = " " + disabledSt.Render(actionLabel) + " "
	}
	cancelRendered := " " + cancelSt.Render("Esc Cancelar") + " "

	actionW := lipgloss.Width(actionRendered)
	cancelW := lipgloss.Width(cancelRendered)
	// innerW used for fill: subtract action, cancel, and the two "── " side pads (3 each = 6 total)
	// Note: ╰ and ╯ are outside innerW, so full bottom width = 1 + 3 + actionW + fillW + cancelW + 3 + 1 = modalW
	// => fillW = innerW + 2 - actionW - cancelW - 8 = innerW - actionW - cancelW - 6
	fillW := innerW - actionW - cancelW - 6
	if fillW < 0 {
		fillW = 0
	}

	var mid string
	if m.mode == FilePickerOpen && fillW > 0 {
		// Insert ┴ at treeW position in the fill to close the panel separator.
		// The │ separator sits at column (1 + treeW) from the left edge of the modal.
		// The mid region starts at column (1 + 3 + actionW).
		// So leftFill = (1 + treeW) - (1 + 3 + actionW) = treeW - 3 - actionW.
		leftFill := treeW - 3 - actionW
		if leftFill < 0 {
			leftFill = 0
		}
		rightFill := fillW - leftFill - 1
		if rightFill < 0 {
			rightFill = 0
			leftFill = fillW - 1
			if leftFill < 0 {
				leftFill = 0
			}
		}
		mid = sepSt.Render(strings.Repeat("─", leftFill) + "┴" + strings.Repeat("─", rightFill))
	} else {
		mid = sepSt.Render(strings.Repeat("─", fillW))
	}

	return borderSt.Render("╰") +
		sepSt.Render("── ") +
		actionRendered +
		mid +
		cancelRendered +
		sepSt.Render(" ──") +
		borderSt.Render("╯")
}

// View renders the spec-accurate two-panel file picker modal (D-08, D-09, D-20).
func (m *filePickerModal) View(maxWidth, maxHeight int) string {
	if maxWidth == 0 || maxHeight == 0 {
		panic(fmt.Sprintf("filePickerModal.View() called without maxWidth/maxHeight: maxWidth=%d maxHeight=%d", maxWidth, maxHeight))
	}

	// Calculate and store viewportHeight for use in Update()
	modalH := maxHeight * 8 / 10
	if modalH < 6 {
		modalH = 6
	}
	visibleH := modalH - 4 // top border + Caminho + panel sep + bottom border
	if m.mode == FilePickerSave {
		visibleH -= 3 // field separator + field row + bottom padding
	}
	if visibleH < 1 {
		visibleH = 1
	}

	// First render: reset scroll so cursor is visible with real viewport height.
	// (Previously done in SetAvailableSize; safe here because rootModel guarantees
	// View() is called before any Update() keypress.)
	// Check if this is the first render (viewportHeight was 0 before assignment)
	isFirstRender := m.viewportHeight == 0

	// Save viewport height for Update() (pgup/pgdown scroll calculations)
	m.viewportHeight = visibleH

	if isFirstRender {
		m.treeScroll = 0
		m.adjustTreeScroll()
	}

	// Nil-safe theme fallback (D-17)
	theme := m.theme
	if theme == nil {
		theme = ThemeTokyoNight
	}

	// Layout dimensions (D-08)
	modalW := maxWidth * 95 / 100
	if modalW < 60 {
		modalW = 60
	}
	innerW := modalW - 2
	treeW := innerW * 40 / 100
	if treeW < 8 {
		treeW = 8
	}
	filesW := innerW - treeW - 1
	if filesW < 8 {
		filesW = 8
	}

	var lines []string

	// 1. Top border
	lines = append(lines, m.renderTopBorder(modalW, theme))

	// 2. Caminho header
	lines = append(lines, m.renderCaminhoHeader(innerW, theme))

	// 3. Panel separator ├─ Estrutura ─┬─ Arquivos ─┤
	lines = append(lines, m.renderPanelSeparator(innerW, treeW, filesW, theme))

	// 4. Panel content rows
	treeLines := m.renderTreeContent(treeW, visibleH, theme)
	filesLines := m.renderFilesContent(filesW, visibleH, theme)
	borderSt := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorBorderFocused))
	for i := 0; i < visibleH; i++ {
		tl := ""
		fl := ""
		if i < len(treeLines) {
			tl = treeLines[i]
		}
		if i < len(filesLines) {
			fl = filesLines[i]
		}
		tl = padRight(tl, treeW)
		fl = padRight(fl, filesW)
		sep := renderTreeSepChar(m.treeScroll, len(m.visibleNodes), visibleH, i)
		rightBorder := renderFileSepChar(m.fileScroll, len(m.files), visibleH, i)
		lines = append(lines, borderSt.Render("│")+tl+sep+fl+rightBorder)
	}

	// 5. Save mode campo nome section
	if m.mode == FilePickerSave {
		lines = append(lines, m.renderFieldSeparator(innerW, treeW))
		lines = append(lines, m.renderFieldRow(innerW, theme))
	}

	// 6. Bottom border
	lines = append(lines, m.renderBottomBorder(innerW, treeW, theme))

	return strings.Join(lines, "\n")
}
