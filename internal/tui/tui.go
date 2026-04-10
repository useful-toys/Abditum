// Package tui provides the terminal user interface for Abditum.
package tui

import (
	"path/filepath"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/useful-toys/abditum/internal/vault"
)

// Option is a functional option for configuring the RootModel.
type Option func(*RootModel)

// WithVersion sets the application version displayed in the UI.
func WithVersion(v string) Option {
	return func(m *RootModel) {
		m.version = v
	}
}

// WithVaultPath sets the optional initial vault file path.
func WithVaultPath(p string) Option {
	return func(m *RootModel) {
		m.vaultPath = p
	}
}

// panelFocus tracks which panel has keyboard focus.
type panelFocus int

const (
	focusTree panelFocus = iota
	focusDetail
)

// tickMsg triggers spinner animation.
type tickMsg time.Time

// clearMsgMsg clears the message bar after TTL.
type clearMsgMsg struct{}

// RootModel is the top-level Bubble Tea model for the Abditum TUI.
type RootModel struct {
	version   string
	vaultPath string

	// Layout
	width  int
	height int

	// Theme
	theme  Theme
	styles styles

	// Mode
	mode appMode

	// Vault state
	manager   *vault.Manager
	vaultName string // display name (filename without extension)

	// Tree panel
	tree        treeModel
	treeFocused panelFocus

	// Detail panel
	detail detailModel

	// Search
	searchActive bool
	searchQuery  string

	// Message bar
	msg       barMessage
	spinFrame int

	// Active dialog (nil = no dialog)
	activeDialog *dialog
	dialogCb     func(actionIdx int) // callback when user selects a dialog action

	// Actions for current context
	actions []Action
}

// NewRootModel creates a new RootModel with the given options.
func NewRootModel(opts ...Option) *RootModel {
	m := &RootModel{
		theme: ThemeTokyoNight,
		width: 80,
		height: 24,
		mode:  modeWelcome,
		tree:  newTreeModel(),
	}
	m.styles = newStyles(tokyoNight)
	for _, o := range opts {
		o(m)
	}
	m.rebuildActions()
	return m
}

// Init implements tea.Model.
func (m *RootModel) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
	)
}

// Update implements tea.Model.
func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		m.spinFrame++
		// Check TTL
		if m.msg.kind != msgNone && m.msg.kind != msgBusy &&
			m.msg.kind != msgHint && !m.msg.expires.IsZero() &&
			time.Now().After(m.msg.expires) {
			m.msg = noMessage
		}
		return m, tickCmd()

	case clearMsgMsg:
		m.msg = noMessage
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// handleKey dispatches key events.
func (m *RootModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Global keys — always active
	switch key {
	case "ctrl+q":
		return m.handleQuit()
	case "f12":
		m.toggleTheme()
		return m, nil
	case "f1":
		// TODO: open help
		return m, nil
	}

	// Dialog active — route to dialog handler
	if m.activeDialog != nil {
		return m.handleDialogKey(key)
	}

	// Mode-specific keys
	switch m.mode {
	case modeWelcome:
		return m.handleWelcomeKey(key)
	case modeVault:
		return m.handleVaultKey(key)
	}

	return m, nil
}

func (m *RootModel) handleQuit() (tea.Model, tea.Cmd) {
	if m.manager != nil && m.manager.IsModified() {
		// Show confirmation dialog
		m.openDialog(dialog{
			title:    "Sair do Abditum",
			body:     "Cofre modificado. Salvar ou descartar?",
			severity: severityDestructive,
			actions: []dialogAction{
				{Key: "Enter", Label: "Salvar", IsEnter: true},
				{Key: "D", Label: "Descartar"},
				{Key: "Esc", Label: "Voltar", IsESC: true},
			},
		}, func(idx int) {
			// idx 0 = Salvar, 1 = Descartar, 2 = Esc/Voltar
			// For now just quit on any action
		})
		return m, nil
	}
	return m, tea.Quit
}

func (m *RootModel) handleWelcomeKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "f5":
		// TODO: create new vault flow
	case "f6":
		// TODO: open vault flow
	}
	return m, nil
}

