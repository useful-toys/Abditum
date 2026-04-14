package screen

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// WelcomeView é exibida quando nenhum cofre está aberto.
// Apresenta uma tela de boas-vindas com orientações para o usuário começar.
type WelcomeView struct{}

// NewWelcomeView cria uma nova instância da tela de boas-vindas.
func NewWelcomeView() *WelcomeView {
	return &WelcomeView{}
}

// Render retorna a tela de boas-vindas preenchendo as dimensões fornecidas com o tema ativo.
func (v *WelcomeView) Render(height, width int, theme *design.Theme) string {
	content := "Welcome"
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Base))
	return style.Render(content)
}

// HandleKey não processa nenhuma tecla nesta view.
func (v *WelcomeView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }

// HandleEvent não processa eventos externos nesta view.
func (v *WelcomeView) HandleEvent(event any) {}

// HandleTeaMsg não processa mensagens do framework nesta view.
func (v *WelcomeView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }

// Update não altera o estado desta view em resposta a mensagens.
func (v *WelcomeView) Update(msg tea.Msg) tea.Cmd { return nil }

// Actions retorna nil — WelcomeView não possui actions próprias nesta sprint.
func (v *WelcomeView) Actions() []tui.Action { return nil }
