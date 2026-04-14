package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// helpModalFactory is a function variable that can be set to create HelpModals.
// This exists to break the circular import cycle: tui imports tui/modal, which imports tui.
// The factory is set by the modal package during init() or by the caller before use.
var helpModalFactory func([]Action, []ActionGroup) ModalView

// SetupApplicationActions registra os grupos e actions de aplicação no root.
// Application actions são avaliadas apenas quando nenhum modal está ativo.
// Deve ser chamado após NewRootModel, antes de iniciar o loop do Bubble Tea.
func (r *RootModel) SetupApplicationActions() {
	r.RegisterActionGroup(ActionGroup{
		ID:    "app",
		Label: "Aplicação",
	})
	r.RegisterApplicationActions([]Action{
		{
			Keys:        []design.Key{design.Shortcuts.Help},
			Label:       "Ajuda",
			Description: "Abre o diálogo de ajuda com todos os atalhos disponíveis.",
			GroupID:     "app",
			Priority:    10,
			Visible:     true,
			OnExecute: func() tea.Cmd {
				if helpModalFactory == nil {
					return nil // Factory not initialized
				}
				var viewActions []Action
				if r.activeView != nil {
					viewActions = r.activeView.Actions()
				}
				// Use slices.Concat-like logic inline to avoid import
				allActions := make([]Action, 0, len(r.systemActions)+len(r.applicationActions)+len(viewActions))
				allActions = append(allActions, r.systemActions...)
				allActions = append(allActions, r.applicationActions...)
				allActions = append(allActions, viewActions...)
				return OpenModal(helpModalFactory(allActions, r.actionGroups))
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
