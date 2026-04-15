package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/screen"
	"github.com/useful-toys/abditum/internal/tui/secret"
	"github.com/useful-toys/abditum/internal/tui/settings"
	tmpl "github.com/useful-toys/abditum/internal/tui/template"
	"github.com/useful-toys/abditum/internal/vault"
)

// WorkArea representa qual área de trabalho está ativa na tela principal.
// É usada por RootModel para decidir qual ChildView exibir.
type WorkArea int

const (
	// WorkAreaWelcome exibe a tela de boas-vindas, para usuários sem cofre aberto.
	WorkAreaWelcome WorkArea = iota
	// WorkAreaSettings exibe as configurações da aplicação.
	WorkAreaSettings
	// WorkAreaVault exibe a área de gerenciamento do cofre de segredos.
	WorkAreaVault
	// WorkAreaTemplates exibe a área de gerenciamento de templates de segredos.
	WorkAreaTemplates
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
	workArea WorkArea
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

	// systemActions are evaluated in any context, including with active modal.
	systemActions []Action
	// applicationActions are evaluated only when no modal is active.
	applicationActions []Action
	// actionGroups groups actions for display in help modal.
	actionGroups []ActionGroup

	// lastActionAt records the time of last user interaction.
	lastActionAt time.Time
	// version is the application version, normally injected via ldflags in build.
	version string
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

// ActiveViewActions returns all actions applicable to the current view context.
// Includes system actions, application actions, and activeView actions.
func (r *RootModel) ActiveViewActions() []Action {
	allActions := make([]Action, 0, len(r.systemActions)+len(r.applicationActions))
	allActions = append(allActions, r.systemActions...)
	allActions = append(allActions, r.applicationActions...)
	// Note: view actions are returned as interface{} so we can't directly append them
	return allActions
}

// GetActionGroups returns the list of registered action groups.
// Exported for use by the actions package.
func (r *RootModel) GetActionGroups() []ActionGroup {
	return r.actionGroups
}

// RegisterActionGroup adds an action group to root.
func (r *RootModel) RegisterActionGroup(group ActionGroup) {
	r.actionGroups = append(r.actionGroups, group)
}

// RegisterSystemActions adds system actions to root.
func (r *RootModel) RegisterSystemActions(actions []Action) {
	r.systemActions = append(r.systemActions, actions...)
}

// RegisterApplicationActions adds application actions to root.
func (r *RootModel) RegisterApplicationActions(actions []Action) {
	r.applicationActions = append(r.applicationActions, actions...)
}

// evalActions iterates through an action list and executes the first one that matches
// the pressed key and whose precondition is satisfied.
func (r *RootModel) evalActions(msg tea.KeyMsg, actions []Action) (tea.Cmd, bool) {
	for _, action := range actions {
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
}

// renderWorkArea returns the rendered string of the active work area.
// Uses design height constants to calculate available space.
// Note: Render argument order is (height, width) — do not swap.
func (r *RootModel) renderWorkArea() string {
	h := r.height - design.HeaderHeight - design.MessageHeight - design.ActionHeight
	w := r.width

	switch r.workArea {
	case WorkAreaWelcome:
		return r.welcomeView.Render(h, w, r.theme)
	case WorkAreaSettings:
		return r.settingsView.Render(h, w, r.theme)
	case WorkAreaVault:
		treeWidth := int(float64(w) * design.PanelTreeRatio)
		detailWidth := w - treeWidth
		return lipgloss.JoinHorizontal(lipgloss.Top,
			r.secretTree.Render(h, treeWidth, r.theme),
			r.secretDetail.Render(h, detailWidth, r.theme),
		)
	case WorkAreaTemplates:
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
func (r *RootModel) View() tea.View {
	if r.width == 0 || r.height == 0 {
		return tea.NewView("Aguarde...")
	}
	if r.width < design.MinWidth {
		return tea.NewView("Aumente a largura do terminal!")
	}
	if r.height < design.MinHeight {
		return tea.NewView("Aumente a altura do terminal!")
	}

	allActions := r.ActiveViewActions()
	actionsInterface := make([]interface{}, len(allActions))
	for i, a := range allActions {
		actionsInterface[i] = a
	}

	base := lipgloss.JoinVertical(lipgloss.Left,
		r.headerView.Render(design.HeaderHeight, r.width, r.theme),
		r.renderWorkArea(),
		r.messageLineView.Render(r.width, r.theme),
		r.actionLineView.Render(design.ActionHeight, r.width, r.theme, actionsInterface),
	)

	if len(r.modals) > 0 {
		top := r.modals[len(r.modals)-1]
		// 1 line padding above and below modal on screen.
		modalH := r.height - 2
		modalContent := top.Render(modalH, r.width, r.theme)
		// Center modal content horizontally within available space.
		centeredModal := lipgloss.Place(r.width, modalH, lipgloss.Center, lipgloss.Center, modalContent)
		// Compose modal (z=1) over base layout (z=0) using lipgloss v2 compositor.
		result := lipgloss.NewCompositor(
			lipgloss.NewLayer(base),
			lipgloss.NewLayer(centeredModal).Y(1).Z(1),
		).Render()
		v := tea.NewView(result)
		v.AltScreen = true
		v.BackgroundColor = lipgloss.Color(r.theme.Surface.Base)
		return v
	}

	return tea.NewView(base)
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

	case ModalReadyMsg:
		if len(r.modals) > 1 {
			parent := r.modals[len(r.modals)-2]
			return r, parent.Update(msg)
		}
		return r, r.activeView.Update(msg)

	case TickMsg:
		// Animate the message bar spinner and decrement TTL.
		return r, r.messageLineView.Update(msg)

	case tea.KeyMsg:
		// 1. System actions — evaluated always, including with active modal.
		if cmd, ok := r.evalActions(msg, r.systemActions); ok {
			return r, cmd
		}

		// 2. Active modal receives the key.
		if len(r.modals) > 0 {
			top := len(r.modals) - 1
			return r, r.modals[top].Update(msg)
		}

		// 3. View actions — evaluated only without active modal.
		viewActions := r.activeView.Actions()
		if viewActions != nil {
			for _, a := range viewActions {
				if action, ok := a.(Action); ok {
					if !action.Matches(msg) {
						continue
					}
					if action.AvailableWhen != nil && !action.AvailableWhen(r, r.activeView) {
						continue
					}
					return r, action.OnExecute()
				}
			}
		}

		// 4. Application actions — evaluated after view actions.
		if cmd, ok := r.evalActions(msg, r.applicationActions); ok {
			return r, cmd
		}

		return r, nil
	}

	if len(r.modals) > 0 {
		top := len(r.modals) - 1
		return r, r.modals[top].Update(msg)
	}

	var cmds []tea.Cmd
	cmds = append(cmds, r.activeView.Update(msg))
	cmds = append(cmds, r.headerView.Update(msg))
	return r, tea.Batch(cmds...)
}

// Init is called once at application startup.
// Returns a command that emits TickMsg every second for spinner animation and TTL decrement.
func (r *RootModel) Init() tea.Cmd {
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
		workArea: WorkAreaWelcome,
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
