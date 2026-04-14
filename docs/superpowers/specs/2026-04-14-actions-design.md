# Design: Sistema de Actions

**Data:** 2026-04-14  
**Status:** Aguardando aprovação

---

## 1. Visão Geral

Actions associam teclas a comportamentos da aplicação. Cada action tem pré-condições que controlam disponibilidade e visibilidade, e uma função que produz um `tea.Cmd`. O `RootModel` é responsável por avaliar e despachar as actions a cada tecla pressionada.

---

## 2. Estruturas de Dados

### 2.1 ActionGroup

Agrupa actions relacionadas no modal de ajuda.

```go
// ActionGroup agrupa actions relacionadas para exibição no modal de ajuda.
type ActionGroup struct {
    ID          string // identificador único do grupo
    Label       string // cabeçalho exibido no modal de ajuda
    Description string // texto descritivo do grupo
}
```

### 2.2 Action

```go
// Action associa teclas a um comportamento da aplicação.
type Action struct {
    Keys          []design.Key // Keys[0] é a tecla principal exibida; demais são aliases funcionais
    Label         string       // texto curto para a linha de status
    Description   string       // texto longo para o modal de ajuda
    GroupID       string       // referencia um ActionGroup registrado
    Priority      int          // ordenação na linha de status; menor valor = mais destaque
    Visible       bool         // false: nunca aparece na linha de status, independente da pré-condição
    AvailableWhen func(app AppState, view ChildView) bool // nil = sempre disponível
    OnExecute     func() tea.Cmd
}
```

**Sobre `Keys`:** `Keys[0]` é a tecla principal exibida na UI; as demais são aliases funcionais. O matching é feito via o método `Matches`:

```go
// Matches retorna true se o evento de teclado corresponde a qualquer tecla declarada na action.
func (a Action) Matches(msg tea.KeyMsg) bool {
    for _, k := range a.Keys {
        if k.Matches(msg) {
            return true
        }
    }
    return false
}
```

No fluxo de avaliação do root, o `tea.KeyMsg` é passado diretamente para `action.Matches(msg)`, que internamente delega para `design.Key.Matches(msg)` em cada tecla do slice.

**Sobre `AvailableWhen`:** `nil` significa que a action está sempre disponível. Quando não-nil, recebe `AppState` (interface fina sobre o root) e `ChildView` atual (sem tipo concreto — type assertions para capacidades específicas).

**Sobre `OnExecute`:** retorna `tea.Cmd`, que pode ser `nil` (no-op válido). O root não inspeciona o resultado — repassa de forma transparente ao loop do Bubble Tea.

---

## 3. Interface AppState

Interface fina implementada pelo `RootModel`, passada às pré-condições para desacoplar actions de detalhes internos do root.

```go
// AppState expõe o estado da aplicação necessário para avaliar pré-condições de actions.
// Implementado por RootModel.
type AppState interface {
    Manager() *vault.Manager // nil se nenhum cofre estiver carregado
}
```

`vault.Manager.Vault()` já retorna `nil` quando bloqueado — a pré-condição pode checar `app.Manager() != nil` para cofre carregado e `app.Manager().Vault() != nil` para cofre desbloqueado.

---

## 4. Interface ChildView (atualizada)

A interface atual em `internal/tui/view.go` possui 5 métodos. Duas mudanças são aplicadas:

1. `Render` passa a receber `*design.Theme` (ponteiro) em vez de `design.Theme` (valor) — `Theme` tem 400 bytes; cópia por valor a cada frame é desperdício desnecessário.
2. `Actions()` é adicionado.

```go
// ChildView define o contrato para views filhas gerenciadas pelo RootModel.
type ChildView interface {
    Render(height, width int, theme *design.Theme) string // theme por ponteiro — 400 bytes
    HandleKey(msg tea.KeyMsg) tea.Cmd
    HandleEvent(event any)
    HandleTeaMsg(msg tea.Msg) tea.Cmd
    Update(msg tea.Msg) tea.Cmd
    Actions() []Action // NOVO — pode retornar nil se a view não possuir actions
}
```

A mesma correção se aplica a `ModalView.Render` em `internal/tui/modal.go`:

```go
type ModalView interface {
    Render(maxHeight, maxWidth int, theme *design.Theme) string // theme por ponteiro
    HandleKey(msg tea.KeyMsg) tea.Cmd
    Update(msg tea.Msg) tea.Cmd
}
```

`root.go` passa `r.theme` diretamente (sem desreferenciar com `*`).

As View Actions são estáveis por natureza. Disponibilidade contextual é controlada por `AvailableWhen`, não por alternar o slice dinamicamente.

### Views existentes que precisam ser atualizadas

