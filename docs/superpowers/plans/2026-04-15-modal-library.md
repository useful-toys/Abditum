# Modal Library Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the stub modals with a fully-featured, design-system-compliant modal library that provides `DialogFrame`, `ScrollState`, and `KeyHandler` as reusable foundations for all future dialog types.

**Architecture:** `DialogFrame` is a stateless renderer (pure function wrapped in a config struct); `ScrollState` is mutable state owned by the caller modal; `KeyHandler` is a composable dispatch helper. `design/design_modal.go` centralises all colour/symbol helpers following the same pattern as `design_action.go`. The `ConfirmModal` and `HelpModal` rewrites demonstrate the foundations.

**Tech Stack:** Go, `charm.land/bubbletea/v2`, `charm.land/lipgloss/v2`, `github.com/useful-toys/abditum/internal/tui/design`, `github.com/useful-toys/abditum/internal/tui/actions`.

---

## File Map

| Status | File | Responsibility |
|--------|------|----------------|
| Modify | `internal/tui/design/keys.go` | Add `Up` and `Down` arrow keys to `Keys` struct |
| Modify | `internal/tui/design/keys_test.go` | Tests for `Keys.Up` and `Keys.Down` |
| Modify | `internal/tui/actions/action.go` | Add `Order int` field to `ActionGroup` |
| Modify | `internal/tui/modal/modal_base.go` | Change `ModalOption.Keys` from `[]string` to `[]design.Key` |
| Create | `internal/tui/design/design_modal.go` | `Severity` type + rendering helpers for dialogs |
| Create | `internal/tui/design/design_modal_test.go` | Unit tests for all helpers in `design_modal.go` |
| Create | `internal/tui/modal/scroll_state.go` | `ScrollState` struct: Up/Down/PageUp/PageDown/Home/End/ThumbLine |
| Create | `internal/tui/modal/scroll_state_test.go` | Unit tests for `ScrollState` including `ThumbLine` priority rules |
| Create | `internal/tui/modal/key_handler.go` | `KeyHandler` struct: dispatches option keys + scroll navigation |
| Create | `internal/tui/modal/key_handler_test.go` | Unit tests for `KeyHandler.Handle` |
| Create | `internal/tui/modal/frame.go` | `DialogFrame.Render` — stateless dialog renderer |
| Create | `internal/tui/modal/frame_test.go` | Golden file tests for all frame variants |
| Rewrite | `internal/tui/modal/confirm_modal.go` | `ConfirmModal` using `DialogFrame` + `KeyHandler` |
| Create | `internal/tui/modal/confirm_modal_test.go` | Golden file test for ConfirmModal |
| Rewrite | `internal/tui/modal/help_modal.go` | `HelpModal` using `DialogFrame` + `ScrollState` + `KeyHandler` |
| Create | `internal/tui/modal/help_modal_test.go` | Golden file tests for HelpModal (no-scroll + with-scroll) |
| Modify | `cmd/abditum/setup.go` | Update `NewHelpModal` caller (signature unchanged, but HelpModal.Update now delegates to HandleKey) |

---

## Task 1: Add arrow keys (`Up` / `Down`) to `design.Keys`

`KeyHandler` needs `design.Keys.Up` and `design.Keys.Down` for scroll navigation. They don't exist yet.

**Files:**
- Modify: `internal/tui/design/keys.go`
- Modify: `internal/tui/design/keys_test.go`

- [ ] **Step 1.1: Write the failing test**

In `internal/tui/design/keys_test.go`, add inside `TestKeys_Labels`:

```go
{"Up",   Keys.Up.Label,   "↑"},
{"Down", Keys.Down.Label, "↓"},
```

And inside `TestKeys_Matches_SimpleKeys`:

```go
{"Up",   Keys.Up,   tea.KeyPressMsg{Code: tea.KeyUp}},
{"Down", Keys.Down, tea.KeyPressMsg{Code: tea.KeyDown}},
```

- [ ] **Step 1.2: Run test to verify it fails**

```
cd internal/tui/design && go test -run "TestKeys_Labels|TestKeys_Matches_SimpleKeys" -v
```

Expected: FAIL — `Keys.Up` and `Keys.Down` are undefined.

- [ ] **Step 1.3: Add `Up` and `Down` to `keys.go`**

Change the `Keys` struct declaration from:

```go
var Keys = struct {
	Enter, Esc, Tab, Del, Ins, Home, End, PgUp, PgDn    Key
	F1, F2, F3, F4, F5, F6, F7, F8, F9, F10, F11, F12 Key
}{
```

to:

```go
var Keys = struct {
	Enter, Esc, Tab, Del, Ins, Home, End, PgUp, PgDn    Key
	Up, Down                                             Key
	F1, F2, F3, F4, F5, F6, F7, F8, F9, F10, F11, F12 Key
}{
```

And add after `PgDn`:

```go
	Up:    Key{Label: "↑", Code: tea.KeyUp},
	Down:  Key{Label: "↓", Code: tea.KeyDown},
```

- [ ] **Step 1.4: Run test to verify it passes**

```
cd internal/tui/design && go test -run "TestKeys_Labels|TestKeys_Matches_SimpleKeys" -v
```

Expected: PASS

- [ ] **Step 1.5: Run the full test suite to check for regressions**

```
go test ./...
```

Expected: PASS (no regressions — only additions).

- [ ] **Step 1.6: Commit**

```bash
git add internal/tui/design/keys.go internal/tui/design/keys_test.go
git commit -m "feat(design): add Up and Down arrow keys to design.Keys"
```

---

## Task 2: Add `Order int` to `ActionGroup`

**Files:**
- Modify: `internal/tui/actions/action.go`

- [ ] **Step 2.1: Verify existing tests pass (baseline)**

```
go test ./internal/tui/actions/...
```

Expected: PASS

- [ ] **Step 2.2: Add the `Order` field to `ActionGroup`**

In `internal/tui/actions/action.go`, change:

```go
type ActionGroup struct {
	ID          string // identificador único do grupo
	Label       string // cabeçalho exibido no modal de ajuda
	Description string // texto descritivo do grupo
}
```

to:

```go
type ActionGroup struct {
	ID          string // identificador único do grupo
	Label       string // cabeçalho exibido no modal de ajuda
	Description string // texto descritivo do grupo
	Order       int    // ordem de exibição no modal de ajuda; menor valor aparece primeiro
}
```

- [ ] **Step 2.3: Verify nothing broke**

```
go test ./...
```

Expected: PASS — callers that don't set `Order` get zero-value `0`, behaviour unchanged.

- [ ] **Step 2.4: Commit**

```bash
git add internal/tui/actions/action.go
git commit -m "feat(actions): add Order field to ActionGroup for help modal ordering"
```

---

## Task 3: Change `ModalOption.Keys` from `[]string` to `[]design.Key`

**Files:**
- Modify: `internal/tui/modal/modal_base.go`

> **Note:** `confirm_modal.go` currently uses `[]string{"Enter"}` etc. inside `NewConfirmModal`. Those will be removed when `confirm_modal.go` is rewritten in Task 8. For now, the existing `confirm_modal.go` will fail to compile after this change — that is expected and intentional. Fix it in the same commit by updating the old `HandleKey` loop so the project builds.

- [ ] **Step 3.1: Update `modal_base.go`**

Replace the entire file content with:

```go
package modal

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

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
// Cada opção é ativada por uma ou mais teclas tipadas.
type ModalOption struct {
	// Keys lista as teclas que ativam esta opção.
	// Keys[0].Label é exibido no rodapé do diálogo.
	// Demais Keys são aliases funcionais (ex: Enter como alias de "S Sobrescrever").
	Keys []design.Key
	// Label é o texto exibido ao usuário para descrever a ação.
	Label string
	// Intent classifica a intenção semântica desta ação.
	Intent Intent
	// Action é a função executada quando a opção é escolhida pelo usuário.
	Action func() tea.Cmd
}
```

- [ ] **Step 3.2: Fix `confirm_modal.go` so the project compiles**

`confirm_modal.go` currently uses `[]string{"Enter"}` and `msg.String() == key`. Replace those usages with `[]design.Key{design.Keys.Enter}` and `k.Matches(msg)`. The full rewrite happens in Task 8, but a minimal fix is needed now.

In `internal/tui/modal/confirm_modal.go`, change the import block and `NewConfirmModal` + `HandleKey` so it compiles. Replace the entire file with this transitional version (to be superseded in Task 8):

