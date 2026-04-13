package tui

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/storage"
	"github.com/useful-toys/abditum/internal/vault"
)

// rootModel is the sole tea.Model in the tui package.
// It owns the workArea state machine, modal stack, active flow slot,
// shared services, and the frame compositor.
// All child models are stored as concrete pointer fields - nil means inactive.
// Never store children as childModel interface values (typed nil trap).
type rootModel struct {
	// State machine
	area          workArea
	mgr           *vault.Manager
	vaultPath     string
	vaultMetadata storage.FileMetadata // snapshot for external change detection
	initialPath   string               // Path passed via CLI, for fast-path
	isDirty       bool
	width         int
	height        int
	theme         *Theme
	header        headerModel
	version       string // Application version, injected at build time

	// Child models - nil = inactive. NEVER store as childModel interface.
	welcome        *welcomeModel
	vaultTree      *vaultTreeModel
	secretDetail   *secretDetailModel
	templateList   *templateListModel
	templateDetail *templateDetailModel
	settings       *settingsModel

	// Modal stack - LIFO; last element = topmost/active. []modalView per D-02.
	modals []modalView

	// Active flow - nil = no flow running.
	activeFlow flowHandler

	// Shared services.
	actions  *ActionManager
	messages *MessageManager

	// Timer fields.
	lastActionAt time.Time
	lastCopyAt   time.Time // D-12: reset when a field is copied to clipboard
	lastRevealAt time.Time // D-12: reset when a sensitive field is revealed
}

// Compile-time assertion: rootModel satisfies tea.Model.
var _ tea.Model = &rootModel{}

// RootModelOption is a functional option for configuring rootModel
type RootModelOption func(*rootModel)

// WithVersion sets the application version for display in the welcome screen
func WithVersion(version string) RootModelOption {
	return func(m *rootModel) {
		m.version = version
	}
}

// WithInitialPath sets the initial vault path for CLI fast-path
func WithInitialPath(path string) RootModelOption {
	return func(m *rootModel) {
		m.initialPath = path
	}
}

// NewRootModel is the exported constructor for main.go (PoC mode, D-04).
// Accepts functional options for configuration.
func NewRootModel(opts ...RootModelOption) *rootModel {
	m := newRootModel("")
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// newRootModel constructs a fully initialized rootModel in PoC mode.
// mgr is nil, vaultPath is "", area is workAreaWelcome (D-02, D-03, D-05).
// If initialPath is non-empty, the CLI fast-path will be initiated in Init().
// version defaults to "dev" if not set via WithVersion option.
func newRootModel(initialPath string) *rootModel {
	actions := NewActionManager()
	messages := NewMessageManager()

	m := &rootModel{
		area:         workAreaWelcome,
		mgr:          nil, // PoC mode — no vault (D-02)
		vaultPath:    "",
		initialPath:  initialPath,
		isDirty:      false,
		actions:      actions,
		messages:     messages,
		lastActionAt: time.Now(),
		theme:        TokyoNight,
		header:       headerModel{},
		version:      "dev",
	}

	m.welcome = newWelcomeModel(m.version)

	// Register production actions: global F1 Help, F12 theme toggle, and vault open/create flows.
	actions.Register(m,
		// Navigation/Global actions
		Action{Keys: []string{"ctrl+q"}, Label: "Sair", Description: "Sair do Abditum",
			Group: 0, Scope: ScopeLocal, Priority: 10, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { return tea.Quit }},
		// Vault actions (pre-vault scope)
		Action{Keys: []string{"f6"}, Label: "Abrir", Description: "Abrir cofre existente",
			Group: 4, Scope: ScopeLocal, Priority: 94, HideFromBar: false,
			Enabled: func() bool { return m.area == workAreaWelcome },
			Handler: func() tea.Cmd {
				flow := newOpenVaultFlow(m.mgr, m.messages, actions, m.theme)
				return func() tea.Msg { return startFlowMsg{flow: flow} }
			}},
		Action{Keys: []string{"f5"}, Label: "Novo", Description: "Criar novo cofre",
			Group: 4, Scope: ScopeLocal, Priority: 95, HideFromBar: false,
			Enabled: func() bool { return m.area == workAreaWelcome },
			Handler: func() tea.Cmd {
				flow := newCreateVaultFlow(m.mgr, m.messages, actions, m.theme)
				return func() tea.Msg { return startFlowMsg{flow: flow} }
			}},
		// Global action for F12 to toggle theme
		Action{Keys: []string{"f12"}, Label: "Toggle Theme", Description: "Alternar tema visual (Tokyo Night / Cyberpunk)",
			Group: 0, Scope: ScopeGlobal, Priority: 100, HideFromBar: true,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { return func() tea.Msg { return toggleThemeMsg{} } }},
		Action{Keys: []string{"f1"}, Label: "Ajuda", Description: "Mostrar atalhos de teclado",
			Group: 1, Scope: ScopeGlobal, Priority: 0, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd {
				return func() tea.Msg { return pushModalMsg{modal: newHelpModal(actions.All(), actions.GroupLabel)} }
			}},
	)
	actions.RegisterGroupLabel(4, "Cofre")

	return m
}

// Init satisfies tea.Model. Always starts global tick for message TTL (D-10, D-11).
// If initialPath is set (CLI fast-path), route to the appropriate flow per D-CLI-01.
func (m *rootModel) Init() tea.Cmd {
	tickCmd := tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })

	// D-CLI-01: 3-case routing when initialPath is non-empty
	if m.initialPath != "" {
		return tea.Batch(
			tickCmd,
			func() tea.Msg {
				info, err := os.Stat(m.initialPath)
				if err == nil && !info.IsDir() {
					// Case 1: file exists → openVaultFlow fast-path (step 3: password entry)
					flow := newOpenVaultFlow(nil, m.messages, m.actions, m.theme)
					flow.cliPath = m.initialPath
					return startFlowMsg{flow: flow}
				}
				if os.IsNotExist(err) {
					// Check if parent directory exists
					parentDir := filepath.Dir(m.initialPath)
					if _, parentErr := os.Stat(parentDir); parentErr == nil {
						// Case 2: file missing but parent exists → createVaultFlow fast-path (step 3)
						flow := newCreateVaultFlow(nil, m.messages, m.actions, m.theme)
						flow.cliPath = m.initialPath
						return startFlowMsg{flow: flow}
					}
				}
				// Case 3: parent dir doesn't exist → no flow, show welcome screen normally
				return nil
			},
		)
	}

	return tickCmd
}