Todas as implementações atuais precisam de duas mudanças: ajustar `Render` para `*design.Theme` e adicionar `Actions() []Action` retornando `nil`:

| Arquivo                                        | Struct         |
|------------------------------------------------|----------------|
| `internal/tui/screen/welcome_view.go`          | `WelcomeView`  |
| `internal/tui/secret/vault_tree.go`            | `TreeView`     |
| `internal/tui/secret/secret_detail.go`         | `DetailView`   |
| `internal/tui/settings/settings_view.go`       | `SettingsView` |
| `internal/tui/template/template_list.go`       | `ListView`     |
| `internal/tui/template/template_detail.go`     | `DetailView`   |

As implementações de `ModalView` também precisam ajustar `Render`:

| Arquivo                                        | Struct           |
|------------------------------------------------|------------------|
| `internal/tui/modal/confirm_modal.go`          | `ConfirmModal`   |
| `internal/tui/modal/help_modal.go`             | `HelpModal` (novo)|

---

## 5. Hierarquia e Escopo

| Camada      | Definidas por                    | Ativas durante modal? |
|-------------|----------------------------------|-----------------------|
| System      | `RootModel.systemActions`        | Sim                   |
| View        | `r.activeView.Actions()`         | Não                   |
| Application | `RootModel.applicationActions`   | Não                   |

**System Actions** transcendem qualquer contexto — são avaliadas sempre, inclusive com modal ativo. Exemplo: F12 (trocar tema).

**View Actions** são definidas pela view ativa. Avaliadas apenas sem modal.

**Application Actions** são globais à aplicação, definidas no root. Avaliadas apenas sem modal, após View Actions. Exemplos: F1 (ajuda), Ctrl+Q (sair).

Conflitos de tecla são resolvidos silenciosamente por precedência: **System > View > Application**.

---

## 6. Fluxo de Avaliação de Teclas

A cada `tea.KeyMsg`, o root percorre as camadas em ordem. A **primeira** action cuja tecla corresponda **e** cuja pré-condição esteja satisfeita é executada; as demais não são avaliadas. A lista é recomputada a cada keypress — sem cache.

```
1. Percorre systemActions
   → action.Matches(msg) + AvailableWhen satisfeita (passa r como AppState, r.activeView como ChildView)?
     → sim: chama OnExecute(), retorna Cmd ao loop. Para aqui.

2. Há modal ativo?
   → sim: repassa KeyMsg ao modal. Para aqui.

3. r.activeView != nil? Percorre r.activeView.Actions()
   → action.Matches(msg) + AvailableWhen satisfeita?
     → sim: chama OnExecute(), retorna Cmd ao loop. Para aqui.

4. Percorre applicationActions
   → action.Matches(msg) + AvailableWhen satisfeita?
     → sim: chama OnExecute(), retorna Cmd ao loop. Para aqui.

5. Nenhuma correspondência → tecla descartada silenciosamente.
```

**Nota sobre `activeView nil`:** quando `r.activeView == nil`, os passos 3 e 4 ainda são executados, mas `nil` é passado como `ChildView` às pré-condições. Pré-condições que fazem type assertion em `view` devem checar `view != nil` antes.

---

## 7. Registro no RootModel

### 7.1 Campos adicionados ao RootModel

```go
type RootModel struct {
    // ... campos existentes ...
    vaultManager       *vault.Manager
    systemActions      []Action
    applicationActions []Action
    actionGroups       []ActionGroup
}
```

### 7.2 Métodos de registro

Os métodos de registro fazem **append** — chamar múltiplas vezes acumula as listas.

```go
func (r *RootModel) RegisterActionGroup(group ActionGroup)
func (r *RootModel) RegisterSystemActions(actions []Action)
func (r *RootModel) RegisterApplicationActions(actions []Action)
```

### 7.3 Método Manager() no RootModel

`RootModel` implementa `AppState` expondo o vault manager:

```go
// Manager retorna o vault manager ativo, ou nil se nenhum cofre estiver carregado.
func (r *RootModel) Manager() *vault.Manager {
    return r.vaultManager
}
```

### 7.4 Arquivo actions_system.go

Define e registra as System Actions. O registro do grupo é feito junto ao registro da action correspondente.

```go
// SetupSystemActions registra os grupos e actions de sistema no root.
func (r *RootModel) SetupSystemActions() {
    r.RegisterActionGroup(ActionGroup{
        ID:    "system",
        Label: "Sistema",
    })
    r.RegisterSystemActions([]Action{
        {
            Keys:      []design.Key{design.Shortcuts.ThemeToggle},
            Label:     "Tema",
            GroupID:   "system",
            Priority:  100,
            Visible:   false, // atalho existe mas não é anunciado na barra
            OnExecute: func() tea.Cmd { r.toggleTheme(); return nil },
        },
    })
}
```