```go
package modal

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// ConfirmModal exibe um diálogo de confirmação com título, mensagem e ações.
// Implementa tui.ModalView. Será reescrito com DialogFrame em Task 8.
type ConfirmModal struct {
	title   string
	message string
	options []ModalOption
}

// NewConfirmModal cria um ConfirmModal com ações padrão (Enter=confirmar, Esc=cancelar).
// Aceita opts opcionais; se nil, usa as ações padrão.
func NewConfirmModal(title, message string, opts ...[]ModalOption) *ConfirmModal {
	defaultOpts := []ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "Confirmar",
			Intent: IntentConfirm,
			Action: func() tea.Cmd {
				return func() tea.Msg { return tui.ModalReadyMsg{} }
			},
		},
		{
			Keys:   []design.Key{design.Keys.Esc},
			Label:  "Cancelar",
			Intent: IntentCancel,
			Action: func() tea.Cmd {
				return func() tea.Msg { return tui.CloseModalMsg{} }
			},
		},
	}
	options := defaultOpts
	if len(opts) > 0 && opts[0] != nil {
		options = opts[0]
	}
	return &ConfirmModal{title: title, message: message, options: options}
}

// Render retorna o modal estilizado — stub; reescrito em Task 8.
func (m *ConfirmModal) Render(maxHeight, maxWidth int, theme *design.Theme) string {
	content := m.title + "\n\n" + m.message
	style := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(lipgloss.Color(theme.Border.Default)).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Raised))
	return style.Render(content)
}

// HandleKey verifica se a tecla corresponde a alguma opção e executa sua ação.
func (m *ConfirmModal) HandleKey(msg tea.KeyMsg) tea.Cmd {
	for _, opt := range m.options {
		for _, k := range opt.Keys {
			if k.Matches(msg) {
				return opt.Action()
			}
		}
	}
	return nil
}

// Update processa mensagens do Bubble Tea, delegando teclas para HandleKey.
func (m *ConfirmModal) Update(msg tea.Msg) tea.Cmd {
	if key, ok := msg.(tea.KeyMsg); ok {
		return m.HandleKey(key)
	}
	return nil
}
```

- [ ] **Step 3.3: Verify the project builds**

```
go build ./...
```

Expected: SUCCESS

- [ ] **Step 3.4: Run all tests**

```
go test ./...
```

Expected: PASS

- [ ] **Step 3.5: Commit**

```bash
git add internal/tui/modal/modal_base.go internal/tui/modal/confirm_modal.go
git commit -m "refactor(modal): change ModalOption.Keys to []design.Key (typed, pattern-consistent)"
```

---

## Task 4: Create `design/design_modal.go` — Severity type + render helpers

**Files:**
- Create: `internal/tui/design/design_modal.go`
- Create: `internal/tui/design/design_modal_test.go`

- [ ] **Step 4.1: Write the failing tests first**

Create `internal/tui/design/design_modal_test.go`:

```go
package design

import (
	"testing"

	"charm.land/lipgloss/v2"
)

// --- Severity ---

func TestSeverity_Symbol(t *testing.T) {
	tests := []struct {
		sev  Severity
		want string
	}{
		{SeverityNeutral, ""},
		{SeverityInformative, SymInfo},
		{SeverityAlert, SymWarning},
		{SeverityDestructive, SymWarning},
		{SeverityError, SymError},
	}
	for _, tt := range tests {
		if got := tt.sev.Symbol(); got != tt.want {
			t.Errorf("Severity(%d).Symbol() = %q, want %q", tt.sev, got, tt.want)
		}
	}
}

func TestSeverity_BorderColor_NotEmpty(t *testing.T) {
	theme := TokyoNight
	sevs := []Severity{SeverityNeutral, SeverityInformative, SeverityAlert, SeverityDestructive, SeverityError}
	for _, s := range sevs {
		if c := s.BorderColor(theme); c == "" {
			t.Errorf("Severity(%d).BorderColor(TokyoNight) = empty string", s)
		}
	}
}

func TestSeverity_DefaultKeyColor_NotEmpty(t *testing.T) {
	theme := TokyoNight
	sevs := []Severity{SeverityNeutral, SeverityInformative, SeverityAlert, SeverityDestructive, SeverityError}
	for _, s := range sevs {
		if c := s.DefaultKeyColor(theme); c == "" {
			t.Errorf("Severity(%d).DefaultKeyColor(TokyoNight) = empty string", s)
		}
	}
}

func TestSeverityDestructive_DefaultKeyColor_IsSemanticError(t *testing.T) {
	theme := TokyoNight
	got := SeverityDestructive.DefaultKeyColor(theme)
	want := theme.Semantic.Error
	if got != want {
		t.Errorf("SeverityDestructive.DefaultKeyColor = %q, want Semantic.Error = %q", got, want)
	}
}

// --- RenderDialogTitle ---

func TestRenderDialogTitle_WithoutSymbol(t *testing.T) {
	theme := TokyoNight
	text, width := RenderDialogTitle("Título", "", "", theme)
	if width <= 0 {
		t.Errorf("RenderDialogTitle width = %d, want > 0", width)
	}
	if lipgloss.Width(text) != width {
		t.Errorf("RenderDialogTitle: returned width %d != lipgloss.Width %d", width, lipgloss.Width(text))
	}
	clean := ansiEscapeRe.ReplaceAllString(text, "")
	if !containsSubstring(clean, "Título") {
		t.Errorf("RenderDialogTitle: clean text %q does not contain title", clean)
	}
}

func TestRenderDialogTitle_WithSymbol(t *testing.T) {
	theme := TokyoNight
	text, width := RenderDialogTitle("Aviso", SymWarning, theme.Semantic.Warning, theme)
	clean := ansiEscapeRe.ReplaceAllString(text, "")
	if !containsSubstring(clean, SymWarning) {
		t.Errorf("RenderDialogTitle with symbol: clean text %q does not contain symbol", clean)
	}
	if !containsSubstring(clean, "Aviso") {
		t.Errorf("RenderDialogTitle with symbol: clean text %q does not contain title", clean)
	}
	if lipgloss.Width(text) != width {
		t.Errorf("RenderDialogTitle: returned width %d != lipgloss.Width %d", width, lipgloss.Width(text))
	}
}

// --- RenderDialogAction ---

func TestRenderDialogAction_ContainsKeyAndLabel(t *testing.T) {
	theme := TokyoNight
	text, width := RenderDialogAction("Enter", "Confirmar", theme.Accent.Primary, theme)
	clean := ansiEscapeRe.ReplaceAllString(text, "")
	if !containsSubstring(clean, "Enter") {
		t.Errorf("RenderDialogAction: clean %q does not contain key", clean)
	}
	if !containsSubstring(clean, "Confirmar") {
		t.Errorf("RenderDialogAction: clean %q does not contain label", clean)
	}
	if lipgloss.Width(text) != width {
		t.Errorf("RenderDialogAction: returned width %d != lipgloss.Width %d", width, lipgloss.Width(text))
	}
}

// --- RenderScrollArrow ---

func TestRenderScrollArrow_Width(t *testing.T) {
	theme := TokyoNight
	upText, upW := RenderScrollArrow(true, theme)
	if upW != 1 {
		t.Errorf("RenderScrollArrow(up) width = %d, want 1", upW)
	}
	if lipgloss.Width(upText) != upW {
		t.Errorf("RenderScrollArrow(up) lipgloss.Width mismatch")
	}
	downText, downW := RenderScrollArrow(false, theme)
	if downW != 1 {
		t.Errorf("RenderScrollArrow(down) width = %d, want 1", downW)
	}
	if lipgloss.Width(downText) != downW {
		t.Errorf("RenderScrollArrow(down) lipgloss.Width mismatch")
	}
}

func TestRenderScrollArrow_CorrectSymbol(t *testing.T) {
	theme := TokyoNight
	upText, _ := RenderScrollArrow(true, theme)
	downText, _ := RenderScrollArrow(false, theme)
	cleanUp := ansiEscapeRe.ReplaceAllString(upText, "")
	cleanDown := ansiEscapeRe.ReplaceAllString(downText, "")
	if cleanUp != SymScrollUp {
		t.Errorf("RenderScrollArrow(true) = %q, want %q", cleanUp, SymScrollUp)
	}
	if cleanDown != SymScrollDown {
		t.Errorf("RenderScrollArrow(false) = %q, want %q", cleanDown, SymScrollDown)
	}
}

// --- RenderScrollThumb ---

func TestRenderScrollThumb_Width(t *testing.T) {
	theme := TokyoNight
	text, w := RenderScrollThumb(theme)
	if w != 1 {
		t.Errorf("RenderScrollThumb width = %d, want 1", w)
	}
	if lipgloss.Width(text) != w {
		t.Errorf("RenderScrollThumb lipgloss.Width mismatch")
	}
	clean := ansiEscapeRe.ReplaceAllString(text, "")
	if clean != SymScrollThumb {
		t.Errorf("RenderScrollThumb = %q, want %q", clean, SymScrollThumb)
	}
}
```

- [ ] **Step 4.2: Run to verify they fail**

```
go test ./internal/tui/design/... -run "TestSeverity|TestRenderDialog|TestRenderScroll" -v
```

Expected: FAIL — `Severity`, `RenderDialogTitle`, `RenderDialogAction`, `RenderScrollArrow`, `RenderScrollThumb` undefined.

- [ ] **Step 4.3: Create `internal/tui/design/design_modal.go`**

