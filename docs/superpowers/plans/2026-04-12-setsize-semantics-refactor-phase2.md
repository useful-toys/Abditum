# SetSize Semantics Refactor — Phase 2 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Eliminar `SetSize`/`SetAvailableSize` das interfaces `childModel` e `modalView`, fazendo `View()` receber dimensões como parâmetros — removendo o acoplamento implícito de estado de tamanho dos modelos.

**Architecture:** A interface `childModel.View()` passa a ser `View(width, height int) string`; a interface `modalView.View()` passa a ser `View(maxWidth, maxHeight int) string`. O `rootModel` passa as dimensões diretamente ao chamar `View()` em `renderFrame()`. Dois modelos com scroll (`helpModal`, `filePickerModal`) mantêm um campo `viewportHeight int` no estado, calculado em `View()` e usado em `Update()`. O `rootModel` protege `child.Update()` e `modal.Update()` de serem chamados antes de `m.width > 0`.

**Tech Stack:** Go, BubbleTea v2 (`charm.land/bubbletea/v2`), Lipgloss v2 (`charm.land/lipgloss/v2`)

---

## Convenções

- Todos os comandos rodam de `C:\git\Abditum-T2`
- Rodar testes: `go test ./internal/tui/... -count=1`
- Rodar testes com update de golden: `go test ./internal/tui/... -count=1 -update`
- Compilar: `go build ./...`

---

## Visão geral das mudanças por arquivo

| Arquivo | O que muda |
|---|---|
| `flows.go` | `childModel.View() string` → `View(width, height int) string`; remover `SetSize`; `modalView.View() string` → `View(maxWidth, maxHeight int) string`; remover `SetAvailableSize` |
| `welcome.go` | Remover campos `width`, `height`; remover `SetSize`; `View()` → `View(width, height int)` |
| `vaulttree.go` | Remover campos `width`, `height`; remover `SetSize`; `View()` → `View(width, height int)` |
| `secretdetail.go` | Remover campos `width`, `height`; remover `SetSize`; `View()` → `View(width, height int)` |
| `templatelist.go` | Remover campos `width`, `height`; remover `SetSize`; `View()` → `View(width, height int)` |
| `templatedetail.go` | Remover campos `width`, `height`; remover `SetSize`; `View()` → `View(width, height int)` |
| `settings.go` | Remover campos `width`, `height`; remover `SetSize`; `View()` → `View(width, height int)` |
| `modal.go` | Remover `SetAvailableSize`; `View()` → `View(maxWidth, maxHeight int)` (não usa os params) |
| `help.go` | Remover campos `width`, `height`; renomear para `viewportHeight int`; remover `SetAvailableSize`; `View()` → `View(maxWidth, maxHeight int)`; `contentHeight()` usa `m.viewportHeight` |
| `decision.go` | Remover campo `height`; manter `width`; remover `SetAvailableSize`; `View()` → `View(maxWidth, maxHeight int)` |
| `filepicker.go` | Remover campos `width`, `height`; adicionar `viewportHeight int`; remover `SetAvailableSize` (migrar lógica de reset de scroll); `View()` → `View(maxWidth, maxHeight int)` |
| `passwordcreate.go` | Remover campo `height`; manter `width`; remover `SetAvailableSize`; `View()` → `View(maxWidth, maxHeight int)` |
| `passwordentry.go` | Remover campo `height`; manter `width`; remover `SetAvailableSize`; `View()` → `View(maxWidth, maxHeight int)` |
| `root.go` | Remover loop `SetSize` do `WindowSizeMsg`; remover `SetAvailableSize` de `View()`; remover `SetSize` de `renderFrame()`/`renderVaultArea()`/`renderTemplatesArea()`/`enterVault()`; passar dims para `child.View(w,h)` e `modal.View(w,h)`; guard `m.width > 0` antes de propagar para `child.Update()`/`modal.Update()` |
| `welcome_test.go` | `SetSize` → params em `View()`; remover `TestWelcomeModel_SetSize`; remover `TestWelcomeModel_ViewPanicsWithoutSetSize` |
| `vaulttree_test.go` | Remover `TestVaultTreeModel_ViewPanicsWithoutSetSize` |
| `secretdetail_test.go` | Remover `TestSecretDetailModel_ViewPanicsWithoutSetSize` |
| `templatelist_test.go` | Remover `TestTemplateListModel_ViewPanicsWithoutSetSize` |
| `templatedetail_test.go` | Remover `TestTemplateDetailModel_ViewPanicsWithoutSetSize` |
| `settings_test.go` | Remover `TestSettingsModel_ViewPanicsWithoutSetSize` |
| `help_test.go` | `SetAvailableSize(w,h)` → params em `View(w,h)`; `contentHeight()` chamado após `View()`; remover `TestHelpModal_ViewPanicsWithoutSetAvailableSize` |
| `decision_test.go` | `SetAvailableSize(w,h)` → params em `View(w,h)`; remover `TestDecisionDialog_ViewPanicsWithoutSetAvailableSize` |
| `filepicker_test.go` | `SetAvailableSize(w,h)` → params em `View(w,h)`; remover `TestFilePickerModal_ViewPanicsWithoutSetAvailableSize` |
| `passwordcreate_test.go` | `SetAvailableSize(w,h)` → params em `View(w,h)`; remover `TestPasswordCreateModal_ViewPanicsWithoutSetAvailableSize` |
| `passwordentry_test.go` | `SetAvailableSize(w,h)` → params em `View(w,h)`; remover `TestPasswordEntryModal_ViewPanicsWithoutSetAvailableSize` |
| `root_test.go` | `stubModal.View()` → `View(maxWidth, maxHeight int)`; remover `stubModal.SetAvailableSize`; atualizar `TestWindowSizeMsg_NoModalSetSize` |

---

## Task 1: Atualizar interfaces em `flows.go`

**Files:**
- Modify: `internal/tui/flows.go:18-52`

Esta é a mudança raiz. Todos os outros tasks dependem dela; por isso vem primeiro — mas só compila após todos os implementadores serem atualizados (Tasks 2–9). Faça esta task junto com as tasks 2–9 antes de compilar.

