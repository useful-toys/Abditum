# Password Modals Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implementar `PasswordEntry` e `PasswordCreate` — dois modais de senha com campo seguro customizado (`[]byte` + `crypto.Wipe`), cursor real (`tea.Cursor`), medidor de força de senha, e testes golden, sem dependências externas além do que já existe no projeto.

**Architecture:** Campo de senha (`PasswordField`) é um struct isolado com buffer `[]byte`; os dois modais o usam e implementam a interface `ModalView` estendida com `Cursor()`. O `RootModel.View()` propaga o cursor do modal ativo usando `lipgloss.Size` para calcular a posição absoluta na tela.

**Tech Stack:** Go, BubbleTea v2 (`charm.land/bubbletea/v2`), Lipgloss v2 (`charm.land/lipgloss/v2`), `internal/crypto.Wipe`, framework de golden files customizado em `internal/tui/testdata`.

---

## File Map

### Novos arquivos
- `internal/tui/modal/password_strength.go` — `StrengthLevel`, `EvaluateStrength`, `RenderStrengthMeter`
- `internal/tui/modal/password_field.go` — `PasswordField` (buffer seguro, HandleKey, Render, Wipe)
- `internal/tui/modal/password_entry_modal.go` — `PasswordEntryModal` (1 campo, cursor, mensagens)
- `internal/tui/modal/password_create_modal.go` — `PasswordCreateModal` (2 campos, medidor, cursor, mensagens)
- `internal/tui/modal/password_strength_test.go` — testes unitários + golden do medidor
- `internal/tui/modal/password_field_test.go` — testes unitários + golden do campo
- `internal/tui/modal/password_entry_modal_test.go` — testes golden do modal de entrada
- `internal/tui/modal/password_create_modal_test.go` — testes golden do modal de criação

### Arquivos modificados
- `internal/tui/modal.go` — adicionar método `Cursor(topY, leftX int) *tea.Cursor` à interface `ModalView`
- `internal/tui/modal/confirm_modal.go` — implementar `Cursor()` retornando `nil`
- `internal/tui/modal/help_modal.go` — implementar `Cursor()` retornando `nil`
- `internal/tui/root.go` — propagar cursor do modal ativo em `View()`

---

## Task 1: Estender interface ModalView com Cursor()

**Files:**
- Modify: `internal/tui/modal.go`
- Modify: `internal/tui/modal/confirm_modal.go`
- Modify: `internal/tui/modal/help_modal.go`

- [ ] **Step 1: Adicionar `Cursor()` à interface `ModalView`**

Em `internal/tui/modal.go`, substituir a interface:

```go
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
```

- [ ] **Step 2: Implementar `Cursor()` em `ConfirmModal`**

Adicionar ao final de `internal/tui/modal/confirm_modal.go`:

```go
// Cursor retorna nil — ConfirmModal não tem campo de texto com cursor.
func (m *ConfirmModal) Cursor(_, _ int) *tea.Cursor {
	return nil
}
```

- [ ] **Step 3: Implementar `Cursor()` em `HelpModal`**

Adicionar ao final de `internal/tui/modal/help_modal.go`:

```go
// Cursor retorna nil — HelpModal não tem campo de texto com cursor.
func (m *HelpModal) Cursor(_, _ int) *tea.Cursor {
	return nil
}
```

- [ ] **Step 4: Verificar que o projeto compila**

```
go build ./internal/tui/...
```

Resultado esperado: nenhuma saída (zero erros).

- [ ] **Step 5: Commit**

```
git add internal/tui/modal.go internal/tui/modal/confirm_modal.go internal/tui/modal/help_modal.go
git commit -m "feat(tui): extend ModalView interface with Cursor() method"
```

---

## Task 2: Propagar cursor do modal ativo em RootModel.View()

**Files:**
- Modify: `internal/tui/root.go`

- [ ] **Step 1: Modificar o bloco `if len(r.modals) > 0` em `View()`**

Localizar em `internal/tui/root.go` o bloco que começa em `if len(r.modals) > 0 {` (linha ~253) e substituí-lo por:

```go
	if len(r.modals) > 0 {
		top := r.modals[len(r.modals)-1]
		// 1 line padding above and below modal on screen.
		modalH := r.height - 2
		modalContent := top.Render(modalH, r.width, r.theme)
		// Center modal content horizontally within available space.
		centeredModal := lipgloss.Place(r.width, modalH, lipgloss.Center, lipgloss.Center, modalContent)

		// Calcular posição absoluta do modal para repassar ao método Cursor do modal.
		// modalContent tem dimensões reais; o compositor aplica Y(1) — offset de 1 linha.
		modalW, modalActualH := lipgloss.Size(modalContent)
		topY := 1 + (modalH-modalActualH)/2  // 1 = offset Y do layer compositor
		leftX := (r.width - modalW) / 2

		// Compose modal (z=1) over base layout (z=0) using lipgloss v2 compositor.
		result := lipgloss.NewCompositor(
			lipgloss.NewLayer(base),
			lipgloss.NewLayer(centeredModal).Y(1).Z(1),
		).Render()
		v := tea.NewView(result)
		v.AltScreen = true
		v.BackgroundColor = lipgloss.Color(r.theme.Surface.Base)
		if c := top.Cursor(topY, leftX); c != nil {
			v.Cursor = c
		}
		return v
	}
```

- [ ] **Step 2: Verificar que o projeto compila**

```
go build ./internal/tui/...
```

Resultado esperado: nenhuma saída.

- [ ] **Step 3: Commit**

```
git add internal/tui/root.go
git commit -m "feat(tui): propagate modal cursor to tea.View in RootModel"
```

---

## Task 3: Implementar StrengthScore (password_strength.go)

**Files:**
- Create: `internal/tui/modal/password_strength.go`
- Create: `internal/tui/modal/password_strength_test.go`

- [ ] **Step 1: Escrever o teste unitário de `EvaluateStrength`**

Criar `internal/tui/modal/password_strength_test.go`:

```go
package modal_test

import (
	"testing"

	"github.com/useful-toys/abditum/internal/tui/modal"
)

func TestEvaluateStrength(t *testing.T) {
	cases := []struct {
		name     string
		password []byte
		want     modal.StrengthLevel
	}{
		{"empty", []byte{}, modal.StrengthWeak},
		{"short only lower", []byte("abc"), modal.StrengthWeak},
		{"short with upper", []byte("Abc"), modal.StrengthWeak},
		{"short with digit", []byte("Ab1"), modal.StrengthWeak},
		// 2 critérios → StrengthFair
		{"long lower only", []byte("abcdefghijkl"), modal.StrengthFair},        // comprimento >= 12
		{"short upper+digit", []byte("Ab1"), modal.StrengthWeak},
		{"upper+digit", []byte("Abc1"), modal.StrengthFair},                    // upper + digit (sem len, sem symbol)
		{"len+upper", []byte("Abcdefghijkl"), modal.StrengthFair},              // len + upper
		// 3 critérios → StrengthFair
		{"len+upper+digit", []byte("Abcdefghijk1"), modal.StrengthFair},        // len + upper + digit
		{"upper+digit+symbol", []byte("Ab1!"), modal.StrengthFair},             // upper + digit + symbol
		// 4 critérios → StrengthStrong
		{"all criteria", []byte("Abcdefghijk1!"), modal.StrengthStrong},        // len + upper + digit + symbol
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := modal.EvaluateStrength(tc.password)
			if got != tc.want {
				t.Errorf("EvaluateStrength(%q) = %d, want %d", tc.password, got, tc.want)
			}
		})
	}
}
```

