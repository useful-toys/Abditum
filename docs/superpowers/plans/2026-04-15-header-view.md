# HeaderView Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implementar o `HeaderView` completo — cabeçalho de 2 linhas com abas clicáveis, indicador dirty, busca e suporte a mouse — substituindo o stub vazio atual.

**Architecture:** `WorkArea` é movida para `package design` para evitar import cycle. `HeaderView` é um consumidor passivo de contexto: setters injetados pelo `RootModel`, sem estado interno além de `tabPositions` para detecção de clique. Helpers de render são funções livres exportadas no `package screen`, cada uma retornando `(string, int)`.

**Tech Stack:** Go, `charm.land/bubbletea/v2`, `charm.land/lipgloss/v2`, golden file tests via `internal/tui/testdata`.

**Spec:** `docs/superpowers/specs/2026-04-15-header-view-design.md`

---

## Mapa de arquivos

| Arquivo | Ação | Responsabilidade |
|---|---|---|
| `internal/tui/design/workspace.go` | **Criar** | Tipo `WorkArea` e suas 4 constantes |
| `internal/tui/root.go` | **Modificar** | Remover `WorkArea` local; usar `design.WorkArea`; adicionar `setWorkArea`; handler `WorkAreaChangedMsg`; chamadas aos setters do header |
| `internal/tui/message.go` | **Modificar** | Adicionar `WorkAreaChangedMsg` |
| `internal/vault/manager.go` | **Modificar** | Adicionar `FilePath() string` getter |
| `internal/vault/manager_test.go` | **Criar** | `NewManagerForTest` — construtor de teste com caminho |
| `internal/tui/screen/header_view.go` | **Substituir** | Implementação completa: estado, setters, `RenderTab`, `RenderTabConnector`, `RenderTitleLine`, `RenderSeparatorLine`, `Render`, `Update` |
| `internal/tui/screen/header_view_test.go` | **Criar** | Golden tests: helpers individuais + componente completo (10 variantes) |
| `internal/tui/screen/testdata/golden/` | **Criar** | Arquivos `.golden.txt` e `.golden.json` gerados automaticamente |

---

## Task 1: Mover WorkArea para package design

**Files:**
- Create: `internal/tui/design/workspace.go`
- Modify: `internal/tui/root.go`

### Por quê esta task vem primeiro

`WorkArea` precisa estar em `package design` antes de qualquer outro arquivo novo ser escrito — `header_view.go` e `message.go` dependem dela.

- [ ] **Passo 1.1: Criar `internal/tui/design/workspace.go`**

```go
package design

// WorkArea representa qual área de trabalho está ativa na tela principal.
// É usada por RootModel para decidir qual ChildView exibir e pelo HeaderView
// para renderizar a aba ativa no cabeçalho.
type WorkArea int

const (
	// WorkAreaWelcome exibe a tela de boas-vindas, para usuários sem cofre aberto.
	WorkAreaWelcome WorkArea = iota
	// WorkAreaSettings exibe as configurações da aplicação.
	WorkAreaSettings
	// WorkAreaVault exibe a área de gerenciamento do cofre de segredos.
	WorkAreaVault
	// WorkAreaTemplates exibe a área de gerenciamento de templates de segredos.
	WorkAreaTemplates
)
```

- [ ] **Passo 1.2: Remover declaração local de `WorkArea` em `internal/tui/root.go`**

Remover as linhas 17–30 de `root.go` (o bloco `type WorkArea int` e as constantes `WorkAreaWelcome`, `WorkAreaSettings`, `WorkAreaVault`, `WorkAreaTemplates`).

- [ ] **Passo 1.3: Atualizar referências em `root.go`**

No campo da struct `RootModel`:
```go
// Antes:
workArea WorkArea
// Depois:
workArea design.WorkArea
```

No método `renderWorkArea()`, substituir os 4 `case` — de `WorkAreaWelcome`, `WorkAreaSettings`, `WorkAreaVault`, `WorkAreaTemplates` — para `design.WorkAreaWelcome`, `design.WorkAreaSettings`, `design.WorkAreaVault`, `design.WorkAreaTemplates`.

No método `NewRootModel()`, substituir `WorkAreaWelcome` por `design.WorkAreaWelcome`.

- [ ] **Passo 1.4: Verificar compilação**

```
go build ./...
```

Esperado: sem erros. Se aparecer `undefined: WorkAreaXxx`, alguma referência foi perdida — procure com `Select-String "WorkArea" internal\tui\root.go`.

- [ ] **Passo 1.5: Commit**

```
git add internal/tui/design/workspace.go internal/tui/root.go
git commit -m "refactor: move WorkArea to package design to avoid import cycle"
```

---

## Task 2: Adicionar WorkAreaChangedMsg

**Files:**
- Modify: `internal/tui/message.go`

- [ ] **Passo 2.1: Adicionar import e mensagem em `internal/tui/message.go`**

Adicionar import de `design` e o novo tipo no final do arquivo:

