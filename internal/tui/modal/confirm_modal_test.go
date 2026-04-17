package modal_test

import (
	"strings"
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

// Golden file tests for different option counts

func TestConfirmModal_SingleOption(t *testing.T) {
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "OK",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd { return tui.CloseModal() },
		},
	}
	m := modal.NewConfirmModal("Confirmação", "Tem certeza?", opts)
	testdata.TestRenderManaged(t, "confirm_modal", "single_option", []string{"50x8"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestConfirmModal_ThreeOptions(t *testing.T) {
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "Salvar",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd { return tui.CloseModal() },
		},
		{
			Keys:   []design.Key{design.Keys.Tab},
			Label:  "Descartar",
			Intent: modal.IntentOther,
			Action: func() tea.Cmd { return tui.CloseModal() },
		},
		{
			Keys:   []design.Key{design.Keys.Esc},
			Label:  "Cancelar",
			Intent: modal.IntentCancel,
			Action: func() tea.Cmd { return tui.CloseModal() },
		},
	}
	m := modal.NewConfirmModal("Alterações pendentes", "Deseja salvar as alterações?", opts)
	testdata.TestRenderManaged(t, "confirm_modal", "three_options", []string{"70x10"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

// Tests with custom keys (S Sim / N Não)

func TestConfirmModal_CustomKeysSimNao(t *testing.T) {
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Letter('s')},
			Label:  "Sim",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd { return tui.CloseModal() },
		},
		{
			Keys:   []design.Key{design.Letter('n')},
			Label:  "Não",
			Intent: modal.IntentCancel,
			Action: func() tea.Cmd { return tui.CloseModal() },
		},
	}
	m := modal.NewConfirmModalSeverity(design.SeverityAlert,
		"Limpar histórico",
		"Tem certeza que deseja limpar todo o histórico?",
		opts,
	)
	testdata.TestRenderManaged(t, "confirm_modal", "custom_keys_sim_nao", []string{"60x10"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

// Tests with text length variations

func TestConfirmModal_ShortText(t *testing.T) {
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "OK",
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
	m := modal.NewConfirmModal("Aviso", "Continuar?", opts)
	testdata.TestRenderManaged(t, "confirm_modal", "short_text", []string{"50x8"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestConfirmModal_LongText(t *testing.T) {
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "Confirmar",
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
	m := modal.NewConfirmModal(
		"Exclusão permanente de dados",
		"Esta ação irá remover permanentemente todos os dados associados à sua conta. Esta operação não pode ser desfeita e resultará em perda total de todos os seus arquivos, configurações e histórico.",
		opts,
	)
	testdata.TestRenderManaged(t, "confirm_modal", "long_text", []string{"80x12"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

// Tests with multiline text

func TestConfirmModal_TwoLineText(t *testing.T) {
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "Continuar",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd { return tui.CloseModal() },
		},
		{
			Keys:   []design.Key{design.Keys.Esc},
			Label:  "Voltar",
			Intent: modal.IntentCancel,
			Action: func() tea.Cmd { return tui.CloseModal() },
		},
	}
	m := modal.NewConfirmModal(
		"Operação",
		"Primeira linha de texto\nSegunda linha de texto",
		opts,
	)
	testdata.TestRenderManaged(t, "confirm_modal", "two_line_text", []string{"60x10"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestConfirmModal_ThreeLineText(t *testing.T) {
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "Prosseguir",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd { return tui.CloseModal() },
		},
		{
			Keys:   []design.Key{design.Keys.Esc},
			Label:  "Sair",
			Intent: modal.IntentCancel,
			Action: func() tea.Cmd { return tui.CloseModal() },
		},
	}
	m := modal.NewConfirmModal(
		"Confirmação",
		"Primeira linha de informação\nSegunda linha com mais detalhes\nTerceira linha final",
		opts,
	)
	testdata.TestRenderManaged(t, "confirm_modal", "three_line_text", []string{"70x12"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

// Tests with title variations

func TestConfirmModal_ShortTitle(t *testing.T) {
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "OK",
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
	m := modal.NewConfirmModal("OK?", "Prosseguir?", opts)
	testdata.TestRenderManaged(t, "confirm_modal", "short_title", []string{"45x8"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestConfirmModal_LongTitle(t *testing.T) {
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "Confirmar",
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
	m := modal.NewConfirmModal(
		"Este é um título muito longo que descreve a operação em detalhes",
		"Deseja confirmar esta ação?",
		opts,
	)
	testdata.TestRenderManaged(t, "confirm_modal", "long_title", []string{"80x10"},
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

func TestConfirmModal_SingleOption_EscDispatchesAction(t *testing.T) {
	// Quando há apenas 1 opção, ESC deve automaticamente disparar a mesma ação.
	// Essa é uma regra de UX: com uma única ação, tanto ENTER quanto ESC devem executá-la.
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

	// ESC deve disparar a ação mesmo que não tenha sido originalmente registrada
	_ = m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEscape})
	if !called {
		t.Error("HandleKey(Esc) with single option: action not called — Esc should trigger the only action")
	}
}

func TestConfirmModal_SingleOption_EnterStillWorks(t *testing.T) {
	// Mesmo com apenas 1 opção, ENTER deve continuar funcionando como antes.
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
		t.Error("HandleKey(Enter) with single option: action not called")
	}
}

func TestConfirmModal_MultipleOptions_EscNotAutoAdded(t *testing.T) {
	// Com múltiplas opções, ESC NÃO deve ser automaticamente adicionado se não estiver registrada.
	// ESC só funciona se explicitamente mapeada a uma opção.
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "Confirmar",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd { return nil },
		},
		{
			Keys:   []design.Key{design.Keys.Tab},
			Label:  "Outra",
			Intent: modal.IntentOther,
			Action: func() tea.Cmd { return nil },
		},
	}
	m := modal.NewConfirmModal("Título", "Mensagem", opts)

	// ESC não foi registrada explicitamente, então não deve funcionar.
	cmd := m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEscape})
	if cmd != nil {
		t.Error("HandleKey(Esc) with multiple options: should return nil — Esc not explicitly registered")
	}
}

func TestConfirmModal_SingleOption_EscAlreadyRegistered(t *testing.T) {
	// Se ESC já está registrada, ela não deve ser adicionada novamente (evitar duplicata).
	callCount := 0
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter, design.Keys.Esc},
			Label:  "OK",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd {
				callCount++
				return nil
			},
		},
	}
	m := modal.NewConfirmModal("Título", "Mensagem", opts)

	// ESC deve funcionar apenas uma vez, não duplicada.
	_ = m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEscape})
	if callCount != 1 {
		t.Errorf("HandleKey(Esc) with pre-registered Esc: callCount = %d, want 1 (should not duplicate)", callCount)
	}
}

// --- Testes de Integração ---

func TestConfirmModal_Update_WithKeyMsg_DelegatesCorrectly(t *testing.T) {
	// Update deve processar mensagens de teclado através de HandleKey
	called := false
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "OK",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd {
				called = true
				return nil // Action pode retornar nil ou um comando
			},
		},
	}
	m := modal.NewConfirmModal("Teste", "Mensagem", opts)

	// Update deve aceitar tea.KeyMsg e repassar para HandleKey
	cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if !called {
		t.Error("Update(KeyMsg): action not executed")
	}
	// cmd pode ser nil se a action retornar nil, ou um comando válido
	_ = cmd
}

func TestConfirmModal_Update_WithNonKeyMsg_ReturnsNil(t *testing.T) {
	// Update deve ignorar mensagens que não são KeyMsg
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "OK",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd { return nil },
		},
	}
	m := modal.NewConfirmModal("Teste", "Mensagem", opts)

	// Simular uma mensagem de outro tipo (ex: tea.WindowSizeMsg)
	cmd := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	if cmd != nil {
		t.Error("Update(non-KeyMsg): should return nil")
	}
}

