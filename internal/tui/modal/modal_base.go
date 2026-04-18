package modal

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// ModalOption representa uma ação disponível ao usuário dentro de um modal.
type ModalOption struct {
	// Keys lista as teclas que ativam esta opção.
	// Keys[0].Label é exibido no rodapé do diálogo.
	// Outras Keys são aliases funcionais (ex: Enter como alias para "S Sobrescrever").
	//
	// Keys é opcional (nil ou vazio). Quando omitido, teclas implícitas são aplicadas
	// pelo KeyHandler e pelo DialogFrame:
	//   - Primeira option → Enter
	//   - Última option   → Esc
	//   - Option única    → Enter e Esc
	// Quando Keys está preenchido, as teclas implícitas são adicionadas como aliases.
	Keys []design.Key
	// Label é o texto exibido ao usuário descrevendo a ação.
	Label string
	// Action é a função executada quando a opção é escolhida.
	Action func() tea.Cmd
}
