package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

// ClearScreenMsg é uma mensagem enviada para limpar a tela antes de sair.
// Usado para garantir que o terminal seja restaurado para um estado limpo.
type ClearScreenMsg struct{}

// QuitWithCleanup cria um comando Bubble Tea que:
// 1. Limpa a tela da aplicação (zera o conteúdo)
// 2. Restaura o terminal para o estado normal
// 3. Encerra a aplicação
//
// Isso garante que quando a aplicação sai, o terminal não mostra mais
// o conteúdo da aplicação TUI.
func QuitWithCleanup() tea.Cmd {
	return func() tea.Msg {
		// Limpa a tela usando escape codes ANSI
		// ESC[2J = clear entire screen
		// ESC[H = move cursor to home (0,0)
		fmt.Print("\033[2J\033[H")

		// Retorna o comando de quit do Bubble Tea
		return tea.QuitMsg{}
	}
}