```go
package tui

import "github.com/useful-toys/abditum/internal/tui/design"

// TickMsg é emitido 1 vez por segundo pelo timer global do RootModel.
// ... (resto do arquivo sem alteração)

// WorkAreaChangedMsg é emitido pelo HeaderView quando o usuário clica numa aba.
// O RootModel processa e chama SetActiveMode no HeaderView e troca a WorkArea ativa.
type WorkAreaChangedMsg struct {
	Area design.WorkArea
}
```

- [ ] **Passo 2.2: Verificar compilação**

```
go build ./internal/tui/...
```

- [ ] **Passo 2.3: Commit**

```
git add internal/tui/message.go
git commit -m "feat: add WorkAreaChangedMsg for header tab click events"
```

---

## Task 3: Adicionar FilePath() ao vault.Manager e NewManagerForTest

**Files:**
- Modify: `internal/vault/manager.go`
- Create: `internal/vault/manager_export_test.go`

### Contexto

O header precisa de `m.FilePath()` para extrair o nome do cofre. Os golden tests precisam de um `*vault.Manager` com `caminho` não-vazio, mas `NewManager` sempre cria com `caminho = ""`. A solução é um construtor de teste exportado apenas no escopo de testes (`_test.go` no mesmo pacote).

- [ ] **Passo 3.1: Adicionar `FilePath()` em `internal/vault/manager.go`**

Logo após o método `IsModified()` (linha 180):

```go
// FilePath retorna o caminho completo do arquivo do cofre.
// Retorna "" se nenhum cofre estiver aberto (caminho não inicializado).
func (m *Manager) FilePath() string {
	return m.caminho
}
```

- [ ] **Passo 3.2: Criar `internal/vault/manager_export_test.go`**

```go
package vault

// NewManagerForTest cria um Manager com caminho explícito para uso em testes.
// Permite golden tests do HeaderView que precisam de um vault com nome de arquivo real.
func NewManagerForTest(cofre *Cofre, repositorio RepositorioCofre, caminho string) *Manager {
	return &Manager{
		cofre:       cofre,
		repositorio: repositorio,
		caminho:     caminho,
		bloqueado:   false,
	}
}
```

> **Nota:** O arquivo usa `package vault` (sem `_test`), então `NewManagerForTest` é visível apenas quando compilado como parte de testes (`go test`). Isso é o padrão Go para "export for test" — não polui a API pública.

- [ ] **Passo 3.3: Verificar compilação**

```
go build ./internal/vault/...
go test ./internal/vault/... -run ^$ -count=1
```

Esperado: sem erros.

- [ ] **Passo 3.4: Commit**

```
git add internal/vault/manager.go internal/vault/manager_export_test.go
git commit -m "feat: add FilePath getter and NewManagerForTest to vault.Manager"
```

---

## Task 4: Implementar header_view.go

**Files:**
- Modify: `internal/tui/screen/header_view.go` (substitui o stub)

### Contexto

O arquivo atual (`header_view.go`) é um stub com `HeaderView struct{}` e métodos vazios. Esta task o substitui completamente. Nenhum literal de cor, símbolo ou atributo tipográfico — sempre `design.SymXxx`, `theme.Xxx.Yyy`, e helpers do lipgloss.

- [ ] **Passo 4.1: Substituir `internal/tui/screen/header_view.go` com a implementação completa**

