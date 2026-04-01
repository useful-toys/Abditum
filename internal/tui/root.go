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
// All child models (preVault, vaultTree, etc.) are stored as concrete
// pointer fields — nil means inactive. Never store children as childModel
// interface values in struct fields (typed nil trap).
type rootModel struct {
	// State machine
	area      workArea
	mgr       *vault.Manager
	vaultPath string
	width     int
	height    int

	// Child models — nil = inactive
	preVault       *preVaultModel
	vaultTree      *vaultTreeModel
	secretDetail   *secretDetailModel
	templateList   *templateListModel
	templateDetail *templateDetailModel
	settings       *settingsModel

	// Modal stack — LIFO; last element = topmost/active.
	// Stored as childModel interface to allow heterogeneous modal types
	// (e.g., *modalModel, *helpModal) without a typed-nil trap at the
	// concrete field level. All elements are non-nil when appended.
	modals []childModel

	// Active flow — nil = no flow running
	activeFlow flowHandler
	flows      *FlowRegistry

	// Shared services
	actions  *ActionManager
	messages *MessageManager

	// Tick tracking
	lastActionAt time.Time
}

// Compile-time assertion: rootModel satisfies tea.Model.
var _ tea.Model = &rootModel{}

// NewRootModel is the exported constructor for main.go (package main cannot
// access unexported newRootModel). It delegates directly to newRootModel.
func NewRootModel(mgr *vault.Manager, initialPath string) *rootModel {
	return newRootModel(mgr, initialPath)
}

// newRootModel constructs a fully initialized rootModel.
// mgr may be nil during tests (Phase 5 has no open vault).
// initialPath is the optional vault path from os.Args (may be empty).
func newRootModel(mgr *vault.Manager, initialPath string) *rootModel {
	actions := NewActionManager()
	messages := NewMessageManager()
	flows := &FlowRegistry{}

	// Register global flows
	flows.Register(openVaultDescriptor{})
	flows.Register(createVaultDescriptor{})

	m := &rootModel{
		area:         workAreaPreVault,
		mgr:          mgr,
		vaultPath:    initialPath,
		flows:        flows,
		actions:      actions,
		messages:     messages,
		lastActionAt: time.Now(),
	}

	// Mount initial work area
	m.preVault = newPreVaultModel(actions)

	// Register global shortcuts into ActionManager
	actions.Register(Action{Key: "ctrl+q", Label: "Quit", Description: "Quit Abditum (confirms if unsaved)", Group: "Global", Priority: 100})
	actions.Register(Action{Key: "?", Label: "Help", Description: "Show keyboard shortcuts", Group: "Global", Priority: 90})

	return m
}

// Init satisfies tea.Model. Returns nil — the global tick does NOT start here.
// Tick starts only on transition to workAreaVault (see enterVault).
func (m *rootModel) Init() tea.Cmd {
	return nil
}

// Update satisfies tea.Model. Implements the dispatch priority rules from D-06:
//  1. Global shortcuts (ctrl+Q, ?)
//  2. activeFlow (if non-nil) receives input
//  3. Topmost modal receives input (if stack non-empty)
//  4. Flow dispatch via FlowRegistry
//  5. Active base child model
//
// Domain messages (tick, vault events, pushModal, popModal) are handled
// before routing to avoid double-dispatch.
func (m *rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// --- Window resize: propagate to all live children ---
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		for _, child := range m.liveModels() {
			child.SetSize(msg.Width, msg.Height)
		}
		return m, nil

	// --- Modal stack management ---
	case pushModalMsg:
		if msg.modal != nil {
			msg.modal.SetSize(m.width, m.height)
			m.modals = append(m.modals, msg.modal)
		}
		return m, nil

	case popModalMsg:
		if len(m.modals) > 0 {
			m.modals = m.modals[:len(m.modals)-1]
		}
		return m, nil

	// --- Global tick ---
	case tickMsg:
		var cmds []tea.Cmd
		// Broadcast tick to all live models for periodic UI updates (e.g., clock)
		for _, child := range m.liveModels() {
			if cmd := child.Update(msg); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		// Re-issue tick to keep loop alive
		cmds = append(cmds, tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) }))
		return m, tea.Batch(cmds...)

	// --- Flow chaining ---
	case chainFlowMsg:
		ctx := m.buildFlowContext()
		if d := m.flows.ForKey(msg.key, ctx); d != nil {
			m.activeFlow = d.New(ctx)
		}
		return m, nil

	// --- Domain messages: broadcast to all live models ---
	case secretAddedMsg, secretDeletedMsg, secretRestoredMsg, secretModifiedMsg,
		secretMovedMsg, secretReorderedMsg, folderStructureChangedMsg,
		vaultSavedMsg, vaultReloadedMsg, vaultClosedMsg, vaultChangedMsg:
		var cmds []tea.Cmd
		for _, child := range m.liveModels() {
			if cmd := child.Update(msg); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		return m, tea.Batch(cmds...)

	// --- Input: apply dispatch priority ---
	case tea.KeyPressMsg:
		m.lastActionAt = time.Now()
		return m.dispatchKey(msg)
	}

	return m, nil
}

