package actions

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/vault"
)

// ActionGroup agrupa actions relacionadas para exibição no modal de ajuda.
type ActionGroup struct {
	ID          string // identificador único do grupo
	Label       string // cabeçalho exibido no modal de ajuda
	Description string // texto descritivo do grupo
	Order       int    // ordem de exibição no modal de ajuda; menor valor aparece primeiro
}

// Action associa teclas a um comportamento da aplicação.
type Action struct {
	Keys          []design.Key                            // Keys[0] é a tecla principal exibida; demais são aliases funcionais
	Label         string                                  // texto curto para a linha de status
	Description   string                                  // texto longo para o modal de ajuda
	GroupID       string                                  // referencia um ActionGroup registrado
	Priority      int                                     // ordenação na linha de status; menor valor = mais destaque
	Visible       bool                                    // false: nunca aparece na linha de status
	AvailableWhen func(app AppState, view ChildView) bool // nil = sempre disponível
	OnExecute     func() tea.Cmd
}

// Matches retorna true se o evento de teclado corresponde a qualquer tecla declarada na action.
func (a Action) Matches(msg tea.KeyMsg) bool {
	for _, k := range a.Keys {
		if k.Matches(msg) {
			return true
		}
	}
	return false
}

// AppState expõe o estado da aplicação necessário para avaliar pré-condições de actions.
// Implementado por RootModel.
type AppState interface {
	Manager() *vault.Manager // nil se nenhum cofre estiver carregado
}

// ChildView é o subconjunto da interface tui.ChildView necessário para AvailableWhen.
// Definida aqui para evitar import cycle. Nenhuma action atual inspeciona a view —
// a interface está vazia agora, mas nomear o tipo documenta a intenção e permite
// adicionar métodos futuramente sem alterar a assinatura de AvailableWhen.
type ChildView interface{}
