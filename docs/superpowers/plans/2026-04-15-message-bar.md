# Message Bar Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implementar a barra de mensagens de status (`MessageLineView`) completa com suporte a tipos semânticos, spinner animado via timer global, e TTL automático.

**Architecture:** Quatro arquivos são tocados em ordem de dependência: (1) `design/design_message.go` define os tipos visuais, (2) `internal/tui/message.go` define a interface pública, (3) `screen/message_view.go` implementa estado + render, (4) `root.go` integra o timer e expõe o controller. Cada task compila e tem testes antes de avançar.

**Tech Stack:** Go, charm.land/bubbletea/v2, charm.land/lipgloss/v2

---

## File Map

| Arquivo | Ação | Responsabilidade |
|---|---|---|
| `internal/tui/design/design_message.go` | Criar | `MessageKind`, `Message`, helpers construtores, `Message.Render` |
| `internal/tui/design/design_message_test.go` | Criar | Testes de `Symbol()`, `DefaultTTL()`, `Message.Render` |
| `internal/tui/message.go` | Criar | `TickMsg`, interface `MessageController` |
| `internal/tui/screen/message_view.go` | Substituir | `MessageLineView` — estado + render + `MessageController` |
| `internal/tui/screen/message_view_test.go` | Criar | Testes de `Update`, `Render` e métodos `SetXxx` |
| `internal/tui/root.go` | Modificar | Timer em `Init`, `case TickMsg` em `Update`, campo, `MessageController()`, `View()` |

---

## Task 1: `MessageKind` e métodos — tipos visuais puros

**Files:**
- Create: `internal/tui/design/design_message.go`
- Create: `internal/tui/design/design_message_test.go`

- [ ] **Step 1.1: Escrever os testes de `MessageKind`**

Criar `internal/tui/design/design_message_test.go`:

```go
package design

import (
	"testing"
)

func TestMessageKind_DefaultTTL(t *testing.T) {
	tests := []struct {
		kind MessageKind
		want int
	}{
		{MsgSuccess, 5},
		{MsgInfo, 5},
		{MsgWarning, 5},
		{MsgError, 5},
		{MsgBusy, 0},
		{MsgHintField, 0},
		{MsgHintUsage, 0},
	}
	for _, tt := range tests {
		if got := tt.kind.DefaultTTL(); got != tt.want {
			t.Errorf("MessageKind(%d).DefaultTTL() = %d, want %d", tt.kind, got, tt.want)
		}
	}
}

func TestMessageKind_Symbol(t *testing.T) {
	tests := []struct {
		kind MessageKind
		want string
	}{
		{MsgSuccess, SymSuccess},    // "✓"
		{MsgInfo, SymInfo},          // "ℹ"
		{MsgWarning, SymWarning},    // "⚠"
		{MsgError, SymError},        // "✕"
		{MsgBusy, SpinnerFrames[0]}, // "◐"
		{MsgHintField, SymBullet},   // "•"
		{MsgHintUsage, SymBullet},   // "•"
	}
	for _, tt := range tests {
		if got := tt.kind.Symbol(); got != tt.want {
			t.Errorf("MessageKind(%d).Symbol() = %q, want %q", tt.kind, got, tt.want)
		}
	}
}

func TestMessageKind_Color_TokyoNight(t *testing.T) {
	theme := TokyoNight
	tests := []struct {
		kind MessageKind
		want string
	}{
		{MsgSuccess, theme.Semantic.Success},
		{MsgInfo, theme.Semantic.Info},
		{MsgWarning, theme.Semantic.Warning},
		{MsgError, theme.Semantic.Error},
		{MsgBusy, theme.Accent.Primary},
		{MsgHintField, theme.Text.Secondary},
		{MsgHintUsage, theme.Text.Secondary},
	}
	for _, tt := range tests {
		if got := tt.kind.Color(theme); got != tt.want {
			t.Errorf("MessageKind(%d).Color(TokyoNight) = %q, want %q", tt.kind, got, tt.want)
		}
	}
}

func TestMessageHelpers(t *testing.T) {
	tests := []struct {
		name string
		msg  Message
		kind MessageKind
		text string
	}{
		{"Success", Success("ok"), MsgSuccess, "ok"},
		{"Error", Error("fail"), MsgError, "fail"},
		{"Info", Info("note"), MsgInfo, "note"},
		{"Warning", Warning("warn"), MsgWarning, "warn"},
		{"HintField", HintField("hint"), MsgHintField, "hint"},
		{"HintUsage", HintUsage("usage"), MsgHintUsage, "usage"},
	}
	for _, tt := range tests {
		if tt.msg.Kind != tt.kind {
			t.Errorf("%s: Kind = %d, want %d", tt.name, tt.msg.Kind, tt.kind)
		}
		if tt.msg.Text != tt.text {
			t.Errorf("%s: Text = %q, want %q", tt.name, tt.msg.Text, tt.text)
		}
	}
}

func TestBusyHelper_WithText(t *testing.T) {
	msg := Busy("carregando...")
	if msg.Kind != MsgBusy {
		t.Errorf("Busy().Kind = %d, want %d", msg.Kind, MsgBusy)
	}
	if msg.Text != "carregando..." {
		t.Errorf("Busy().Text = %q, want %q", msg.Text, "carregando...")
	}
}

func TestBusyHelper_WithoutText(t *testing.T) {
	msg := Busy()
	if msg.Kind != MsgBusy {
		t.Errorf("Busy().Kind = %d, want %d", msg.Kind, MsgBusy)
	}
	if msg.Text != "" {
		t.Errorf("Busy() without text: Text = %q, want empty", msg.Text)
	}
}

func TestBusyHelper_SpinnerFrame_ZeroValue(t *testing.T) {
	msg := Busy("teste")
	if msg.SpinnerFrame != 0 {
		t.Errorf("Busy(): SpinnerFrame = %d, want 0 (zero value)", msg.SpinnerFrame)
	}
}
```

