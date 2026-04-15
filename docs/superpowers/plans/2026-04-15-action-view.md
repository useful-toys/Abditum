# ActionLineView Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implementar `ActionLineView` — a barra de comandos (última linha) que exibe ações disponíveis no contexto atual, com testes golden file.

**Architecture:** Mover `Action`/`ActionGroup`/`AppState` para `tui/actions/action.go` para eliminar o import cycle; criar `design/design_action.go` com helper de renderização; implementar `ActionLineView.Render` com layout de ancoragem F1 à direita.

**Tech Stack:** Go, charm.land/lipgloss/v2, charm.land/bubbletea/v2, testdata golden helpers

---

## File Map

| Arquivo | Operação | Responsabilidade |
|---|---|---|
| `internal/tui/actions/action.go` | **Criar** | Recebe `Action`, `ActionGroup`, `AppState`, `ChildView` movidos de `tui/action.go` |
| `internal/tui/action.go` | **Deletar** | Conteúdo movido |
| `internal/tui/view.go` | **Editar** | `Actions() []actions.Action` |
| `internal/tui/root.go` | **Editar** | Imports, `ActiveViewActions()` com três fontes, `ActiveViewActionsForBar()` novo, `View()` sem conversão |
| `internal/tui/actions/setup.go` | **Editar** | Trocar `tui.Action`/`tui.ActionGroup` por tipos locais `Action`/`ActionGroup` |
| `internal/tui/modal/help_modal.go` | **Editar** | Trocar `tui.Action`/`tui.ActionGroup` por `actions.Action`/`actions.ActionGroup` |
| `internal/tui/screen/types.go` | **Editar** | Remover `type Action = interface{}` e atualizar `ChildView.Actions()` |
| `internal/tui/screen/header_view.go` | **Editar** | `Actions() []actions.Action` |
| `internal/tui/screen/welcome_view.go` | **Editar** | `Actions() []actions.Action` |
| `internal/tui/secret/vault_tree.go` | **Editar** | `Actions() []actions.Action` |
| `internal/tui/secret/secret_detail.go` | **Editar** | `Actions() []actions.Action` |
| `internal/tui/settings/settings_view.go` | **Editar** | `Actions() []actions.Action` |
| `internal/tui/template/template_list.go` | **Editar** | `Actions() []actions.Action` |
| `internal/tui/template/template_detail.go` | **Editar** | `Actions() []actions.Action` |
| `internal/tui/screen/action_view.go` | **Editar** | Implementação completa de `Render` |
| `internal/tui/design/design_action.go` | **Criar** | `RenderAction` e `ActionSeparator` helpers |
| `internal/tui/design/design_action_test.go` | **Criar** | Testes unitários do helper |
| `internal/tui/screen/action_view_test.go` | **Criar** | Testes unitários + golden file |
| `internal/tui/screen/testdata/golden/actions-*.golden.*` | **Criar** | Golden files gerados com `-update-golden` |

---

## Task 1: Criar `tui/actions/action.go` com os tipos movidos

**Files:**
- Create: `internal/tui/actions/action.go`

- [ ] **Step 1: Criar `internal/tui/actions/action.go`**

```go
package actions

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/vault"
)

// ActionGroup agrupa actions relacionadas para exibição no modal de ajuda.
type ActionGroup struct {
	ID          string // identificador único do grupo
	Label       string // cabeçalho exibido no modal de ajuda
	Description string // texto descritivo do grupo
}

// Action associa teclas a um comportamento da aplicação.
type Action struct {
	Keys          []design.Key                            // Keys[0] é a tecla principal exibida; demais são aliases funcionais
	Label         string                                  // texto curto para a linha de status
	Description   string                                  // texto longo para o modal de ajuda
	GroupID       string                                  // referencia um ActionGroup registrado
	Priority      int                                     // ordenação na linha de status; menor valor = mais destaque
	Visible       bool                                    // false: nunca aparece na linha de status
	AvailableWhen func(app AppState, view ChildView) bool // nil = sempre disponível
	OnExecute     func() tea.Cmd
}

// Matches retorna true se o evento de teclado corresponde a qualquer tecla declarada na action.
func (a Action) Matches(msg tea.KeyMsg) bool {
	for _, k := range a.Keys {
		if k.Matches(msg) {
			return true
		}
	}
	return false
}

// AppState expõe o estado da aplicação necessário para avaliar pré-condições de actions.
// Implementado por RootModel.
type AppState interface {
	Manager() *vault.Manager // nil se nenhum cofre estiver carregado
}

// ChildView é o subconjunto da interface tui.ChildView necessário para avaliar AvailableWhen.
// Definida aqui para evitar import cycle. Nenhuma action atual inspeciona a view —
// a interface está vazia agora, mas nomear o tipo documenta a intenção e permite
// adicionar métodos futuramente sem alterar a assinatura de AvailableWhen.
type ChildView interface{}
```

- [ ] **Step 2: Verificar que o arquivo compila**

```
go build ./internal/tui/actions/...
```
Esperado: sem erros.

