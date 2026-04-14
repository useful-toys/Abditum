package secret

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// TreeView exibe a árvore de pastas e segredos do cofre aberto.
// É o painel de navegação lateral da área de trabalho do cofre.
type TreeView struct{}

// NewTreeView cria uma nova instância da árvore do cofre.
func NewTreeView() *TreeView {
	return &TreeView{}
}

// Render retorna a árvore de segredos preenchendo as dimensões fornecidas com o tema ativo.
func (v *TreeView) Render(height, width int, theme design.Theme) string {
	content := "Vault Tree"
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Base))
	return style.Render(content)
}

// HandleKey não processa teclas nesta view.
func (v *TreeView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }

// HandleEvent não processa eventos externos nesta view.
func (v *TreeView) HandleEvent(event any) {}

// HandleTeaMsg não processa mensagens do framework nesta view.
func (v *TreeView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }

// Update não altera o estado desta view em resposta a mensagens.
func (v *TreeView) Update(msg tea.Msg) tea.Cmd { return nil }
