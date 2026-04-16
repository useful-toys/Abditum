package modal_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/tui/testdata"
)

func TestConfirmModal_Destructive(t *testing.T) {
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "Excluir",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd { return tui.CloseModal() },
		},
		{
			Keys:   []design.Key{design.Keys.Esc},
			Label:  "Cancelar",
			Intent: modal.IntentCancel,
			Action: func() tea.Cmd { return tui.CloseModal() },
		},
	}
	m := modal.NewConfirmModalSeverity(design.SeverityDestructive,
		"Excluir cofre",
		"Esta ação é permanente e não pode ser desfeita.",
		opts,
	)
	testdata.TestRenderManaged(t, "confirm_modal", "destructive", []string{"60x10"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestConfirmModal_HandleKey_Enter_ExecutesAction(t *testing.T) {
	called := false
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "OK",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd {
				called = true
				return nil
			},
		},
	}
	m := modal.NewConfirmModal("Título", "Mensagem", opts)
	_ = m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEnter})
	if !called {
		t.Error("HandleKey(Enter): action not called")
	}
}

func TestConfirmModal_Update_DelegatesKeys(t *testing.T) {
	called := false
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Esc},
			Label:  "Cancelar",
			Intent: modal.IntentCancel,
			Action: func() tea.Cmd {
				called = true
				return nil
			},
		},
	}
	m := modal.NewConfirmModal("Título", "Mensagem", opts)
	_ = m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if !called {
		t.Error("Update(KeyEsc): action not called — Update must delegate to HandleKey")
	}
}