- [ ] **Step 3: Commit**

```
git add internal/tui/actions/action.go
git commit -m "feat: move Action, ActionGroup, AppState, ChildView to actions package"
```

---

## Task 2: Atualizar `tui/view.go`, `tui/root.go` e `tui/modal/help_modal.go`

**Files:**
- Modify: `internal/tui/view.go`
- Modify: `internal/tui/root.go`
- Modify: `internal/tui/modal/help_modal.go`
- Delete: `internal/tui/action.go`

- [ ] **Step 1: Substituir todo o conteúdo de `tui/view.go`**

```go
package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/actions"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// ChildView define o contrato para componentes renderizáveis da tela principal.
type ChildView interface {
	// Render retorna a representação em string do componente para exibição.
	// height e width definem as dimensões disponíveis. theme é passado por ponteiro
	// para evitar cópia desnecessária — design.Theme tem 400 bytes.
	Render(height, width int, theme *design.Theme) string

	// HandleKey processa eventos de teclado e retorna um comando ou nil.
	HandleKey(msg tea.KeyMsg) tea.Cmd

	// HandleEvent processa eventos customizados para o componente.
	HandleEvent(event any)

	// HandleTeaMsg processa mensagens do Bubble Tea framework.
	HandleTeaMsg(msg tea.Msg) tea.Cmd

	// Update é chamado para atualizar o estado do componente.
	Update(msg tea.Msg) tea.Cmd

	// Actions retorna as actions disponíveis nesta view.
	// Pode retornar nil se a view não possuir actions próprias.
	Actions() []actions.Action
}
```

- [ ] **Step 2: Editar `tui/root.go` — campos, métodos de registro, ActiveViewActions, ActiveViewActionsForBar, View, Update**

Adicionar import `"github.com/useful-toys/abditum/internal/tui/actions"` no bloco de imports existente.

Alterar os três campos em `RootModel`:
```go
systemActions      []actions.Action
applicationActions []actions.Action
actionGroups       []actions.ActionGroup
```

Alterar os quatro métodos de registro (assinatura e corpo permanecem idênticos, apenas o tipo muda):
```go
func (r *RootModel) RegisterActionGroup(group actions.ActionGroup) {
	r.actionGroups = append(r.actionGroups, group)
}

func (r *RootModel) RegisterSystemActions(acts []actions.Action) {
	r.systemActions = append(r.systemActions, acts...)
}

func (r *RootModel) RegisterApplicationActions(acts []actions.Action) {
	r.applicationActions = append(r.applicationActions, acts...)
}

func (r *RootModel) GetActionGroups() []actions.ActionGroup {
	return r.actionGroups
}
```

Substituir `ActiveViewActions()` — agora combina três fontes (system + application + view):
```go
// ActiveViewActions retorna todas as actions do contexto atual.
// Combina system actions, application actions e actions da view ativa.
// Usado pelo modal de Ajuda para exibir todos os atalhos disponíveis.
func (r *RootModel) ActiveViewActions() []actions.Action {
	viewActions := r.activeView.Actions()
	all := make([]actions.Action, 0,
		len(r.systemActions)+len(r.applicationActions)+len(viewActions))
	all = append(all, r.systemActions...)
	all = append(all, r.applicationActions...)
	all = append(all, viewActions...)
	return all
}
```

Adicionar novo método `ActiveViewActionsForBar()` logo abaixo de `ActiveViewActions()`:
```go
// ActiveViewActionsForBar retorna as actions filtradas e ordenadas para exibição
// na barra de comandos. Filtra por Visible e AvailableWhen; ordena por Priority crescente.
func (r *RootModel) ActiveViewActionsForBar() []actions.Action {
	all := r.ActiveViewActions()

	// Filtrar: Visible == true E (AvailableWhen == nil OU AvailableWhen satisfeita)
	filtered := all[:0]
	for _, a := range all {
		if !a.Visible {
			continue
		}
		if a.AvailableWhen != nil && !a.AvailableWhen(r, r.activeView) {
			continue
		}
		filtered = append(filtered, a)
	}

	// Ordenar por Priority crescente (insertion sort — lista pequena)
	for i := 1; i < len(filtered); i++ {
		for j := i; j > 0 && filtered[j].Priority < filtered[j-1].Priority; j-- {
			filtered[j], filtered[j-1] = filtered[j-1], filtered[j]
		}
	}
	return filtered
}
```

Alterar `evalActions` (apenas tipo do parâmetro):
```go
func (r *RootModel) evalActions(msg tea.KeyMsg, acts []actions.Action) (tea.Cmd, bool) {
	for _, action := range acts {
		if !action.Matches(msg) {
			continue
		}
		if action.AvailableWhen != nil && !action.AvailableWhen(r, r.activeView) {
			continue
		}
		return action.OnExecute(), true
	}
	return nil, false
}
```