```go
package screen

import (
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/actions"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/vault"
)

// TabSpec descreve uma aba do cabeçalho.
type TabSpec struct {
	Label string
	Area  design.WorkArea
}

// DefaultTabs é a lista fixa de abas da aplicação, na ordem de exibição.
var DefaultTabs = []TabSpec{
	{Label: "Cofre", Area: design.WorkAreaVault},
	{Label: "Modelos", Area: design.WorkAreaTemplates},
	{Label: "Config", Area: design.WorkAreaSettings},
}

// tabSpacing é o número de espaços entre abas consecutivas na linha 1.
const tabSpacing = 2

// tabHitArea registra o intervalo horizontal de uma aba para detecção de clique.
type tabHitArea struct {
	area     design.WorkArea
	colStart int // coluna inicial (inclusive)
	colEnd   int // coluna final (exclusive)
}

// HeaderView renderiza o cabeçalho fixo de 2 linhas da aplicação.
// É um consumidor passivo: o RootModel injeta contexto via setters antes de cada render.
// Não possui eventos internos, filas nem timers.
type HeaderView struct {
	vault        *vault.Manager  // nil = sem cofre (estado boas-vindas)
	searchQuery  *string         // nil = busca inativa; non-nil = ativa (pode ser "")
	activeMode   design.WorkArea // relevante só quando vault != nil
	tabPositions []tabHitArea    // posições calculadas no último Render, para mouse
}

// NewHeaderView cria uma nova instância do cabeçalho.
func NewHeaderView() *HeaderView {
	return &HeaderView{}
}

// SetVault informa qual cofre está aberto.
// nil = sem cofre — renderiza apenas "  Abditum" sem abas.
func (v *HeaderView) SetVault(m *vault.Manager) {
	v.vault = m
}

// SetSearchQuery informa o estado de busca.
// nil = busca inativa. Non-nil = ativa (string pode ser "").
func (v *HeaderView) SetSearchQuery(q *string) {
	v.searchQuery = q
}

// SetActiveMode informa qual aba está ativa.
// Só tem efeito visual quando vault != nil.
func (v *HeaderView) SetActiveMode(mode design.WorkArea) {
	v.activeMode = mode
}

// Render retorna as 2 linhas do cabeçalho concatenadas com "\n",
// cada linha com exatamente `width` colunas.
// Atualiza v.tabPositions como efeito colateral para uso em Update.
func (v *HeaderView) Render(height, width int, theme *design.Theme) string {
	vaultName := vaultDisplayName(v.vault)
	isDirty := v.vault != nil && v.vault.IsModified()

	var tabs []TabSpec
	if v.vault != nil {
		tabs = DefaultTabs
	}

	// Linha 1: título + abas
	line1, _ := RenderTitleLine(vaultName, isDirty, tabs, v.activeMode, width, theme)

	// Calcular posições das abas para detecção de clique (efeito colateral).
	v.tabPositions = computeTabPositions(tabs, v.activeMode, width, theme)

	// Linha 2: separador com conector da aba ativa e busca (se ativa)
	line2, _ := RenderSeparatorLine(tabs, v.activeMode, v.searchQuery, width, theme)

	return line1 + "\n" + line2
}

// HandleKey não processa teclas nesta view.
func (v *HeaderView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }

// HandleEvent não processa eventos externos nesta view.
func (v *HeaderView) HandleEvent(event any) {}

// HandleTeaMsg não processa mensagens do framework nesta view.
func (v *HeaderView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }

// Update detecta cliques de mouse nas linhas Y=0 ou Y=1 e emite WorkAreaChangedMsg.
// Não atualiza v.activeMode — o RootModel faz isso via SetActiveMode.
func (v *HeaderView) Update(msg tea.Msg) tea.Cmd {
	switch m := msg.(type) {
	case tea.MouseClickMsg:
		if m.Button == tea.MouseLeft && (m.Y == 0 || m.Y == 1) {
			for _, hit := range v.tabPositions {
				if m.X >= hit.colStart && m.X < hit.colEnd {
					area := hit.area
					return func() tea.Msg {
						return WorkAreaChangedMsg{Area: area}
					}
				}
			}
		}
	}
	return nil
}

// Actions retorna nil — HeaderView não possui actions próprias.
func (v *HeaderView) Actions() []actions.Action { return nil }

// --- Helpers de render exportados ---

// RenderTab retorna a representação visual de uma aba na linha 1.
//
// Inativa: "╭ Label ╮" — ╭╮ em theme.Border.Default, label em theme.Text.Secondary.
// Ativa:   "╭──────╮" — ╭╮ e ─ em theme.Border.Default (mesmo ─ que a linha separadora).
//
// A largura total é idêntica nos dois estados.
func RenderTab(label string, active bool, theme *design.Theme) (string, int) {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Border.Default))
	tl := borderStyle.Render(design.SymCornerTL)
	tr := borderStyle.Render(design.SymCornerTR)

	// A largura interna é: 1 espaço + len(label) + 1 espaço = lipgloss.Width(label) + 2
	innerWidth := lipgloss.Width(label) + 2

	var inner string
	if active {
		fill := strings.Repeat(design.SymBorderH, innerWidth)
		inner = borderStyle.Render(fill)
	} else {
		labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Secondary))
		inner = " " + labelStyle.Render(label) + " "
	}

	rendered := tl + inner + tr
	return rendered, lipgloss.Width(rendered)
}

// RenderTabConnector retorna o fragmento da linha 2 que "suspende" a aba ativa.
//
// Formato: ╯[espaço + Label + espaço]╰
//   - ╯ e ╰ em theme.Border.Default, sem fundo
//   - conteúdo interno com fundo theme.Special.Highlight
//   - Label em theme.Accent.Primary bold
func RenderTabConnector(label string, theme *design.Theme) (string, int) {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Border.Default))
	br := borderStyle.Render(design.SymCornerBR) // ╯
	bl := borderStyle.Render(design.SymCornerBL) // ╰

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Accent.Primary)).
		Bold(true).
		Background(lipgloss.Color(theme.Special.Highlight))
	spaceStyle := lipgloss.NewStyle().Background(lipgloss.Color(theme.Special.Highlight))

	inner := spaceStyle.Render(" ") + labelStyle.Render(label) + spaceStyle.Render(" ")

	rendered := br + inner + bl
	return rendered, lipgloss.Width(rendered)
}

// RenderTitleLine monta a linha 1 completa com exatamente `width` colunas.
//
// Sem cofre (vaultName == ""):
//
//	"  Abditum" + padding até width
//
// Com cofre:
//
//	"  Abditum · nome •   [tab1]  [tab2]  [tab3]"
func RenderTitleLine(
	vaultName string,
	isDirty bool,
	tabs []TabSpec,
	activeMode design.WorkArea,
	width int,
	theme *design.Theme,
) (string, int) {
	appStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Accent.Primary)).Bold(true)
	appName := "  " + appStyle.Render("Abditum")

	if vaultName == "" {
		// Sem cofre: só o nome da aplicação, preenchido até width
		line := appName + strings.Repeat(" ", max(0, width-lipgloss.Width(appName)))
		return line, lipgloss.Width(line)
	}

	// Bloco de abas
	tabsBlock, tabsWidth := renderTabsBlock(tabs, activeMode, theme)

	// Cálculo de espaço disponível para o nome do cofre
	// Prefixo: "  Abditum · " = 2 espaços + "Abditum" (7) + " · " (3) = 12 colunas
	prefixWidth := lipgloss.Width("  Abditum · ")
	dirtyWidth := 0
	if isDirty {
		dirtyWidth = lipgloss.Width(" " + design.SymBullet)
	}
	const paddingMin = 1
	disponivel := width - prefixWidth - dirtyWidth - paddingMin - tabsWidth

	// Truncar nome do cofre se necessário
	name := truncateRight(vaultName, disponivel)

	// Montar prefixo completo
	sepStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Border.Default))
	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Secondary))
	dirtyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Semantic.Warning))

	prefix := appName +
		" " + sepStyle.Render(design.SymHeaderSep) + " " +
		nameStyle.Render(name)

	if isDirty {
		prefix += " " + dirtyStyle.Render(design.SymBullet)
	}

	// Padding entre nome/dirty e abas
	usedLeft := lipgloss.Width(prefix)
	padding := width - usedLeft - tabsWidth
	if padding < paddingMin {
		padding = paddingMin
	}

	line := prefix + strings.Repeat(" ", padding) + tabsBlock
	// Garantir largura exata
	lineWidth := lipgloss.Width(line)
	if lineWidth < width {
		line += strings.Repeat(" ", width-lineWidth)
	}
	return line, lipgloss.Width(line)
}

// RenderSeparatorLine monta a linha 2 completa com exatamente `width` colunas.
//
// Sem abas: SymBorderH repetido width vezes.
// Com abas, busca inativa: ─── + conector da aba ativa + ───
// Com abas, busca ativa: " ─ Busca: " + query + ─── + conector + ───
func RenderSeparatorLine(
	tabs []TabSpec,
	activeMode design.WorkArea,
	searchQuery *string,
	width int,
	theme *design.Theme,
) (string, int) {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Border.Default))
	hBorder := func(n int) string {
		if n <= 0 {
			return ""
		}
		return borderStyle.Render(strings.Repeat(design.SymBorderH, n))
	}

	// Sem abas: linha de borda horizontal simples
	if len(tabs) == 0 {
		line := hBorder(width)
		return line, lipgloss.Width(line)
	}

	// Encontrar aba ativa e sua posição horizontal na linha 1
	activeLabel := ""
	activeColStart := 0
	col := 0
	for i, tab := range tabs {
		if i > 0 {
			col += tabSpacing
		}
		_, tw := RenderTab(tab.Label, tab.Area == activeMode, theme)
		if tab.Area == activeMode {
			activeLabel = tab.Label
			activeColStart = col
		}
		col += tw
	}

	connector, connectorWidth := RenderTabConnector(activeLabel, theme)

	// Posição do conector na linha 2: alinhado com ╭ da aba ativa na linha 1.
	// O prefixo "  Abditum · nome • " tem uma largura variável, mas as abas começam
	// num offset que depende do width e da largura do bloco de abas.
	// Para simplificar: o conector começa na mesma coluna que a aba ativa.
	// Calculamos o offset das abas em relação ao início da linha usando a mesma
	// lógica de renderTabsBlock — as abas ficam à direita na linha 1.
	_, tabsWidth := renderTabsBlock(tabs, activeMode, theme)
	tabsStartCol := width - tabsWidth
	connectorStartCol := tabsStartCol + activeColStart
	// Ajuste para "  " que precede a primeira aba (já incluído em tabsBlock se
	// usarmos tabSpacing entre abas; aqui as abas não têm prefixo, então o
	// connectorStartCol é direto).

	if searchQuery == nil {
		// Busca inativa: ─── + conector + ───
		leftCount := connectorStartCol
		rightCount := width - connectorStartCol - connectorWidth
		if rightCount < 0 {
			rightCount = 0
		}
		line := hBorder(leftCount) + connector + hBorder(rightCount)
		lineWidth := lipgloss.Width(line)
		if lineWidth < width {
			line += hBorder(width - lineWidth)
		}
		return line, lipgloss.Width(line)
	}

	// Busca ativa
	// Prefixo: " ─ Busca: " (10 colunas)
	searchPrefix := " " + design.SymBorderH + " Busca: "
	searchPrefixRendered := borderStyle.Render(searchPrefix)
	const searchPrefixCols = 10

	disponivel := width - searchPrefixCols - connectorWidth
	query := truncateLeft(*searchQuery, disponivel)
	queryStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Accent.Primary)).Bold(true)
	queryRendered := queryStyle.Render(query)
	queryWidth := lipgloss.Width(queryRendered)

	fillCount := disponivel - queryWidth
	if fillCount < 0 {
		fillCount = 0
	}

	line := searchPrefixRendered + queryRendered + hBorder(fillCount) + connector
	lineWidth := lipgloss.Width(line)
	if lineWidth < width {
		line += hBorder(width - lineWidth)
	}
	return line, lipgloss.Width(line)
}

// --- Helpers internos ---

// vaultDisplayName extrai o nome de exibição do cofre a partir do caminho do arquivo.
// Retorna "" se m for nil.
func vaultDisplayName(m *vault.Manager) string {
	if m == nil {
		return ""
	}
	base := filepath.Base(m.FilePath())
	return strings.TrimSuffix(base, ".abditum")
}

// renderTabsBlock renderiza todas as abas em sequência, com tabSpacing entre elas.
// Retorna o bloco renderizado e sua largura visual total.
func renderTabsBlock(tabs []TabSpec, activeMode design.WorkArea, theme *design.Theme) (string, int) {
	var sb strings.Builder
	total := 0
	for i, tab := range tabs {
		if i > 0 {
			sb.WriteString(strings.Repeat(" ", tabSpacing))
			total += tabSpacing
		}
		rendered, w := RenderTab(tab.Label, tab.Area == activeMode, theme)
		sb.WriteString(rendered)
		total += w
	}
	return sb.String(), total
}

// computeTabPositions calcula os intervalos horizontais de cada aba na linha 1
// para uso na detecção de cliques do mouse.
func computeTabPositions(tabs []TabSpec, activeMode design.WorkArea, width int, theme *design.Theme) []tabHitArea {
	if len(tabs) == 0 {
		return nil
	}
	_, tabsWidth := renderTabsBlock(tabs, activeMode, theme)
	tabsStartCol := width - tabsWidth

	var positions []tabHitArea
	col := tabsStartCol
	for i, tab := range tabs {
		if i > 0 {
			col += tabSpacing
		}
		_, tw := RenderTab(tab.Label, tab.Area == activeMode, theme)
		positions = append(positions, tabHitArea{
			area:     tab.Area,
			colStart: col,
			colEnd:   col + tw,
		})
		col += tw
	}
	return positions
}

// truncateRight trunca s à direita para caber em maxCols colunas visuais.
// Adiciona design.SymEllipsis se necessário.
func truncateRight(s string, maxCols int) string {
	runes := []rune(s)
	if lipgloss.Width(s) <= maxCols {
		return s
	}
	if maxCols < 2 {
		return design.SymEllipsis
	}
	// Reduzir runa por runa até caber com ellipsis
	for len(runes) > 0 {
		candidate := string(runes) + design.SymEllipsis
		if lipgloss.Width(candidate) <= maxCols {
			return candidate
		}
		runes = runes[:len(runes)-1]
	}
	return design.SymEllipsis
}

// truncateLeft trunca s à esquerda para caber em maxCols colunas visuais.
// Adiciona design.SymEllipsis no início se necessário.
func truncateLeft(s string, maxCols int) string {
	runes := []rune(s)
	if lipgloss.Width(s) <= maxCols {
		return s
	}
	if maxCols < 2 {
		return design.SymEllipsis
	}
	// Reduzir do início até caber com ellipsis
	for len(runes) > 0 {
		candidate := design.SymEllipsis + string(runes)
		if lipgloss.Width(candidate) <= maxCols {
			return candidate
		}
		runes = runes[1:]
	}
	return design.SymEllipsis
}

// max retorna o maior de dois inteiros.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// WorkAreaChangedMsg é referenciada aqui para que o compilador valide o tipo.
// A definição real está em internal/tui/message.go (package tui).
// Importada indiretamente via root.go — o package screen não importa package tui.
// O tipo é definido em package tui mas emitido aqui; o RootModel faz o type switch.
// Para compilar sem import cycle, usamos interface{} internamente e o RootModel
// reconhece pela estrutura. Na prática: o package screen emite o tipo correto
// porque root.go importa screen e passa a mensagem como tea.Msg.
//
// ATENÇÃO: WorkAreaChangedMsg deve ser definido em package tui (message.go),
// não aqui. O Update acima retorna um func() tea.Msg que retorna tui.WorkAreaChangedMsg.
// Para isso funcionar sem import cycle, redefinimos o tipo localmente:

// WorkAreaChangedMsg é emitido pelo HeaderView quando o usuário clica numa aba.
// Definido aqui (package screen) para evitar import cycle com package tui.
// O RootModel em package tui faz o type switch neste tipo (importando package screen).
type WorkAreaChangedMsg struct {
	Area design.WorkArea
}
```

