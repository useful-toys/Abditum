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
	area        workArea
	mgr         *vault.Manager
	vaultPath   string
	initialPath string // Path passed via CLI, for fast-path
	isDirty     bool
	width       int
	height      int
	theme       *Theme
	header      headerModel

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
// Optional initialPath parameter enables CLI fast-path for vault opening.
func NewRootModel(initialPath ...string) *rootModel {
	path := ""
	if len(initialPath) > 0 {
		path = initialPath[0]
	}
	return newRootModel(path)
}

// newRootModel constructs a fully initialized rootModel in PoC mode.
// mgr is nil, vaultPath is "", area is workAreaWelcome (D-02, D-03, D-05).
// If initialPath is non-empty, the CLI fast-path will be initiated in Init().
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
		theme:        ThemeTokyoNight,
		header:       headerModel{},
	}

	m.welcome = newWelcomeModel(actions, m.theme)

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
		// Vault actions (pre-vault scope)
		Action{Keys: []string{"o"}, Label: "Abrir", Description: "Abrir cofre existente",
			Group: 4, Scope: ScopeLocal, Priority: 95, HideFromBar: false,
			Enabled: func() bool { return m.area == workAreaWelcome },
			Handler: func() tea.Cmd {
				flow := newOpenVaultFlow(m.mgr, m.messages, actions, m.theme)
				return func() tea.Msg { return startFlowMsg{flow: flow} }
			}},
		Action{Keys: []string{"n"}, Label: "Novo", Description: "Criar novo cofre",
			Group: 4, Scope: ScopeLocal, Priority: 94, HideFromBar: false,
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
	actions.RegisterGroupLabel(1, "Mensagens")
	actions.RegisterGroupLabel(2, "Status")
	actions.RegisterGroupLabel(3, "Diálogos")
	actions.RegisterGroupLabel(4, "Cofre")

	// Group 3 — Dialog PoC (5 severidades × 3 nº ações)
	actions.Register(m,
		// Destrutivo — 1, 2, 3 ações
		Action{Keys: []string{"1"}, Label: "Dest·1", Description: "Destrutivo 1 ação: Exclusão concluída",
			Group: 3, Scope: ScopeLocal, Priority: 24, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd {
				return Acknowledge(SeverityDestructive, "Exclusão concluída", "Gmail foi excluído permanentemente.", nil)
			}},
		Action{Keys: []string{"2"}, Label: "Dest·2", Description: "Destrutivo 2 ações: Excluir segredo",
			Group: 3, Scope: ScopeLocal, Priority: 23, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd {
				return Decision(SeverityDestructive, "Excluir segredo",
					"Gmail será excluído permanentemente. Esta ação não pode ser desfeita.",
					DecisionAction{Key: "Enter", Label: "Excluir", Default: true},
					nil,
					DecisionAction{Key: "Esc", Label: "Cancelar"})
			}},
		Action{Keys: []string{"3"}, Label: "Dest·3", Description: "Destrutivo 3 ações: Excluir pasta",
			Group: 3, Scope: ScopeLocal, Priority: 22, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd {
				return Decision(SeverityDestructive, "Excluir pasta",
					"Financeiro e todos os seus segredos serão excluídos permanentemente.",
					DecisionAction{Key: "Enter", Label: "Excluir", Default: true},
					[]DecisionAction{{Key: "M", Label: "Mover conteúdo"}},
					DecisionAction{Key: "Esc", Label: "Cancelar"})
			}},

		// Erro — 1, 2, 3 ações
		Action{Keys: []string{"4"}, Label: "Err·1", Description: "Erro 1 ação: Falha ao salvar",
			Group: 3, Scope: ScopeLocal, Priority: 21, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd {
				return Acknowledge(SeverityError, "Falha ao salvar", "Não foi possível salvar o cofre. O arquivo pode estar em uso por outro processo.", nil)
			}},
		Action{Keys: []string{"5"}, Label: "Err·2", Description: "Erro 2 ações: Senha incorreta",
			Group: 3, Scope: ScopeLocal, Priority: 20, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd {
				return Decision(SeverityError, "Senha incorreta",
					"A senha está incorreta. O cofre não pôde ser aberto.",
					DecisionAction{Key: "Enter", Label: "Tentar novamente", Default: true},
					nil,
					DecisionAction{Key: "Esc", Label: "Cancelar"})
			}},
		Action{Keys: []string{"6"}, Label: "Err·3", Description: "Erro 3 ações: Cofre corrompido",
			Group: 3, Scope: ScopeLocal, Priority: 19, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd {
				return Decision(SeverityError, "Cofre corrompido",
					"O arquivo está corrompido. Deseja tentar recuperar a partir do backup?",
					DecisionAction{Key: "Enter", Label: "Recuperar", Default: true},
					[]DecisionAction{{Key: "A", Label: "Abrir backup"}},
					DecisionAction{Key: "Esc", Label: "Cancelar"})
			}},

		// Alerta — 1, 2, 3 ações
		Action{Keys: []string{"7"}, Label: "Ale·1", Description: "Alerta 1 ação: Sessão bloqueada",
			Group: 3, Scope: ScopeLocal, Priority: 18, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd {
				return Acknowledge(SeverityAlert, "Sessão bloqueada", "O cofre foi bloqueado após 5 minutos de inatividade.", nil)
			}},
		Action{Keys: []string{"8"}, Label: "Ale·2", Description: "Alerta 2 ações: Alterações não salvas",
			Group: 3, Scope: ScopeLocal, Priority: 17, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd {
				return Decision(SeverityAlert, "Alterações não salvas",
					"Existem alterações não salvas. Sair irá descartá-las.",
					DecisionAction{Key: "Enter", Label: "Descartar", Default: true},
					nil,
					DecisionAction{Key: "Esc", Label: "Voltar"})
			}},
		Action{Keys: []string{"9"}, Label: "Ale·3", Description: "Alerta 3 ações: Senha fraca",
			Group: 3, Scope: ScopeLocal, Priority: 16, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd {
				return Decision(SeverityAlert, "Senha fraca",
					"A senha mestra é fraca e pode ser facilmente descoberta.",
					DecisionAction{Key: "Enter", Label: "Usar assim mesmo", Default: true},
					[]DecisionAction{{Key: "T", Label: "Trocar senha"}},
					DecisionAction{Key: "Esc", Label: "Cancelar"})
			}},

		// Informativo — 1, 2, 3 ações
		Action{Keys: []string{"a"}, Label: "Inf·1", Description: "Informativo 1 ação: Cofre criado",
			Group: 3, Scope: ScopeLocal, Priority: 15, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd {
				return Acknowledge(SeverityInformative, "Cofre criado", "O cofre foi criado com sucesso em ~/documentos/pessoal.abditum.", nil)
			}},
		Action{Keys: []string{"b"}, Label: "Inf·2", Description: "Informativo 2 ações: Conflito detectado",
			Group: 3, Scope: ScopeLocal, Priority: 14, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd {
				return Decision(SeverityInformative, "Conflito detectado",
					"O arquivo foi modificado externamente desde a última abertura.",
					DecisionAction{Key: "Enter", Label: "Sobrescrever", Default: true},
					nil,
					DecisionAction{Key: "Esc", Label: "Cancelar"})
			}},
		Action{Keys: []string{"c"}, Label: "Inf·3", Description: "Informativo 3 ações: Importação concluída",
			Group: 3, Scope: ScopeLocal, Priority: 13, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd {
				return Decision(SeverityInformative, "Importação concluída",
					"12 segredos importados. 3 entradas já existentes foram atualizadas.",
					DecisionAction{Key: "Enter", Label: "Ver detalhes", Default: true},
					[]DecisionAction{{Key: "F", Label: "Fechar"}},
					DecisionAction{Key: "Esc", Label: "OK"})
			}},

		// Neutro — 1, 2, 3 ações
		Action{Keys: []string{"d"}, Label: "Neu·1", Description: "Neutro 1 ação: Operação concluída",
			Group: 3, Scope: ScopeLocal, Priority: 12, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd {
				return Acknowledge(SeverityNeutral, "Operação concluída", "A exportação foi salva em ~/documentos/backup-2026-04-05.json.", nil)
			}},
		Action{Keys: []string{"e"}, Label: "Neu·2", Description: "Neutro 2 ações: Sair do Abditum",
			Group: 3, Scope: ScopeLocal, Priority: 11, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd {
				return Decision(SeverityNeutral, "Sair do Abditum",
					"Deseja sair? Todas as alterações não salvas serão descartadas.",
					DecisionAction{Key: "Enter", Label: "Sair", Default: true},
					nil,
					DecisionAction{Key: "Esc", Label: "Cancelar"})
			}},
		Action{Keys: []string{"f"}, Label: "Neu·3", Description: "Neutro 3 ações: Salvar cofre",
			Group: 3, Scope: ScopeLocal, Priority: 10, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd {
				return Decision(SeverityNeutral, "Salvar cofre",
					"Deseja salvar as alterações antes de continuar?",
					DecisionAction{Key: "Enter", Label: "Salvar", Default: true},
					[]DecisionAction{{Key: "N", Label: "Não salvar"}},
					DecisionAction{Key: "Esc", Label: "Voltar"})
			}},
	)

	return m
}

