package modal

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

func makeSpecialKeyMsg(code rune) tea.KeyMsg {
	return tea.KeyPressMsg{Code: code}
}

func TestKeyHandler_OptionMatch_ExecutesAction(t *testing.T) {
	called := false
	opts := []ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "OK",
			Action: func() tea.Cmd {
				called = true
				return nil
			},
		},
	}
	h := KeyHandler{Options: opts}
	cmd, handled := h.Handle(makeSpecialKeyMsg(tea.KeyEnter))
	if !handled {
		t.Error("Handle(Enter): handled = false, want true")
	}
	_ = cmd
	if !called {
		t.Error("Handle(Enter): action was not called")
	}
}

func TestKeyHandler_MultipleKeys_AnyActivatesAction(t *testing.T) {
	callCount := 0
	opts := []ModalOption{
		{
			Keys:  []design.Key{design.Keys.Enter, design.Keys.Esc},
			Label: "OK",
			Action: func() tea.Cmd {
				callCount++
				return nil
			},
		},
	}
	h := KeyHandler{Options: opts}

	h.Handle(makeSpecialKeyMsg(tea.KeyEnter))
	h.Handle(makeSpecialKeyMsg(tea.KeyEscape))
	if callCount != 2 {
		t.Errorf("Handle: callCount = %d, want 2 (both keys should trigger)", callCount)
	}
}

func TestKeyHandler_UnrecognizedKey_ReturnsNotHandled(t *testing.T) {
	// Tecla completamente fora do mapeamento (nem Enter/Esc nem explícitas).
	h := KeyHandler{Options: []ModalOption{
		{
			Keys:   []design.Key{design.Keys.Tab},
			Label:  "OK",
			Action: func() tea.Cmd { return nil },
		},
	}}
	// 'a' não é Enter, Esc nem Tab — não deve ser tratada.
	_, handled := h.Handle(tea.KeyPressMsg{Code: 'a'})
	if handled {
		t.Error("Handle('a' quando apenas Tab registrado): handled = true, want false")
	}
}

func TestKeyHandler_ScrollKeys_WithScroll(t *testing.T) {
	scroll := &ScrollState{Offset: 5, Total: 30, Viewport: 10}
	h := KeyHandler{Scroll: scroll}

	_, handled := h.Handle(makeSpecialKeyMsg(tea.KeyUp))
	if !handled {
		t.Error("Handle(Up) with Scroll: handled = false, want true")
	}
	if scroll.Offset != 4 {
		t.Errorf("Handle(Up): Offset = %d, want 4", scroll.Offset)
	}

	_, handled = h.Handle(makeSpecialKeyMsg(tea.KeyDown))
	if !handled {
		t.Error("Handle(Down) with Scroll: handled = false, want true")
	}
	if scroll.Offset != 5 {
		t.Errorf("Handle(Down): Offset = %d, want 5", scroll.Offset)
	}

	scroll.Offset = 5
	_, handled = h.Handle(makeSpecialKeyMsg(tea.KeyPgUp))
	if !handled {
		t.Error("Handle(PgUp): handled = false, want true")
	}
	if scroll.Offset != 0 {
		t.Errorf("Handle(PgUp): Offset = %d, want 0 (5-viewport=5-10→0)", scroll.Offset)
	}

	_, handled = h.Handle(makeSpecialKeyMsg(tea.KeyPgDown))
	if !handled {
		t.Error("Handle(PgDown): handled = false, want true")
	}
	if scroll.Offset != 10 {
		t.Errorf("Handle(PgDown): Offset = %d, want 10", scroll.Offset)
	}

	_, handled = h.Handle(makeSpecialKeyMsg(tea.KeyHome))
	if !handled {
		t.Error("Handle(Home): handled = false, want true")
	}
	if scroll.Offset != 0 {
		t.Errorf("Handle(Home): Offset = %d, want 0", scroll.Offset)
	}

	_, handled = h.Handle(makeSpecialKeyMsg(tea.KeyEnd))
	if !handled {
		t.Error("Handle(End): handled = false, want true")
	}
	if scroll.Offset != 20 {
		t.Errorf("Handle(End): Offset = %d, want 20 (total-viewport=30-10)", scroll.Offset)
	}
}

func TestKeyHandler_ScrollKeys_WithoutScroll_NotHandled(t *testing.T) {
	h := KeyHandler{} // Scroll == nil
	_, handled := h.Handle(makeSpecialKeyMsg(tea.KeyUp))
	if handled {
		t.Error("Handle(Up) with nil Scroll: handled = true, want false")
	}
	_, handled = h.Handle(makeSpecialKeyMsg(tea.KeyDown))
	if handled {
		t.Error("Handle(Down) with nil Scroll: handled = true, want false")
	}
}

