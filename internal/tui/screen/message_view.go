package screen

import "github.com/useful-toys/abditum/internal/tui/design"

// MessageLineView renderiza a linha de mensagem de sistema.
// Não implementa ChildView — é um renderizador stateless chamado diretamente pelo root.
type MessageLineView struct{}

// NewMessageLineView cria uma nova instância da linha de mensagem.
func NewMessageLineView() *MessageLineView {
	return &MessageLineView{}
}

// Render retorna a linha de mensagem para exibição.
// message é o texto atual a ser exibido (pode ser vazio).
// Stub — retorna string vazia até implementação visual completa.
func (v *MessageLineView) Render(height, width int, theme *design.Theme, message string) string {
	return ""
}
