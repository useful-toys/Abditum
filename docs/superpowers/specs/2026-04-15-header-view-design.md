# Design — HeaderView

**Data:** 2026-04-15  
**Spec de referência:** `golden/tui-spec-cabecalho.md`  
**Design system:** `golden/tui-design-system.md`

---

## Contexto

O `HeaderView` em `internal/tui/screen/header_view.go` é atualmente um stub vazio. Este documento especifica sua implementação completa.

**Padrão adotado:** helpers de render exportados como funções livres no pacote `screen`, cada um retornando `(string, int)` (texto renderizado + largura visual). Esse padrão é novo em relação a `ActionLineView` e `MessageLineView` — que encapsulam toda a lógica em métodos privados — e é introduzido aqui porque as sub-partes do header (abas, conectores, linha de título, linha separadora) têm valor de teste independente e são reutilizáveis.

**Restrição absoluta:** nenhum literal de cor, símbolo ou atributo tipográfico na implementação. Sempre `design.SymXxx`, `theme.Xxx.Yyy`, e helpers do lipgloss.

---

## Resolução do import cycle

`WorkArea` está atualmente em `package tui` (`internal/tui/root.go`). `HeaderView` está em `package screen` (`internal/tui/screen/`). Como `package tui` importa `package screen`, usar `tui.WorkArea` em `screen` criaria um ciclo.

**Solução:** mover `WorkArea` e suas constantes para `package design` (`internal/tui/design/`), que já é importado por ambos os pacotes sem ciclo.

```go
// Em internal/tui/design/workspace.go (arquivo novo)
type WorkArea int

const (
    WorkAreaWelcome   WorkArea = iota
    WorkAreaSettings
    WorkAreaVault
    WorkAreaTemplates
)
```

`internal/tui/root.go` passa a importar `design.WorkArea` em vez de declarar localmente. Sem quebra de API externa — apenas mudança de pacote origem.

---

## Estrutura do componente

### Estado

O `HeaderView` é quase stateless. O RootModel injeta contexto via setters antes de cada ciclo de render. Não há eventos internos, filas nem timers.

```go
type HeaderView struct {
    vault        *vault.Manager   // nil = sem cofre (estado boas-vindas)
    searchQuery  *string          // nil = busca inativa; non-nil = ativa (pode ser "")
    activeMode   design.WorkArea  // relevante só quando vault != nil
    tabPositions []tabHitArea     // posições calculadas no último Render, para mouse
}

// tabHitArea registra o intervalo horizontal de uma aba para detecção de clique.
type tabHitArea struct {
    area       design.WorkArea
    colStart   int // coluna inicial (inclusive)
    colEnd     int // coluna final (exclusive)
}
```

**Invariante:** quando `vault == nil`, `activeMode` é ignorado — nenhuma aba é exibida. `SetActiveMode` só tem efeito observável quando `vault != nil`.

### Setters (chamados pelo RootModel)

```go
// SetVault informa qual cofre está aberto.
// nil = sem cofre — renderiza apenas "  Abditum" sem abas.
// Quando nil, activeMode é ignorado no render.
func (v *HeaderView) SetVault(m *vault.Manager)

// SetSearchQuery informa o estado de busca.
// nil = busca inativa. Non-nil = ativa (string pode ser "").
func (v *HeaderView) SetSearchQuery(q *string)

// SetActiveMode informa qual aba está ativa.
// Só tem efeito visual quando vault != nil.
func (v *HeaderView) SetActiveMode(mode design.WorkArea)
```

### Render

```go
// Render retorna as 2 linhas do cabeçalho concatenadas com "\n",
// cada linha com exatamente `width` colunas.
// Atualiza v.tabPositions como efeito colateral para uso em Update.
// A constante design.HeaderHeight (= 2) define a altura esperada.
func (v *HeaderView) Render(height, width int, theme *design.Theme) string
```

`Render` atualiza `v.tabPositions` como efeito colateral — registra as posições horizontais de cada aba para que `Update` possa mapear cliques. Isso é seguro no modelo single-threaded do bubbletea (Render e Update nunca executam concorrentemente).