func (m *RootModel) handleVaultKey(key string) (tea.Model, tea.Cmd) {
	if m.manager == nil {
		return m, nil
	}

	switch key {
	case "f2":
		// already in vault mode
	case "f3":
		m.switchMode(modeModels)
	case "f4":
		m.switchMode(modeConfig)
	case "f10", "ctrl+f":
		m.toggleSearch()
	case "ctrl+alt+shift+q":
		// Emergency lock — no confirmation
		if m.manager != nil {
			m.manager.Lock()
			m.mode = modeWelcome
			m.vaultName = ""
			m.manager = nil
			m.tree = newTreeModel()
			m.detail = detailModel{}
		}
		return m, nil
	}

	// Search-active input routing
	if m.searchActive {
		return m.handleSearchKey(key)
	}

	// Tree / detail navigation
	if m.treeFocused == focusTree {
		return m.handleTreeKey(key)
	}
	return m.handleDetailKey(key)
}

func (m *RootModel) handleSearchKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc", "ctrl+f", "f10":
		m.searchActive = false
		m.searchQuery = ""
		m.setHint("• ↑↓ para navegar")
		m.rebuildActions()
	case "delete":
		m.searchQuery = ""
		m.setHint("• Digite para filtrar os segredos")
	case "backspace":
		if len(m.searchQuery) > 0 {
			runes := []rune(m.searchQuery)
			m.searchQuery = string(runes[:len(runes)-1])
		}
		if m.searchQuery == "" {
			m.setHint("• Digite para filtrar os segredos")
		}
	case "up":
		m.tree.moveUp()
	case "down":
		m.tree.moveDown()
	case "home":
		m.tree.moveHome()
	case "end":
		m.tree.moveEnd()
	default:
		// printable character → append to query
		if len(key) == 1 && key[0] >= 32 {
			m.searchQuery += key
		}
	}
	m.syncDetailFromTree()
	return m, nil
}

func (m *RootModel) handleTreeKey(key string) (tea.Model, tea.Cmd) {
	cofre := m.manager.Vault()

	switch key {
	case "up":
		m.tree.moveUp()
	case "down":
		m.tree.moveDown()
	case "home":
		m.tree.moveHome()
	case "end":
		m.tree.moveEnd()
	case "right", "→":
		m.tree.expandFolder(cofre)
	case "left", "←":
		m.tree.collapseFolder(cofre)
	case "enter":
		m.tree.expandFolder(cofre)
	case "tab":
		m.treeFocused = focusDetail
		m.detail.focused = true
		m.rebuildActions()
		return m, nil
	case "ctrl+s":
		return m.savekVault()
	case "f7":
		return m.savekVault()
	}

	m.syncDetailFromTree()
	m.rebuildActions()
	return m, nil
}

func (m *RootModel) handleDetailKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "tab":
		m.treeFocused = focusTree
		m.detail.focused = false
		m.setHint("• ↑↓ para navegar")
		m.rebuildActions()
	case "up":
		if m.detail.scrollTop > 0 {
			m.detail.scrollTop--
		}
	case "down":
		m.detail.scrollTop++
	case "ctrl+r":
		// cycle reveal
		switch m.detail.reveal {
		case revealMasked:
			m.detail.reveal = revealHint
		case revealHint:
			m.detail.reveal = revealFull
		case revealFull:
			m.detail.reveal = revealMasked
		}
	case "esc":
		m.treeFocused = focusTree
		m.detail.focused = false
		m.rebuildActions()
	}
	return m, nil
}

func (m *RootModel) handleDialogKey(key string) (tea.Model, tea.Cmd) {
	if m.activeDialog == nil {
		return m, nil
	}
	actions := m.activeDialog.actions
	switch key {
	case "esc":
		for i, a := range actions {
			if a.IsESC {
				m.closeDialog(i)
				return m, nil
			}
		}
		m.closeDialog(-1)
	case "enter":
		for i, a := range actions {
			if a.IsEnter {
				m.closeDialog(i)
				return m, nil
			}
		}
		m.closeDialog(0)
	default:
		for i, a := range actions {
			if strings.EqualFold(a.Key, key) && !a.IsEnter && !a.IsESC {
				m.closeDialog(i)
				return m, nil
			}
		}
	}
	return m, nil
}

