package tui

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

// ClearScreenMsg é uma mensagem enviada para limpar a tela antes de sair.
// Usado para garantir que o terminal seja restaurado para um estado limpo.
type ClearScreenMsg struct{}

// QuitWithCleanup cria um comando Bubble Tea que:
// 1. Desabilita o alternate screen buffer do Bubble Tea
// 2. Limpa toda a tela (incluindo scroll back)
// 3. Move o cursor para home (0,0)
// 4. Restaura o terminal para o estado normal
// 5. Encerra a aplicação
//
// Isso garante que quando a aplicação sai, o terminal não mostra mais
// o conteúdo da aplicação TUI, mesmo quando rolando para cima.
func QuitWithCleanup() tea.Cmd {
	return func() tea.Msg {
		// Sequência de escape ANSI para limpeza completa:
		// ESC[?1049l = Desabilita o alternate screen buffer (volta ao screen normal)
		// ESC[2J = Apaga toda a tela
		// ESC[3J = Apaga o scroll back buffer (apenas em alguns terminais)
		// ESC[H = Move cursor para home (0,0)
		fmt.Fprint(os.Stdout, "\033[?1049l") // Desabilita alternate screen
		fmt.Fprint(os.Stdout, "\033[2J")     // Apaga tela
		fmt.Fprint(os.Stdout, "\033[3J")     // Apaga scroll back (Linux/Unix)
		fmt.Fprint(os.Stdout, "\033[H")      // Cursor para home
		fmt.Fprint(os.Stdout, "\033[?25h")   // Mostra o cursor
		os.Stdout.Sync()                     // Força sincronização com o terminal

		// Retorna o comando de quit do Bubble Tea
		return tea.QuitMsg{}
	}
}