- [ ] **Step 1.2: Verificar que os testes falham (tipos não existem ainda)**

```
go test ./internal/tui/design/... 2>&1
```

Resultado esperado: erro de compilação — `MessageKind undefined`.

- [ ] **Step 1.3: Criar `internal/tui/design/design_message.go`**

```go
package design

import "charm.land/lipgloss/v2"

// MessageKind representa o tipo semântico de uma mensagem de status.
// Define o símbolo, cor e TTL padrão de cada tipo.
type MessageKind int

const (
	// MsgSuccess indica operação concluída com êxito.
	MsgSuccess MessageKind = iota
	// MsgInfo é uma mensagem informativa neutra.
	MsgInfo
	// MsgWarning indica situação que requer atenção.
	MsgWarning
	// MsgError indica falha ou estado inválido.
	MsgError
	// MsgBusy indica operação em andamento — exibe spinner animado.
	MsgBusy
	// MsgHintField exibe dica para o campo com foco atual.
	MsgHintField
	// MsgHintUsage exibe dica de uso geral do contexto.
	MsgHintUsage
)

// Symbol retorna o símbolo Unicode estático do tipo de mensagem.
// Para MsgBusy retorna SpinnerFrames[0] como placeholder — o frame real é
// lido de Message.SpinnerFrame em Message.Render.
func (k MessageKind) Symbol() string {
	switch k {
	case MsgSuccess:
		return SymSuccess
	case MsgInfo:
		return SymInfo
	case MsgWarning:
		return SymWarning
	case MsgError:
		return SymError
	case MsgBusy:
		return SpinnerFrames[0]
	default: // MsgHintField, MsgHintUsage
		return SymBullet
	}
}

// Color retorna o token de cor hex para o tipo, retirado do tema fornecido.
func (k MessageKind) Color(theme *Theme) string {
	switch k {
	case MsgSuccess:
		return theme.Semantic.Success
	case MsgInfo:
		return theme.Semantic.Info
	case MsgWarning:
		return theme.Semantic.Warning
	case MsgError:
		return theme.Semantic.Error
	case MsgBusy:
		return theme.Accent.Primary
	default: // MsgHintField, MsgHintUsage
		return theme.Text.Secondary
	}
}

// Style retorna o estilo lipgloss para o tipo, incluindo cor e atributos tipográficos.
// MsgError adiciona bold; MsgHintField e MsgHintUsage adicionam italic.
func (k MessageKind) Style(theme *Theme) lipgloss.Style {
	base := lipgloss.NewStyle().Foreground(lipgloss.Color(k.Color(theme)))
	switch k {
	case MsgError:
		return base.Bold(true)
	case MsgHintField, MsgHintUsage:
		return base.Italic(true)
	default:
		return base
	}
}

// DefaultTTL retorna o TTL padrão em ticks (segundos) para o tipo.
// 0 significa permanente — a mensagem só é removida por SetXxx ou Clear explícito.
func (k MessageKind) DefaultTTL() int {
	switch k {
	case MsgSuccess, MsgInfo, MsgWarning, MsgError:
		return 5
	default: // MsgBusy, MsgHintField, MsgHintUsage
		return 0
	}
}

// Message representa uma mensagem de status a ser exibida na barra.
// O TTL não é armazenado aqui — é gerenciado pelo MessageLineView.
// SpinnerFrame guarda o frame atual da animação quando Kind == MsgBusy;
// é irrelevante para outros tipos.
type Message struct {
	Kind         MessageKind
	Text         string
	SpinnerFrame int // frame atual da animação; só relevante quando Kind == MsgBusy
}

// Render produz a string ANSI renderizada do conteúdo da mensagem.
// maxWidth é o espaço disponível em colunas (calculado pelo caller como terminalWidth-6).
// Para MsgBusy, usa m.SpinnerFrame para determinar o frame da animação.
// Retorna a string renderizada e o número real de colunas ocupadas.
func (m Message) Render(theme *Theme, maxWidth int) (string, int) {
	style := m.Kind.Style(theme)

	var sym string
	if m.Kind == MsgBusy {
		// SpinnerFrame pode ter largura ambígua em alguns locales — medir com lipgloss.Width.
		sym = SpinnerFrame(m.SpinnerFrame)
	} else {
		sym = m.Kind.Symbol()
	}

	// MsgBusy com texto: "spinner  texto"; sem texto: apenas "spinner".
	// Demais tipos: "símbolo  texto".
	var raw string
	if m.Kind == MsgBusy && m.Text == "" {
		raw = sym
	} else {
		raw = sym + "  " + m.Text
	}

	// Truncar se necessário, garantindo que caiba em maxWidth colunas.
	// Usar lipgloss.Width para suporte correto a caracteres multibyte — nunca len().
	if lipgloss.Width(raw) > maxWidth {
		runes := []rune(raw)
		for len(runes) > 0 {
			candidate := string(runes) + SymEllipsis
			if lipgloss.Width(candidate) <= maxWidth {
				raw = candidate
				break
			}
			runes = runes[:len(runes)-1]
		}
		if len(runes) == 0 {
			raw = ""
		}
	}

	rendered := style.Render(raw)
	cols := lipgloss.Width(rendered)
	return rendered, cols
}

// Success cria uma mensagem de sucesso.
func Success(text string) Message { return Message{Kind: MsgSuccess, Text: text} }

// Error cria uma mensagem de erro.
func Error(text string) Message { return Message{Kind: MsgError, Text: text} }

// Info cria uma mensagem informativa.
func Info(text string) Message { return Message{Kind: MsgInfo, Text: text} }

// Warning cria uma mensagem de aviso.
func Warning(text string) Message { return Message{Kind: MsgWarning, Text: text} }

// Busy cria uma mensagem de operação em andamento.
// text é opcional — se fornecido, é exibido ao lado do spinner.
// SpinnerFrame é 0 (zero value) — MessageLineView avança o frame a cada TickMsg.
func Busy(text ...string) Message {
	msg := Message{Kind: MsgBusy}
	if len(text) > 0 {
		msg.Text = text[0]
	}
	return msg
}

// HintField cria uma dica para o campo com foco atual.
func HintField(text string) Message { return Message{Kind: MsgHintField, Text: text} }

// HintUsage cria uma dica de uso geral do contexto.
func HintUsage(text string) Message { return Message{Kind: MsgHintUsage, Text: text} }
```

