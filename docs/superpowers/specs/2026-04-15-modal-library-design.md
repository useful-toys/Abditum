# Design — Biblioteca de Diálogos (`internal/tui/modal`)

**Data:** 2026-04-15  
**Revisado:** 2026-04-15 (pós-review)  
**Escopo:** Fundação reutilizável para renderização e tratamento de eventos de diálogos modais + reescrita de `ConfirmModal` e `HelpModal` como demonstrações.  
**Referências:**
- [`golden/tui-design-system.md`](../../../golden/tui-design-system.md)
- [`golden/tui-spec-dialogos.md`](../../../golden/tui-spec-dialogos.md)
- [`golden/tui-spec-dialog-help.md`](../../../golden/tui-spec-dialog-help.md)
- [`golden/tui-spec-dialog-senha.md`](../../../golden/tui-spec-dialog-senha.md)

---

## Contexto

O pacote `internal/tui/modal/` já contém `ConfirmModal` e `HelpModal` como stubs mínimos. Eles usam `lipgloss.Border()` diretamente e ignoram toda a anatomia visual definida no design system: título com símbolo, rodapé com ações posicionadas, scroll na borda lateral, severidade, tokens de cor. Esta biblioteca substitui os stubs por implementações corretas e fornece a fundação reutilizável para todos os diálogos futuros (PasswordEntry, PasswordCreate, FilePicker etc.).

O pacote `modal/` coexiste com os stubs durante o desenvolvimento — a migração é feita ao reescrever os arquivos concretos, sem alterar a interface `tui.ModalView`.

---

## Decisões de design

| Decisão | Escolha | Justificativa |
|---|---|---|
| Padrão arquitetural | `DialogFrame` como struct renderizadora sem estado (Abordagem A) | Alinha com `ActionLineView`, `MessageLineView` — renderizadores puros no projeto |
| Localização | Expandir `internal/tui/modal/` | Não cria novo subpacote; mantém tudo coeso |
| Cores/símbolos | Nunca hardcoded — sempre via `design.Theme` ou constantes de `design/symbols.go` | Garante suporte a múltiplos temas e NO_COLOR |
| Assinatura de helpers | `func(...) (string, int)` — string ANSI + largura em colunas | Padrão do projeto (`RenderedAction` em `design_action.go`) |
| Scroll | `ScrollState` como struct separada em `modal/scroll_state.go` | Estado mutável do modal, não da fundação |
| Tratamento de teclas | `KeyHandler` como struct utilitária em `modal/key_handler.go` | Composição explícita; modal decide quando e se delega |
| `ModalOption.Keys` | `[]design.Key` (tipado) | Padrão do projeto; `design.Key.Matches(msg)` é robusto. `[]string` + `msg.String()` é anti-padrão já abandonado |
| `ActionGroup.Order` | Campo `Order int` adicionado a `ActionGroup` | Ordem declarativa. Callers existentes com zero-value mantêm comportamento atual |
| Corpo do diálogo | `string` com linhas separadas por `\n` | Mais natural em Go; `strings.Split(body, "\n")` no Render quando necessário |
| Fatiagem do viewport | Feita pelo `HelpModal` antes de chamar `Render` | `DialogFrame.Render` permanece stateless — mesmos argumentos = mesmo output |
| Testes | Golden files em `modal/testdata/golden/` + unit tests para `KeyHandler` | Padrão do projeto (`screen/testdata/golden/`) |

---

## Arquivos alterados/criados

```
internal/tui/actions/
  action.go              ← ALTERADO: adiciona Order int a ActionGroup

internal/tui/modal/
  modal_base.go          ← ALTERADO: ModalOption.Keys muda de []string para []design.Key
  frame.go               ← NOVO
  scroll_state.go        ← NOVO
  key_handler.go         ← NOVO
  key_handler_test.go    ← NOVO
  confirm_modal.go       ← reescrito
  help_modal.go          ← reescrito
  testdata/
    golden/
      frame_no_scroll.txt
      frame_with_scroll_top.txt
      frame_with_scroll_middle.txt
      frame_with_scroll_bottom.txt
      frame_severity_destructive.txt
      frame_severity_error.txt
      confirm_modal.txt
      help_modal_no_scroll.txt
      help_modal_with_scroll.txt

internal/tui/design/
  design_modal.go        ← NOVO
  design_modal_test.go   ← NOVO
```