### Mouse (Update) — comportamento command-driven

`Update(msg tea.Msg)` detecta `tea.MouseClickMsg` (botão esquerdo) nas linhas Y=0 ou Y=1 (0-indexed; o header ocupa sempre as duas primeiras linhas da tela). Usa `v.tabPositions` para determinar qual aba foi clicada.

Ao detectar clique válido numa aba, **emite `tea.Cmd` com `WorkAreaChangedMsg`** e **não atualiza `v.activeMode` internamente**. O header aguarda o RootModel processar a mensagem e chamar `SetActiveMode` na próxima iteração. Isso evita estado duplicado e mantém o header como consumidor passivo de contexto.

```go
func (v *HeaderView) Update(msg tea.Msg) tea.Cmd {
    switch m := msg.(type) {
    case tea.MouseClickMsg:
        if m.Button == tea.MouseLeft && (m.Y == 0 || m.Y == 1) {
            for _, hit := range v.tabPositions {
                if m.X >= hit.colStart && m.X < hit.colEnd {
                    return func() tea.Msg {
                        return WorkAreaChangedMsg{Area: hit.area}
                    }
                }
            }
        }
    }
    return nil
}
```

`HandleKey`, `HandleEvent`, `HandleTeaMsg`, `Actions` permanecem sem implementação (retornam nil/zero).

---

## Mensagem nova: WorkAreaChangedMsg

Definir em `internal/tui/message.go` (junto com as demais mensagens):

```go
// WorkAreaChangedMsg é emitido pelo HeaderView quando o usuário clica numa aba.
// O RootModel processa e chama SetActiveMode no HeaderView e troca a WorkArea ativa.
type WorkAreaChangedMsg struct {
    Area design.WorkArea
}
```

O `RootModel.Update` precisa de um novo case para essa mensagem:

```go
case WorkAreaChangedMsg:
    r.setWorkArea(m.Area) // ou equivalente — troca a work area ativa e atualiza o header
```

---

## Adição necessária ao vault.Manager

O `vault.Manager` não expõe o caminho do arquivo (campo privado `caminho`). Adicionar getter:

```go
// FilePath retorna o caminho completo do arquivo do cofre.
// Retorna "" se nenhum cofre estiver aberto (caminho não inicializado).
func (m *Manager) FilePath() string {
    return m.caminho
}
```

O header extrai o nome de exibição assim (guardando nil):

```go
func vaultDisplayName(m *vault.Manager) string {
    if m == nil {
        return ""
    }
    base := filepath.Base(m.FilePath())
    return strings.TrimSuffix(base, ".abditum")
}
```

O método relevante para o indicador dirty é `m.IsModified() bool` — **não** `IsDirty()`, que não existe.

---

## Helpers de render

Todas as funções retornam `(string, int)`. Nenhum literal hardcoded.

### RenderTab

```go
// RenderTab retorna a representação visual de uma aba na linha 1.
//
// Inativa: "╭ Label ╮"
//   ╭╮ em theme.Border.Default
//   Label em theme.Text.Secondary
//
// Ativa: "╭──────╮"  (preenchimento SymBorderH no lugar do label)
//   ╭╮ em theme.Border.Default
//   ─ (SymBorderH) em theme.Border.Default
//   Quantidade de ─: lipgloss.Width("╭ Label ╮") - 2 (os dois cantos)
//
// A largura total é idêntica nos dois estados.
func RenderTab(label string, active bool, theme *design.Theme) (string, int)
```

### RenderTabConnector

```go
// RenderTabConnector retorna o fragmento da linha 2 que "suspende" a aba ativa.
//
// Formato: ╯[espaço + Label + espaço]╰
//   ╯ e ╰ em theme.Border.Default, sem fundo
//   conteúdo interno (espaço + Label + espaço) com fundo theme.Special.Highlight
//   Label em theme.Accent.Primary bold
//
// O fundo Special.Highlight cobre APENAS o conteúdo entre os cantos,
// não os próprios cantos ╯ e ╰.
func RenderTabConnector(label string, theme *design.Theme) (string, int)
```

