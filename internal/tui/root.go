package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/actions"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/screen"
	"github.com/useful-toys/abditum/internal/tui/secret"
	"github.com/useful-toys/abditum/internal/tui/settings"
	tmpl "github.com/useful-toys/abditum/internal/tui/template"
	"github.com/useful-toys/abditum/internal/vault"
)

// RootModel is the main Bubble Tea model for the application.
// Coordinates 4 fixed screen regions, active work area, and modal stack.
type RootModel struct {
	// width and height are current terminal dimensions, updated in real-time.
	width  int
	height int
	// theme is the active visual theme, applied to all child components.
	theme *design.Theme

	// headerView is the header region — always present, implements ChildView.
	headerView screen.HeaderView
	// messageLineView is the status message bar — manages state and renders messages.
	messageLineView screen.MessageLineView
	// actionLineView is the context action bar — stateless renderer.
	actionLineView screen.ActionLineView

	// workArea indicates which work area is currently displayed.
	workArea design.WorkArea
	// activeView points to the ChildView with focus in the current work area.
	// Never nil after NewRootModel — initialized with &welcomeView.
	activeView ChildView

	// welcomeView is displayed when no vault is open (vaultManager == nil).
	// Stored as direct value (not pointer) — addressable for assignment to activeView.
	welcomeView screen.WelcomeView

	// Views below depend on vaultManager and are nil until initVaultViews is called.
	settingsView   *settings.SettingsView
	secretTree     *secret.VaultTreeView
	secretDetail   *secret.SecretDetailView
	templateList   *tmpl.TemplateListView
	templateDetail *tmpl.TemplateDetailView

	// vaultManager is the active vault manager, or nil if no vault is loaded.
	vaultManager *vault.Manager

	// modals is the stack of open modals; the top of stack is the active modal.
	modals []ModalView

	// activeOperation é a operação em andamento, ou nil quando ocioso.
	// Recebe todas as mensagens não tratadas pelo switch do root.
	activeOperation Operation

	// systemActions stores global/system-level actions evaluated in any context.
	systemActions []actions.Action
	// appActions stores application-wide actions evaluated only without active modal.
	appActions []actions.Action
	// actionGroups groups actions for display in help modal.
	actionGroups []actions.ActionGroup

	// lastActionAt records the time of last user interaction.
	lastActionAt time.Time
	// version is the application version, normally injected via ldflags in build.
	version string

	// initialCmd é um comando opcional a ser emitido junto com Init().
	// Configurado via SetInitialCommand antes de iniciar o loop Tea.
	initialCmd tea.Cmd
}

// Manager returns the active vault manager, or nil if no vault is loaded.
// Implements the AppState interface.
func (r *RootModel) Manager() *vault.Manager {
	return r.vaultManager
}

// ToggleTheme alternates active theme between TokyoNight and Cyberpunk.
// Exported for use by the actions package.
func (r *RootModel) ToggleTheme() {
	if r.theme == design.TokyoNight {
		r.theme = design.Cyberpunk
	} else {
		r.theme = design.TokyoNight
	}
}

// MessageController returns the interface for controlling status messages.
// Allows views and actions to set busy, success, error, warning, info, and hint messages.
func (r *RootModel) MessageController() MessageController {
	return &r.messageLineView
}

// ActiveViewActions retorna todas as actions do contexto atual.
// Combina system actions, application actions e actions da view ativa.
// Usado pelo modal de Ajuda para exibir todos os atalhos disponíveis.
func (r *RootModel) ActiveViewActions() []actions.Action {
	viewActions := r.activeView.Actions()
	all := make([]actions.Action, 0,
		len(r.systemActions)+len(r.appActions)+len(viewActions))
	all = append(all, r.systemActions...)
	all = append(all, r.appActions...)
	all = append(all, viewActions...)
	return all
}

