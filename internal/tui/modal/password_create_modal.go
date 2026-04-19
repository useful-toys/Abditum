package modal

import (
	"bytes"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// passwordCreateTitle é o título do modal conforme a spec.
const passwordCreateTitle = "Definir senha mestra"

// passwordCreateWidth é a largura fixa do modal em colunas.
const passwordCreateWidth = 50

// focusedField é uma enumeração para rastrear qual campo está em foco.
type focusedField int

const (
	// fieldNew representa o campo "Nova senha".
	fieldNew focusedField = iota
	// fieldConfirm representa o campo "Confirmação".
	fieldConfirm
)

// PasswordCreateModal exibe o diálogo de criação de senha mestra com dois campos de senha
// (Nova senha e Confirmação) e um medidor de força de senha.
// Implementa tui.ModalView.
type PasswordCreateModal struct {
	mc        tui.MessageController
	fieldNew  *PasswordField
	fieldConf *PasswordField
	focused   focusedField
	onConfirm func(password []byte) tea.Cmd
	onCancel  func() tea.Cmd
}

// NewPasswordCreateModal cria o modal e emite a dica inicial na barra de status.
func NewPasswordCreateModal(
	mc tui.MessageController,
	onConfirm func(password []byte) tea.Cmd,
	onCancel func() tea.Cmd,
) *PasswordCreateModal {
	m := &PasswordCreateModal{
		mc:        mc,
		fieldNew:  NewPasswordField("Nova senha"),
		fieldConf: NewPasswordField("Confirmação"),
		focused:   fieldNew,
		onConfirm: onConfirm,
		onCancel:  onCancel,
	}
	mc.SetHintField("• A senha mestra protege todo o cofre — use 12+ caracteres")
	return m
}

// FocusedOnNew retorna true se o foco está no campo Nova senha.
func (m *PasswordCreateModal) FocusedOnNew() bool {
	return m.focused == fieldNew
}

// canConfirm verifica se ambos os campos estão preenchidos e se as senhas são iguais.
func (m *PasswordCreateModal) canConfirm() bool {
	if m.fieldNew.Len() == 0 || m.fieldConf.Len() == 0 {
		return false
	}
	vNew := m.fieldNew.Value()
	vConf := m.fieldConf.Value()
	defer func() {
		for i := range vNew {
			vNew[i] = 0
		}
		for i := range vConf {
			vConf[i] = 0
		}
	}()
	return bytes.Equal(vNew, vConf)
}

// switchFocus alterna o foco entre os dois campos e emite a mensagem apropriada.
func (m *PasswordCreateModal) switchFocus() {
	if m.focused == fieldNew {
		m.focused = fieldConfirm
		m.validateConfirmation()
	} else {
		m.focused = fieldNew
		m.mc.SetHintField("• A senha mestra protege todo o cofre — use 12+ caracteres")
	}
}

// validateConfirmation valida o campo Confirmação quando ele está em foco
// e emite a mensagem apropriada baseada no match ou não.
func (m *PasswordCreateModal) validateConfirmation() {
	if m.fieldConf.Len() == 0 {
		m.mc.SetHintField("• Redigite a senha para confirmar")
	} else {
		vNew := m.fieldNew.Value()
		vConf := m.fieldConf.Value()
		defer func() {
			for i := range vNew {
				vNew[i] = 0
			}
			for i := range vConf {
				vConf[i] = 0
			}
		}()
		if bytes.Equal(vNew, vConf) {
			m.mc.SetHintField("• Redigite a senha para confirmar")
		} else {
			m.mc.SetHintField("✕ As senhas não conferem — digite novamente")
		}
	}
}

// Render gera a representação visual do modal.
// Altura: 9 linhas de corpo sem meter, 11 linhas com meter (quando Nova senha não vazio) + 2 bordas.
func (m *PasswordCreateModal) Render(maxHeight, maxWidth int, theme *design.Theme) string {
	innerWidth := passwordCreateWidth - 2 - 2*design.DialogPaddingH

	// Render both fields
	fieldNewRendered := m.fieldNew.Render(innerWidth, m.focused == fieldNew, theme)
	fieldConfRendered := m.fieldConf.Render(innerWidth, m.focused == fieldConfirm, theme)

	// Body structure:
	// Linha 0: vazia (padding)
	// Linha 1: label "Nova senha"
	// Linha 2: input area Nova senha
	// Linha 3: vazia (single blank line separator)
	// Linha 4: label "Confirmação"
	// Linha 5: input area Confirmação
	// Linha 6: vazia
	// [Linha 7: strength meter]     ← only when Nova senha not empty
	// [Linha 8: empty after meter]  ← only when Nova senha not empty
	body := "\n" + fieldNewRendered + "\n" + fieldConfRendered + "\n"

	// Add strength meter if Nova senha is not empty
	if m.fieldNew.Len() > 0 {
		meter := RenderStrengthMeter(m.fieldNew.Value(), innerWidth, theme)
		body += meter + "\n"
	}

	// Determine button colors
	confirmColor := theme.Text.Disabled
	if m.canConfirm() {
		confirmColor = theme.Accent.Primary
	}

	opts := []ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "Confirmar",
			Action: func() tea.Cmd {
				if !m.canConfirm() {
					if m.fieldConf.Len() > 0 {
						m.mc.SetError("As senhas não conferem")
					}
					return nil
				}
				pwd := m.fieldNew.Value()
				m.fieldConf.Wipe()
				return m.onConfirm(pwd)
			},
		},
		{
			Keys:   []design.Key{design.Keys.Esc},
			Label:  "Cancelar",
			Action: func() tea.Cmd {
				m.fieldNew.Wipe()
				m.fieldConf.Wipe()
				return m.onCancel()
			},
		},
	}

	frame := DialogFrame{
		Title:           passwordCreateTitle,
		TitleColor:      theme.Text.Primary,
		Symbol:          "",
		SymbolColor:     "",
		BorderColor:     theme.Border.Focused,
		Options:         opts,
		DefaultKeyColor: confirmColor,
		Scroll:          nil,
	}
	return frame.Render(body, passwordCreateWidth, theme)
}

