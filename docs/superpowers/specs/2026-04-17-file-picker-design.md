# Design — FilePicker Modal

**Data:** 2026-04-17  
**Spec de referência:** `golden/tui-spec-dialog-filepicker.md`  
**Design system:** `golden/tui-design-system.md`

---

## Decisões de alto nível

| Decisão | Escolha | Justificativa |
|---|---|---|
| Abordagem | Porta direta do legado com adaptação de integração | Lógica de negócio comprovada; menor risco de regressão |
| Retorno de resultado | Callback `OnResult func(path string) tea.Cmd` | Cancelamento: `""`. Confirmação: caminho completo |
| Cursor no campo Arquivo: | `Cursor()` real via `ModalView` | Compatibilidade com arquitetura atual |
| Campo Arquivo: | `charm.land/bubbles/v2/textinput` | Reutiliza componente já usado no legado |
| Reutilização de DialogFrame | Não reutilizado — exceção documentada no DS | Layout de dois painéis com scrolls independentes incompatível com DialogFrame |
| Funções atômicas do design | Reutilizadas: `RenderDialogTitle`, `RenderDialogAction`, `RenderScrollArrow`, `RenderScrollThumb`, todos os `Sym*` | Zero hardcode de cores, fontes, símbolos |
| Localização | `internal/tui/modal/file_picker.go` (arquivo único) | Coerente com `confirm_modal.go` e `help_modal.go` |
| Modificações em arquivos existentes | Nenhuma | FilePicker é autossuficiente |

---

## Contrato público

### Tipos

```go
// FilePickerMode controla o comportamento do picker.
type FilePickerMode int

const (
    FilePickerOpen FilePickerMode = iota // abrir arquivo existente
    FilePickerSave                       // salvar / nomear novo arquivo
)

// FilePickerOptions são os parâmetros de construção do FilePicker.
type FilePickerOptions struct {
    Mode       FilePickerMode          // Open | Save
    Extension  string                 // ex: ".abditum" — inclui o ponto
    InitialDir string                 // "" → CWD → fallback ~
    Suggested  string                 // Save: valor inicial do campo Arquivo:
    OnResult   func(path string) tea.Cmd // path="" se cancelado
    Messages   tui.MessageController  // nil tolerado (sem hints/erros)
}
```

### Construtor

```go
// NewFilePicker cria e inicializa o modal.
// Chamar tui.OpenModal(NewFilePicker(opts)) para exibir.
func NewFilePicker(opts FilePickerOptions) *FilePickerModal
```

### Interface ModalView

```go
func (m *FilePickerModal) Render(maxHeight, maxWidth int, theme *design.Theme) string
func (m *FilePickerModal) HandleKey(msg tea.KeyMsg) tea.Cmd
func (m *FilePickerModal) Update(msg tea.Msg) tea.Cmd
func (m *FilePickerModal) Cursor(topY, leftX int) *tea.Cursor
```

---

## Estrutura interna

### Tipos privados

```go
// treeNode representa uma entrada de diretório na árvore lazy (carregamento sob demanda).
type treeNode struct {
    path       string
    name       string
    depth      int
    expanded   bool
    loaded     bool
    children   []*treeNode
    hasSubdirs bool
}

// visibleNode é uma entrada achatada da árvore — visível no painel Estrutura.
type visibleNode struct {
    node *treeNode
}
```

### Struct principal

```go
type FilePickerModal struct {
    // Parâmetros de construção
    mode       FilePickerMode
    ext        string
    onResult   func(string) tea.Cmd
    messages   tui.MessageController  // nil tolerado

    // Estado da árvore
    root         *treeNode
    visibleNodes []visibleNode
    treeCursor   int
    treeScroll   int

    // Diretório e lista de arquivos atuais
    currentPath string
    files       []string
    fileInfos   []os.FileInfo
    fileCursor  int  // -1 quando vazio
    fileScroll  int

    // Foco: 0=árvore, 1=arquivos, 2=campo nome (Save apenas)
    focusPanel int

    // Campo Arquivo: (Save mode)
    nameField textinput.Model

    // Injeção de teste
    timeFmt func(time.Time) string  // nil → Local "02/01/06 15:04"
}
```

---

## Renderização

### Estratégia

O `FilePickerModal` renderiza o modal inteiro em `Render()`. O `DialogFrame` existente **não é usado** — o layout de dois painéis com dois scrolls independentes e separadores internos é incompatível com a estrutura do frame (exceção documentada no DS).

Funções atômicas do `design` package reutilizadas:

| Função/Constante | Uso |
|---|---|
| `design.RenderDialogTitle()` | Borda superior — título formatado |
| `design.RenderDialogAction()` | Rodapé — ações Enter/Esc com cores corretas |
| `design.RenderScrollArrow()` | Indicadores ↑/↓ nos painéis |
| `design.RenderScrollThumb()` | Thumb ■ nos painéis |
| `design.Sym*` (todos) | Bordas, junctions, indicadores, ícones de pasta/arquivo |

**Zero hardcode:** nenhuma cor, símbolo ou string de UI é definida como literal no arquivo. Tudo via `theme.*` ou `design.Sym*`.

### Layout vertical (linhas a partir do topo do modal)

