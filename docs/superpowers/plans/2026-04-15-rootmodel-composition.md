# RootModel Composition Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implementar a arquitetura de composição do RootModel com 4 regiões fixas, modal overlay via lipgloss compositor, e nomes descritivos para os tipos de view.

**Architecture:** O `RootModel` ganha campos individuais para cada view, `renderWorkArea()` separa a lógica de layout, e `View()` usa `lipgloss.NewCompositor` para sobrepor modais sobre o base layout. Os tipos `secret.TreeView`, `secret.DetailView`, `template.ListView` e `template.DetailView` são renomeados para nomes descritivos antes de qualquer outra mudança.

**Tech Stack:** Go, Bubble Tea v2 (`charm.land/bubbletea/v2`), lipgloss v2 (`charm.land/lipgloss/v2`)

---

## Mapa de Arquivos

| Arquivo | Ação | O que muda |
|---|---|---|
| `internal/tui/design/design.go` | Modificar | Adicionar constantes `HeaderHeight`, `MessageHeight`, `ActionHeight` |
| `internal/tui/secret/vault_tree.go` | Modificar | Renomear `TreeView` → `VaultTreeView`, `NewTreeView` → `NewVaultTreeView`, adicionar parâmetro `vm *vault.Manager` |
| `internal/tui/secret/secret_detail.go` | Modificar | Renomear `DetailView` → `SecretDetailView`, `NewDetailView` → `NewSecretDetailView`, adicionar parâmetro `vm *vault.Manager` |
| `internal/tui/template/template_list.go` | Modificar | Renomear `ListView` → `TemplateListView`, `NewListView` → `NewTemplateListView`, adicionar parâmetro `vm *vault.Manager` |
| `internal/tui/template/template_detail.go` | Modificar | Renomear `DetailView` → `TemplateDetailView`, `NewDetailView` → `NewTemplateDetailView`, adicionar parâmetro `vm *vault.Manager` |
| `internal/tui/settings/settings_view.go` | Modificar | Adicionar parâmetro `vm *vault.Manager` ao construtor |
| `internal/tui/screen/header_view.go` | Modificar | Implementar métodos `ChildView` como stubs, corrigir comentários, adicionar imports |
| `internal/tui/root.go` | Modificar | Adicionar campos de view, `initVaultViews`, `renderWorkArea`, `SetMessage`; reescrever `NewRootModel`, `View`, `Update`; remover guards nil; corrigir bug `Render(width, height)` |

Não há criação de arquivos novos.

---

## Task 1: Constantes de altura em design/design.go

**Files:**
- Modify: `internal/tui/design/design.go`
- Test: `internal/tui/design/design_test.go`

- [ ] **Step 1: Escrever teste que verifica as novas constantes**

Abra `internal/tui/design/design_test.go` e adicione ao final:

```go
func TestLayoutHeightConstants(t *testing.T) {
    if design.HeaderHeight != 2 {
        t.Errorf("HeaderHeight = %d, want 2", design.HeaderHeight)
    }
    if design.MessageHeight != 1 {
        t.Errorf("MessageHeight = %d, want 1", design.MessageHeight)
    }
    if design.ActionHeight != 1 {
        t.Errorf("ActionHeight = %d, want 1", design.ActionHeight)
    }
    // Verificação de sanidade: a soma deve ser menor que MinHeight
    fixed := design.HeaderHeight + design.MessageHeight + design.ActionHeight
    if fixed >= design.MinHeight {
        t.Errorf("soma das regiões fixas %d >= MinHeight %d, não sobraria espaço para work area", fixed, design.MinHeight)
    }
}
```

- [ ] **Step 2: Rodar o teste — deve falhar**

```
go test ./internal/tui/design/...
```

Esperado: `FAIL — design.HeaderHeight undefined`

- [ ] **Step 3: Adicionar as constantes em design/design.go**

Ao final do arquivo `internal/tui/design/design.go`, após a linha `const PanelTreeRatio = 0.35`, adicione:

```go
// HeaderHeight é a altura em linhas da região de cabeçalho da tela principal.
const HeaderHeight = 2

// MessageHeight é a altura em linhas da barra de mensagens de status.
const MessageHeight = 1

// ActionHeight é a altura em linhas da barra de ações do contexto atual.
const ActionHeight = 1
```

- [ ] **Step 4: Rodar o teste — deve passar**

