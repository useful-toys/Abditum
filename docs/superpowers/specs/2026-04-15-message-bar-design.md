# Design — Barra de Mensagens

> Implementação da barra de mensagens conforme `golden/tui-design-system.md` e `golden/tui-spec-barras.md`.

## Escopo

- Tipos visuais de mensagem (`MessageKind`, `Message`, helpers, `Render`)
- Estado e renderização da barra (`MessageLineView`)
- Interface `MessageController` para controle externo
- Integração com `RootModel` via timer global e `TickMsg`

## Arquivos

| Arquivo | Conteúdo |
|---|---|
| `internal/tui/design/design_message.go` | `MessageKind`, `Message`, helpers, `Message.Render` |
| `internal/tui/screen/message_view.go` | `MessageLineView` — estado + render + implementa `MessageController` |
| `internal/tui/root.go` | Timer global, trata `TickMsg`, expõe `MessageController()` |

## Arquitetura

### Interface `MessageController`

Definida em `package tui` (evita ciclo de imports com `screen`).

```go
type MessageController interface {
    SetBusy(text string)
    SetSuccess(text string)
    SetError(text string)
    SetWarning(text string)
    SetInfo(text string)
    SetHintField(text string)
    SetHintUsage(text string)
    Clear()
}
```

Usada por ações via closure: `r.MessageController().SetSuccess("Cofre salvo")`.

### `MessageLineView` (em `screen/message_view.go`)

Struct única — estado + renderização + implementação da interface:

```go
type MessageLineView struct {
    current      design.Message
    ttl          int  // segundos restantes; 0 = sem TTL (permanente)
    spinnerFrame int
}
```

**Métodos da interface:**
- `SetXxx(text string)` — popula `current` e seta `ttl = kind.DefaultTTL()`
- `Clear()` — zera `current` e `ttl`

**Tick:**
- Se `current.Kind == MsgBusy` → avança `spinnerFrame = (spinnerFrame + 1) % 4`
- Se `ttl > 0` → decrementa; se chega a 0 → `current = design.Message{}`
- As duas operações acontecem no mesmo tick

**Render:**
- Delega para `current.Render(theme, maxWidth, spinnerFrame)` para obter conteúdo
- Monta a estrutura visual conforme spec (ver seção abaixo)

### `RootModel` (em `root.go`)

- Campo `messageLineView *screen.MessageLineView` substitui `currentMessage string`
- `SetMessage(string)` removido
- `MessageController() MessageController` exposto — retorna `r.messageLineView`
- `Init()` retorna timer global: `tea.Every(1*time.Second, func(time.Time) tea.Msg { return TickMsg{} })`
- `Update` trata `TickMsg` → `r.messageLineView.Tick()`
- `View()` chama `r.messageLineView.Render(r.width, r.theme)`

### Mensagens Bubble Tea

```go
// TickMsg é emitido 1 vez por segundo pelo timer global.
// Avança animação do spinner e decrementa TTL de mensagens temporárias.
type TickMsg struct{}
```

---

## `design/design_message.go`

### `MessageKind`

```go
type MessageKind int

const (
    MsgSuccess MessageKind = iota
    MsgInfo
    MsgWarning
    MsgError
    MsgBusy
    MsgHintField
    MsgHintUsage
)
```

### Métodos de `MessageKind`

| Método | Descrição |
|---|---|
| `Symbol() string` | Símbolo Unicode do tipo |
| `Color(theme *Theme) string` | Token de cor hex |
| `Style(theme *Theme) lipgloss.Style` | Estilo com cor + atributos |
| `DefaultTTL() int` | TTL padrão em segundos (0 = permanente) |

Tabela de valores:

| Kind | Símbolo | Token | Atributo | TTL |
|---|---|---|---|---|
| MsgSuccess | `✓` | `semantic.success` | — | 5s |
| MsgInfo | `ℹ` | `semantic.info` | — | 5s |
| MsgWarning | `⚠` | `semantic.warning` | — | 5s |
| MsgError | `✕` | `semantic.error` | **bold** | 5s |
| MsgBusy | frame do spinner | `accent.primary` | — | 0 |
| MsgHintField | `•` | `text.secondary` | *italic* | 0 |
| MsgHintUsage | `•` | `text.secondary` | *italic* | 0 |

### `Message`

```go
type Message struct {
    Kind MessageKind
    Text string
}
```

Sem campo `TTL` — o TTL é sempre `kind.DefaultTTL()` (padrão) ou sobrescrito no `SetXxx` se necessário no futuro.

### Helpers

```go
func Success(text string) Message
func Error(text string) Message
func Info(text string) Message
func Warning(text string) Message
func Busy(text ...string) Message   // texto é opcional para MsgBusy
func HintField(text string) Message
func HintUsage(text string) Message
```

### `Message.Render`

```go
func (m Message) Render(theme *Theme, maxWidth int, spinnerFrame int) (output string, columns int)
```

- Para `MsgBusy`: símbolo = `SpinnerFrame(spinnerFrame)`, sem texto
- Para demais: símbolo + `"  "` (2 espaços) + texto, truncado com `…` se necessário
- Cor aplicada via `m.Kind.Style(theme)` a símbolo e texto juntos
- Usa `lipgloss.Width()` para medir colunas — nunca `len()`
- `maxWidth` para tipos sem símbolo: `width - 6`; com símbolo: `width - 9` (conforme spec)

---

## Estrutura Visual da Barra

Conforme `golden/tui-spec-barras.md`:

**Sem mensagem:**
```
────────────────────────────────────────────────────── (largura total)
```

**Com mensagem:**
```
─── ✓  Cofre salvo ────────────────────────────────────
```

Estrutura char a char:
- `───` (3× `SymBorderH`) + ` ` (1 espaço) + conteúdo + ` ` (1 espaço) + `─`×N (mínimo 1)
- Borda sempre em `border.default`
- Conteúdo: símbolo + `"  "` + texto (para tipos com símbolo)
- Truncamento: texto truncado com `…` se excede espaço disponível

**Cálculo de largura máxima do texto** (conforme spec):
- Com símbolo: `width - 9`
- Sem símbolo (extensibilidade futura): `width - 6`