### RenderTitleLine

```go
// RenderTitleLine monta a linha 1 completa com exatamente `width` colunas.
//
// Sem cofre (vaultName == ""):
//   "  Abditum" + padding até width
//   Abditum em theme.Accent.Primary bold
//
// Com cofre:
//   "  Abditum · nome •   [tab1]  [tab2]  [tab3]"
//   Abditum → theme.Accent.Primary bold
//   · (SymHeaderSep) → theme.Border.Default
//   nome → theme.Text.Secondary (truncado com SymEllipsis se necessário)
//   • (SymBullet) → theme.Semantic.Warning (omitido se !isDirty)
//   Abas via RenderTab; 2 espaços entre abas (constante tabSpacing = 2)
//
// O padding mínimo entre o bloco nome/dirty e o bloco de abas é 1 coluna.
// O nome é o primeiro a ceder espaço; abas nunca truncam.
func RenderTitleLine(
    vaultName  string,
    isDirty    bool,
    tabs       []TabSpec,
    activeMode design.WorkArea,
    width      int,
    theme      *design.Theme,
) (string, int)
```

### RenderSeparatorLine

```go
// RenderSeparatorLine monta a linha 2 completa com exatamente `width` colunas.
//
// Sem abas (sem cofre):
//   SymBorderH repetido width vezes, em theme.Border.Default
//
// Com abas, busca inativa (searchQuery == nil):
//   SymBorderH preenchendo à esquerda até a posição da aba ativa
//   RenderTabConnector(activeLabel)
//   SymBorderH preenchendo à direita até width
//
// Com abas, busca ativa (searchQuery != nil):
//   " ─ Busca: " (10 colunas — espaço + SymBorderH + " Busca: ") em theme.Border.Default
//   query em theme.Accent.Primary bold (truncada à esquerda se necessária — ver algoritmo)
//   SymBorderH preenchendo entre query e conector da aba ativa
//   RenderTabConnector(activeLabel)
//   SymBorderH preenchendo à direita até width (pode ser zero se aba encosta na borda)
//
// Tokens: SymBorderH → theme.Border.Default
func RenderSeparatorLine(
    tabs        []TabSpec,
    activeMode  design.WorkArea,
    searchQuery *string,
    width       int,
    theme       *design.Theme,
) (string, int)
```

**Atenção:** o prefixo de busca tem **10 colunas** (` ─ Busca: ` = espaço + `─` + ` Busca: `), não 9 como indicado na spec visual. A contagem correta é: `" "` (1) + `"─"` (1) + `" "` (1) + `"Busca:"` (6) + `" "` (1) = 10.

### TabSpec e constantes de layout

```go
// TabSpec descreve uma aba do cabeçalho.
type TabSpec struct {
    Label string
    Area  design.WorkArea
}

// DefaultTabs é a lista fixa de abas da aplicação, na ordem de exibição.
var DefaultTabs = []TabSpec{
    {Label: "Cofre",   Area: design.WorkAreaVault},
    {Label: "Modelos", Area: design.WorkAreaTemplates},
    {Label: "Config",  Area: design.WorkAreaSettings},
}

const tabSpacing = 2 // espaços entre abas consecutivas na linha 1
```

---

## Algoritmo de truncamento do nome do cofre (linha 1)

```
prefixoCols = lipgloss.Width("  Abditum · ")   // 12 colunas (fixo)
dirtyCols   = lipgloss.Width(" •")              // 2 colunas se isDirty, 0 se não
              // o espaço faz parte do dirtyCols — é o espaço entre nome e •
tabsWidth   = soma de RenderTab(label, active, theme).width para cada aba
              + tabSpacing * (len(tabs) - 1)    // espaços entre abas
paddingMin  = 1                                 // mínimo entre dirty/nome e abas

disponível = width - prefixoCols - dirtyCols - paddingMin - tabsWidth
```

