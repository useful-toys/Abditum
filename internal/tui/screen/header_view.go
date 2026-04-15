package screen

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// HeaderView renderiza o cabeçalho fixo de 2 linhas da aplicação.
// Implementa ChildView — suportará mouse e actions no futuro.
type HeaderView struct{}

// NewHeaderView cria uma nova instância do cabeçalho.
func NewHeaderView() *HeaderView {
	return &HeaderView{}
}

// Render retorna o cabeçalho com as dimensões fornecidas.
// Stub — retorna string vazia até implementação visual completa.
func (v *HeaderView) Render(height, width int, theme *design.Theme) string {
	return ""
}

// HandleKey não processa teclas nesta view.
func (v *HeaderView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }

// HandleEvent não processa eventos externos nesta view.
func (v *HeaderView) HandleEvent(event any) {}

// HandleTeaMsg não processa mensagens do framework nesta view.
func (v *HeaderView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }

// Update não altera o estado desta view em resposta a mensagens.
func (v *HeaderView) Update(msg tea.Msg) tea.Cmd { return nil }

// Actions retorna nil — HeaderView não possui actions próprias nesta sprint.
func (v *HeaderView) Actions() []tui.Action { return nil }