- [ ] **Step 2: Rodar o teste para confirmar que falha**

```
go test ./internal/tui/modal/... -run TestEvaluateStrength -v
```

Resultado esperado: erro de compilação (tipo/função não definidos).

- [ ] **Step 3: Criar `password_strength.go`**

Criar `internal/tui/modal/password_strength.go`:

```go
package modal

import (
	"strings"
	"unicode"

	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// StrengthLevel classifica a força de uma senha em três níveis.
type StrengthLevel int

const (
	// StrengthWeak corresponde a 0 ou 1 critério satisfeito.
	StrengthWeak StrengthLevel = iota
	// StrengthFair corresponde a 2 ou 3 critérios satisfeitos.
	StrengthFair
	// StrengthStrong corresponde a todos os 4 critérios satisfeitos.
	StrengthStrong
)

// symbolChars são os símbolos reconhecidos pelo critério de símbolo da heurística.
const symbolChars = "!@#$%^&*()-_=+[]{}|;:,.<>?/~"

// EvaluateStrength avalia a força de uma senha usando 4 critérios (1 ponto cada):
//   - Comprimento >= 12 caracteres
//   - Contém pelo menos 1 letra maiúscula
//   - Contém pelo menos 1 dígito numérico
//   - Contém pelo menos 1 símbolo de symbolChars
//
// Retorna StrengthWeak (0–1 pts), StrengthFair (2–3 pts) ou StrengthStrong (4 pts).
func EvaluateStrength(password []byte) StrengthLevel {
	if len(password) == 0 {
		return StrengthWeak
	}

	points := 0

	if len(password) >= 12 {
		points++
	}

	hasUpper := false
	hasDigit := false
	hasSymbol := false
	for _, b := range password {
		r := rune(b)
		if unicode.IsUpper(r) {
			hasUpper = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
		if strings.ContainsRune(symbolChars, r) {
			hasSymbol = true
		}
	}
	if hasUpper {
		points++
	}
	if hasDigit {
		points++
	}
	if hasSymbol {
		points++
	}

	switch {
	case points >= 4:
		return StrengthStrong
	case points >= 2:
		return StrengthFair
	default:
		return StrengthWeak
	}
}

// strengthMeterBlocks é o total de blocos na barra de progresso.
const strengthMeterBlocks = 10

// RenderStrengthMeter renderiza a linha do medidor de força de senha.
// Formato: "Força: ████████░░ Boa" — 10 blocos, label de acordo com o nível.
// Deve ser chamado apenas quando a senha não está vazia (len(password) > 0).
// innerWidth é a largura disponível (sem bordas e sem padding do modal).
func RenderStrengthMeter(password []byte, innerWidth int, theme *design.Theme) string {
	level := EvaluateStrength(password)

	filled := filledBlocks(level)
	empty := strengthMeterBlocks - filled

	var barColor string
	var label string
	switch level {
	case StrengthStrong:
		barColor = theme.Semantic.Success
		label = "✓ Forte"
	case StrengthFair:
		barColor = theme.Semantic.Success
		label = "Boa"
	default: // StrengthWeak
		barColor = theme.Semantic.Warning
		label = "⚠ Fraca"
	}

	filledStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(barColor))
	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Disabled))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(barColor))
	prefixStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Secondary))

	bar := filledStyle.Render(strings.Repeat("█", filled)) +
		emptyStyle.Render(strings.Repeat("░", empty))

	return prefixStyle.Render("Força: ") + bar + " " + labelStyle.Render(label)
}

// filledBlocks calcula quantos blocos preenchidos exibir para cada nível de força.
// Escala: Weak=2, Fair=8, Strong=10 — representação visual intuitiva.
func filledBlocks(level StrengthLevel) int {
	switch level {
	case StrengthStrong:
		return 10
	case StrengthFair:
		return 8
	default: // StrengthWeak (0 ou 1 ponto)
		return 2
	}
}
```

- [ ] **Step 4: Rodar o teste unitário**

```
go test ./internal/tui/modal/... -run TestEvaluateStrength -v
```

Resultado esperado: `PASS` para todos os casos.

- [ ] **Step 5: Escrever o teste golden do medidor**

Adicionar em `internal/tui/modal/password_strength_test.go`:

```go
func TestStrengthMeter_Weak(t *testing.T) {
	testdata.TestRenderManaged(t, "strength_meter", "weak", []string{"44x1"},
		func(w, h int, theme *design.Theme) string {
			// Senha com 1 critério (apenas maiúscula) → StrengthWeak
			return modal.RenderStrengthMeter([]byte("A"), w, theme)
		})
}

func TestStrengthMeter_Fair(t *testing.T) {
	testdata.TestRenderManaged(t, "strength_meter", "fair", []string{"44x1"},
		func(w, h int, theme *design.Theme) string {
			// Senha com 2 critérios (maiúscula + dígito) → StrengthFair
			return modal.RenderStrengthMeter([]byte("Ab1"), w, theme)
		})
}

func TestStrengthMeter_Strong(t *testing.T) {
	testdata.TestRenderManaged(t, "strength_meter", "strong", []string{"44x1"},
		func(w, h int, theme *design.Theme) string {
			// Senha com todos os critérios → StrengthStrong
			return modal.RenderStrengthMeter([]byte("Abcdefghijk1!"), w, theme)
		})
}
```

Adicionar o import necessário no topo do arquivo de teste:

```go
import (
	"testing"

	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/tui/testdata"
)
```

- [ ] **Step 6: Gerar os golden files**

```
go test ./internal/tui/modal/... -run TestStrengthMeter -update-golden -v
```

Resultado esperado: 3 subpasses, nenhum erro. Arquivos criados em `internal/tui/modal/testdata/golden/`.

- [ ] **Step 7: Verificar os golden files visualmente**

Ler `internal/tui/modal/testdata/golden/strength_meter-weak-44x1.golden.txt` e verificar:
- Linha única com `Força: ██░░░░░░░░ ⚠ Fraca` (2 blocos preenchidos, 8 vazios)
- `internal/tui/modal/testdata/golden/strength_meter-fair-44x1.golden.txt`: `Força: ████████░░ Boa` (8 blocos, 2 vazios)
- `internal/tui/modal/testdata/golden/strength_meter-strong-44x1.golden.txt`: `Força: ██████████ ✓ Forte` (10 blocos)

- [ ] **Step 8: Rodar todos os testes do pacote modal**

```
go test ./internal/tui/modal/... -v
```

Resultado esperado: todos passam.

- [ ] **Step 9: Commit**

```
git add internal/tui/modal/password_strength.go internal/tui/modal/password_strength_test.go internal/tui/modal/testdata/
git commit -m "feat(modal): add password strength evaluator and strength meter"
```

---

## Task 4: Implementar PasswordField

**Files:**
- Create: `internal/tui/modal/password_field.go`
- Create: `internal/tui/modal/password_field_test.go`

- [ ] **Step 1: Escrever testes unitários do campo**

Criar `internal/tui/modal/password_field_test.go`:

```go
package modal_test

import (
	"bytes"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/tui/testdata"
)

func TestPasswordField_InitiallyEmpty(t *testing.T) {
	f := modal.NewPasswordField("Senha")
	if f.Len() != 0 {
		t.Errorf("NewPasswordField: expected Len()=0, got %d", f.Len())
	}
	if len(f.Value()) != 0 {
		t.Errorf("NewPasswordField: expected empty Value(), got %q", f.Value())
	}
}

func TestPasswordField_HandleKey_Rune(t *testing.T) {
	f := modal.NewPasswordField("Senha")
	consumed := f.HandleKey(tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{'a'}})
	if !consumed {
		t.Error("HandleKey: rune should be consumed")
	}
	if f.Len() != 1 {
		t.Errorf("HandleKey: expected Len()=1, got %d", f.Len())
	}
	if !bytes.Equal(f.Value(), []byte("a")) {
		t.Errorf("HandleKey: expected Value()=%q, got %q", "a", f.Value())
	}
}

func TestPasswordField_HandleKey_Backspace(t *testing.T) {
	f := modal.NewPasswordField("Senha")
	f.HandleKey(tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{'a'}})
	f.HandleKey(tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{'b'}})
	consumed := f.HandleKey(tea.KeyMsg{Type: tea.KeyBackspace})
	if !consumed {
		t.Error("HandleKey: backspace should be consumed")
	}
	if f.Len() != 1 {
		t.Errorf("HandleKey: after backspace expected Len()=1, got %d", f.Len())
	}
	if !bytes.Equal(f.Value(), []byte("a")) {
		t.Errorf("HandleKey: after backspace expected Value()=%q, got %q", "a", f.Value())
	}
}

func TestPasswordField_HandleKey_BackspaceEmpty(t *testing.T) {
	f := modal.NewPasswordField("Senha")
	consumed := f.HandleKey(tea.KeyMsg{Type: tea.KeyBackspace})
	if !consumed {
		t.Error("HandleKey: backspace on empty should still be consumed")
	}
	if f.Len() != 0 {
		t.Errorf("HandleKey: Len() should remain 0, got %d", f.Len())
	}
}

func TestPasswordField_HandleKey_ArrowNotConsumed(t *testing.T) {
	f := modal.NewPasswordField("Senha")
	consumed := f.HandleKey(tea.KeyMsg{Type: tea.KeyRight})
	if consumed {
		t.Error("HandleKey: arrow key should not be consumed")
	}
}

func TestPasswordField_Clear(t *testing.T) {
	f := modal.NewPasswordField("Senha")
	f.HandleKey(tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{'x'}})
	f.Clear()
	if f.Len() != 0 {
		t.Errorf("Clear: expected Len()=0, got %d", f.Len())
	}
}

func TestPasswordField_Value_ReturnsCopy(t *testing.T) {
	f := modal.NewPasswordField("Senha")
	f.HandleKey(tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{'a'}})
	v1 := f.Value()
	v1[0] = 'z'  // modificar a cópia não deve afetar o campo
	if f.Value()[0] != 'a' {
		t.Error("Value() must return a copy, not a reference to internal buffer")
	}
}

func TestPasswordField_Wipe(t *testing.T) {
	f := modal.NewPasswordField("Senha")
	f.HandleKey(tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{'s'}})
	f.HandleKey(tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{'e'}})
	f.Wipe()
	if f.Len() != 0 {
		t.Errorf("Wipe: expected Len()=0, got %d", f.Len())
	}
}

// Testes golden de renderização

func TestPasswordField_EmptyFocused(t *testing.T) {
	testdata.TestRenderManaged(t, "password_field", "empty_focused", []string{"44x2"},
		func(w, h int, theme *design.Theme) string {
			f := modal.NewPasswordField("Senha")
			return f.Render(w, true, theme)
		})
}

func TestPasswordField_EmptyBlurred(t *testing.T) {
	testdata.TestRenderManaged(t, "password_field", "empty_blurred", []string{"44x2"},
		func(w, h int, theme *design.Theme) string {
			f := modal.NewPasswordField("Senha")
			return f.Render(w, false, theme)
		})
}

func TestPasswordField_ContentFocused(t *testing.T) {
	testdata.TestRenderManaged(t, "password_field", "content_focused", []string{"44x2"},
		func(w, h int, theme *design.Theme) string {
			f := modal.NewPasswordField("Senha")
			for _, r := range "12345678" {
				f.HandleKey(tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{r}})
			}
			return f.Render(w, true, theme)
		})
}

func TestPasswordField_ContentBlurred(t *testing.T) {
	testdata.TestRenderManaged(t, "password_field", "content_blurred", []string{"44x2"},
		func(w, h int, theme *design.Theme) string {
			f := modal.NewPasswordField("Senha")
			for _, r := range "12345678" {
				f.HandleKey(tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{r}})
			}
			return f.Render(w, false, theme)
		})
}
```

- [ ] **Step 2: Rodar para confirmar que falha**

```
go test ./internal/tui/modal/... -run TestPasswordField -v
```

Resultado esperado: erro de compilação (tipo não definido).

- [ ] **Step 3: Criar `password_field.go`**

Criar `internal/tui/modal/password_field.go`:

```go
package modal

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/crypto"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// PasswordField é um campo de entrada de senha seguro.
// Gerencia o buffer como []byte e nunca o expõe como string.
// Suporta apenas append (rune imprimível) e remoção do último byte (Backspace).
// Navegação interna (←→, Home, End, Delete) não é suportada — o cursor fica sempre ao final.
type PasswordField struct {
	label string // texto exibido acima da área digitável
	value []byte // buffer da senha — nunca convertido para string
}

// NewPasswordField cria um PasswordField com o label fornecido e buffer vazio.
func NewPasswordField(label string) *PasswordField {
	return &PasswordField{label: label}
}

// Value retorna uma cópia do buffer interno.
// O caller é responsável por zerar a cópia após o uso.
func (f *PasswordField) Value() []byte {
	if len(f.value) == 0 {
		return []byte{}
	}
	cp := make([]byte, len(f.value))
	copy(cp, f.value)
	return cp
}

// Len retorna o comprimento atual do buffer.
func (f *PasswordField) Len() int {
	return len(f.value)
}

// Clear zera logicamente o buffer sem wipe criptográfico (realoca).
// Use Wipe quando o conteúdo for sensível e deva ser apagado da memória.
func (f *PasswordField) Clear() {
	f.value = nil
}

// Wipe apaga o conteúdo do buffer com crypto.Wipe (zeragem byte a byte) e depois limpa.
// Deve ser chamado ao cancelar ou ao descartar o valor após uso.
func (f *PasswordField) Wipe() {
	crypto.Wipe(f.value)
	f.value = nil
}

// HandleKey processa um evento de teclado.
// Retorna true se a tecla foi consumida pelo campo; false caso contrário.
//   - Rune imprimível: faz append ao buffer.
//   - Backspace: remove o último byte (no-op se vazio).
//   - Qualquer outra tecla: não consome (retorna false).
func (f *PasswordField) HandleKey(msg tea.KeyMsg) bool {
	switch msg.Type {
	case tea.KeyRune:
		f.value = append(f.value, []byte(string(msg.Runes))...)
		return true
	case tea.KeyBackspace:
		if len(f.value) > 0 {
			f.value = f.value[:len(f.value)-1]
		}
		return true
	}
	return false
}

// Render retorna duas linhas separadas por "\n":
//   - Linha 1: label com cor dependente do foco (accent.primary bold se focado, text.secondary se não)
//   - Linha 2: área digitável com largura innerWidth e fundo surface.input
//     — vazia: somente espaços (fundo visível)
//     — com conteúdo: f.Len() bullets "•" em text.secondary + espaços até a largura total
func (f *PasswordField) Render(innerWidth int, focused bool, theme *design.Theme) string {
	// Linha 1: label
	var labelStyle lipgloss.Style
	if focused {
		labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Accent.Primary)).
			Bold(true)
	} else {
		labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Text.Secondary))
	}
	labelLine := labelStyle.Render(f.label)

	// Linha 2: área digitável
	bulletStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Text.Secondary))
	inputBg := lipgloss.NewStyle().
		Background(lipgloss.Color(theme.Surface.Input))

	bullets := strings.Repeat("•", f.Len())
	padding := strings.Repeat(" ", innerWidth-f.Len())
	if f.Len() > innerWidth {
		// Truncar visualmente se exceder a largura (segurança: não mostra mais que innerWidth bullets)
		bullets = strings.Repeat("•", innerWidth)
		padding = ""
	}
	inputContent := bulletStyle.Render(bullets) + padding
	inputLine := inputBg.Render(inputContent)

	return labelLine + "\n" + inputLine
}
```