- [ ] **Step 1: Atualizar `childModel` — remover `SetSize`, mudar assinatura de `View`**

```go
// childModel represents a UI component managed by rootModel.
type childModel interface {
    Update(tea.Msg) tea.Cmd
    View(width, height int) string
    ApplyTheme(*Theme)
}
```

- [ ] **Step 2: Atualizar `modalView` — remover `SetAvailableSize`, mudar assinatura de `View`**

```go
// modalView represents an overlay modal dialog.
type modalView interface {
    Update(tea.Msg) tea.Cmd
    View(maxWidth, maxHeight int) string
    Shortcuts() []Shortcut
}
```

> **Não compilar ainda** — implementadores precisam ser atualizados (Tasks 2–9).

---

## Task 2: Atualizar `welcome.go`

**Files:**
- Modify: `internal/tui/welcome.go`

- [ ] **Step 1: Remover campos `width` e `height` do struct; remover `SetSize`; mudar assinatura de `View`**

Substituir o arquivo inteiro por:

```go
package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// welcomeModel renders the welcome background (ASCII art logo + version + action hints).
// It is active during workAreaWelcome and has no sub-states.
// Open/create vault flows are orchestrated via the modal stack, not this model.
type welcomeModel struct {
	actions *ActionManager
	theme   *Theme
	version string // Application version to display below logo
}

// ApplyTheme applies the given theme to the welcomeModel.
func (m *welcomeModel) ApplyTheme(t *Theme) {
	m.theme = t
}

// Compile-time assertion: welcomeModel satisfies childModel.
var _ childModel = &welcomeModel{}

// newWelcomeModel creates a new welcome screen model.
func newWelcomeModel(actions *ActionManager, theme *Theme, version string) *welcomeModel {
	return &welcomeModel{actions: actions, theme: theme, version: version}
}

// Update processes messages for the welcome screen.
// Phase 5.1: welcomeModel is display-only. No input handling until Phase 6.
func (m *welcomeModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// View renders the ASCII art logo centered on screen.
// Per spec (tui-specification-novo.md § Boas-vindas), the logo and version
// are centered horizontally and vertically via lipgloss.Place().
// Logo width is hardcoded to 43 columns matching the ASCII art width.
// Version is displayed below the logo in text.secondary color.
func (m *welcomeModel) View(width, height int) string {
	// 43 = width of AsciiArt (const in ascii.go) — each line is exactly 43 characters.
	// No background is set here: the root workAreaStyle already applies SurfaceBase
	// to the entire work area. Setting background here would emit redundant SGR codes
	// that may conflict with the terminal's own background rendering.
	logoBlock := lipgloss.NewStyle().Width(43).Render(RenderLogo(m.theme))

	// Format version with semantic.secondary color (from theme)
	// Per spec: version token = text.secondary
	versionStyle := lipgloss.NewStyle().Foreground(m.theme.TextSecondary)
	versionLine := versionStyle.Render(m.version)

	content := lipgloss.JoinVertical(lipgloss.Center, logoBlock, "", versionLine)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}
```

> Note: o import de `"fmt"` é removido (não há mais panic guard).

---

## Task 3: Atualizar os 4 stubs childModel (`vaulttree`, `secretdetail`, `templatelist`, `templatedetail`, `settings`)

**Files:**
- Modify: `internal/tui/vaulttree.go`
- Modify: `internal/tui/secretdetail.go`
- Modify: `internal/tui/templatelist.go`
- Modify: `internal/tui/templatedetail.go`
- Modify: `internal/tui/settings.go`

Todos têm a mesma estrutura. Para cada um:
1. Remover campos `width int` e `height int` do struct
2. Remover o método `SetSize`
3. Mudar assinatura de `View() string` para `View(width, height int) string`
4. Remover o panic guard de `View()`
5. Remover o import de `"fmt"` (não é mais necessário)

- [ ] **Step 1: Atualizar `vaulttree.go`**

```go
package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/vault"
)

// vaultTreeModel is the left-panel child model active during workAreaVault.
// Stub in Phase 5 - real tree in Phase 7.
type vaultTreeModel struct {
	mgr     *vault.Manager
	actions *ActionManager
	msgs    *MessageManager
	theme   *Theme
}

// ApplyTheme applies the given theme to the vaultTreeModel.
func (m *vaultTreeModel) ApplyTheme(t *Theme) {
	m.theme = t
}

// Compile-time assertion: vaultTreeModel satisfies childModel.
var _ childModel = &vaultTreeModel{}

// newVaultTreeModel creates a new vault tree stub.
func newVaultTreeModel(mgr *vault.Manager, actions *ActionManager, msgs *MessageManager, theme *Theme) *vaultTreeModel {
	return &vaultTreeModel{mgr: mgr, actions: actions, msgs: msgs, theme: theme}
}

// Update processes messages for the vault tree.
func (m *vaultTreeModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// View renders a placeholder for the vault tree panel.
func (m *vaultTreeModel) View(width, height int) string {
	return lipgloss.NewStyle().Foreground(m.theme.SemanticInfo).
		Render("[vault tree - Phase 7]")
}
```

- [ ] **Step 2: Atualizar `secretdetail.go`** (idêntico, troca render string para `"[secret detail - Phase 8]"`)

```go
package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/vault"
)

type secretDetailModel struct {
	mgr     *vault.Manager
	actions *ActionManager
	msgs    *MessageManager
	theme   *Theme
}

func (m *secretDetailModel) ApplyTheme(t *Theme) { m.theme = t }

var _ childModel = &secretDetailModel{}

func newSecretDetailModel(mgr *vault.Manager, actions *ActionManager, msgs *MessageManager, theme *Theme) *secretDetailModel {
	return &secretDetailModel{mgr: mgr, actions: actions, msgs: msgs, theme: theme}
}

func (m *secretDetailModel) Update(msg tea.Msg) tea.Cmd { return nil }

func (m *secretDetailModel) View(width, height int) string {
	return lipgloss.NewStyle().Foreground(m.theme.SemanticInfo).
		Render("[secret detail - Phase 8]")
}
```