Alterar `View()` — remover bloco de conversão `[]interface{}` e chamar `ActiveViewActionsForBar()`:
```go
func (r *RootModel) View() tea.View {
	if r.width == 0 || r.height == 0 {
		return tea.NewView("Aguarde...")
	}
	if r.width < design.MinWidth {
		return tea.NewView("Aumente a largura do terminal!")
	}
	if r.height < design.MinHeight {
		return tea.NewView("Aumente a altura do terminal!")
	}

	base := lipgloss.JoinVertical(lipgloss.Left,
		r.headerView.Render(design.HeaderHeight, r.width, r.theme),
		r.renderWorkArea(),
		r.messageLineView.Render(r.width, r.theme),
		r.actionLineView.Render(r.width, r.theme, r.ActiveViewActionsForBar()),
	)

	if len(r.modals) > 0 {
		top := r.modals[len(r.modals)-1]
		modalH := r.height - 2
		modalContent := top.Render(modalH, r.width, r.theme)
		centeredModal := lipgloss.Place(r.width, modalH, lipgloss.Center, lipgloss.Center, modalContent)
		result := lipgloss.NewCompositor(
			lipgloss.NewLayer(base),
			lipgloss.NewLayer(centeredModal).Y(1).Z(1),
		).Render()
		v := tea.NewView(result)
		v.AltScreen = true
		v.BackgroundColor = lipgloss.Color(r.theme.Surface.Base)
		return v
	}

	return tea.NewView(base)
}
```

Alterar o bloco `case tea.KeyMsg:` em `Update` — remover type assertion `.(Action)` (agora `viewActions` já é `[]actions.Action`):
```go
case tea.KeyMsg:
	// 1. System actions — avaliadas sempre, inclusive com modal ativo.
	if cmd, ok := r.evalActions(msg, r.systemActions); ok {
		return r, cmd
	}

	// 2. Modal ativo recebe a tecla.
	if len(r.modals) > 0 {
		top := len(r.modals) - 1
		return r, r.modals[top].Update(msg)
	}

	// 3. View actions — avaliadas apenas sem modal ativo.
	viewActions := r.activeView.Actions()
	for _, action := range viewActions {
		if !action.Matches(msg) {
			continue
		}
		if action.AvailableWhen != nil && !action.AvailableWhen(r, r.activeView) {
			continue
		}
		return r, action.OnExecute()
	}

	// 4. Application actions — avaliadas após view actions.
	if cmd, ok := r.evalActions(msg, r.applicationActions); ok {
		return r, cmd
	}

	return r, nil
```

- [ ] **Step 3: Editar `internal/tui/modal/help_modal.go` — trocar `tui.Action`/`tui.ActionGroup` por `actions.Action`/`actions.ActionGroup`**

Adicionar import `"github.com/useful-toys/abditum/internal/tui/actions"`. Alterar os campos e a assinatura de `NewHelpModal`:

```go
package modal

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/actions"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// HelpModal exibe todas as actions registradas, agrupadas por ActionGroup.
// Implementa tui.ModalView.
type HelpModal struct {
	actions []actions.Action
	groups  []actions.ActionGroup
}

// NewHelpModal cria um HelpModal com as actions e grupos fornecidos.
func NewHelpModal(acts []actions.Action, groups []actions.ActionGroup) *HelpModal {
	return &HelpModal{
		actions: acts,
		groups:  groups,
	}
}

// Render retorna o modal de ajuda — stub minimalista.
func (m *HelpModal) Render(maxHeight, maxWidth int, theme *design.Theme) string {
	style := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(lipgloss.Color(theme.Border.Default)).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Raised))
	return style.Render("Ajuda")
}

// HandleKey fecha o modal quando Esc é pressionado.
func (m *HelpModal) HandleKey(msg tea.KeyMsg) tea.Cmd {
	if msg.String() == "esc" {
		return tui.CloseModal()
	}
	return nil
}

// Update processa mensagens do Bubble Tea delegando a HandleKey para eventos de teclado.
func (m *HelpModal) Update(msg tea.Msg) tea.Cmd {
	if key, ok := msg.(tea.KeyMsg); ok {
		return m.HandleKey(key)
	}
	return nil
}
```

- [ ] **Step 4: Deletar `tui/action.go`**

```
git rm internal/tui/action.go
```

- [ ] **Step 5: Verificar compilação**

```
go build ./internal/tui/...
```
Esperado: sem erros.

- [ ] **Step 6: Commit**

```
git add -A internal/tui/
git commit -m "refactor: use actions.Action throughout tui package, add ActiveViewActionsForBar"
```

---

## Task 3: Atualizar `actions/setup.go` para usar tipos locais

**Files:**
- Modify: `internal/tui/actions/setup.go`

O package `actions` importava `tui` para usar `tui.Action` e `tui.ActionGroup`. Agora esses tipos estão em `actions` (o próprio package) — o import de `tui` pode ser mantido apenas para `*tui.RootModel`.

- [ ] **Step 1: Substituir todo o conteúdo de `setup.go`**