```go
package design

import "charm.land/lipgloss/v2"

// DialogPaddingH é o padding horizontal interno dos diálogos (colunas em cada lado).
const DialogPaddingH = 2

// DialogPaddingV é o padding vertical interno dos diálogos (linhas acima e abaixo do corpo).
// O HelpModal usa 0 linhas de padding vertical — não inclui linhas em branco no corpo.
const DialogPaddingV = 1

// Severity representa a severidade visual de um diálogo de Notificação ou Confirmação.
// Diálogos de Ajuda e Funcionais não usam severidade.
type Severity int

const (
	// SeverityNeutral não tem símbolo; usa border.focused; tecla default em accent.primary.
	SeverityNeutral Severity = iota
	// SeverityInformative tem símbolo ℹ e usa semantic.info como cor de borda.
	SeverityInformative
	// SeverityAlert tem símbolo ⚠ e usa semantic.warning como cor de borda.
	SeverityAlert
	// SeverityDestructive tem símbolo ⚠, usa semantic.warning como cor de borda,
	// mas a tecla default usa semantic.error (ação destrutiva irrecuperável).
	SeverityDestructive
	// SeverityError tem símbolo ✕ e usa semantic.error como cor de borda.
	SeverityError
)

// Symbol retorna o símbolo Unicode da severidade.
// Retorna "" para SeverityNeutral.
// SeverityAlert e SeverityDestructive retornam ambos SymWarning — a distinção
// visual está na cor da tecla default (DefaultKeyColor), não no símbolo.
func (s Severity) Symbol() string {
	switch s {
	case SeverityInformative:
		return SymInfo
	case SeverityAlert, SeverityDestructive:
		return SymWarning
	case SeverityError:
		return SymError
	default: // SeverityNeutral
		return ""
	}
}

// BorderColor retorna a cor de borda para a severidade a partir do tema.
func (s Severity) BorderColor(theme *Theme) string {
	switch s {
	case SeverityInformative:
		return theme.Semantic.Info
	case SeverityAlert, SeverityDestructive:
		return theme.Semantic.Warning
	case SeverityError:
		return theme.Semantic.Error
	default: // SeverityNeutral
		return theme.Border.Focused
	}
}

// DefaultKeyColor retorna a cor da tecla da ação default (primeira opção) para a severidade.
// Para SeverityDestructive, a tecla default usa semantic.error para enfatizar o risco.
// Todas as demais severidades usam accent.primary.
func (s Severity) DefaultKeyColor(theme *Theme) string {
	if s == SeverityDestructive {
		return theme.Semantic.Error
	}
	return theme.Accent.Primary
}

// RenderDialogTitle renderiza o bloco título do diálogo.
// Se symbol != "", inclui "symbol  title" (símbolo + 2 espaços + título).
// Se symbol == "", inclui apenas "title".
// Cores: symbol em symbolColor, título em theme.Text.Primary + bold.
// Retorna texto ANSI e largura visual em colunas.
func RenderDialogTitle(title, symbol, symbolColor string, theme *Theme) (string, int) {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Bold(true)

	var rendered string
	if symbol != "" {
		symStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(symbolColor))
		rendered = symStyle.Render(symbol) + "  " + titleStyle.Render(title)
	} else {
		rendered = titleStyle.Render(title)
	}
	return rendered, lipgloss.Width(rendered)
}

// RenderDialogAction renderiza uma ação do rodapé: "key label".
// key é o Label da tecla (Keys[0].Label da ModalOption — ex: "Enter", "S", "Esc").
// key é renderizada em keyColor; label é renderizada em theme.Text.Primary.
// Retorna texto ANSI e largura visual em colunas.
func RenderDialogAction(key, label, keyColor string, theme *Theme) (string, int) {
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(keyColor)).
		Bold(true)
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Text.Primary))

	rendered := keyStyle.Render(key) + " " + labelStyle.Render(label)
	return rendered, lipgloss.Width(rendered)
}

// RenderScrollArrow renderiza ↑ (up=true) ou ↓ (up=false) em theme.Text.Secondary.
// Sempre retorna width == 1.
func RenderScrollArrow(up bool, theme *Theme) (string, int) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Secondary))
	sym := SymScrollDown
	if up {
		sym = SymScrollUp
	}
	return style.Render(sym), 1
}

// RenderScrollThumb renderiza ■ em theme.Text.Secondary.
// Sempre retorna width == 1.
func RenderScrollThumb(theme *Theme) (string, int) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Secondary))
	return style.Render(SymScrollThumb), 1
}
```

- [ ] **Step 4.4: Run tests to verify they pass**

```
go test ./internal/tui/design/... -run "TestSeverity|TestRenderDialog|TestRenderScroll" -v
```

Expected: PASS

- [ ] **Step 4.5: Run full suite**

```
go test ./...
```

Expected: PASS

- [ ] **Step 4.6: Commit**

```bash
git add internal/tui/design/design_modal.go internal/tui/design/design_modal_test.go
git commit -m "feat(design): add Severity type and dialog render helpers"
```

---

## Task 5: Create `modal/scroll_state.go`

**Files:**
- Create: `internal/tui/modal/scroll_state.go`
- Create: `internal/tui/modal/scroll_state_test.go`

- [ ] **Step 5.1: Write the failing tests**

Create `internal/tui/modal/scroll_state_test.go`:

```go
package modal

import "testing"

func TestScrollState_CanScrollUp_Down(t *testing.T) {
	s := ScrollState{Offset: 0, Total: 20, Viewport: 10}
	if s.CanScrollUp() {
		t.Error("CanScrollUp should be false when Offset=0")
	}
	if !s.CanScrollDown() {
		t.Error("CanScrollDown should be true when Total > Viewport")
	}

	s.Offset = 10
	if !s.CanScrollUp() {
		t.Error("CanScrollUp should be true when Offset > 0")
	}
	if !s.CanScrollDown() {
		t.Error("CanScrollDown should be true when offset+viewport < total")
	}

	s.Offset = 10
	s.Total = 20
	s.Viewport = 10
	// offset + viewport == total → cannot scroll down
	if s.CanScrollDown() {
		t.Error("CanScrollDown should be false when offset+viewport == total")
	}
}

func TestScrollState_Up_Down(t *testing.T) {
	s := ScrollState{Offset: 5, Total: 20, Viewport: 10}
	s.Up()
	if s.Offset != 4 {
		t.Errorf("Up: Offset = %d, want 4", s.Offset)
	}
	s.Offset = 0
	s.Up()
	if s.Offset != 0 {
		t.Errorf("Up at 0: Offset = %d, want 0 (clamped)", s.Offset)
	}
	s.Offset = 5
	s.Down()
	if s.Offset != 6 {
		t.Errorf("Down: Offset = %d, want 6", s.Offset)
	}
	s.Offset = 10 // offset + viewport == total → at bottom
	s.Down()
	if s.Offset != 10 {
		t.Errorf("Down at bottom: Offset = %d, want 10 (clamped)", s.Offset)
	}
}

func TestScrollState_PageUp_PageDown(t *testing.T) {
	s := ScrollState{Offset: 15, Total: 40, Viewport: 10}
	s.PageUp()
	if s.Offset != 5 {
		t.Errorf("PageUp: Offset = %d, want 5", s.Offset)
	}
	s.PageUp()
	if s.Offset != 0 {
		t.Errorf("PageUp clamp: Offset = %d, want 0", s.Offset)
	}

	s.Offset = 15
	s.PageDown()
	if s.Offset != 25 {
		t.Errorf("PageDown: Offset = %d, want 25", s.Offset)
	}
	s.PageDown()
	if s.Offset != 30 {
		t.Errorf("PageDown clamp: Offset = %d, want 30 (total-viewport)", s.Offset)
	}
}

func TestScrollState_Home_End(t *testing.T) {
	s := ScrollState{Offset: 15, Total: 40, Viewport: 10}
	s.Home()
	if s.Offset != 0 {
		t.Errorf("Home: Offset = %d, want 0", s.Offset)
	}
	s.End()
	if s.Offset != 30 {
		t.Errorf("End: Offset = %d, want 30 (total-viewport)", s.Offset)
	}
}

func TestScrollState_ThumbLine_InactiveWhenNoScroll(t *testing.T) {
	// Total <= Viewport → no scroll → ThumbLine returns -1
	s := ScrollState{Offset: 0, Total: 5, Viewport: 10}
	if got := s.ThumbLine(); got != -1 {
		t.Errorf("ThumbLine (no scroll): got %d, want -1", got)
	}
	s.Total = 10
	if got := s.ThumbLine(); got != -1 {
		t.Errorf("ThumbLine (total==viewport): got %d, want -1", got)
	}
}

func TestScrollState_ThumbLine_ArrowsHavePriority(t *testing.T) {
	// At top: ↑ inactive, ↓ active on last line.
	// Thumb must NOT be on line viewport (last line — occupied by ↓).
	s := ScrollState{Offset: 0, Total: 30, Viewport: 10}
	thumb := s.ThumbLine()
	if thumb == s.Viewport {
		t.Errorf("ThumbLine at top: thumb on last line %d (occupied by ↓)", thumb)
	}
	if thumb < 1 || thumb > s.Viewport {
		t.Errorf("ThumbLine at top: got %d, want in [1..%d]", thumb, s.Viewport)
	}

	// At bottom: ↑ active on line 1, ↓ inactive.
	// Thumb must NOT be on line 1.
	s.Offset = 20 // offset + viewport == total
	thumb = s.ThumbLine()
	if thumb == 1 {
		t.Errorf("ThumbLine at bottom: thumb on line 1 (occupied by ↑)")
	}
	if thumb < 1 || thumb > s.Viewport {
		t.Errorf("ThumbLine at bottom: got %d, want in [1..%d]", thumb, s.Viewport)
	}

	// In middle: both ↑ and ↓ active.
	// Thumb must NOT be on line 1 or line viewport.
	s.Offset = 10
	thumb = s.ThumbLine()
	if thumb == 1 {
		t.Errorf("ThumbLine in middle: thumb on line 1 (occupied by ↑)")
	}
	if thumb == s.Viewport {
		t.Errorf("ThumbLine in middle: thumb on last line (occupied by ↓)")
	}
	if thumb < 2 || thumb > s.Viewport-1 {
		t.Errorf("ThumbLine in middle: got %d, want in [2..%d]", thumb, s.Viewport-1)
	}
}

func TestScrollState_ThumbLine_ReturnsMinus1WhenNoSpace(t *testing.T) {
	// Viewport = 2, both arrows active → no room for thumb.
	s := ScrollState{Offset: 5, Total: 20, Viewport: 2}
	if got := s.ThumbLine(); got != -1 {
		t.Errorf("ThumbLine tiny viewport: got %d, want -1", got)
	}
}
```