```
go test ./internal/tui/design/...
```

Esperado: `ok  github.com/useful-toys/abditum/internal/tui/design`

- [ ] **Step 5: Verificar build completo**

```
go build ./...
```

Esperado: sem erros.

- [ ] **Step 6: Commit**

```
git add internal/tui/design/design.go internal/tui/design/design_test.go
git commit -m "feat: add HeaderHeight, MessageHeight, ActionHeight constants to design"
```

---

## Task 2: Renomear tipos em secret/ e template/

**Files:**
- Modify: `internal/tui/secret/vault_tree.go`
- Modify: `internal/tui/secret/secret_detail.go`
- Modify: `internal/tui/template/template_list.go`
- Modify: `internal/tui/template/template_detail.go`

Não há testes unitários existentes para estes arquivos (são stubs puros). O "teste" aqui é o build — se o Go compilar sem erros, a renomeação está correta.

- [ ] **Step 1: Substituir vault_tree.go inteiro**

Reescreva `internal/tui/secret/vault_tree.go` com o seguinte conteúdo (todos os receptores e o tipo renomeados):

```go
package secret

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/vault"
)

// VaultTreeView exibe a árvore de pastas e segredos do cofre aberto.
type VaultTreeView struct{}

// NewVaultTreeView cria uma nova instância da árvore do cofre.
func NewVaultTreeView(vm *vault.Manager) *VaultTreeView {
	return &VaultTreeView{}
}

// Render retorna a árvore de segredos preenchendo as dimensões fornecidas com o tema ativo.
func (v *VaultTreeView) Render(height, width int, theme *design.Theme) string {
	content := "Vault Tree"
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Base))
	return style.Render(content)
}

// HandleKey não processa teclas nesta view.
func (v *VaultTreeView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }

// HandleEvent não processa eventos externos nesta view.
func (v *VaultTreeView) HandleEvent(event any) {}

// HandleTeaMsg não processa mensagens do framework nesta view.
func (v *VaultTreeView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }

// Update não altera o estado desta view em resposta a mensagens.
func (v *VaultTreeView) Update(msg tea.Msg) tea.Cmd { return nil }

// Actions retorna nil — VaultTreeView não possui actions próprias nesta sprint.
func (v *VaultTreeView) Actions() []tui.Action { return nil }
```

- [ ] **Step 2: Substituir secret_detail.go inteiro**

Reescreva `internal/tui/secret/secret_detail.go`:

```go
package secret

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/vault"
)

// SecretDetailView exibe os campos e o valor do segredo selecionado na árvore do cofre.
type SecretDetailView struct{}

// NewSecretDetailView cria uma nova instância do painel de detalhes de segredo.
func NewSecretDetailView(vm *vault.Manager) *SecretDetailView {
	return &SecretDetailView{}
}

// Render retorna os detalhes do segredo selecionado pelas dimensões fornecidas com o tema ativo.
func (v *SecretDetailView) Render(height, width int, theme *design.Theme) string {
	content := "Secret Detail"
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Base))
	return style.Render(content)
}

// HandleKey não processa teclas nesta view.
func (v *SecretDetailView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }

// HandleEvent não processa eventos externos nesta view.
func (v *SecretDetailView) HandleEvent(event any) {}

// HandleTeaMsg não processa mensagens do framework nesta view.
func (v *SecretDetailView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }

// Update não altera o estado desta view em resposta a mensagens.
func (v *SecretDetailView) Update(msg tea.Msg) tea.Cmd { return nil }

// Actions retorna nil — SecretDetailView não possui actions próprias nesta sprint.
func (v *SecretDetailView) Actions() []tui.Action { return nil }
```

- [ ] **Step 3: Substituir template_list.go inteiro**

Reescreva `internal/tui/template/template_list.go`:

```go
package template

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/vault"
)

// TemplateListView exibe a lista de templates disponíveis para criação de segredos.
type TemplateListView struct{}

// NewTemplateListView cria uma nova instância da lista de templates.
func NewTemplateListView(vm *vault.Manager) *TemplateListView {
	return &TemplateListView{}
}

// Render retorna a lista de templates preenchendo as dimensões fornecidas com o tema ativo.
func (v *TemplateListView) Render(height, width int, theme *design.Theme) string {
	content := "Template List"
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Base))
	return style.Render(content)
}

// HandleKey não processa teclas nesta view.
func (v *TemplateListView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }

// HandleEvent não processa eventos externos nesta view.
func (v *TemplateListView) HandleEvent(event any) {}

// HandleTeaMsg não processa mensagens do framework nesta view.
func (v *TemplateListView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }

// Update não altera o estado desta view em resposta a mensagens.
func (v *TemplateListView) Update(msg tea.Msg) tea.Cmd { return nil }

// Actions retorna nil — TemplateListView não possui actions próprias nesta sprint.
func (v *TemplateListView) Actions() []tui.Action { return nil }
```