> **IMPORTANTE — Import cycle:** O `WorkAreaChangedMsg` deve ser definido em `package screen` (não em `package tui`), porque `package tui` importa `package screen`. Se fosse definido em `package tui`, o `header_view.go` precisaria importar `package tui`, criando um ciclo. O `RootModel` faz `case screen.WorkAreaChangedMsg:` no seu `Update`. **Remova o `WorkAreaChangedMsg` de `internal/tui/message.go` — não adicione lá.** (A Task 2 deste plano está cancelada.)

- [ ] **Passo 4.2: Verificar compilação**

```
go build ./internal/tui/screen/...
go build ./...
```

Esperado: sem erros.

- [ ] **Passo 4.3: Commit**

```
git add internal/tui/screen/header_view.go
git commit -m "feat: implement HeaderView with tab rendering, mouse click, and search support"
```

---

## Task 5: Integrar WorkAreaChangedMsg no RootModel

**Files:**
- Modify: `internal/tui/root.go`

### Contexto

Com `WorkAreaChangedMsg` definido em `package screen`, o `RootModel` precisa:
1. Adicionar `case screen.WorkAreaChangedMsg:` no `Update`
2. Implementar `setWorkArea` que troca a work area e chama `r.headerView.SetActiveMode`
3. Chamar `r.headerView.SetVault(...)` quando o cofre abre/fecha

