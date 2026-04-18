package modal_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/actions"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/tui/testdata"
)

func sampleActionsAndGroups() ([]actions.Action, []actions.ActionGroup) {
	groups := []actions.ActionGroup{
		{ID: "app", Label: "Aplicação", Order: 0},
		{ID: "nav", Label: "Navegação", Order: 1},
	}
	acts := []actions.Action{
		{Keys: []design.Key{design.Shortcuts.Help}, Label: "Ajuda", Description: "Abre o diálogo de ajuda.", GroupID: "app", Priority: 10},
		{Keys: []design.Key{design.Shortcuts.Quit}, Label: "Sair", Description: "Encerra a aplicação.", GroupID: "app", Priority: 20},
		{Keys: []design.Key{design.Keys.Enter}, Label: "Abrir", Description: "Abre o item selecionado.", GroupID: "nav", Priority: 10},
		{Keys: []design.Key{design.Keys.Esc}, Label: "Voltar", Description: "Retorna ao nível anterior.", GroupID: "nav", Priority: 20},
	}
	return acts, groups
}

func TestHelpModal_NoScroll(t *testing.T) {
	acts, groups := sampleActionsAndGroups()
	m := modal.NewHelpModal(acts, groups)
	testdata.TestRenderManaged(t, "help_modal", "no_scroll", []string{"60x20"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestHelpModal_WithScroll(t *testing.T) {
	acts, groups := sampleActionsAndGroups()
	m := modal.NewHelpModal(acts, groups)
	// Small height forces scroll
	testdata.TestRenderManaged(t, "help_modal", "with_scroll", []string{"60x6"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestHelpModal_HandleKey_Esc_ClosesModal(t *testing.T) {
	acts, groups := sampleActionsAndGroups()
	m := modal.NewHelpModal(acts, groups)
	cmd := m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEscape})
	if cmd == nil {
		t.Error("HandleKey(Esc): cmd = nil, want CloseModal command")
	}
}

func TestHelpModal_HandleKey_WithUnmatchedKey_ReturnsNil(t *testing.T) {
	acts, groups := sampleActionsAndGroups()
	m := modal.NewHelpModal(acts, groups)
	// HandleKey com tecla que não fecha o modal deve retornar nil
	cmd := m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	if cmd != nil {
		t.Error("HandleKey(unmatched key): expected nil cmd")
	}
}
