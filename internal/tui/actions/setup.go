package actions

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
)

// Setup registra todos os grupos e actions na aplicação (system e application).
// Deve ser chamado após NewRootModel, antes de iniciar o loop do Bubble Tea.
func Setup(r *tui.RootModel) {
	SetupSystem(r)
	SetupApplication(r)
}

// SetupSystem registra os grupos e actions de sistema no root.
// System actions são avaliadas em qualquer contexto, inclusive com modal ativo.
// Deve ser chamado após NewRootModel, antes de iniciar o loop do Bubble Tea.
func SetupSystem(r *tui.RootModel) {
	r.RegisterActionGroup(tui.ActionGroup{
		ID:    "system",
		Label: "Sistema",
	})
	r.RegisterSystemActions([]tui.Action{
		{
			Keys:        []design.Key{design.Shortcuts.ThemeToggle},
			Label:       "Tema",
			Description: "Alterna entre os temas Tokyo Night e Cyberpunk.",
			GroupID:     "system",
			Priority:    100,
			Visible:     false, // atalho funciona mas não aparece na barra de status
			OnExecute:   func() tea.Cmd { r.ToggleTheme(); return nil },
		},
	})
}

// SetupApplication registra os grupos e actions de aplicação no root.
// Application actions são avaliadas apenas quando nenhum modal está ativo.
// Deve ser chamado após NewRootModel, antes de iniciar o loop do Bubble Tea.
func SetupApplication(r *tui.RootModel) {
	r.RegisterActionGroup(tui.ActionGroup{
		ID:    "app",
		Label: "Aplicação",
	})
	r.RegisterApplicationActions([]tui.Action{
		{
			Keys:        []design.Key{design.Shortcuts.Help},
			Label:       "Ajuda",
			Description: "Abre o diálogo de ajuda com todos os atalhos disponíveis.",
			GroupID:     "app",
			Priority:    10,
			Visible:     true,
			OnExecute: func() tea.Cmd {
				return tui.OpenModal(modal.NewHelpModal(r.ActiveViewActions(), r.GetActionGroups()))
			},
		},
		{
			Keys:        []design.Key{design.Shortcuts.Quit},
			Label:       "Sair",
			Description: "Encerra a aplicação.",
			GroupID:     "app",
			Priority:    20,
			Visible:     true,
			OnExecute:   func() tea.Cmd { return tea.Quit },
		},
	})
}