---

## Alterações em arquivos existentes

### `internal/tui/actions/action.go`

Adicionar campo `Order int` a `ActionGroup`:

```go
// ActionGroup agrupa actions relacionadas para exibição no modal de ajuda.
type ActionGroup struct {
    ID          string // identificador único do grupo
    Label       string // cabeçalho exibido no modal de ajuda
    Description string // texto descritivo do grupo
    Order       int    // ordem de exibição no modal de ajuda; menor valor aparece primeiro
}
```

> Callers existentes que não preenchem `Order` ficam com zero-value `0` — comportamento atual preservado, todos na mesma "ordem 0", serão exibidos na ordem em que aparecerem no slice.

### `internal/tui/modal/modal_base.go`

`ModalOption.Keys` muda de `[]string` para `[]design.Key`:

```go
import "github.com/useful-toys/abditum/internal/tui/design"

// ModalOption representa uma ação disponível ao usuário dentro de um modal.
type ModalOption struct {
    // Keys lista as teclas que ativam esta opção.
    // Keys[0].Label é exibido no rodapé do diálogo.
    // Demais Keys são aliases funcionais (ex: Enter como alias de "S Sobrescrever").
    Keys   []design.Key
    // Label é o texto exibido ao usuário para descrever a ação.
    Label  string
    // Intent classifica a intenção semântica desta ação.
    Intent Intent
    // Action é a função executada quando a opção é escolhida.
    Action func() tea.Cmd
}
```

---

## `design/design_modal.go`

Centraliza todas as helpers de renderização de diálogos. Segue o padrão de `design_action.go` e `design_message.go`: nenhuma cor literal, nenhum símbolo literal, sempre tokens do `Theme` ou constantes de `symbols.go`.

### Constantes de layout

```go
// Padding interno dos diálogos — conforme DS.
const (
    DialogPaddingH = 2  // colunas entre │ e o conteúdo (cada lado)
    DialogPaddingV = 1  // linhas acima e abaixo do conteúdo (Notificação, Confirmação, Ajuda)
    // Funcional usa 0 linhas de padding vertical — documentar na spec de cada subtipo
)
```

### `Severity`

```go
// Severity representa a severidade visual de um diálogo de Notificação ou Confirmação.
// Diálogos de Ajuda e Funcionais não usam severidade.
type Severity int

const (
    SeverityNeutral     Severity = iota // sem símbolo, border.focused, key default em accent.primary
    SeverityInformative                 // ℹ, semantic.info, key default em accent.primary
    SeverityAlert                       // ⚠, semantic.warning, key default em accent.primary
    SeverityDestructive                 // ⚠, semantic.warning, key default em semantic.error
    SeverityError                       // ✕, semantic.error, key default em accent.primary
)

// Symbol retorna o símbolo Unicode da severidade usando constantes de symbols.go.
// Retorna "" para SeverityNeutral.
// SeverityAlert e SeverityDestructive retornam ambos SymWarning — a distinção
// visual está na cor da tecla default (DefaultKeyColor), não no símbolo.
func (s Severity) Symbol() string

// BorderColor retorna a cor de borda da severidade a partir do Theme.
func (s Severity) BorderColor(theme *Theme) string

// DefaultKeyColor retorna a cor da tecla default (1ª ação, ação principal) para a severidade.
// Todas as ações secundárias e de cancelamento usam BorderColor — não DefaultKeyColor.
func (s Severity) DefaultKeyColor(theme *Theme) string
```

### Helpers de renderização

Todas retornam `(string, int)` — string ANSI + largura em colunas.