**Algoritmo de truncamento:**
1. Se `runeCount(radical) <= disponível` → exibir como está
2. Se não cabe e `disponível >= 2` → `radical[0..disponível-1] + design.SymEllipsis`
3. Se `disponível < 2` → exibir apenas `design.SymEllipsis`

> Usar `runeCount` (não `len`) para contagem de caracteres; usar `lipgloss.Width` para largura visual ao comparar com `disponível`.

---

## Algoritmo de truncamento da query de busca (linha 2, à esquerda)

A query é truncada **à esquerda** — a parte mais recente sempre fica visível.

```
prefixoBuscaCols = 10                           // " ─ Busca: "
connectorWidth   = RenderTabConnector(label).width
disponívelQuery  = width - prefixoBuscaCols - connectorWidth
                   // sem margem direita: aba pode encostar na borda direita
```

**Algoritmo:**
1. Se `lipgloss.Width(query) <= disponívelQuery` → exibir como está
2. Se não cabe e `disponívelQuery >= 2` → `design.SymEllipsis + query[len-n:]`
   onde `n` é o número de runas que cabem em `disponívelQuery - 1` colunas
3. Se `disponívelQuery < 2` → exibir apenas `design.SymEllipsis`

O preenchimento `SymBorderH` entre query e conector ocupa `disponívelQuery - lipgloss.Width(queryExibida)` colunas.

---

## Mecânica visual da aba ativa

| Linha | Estado inativo | Estado ativo |
|---|---|---|
| **1** | `╭ Label ╮` (bordas + texto) | `╭──────╮` (bordas + SymBorderH no lugar do texto) |
| **2** | `─────────` (SymBorderH contínuo) | `╯[espaço Label espaço]╰` (conector com fundo highlight) |

Regras de alinhamento:
- Largura total de `RenderTab` é idêntica nos estados ativo e inativo
- `╯` alinha-se verticalmente com `╭` da linha acima (mesma coluna)
- `╰` alinha-se verticalmente com `╮` da linha acima (mesma coluna)
- O fundo `theme.Special.Highlight` cobre **apenas** o conteúdo interno entre `╯` e `╰` — os cantos em si ficam sem fundo highlight, apenas com `theme.Border.Default` no foreground

---

## Identidade visual

| Elemento | Token | Atributo |
|---|---|---|
| `Abditum` | `theme.Accent.Primary` | bold |
| `·` (`design.SymHeaderSep`) | `theme.Border.Default` | — |
| Nome do cofre | `theme.Text.Secondary` | — |
| `•` dirty (`design.SymBullet`) | `theme.Semantic.Warning` | — |
| Bordas das abas (`╭╮╯╰` = `SymCornerTL SymCornerTR SymCornerBR SymCornerBL`) | `theme.Border.Default` | — |
| Fill ativo `─` (`design.SymBorderH`) | `theme.Border.Default` | — |
| Aba ativa — fundo interno | `theme.Special.Highlight` | — |
| Aba ativa — texto | `theme.Accent.Primary` | bold |
| Aba inativa — texto | `theme.Text.Secondary` | — |
| Linha separadora (`design.SymBorderH`) | `theme.Border.Default` | — |
| ` ─ Busca: ` rótulo | `theme.Border.Default` | — |
| Texto da query | `theme.Accent.Primary` | bold |
| `…` truncamento (`design.SymEllipsis`) | mesmo estilo do elemento truncado | — |

> **Regra:** nenhum literal de cor, símbolo ou atributo tipográfico na implementação. Sempre `design.SymXxx`, `theme.Xxx.Yyy`, e helpers do lipgloss.

---

## Testes com golden files

**Arquivo:** `internal/tui/screen/header_view_test.go`  
**Fixtures:** `internal/tui/screen/testdata/golden/header-*`

### Adapter

```go
func headerRenderFn(setup func(v *HeaderView)) testdata.RenderFn {
    return func(w, h int, theme *design.Theme) string {
        v := NewHeaderView()
        setup(v)
        return v.Render(h, w, theme)
    }
}
```

