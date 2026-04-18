package design

import (
	"unicode"

	tea "charm.land/bubbletea/v2"
)

// Key associa o rótulo de UI da tecla com sua representação tipada no bubbletea.
// Label é exibido na barra de comandos e no diálogo de ajuda.
// Code e Mod são usados para comparação tipada com tea.KeyMsg — sem strings mágicas.
type Key struct {
	// Label é o rótulo exibido na UI: "Enter", "⌃Q", "F1".
	Label string
	// Code é o código da tecla: tea.KeyEnter, tea.KeyF1, 'q', etc.
	Code rune
	// Mod são os modificadores: tea.ModCtrl, tea.ModShift, etc.
	Mod tea.KeyMod
}

// Matches reporta se o evento de teclado corresponde a esta Key.
// Compara Code e Mod exatamente — sem parsing de string.
func (k Key) Matches(msg tea.KeyMsg) bool {
	key := msg.Key()
	return key.Code == k.Code && key.Mod == k.Mod
}

// Constantes de modificadores usadas para montar rótulos na barra de comandos e no diálogo de ajuda.
const (
	// ModLabelCtrl é a notação do modificador Ctrl conforme o design system.
	ModLabelCtrl = "⌃" // U+2303
	// ModLabelShift é a notação do modificador Shift conforme o design system.
	ModLabelShift = "⇧" // U+21E7
	// ModLabelAlt é a notação do modificador Alt conforme o design system.
	ModLabelAlt = "!" // sem glifo Unicode dedicado
)

// Keys contém as teclas simples pré-definidas, sem modificadores.
// Para teclas com modificadores, use as funções WithCtrl, WithShift, WithAlt e Letter.
var Keys = struct {
	Enter, Esc, Tab, Del, Ins, Home, End, PgUp, PgDn, Up, Down Key
	F1, F2, F3, F4, F5, F6, F7, F8, F9, F10, F11, F12          Key
}{
	Enter: Key{Label: "Enter", Code: tea.KeyEnter},
	Esc:   Key{Label: "Esc", Code: tea.KeyEscape},
	Tab:   Key{Label: "Tab", Code: tea.KeyTab},
	Del:   Key{Label: "Del", Code: tea.KeyDelete},
	Ins:   Key{Label: "Ins", Code: tea.KeyInsert},
	Home:  Key{Label: "Home", Code: tea.KeyHome},
	End:   Key{Label: "End", Code: tea.KeyEnd},
	PgUp:  Key{Label: "PgUp", Code: tea.KeyPgUp},
	PgDn:  Key{Label: "PgDn", Code: tea.KeyPgDown},
	Up:    Key{Label: "↑", Code: tea.KeyUp},
	Down:  Key{Label: "↓", Code: tea.KeyDown},
	F1:    Key{Label: "F1", Code: tea.KeyF1},
	F2:    Key{Label: "F2", Code: tea.KeyF2},
	F3:    Key{Label: "F3", Code: tea.KeyF3},
	F4:    Key{Label: "F4", Code: tea.KeyF4},
	F5:    Key{Label: "F5", Code: tea.KeyF5},
	F6:    Key{Label: "F6", Code: tea.KeyF6},
	F7:    Key{Label: "F7", Code: tea.KeyF7},
	F8:    Key{Label: "F8", Code: tea.KeyF8},
	F9:    Key{Label: "F9", Code: tea.KeyF9},
	F10:   Key{Label: "F10", Code: tea.KeyF10},
	F11:   Key{Label: "F11", Code: tea.KeyF11},
	F12:   Key{Label: "F12", Code: tea.KeyF12},
}

// Letter cria uma Key para uma tecla de letra sem modificadores.
// O Label usa a letra maiúscula por convenção de notação do design system.
func Letter(r rune) Key {
	return Key{
		Label: string(unicode.ToUpper(r)),
		Code:  r,
	}
}

// WithCtrl adiciona o modificador Ctrl à tecla base, prefixando ModLabelCtrl ao Label.
func WithCtrl(base Key) Key {
	return Key{Label: ModLabelCtrl + base.Label, Code: base.Code, Mod: base.Mod | tea.ModCtrl}
}

// WithShift adiciona o modificador Shift à tecla base, prefixando ModLabelShift ao Label.
func WithShift(base Key) Key {
	return Key{Label: ModLabelShift + base.Label, Code: base.Code, Mod: base.Mod | tea.ModShift}
}

// WithAlt adiciona o modificador Alt à tecla base, prefixando ModLabelAlt ao Label.
func WithAlt(base Key) Key {
	return Key{Label: ModLabelAlt + base.Label, Code: base.Code, Mod: base.Mod | tea.ModAlt}
}

// Shortcuts contém os atalhos globais do design system, ativos em qualquer contexto da aplicação.
var Shortcuts = struct {
	// Help abre e fecha o diálogo de ajuda.
	Help Key
	// ThemeToggle alterna entre os temas Tokyo Night e Cyberpunk.
	// Não é exibido na barra de comandos.
	ThemeToggle Key
	// Quit sai da aplicação (com confirmação quando há alterações não salvas).
	Quit Key
	// LockVault bloqueia o cofre imediatamente, descartando alterações sem confirmação.
	// O atalho complexo (⌃!⇧Q) é intencional para evitar acionamento acidental.
	LockVault Key
	// NewVault inicia o fluxo de criação de novo cofre (Fluxo 2).
	NewVault Key
	// OpenVault inicia o fluxo de abertura de cofre existente (Fluxo 1).
	OpenVault Key
}{
	Help:        Keys.F1,
	ThemeToggle: Keys.F12,
	Quit:        WithCtrl(Letter('q')),
	LockVault:   WithCtrl(WithAlt(WithShift(Letter('q')))),
	NewVault:    WithCtrl(Letter('n')),
	OpenVault:   WithCtrl(Letter('o')),
}
