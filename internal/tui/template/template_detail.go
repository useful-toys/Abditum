package template

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/vault"
)

// TemplateDetailView exibe os detalhes e campos de um template de segredo selecionado.
type TemplateDetailView struct{}

// NewTemplateDetailView cria uma nova instância do painel de detalhes de template.
func NewTemplateDetailView(vm *vault.Manager) *TemplateDetailView {
	return &TemplateDetailView{}
}

// Render retorna os detalhes do template selecionado pelas dimensões fornecidas com o tema ativo.
func (v *TemplateDetailView) Render(height, width int, theme *design.Theme) string {
	content := "Template Detail"
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Base))
	return style.Render(content)
}

// HandleKey não processa teclas nesta view.
func (v *TemplateDetailView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }

// HandleEvent não processa eventos externos nesta view.
func (v *TemplateDetailView) HandleEvent(event any) {}

// HandleTeaMsg não processa mensagens do framework nesta view.
func (v *TemplateDetailView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }

// Update não altera o estado desta view em resposta a mensagens.
func (v *TemplateDetailView) Update(msg tea.Msg) tea.Cmd { return nil }

// Actions retorna nil — TemplateDetailView não possui actions próprias nesta sprint.
func (v *TemplateDetailView) Actions() []tui.Action { return nil }