```go
package actions

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
)

// Setup registra todos os grupos e actions na aplicação (system e application).
// Deve ser chamado após NewRootModel, antes de iniciar o loop do Bubble Tea.
func Setup(r *tui.RootModel) {
	SetupSystem(r)
	SetupApplication(r)
}

// SetupSystem registra os grupos e actions de sistema no root.
// System actions são avaliadas em qualquer contexto, inclusive com modal ativo.
func SetupSystem(r *tui.RootModel) {
	r.RegisterActionGroup(ActionGroup{
		ID:    "system",
		Label: "Sistema",
	})
	r.RegisterSystemActions([]Action{
		{
			Keys:        []design.Key{design.Shortcuts.ThemeToggle},
			Label:       "Tema",
			Description: "Alterna entre os temas Tokyo Night e Cyberpunk.",
			GroupID:     "system",
			Priority:    100,
			Visible:     false,
			OnExecute:   func() tea.Cmd { r.ToggleTheme(); return nil },
		},
	})
}

// SetupApplication registra os grupos e actions de aplicação no root.
// Application actions são avaliadas apenas quando nenhum modal está ativo.
func SetupApplication(r *tui.RootModel) {
	r.RegisterActionGroup(ActionGroup{
		ID:    "app",
		Label: "Aplicação",
	})
	r.RegisterApplicationActions([]Action{
		{
			Keys:        []design.Key{design.Shortcuts.Help},
			Label:       "Ajuda",
			Description: "Abre o diálogo de ajuda com todos os atalhos disponíveis.",
			GroupID:     "app",
			Priority:    10,
			Visible:     true,
			OnExecute: func() tea.Cmd {
				return tui.OpenModal(modal.NewHelpModal(r.ActiveViewActions(), r.GetActionGroups()))
			},
		},
		{
			Keys:        []design.Key{design.Shortcuts.Quit},
			Label:       "Sair",
			Description: "Encerra a aplicação.",
			GroupID:     "app",
			Priority:    20,
			Visible:     true,
			OnExecute:   func() tea.Cmd { return tea.Quit },
		},
	})
}
```

- [ ] **Step 2: Verificar compilação**

```
go build ./internal/tui/actions/...
```
Esperado: sem erros.

- [ ] **Step 3: Commit**

```
git add internal/tui/actions/setup.go
git commit -m "refactor: use local Action/ActionGroup types in actions/setup.go"
```

---

## Task 4: Atualizar `screen/types.go` e todas as views existentes

**Files:**
- Modify: `internal/tui/screen/types.go`
- Modify: `internal/tui/screen/header_view.go`
- Modify: `internal/tui/screen/welcome_view.go`
- Modify: `internal/tui/secret/vault_tree.go`
- Modify: `internal/tui/secret/secret_detail.go`
- Modify: `internal/tui/settings/settings_view.go`
- Modify: `internal/tui/template/template_list.go`
- Modify: `internal/tui/template/template_detail.go`

- [ ] **Step 1: Substituir todo o conteúdo de `screen/types.go`**

```go
package screen

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/actions"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// ChildView matches the interface from tui package to avoid importing it.
// Defines the contract for renderable components on the main screen.
type ChildView interface {
	Render(height, width int, theme *design.Theme) string
	HandleKey(msg tea.KeyMsg) tea.Cmd
	HandleEvent(event any)
	HandleTeaMsg(msg tea.Msg) tea.Cmd
	Update(msg tea.Msg) tea.Cmd
	// Actions retorna as actions disponíveis nesta view.
	// Pode retornar nil se a view não possuir actions próprias.
	Actions() []actions.Action
}
```

- [ ] **Step 2: Atualizar `Actions()` em `screen/header_view.go`**

Adicionar import `"github.com/useful-toys/abditum/internal/tui/actions"` e alterar a assinatura:
```go
func (v *HeaderView) Actions() []actions.Action { return nil }
```

- [ ] **Step 3: Atualizar `Actions()` em `screen/welcome_view.go`**

Adicionar import `"github.com/useful-toys/abditum/internal/tui/actions"` e alterar a assinatura:
```go
func (v *WelcomeView) Actions() []actions.Action { return nil }
```

- [ ] **Step 4: Atualizar `Actions()` em `secret/vault_tree.go`**

Adicionar import `"github.com/useful-toys/abditum/internal/tui/actions"` e alterar a assinatura:
```go
func (v *VaultTreeView) Actions() []actions.Action { return nil }
```

- [ ] **Step 5: Atualizar `Actions()` em `secret/secret_detail.go`**

Adicionar import `"github.com/useful-toys/abditum/internal/tui/actions"` e alterar a assinatura:
```go
func (v *SecretDetailView) Actions() []actions.Action { return nil }
```

- [ ] **Step 6: Atualizar `Actions()` em `settings/settings_view.go`**

Adicionar import `"github.com/useful-toys/abditum/internal/tui/actions"` e alterar a assinatura:
```go
func (v *SettingsView) Actions() []actions.Action { return nil }
```

- [ ] **Step 7: Atualizar `Actions()` em `template/template_list.go`**

