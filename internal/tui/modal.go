package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// ModalView define o contrato para componentes de modal da interface.
// Modais são exibidos sobrepostos à área de trabalho e gerenciados por RootModel.
type ModalView interface {
	// Render retorna a representação visual do modal dentro dos limites fornecidos.
	// theme é passado por ponteiro para evitar cópia — design.Theme tem 400 bytes.
	Render(maxHeight, maxWidth int, theme *design.Theme) string
	// HandleKey processa eventos de teclado e retorna um comando ou nil.
	HandleKey(msg tea.KeyMsg) tea.Cmd
	// Update processa mensagens do Bubble Tea e atualiza o estado interno do modal.
	Update(msg tea.Msg) tea.Cmd
	// Cursor retorna a posição do cursor real para o modal ativo, ou nil se não houver cursor.
	// topY e leftX são as coordenadas absolutas do canto superior esquerdo do modal na tela.
	Cursor(topY, leftX int) *tea.Cursor
}

// OpenModalMsg é enviada para empilhar um novo modal na pilha de modais.
// Use OpenModal para criar o comando correspondente.
type OpenModalMsg struct {
	// Modal é o componente a ser exibido sobreposto à tela atual.
	Modal ModalView
}

// CloseModalMsg é enviada para fechar e remover o modal no topo da pilha.
// Use CloseModal para criar o comando correspondente.
type CloseModalMsg struct{}

// ModalReadyMsg é enviada pelo modal quando sua operação está concluída.
// O componente pai pode então coletar o resultado e encerrar o modal.
type ModalReadyMsg struct{}

// OpenModal cria um comando Bubble Tea para empilhar o modal fornecido.
func OpenModal(modal ModalView) tea.Cmd {
	return func() tea.Msg { return OpenModalMsg{Modal: modal} }
}

// CloseModal cria um comando Bubble Tea para remover o modal do topo da pilha.
func CloseModal() tea.Cmd {
	return func() tea.Msg { return CloseModalMsg{} }
}