// HandleKey processa eventos de teclado.
// Tab: alterna o foco entre os campos.
// Enter: confirma se canConfirm() retornar true; emite erro se mismatch.
// Esc: limpa ambos os campos e cancela.
// Outros: delega para o campo em foco.
func (m *PasswordCreateModal) HandleKey(msg tea.KeyMsg) tea.Cmd {
	key := msg.Key()

	switch key.Code {
	case tea.KeyTab:
		m.switchFocus()
		return nil
	case tea.KeyEnter:
		if !m.canConfirm() {
			if m.fieldConf.Len() > 0 {
				m.mc.SetError("As senhas não conferem")
			}
			return nil
		}
		pwd := m.fieldNew.Value()
		m.fieldConf.Wipe()
		return m.onConfirm(pwd)
	case tea.KeyEsc:
		m.fieldNew.Wipe()
		m.fieldConf.Wipe()
		return m.onCancel()
	}

	// Delegate to focused field for other keys (printable characters, backspace)
	if m.focused == fieldNew {
		if pressMsg, ok := msg.(tea.KeyPressMsg); ok {
			m.fieldNew.HandleKey(pressMsg)
		}
	} else {
		if pressMsg, ok := msg.(tea.KeyPressMsg); ok {
			m.fieldConf.HandleKey(pressMsg)
			m.validateConfirmation()
		}
	}
	return nil
}

// HandleMouse processa eventos de mouse. PasswordCreateModal não reage a mouse — retorna nil.
func (m *PasswordCreateModal) HandleMouse(_ tea.MouseMsg) tea.Cmd {
	return nil
}

// Cursor retorna a posição do cursor real para o campo em foco.
// topY e leftX são as coordenadas absolutas do canto superior esquerdo do modal na tela.
//
// Mapa de linhas do body (0-indexed):
//
//	Linha 0: vazia (padding)
//	Linha 1: label "Nova senha"
//	Linha 2: área digitável do Nova senha  ← cursor aqui se focused == fieldNew
//	Linha 3: label "Confirmação"
//	Linha 4: área digitável do Confirmação  ← cursor aqui se focused == fieldConfirm
//	Linha 5: vazia
//
// Fórmula (Nova senha):
//
//	cursorY = topY + 1 (borda superior) + 2 (linha do field no body)
//	cursorX = leftX + 1 (borda esquerda) + DialogPaddingH + field.Len()
//
// Fórmula (Confirmação):
//
//	cursorY = topY + 1 (borda superior) + 4 (linha do field no body)
//	cursorX = leftX + 1 (borda esquerda) + DialogPaddingH + field.Len()
func (m *PasswordCreateModal) Cursor(topY, leftX int) *tea.Cursor {
	var lineOffset int
	var field *PasswordField

	if m.focused == fieldNew {
		lineOffset = 2
		field = m.fieldNew
	} else {
		lineOffset = 4
		field = m.fieldConf
	}

	y := topY + 1 + lineOffset
	x := leftX + 1 + design.DialogPaddingH + field.Len()
	return tea.NewCursor(x, y)
}