func (m *RootModel) openDialog(d dialog, cb func(int)) {
	m.activeDialog = &d
	m.dialogCb = cb
}

func (m *RootModel) closeDialog(actionIdx int) {
	m.activeDialog = nil
	if m.dialogCb != nil && actionIdx >= 0 {
		m.dialogCb(actionIdx)
		m.dialogCb = nil
	}
}

func (m *RootModel) toggleTheme() {
	if m.theme == ThemeTokyoNight {
		m.theme = ThemeCyberpunk
		m.styles = newStyles(cyberpunk)
	} else {
		m.theme = ThemeTokyoNight
		m.styles = newStyles(tokyoNight)
	}
}

func (m *RootModel) toggleSearch() {
	m.searchActive = !m.searchActive
	if m.searchActive {
		m.searchQuery = ""
		m.setHint("• Digite para filtrar os segredos")
	} else {
		m.searchQuery = ""
		m.setHint("• ↑↓ para navegar")
	}
	m.rebuildActions()
}

func (m *RootModel) switchMode(mode appMode) {
	m.mode = mode
	m.rebuildActions()
}

func (m *RootModel) setMsg(kind msgType, text string, ttl time.Duration) {
	expires := time.Time{}
	if ttl > 0 {
		expires = time.Now().Add(ttl)
	}
	m.msg = barMessage{kind: kind, text: text, expires: expires}
}

func (m *RootModel) setHint(text string) {
	m.msg = barMessage{kind: msgHint, text: text}
}

func (m *RootModel) savekVault() (tea.Model, tea.Cmd) {
	if m.manager == nil {
		return m, nil
	}
	// TODO: actual save with file picker if no path
	m.setMsg(msgSuccess, "✓ Cofre salvo", 5*time.Second)
	return m, nil
}

func (m *RootModel) syncDetailFromTree() {
	s := m.tree.selectedSecret()
	if s != m.detail.secret {
		m.detail.secret = s
		m.detail.scrollTop = 0
		m.detail.reveal = revealMasked
	}
}

func (m *RootModel) rebuildActions() {
	st := m.styles
	_ = st

	if m.activeDialog != nil {
		// During a dialog, no main actions shown
		m.actions = nil
		return
	}

	if m.searchActive {
		m.actions = []Action{
			{Key: "⌃F", Label: "Fechar", Priority: 100, Enabled: true},
			{Key: "Del", Label: "Limpar", Priority: 90, Enabled: true},
		}
		return
	}

	switch m.mode {
	case modeWelcome:
		m.actions = []Action{
			{Key: "F5", Label: "Novo cofre", Priority: 100, Enabled: true},
			{Key: "F6", Label: "Abrir cofre", Priority: 90, Enabled: true},
		}
	case modeVault:
		if m.treeFocused == focusTree {
			m.actions = treeActions(m)
		} else {
			m.actions = detailActions(m)
		}
	case modeModels:
		m.actions = []Action{
			{Key: "F2", Label: "Cofre", Priority: 100, Enabled: true},
			{Key: "F4", Label: "Config", Priority: 90, Enabled: true},
		}
	case modeConfig:
		m.actions = []Action{
			{Key: "F2", Label: "Cofre", Priority: 100, Enabled: true},
			{Key: "F3", Label: "Modelos", Priority: 90, Enabled: true},
		}
	}
}