- [ ] **Passo 5.1: Adicionar método `setWorkArea` em `root.go`**

Adicionar antes do método `Init`:

```go
// setWorkArea troca a área de trabalho ativa e sincroniza o estado do cabeçalho.
// Deve ser chamado sempre que a work area mudar — inclusive na abertura de cofre.
func (r *RootModel) setWorkArea(area design.WorkArea) {
	r.workArea = area
	r.headerView.SetActiveMode(area)
}
```

- [ ] **Passo 5.2: Adicionar `case screen.WorkAreaChangedMsg` no `Update`**

No método `Update`, após o `case CloseModalMsg` (ou antes do fall-through final), adicionar:

```go
case screen.WorkAreaChangedMsg:
    r.setWorkArea(msg.Area)
    return r, nil
```

- [ ] **Passo 5.3: Chamar `SetVault` quando cofre abre/fecha**

Localizar onde `r.vaultManager` é atribuído (abertura de cofre) e onde é zerado (fechamento). Em cada um desses pontos, adicionar chamada a `r.headerView.SetVault(r.vaultManager)`.

> Dica: buscar `r.vaultManager =` em `root.go` para encontrar os pontos de atribuição.

- [ ] **Passo 5.4: Substituir chamadas diretas de `r.workArea =` por `r.setWorkArea`**