// ActiveViewActionsForBar retorna as actions filtradas e ordenadas para exibição
// na barra de comandos. Filtra por Visible e AvailableWhen; ordena por Priority crescente.
func (r *RootModel) ActiveViewActionsForBar() []actions.Action {
	all := r.ActiveViewActions()

	// Filtrar: Visible == true E (AvailableWhen == nil OU AvailableWhen satisfeita)
	filtered := all[:0]
	for _, a := range all {
		if !a.Visible {
			continue
		}
		if a.AvailableWhen != nil && !a.AvailableWhen(r, r.activeView) {
			continue
		}
		filtered = append(filtered, a)
	}

	// Ordenar por Priority crescente (insertion sort — lista pequena)
	for i := 1; i < len(filtered); i++ {
		for j := i; j > 0 && filtered[j].Priority < filtered[j-1].Priority; j-- {
			filtered[j], filtered[j-1] = filtered[j-1], filtered[j]
		}
	}
	return filtered
}

// GetActionGroups retorna a lista de grupos de ações registrados.
// Exportado para uso pelo package actions.
func (r *RootModel) GetActionGroups() []actions.ActionGroup {
	return r.actionGroups
}

// RegisterActionGroup adiciona um grupo de ações ao root.
func (r *RootModel) RegisterActionGroup(group actions.ActionGroup) {
	r.actionGroups = append(r.actionGroups, group)
}

// RegisterSystemActions adiciona actions de sistema ao root.
func (r *RootModel) RegisterSystemActions(acts []actions.Action) {
	r.systemActions = append(r.systemActions, acts...)
}

// RegisterApplicationActions adiciona actions de aplicação ao root.
func (r *RootModel) RegisterApplicationActions(acts []actions.Action) {
	r.appActions = append(r.appActions, acts...)
}

// evalActions itera através de uma lista de actions e executa a primeira que corresponde
// à tecla pressionada e cuja pré-condição é satisfeita.
func (r *RootModel) evalActions(msg tea.KeyMsg, acts []actions.Action) (tea.Cmd, bool) {
	for _, action := range acts {
		if !action.Matches(msg) {
			continue
		}
		if action.AvailableWhen != nil && !action.AvailableWhen(r, r.activeView) {
			continue
		}
		return action.OnExecute(), true
	}
	return nil, false
}

// initVaultViews creates views that depend on vaultManager.
// Called when vaultManager is available — at initialization or during lifecycle.
func (r *RootModel) initVaultViews() {
	r.settingsView = settings.NewSettingsView(r.vaultManager)
	r.secretTree = secret.NewVaultTreeView(r.vaultManager)
	r.secretDetail = secret.NewSecretDetailView(r.vaultManager)
	r.templateList = tmpl.NewTemplateListView(r.vaultManager)
	r.templateDetail = tmpl.NewTemplateDetailView(r.vaultManager)
	r.headerView.SetVault(r.vaultManager)
}

// setVaultManager sets the active vault manager and synchronizes the header.
// When vault is set, initVaultViews is called. When vault is cleared (nil), header is notified.
func (r *RootModel) setVaultManager(vm *vault.Manager) {
	r.vaultManager = vm
	if vm != nil {
		r.initVaultViews()
	} else {
		r.headerView.SetVault(nil)
	}
}

// renderWorkArea returns the rendered string of the active work area.
// Uses design height constants to calculate available space.
// Note: Render argument order is (height, width) — do not swap.
func (r *RootModel) renderWorkArea() string {
	h := r.height - design.HeaderHeight - design.MessageHeight - design.ActionHeight
	w := r.width

	switch r.workArea {
	case design.WorkAreaWelcome:
		return r.welcomeView.Render(h, w, r.theme)
	case design.WorkAreaSettings:
		return r.settingsView.Render(h, w, r.theme)
	case design.WorkAreaVault:
		treeWidth := int(float64(w) * design.PanelTreeRatio)
		detailWidth := w - treeWidth
		return lipgloss.JoinHorizontal(lipgloss.Top,
			r.secretTree.Render(h, treeWidth, r.theme),
			r.secretDetail.Render(h, detailWidth, r.theme),
		)
	case design.WorkAreaTemplates:
		listWidth := int(float64(w) * design.PanelTreeRatio)
		detailWidth := w - listWidth
		return lipgloss.JoinHorizontal(lipgloss.Top,
			r.templateList.Render(h, listWidth, r.theme),
			r.templateDetail.Render(h, detailWidth, r.theme),
		)
	default:
		return r.welcomeView.Render(h, w, r.theme)
	}
}