- [ ] **Step 1.4: Verificar que os testes passam**

```
go test ./internal/tui/design/...
```

Resultado esperado: `ok  github.com/useful-toys/abditum/internal/tui/design`

- [ ] **Step 1.5: Commit**

```
git add internal/tui/design/design_message.go internal/tui/design/design_message_test.go
git commit -m "feat(design): add MessageKind, Message types and Render"
```

---

## Task 2: `Message.Render` — testes de renderização

**Files:**
- Modify: `internal/tui/design/design_message_test.go`

Testes da Task 1 cobrem helpers e metadados. Aqui adicionamos testes de `Render`.

- [ ] **Step 2.1: Adicionar testes de `Message.Render` ao arquivo de teste**

Adicionar ao final de `internal/tui/design/design_message_test.go`:

```go
import "charm.land/lipgloss/v2"
```

Adicionar ao bloco de imports existente (junto com `"testing"`), depois adicionar ao final do arquivo:

```go
func TestMessage_Render_Basic(t *testing.T) {
	theme := TokyoNight

	msg := Success("Cofre salvo")
	output, cols := msg.Render(theme, 100)

	if lipgloss.Width(output) != cols {
		t.Errorf("Render: lipgloss.Width(output) = %d, mas cols retornado = %d", lipgloss.Width(output), cols)
	}
	if cols > 100 {
		t.Errorf("Render: cols = %d, excede maxWidth = 100", cols)
	}
}

func TestMessage_Render_BusyWithText(t *testing.T) {
	theme := TokyoNight

	msg := Busy("Salvando cofre...")
	msg.SpinnerFrame = 1 // simular frame 1 da animação
	output, cols := msg.Render(theme, 100)

	if cols > 100 {
		t.Errorf("Render Busy: cols = %d, excede maxWidth = 100", cols)
	}
	if lipgloss.Width(output) != cols {
		t.Errorf("Render Busy: lipgloss.Width(output) != cols: %d != %d", lipgloss.Width(output), cols)
	}
}

func TestMessage_Render_BusyWithoutText(t *testing.T) {
	theme := TokyoNight

	msg := Busy()
	_, cols := msg.Render(theme, 100)

	if cols > 100 {
		t.Errorf("Render Busy sem texto: cols = %d, excede maxWidth = 100", cols)
	}
}

func TestMessage_Render_TruncatesLongText(t *testing.T) {
	theme := TokyoNight

	long := "Este texto é muito longo e deve ser truncado com reticências para caber na barra de mensagem da interface"
	msg := Info(long)
	_, cols := msg.Render(theme, 30)

	if cols > 30 {
		t.Errorf("Render: texto longo não foi truncado — cols = %d, maxWidth = 30", cols)
	}
}

func TestMessage_Render_SpinnerFrameVaries(t *testing.T) {
	theme := TokyoNight

	// Os 4 frames devem produzir outputs de mesma largura (frames têm mesma largura visual).
	widths := make([]int, 4)
	for i := 0; i < 4; i++ {
		msg := Busy("teste")
		msg.SpinnerFrame = i
		_, cols := msg.Render(theme, 100)
		widths[i] = cols
	}
	for i := 1; i < 4; i++ {
		if widths[i] != widths[0] {
			t.Errorf("frame %d produziu cols = %d, esperado %d (mesmo que frame 0)", i, widths[i], widths[0])
		}
	}
}
```

