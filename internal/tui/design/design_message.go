package design

import "charm.land/lipgloss/v2"

// MessageKind representa o tipo semântico de uma mensagem de status.
// Define o símbolo, cor e TTL padrão de cada tipo.
type MessageKind int

const (
	// MsgSuccess indica operação concluída com êxito.
	MsgSuccess MessageKind = iota
	// MsgInfo é uma mensagem informativa neutra.
	MsgInfo
	// MsgWarning indica situação que requer atenção.
	MsgWarning
	// MsgError indica falha ou estado inválido.
	MsgError
	// MsgBusy indica operação em andamento — exibe spinner animado.
	MsgBusy
	// MsgHintField exibe dica para o campo com foco atual.
	MsgHintField
	// MsgHintUsage exibe dica de uso geral do contexto.
	MsgHintUsage
)

// Symbol retorna o símbolo Unicode estático do tipo de mensagem.
// Para MsgBusy retorna SpinnerFrames[0] como placeholder — o frame real é
// lido de Message.SpinnerFrame em Message.Render.
func (k MessageKind) Symbol() string {
	switch k {
	case MsgSuccess:
		return SymSuccess
	case MsgInfo:
		return SymInfo
	case MsgWarning:
		return SymWarning
	case MsgError:
		return SymError
	case MsgBusy:
		return SpinnerFrames[0]
	default: // MsgHintField, MsgHintUsage
		return SymBullet
	}
}

// Color retorna o token de cor hex para o tipo, retirado do tema fornecido.
func (k MessageKind) Color(theme *Theme) string {
	switch k {
	case MsgSuccess:
		return theme.Semantic.Success
	case MsgInfo:
		return theme.Semantic.Info
	case MsgWarning:
		return theme.Semantic.Warning
	case MsgError:
		return theme.Semantic.Error
	case MsgBusy:
		return theme.Accent.Primary
	default: // MsgHintField, MsgHintUsage
		return theme.Text.Secondary
	}
}

// Style retorna o estilo lipgloss para o tipo, incluindo cor e atributos tipográficos.
// MsgError adiciona bold; MsgHintField e MsgHintUsage adicionam italic.
func (k MessageKind) Style(theme *Theme) lipgloss.Style {
	base := lipgloss.NewStyle().Foreground(lipgloss.Color(k.Color(theme)))
	switch k {
	case MsgError:
		return base.Bold(true)
	case MsgHintField, MsgHintUsage:
		return base.Italic(true)
	default:
		return base
	}
}

// DefaultTTL retorna o TTL padrão em ticks (segundos) para o tipo.
// 0 significa permanente — a mensagem só é removida por SetXxx ou Clear explícito.
func (k MessageKind) DefaultTTL() int {
	switch k {
	case MsgSuccess, MsgInfo, MsgWarning, MsgError:
		return 5
	default: // MsgBusy, MsgHintField, MsgHintUsage
		return 0
	}
}

// Message representa uma mensagem de status a ser exibida na barra.
// O TTL não é armazenado aqui — é gerenciado pelo MessageLineView.
// SpinnerFrame guarda o frame atual da animação quando Kind == MsgBusy;
// é irrelevante para outros tipos.
type Message struct {
	Kind         MessageKind
	Text         string
	SpinnerFrame int // frame atual da animação; só relevante quando Kind == MsgBusy
}

// Render produz a string ANSI renderizada do conteúdo da mensagem.
// maxWidth é o espaço disponível em colunas (calculado pelo caller como terminalWidth-6).
// Para MsgBusy, usa m.SpinnerFrame para determinar o frame da animação.
// Retorna a string renderizada e o número real de colunas ocupadas.
func (m Message) Render(theme *Theme, maxWidth int) (string, int) {
	style := m.Kind.Style(theme)

	var sym string
	if m.Kind == MsgBusy {
		// SpinnerFrame pode ter largura ambígua em alguns locales — medir com lipgloss.Width.
		sym = SpinnerFrame(m.SpinnerFrame)
	} else {
		sym = m.Kind.Symbol()
	}

	// MsgBusy com texto: "spinner  texto"; sem texto: apenas "spinner".
	// Demais tipos: "símbolo  texto".
	var raw string
	if m.Kind == MsgBusy && m.Text == "" {
		raw = sym
	} else {
		raw = sym + "  " + m.Text
	}

	// Truncar se necessário, garantindo que caiba em maxWidth colunas.
	// Usar lipgloss.Width para suporte correto a caracteres multibyte — nunca len().
	if lipgloss.Width(raw) > maxWidth {
		runes := []rune(raw)
		for len(runes) > 0 {
			candidate := string(runes) + SymEllipsis
			if lipgloss.Width(candidate) <= maxWidth {
				raw = candidate
				break
			}
			runes = runes[:len(runes)-1]
		}
		if len(runes) == 0 {
			raw = ""
		}
	}

	rendered := style.Render(raw)
	cols := lipgloss.Width(rendered)
	return rendered, cols
}

// Success cria uma mensagem de sucesso.
func Success(text string) Message { return Message{Kind: MsgSuccess, Text: text} }

// Error cria uma mensagem de erro.
func Error(text string) Message { return Message{Kind: MsgError, Text: text} }

// Info cria uma mensagem informativa.
func Info(text string) Message { return Message{Kind: MsgInfo, Text: text} }

// Warning cria uma mensagem de aviso.
func Warning(text string) Message { return Message{Kind: MsgWarning, Text: text} }

// Busy cria uma mensagem de operação em andamento.
// text é opcional — se fornecido, é exibido ao lado do spinner.
// SpinnerFrame é 0 (zero value) — MessageLineView avança o frame a cada TickMsg.
func Busy(text ...string) Message {
	msg := Message{Kind: MsgBusy}
	if len(text) > 0 {
		msg.Text = text[0]
	}
	return msg
}

// HintField cria uma dica para o campo com foco atual.
func HintField(text string) Message { return Message{Kind: MsgHintField, Text: text} }

// HintUsage cria uma dica de uso geral do contexto.
func HintUsage(text string) Message { return Message{Kind: MsgHintUsage, Text: text} }
