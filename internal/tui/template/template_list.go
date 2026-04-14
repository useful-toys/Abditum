package template

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// ListView exibe a lista de templates disponíveis para criação de segredos.
type ListView struct{}

// NewListView cria uma nova instância da lista de templates.
func NewListView() *ListView {
	return &ListView{}
}

// Render retorna a lista de templates preenchendo as dimensões fornecidas com o tema ativo.
func (v *ListView) Render(height, width int, theme *design.Theme) string {
	content := "Template List"
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Base))
	return style.Render(content)
}

// HandleKey não processa teclas nesta view.
func (v *ListView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }

// HandleEvent não processa eventos externos nesta view.
func (v *ListView) HandleEvent(event any) {}

// HandleTeaMsg não processa mensagens do framework nesta view.
func (v *ListView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }

// Update não altera o estado desta view em resposta a mensagens.
func (v *ListView) Update(msg tea.Msg) tea.Cmd { return nil }

// Actions retorna nil — ListView não possui actions próprias nesta sprint.
func (v *ListView) Actions() []tui.Action { return nil }