- [ ] **Step 5.2: Run to verify failures**

```
go test ./internal/tui/modal/... -run "TestScrollState" -v
```

Expected: FAIL — `ScrollState` undefined.

- [ ] **Step 5.3: Create `internal/tui/modal/scroll_state.go`**

```go
package modal

// ScrollState mantém a posição do viewport em conteúdo que pode ser maior que a tela.
// É um estado mutável — pertence ao modal que o utiliza, não ao DialogFrame.
type ScrollState struct {
	// Offset é o índice da primeira linha visível no conteúdo (0-based).
	Offset int
	// Total é o número total de linhas do conteúdo.
	Total int
	// Viewport é o número de linhas visíveis (definido pelo modal em cada Render).
	Viewport int
}

// CanScrollUp retorna true se há conteúdo acima do viewport (Offset > 0).
func (s *ScrollState) CanScrollUp() bool {
	return s.Offset > 0
}

// CanScrollDown retorna true se há conteúdo abaixo do viewport.
func (s *ScrollState) CanScrollDown() bool {
	return s.Offset+s.Viewport < s.Total
}

// Up move o viewport uma linha para cima (sem ultrapassar o início).
func (s *ScrollState) Up() {
	if s.Offset > 0 {
		s.Offset--
	}
}

// Down move o viewport uma linha para baixo (sem ultrapassar o fim).
func (s *ScrollState) Down() {
	maxOffset := s.Total - s.Viewport
	if maxOffset < 0 {
		maxOffset = 0
	}
	if s.Offset < maxOffset {
		s.Offset++
	}
}

// PageUp move o viewport um viewport inteiro para cima (sem ultrapassar o início).
func (s *ScrollState) PageUp() {
	s.Offset -= s.Viewport
	if s.Offset < 0 {
		s.Offset = 0
	}
}

// PageDown move o viewport um viewport inteiro para baixo (sem ultrapassar o fim).
func (s *ScrollState) PageDown() {
	s.Offset += s.Viewport
	maxOffset := s.Total - s.Viewport
	if maxOffset < 0 {
		maxOffset = 0
	}
	if s.Offset > maxOffset {
		s.Offset = maxOffset
	}
}

// Home move o viewport para o início do conteúdo.
func (s *ScrollState) Home() {
	s.Offset = 0
}

// End move o viewport para o fim do conteúdo.
func (s *ScrollState) End() {
	maxOffset := s.Total - s.Viewport
	if maxOffset < 0 {
		maxOffset = 0
	}
	s.Offset = maxOffset
}

// ThumbLine calcula a linha (1-based dentro do viewport) onde o thumb ■ deve aparecer.
//
// Regras:
//   - Retorna -1 se o conteúdo não excede o viewport (scroll inativo).
//   - Setas têm prioridade absoluta:
//     • Se CanScrollUp() == true, a linha 1 do viewport está ocupada por ↑.
//     • Se CanScrollDown() == true, a última linha do viewport está ocupada por ↓.
//   - O thumb é posicionado proporcionalmente nas linhas restantes.
//   - Se o intervalo disponível para o thumb for zero (viewport muito pequeno), retorna -1.
func (s *ScrollState) ThumbLine() int {
	if s.Total <= s.Viewport {
		return -1
	}

	// Determinar linhas reservadas pelas setas.
	firstAvailable := 1
	lastAvailable := s.Viewport
	if s.CanScrollUp() {
		firstAvailable = 2 // linha 1 ocupada pela seta ↑
	}
	if s.CanScrollDown() {
		lastAvailable = s.Viewport - 1 // última linha ocupada pela seta ↓
	}

	available := lastAvailable - firstAvailable + 1
	if available <= 0 {
		return -1
	}

	// Posição proporcional do thumb dentro do intervalo disponível.
	// scrollable é o número máximo de passos de scroll.
	scrollable := s.Total - s.Viewport
	if scrollable == 0 {
		return firstAvailable
	}

	// Mapeia Offset → posição dentro de [0, available-1].
	thumbIndex := (s.Offset * (available - 1)) / scrollable
	return firstAvailable + thumbIndex
}
```

- [ ] **Step 5.4: Run tests to verify they pass**

```
go test ./internal/tui/modal/... -run "TestScrollState" -v
```

Expected: PASS

- [ ] **Step 5.5: Commit**

```bash
git add internal/tui/modal/scroll_state.go internal/tui/modal/scroll_state_test.go
git commit -m "feat(modal): add ScrollState with ThumbLine — arrows always have priority"
```

---

## Task 6: Create `modal/key_handler.go`

**Files:**
- Create: `internal/tui/modal/key_handler.go`
- Create: `internal/tui/modal/key_handler_test.go`

- [ ] **Step 6.1: Write the failing tests**

Create `internal/tui/modal/key_handler_test.go`:

```go
package modal

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

func makeKeyMsg(code rune) tea.KeyMsg {
	return tea.KeyPressMsg{Code: code}
}

func makeSpecialKeyMsg(code tea.Key) tea.KeyMsg {
	return tea.KeyPressMsg{Code: code}
}

func TestKeyHandler_OptionMatch_ExecutesAction(t *testing.T) {
	called := false
	opts := []ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "OK",
			Intent: IntentConfirm,
			Action: func() tea.Cmd {
				called = true
				return nil
			},
		},
	}
	h := KeyHandler{Options: opts}
	cmd, handled := h.Handle(makeSpecialKeyMsg(tea.KeyEnter))
	if !handled {
		t.Error("Handle(Enter): handled = false, want true")
	}
	_ = cmd
	if !called {
		t.Error("Handle(Enter): action was not called")
	}
}

func TestKeyHandler_MultipleKeys_AnyActivatesAction(t *testing.T) {
	callCount := 0
	opts := []ModalOption{
		{
			Keys:  []design.Key{design.Keys.Enter, design.Keys.Esc},
			Label: "OK",
			Action: func() tea.Cmd {
				callCount++
				return nil
			},
		},
	}
	h := KeyHandler{Options: opts}

	h.Handle(makeSpecialKeyMsg(tea.KeyEnter))
	h.Handle(makeSpecialKeyMsg(tea.KeyEscape))
	if callCount != 2 {
		t.Errorf("Handle: callCount = %d, want 2 (both keys should trigger)", callCount)
	}
}

func TestKeyHandler_UnrecognizedKey_ReturnsNotHandled(t *testing.T) {
	h := KeyHandler{Options: []ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "OK",
			Action: func() tea.Cmd { return nil },
		},
	}}
	_, handled := h.Handle(makeSpecialKeyMsg(tea.KeyEscape))
	if handled {
		t.Error("Handle(Esc when only Enter registered): handled = true, want false")
	}
}

func TestKeyHandler_ScrollKeys_WithScroll(t *testing.T) {
	scroll := &ScrollState{Offset: 5, Total: 30, Viewport: 10}
	h := KeyHandler{Scroll: scroll}

	_, handled := h.Handle(makeSpecialKeyMsg(tea.KeyUp))
	if !handled {
		t.Error("Handle(Up) with Scroll: handled = false, want true")
	}
	if scroll.Offset != 4 {
		t.Errorf("Handle(Up): Offset = %d, want 4", scroll.Offset)
	}

	_, handled = h.Handle(makeSpecialKeyMsg(tea.KeyDown))
	if !handled {
		t.Error("Handle(Down) with Scroll: handled = false, want true")
	}
	if scroll.Offset != 5 {
		t.Errorf("Handle(Down): Offset = %d, want 5", scroll.Offset)
	}

	scroll.Offset = 5
	_, handled = h.Handle(makeSpecialKeyMsg(tea.KeyPgUp))
	if !handled {
		t.Error("Handle(PgUp): handled = false, want true")
	}
	if scroll.Offset != 0 {
		t.Errorf("Handle(PgUp): Offset = %d, want 0 (5-viewport=5-10→0)", scroll.Offset)
	}

	_, handled = h.Handle(makeSpecialKeyMsg(tea.KeyPgDown))
	if !handled {
		t.Error("Handle(PgDown): handled = false, want true")
	}
	if scroll.Offset != 10 {
		t.Errorf("Handle(PgDown): Offset = %d, want 10", scroll.Offset)
	}

	_, handled = h.Handle(makeSpecialKeyMsg(tea.KeyHome))
	if !handled {
		t.Error("Handle(Home): handled = false, want true")
	}
	if scroll.Offset != 0 {
		t.Errorf("Handle(Home): Offset = %d, want 0", scroll.Offset)
	}

	_, handled = h.Handle(makeSpecialKeyMsg(tea.KeyEnd))
	if !handled {
		t.Error("Handle(End): handled = false, want true")
	}
	if scroll.Offset != 20 {
		t.Errorf("Handle(End): Offset = %d, want 20 (total-viewport=30-10)", scroll.Offset)
	}
}

func TestKeyHandler_ScrollKeys_WithoutScroll_NotHandled(t *testing.T) {
	h := KeyHandler{} // Scroll == nil
	_, handled := h.Handle(makeSpecialKeyMsg(tea.KeyUp))
	if handled {
		t.Error("Handle(Up) with nil Scroll: handled = true, want false")
	}
	_, handled = h.Handle(makeSpecialKeyMsg(tea.KeyDown))
	if handled {
		t.Error("Handle(Down) with nil Scroll: handled = true, want false")
	}
}

func TestKeyHandler_EmptyOptions_ScrollStillWorks(t *testing.T) {
	scroll := &ScrollState{Offset: 3, Total: 20, Viewport: 5}
	h := KeyHandler{Scroll: scroll}
	_, handled := h.Handle(makeSpecialKeyMsg(tea.KeyUp))
	if !handled {
		t.Error("Handle(Up) with empty Options but Scroll != nil: not handled")
	}
	if scroll.Offset != 2 {
		t.Errorf("Offset = %d, want 2", scroll.Offset)
	}
}
```

