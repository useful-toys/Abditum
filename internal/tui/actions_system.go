package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// SetupSystemActions registra os grupos e actions de sistema no root.
// System actions são avaliadas em qualquer contexto, inclusive com modal ativo.
// Deve ser chamado após NewRootModel, antes de iniciar o loop do Bubble Tea.
func (r *RootModel) SetupSystemActions() {
	r.RegisterActionGroup(ActionGroup{
		ID:    "system",
		Label: "Sistema",
	})
	r.RegisterSystemActions([]Action{
		{
			Keys:        []design.Key{design.Shortcuts.ThemeToggle},
			Label:       "Tema",
			Description: "Alterna entre os temas Tokyo Night e Cyberpunk.",
			GroupID:     "system",
			Priority:    100,
			Visible:     false, // atalho funciona mas não aparece na barra de status
			OnExecute:   func() tea.Cmd { r.toggleTheme(); return nil },
		},
	})
}
