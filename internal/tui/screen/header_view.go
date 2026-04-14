package screen

import "github.com/useful-toys/abditum/internal/tui/design"

// HeaderView renderiza o cabeçalho fixo de 2 linhas da aplicação.
// Não implementa ChildView — é um renderizador stateless chamado diretamente pelo root.
type HeaderView struct{}

// NewHeaderView cria uma nova instância do cabeçalho.
func NewHeaderView() *HeaderView {
	return &HeaderView{}
}

// Render retorna o cabeçalho com 2 linhas para exibição.
// Stub — retorna string vazia até implementação visual completa.
func (v *HeaderView) Render(height, width int, theme *design.Theme) string {
	return ""
}