// Init satisfies tea.Model. Always starts global tick for message TTL (D-10, D-11).
// If initialPath is set (CLI fast-path), start openVaultFlow immediately.
func (m *rootModel) Init() tea.Cmd {
	tickCmd := tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })

	// CLI fast-path: if initialPath is non-empty, start openVaultFlow
	if m.initialPath != "" {
		return tea.Batch(
			tickCmd,
			func() tea.Msg {
				// Create temporary vault manager for the flow
				// It will be populated when vault is opened
				flow := newOpenVaultFlow(nil, m.messages, m.actions, m.theme)
				flow.cliPath = m.initialPath // Set CLI path for fast-path
				return startFlowMsg{flow: flow}
			},
		)
	}

	return tickCmd
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

	// --- Vault opened: transition to work area and store vault (D-08) ---
	case vaultOpenedMsg:
		// TODO: In Phase 9+, populate m.mgr with the opened vault
		// For now, just transition to vault area
		m.area = workAreaVault
		m.vaultPath = msg.Path
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

	// --- Keyboard input: D-09 dispatch order ---
	case tea.KeyPressMsg:
		m.messages.HandleInput()
		m.lastActionAt = time.Now()
		key := msg.String()
		inFlowOrModal := m.activeFlow != nil || len(m.modals) > 0

		// Check for Ctrl+Q (exit flow) before any other key handling
		if key == "ctrl+q" {
			// If there are unsaved changes, prompt the user
			if m.mgr != nil && m.mgr.IsModified() {
				return m, func() tea.Msg {
					return Decision(SeverityNeutral, "Alterações não salvas",
						"Deseja salvar as alterações antes de sair?",
						DecisionAction{Key: "Enter", Label: "Salvar", Default: true},
						[]DecisionAction{{Key: "D", Label: "Descartar"}},
						DecisionAction{Key: "Esc", Label: "Voltar"})
				}
			}
			// No unsaved changes, exit immediately
			return m, tea.Quit
		}

		// Check for F12 theme toggle before any other key handling
		if key == "f12" {
			if m.theme == ThemeTokyoNight {
				m.theme = ThemeCyberpunk
			} else {
				m.theme = ThemeTokyoNight
			}
			m.applyTheme()
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

	cmdBarStyle := lipgloss.NewStyle().Width(m.width).Background(m.theme.SurfaceBase)
	workAreaStyle := lipgloss.NewStyle().Width(m.width).Background(m.theme.SurfaceBase)

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
		cmdBarContent = RenderCommandBar(m.actions.Visible(), m.width)
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
	m.vaultTree = newVaultTreeModel(m.mgr, m.actions, m.messages, m.theme)
	m.secretDetail = newSecretDetailModel(m.mgr, m.actions, m.messages, m.theme)
	m.vaultTree.SetSize(m.width/2, m.height-4)
	m.secretDetail.SetSize(m.width-m.width/2, m.height-4)
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}

// applyTheme propagates the current theme to all active child models and modals.
func (m *rootModel) applyTheme() {
	for _, child := range m.liveWorkChildren() {
		child.ApplyTheme(m.theme)
	}
	for _, modal := range m.modals {
		// Modals should implement ApplyTheme if they need theme changes.
		if themeableModal, ok := modal.(interface{ ApplyTheme(*Theme) }); ok {
			themeableModal.ApplyTheme(m.theme)
		}
	}
}
