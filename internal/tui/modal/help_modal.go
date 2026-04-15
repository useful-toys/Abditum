package modal

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/actions"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// HelpModal exibe todas as actions registradas, agrupadas por ActionGroup.
// Implementa tui.ModalView.
type HelpModal struct {
	actions []actions.Action
	groups  []actions.ActionGroup
}

// NewHelpModal cria um HelpModal com as actions e grupos fornecidos.
func NewHelpModal(acts []actions.Action, groups []actions.ActionGroup) *HelpModal {
	return &HelpModal{
		actions: acts,
		groups:  groups,
	}
}

// Render retorna o modal de ajuda — stub minimalista que lista as actions por grupo.
func (m *HelpModal) Render(maxHeight, maxWidth int, theme *design.Theme) string {
	style := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(lipgloss.Color(theme.Border.Default)).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Raised))
	return style.Render("Ajuda")
}

// HandleKey fecha o modal quando Esc é pressionado.
func (m *HelpModal) HandleKey(msg tea.KeyMsg) tea.Cmd {
	if msg.String() == "esc" {
		return tui.CloseModal()
	}
	return nil
}

// Update processa mensagens do Bubble Tea delegando a HandleKey para eventos de teclado.
func (m *HelpModal) Update(msg tea.Msg) tea.Cmd {
	if key, ok := msg.(tea.KeyMsg); ok {
		return m.HandleKey(key)
	}
	return nil
}