```go
// RenderDialogTitle renderiza o bloco título da borda superior.
// Se symbol != "", inclui "symbol  title" (símbolo + 2 espaços + título).
// Se symbol == "", inclui apenas "title".
// Cores: symbol em symbolColor, título em theme.Text.Primary + bold.
func RenderDialogTitle(title, symbol, symbolColor string, theme *Theme) (string, int)

// RenderDialogAction renderiza uma ação do rodapé: "key label".
// key é o Label da tecla (Keys[0].Label da ModalOption — ex: "Enter", "S", "Esc").
// key é renderizada em keyColor (ex: Severity.DefaultKeyColor para a 1ª ação,
// Severity.BorderColor para as demais).
// label é renderizada em theme.Text.Primary.
func RenderDialogAction(key, label, keyColor string, theme *Theme) (string, int)

// RenderScrollArrow renderiza ↑ ou ↓ (SymScrollUp / SymScrollDown) em theme.Text.Secondary.
func RenderScrollArrow(up bool, theme *Theme) (string, int)

// RenderScrollThumb renderiza ■ (SymScrollThumb) em theme.Text.Secondary.
func RenderScrollThumb(theme *Theme) (string, int)
```

---

## `modal/scroll_state.go`

```go
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

func (s *ScrollState) Up()
func (s *ScrollState) Down()
func (s *ScrollState) PageUp()
func (s *ScrollState) PageDown()
func (s *ScrollState) Home()
func (s *ScrollState) End()

// ThumbLine calcula a linha (1-based dentro do viewport) onde o thumb ■ deve aparecer.
//
// Regras:
//   - Retorna -1 se o conteúdo não excede o viewport (scroll inativo).
//   - O thumb ocupa posições entre a 1ª e a última linha do viewport, mas
//     NUNCA sobrepõe uma seta ativa:
//     • Se CanScrollUp() == true, a linha 1 do viewport está ocupada por ↑.
//     • Se CanScrollDown() == true, a última linha do viewport está ocupada por ↓.
//     • O thumb é posicionado proporcionalmente nas linhas restantes (entre as setas ativas).
//   - Se o intervalo disponível para o thumb for zero (viewport muito pequeno), retorna -1.
func (s *ScrollState) ThumbLine() int

// CanScrollUp retorna true se há conteúdo acima do viewport (Offset > 0).
func (s *ScrollState) CanScrollUp() bool

// CanScrollDown retorna true se há conteúdo abaixo do viewport.
func (s *ScrollState) CanScrollDown() bool
```

**Invariante de ThumbLine:** as setas `↑`/`↓` **sempre** têm prioridade sobre o thumb.
O thumb só aparece nas linhas que não estão ocupadas por setas ativas.
Alguns wireframes nos golden files estão incorretos (mostram thumb na linha 1 sem seta ↑);
a implementação deve seguir esta spec, não os wireframes.

---

## `modal/frame.go`

`DialogFrame` é uma struct de configuração sem estado. Cada chamada a `Render()` é independente e idempotente.

```go
// DialogFrame define a aparência visual de um diálogo: borda superior com título,
// bordas laterais com indicadores de scroll opcionais, e borda de rodapé com ações posicionadas.
// Não tem estado próprio — quem o usa mantém o ScrollState externamente.
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
    // A 1ª opção é sempre a ação default (tecla: DefaultKeyColor da severidade).
    // A última opção é sempre a ação de cancelamento (tecla: BorderColor).
    // Se há apenas 1 opção, ela é tratada como default.
    Options []ModalOption
    // DefaultKeyColor é a cor da tecla da 1ª opção (ação principal).
    // Ex: Severity.DefaultKeyColor(theme).
    DefaultKeyColor string
    // Scroll é o estado de scroll para exibir indicadores na borda lateral direita.
    // nil = sem scroll.
    Scroll *ScrollState
}

// Render monta a string completa do diálogo a partir do corpo fornecido.
//
// body é uma string com linhas separadas por \n. Cada linha já deve estar renderizada
// com ANSI e ter largura visual de (maxWidth - 2 - 2*DialogPaddingH) colunas —
// o frame não reaplica padding horizontal, apenas adiciona as bordas laterais.
// O caller (ex: HelpModal) é responsável por fatiar o body ao viewport antes de chamar
// Render quando há scroll ativo.
//
// maxWidth é a largura máxima disponível (em colunas).
// theme fornece os tokens de cor para preenchimento e caracteres estruturais.
//
// Fundo: cada linha do corpo recebe background theme.Surface.Raised.
//
// Algoritmo:
//  1. Borda superior: ╭── [símbolo  ]título ───╮
//     (título truncado com … se necessário)
//  2. Para cada linha do body (após Split por \n):
//     │ linha │  com background surface.raised
//     Se Scroll != nil, substitui o │ direito:
//       - linha 1 do body + CanScrollUp()  → ↑ (RenderScrollArrow)
//       - linha N do body + CanScrollDown() → ↓ (RenderScrollArrow)
//       - linha == ThumbLine()              → ■ (RenderScrollThumb)
//       - caso contrário                   → │ normal
//     Prioridade: setas têm prioridade sobre thumb (ThumbLine garante isso).
//  3. Borda de rodapé: ╰─ [ação1] ── [ação2] ── [ação3] ─╯
//     Posicionamento por número de ações:
//       1 ação  → alinhada à direita
//       2 ações → 1ª à esquerda, 2ª à direita
//       3 ações → 1ª à esquerda, 2ª ao centro, 3ª à direita
//     Cor da tecla: 1ª opção usa DefaultKeyColor; demais usam BorderColor.
func (f DialogFrame) Render(body string, maxWidth int, theme *design.Theme) string
```

