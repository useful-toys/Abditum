package secret

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// DetailView exibe os campos e o valor do segredo selecionado na árvore do cofre.
type DetailView struct{}

// NewDetailView cria uma nova instância do painel de detalhes de segredo.
func NewDetailView() *DetailView {
	return &DetailView{}
}

// Render retorna os detalhes do segredo selecionado pelas dimensões fornecidas com o tema ativo.
func (v *DetailView) Render(height, width int, theme *design.Theme) string {
	content := "Secret Detail"
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Base))
	return style.Render(content)
}

// HandleKey não processa teclas nesta view.
func (v *DetailView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }

// HandleEvent não processa eventos externos nesta view.
func (v *DetailView) HandleEvent(event any) {}

// HandleTeaMsg não processa mensagens do framework nesta view.
func (v *DetailView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }

// Update não altera o estado desta view em resposta a mensagens.
func (v *DetailView) Update(msg tea.Msg) tea.Cmd { return nil }

// Actions retorna nil — DetailView não possui actions próprias nesta sprint.
func (v *DetailView) Actions() []tui.Action { return nil }