// Update satisfies tea.Model. Implements the D-09 dispatch order.
func (m *rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// --- Window resize: store dimensions only (D-02) ---
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	// --- Modal stack: push ---
	case pushModalMsg:
		if msg.modal != nil {
			m.modals = append(m.modals, msg.modal)
			type initer interface{ Init() tea.Cmd }
			if mi, ok := msg.modal.(initer); ok {
				return m, mi.Init()
			}
		}
		return m, nil

	// --- Modal stack: pop ---
	case popModalMsg:
		if len(m.modals) > 0 {
			m.modals = m.modals[:len(m.modals)-1]
		}
		return m, nil

	// --- Flow lifecycle: start (D-08) ---
	case startFlowMsg:
		m.modals = m.modals[:0] // clear orphan modals from any previous flow
		m.activeFlow = msg.flow
		return m, m.activeFlow.Init()

	// --- Flow lifecycle: end (D-08) ---
	case endFlowMsg:
		m.activeFlow = nil
		return m, nil

	// --- Vault opened: transition to work area and store vault (D-08) ---
	case vaultOpenedMsg:
		// TODO: In Phase 9+, populate m.mgr with the opened vault
		// For now, just transition to vault area
		m.area = workAreaVault
		m.vaultPath = msg.Path
		m.vaultMetadata = msg.Metadata
		m.isDirty = false
		return m, nil

	// --- Modal result: route ONLY to activeFlow (D-03) ---
	case modalResult:
		if m.activeFlow != nil {
			return m, m.activeFlow.Update(msg)
		}
		return m, nil

	// --- Global tick (D-07): advance message TTL + re-issue tick ---
	case tickMsg:
		m.messages.Tick()
		cmds := m.broadcast(msg)
		cmds = append(cmds, tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) }))
		return m, tea.Batch(cmds...)

	// --- Domain messages: broadcast to all live models ---
	case secretAddedMsg, secretDeletedMsg, secretRestoredMsg, secretModifiedMsg,
		secretMovedMsg, secretReorderedMsg, folderStructureChangedMsg,
		vaultSavedMsg, vaultReloadedMsg, vaultClosedMsg, vaultChangedMsg:
		return m, tea.Batch(m.broadcast(msg)...)

	// --- Save-and-exit flow internal messages: route to activeFlow ---
	case saveAndExitReadyMsg, saveAndExitOKMsg:
		if m.activeFlow != nil {
			return m, m.activeFlow.Update(msg)
		}
		return m, nil

	// --- Keyboard input: D-09 dispatch order ---
	case tea.KeyPressMsg:
		m.messages.HandleInput()
		m.lastActionAt = time.Now()
		key := msg.String()
		inFlowOrModal := m.activeFlow != nil || len(m.modals) > 0

		// Fluxo 6: Emergency vault lock — discards all unsaved changes immediately.
		// Must be checked BEFORE ctrl+q to guarantee it is always handled first.
		if key == "ctrl+alt+shift+q" {
			m.mgr = nil
			m.vaultPath = ""
			m.vaultMetadata = storage.FileMetadata{}
			m.isDirty = false
			m.activeFlow = nil
			m.modals = m.modals[:0]
			m.area = workAreaWelcome
			return m, nil
		}

		// Check for Ctrl+Q (exit flow) before any other key handling
		if key == "ctrl+q" {
			// Fluxo 5: vault is open AND has unsaved changes
			if m.mgr != nil && m.mgr.IsModified() {
				return m, Decision(SeverityAlert, "Sair do Abditum",
					"Cofre modificado. Salvar ou descartar?",
					DecisionAction{Key: "S", Label: "Salvar", Default: true,
						Cmd: func() tea.Msg {
							flow := newSaveAndExitFlow(m.mgr, m.vaultPath, m.vaultMetadata, m.messages, m.theme)
							return startFlowMsg{flow: flow}
						}},
					[]DecisionAction{{Key: "D", Label: "Descartar", Cmd: tea.Quit}},
					DecisionAction{Key: "Esc", Label: "Voltar"})
			}
			// Fluxos 3 & 4: no vault open (or vault open but clean) — confirm exit
			return m, Decision(SeverityNeutral, "Sair do Abditum",
				"Sair do Abditum?",
				DecisionAction{Key: "Enter", Label: "Sair", Default: true, Cmd: tea.Quit},
				nil,
				DecisionAction{Key: "Esc", Label: "Voltar"})
		}

		// Check for F12 theme toggle before any other key handling
		if key == "f12" {
			if m.theme == TokyoNight {
				m.theme = Cyberpunk
			} else {
				m.theme = TokyoNight
			}
			return m, nil
		}

		// 1. If help modal is open, let it handle ALL keys (including F1 and ESC)
		if len(m.modals) > 0 {
			if _, isHelp := m.modals[len(m.modals)-1].(*helpModal); isHelp {
				return m, m.modals[len(m.modals)-1].Update(msg)
			}
		}
		// 2. ActionManager.Dispatch - handles ScopeGlobal and ScopeLocal
		if cmd := m.actions.Dispatch(key, inFlowOrModal); cmd != nil {
			return m, cmd
		}
		// 3. Other modals
		if len(m.modals) > 0 {
			return m, m.modals[len(m.modals)-1].Update(msg)
		}
		// 3. Active flow (fallback - flow without modal)
		if m.activeFlow != nil {
			return m, m.activeFlow.Update(msg)
		}
		// 4. Active work-area child - only if width > 0 (guard against Update before first View)
		if m.width > 0 {
			if child := m.activeChild(); child != nil {
				return m, child.Update(msg)
			}
		}
		return m, nil
	}

	return m, nil
}

