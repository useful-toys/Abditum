# Design: RootModel com Composição Direta de Views

**Data:** 2026-04-14  
**Status:** Aprovado

---

## Contexto e Motivação

O `RootModel` atual tem três problemas críticos:

1. **`activeView` pode ser nil** — `NewRootModel()` inicializa `activeView = nil`, causando panic em `View()` (linha 135) e na chamada a `r.activeView.Actions()` dentro de `ActiveViewActions()` (linha 76) e `Update()` (linha 191).
2. **`workArea` é ignorado** — o enum `WorkArea` existe mas nunca é usado para decidir o que renderizar.
3. **Arquitetura incompleta** — a intenção de ter 4 regiões fixas (Header, WorkArea, Message, Action) não está implementada.

Adicionalmente, o `root.go` atual tem um bug: `View()` chama `r.activeView.Render(r.width, r.height, r.theme)` com os argumentos `width` e `height` trocados em relação à assinatura `Render(height, width int, ...)`. Este bug também é corrigido neste design.

---

## Arquitetura

### 4 Regiões Fixas

A tela é composta de 4 regiões sempre presentes, empilhadas verticalmente com `lipgloss.JoinVertical`:

```
┌─────────────────────────────┐
│ Header          (2 linhas)  │
├─────────────────────────────┤
│ Work Area       (dinâmico)  │
├─────────────────────────────┤
│ Message         (1 linha)   │
├─────────────────────────────┤
│ Action          (1 linha)   │
└─────────────────────────────┘
```

As alturas das regiões fixas e a proporção do painel esquerdo são constantes definidas em `design/design.go`:

```go
// Em design/design.go — a adicionar
const (
    HeaderHeight  = 2  // linhas ocupadas pelo header
    MessageHeight = 1  // linhas ocupadas pela barra de mensagem
    ActionHeight  = 1  // linhas ocupadas pela barra de actions
)
```

`PanelTreeRatio`, `MinWidth` e `MinHeight` já existem em `design/design.go`.

Essas constantes são usadas tanto no cálculo da work area quanto nas chamadas `Render()` de cada região.

---

## Interfaces — Arquivo view.go

A interface `ChildView` já existe e já inclui `HandleEvent`. Ela não muda de contrato — apenas `HeaderView` passa a implementá-la (hoje só tem `Render`):

```go
// ChildView — contrato atual, não muda
type ChildView interface {
    Render(height, width int, theme *design.Theme) string
    HandleKey(msg tea.KeyMsg) tea.Cmd
    HandleEvent(event any)
    HandleTeaMsg(msg tea.Msg) tea.Cmd
    Update(msg tea.Msg) tea.Cmd
    Actions() []Action
}
```

`MessageLineView` e `ActionLineView` são **casos especiais** — não implementam `ChildView`. São renderizadores sem estado próprio cujos dados de renderização são fornecidos pelo RootModel como parâmetros de `Render`.

---

## Renomeação de Tipos

Os tipos nos packages `secret` e `template` usam nomes genéricos que conflitam entre si e não expressam seu propósito. Devem ser renomeados **nos seus arquivos fonte** antes de qualquer outra mudança:

| Arquivo | Tipo atual | Tipo novo | Construtor atual | Construtor novo |
|---|---|---|---|---|
| `internal/tui/secret/vault_tree.go` | `TreeView` | `VaultTreeView` | `NewTreeView()` | `NewVaultTreeView()` |
| `internal/tui/secret/secret_detail.go` | `DetailView` | `SecretDetailView` | `NewDetailView()` | `NewSecretDetailView()` |
| `internal/tui/template/template_list.go` | `ListView` | `TemplateListView` | `NewListView()` | `NewTemplateListView()` |
| `internal/tui/template/template_detail.go` | `DetailView` | `TemplateDetailView` | `NewDetailView()` | `NewTemplateDetailView()` |

A renomeação inclui: o nome do tipo (`type X struct`), o nome do construtor (`func NewX()`), o tipo de retorno do construtor, e todos os receptores dos métodos (`func (v *X)`). Os comentários do tipo e do construtor também devem usar o novo nome.

---

## Nomes dos Tipos Após Renomeação

Os tipos que o RootModel referenciará após a renomeação acima:

| Package | Tipo | Construtor |
|---|---|---|
| `screen` | `HeaderView` | `NewHeaderView()` |
| `screen` | `MessageLineView` | `NewMessageLineView()` |
| `screen` | `ActionLineView` | `NewActionLineView()` |
| `screen` | `WelcomeView` | `NewWelcomeView()` |
| `settings` | `SettingsView` | `NewSettingsView()` |
| `secret` | `VaultTreeView` | `NewVaultTreeView()` |
| `secret` | `SecretDetailView` | `NewSecretDetailView()` |
| `template` | `TemplateListView` | `NewTemplateListView()` |
| `template` | `TemplateDetailView` | `NewTemplateDetailView()` |

