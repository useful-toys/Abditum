package secret

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/vault"
)

// VaultTreeView exibe a árvore de pastas e segredos do cofre aberto.
type VaultTreeView struct{}

// NewVaultTreeView cria uma nova instância da árvore do cofre.
func NewVaultTreeView(vm *vault.Manager) *VaultTreeView {
	return &VaultTreeView{}
}

// Render retorna a árvore de segredos preenchendo as dimensões fornecidas com o tema ativo.
func (v *VaultTreeView) Render(height, width int, theme *design.Theme) string {
	content := "Vault Tree"
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Base))
	return style.Render(content)
}

// HandleKey não processa teclas nesta view.
func (v *VaultTreeView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }

// HandleEvent não processa eventos externos nesta view.
func (v *VaultTreeView) HandleEvent(event any) {}

// HandleTeaMsg não processa mensagens do framework nesta view.
func (v *VaultTreeView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }

// Update não altera o estado desta view em resposta a mensagens.
func (v *VaultTreeView) Update(msg tea.Msg) tea.Cmd { return nil }

// Actions retorna nil — VaultTreeView não possui actions próprias nesta sprint.
func (v *VaultTreeView) Actions() []interface{} { return nil }