- [ ] **Step 3: Atualizar `templatelist.go`** (idêntico, render string `"[template list - Phase 8]"`)

```go
package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/vault"
)

type templateListModel struct {
	mgr     *vault.Manager
	actions *ActionManager
	msgs    *MessageManager
	theme   *Theme
}

func (m *templateListModel) ApplyTheme(t *Theme) { m.theme = t }

var _ childModel = &templateListModel{}

func newTemplateListModel(mgr *vault.Manager, actions *ActionManager, msgs *MessageManager, theme *Theme) *templateListModel {
	return &templateListModel{mgr: mgr, actions: actions, msgs: msgs, theme: theme}
}

func (m *templateListModel) Update(msg tea.Msg) tea.Cmd { return nil }

func (m *templateListModel) View(width, height int) string {
	return lipgloss.NewStyle().Foreground(m.theme.SemanticInfo).
		Render("[template list - Phase 8]")
}
```

- [ ] **Step 4: Atualizar `templatedetail.go`** (idêntico, render string `"[template detail - Phase 8]"`)

```go
package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/vault"
)

type templateDetailModel struct {
	mgr     *vault.Manager
	actions *ActionManager
	msgs    *MessageManager
	theme   *Theme
}

func (m *templateDetailModel) ApplyTheme(t *Theme) { m.theme = t }

var _ childModel = &templateDetailModel{}

func newTemplateDetailModel(mgr *vault.Manager, actions *ActionManager, msgs *MessageManager, theme *Theme) *templateDetailModel {
	return &templateDetailModel{mgr: mgr, actions: actions, msgs: msgs, theme: theme}
}

func (m *templateDetailModel) Update(msg tea.Msg) tea.Cmd { return nil }

func (m *templateDetailModel) View(width, height int) string {
	return lipgloss.NewStyle().Foreground(m.theme.SemanticInfo).
		Render("[template detail - Phase 8]")
}
```

- [ ] **Step 5: Atualizar `settings.go`** — ler o arquivo real primeiro para preservar campos extras

> `settings.go` tem a mesma estrutura dos outros stubs. Remover `width`, `height`, `SetSize`, panic guard; mudar `View() string` → `View(width, height int) string`; remover import `"fmt"`.

---

## Task 4: Atualizar `modal.go`

**Files:**
- Modify: `internal/tui/modal.go:66,99-108`

`modalModel.View()` usa `boxW := 50` hardcoded — não precisa dos parâmetros, mas precisa aceitar a nova assinatura.

- [ ] **Step 1: Mudar assinatura de `View` e remover `SetAvailableSize`**

Substituir as linhas 64–108:

```go
// View renders the modal box. Returns only the box - rootModel positions it.
func (m *modalModel) View(maxWidth, maxHeight int) string {
	boxW := 50

	var content strings.Builder
	if m.body != "" {
		content.WriteString(m.body + "\n\n")
	}

	selectedStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11"))
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

	for i, opt := range m.options {
		prefix := "  "
		if i == m.selectedIndex {
			content.WriteString(selectedStyle.Render(fmt.Sprintf("\u25b6 %s", opt)) + "\n")
		} else {
			content.WriteString(normalStyle.Render(prefix+opt) + "\n")
		}
	}
	if len(m.options) == 0 {
		content.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("245")).
			Render("  Press Enter or ESC to close") + "\n")
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Width(boxW).
		Render(content.String())
}

// Shortcuts returns nil - basic dialog modals show no command bar shortcuts.
func (m *modalModel) Shortcuts() []Shortcut { return nil }

// Compile-time assertion: modalModel must satisfy the modalView interface.
var _ modalView = &modalModel{}
```

> O método `SetAvailableSize` (linha 105) é simplesmente removido.

---

## Task 5: Atualizar `help.go`

**Files:**
- Modify: `internal/tui/help.go:100-123,158-207,284-293`

`helpModal` mantém `viewportHeight int` no estado porque `Update()` (linhas 139,141,145,148) chama `contentHeight()` que usa esse valor para calcular scroll em PgUp/PgDown/End/clamp.

- [ ] **Step 1: Atualizar o struct — trocar `width`/`height` por `viewportHeight`**

Substituir linhas 100–106:

```go
type helpModal struct {
	actions       []Action         // all registered actions for the help overlay
	groupLabel    func(int) string // resolves a group int to a display label
	viewportHeight int             // usable content lines; set by View(), used by Update()
	scroll        int              // current scroll offset
}
```

- [ ] **Step 2: Remover `SetAvailableSize` (linhas 119–123)**

Deletar completamente:

```go
// SetAvailableSize sets the maximum available dimensions for dynamic modal sizing.
func (m *helpModal) SetAvailableSize(maxWidth, maxHeight int) {
	m.width = maxWidth
	m.height = maxHeight
}
```

- [ ] **Step 3: Atualizar `contentHeight()` para usar `m.viewportHeight` diretamente (linhas 284–293)**

```go
// contentHeight returns the visible content height (usable action lines).
// Value is set by View() each render; safe to use in Update() per rootModel invariant.
func (m *helpModal) contentHeight() int {
	return m.viewportHeight
}
```

- [ ] **Step 4: Atualizar `View()` — nova assinatura, calcular `viewportHeight`, salvar no estado**

Substituir linhas 156–207:

```go
// View renders the help modal with title in top border and action in bottom border.
// Follows DS dialog anatomy (§436-458): title embedded in top border, action bar in bottom border.
// Calculates and saves viewportHeight for use by Update() (pgup/pgdown/end/clamp).
func (m *helpModal) View(maxWidth, maxHeight int) string {
	// Dynamic sizing per DS: max 60 cols or 70% of terminal
	maxW := 60
	if maxWidth > 0 {
		pctW := int(float64(maxWidth) * 0.7)
		if pctW < maxW {
			maxW = pctW
		}
	}
	boxW := maxW

	allActions := m.actions
	lines := m.buildContentLines(allActions)
	totalLines := len(lines)

	// Dialog layout: top border(1) + content(innerH) + bottom border(1)
	// Content area: top padding(1) + action lines(usableH) + bottom padding(1)
	// Total dialog = usableH + 4 lines — must fit in terminal.
	maxUsable := maxHeight - 4
	if maxUsable > 20 {
		maxUsable = 20 // cap for large terminals
	}
	if maxUsable < 3 {
		maxUsable = 3 // minimum usable action lines
	}
	usableH := maxUsable

	// Save viewport height for Update() (pgup/pgdown/end/clamp)
	m.viewportHeight = usableH

	innerH := usableH + 2 // content area includes padding lines

	// Clamp visible window to available content
	start := m.scroll
	if start > totalLines-usableH {
		start = totalLines - usableH
	}
	if start < 0 {
		start = 0
	}
	end := start + usableH
	if end > totalLines {
		end = totalLines
	}
	visibleLines := lines[start:end]

	hasAbove := start > 0
	hasBelow := end < totalLines

	return m.renderDialog(visibleLines, boxW, innerH, hasAbove, hasBelow, totalLines, start, usableH)
}
```

- [ ] **Step 5: Remover import `"fmt"` se não for mais usado**

Verificar os imports de `help.go` e remover `"fmt"` (o panic guard que o usava foi removido).

---

## Task 6: Atualizar `decision.go`

**Files:**
- Modify: `internal/tui/decision.go:57-65,121-124,171-230`

`DecisionDialog` usa `d.width` em `boxWidth()` (linha 320+). `d.height` nunca é lido — só armazenado. Remover `height`, manter `width`, remover `SetAvailableSize`, mudar assinatura de `View`.

- [ ] **Step 1: Remover campo `height` do struct (linha 64)**

```go
type DecisionDialog struct {
	title     string
	body      string
	severity  Severity
	intention Intention
	actions   []DecisionAction // ordered: default first, cancel last
	width     int
}
```

- [ ] **Step 2: Remover `SetAvailableSize` (linhas 121–124)**

Deletar completamente:

```go
func (m *DecisionDialog) SetAvailableSize(maxWidth, maxHeight int) {
	m.width = maxWidth
	m.height = maxHeight
}
```

- [ ] **Step 3: Atualizar `View()` — nova assinatura, usar `maxWidth` em vez de `m.width`**

Substituir a assinatura e a primeira linha de uso de `m.width`:

Linha original (171): `func (m *DecisionDialog) View() string {`
→ `func (m *DecisionDialog) View(maxWidth, maxHeight int) string {`

Em seguida, logo no início de `View()`, adicionar:
```go
m.width = maxWidth
```

> Isso mantém `boxWidth()` funcionando sem alteração, pois ela continua lendo `m.width`.

---

## Task 7: Atualizar `filepicker.go`

**Files:**
- Modify: `internal/tui/filepicker.go:42-76,82-90,402-414,1118-1122`

`filePickerModal` mantém `viewportHeight int` no estado porque `Update()` chama `visibleTreeHeight()` que usa `m.height`. A lógica de reset de scroll que estava em `SetAvailableSize` (linhas 88–89) migra para `View()`.

- [ ] **Step 1: Atualizar o struct — trocar `width`/`height` por `viewportHeight` (linhas 73–75)**

```go
	// Dimensions — viewportHeight is computed in View() and stored for Update() (scroll)
	viewportHeight int
```

(Remover as duas linhas `width int` e `height int`; adicionar `viewportHeight int`.)

- [ ] **Step 2: Remover `SetAvailableSize` (linhas 81–90)**

Deletar completamente:

```go
// SetAvailableSize stores the maximum available dimensions for use in View().
func (m *filePickerModal) SetAvailableSize(maxWidth, maxHeight int) {
	m.width = maxWidth
	m.height = maxHeight
	// Reset scroll to 0 then re-adjust so the cursor is visible with the real
	// viewport height. Without the reset, a bogus scroll computed during Init
	// (when height was 0) would persist and hide parent nodes above the cursor.
	m.treeScroll = 0
	m.adjustTreeScroll()
}
```

- [ ] **Step 3: Atualizar `visibleTreeHeight()` para usar `viewportHeight` (linhas 402–414)**

```go
// visibleTreeHeight returns the number of rows available for tree/file rows inside the modal.
// Returns m.viewportHeight which is calculated and stored by View() each render.
func (m *filePickerModal) visibleTreeHeight() int {
	return m.viewportHeight
}
```

> `visibleFilesHeight()` (linha 416) já delega para `visibleTreeHeight()` — não precisa mudar.

- [ ] **Step 4: Atualizar `View()` — nova assinatura, calcular dimensões, salvar `viewportHeight`, fazer reset de scroll na primeira vez**

Substituir linhas 1118–1153 (abertura de `View()` até o fim do bloco de layout dimensions):

```go
func (m *filePickerModal) View(maxWidth, maxHeight int) string {
	// Nil-safe theme fallback (D-17)
	theme := m.theme
	if theme == nil {
		theme = ThemeTokyoNight
	}

	// Layout dimensions (D-08)
	modalW := maxWidth * 95 / 100
	if modalW < 60 {
		modalW = 60
	}
	modalH := maxHeight * 8 / 10
	if modalH < 6 {
		modalH = 6
	}
	innerW := modalW - 2
	treeW := innerW * 40 / 100
	if treeW < 8 {
		treeW = 8
	}
	filesW := innerW - treeW - 1
	if filesW < 8 {
		filesW = 8
	}
	visibleH := modalH - 4 // top border + Caminho + panel sep + bottom border
	if m.mode == FilePickerSave {
		visibleH -= 3 // field separator + field row + bottom padding
	}
	if visibleH < 1 {
		visibleH = 1
	}

	// First render: reset scroll so cursor is visible with real viewport height.
	// (Previously done in SetAvailableSize; safe here because rootModel guarantees
	// View() is called before any Update() keypress.)
	if m.viewportHeight == 0 {
		m.treeScroll = 0
		m.adjustTreeScroll()
	}

	// Save viewport height for Update() (pgup/pgdown scroll calculations)
	m.viewportHeight = visibleH
```