- [ ] **Step 2.2: Verificar que os testes passam**

```
go test ./internal/tui/design/...
```

Resultado esperado: `ok  github.com/useful-toys/abditum/internal/tui/design`

- [ ] **Step 2.3: Commit**

```
git add internal/tui/design/design_message_test.go
git commit -m "test(design): add Message.Render tests"
```

---

## Task 3: `internal/tui/message.go` — interface pública

**Files:**
- Create: `internal/tui/message.go`

Sem testes nesta task — o arquivo contém apenas declarações de tipos/interface.

- [ ] **Step 3.1: Criar `internal/tui/message.go`**

```go
package tui

// TickMsg é emitido 1 vez por segundo pelo timer global do RootModel.
// Avança a animação do spinner (MsgBusy) e decrementa o TTL de mensagens temporárias.
type TickMsg struct{}

// MessageController é a interface de controle da barra de mensagens.
// Implementada por MessageLineView em screen/. Exposta pelo RootModel via MessageController().
// Usada por ações via closure: r.MessageController().SetSuccess("Cofre salvo").
type MessageController interface {
	// SetBusy exibe spinner com texto opcional de status (ex: "Salvando...").
	// Permanente até SetXxx ou Clear explícito.
	SetBusy(text string)
	// SetSuccess exibe mensagem de sucesso. Desaparece após 5 segundos.
	SetSuccess(text string)
	// SetError exibe mensagem de erro em destaque (bold). Desaparece após 5 segundos.
	SetError(text string)
	// SetWarning exibe aviso. Desaparece após 5 segundos.
	SetWarning(text string)
	// SetInfo exibe mensagem informativa. Desaparece após 5 segundos.
	SetInfo(text string)
	// SetHintField exibe dica para o campo focado. Permanente até substituição.
	SetHintField(text string)
	// SetHintUsage exibe dica de uso geral. Permanente até substituição.
	SetHintUsage(text string)
	// Clear remove a mensagem atual imediatamente.
	Clear()
}
```

- [ ] **Step 3.2: Verificar compilação**

```
go build ./internal/tui/...
```

Resultado esperado: sem erros.

- [ ] **Step 3.3: Commit**

```
git add internal/tui/message.go
git commit -m "feat(tui): add TickMsg type and MessageController interface"
```

---

## Task 4: `screen/message_view.go` — estado + render

**Files:**
- Replace: `internal/tui/screen/message_view.go`
- Create: `internal/tui/screen/message_view_test.go`

- [ ] **Step 4.1: Escrever os testes de `MessageLineView`**

Criar `internal/tui/screen/message_view_test.go`:

```go
package screen

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

func TestMessageLineView_ZeroValue(t *testing.T) {
	var v MessageLineView
	// Zero value deve renderizar sem pânico e retornar linha de borda com largura correta.
	output := v.Render(80, design.TokyoNight)
	if lipgloss.Width(output) == 0 {
		t.Error("Render de zero value retornou string vazia, esperado linha de borda")
	}
}

func TestMessageLineView_SetSuccess(t *testing.T) {
	var v MessageLineView
	v.SetSuccess("Cofre salvo")

	if v.current.Kind != design.MsgSuccess {
		t.Errorf("SetSuccess: Kind = %d, want %d", v.current.Kind, design.MsgSuccess)
	}
	if v.current.Text != "Cofre salvo" {
		t.Errorf("SetSuccess: Text = %q, want %q", v.current.Text, "Cofre salvo")
	}
	if v.ttl != design.MsgSuccess.DefaultTTL() {
		t.Errorf("SetSuccess: ttl = %d, want %d", v.ttl, design.MsgSuccess.DefaultTTL())
	}
}

func TestMessageLineView_SetBusy_ResetsSpinner(t *testing.T) {
	var v MessageLineView
	v.current.SpinnerFrame = 3 // simular que já havia animação em curso
	v.SetBusy("Salvando...")

	if v.current.SpinnerFrame != 0 {
		t.Errorf("SetBusy deve zerar SpinnerFrame, got %d", v.current.SpinnerFrame)
	}
	if v.current.Kind != design.MsgBusy {
		t.Errorf("SetBusy: Kind = %d, want %d", v.current.Kind, design.MsgBusy)
	}
	if v.ttl != 0 {
		t.Errorf("SetBusy: ttl = %d, want 0 (permanente)", v.ttl)
	}
}

func TestMessageLineView_Clear(t *testing.T) {
	var v MessageLineView
	v.SetSuccess("algo")
	v.Clear()

	var zero design.Message
	if v.current != zero {
		t.Errorf("Clear: current = %+v, want zero value", v.current)
	}
	if v.ttl != 0 {
		t.Errorf("Clear: ttl = %d, want 0", v.ttl)
	}
}

func TestMessageLineView_SetWarning(t *testing.T) {
	var v MessageLineView
	v.SetWarning("atenção")
	if v.current.Kind != design.MsgWarning {
		t.Errorf("SetWarning: Kind = %d, want %d", v.current.Kind, design.MsgWarning)
	}
	if v.ttl != 5 {
		t.Errorf("SetWarning: ttl = %d, want 5", v.ttl)
	}
}

func TestMessageLineView_SetError(t *testing.T) {
	var v MessageLineView
	v.SetError("falha")
	if v.current.Kind != design.MsgError {
		t.Errorf("SetError: Kind = %d, want %d", v.current.Kind, design.MsgError)
	}
	if v.ttl != 5 {
		t.Errorf("SetError: ttl = %d, want 5", v.ttl)
	}
}

func TestMessageLineView_SetInfo(t *testing.T) {
	var v MessageLineView
	v.SetInfo("info")
	if v.current.Kind != design.MsgInfo {
		t.Errorf("SetInfo: Kind = %d, want %d", v.current.Kind, design.MsgInfo)
	}
}

func TestMessageLineView_SetHintField(t *testing.T) {
	var v MessageLineView
	v.SetHintField("pressione Tab")
	if v.current.Kind != design.MsgHintField {
		t.Errorf("SetHintField: Kind = %d, want %d", v.current.Kind, design.MsgHintField)
	}
	if v.ttl != 0 {
		t.Errorf("SetHintField: ttl = %d, want 0 (permanente)", v.ttl)
	}
}

func TestMessageLineView_SetHintUsage(t *testing.T) {
	var v MessageLineView
	v.SetHintUsage("use ctrl+s para salvar")
	if v.current.Kind != design.MsgHintUsage {
		t.Errorf("SetHintUsage: Kind = %d, want %d", v.current.Kind, design.MsgHintUsage)
	}
	if v.ttl != 0 {
		t.Errorf("SetHintUsage: ttl = %d, want 0 (permanente)", v.ttl)
	}
}

func TestMessageLineView_Update_UnknownMsg(t *testing.T) {
	var v MessageLineView
	v.SetSuccess("algo")
	initialTTL := v.ttl

	type unknownMsg struct{}
	cmd := v.Update(unknownMsg{})

	if cmd != nil {
		t.Error("Update com msg desconhecida deve retornar nil cmd")
	}
	if v.ttl != initialTTL {
		t.Errorf("Update com msg desconhecida não deve alterar ttl: got %d, want %d", v.ttl, initialTTL)
	}
}

func TestMessageLineView_Render_WithMessage(t *testing.T) {
	var v MessageLineView
	v.SetSuccess("Cofre salvo")
	output := v.Render(80, design.TokyoNight)

	w := lipgloss.Width(output)
	if w != 80 {
		t.Errorf("Render: largura = %d, want 80", w)
	}
}

func TestMessageLineView_Render_ZeroValue_Width(t *testing.T) {
	var v MessageLineView
	output := v.Render(80, design.TokyoNight)

	w := lipgloss.Width(output)
	if w != 80 {
		t.Errorf("Render zero value: largura = %d, want 80", w)
	}
}

func TestMessageLineView_Render_ReturnsNoNewline(t *testing.T) {
	var v MessageLineView
	v.SetInfo("teste")
	output := v.Render(80, design.TokyoNight)

	for _, r := range output {
		if r == '\n' {
			t.Error("Render não deve conter newline — barra é linha única")
			break
		}
	}
}

// _testTick é um helper de teste que dispara a lógica de TickMsg diretamente,
// sem precisar importar package tui (o que causaria import cycle).
func (v *MessageLineView) _testTick() tea.Cmd {
	return v.tick()
}

func TestMessageLineView_SpinnerAdvances(t *testing.T) {
	var v MessageLineView
	v.SetBusy("carregando")

	for i := 1; i <= 8; i++ {
		v._testTick()
		want := i % 4
		if v.current.SpinnerFrame != want {
			t.Errorf("após %d ticks: SpinnerFrame = %d, want %d", i, v.current.SpinnerFrame, want)
		}
	}
}

func TestMessageLineView_TTL_Decrements(t *testing.T) {
	var v MessageLineView
	v.SetSuccess("ok") // ttl = 5

	for i := 4; i >= 1; i-- {
		v._testTick()
		if v.ttl != i {
			t.Errorf("após tick: ttl = %d, want %d", v.ttl, i)
		}
	}
	// Último tick: ttl chega a 0, mensagem é zerada.
	v._testTick()
	var zero design.Message
	if v.current != zero {
		t.Errorf("após ttl=0: current = %+v, want zero value", v.current)
	}
	if v.ttl != 0 {
		t.Errorf("após ttl=0: ttl = %d, want 0", v.ttl)
	}
}

func TestMessageLineView_BusyTTL_NeverExpires(t *testing.T) {
	var v MessageLineView
	v.SetBusy("operando")

	for i := 0; i < 10; i++ {
		v._testTick()
	}
	// Kind ainda deve ser MsgBusy — ttl=0 significa permanente.
	if v.current.Kind != design.MsgBusy {
		t.Errorf("MsgBusy não deve expirar: Kind = %d, want %d", v.current.Kind, design.MsgBusy)
	}
}

func TestMessageLineView_HintField_NeverExpires(t *testing.T) {
	var v MessageLineView
	v.SetHintField("pressione Enter")

	for i := 0; i < 10; i++ {
		v._testTick()
	}
	if v.current.Kind != design.MsgHintField {
		t.Errorf("MsgHintField não deve expirar: Kind = %d, want %d", v.current.Kind, design.MsgHintField)
	}
}
```

