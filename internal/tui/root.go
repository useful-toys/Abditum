package tui

import (
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/vault"
)

// rootModel is the sole tea.Model in the tui package.
// It owns the workArea state machine, modal stack, active flow slot,
// shared services, and the frame compositor.
// All child models are stored as concrete pointer fields - nil means inactive.
// Never store children as childModel interface values (typed nil trap).
type rootModel struct {
	// State machine
	area      workArea
	mgr       *vault.Manager
	vaultPath string
	width     int
	height    int

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

// NewRootModel is the exported constructor for main.go.
func NewRootModel(mgr *vault.Manager, initialPath string) *rootModel {
	return newRootModel(mgr, initialPath)
}

// newRootModel constructs a fully initialized rootModel.
func newRootModel(mgr *vault.Manager, initialPath string) *rootModel {
	actions := NewActionManager()
	messages := NewMessageManager()

	m := &rootModel{
		area:         workAreaWelcome,
		mgr:          mgr,
		vaultPath:    initialPath,
		actions:      actions,
		messages:     messages,
		lastActionAt: time.Now(),
	}

	m.welcome = newWelcomeModel(actions)

	// Register global actions on rootModel as owner (D-06).
	actions.Register(m,
		Action{
			Keys:        []string{"ctrl+q"},
			Label:       "Quit",
			Description: "Quit Abditum (confirms if unsaved changes)",
			Group:       "Global",
			Scope:       ScopeLocal,
			Enabled:     func() bool { return true },
			Handler: func() tea.Cmd {
				if m.mgr != nil && m.mgr.IsModified() {
					return Confirm(DialogAlert, "Sair", "Ha alteracoes nao salvas. Deseja sair mesmo assim?", tea.Quit, nil)
				}
				return tea.Quit
			},
		},
		Action{
			Keys:        []string{"?"},
			Label:       "Ajuda",
			Description: "Mostrar atalhos de teclado",
			Group:       "Global",
			Scope:       ScopeGlobal,
			Enabled:     func() bool { return true },
			Handler: func() tea.Cmd {
				return func() tea.Msg {
					return pushModalMsg{modal: newHelpModal(actions)}
				}
			},
		},
	)

	return m
}

// Init satisfies tea.Model. Returns nil - the global tick does NOT start here.
func (m *rootModel) Init() tea.Cmd {
	return nil
}

// Update satisfies tea.Model. Implements the D-09 dispatch order.
func (m *rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// --- Window resize: propagate to work-area children only (not modals - D-02) ---
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		for _, child := range m.liveWorkChildren() {
			child.SetSize(msg.Width, msg.Height)
		}
		return m, nil

	// --- Modal stack: push ---
	case pushModalMsg:
		if msg.modal != nil {
			m.modals = append(m.modals, msg.modal)
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

	// --- Keyboard input: D-09 dispatch order ---
	case tea.KeyPressMsg:
		m.messages.HandleInput()
		m.lastActionAt = time.Now()
		key := msg.String()
		inFlowOrModal := m.activeFlow != nil || len(m.modals) > 0

		// 1. ActionManager.Dispatch - handles ScopeGlobal and ScopeLocal
		if cmd := m.actions.Dispatch(key, inFlowOrModal); cmd != nil {
			return m, cmd
		}
		// 2. Topmost modal receives input
		if len(m.modals) > 0 {
			return m, m.modals[len(m.modals)-1].Update(msg)
		}
		// 3. Active flow (fallback - flow without modal)
		if m.activeFlow != nil {
			return m, m.activeFlow.Update(msg)
		}
		// 4. Active work-area child
		if child := m.activeChild(); child != nil {
			return m, child.Update(msg)
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
	if m.welcome != nil        { live = append(live, m.welcome) }
	if m.vaultTree != nil      { live = append(live, m.vaultTree) }
	if m.secretDetail != nil   { live = append(live, m.secretDetail) }
	if m.templateList != nil   { live = append(live, m.templateList) }
	if m.templateDetail != nil { live = append(live, m.templateDetail) }
	if m.settings != nil       { live = append(live, m.settings) }
	return live
}

// activeChild returns the single base child model for the current workArea.
func (m *rootModel) activeChild() childModel {
	switch m.area {
	case workAreaWelcome:
		if m.welcome != nil { return m.welcome }
	case workAreaVault:
		if m.vaultTree != nil { return m.vaultTree }
	case workAreaTemplates:
		if m.templateList != nil { return m.templateList }
	case workAreaSettings:
		if m.settings != nil { return m.settings }
	}
	return nil
}

// View satisfies tea.Model. Composes the full frame and overlays any active modal.
func (m *rootModel) View() tea.View {
	content := m.renderFrame()

	if len(m.modals) > 0 {
		top := m.modals[len(m.modals)-1]
		content = lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center, top.View())
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

// renderShortcuts renders a command bar from modal shortcuts.
func renderShortcuts(shortcuts []Shortcut, width int) string {
	if len(shortcuts) == 0 {
		return ""
	}
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	sepStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	var parts []string
	for _, s := range shortcuts {
		parts = append(parts, keyStyle.Render(s.Key)+" "+labelStyle.Render(s.Label))
	}
	return "  " + strings.Join(parts, sepStyle.Render("  |  ")+"  ")
}

// renderFrame composes the base frame zones: header + message bar + work area + command bar.
func (m *rootModel) renderFrame() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	headerStyle := lipgloss.NewStyle().Width(m.width).Background(lipgloss.Color("236")).Foreground(lipgloss.Color("255"))
	msgBarStyle := lipgloss.NewStyle().Width(m.width).Foreground(lipgloss.Color("245"))
	cmdBarStyle := lipgloss.NewStyle().Width(m.width).Background(lipgloss.Color("236"))
	workAreaStyle := lipgloss.NewStyle().Width(m.width)

	const headerH = 1
	const msgBarH = 1
	const cmdBarH = 1
	workH := m.height - headerH - msgBarH - cmdBarH
	if workH < 0 {
		workH = 0
	}

	// Header: app name + vault path + dirty indicator
	vaultName := "No vault open"
	if m.vaultPath != "" {
		vaultName = m.vaultPath
	}
	dirty := ""
	if m.mgr != nil && m.mgr.IsModified() {
		dirty = " ●"
	}
	header := headerStyle.Render("  Abditum  " + vaultName + dirty)

	// Message bar
	var msgText string
	if msg := m.messages.Current(); msg != nil {
		msgText = msg.Text
	}
	msgBar := msgBarStyle.Render("  " + msgText)

	// Work area: delegate to active child
	var workContent string
	switch m.area {
	case workAreaWelcome:
		if m.welcome != nil {
			m.welcome.SetSize(m.width, workH)
			workContent = m.welcome.View()
		}
	case workAreaVault:
		workContent = m.renderVaultArea(workH)
	case workAreaTemplates:
		workContent = m.renderTemplatesArea(workH)
	case workAreaSettings:
		if m.settings != nil {
			m.settings.SetSize(m.width, workH)
			workContent = m.settings.View()
		}
	}
	workArea := workAreaStyle.Height(workH).Render(workContent)

	// Command bar: use modal shortcuts when modal active
	var cmdBarContent string
	if len(m.modals) > 0 {
		cmdBarContent = renderShortcuts(m.modals[len(m.modals)-1].Shortcuts(), m.width)
	} else {
		cmdBarContent = m.actions.RenderCommandBar(m.width)
	}
	cmdBar := cmdBarStyle.Render(cmdBarContent)

	return strings.Join([]string{header, msgBar, workArea, cmdBar}, "\n")
}

// renderVaultArea renders workAreaVault: vaultTree (left) + secretDetail (right).
func (m *rootModel) renderVaultArea(workH int) string {
	halfW := m.width / 2
	if m.vaultTree != nil {
		m.vaultTree.SetSize(halfW, workH)
	}
	if m.secretDetail != nil {
		m.secretDetail.SetSize(m.width-halfW, workH)
	}

	left := "[vault tree - Phase 7]"
	right := "[secret detail - Phase 8]"
	if m.vaultTree != nil {
		left = m.vaultTree.View()
	}
	if m.secretDetail != nil {
		right = m.secretDetail.View()
	}

	leftStyle := lipgloss.NewStyle().Width(halfW).Height(workH)
	rightStyle := lipgloss.NewStyle().Width(m.width - halfW).Height(workH)
	return lipgloss.JoinHorizontal(lipgloss.Top, leftStyle.Render(left), rightStyle.Render(right))
}

// renderTemplatesArea renders workAreaTemplates: templateList (left) + templateDetail (right).
func (m *rootModel) renderTemplatesArea(workH int) string {
	halfW := m.width / 2
	if m.templateList != nil {
		m.templateList.SetSize(halfW, workH)
	}
	if m.templateDetail != nil {
		m.templateDetail.SetSize(m.width-halfW, workH)
	}

	left := "[template list - Phase 8]"
	right := "[template detail - Phase 8]"
	if m.templateList != nil {
		left = m.templateList.View()
	}
	if m.templateDetail != nil {
		right = m.templateDetail.View()
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
	m.vaultTree.SetSize(m.width/2, m.height-4)
	m.secretDetail.SetSize(m.width-m.width/2, m.height-4)
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}
