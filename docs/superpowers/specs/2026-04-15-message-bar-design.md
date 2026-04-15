# Design — Barra de Mensagens

> Implementação da barra de mensagens conforme `golden/tui-design-system.md` e `golden/tui-spec-barras.md`.

## Escopo

- Tipos visuais de mensagem (`MessageKind`, `Message`, helpers, `Render`)
- Estado e renderização da barra (`MessageLineView`)
- Interface `MessageController` e tipo `TickMsg` em `package tui`
- Integração com `RootModel` via timer global e `TickMsg`

## Arquivos

| Arquivo | Conteúdo |
|---|---|
| `internal/tui/design/design_message.go` | `MessageKind`, `Message`, helpers, `Message.Render` |
| `internal/tui/screen/message_view.go` | `MessageLineView` — estado + render + implementa `MessageController` |
| `internal/tui/message.go` | Interface `MessageController` e tipo `TickMsg` |
| `internal/tui/root.go` | Timer global em `Init()`, trata `TickMsg`, expõe `MessageController()` |

## Migração do código atual

`MessageLineView` hoje é um stub stateless com assinatura:

```go
func (v *MessageLineView) Render(height, width int, theme *design.Theme, message string) string
```

Após esta implementação, a assinatura muda para:

```go
func (v *MessageLineView) Render(width int, theme *design.Theme) string
```

O campo `currentMessage string` em `RootModel` é removido. O método `SetMessage(string)` em `RootModel` é removido. A chamada em `root.go:View()` passa de `r.messageLineView.Render(design.MessageHeight, r.width, r.theme, r.currentMessage)` para `r.messageLineView.Render(r.width, r.theme)`.

---

## `internal/tui/message.go`

```go
package tui

// TickMsg é emitido 1 vez por segundo pelo timer global do RootModel.
// Avança a animação do spinner (MsgBusy) e decrementa o TTL de mensagens temporárias.
type TickMsg struct{}

// MessageController é a interface de controle da barra de mensagens.
// Implementada por MessageLineView em screen/. Exposta pelo RootModel.
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
| `Style(theme *Theme) lipgloss.Style` | Estilo com cor + atributos tipográficos |
| `DefaultTTL() int` | TTL padrão em segundos; 0 = permanente |

Tabela de valores:

| Kind | Símbolo | Token | Atributo | TTL |
|---|---|---|---|---|
| MsgSuccess | `✓` (`SymSuccess`) | `semantic.success` | — | 5 |
| MsgInfo | `ℹ` (`SymInfo`) | `semantic.info` | — | 5 |
| MsgWarning | `⚠` (`SymWarning`) | `semantic.warning` | — | 5 |
| MsgError | `✕` (`SymError`) | `semantic.error` | **bold** | 5 |
| MsgBusy | `SpinnerFrames[0]` (`◐`) via `Symbol()` | `accent.primary` | — | 0 |
| MsgHintField | `•` (`SymBullet`) | `text.secondary` | *italic* | 0 |
| MsgHintUsage | `•` (`SymBullet`) | `text.secondary` | *italic* | 0 |

> `Symbol()` retorna o símbolo estático do tipo. Para `MsgBusy`, retorna `SpinnerFrames[0]`
> (`"◐"`) como placeholder — o frame real da animação só aparece em `Message.Render` via
> `spinnerFrame`. `MsgBusy`, `MsgHintField` e `MsgHintUsage` têm `DefaultTTL() == 0` — são permanentes. Só são
> limpos por substituição explícita (outro `SetXxx` ou `Clear()`). O `Update(TickMsg)` nunca os limpa.

### `Message`

```go
type Message struct {
    Kind         MessageKind
    Text         string
    SpinnerFrame int  // frame atual da animação; só relevante quando Kind == MsgBusy
}
```

Sem campo `TTL` na struct — o TTL é sempre `kind.DefaultTTL()`, gerenciado pelo `MessageLineView`.

`SpinnerFrame` é parte do estado de `Message` para que `Render` não precise receber o frame como parâmetro. O zero value (`0`) é o frame inicial correto — `SetBusy` produz `design.Busy()` que retorna zero value de `Message`, zerando o campo implicitamente.

### Helpers

```go
func Success(text string) Message
func Error(text string) Message
func Info(text string) Message
func Warning(text string) Message
func Busy(text ...string) Message   // texto é opcional; exibido ao lado do spinner se fornecido
func HintField(text string) Message
func HintUsage(text string) Message
```

### `Message.Render`

```go
func (m Message) Render(theme *Theme, maxWidth int) (output string, columns int)
```

**Contrato de `maxWidth`:** o caller (`MessageLineView.Render`) é responsável por calcular o espaço
disponível e passá-lo aqui. A fórmula usada pelo caller é:

```
maxWidth = terminalWidth - 4 - 2
         = terminalWidth - 6