- [ ] **Step 4.2: Verificar que os testes falham (MessageLineView ainda é stub)**

```
go test ./internal/tui/screen/... 2>&1
```

Resultado esperado: erro de compilação — `v.current`, `v.ttl`, `v.tick()` não existem.

- [ ] **Step 4.3: Substituir `internal/tui/screen/message_view.go` completo**

```go
package screen

import (
	"reflect"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// MessageLineView gerencia estado e renderização da barra de mensagens de status.
// Implementa tui.MessageController via pointer receivers.
// O zero value é válido — representa estado sem mensagem ativa.
type MessageLineView struct {
	// current é a mensagem atualmente exibida. Zero value = sem mensagem.
	// Para MsgBusy, current.SpinnerFrame guarda o frame atual da animação.
	current design.Message
	// ttl é o número de ticks restantes antes da mensagem ser removida.
	// 0 significa permanente (MsgBusy, MsgHintField, MsgHintUsage nunca expiram).
	ttl int
}

// SetBusy exibe o spinner com texto opcional de status.
// design.Busy() retorna zero value de Message — SpinnerFrame começa em 0.
func (v *MessageLineView) SetBusy(text string) {
	v.current = design.Busy(text)
	v.ttl = 0
}

// SetSuccess exibe mensagem de sucesso. Desaparece após 5 ticks.
func (v *MessageLineView) SetSuccess(text string) {
	v.current = design.Success(text)
	v.ttl = v.current.Kind.DefaultTTL()
}

// SetError exibe mensagem de erro em bold. Desaparece após 5 ticks.
func (v *MessageLineView) SetError(text string) {
	v.current = design.Error(text)
	v.ttl = v.current.Kind.DefaultTTL()
}

// SetWarning exibe mensagem de aviso. Desaparece após 5 ticks.
func (v *MessageLineView) SetWarning(text string) {
	v.current = design.Warning(text)
	v.ttl = v.current.Kind.DefaultTTL()
}

// SetInfo exibe mensagem informativa. Desaparece após 5 ticks.
func (v *MessageLineView) SetInfo(text string) {
	v.current = design.Info(text)
	v.ttl = v.current.Kind.DefaultTTL()
}

// SetHintField exibe dica para o campo focado. Permanente até substituição.
func (v *MessageLineView) SetHintField(text string) {
	v.current = design.HintField(text)
	v.ttl = 0
}

// SetHintUsage exibe dica de uso geral. Permanente até substituição.
func (v *MessageLineView) SetHintUsage(text string) {
	v.current = design.HintUsage(text)
	v.ttl = 0
}

// Clear remove a mensagem atual imediatamente.
func (v *MessageLineView) Clear() {
	v.current = design.Message{}
	v.ttl = 0
}

// tick executa a lógica de um TickMsg: avança spinner e decrementa TTL.
// Separado de Update para ser testável sem importar package tui.
func (v *MessageLineView) tick() tea.Cmd {
	if v.current.Kind == design.MsgBusy {
		v.current.SpinnerFrame = (v.current.SpinnerFrame + 1) % 4
	}
	if v.ttl > 0 {
		v.ttl--
		if v.ttl == 0 {
			v.current = design.Message{}
		}
	}
	return nil
}

// Update trata mensagens do Bubble Tea. Delega tui.TickMsg para tick().
// Usa reflect para identificar o tipo e evitar import cycle com package tui.
func (v *MessageLineView) Update(msg tea.Msg) tea.Cmd {
	if reflect.TypeOf(msg).String() == "tui.TickMsg" {
		return v.tick()
	}
	return nil
}

// Render produz a string ANSI da barra de mensagens com exatamente `width` colunas.
//
// Estrutura:
//
//	Sem mensagem: ────────────────── (SymBorderH × width)
//	Com mensagem: ─── ✓  texto ─────
func (v *MessageLineView) Render(width int, theme *design.Theme) string {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Border.Default))

	// Sem mensagem ativa — linha de borda completa.
	if v.current == (design.Message{}) {
		return borderStyle.Render(strings.Repeat(design.SymBorderH, width))
	}

	// Prefixo: 3× SymBorderH + 1 espaço = 4 colunas fixas.
	prefix := borderStyle.Render(strings.Repeat(design.SymBorderH, 3)) + " "

	// Espaço disponível para o conteúdo.
	// 4 = prefixo, 2 = sufixo mínimo (1 espaço + 1× SymBorderH).
	maxWidth := width - 6
	if maxWidth < 1 {
		maxWidth = 1
	}

	content, contentCols := v.current.Render(theme, maxWidth)

	// Sufixo: espaço separador + SymBorderH repetido até preencher.
	sufixoCols := width - 4 - contentCols - 1
	if sufixoCols < 1 {
		sufixoCols = 1
	}
	suffix := " " + borderStyle.Render(strings.Repeat(design.SymBorderH, sufixoCols))

	return prefix + content + suffix
}
```