- [ ] **Step 6.2: Run to verify failures**

```
go test ./internal/tui/modal/... -run "TestKeyHandler" -v
```

Expected: FAIL — `KeyHandler` undefined.

- [ ] **Step 6.3: Create `internal/tui/modal/key_handler.go`**

```go
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
```

- [ ] **Step 6.4: Run tests to verify they pass**

```
go test ./internal/tui/modal/... -run "TestKeyHandler" -v
```

Expected: PASS

- [ ] **Step 6.5: Commit**

```bash
git add internal/tui/modal/key_handler.go internal/tui/modal/key_handler_test.go
git commit -m "feat(modal): add KeyHandler — centralised typed-key dispatch for modal actions and scroll"
```

---

## Task 7: Create `modal/frame.go` with `DialogFrame.Render`

This is the most complex task. `DialogFrame.Render` builds the complete dialog string character-by-character following the DS anatomy.

**Files:**
- Create: `internal/tui/modal/frame.go`
- Create: `internal/tui/modal/frame_test.go`
- Create: `internal/tui/modal/testdata/golden/` (directory — created automatically by `-update-golden`)

- [ ] **Step 7.1: Write frame_test.go**

Create `internal/tui/modal/frame_test.go`:

```go
package modal_test

import (
	"testing"

	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/tui/testdata"
)

// renderFrame is a helper that constructs a DialogFrame and renders it.
func renderFrame(
	title, symbol, symbolColor, borderColor, defaultKeyColor string,
	options []modal.ModalOption,
	scroll *modal.ScrollState,
	body string,
) testdata.RenderFn {
	return func(w, h int, theme *design.Theme) string {
		f := modal.DialogFrame{
			Title:           title,
			TitleColor:      theme.Text.Primary,
			Symbol:          symbol,
			SymbolColor:     symbolColor,
			BorderColor:     borderColor,
			Options:         options,
			DefaultKeyColor: defaultKeyColor,
			Scroll:          scroll,
		}
		return f.Render(body, w, theme)
	}
}

func twoOptions(theme *design.Theme) []modal.ModalOption {
	return []modal.ModalOption{
		{Keys: []design.Key{design.Keys.Enter}, Label: "Confirmar", Intent: modal.IntentConfirm, Action: func() tea.Cmd { return nil }},
		{Keys: []design.Key{design.Keys.Esc}, Label: "Cancelar", Intent: modal.IntentCancel, Action: func() tea.Cmd { return nil }},
	}
}

func TestDialogFrame_NoScroll(t *testing.T) {
	theme := design.TokyoNight
	opts := twoOptions(theme)
	body := "Linha 1\nLinha 2\nLinha 3"
	testdata.TestRenderManaged(t, "frame", "no_scroll", []string{"60x10"},
		renderFrame("Título do Diálogo", "", "", theme.Border.Focused, theme.Accent.Primary, opts, nil, body))
}

func TestDialogFrame_WithScrollTop(t *testing.T) {
	theme := design.TokyoNight
	opts := twoOptions(theme)
	scroll := &modal.ScrollState{Offset: 0, Total: 30, Viewport: 8}
	// Body has 8 visible lines (viewport size)
	var lines []string
	for i := 1; i <= 8; i++ {
		lines = append(lines, fmt.Sprintf("Linha %d de 30", i))
	}
	body := strings.Join(lines, "\n")
	testdata.TestRenderManaged(t, "frame", "scroll_top", []string{"60x10"},
		renderFrame("Diálogo com Scroll", "", "", theme.Border.Focused, theme.Accent.Primary, opts, scroll, body))
}

func TestDialogFrame_WithScrollMiddle(t *testing.T) {
	theme := design.TokyoNight
	opts := twoOptions(theme)
	scroll := &modal.ScrollState{Offset: 11, Total: 30, Viewport: 8}
	var lines []string
	for i := 12; i <= 19; i++ {
		lines = append(lines, fmt.Sprintf("Linha %d de 30", i))
	}
	body := strings.Join(lines, "\n")
	testdata.TestRenderManaged(t, "frame", "scroll_middle", []string{"60x10"},
		renderFrame("Diálogo com Scroll", "", "", theme.Border.Focused, theme.Accent.Primary, opts, scroll, body))
}

func TestDialogFrame_WithScrollBottom(t *testing.T) {
	theme := design.TokyoNight
	opts := twoOptions(theme)
	scroll := &modal.ScrollState{Offset: 22, Total: 30, Viewport: 8}
	var lines []string
	for i := 23; i <= 30; i++ {
		lines = append(lines, fmt.Sprintf("Linha %d de 30", i))
	}
	body := strings.Join(lines, "\n")
	testdata.TestRenderManaged(t, "frame", "scroll_bottom", []string{"60x10"},
		renderFrame("Diálogo com Scroll", "", "", theme.Border.Focused, theme.Accent.Primary, opts, scroll, body))
}

func TestDialogFrame_SeverityDestructive(t *testing.T) {
	theme := design.TokyoNight
	sev := design.SeverityDestructive
	opts := []modal.ModalOption{
		{Keys: []design.Key{design.Keys.Enter}, Label: "Excluir", Intent: modal.IntentConfirm, Action: func() tea.Cmd { return nil }},
		{Keys: []design.Key{design.Keys.Esc}, Label: "Cancelar", Intent: modal.IntentCancel, Action: func() tea.Cmd { return nil }},
	}
	body := "Esta ação não pode ser desfeita."
	testdata.TestRenderManaged(t, "frame", "severity_destructive", []string{"60x8"},
		renderFrame("Excluir item", sev.Symbol(), sev.BorderColor(theme), sev.BorderColor(theme), sev.DefaultKeyColor(theme), opts, nil, body))
}

func TestDialogFrame_SeverityError(t *testing.T) {
	theme := design.TokyoNight
	sev := design.SeverityError
	opts := []modal.ModalOption{
		{Keys: []design.Key{design.Keys.Enter}, Label: "OK", Intent: modal.IntentConfirm, Action: func() tea.Cmd { return nil }},
	}
	body := "Ocorreu um erro inesperado."
	testdata.TestRenderManaged(t, "frame", "severity_error", []string{"60x6"},
		renderFrame("Erro", sev.Symbol(), sev.BorderColor(theme), sev.BorderColor(theme), sev.DefaultKeyColor(theme), opts, nil, body))
}
```

> **Note:** `frame_test.go` uses `package modal_test` (external test package). Add `import ("fmt"; "strings"; tea "charm.land/bubbletea/v2")` at the top.

- [ ] **Step 7.2: Verify test file compiles (before Render exists)**

```
go vet ./internal/tui/modal/...
```

Expected: FAIL — `modal.DialogFrame` undefined. That's expected at this stage.

- [ ] **Step 7.3: Create `internal/tui/modal/frame.go`**

