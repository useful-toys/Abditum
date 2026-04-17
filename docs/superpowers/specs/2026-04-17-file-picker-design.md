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

**Comportamento do construtor:**

1. Resolve `InitialDir`: se vazio, usa `os.Getwd()`. Se `os.Getwd()` falha ou o diretório não existe, usa `os.UserHomeDir()` e emite `SetWarning` ao primeiro `Render()` (não dentro do construtor, pois `messages` pode estar disponível apenas depois).
2. Constrói a árvore da raiz até `InitialDir` via `buildTreeChain()` — todos os ancestors ficam `expanded=true`, o nó `InitialDir` fica selecionado (`treeCursor` aponta para ele).
3. Carrega os arquivos do `InitialDir` via `loadFiles()`.
4. Pré-seleciona o primeiro arquivo se houver: `fileCursor = 0`, caso contrário `fileCursor = -1`.
5. Se `opts.Suggested != ""`, preenche `nameField.SetValue(opts.Suggested)`.
6. `focusPanel = 0` (árvore).
7. `readDir = nil` (usa `os.ReadDir` real); pode ser sobrescrito em testes após construção via campo exportado de teste — **alternativa**: adicionar `WithReadDir(fn)` option func. Como os outros campos de teste (`timeFmt`) são configurados diretamente, `readDir` segue o mesmo padrão: campo privado acessível via função setter de teste ou via struct literal em `_test.go` dentro do mesmo pacote. Como o arquivo de teste é `package modal_test` (externo), os campos de injeção precisam ser **exportados** ou expostos via método setter. Decisão: **métodos setter com sufixo `ForTest`** — `SetReadDirForTest(fn)` e `SetTimeFmtForTest(fn)` — que existem apenas para facilitar testes sem poluir a API pública.

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

    // Injeção de teste — nil usa os.ReadDir real.
    // Testes injetam uma função que retorna uma árvore hipotética fixa, sem disco.
    readDir func(path string) ([]os.DirEntry, error)

    // hintEmitted controla se o hint inicial já foi emitido.
    // O hint é enviado na primeira oportunidade em Update() ou HandleKey().
    hintEmitted bool

    // Injeção de teste — nil → Local "02/01/06 15:04"
    timeFmt func(time.Time) string
}
```

**Toda leitura de diretório no código de produção passa por `m.readDir(path)`.** Quando `readDir == nil`, usa `os.ReadDir`. Isso inclui `expandNode()` e a listagem de arquivos ao mudar de pasta.

**Acesso nos testes:** como o arquivo de teste é `package modal_test` (externo), os campos de injeção são expostos via métodos setter:

```go
// SetReadDirForTest injeta um filesystem fictício — usado exclusivamente em testes.
func (m *FilePickerModal) SetReadDirForTest(fn func(string) ([]os.DirEntry, error)) { m.readDir = fn }

// SetTimeFmtForTest injeta formatação de tempo fixa — usado exclusivamente em testes.
func (m *FilePickerModal) SetTimeFmtForTest(fn func(time.Time) string) { m.timeFmt = fn }
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

`overheadLines` = número de linhas fixas fora do conteúdo dos painéis:
- Modo Open: 3 linhas fixas (borda superior + linha de caminho + separador de painéis + borda inferior = 4, mas a borda inferior não consome `visibleH` → `overheadLines = 3`)
- Modo Save: 5 linhas fixas (borda superior + caminho + separador painéis + separador campo + campo Arquivo: + borda inferior = 6, menos a borda inferior → `overheadLines = 5`)

Fórmula completa:
```
modalH     = maxHeight * 8 / 10
overhead   = 3          // Open; 5 para Save
visibleH   = modalH - overhead
if visibleH < 3 { visibleH = 3 }  // mínimo operacional
```

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