```
linha 0:           borda superior ╭── título ──╮
linha 1:           caminho │ /path/to/dir       │
linha 2:           separador ├─ Estrutura ─┬─ Arquivos ─┤
linhas 3..2+H:     conteúdo dos painéis lado a lado
linha 3+H:         (Save) separador campo ├────┴────┤
linha 4+H:         (Save) campo │ Arquivo: ░...░ │
linha 3+H ou 5+H:  borda inferior ╰── Enter ──── Esc ──╯
```

Onde `H = visibleH` calculado a partir de `maxHeight * 8/10 - overheadLines`.

### Largura dos painéis

```
innerW = modalW - 2
treeW  = innerW * 40 / 100  (mínimo 8)
filesW = innerW - treeW - 1  (mínimo 8)
modalW = min(maxWidth, 70)  — seguindo DS: máximo 70 ou 80% terminal, o menor
```

### Scroll duplo independente

Cada painel mantém seu próprio par `(cursor, scroll)`:

- **Árvore:** `treeCursor int`, `treeScroll int`
- **Arquivos:** `fileCursor int`, `fileScroll int`

O indicador de scroll de cada painel substitui o caractere de borda do respectivo lado:
- Árvore: o `│` do separador entre painéis é substituído por `↑`/`■`/`↓`
- Arquivos: o `│` da borda direita do modal é substituído por `↑`/`■`/`↓`

Funções `renderTreeSepChar()` e `renderFileSepChar()` calculam qual caractere usar em cada linha.

---

## Cursor de terminal (campo Arquivo:)

`Cursor(topY, leftX int) *tea.Cursor`:

- Retorna `nil` quando `focusPanel != 2`
- Quando `focusPanel == 2`:
  - `Y = topY + 4 + visibleH` (linha do campo Arquivo:)
  - `X = leftX + 1 + len("Arquivo: ") + len([]rune(nameField.Value()))` (após o valor digitado)

O `visibleH` é calculado deterministicamente a partir de `maxHeight` e `mode` — mesmo algoritmo de `Render()`.

---

## Integração com MessageController

Hints e erros são emitidos via `tui.MessageController`:

| Situação | Método | Permanência |
|---|---|---|
| Foco na árvore | `SetHintField()` | Permanente |
| Foco no painel de arquivos | `SetHintField()` | Permanente |
| Foco no campo Arquivo: | `SetHintField()` | Permanente |
| Sem permissão para expandir pasta | `SetError()` | 5s |
| Fallback de CWD inacessível | `SetWarning()` | Permanente |

`messages` é verificado contra `nil` antes de cada chamada — tolerância a testes sem message bar.

O `emitHint()` retorna `tea.Cmd` e é chamado ao final de cada handler de teclado e em `Init()`.

---

## Comportamentos portados do legado

Todos os comportamentos da spec são portados integralmente. Os mais críticos (nuances do legado):

| Comportamento | Detalhe |
|---|---|
| `buildTreeChain` | Constrói a cadeia da raiz até `initialDir` com ancestors expandidos. Funciona em Unix (`/`) e Windows (`C:\`) via `filepath.VolumeName` |
| Lazy loading | `expandNode()` carrega filhos sob demanda em `buildVisibleNodes()` |
| `adjustTreeScroll` em `SetSize` | Scroll resetado para 0 em `SetSize` e recalculado — evita scroll incorreto computado com `height=0` durante `Init` |
| Cursor passivo na árvore | Quando foco está em outro painel, a pasta selecionada usa bold+accent sem fundo `special.highlight` |
| `Enter` na árvore (Open/Save) | Avança foco para o primeiro arquivo se a pasta contém `<ext>`; sem efeito se vazia |
| `Enter` no painel de arquivos (Save) | Copia nome para campo, move foco — não confirma o diálogo |
| `Tab` em painel de arquivos vazio | Open: volta para árvore; Save: pula para campo |
| Caracteres inválidos no campo | `/\:*?"<>|` bloqueados silenciosamente |
| Extensão automática no Save | Adicionada silenciosamente ao caminho de retorno se não presente |
| `←` em pasta expandida | Recolhe; em pasta já recolhida, navega para o pai (`parentCursor`) |
| `→` em `▷` (sem subdiretórios) | Sem efeito |
| Indicador de pasta raiz | Raiz não exibe indicador (`▶`/`▼`/`▷`) — apenas nome |

---

## Funções utilitárias privadas

Implementadas em `file_picker.go` como funções privadas de pacote:

```go
// formatFileSize formata bytes em KB/MB/GB (base 1024, 1 casa decimal).
func formatFileSize(bytes int64) string

// padRight pads s até width colunas visuais (ANSI-aware via lipgloss.Width).
func padRight(s string, width int) string
```

Verificar antes da implementação se `padRight` já existe em outro arquivo do pacote `modal` para evitar duplicação.

---

## O que NÃO muda

- Nenhum arquivo existente é modificado
- `DialogFrame` permanece inalterado
- `design` package não recebe novas funções
- O sistema de `Operation` não é envolvido — o FilePicker é um `ModalView` puro; o orquestrador (Operation ou ChildView que chamou `NewFilePicker`) trata o resultado via `OnResult`
