package settings

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/vault"
)

// SettingsView exibe as opções de configuração da aplicação.
type SettingsView struct{}

// NewSettingsView cria uma nova instância da tela de configurações.
// vm é o gerenciador do cofre ativo — pode ser nil durante inicialização.
func NewSettingsView(vm *vault.Manager) *SettingsView {
	return &SettingsView{}
}

// Render retorna a tela de configurações preenchendo as dimensões fornecidas com o tema ativo.
func (v *SettingsView) Render(height, width int, theme *design.Theme) string {
	content := "Settings"
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Base))
	return style.Render(content)
}

// HandleKey não processa teclas nesta view.
func (v *SettingsView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }

// HandleEvent não processa eventos externos nesta view.
func (v *SettingsView) HandleEvent(event any) {}

// HandleTeaMsg não processa mensagens do framework nesta view.
func (v *SettingsView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }

// Update não altera o estado desta view em resposta a mensagens.
func (v *SettingsView) Update(msg tea.Msg) tea.Cmd { return nil }

// Actions retorna nil — SettingsView não possui actions próprias nesta sprint.
func (v *SettingsView) Actions() []tui.Action { return nil }