- [ ] **Step 4: Rodar testes unitários do campo**

```
go test ./internal/tui/modal/... -run TestPasswordField_Init -v
go test ./internal/tui/modal/... -run TestPasswordField_HandleKey -v
go test ./internal/tui/modal/... -run TestPasswordField_Clear -v
go test ./internal/tui/modal/... -run TestPasswordField_Value -v
go test ./internal/tui/modal/... -run TestPasswordField_Wipe -v
```

Resultado esperado: todos `PASS`.

- [ ] **Step 5: Gerar golden files do campo**

```
go test ./internal/tui/modal/... -run TestPasswordField_Empty -update-golden -v
go test ./internal/tui/modal/... -run TestPasswordField_Content -update-golden -v
```

- [ ] **Step 6: Verificar os golden files**

Ler os 4 arquivos `.golden.txt` e confirmar:
- `password_field-empty_focused-44x2.golden.txt`: Linha 1 = `Senha` (sem ANSI), linha 2 = 44 espaços
- `password_field-empty_blurred-44x2.golden.txt`: igual mas label em cor diferente (verificável no `.json`)
- `password_field-content_focused-44x2.golden.txt`: Linha 2 = `••••••••` + 36 espaços
- `password_field-content_blurred-44x2.golden.txt`: mesmo conteúdo, label sem bold

- [ ] **Step 7: Rodar todos os testes do pacote**

```
go test ./internal/tui/modal/... -v
```

Resultado esperado: todos passam.

- [ ] **Step 8: Commit**

```
git add internal/tui/modal/password_field.go internal/tui/modal/password_field_test.go internal/tui/modal/testdata/
git commit -m "feat(modal): add secure PasswordField with bullet masking and wipe"
```

---

## Task 5: Implementar PasswordEntryModal

**Files:**
- Create: `internal/tui/modal/password_entry_modal.go`
- Create: `internal/tui/modal/password_entry_modal_test.go`

- [ ] **Step 1: Escrever teste golden do modal vazio**

Criar `internal/tui/modal/password_entry_modal_test.go`:

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

// stubMessageController implementa tui.MessageController com métodos no-op.
// Reutilizado em todos os testes de modal de senha.
type stubMessageController struct{}

func (s *stubMessageController) SetBusy(string)      {}
func (s *stubMessageController) SetSuccess(string)   {}
func (s *stubMessageController) SetError(string)     {}
func (s *stubMessageController) SetWarning(string)   {}
func (s *stubMessageController) SetInfo(string)      {}
func (s *stubMessageController) SetHintField(string) {}
func (s *stubMessageController) SetHintUsage(string) {}
func (s *stubMessageController) Clear()              {}

var _ tui.MessageController = (*stubMessageController)(nil)

