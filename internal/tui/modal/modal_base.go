package modal

import tea "charm.land/bubbletea/v2"

// Intent classifica a intenção semântica de uma opção de modal.
// Permite que o componente pai interprete o resultado sem inspecionar o label.
type Intent int

const (
	// IntentConfirm indica que a ação confirma a operação em andamento.
	IntentConfirm Intent = iota
	// IntentCancel indica que a ação cancela e retorna ao estado anterior.
	IntentCancel
	// IntentOther indica uma ação auxiliar sem intenção semântica predefinida.
	IntentOther
)

// ModalOption representa uma ação disponível ao usuário dentro de um modal.
// Cada opção é ativada por uma ou mais teclas de atalho.
type ModalOption struct {
	// Keys lista as teclas que ativam esta opção, ex: []string{"Enter", "y"}.
	Keys []string
	// Label é o texto exibido ao usuário para descrever a ação.
	Label string
	// Intent classifica a intenção semântica desta ação.
	Intent Intent
	// Action é a função executada quando a opção é escolhida pelo usuário.
	Action func() tea.Cmd
}