`emitHint()` retorna `tea.Cmd` e é chamado ao final de cada handler de teclado. O hint inicial (foco na árvore ao abrir) é emitido no primeiro `HandleKey` ou pode ser disparado pelo orquestrador via `Update(tui.ModalReadyMsg{})` — a interface `ModalView` não tem `Init()`. Alternativa adotada: `Update()` verifica se é o primeiro update e emite o hint inicial; ou o hint é emitido imediatamente no construtor como `tea.Cmd` retornado por uma função de setup. **Decisão:** `NewFilePicker` não retorna `tea.Cmd`. O hint inicial é emitido no primeiro `HandleKey` ou `Update` que processar qualquer mensagem. Para garantir que o hint apareça sem interação do usuário, o orquestrador deve chamar `Update(nil)` logo após abrir o modal — ou o `RootModel` passa um `tea.WindowSizeMsg` inicial que dispara `Update()`. Como outros modais não têm esse problema (não emitem hints), a solução mais simples é: emitir o hint inicial dentro de `HandleKey` apenas quando for o primeiro evento de teclado recebido, e também em `Update` quando recebe `tea.WindowSizeMsg`. **Decisão final:** `emitHint()` é chamado em `HandleKey` e em `Update` para qualquer msg não-nula. O campo `hintEmitted bool` controla se o hint inicial já foi enviado — na primeira oportunidade é emitido.

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
| `←` em pasta expandida | Recolhe; em pasta já recolhida, navega para o pai: percorre `visibleNodes` para trás buscando o primeiro nó com `depth < nóAtual.depth` |
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

---

## Testes

### Estratégia

Os testes seguem o mesmo padrão do restante do pacote `modal`:

- **Golden files** via `testdata.TestRenderManaged` (produz `.golden.txt` + `.golden.json`)
- **Injeção de filesystem** via campo `readDir` na struct — nil usa `os.ReadDir`; testes injetam uma função que retorna uma árvore hipotética fixa sem tocar disco
- **Injeção de tempo** via campo `timeFmt` — testes passam uma função que retorna datas fixas para output determinístico
- Arquivo de teste: `internal/tui/modal/file_picker_test.go` (pacote `modal_test`)

### Árvore hipotética para testes

Todos os testes de render e comportamento usam a mesma árvore fictícia. Ela deve ter profundidade e largura suficientes para forçar scroll na árvore com janelas de altura pequena (`visibleH ≈ 8`). Mínimo: 14 nós visíveis quando totalmente expandida.

```
/
└── home/
    └── usuario/
        ├── documentos/
        │   ├── contratos/
        │   │   ├── 2024/
        │   │   └── 2025/
        │   └── relatorios/
        ├── downloads/
        │   ├── instaladores/
        │   └── temporarios/
        ├── projetos/
        │   ├── abditum/
        │   │   ├── docs/
        │   │   └── src/
        │   └── site/
        └── fotos/
```

Arquivos `.abditum` presentes em:
- `/home/usuario/projetos/abditum/` → `database.abditum` (25_800_000 bytes, 2025-03-15 14:32), `config.abditum` (1_229 bytes, 2025-01-02 09:15), `backup.abditum` (18_400_000 bytes, 2025-04-04 18:47)
- `/home/usuario/documentos/contratos/2025/` → `cofre.abditum` (512_000 bytes, 2025-04-01 10:00)

Todos os outros diretórios: sem arquivos `.abditum`.

**Implementação nos testes:** `makeTestReadDir()` retorna um `func(path string) ([]os.DirEntry, error)` que mapeia cada caminho acima para entradas fixas. Retorna `fs.ErrPermission` para o caminho especial `/home/usuario/documentos/contratos/2024/` — usado nos testes de permissão negada.

Configuração padrão de um modal de teste:
```go
func newOpenPicker(t *testing.T) *modal.FilePickerModal {
    m := modal.NewFilePicker(modal.FilePickerOptions{
        Mode:       modal.FilePickerOpen,
        Extension:  ".abditum",
        InitialDir: "/home/usuario/projetos/abditum",
        OnResult:   func(string) tea.Cmd { return nil },
    })
    m.SetReadDirForTest(makeTestReadDir())
    m.SetTimeFmtForTest(func(tm time.Time) string { return tm.Format("02/01/06 15:04") })
    return m
}
```

