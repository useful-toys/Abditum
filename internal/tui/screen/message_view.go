package screen

import (
	"reflect"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// MessageLineView gerencia estado e renderização da barra de mensagens de status.
// Implementa tui.MessageController via pointer receivers.
// O zero value é válido — representa estado sem mensagem ativa.
type MessageLineView struct {
	// current é a mensagem atualmente exibida. Zero value = sem mensagem.
	// Para MsgBusy, current.SpinnerFrame guarda o frame atual da animação.
	current design.Message
	// ttl é o número de ticks restantes antes da mensagem ser removida.
	// 0 significa permanente (MsgBusy, MsgHintField, MsgHintUsage nunca expiram).
	ttl int
}

// SetBusy exibe o spinner com texto opcional de status.
// design.Busy() retorna zero value de Message — SpinnerFrame começa em 0.
func (v *MessageLineView) SetBusy(text string) {
	v.current = design.Busy(text)
	v.ttl = 0
}

// SetSuccess exibe mensagem de sucesso. Desaparece após 5 ticks.
func (v *MessageLineView) SetSuccess(text string) {
	v.current = design.Success(text)
	v.ttl = v.current.Kind.DefaultTTL()
}

// SetError exibe mensagem de erro em bold. Desaparece após 5 ticks.
func (v *MessageLineView) SetError(text string) {
	v.current = design.Error(text)
	v.ttl = v.current.Kind.DefaultTTL()
}

// SetWarning exibe mensagem de aviso. Desaparece após 5 ticks.
func (v *MessageLineView) SetWarning(text string) {
	v.current = design.Warning(text)
	v.ttl = v.current.Kind.DefaultTTL()
}

// SetInfo exibe mensagem informativa. Desaparece após 5 ticks.
func (v *MessageLineView) SetInfo(text string) {
	v.current = design.Info(text)
	v.ttl = v.current.Kind.DefaultTTL()
}

// SetHintField exibe dica para o campo focado. Permanente até substituição.
func (v *MessageLineView) SetHintField(text string) {
	v.current = design.HintField(text)
	v.ttl = 0
}

// SetHintUsage exibe dica de uso geral. Permanente até substituição.
func (v *MessageLineView) SetHintUsage(text string) {
	v.current = design.HintUsage(text)
	v.ttl = 0
}

// Clear remove a mensagem atual imediatamente.
func (v *MessageLineView) Clear() {
	v.current = design.Message{}
	v.ttl = 0
}

// tick executa a lógica de um TickMsg: avança spinner e decrementa TTL.
// Separado de Update para ser testável sem importar package tui.
func (v *MessageLineView) tick() tea.Cmd {
	if v.current.Kind == design.MsgBusy {
		v.current.SpinnerFrame = (v.current.SpinnerFrame + 1) % 4
	}
	if v.ttl > 0 {
		v.ttl--
		if v.ttl == 0 {
			v.current = design.Message{}
		}
	}
	return nil
}

// Update trata mensagens do Bubble Tea. Delega tui.TickMsg para tick().
// Usa reflect para identificar o tipo e evitar import cycle com package tui.
func (v *MessageLineView) Update(msg tea.Msg) tea.Cmd {
	if reflect.TypeOf(msg).String() == "tui.TickMsg" {
		return v.tick()
	}
	return nil
}

// Render produz a string ANSI da barra de mensagens com exatamente `width` colunas.
//
// Estrutura:
//
//	Sem mensagem: ────────────────── (SymBorderH × width)
//	Com mensagem: ─── ✓  texto ─────
func (v *MessageLineView) Render(width int, theme *design.Theme) string {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Border.Default))

	// Sem mensagem ativa — linha de borda completa.
	if v.current == (design.Message{}) {
		return borderStyle.Render(strings.Repeat(design.SymBorderH, width))
	}

	// Prefixo: 3× SymBorderH + 1 espaço = 4 colunas fixas.
	prefix := borderStyle.Render(strings.Repeat(design.SymBorderH, 3)) + " "

	// Espaço disponível para o conteúdo.
	// 4 = prefixo, 2 = sufixo mínimo (1 espaço + 1× SymBorderH).
	maxWidth := width - 6
	if maxWidth < 1 {
		maxWidth = 1
	}

	content, contentCols := v.current.Render(theme, maxWidth)

	// Sufixo: espaço separador + SymBorderH repetido até preencher.
	sufixoCols := width - 4 - contentCols - 1
	if sufixoCols < 1 {
		sufixoCols = 1
	}
	suffix := " " + borderStyle.Render(strings.Repeat(design.SymBorderH, sufixoCols))

	return prefix + content + suffix
}
