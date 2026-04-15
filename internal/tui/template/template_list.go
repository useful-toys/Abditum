package template

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/actions"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/vault"
)

// TemplateListView exibe a lista de templates disponíveis para criação de segredos.
type TemplateListView struct{}

// NewTemplateListView cria uma nova instância da lista de templates.
func NewTemplateListView(vm *vault.Manager) *TemplateListView {
	return &TemplateListView{}
}

// Render retorna a lista de templates preenchendo as dimensões fornecidas com o tema ativo.
func (v *TemplateListView) Render(height, width int, theme *design.Theme) string {
	content := "Template List"
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Base))
	return style.Render(content)
}

// HandleKey não processa teclas nesta view.
func (v *TemplateListView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }

// HandleEvent não processa eventos externos nesta view.
func (v *TemplateListView) HandleEvent(event any) {}

// HandleTeaMsg não processa mensagens do framework nesta view.
func (v *TemplateListView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }

// Update não altera o estado desta view em resposta a mensagens.
func (v *TemplateListView) Update(msg tea.Msg) tea.Cmd { return nil }

// Actions retorna nil — TemplateListView não possui actions próprias nesta sprint.
func (v *TemplateListView) Actions() []actions.Action { return nil }