func treeActions(m *RootModel) []Action {
	actions := []Action{
		{Key: "F2", Label: "Cofre", Priority: 200, Enabled: true, HideFromBar: true},
		{Key: "F3", Label: "Modelos", Priority: 190, Enabled: true},
		{Key: "F4", Label: "Config", Priority: 180, Enabled: true},
		{Key: "F7", Label: "Salvar", Priority: 150, Enabled: m.manager != nil && m.manager.IsModified()},
		{Key: "F10", Label: "Busca", Priority: 140, Enabled: true},
	}

	s := m.tree.selectedSecret()
	if s != nil {
		hasSensitive := false
		for _, c := range s.Campos() {
			if c.Tipo() == vault.TipoCampoSensivel {
				hasSensitive = true
				break
			}
		}
		actions = append(actions,
			Action{Key: "Enter", Label: "Detalhes", Priority: 300, Enabled: true},
			Action{Key: "Ins", Label: "Novo", Priority: 250, Enabled: true},
			Action{Key: "Del", Label: "Excluir", Priority: 120, Enabled: true},
		)
		if hasSensitive {
			revealLabel := "Revelar"
			if m.detail.reveal == revealHint {
				revealLabel = "Mostrar tudo"
			} else if m.detail.reveal == revealFull {
				revealLabel = "Ocultar"
			}
			actions = append(actions,
				Action{Key: "⌃R", Label: revealLabel, Priority: 130, Enabled: true},
				Action{Key: "⌃C", Label: "Copiar", Priority: 125, Enabled: true},
			)
		}
	}
	return actions
}

func detailActions(m *RootModel) []Action {
	return []Action{
		{Key: "Tab", Label: "Árvore", Priority: 300, Enabled: true},
		{Key: "⌃R", Label: "Revelar", Priority: 200, Enabled: m.detail.secret != nil},
	}
}

// View implements tea.Model.
func (m *RootModel) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	st := m.styles
	w := m.width
	h := m.height

	// Work area height: total - header(2) - msgbar(1) - cmdbar(1)
	workH := h - 4
	if workH < 0 {
		workH = 0
	}

	// --- Header (2 lines) ---
	header := renderHeader(st, w, m.mode, m.vaultName, m.isDirty(), m.searchQuery, m.searchActive)

	// --- Work area ---
	var workArea string
	switch m.mode {
	case modeWelcome:
		workArea = renderWelcome(st, w, workH, m.version)
	case modeVault:
		workArea = m.renderVaultMode(w, workH)
	case modeModels:
		workArea = renderPlaceholder(st, w, workH, "Modo Modelos — em construção")
	case modeConfig:
		workArea = renderPlaceholder(st, w, workH, "Modo Configurações — em construção")
	default:
		workArea = renderWelcome(st, w, workH, m.version)
	}

	// Overlay dialog if active
	if m.activeDialog != nil {
		dlgStr := renderDialog(st, *m.activeDialog, w, workH)
		workArea = overlayDialog(workArea, dlgStr, w, workH)
	}

	// --- Message bar (1 line) ---
	msgBar := renderMsgBar(st, w, m.msg, m.spinFrame)

	// --- Command bar (1 line) ---
	cmdBar := renderCmdBar(st, w, m.actions)

	return header + "\n" + workArea + "\n" + msgBar + "\n" + cmdBar
}

// renderVaultMode renders the vault mode (tree + detail panels).
func (m *RootModel) renderVaultMode(width, height int) string {
	treeW := width * 35 / 100
	if treeW < 15 {
		treeW = 15
	}
	detailW := width - treeW

	m.tree.clampScroll(height)

	var detailSecret *vault.Segredo
	if m.treeFocused == focusTree || m.detail.secret != nil {
		detailSecret = m.detail.secret
	}

	treeStr := renderTree(m.styles, m.tree, height, treeW+1, treeW, detailSecret, m.treeFocused == focusTree)
	detailStr := renderDetail(m.styles, m.detail, detailW, height)

	treeLines := strings.Split(treeStr, "\n")
	detailLines := strings.Split(detailStr, "\n")

	combined := make([]string, height)
	for i := 0; i < height; i++ {
		tl := ""
		if i < len(treeLines) {
			tl = treeLines[i]
		}
		dl := ""
		if i < len(detailLines) {
			dl = detailLines[i]
		}
		combined[i] = tl + dl
	}
	return strings.Join(combined, "\n")
}

func (m *RootModel) isDirty() bool {
	return m.manager != nil && m.manager.IsModified()
}