> O restante de `View()` (linhas 1154–1195) permanece inalterado, pois já usa a variável local `visibleH`.

---

## Task 8: Atualizar `passwordcreate.go` e `passwordentry.go`

**Files:**
- Modify: `internal/tui/passwordcreate.go:34-44,177-313,316-319`
- Modify: `internal/tui/passwordentry.go:27-36,106-215,218-221`

Ambos usam só `width` (para clamp de `boxW`). `height` nunca é lido.

- [ ] **Step 1: Atualizar `passwordcreate.go` — remover `height` do struct, remover `SetAvailableSize`, mudar `View`**

No struct (linha 34–44), remover `height int`:
```go
type passwordCreateModal struct {
	// ... outros campos ...
	width int
}
```

Remover `SetAvailableSize` (linhas 316–319) completamente.

Mudar assinatura de `View()` (linha 177):
```go
func (m *passwordCreateModal) View(maxWidth, maxHeight int) string {
```

Logo no início de `View()`, adicionar:
```go
m.width = maxWidth
```

- [ ] **Step 2: Atualizar `passwordentry.go` — remover `height` do struct, remover `SetAvailableSize`, mudar `View`**

No struct (linhas 27–36), remover `height int`.

Remover `SetAvailableSize` (linhas 218–221) completamente.

Mudar assinatura de `View()` (linha 106):
```go
func (m *passwordEntryModal) View(maxWidth, maxHeight int) string {
```

Logo no início de `View()`, adicionar:
```go
m.width = maxWidth
```

---

## Task 9: Atualizar `root.go`

**Files:**
- Modify: `internal/tui/root.go:192-198,416-502,521-577`

Esta é a task mais complexa. São 4 mudanças independentes no mesmo arquivo.

- [ ] **Step 1: Remover loop `SetSize` do handler `WindowSizeMsg` (linhas 192–198)**

```go
case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
    return m, nil
```

(Remover as 3 linhas do `for _, child := range m.liveWorkChildren()` que chamavam `child.SetSize`.)

- [ ] **Step 2: Adicionar guard `m.width > 0` nos pontos de propagação para `child.Update()` e `modal.Update()`**

Em `Update()`, localizar todos os pontos onde o código propaga para `child.Update()` ou `modal.Update()`:

**Keypress — helpModal (linha 322):**
```go
if len(m.modals) > 0 {
    if _, isHelp := m.modals[len(m.modals)-1].(*helpModal); isHelp {
        if m.width > 0 {
            return m, m.modals[len(m.modals)-1].Update(msg)
        }
        return m, nil
    }
}
```

**Keypress — outros modais (linha 331):**
```go
if len(m.modals) > 0 {
    if m.width > 0 {
        return m, m.modals[len(m.modals)-1].Update(msg)
    }
    return m, nil
}
```

**Keypress — active child (linha 338–340):**
```go
if child := m.activeChild(); child != nil {
    if m.width > 0 {
        return m, child.Update(msg)
    }
    return m, nil
}
```

> **Nota:** `activeFlow.Update()` (linhas 242, 262, 334–335) **não** recebe guard — flows podem precisar rodar antes do terminal dimensionar. Só `child.Update()` e `modal.Update()` são protegidos.

- [ ] **Step 3: Atualizar `View()` — remover `SetAvailableSize`, passar dims para `modal.View()`**

Substituir linhas 416–436:

```go
func (m *rootModel) View() tea.View {
    content := m.renderFrame(nil)

    if len(m.modals) > 0 {
        const headerH = 2
        const msgBarH = 1
        const cmdBarH = 1
        workH := m.height - headerH - msgBarH - cmdBarH
        if workH < 0 {
            workH = 0
        }
        top := m.modals[len(m.modals)-1]
        content = m.renderFrame(top)
    }

    v := tea.NewView(content)
    v.AltScreen = true
    v.BackgroundColor = m.theme.SurfaceBase
    return v
}
```

- [ ] **Step 4: Atualizar `renderFrame()` — remover `SetSize` inline, passar dims para `child.View()` e `modal.View()`**

Substituir linhas 478–502:

```go
// Work area
var workContent string
switch m.area {
case workAreaWelcome:
    if m.welcome != nil {
        workContent = m.welcome.View(m.width, workH)
    }
case workAreaVault:
    workContent = m.renderVaultArea(workH)
case workAreaTemplates:
    workContent = m.renderTemplatesArea(workH)
case workAreaSettings:
    if m.settings != nil {
        workContent = m.settings.View(m.width, workH)
    }
}
workArea := workAreaStyle.Height(workH).Render(workContent)

// Overlay modal centered inside work area using lipgloss.Place.
if modal != nil {
    modalStr := modal.View(m.width, workH)
    workArea = lipgloss.Place(m.width, workH, lipgloss.Center, lipgloss.Center, modalStr)
}
```

> O comentário na linha 458 que menciona `SetAvailableSize` também deve ser atualizado:
> `// The modal receives dimensions via View(width, workH).`

- [ ] **Step 5: Atualizar `renderVaultArea()` — remover `SetSize`, usar `child.View()` com params (linhas 521–542)**

```go
func (m *rootModel) renderVaultArea(workH int) string {
    halfW := m.width / 2

    left := "[vault tree - Phase 7]"
    right := "[secret detail - Phase 8]"
    if m.vaultTree != nil {
        left = m.vaultTree.View(halfW, workH)
    }
    if m.secretDetail != nil {
        right = m.secretDetail.View(m.width-halfW, workH)
    }

    leftStyle := lipgloss.NewStyle().Width(halfW).Height(workH)
    rightStyle := lipgloss.NewStyle().Width(m.width - halfW).Height(workH)
    return lipgloss.JoinHorizontal(lipgloss.Top, leftStyle.Render(left), rightStyle.Render(right))
}
```

- [ ] **Step 6: Atualizar `renderTemplatesArea()` — remover `SetSize`, usar `child.View()` com params (linhas 545–566)**