// 4 = prefixo (3× SymBorderH + 1 espaço)
// 2 = sufixo mínimo (1 espaço + 1× SymBorderH)
```

`Message.Render` garante que `símbolo + "  " + texto` (ou `spinner + "  " + texto` para MsgBusy)
caiba em `maxWidth`, truncando o texto com `SymEllipsis` se necessário. Retorna a string ANSI
renderizada e o número real de colunas ocupadas.

**Regras de renderização:**

- `MsgBusy`: símbolo = `SpinnerFrame(m.SpinnerFrame)`. Se `Text` não é vazio, exibe `spinner + "  " + texto`. Se vazio, exibe só o spinner. O frame do spinner pode ter largura ambígua em alguns locales — usar `lipgloss.Width(SpinnerFrame(m.SpinnerFrame))` para medir, nunca `len()`.
- Demais tipos: símbolo fixo + `"  "` + texto.
- Cor aplicada via `m.Kind.Style(theme)` ao bloco inteiro (símbolo e texto com a mesma cor).
- Usa `lipgloss.Width()` para medir colunas em toda medição — nunca `len()`.

---

## `screen/message_view.go` — `MessageLineView`

Struct única — estado + renderização + implementação de `MessageController`:

```go
type MessageLineView struct {
    current      design.Message
    ttl          int  // segundos restantes; 0 = sem TTL (permanente)
}
```

O zero value de `MessageLineView` é um estado válido — sem mensagem ativa, `ttl = 0`. Não requer construtor.

O `SpinnerFrame` vive dentro de `current.SpinnerFrame` — não há campo separado em `MessageLineView`.

### Métodos da interface `MessageController`

- `SetBusy(text string)` — popula `current` com `design.Busy(text)`, seta `ttl = 0`. `current.SpinnerFrame` é 0 implicitamente (zero value de `Message`).
- `SetXxx(text string)` (demais tipos) — popula `current`, seta `ttl = current.Kind.DefaultTTL()`
- `Clear()` — zera `current` (zero value de `design.Message`) e `ttl = 0`

### Update

```go
func (v *MessageLineView) Update(msg tea.Msg) tea.Cmd
```

Chamado pelo `RootModel` via delegação. Trata `TickMsg`:

1. Se `current.Kind == MsgBusy` → avança `current.SpinnerFrame = (current.SpinnerFrame + 1) % 4`
2. Se `ttl > 0` → decrementa `ttl`; se `ttl` chega a 0 → `current = design.Message{}`

As duas operações acontecem no mesmo `TickMsg`. Se não há mensagem ativa (`current` é zero value), `Update` é no-op — `ttl` é 0, `Kind` é `MsgSuccess` (iota 0, não `MsgBusy`), nenhuma condição dispara. Retorna sempre `nil`.

### Render

```go
func (v *MessageLineView) Render(width int, theme *design.Theme) string
```

Estrutura visual (conforme `golden/tui-spec-barras.md`):

**Sem mensagem (`current` é zero value):**
```
────────────────────────────────────────  (SymBorderH × width)
```

**Com mensagem:**
```
─── ✓  Cofre salvo ────────────────────
```

Composição:
1. `prefixo` = `borderStyle.Render(SymBorderH × 3)` + `" "` — 4 colunas fixas
2. `maxWidth` = `width - 6` (4 do prefixo + 2 do sufixo mínimo)
3. `content, contentCols` = `v.current.Render(theme, maxWidth)`
4. `sufixoCols` = `width - 4 - contentCols - 1` (o `- 1` é o espaço antes do sufixo)
5. resultado = `prefixo + content + " " + borderStyle.Render(SymBorderH × sufixoCols)`

`borderStyle` aplica `border.default` às bordas. O conteúdo carrega sua própria cor via `Message.Render`.

---

## `root.go` — integração

### Campo substituído

```go
// Antes:
currentMessage string

// Depois:
messageLineView screen.MessageLineView  // value type, como os demais campos de view
```

### NewRootModel — inicialização

`messageLineView` é inicializado em `NewRootModel()`:

```go
m := &RootModel{
    theme:           design.TokyoNight,
    workArea:        WorkAreaWelcome,
    version:         "dev",
    messageLineView: screen.MessageLineView{},  // zero value é estado válido (sem mensagem)
}
```

O zero value de `MessageLineView` é um estado válido — sem mensagem ativa, `ttl = 0`, `spinnerFrame = 0`.

### Init

```go
func (r *RootModel) Init() tea.Cmd {
    ticker := tea.Every(1*time.Second, func(time.Time) tea.Msg {
        return TickMsg{}
    })
    return tea.Batch(ticker) // usar Batch para acomodar futuros Cmds iniciais
}
```

### Update — novo case

```go
case TickMsg:
    return r, r.messageLineView.Update(msg)
```

### MessageController exposto

```go
func (r *RootModel) MessageController() MessageController {
    return &r.messageLineView  // ponteiro para satisfazer interface com pointer receivers
}
```

### View — chamada atualizada

```go
// Antes:
r.messageLineView.Render(design.MessageHeight, r.width, r.theme, r.currentMessage)

// Depois:
r.messageLineView.Render(r.width, r.theme)
```