func TestKeyHandler_EmptyOptions_ScrollStillWorks(t *testing.T) {
	scroll := &ScrollState{Offset: 3, Total: 20, Viewport: 5}
	h := KeyHandler{Scroll: scroll}
	_, handled := h.Handle(makeSpecialKeyMsg(tea.KeyUp))
	if !handled {
		t.Error("Handle(Up) with empty Options but Scroll != nil: not handled")
	}
	if scroll.Offset != 2 {
		t.Errorf("Offset = %d, want 2", scroll.Offset)
	}
}

func TestKeyHandler_ImplicitEnter_FirstOption_NoKeys(t *testing.T) {
	// Primeira option sem Keys: Enter deve disparar sua ação.
	called := false
	opts := []ModalOption{
		{
			Label:  "Confirmar",
			Action: func() tea.Cmd { called = true; return nil },
		},
		{
			Keys:   []design.Key{design.Keys.Esc},
			Label:  "Cancelar",
			Action: func() tea.Cmd { return nil },
		},
	}
	h := KeyHandler{Options: opts}
	_, handled := h.Handle(makeSpecialKeyMsg(tea.KeyEnter))
	if !handled {
		t.Error("Handle(Enter): handled = false, want true")
	}
	if !called {
		t.Error("Handle(Enter): action da primeira option não foi chamada")
	}
}

func TestKeyHandler_ImplicitEsc_LastOption_NoKeys(t *testing.T) {
	// Última option sem Keys: Esc deve disparar sua ação.
	called := false
	opts := []ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "Confirmar",
			Action: func() tea.Cmd { return nil },
		},
		{
			Label:  "Cancelar",
			Action: func() tea.Cmd { called = true; return nil },
		},
	}
	h := KeyHandler{Options: opts}
	_, handled := h.Handle(makeSpecialKeyMsg(tea.KeyEscape))
	if !handled {
		t.Error("Handle(Esc): handled = false, want true")
	}
	if !called {
		t.Error("Handle(Esc): action da última option não foi chamada")
	}
}

func TestKeyHandler_ImplicitBoth_SingleOption_NoKeys(t *testing.T) {
	// Option única sem Keys: Enter e Esc devem disparar a mesma ação.
	callCount := 0
	opts := []ModalOption{
		{
			Label:  "OK",
			Action: func() tea.Cmd { callCount++; return nil },
		},
	}
	h := KeyHandler{Options: opts}

	h.Handle(makeSpecialKeyMsg(tea.KeyEnter))
	h.Handle(makeSpecialKeyMsg(tea.KeyEscape))
	if callCount != 2 {
		t.Errorf("Single option sem Keys: callCount = %d, want 2 (Enter e Esc ambos devem disparar)", callCount)
	}
}

func TestKeyHandler_ImplicitEnter_AddsToExplicitKeys(t *testing.T) {
	// Primeira option com Keys: [letter('s')].
	// Enter deve ser adicionado como alias — action chamada por 's' e por Enter.
	callCount := 0
	opts := []ModalOption{
		{
			Keys:   []design.Key{design.Letter('s')},
			Label:  "Sim",
			Action: func() tea.Cmd { callCount++; return nil },
		},
		{
			Keys:   []design.Key{design.Letter('n')},
			Label:  "Não",
			Action: func() tea.Cmd { return nil },
		},
	}
	h := KeyHandler{Options: opts}

	h.Handle(tea.KeyPressMsg{Code: 's'})      // tecla explícita
	h.Handle(makeSpecialKeyMsg(tea.KeyEnter)) // alias implícito
	if callCount != 2 {
		t.Errorf("First option with explicit key + implicit Enter: callCount = %d, want 2", callCount)
	}
}

func TestKeyHandler_ImplicitKeys_DoNotOverrideExplicit(t *testing.T) {
	// Se Enter já está declarado explicitamente, não deve ser disparado duas vezes.
	callCount := 0
	opts := []ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "OK",
			Action: func() tea.Cmd { callCount++; return nil },
		},
	}
	h := KeyHandler{Options: opts}

	h.Handle(makeSpecialKeyMsg(tea.KeyEnter))
	if callCount != 1 {
		t.Errorf("Enter explícito + implícito: callCount = %d, want 1 (não deve duplicar)", callCount)
	}
}
