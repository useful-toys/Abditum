# TUI Root Model Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implementar a arquitetura TUI com RootModel, ChildViews stub, e sistema de modais baseada na specapproved.

**Architecture:** Arquitetura onde RootModel é o orquestrador central que coordena Views isoladas. Views não acessam estado de outras views — comunicação via Cmd/msg. Modal via pilha, view guarda referência e consulta via getters.

**Tech Stack:** Go, Bubble Tea (charm.land/bubbletea/v2), Lipgloss (charm.land/lipgloss/v2)

---

## Estrutura de Arquivos a Criar

```
tui/
├── view.go              # interfaces + Theme + WorkArea + Mensagens de modais
├── root.go             # RootModel (tea.Model)
├── welcome/
│   └── welcome_view.go
├── settings/
│   └── settings_view.go
├── secret/
│   ├── vault_tree.go
│   └── secret_detail.go
├── template/
│   ├── template_list.go
│   └── template_detail.go
└── modal/
    ├── modal_base.go
    ├── password_modal.go
    ├── confirm_modal.go
    ├── filepicker_modal.go
    └── help_modal.go
```

---

### Task 1: Tipos base (Theme, WorkArea, Mensagens)

**Files:**
- Create: `tui/view.go`

- [ ] **Step 1: Criar arquivo tui/view.go com interfaces e tipos base**

```go
package tui

import tea "charm.land/bubbletea/v2"

// WorkArea representa o estado da área de trabalho
type WorkArea int

const (
    WorkAreaWelcome WorkArea = iota
    WorkAreaSettings
    WorkAreaVault
    WorkAreaTemplates
)

// Theme define cores, tipografia e símbolos
type Theme struct {
    Name     string
    Surface  SurfaceTokens
    Text     TextTokens
    Accent   AccentTokens
    Border   BorderTokens
    Semantic SemanticTokens
    Special  SpecialTokens
}

type SurfaceTokens struct {
    Base   string
    Raised string
    Input  string
}

type TextTokens struct {
    Primary   string
    Secondary string
    Disabled  string
    Link      string
}

type AccentTokens struct {
    Primary   string
    Secondary string
}

type BorderTokens struct {
    Default string
    Focused string
}

type SemanticTokens struct {
    Success string
    Warning string
    Error   string
    Info    string
    Off     string
}

type SpecialTokens struct {
    Muted     string
    Highlight string
    Match     string
}

// ChildView interface para componentes da tela principal
type ChildView interface {
    Render(height, width int, theme Theme) string
    HandleKey(msg tea.KeyMsg) tea.Cmd
    HandleEvent(event any)
    HandleTeaMsg(msg tea.Msg)
}

// ModalView interface para modais
type ModalView interface {
    Render(maxHeight, maxWidth int, theme Theme) string
    HandleKey(msg tea.KeyMsg) tea.Cmd
}
```

- [ ] **Step 2: Adicionar mensagens de modais no mesmo arquivo**

```go
// OpenModalMsg empilha um novo modal
type OpenModalMsg struct {
    Modal ModalView
}

// CloseModalMsg desempilha o modal do topo
type CloseModalMsg struct{}

// ModalReadyMsg indica que modal tem resultado
type ModalReadyMsg struct{}

// Funções auxiliares para criar comandos
func OpenModal(modal ModalView) tea.Cmd {
    return func() tea.Msg { return OpenModalMsg{Modal: modal} }
}

func CloseModal() tea.Cmd {
    return func() tea.Msg { return CloseModalMsg{} }
}
```

- [ ] **Step 3: Adicionar temas TokyoNight e Cyberpunk**

```go
var TokyoNight = &Theme{
    Name: "Tokyo Night",
    Surface:  SurfaceTokens{"#1a1b26", "#24283b", "#1e1f2e"},
    Text:     TextTokens{"#a9b1d6", "#565f89", "#3b4261", "#7aa2f7"},
    Accent:   AccentTokens{"#7aa2f7", "#bb9af7"},
    Border:   BorderTokens{"#414868", "#7aa2f7"},
    Semantic: SemanticTokens{"#9ece6a", "#e0af68", "#f7768e", "#7dcfff", "#737aa2"},
    Special:  SpecialTokens{"#8690b5", "#283457", "#f7c67a"},
}

var Cyberpunk = &Theme{
    Name: "Cyberpunk",
    Surface:  SurfaceTokens{"#0a0a1a", "#1a1a2e", "#0e0e22"},
    Text:     TextTokens{"#e0e0ff", "#8888aa", "#444466", "#ff2975"},
    Accent:   AccentTokens{"#ff2975", "#00fff5"},
    Border:   BorderTokens{"#3a3a5c", "#ff2975"},
    Semantic: SemanticTokens{"#05ffa1", "#ffe900", "#ff3860", "#00b4d8", "#9999cc"},
    Special:  SpecialTokens{"#666688", "#2a1533", "#ffc107"},
}
```