// dispatchKey applies the D-06 priority rules for keyboard events.
func (m *rootModel) dispatchKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Priority 1: Global shortcuts — always intercepted first
	switch key {
	case "ctrl+q":
		if m.mgr != nil && m.mgr.IsModified() {
			// Unsaved changes: push confirmation modal
			cmd := NewConfirm(
				"You have unsaved changes. Quit anyway?",
				tea.Quit,
				nil,
			)
			return m, cmd
		}
		return m, tea.Quit

	case "?":
		// Push help modal
		help := newHelpModal(m.actions)
		help.SetSize(m.width, m.height)
		m.modals = append(m.modals, help)
		return m, nil
	}

	// Priority 2: Active flow receives input
	if m.activeFlow != nil {
		cmd := m.activeFlow.Update(msg)
		return m, cmd
	}

	// Priority 3: Topmost modal receives input
	if len(m.modals) > 0 {
		top := m.modals[len(m.modals)-1]
		cmd := top.Update(msg)
		return m, cmd
	}

	// Priority 4: Flow dispatch via FlowRegistry (and child flows)
	ctx := m.buildFlowContext()
	// Check child flows first (escape hatch), then global registry
	var matched flowDescriptor
	if child := m.activeChild(); child != nil {
		for _, fd := range child.ChildFlows() {
			if fd.Key() == key && fd.IsApplicable(ctx) {
				matched = fd
				break
			}
		}
	}
	if matched == nil {
		matched = m.flows.ForKey(key, ctx)
	}
	if matched != nil {
		m.activeFlow = matched.New(ctx)
		cmd := m.activeFlow.Update(msg)
		return m, cmd
	}

	// Priority 5: Active base child model
	if child := m.activeChild(); child != nil {
		cmd := child.Update(msg)
		return m, cmd
	}

	return m, nil
}

// buildFlowContext assembles FlowContext from rootModel + active child state.
func (m *rootModel) buildFlowContext() FlowContext {
	ctx := FlowContext{}
	if m.mgr != nil {
		ctx.VaultOpen = !m.mgr.IsLocked()
		ctx.VaultDirty = m.mgr.IsModified()
	}
	if child := m.activeChild(); child != nil {
		childCtx := child.Context()
		ctx.FocusedFolder = childCtx.FocusedFolder
		ctx.FocusedSecret = childCtx.FocusedSecret
		ctx.SecretOpen = childCtx.SecretOpen
		ctx.FocusedField = childCtx.FocusedField
		ctx.FocusedTemplate = childCtx.FocusedTemplate
		ctx.Mode = childCtx.Mode
	}
	return ctx
}

// activeChild returns the single "base" child model for the current workArea
// (not modals). Returns nil if none is mounted.
func (m *rootModel) activeChild() childModel {
	switch m.area {
	case workAreaPreVault:
		if m.preVault != nil {
			return m.preVault
		}
	case workAreaVault:
		// Primary focus child: vaultTree (left panel)
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

// liveModels returns all non-nil child models (work area + modals) as a
// flat []childModel slice. Uses explicit nil checks on concrete pointer
// fields — never stores as interface to avoid the typed-nil trap.
func (m *rootModel) liveModels() []childModel {
	var live []childModel
	if m.preVault != nil {
		live = append(live, m.preVault)
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
	for _, modal := range m.modals {
		live = append(live, modal)
	}
	return live
}

// View satisfies tea.Model. Composes the full frame and overlays any active modal.
// Returns tea.View (not string) — critical for Bubble Tea v2.
func (m *rootModel) View() tea.View {
	content := m.renderFrame()

	// Overlay topmost modal if stack is non-empty
	if len(m.modals) > 0 {
		top := m.modals[len(m.modals)-1]
		content = top.View()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

// renderFrame composes the base frame zones: header + message bar + work area + command bar.
// Phase 5: placeholder styles; real visual design in Phase 6+.
func (m *rootModel) renderFrame() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	headerStyle := lipgloss.NewStyle().Width(m.width).Background(lipgloss.Color("236")).Foreground(lipgloss.Color("255"))
	msgBarStyle := lipgloss.NewStyle().Width(m.width).Foreground(lipgloss.Color("245"))
	cmdBarStyle := lipgloss.NewStyle().Width(m.width).Background(lipgloss.Color("236"))
	workAreaStyle := lipgloss.NewStyle().Width(m.width)

	// Fixed zone heights
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
	msgBar := msgBarStyle.Render("  " + m.messages.Current())

	// Work area: delegate to active child
	var workContent string
	switch m.area {
	case workAreaPreVault:
		if m.preVault != nil {
			m.preVault.SetSize(m.width, workH)
			workContent = m.preVault.View()
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

	// Command bar
	cmdBar := cmdBarStyle.Render(m.actions.RenderCommandBar(m.width))

	return strings.Join([]string{header, msgBar, workArea, cmdBar}, "\n")
}

// renderVaultArea renders workAreaVault: vaultTree (left) + secretDetail (right) side by side.
func (m *rootModel) renderVaultArea(workH int) string {
	halfW := m.width / 2
	if m.vaultTree != nil {
		m.vaultTree.SetSize(halfW, workH)
	}
	if m.secretDetail != nil {
		m.secretDetail.SetSize(m.width-halfW, workH)
	}

	left := "[vault tree — Phase 7]"
	right := "[secret detail — Phase 8]"
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

	left := "[template list — Phase 8]"
	right := "[template detail — Phase 8]"
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
// Phase 5: exported for testing; not called by any flow yet.
func (m *rootModel) enterVault() tea.Cmd {
	m.area = workAreaVault
	m.preVault = nil // GC old model
	m.vaultTree = newVaultTreeModel(m.mgr, m.actions, m.messages)
	m.secretDetail = newSecretDetailModel(m.mgr, m.actions, m.messages)
	m.vaultTree.SetSize(m.width/2, m.height-4)
	m.secretDetail.SetSize(m.width-m.width/2, m.height-4)
	// Start global tick
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}