`makeTestTimeFmt()` retorna `func(time.Time) string` que usa `t.Format("02/01/06 15:04")` — deterministico porque `fileInfos` terá `ModTime()` fixos via `fakeFileInfo`.

`fakeFileInfo` implementa `os.FileInfo` com campos fixos (nome, tamanho, `ModTime`, `IsDir`).
`fakeDirEntry` implementa `os.DirEntry` wrappando `fakeFileInfo`.

### Golden files de render

Localização: `internal/tui/modal/testdata/golden/file_picker-<variant>-<size>.golden.{txt,json}`

**Tamanho padrão para golden files de render: `88x30`** (terminal amplo, `modalW=70`, `visibleH≈18`).

| Arquivo golden | Variante | Estado do modal | Foco | Scroll árvore | Scroll arquivos |
|---|---|---|---|---|---|
| `file_picker-open_tree_initial-88x30` | Open, árvore inicial | Árvore aberta com `initialDir=/home/usuario/projetos/abditum`, ancestors expandidos (`/`, `home`, `usuario`, `projetos`, `abditum`), cursor em `abditum`. Painel de arquivos mostra os 3 arquivos `.abditum`, primeiro pré-selecionado. Sem scroll (tudo cabe). | árvore (0) | sem scroll | sem scroll |
| `file_picker-open_files_noscroll-88x30` | Open, arquivos sem scroll | Mesmo estado acima, foco movido para painel de arquivos. 3 arquivos visíveis, `database.abditum` selecionado (highlight). | arquivos (1) | sem scroll | sem scroll |
| `file_picker-open_files_scroll_top-88x30` | Open, scroll arquivos topo | Pasta com 12 arquivos fictícios (lista longa o suficiente para scroll, adicionar ao `makeTestReadDir` para `/home/usuario/projetos/abditum/` quando solicitado em testes de scroll). `fileScroll=0`, cursor no item 0. Indicador `↓` visível. | arquivos (1) | qualquer | scroll início |
| `file_picker-open_files_scroll_mid-88x30` | Open, scroll arquivos meio | Mesma lista longa. `fileScroll` posicionado para que cursor esteja no meio da lista. Indicadores `↑` e `↓` visíveis. | arquivos (1) | qualquer | scroll meio |
| `file_picker-open_files_scroll_end-88x30` | Open, scroll arquivos fim | Mesma lista longa. `fileScroll` no máximo. Indicador `↑` visível. | arquivos (1) | qualquer | scroll fim |
| `file_picker-open_tree_scroll_top-88x30` | Open, scroll árvore topo | Árvore totalmente expandida (≥14 nós), janela pequena `88x14` para forçar scroll. `treeScroll=0`. Indicador `↓` no separador. | árvore (0) | scroll início | qualquer |
| `file_picker-open_tree_scroll_mid-88x30` | Open, scroll árvore meio | Idem. `treeScroll` no meio. Indicadores `↑` e `↓`. | árvore (0) | scroll meio | qualquer |
| `file_picker-open_tree_scroll_end-88x30` | Open, scroll árvore fim | Idem. `treeScroll` no máximo. Indicador `↑`. | árvore (0) | scroll fim | qualquer |
| `file_picker-save_name_empty-88x30` | Save, campo vazio | `initialDir=/home/usuario/projetos/abditum`, `Suggested=""`. Campo Arquivo: vazio. Ação Enter em `text.disabled`. | campo (2) | sem scroll | sem scroll |
| `file_picker-save_name_filled-88x30` | Save, campo preenchido | Mesmo acima, `nameField` com valor `"meu-cofre"`. Ação Enter em `accent.primary`. Cursor `▌` visível no campo. | campo (2) | sem scroll | sem scroll |
| `file_picker-open_empty_dir-88x30` | Open, pasta vazia | Cursor em `/home/usuario/downloads/temporarios/` (sem `.abditum`). Painel de arquivos exibe `Nenhum cofre neste diretório`. Ação Enter bloqueada. | árvore (0) | sem scroll | N/A |