**Regras de implementação:**
- Caracteres estruturais (`╭╮╰╯│─`) vêm exclusivamente de `design.SymCornerTL` etc.
- O frame **não usa `lipgloss.Border()`** — o border do lipgloss não permite customização char a char do título e do rodapé conforme o DS.
- Coloração é feita com `lipgloss.NewStyle().Foreground(...).Render(char)` por caractere/segmento.
- O preenchimento `─` é colorido com `BorderColor`.
- Background de cada linha do corpo: `lipgloss.NewStyle().Background(lipgloss.Color(theme.Surface.Raised))`.
- A largura mínima do frame é `max(len_visual(title) + overhead, len_visual(maior_acao_rodape) + overhead, 20)`.

---

## `modal/key_handler.go`

`KeyHandler` centraliza o despacho das teclas comuns a todos os diálogos, eliminando a lógica repetitiva que cada modal precisaria reimplementar: iterar opções do rodapé e atualizar o scroll.

```go
// KeyHandler centraliza o despacho de teclas comuns a todos os diálogos:
// ações do rodapé (Options) e navegação de scroll (ScrollState).
//
// O modal concreto compõe um KeyHandler como campo e chama Handle() no momento
// que escolher — podendo interceptar teclas específicas antes ou depois da delegação.
// A chamada a Handle() é sempre explícita e opcional por tecla.
type KeyHandler struct {
    // Options lista as ações cujas teclas devem ser despachadas automaticamente.
    // Corresponde às mesmas Options usadas pelo DialogFrame — geralmente o mesmo slice.
    Options []ModalOption
    // Scroll é o ScrollState a ser atualizado pelas teclas de navegação.
    // nil = sem scroll; teclas de scroll não serão consumidas pelo handler.
    Scroll *ScrollState
}

// Handle processa a tecla fornecida.
//
// Retorna (cmd, true) se a tecla foi consumida — execução de ação ou movimento de scroll.
// Retorna (nil, false) se a tecla não foi reconhecida — o modal decide o que fazer.
//
// Ordem de despacho:
//  1. Opções: itera Options em ordem, compara com cada Key em opt.Keys usando key.Matches(msg).
//     No primeiro match, executa opt.Action() e retorna (cmd, true).
//  2. Scroll (apenas se Scroll != nil):
//     Keys.Up     → Scroll.Up()
//     Keys.Down   → Scroll.Down()
//     Keys.PgUp   → Scroll.PageUp()
//     Keys.PgDn   → Scroll.PageDown()
//     Keys.Home   → Scroll.Home()
//     Keys.End    → Scroll.End()
//     Após atualizar o estado, retorna (nil, true). O ciclo do Bubble Tea
//     re-renderiza o modal naturalmente; nenhum comando extra é necessário.
func (h *KeyHandler) Handle(msg tea.KeyMsg) (tea.Cmd, bool)
```

> As teclas de scroll são comparadas com `design.Keys.Up`, `design.Keys.Down` etc. — nunca com strings.

