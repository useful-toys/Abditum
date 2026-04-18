package main

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/actions"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/tui/operation"
)

// setupActions registers all action groups and actions in the application (system and application).
// Must be called after NewRootModel, before starting the Bubble Tea loop.
func setupActions(r *tui.RootModel) {
	setupSystem(r)
	setupApplication(r)
}

// setupSystem registers system-level action groups and actions in root.
// System actions are evaluated in any context, including with active modal.
func setupSystem(r *tui.RootModel) {
	r.RegisterActionGroup(actions.ActionGroup{
		ID:    "system",
		Label: "Sistema",
	})
	r.RegisterSystemActions([]actions.Action{
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

// setupApplication registers application-level action groups and actions in root.
// Application actions are evaluated only when no modal is active.
func setupApplication(r *tui.RootModel) {
	r.RegisterActionGroup(actions.ActionGroup{
		ID:    "app",
		Label: "Aplicação",
	})
	r.RegisterApplicationActions([]actions.Action{
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
			Description: "Encerra a aplicação, com opção de salvar se houver alterações.",
			GroupID:     "app",
			Priority:    20,
			Visible:     true,
			OnExecute: func() tea.Cmd {
				return tui.StartOperation(operation.NewQuitOperation(r.MessageController(), r.Manager()))
			},
		},
		{
			Keys:        []design.Key{design.Shortcuts.NewVault},
			Label:       "Criar cofre",
			Description: "Cria um novo cofre protegido por senha.",
			GroupID:     "app",
			Priority:    30,
			Visible:     true,
			AvailableWhen: func(app actions.AppState, _ actions.ChildView) bool {
				// Criar cofre só faz sentido quando nenhum cofre está aberto.
				return app.Manager() == nil
			},
			OnExecute: func() tea.Cmd {
				return tui.StartOperation(operation.NewCriarCofreOperation(r.MessageController(), r.Manager(), ""))
			},
		},
		{
			Keys:        []design.Key{design.Shortcuts.OpenVault},
			Label:       "Abrir cofre",
			Description: "Abre um cofre existente a partir de um arquivo.",
			GroupID:     "app",
			Priority:    31,
			Visible:     true,
			AvailableWhen: func(app actions.AppState, _ actions.ChildView) bool {
				// Abrir cofre só faz sentido quando nenhum cofre está aberto.
				return app.Manager() == nil
			},
			OnExecute: func() tea.Cmd {
				return tui.StartOperation(operation.NewAbrirCofreOperation(r.MessageController(), r.Manager(), ""))
			},
		},
		{
			Keys:        []design.Key{design.Keys.F2},
			Label:       "Operação Fake",
			Description: "Demonstração do padrão Operation — confirmação + trabalho assíncrono.",
			GroupID:     "app",
			Priority:    99,
			Visible:     true,
			OnExecute: func() tea.Cmd {
				return tui.StartOperation(operation.NewFakeOperation(r.MessageController()))
			},
		},
	})
}