- [ ] **Step 4: Substituir template_detail.go inteiro**

Reescreva `internal/tui/template/template_detail.go`:

```go
package template

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/vault"
)

// TemplateDetailView exibe os detalhes e campos de um template de segredo selecionado.
type TemplateDetailView struct{}

// NewTemplateDetailView cria uma nova instância do painel de detalhes de template.
func NewTemplateDetailView(vm *vault.Manager) *TemplateDetailView {
	return &TemplateDetailView{}
}

// Render retorna os detalhes do template selecionado pelas dimensões fornecidas com o tema ativo.
func (v *TemplateDetailView) Render(height, width int, theme *design.Theme) string {
	content := "Template Detail"
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Base))
	return style.Render(content)
}

// HandleKey não processa teclas nesta view.
func (v *TemplateDetailView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }

// HandleEvent não processa eventos externos nesta view.
func (v *TemplateDetailView) HandleEvent(event any) {}

// HandleTeaMsg não processa mensagens do framework nesta view.
func (v *TemplateDetailView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }

// Update não altera o estado desta view em resposta a mensagens.
func (v *TemplateDetailView) Update(msg tea.Msg) tea.Cmd { return nil }

// Actions retorna nil — TemplateDetailView não possui actions próprias nesta sprint.
func (v *TemplateDetailView) Actions() []tui.Action { return nil }
```

- [ ] **Step 5: Verificar build**

```
go build ./...
```

Esperado: sem erros. Se houver erros de "undefined: TreeView" ou similar, há algum caller externo dos construtores antigos — localize com `grep -r "NewTreeView\|NewDetailView\|NewListView" .` e corrija.

- [ ] **Step 6: Commit**

```
git add internal/tui/secret/vault_tree.go internal/tui/secret/secret_detail.go internal/tui/template/template_list.go internal/tui/template/template_detail.go
git commit -m "refactor: rename secret/template view types to descriptive names"
```

---

## Task 3: Atualizar construtor de SettingsView

**Files:**
- Modify: `internal/tui/settings/settings_view.go`

- [ ] **Step 1: Atualizar settings_view.go**

Substitua as linhas do construtor em `internal/tui/settings/settings_view.go`:

```go
// Antes:
// NewSettingsView cria uma nova instância da tela de configurações.
func NewSettingsView() *SettingsView {
	return &SettingsView{}
}

// Depois:
// NewSettingsView cria uma nova instância da tela de configurações.
// vm é o gerenciador do cofre ativo — pode ser nil durante inicialização.
func NewSettingsView(vm *vault.Manager) *SettingsView {
	return &SettingsView{}
}
```

Também adicione o import de `vault` no bloco de imports:

```go
import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/vault"
)
```

- [ ] **Step 2: Verificar build**

```
go build ./...
```

Esperado: sem erros.

- [ ] **Step 3: Commit**

```
git add internal/tui/settings/settings_view.go
git commit -m "feat: add vault.Manager parameter to NewSettingsView"
```

---

## Task 4: HeaderView implementa ChildView

**Files:**
- Modify: `internal/tui/screen/header_view.go`

- [ ] **Step 1: Reescrever header_view.go**

Substitua o conteúdo completo de `internal/tui/screen/header_view.go`:

```go
package screen

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// HeaderView renderiza o cabeçalho fixo de 2 linhas da aplicação.
// Implementa ChildView — suportará mouse e actions no futuro.
type HeaderView struct{}

// NewHeaderView cria uma nova instância do cabeçalho.
func NewHeaderView() *HeaderView {
	return &HeaderView{}
}

// Render retorna o cabeçalho com as dimensões fornecidas.
// Stub — retorna string vazia até implementação visual completa.
func (v *HeaderView) Render(height, width int, theme *design.Theme) string {
	return ""
}

// HandleKey não processa teclas nesta view.
func (v *HeaderView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }

// HandleEvent não processa eventos externos nesta view.
func (v *HeaderView) HandleEvent(event any) {}

// HandleTeaMsg não processa mensagens do framework nesta view.
func (v *HeaderView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }

// Update não altera o estado desta view em resposta a mensagens.
func (v *HeaderView) Update(msg tea.Msg) tea.Cmd { return nil }

// Actions retorna nil — HeaderView não possui actions próprias nesta sprint.
func (v *HeaderView) Actions() []tui.Action { return nil }
```