Buscar `r.workArea =` no arquivo e substituir por `r.setWorkArea(...)`.

- [ ] **Passo 5.5: Verificar compilação e testes**

```
go build ./...
go test ./internal/tui/... -count=1
```

- [ ] **Passo 5.6: Commit**

```
git add internal/tui/root.go
git commit -m "feat: integrate WorkAreaChangedMsg handler and setWorkArea in RootModel"
```

---

## Task 6: Criar golden tests do HeaderView

**Files:**
- Create: `internal/tui/screen/header_view_test.go`

### Contexto

Os golden tests verificam tanto o texto limpo (`.golden.txt`) quanto as transições de estilo ANSI (`.golden.json`). Na primeira execução com `-update-golden`, os arquivos são criados. Em execuções subsequentes, a saída é comparada.

O `NewManagerForTest` está em `package vault` e é acessível em testes do `package screen` (que importa `package vault`).

- [ ] **Passo 6.1: Criar `internal/tui/screen/header_view_test.go`**

```go
package screen

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/testdata"
	"github.com/useful-toys/abditum/internal/vault"
)

// ptr é um helper para criar ponteiro de string em testes.
func ptr(s string) *string { return &s }

// newTestVault cria um vault.Manager mínimo para uso em golden tests.
// O cofre tem um nome (via caminho) mas sem segredos.
func newTestVault(caminho string) *vault.Manager {
	cofre := vault.NewCofre()
	return vault.NewManagerForTest(cofre, nil, caminho)
}

// headerRenderFn adapta HeaderView.Render à assinatura testdata.RenderFn.
func headerRenderFn(setup func(v *HeaderView)) testdata.RenderFn {
	return func(w, h int, theme *design.Theme) string {
		v := NewHeaderView()
		setup(v)
		return v.Render(h, w, theme)
	}
}

// --- Testes dos helpers individuais ---

func TestRenderTab_Inactive(t *testing.T) {
	testdata.TestRenderManaged(t, "header-tab", "inactive", []string{"10x1"},
		func(w, h int, theme *design.Theme) string {
			rendered, _ := RenderTab("Cofre", false, theme)
			// Garantir largura correta
			got := lipgloss.Width(rendered)
			_ = got
			return rendered
		},
	)
}

func TestRenderTab_Active(t *testing.T) {
	testdata.TestRenderManaged(t, "header-tab", "active", []string{"10x1"},
		func(w, h int, theme *design.Theme) string {
			rendered, _ := RenderTab("Cofre", true, theme)
			return rendered
		},
	)
}

func TestRenderTabConnector_Vault(t *testing.T) {
	testdata.TestRenderManaged(t, "header-connector", "vault", []string{"9x1"},
		func(w, h int, theme *design.Theme) string {
			rendered, _ := RenderTabConnector("Cofre", theme)
			return rendered
		},
	)
}

func TestRenderTabConnector_Models(t *testing.T) {
	testdata.TestRenderManaged(t, "header-connector", "models", []string{"11x1"},
		func(w, h int, theme *design.Theme) string {
			rendered, _ := RenderTabConnector("Modelos", theme)
			return rendered
		},
	)
}

func TestRenderTabConnector_Config(t *testing.T) {
	testdata.TestRenderManaged(t, "header-connector", "config", []string{"10x1"},
		func(w, h int, theme *design.Theme) string {
			rendered, _ := RenderTabConnector("Config", theme)
			return rendered
		},
	)
}

// --- Testes do componente completo ---

var headerGoldenSizes = []string{"80x2"}

func TestHeader_NoVault(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "no-vault", headerGoldenSizes,
		headerRenderFn(func(v *HeaderView) {
			v.SetVault(nil)
		}),
	)
}

func TestHeader_VaultClean(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "vault-clean", headerGoldenSizes,
		headerRenderFn(func(v *HeaderView) {
			v.SetVault(newTestVault("/home/user/meu-cofre.abditum"))
			v.SetActiveMode(design.WorkAreaVault)
		}),
	)
}

func TestHeader_VaultDirty(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "vault-dirty", headerGoldenSizes,
		headerRenderFn(func(v *HeaderView) {
			m := newTestVault("/home/user/meu-cofre.abditum")
			// Marcar como modificado via operação real no cofre
			// Como não há segredos, usamos SetSearchQuery para simular estado dirty?
			// NÃO — IsModified() reflete m.cofre.modificado.
			// Precisamos criar um segredo para marcar como dirty.
			// Alternativa: usar um modelo com modificado = true via NewManagerForTest.
			// SOLUÇÃO: NewManagerForTest aceita cofre já modificado.
			// Criar cofre e marcá-lo como modificado:
			cofre := vault.NewCofre()
			pasta, _ := cofre.CriarPasta("Test")
			_ = pasta
			m2 := vault.NewManagerForTest(cofre, nil, "/home/user/meu-cofre.abditum")
			// Criar segredo para marcar como modificado
			modelo := &vault.ModeloSegredo{}
			_ = modelo
			// Se não houver API simples para marcar dirty, usar campo direto via test helper.
			// Por ora, testar com vault clean (isDirty=false) — o golden file refletirá isso.
			// TODO: revisar quando houver API de marcação dirty de teste.
			_ = m2
			v.SetVault(m)
			v.SetActiveMode(design.WorkAreaVault)
		}),
	)
}
```