func TestConfirmModal_SequentialKeyPresses_AllProcessed(t *testing.T) {
	// Teste que simula múltiplas interações sequenciais
	callCount := 0
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "OK",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd {
				callCount++
				return nil
			},
		},
	}
	m := modal.NewConfirmModal("Teste", "Mensagem", opts)

	// Simular múltiplas pressões de tecla
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEnter})
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEnter})
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEscape})

	if callCount != 3 {
		t.Errorf("Sequential keypresses: callCount = %d, want 3", callCount)
	}
}

func TestConfirmModal_DifferentSeverities_AllWorkCorrectly(t *testing.T) {
	// Teste que verifica se a adição automática de ESC funciona em todas as severidades
	severities := []design.Severity{
		design.SeverityNeutral,
		design.SeverityInformative,
		design.SeverityAlert,
		design.SeverityDestructive,
		design.SeverityError,
	}

	for _, sev := range severities {
		called := false
		opts := []modal.ModalOption{
			{
				Keys:   []design.Key{design.Keys.Enter},
				Label:  "Ação",
				Intent: modal.IntentConfirm,
				Action: func() tea.Cmd {
					called = true
					return nil
				},
			},
		}
		m := modal.NewConfirmModalSeverity(sev, "Título", "Mensagem", opts)

		// ESC deve funcionar em qualquer severidade
		_ = m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEscape})
		if !called {
			t.Errorf("Severity(%d): ESC should trigger action", sev)
		}
	}
}