func TestPasswordEntryModal_Empty(t *testing.T) {
	mc := &stubMessageController{}
	m := modal.NewPasswordEntryModal(mc,
		func(_ []byte) tea.Cmd { return nil },
		func() tea.Cmd { return nil },
	)
	testdata.TestRenderManaged(t, "password_entry", "empty", []string{"50x7"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestPasswordEntryModal_WithContent(t *testing.T) {
	mc := &stubMessageController{}
	m := modal.NewPasswordEntryModal(mc,
		func(_ []byte) tea.Cmd { return nil },
		func() tea.Cmd { return nil },
	)
	for _, r := range "12345678" {
		m.Update(tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{r}})
	}
	testdata.TestRenderManaged(t, "password_entry", "with_content", []string{"50x7"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestPasswordEntryModal_NotifyWrongPassword(t *testing.T) {
	mc := &stubMessageController{}
	m := modal.NewPasswordEntryModal(mc,
		func(_ []byte) tea.Cmd { return nil },
		func() tea.Cmd { return nil },
	)
	for _, r := range "12345678" {
		m.Update(tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{r}})
	}
	m.NotifyWrongPassword()
	if m.Len() != 0 {
		t.Errorf("NotifyWrongPassword: expected field cleared, got Len()=%d", m.Len())
	}
}

func TestPasswordEntryModal_Cursor_Empty(t *testing.T) {
	mc := &stubMessageController{}
	m := modal.NewPasswordEntryModal(mc,
		func(_ []byte) tea.Cmd { return nil },
		func() tea.Cmd { return nil },
	)
	c := m.Cursor(1, 0)
	if c == nil {
		t.Fatal("Cursor: expected non-nil cursor")
	}
	// topY=1, leftX=0, borda=1, linhaDoFieldNoBody=2, paddingH=2, len=0
	// Y = 1 + 1 + 2 = 4
	// X = 0 + 1 + 2 + 0 = 3
	if c.Y != 4 {
		t.Errorf("Cursor.Y: expected 4, got %d", c.Y)
	}
	if c.X != 3 {
		t.Errorf("Cursor.X: expected 3, got %d", c.X)
	}
}

func TestPasswordEntryModal_Cursor_WithContent(t *testing.T) {
	mc := &stubMessageController{}
	m := modal.NewPasswordEntryModal(mc,
		func(_ []byte) tea.Cmd { return nil },
		func() tea.Cmd { return nil },
	)
	for _, r := range "12345" {
		m.Update(tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{r}})
	}
	c := m.Cursor(1, 0)
	if c == nil {
		t.Fatal("Cursor: expected non-nil cursor")
	}
	// Y = 1 + 1 + 2 = 4
	// X = 0 + 1 + 2 + 5 = 8
	if c.Y != 4 {
		t.Errorf("Cursor.Y: expected 4, got %d", c.Y)
	}
	if c.X != 8 {
		t.Errorf("Cursor.X: expected 8, got %d", c.X)
	}
}
```

- [ ] **Step 2: Rodar para confirmar que falha**

```
go test ./internal/tui/modal/... -run TestPasswordEntry -v
```

Resultado esperado: erro de compilação.

- [ ] **Step 3: Criar `password_entry_modal.go`**

Criar `internal/tui/modal/password_entry_modal.go`:

```go
package modal

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// passwordEntryTitle é o título do modal conforme a spec.
const passwordEntryTitle = "Senha mestra"

// passwordEntryWidth é a largura fixa do modal em colunas.
const passwordEntryWidth = 50

// PasswordEntryModal exibe o diálogo de entrada de senha para abrir o cofre.
// Tem um único campo de senha. O modal notifica o orquestrador via callbacks —
// não sabe se a senha foi aceita ou rejeitada.
// Implementa tui.ModalView.
type PasswordEntryModal struct {
	mc        tui.MessageController
	field     *PasswordField
	onConfirm func(password []byte) tea.Cmd
	onCancel  func() tea.Cmd
}

// NewPasswordEntryModal cria o modal e emite a dica inicial na barra de status.
func NewPasswordEntryModal(
	mc tui.MessageController,
	onConfirm func(password []byte) tea.Cmd,
	onCancel func() tea.Cmd,
) *PasswordEntryModal {
	m := &PasswordEntryModal{
		mc:        mc,
		field:     NewPasswordField("Senha"),
		onConfirm: onConfirm,
		onCancel:  onCancel,
	}
	mc.SetHintField("• Digite a senha para desbloquear o cofre")
	return m
}

// Len retorna o comprimento atual do campo de senha.
// Usado pelos testes para verificar o estado do campo.
func (m *PasswordEntryModal) Len() int {
	return m.field.Len()
}

// NotifyWrongPassword limpa o campo para que o usuário possa tentar novamente.
// O orquestrador é responsável por exibir a mensagem de tentativa na barra de status.
func (m *PasswordEntryModal) NotifyWrongPassword() {
	m.field.Wipe()
}

// Render gera a representação visual do modal.
// Altura fixa: 5 linhas de corpo + 2 bordas = 7 linhas totais.
func (m *PasswordEntryModal) Render(maxHeight, maxWidth int, theme *design.Theme) string {
	innerWidth := passwordEntryWidth - 2 - 2*design.DialogPaddingH

	fieldRendered := m.field.Render(innerWidth, true, theme)

	// Body: padding + campo + padding inferior
	// Linha 0: vazia, Linha 1: label, Linha 2: área digitável, Linha 3: vazia, Linha 4: vazia
	body := "\n" + fieldRendered + "\n\n"

	confirmColor := theme.Text.Disabled
	if m.field.Len() > 0 {
		confirmColor = theme.Accent.Primary
	}

	opts := []ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "Confirmar",
			Intent: IntentConfirm,
			Action: func() tea.Cmd {
				if m.field.Len() == 0 {
					return nil
				}
				return m.onConfirm(m.field.Value())
			},
		},
		{
			Keys:   []design.Key{design.Keys.Esc},
			Label:  "Cancelar",
			Intent: IntentCancel,
			Action: func() tea.Cmd {
				m.field.Wipe()
				return m.onCancel()
			},
		},
	}

	frame := DialogFrame{
		Title:           passwordEntryTitle,
		TitleColor:      theme.Text.Primary,
		Symbol:          "",
		SymbolColor:     "",
		BorderColor:     theme.Border.Focused,
		Options:         opts,
		DefaultKeyColor: confirmColor,
		Scroll:          nil,
	}
	return frame.Render(body, passwordEntryWidth, theme)
}

// HandleKey processa eventos de teclado.
// Enter: confirma se campo não vazio. Esc: cancela e faz wipe. Tab: no-op (campo único).
func (m *PasswordEntryModal) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyEnter:
		if m.field.Len() == 0 {
			return nil
		}
		return m.onConfirm(m.field.Value())
	case tea.KeyEsc:
		m.field.Wipe()
		return m.onCancel()
	case tea.KeyTab:
		return nil // campo único — Tab não faz nada
	}
	m.field.HandleKey(msg)
	return nil
}

// Update processa mensagens Bubble Tea. Delega para HandleKey em tea.KeyMsg.
func (m *PasswordEntryModal) Update(msg tea.Msg) tea.Cmd {
	if key, ok := msg.(tea.KeyMsg); ok {
		return m.HandleKey(key)
	}
	return nil
}

// Cursor retorna a posição do cursor real para o campo de senha.
// topY e leftX são as coordenadas absolutas do canto superior esquerdo do modal na tela.
//
// Mapa de linhas do body (0-indexed):
//
//	Linha 0: vazia (padding)
//	Linha 1: label "Senha"
//	Linha 2: área digitável  ← cursor aqui
//	Linha 3: vazia
//	Linha 4: vazia
//
// Fórmula:
//
//	cursorY = topY + 1 (borda superior) + 2 (linha do field no body)
//	cursorX = leftX + 1 (borda esquerda) + DialogPaddingH + field.Len()
func (m *PasswordEntryModal) Cursor(topY, leftX int) *tea.Cursor {
	y := topY + 1 + 2
	x := leftX + 1 + design.DialogPaddingH + m.field.Len()
	return &tea.Cursor{X: x, Y: y}
}

// renderBodyLine é uma constante auxiliar para legibilidade.
// linhaDoFieldNoBody = 2 (linha da área digitável dentro do body, 0-indexed).
const _ = strings.Builder{} // força import de "strings" via body
```

> **Nota:** O `const _ = strings.Builder{}` é um workaround para manter o import `"strings"` que é usado no body string. Remova-o se o body for construído sem uso direto do pacote strings — o compilador irá avisar se o import for desnecessário.

Na prática, o import `"strings"` pode não ser necessário se o body for construído com literais `\n`. Remova o `const _` e o import se o compilador reclamar.

- [ ] **Step 4: Compilar e ajustar imports**

```
go build ./internal/tui/modal/...
```

Ajustar imports conforme erros do compilador (remover `"strings"` se não usado, etc.).

- [ ] **Step 5: Rodar testes unitários**

```
go test ./internal/tui/modal/... -run TestPasswordEntry -v
```

Resultado esperado: testes unitários (NotifyWrongPassword, Cursor) passam. Testes golden falham com "golden file not found".

- [ ] **Step 6: Gerar golden files**

```
go test ./internal/tui/modal/... -run TestPasswordEntryModal_Empty -update-golden -v
go test ./internal/tui/modal/... -run TestPasswordEntryModal_WithContent -update-golden -v
```

- [ ] **Step 7: Verificar os golden files**

Ler `internal/tui/modal/testdata/golden/password_entry-empty-50x7.golden.txt` e confirmar:
- Largura total da borda = 50 colunas
- Título `Senha mestra` centralizado entre os `──`
- Label `Senha` com 2 espaços de indentação (DialogPaddingH)
- Linha da área digitável: 44 espaços (innerWidth = 44)
- Rodapé: `Enter Confirmar` à esquerda, `Esc Cancelar` à direita

Ler `password_entry-with_content-50x7.golden.txt` e confirmar:
- Linha digitável: `••••••••` + 36 espaços

- [ ] **Step 8: Verificar o golden JSON**

Ler `password_entry-empty-50x7.golden.json` e confirmar:
- Posições da borda com cor `border.focused`
- `Enter Confirmar`: cor `text.disabled` (campo vazio)
- Ler `password_entry-with_content-50x7.golden.json`: `Enter Confirmar` com `accent.primary` + bold

- [ ] **Step 9: Rodar todos os testes**

```
go test ./internal/tui/modal/... -v
```

Resultado esperado: todos passam.

- [ ] **Step 10: Commit**

```
git add internal/tui/modal/password_entry_modal.go internal/tui/modal/password_entry_modal_test.go internal/tui/modal/testdata/
git commit -m "feat(modal): add PasswordEntryModal with cursor support"
```

---

## Task 6: Implementar PasswordCreateModal

**Files:**
- Create: `internal/tui/modal/password_create_modal.go`
- Create: `internal/tui/modal/password_create_modal_test.go`

- [ ] **Step 1: Escrever testes golden e unitários**

Criar `internal/tui/modal/password_create_modal_test.go`:

```go
package modal_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/tui/testdata"
)