---

## Estrutura do RootModel

```go
type RootModel struct {
    // Dimensões do terminal
    width  int
    height int

    // Tema visual ativo
    theme *design.Theme

    // Região fixa com estado (ChildView — suportará mouse/actions no futuro)
    headerView screen.HeaderView

    // Regiões fixas sem estado — casos especiais, não implementam ChildView
    messageLineView screen.MessageLineView
    actionLineView  screen.ActionLineView

    // Mensagem de status atual — estado simples gerenciado pelo RootModel
    currentMessage string

    // Controle da work area
    workArea   WorkArea
    activeView ChildView  // qual ChildView tem o foco dentro da work area

    // View presente quando vaultManager == nil (valor direto, nunca nil)
    welcomeView screen.WelcomeView

    // Views presentes quando vaultManager != nil (ponteiros, podem ser nil)
    settingsView   *settings.SettingsView
    secretTree     *secret.VaultTreeView
    secretDetail   *secret.SecretDetailView
    templateList   *template.TemplateListView
    templateDetail *template.TemplateDetailView

    // vaultManager — nil até ser criado no ciclo de vida da aplicação
    vaultManager *vault.Manager

    // Pilha de modais
    modals []ModalView

    // Sistema de actions
    systemActions      []Action
    applicationActions []Action
    actionGroups       []ActionGroup

    // Metadados
    lastActionAt time.Time
    version      string
}
```

Os campos `width`, `height`, `theme`, `workArea`, `activeView`, `modals`, `lastActionAt`, `version`, `vaultManager`, `systemActions`, `applicationActions` e `actionGroups` já existem no `RootModel` atual. Os demais (`headerView`, `messageLineView`, `actionLineView`, `currentMessage`, `welcomeView`, `settingsView`, `secretTree`, `secretDetail`, `templateList`, `templateDetail`) são novos.

### Imports novos necessários em root.go

```go
"github.com/useful-toys/abditum/internal/tui/screen"
"github.com/useful-toys/abditum/internal/tui/settings"
"github.com/useful-toys/abditum/internal/tui/secret"
"github.com/useful-toys/abditum/internal/tui/template"
```

---

## Inicialização

`NewRootModel()` já existe com a assinatura `func NewRootModel(opts ...RootModelOption) *RootModel`. O corpo deve ser substituído para:

- Definir `workArea = WorkAreaWelcome`
- Definir `activeView = &m.welcomeView` — garante que `activeView` nunca é nil após a inicialização
- `welcomeView`, `headerView`, `messageLineView`, `actionLineView` são valores diretos; zero value é suficiente para stubs
- `currentMessage = ""` (zero value, não precisa ser explícito)
- Views dependentes de `vaultManager` ficam como ponteiros nil