// broadcast sends msg to all live work-area children and all active modals.
func (m *rootModel) broadcast(msg tea.Msg) []tea.Cmd {
	var cmds []tea.Cmd
	for _, c := range m.liveWorkChildren() {
		if cmd := c.Update(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	for _, modal := range m.modals {
		if cmd := modal.Update(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return cmds
}

// liveWorkChildren returns all non-nil work-area children as []childModel.
// Modals are NOT included - they are iterated via m.modals separately.
// Uses explicit nil checks on concrete pointer fields to avoid typed-nil trap.
func (m *rootModel) liveWorkChildren() []childModel {
	var live []childModel
	if m.welcome != nil {
		live = append(live, m.welcome)
	}
	if m.vaultTree != nil {
		live = append(live, m.vaultTree)
	}
	if m.secretDetail != nil {
		live = append(live, m.secretDetail)
	}
	if m.templateList != nil {
		live = append(live, m.templateList)
	}
	if m.templateDetail != nil {
		live = append(live, m.templateDetail)
	}
	if m.settings != nil {
		live = append(live, m.settings)
	}
	return live
}

// activeChild returns the single base child model for the current workArea.
func (m *rootModel) activeChild() childModel {
	switch m.area {
	case workAreaWelcome:
		if m.welcome != nil {
			return m.welcome
		}
	case workAreaVault:
		if m.vaultTree != nil {
			return m.vaultTree
		}
	case workAreaTemplates:
		if m.templateList != nil {
			return m.templateList
		}
	case workAreaSettings:
		if m.settings != nil {
			return m.settings
		}
	}
	return nil
}

// View satisfies tea.Model. Composes the full frame and overlays any active modal.
// View satisfies tea.Model. Composes the full frame and overlays any active modal.
// The modal is placed inside the work area only — header, msgBar, and cmdBar are
// never covered.
func (m *rootModel) View() tea.View {
	content := m.renderFrame(nil)

	if len(m.modals) > 0 {
		const headerH = 2
		const msgBarH = 1
		const cmdBarH = 1
		workH := m.height - headerH - msgBarH - cmdBarH
		if workH < 0 {
			workH = 0
		}
		top := m.modals[len(m.modals)-1]
		// Pass dimensions to renderFrame which will pass them to modal.View()
		content = m.renderFrameWithModal(top, m.width, workH)
	}

	v := tea.NewView(content)
	v.AltScreen = true
	v.BackgroundColor = lipgloss.Color(m.theme.Surface.Base)
	return v
}

// renderFrameWithModal renders a frame with a modal at specified dimensions.
// This is used when pushing a modal to pass dimensions to its View() method.
func (m *rootModel) renderFrameWithModal(modal modalView, modalWidth, modalHeight int) string {
	cmdBarStyle := lipgloss.NewStyle().Width(m.width).Background(lipgloss.Color(m.theme.Surface.Base))
	workAreaStyle := lipgloss.NewStyle().Width(m.width).Background(lipgloss.Color(m.theme.Surface.Base))

	const headerH = 2
	const msgBarH = 1
	const cmdBarH = 1
	workH := m.height - headerH - msgBarH - cmdBarH
	if workH < 0 {
		workH = 0
	}

	// Header
	header := m.header.Render(m.width, m.vaultPath, m.mgr != nil && m.mgr.IsModified(), m.area, m.theme)

	// Message bar
	msgBar := RenderMessageBar(m.messages.Current(), m.width, m.theme)

	// Work area - render normally (no modal overlay)
	var workContent string
	switch m.area {
	case workAreaWelcome:
		if m.welcome != nil {
			workContent = m.welcome.View(m.width, workH, m.theme)
		}
	case workAreaVault:
		workContent = m.renderVaultArea(workH)
	case workAreaTemplates:
		workContent = m.renderTemplatesArea(workH)
	case workAreaSettings:
		if m.settings != nil {
			workContent = m.settings.View(m.width, workH, m.theme)
		}
	}
	workArea := workAreaStyle.Height(workH).Render(workContent)

	// Overlay modal centered inside work area using lipgloss.Place.
	// Pass dimensions to modal.View() as per new interface
	if modal != nil {
		modalStr := modal.View(modalWidth, modalHeight, m.theme)
		workArea = lipgloss.Place(m.width, workH, lipgloss.Center, lipgloss.Center, modalStr)
	}

	// Command bar: always render so the frame occupies exactly `height` lines.
	// When no shortcuts, render a blank background line.
	var cmdBarContent string
	if modal != nil {
		cmdBarContent = renderShortcuts(modal.Shortcuts(), m.width, m.theme)
	} else {
		cmdBarContent = RenderCommandBar(m.actions.Visible(), m.width, m.theme)
	}
	if cmdBarContent == "" {
		cmdBarContent = strings.Repeat(" ", m.width)
	}
	cmdBar := cmdBarStyle.Render(cmdBarContent)

	return strings.Join([]string{header, workArea, msgBar, cmdBar}, "\n")
}

// overlayModal renders a modal dialog centered over the existing frame content.
// The frame remains visible behind/around the modal — only the modal region is replaced.
// Deprecated: kept for reference; replaced by lipgloss.Place inside renderFrame.
// renderShortcuts renders a command bar from modal shortcuts.
func renderShortcuts(shortcuts []Shortcut, width int, theme *Theme) string {
	if len(shortcuts) == 0 {
		return ""
	}
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Accent.Primary)).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Primary))
	sepStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Secondary))
	var parts []string
	for _, s := range shortcuts {
		parts = append(parts, keyStyle.Render(s.Key)+" "+labelStyle.Render(s.Label))
	}
	return "  " + strings.Join(parts, sepStyle.Render("  |  ")+"  ")
}