- [ ] **Step 4: Commit**

```bash
git add tui/view.go
git commit -m "feat(tui): add base types, interfaces, and themes"
```

---

### Task 2: RootModel stub

**Files:**
- Create: `tui/root.go`

- [ ] **Step 1: Criar RootModel com estado mínimo**

```go
package tui

import (
    "time"

    tea "charm.land/bubbletea/v2"
    "charm.land/lipgloss/v2"
)

type RootModel struct {
    width         int
    height        int
    theme         *Theme
    workArea      WorkArea
    activeView    ChildView
    modals        []ModalView
    lastActionAt  time.Time
}
```

- [ ] **Step 2: Implementar View()**

```go
func (r *RootModel) View() string {
    if r.width == 0 || r.height == 0 {
        return "Aguarde..."
    }

    base := r.activeView.Render(r.width, r.height, *r.theme)

    if len(r.modals) == 0 {
        return base
    }

    top := r.modals[len(r.modals)-1]
    modalView := top.Render(r.width, r.height, *r.theme)

    workH := r.height - 4 // 2 header + 1 msg + 1 action
    return lipgloss.Place(r.width, workH, lipgloss.Center, lipgloss.Center, modalView)
}
```

- [ ] **Step 3: Implementar Update() com roteamento**

```go
func (r *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        r.width = msg.Width
        r.height = msg.Height
        return r, nil

    case OpenModalMsg:
        r.modals = append(r.modals, msg.Modal)
        return r, nil

    case CloseModalMsg:
        if len(r.modals) > 0 {
            r.modals = r.modals[:len(r.modals)-1]
        }
        return r, nil

    case ModalReadyMsg:
        if len(r.modals) > 1 {
            parent := r.modals[len(r.modals)-2]
            return r, parent.Update(msg)
        }
        return r, r.activeView.Update(msg)
    }

    if len(r.modals) > 0 {
        top := len(r.modals) - 1
        return r, r.modals[top].Update(msg)
    }

    return r, r.activeView.Update(msg)
}
```

- [ ] **Step 4: Adicionar Init() e método NewRootModel()**

```go
func (r *RootModel) Init() tea.Cmd {
    return nil
}

func NewRootModel(view ChildView) *RootModel {
    return &RootModel{
        theme:      TokyoNight,
        workArea:   WorkAreaWelcome,
        activeView: view,
    }
}
```

- [ ] **Step 5: Commit**

```bash
git add tui/root.go
git commit -m "feat(tui): add RootModel stub with View and Update"
```

---

### Task 3: WelcomeView stub

**Files:**
- Create: `tui/welcome/welcome_view.go`

- [ ] **Step 1: Criar WelcomeView stub**

```go
package welcome

import (
    tea "charm.land/bubbletea/v2"
    "charm.land/lipgloss/v2"
)

type Theme struct {
    Text struct {
        Primary   string
        Secondary string
    }
    Surface struct {
        Base string
    }
}

// WelcomeView exibe a tela de boas-vindas
type WelcomeView struct {
    state int
}

const (
    StateNormal ViewState = iota
    StateAwaitingModal
)

func NewWelcomeView() *WelcomeView {
    return &WelcomeView{state: StateNormal}
}

func (v *WelcomeView) Render(height, width int, theme Theme) string {
    content := "Welcome"
    style := lipgloss.NewStyle().
        Width(width).
        Height(height).
       Foreground(lipgloss.Color(theme.Text.Primary)).
        Background(lipgloss.Color(theme.Surface.Base))
    return style.Render(content)
}

func (v *WelcomeView) HandleKey(msg tea.KeyMsg) tea.Cmd {
    return nil
}

func (v *WelcomeView) HandleEvent(event any) {}

func (v *WelcomeView) HandleTeaMsg(msg tea.Msg) {}
```

- [ ] **Step 2: Commit**