### Padrão de uso

**Modal sem interceptação (ConfirmModal):**

```go
type ConfirmModal struct {
    // ...
    keys KeyHandler  // composto, não embedded — intenção explícita
}

func (m *ConfirmModal) HandleKey(msg tea.KeyMsg) tea.Cmd {
    if cmd, handled := m.keys.Handle(msg); handled {
        return cmd
    }
    return nil
}
```

**Modal com scroll (HelpModal):**

```go
type HelpModal struct {
    scroll ScrollState
    keys   KeyHandler  // KeyHandler.Scroll aponta para &m.scroll
}

func (m *HelpModal) HandleKey(msg tea.KeyMsg) tea.Cmd {
    if cmd, handled := m.keys.Handle(msg); handled {
        return cmd
    }
    return nil
}
```

**Modal com interceptação antes da delegação (FilePicker, futuro):**

```go
func (m *FilePicker) HandleKey(msg tea.KeyMsg) tea.Cmd {
    // Enter só vai para KeyHandler se o arquivo já está selecionado no campo de texto
    if design.Keys.Enter.Matches(msg) && !m.fileSelected {
        return m.navigateIntoDir()
    }
    if cmd, handled := m.keys.Handle(msg); handled {
        return cmd
    }
    return nil
}
```

### Testes (`modal/key_handler_test.go`)

Testes unitários, sem golden files — comportamento determinístico puro:

| Caso | Expectativa |
|---|---|
| Tecla de opção reconhecida | Retorna `(cmd, true)` — `cmd` é o retorno de `opt.Action()` |
| Tecla de opção com múltiplas Keys | Qualquer das teclas aciona a ação |
| Tecla de scroll com `Scroll != nil` | Retorna `(nil, true)`, estado de `ScrollState` atualizado |
| Tecla de scroll com `Scroll == nil` | Retorna `(nil, false)` — tecla não consumida |
| Tecla não reconhecida em nenhuma opção nem scroll | Retorna `(nil, false)` |
| `Options` vazia | Teclas de scroll ainda funcionam se `Scroll != nil` |

---

## `modal/confirm_modal.go` — reescrito

`ConfirmModal` usa `DialogFrame` e `design.Severity` para renderizar conforme o DS.

```go
// ConfirmModal exibe um diálogo de confirmação com título, mensagem, severidade e ações.
// Implementa tui.ModalView. Criado via NewConfirmModal ou NewConfirmModalSeverity.
type ConfirmModal struct {
    severity design.Severity
    title    string
    message  string
    options  []ModalOption
    keys     KeyHandler  // despacha teclas das opções; sem scroll
}

// NewConfirmModal cria um ConfirmModal de severidade Neutra com as opções fornecidas.
// opts define as ações disponíveis — o caller injeta os closures corretos.
// Convenção: 1ª opção é a ação principal (Enter); última é o cancelamento (Esc).
func NewConfirmModal(title, message string, opts []ModalOption) *ConfirmModal

// NewConfirmModalSeverity cria um ConfirmModal com severidade visual explícita.
func NewConfirmModalSeverity(severity design.Severity, title, message string, opts []ModalOption) *ConfirmModal

// Render constrói um DialogFrame com cores e símbolo derivados da severidade,
// e passa o corpo (mensagem com padding) para o frame renderizar.
func (m *ConfirmModal) Render(maxHeight, maxWidth int, theme *design.Theme) string

// HandleKey delega para m.keys.Handle(msg).
func (m *ConfirmModal) HandleKey(msg tea.KeyMsg) tea.Cmd

// Update processa mensagens Bubble Tea. Delega para HandleKey em tea.KeyMsg;
// ignora demais mensagens.
func (m *ConfirmModal) Update(msg tea.Msg) tea.Cmd
```

**Regras de `Render()`:**
- `DialogFrame.Symbol` = `m.severity.Symbol()`
- `DialogFrame.SymbolColor` = `m.severity.BorderColor(theme)`
- `DialogFrame.BorderColor` = `m.severity.BorderColor(theme)`
- `DialogFrame.TitleColor` = `theme.Text.Primary`
- `DialogFrame.DefaultKeyColor` = `m.severity.DefaultKeyColor(theme)`
- Corpo: `strings.Repeat("\n", DialogPaddingV)` + mensagem + `strings.Repeat("\n", DialogPaddingV)`
- Nenhuma cor hardcoded