```go
func (m *rootModel) renderTemplatesArea(workH int) string {
    halfW := m.width / 2

    left := "[template list - Phase 8]"
    right := "[template detail - Phase 8]"
    if m.templateList != nil {
        left = m.templateList.View(halfW, workH)
    }
    if m.templateDetail != nil {
        right = m.templateDetail.View(m.width-halfW, workH)
    }

    leftStyle := lipgloss.NewStyle().Width(halfW).Height(workH)
    rightStyle := lipgloss.NewStyle().Width(m.width - halfW).Height(workH)
    return lipgloss.JoinHorizontal(lipgloss.Top, leftStyle.Render(left), rightStyle.Render(right))
}
```

- [ ] **Step 7: Atualizar `enterVault()` — remover `SetSize` (linhas 570–578)**

```go
func (m *rootModel) enterVault() tea.Cmd {
    m.area = workAreaVault
    m.welcome = nil // GC old model
    m.vaultTree = newVaultTreeModel(m.mgr, m.actions, m.messages, m.theme)
    m.secretDetail = newSecretDetailModel(m.mgr, m.actions, m.messages, m.theme)
    return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}
```

- [ ] **Step 8: Verificar que compila**

```
go build ./...
```

Expected: sem erros de compilação.

---

## Task 10: Atualizar testes dos stubs `childModel`

**Files:**
- Modify: `internal/tui/vaulttree_test.go`
- Modify: `internal/tui/secretdetail_test.go`
- Modify: `internal/tui/templatelist_test.go`
- Modify: `internal/tui/templatedetail_test.go`
- Modify: `internal/tui/settings_test.go`

Cada arquivo tem apenas um teste (`TestXxx_ViewPanicsWithoutSetSize`) que testa o panic guard que foi removido. Substituir cada arquivo pelo teste equivalente para a nova semântica.

- [ ] **Step 1: Substituir `vaulttree_test.go`**

```go
package tui

import "testing"

// TestVaultTreeModel_View verifies that View() returns a non-empty string when called with dimensions.
func TestVaultTreeModel_View(t *testing.T) {
    vm := newVaultTreeModel(nil, nil, nil, ThemeTokyoNight)
    out := vm.View(80, 24)
    if out == "" {
        t.Error("View(80, 24) returned empty string")
    }
}
```

- [ ] **Step 2: Substituir `secretdetail_test.go`**

```go
package tui

import "testing"

// TestSecretDetailModel_View verifies that View() returns a non-empty string when called with dimensions.
func TestSecretDetailModel_View(t *testing.T) {
    sdm := newSecretDetailModel(nil, nil, nil, ThemeTokyoNight)
    out := sdm.View(80, 24)
    if out == "" {
        t.Error("View(80, 24) returned empty string")
    }
}
```

- [ ] **Step 3: Substituir `templatelist_test.go`**

```go
package tui

import "testing"

// TestTemplateListModel_View verifies that View() returns a non-empty string when called with dimensions.
func TestTemplateListModel_View(t *testing.T) {
    tlm := newTemplateListModel(nil, nil, nil, ThemeTokyoNight)
    out := tlm.View(80, 24)
    if out == "" {
        t.Error("View(80, 24) returned empty string")
    }
}
```

- [ ] **Step 4: Substituir `templatedetail_test.go`**

```go
package tui

import "testing"

// TestTemplateDetailModel_View verifies that View() returns a non-empty string when called with dimensions.
func TestTemplateDetailModel_View(t *testing.T) {
    tdm := newTemplateDetailModel(nil, nil, nil, ThemeTokyoNight)
    out := tdm.View(80, 24)
    if out == "" {
        t.Error("View(80, 24) returned empty string")
    }
}
```

- [ ] **Step 5: Substituir `settings_test.go`**

```go
package tui

import "testing"

// TestSettingsModel_View verifies that View() returns a non-empty string when called with dimensions.
func TestSettingsModel_View(t *testing.T) {
    sm := newSettingsModel(nil, nil, nil, ThemeTokyoNight)
    out := sm.View(80, 24)
    if out == "" {
        t.Error("View(80, 24) returned empty string")
    }
}
```

---

## Task 11: Atualizar `welcome_test.go`

**Files:**
- Modify: `internal/tui/welcome_test.go`

Três mudanças:
1. `wm.SetSize(80, 24)` → removido; `wm.View()` → `wm.View(80, 24)`
2. `TestWelcomeModel_SetSize` → remover (método não existe mais)
3. `TestWelcomeModel_ViewPanicsWithoutSetSize` → remover

- [ ] **Step 1: Atualizar todos os `SetSize` + `View()` no arquivo**

Substituir `wm.SetSize(80, 24)` por nada (remover a linha), e `wm.View()` por `wm.View(80, 24)` em todo o arquivo.

- [ ] **Step 2: Remover `TestWelcomeModel_SetSize` (linhas 37–45)**

Deletar completamente:

```go
// TestWelcomeModel_SetSize stores terminal dimensions.
func TestWelcomeModel_SetSize(t *testing.T) {
    wm := newWelcomeModel(nil, ThemeTokyoNight, "v0.1.0")
    wm.SetSize(80, 24)

    if wm.width != 80 || wm.height != 24 {
        t.Errorf("SetSize failed: expected 80x24, got %dx%d", wm.width, wm.height)
    }
}
```

- [ ] **Step 3: Remover `TestWelcomeModel_ViewPanicsWithoutSetSize` (linhas 130–140)**

Deletar completamente.

- [ ] **Step 4: Rodar testes de welcome**

```
go test ./internal/tui/... -run TestWelcomeModel -count=1
```

Expected: PASS.

---

## Task 12: Atualizar `help_test.go`

**Files:**
- Modify: `internal/tui/help_test.go`

`SetAvailableSize(w, h)` é chamado em múltiplos lugares. A nova abordagem: `m.View(w, h)` é chamado antes de qualquer operação que dependa de `viewportHeight`, e os testes de Update já chamam `View()` primeiro para inicializar o estado.

- [ ] **Step 1: Atualizar `TestHelpModal_Golden` — substituir `SetAvailableSize` por `View` com params**

