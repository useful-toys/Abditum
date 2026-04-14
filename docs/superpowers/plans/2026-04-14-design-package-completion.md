# Design Package Completion — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Completar o pacote `internal/tui/design` com `LogoGradient` + constantes de layout, inventário de símbolos Unicode, e o tipo `Key` com teclas pré-definidas e funções de composição.

**Architecture:** Três adições ao pacote `design`: extensão de `design.go` (LogoGradient + layout), novo `symbols.go` (constantes Unicode), novo `keys.go` (tipo `Key` com `Matches`, `var Keys`, funções `WithCtrl`/`WithShift`/`WithAlt`/`Letter`, `var Shortcuts`). Nenhuma interface pública existente é alterada.

**Tech Stack:** Go 1.26, `charm.land/bubbletea/v2` (para `tea.KeyMod` e constantes de código de tecla em `keys.go`), `unicode` stdlib.

**Spec:** `docs/superpowers/specs/2026-04-14-design-package-completion-design.md`

---

## File Map

| Arquivo | Ação | Responsabilidade |
|---|---|---|
| `internal/tui/design/design.go` | Modificar | Adicionar `LogoGradient [5]string` ao `Theme` + valores nos dois temas + constantes `MinWidth`, `MinHeight`, `PanelTreeRatio` |
| `internal/tui/design/design_test.go` | Criar | Testes de `LogoGradient` e das constantes de layout |
| `internal/tui/design/symbols.go` | Criar | Constantes de símbolos Unicode + `SpinnerFrames` + `SpinnerFrame()` |
| `internal/tui/design/symbols_test.go` | Criar | Testes de todos os símbolos e do spinner |
| `internal/tui/design/keys.go` | Criar | `type Key`, `Matches()`, `ModLabel*`, `var Keys`, `Letter`/`WithCtrl`/`WithShift`/`WithAlt`, `var Shortcuts` |
| `internal/tui/design/keys_test.go` | Criar | Testes de labels, `Matches`, composição e `Shortcuts` |

---

## Task 1: Estender design.go com LogoGradient e constantes de layout

**Files:**
- Modify: `internal/tui/design/design.go`
- Create: `internal/tui/design/design_test.go`

- [ ] **Step 1: Criar design_test.go com os testes**

Criar `internal/tui/design/design_test.go`:

```go
package design

import "testing"

func TestTokyoNight_LogoGradient(t *testing.T) {
	want := [5]string{"#9d7cd8", "#89ddff", "#7aa2f7", "#7dcfff", "#bb9af7"}
	if TokyoNight.LogoGradient != want {
		t.Errorf("TokyoNight.LogoGradient = %v, want %v", TokyoNight.LogoGradient, want)
	}
}

func TestCyberpunk_LogoGradient(t *testing.T) {
	want := [5]string{"#ff2975", "#b026ff", "#00fff5", "#05ffa1", "#ff2975"}
	if Cyberpunk.LogoGradient != want {
		t.Errorf("Cyberpunk.LogoGradient = %v, want %v", Cyberpunk.LogoGradient, want)
	}
}

func TestLayoutConstants(t *testing.T) {
	if MinWidth != 80 {
		t.Errorf("MinWidth = %d, want 80", MinWidth)
	}
	if MinHeight != 24 {
		t.Errorf("MinHeight = %d, want 24", MinHeight)
	}
	const wantRatio = 0.35
	if PanelTreeRatio != wantRatio {
		t.Errorf("PanelTreeRatio = %v, want %v", PanelTreeRatio, wantRatio)
	}
}
```

- [ ] **Step 2: Rodar para confirmar falha de compilação**

```
go test ./internal/tui/design/
```

Esperado: erro de compilação com `undefined: LogoGradient`, `undefined: MinWidth`, etc.

- [ ] **Step 3: Adicionar LogoGradient ao struct Theme**

Em `internal/tui/design/design.go`, adicionar o campo ao final do struct `Theme`:

```go
type Theme struct {
	Name     string
	Surface  SurfaceTokens
	Text     TextTokens
	Accent   AccentTokens
	Border   BorderTokens
	Semantic SemanticTokens
	Special  SpecialTokens
	// LogoGradient são as 5 cores para as linhas do logo ASCII, de cima para baixo.
	LogoGradient [5]string
}
```