**Exemplo de caller** (como um caller constrói o ConfirmModal):
```go
modal.NewConfirmModal("Sair do Abditum", "Sair do Abditum?", []modal.ModalOption{
    {
        Keys:   []design.Key{design.Keys.Enter},
        Label:  "Sair",
        Intent: modal.IntentConfirm,
        Action: func() tea.Cmd { return tui.CloseModal() },
    },
    {
        Keys:   []design.Key{design.Keys.Esc},
        Label:  "Voltar",
        Intent: modal.IntentCancel,
        Action: func() tea.Cmd { return tui.CloseModal() },
    },
})
```

---

## `modal/help_modal.go` — reescrito

`HelpModal` usa `DialogFrame` + `ScrollState` e gera o corpo dinamicamente a partir das `Action`s registradas.

```go
// HelpModal exibe todas as actions registradas, agrupadas por ActionGroup.
// Suporta scroll quando o conteúdo excede o espaço disponível.
// Implementa tui.ModalView.
type HelpModal struct {
    actions []actions.Action
    groups  []actions.ActionGroup
    scroll  ScrollState  // estado de scroll — mutável, começa em Offset=0
    keys    KeyHandler   // despacha scroll (↑↓PgUp/PgDn/Home/End) e Esc (fechar)
}

// NewHelpModal cria o HelpModal com as actions e grupos fornecidos.
// Scroll começa no topo (Offset = 0).
// KeyHandler é configurado com Scroll apontando para &m.scroll e
// Options = [{Keys: [Esc], Label: "Fechar", Intent: IntentCancel, Action: CloseModal}].
func NewHelpModal(acts []actions.Action, groups []actions.ActionGroup) *HelpModal

// Render gera o corpo dinamicamente, fatia o viewport conforme scroll,
// e passa para DialogFrame.Render.
//
// Algoritmo:
//  1. Gerar allLines []string: todas as linhas do conteúdo (grupos + ações).
//     Ver "Geração do corpo" abaixo.
//  2. Calcular viewport = maxHeight - 2 (subtraindo borda superior e borda de rodapé).
//     Atualizar m.scroll.Total = len(allLines) e m.scroll.Viewport = viewport.
//  3. Fatiar: visibleLines = allLines[m.scroll.Offset : m.scroll.Offset+viewport]
//     (clampado para não ultrapassar Total).
//  4. body = strings.Join(visibleLines, "\n")
//  5. Chamar DialogFrame{...}.Render(body, maxWidth, theme).
func (m *HelpModal) Render(maxHeight, maxWidth int, theme *design.Theme) string

// HandleKey delega para m.keys.Handle(msg).
func (m *HelpModal) HandleKey(msg tea.KeyMsg) tea.Cmd

// Update processa mensagens Bubble Tea. Delega para HandleKey em tea.KeyMsg;
// ignora demais mensagens.
func (m *HelpModal) Update(msg tea.Msg) tea.Cmd
```

**Geração do corpo (`allLines`):**
- Grupos ordenados por `ActionGroup.Order` crescente; empate → ordem original do slice.
- Para cada grupo:
  - Linha de cabeçalho: `lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Secondary)).Bold(true).Render(group.Label)`
  - Para cada action do grupo, ordenada por `Action.Priority` crescente (menor = mais destaque):
    - Linha de ação: tecla (`Keys[0].Label`) em `theme.Accent.Primary`; descrição (`Description`) em `theme.Text.Primary`; colunas alinhadas com espaçamento fixo de 14 colunas para a coluna de teclas.
  - Linha em branco **entre** grupos (não antes do primeiro, não após o último).
- As linhas geradas **não** incluem padding vertical (DialogPaddingV) — o HelpModal não usa padding V, conforme o wireframe do DS.

