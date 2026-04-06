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

// NewRootModel is the exported constructor for main.go (PoC mode, D-04).
func NewRootModel() *rootModel {
	return newRootModel()
}

// newRootModel constructs a fully initialized rootModel in PoC mode.
// mgr is nil, vaultPath is "", area is workAreaWelcome (D-02, D-03, D-05).
func newRootModel() *rootModel {
	actions := NewActionManager()
	messages := NewMessageManager()

	m := &rootModel{
		area:         workAreaWelcome,
		mgr:          nil, // PoC mode — no vault (D-02)
		vaultPath:    "",
		actions:      actions,
		messages:     messages,
		lastActionAt: time.Now(),
	}

	m.welcome = newWelcomeModel(actions)

	// Register all 15 PoC actions on rootModel as owner (D-19 through D-23).
	actions.Register(m,
		// Group 1 "Mensagens" — message type demonstrations
		Action{Keys: []string{"f2"}, Label: "Dica uso", Description: "Mostrar MsgHint permanente",
			Group: 1, Scope: ScopeLocal, Priority: 15, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.messages.Show(MsgHint, "Dica de uso permanente", 0, false); return nil }},
		Action{Keys: []string{"f3"}, Label: "Dica campo", Description: "Mostrar MsgHint de campo",
			Group: 1, Scope: ScopeLocal, Priority: 20, HideFromBar: true,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.messages.Show(MsgHint, "Dica de campo permanente", 0, false); return nil }},
		Action{Keys: []string{"f4"}, Label: "Info", Description: "Mostrar MsgInfo",
			Group: 1, Scope: ScopeLocal, Priority: 30, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.messages.Show(MsgInfo, "Informação neutra", 0, false); return nil }},
		Action{Keys: []string{"f5"}, Label: "Alerta", Description: "Mostrar MsgWarn",
			Group: 1, Scope: ScopeLocal, Priority: 40, HideFromBar: true,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.messages.Show(MsgWarn, "Alerta de atenção", 0, false); return nil }},
		Action{Keys: []string{"f6"}, Label: "Erro", Description: "Mostrar MsgError",
			Group: 1, Scope: ScopeLocal, Priority: 50, HideFromBar: false,
			Enabled: func() bool { return false },
			Handler: func() tea.Cmd { m.messages.Show(MsgError, "Erro de operação", 0, false); return nil }},

		// Group 2 "Status" — operation state demonstrations
		Action{Keys: []string{"f7"}, Label: "Ocupado", Description: "Mostrar MsgBusy com spinner",
			Group: 2, Scope: ScopeLocal, Priority: 60, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.messages.Show(MsgBusy, "Processando...", 0, false); return nil }},
		Action{Keys: []string{"f8"}, Label: "Sucesso", Description: "Mostrar MsgSuccess",
			Group: 2, Scope: ScopeLocal, Priority: 70, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.messages.Show(MsgSuccess, "Operação concluída", 0, false); return nil }},

		// No group — utility actions
		Action{Keys: []string{"f9"}, Label: "Limpar", Description: "Limpar barra de mensagens",
			Group: 0, Scope: ScopeLocal, Priority: 80, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.messages.Clear(); return nil }},
		Action{Keys: []string{"f10", "f11"}, Label: "Truncar", Description: "Testar truncação de mensagem longa",
			Group: 0, Scope: ScopeLocal, Priority: 90, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd {
				m.messages.Show(MsgInfo, "Esta é uma mensagem muito longa com mais de cem caracteres para testar o sistema de truncação com reticências quando a mensagem excede a largura disponível da barra de mensagens do sistema.", 0, false)
				return nil
			}},

		// Shift+Fx variants — HideFromBar: true (visible in Help modal only, TTL 5s)
		Action{Keys: []string{"shift+f2"}, Label: "Dica uso 5s", Description: "MsgHint com TTL de 5s",
			Group: 1, Scope: ScopeLocal, Priority: 14, HideFromBar: true,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.messages.Show(MsgHint, "Dica de uso (5s)", 5, false); return nil }},
		Action{Keys: []string{"shift+f3"}, Label: "Dica campo 5s", Description: "MsgHint campo com TTL de 5s",
			Group: 1, Scope: ScopeLocal, Priority: 19, HideFromBar: true,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.messages.Show(MsgHint, "Dica de campo (5s)", 5, false); return nil }},
		Action{Keys: []string{"shift+f4"}, Label: "Info 5s", Description: "MsgInfo com TTL de 5s",
			Group: 1, Scope: ScopeLocal, Priority: 29, HideFromBar: true,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.messages.Show(MsgInfo, "Informação (5s)", 5, false); return nil }},
		Action{Keys: []string{"shift+f5"}, Label: "Alerta 5s", Description: "MsgWarn com TTL de 5s",
			Group: 1, Scope: ScopeLocal, Priority: 39, HideFromBar: true,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.messages.Show(MsgWarn, "Alerta (5s)", 5, false); return nil }},
		Action{Keys: []string{"shift+f6"}, Label: "Erro 5s", Description: "MsgError com TTL de 5s",
			Group: 1, Scope: ScopeLocal, Priority: 49, HideFromBar: true,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.messages.Show(MsgError, "Erro (5s)", 5, false); return nil }},

		// Navigation/Global actions
		Action{Keys: []string{"ctrl+q"}, Label: "Sair", Description: "Sair do Abditum",
			Group: 0, Scope: ScopeLocal, Priority: 10, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { return tea.Quit }},
		Action{Keys: []string{"f1"}, Label: "Ajuda", Description: "Mostrar atalhos de teclado",
			Group: 1, Scope: ScopeGlobal, Priority: 0, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd {
				return func() tea.Msg { return pushModalMsg{modal: newHelpModal(actions)} }
			}},
	)
	actions.RegisterGroupLabel(1, "Mensagens")
	actions.RegisterGroupLabel(2, "Status")

	return m
}