Adicionar import `"github.com/useful-toys/abditum/internal/tui/actions"` e alterar a assinatura:
```go
func (v *TemplateListView) Actions() []actions.Action { return nil }
```

- [ ] **Step 8: Atualizar `Actions()` em `template/template_detail.go`**

Adicionar import `"github.com/useful-toys/abditum/internal/tui/actions"` e alterar a assinatura:
```go
func (v *TemplateDetailView) Actions() []actions.Action { return nil }
```

- [ ] **Step 9: Verificar compilação e testes**

```
go build ./internal/tui/...
go test ./internal/tui/...
```
Esperado: sem erros de compilação, todos os testes passam.

- [ ] **Step 10: Commit**

```
git add -A internal/tui/
git commit -m "refactor: update all ChildView.Actions() to []actions.Action"
```

---

## Task 5: Criar `design/design_action.go` com helper de renderização

**Files:**
- Create: `internal/tui/design/design_action.go`
- Create: `internal/tui/design/design_action_test.go`

- [ ] **Step 1: Criar `internal/tui/design/design_action_test.go`**

```go
package design

import (
	"regexp"
	"testing"

	"charm.land/lipgloss/v2"
)

var ansiEscapeRe = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func TestRenderAction_WidthMatchesLipglossWidth(t *testing.T) {
	theme := TokyoNight
	ra := RenderAction("⌃S", "Salvar", theme)

	measured := lipgloss.Width(ra.Text)
	if measured != ra.Width {
		t.Errorf("RenderAction: Width = %d, lipgloss.Width(Text) = %d", ra.Width, measured)
	}
}

func TestRenderAction_ContainsKeyAndLabel(t *testing.T) {
	theme := TokyoNight
	ra := RenderAction("F1", "Ajuda", theme)

	clean := ansiEscapeRe.ReplaceAllString(ra.Text, "")
	if len(clean) == 0 {
		t.Fatal("RenderAction: texto limpo está vazio")
	}
	if !containsSubstring(clean, "F1") {
		t.Errorf("RenderAction: texto limpo %q não contém tecla 'F1'", clean)
	}
	if !containsSubstring(clean, "Ajuda") {
		t.Errorf("RenderAction: texto limpo %q não contém rótulo 'Ajuda'", clean)
	}
}

func TestActionSeparator_Width(t *testing.T) {
	theme := TokyoNight
	sep := ActionSeparator(theme)

	if sep.Width != 3 {
		t.Errorf("ActionSeparator: Width = %d, want 3", sep.Width)
	}
	if lipgloss.Width(sep.Text) != 3 {
		t.Errorf("ActionSeparator: lipgloss.Width = %d, want 3", lipgloss.Width(sep.Text))
	}
}

func TestRenderAction_WidthPositive(t *testing.T) {
	theme := TokyoNight
	ra := RenderAction("⌃Q", "Sair", theme)
	if ra.Width <= 0 {
		t.Errorf("RenderAction: Width = %d, deve ser > 0", ra.Width)
	}
}

// containsSubstring verifica se s contém sub.
func containsSubstring(s, sub string) bool {
	if len(sub) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
```

- [ ] **Step 2: Rodar o teste — deve falhar**

```
go test ./internal/tui/design/ -run TestRenderAction -v
```
Esperado: `FAIL` — `RenderAction` não existe ainda.

- [ ] **Step 3: Criar `internal/tui/design/design_action.go`**

```go
package design

import "charm.land/lipgloss/v2"

// RenderedAction encapsula o texto ANSI estilizado de uma Action e sua largura em colunas.
// Usar Width (em vez de len) porque sequências ANSI não contribuem para a largura visual.
type RenderedAction struct {
	Text  string // texto estilizado com sequências ANSI
	Width int    // largura visual em colunas (lipgloss.Width — nunca len)
}

// RenderAction renderiza uma action: tecla (Accent.Primary + bold) + espaço + rótulo (Text.Primary).
// key é o rótulo da tecla (ex: "⌃S", "F1"); label é o texto descritivo (ex: "Salvar", "Ajuda").
// Retorna RenderedAction com texto ANSI e largura em colunas.
func RenderAction(key, label string, theme *Theme) RenderedAction {
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Accent.Primary)).
		Bold(true)
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Text.Primary))

	rendered := keyStyle.Render(key) + " " + labelStyle.Render(label)
	return RenderedAction{
		Text:  rendered,
		Width: lipgloss.Width(rendered),
	}
}

// ActionSeparator retorna o separador " · " (espaço + SymHeaderSep + espaço) estilizado
// com Text.Secondary. Sempre tem Width == 3.
func ActionSeparator(theme *Theme) RenderedAction {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Secondary))
	rendered := style.Render(" " + SymHeaderSep + " ")
	return RenderedAction{
		Text:  rendered,
		Width: lipgloss.Width(rendered),
	}
}
```

- [ ] **Step 4: Rodar os testes**