```go
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

`RootModelOption` já existe como `type RootModelOption func(*RootModel)`.

---

## Views a Adaptar

### HeaderView — adaptar para ChildView

Existe em `internal/tui/screen/header_view.go` mas hoje só tem `Render`. O comentário do arquivo diz erroneamente "não implementa ChildView". Precisa implementar os métodos restantes como stubs e adicionar os imports necessários (`bubbletea/v2` e `tui`):

```go
func (v *HeaderView) HandleKey(msg tea.KeyMsg) tea.Cmd   { return nil }
func (v *HeaderView) HandleEvent(event any)              {}
func (v *HeaderView) HandleTeaMsg(msg tea.Msg) tea.Cmd   { return nil }
func (v *HeaderView) Update(msg tea.Msg) tea.Cmd         { return nil }
func (v *HeaderView) Actions() []tui.Action              { return nil }
```

O comentário do tipo e do construtor deve ser corrigido para refletir que implementa `ChildView`.

### MessageLineView e ActionLineView — sem mudança de interface

Já existem com a assinatura correta. Não implementam `ChildView` — isso é intencional:

```go
func (v *MessageLineView) Render(height, width int, theme *design.Theme, message string) string
func (v *ActionLineView) Render(height, width int, theme *design.Theme, actions []tui.Action) string
```

O comentário de cada arquivo diz "não implementa ChildView — é um renderizador stateless". Isso permanece correto.

### Views da work area — sem mudança de interface

`WelcomeView`, `SettingsView`, `VaultTreeView`, `SecretDetailView`, `TemplateListView`, `TemplateDetailView` já implementam `ChildView` completo (incluindo `HandleEvent`). Nenhuma mudança de interface necessária.

### Construtores das views dependentes de vaultManager

Os construtores atuais não recebem argumentos. Precisam ser atualizados para receber `*vault.Manager` — mesmo sendo stubs hoje, para manter a assinatura consistente com o futuro:

| Package | Construtor atual (no código) | Novo construtor |
|---|---|---|
| `settings` | `NewSettingsView()` | `NewSettingsView(vm *vault.Manager)` |
| `secret` | `NewTreeView()` → renomeado para `NewVaultTreeView()` | `NewVaultTreeView(vm *vault.Manager)` |
| `secret` | `NewDetailView()` → renomeado para `NewSecretDetailView()` | `NewSecretDetailView(vm *vault.Manager)` |
| `template` | `NewListView()` → renomeado para `NewTemplateListView()` | `NewTemplateListView(vm *vault.Manager)` |
| `template` | `NewDetailView()` → renomeado para `NewTemplateDetailView()` | `NewTemplateDetailView(vm *vault.Manager)` |

`WelcomeView`, `HeaderView`, `MessageLineView`, `ActionLineView` **não** recebem `vaultManager`. Construtores não mudam.

---

## Criação das Views Dependentes de VaultManager

`initVaultViews` é um método novo, não existe no código atual:

```go
// initVaultViews cria as views que dependem do vaultManager.
// Chamado quando vaultManager é definido — na inicialização ou no ciclo de vida.
func (r *RootModel) initVaultViews() {
    r.settingsView  = settings.NewSettingsView(r.vaultManager)
    r.secretTree    = secret.NewVaultTreeView(r.vaultManager)
    r.secretDetail  = secret.NewSecretDetailView(r.vaultManager)
    r.templateList  = template.NewTemplateListView(r.vaultManager)
    r.templateDetail = template.NewTemplateDetailView(r.vaultManager)
}
```

Quando `workArea` muda, `activeView` é atualizado para a view padrão da nova área:

| `workArea`          | `activeView` padrão |
|---------------------|---------------------|
| `WorkAreaWelcome`   | `&m.welcomeView`    |
| `WorkAreaSettings`  | `r.settingsView`    |
| `WorkAreaVault`     | `r.secretTree`      |
| `WorkAreaTemplates` | `r.templateList`    |

---

## Renderização

### Verificação de dimensões mínimas

`View()` verifica dimensões antes de qualquer renderização:

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
    // ...
}
```

### View()

O base layout é sempre renderizado. Quando há modal ativo, ele é sobreposto ao base via `lipgloss.NewLayer` + `lipgloss.NewCompositor` (API de compositing do lipgloss v2):