// renderFrame composes the full frame: header + work area + msg bar + cmd bar.
// If modal is non-nil, it is centered inside the work area using lipgloss.Place.
func (m *rootModel) renderFrame(modal modalView) string {
	cmdBarStyle := lipgloss.NewStyle().Width(m.width).Background(lipgloss.Color(m.theme.Surface.Base))
	workAreaStyle := lipgloss.NewStyle().Width(m.width).Background(lipgloss.Color(m.theme.Surface.Base))

	const headerH = 2
	const msgBarH = 1
	const cmdBarH = 1
	workH := m.height - headerH - msgBarH - cmdBarH
	if workH < 0 {
		workH = 0
	}

	// Header
	header := m.header.Render(m.width, m.vaultPath, m.mgr != nil && m.mgr.IsModified(), m.area, m.theme)

	// Message bar
	msgBar := RenderMessageBar(m.messages.Current(), m.width, m.theme)

	// Work area - pass dimensions to View() calls
	var workContent string
	switch m.area {
	case workAreaWelcome:
		if m.welcome != nil {
			workContent = m.welcome.View(m.width, workH, m.theme)
		}
	case workAreaVault:
		workContent = m.renderVaultArea(workH)
	case workAreaTemplates:
		workContent = m.renderTemplatesArea(workH)
	case workAreaSettings:
		if m.settings != nil {
			workContent = m.settings.View(m.width, workH, m.theme)
		}
	}
	workArea := workAreaStyle.Height(workH).Render(workContent)

	// Overlay modal centered inside work area using lipgloss.Place.
	// lipgloss.Place handles ANSI correctly and never touches header/msgBar/cmdBar.
	if modal != nil {
		modalStr := modal.View(m.width, workH, m.theme)
		workArea = lipgloss.Place(m.width, workH, lipgloss.Center, lipgloss.Center, modalStr)
	}

	// Command bar: always render so the frame occupies exactly `height` lines.
	// When no shortcuts, render a blank background line.
	var cmdBarContent string
	if modal != nil {
		cmdBarContent = renderShortcuts(modal.Shortcuts(), m.width, m.theme)
	} else {
		cmdBarContent = RenderCommandBar(m.actions.Visible(), m.width, m.theme)
	}
	if cmdBarContent == "" {
		cmdBarContent = strings.Repeat(" ", m.width)
	}
	cmdBar := cmdBarStyle.Render(cmdBarContent)

	return strings.Join([]string{header, workArea, msgBar, cmdBar}, "\n")
}