// View generates the current visual representation of the application.
// The base layout (4 regions) is always rendered.
// If an active modal exists, it is overlaid on the base via lipgloss v2 compositor.
// AltScreen is always enabled to isolate the TUI from the terminal scroll history.
func (r *RootModel) View() tea.View {
	// Guard views shown before the terminal size is known use AltScreen too,
	// so the application never leaks content into the terminal scroll buffer.
	makeView := func(content string) tea.View {
		v := tea.NewView(content)
		v.AltScreen = true
		return v
	}

	if r.width == 0 || r.height == 0 {
		return makeView("Aguarde...")
	}
	if r.width < design.MinWidth {
		return makeView("Aumente a largura do terminal!")
	}
	if r.height < design.MinHeight {
		return makeView("Aumente a altura do terminal!")
	}

	base := lipgloss.JoinVertical(lipgloss.Left,
		r.headerView.Render(design.HeaderHeight, r.width, r.theme),
		r.renderWorkArea(),
		r.messageLineView.Render(r.width, r.theme),
		r.actionLineView.Render(r.width, r.theme, r.ActiveViewActionsForBar()),
	)

	if len(r.modals) > 0 {
		top := r.modals[len(r.modals)-1]
		// 1 line padding above and below modal on screen.
		modalH := r.height - 2
		modalContent := top.Render(modalH, r.width, r.theme)
		// Center modal content horizontally within available space.
		centeredModal := lipgloss.Place(r.width, modalH, lipgloss.Center, lipgloss.Center, modalContent)

		// Calcular posição absoluta do modal para repassar ao método Cursor do modal.
		// modalContent tem dimensões reais; o compositor aplica Y(1) — offset de 1 linha.
		modalW, modalActualH := lipgloss.Size(modalContent)
		topY := 1 + (modalH-modalActualH)/2 // 1 = offset Y do layer compositor
		leftX := (r.width - modalW) / 2

		// Compose modal (z=1) over base layout (z=0) using lipgloss v2 compositor.
		result := lipgloss.NewCompositor(
			lipgloss.NewLayer(base),
			lipgloss.NewLayer(centeredModal).Y(1).Z(1),
		).Render()
		v := tea.NewView(result)
		v.AltScreen = true
		v.BackgroundColor = lipgloss.Color(r.theme.Surface.Base)
		if c := top.Cursor(topY, leftX); c != nil {
			v.Cursor = c
		}
		return v
	}

	return makeView(base)
}

