package screen

import (
	"github.com/useful-toys/abditum/internal/tui/design"
)

// ActionLineView renderiza a linha de ações disponíveis no contexto atual.
// Não implementa ChildView — é um renderizador stateless chamado diretamente pelo root.
type ActionLineView struct{}

// NewActionLineView cria uma nova instância da linha de ações.
func NewActionLineView() *ActionLineView {
	return &ActionLineView{}
}

// Render retorna a linha de ações para exibição.
// actions é a lista já filtrada (Visible + AvailableWhen satisfeita) e ordenada por Priority.
// Stub — retorna string vazia até implementação visual completa.
func (v *ActionLineView) Render(height, width int, theme *design.Theme, actions []interface{}) string {
	return ""
}
