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

    // Última dimensão recebida em Render() — necessária para Cursor() calcular
    // a linha Y do campo Arquivo: fora do contexto de Render().
    lastMaxHeight int
    lastMaxWidth  int

    // Injeção de teste
    timeFmt func(time.Time) string  // nil → Local "02/01/06 15:04"
}
```

> **Nota sobre `lastMaxHeight`/`lastMaxWidth`:** `Cursor(topY, leftX)` precisa saber `visibleH` para calcular a linha Y do campo `Arquivo:`. Como `Cursor()` não recebe `maxHeight`, esses campos são atualizados no início de cada `Render()` para que `Cursor()` possa reutilizá-los. O mesmo padrão se aplica ao `lastMaxWidth` para consistência (não é usado em `Cursor()` mas é necessário para `visibleH` no modo Save).

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
modalW = min(70, maxWidth * 80 / 100)  — DS: 70 colunas ou 80% terminal, o menor
```

### Scroll duplo independente

Cada painel mantém seu próprio par `(cursor, scroll)`:

- **Árvore:** `treeCursor int`, `treeScroll int`
- **Arquivos:** `fileCursor int`, `fileScroll int`

O indicador de scroll de cada painel substitui o caractere de borda do respectivo lado:
- Árvore: o `│` do separador entre painéis é substituído por `↑`/`■`/`↓`
- Arquivos: o `│` da borda direita do modal é substituído por `↑`/`■`/`↓`

Funções `renderTreeSepChar()` e `renderFileSepChar()` calculam qual caractere usar em cada linha.

> **`KeyHandler` e `ScrollState` do pacote `modal`:** o pacote já oferece esses helpers para scroll vertical simples. Eles **não são usados** aqui porque o FilePicker tem dois painéis de scroll completamente independentes, cada um com seus próprios cursores e viewports — o `ScrollState` existente é projetado para um único fluxo de scroll por modal. O scroll duplo é mantido manual, exatamente como no legado.

### Suporte a mouse

O DS especifica clique simples e duplo-clique nos painéis. Esta implementação **não inclui suporte a mouse** — está fora do escopo desta iteração. A interface é completamente operável por teclado, que é o canal primário conforme o DS.

### Truncamento de metadados no painel de arquivos

A regra de prioridade de truncamento (conforme golden spec, linha 325):

- Em terminais estreitos, metadados (tamanho + data/hora) são **truncados primeiro** (da direita para a esquerda)
- O nome do arquivo tem prioridade — só trunca com `…` se não houver espaço mesmo para ele

Implementação: calcular `nameW = filesW - 1(bullet) - 1(espaço) - sizeW - colSep - dateW - colSep`. Se `nameW < mínimo aceitável` (ex: 4), truncar `dateW` primeiro, depois `sizeW`, até que `nameW` seja viável. O nome nunca é sacrificado antes dos metadados.

---

## Cursor de terminal (campo Arquivo:)

`Cursor(topY, leftX int) *tea.Cursor`:

- Retorna `nil` quando `focusPanel != 2`
- Quando `focusPanel == 2`:
  - `visibleH` é recalculado usando `lastMaxHeight` e `mode` — mesma fórmula de `Render()`
  - `Y = topY + 4 + visibleH` (linha do campo Arquivo:)
  - `X = leftX + 1 + len("Arquivo: ") + nameField.Position()` — `textinput.Position()` retorna a coluna do cursor dentro do valor (não o final — suporta navegação com `←`/`→`)

`lastMaxHeight` e `lastMaxWidth` são atualizados no início de cada `Render()`. Antes do primeiro `Render()`, `Cursor()` retorna `nil` (campo não visível ainda).

---

## Integração com MessageController

Hints e erros são emitidos via `tui.MessageController`:

| Situação | Método | Texto | Permanência |
|---|---|---|---|
| Foco na árvore (Open) | `SetHintField()` | `• Navegue pelas pastas e selecione um cofre` | Permanente |
| Foco na árvore (Save) | `SetHintField()` | `• Navegue pelas pastas e escolha onde salvar` | Permanente |
| Foco no painel de arquivos (Open, com arquivos) | `SetHintField()` | `• Selecione o cofre para abrir` | Permanente |
| Foco no painel de arquivos (Open, vazio) | `SetHintField()` | `• Nenhum cofre neste diretório — navegue para outra pasta` | Permanente |
| Foco no painel de arquivos (Save) | `SetHintField()` | `• Arquivos existentes neste diretório` | Permanente |
| Foco no campo Arquivo: (vazio) | `SetHintField()` | `• Digite o nome do arquivo — <ext> será adicionado automaticamente` | Permanente |
| Foco no campo Arquivo: (preenchido) | `SetHintField()` | `• Confirme para salvar o cofre` | Permanente |
| Sem permissão para expandir pasta | `SetError()` | `✕ Sem permissão para acessar <nomePasta>` | 5s |
| Fallback de CWD inacessível | `SetWarning()` | `⚠ Diretório atual inacessível — navegando para home` | Permanente |

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
| `←` na raiz (`depth == 0`) | Sem efeito — a seleção permanece na raiz |
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

## Fora de escopo

- **Suporte a mouse** (clique simples, duplo-clique): não implementado nesta iteração. O DS especifica esses eventos, mas o canal primário é teclado.
- **Atalhos de teclado adicionais**: apenas os definidos na spec visual são implementados.