// Update processes Bubble Tea messages and updates model state.
func (r *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.width = msg.Width
		r.height = msg.Height
		return r, nil

	case OpenModalMsg:
		r.modals = append(r.modals, msg.Modal)
		return r, nil

	case CloseModalMsg:
		if len(r.modals) > 0 {
			r.modals = r.modals[:len(r.modals)-1]
		}
		return r, nil

	case screen.WorkAreaChangedMsg:
		r.setWorkArea(msg.Area)
		return r, nil

	case StartOperationMsg:
		// TODO: quando Operation ganhar Cancel(), chamar r.activeOperation.Cancel() aqui
		// antes de substituir — operações com goroutines de IO/criptografia podem vazar.
		r.activeOperation = msg.Op
		return r, msg.Op.Init()

	case OperationCompletedMsg:
		r.activeOperation = nil
		return r, nil

	case VaultOpenedMsg:
		if msg.Manager == nil {
			return r, nil
		}
		r.setVaultManager(msg.Manager)
		r.setWorkArea(design.WorkAreaVault)
		return r, nil

	case SecretExportedMsg:
		// TODO: mostrar confirmação ao usuário após exportação de segredo.
		return r, nil

	case ModalReadyMsg:
		if len(r.modals) > 1 {
			// Modais pai não recebem ModalReadyMsg — apenas a view ativa recebe.
			return r, nil
		}
		return r, r.activeView.Update(msg)

	case tea.KeyMsg:
		// 1. System actions — evaluated always, including with active modal.
		if cmd, ok := r.evalActions(msg, r.systemActions); ok {
			return r, cmd
		}

		// 2. Active modal receives the key.
		if len(r.modals) > 0 {
			top := len(r.modals) - 1
			return r, r.modals[top].HandleKey(msg)
		}

		// 3. View actions — evaluated only without active modal.
		viewActions := r.activeView.Actions()
		for _, action := range viewActions {
			if !action.Matches(msg) {
				continue
			}
			if action.AvailableWhen != nil && !action.AvailableWhen(r, r.activeView) {
				continue
			}
			return r, action.OnExecute()
		}

		// 4. Application actions — evaluated after view actions.
		if cmd, ok := r.evalActions(msg, r.appActions); ok {
			return r, cmd
		}

		return r, nil
	}

	var cmds []tea.Cmd

	// Roteia para activeOperation ANTES da bifurcação modal/view.
	// tea.Batch não garante ordem: mensagens privadas da operação (ex: mensagens de
	// confirmação emitidas internamente pela operação) podem chegar antes de CloseModalMsg
	// no mesmo Batch. Se aguardarmos o modal sair
	// da pilha para rotear, a operação nunca receberia essas mensagens.
	if r.activeOperation != nil {
		if cmd := r.activeOperation.Update(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// TickMsg sempre vai para messageLineView para animar o spinner.
	// Também re-agenda o timer — tea.Every dispara uma única vez e para.
	// Sem re-agendar, o spinner animaria apenas uma vez e pararia.
	if _, isTickMsg := msg.(TickMsg); isTickMsg {
		r.messageLineView.Update(msg)
		cmds = append(cmds, tickCmd())
	}

	if len(r.modals) > 0 {
		top := len(r.modals) - 1
		// tea.KeyMsg já foi tratado no case acima — aqui chegam mensagens não-key.
		// Modais só recebem eventos de mouse além de teclas.
		if mouseMsg, ok := msg.(tea.MouseMsg); ok {
			cmds = append(cmds, r.modals[top].HandleMouse(mouseMsg))
		}
	} else {
		cmds = append(cmds, r.activeView.Update(msg))
		cmds = append(cmds, r.headerView.Update(msg))
	}
	return r, tea.Batch(cmds...)
}

// setWorkArea troca a área de trabalho ativa e sincroniza o estado do cabeçalho.
// Deve ser chamado sempre que a work area mudar — inclusive na abertura de cofre.
func (r *RootModel) setWorkArea(area design.WorkArea) {
	r.workArea = area
	r.headerView.SetActiveMode(area)
}

// Init is called once at application startup.
// Returns a command that emits TickMsg every second for spinner animation and TTL decrement.
// Se initialCmd foi configurado, é emitido junto para disparar a operação inicial.
func (r *RootModel) Init() tea.Cmd {
	cmds := []tea.Cmd{tickCmd()}
	if r.initialCmd != nil {
		cmds = append(cmds, r.initialCmd)
	}
	return tea.Batch(cmds...)
}

// SetInitialCommand define um comando a ser emitido junto com Init().
// Deve ser chamado antes de iniciar o loop Tea.
// Usado por main.go para disparar automaticamente uma operação via --vault.
func (r *RootModel) SetInitialCommand(cmd tea.Cmd) {
	r.initialCmd = cmd
}

// tickCmd agenda um único TickMsg para daqui a 1 segundo.
// Deve ser re-agendado a cada TickMsg recebido — tea.Every dispara apenas uma vez.
func tickCmd() tea.Cmd {
	return tea.Every(1*time.Second, func(time.Time) tea.Msg {
		return TickMsg{}
	})
}

// RootModelOption is a configuration function applied to RootModel at creation.
type RootModelOption func(*RootModel)

// WithVersion defines the application version displayed in interface.
func WithVersion(version string) RootModelOption {
	return func(m *RootModel) {
		m.version = version
	}
}

// NewRootModel creates and initializes a RootModel with default TokyoNight theme.
// activeView is initialized with &welcomeView — never nil after this function.
func NewRootModel(opts ...RootModelOption) *RootModel {
	m := &RootModel{
		theme:    design.TokyoNight,
		workArea: design.WorkAreaWelcome,
		version:  "dev",
	}
	m.activeView = &m.welcomeView
	for _, opt := range opts {
		opt(m)
	}
	if m.vaultManager != nil {
		m.initVaultViews()
	}
	return m
}
