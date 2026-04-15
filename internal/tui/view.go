package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/actions"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// ChildView define o contrato para componentes renderizáveis da tela principal.
type ChildView interface {
	// Render retorna a representação em string do componente para exibição.
	// height e width definem as dimensões disponíveis. theme é passado por ponteiro
	// para evitar cópia desnecessária — design.Theme tem 400 bytes.
	Render(height, width int, theme *design.Theme) string

	// HandleKey processa eventos de teclado e retorna um comando ou nil.
	HandleKey(msg tea.KeyMsg) tea.Cmd

	// HandleEvent processa eventos customizados para o componente.
	HandleEvent(event any)

	// HandleTeaMsg processa mensagens do Bubble Tea framework.
	HandleTeaMsg(msg tea.Msg) tea.Cmd

	// Update é chamado para atualizar o estado do componente.
	Update(msg tea.Msg) tea.Cmd

	// Actions retorna as actions disponíveis nesta view.
	// Pode retornar nil se a view não possuir actions próprias.
	Actions() []actions.Action
}