### 7.5 Arquivo actions_application.go

Define e registra as Application Actions.

```go
// SetupApplicationActions registra os grupos e actions de aplicação no root.
func (r *RootModel) SetupApplicationActions() {
    r.RegisterActionGroup(ActionGroup{
        ID:    "app",
        Label: "Aplicação",
    })
    r.RegisterApplicationActions([]Action{
        {
            Keys:      []design.Key{design.Shortcuts.Help},
            Label:     "Ajuda",
            GroupID:   "app",
            Priority:  10,
            Visible:   true,
            OnExecute: func() tea.Cmd {
                var viewActions []Action
                if r.activeView != nil {
                    viewActions = r.activeView.Actions()
                }
                actions := slices.Concat(r.systemActions, r.applicationActions, viewActions)
                return OpenModal(NewHelpModal(actions, r.actionGroups))
            },
        },
        {
            Keys:      []design.Key{design.Shortcuts.Quit},
            Label:     "Sair",
            GroupID:   "app",
            Priority:  20,
            Visible:   true,
            OnExecute: func() tea.Cmd { return tea.Quit },
        },
    })
}
```

### 7.6 Orquestração

O caller (main ou builder) cria o RootModel e chama os métodos de setup:

```go
root := tui.NewRootModel(tui.WithVersion(version))
root.SetupSystemActions()
root.SetupApplicationActions()
```

---

## 8. toggleTheme no RootModel

Alterna entre os dois temas disponíveis. Implementado em `root.go`.

```go
// toggleTheme alterna o tema ativo entre TokyoNight e Cyberpunk.
func (r *RootModel) toggleTheme() {
    if r.theme == design.TokyoNight {
        r.theme = design.Cyberpunk
    } else {
        r.theme = design.TokyoNight
    }
}
```

---

## 9. Linha de Actions e Linha de Mensagem

### 9.1 Natureza dos componentes

`ActionLineView` e `MessageLineView` **não implementam `ChildView`**. São renderizadores stateless — não têm estado interno, não processam eventos, e o conteúdo a renderizar é passado como parâmetro direto ao `Render`. O root os chama diretamente, sem passar por qualquer interface.

### 9.2 ActionLineView

Exibe as actions disponíveis no contexto atual, da esquerda para direita, ordenadas por `Priority` (menor valor primeiro).

Uma action aparece se e somente se **ambas** as condições forem verdadeiras:
1. `Visible == true`
2. `AvailableWhen` satisfeita (ou `nil`)

Durante modal ativo, a linha fica em branco — nenhuma action é exibida.

```go
// ActionLineView renderiza a linha de ações disponíveis no contexto atual.
// Não implementa ChildView — é um renderizador stateless chamado diretamente pelo root.
type ActionLineView struct{}

// Render retorna a linha de ações para exibição.
// actions é a lista já filtrada (Visible + AvailableWhen) e ordenada por Priority.
func (v *ActionLineView) Render(height, width int, theme *design.Theme, actions []Action) string {
    // stub — retorna string vazia
    return ""
}
```

### 9.3 MessageLineView

```go
// MessageLineView renderiza a linha de mensagem de sistema.
// Não implementa ChildView — é um renderizador stateless chamado diretamente pelo root.
type MessageLineView struct{}

// Render retorna a linha de mensagem para exibição.
func (v *MessageLineView) Render(height, width int, theme *design.Theme, message string) string {
    // stub — retorna string vazia
    return ""
}
```

### 9.4 HeaderView

`HeaderView` é igualmente stateless — não implementa `ChildView`.

```go
// HeaderView renderiza o cabeçalho fixo de 2 linhas da aplicação.
// Não implementa ChildView — é um renderizador stateless chamado diretamente pelo root.
type HeaderView struct{}

// Render retorna o cabeçalho com 2 linhas para exibição.
func (v *HeaderView) Render(height, width int, theme *design.Theme) string {
    // stub — retorna string vazia
    return ""
}
```

---

## 10. Modal de Ajuda (HelpModal)

### 10.1 Comportamento

- Aberto pela action F1 via `OpenModal(NewHelpModal(actions, r.actionGroups))`
- Exibe **todas** as actions (System + View + Application), independente de `Visible` ou pré-condição
- Actions com pré-condição não satisfeita são exibidas esmaecidas, sem nota adicional
- Agrupadas por `GroupID`, com `Label` e `Description` do `ActionGroup` como cabeçalho
- Implementação minimalista: lista simples, agrupada

### 10.2 Interface com o root