```
go test ./internal/tui/design/ -run "TestRenderAction|TestActionSeparator" -v
```
Esperado: todos `PASS`.

- [ ] **Step 5: Commit**

```
git add internal/tui/design/design_action.go internal/tui/design/design_action_test.go
git commit -m "feat: add RenderAction and ActionSeparator helpers in design package"
```

---

## Task 6: Implementar `ActionLineView.Render` com testes

**Files:**
- Modify: `internal/tui/screen/action_view.go`
- Create: `internal/tui/screen/action_view_test.go`

- [ ] **Step 1: Criar `internal/tui/screen/action_view_test.go`**

```go
package screen

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/actions"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/testdata"
)

// buildHelpAction cria a action F1 Ajuda (âncora da barra).
func buildHelpAction() actions.Action {
	return actions.Action{
		Keys:     []design.Key{design.Shortcuts.Help},
		Label:    "Ajuda",
		Priority: 10,
		Visible:  true,
	}
}

// buildAction cria uma action simples para testes, com rótulo de tecla arbitrário.
func buildAction(keyLabel, label string, priority int) actions.Action {
	return actions.Action{
		Keys:     []design.Key{{Label: keyLabel}},
		Label:    label,
		Priority: priority,
		Visible:  true,
	}
}

// stripANSI remove sequências de escape ANSI de uma string, retornando apenas texto visível.
func stripANSI(s string) string {
	var result strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			j := i + 2
			for j < len(s) && (s[j] < 'A' || s[j] > 'Z') && (s[j] < 'a' || s[j] > 'z') {
				j++
			}
			if j < len(s) {
				j++
			}
			i = j
		} else {
			result.WriteByte(s[i])
			i++
		}
	}
	return result.String()
}

func TestActionLineView_ZeroValue(t *testing.T) {
	var v ActionLineView
	output := v.Render(80, design.TokyoNight, nil)
	if lipgloss.Width(output) == 0 {
		t.Error("Render de zero value retornou string vazia")
	}
}

func TestActionLineView_Render_ExactWidth(t *testing.T) {
	var v ActionLineView
	acts := []actions.Action{buildHelpAction()}
	output := v.Render(80, design.TokyoNight, acts)
	w := lipgloss.Width(output)
	if w != 80 {
		t.Errorf("Render: largura = %d, want 80", w)
	}
}

func TestActionLineView_Render_NoNewline(t *testing.T) {
	var v ActionLineView
	acts := []actions.Action{buildHelpAction()}
	output := v.Render(80, design.TokyoNight, acts)
	if strings.Contains(output, "\n") {
		t.Error("Render não deve conter newline — barra é linha única")
	}
}

func TestActionLineView_Render_EmptyList_CorrectWidth(t *testing.T) {
	var v ActionLineView
	output := v.Render(80, design.TokyoNight, nil)
	w := lipgloss.Width(output)
	if w != 80 {
		t.Errorf("Render lista vazia: largura = %d, want 80", w)
	}
}

func TestActionLineView_Render_F1IsAnchor(t *testing.T) {
	var v ActionLineView
	acts := []actions.Action{buildHelpAction()}
	output := v.Render(80, design.TokyoNight, acts)
	clean := stripANSI(output)
	if !strings.Contains(clean, "F1") {
		t.Errorf("Render: âncora F1 ausente; output limpo: %q", clean)
	}
	if !strings.Contains(clean, "Ajuda") {
		t.Errorf("Render: label 'Ajuda' ausente; output limpo: %q", clean)
	}
}

func TestActionLineView_Render_TruncatesLowPriority(t *testing.T) {
	var v ActionLineView
	acts := []actions.Action{
		buildHelpAction(),
		buildAction("⌃S", "Salvar Cofre Atual", 1),
		buildAction("⌃O", "Abrir Cofre Existente", 2),
		buildAction("⌃N", "Novo Cofre Vazio", 3),
		buildAction("⌃E", "Exportar Todos os Dados do Cofre", 4),
		buildAction("⌃I", "Importar Dados para o Cofre", 5),
	}
	output := v.Render(80, design.TokyoNight, acts)
	w := lipgloss.Width(output)
	if w != 80 {
		t.Errorf("Render overflow: largura = %d, want 80", w)
	}
}

func TestActionLineView_Render_MultipleWidths(t *testing.T) {
	var v ActionLineView
	acts := []actions.Action{
		buildHelpAction(),
		buildAction("⌃S", "Salvar", 1),
	}
	for _, width := range []int{80, 100, 120} {
		output := v.Render(width, design.TokyoNight, acts)
		w := lipgloss.Width(output)
		if w != width {
			t.Errorf("Render width=%d: largura = %d, want %d", width, w, width)
		}
	}
}

// --- Golden file tests ---

var actionGoldenSizes = []string{"80x1"}

func actionRenderFn(acts []actions.Action) testdata.RenderFn {
	return func(w, _ int, theme *design.Theme) string {
		var v ActionLineView
		return v.Render(w, theme, acts)
	}
}

func TestActionLineView_Golden_Empty(t *testing.T) {
	testdata.TestRenderManaged(t, "actions", "empty", actionGoldenSizes,
		actionRenderFn(nil),
	)
}

func TestActionLineView_Golden_SingleAction(t *testing.T) {
	testdata.TestRenderManaged(t, "actions", "single-action", actionGoldenSizes,
		actionRenderFn([]actions.Action{buildHelpAction()}),
	)
}

func TestActionLineView_Golden_MultipleActions(t *testing.T) {
	testdata.TestRenderManaged(t, "actions", "multiple-actions", actionGoldenSizes,
		actionRenderFn([]actions.Action{
			buildHelpAction(),
			buildAction("⌃S", "Salvar", 1),
			buildAction("⌃Q", "Sair", 2),
		}),
	)
}

func TestActionLineView_Golden_Overflow(t *testing.T) {
	testdata.TestRenderManaged(t, "actions", "overflow", actionGoldenSizes,
		actionRenderFn([]actions.Action{
			buildHelpAction(),
			buildAction("⌃S", "Salvar", 1),
			buildAction("⌃O", "Abrir Cofre Existente", 2),
			buildAction("⌃N", "Novo Cofre", 3),
			buildAction("⌃E", "Exportar Todos os Dados do Cofre Para Arquivo Externo", 4),
		}),
	)
}

func TestActionLineView_Golden_NoF1(t *testing.T) {
	testdata.TestRenderManaged(t, "actions", "no-f1", actionGoldenSizes,
		actionRenderFn([]actions.Action{
			buildAction("⌃S", "Salvar", 1),
			buildAction("⌃Q", "Sair", 2),
		}),
	)
}
```