- [ ] **Step 4.4: Verificar que os testes passam**

```
go test ./internal/tui/screen/...
```

Resultado esperado: `ok  github.com/useful-toys/abditum/internal/tui/screen`

- [ ] **Step 4.5: Commit**

```
git add internal/tui/screen/message_view.go internal/tui/screen/message_view_test.go
git commit -m "feat(screen): implement MessageLineView with state, spinner, TTL and render"
```

---

## Task 5: Integração em `root.go`

**Files:**
- Modify: `internal/tui/root.go`

- [ ] **Step 5.1: Remover `currentMessage` e `SetMessage`**

Em `internal/tui/root.go`, remover o campo:

```go
// currentMessage is the current status message, displayed in messageLineView.
currentMessage string
```

E remover o método:

```go
// SetMessage defines the status message displayed in the bottom bar.
// Pass empty string to clear the message.
func (r *RootModel) SetMessage(msg string) {
	r.currentMessage = msg
}
```

- [ ] **Step 5.2: Atualizar `Init()` para iniciar o timer global**

Substituir:

```go
// Init is called once at application startup. No initial commands.
func (r *RootModel) Init() tea.Cmd {
	return nil
}
```

Por:

```go
// Init é chamado uma vez na inicialização. Inicia o timer global de 1s para TickMsg.
func (r *RootModel) Init() tea.Cmd {
	ticker := tea.Every(1*time.Second, func(time.Time) tea.Msg {
		return TickMsg{}
	})
	return tea.Batch(ticker)
}
```