```go
func (r *RootModel) View() tea.View {
    // ... verificações de dimensão acima ...

    base := lipgloss.JoinVertical(lipgloss.Left,
        r.headerView.Render(design.HeaderHeight, r.width, r.theme),
        r.renderWorkArea(),
        r.messageLineView.Render(design.MessageHeight, r.width, r.theme, r.currentMessage),
        r.actionLineView.Render(design.ActionHeight, r.width, r.theme, r.ActiveViewActions()),
    )

    if len(r.modals) > 0 {
        top := r.modals[len(r.modals)-1]
        // Padding de 1 linha acima e abaixo do modal
        modalH := r.height - 2
        modalContent := top.Render(modalH, r.width, r.theme)
        // Centraliza o modal sobre o base usando o compositor de camadas do lipgloss v2.
        // NewLayer(base) é a camada de fundo; NewLayer(modal).Y(1).Z(1) flutua sobre ela.
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

### renderWorkArea()

**Correção de bug:** o código atual chama `r.activeView.Render(r.width, r.height, r.theme)` com `width` e `height` trocados. `renderWorkArea()` usa a ordem correta: `Render(height, width, theme)`.

```go
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
```

---

## Propagação de Eventos

### Teclado

O RootModel mantém a lógica atual de `evalActions` e `Update`:

1. System actions (sempre avaliadas)
2. Modal ativo recebe a tecla (se houver)
3. View actions de `activeView`
4. Application actions

`evalActions` é um método em `*RootModel` com assinatura `evalActions(msg tea.KeyMsg, actions []Action) (tea.Cmd, bool)`. Não muda. Internamente acessa `r` como receiver (que satisfaz `AppState`) e `r.activeView` diretamente do struct.

`activeView` nunca é nil após a inicialização. Os guards `if r.activeView != nil` presentes no código atual devem ser **todos removidos**. Há três ocorrências:

- `root.go` — em `ActiveViewActions()` (guarda `r.activeView.Actions()`)
- `root.go` — em `Update()`, branch de `tea.KeyMsg` (guarda a chamada a `evalActions`)
- `root.go` — em `Update()`, path de mensagens gerais (guarda `r.activeView.Update(msg)`)

`ActiveViewActions()` não muda de contrato — agrega `r.systemActions + r.applicationActions + r.activeView.Actions()`.

### Propagação de tea.Msg para ChildView

O método `Update` propaga mensagens para `activeView` e `headerView` (ambos `ChildView`). `MessageLineView` e `ActionLineView` não participam do loop do Bubble Tea:

```go
// No Update, ao final (sem modal ativo):
var cmds []tea.Cmd
cmds = append(cmds, r.activeView.Update(msg))
cmds = append(cmds, r.headerView.Update(msg))
return r, tea.Batch(cmds...)
```

---

## Estado currentMessage

`currentMessage string` é um campo novo no RootModel que representa a mensagem de status exibida na `MessageLineView`.

**Ciclo de vida:**

- Inicializado como `""` (zero value)
- Atualizado via `SetMessage(msg string)` — método novo, chamado sincronamente por qualquer parte da aplicação que já tenha referência ao `RootModel`
- Limpo via `SetMessage("")`

```go
// SetMessage define a mensagem de status exibida na barra inferior.
// Passe string vazia para limpar.
func (r *RootModel) SetMessage(msg string) {
    r.currentMessage = msg
}
```

Timer de auto-fade é trabalho futuro — usará `tea.Cmd` com `time.After`.

---

## Constantes a Adicionar em design/design.go

```go
// HeaderHeight é a altura em linhas da região de cabeçalho.
const HeaderHeight = 2

// MessageHeight é a altura em linhas da barra de mensagens.
const MessageHeight = 1

// ActionHeight é a altura em linhas da barra de ações.
const ActionHeight = 1
```

---

## O que NÃO está neste design

- **Limpeza automática de `currentMessage`** — timer de auto-fade implementar quando houver casos de uso concretos.
- **Transição de WorkArea** — a mudança de `workArea` e atualização de `activeView` será implementada junto com os casos de uso concretos (ex: abrir cofre).
- **Estado concreto de HeaderView** — stub por enquanto.
- **Ciclo de vida do vaultManager** — implementar em sprint separada.
- **Split ratio configurável** — `design.PanelTreeRatio` é usado diretamente; configuração por usuário é trabalho futuro.

---

## Resumo das Mudanças em Relação ao Código Atual

| Aspecto | Antes | Depois |
|---|---|---|
| `activeView` na inicialização | `nil` (panic) | `&m.welcomeView` |
| `workArea` | Ignorado | Controla `renderWorkArea()` |
| `headerView` | Não existe no RootModel | Campo `screen.HeaderView`, implementa `ChildView` |
| `messageLineView` | Não existe no RootModel | Campo `screen.MessageLineView`, caso especial sem interface |
| `actionLineView` | Não existe no RootModel | Campo `screen.ActionLineView`, caso especial sem interface |
| `currentMessage` | Não existe | Campo `string` no RootModel, atualizado via `SetMessage()` |
| Guards `activeView != nil` | 3 ocorrências em `root.go` | Todos removidos — `activeView` nunca é nil pós-init |
| Bug `Render(width, height)` | Argumentos trocados em `View()` | Corrigido em `renderWorkArea()` com ordem `(height, width)` |
| Composição vertical | Não implementada | `lipgloss.JoinVertical` |
| Composição horizontal (Vault/Templates) | Não implementada | `lipgloss.JoinHorizontal` com `design.PanelTreeRatio` |
| Modal renderizado | Sobre work area, tela limpa | Compositor `lipgloss.NewLayer` + `NewCompositor` sobre o base |
| Alturas de regiões | Magic number `-4` | Constantes `design.HeaderHeight` etc. |
| Verificação de terminal pequeno | Não existe | Mensagens "Aumente largura/altura do terminal!" |
| Construtores dependentes de vault | Sem `vaultManager` | Todos recebem `*vault.Manager` |
| Nomes dos tipos em `secret` | `TreeView`, `DetailView` | `VaultTreeView`, `SecretDetailView` |
| Nomes dos tipos em `template` | `ListView`, `DetailView` | `TemplateListView`, `TemplateDetailView` |
| Imports em `root.go` | Sem `screen`, `settings`, `secret`, `template` | Todos adicionados |