- [ ] **Step 2: Rodar os testes — devem falhar (stub retorna "")**

```
go test ./internal/tui/screen/ -run TestActionLineView -v
```
Esperado: `FAIL` — `Render` ainda é stub que retorna `""`.

- [ ] **Step 3: Substituir o conteúdo de `action_view.go`**

```go
package screen

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/actions"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// ActionLineView renderiza a linha de ações disponíveis no contexto atual.
// Não implementa ChildView — é um renderizador stateless chamado diretamente pelo root.
// O zero value é válido e produz uma linha com apenas espaços na largura correta.
type ActionLineView struct{}

// NewActionLineView cria uma nova instância da linha de ações.
func NewActionLineView() *ActionLineView {
	return &ActionLineView{}
}

// Render retorna a linha de ações com exatamente `width` colunas.
//
// Layout:
//
//	[2 espaços][ação₁][ · ][ação₂][ · ]…[padding][F1 Ajuda]
//
// A âncora F1 é identificada por design.Shortcuts.Help e fixada à direita.
// Ações que não cabem no espaço disponível são descartadas (as de maior Priority, mais à direita).
// `acts` deve estar pré-ordenada por Priority crescente (menor = mais à esquerda).
func (v *ActionLineView) Render(width int, theme *design.Theme, acts []actions.Action) string {
	const (
		prefixCols = 2 // 2 espaços à esquerda
		anchorCols = 8 // reservado para "F1 Ajuda" ou espaços quando âncora ausente
		minPadding = 1 // pelo menos 1 espaço entre ações normais e âncora
	)

	// Separar âncora (F1) das demais ações.
	var anchor *actions.Action
	var normal []actions.Action
	for i := range acts {
		a := acts[i]
		if len(a.Keys) > 0 &&
			a.Keys[0].Code == design.Shortcuts.Help.Code &&
			a.Keys[0].Mod == design.Shortcuts.Help.Mod {
			anchor = &a
		} else {
			normal = append(normal, a)
		}
	}

	// Espaço disponível para ações normais: total menos prefixo, padding mínimo e âncora.
	availableCols := width - prefixCols - minPadding - anchorCols

	// Renderizar ações normais que cabem no espaço disponível.
	sep := design.ActionSeparator(theme)
	var renderedNormal []design.RenderedAction
	usedCols := 0
	for _, a := range normal {
		if len(a.Keys) == 0 {
			continue
		}
		ra := design.RenderAction(a.Keys[0].Label, a.Label, theme)
		needed := ra.Width
		if len(renderedNormal) > 0 {
			needed += sep.Width // separador antes de cada ação (exceto a primeira)
		}
		if usedCols+needed > availableCols {
			break // ações restantes não cabem — descartar
		}
		renderedNormal = append(renderedNormal, ra)
		usedCols += needed
	}

	// Montar bloco de ações normais com separadores.
	var normalBuilder strings.Builder
	for i, ra := range renderedNormal {
		if i > 0 {
			normalBuilder.WriteString(sep.Text)
		}
		normalBuilder.WriteString(ra.Text)
	}
	normalText := normalBuilder.String()

	// Calcular padding entre ações normais e âncora.
	paddingCols := width - prefixCols - usedCols - anchorCols
	if paddingCols < minPadding {
		paddingCols = minPadding
	}

	// Renderizar âncora ou preencher com espaços quando ausente.
	var anchorText string
	if anchor != nil && len(anchor.Keys) > 0 {
		ra := design.RenderAction(anchor.Keys[0].Label, anchor.Label, theme)
		anchorText = ra.Text
	} else {
		anchorText = strings.Repeat(" ", anchorCols)
	}

	return strings.Repeat(" ", prefixCols) +
		normalText +
		strings.Repeat(" ", paddingCols) +
		anchorText
}
```