// renderPlaceholder renders a centered placeholder text for unimplemented modes.
func renderPlaceholder(st styles, width, height int, text string) string {
	lines := make([]string, height)
	if height > 0 {
		row := height / 2
		lines[row] = centerText(st.TextDisabled.Italic(true).Render(text), width, len([]rune(text)))
	}
	return strings.Join(lines, "\n")
}

// vaultDisplayName extracts the display name from a file path.
// Removes directory path and ".abditum" extension.
func vaultDisplayName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext == ".abditum" {
		return base[:len(base)-len(ext)]
	}
	return base
}

	// Layout
	width  int
	height int

	// Theme
	theme  Theme
	styles styles

	// Mode
	mode appMode

	// Vault state
	manager   *vault.Manager
	vaultName string // display name (filename without extension)

	// Tree panel
	tree        treeModel
	treeFocused panelFocus

	// Detail panel
	detail detailModel

	// Search
	searchActive bool
	searchQuery  string

	// Message bar
	msg       barMessage
	spinFrame int

	// Active dialog (nil = no dialog)
	activeDialog *dialog
	dialogCb     func(actionIdx int) // callback when user selects a dialog action

	// Actions for current context
	actions []Action
}

// NewRootModel creates a new RootModel with the given options.
func NewRootModel(opts ...Option) *RootModel {
	m := &RootModel{
		theme: ThemeTokyoNight,
		width: 80,
		height: 24,
		mode:  modeWelcome,
		tree:  newTreeModel(),
	}
	m.styles = newStyles(tokyoNight)
	for _, o := range opts {
		o(m)
	}
	m.rebuildActions()
	return m
}

// Init implements tea.Model.
func (m *RootModel) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
	)
}

// Update implements tea.Model.
func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		m.spinFrame++
		// Check TTL
		if m.msg.kind != msgNone && m.msg.kind != msgBusy &&
			m.msg.kind != msgHint && !m.msg.expires.IsZero() &&
			time.Now().After(m.msg.expires) {
			m.msg = noMessage
		}
		return m, tickCmd()

	case clearMsgMsg:
		m.msg = noMessage
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// handleKey dispatches key events.
func (m *RootModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Global keys — always active
	switch key {
	case "ctrl+q":
		return m.handleQuit()
	case "f12":
		m.toggleTheme()
		return m, nil
	case "f1":
		// TODO: open help
		return m, nil
	}

	// Dialog active — route to dialog handler
	if m.activeDialog != nil {
		return m.handleDialogKey(key)
	}

	// Mode-specific keys
	switch m.mode {
	case modeWelcome:
		return m.handleWelcomeKey(key)
	case modeVault:
		return m.handleVaultKey(key)
	}

	return m, nil
}

func (m *RootModel) handleQuit() (tea.Model, tea.Cmd) {
	if m.manager != nil && m.manager.IsModified() {
		// Show confirmation dialog
		m.openDialog(dialog{
			title:    "Sair do Abditum",
			body:     "Cofre modificado. Salvar ou descartar?",
			severity: severityDestructive,
			actions: []dialogAction{
				{Key: "Enter", Label: "Salvar", IsEnter: true},
				{Key: "D", Label: "Descartar"},
				{Key: "Esc", Label: "Voltar", IsESC: true},
			},
		}, func(idx int) {
			// idx 0 = Salvar, 1 = Descartar, 2 = Esc/Voltar
			// For now just quit on any action
		})
		return m, nil
	}
	return m, tea.Quit
}

func (m *RootModel) handleWelcomeKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "f5":
		// TODO: create new vault flow
	case "f6":
		// TODO: open vault flow
	}
	return m, nil
}

func (m *RootModel) handleVaultKey(key string) (tea.Model, tea.Cmd) {
	if m.manager == nil {
		return m, nil
	}

	switch key {
	case "f2":
		// already in vault mode
	case "f3":
		m.switchMode(modeModels)
	case "f4":
		m.switchMode(modeConfig)
	case "f10", "ctrl+f":
		m.toggleSearch()
	case "ctrl+alt+shift+q":
		// Emergency lock — no confirmation
		if m.manager != nil {
			m.manager.Lock()
			m.mode = modeWelcome
			m.vaultName = ""
			m.manager = nil
			m.tree = newTreeModel()
			m.detail = detailModel{}
		}
		return m, nil
	}

	// Search-active input routing
	if m.searchActive {
		return m.handleSearchKey(key)
	}

	// Tree / detail navigation
	if m.treeFocused == focusTree {
		return m.handleTreeKey(key)
	}
	return m.handleDetailKey(key)
}