// Init satisfies tea.Model. Always starts global tick for message TTL (D-10, D-11).
func (m *rootModel) Init() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
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
		top.SetSize(m.width, workH) // pass workH, not total height
		content = m.renderFrame(top)
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

// overlayModal renders a modal dialog centered over the existing frame content.
// The frame remains visible behind/around the modal — only the modal region is replaced.
// Deprecated: kept for reference; replaced by lipgloss.Place inside renderFrame.
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

// renderFrame composes the full frame: header + work area + msg bar + cmd bar.
// If modal is non-nil, it is centered inside the work area using lipgloss.Place.
// The modal must have been sized with SetSize(width, workH) before calling.
func (m *rootModel) renderFrame(modal modalView) string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	headerStyle := lipgloss.NewStyle().Width(m.width).Foreground(lipgloss.Color(ColorAccentPrimary)).Bold(true)
	separatorStyle := lipgloss.NewStyle().Width(m.width).Foreground(lipgloss.Color(ColorBorderDefault))
	cmdBarStyle := lipgloss.NewStyle().Width(m.width)
	workAreaStyle := lipgloss.NewStyle().Width(m.width)

	const headerH = 2
	const msgBarH = 1
	const cmdBarH = 1
	workH := m.height - headerH - msgBarH - cmdBarH
	if workH < 0 {
		workH = 0
	}

	// Header
	header := headerStyle.Render("  Abditum") + "\n" + separatorStyle.Render(strings.Repeat("─", m.width))

	// Message bar
	msgBar := RenderMessageBar(m.messages.Current(), m.width)

	// Work area
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

	// Overlay modal centered inside work area using lipgloss.Place.
	// lipgloss.Place handles ANSI correctly and never touches header/msgBar/cmdBar.
	if modal != nil {
		modalStr := modal.View()
		workArea = lipgloss.Place(m.width, workH, lipgloss.Center, lipgloss.Center, modalStr)
	}

	// Command bar: always render so the frame occupies exactly `height` lines.
	// When no shortcuts, render a blank background line.
	var cmdBarContent string
	if modal != nil {
		cmdBarContent = renderShortcuts(modal.Shortcuts(), m.width)
	} else {
		cmdBarContent = m.actions.RenderCommandBar(m.width)
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