- [ ] **Step 4: Rodar os testes unitários**

```
go test ./internal/tui/screen/ -run TestActionLineView -v
```
Esperado: todos `PASS`.

Se `TestActionLineView_Render_ExactWidth` falhar por largura incorreta, a causa provável é que `lipgloss.Width` dos fragmentos estilizados inclui terminadores de reset. Corrigir medindo `lipgloss.Width(normalText)` em vez de `usedCols` ao calcular o padding:
```go
paddingCols := width - prefixCols - lipgloss.Width(normalText) - anchorCols
```

- [ ] **Step 5: Gerar golden files**

```
go test ./internal/tui/screen/ -run TestActionLineView_Golden -update-golden -v
```
Esperado: golden files criados em `screen/testdata/golden/actions-*.golden.txt` e `actions-*.golden.json`.

- [ ] **Step 6: Confirmar que os golden tests passam**

```
go test ./internal/tui/screen/ -run TestActionLineView_Golden -v
```
Esperado: todos `PASS`.

- [ ] **Step 7: Commit**

```
git add internal/tui/screen/action_view.go internal/tui/screen/action_view_test.go internal/tui/screen/testdata/
git commit -m "feat: implement ActionLineView.Render with golden file tests"
```

---

## Task 7: Executar todos os testes e corrigir regressões

- [ ] **Step 1: Rodar todos os testes**

```
go test ./...
```
Esperado: todos passam. Se houver falhas, corrigi-las antes de prosseguir.

- [ ] **Step 2: Verificar que o build completo funciona**

```
go build ./...
```
Esperado: sem erros.

- [ ] **Step 3: Commit final (apenas se houver correções)**

```
git add -A
git commit -m "fix: resolve regressions after ActionLineView implementation"
```

---

## Self-Review

### Spec coverage

| Seção do spec | Tarefa que implementa |
|---|---|
| Prerequisito: reorganização de tipos | Tasks 1–4 |
| `ChildView interface{}` nomeada em `actions` | Task 1 Step 1 |
| `ActionLineView` stateless | Task 6 Step 3 |
| Assinatura `Render(width, theme, actions)` sem `height` | Task 6 Step 3 |
| Identificação da âncora F1 por `Keys[0].Code` e `.Mod` | Task 6 Step 3 |
| Algoritmo de layout (prefixo 2, âncora 8, padding mínimo 1) | Task 6 Step 3 |
| `lipgloss.Width` (não `len`) | Task 6 Step 3 |
| Identidade visual: Accent.Primary + bold / Text.Primary / Text.Secondary | Task 5 Step 3 |
| Separador `" " + SymHeaderSep + " "` (não literal `" · "`) | Task 5 Step 3 |
| Sem background na barra | Task 6 Step 3 (nenhum `.Background()`) |
| `design_action.go` com `RenderAction` e `ActionSeparator` | Task 5 |
| `ActiveViewActions()` inclui três fontes | Task 2 Step 2 |
| `ActiveViewActionsForBar()` chama `ActiveViewActions()` internamente | Task 2 Step 2 |
| `View()` chama `ActiveViewActionsForBar()` sem conversão `[]interface{}` | Task 2 Step 2 |
| `modal/help_modal.go` atualizado para `actions.Action` | Task 2 Step 3 |
| Testes unitários (ZeroValue, ExactWidth, NoNewline, F1Anchor, Truncate, MultipleWidths) | Task 6 Step 1 |
| Testes golden (empty, single-action, multiple-actions, overflow, no-f1) | Task 6 Steps 5–6 |

### Placeholder scan

Nenhum "TBD" ou "TODO" encontrado. Todos os steps contêm código completo.

### Type consistency

- `actions.Action` — definido em Task 1, usado em Tasks 2–4, 6.
- `actions.ActionGroup` — definido em Task 1, usado em Tasks 2–3.
- `actions.ChildView` — definido em Task 1 (interface vazia nomeada).
- `actions.AppState` — definido em Task 1, implementado por `RootModel`.
- `design.RenderedAction` — definido em Task 5, usado em Task 6.
- `design.RenderAction(key, label string, theme *Theme) RenderedAction` — definido em Task 5 Step 3, chamado em Task 6 Step 3.
- `design.ActionSeparator(theme *Theme) RenderedAction` — definido em Task 5 Step 3, chamado em Task 6 Step 3.
- `ActionLineView.Render(width int, theme *design.Theme, acts []actions.Action) string` — stub atual tem `height` extra e `[]interface{}`; Task 2 corrige a chamada em `View()` e Task 6 corrige a implementação.