func (m *RootModel) handleSearchKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc", "ctrl+f", "f10":
		m.searchActive = false
		m.searchQuery = ""
		m.setHint("• ↑↓ para navegar")
		m.rebuildActions()
	case "delete":
		m.searchQuery = ""
		m.setHint("• Digite para filtrar os segredos")
	case "backspace":
		if len(m.searchQuery) > 0 {
			runes := []rune(m.searchQuery)
			m.searchQuery = string(runes[:len(runes)-1])
		}
		if m.searchQuery == "" {
			m.setHint("• Digite para filtrar os segredos")
		}
	case "up":
		m.tree.moveUp()
	case "down":
		m.tree.moveDown()
	case "home":
		m.tree.moveHome()
	case "end":
		m.tree.moveEnd()
	default:
		// printable character → append to query
		if len(key) == 1 && key[0] >= 32 {
			m.searchQuery += key
		}
	}
	m.syncDetailFromTree()
	return m, nil
}

func (m *RootModel) handleTreeKey(key string) (tea.Model, tea.Cmd) {
	cofre := m.manager.Vault()

	switch key {
	case "up":
		m.tree.moveUp()
	case "down":
		m.tree.moveDown()
	case "home":
		m.tree.moveHome()
	case "end":
		m.tree.moveEnd()
	case "right", "→":
		m.tree.expandFolder(cofre)
	case "left", "←":
		m.tree.collapseFolder(cofre)
	case "enter":
		m.tree.expandFolder(cofre)
	case "tab":
		m.treeFocused = focusDetail
		m.detail.focused = true
		m.rebuildActions()
		return m, nil
	case "ctrl+s":
		return m.savekVault()
	case "f7":
		return m.savekVault()
	}

	m.syncDetailFromTree()
	m.rebuildActions()
	return m, nil
}

func (m *RootModel) handleDetailKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "tab":
		m.treeFocused = focusTree
		m.detail.focused = false
		m.setHint("• ↑↓ para navegar")
		m.rebuildActions()
	case "up":
		if m.detail.scrollTop > 0 {
			m.detail.scrollTop--
		}
	case "down":
		m.detail.scrollTop++
	case "ctrl+r":
		// cycle reveal
		switch m.detail.reveal {
		case revealMasked:
			m.detail.reveal = revealHint
		case revealHint:
			m.detail.reveal = revealFull
		case revealFull:
			m.detail.reveal = revealMasked
		}
	case "esc":
		m.treeFocused = focusTree
		m.detail.focused = false
		m.rebuildActions()
	}
	return m, nil
}

func (m *RootModel) handleDialogKey(key string) (tea.Model, tea.Cmd) {
	if m.activeDialog == nil {
		return m, nil
	}
	actions := m.activeDialog.actions
	switch key {
	case "esc":
		for i, a := range actions {
			if a.IsESC {
				m.closeDialog(i)
				return m, nil
			}
		}
		m.closeDialog(-1)
	case "enter":
		for i, a := range actions {
			if a.IsEnter {
				m.closeDialog(i)
				return m, nil
			}
		}
		m.closeDialog(0)
	default:
		for i, a := range actions {
			if strings.EqualFold(a.Key, key) && !a.IsEnter && !a.IsESC {
				m.closeDialog(i)
				return m, nil
			}
		}
	}
	return m, nil
}

func (m *RootModel) openDialog(d dialog, cb func(int)) {
	m.activeDialog = &d
	m.dialogCb = cb
}

func (m *RootModel) closeDialog(actionIdx int) {
	m.activeDialog = nil
	if m.dialogCb != nil && actionIdx >= 0 {
		m.dialogCb(actionIdx)
		m.dialogCb = nil
	}
}