No helper `make15` (linhas 108–122):
```go
make15 := func(w, h, scroll int) *helpModal {
    m := newHelpModal(help15actions(), helpGroupLabel)
    m.View(w, h) // initialize viewportHeight
    if scroll < 0 {
        maxScroll := m.totalLines() - m.contentHeight()
        if maxScroll < 0 {
            maxScroll = 0
        }
        m.scroll = maxScroll
    } else {
        m.scroll = scroll
    }
    return m
}
```

Nas cases inline que usam `m.SetAvailableSize(30, 12)` / `m.SetAvailableSize(60, 12)`:
```go
{
    variant: "3actions-30x12",
    modal: func() *helpModal {
        m := newHelpModal(help3actions(), helpGroupLabel)
        m.View(30, 12) // initialize viewportHeight
        return m
    }(),
},
```

No loop do test:
```go
out := tc.modal.View(/* já foi chamado no setup — usar as dims salvas */
```

> **Atenção:** Como `View()` agora recebe params e não usa estado, o golden test pode chamar `tc.modal.View(w, h)` passando as mesmas dims. Mas o helper já chamou `View(w, h)` para inicializar `viewportHeight`. A segunda chamada (`out := tc.modal.View(...)`) precisa das mesmas dimensões.

Melhor abordagem: guardar as dims no helper de setup e passar para o loop. Refatorar o test case struct para incluir `width, height int`:

```go
type testCase struct {
    variant string
    modal   *helpModal
    width   int
    height  int
}
```

E nos casos:
```go
{variant: "3actions-30x12", modal: newHelpModal(help3actions(), helpGroupLabel), width: 30, height: 12},
{variant: "3actions-60x12", modal: newHelpModal(help3actions(), helpGroupLabel), width: 60, height: 12},
{variant: "15actions-top-30x16", modal: make15base(30, 16, 0), width: 30, height: 16},
// etc.
```

Com `make15base` que NÃO chama `View()`:
```go
make15base := func(w, h, scroll int) *helpModal {
    m := newHelpModal(help15actions(), helpGroupLabel)
    // Initialize viewportHeight first to compute contentHeight for maxScroll
    m.View(w, h)
    if scroll < 0 {
        maxScroll := m.totalLines() - m.contentHeight()
        if maxScroll < 0 { maxScroll = 0 }
        m.scroll = maxScroll
    } else {
        m.scroll = scroll
    }
    return m
}
```

No loop:
```go
out := tc.modal.View(tc.width, tc.height)
```

- [ ] **Step 2: Atualizar `TestHelpModal_Update_*` — substituir `SetAvailableSize` por `View` chamado antes**

Em cada teste de Update que usa `m.SetAvailableSize(60, 16)`:

```go
// Antes:
m := newHelpModal(help15actions(), helpGroupLabel)
m.SetAvailableSize(60, 16)
m.Update(...)

// Depois:
m := newHelpModal(help15actions(), helpGroupLabel)
m.View(60, 16) // initialize viewportHeight
m.Update(...)
```

Isso aplica-se a todos os `TestHelpModal_Update_*` (linhas 175–270).

Para `TestHelpModal_Update_EndScrollsToBottom` e `TestHelpModal_Update_ScrollClampedAtMax`, que chamam `m.contentHeight()` após `SetAvailableSize`:

```go
m := newHelpModal(help15actions(), helpGroupLabel)
m.View(60, 16) // initialize viewportHeight
maxScroll := m.totalLines() - m.contentHeight()
```

Isso funciona porque `View()` salva `viewportHeight`, e `contentHeight()` agora só retorna `m.viewportHeight`.

- [ ] **Step 3: Remover `TestHelpModal_ViewPanicsWithoutSetAvailableSize` (linha 322)**

Deletar completamente.

- [ ] **Step 4: Rodar testes de help**

```
go test ./internal/tui/... -run TestHelpModal -count=1
```

Expected: PASS.

---

## Task 13: Atualizar `decision_test.go`

**Files:**
- Modify: `internal/tui/decision_test.go`

- [ ] **Step 1: Substituir `SetAvailableSize(w, h)` + `View()` por `View(w, h)` em todas as ocorrências**

Buscar no arquivo por `SetAvailableSize` e substituir o padrão:
```go
// Antes:
d.SetAvailableSize(80, 24)
out := d.View()

// Depois:
out := d.View(80, 24)
```

Para casos onde `SetAvailableSize` é chamado e `View()` é chamado em linha diferente, fazer o mesmo. Verificar também `TestDecisionDialog_SmallSizeUsesMinWidth` (linha 408):
```go
// Antes:
d := pocKey4()
d.SetAvailableSize(20, 10)
out := d.View()

// Depois:
d := pocKey4()
out := d.View(20, 10)
```

Para os fixture constructors `pocKey1`…`pocKeyF` que chamam `d.SetAvailableSize(80, 24)` internamente, atualizar para simplesmente não chamar mais (a dimensão será passada no `View()`).

- [ ] **Step 2: Atualizar golden tests — `View()` → `View(w, h)`**

`TestDecisionDialog_Golden` (linha 466) provavelmente chama `d.View()` no loop interno. Substituir por `d.View(width, 24)` conforme cada cenário.

- [ ] **Step 3: Remover `TestDecisionDialog_ViewPanicsWithoutSetAvailableSize` (linha 715)**

Deletar completamente.

- [ ] **Step 4: Rodar testes de decision**

```
go test ./internal/tui/... -run TestDecisionDialog -count=1
```

Expected: PASS.

---

## Task 14: Atualizar `filepicker_test.go`

**Files:**
- Modify: `internal/tui/filepicker_test.go`

- [ ] **Step 1: Substituir `fpk.SetAvailableSize(80, 24)` → usar `fpk.View(80, 24)` onde View é chamado logo após**

Para `TestFilePickerModalView` (linha 48–55):
```go
func TestFilePickerModalView(t *testing.T) {
    fpk := newTestFilePickerModal()
    view := fpk.View(80, 24)
    if view == "" {
        t.Log("View() returned empty string (acceptable for initial render)")
    }
}
```

- [ ] **Step 2: Atualizar `runFPKGolden` helper (linha 940+) — passar dims para `View()`**