- [ ] **Step 4: Adicionar LogoGradient a TokyoNight**

Na var `TokyoNight`, adicionar antes do `}`:

```go
	LogoGradient: [5]string{"#9d7cd8", "#89ddff", "#7aa2f7", "#7dcfff", "#bb9af7"},
```

- [ ] **Step 5: Adicionar LogoGradient a Cyberpunk**

Na var `Cyberpunk`, adicionar antes do `}`:

```go
	LogoGradient: [5]string{"#ff2975", "#b026ff", "#00fff5", "#05ffa1", "#ff2975"},
```

- [ ] **Step 6: Adicionar constantes de layout ao final de design.go**

```go
// MinWidth é a largura mínima suportada do terminal em colunas (padrão POSIX).
const MinWidth = 80

// MinHeight é a altura mínima suportada do terminal em linhas (padrão POSIX).
const MinHeight = 24

// PanelTreeRatio é a proporção nominal do painel de navegação esquerdo em modos de dois painéis.
// O painel direito (detalhe) ocupa o restante. A implementação pode ajustar ±5%.
const PanelTreeRatio = 0.35
```

- [ ] **Step 7: Rodar testes e confirmar passagem**

```
go test ./internal/tui/design/
```

Esperado: `ok  github.com/useful-toys/abditum/internal/tui/design`

- [ ] **Step 8: Commit**

```
git add internal/tui/design/design.go internal/tui/design/design_test.go
git commit -m "feat(design): add LogoGradient to Theme and layout constants

- LogoGradient [5]string em Theme com valores para TokyoNight e Cyberpunk
- Constantes MinWidth=80, MinHeight=24, PanelTreeRatio=0.35

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

---

## Task 2: Criar symbols.go com o inventário de símbolos Unicode

**Files:**
- Create: `internal/tui/design/symbols.go`
- Create: `internal/tui/design/symbols_test.go`

- [ ] **Step 1: Criar symbols_test.go com os testes**

Criar `internal/tui/design/symbols_test.go`:

```go
package design

import "testing"