func (m *RootModel) toggleTheme() {
	if m.theme == ThemeTokyoNight {
		m.theme = ThemeCyberpunk
		m.styles = newStyles(cyberpunk)
	} else {
		m.theme = ThemeTokyoNight
		m.styles = newStyles(tokyoNight)
	}
}

func (m *RootModel) toggleSearch() {
	m.searchActive = !m.searchActive
	if m.searchActive {
		m.searchQuery = ""
		m.setHint("• Digite para filtrar os segredos")
	} else {
		m.searchQuery = ""
		m.setHint("• ↑↓ para navegar")
	}
	m.rebuildActions()
}

func (m *RootModel) switchMode(mode appMode) {
	m.mode = mode
	m.rebuildActions()
}

func (m *RootModel) setMsg(kind msgType, text string, ttl time.Duration) {
	expires := time.Time{}
	if ttl > 0 {
		expires = time.Now().Add(ttl)
	}
	m.msg = barMessage{kind: kind, text: text, expires: expires}
}

func (m *RootModel) setHint(text string) {
	m.msg = barMessage{kind: msgHint, text: text}
}

func (m *RootModel) savekVault() (tea.Model, tea.Cmd) {
	if m.manager == nil {
		return m, nil
	}
	// TODO: actual save with file picker if no path
	m.setMsg(msgSuccess, "✓ Cofre salvo", 5*time.Second)
	return m, nil
}

func (m *RootModel) syncDetailFromTree() {
	s := m.tree.selectedSecret()
	if s != m.detail.secret {
		m.detail.secret = s
		m.detail.scrollTop = 0
		m.detail.reveal = revealMasked
	}
}

func (m *RootModel) rebuildActions() {
	st := m.styles
	_ = st

	if m.activeDialog != nil {
		// During a dialog, no main actions shown
		m.actions = nil
		return
	}

	if m.searchActive {
		m.actions = []Action{
			{Key: "⌃F", Label: "Fechar", Priority: 100, Enabled: true},
			{Key: "Del", Label: "Limpar", Priority: 90, Enabled: true},
		}
		return
	}

	switch m.mode {
	case modeWelcome:
		m.actions = []Action{
			{Key: "F5", Label: "Novo cofre", Priority: 100, Enabled: true},
			{Key: "F6", Label: "Abrir cofre", Priority: 90, Enabled: true},
		}
	case modeVault:
		if m.treeFocused == focusTree {
			m.actions = treeActions(m)
		} else {
			m.actions = detailActions(m)
		}
	case modeModels:
		m.actions = []Action{
			{Key: "F2", Label: "Cofre", Priority: 100, Enabled: true},
			{Key: "F4", Label: "Config", Priority: 90, Enabled: true},
		}
	case modeConfig:
		m.actions = []Action{
			{Key: "F2", Label: "Cofre", Priority: 100, Enabled: true},
			{Key: "F3", Label: "Modelos", Priority: 90, Enabled: true},
		}
	}
}

func treeActions(m *RootModel) []Action {
	actions := []Action{
		{Key: "F2", Label: "Cofre", Priority: 200, Enabled: true, HideFromBar: true},
		{Key: "F3", Label: "Modelos", Priority: 190, Enabled: true},
		{Key: "F4", Label: "Config", Priority: 180, Enabled: true},
		{Key: "F7", Label: "Salvar", Priority: 150, Enabled: m.manager != nil && m.manager.IsModified()},
		{Key: "F10", Label: "Busca", Priority: 140, Enabled: true},
	}

	s := m.tree.selectedSecret()
	if s != nil {
		hasSensitive := false
		for _, c := range s.Campos() {
			if c.Tipo() == vault.TipoCampoSensivel {
				hasSensitive = true
				break
			}
		}
		actions = append(actions,
			Action{Key: "Enter", Label: "Detalhes", Priority: 300, Enabled: true},
			Action{Key: "Ins", Label: "Novo", Priority: 250, Enabled: true},
			Action{Key: "Del", Label: "Excluir", Priority: 120, Enabled: true},
		)
		if hasSensitive {
			revealLabel := "Revelar"
			if m.detail.reveal == revealHint {
				revealLabel = "Mostrar tudo"
			} else if m.detail.reveal == revealFull {
				revealLabel = "Ocultar"
			}
			actions = append(actions,
				Action{Key: "⌃R", Label: revealLabel, Priority: 130, Enabled: true},
				Action{Key: "⌃C", Label: "Copiar", Priority: 125, Enabled: true},
			)
		}
	}
	return actions
}