`HelpModal` não tem referência ao `RootModel`. A action F1 monta os dados no momento da execução e os passa ao construtor:

```go
OnExecute: func() tea.Cmd {
    var viewActions []Action
    if r.activeView != nil {
        viewActions = r.activeView.Actions()
    }
    actions := slices.Concat(
        r.systemActions,
        r.applicationActions,
        viewActions,
    )
    return OpenModal(NewHelpModal(actions, r.actionGroups))
},
```

`NewHelpModal(actions []Action, groups []ActionGroup)` — sem acoplamento ao root.

---

## 11. Layout de Tela

```
┌──────────────────────────────────┐
│ headerView         (2 linhas)    │
├──────────────────────────────────┤
│                                  │
│ activeView    (height - 4 linhas)│
│                                  │
├──────────────────────────────────┤
│ messageLineView    (1 linha)     │
├──────────────────────────────────┤
│ actionLineView     (1 linha)     │
└──────────────────────────────────┘
```

### 11.1 Stubs em tui/screen/

Todos os componentes de tela fixos ficam no pacote `screen`. `WelcomeView` implementa `ChildView`. `HeaderView`, `MessageLineView` e `ActionLineView` são renderizadores stateless independentes — **não implementam `ChildView`**.

| Arquivo                  | Componente        | Implementa ChildView? | Stub mínimo                              |
|--------------------------|-------------------|-----------------------|------------------------------------------|
| `screen/welcome_view.go` | `WelcomeView`     | Sim                   | já existe (ajustar pacote e assinatura)  |
| `screen/header_view.go`  | `HeaderView`      | Não                   | `Render(height, width int, theme *design.Theme) string` retorna `""` |
| `screen/message_view.go` | `MessageLineView` | Não                   | `Render(height, width int, theme *design.Theme, message string) string` retorna `""` |
| `screen/action_view.go`  | `ActionLineView`  | Não                   | `Render(height, width int, theme *design.Theme, actions []Action) string` retorna `""` |

**Nota:** `welcome_view.go` está no pacote `welcome` — deverá ser ajustado para `package screen`.

---

## 12. Arquivos a Criar/Modificar

| Arquivo                              | Ação       | Conteúdo                                                                  |
|--------------------------------------|------------|---------------------------------------------------------------------------|
| `internal/tui/action.go`             | Criar      | `Action`, `ActionGroup`, `AppState`                                       |
| `internal/tui/actions_system.go`     | Criar      | `SetupSystemActions()`                                                    |
| `internal/tui/actions_application.go`| Criar      | `SetupApplicationActions()`                                               |
| `internal/tui/root.go`               | Modificar  | campos, `toggleTheme()`, métodos de registro, `Manager()`, passa `r.theme` sem `*` |
| `internal/tui/view.go`               | Modificar  | `Render` recebe `*design.Theme`; adiciona `Actions() []Action`            |
| `internal/tui/modal.go`              | Modificar  | `ModalView.Render` recebe `*design.Theme`                                 |
| `internal/tui/modal/confirm_modal.go`| Modificar  | `Render` recebe `*design.Theme`                                           |
| `internal/tui/modal/help_modal.go`   | Criar      | `HelpModal` (stub minimalista)                                            |
| `internal/tui/screen/welcome_view.go`| Modificar  | `package screen`; `Render` recebe `*design.Theme`; adiciona `Actions()`  |
| `internal/tui/secret/vault_tree.go`  | Modificar  | `Render` recebe `*design.Theme`; adiciona `Actions()`                     |
| `internal/tui/secret/secret_detail.go`| Modificar| `Render` recebe `*design.Theme`; adiciona `Actions()`                     |
| `internal/tui/settings/settings_view.go`| Modificar| `Render` recebe `*design.Theme`; adiciona `Actions()`                    |
| `internal/tui/template/template_list.go`| Modificar| `Render` recebe `*design.Theme`; adiciona `Actions()`                    |
| `internal/tui/template/template_detail.go`| Modificar| `Render` recebe `*design.Theme`; adiciona `Actions()`                  |
| `internal/tui/screen/header_view.go` | Criar      | `HeaderView` stub                                                         |
| `internal/tui/screen/message_view.go`| Criar      | `MessageLineView` stub                                                    |
| `internal/tui/screen/action_view.go` | Criar      | `ActionLineView` stub                                                     |

---

## 13. Fora do Escopo desta Sprint

- `design.Shortcuts.LockVault` (bloqueio de emergência) — action não implementada nesta sprint
- Implementação visual real de `HeaderView`, `MessageLineView`, `ActionLineView`
- Fluxos multi-step iniciados por actions
- Cache de avaliação de actions
- Conflitos de tecla com erro explícito