**Frame:**
- `BorderColor` = `theme.Border.Default`
- `Symbol` = `""`
- `TitleColor` = `theme.Text.Primary`
- `Title` = `"Ajuda — Atalhos e Ações"` (conforme `tui-spec-dialog-help.md`)
- `DefaultKeyColor` = `theme.Accent.Primary`
- `Options` = `[]ModalOption{{Keys: []design.Key{design.Keys.Esc}, Label: "Fechar", Intent: IntentCancel, Action: tui.CloseModal}}`

> **Nota sobre o título:** `tui-spec-dialogos.md` usa `"Ajuda"` como exemplo simplificado; `tui-spec-dialog-help.md` é o documento autoritativo e usa `"Ajuda — Atalhos e Ações"`.

---

## Testes

### `design/design_modal_test.go`
Testes unitários para cada helper:
- `RenderDialogTitle` — com e sem símbolo; largura retornada == `lipgloss.Width(text)`
- `RenderDialogAction` — cor de tecla default vs secundária; largura correta
- `RenderScrollArrow` e `RenderScrollThumb` — largura == 1

### `modal/key_handler_test.go`
Testes unitários para `KeyHandler.Handle()` — ver tabela na seção `key_handler.go`.

### `modal/` — golden files

| Arquivo | Cenário |
|---|---|
| `frame_no_scroll.txt` | Borda simples, 2 opções, sem scroll |
| `frame_with_scroll_top.txt` | Scroll no topo: apenas `↓` visível na última linha; thumb proporcionalmente nas linhas intermediárias |
| `frame_with_scroll_middle.txt` | Scroll no meio: `↑` na 1ª linha, `↓` na última, thumb no meio |
| `frame_with_scroll_bottom.txt` | Scroll no final: apenas `↑` visível na 1ª linha; thumb proporcionalmente |
| `frame_severity_destructive.txt` | Borda `semantic.warning`, símbolo `⚠`, tecla default em `semantic.error` |
| `frame_severity_error.txt` | Borda `semantic.error`, símbolo `✕`, tecla default em `accent.primary` |
| `confirm_modal.txt` | ConfirmModal completo com severidade Destrutiva |
| `help_modal_no_scroll.txt` | HelpModal sem overflow |
| `help_modal_with_scroll.txt` | HelpModal com scroll ativo |

Os golden files são gerados com `go test -update-golden ./...` e verificados em CI sem a flag. Padrão idêntico ao usado em `screen/testdata/golden/` (helper `testdata.TestRenderManaged`).

---

## Restrições e invariantes

1. Nenhuma cor literal (`"#7aa2f7"` etc.) nos arquivos do pacote `modal/` — todas as cores passam pelo `design.Theme`.
2. Nenhum símbolo literal (`"⚠"`, `"│"` etc.) nos arquivos do pacote `modal/` — todos os símbolos vêm de constantes em `design/symbols.go`.
3. `DialogFrame` não tem estado — qualquer chamada a `Render()` com os mesmos argumentos produz a mesma string.
4. `ScrollState` é externo ao `DialogFrame` — é responsabilidade do modal que usa scroll.
5. `KeyHandler` é composto (não embedded) nos modais concretos — a delegação é sempre explícita.
6. A interface `tui.ModalView` não muda — `RootModel` não precisa de alteração.
7. O pacote `modal/` não importa nada além de `charm.land/bubbletea/v2`, `charm.land/lipgloss/v2`, e os pacotes internos `tui`, `tui/design`, `tui/actions`.
8. Máximo de 3 ações na `Options` do `DialogFrame` — conforme DS.
9. `ThumbLine()` nunca retorna a posição de uma seta ativa. Setas têm prioridade absoluta sobre o thumb.
10. O thumb `■` só aparece quando o scroll está ativo (`Total > Viewport`).
11. `ModalOption.Keys` é `[]design.Key` — nunca `[]string`. Comparação via `key.Matches(msg)`.

---

## Fora do escopo desta entrega

- `PasswordEntry` e `PasswordCreate` — documentados em `golden/tui-spec-dialog-senha.md`
- `FilePicker` — documentado em `golden/tui-spec-dialog-filepicker.md`
- Word-wrap automático do corpo — o caller é responsável por quebrar linhas
- Suporte a divisores internos do corpo (`├─┤`) — reservado para diálogos Funcionais futuros