func detailActions(m *RootModel) []Action {
	return []Action{
		{Key: "Tab", Label: "Árvore", Priority: 300, Enabled: true},
		{Key: "⌃R", Label: "Revelar", Priority: 200, Enabled: m.detail.secret != nil},
	}
}

// View implements tea.Model.
func (m *RootModel) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	st := m.styles
	w := m.width
	h := m.height

	// Work area height: total - header(2) - msgbar(1) - cmdbar(1)
	workH := h - 4
	if workH < 0 {
		workH = 0
	}

	// --- Header (2 lines) ---
	header := renderHeader(st, w, m.mode, m.vaultName, m.isDirty(), m.searchQuery, m.searchActive)

	// --- Work area ---
	var workArea string
	switch m.mode {
	case modeWelcome:
		workArea = renderWelcome(st, w, workH, m.version)
	case modeVault:
		workArea = m.renderVaultMode(w, workH)
	case modeModels:
		workArea = renderPlaceholder(st, w, workH, "Modo Modelos — em construção")
	case modeConfig:
		workArea = renderPlaceholder(st, w, workH, "Modo Configurações — em construção")
	default:
		workArea = renderWelcome(st, w, workH, m.version)
	}

	// Overlay dialog if active
	if m.activeDialog != nil {
		dlgStr := renderDialog(st, *m.activeDialog, w, workH)
		workArea = overlayDialog(workArea, dlgStr, w, workH)
	}

	// --- Message bar (1 line) ---
	msgBar := renderMsgBar(st, w, m.msg, m.spinFrame)

	// --- Command bar (1 line) ---
	cmdBar := renderCmdBar(st, w, m.actions)

	return header + "\n" + workArea + "\n" + msgBar + "\n" + cmdBar
}

// renderVaultMode renders the vault mode (tree + detail panels).
func (m *RootModel) renderVaultMode(width, height int) string {
	treeW := width * 35 / 100
	if treeW < 15 {
		treeW = 15
	}
	detailW := width - treeW

	m.tree.clampScroll(height)

	var detailSecret *vault.Segredo
	if m.treeFocused == focusTree || m.detail.secret != nil {
		detailSecret = m.detail.secret
	}

	treeStr := renderTree(m.styles, m.tree, height, treeW+1, treeW, detailSecret, m.treeFocused == focusTree)
	detailStr := renderDetail(m.styles, m.detail, detailW, height)

	treeLines := strings.Split(treeStr, "\n")
	detailLines := strings.Split(detailStr, "\n")

	combined := make([]string, height)
	for i := 0; i < height; i++ {
		tl := ""
		if i < len(treeLines) {
			tl = treeLines[i]
		}
		dl := ""
		if i < len(detailLines) {
			dl = detailLines[i]
		}
		combined[i] = tl + dl
	}
	return strings.Join(combined, "\n")
}

func (m *RootModel) isDirty() bool {
	return m.manager != nil && m.manager.IsModified()
}

// renderPlaceholder renders a centered placeholder text for unimplemented modes.
func renderPlaceholder(st styles, width, height int, text string) string {
	lines := make([]string, height)
	if height > 0 {
		row := height / 2
		lines[row] = centerText(st.TextDisabled.Italic(true).Render(text), width, len([]rune(text)))
	}
	return strings.Join(lines, "\n")
}

// vaultDisplayName extracts the display name from a file path.
// Removes directory path and ".abditum" extension.
func vaultDisplayName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext == ".abditum" {
		return base[:len(base)-len(ext)]
	}
	return base
}

