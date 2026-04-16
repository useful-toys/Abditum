package modal

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// KeyHandler centraliza o despacho de teclas comuns a todos os diálogos:
// ações do rodapé (Options) e navegação de scroll (ScrollState).
//
// O modal concreto compõe um KeyHandler como campo (não embedded) e chama
// Handle() explicitamente — podendo interceptar teclas antes ou depois.
type KeyHandler struct {
	// Options lista as ações cujas teclas devem ser despachadas automaticamente.
	Options []ModalOption
	// Scroll é o ScrollState a ser atualizado pelas teclas de navegação.
	// nil = sem scroll; teclas de scroll não serão consumidas pelo handler.
	Scroll *ScrollState
}

// Handle processa a tecla fornecida.
//
// Retorna (cmd, true) se a tecla foi consumida — execução de ação ou movimento de scroll.
// Retorna (nil, false) se a tecla não foi reconhecida.
//
// Ordem de despacho:
//  1. Opções: itera Options, compara com cada Key em opt.Keys usando key.Matches(msg).
//     No primeiro match, executa opt.Action() e retorna (cmd, true).
//  2. Scroll (apenas se Scroll != nil):
//     ↑ → Scroll.Up(), ↓ → Scroll.Down()
//     PgUp → Scroll.PageUp(), PgDn → Scroll.PageDown()
//     Home → Scroll.Home(), End → Scroll.End()
//     Após atualizar o estado, retorna (nil, true).
func (h *KeyHandler) Handle(msg tea.KeyMsg) (tea.Cmd, bool) {
	// 1. Despachar ações registradas.
	for _, opt := range h.Options {
		for _, k := range opt.Keys {
			if k.Matches(msg) {
				return opt.Action(), true
			}
		}
	}

	// 2. Navegar scroll (se configurado).
	if h.Scroll == nil {
		return nil, false
	}
	switch {
	case design.Keys.Up.Matches(msg):
		h.Scroll.Up()
		return nil, true
	case design.Keys.Down.Matches(msg):
		h.Scroll.Down()
		return nil, true
	case design.Keys.PgUp.Matches(msg):
		h.Scroll.PageUp()
		return nil, true
	case design.Keys.PgDn.Matches(msg):
		h.Scroll.PageDown()
		return nil, true
	case design.Keys.Home.Matches(msg):
		h.Scroll.Home()
		return nil, true
	case design.Keys.End.Matches(msg):
		h.Scroll.End()
		return nil, true
	}
	return nil, false
}