Localizar onde `runFPKGolden` chama `fpk.View()` e substituir por `fpk.View(80, 24)` (ou parâmetro de dimensão do helper).

- [ ] **Step 3: Atualizar `newGoldenFPK` (linha 907+) — remover `SetAvailableSize`, passar dims no `View()`**

Qualquer chamada `fpk.SetAvailableSize(...)` dentro de helpers de golden deve ser removida.

- [ ] **Step 4: Remover `TestFilePickerModal_ViewPanicsWithoutSetAvailableSize` (linha 1060)**

Deletar completamente.

- [ ] **Step 5: Rodar testes de filepicker**

```
go test ./internal/tui/... -run TestFilePicker -count=1
```

Expected: PASS.

---

## Task 15: Atualizar `passwordcreate_test.go` e `passwordentry_test.go`

**Files:**
- Modify: `internal/tui/passwordcreate_test.go`
- Modify: `internal/tui/passwordentry_test.go`

- [ ] **Step 1: Atualizar `passwordcreate_test.go` — substituir `SetAvailableSize` + `View()` por `View(w, h)`**

Buscar por `SetAvailableSize` no arquivo. Padrão a substituir:
```go
// Antes:
m.SetAvailableSize(80, 24)
out := m.View()

// Depois:
out := m.View(80, 24)
```

Remover `TestPasswordCreateModal_ViewPanicsWithoutSetAvailableSize` (linha 426).

- [ ] **Step 2: Atualizar `passwordentry_test.go` — substituir `SetAvailableSize` + `View()` por `View(w, h)`**

Mesmo padrão. Remover `TestPasswordEntryModal_ViewPanicsWithoutSetAvailableSize` (linha 319).

- [ ] **Step 3: Rodar testes de password**

```
go test ./internal/tui/... -run TestPassword -count=1
```

Expected: PASS.

---

## Task 16: Atualizar `root_test.go`

**Files:**
- Modify: `internal/tui/root_test.go`

O `stubModal` precisa implementar a nova interface `modalView`. Dois lugares a atualizar:

- [ ] **Step 1: Atualizar `stubModal.View()` para `View(maxWidth, maxHeight int) string`**

Localizar a definição de `stubModal` (próximo de linha 56 ou no início do arquivo) e atualizar:
```go
func (s *stubModal) View(maxWidth, maxHeight int) string { return "" }
```

- [ ] **Step 2: Remover `stubModal.SetAvailableSize`**

Localizar e remover o método `SetAvailableSize` de `stubModal`.

- [ ] **Step 3: Atualizar `TestWindowSizeMsg_NoModalSetSize`**

O teste verifica que modais não recebem `SetSize` via `WindowSizeMsg`. Com a nova arquitetura, `SetAvailableSize` não existe mais — o teste verifica que `modal.Update()` não é chamado no `WindowSizeMsg`. O invariante continua válido.

O teste atual (linha 189) já testa `modal.updateCalls` — não precisa mudar a lógica, só garantir que `stubModal` ainda conta chamadas a `Update()`. O nome do teste pode ser atualizado:

```go
// TestWindowSizeMsg_NoModalUpdate verifies that WindowSizeMsg does not
// propagate to modal.Update() — rootModel handles it directly.
func TestWindowSizeMsg_NoModalUpdate(t *testing.T) {
```

- [ ] **Step 4: Rodar todos os testes de root**

```
go test ./internal/tui/... -run TestRootModel -count=1
go test ./internal/tui/... -run TestWindowSizeMsg -count=1
```

Expected: PASS.

---

## Task 17: Rodar suite completa e verificar golden files

- [ ] **Step 1: Rodar toda a suite**

```
go test ./internal/tui/... -count=1
```

Expected: PASS. Se houver falhas de golden files, prosseguir com Step 2.

- [ ] **Step 2: Regenerar golden files se necessário**

Se houver falhas em golden tests (output mudou por causa da refatoração):

```
go test ./internal/tui/... -count=1 -update
```

Depois verificar manualmente que os golden files regenerados fazem sentido (mesma saída visual — só mudou a assinatura de `View()`, não o output).

- [ ] **Step 3: Rodar suite novamente para confirmar PASS**

```
go test ./internal/tui/... -count=1
```

Expected: PASS com 0 falhas.

- [ ] **Step 4: Verificar que não há referências a `SetSize` ou `SetAvailableSize` no código de produção**

```
go build ./...
```

Confirmar zero erros de compilação.

- [ ] **Step 5: Commit final**

```
git add internal/tui/
git commit -m "refactor(tui): phase 2 — View() receives dimensions as params, remove SetSize/SetAvailableSize"
```

---

## Notas de implementação

### `filepicker.go` — lógica de reset de scroll na primeira renderização

O reset `m.treeScroll = 0; m.adjustTreeScroll()` que estava em `SetAvailableSize` migra para `View()` com guard `if m.viewportHeight == 0`. Isso é correto porque:
- `viewportHeight == 0` significa que `View()` ainda não foi chamado com dimensões reais
- Após o primeiro `View()`, `viewportHeight > 0` e o reset não acontece mais
- O `rootModel` garante que nenhum keypress chega antes do primeiro `View()`

### `help.go` — `contentHeight()` agora retorna `m.viewportHeight`

Os testes de Update que chamam `m.contentHeight()` para calcular `maxScroll` precisam chamar `m.View(w, h)` primeiro. O `viewportHeight` fica válido após o primeiro `View()`.

### `decision.go` e `passwordcreate/entry.go` — `m.width` preservado no estado

Esses modelos ainda usam `m.width` internamente em `boxWidth()` / clamp de `boxW`. A estratégia é: `View(maxWidth, maxHeight)` atribui `m.width = maxWidth` na primeira linha — simples e sem quebrar o código existente de `boxWidth()`.

### Guard `m.width > 0` em `root.go`

Só protege `child.Update()` e `modal.Update()`. `activeFlow.Update()` **não** é protegido — flows podem precisar processar mensagens antes do terminal dimensionar (ex: comandos de linha de CLI passados via `WithInitialPath`).