`RenderFn` tem assinatura `func(w, h int, theme *design.Theme) string`. O adapter chama `v.Render(h, w, theme)` — nota a inversão: `Render` recebe `(height, width)` enquanto `RenderFn` fornece `(w, h)`.

### Helpers individuais

Cada variante usa seu próprio size slice (não uma variável compartilhada):

| Componente | Variante | Size | Label usado |
|---|---|---|---|
| `header-tab` | `inactive` | `10x1` | `"Cofre"` (9 cols: `╭ Cofre ╮`) |
| `header-tab` | `active` | `10x1` | `"Cofre"` (9 cols: `╭─────╮`) |
| `header-connector` | `vault` | `9x1` | `"Cofre"` (9 cols: `╯ Cofre ╰`) |
| `header-connector` | `models` | `11x1` | `"Modelos"` |
| `header-connector` | `config` | `10x1` | `"Config"` |

> **Nota sobre sizes:** os tamanhos acima são aproximados — o implementador deve calcular a largura real de cada helper com `lipgloss.Width` e usar esse valor como size nos golden tests. Os valores acima ilustram a ordem de grandeza.

### Componente completo

```go
var headerGoldenSizes = []string{"80x2"}
```

| Componente | Variante | Setup |
|---|---|---|
| `header` | `no-vault` | `SetVault(nil)` |
| `header` | `vault-clean` | `SetVault(m)`, `SetActiveMode(WorkAreaVault)`, `m.IsModified() == false` |
| `header` | `vault-dirty` | `SetVault(m)`, `SetActiveMode(WorkAreaVault)`, vault com alterações |
| `header` | `mode-vault` | `SetVault(m)`, `SetActiveMode(design.WorkAreaVault)` |
| `header` | `mode-models` | `SetVault(m)`, `SetActiveMode(design.WorkAreaTemplates)` |
| `header` | `mode-config` | `SetVault(m)`, `SetActiveMode(design.WorkAreaSettings)` |
| `header` | `vault-name-long` | `SetVault(m)` com caminho longo (ex: `meu-cofre-pessoal-muito-longo.abditum`) |
| `header` | `search-empty` | `SetVault(m)`, `SetActiveMode(design.WorkAreaVault)`, `SetSearchQuery(ptr(""))` |
| `header` | `search-with-query` | `SetVault(m)`, `SetActiveMode(design.WorkAreaVault)`, `SetSearchQuery(ptr("gmail"))` |
| `header` | `search-query-long` | `SetVault(m)`, `SetActiveMode(design.WorkAreaVault)`, `SetSearchQuery(ptr("query muito longa que não cabe"))` |

Helper para ponteiro de string nos testes:
```go
func ptr(s string) *string { return &s }
```

Golden files gerados com:
```sh
go test ./internal/tui/screen/... -run TestHeader -update-golden
```

---

## Arquivos afetados

| Arquivo | Ação |
|---|---|
| `internal/tui/design/workspace.go` | **Criar** — move `WorkArea` e constantes de `root.go` para cá |
| `internal/tui/root.go` | Remover declaração local de `WorkArea`; usar `design.WorkArea` |
| `internal/vault/manager.go` | Adicionar `FilePath() string` |
| `internal/tui/message.go` | Adicionar `WorkAreaChangedMsg` |
| `internal/tui/root.go` | Adicionar handler `case WorkAreaChangedMsg` em `Update` |
| `internal/tui/screen/header_view.go` | Implementação completa (substitui stub) |
| `internal/tui/screen/header_view_test.go` | Criar — testes golden |
| `internal/tui/screen/testdata/golden/header-*.golden.*` | Criar — 10 variantes × 2 formatos |
| `internal/tui/screen/testdata/golden/header-tab-*.golden.*` | Criar — 2 variantes × 2 formatos |
| `internal/tui/screen/testdata/golden/header-connector-*.golden.*` | Criar — 3 variantes × 2 formatos |