func TestSymbols_TreeNavigation(t *testing.T) {
	tests := []struct{ name, got, want string }{
		{"FolderCollapsed", SymFolderCollapsed, "▶"},
		{"FolderExpanded", SymFolderExpanded, "▼"},
		{"FolderEmpty", SymFolderEmpty, "▷"},
		{"Leaf", SymLeaf, "●"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestSymbols_ItemStates(t *testing.T) {
	tests := []struct{ name, got, want string }{
		{"Favorite", SymFavorite, "★"},
		{"Deleted", SymDeleted, "✗"},
		{"Created", SymCreated, "✦"},
		{"Modified", SymModified, "✎"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestSymbols_Semantic(t *testing.T) {
	tests := []struct{ name, got, want string }{
		{"Success", SymSuccess, "✓"},
		{"Info", SymInfo, "ℹ"},
		{"Warning", SymWarning, "⚠"},
		{"Error", SymError, "✕"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestSymbols_UI(t *testing.T) {
	tests := []struct{ name, got, want string }{
		{"Revealable", SymRevealable, "◉"},
		{"Mask", SymMask, "•"},
		{"Cursor", SymCursor, "▌"},
		{"ScrollUp", SymScrollUp, "↑"},
		{"ScrollDown", SymScrollDown, "↓"},
		{"ScrollThumb", SymScrollThumb, "■"},
		{"Ellipsis", SymEllipsis, "…"},
		{"Bullet", SymBullet, "•"},
		{"HeaderSep", SymHeaderSep, "·"},
		{"TreeConnector", SymTreeConnector, "<╡"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestSymbols_BoxDrawing(t *testing.T) {
	tests := []struct{ name, got, want string }{
		{"BorderH", SymBorderH, "─"},
		{"BorderV", SymBorderV, "│"},
		{"CornerTL", SymCornerTL, "╭"},
		{"CornerTR", SymCornerTR, "╮"},
		{"CornerBL", SymCornerBL, "╰"},
		{"CornerBR", SymCornerBR, "╯"},
		{"JunctionL", SymJunctionL, "├"},
		{"JunctionT", SymJunctionT, "┬"},
		{"JunctionB", SymJunctionB, "┴"},
		{"JunctionR", SymJunctionR, "┤"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestSpinnerFrames_Values(t *testing.T) {
	want := [4]string{"◐", "◓", "◑", "◒"}
	if SpinnerFrames != want {
		t.Errorf("SpinnerFrames = %v, want %v", SpinnerFrames, want)
	}
}

func TestSpinnerFrame_Wraps(t *testing.T) {
	// frame%4 deve iterar pelos 4 frames e reiniciar
	for i := 0; i < 8; i++ {
		got := SpinnerFrame(i)
		want := SpinnerFrames[i%4]
		if got != want {
			t.Errorf("SpinnerFrame(%d) = %q, want %q", i, got, want)
		}
	}
}
```

- [ ] **Step 2: Rodar para confirmar falha de compilação**

```
go test ./internal/tui/design/
```

Esperado: erro de compilação com `undefined: SymFolderCollapsed`, etc.

- [ ] **Step 3: Criar symbols.go**

Criar `internal/tui/design/symbols.go`:

```go
package design

// Símbolos de navegação em árvore.
const (
	// SymFolderCollapsed é o símbolo de pasta recolhida.
	SymFolderCollapsed = "▶" // U+25B6
	// SymFolderExpanded é o símbolo de pasta expandida.
	SymFolderExpanded = "▼" // U+25BC
	// SymFolderEmpty é o símbolo de pasta vazia.
	SymFolderEmpty = "▷" // U+25B7
	// SymLeaf é o símbolo de item folha (segredo).
	SymLeaf = "●" // U+25CF
)

// Símbolos de estado de item.
const (
	// SymFavorite é o símbolo de item favorito.
	SymFavorite = "★" // U+2605
	// SymDeleted é o símbolo de item marcado para exclusão.
	SymDeleted = "✗" // U+2717
	// SymCreated é o símbolo de item recém-criado, ainda não salvo.
	SymCreated = "✦" // U+2726
	// SymModified é o símbolo de item modificado, ainda não salvo.
	SymModified = "✎" // U+270E
)

// Símbolos semânticos — usados na barra de mensagens e em badges de estado.
const (
	// SymSuccess indica operação concluída com êxito.
	SymSuccess = "✓" // U+2713
	// SymInfo indica informação contextual.
	SymInfo = "ℹ" // U+2139
	// SymWarning indica alerta ou aviso.
	SymWarning = "⚠" // U+26A0
	// SymError indica falha ou estado inválido. Distinto de SymDeleted (✗).
	SymError = "✕" // U+2715
)

// Símbolos de elementos de UI.
const (
	// SymRevealable é o indicador de campo com conteúdo revelável.
	SymRevealable = "◉" // U+25C9
	// SymMask é o caractere de máscara para conteúdo sensível.
	// Mesmo glifo que SymBullet (U+2022); papéis semânticos distintos.
	SymMask = "•" // U+2022
	// SymCursor é o cursor em campos de entrada de texto.
	SymCursor = "▌" // U+258C
	// SymScrollUp é o indicador de scroll disponível para cima.
	SymScrollUp = "↑" // U+2191
	// SymScrollDown é o indicador de scroll disponível para baixo.
	SymScrollDown = "↓" // U+2193
	// SymScrollThumb é o thumb que indica a posição atual do scroll.
	SymScrollThumb = "■" // U+25A0
	// SymEllipsis é o caractere de truncamento de texto.
	SymEllipsis = "…" // U+2026
	// SymBullet é o indicador contextual — dirty marker no cabeçalho e marcador de dica.
	// Mesmo glifo que SymMask (U+2022); papéis semânticos distintos.
	SymBullet = "•" // U+2022
	// SymHeaderSep é o separador usado no cabeçalho entre elementos.
	SymHeaderSep = "·" // U+00B7
	// SymTreeConnector é o conector entre a linha selecionada na árvore e o painel de detalhe.
	// Ocupa 2 colunas (Basic Latin '<' + U+2561 '╡').
	SymTreeConnector = "<╡"
)

// Símbolos de estrutura — caracteres Box Drawing usados em bordas e separadores.
const (
	// SymBorderH é o separador horizontal.
	SymBorderH = "─" // U+2500
	// SymBorderV é o separador vertical.
	SymBorderV = "│" // U+2502
	// SymCornerTL é o canto superior esquerdo arredondado.
	SymCornerTL = "╭" // U+256D
	// SymCornerTR é o canto superior direito arredondado.
	SymCornerTR = "╮" // U+256E
	// SymCornerBL é o canto inferior esquerdo arredondado.
	SymCornerBL = "╰" // U+2570
	// SymCornerBR é o canto inferior direito arredondado.
	SymCornerBR = "╯" // U+256F
	// SymJunctionL é o ponto de junção em T voltado para a esquerda.
	SymJunctionL = "├" // U+251C
	// SymJunctionT é o ponto de junção em T voltado para cima.
	SymJunctionT = "┬" // U+252C
	// SymJunctionB é o ponto de junção em T voltado para baixo.
	SymJunctionB = "┴" // U+2534
	// SymJunctionR é o ponto de junção em T voltado para a direita.
	SymJunctionR = "┤" // U+2524
)

// SpinnerFrames são os 4 frames da animação de atividade, em ordem de exibição.
// Nota: caracteres Geometric Shapes podem renderizar em 2 colunas em alguns locales.
// Reserve 2 colunas no layout adjacente para evitar jitter.
var SpinnerFrames = [4]string{"◐", "◓", "◑", "◒"}

// SpinnerFrame retorna o frame do spinner para o índice de animação fornecido.
// O índice é reduzido com frame%4, portanto aceita qualquer valor inteiro não-negativo.
func SpinnerFrame(frame int) string {
	return SpinnerFrames[frame%4]
}
```

- [ ] **Step 4: Rodar testes e confirmar passagem**

```
go test ./internal/tui/design/
```

Esperado: `ok  github.com/useful-toys/abditum/internal/tui/design`

- [ ] **Step 5: Commit**

```
git add internal/tui/design/symbols.go internal/tui/design/symbols_test.go
git commit -m "feat(design): add symbols inventory

Inventario completo de simbolos Unicode como constantes nomeadas.
Inclui SpinnerFrames e SpinnerFrame() para animacao de atividade.

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

---

## Task 3: Criar keys.go com o tipo Key e teclas pré-definidas

**Files:**
- Create: `internal/tui/design/keys.go`
- Create: `internal/tui/design/keys_test.go`

- [ ] **Step 1: Criar keys_test.go com os testes**

Criar `internal/tui/design/keys_test.go`:

```go
package design

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestKeys_Labels(t *testing.T) {
	tests := []struct{ name, got, want string }{
		{"Enter", Keys.Enter.Label, "Enter"},
		{"Esc", Keys.Esc.Label, "Esc"},
		{"Tab", Keys.Tab.Label, "Tab"},
		{"Del", Keys.Del.Label, "Del"},
		{"Ins", Keys.Ins.Label, "Ins"},
		{"Home", Keys.Home.Label, "Home"},
		{"End", Keys.End.Label, "End"},
		{"PgUp", Keys.PgUp.Label, "PgUp"},
		{"PgDn", Keys.PgDn.Label, "PgDn"},
		{"F1", Keys.F1.Label, "F1"},
		{"F6", Keys.F6.Label, "F6"},
		{"F12", Keys.F12.Label, "F12"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("Keys.%s.Label = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestKeys_Matches_SimpleKeys(t *testing.T) {
	tests := []struct {
		name string
		key  Key
		msg  tea.KeyPressMsg
	}{
		{"Enter", Keys.Enter, tea.KeyPressMsg{Code: tea.KeyEnter}},
		{"Esc", Keys.Esc, tea.KeyPressMsg{Code: tea.KeyEscape}},
		{"Tab", Keys.Tab, tea.KeyPressMsg{Code: tea.KeyTab}},
		{"Del", Keys.Del, tea.KeyPressMsg{Code: tea.KeyDelete}},
		{"Ins", Keys.Ins, tea.KeyPressMsg{Code: tea.KeyInsert}},
		{"Home", Keys.Home, tea.KeyPressMsg{Code: tea.KeyHome}},
		{"End", Keys.End, tea.KeyPressMsg{Code: tea.KeyEnd}},
		{"PgUp", Keys.PgUp, tea.KeyPressMsg{Code: tea.KeyPgUp}},
		{"PgDn", Keys.PgDn, tea.KeyPressMsg{Code: tea.KeyPgDown}},
		{"F1", Keys.F1, tea.KeyPressMsg{Code: tea.KeyF1}},
		{"F6", Keys.F6, tea.KeyPressMsg{Code: tea.KeyF6}},
		{"F12", Keys.F12, tea.KeyPressMsg{Code: tea.KeyF12}},
	}
	for _, tt := range tests {
		if !tt.key.Matches(tt.msg) {
			t.Errorf("Keys.%s.Matches(correct msg) = false, want true", tt.name)
		}
	}
}

func TestKeys_Matches_DoesNotMatchWrong(t *testing.T) {
	// Enter não deve casar com Esc
	if Keys.Enter.Matches(tea.KeyPressMsg{Code: tea.KeyEscape}) {
		t.Error("Keys.Enter.Matches(Esc) = true, want false")
	}
	// F1 sem modificador não deve casar com Ctrl+F1
	if Keys.F1.Matches(tea.KeyPressMsg{Code: tea.KeyF1, Mod: tea.ModCtrl}) {
		t.Error("Keys.F1.Matches(ctrl+F1) = true, want false")
	}
}

func TestLetter(t *testing.T) {
	k := Letter('q')
	if k.Label != "Q" {
		t.Errorf("Letter('q').Label = %q, want \"Q\"", k.Label)
	}
	if k.Code != 'q' {
		t.Errorf("Letter('q').Code = %v, want 'q'", k.Code)
	}
	if k.Mod != 0 {
		t.Errorf("Letter('q').Mod = %v, want 0", k.Mod)
	}
	if !k.Matches(tea.KeyPressMsg{Code: 'q'}) {
		t.Error("Letter('q').Matches({Code:'q'}) = false, want true")
	}
}

func TestWithCtrl(t *testing.T) {
	k := WithCtrl(Letter('q'))
	if k.Label != "⌃Q" {
		t.Errorf("WithCtrl(Letter('q')).Label = %q, want \"⌃Q\"", k.Label)
	}
	if k.Code != 'q' {
		t.Errorf("WithCtrl(Letter('q')).Code = %v, want 'q'", k.Code)
	}
	if !k.Mod.Contains(tea.ModCtrl) {
		t.Error("WithCtrl(Letter('q')).Mod deve conter ModCtrl")
	}
	if !k.Matches(tea.KeyPressMsg{Code: 'q', Mod: tea.ModCtrl}) {
		t.Error("WithCtrl(Letter('q')).Matches(ctrl+q) = false, want true")
	}
	// sem modificador não deve casar
	if k.Matches(tea.KeyPressMsg{Code: 'q'}) {
		t.Error("WithCtrl(Letter('q')).Matches(q) = true, want false")
	}
}

func TestWithShift(t *testing.T) {
	k := WithShift(Keys.F6)
	if k.Label != "⇧F6" {
		t.Errorf("WithShift(Keys.F6).Label = %q, want \"⇧F6\"", k.Label)
	}
	if !k.Matches(tea.KeyPressMsg{Code: tea.KeyF6, Mod: tea.ModShift}) {
		t.Error("WithShift(Keys.F6).Matches(shift+F6) = false, want true")
	}
	// sem modificador não deve casar
	if k.Matches(tea.KeyPressMsg{Code: tea.KeyF6}) {
		t.Error("WithShift(Keys.F6).Matches(F6) = true, want false")
	}
}

func TestWithAlt(t *testing.T) {
	k := WithAlt(Letter('x'))
	if k.Label != "!X" {
		t.Errorf("WithAlt(Letter('x')).Label = %q, want \"!X\"", k.Label)
	}
	if !k.Matches(tea.KeyPressMsg{Code: 'x', Mod: tea.ModAlt}) {
		t.Error("WithAlt(Letter('x')).Matches(alt+x) = false, want true")
	}
}

func TestComposition_CtrlAltShift(t *testing.T) {
	// Ordem de composição: WithCtrl(WithAlt(WithShift(base))) produz "⌃!⇧Q"
	k := WithCtrl(WithAlt(WithShift(Letter('q'))))
	if k.Label != "⌃!⇧Q" {
		t.Errorf("label = %q, want \"⌃!⇧Q\"", k.Label)
	}
	msg := tea.KeyPressMsg{Code: 'q', Mod: tea.ModCtrl | tea.ModAlt | tea.ModShift}
	if !k.Matches(msg) {
		t.Error("Matches(ctrl+alt+shift+q) = false, want true")
	}
}

func TestShortcuts_Labels(t *testing.T) {
	tests := []struct{ name, got, want string }{
		{"Help", Shortcuts.Help.Label, "F1"},
		{"ThemeToggle", Shortcuts.ThemeToggle.Label, "F12"},
		{"Quit", Shortcuts.Quit.Label, "⌃Q"},
		{"LockVault", Shortcuts.LockVault.Label, "⌃!⇧Q"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("Shortcuts.%s.Label = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestShortcuts_Matches(t *testing.T) {
	if !Shortcuts.Help.Matches(tea.KeyPressMsg{Code: tea.KeyF1}) {
		t.Error("Shortcuts.Help.Matches(F1) = false, want true")
	}
	if !Shortcuts.ThemeToggle.Matches(tea.KeyPressMsg{Code: tea.KeyF12}) {
		t.Error("Shortcuts.ThemeToggle.Matches(F12) = false, want true")
	}
	if !Shortcuts.Quit.Matches(tea.KeyPressMsg{Code: 'q', Mod: tea.ModCtrl}) {
		t.Error("Shortcuts.Quit.Matches(ctrl+q) = false, want true")
	}
	emergency := tea.KeyPressMsg{Code: 'q', Mod: tea.ModCtrl | tea.ModAlt | tea.ModShift}
	if !Shortcuts.LockVault.Matches(emergency) {
		t.Error("Shortcuts.LockVault.Matches(ctrl+alt+shift+q) = false, want true")
	}
}
```

- [ ] **Step 2: Rodar para confirmar falha de compilação**

```
go test ./internal/tui/design/
```

Esperado: erro de compilação com `undefined: Keys`, `undefined: Letter`, etc.

- [ ] **Step 3: Criar keys.go**

Criar `internal/tui/design/keys.go`:

```go
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
	Enter, Esc, Tab, Del, Ins, Home, End, PgUp, PgDn    Key
	F1, F2, F3, F4, F5, F6, F7, F8, F9, F10, F11, F12 Key
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

// Shortcuts contém os 4 atalhos globais do design system, ativos em qualquer contexto da aplicação.
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
}{
	Help:        Keys.F1,
	ThemeToggle: Keys.F12,
	Quit:        WithCtrl(Letter('q')),
	LockVault:   WithCtrl(WithAlt(WithShift(Letter('q')))),
}
```

- [ ] **Step 4: Rodar testes e confirmar passagem**

```
go test ./internal/tui/design/
```

Esperado: `ok  github.com/useful-toys/abditum/internal/tui/design`

- [ ] **Step 5: Rodar a suite completa**

```
go test ./...
```

Esperado: todos os pacotes passam sem erros.

- [ ] **Step 6: Commit**

```
git add internal/tui/design/keys.go internal/tui/design/keys_test.go
git commit -m "feat(design): add Key type with typed key matching and shortcuts

- type Key com Label, Code, Mod e metodo Matches (sem string matching)
- ModLabelCtrl/Shift/Alt para montar rotulos de atalhos
- var Keys com todas as teclas simples (Enter, Esc, F1-F12, etc.)
- Funcoes Letter, WithCtrl, WithShift, WithAlt para composicao
- var Shortcuts com os 4 atalhos globais (F1, F12, ctrl+q, ctrl+alt+shift+q)

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```