// renderVaultArea renders workAreaVault: vaultTree (left) + secretDetail (right).
func (m *rootModel) renderVaultArea(workH int) string {
	halfW := m.width / 2

	left := "[vault tree - Phase 7]"
	right := "[secret detail - Phase 8]"
	if m.vaultTree != nil {
		left = m.vaultTree.View(halfW, workH, m.theme)
	}
	if m.secretDetail != nil {
		right = m.secretDetail.View(m.width-halfW, workH, m.theme)
	}

	leftStyle := lipgloss.NewStyle().Width(halfW).Height(workH)
	rightStyle := lipgloss.NewStyle().Width(m.width - halfW).Height(workH)
	return lipgloss.JoinHorizontal(lipgloss.Top, leftStyle.Render(left), rightStyle.Render(right))
}

// renderTemplatesArea renders workAreaTemplates: templateList (left) + templateDetail (right).
func (m *rootModel) renderTemplatesArea(workH int) string {
	halfW := m.width / 2

	left := "[template list - Phase 8]"
	right := "[template detail - Phase 8]"
	if m.templateList != nil {
		left = m.templateList.View(halfW, workH, m.theme)
	}
	if m.templateDetail != nil {
		right = m.templateDetail.View(m.width-halfW, workH, m.theme)
	}

	leftStyle := lipgloss.NewStyle().Width(halfW).Height(workH)
	rightStyle := lipgloss.NewStyle().Width(m.width - halfW).Height(workH)
	return lipgloss.JoinHorizontal(lipgloss.Top, leftStyle.Render(left), rightStyle.Render(right))
}

// enterVault transitions to workAreaVault, mounts the vault children, and
// starts the global 1-second tick. Called by domain message handlers in Phase 6+.
func (m *rootModel) enterVault() tea.Cmd {
	m.area = workAreaVault
	m.welcome = nil // GC old model
	m.vaultTree = newVaultTreeModel(m.mgr, m.actions, m.messages)
	m.secretDetail = newSecretDetailModel(m.mgr, m.actions, m.messages)
	// Dimensions are passed to View() when rendering, no SetSize needed
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}