```go
package modal

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// DialogFrame define a aparência visual de um diálogo: borda superior com título,
// bordas laterais com indicadores de scroll opcionais, e borda de rodapé com ações.
// Não tem estado próprio — mesmos argumentos = mesmo output (idempotente).
type DialogFrame struct {
	// Title é o texto do cabeçalho.
	Title string
	// TitleColor é a cor do título — ex: theme.Text.Primary. Nunca hardcoded.
	TitleColor string
	// Symbol é o símbolo de severidade ou "" para omitir — ex: design.SymWarning.
	Symbol string
	// SymbolColor é a cor do símbolo — ex: design.SeverityDestructive.BorderColor(theme).
	SymbolColor string
	// BorderColor é a cor de toda a borda — ex: theme.Border.Focused.
	BorderColor string
	// Options lista as ações do rodapé (máximo 3 conforme DS).
	// A 1ª opção usa DefaultKeyColor; as demais usam BorderColor.
	Options []ModalOption
	// DefaultKeyColor é a cor da tecla da 1ª opção (ação principal).
	DefaultKeyColor string
	// Scroll é o estado de scroll para exibir indicadores na borda lateral direita.
	// nil = sem scroll.
	Scroll *ScrollState
}

// Render monta a string completa do diálogo a partir do corpo fornecido.
//
// body é uma string com linhas separadas por \n. Cada linha já deve estar renderizada
// com ANSI e ter largura visual de (maxWidth - 2 - 2*DialogPaddingH) colunas.
// O frame não reaplica padding horizontal — apenas adiciona as bordas laterais.
//
// Algoritmo:
//  1. Borda superior: ╭── [símbolo  ]título ───╮
//  2. Para cada linha do body:
//     │ [padding] linha [padding] │  com background surface.raised
//     Se Scroll != nil, substitui │ direito por indicador de scroll.
//  3. Borda de rodapé com ações posicionadas.
func (f DialogFrame) Render(body string, maxWidth int, theme *design.Theme) string {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(f.BorderColor))
	bgStyle := lipgloss.NewStyle().Background(lipgloss.Color(theme.Surface.Raised))

	innerWidth := maxWidth - 2 // subtrair as duas bordas verticais

	// --- Borda superior ---
	topLine := f.renderTopBorder(innerWidth, borderStyle, theme)

	// --- Linhas do corpo ---
	lines := strings.Split(body, "\n")
	var bodyLines []string
	for i, line := range lines {
		bodyLines = append(bodyLines, f.renderBodyLine(i+1, len(lines), line, innerWidth, borderStyle, bgStyle, theme))
	}

	// --- Borda de rodapé ---
	bottomLine := f.renderBottomBorder(innerWidth, borderStyle, theme)

	var sb strings.Builder
	sb.WriteString(topLine)
	sb.WriteRune('\n')
	for _, l := range bodyLines {
		sb.WriteString(l)
		sb.WriteRune('\n')
	}
	sb.WriteString(bottomLine)
	return sb.String()
}

// renderTopBorder gera a linha: ╭── [símbolo  título] ───╮
func (f DialogFrame) renderTopBorder(innerWidth int, borderStyle lipgloss.Style, theme *design.Theme) string {
	titleText, titleWidth := design.RenderDialogTitle(f.Title, f.Symbol, f.SymbolColor, theme)

	// " título " — 1 espaço antes, 1 espaço depois
	titleSegmentWidth := titleWidth + 2
	// preenchimento ─ à esquerda (mínimo 1) e à direita
	fillLeft := 1
	fillRight := innerWidth - fillLeft - titleSegmentWidth
	if fillRight < 1 {
		fillRight = 1
	}

	tl := borderStyle.Render(design.SymCornerTL)
	tr := borderStyle.Render(design.SymCornerTR)
	dashL := borderStyle.Render(strings.Repeat(design.SymBorderH, fillLeft))
	dashR := borderStyle.Render(strings.Repeat(design.SymBorderH, fillRight))
	space := " "

	return tl + dashL + space + titleText + space + dashR + tr
}

// renderBodyLine gera uma linha do corpo com bordas e indicador de scroll.
// lineNum é 1-based dentro das linhas visíveis (para determinar posição das setas).
func (f DialogFrame) renderBodyLine(lineNum, totalLines int, content string, innerWidth int, borderStyle lipgloss.Style, bgStyle lipgloss.Style, theme *design.Theme) string {
	lBorder := borderStyle.Render(design.SymBorderV)
	paddingH := strings.Repeat(" ", design.DialogPaddingH)

	// Conteúdo com background
	lineContent := bgStyle.Render(paddingH + content + paddingH)

	// Borda direita: seta ou thumb ou │ normal
	var rBorder string
	if f.Scroll != nil {
		isFirstLine := lineNum == 1
		isLastLine := lineNum == totalLines
		thumbLine := f.Scroll.ThumbLine()

		switch {
		case isFirstLine && f.Scroll.CanScrollUp():
			arrow, _ := design.RenderScrollArrow(true, theme)
			rBorder = arrow
		case isLastLine && f.Scroll.CanScrollDown():
			arrow, _ := design.RenderScrollArrow(false, theme)
			rBorder = arrow
		case thumbLine != -1 && lineNum == thumbLine:
			thumb, _ := design.RenderScrollThumb(theme)
			rBorder = thumb
		default:
			rBorder = borderStyle.Render(design.SymBorderV)
		}
	} else {
		rBorder = borderStyle.Render(design.SymBorderV)
	}

	return lBorder + lineContent + rBorder
}

// renderBottomBorder gera a linha de rodapé: ╰─ [ação1] ── [ação2] ── [ação3] ─╯
// Posicionamento:
//   1 ação  → alinhada à direita
//   2 ações → 1ª à esquerda, 2ª à direita
//   3 ações → 1ª à esquerda, 2ª ao centro, 3ª à direita
func (f DialogFrame) renderBottomBorder(innerWidth int, borderStyle lipgloss.Style, theme *design.Theme) string {
	bl := borderStyle.Render(design.SymCornerBL)
	br := borderStyle.Render(design.SymCornerBR)
	dash := borderStyle.Render(design.SymBorderH)

	// Renderizar as ações
	type renderedOpt struct {
		text  string
		width int
	}
	var rendered []renderedOpt
	for i, opt := range f.Options {
		if len(opt.Keys) == 0 {
			continue
		}
		keyColor := f.BorderColor
		if i == 0 {
			keyColor = f.DefaultKeyColor
		}
		text, w := design.RenderDialogAction(opt.Keys[0].Label, opt.Label, keyColor, theme)
		// Adicionar espaços " " em torno da ação (1 espaço antes e depois)
		rendered = append(rendered, renderedOpt{text: " " + text + " ", width: w + 2})
	}

	if len(rendered) == 0 {
		// Sem ações: linha ─ completa
		return bl + borderStyle.Render(strings.Repeat(design.SymBorderH, innerWidth)) + br
	}

	// Montar a linha de rodapé conforme o número de ações
	line := make([]byte, 0, innerWidth*4)
	writeDashes := func(count int) {
		for i := 0; i < count; i++ {
			line = append(line, []byte(dash)...)
		}
	}
	writeAction := func(r renderedOpt) {
		line = append(line, []byte(r.text)...)
	}

	switch len(rendered) {
	case 1:
		// Ação única à direita
		totalActionWidth := rendered[0].width
		fill := innerWidth - totalActionWidth
		if fill < 0 {
			fill = 0
		}
		writeDashes(fill)
		writeAction(rendered[0])
	case 2:
		// 1ª à esquerda, 2ª à direita
		gap := innerWidth - rendered[0].width - rendered[1].width
		if gap < 1 {
			gap = 1
		}
		writeAction(rendered[0])
		writeDashes(gap)
		writeAction(rendered[1])
	case 3:
		// 1ª à esquerda, 2ª ao centro, 3ª à direita
		remaining := innerWidth - rendered[0].width - rendered[1].width - rendered[2].width
		if remaining < 2 {
			remaining = 2
		}
		gapLeft := remaining / 2
		gapRight := remaining - gapLeft
		writeAction(rendered[0])
		writeDashes(gapLeft)
		writeAction(rendered[1])
		writeDashes(gapRight)
		writeAction(rendered[2])
	default:
		// Mais de 3: renderizar apenas as 3 primeiras (DS: máximo 3)
		remaining := innerWidth - rendered[0].width - rendered[1].width - rendered[2].width
		if remaining < 2 {
			remaining = 2
		}
		gapLeft := remaining / 2
		gapRight := remaining - gapLeft
		writeAction(rendered[0])
		writeDashes(gapLeft)
		writeAction(rendered[1])
		writeDashes(gapRight)
		writeAction(rendered[2])
	}

	return bl + string(line) + br
}
```

- [ ] **Step 7.4: Fix frame_test.go imports**

The test file needs `fmt`, `strings`, and `tea`. Update its import block:

```go
import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/tui/testdata"
)
```

> The `twoOptions` helper also needs `tea` for the `tea.Cmd` return type. Move `twoOptions` to use a nil cmd: `Action: func() tea.Cmd { return nil }`.

- [ ] **Step 7.5: Generate golden files**

```
go test ./internal/tui/modal/... -run "TestDialogFrame" -update-golden -v
```

Expected: PASS — golden files created in `internal/tui/modal/testdata/golden/`.

- [ ] **Step 7.6: Verify golden files were created**

```
Get-ChildItem internal/tui/modal/testdata/golden
```

Expected: `.txt` and `.json` files for each test variant.

- [ ] **Step 7.7: Run without -update-golden to verify tests pass normally**

```
go test ./internal/tui/modal/... -run "TestDialogFrame" -v
```

Expected: PASS

- [ ] **Step 7.8: Run full suite**

```
go test ./...
```

Expected: PASS

- [ ] **Step 7.9: Commit**