func TestConfirmModal_MultipleOptions_BothEscAndEnter_DifferentActions(t *testing.T) {
	// Teste que verifica comportamento correto com múltiplas opções onde ambas ESC e ENTER estão mapeadas
	confirmCalled := false
	cancelCalled := false

	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "Confirmar",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd {
				confirmCalled = true
				return nil
			},
		},
		{
			Keys:   []design.Key{design.Keys.Esc},
			Label:  "Cancelar",
			Intent: modal.IntentCancel,
			Action: func() tea.Cmd {
				cancelCalled = true
				return nil
			},
		},
	}
	m := modal.NewConfirmModal("Teste", "Mensagem", opts)

	// ENTER deve executar a primeira ação
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEnter})
	if !confirmCalled {
		t.Error("Multiple options: ENTER should trigger confirm action")
	}

	// ESC deve executar a segunda ação
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEscape})
	if !cancelCalled {
		t.Error("Multiple options: ESC should trigger cancel action")
	}
}

func TestConfirmModal_UnhandledKey_ReturnsNil(t *testing.T) {
	// Teste que verifica se teclas não mapeadas retornam nil
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "OK",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd { return nil },
		},
	}
	m := modal.NewConfirmModal("Teste", "Mensagem", opts)

	// Tecla 'a' não está mapeada
	cmd := m.HandleKey(tea.KeyPressMsg{Code: 'a'})
	if cmd != nil {
		t.Error("Unhandled key: should return nil")
	}
}

func TestConfirmModal_Cursor_AlwaysNil(t *testing.T) {
	// ConfirmModal não tem campo de entrada, logo Cursor deve sempre retornar nil
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "OK",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd { return nil },
		},
	}
	m := modal.NewConfirmModal("Teste", "Mensagem", opts)

	cursor := m.Cursor(80, 24)
	if cursor != nil {
		t.Error("Cursor: should always return nil")
	}
}

func TestConfirmModal_RenderWithDifferentDimensions(t *testing.T) {
	// Teste de renderização com diferentes dimensões
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "OK",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd { return nil },
		},
	}
	theme := design.TokyoNight

	dimensions := []struct {
		width  int
		height int
	}{
		{40, 8},
		{60, 10},
		{80, 15},
		{100, 20},
	}

	for _, dim := range dimensions {
		m := modal.NewConfirmModal("Título", "Mensagem de teste", opts)
		rendered := m.Render(dim.height, dim.width, theme)

		if rendered == "" {
			t.Errorf("Render(%d, %d): empty output", dim.width, dim.height)
		}
		// Verificar se o conteúdo está presente
		if !strings.Contains(rendered, "Título") && !strings.Contains(rendered, "Mensagem") {
			t.Errorf("Render(%d, %d): output missing title or message", dim.width, dim.height)
		}
	}
}

func TestConfirmModal_SingleOption_WithCustomKeys(t *testing.T) {
	// Teste com uma opção única usando teclas customizadas
	called := false
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Tab},
			Label:  "Ação",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd {
				called = true
				return nil
			},
		},
	}
	m := modal.NewConfirmModal("Título", "Mensagem", opts)

	// TAB deve funcionar (tecla original)
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	if !called {
		t.Error("Custom single key: TAB should work")
	}

	// ESC deve ser adicionado automaticamente mesmo com TAB como tecla principal
	called = false
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEscape})
	if !called {
		t.Error("Custom single key: ESC should be auto-added and work")
	}
}

func TestConfirmModal_IntentTypes_Preserved(t *testing.T) {
	// Verificar se os tipos de Intent são preservados corretamente
	// Não podemos acessar campos privados de ConfirmModal diretamente,
	// então testamos através do comportamento observável
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "Confirmar",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd { return nil },
		},
		{
			Keys:   []design.Key{design.Keys.Esc},
			Label:  "Cancelar",
			Intent: modal.IntentCancel,
			Action: func() tea.Cmd { return nil },
		},
	}
	m := modal.NewConfirmModal("Título", "Mensagem", opts)

	// Verificar que o modal foi criado com sucesso
	if m == nil {
		t.Fatal("modal is nil")
	}

	// Render deve funcionar, indicando que as opções estão presentes
	rendered := m.Render(10, 60, design.TokyoNight)
	if rendered == "" {
		t.Error("Modal with IntentTypes should render successfully")
	}
}

func TestConfirmModal_Update_ChainedMessages(t *testing.T) {
	// Teste simulando uma sequência de mensagens de diferentes tipos
	callCount := 0
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "OK",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd {
				callCount++
				return nil
			},
		},
	}
	m := modal.NewConfirmModal("Teste", "Mensagem", opts)

	// Sequência: WindowSize -> Key -> WindowSize -> Key
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})

	if callCount != 2 {
		t.Errorf("Chained messages: callCount = %d, want 2", callCount)
	}
}

func TestConfirmModal_EmptyOptions_HandleKey_ReturnsNil(t *testing.T) {
	// Teste com modal sem opções (edge case)
	m := modal.NewConfirmModal("Título", "Mensagem", []modal.ModalOption{})

	cmd := m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd != nil {
		t.Error("Empty options: HandleKey should return nil")
	}
}