> **PARAR AQUI** — o bloco `TestHeader_VaultDirty` acima tem um problema: criar um vault "dirty" nos testes requer ou (a) uma API de teste que exponha `modificado`, ou (b) executar uma operação real que marque o vault como modificado. Leia o `vault.Manager` para encontrar a forma correta antes de continuar.

**Antes de escrever o arquivo final de teste, verificar:**
```
go doc github.com/useful-toys/abditum/internal/vault Cofre
go doc github.com/useful-toys/abditum/internal/vault Manager.CriarPasta
```

Se `CriarPasta` existe e retorna `(*Pasta, error)`, pode-se criar uma pasta para marcar o vault como dirty. Se não, pode-se expor um setter de test em `manager_export_test.go`.

Reescrever `TestHeader_VaultDirty` com a API real:

```go
func TestHeader_VaultDirty(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "vault-dirty", headerGoldenSizes,
		headerRenderFn(func(v *HeaderView) {
			m := newTestVaultDirty("/home/user/meu-cofre.abditum")
			v.SetVault(m)
			v.SetActiveMode(design.WorkAreaVault)
		}),
	)
}

// newTestVaultDirty cria um Manager com IsModified() == true.
// Usa a API pública do vault para criar um segredo e marcar como modificado.
func newTestVaultDirty(caminho string) *vault.Manager {
	cofre := vault.NewCofre()
	m := vault.NewManagerForTest(cofre, nil, caminho)
	// CriarPasta marca o vault como modificado (estadoSessao = EstadoModificado)
	_, _ = m.CriarPasta("_test_pasta")
	return m
}
```

- [ ] **Passo 6.2: Verificar API do vault para marcar dirty**

```
go doc github.com/useful-toys/abditum/internal/vault Manager
```

Confirmar se `CriarPasta` existe e marca como modificado. Ajustar `newTestVaultDirty` conforme necessário.

- [ ] **Passo 6.3: Completar o arquivo de teste com todas as variantes**

Adicionar os testes restantes após `TestHeader_VaultDirty`:

```go
func TestHeader_ModeVault(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "mode-vault", headerGoldenSizes,
		headerRenderFn(func(v *HeaderView) {
			v.SetVault(newTestVault("/home/user/meu-cofre.abditum"))
			v.SetActiveMode(design.WorkAreaVault)
		}),
	)
}

func TestHeader_ModeModels(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "mode-models", headerGoldenSizes,
		headerRenderFn(func(v *HeaderView) {
			v.SetVault(newTestVault("/home/user/meu-cofre.abditum"))
			v.SetActiveMode(design.WorkAreaTemplates)
		}),
	)
}

func TestHeader_ModeConfig(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "mode-config", headerGoldenSizes,
		headerRenderFn(func(v *HeaderView) {
			v.SetVault(newTestVault("/home/user/meu-cofre.abditum"))
			v.SetActiveMode(design.WorkAreaSettings)
		}),
	)
}

func TestHeader_VaultNameLong(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "vault-name-long", headerGoldenSizes,
		headerRenderFn(func(v *HeaderView) {
			v.SetVault(newTestVault("/home/user/meu-cofre-pessoal-muito-longo-que-nao-cabe.abditum"))
			v.SetActiveMode(design.WorkAreaVault)
		}),
	)
}

func TestHeader_SearchEmpty(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "search-empty", headerGoldenSizes,
		headerRenderFn(func(v *HeaderView) {
			v.SetVault(newTestVault("/home/user/meu-cofre.abditum"))
			v.SetActiveMode(design.WorkAreaVault)
			v.SetSearchQuery(ptr(""))
		}),
	)
}

func TestHeader_SearchWithQuery(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "search-with-query", headerGoldenSizes,
		headerRenderFn(func(v *HeaderView) {
			v.SetVault(newTestVault("/home/user/meu-cofre.abditum"))
			v.SetActiveMode(design.WorkAreaVault)
			v.SetSearchQuery(ptr("gmail"))
		}),
	)
}

func TestHeader_SearchQueryLong(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "search-query-long", headerGoldenSizes,
		headerRenderFn(func(v *HeaderView) {
			v.SetVault(newTestVault("/home/user/meu-cofre.abditum"))
			v.SetActiveMode(design.WorkAreaVault)
			v.SetSearchQuery(ptr("query muito longa que não cabe de jeito nenhum"))
		}),
	)
}
```

- [ ] **Passo 6.4: Verificar compilação dos testes**

```
go test ./internal/tui/screen/... -run ^$ -count=1
```

Esperado: compilação OK, 0 testes executados (nenhum golden file ainda).

- [ ] **Passo 6.5: Gerar golden files**

```
go test ./internal/tui/screen/... -run TestHeader -update-golden -v
go test ./internal/tui/screen/... -run TestRenderTab -update-golden -v
go test ./internal/tui/screen/... -run TestRenderTabConnector -update-golden -v
```

Esperado: arquivos `.golden.txt` e `.golden.json` criados em `internal/tui/screen/testdata/golden/`.

- [ ] **Passo 6.6: Inspecionar golden files gerados**

```
Get-Content internal\tui\screen\testdata\golden\header-no-vault-80x2.golden.txt
Get-Content internal\tui\screen\testdata\golden\header-vault-clean-80x2.golden.txt
```

Verificar visualmente se o layout está correto conforme `golden/tui-spec-cabecalho.md`.

- [ ] **Passo 6.7: Executar testes sem flag update para confirmar que passam**

```
go test ./internal/tui/screen/... -run TestHeader -v
go test ./internal/tui/screen/... -run TestRenderTab -v
go test ./internal/tui/screen/... -run TestRenderTabConnector -v
```

Esperado: todos PASS.

- [ ] **Passo 6.8: Commit**

```
git add internal/tui/screen/header_view_test.go internal/tui/screen/testdata/
git commit -m "test: add golden tests for HeaderView and render helpers"
```

---

## Task 7: Suite completa e validação final

**Files:** nenhum novo

- [ ] **Passo 7.1: Executar suite completa**

```
go test ./... -count=1
```

Esperado: todos PASS. Se houver falhas em outros pacotes, investigar se são causadas pelas mudanças desta sprint (especialmente em `root.go`).

- [ ] **Passo 7.2: Build final**

```
go build ./...
```

Esperado: sem erros.

- [ ] **Passo 7.3: Commit final (se houver ajustes)**

```
git add -A
git commit -m "chore: final integration and test verification for HeaderView"
```

---

## Notas de implementação importantes

### WorkAreaChangedMsg — definição em package screen

Ao contrário do que a spec original indicava, `WorkAreaChangedMsg` deve ser definido em `package screen` (não em `package tui`), porque:
- `package tui` importa `package screen`
- Se `WorkAreaChangedMsg` fosse em `package tui`, `header_view.go` precisaria importar `package tui` → ciclo

O `RootModel.Update` faz `case screen.WorkAreaChangedMsg:`.

### NewManagerForTest — visibilidade

`NewManagerForTest` está em `internal/vault/manager_export_test.go` com `package vault` (sem `_test`). Isso significa que é compilado apenas quando `go test` processa o pacote `vault`, mas é exportado para qualquer pacote de teste que importe `vault`. Os testes em `package screen` podem usá-lo.

### Tamanhos dos golden tests dos helpers

Os tamanhos `10x1`, `9x1`, `11x1`, `10x1` para os helpers de tab/connector são aproximações. O `testdata.TestRenderManaged` usa o tamanho apenas para chamar `render(w, h, theme)` — os helpers não usam `w` nem `h` (retornam largura fixa). O size string aparece apenas no nome do arquivo golden. Use os valores da spec; se a largura visual for diferente, os golden files refletirão o valor real.

### `vault.NewCofre()` — verificar existência

Antes de escrever o código de teste, confirmar que `vault.NewCofre()` existe e é exportado:
```
go doc github.com/useful-toys/abditum/internal/vault NewCofre
```

Se não existir, ajustar `newTestVault` para usar outra forma de criar um `*vault.Cofre`.