```bash
git add internal/tui/modal/frame.go internal/tui/modal/frame_test.go internal/tui/modal/testdata/
git commit -m "feat(modal): add DialogFrame — stateless design-system-compliant dialog renderer"
```

---

## Task 8: Rewrite `modal/confirm_modal.go`

Replace the transitional stub (from Task 3) with the full implementation using `DialogFrame` and `KeyHandler`.

**Files:**
- Rewrite: `internal/tui/modal/confirm_modal.go`
- Create: `internal/tui/modal/confirm_modal_test.go`

- [ ] **Step 8.1: Write the failing golden test**

Create `internal/tui/modal/confirm_modal_test.go`:

```go
package modal_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/tui/testdata"
)

func TestConfirmModal_Destructive(t *testing.T) {
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "Excluir",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd { return tui.CloseModal() },
		},
		{
			Keys:   []design.Key{design.Keys.Esc},
			Label:  "Cancelar",
			Intent: modal.IntentCancel,
			Action: func() tea.Cmd { return tui.CloseModal() },
		},
	}
	m := modal.NewConfirmModalSeverity(design.SeverityDestructive,
		"Excluir cofre",
		"Esta ação é permanente e não pode ser desfeita.",
		opts,
	)
	testdata.TestRenderManaged(t, "confirm_modal", "destructive", []string{"60x10"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestConfirmModal_HandleKey_Enter_ExecutesAction(t *testing.T) {
	called := false
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "OK",
			Intent: modal.IntentConfirm,
			Action: func() tea.Cmd {
				called = true
				return nil
			},
		},
	}
	m := modal.NewConfirmModal("Título", "Mensagem", opts)
	_ = m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEnter})
	if !called {
		t.Error("HandleKey(Enter): action not called")
	}
}

func TestConfirmModal_Update_DelegatesKeys(t *testing.T) {
	called := false
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Esc},
			Label:  "Cancelar",
			Intent: modal.IntentCancel,
			Action: func() tea.Cmd {
				called = true
				return nil
			},
		},
	}
	m := modal.NewConfirmModal("Título", "Mensagem", opts)
	_ = m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if !called {
		t.Error("Update(KeyEsc): action not called — Update must delegate to HandleKey")
	}
}
```

- [ ] **Step 8.2: Run to verify test compiles but golden file is missing**

```
go test ./internal/tui/modal/... -run "TestConfirmModal" -v
```

Expected: Some tests might fail with "golden file not found" — that's fine; the unit tests will also fail until the rewrite.

- [ ] **Step 8.3: Rewrite `confirm_modal.go`**

```go
package modal

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// ConfirmModal exibe um diálogo de confirmação com título, mensagem, severidade e ações.
// Implementa tui.ModalView. Criado via NewConfirmModal ou NewConfirmModalSeverity.
type ConfirmModal struct {
	severity design.Severity
	title    string
	message  string
	options  []ModalOption
	keys     KeyHandler // despacha teclas das opções; sem scroll
}

// NewConfirmModal cria um ConfirmModal de severidade Neutra com as opções fornecidas.
// opts define as ações disponíveis — o caller injeta os closures corretos.
// Convenção: 1ª opção é a ação principal (Enter); última é o cancelamento (Esc).
func NewConfirmModal(title, message string, opts []ModalOption) *ConfirmModal {
	return NewConfirmModalSeverity(design.SeverityNeutral, title, message, opts)
}

// NewConfirmModalSeverity cria um ConfirmModal com severidade visual explícita.
func NewConfirmModalSeverity(severity design.Severity, title, message string, opts []ModalOption) *ConfirmModal {
	m := &ConfirmModal{
		severity: severity,
		title:    title,
		message:  message,
		options:  opts,
	}
	m.keys = KeyHandler{Options: opts}
	return m
}

// Render constrói um DialogFrame com cores e símbolo derivados da severidade,
// e passa o corpo (mensagem com padding) para o frame renderizar.
func (m *ConfirmModal) Render(maxHeight, maxWidth int, theme *design.Theme) string {
	// Corpo: padding vertical acima e abaixo da mensagem.
	padding := strings.Repeat("\n", design.DialogPaddingV)
	body := padding + m.message + padding

	frame := DialogFrame{
		Title:           m.title,
		TitleColor:      theme.Text.Primary,
		Symbol:          m.severity.Symbol(),
		SymbolColor:     m.severity.BorderColor(theme),
		BorderColor:     m.severity.BorderColor(theme),
		Options:         m.options,
		DefaultKeyColor: m.severity.DefaultKeyColor(theme),
		Scroll:          nil,
	}
	return frame.Render(body, maxWidth, theme)
}

// HandleKey delega para m.keys.Handle(msg).
func (m *ConfirmModal) HandleKey(msg tea.KeyMsg) tea.Cmd {
	if cmd, handled := m.keys.Handle(msg); handled {
		return cmd
	}
	return nil
}

// Update processa mensagens Bubble Tea. Delega para HandleKey em tea.KeyMsg.
func (m *ConfirmModal) Update(msg tea.Msg) tea.Cmd {
	if key, ok := msg.(tea.KeyMsg); ok {
		return m.HandleKey(key)
	}
	return nil
}
```

- [ ] **Step 8.4: Generate golden file**

```
go test ./internal/tui/modal/... -run "TestConfirmModal_Destructive" -update-golden -v
```

Expected: PASS — golden file `confirm_modal-destructive-60x10.golden.txt` created.

- [ ] **Step 8.5: Run all ConfirmModal tests**

```
go test ./internal/tui/modal/... -run "TestConfirmModal" -v
```

Expected: PASS

- [ ] **Step 8.6: Run full suite**

```
go test ./...
```

Expected: PASS

- [ ] **Step 8.7: Commit**

```bash
git add internal/tui/modal/confirm_modal.go internal/tui/modal/confirm_modal_test.go internal/tui/modal/testdata/
git commit -m "feat(modal): rewrite ConfirmModal using DialogFrame and KeyHandler"
```

---

## Task 9: Rewrite `modal/help_modal.go`

**Files:**
- Rewrite: `internal/tui/modal/help_modal.go`
- Create: `internal/tui/modal/help_modal_test.go`

- [ ] **Step 9.1: Write the failing tests**

Create `internal/tui/modal/help_modal_test.go`:

```go
package modal_test

import (
	"testing"

	"github.com/useful-toys/abditum/internal/tui/actions"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/tui/testdata"
)

func sampleActionsAndGroups() ([]actions.Action, []actions.ActionGroup) {
	groups := []actions.ActionGroup{
		{ID: "app", Label: "Aplicação", Order: 0},
		{ID: "nav", Label: "Navegação", Order: 1},
	}
	acts := []actions.Action{
		{Keys: []design.Key{design.Shortcuts.Help}, Label: "Ajuda", Description: "Abre o diálogo de ajuda.", GroupID: "app", Priority: 10},
		{Keys: []design.Key{design.Shortcuts.Quit}, Label: "Sair", Description: "Encerra a aplicação.", GroupID: "app", Priority: 20},
		{Keys: []design.Key{design.Keys.Enter}, Label: "Abrir", Description: "Abre o item selecionado.", GroupID: "nav", Priority: 10},
		{Keys: []design.Key{design.Keys.Esc}, Label: "Voltar", Description: "Retorna ao nível anterior.", GroupID: "nav", Priority: 20},
	}
	return acts, groups
}

func TestHelpModal_NoScroll(t *testing.T) {
	acts, groups := sampleActionsAndGroups()
	m := modal.NewHelpModal(acts, groups)
	testdata.TestRenderManaged(t, "help_modal", "no_scroll", []string{"60x20"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestHelpModal_WithScroll(t *testing.T) {
	acts, groups := sampleActionsAndGroups()
	m := modal.NewHelpModal(acts, groups)
	// Small height forces scroll
	testdata.TestRenderManaged(t, "help_modal", "with_scroll", []string{"60x6"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestHelpModal_HandleKey_Esc_ClosesModal(t *testing.T) {
	acts, groups := sampleActionsAndGroups()
	m := modal.NewHelpModal(acts, groups)
	cmd := m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEscape})
	if cmd == nil {
		t.Error("HandleKey(Esc): cmd = nil, want CloseModal command")
	}
}

func TestHelpModal_Update_DelegatesKeys(t *testing.T) {
	acts, groups := sampleActionsAndGroups()
	m := modal.NewHelpModal(acts, groups)
	// Update with non-key message should return nil
	cmd := m.Update("not-a-key")
	if cmd != nil {
		t.Error("Update(non-key): expected nil cmd")
	}
}
```

Add `tea "charm.land/bubbletea/v2"` to the import block.

- [ ] **Step 9.2: Run to verify failures**

```
go test ./internal/tui/modal/... -run "TestHelpModal" -v
```

Expected: FAIL (golden files missing, HelpModal not yet rewritten).

- [ ] **Step 9.3: Rewrite `help_modal.go`**