> **Nota sobre listas longas:** o `makeTestReadDir` deve ter uma segunda configuração (ou um helper `withManyFiles`) que retorna 12 arquivos `.abditum` numerados em `/home/usuario/projetos/abditum/` para os 3 golden files de scroll de arquivos. Os 3 arquivos normais (`database`, `config`, `backup`) são usados nos demais testes.

### Testes de comportamento (sem golden)

| Teste | Cenário | Verificação |
|---|---|---|
| `TestFilePicker_Open_Enter_OnFile` | foco em arquivos, pressionar Enter | `OnResult` chamado com path completo |
| `TestFilePicker_Open_Esc_Cancels` | pressionar Esc em qualquer foco | `OnResult` chamado com `""` |
| `TestFilePicker_Save_Enter_WithName` | campo preenchido com `"meu-cofre"`, pressionar Enter | `OnResult` recebe `"/home/.../meu-cofre.abditum"` (extensão adicionada) |
| `TestFilePicker_Save_Enter_WithExtension` | campo preenchido com `"meu-cofre.abditum"`, pressionar Enter | `OnResult` recebe path sem extensão duplicada |
| `TestFilePicker_Save_Enter_EmptyField` | campo vazio, pressionar Enter | `OnResult` não chamado; cmd == nil |
| `TestFilePicker_Tree_Left_CollapsesNode` | pasta expandida, pressionar `←` | `expanded=false`; `buildVisibleNodes` reflete o colapso |
| `TestFilePicker_Tree_Left_AtRoot` | cursor na raiz, pressionar `←` | sem efeito; cursor permanece 0 |
| `TestFilePicker_Tree_Right_EmptyFolder` | cursor em pasta `▷`, pressionar `→` | sem efeito |
| `TestFilePicker_InvalidChars_Blocked` | foco no campo, pressionar `/` | caractere não inserido |
| `TestFilePicker_Tab_Cycles_Open` | Tab 2x no modo Open | cicla árvore → arquivos → árvore |
| `TestFilePicker_Tab_Cycles_Save` | Tab 3x no modo Save | cicla árvore → arquivos → campo → árvore |
| `TestFilePicker_PermissionDenied` | expandir `/home/usuario/documentos/contratos/2024/` | pasta permanece recolhida; `SetError` chamado no stub |
| `TestFilePicker_Cursor_NilWhenNotSave` | modo Open | `Cursor()` retorna nil |
| `TestFilePicker_Cursor_NilWhenNotFocused` | modo Save, foco=0 | `Cursor()` retorna nil |
| `TestFilePicker_Cursor_Position` | modo Save, foco=2, `nameField` com valor `"abc"` | `Cursor()` retorna posição correta (X e Y) |
| `TestFormatFileSize` | unitário puro | `25800000` → `"25.8 MB"`, `1229` → `"1.2 KB"`, `2000000000` → `"1.9 GB"` |

### Stub de MessageController

```go
// stubMessageController implementa tui.MessageController para testes.
// Grava a última chamada para asserção sem efeitos colaterais.
type stubMessageController struct {
    lastMethod string
    lastText   string
}

func (s *stubMessageController) SetHintField(text string) { s.lastMethod = "HintField"; s.lastText = text }
func (s *stubMessageController) SetError(text string)     { s.lastMethod = "Error"; s.lastText = text }
func (s *stubMessageController) SetWarning(text string)   { s.lastMethod = "Warning"; s.lastText = text }
// ... demais métodos com corpo vazio
```