- [ ] **Step 2: Verificar build**

```
go build ./...
```

Esperado: sem erros.

- [ ] **Step 3: Commit**

```
git add internal/tui/screen/header_view.go
git commit -m "feat: HeaderView implements ChildView interface"
```

---

## Task 5: Reescrever root.go

Esta é a task central. Reescreve `NewRootModel`, `View`, `Update`, e adiciona `renderWorkArea`, `initVaultViews`, `SetMessage`. Remove guards nil. Corrige bug de `width`/`height` trocados.

**Files:**
- Modify: `internal/tui/root.go`

- [ ] **Step 1: Substituir o conteúdo completo de root.go**

Reescreva `internal/tui/root.go` com o seguinte conteúdo:

```go
package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/screen"
	"github.com/useful-toys/abditum/internal/tui/secret"
	tmpl "github.com/useful-toys/abditum/internal/tui/template"
	"github.com/useful-toys/abditum/internal/tui/settings"
	"github.com/useful-toys/abditum/internal/vault"
)

// WorkArea representa qual área de trabalho está ativa na tela principal.
// É usada por RootModel para decidir qual ChildView exibir.
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

// RootModel é o modelo principal da aplicação Bubble Tea.
// Coordena as 4 regiões fixas da tela, a work area ativa, e a pilha de modais.
type RootModel struct {
	// width e height são as dimensões atuais do terminal, atualizadas em tempo real.
	width  int
	height int
	// theme é o tema visual ativo, aplicado a todos os componentes filhos.
	theme *design.Theme

	// headerView é a região de cabeçalho — sempre presente, implementa ChildView.
	headerView screen.HeaderView
	// messageLineView é a barra de mensagens de status — renderizador stateless.
	messageLineView screen.MessageLineView
	// actionLineView é a barra de ações do contexto — renderizador stateless.
	actionLineView screen.ActionLineView
	// currentMessage é a mensagem de status atual, exibida na messageLineView.
	currentMessage string

	// workArea indica qual área de trabalho está sendo exibida no momento.
	workArea WorkArea
	// activeView aponta para a ChildView com foco na work area atual.
	// Nunca é nil após NewRootModel — inicializado com &welcomeView.
	activeView ChildView

	// welcomeView é exibida quando nenhum cofre está aberto (vaultManager == nil).
	// Valor direto (não ponteiro) — addressável para atribuição a activeView.
	welcomeView screen.WelcomeView

	// As views abaixo dependem de vaultManager e são nil até initVaultViews ser chamado.
	settingsView   *settings.SettingsView
	secretTree     *secret.VaultTreeView
	secretDetail   *secret.SecretDetailView
	templateList   *tmpl.TemplateListView
	templateDetail *tmpl.TemplateDetailView

	// vaultManager é o gerenciador do cofre ativo, ou nil se nenhum cofre estiver carregado.
	vaultManager *vault.Manager

	// modals é a pilha de modais abertos; o topo da pilha é o modal ativo.
	modals []ModalView

	// systemActions são avaliadas em qualquer contexto, inclusive com modal ativo.
	systemActions []Action
	// applicationActions são avaliadas apenas quando nenhum modal está ativo.
	applicationActions []Action
	// actionGroups agrupa actions para exibição no modal de ajuda.
	actionGroups []ActionGroup

	// lastActionAt registra o momento da última interação do usuário.
	lastActionAt time.Time
	// version é a versão da aplicação, normalmente injetada via ldflags no build.
	version string
}

// Manager retorna o vault manager ativo, ou nil se nenhum cofre estiver carregado.
// Implementa a interface AppState.
func (r *RootModel) Manager() *vault.Manager {
	return r.vaultManager
}

// ToggleTheme alterna o tema ativo entre TokyoNight e Cyberpunk.
// Exportada para uso pelo package actions.
func (r *RootModel) ToggleTheme() {
	if r.theme == design.TokyoNight {
		r.theme = design.Cyberpunk
	} else {
		r.theme = design.TokyoNight
	}
}

// SetMessage define a mensagem de status exibida na barra inferior.
// Passe string vazia para limpar a mensagem.
func (r *RootModel) SetMessage(msg string) {
	r.currentMessage = msg
}

// ActiveViewActions retorna todas as actions aplicáveis ao contexto da view ativa.
// Inclui system actions, application actions, e view actions da activeView.
func (r *RootModel) ActiveViewActions() []Action {
	viewActions := r.activeView.Actions()
	allActions := make([]Action, 0, len(r.systemActions)+len(r.applicationActions)+len(viewActions))
	allActions = append(allActions, r.systemActions...)
	allActions = append(allActions, r.applicationActions...)
	allActions = append(allActions, viewActions...)
	return allActions
}

// GetActionGroups retorna a lista de action groups registrados.
// Exportada para uso pelo package actions.
func (r *RootModel) GetActionGroups() []ActionGroup {
	return r.actionGroups
}

// RegisterActionGroup adiciona um grupo de actions ao root.
func (r *RootModel) RegisterActionGroup(group ActionGroup) {
	r.actionGroups = append(r.actionGroups, group)
}

// RegisterSystemActions adiciona actions de sistema ao root.
func (r *RootModel) RegisterSystemActions(actions []Action) {
	r.systemActions = append(r.systemActions, actions...)
}

// RegisterApplicationActions adiciona actions de aplicação ao root.
func (r *RootModel) RegisterApplicationActions(actions []Action) {
	r.applicationActions = append(r.applicationActions, actions...)
}

// evalActions percorre uma lista de actions e executa a primeira que corresponda
// à tecla pressionada e cuja pré-condição esteja satisfeita.
func (r *RootModel) evalActions(msg tea.KeyMsg, actions []Action) (tea.Cmd, bool) {
	for _, action := range actions {
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

// initVaultViews cria as views que dependem do vaultManager.
// Chamado quando vaultManager está disponível — na inicialização ou no ciclo de vida.
func (r *RootModel) initVaultViews() {
	r.settingsView   = settings.NewSettingsView(r.vaultManager)
	r.secretTree     = secret.NewVaultTreeView(r.vaultManager)
	r.secretDetail   = secret.NewSecretDetailView(r.vaultManager)
	r.templateList   = tmpl.NewTemplateListView(r.vaultManager)
	r.templateDetail = tmpl.NewTemplateDetailView(r.vaultManager)
}

// renderWorkArea retorna a string renderizada da área de trabalho ativa.
// Usa as constantes de altura de design para calcular o espaço disponível.
// Nota: a ordem dos argumentos em Render é (height, width) — não inverta.
func (r *RootModel) renderWorkArea() string {
	h := r.height - design.HeaderHeight - design.MessageHeight - design.ActionHeight
	w := r.width

	switch r.workArea {
	case WorkAreaWelcome:
		return r.welcomeView.Render(h, w, r.theme)
	case WorkAreaSettings:
		return r.settingsView.Render(h, w, r.theme)
	case WorkAreaVault:
		treeWidth := int(float64(w) * design.PanelTreeRatio)
		detailWidth := w - treeWidth
		return lipgloss.JoinHorizontal(lipgloss.Top,
			r.secretTree.Render(h, treeWidth, r.theme),
			r.secretDetail.Render(h, detailWidth, r.theme),
		)
	case WorkAreaTemplates:
		listWidth := int(float64(w) * design.PanelTreeRatio)
		detailWidth := w - listWidth
		return lipgloss.JoinHorizontal(lipgloss.Top,
			r.templateList.Render(h, listWidth, r.theme),
			r.templateDetail.Render(h, detailWidth, r.theme),
		)
	default:
		return r.welcomeView.Render(h, w, r.theme)
	}
}

// View gera a representação visual atual da aplicação.
// O base layout (4 regiões) é sempre renderizado.
// Se houver modal ativo, ele é sobreposto ao base via compositor de camadas do lipgloss v2.
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
		r.messageLineView.Render(design.MessageHeight, r.width, r.theme, r.currentMessage),
		r.actionLineView.Render(design.ActionHeight, r.width, r.theme, r.ActiveViewActions()),
	)

	if len(r.modals) > 0 {
		top := r.modals[len(r.modals)-1]
		// Padding de 1 linha acima e abaixo do modal dentro da tela.
		modalH := r.height - 2
		modalContent := top.Render(modalH, r.width, r.theme)
		// Centraliza o conteúdo do modal horizontalmente dentro do espaço disponível.
		centeredModal := lipgloss.Place(r.width, modalH, lipgloss.Center, lipgloss.Center, modalContent)
		// Compõe o modal (z=1) sobre o base layout (z=0) usando o compositor do lipgloss v2.
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

// Update processa mensagens do Bubble Tea e atualiza o estado do modelo.
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
		if cmd, ok := r.evalActions(msg, r.activeView.Actions()); ok {
			return r, cmd
		}

		// 4. Application actions — avaliadas após view actions.
		if cmd, ok := r.evalActions(msg, r.applicationActions); ok {
			return r, cmd
		}

		return r, nil
	}

	if len(r.modals) > 0 {
		top := len(r.modals) - 1
		return r, r.modals[top].Update(msg)
	}

	var cmds []tea.Cmd
	cmds = append(cmds, r.activeView.Update(msg))
	cmds = append(cmds, r.headerView.Update(msg))
	return r, tea.Batch(cmds...)
}

// Init é chamado uma vez ao iniciar a aplicação. Não há comandos iniciais.
func (r *RootModel) Init() tea.Cmd {
	return nil
}

// RootModelOption é uma função de configuração aplicada ao RootModel na criação.
type RootModelOption func(*RootModel)

// WithVersion define a versão da aplicação exibida na interface.
func WithVersion(version string) RootModelOption {
	return func(m *RootModel) {
		m.version = version
	}
}

// NewRootModel cria e inicializa um RootModel com o tema padrão TokyoNight.
// activeView é inicializado com &welcomeView — nunca é nil após esta função.
func NewRootModel(opts ...RootModelOption) *RootModel {
	m := &RootModel{
		theme:    design.TokyoNight,
		workArea: WorkAreaWelcome,
		version:  "dev",
	}
	m.activeView = &m.welcomeView
	for _, opt := range opts {
		opt(m)
	}
	if m.vaultManager != nil {
		m.initVaultViews()
	}
	return m
}
```

