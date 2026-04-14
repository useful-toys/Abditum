package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
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

// RootModel é o modelo principal da aplicação Bubble Tea.
// Coordena a exibição da área de trabalho ativa e a pilha de modais abertos.
type RootModel struct {
	// width e height são as dimensões atuais do terminal, atualizadas em tempo real.
	width  int
	height int
	// theme é o tema visual ativo, aplicado a todos os componentes filhos.
	theme *design.Theme
	// workArea indica qual área de trabalho está sendo exibida no momento.
	workArea WorkArea
	// activeView é o componente principal atualmente visível na tela.
	activeView ChildView
	// modals é a pilha de modais abertos; o topo da pilha é o modal ativo.
	modals []ModalView
	// lastActionAt registra o momento da última interação do usuário.
	lastActionAt time.Time
	// version é a versão da aplicação, normalmente injetada via ldflags no build.
	version string
	// vaultManager é o gerenciador do cofre ativo, ou nil se nenhum cofre estiver carregado.
	vaultManager *vault.Manager
	// systemActions são as actions de sistema, avaliadas em qualquer contexto (inclusive com modal ativo).
	systemActions []Action
	// applicationActions são as actions globais da aplicação, avaliadas sem modal ativo.
	applicationActions []Action
	// actionGroups agrupa actions para exibição no modal de ajuda.
	actionGroups []ActionGroup
}

// Manager retorna o vault manager ativo, ou nil se nenhum cofre estiver carregado.
// Implementa a interface AppState.
func (r *RootModel) Manager() *vault.Manager {
	return r.vaultManager
}

// toggleTheme alterna o tema ativo entre TokyoNight e Cyberpunk.
func (r *RootModel) toggleTheme() {
	if r.theme == design.TokyoNight {
		r.theme = design.Cyberpunk
	} else {
		r.theme = design.TokyoNight
	}
}

// RegisterActionGroup adiciona um grupo de actions ao root.
// Chamar múltiplas vezes acumula os grupos — não substitui.
func (r *RootModel) RegisterActionGroup(group ActionGroup) {
	r.actionGroups = append(r.actionGroups, group)
}

// RegisterSystemActions adiciona actions de sistema ao root.
// System actions são avaliadas em qualquer contexto, inclusive com modal ativo.
// Chamar múltiplas vezes acumula as listas — não substitui.
func (r *RootModel) RegisterSystemActions(actions []Action) {
	r.systemActions = append(r.systemActions, actions...)
}

// RegisterApplicationActions adiciona actions de aplicação ao root.
// Application actions são avaliadas apenas quando nenhum modal está ativo.
// Chamar múltiplas vezes acumula as listas — não substitui.
func (r *RootModel) RegisterApplicationActions(actions []Action) {
	r.applicationActions = append(r.applicationActions, actions...)
}

// evalActions percorre uma lista de actions e executa a primeira que corresponda
// à tecla pressionada e cuja pré-condição esteja satisfeita.
// Retorna o Cmd resultante e true se uma action foi executada; nil, false caso contrário.
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

// View gera a representação visual atual da aplicação.
// Se houver modais abertos, centraliza o modal do topo sobre a área de trabalho.
func (r *RootModel) View() tea.View {
	if r.width == 0 || r.height == 0 {
		return tea.NewView("Aguarde...")
	}

	base := r.activeView.Render(r.width, r.height, r.theme)

	if len(r.modals) == 0 {
		return tea.NewView(base)
	}

	top := r.modals[len(r.modals)-1]
	modalView := top.Render(r.width, r.height, r.theme)

	workH := r.height - 4 // 2 header + 1 msg + 1 action
	modalContent := lipgloss.Place(r.width, workH, lipgloss.Center, lipgloss.Center, modalView)
	v := tea.NewView(modalContent)
	v.AltScreen = true
	v.BackgroundColor = lipgloss.Color(r.theme.Surface.Base)
	return v
}

// Update processa mensagens do Bubble Tea e atualiza o estado do modelo.
// Redireciona eventos para o modal ativo ou para a view principal conforme o contexto.
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

	case tea.KeyMsg:
		// 1. System actions — avaliadas sempre, inclusive com modal ativo.
		if cmd, ok := r.evalActions(msg, r.systemActions); ok {
			return r, cmd
		}

		// 2. Modal ativo recebe a tecla.
		if len(r.modals) > 0 {
			top := len(r.modals) - 1
			return r, r.modals[top].Update(msg)
		}

		// 3. View actions — avaliadas apenas sem modal ativo.
		if r.activeView != nil {
			if cmd, ok := r.evalActions(msg, r.activeView.Actions()); ok {
				return r, cmd
			}
		}

		// 4. Application actions — avaliadas apenas sem modal ativo, após View actions.
		if cmd, ok := r.evalActions(msg, r.applicationActions); ok {
			return r, cmd
		}

		// 5. Tecla não reconhecida — descartada silenciosamente.
		return r, nil
	}

	if len(r.modals) > 0 {
		top := len(r.modals) - 1
		return r, r.modals[top].Update(msg)
	}

	if r.activeView != nil {
		return r, r.activeView.Update(msg)
	}
	return r, nil
}

// Init é chamado uma vez ao iniciar a aplicação. Não há comandos iniciais.
func (r *RootModel) Init() tea.Cmd {
	return nil
}

// RootModelOption é uma função de configuração aplicada ao RootModel na criação.
// Use com NewRootModel para personalizar o modelo sem expor campos internos.
type RootModelOption func(*RootModel)

// WithVersion define a versão da aplicação exibida na interface.
// A versão é normalmente injetada via ldflags no momento do build.
func WithVersion(version string) RootModelOption {
	return func(m *RootModel) {
		m.version = version
	}
}

// NewRootModel cria e inicializa um RootModel com o tema padrão TokyoNight.
// Aplique opções funcionais para personalizar o modelo, ex: WithVersion.
func NewRootModel(opts ...RootModelOption) *RootModel {
	m := &RootModel{
		theme:      design.TokyoNight,
		workArea:   WorkAreaWelcome,
		activeView: nil,
		version:    "dev",
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}
