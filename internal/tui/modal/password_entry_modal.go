package modal

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// passwordEntryTitle é o título do modal conforme a spec.
const passwordEntryTitle = "Senha mestra"

// passwordEntryWidth é a largura fixa do modal em colunas.
const passwordEntryWidth = 50

// PasswordEntryModal exibe o diálogo de entrada de senha para abrir o cofre.
// Tem um único campo de senha. O modal notifica o orquestrador via callbacks —
// não sabe se a senha foi aceita ou rejeitada.
// Implementa tui.ModalView.
type PasswordEntryModal struct {
	mc        tui.MessageController
	field     *PasswordField
	onConfirm func(password []byte) tea.Cmd
	onCancel  func() tea.Cmd
}

// NewPasswordEntryModal cria o modal e emite a dica inicial na barra de status.
func NewPasswordEntryModal(
	mc tui.MessageController,
	onConfirm func(password []byte) tea.Cmd,
	onCancel func() tea.Cmd,
) *PasswordEntryModal {
	m := &PasswordEntryModal{
		mc:        mc,
		field:     NewPasswordField("Senha"),
		onConfirm: onConfirm,
		onCancel:  onCancel,
	}
	mc.SetHintField("• Digite a senha para desbloquear o cofre")
	return m
}

// Len retorna o comprimento atual do campo de senha.
// Usado pelos testes para verificar o estado do campo.
func (m *PasswordEntryModal) Len() int {
	return m.field.Len()
}

// NotifyWrongPassword limpa o campo para que o usuário possa tentar novamente.
// O orquestrador é responsável por exibir a mensagem de tentativa na barra de status.
func (m *PasswordEntryModal) NotifyWrongPassword() {
	m.field.Wipe()
}

// Render gera a representação visual do modal.
// Altura fixa: 5 linhas de corpo + 2 bordas = 7 linhas totais.
func (m *PasswordEntryModal) Render(maxHeight, maxWidth int, theme *design.Theme) string {
	innerWidth := passwordEntryWidth - 2 - 2*design.DialogPaddingH

	fieldRendered := m.field.Render(innerWidth, true, theme)

	// Body: padding + campo + padding inferior
	// Linha 0: vazia, Linha 1: label, Linha 2: área digitável, Linha 3: vazia
	body := "\n" + fieldRendered + "\n"

	confirmColor := theme.Text.Disabled
	if m.field.Len() > 0 {
		confirmColor = theme.Accent.Primary
	}

	opts := []ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "Confirmar",
			Intent: IntentConfirm,
			Action: func() tea.Cmd {
				if m.field.Len() == 0 {
					return nil
				}
				return m.onConfirm(m.field.Value())
			},
		},
		{
			Keys:   []design.Key{design.Keys.Esc},
			Label:  "Cancelar",
			Intent: IntentCancel,
			Action: func() tea.Cmd {
				m.field.Wipe()
				return m.onCancel()
			},
		},
	}

	frame := DialogFrame{
		Title:           passwordEntryTitle,
		TitleColor:      theme.Text.Primary,
		Symbol:          "",
		SymbolColor:     "",
		BorderColor:     theme.Border.Focused,
		Options:         opts,
		DefaultKeyColor: confirmColor,
		Scroll:          nil,
	}
	return frame.Render(body, passwordEntryWidth, theme)
}

// HandleKey processa eventos de teclado.
// Enter: confirma se campo não vazio. Esc: cancela e faz wipe. Tab: no-op (campo único).
func (m *PasswordEntryModal) HandleKey(msg tea.KeyMsg) tea.Cmd {
	key := msg.Key()

	switch key.Code {
	case tea.KeyEnter:
		if m.field.Len() == 0 {
			return nil
		}
		return m.onConfirm(m.field.Value())
	case tea.KeyEsc:
		m.field.Wipe()
		return m.onCancel()
	case tea.KeyTab:
		return nil // campo único — Tab não faz nada
	}
	// Delegate to field for other keys (printable characters, backspace)
	if pressMsg, ok := msg.(tea.KeyPressMsg); ok {
		m.field.HandleKey(pressMsg)
	}
	return nil
}

// Update processa mensagens Bubble Tea. Delega para HandleKey em tea.KeyMsg.
func (m *PasswordEntryModal) Update(msg tea.Msg) tea.Cmd {
	if key, ok := msg.(tea.KeyMsg); ok {
		return m.HandleKey(key)
	}
	return nil
}

// Cursor retorna a posição do cursor real para o campo de senha.
// topY e leftX são as coordenadas absolutas do canto superior esquerdo do modal na tela.
//
// Mapa de linhas do body (0-indexed):
//
//	Linha 0: vazia (padding)
//	Linha 1: label "Senha"
//	Linha 2: área digitável  ← cursor aqui
//	Linha 3: vazia
//
// Fórmula:
//
//	cursorY = topY + 1 (borda superior) + 2 (linha do field no body)
//	cursorX = leftX + 1 (borda esquerda) + DialogPaddingH + field.Len()
func (m *PasswordEntryModal) Cursor(topY, leftX int) *tea.Cursor {
	y := topY + 1 + 2
	x := leftX + 1 + design.DialogPaddingH + m.field.Len()
	return tea.NewCursor(x, y)
}
