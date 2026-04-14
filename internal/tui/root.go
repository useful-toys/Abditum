package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
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
}

// View gera a representação visual atual da aplicação.
// Se houver modais abertos, centraliza o modal do topo sobre a área de trabalho.
func (r *RootModel) View() tea.View {
	if r.width == 0 || r.height == 0 {
		return tea.NewView("Aguarde...")
	}

	base := r.activeView.Render(r.width, r.height, *r.theme)

	if len(r.modals) == 0 {
		return tea.NewView(base)
	}

	top := r.modals[len(r.modals)-1]
	modalView := top.Render(r.width, r.height, *r.theme)

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
	}

	if len(r.modals) > 0 {
		top := len(r.modals) - 1
		return r, r.modals[top].Update(msg)
	}

	return r, r.activeView.Update(msg)
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