- [ ] **Step 2: Verificar build**

```
go build ./...
```

Esperado: sem erros. Se houver erro de `import cycle` ou `undefined`, leia a mensagem com atenção — o package `template` é palavra reservada em Go em alguns contextos; o alias `tmpl` já está no código acima para evitar isso.

- [ ] **Step 3: Rodar todos os testes**

```
go test ./...
```

Esperado: todos os pacotes passam ou reportam `[no test files]`. Nenhum `FAIL`.

- [ ] **Step 4: Commit**

```
git add internal/tui/root.go
git commit -m "feat: implement RootModel 4-region composition with modal overlay"
```

---

## Task 6: Verificação final

- [ ] **Step 1: Build limpo**

```
go build ./...
```

Esperado: sem erros, sem warnings.

- [ ] **Step 2: Todos os testes passam**

```
go test ./...
```

Esperado: todos `ok` ou `[no test files]`, zero `FAIL`.

- [ ] **Step 3: Verificar que activeView nunca é nil na inicialização**

Execute o snippet abaixo como teste rápido no terminal:

```
go vet ./...
```

Esperado: sem erros de vet.

---

## Notas para o Implementador

**Import alias para `template`:** O package `internal/tui/template` conflita com o identificador `template` que é comum em Go. Use o alias `tmpl` conforme mostrado em `root.go` na Task 5.

**Ordem dos argumentos em `Render`:** A assinatura é `Render(height, width int, theme *design.Theme)` — altura primeiro, largura depois. O código antigo em `root.go` tinha esses argumentos invertidos (`Render(r.width, r.height, r.theme)`). O novo código em `renderWorkArea()` usa a ordem correta.

**`welcomeView` como valor direto:** O campo `welcomeView screen.WelcomeView` no `RootModel` é um valor, não um ponteiro. `&m.welcomeView` é atribuído a `activeView` para satisfazer a interface `ChildView` (cujos métodos têm receptores ponteiro). Isso é válido em Go pois campos de struct são addressáveis.

**`vm *vault.Manager` nos construtores:** O parâmetro é recebido mas não usado nos stubs desta sprint. Isso é intencional — a assinatura está correta para o futuro.