```go
package modal

import (
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/actions"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// helpTitle é o título do modal de ajuda conforme tui-spec-dialog-help.md.
const helpTitle = "Ajuda — Atalhos e Ações"

// keyColumnWidth é a largura reservada para a coluna de teclas no corpo do HelpModal.
const keyColumnWidth = 14

// HelpModal exibe todas as actions registradas, agrupadas por ActionGroup.
// Suporta scroll quando o conteúdo excede o espaço disponível.
// Implementa tui.ModalView.
type HelpModal struct {
	actions []actions.Action
	groups  []actions.ActionGroup
	scroll  ScrollState // estado de scroll — começa em Offset=0
	keys    KeyHandler  // despacha scroll (↑↓PgUp/PgDn/Home/End) e Esc (fechar)
}

// NewHelpModal cria o HelpModal com as actions e grupos fornecidos.
// Scroll começa no topo (Offset = 0).
func NewHelpModal(acts []actions.Action, groups []actions.ActionGroup) *HelpModal {
	m := &HelpModal{
		actions: acts,
		groups:  groups,
	}
	closeOpts := []ModalOption{
		{
			Keys:   []design.Key{design.Keys.Esc},
			Label:  "Fechar",
			Intent: IntentCancel,
			Action: func() tea.Cmd { return tui.CloseModal() },
		},
	}
	m.keys = KeyHandler{
		Options: closeOpts,
		Scroll:  &m.scroll,
	}
	return m
}

// Render gera o corpo dinamicamente, fatia o viewport conforme scroll,
// e passa para DialogFrame.Render.
func (m *HelpModal) Render(maxHeight, maxWidth int, theme *design.Theme) string {
	allLines := m.buildAllLines(maxWidth, theme)

	// viewport = maxHeight - 2 (borda superior + borda de rodapé)
	viewport := maxHeight - 2
	if viewport < 1 {
		viewport = 1
	}

	// Atualizar estado de scroll.
	m.scroll.Total = len(allLines)
	m.scroll.Viewport = viewport

	// Fatiar linhas visíveis.
	start := m.scroll.Offset
	end := start + viewport
	if end > len(allLines) {
		end = len(allLines)
	}
	visibleLines := allLines[start:end]
	body := strings.Join(visibleLines, "\n")

	// Configurar o frame.
	closeOpts := []ModalOption{
		{
			Keys:   []design.Key{design.Keys.Esc},
			Label:  "Fechar",
			Intent: IntentCancel,
			Action: func() tea.Cmd { return tui.CloseModal() },
		},
	}
	var scrollPtr *ScrollState
	if m.scroll.Total > m.scroll.Viewport {
		scrollPtr = &m.scroll
	}

	frame := DialogFrame{
		Title:           helpTitle,
		TitleColor:      theme.Text.Primary,
		Symbol:          "",
		SymbolColor:     "",
		BorderColor:     theme.Border.Default,
		Options:         closeOpts,
		DefaultKeyColor: theme.Accent.Primary,
		Scroll:          scrollPtr,
	}
	return frame.Render(body, maxWidth, theme)
}

// buildAllLines gera todas as linhas de conteúdo do modal de ajuda.
// Grupos ordenados por ActionGroup.Order crescente; ações ordenadas por Action.Priority crescente.
// Linha em branco entre grupos (não antes do primeiro, não após o último).
func (m *HelpModal) buildAllLines(maxWidth int, theme *design.Theme) []string {
	// Ordenar grupos por Order (estável para empates).
	sortedGroups := make([]actions.ActionGroup, len(m.groups))
	copy(sortedGroups, m.groups)
	sort.SliceStable(sortedGroups, func(i, j int) bool {
		return sortedGroups[i].Order < sortedGroups[j].Order
	})

	// Mapear GroupID → actions, ordenadas por Priority.
	groupActions := make(map[string][]actions.Action)
	for _, a := range m.actions {
		groupActions[a.GroupID] = append(groupActions[a.GroupID], a)
	}
	for id := range groupActions {
		sort.SliceStable(groupActions[id], func(i, j int) bool {
			return groupActions[id][i].Priority < groupActions[id][j].Priority
		})
	}

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Text.Secondary)).
		Bold(true)
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Accent.Primary))
	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Text.Primary))

	var allLines []string
	for i, grp := range sortedGroups {
		if i > 0 {
			allLines = append(allLines, "") // linha em branco entre grupos
		}
		// Cabeçalho do grupo
		allLines = append(allLines, headerStyle.Render(grp.Label))

		// Ações do grupo
		for _, act := range groupActions[grp.ID] {
			keyLabel := ""
			if len(act.Keys) > 0 {
				keyLabel = act.Keys[0].Label
			}
			// Coluna de tecla: largura fixa de keyColumnWidth, pad com espaços
			keyRendered := keyStyle.Render(keyLabel)
			keyVisualWidth := lipgloss.Width(keyRendered)
			pad := keyColumnWidth - keyVisualWidth
			if pad < 1 {
				pad = 1
			}
			line := keyRendered + strings.Repeat(" ", pad) + descStyle.Render(act.Description)
			allLines = append(allLines, line)
		}
	}
	return allLines
}

// HandleKey delega para m.keys.Handle(msg).
func (m *HelpModal) HandleKey(msg tea.KeyMsg) tea.Cmd {
	if cmd, handled := m.keys.Handle(msg); handled {
		return cmd
	}
	return nil
}

// Update processa mensagens Bubble Tea. Delega para HandleKey em tea.KeyMsg.
func (m *HelpModal) Update(msg tea.Msg) tea.Cmd {
	if key, ok := msg.(tea.KeyMsg); ok {
		return m.HandleKey(key)
	}
	return nil
}
```

- [ ] **Step 9.4: Generate golden files**

```
go test ./internal/tui/modal/... -run "TestHelpModal_NoScroll|TestHelpModal_WithScroll" -update-golden -v
```

Expected: PASS — golden files created.

- [ ] **Step 9.5: Run all HelpModal tests**

```
go test ./internal/tui/modal/... -run "TestHelpModal" -v
```

Expected: PASS

- [ ] **Step 9.6: Run full suite**

```
go test ./...
```

Expected: PASS

- [ ] **Step 9.7: Commit**

```bash
git add internal/tui/modal/help_modal.go internal/tui/modal/help_modal_test.go internal/tui/modal/testdata/
git commit -m "feat(modal): rewrite HelpModal using DialogFrame, ScrollState, and KeyHandler"
```

---

## Task 10: Final verification — build and test everything

- [ ] **Step 10.1: Run the complete test suite**

```
go test ./...
```

Expected: PASS — all tests pass.

- [ ] **Step 10.2: Verify the application builds**

```
go build ./...
```

Expected: SUCCESS — no compilation errors.

- [ ] **Step 10.3: Verify no hardcoded colours or symbols in modal/**

```
Select-String -Path "internal/tui/modal/*.go" -Pattern '"#[0-9a-fA-F]{3,8}"'
```

Expected: No matches.

```
Select-String -Path "internal/tui/modal/*.go" -Pattern '"[╭╮╰╯│─■↑↓⚠✕ℹ]"'
```

Expected: No matches (all come through design constants).

- [ ] **Step 10.4: Final commit if any minor fixes were needed**

```bash
git add -A
git commit -m "chore(modal): final cleanup and verification"
```

---

## Self-Review Checklist

**Spec coverage:**

| Requirement | Task |
|---|---|
| `design.Keys.Up` / `Keys.Down` (needed by KeyHandler) | Task 1 |
| `ActionGroup.Order int` | Task 2 |
| `ModalOption.Keys []design.Key` | Task 3 |
| `Severity` type + `Symbol()` / `BorderColor()` / `DefaultKeyColor()` | Task 4 |
| `RenderDialogTitle`, `RenderDialogAction`, `RenderScrollArrow`, `RenderScrollThumb` | Task 4 |
| `DialogPaddingH`, `DialogPaddingV` constants | Task 4 |
| `ScrollState` with all navigation methods | Task 5 |
| `ThumbLine()` with arrow-priority rule | Task 5 |
| `KeyHandler.Handle()` with typed key matching | Task 6 |
| `DialogFrame.Render()` — stateless, character-by-character border | Task 7 |
| Frame golden files (no_scroll, scroll_top/middle/bottom, severity variants) | Task 7 |
| `ConfirmModal` rewrite using `DialogFrame` + `KeyHandler` | Task 8 |
| `ConfirmModal` golden file (destructive severity) | Task 8 |
| `HelpModal` rewrite with scroll, group ordering, column alignment | Task 9 |
| `HelpModal` golden files (no_scroll, with_scroll) | Task 9 |
| `HelpModal` title = `"Ajuda — Atalhos e Ações"` | Task 9 |
| No hardcoded colours or symbols in `modal/` | Task 10 |
| `tui.ModalView` interface unchanged | All tasks (never touched `modal.go`) |

**No placeholders found** — all code blocks are complete.

**Type consistency check:**
- `ScrollState` defined in Task 5, used in Tasks 6, 7, 8, 9 — consistent.
- `KeyHandler` defined in Task 6, used in Tasks 8, 9 — consistent.
- `DialogFrame` defined in Task 7, used in Tasks 8, 9 — consistent.
- `design.Severity` defined in Task 4, used in Tasks 8 — consistent.
- `design.Keys.Up` / `Keys.Down` defined in Task 1, used in Task 6 — consistent.
- `ModalOption.Keys []design.Key` — established in Task 3, consistent throughout.