- [ ] **Step 5.3: Adicionar `case TickMsg` em `Update`**

No método `Update`, logo após `case tea.WindowSizeMsg:` (antes de `case OpenModalMsg:`), adicionar:

```go
case TickMsg:
    return r, r.messageLineView.Update(msg)
```

- [ ] **Step 5.4: Adicionar método `MessageController()`**

Após o método `Manager()`, adicionar:

```go
// MessageController retorna a interface de controle da barra de mensagens.
// Usada por ações para exibir feedback: r.MessageController().SetSuccess("Salvo").
func (r *RootModel) MessageController() MessageController {
	return &r.messageLineView
}
```

- [ ] **Step 5.5: Atualizar chamada em `View()`**

Substituir em `View()`:

```go
r.messageLineView.Render(design.MessageHeight, r.width, r.theme, r.currentMessage),
```

Por:

```go
r.messageLineView.Render(r.width, r.theme),
```

- [ ] **Step 5.6: Verificar compilação e testes**

```
go build ./...
go test ./...
```

Resultado esperado:
- `go build ./...`: sem erros
- `go test ./...`: todos os pacotes passam (ou `[no test files]`)

- [ ] **Step 5.7: Commit**

```
git add internal/tui/root.go
git commit -m "feat(tui): integrate MessageLineView timer, TickMsg and MessageController in RootModel"
```

---

## Self-Review

### Spec Coverage

| Requisito da spec | Task que implementa |
|---|---|
| `TickMsg` struct em `package tui` | Task 3 |
| `MessageController` interface com 8 métodos | Task 3 |
| `MessageKind` com 7 valores iota | Task 1 |
| `Symbol()`, `Color()`, `Style()`, `DefaultTTL()` | Task 1 |
| `Message` struct com `SpinnerFrame int` | Task 1 |
| `Message` sem campo TTL | Task 1 |
| Helpers `Success`, `Error`, `Info`, `Warning`, `Busy`, `HintField`, `HintUsage` | Task 1 |
| `Message.Render(theme, maxWidth)` — sem parâmetro spinnerFrame | Task 1 + Task 2 |
| `MsgBusy.Symbol()` retorna `SpinnerFrames[0]` | Task 1 |
| `Message.Render` usa `m.SpinnerFrame` para MsgBusy | Task 1 |
| `Message.Render` usa `lipgloss.Width()` nunca `len()` | Task 1 |
| `MessageLineView` zero value válido | Task 4 |
| `MessageLineView` sem campo `spinnerFrame` separado | Task 4 |
| `SetBusy` produz zero value de `Message` (SpinnerFrame=0 implícito) | Task 4 |
| `tick()` avança `current.SpinnerFrame` | Task 4 |
| `Update(TickMsg)` decrementa TTL | Task 4 |
| `ttl=0` não expira (MsgBusy, HintField, HintUsage) | Task 4 |
| `Render` sem mensagem = linha de borda completa | Task 4 |
| `Render` chama `v.current.Render(theme, maxWidth)` sem spinnerFrame | Task 4 |
| `maxWidth = width - 6` calculado pelo caller | Task 4 |
| `Init()` retorna `tea.Batch(ticker)` | Task 5 |
| `case TickMsg` em `Update` | Task 5 |
| `MessageController()` método em `RootModel` | Task 5 |
| Remover `currentMessage` e `SetMessage` | Task 5 |
| `View()` chama `Render(r.width, r.theme)` | Task 5 |

### Placeholder Scan

Nenhum "TBD", "TODO" ou descrição sem código encontrada.

### Type Consistency

- `design.Message` com `SpinnerFrame int` — definida em Task 1, usada em testes da Task 4 (`v.current.SpinnerFrame`) ✓
- `Message.Render(theme, maxWidth)` — assinatura sem spinnerFrame usada consistentemente em Task 1, Task 2, Task 4 ✓
- `tick()` incrementa `v.current.SpinnerFrame` — nenhuma referência a `v.spinnerFrame` separado ✓
- `MessageController` interface — 8 métodos definidos em Task 3; `MessageLineView` implementa todos em Task 4 ✓
- `MessageController()` retorna `&r.messageLineView` (pointer) — satisfaz interface com pointer receivers ✓