```bash
git add tui/welcome/welcome_view.go
git commit -m "feat(tui): add WelcomeView stub"
```

---

### Task 4: Demais ChildViews stub

**Files:**
- Create: `tui/settings/settings_view.go`
- Create: `tui/secret/vault_tree.go`
- Create: `tui/secret/secret_detail.go`
- Create: `tui/template/template_list.go`
- Create: `tui/template/template_detail.go`

- [ ] **Step 1: Criar SettingsView, Tree, Detail, List stubs**

Cada arquivo segue o mesmo padrão de WelcomeView, apenas renderizando seu nome como conteúdo.

```go
package settings

import tea "charm.land/bubbletea/v2"

type Theme struct {
    Text struct{ Primary string }
    Surface struct{ Base string }
}

type SettingsView struct{}

func NewSettingsView() *SettingsView {
    return &SettingsView{}
}

func (v *SettingsView) Render(height, width int, theme Theme) string {
    return "Settings"
}

func (v *SettingsView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }
func (v *SettingsView) HandleEvent(event any)           {}
func (v *SettingsView) HandleTeaMsg(msg tea.Msg)         {}
```

Repetir para os outros 4 packages com seus nomes respectivos.

- [ ] **Step 2: Commit**

```bash
git add tui/settings/settings_view.go tui/secret/vault_tree.go tui/secret/secret_detail.go tui/template/template_list.go tui/template/template_detail.go
git commit -m "feat(tui): add remaining ChildViews stubs"
```

---

### Task 5: Modal base e implementação básica

**Files:**
- Create: `tui/modal/modal_base.go`
- Create: `tui/modal/confirm_modal.go`

- [ ] **Step 1: Criar modal_base.go com Intent e ModalOption**

```go
package modal

import tea "charm.land/bubbletea/v2"

type Intent int

const (
    IntentConfirm Intent = iota
    IntentCancel
    IntentOther
)

type ModalOption struct {
    Keys   []string
    Label  string
    Intent Intent
    Action func() tea.Cmd
}
```

- [ ] **Step 2: Criar confirm_modal.go stub**

```go
package modal

import (
    "charm.land/lipgloss/v2"
    tea "charm.land/bubbletea/v2"
)

type Theme struct {
    Text struct{ Primary, Secondary string }
    Border struct{ Default, Focused string }
    Surface struct{ Base, Raised string }
}

type ConfirmModal struct {
    title    string
    message  string
    options  []ModalOption
}

func NewConfirmModal(title, message string) *ConfirmModal {
    return &ConfirmModal{
        title:   title,
        message: message,
        options: []ModalOption{
            {Keys: []string{"Enter"}, Label: "Confirmar", Intent: IntentConfirm},
            {Keys: []string{"Esc"}, Label: "Cancelar", Intent: IntentCancel},
        },
    }
}

func (m *ConfirmModal) Render(maxHeight, maxWidth int, theme Theme) string {
    return m.title + "\n" + m.message
}

func (m *ConfirmModal) HandleKey(msg tea.KeyMsg) tea.Cmd {
    for _, opt := range m.options {
        for _, key := range opt.Keys {
            if msg.String() == key {
                return opt.Action()
            }
        }
    }
    return nil
}
```

- [ ] **Step 3: Commit**

```bash
git add tui/modal/modal_base.go tui/modal/confirm_modal.go
git commit -m "feat(tui): add modal base types and confirm modal stub"
```

---

### Task 6: Integração básica (main.go de teste)

**Files:**
- Create: `cmd/abditum/main.go` (se não existir) ou verificar dependências

- [ ] **Step 1: Verificar se o código compila**

```bash
cd C:/git/Abditum-T2 && go build ./...
```

- [ ] **Step 2: Se compilar, commit final**

```bash
git add .
git commit -m "feat(tui): initial TUI architecture with RootModel and stubs"
```

---

## Execution

**Plan complete and saved to `docs/superpowers/plans/2026-04-14-tui-root-model-impl.md`. Two execution options:**

**1. Subagent-Driven (recommended)** - I dispatch a fresh subagent per task, review between tasks, fast iteration

**2. Inline Execution** - Execute tasks in this session using executing-plans, batch execution with checkpoints

**Which approach?**

- If Subagent-Driven chosen: Use superpowers:subagent-driven-development
- If Inline Execution chosen: Use superpowers:executing-plans