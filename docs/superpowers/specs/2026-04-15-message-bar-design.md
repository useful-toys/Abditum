# Design — Barra de Mensagens

> Implementação da barra de mensagens conforme `golden/tui-spec-barras.md`.

## Escopo

- Renderização visual do componente `MessageLineView`
- Sistema de mensagens (`MessageKind`, `Message`, helpers)
- Subsistema `MessageBar` dedicado
- Integração com `RootModel`

## Arquivos

| Arquivo | Conteúdo |
|---|---|
| `internal/tui/design/design_message.go` | `MessageKind`, `Message`, helpers, `Render()` |
| `internal/tui/screen/message_view.go` | `MessageLineView.Render()` |
| `internal/tui/subsystem/message_bar.go` | Subsistema `MessageBar` |
| `internal/tui/root.go` | Integração com `MessageBar`, timer TTL |

## Arquitetura

### MessageBar (subsistema dedicado)

Responsável exclusivamente por manter o estado da barra de mensagens.

```go
type MessageBar struct {
    current      design.Message
    spinnerFrame int
}

func NewMessageBar() *MessageBar
```

**API:**
```go
// Estado
func (mb *MessageBar) Current() design.Message

// Controles
func (mb *MessageBar) Clear()
func (mb *MessageBar) SetBusy(text string)
func (mb *MessageBar) Success(text string)
func (mb *MessageBar) Error(text string)
func (mb *MessageBar) Warning(text string)
func (mb *MessageBar) Info(text string)
func (mb *MessageBar) HintField(text string)
func (mb *MessageBar) HintUsage(text string)

// Renderização
func (mb *MessageBar) Render(theme *design.Theme, width int) string

// Tick (para animação do spinner)
func (mb *MessageBar) Tick()
```

Nota: `Success`, `Error`, `Warning`, `Info` substituem automaticamente qualquer mensagem anterior (incluindo spinner).

### Timer no RootModel

O TTL é gerenciado pelo RootModel, não pelo MessageBar.

```go
type RootModel struct {
    // ...
    messageBar  *MessageBar
    msgTimer    *time.Timer
}
```

```go
func (r *RootModel) SetMessage(msg design.Message) {
    r.messageBar.SetMessage(msg)
    
    if r.msgTimer != nil {
        r.msgTimer.Stop()
    }
    
    ttl := msg.TTL
    if ttl == 0 {
        ttl = msg.Kind.DefaultTTL()
    }
    if ttl > 0 {
        r.msgTimer = time.AfterFunc(ttl, func() {
            r.Send(ClearMessageMsg{})
        })
    }
}

func (r *RootModel) ClearMessage() {
    r.messageBar.Clear()
    if r.msgTimer != nil {
        r.msgTimer.Stop()
        r.msgTimer = nil
    }
}
```

No `Update`:
```go
case ClearMessageMsg:
    r.messageBar.Clear()
    return r, nil
```

### Tick

Para animação do spinner:
```go
case TickMsg:
    r.messageBar.Tick()
    return r, nil
```

## 1. design_message.go

### MessageKind

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

| Kind | Símbolo | Token | Atributo |
|---|---|---|---|
| MsgSuccess | ✓ | semantic.success | — |
| MsgInfo | ℹ | semantic.info | — |
| MsgWarning | ⚠ | semantic.warning | — |
| MsgError | ✕ | semantic.error | bold |
| MsgBusy | ◐ | accent.primary | — |
| MsgHintField | • | text.secondary | italic |
| MsgHintUsage | • | text.secondary | italic |

### MessageKind Métodos

```go
func (k MessageKind) Symbol() string
func (k MessageKind) Color(theme *Theme) string
func (k MessageKind) Style(theme *Theme) lipgloss.Style
func (k MessageKind) DefaultTTL() time.Duration
```

### Message

```go
type Message struct {
    Kind MessageKind
    Text string
    TTL  time.Duration
}
```

### Helpers

```go
func Success(text string, ttl ...time.Duration) Message
func Error(text string, ttl ...time.Duration) Message
func Info(text string, ttl ...time.Duration) Message
func Warning(text string, ttl ...time.Duration) Message
func Busy(text ...string) Message
func HintField(text string) Message
func HintUsage(text string) Message
```

Nota: `Busy(text ...string)` permite texto opcional de descrição.

### Message.Render

```go
// Retorna símbolo + "  " + texto com cor aplicada, truncado se maxWidth excedido.
// Retorna a string ANSI e o comprimento real em colunas.
func (m Message) Render(theme *Theme, maxWidth int) (output string, columns int)
```

Implementação:
- Para `MsgBusy`: retorna `SpinnerFrame(mb.spinnerFrame)` (sem texto)
- Para demais: símbolo + "  " + texto com cor, truncado se necessário
- Usa `Style()` do Kind para cor + atributos

### MessageLineView.Render

Usa o `MessageBar` para renderizar:

```go
func (v *MessageLineView) Render(height, width int, theme *design.Theme, mb *MessageBar) string {
    return mb.Render(theme, width)
}
```

## TTLs Padrão

| Kind | TTL |
|---|---|
| MsgSuccess | 5s |
| MsgInfo | 5s |
| MsgWarning | 5s |
| MsgError | 5s |
| MsgBusy | — (manual) |
| MsgHintField | permanente |
| MsgHintUsage | permanente |

## Mensagens Bubble Tea

```go
type ClearMessageMsg struct{}

type TickMsg struct{}
```

## Constantes

`design.go` já possui:
```go
const HeaderHeight = 2
const MessageHeight = 1
const ActionHeight = 1
```