func newCreateModal() *modal.PasswordCreateModal {
	mc := &stubMessageController{}
	return modal.NewPasswordCreateModal(mc,
		func(_ []byte) tea.Cmd { return nil },
		func() tea.Cmd { return nil },
	)
}

func TestPasswordCreateModal_Initial(t *testing.T) {
	m := newCreateModal()
	testdata.TestRenderManaged(t, "password_create", "initial", []string{"50x9"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestPasswordCreateModal_WithMeter(t *testing.T) {
	m := newCreateModal()
	for _, r := range "12345678" {
		m.Update(tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{r}})
	}
	testdata.TestRenderManaged(t, "password_create", "with_meter", []string{"50x11"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestPasswordCreateModal_Confirmed(t *testing.T) {
	m := newCreateModal()
	senha := "Abcdefghijk1!"
	for _, r := range senha {
		m.Update(tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{r}})
	}
	// Tab para ir ao campo Confirmação
	m.Update(tea.KeyMsg{Type: tea.KeyTab})
	for _, r := range senha {
		m.Update(tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{r}})
	}
	testdata.TestRenderManaged(t, "password_create", "confirmed", []string{"50x11"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestPasswordCreateModal_Mismatch(t *testing.T) {
	m := newCreateModal()
	for _, r := range "senha123" {
		m.Update(tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{r}})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyTab})
	for _, r := range "senhaXXX" {
		m.Update(tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{r}})
	}
	testdata.TestRenderManaged(t, "password_create", "mismatch", []string{"50x11"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

// Testes unitários de lógica

func TestPasswordCreateModal_TabSwitchesFocus(t *testing.T) {
	m := newCreateModal()
	// Foco inicial: Nova senha
	if !m.FocusedOnNew() {
		t.Error("expected initial focus on fieldNew")
	}
	m.Update(tea.KeyMsg{Type: tea.KeyTab})
	if m.FocusedOnNew() {
		t.Error("after Tab, expected focus on fieldConfirm")
	}
	m.Update(tea.KeyMsg{Type: tea.KeyTab})
	if !m.FocusedOnNew() {
		t.Error("after second Tab, expected focus back on fieldNew")
	}
}

func TestPasswordCreateModal_Cursor_FocusedNew(t *testing.T) {
	m := newCreateModal()
	c := m.Cursor(1, 0)
	if c == nil {
		t.Fatal("Cursor: expected non-nil")
	}
	// topY=1, leftX=0, borda=1, linhaFieldNew=2, paddingH=2, len=0
	// Y = 1 + 1 + 2 = 4
	// X = 0 + 1 + 2 + 0 = 3
	if c.Y != 4 {
		t.Errorf("Cursor.Y (fieldNew): expected 4, got %d", c.Y)
	}
	if c.X != 3 {
		t.Errorf("Cursor.X (fieldNew): expected 3, got %d", c.X)
	}
}

func TestPasswordCreateModal_Cursor_FocusedConfirm(t *testing.T) {
	m := newCreateModal()
	m.Update(tea.KeyMsg{Type: tea.KeyTab}) // mover foco para Confirmação
	c := m.Cursor(1, 0)
	if c == nil {
		t.Fatal("Cursor: expected non-nil")
	}
	// topY=1, leftX=0, borda=1, linhaFieldConfirm=6, paddingH=2, len=0
	// Y = 1 + 1 + 6 = 8
	// X = 0 + 1 + 2 + 0 = 3
	if c.Y != 8 {
		t.Errorf("Cursor.Y (fieldConfirm): expected 8, got %d", c.Y)
	}
	if c.X != 3 {
		t.Errorf("Cursor.X (fieldConfirm): expected 3, got %d", c.X)
	}
}
```

- [ ] **Step 2: Rodar para confirmar que falha**

```
go test ./internal/tui/modal/... -run TestPasswordCreate -v
```

Resultado esperado: erro de compilação.

- [ ] **Step 3: Criar `password_create_modal.go`**

Criar `internal/tui/modal/password_create_modal.go`:

```go
package modal

import (
	"bytes"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// passwordCreateTitle é o título do modal conforme a spec.
const passwordCreateTitle = "Definir senha mestra"

// passwordCreateWidth é a largura fixa do modal em colunas.
const passwordCreateWidth = 50

// focusedField identifica qual campo está ativo no PasswordCreateModal.
type focusedField int

const (
	fieldNew     focusedField = iota // campo "Nova senha"
	fieldConfirm                     // campo "Confirmação"
)

// PasswordCreateModal exibe o diálogo de criação/alteração de senha mestra.
// Tem dois campos: "Nova senha" e "Confirmação". O medidor de força aparece
// quando "Nova senha" não está vazio.
// Implementa tui.ModalView.
type PasswordCreateModal struct {
	mc           tui.MessageController
	fieldNew     *PasswordField
	fieldConfirm *PasswordField
	focused      focusedField
	onConfirm    func(password []byte) tea.Cmd
	onCancel     func() tea.Cmd
}

// NewPasswordCreateModal cria o modal e emite a dica inicial na barra de status.
func NewPasswordCreateModal(
	mc tui.MessageController,
	onConfirm func(password []byte) tea.Cmd,
	onCancel func() tea.Cmd,
) *PasswordCreateModal {
	m := &PasswordCreateModal{
		mc:           mc,
		fieldNew:     NewPasswordField("Nova senha"),
		fieldConfirm: NewPasswordField("Confirmação"),
		focused:      fieldNew,
		onConfirm:    onConfirm,
		onCancel:     onCancel,
	}
	mc.SetHintField("• A senha mestra protege todo o cofre — use 12+ caracteres")
	return m
}

// FocusedOnNew retorna true se o foco estiver em "Nova senha".
// Usado pelos testes para verificar o estado de foco.
func (m *PasswordCreateModal) FocusedOnNew() bool {
	return m.focused == fieldNew
}

// canConfirm retorna true quando ambos os campos estão preenchidos e seus valores são iguais.
func (m *PasswordCreateModal) canConfirm() bool {
	if m.fieldNew.Len() == 0 || m.fieldConfirm.Len() == 0 {
		return false
	}
	vNew := m.fieldNew.Value()
	vConf := m.fieldConfirm.Value()
	defer func() {
		// Wipe das cópias temporárias criadas por Value()
		for i := range vNew {
			vNew[i] = 0
		}
		for i := range vConf {
			vConf[i] = 0
		}
	}()
	return bytes.Equal(vNew, vConf)
}

// Render gera a representação visual do modal.
// O medidor de força aparece abaixo dos campos quando "Nova senha" não está vazio.
func (m *PasswordCreateModal) Render(maxHeight, maxWidth int, theme *design.Theme) string {
	innerWidth := passwordCreateWidth - 2 - 2*design.DialogPaddingH

	newFocused := m.focused == fieldNew
	confirmFocused := m.focused == fieldConfirm

	fieldNewRendered := m.fieldNew.Render(innerWidth, newFocused, theme)
	fieldConfirmRendered := m.fieldConfirm.Render(innerWidth, confirmFocused, theme)

	// Body:
	// Linha 0: vazia
	// Linha 1: label "Nova senha"
	// Linha 2: área digitável Nova senha
	// Linha 3: vazia
	// Linha 4: vazia
	// Linha 5: label "Confirmação"
	// Linha 6: área digitável Confirmação
	// Linha 7: vazia
	// [Linha 8: medidor de força]   ← só quando Nova senha não vazio
	// [Linha 9: vazia após medidor] ← só quando Nova senha não vazio
	body := "\n" + fieldNewRendered + "\n\n" + fieldConfirmRendered + "\n"

	if m.fieldNew.Len() > 0 {
		meter := RenderStrengthMeter(m.fieldNew.Value(), innerWidth, theme)
		body += "\n" + meter + "\n"
	}

	confirmColor := theme.Text.Disabled
	if m.canConfirm() {
		confirmColor = theme.Accent.Primary
	}

	opts := []ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "Confirmar",
			Intent: IntentConfirm,
			Action: func() tea.Cmd {
				if !m.canConfirm() {
					return nil
				}
				pwd := m.fieldNew.Value()
				m.fieldConfirm.Wipe()
				return m.onConfirm(pwd)
			},
		},
		{
			Keys:   []design.Key{design.Keys.Esc},
			Label:  "Cancelar",
			Intent: IntentCancel,
			Action: func() tea.Cmd {
				m.fieldNew.Wipe()
				m.fieldConfirm.Wipe()
				return m.onCancel()
			},
		},
	}

	frame := DialogFrame{
		Title:           passwordCreateTitle,
		TitleColor:      theme.Text.Primary,
		Symbol:          "",
		SymbolColor:     "",
		BorderColor:     theme.Border.Focused,
		Options:         opts,
		DefaultKeyColor: confirmColor,
		Scroll:          nil,
	}
	return frame.Render(body, passwordCreateWidth, theme)
}

// HandleKey processa eventos de teclado.
func (m *PasswordCreateModal) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyTab:
		m.switchFocus()
		return nil

	case tea.KeyEnter:
		if !m.canConfirm() {
			if m.fieldConfirm.Len() > 0 {
				m.mc.SetError("✕ As senhas não conferem — digite novamente")
			}
			return nil
		}
		pwd := m.fieldNew.Value()
		m.fieldConfirm.Wipe()
		return m.onConfirm(pwd)

	case tea.KeyEsc:
		m.fieldNew.Wipe()
		m.fieldConfirm.Wipe()
		return m.onCancel()
	}

	// Delegar para o campo focado
	if m.focused == fieldNew {
		m.fieldNew.HandleKey(msg)
	} else {
		m.fieldConfirm.HandleKey(msg)
		// Validação em tempo real no campo Confirmação
		m.validateConfirmation()
	}
	return nil
}

// switchFocus alterna o foco entre os dois campos e emite a dica adequada na barra.
func (m *PasswordCreateModal) switchFocus() {
	if m.focused == fieldNew {
		// Saindo de Nova senha para Confirmação
		m.focused = fieldConfirm
		// Verificar divergência ao entrar em Confirmação
		if m.fieldConfirm.Len() > 0 && !m.canConfirm() {
			m.mc.SetError("✕ As senhas não conferem — digite novamente")
		} else {
			m.mc.SetHintField("• Redigite a senha para confirmar")
		}
	} else {
		// Saindo de Confirmação para Nova senha com senhas divergentes
		if m.fieldConfirm.Len() > 0 && !m.canConfirm() {
			m.mc.SetError("✕ As senhas não conferem — digite novamente")
		}
		m.focused = fieldNew
		m.mc.SetHintField("• A senha mestra protege todo o cofre — use 12+ caracteres")
	}
}

// validateConfirmation emite mensagem na barra após cada tecla em Confirmação.
func (m *PasswordCreateModal) validateConfirmation() {
	if m.fieldConfirm.Len() == 0 {
		return
	}
	if m.canConfirm() {
		m.mc.SetHintField("• Redigite a senha para confirmar")
	} else {
		m.mc.SetError("✕ As senhas não conferem — digite novamente")
	}
}

// Update processa mensagens Bubble Tea. Delega para HandleKey em tea.KeyMsg.
func (m *PasswordCreateModal) Update(msg tea.Msg) tea.Cmd {
	if key, ok := msg.(tea.KeyMsg); ok {
		return m.HandleKey(key)
	}
	return nil
}

// Cursor retorna a posição do cursor real para o campo focado.
// topY e leftX são as coordenadas absolutas do canto superior esquerdo do modal na tela.
//
// Mapa de linhas do body (0-indexed, independente do medidor):
//
//	Linha 2: área digitável "Nova senha"     ← linhaDoFieldNoBody = 2
//	Linha 6: área digitável "Confirmação"    ← linhaDoFieldNoBody = 6
//
// Fórmula:
//
//	cursorY = topY + 1 (borda) + linhaDoFieldNoBody
//	cursorX = leftX + 1 (borda) + DialogPaddingH + field.Len()
func (m *PasswordCreateModal) Cursor(topY, leftX int) *tea.Cursor {
	var lineOffset int
	var fieldLen int

	if m.focused == fieldNew {
		lineOffset = 2
		fieldLen = m.fieldNew.Len()
	} else {
		lineOffset = 6
		fieldLen = m.fieldConfirm.Len()
	}

	y := topY + 1 + lineOffset
	x := leftX + 1 + design.DialogPaddingH + fieldLen
	return &tea.Cursor{X: x, Y: y}
}
```

- [ ] **Step 4: Compilar**

```
go build ./internal/tui/modal/...
```

Resultado esperado: nenhum erro.

- [ ] **Step 5: Rodar testes unitários**

```
go test ./internal/tui/modal/... -run TestPasswordCreateModal_Tab -v
go test ./internal/tui/modal/... -run TestPasswordCreateModal_Cursor -v
```

Resultado esperado: todos `PASS`.

- [ ] **Step 6: Gerar golden files**

```
go test ./internal/tui/modal/... -run TestPasswordCreateModal_Initial -update-golden -v
go test ./internal/tui/modal/... -run TestPasswordCreateModal_WithMeter -update-golden -v
go test ./internal/tui/modal/... -run TestPasswordCreateModal_Confirmed -update-golden -v
go test ./internal/tui/modal/... -run TestPasswordCreateModal_Mismatch -update-golden -v
```

- [ ] **Step 7: Verificar os golden files**

Ler `password_create-initial-50x9.golden.txt` e verificar:
- Largura total = 50 colunas
- Título `Definir senha mestra` nos `──` da borda superior
- `Nova senha` com 2 espaços de indentação
- Linha digitável de `Nova senha`: 44 espaços
- `Confirmação` com 2 espaços de indentação
- Linha digitável de `Confirmação`: 44 espaços
- Sem medidor (campo `Nova senha` vazio)
- 9 linhas totais (7 corpo + 2 bordas)

Ler `password_create-with_meter-50x11.golden.txt` e verificar:
- 11 linhas totais
- Linha `Força: ██░░░░░░░░ ⚠ Fraca` (senha `12345678` = apenas dígitos, 0 pontos = StrengthWeak)
- `Enter Confirmar` com `text.disabled` (Confirmação vazia)

Ler `password_create-confirmed-50x11.golden.txt`:
- `Força: ██████████ ✓ Forte` (senha `Abcdefghijk1!` = 4 pontos)
- `Enter Confirmar` com `accent.primary` bold

Ler `password_create-mismatch-50x11.golden.txt`:
- `Enter Confirmar` com `text.disabled` (senhas divergentes)

- [ ] **Step 8: Rodar todos os testes**

```
go test ./internal/tui/... -v
```

Resultado esperado: todos passam.

- [ ] **Step 9: Commit**

```
git add internal/tui/modal/password_create_modal.go internal/tui/modal/password_create_modal_test.go internal/tui/modal/testdata/
git commit -m "feat(modal): add PasswordCreateModal with strength meter and cursor support"
```

---

## Task 7: Verificação final

**Files:** nenhum arquivo novo

- [ ] **Step 1: Rodar todos os testes do projeto**

```
go test ./...
```

Resultado esperado: nenhum teste falhando.

- [ ] **Step 2: Verificar que o projeto compila sem warnings**

```
go build ./...
go vet ./...
```

Resultado esperado: nenhuma saída (zero erros, zero warnings).

- [ ] **Step 3: Confirmar golden files existem**

Verificar que os seguintes arquivos existem em `internal/tui/modal/testdata/golden/`:

```
password_field-empty_focused-44x2.golden.txt
password_field-empty_focused-44x2.golden.json
password_field-empty_blurred-44x2.golden.txt
password_field-empty_blurred-44x2.golden.json
password_field-content_focused-44x2.golden.txt
password_field-content_focused-44x2.golden.json
password_field-content_blurred-44x2.golden.txt
password_field-content_blurred-44x2.golden.json
password_entry-empty-50x7.golden.txt
password_entry-empty-50x7.golden.json
password_entry-with_content-50x7.golden.txt
password_entry-with_content-50x7.golden.json
password_create-initial-50x9.golden.txt
password_create-initial-50x9.golden.json
password_create-with_meter-50x11.golden.txt
password_create-with_meter-50x11.golden.json
password_create-confirmed-50x11.golden.txt
password_create-confirmed-50x11.golden.json
password_create-mismatch-50x11.golden.txt
password_create-mismatch-50x11.golden.json
strength_meter-weak-44x1.golden.txt
strength_meter-weak-44x1.golden.json
strength_meter-fair-44x1.golden.txt
strength_meter-fair-44x1.golden.json
strength_meter-strong-44x1.golden.txt
strength_meter-strong-44x1.golden.json
```

- [ ] **Step 4: Commit final de verificação (se houver ajustes)**

```
git add -A
git commit -m "chore(modal): verify all password modal golden files and tests pass"
```

---

## Notas de Implementação

### Import de `"strings"` em `password_entry_modal.go`
O body string é construído com literais `\n` — não há uso direto do pacote `"strings"`. Remova o import e o `const _` placeholder do Step 3 da Task 5.

### `canConfirm()` e wipe de cópias
`canConfirm()` chama `Value()` que retorna cópias. As cópias são zeradas com um `defer` inline dentro do método para não vazar dados sensíveis. O `bytes.Equal` é seguro mesmo com defer pois opera antes do retorno da função.

### `RenderStrengthMeter` recebe `Value()` (cópia)
Em `Render()` do `PasswordCreateModal`, `RenderStrengthMeter(m.fieldNew.Value(), ...)` recebe uma cópia do buffer. O caller não faz wipe desta cópia pois `RenderStrengthMeter` é uma função pura que não armazena a referência. Este é um trade-off intencional de segurança vs. complexidade — aceitável para a renderização (que ocorre a cada frame).

### Altura do modal `PasswordCreate`
- Sem medidor: body = 8 linhas (`\n` + 2 linhas campo1 + `\n\n` + 2 linhas campo2 + `\n`) = 8 linhas + 2 bordas = 10 linhas. Mas a spec diz 9 — revisar a contagem de `\n` se o golden file não bater.
- Com medidor: body += `\n` + 1 linha medidor + `\n` = +2 linhas = 12 total. Spec diz 11.

**Importante:** a contagem exata depende de como `DialogFrame.Render` conta as linhas do body. O body é `strings.Split(body, "\n")` — cada `\n` adiciona uma linha. Conte os `\n` no body literal e ajuste para bater com as alturas da spec (9 sem medidor, 11 com medidor).

Body sem medidor deve ter 7 linhas (indices 0–6):
```
"" (linha 0 vazia)
"Nova senha" (linha 1 - label)  
"••••..." (linha 2 - field)
"" (linha 3 vazia)
"" (linha 4 vazia)
"Confirmação" (linha 5 - label)
"••••..." (linha 6 - field)
"" (linha 7 vazia)
```
= string com 7 `\n`: `"\n" + label + "\n" + field + "\n\n" + label + "\n" + field + "\n"`

Wait — `fieldNew.Render()` já retorna `label\nfield` (1 `\n` interno). Então:
```go
body := "\n" + fieldNewRendered + "\n\n" + fieldConfirmRendered + "\n"
```
= `\n` + `label\nfield` + `\n\n` + `label\nfield` + `\n`
= linhas: `["", "label", "field", "", "", "label", "field", ""]`
= 8 linhas quando splitado por `\n` → body com 7 `\n` → 8 linhas no frame → total 10 com bordas

A spec diz 9. Ajuste: remover um `\n` final:
```go
body := "\n" + fieldNewRendered + "\n\n" + fieldConfirmRendered
```
= `["", "label", "field", "", "", "label", "field"]` = 7 linhas → total 9 com bordas ✓

Com medidor (adicionar 2 linhas):
```go
body += "\n\n" + meter
```
= +2 linhas → total 11 ✓

**Use esta versão corrigida no Step 3 da Task 6:**

```go
body := "\n" + fieldNewRendered + "\n\n" + fieldConfirmRendered

if m.fieldNew.Len() > 0 {
    meter := RenderStrengthMeter(m.fieldNew.Value(), innerWidth, theme)
    body += "\n\n" + meter
}
```

E o mapa de linhas do body (0-indexed) fica:
```
Linha 0: vazia
Linha 1: label "Nova senha"
Linha 2: área digitável Nova senha
Linha 3: vazia
Linha 4: vazia
Linha 5: label "Confirmação"
Linha 6: área digitável Confirmação
[Linha 7: vazia — quando medidor presente]
[Linha 8: medidor — quando medidor presente]
```

O cursor nas linhas 2 e 6 não muda, conforme spec.

Da mesma forma, ajuste o body do `PasswordEntryModal`:

Body deve ter 5 linhas (spec: 5 linhas de corpo):
```
Linha 0: vazia
Linha 1: label "Senha"
Linha 2: área digitável
Linha 3: vazia
Linha 4: vazia
```
= `"\n" + fieldRendered + "\n\n"` 
= `["", "label", "field", "", ""]` = 5 linhas → total 7 com bordas ✓

Verifique a contagem ao gerar os goldens e ajuste se necessário.
