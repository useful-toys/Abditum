# Arquitetura TUI do Abditum

> Documento de referência arquitetural para o pacote `internal/tui`.  
> Baseado nas decisões capturadas em `.planning/phases/05-tui-scaffold-root-model/05-CONTEXT.md`.

---

## Visão Geral

A TUI do Abditum segue o modelo **ELM** (Model-Update-View) do Bubble Tea, mas com uma hierarquia de modelos customizada. Apenas o `rootModel` implementa a interface `tea.Model` do Bubble Tea. Todos os demais modelos implementam interfaces internas mais simples.

```
tea.Program
    └── rootModel          (único tea.Model)
            ├── work children  (interface childModel)
            ├── modal stack    (interface modalView)
            ├── activeFlow     (flowHandler)
            └── shared services
                    ├── *vault.Manager
                    ├── *ActionManager
                    └── *MessageManager
```

---

## Interfaces Centrais

### `childModel`

Interface implementada por modelos filhos de área de trabalho.

```go
type childModel interface {
    Update(tea.Msg) tea.Cmd       // muta in-place; retorna apenas Cmd
    View() string                  // retorna string, NÃO tea.View
    SetSize(w, h int)              // recebe o tamanho alocado pelo compositor
}
```

- `View()` retorna `string` (não `tea.View`). Somente `rootModel.View()` retorna `tea.View`.
- `Update()` muta o próprio estado via pointer receiver — sem retornar `(Model, Cmd)`.
- `rootModel` calcula o tamanho de cada filho e chama `SetSize()` ao receber `tea.WindowSizeMsg`.
- **Children e modais são position-unaware:** renderizam conteúdo preenchendo exatamente o tamanho recebido via `SetSize()`. Posicionamento é exclusivamente responsabilidade do `rootModel`.

### `modalView`

Interface implementada por modais. Separada de `childModel` porque modais têm contrato diferente: não recebem tamanho alocado (se auto-dimensionam por conteúdo), não participam de despacho de ações, e têm ciclo de vida gerenciado pela stack.

```go
type modalView interface {
    Update(tea.Msg) tea.Cmd
    View() string
}
```

### `modalResult`

Interface marcadora implementada por mensagens que carregam o resultado de um modal de volta para o flow que o abriu. `rootModel` roteia essas mensagens **somente** para `activeFlow` — sem broadcast para children ou outros modais.

```go
type modalResult interface {
    isModalResult()
}

// Exemplos de tipos concretos:
type passwordEntryResult struct {
    Password  []byte
    Cancelled bool
}
func (passwordEntryResult) isModalResult() {}

type confirmResult struct {
    Confirmed bool
}
func (confirmResult) isModalResult() {}
```

Dados sensíveis (ex: bytes de senha) **nunca entram em broadcast** — ficam isolados no caminho flow ↔ modal.

### `flowHandler`

Interface implementada por fluxos multi-passo (abrir cofre, criar cofre, etc.).

```go
type flowHandler interface {
    Update(tea.Msg) tea.Cmd  // orquestra modais e operações assíncronas
}
```

- Sem `View()` — fluxos não renderizam nada diretamente; empurram modais na stack.
- Sem `SetSize()` — não necessário.

---

## `Action` — Unidade de Interação

`Action` é o objeto central de despacho de interação. Cada ação encapsula três perguntas: **qual tecla me aciona?**, **estou disponível agora?**, **o que devo fazer?**

```go
type ActionScope int
const (
    ScopeLocal  ActionScope = iota  // só quando não há flow/modal ativo
    ScopeGlobal                      // sempre — mesmo durante flow ou modal
)

type Action struct {
    Keys        []string        // teclas que disparam esta ação
                                // Keys[0] aparece na command bar
                                // todas aparecem no help screen
    Label       string          // nome curto — command bar e help
    Description string          // texto longo — só no help screen
    Group       string          // agrupamento no help screen
    Scope       ActionScope     // quando a ação pode disparar
    Enabled     func() bool     // lê estado do child no momento da chamada
    Handler     func() tea.Cmd  // retorna intenção de execução
}
```

`Enabled` e `Handler` são closures sobre o child que registrou a ação. Em Go, a closure captura o **ponteiro** — então `m.focused` é sempre o valor atual no momento da chamada, não no momento do registro.

`Enabled()` é a implementação do conceito de **"fluxo elegível"** definido em `fluxos.md` — materializa o "contexto necessário" de cada fluxo como predicado avaliado sob demanda.

### Dois tipos de `Handler`

| Tipo | `Handler` retorna | Exemplos |
|---|---|---|
| Operação simples | Cmd factory de `mutations.go` | favoritar, marcar exclusão, duplicar |
| Fluxo orquestrado | `startFlowMsg{flow: ...}` | abrir cofre, salvar como, alterar senha |

```go
// Operação simples — chama Cmd factory:
Handler: func() tea.Cmd {
    return cmdToggleFavorite(m.mgr, m.focused)
}

// Fluxo orquestrado — instancia e inicia um flowHandler:
Handler: func() tea.Cmd {
    return func() tea.Msg {
        return startFlowMsg{flow: newOpenVaultFlow(m.mgr)}
    }
}
```

### Action Factories por Domínio

Ações que operam sobre o mesmo tipo de entidade são definidas **uma única vez** em factories por domínio. Cada factory recebe accessors que o child fornece ao registrar:

```go
// actions_segredo.go
type SecretAccessors struct {
    GetSecret func() *vault.Segredo
    GetCampo  func() *vault.Campo
}

func ActionsSegredo(mgr *vault.Manager, a SecretAccessors) []Action {
    return []Action{
        {
            Keys:    []string{"f"},
            Label:   "Favoritar",
            Group:   "Segredo",
            Scope:   ScopeLocal,
            Enabled: func() bool { return a.GetSecret() != nil },
            Handler: func() tea.Cmd { return cmdToggleFavorite(mgr, a.GetSecret()) },
        },
        {
            Keys:    []string{"d", "delete"},
            Label:   "Excluir",
            Group:   "Segredo",
            Scope:   ScopeLocal,
            Enabled: func() bool { return a.GetSecret() != nil },
            Handler: func() tea.Cmd {
                return dialogs.Confirm("Excluir segredo?",
                    cmdMarkDeleted(mgr, a.GetSecret()), nil)
            },
        },
        // ... demais ações de segredo
    }
}
```

Children consomem a factory com seus accessors:

```go
// vaulttree.go
func newVaultTreeModel(mgr *vault.Manager, actions *ActionManager, ...) *vaultTreeModel {
    m := &vaultTreeModel{mgr: mgr}

    actions.Register(m,
        ActionsSegredo(mgr, SecretAccessors{
            GetSecret: func() *vault.Segredo { return m.focused },
            GetCampo:  func() *vault.Campo   { return m.focusedCampo },
        })...,
        ActionsPasta(mgr, PastaAccessors{
            GetPasta: func() *vault.Pasta { return m.focusedPasta },
        })...,
    )

    return m
}

// secretdetail.go
func newSecretDetailModel(mgr *vault.Manager, actions *ActionManager, ...) *secretDetailModel {
    m := &secretDetailModel{mgr: mgr}

    actions.Register(m,
        ActionsSegredo(mgr, SecretAccessors{
            GetSecret: func() *vault.Segredo { return m.secret },
            GetCampo:  func() *vault.Campo   { return m.campo },
        })...,
    )

    return m
}
```

Adicionar uma nova ação de segredo significa editar `ActionsSegredo` — automaticamente disponível em todos os painéis que a usam.

---

## rootModel — Estrutura

```go
type rootModel struct {
    area          workArea
    mgr           *vault.Manager
    vaultPath     string
    width, height int
    lastActionAt  time.Time

    // Modelos filhos — nil = inativo (GC recolhe)
    welcome        *welcomeModel
    vaultTree      *vaultTreeModel
    secretDetail   *secretDetailModel
    templateList   *templateListModel
    templateDetail *templateDetailModel
    settings       *settingsModel

    // Stack de modais — LIFO
    modals         []modalView

    // Fluxo ativo — nil = nenhum fluxo em andamento
    activeFlow     flowHandler

    // Serviços compartilhados
    actions        *ActionManager
    messages       *MessageManager
}
```

**Regra de nil-safety:** filhos são armazenados sempre como ponteiros concretos, nunca como interface `childModel`. Um `*welcomeModel` nil guardado numa interface `childModel` **não é nil** em Go — isso é uma armadilha de compilação. A interface é usada apenas transitoriamente no helper `liveWorkChildren()`.

---

## Área de Trabalho (`workArea`)

O `rootModel` rastreia `area workArea` — descreve o que está montado na zona central da tela:

```go
type workArea int
const (
    workAreaWelcome   workArea = iota // tela de boas-vindas (ASCII art)
    workAreaVault                     // cofre aberto — árvore + detalhe
    workAreaTemplates                 // editor de modelos — lista + detalhe
    workAreaSettings                  // tela de configurações
)
```

| `workArea` | Conteúdo renderizado |
|---|---|
| `workAreaWelcome` | `welcomeModel` — ASCII art de boas-vindas, sem sub-estados |
| `workAreaVault` | `vaultTreeModel` (esquerda) + `secretDetailModel` (direita) |
| `workAreaTemplates` | `templateListModel` (esquerda) + `templateDetailModel` (direita) |
| `workAreaSettings` | `settingsModel` ocupa toda a área |

A área de trabalho **não muda durante fluxos** — o usuário vê a área atual com modais sobrepostos. A transição só ocorre após um fluxo concluir com sucesso.

---

## Layout do Frame

> **Status das zonas:**
> - **Work area** — decisão fechada. O `rootModel` alternará entre `workAreaWelcome`, `workAreaVault`, `workAreaTemplates` e `workAreaSettings`.
> - **Demais zonas** (header, message bar, command bar) — bastante prováveis, mas não comprometidas. A estrutura exata do frame será definida na fase de implementação.

Layout de referência (intenção atual, sujeito a revisão):

```
┌─────────────────────────────────┐
│ Header                          │  ← nome do app, nome do cofre, indicador de alterações
├─────────────────────────────────┤
│ Message bar                     │  ← MessageManager.Current()
├─────────────────────────────────┤
│                                 │
│ Work area                       │  ← childModel ativo  [COMPROMETIDO]
│                                 │
├─────────────────────────────────┤
│ Command bar                     │  ← ActionManager.Visible()
└─────────────────────────────────┘
```

Modais serão sobrepostos **acima** de todo o frame, provavelmente via `lipgloss.Place()`.

---

## Despacho de Mensagens

### Input do usuário (teclas e mouse)

`rootModel.Update()` despacha input na seguinte ordem:

```
1. actions.Dispatch(key, inFlowOrModal)
        → ScopeGlobal: sempre elegível
        → ScopeLocal: só quando não há flow/modal ativo
        → verifica Enabled() antes de executar Handler()
        ↓ nenhuma ação encontrada
2. activeFlow != nil              → delega ao flowHandler ativo
        ↓ senão
3. stack de modais não vazia      → topmost modal recebe input
        ↓ senão
4. child ativo da área de trabalho → recebe input
```

```go
// rootModel.Update()
case tea.KeyPressMsg:
    m.messages.HandleInput()  // limpa mensagens com clearOnInput (ex: warning de lock)
    m.lastActionAt = time.Now()

    key := msg.String()
    inFlowOrModal := m.activeFlow != nil || len(m.modals) > 0

    if cmd := m.actions.Dispatch(key, inFlowOrModal); cmd != nil {
        return m, cmd
    }

    if m.activeFlow != nil {
        return m, m.activeFlow.Update(msg)
    }

    if len(m.modals) > 0 {
        return m, m.modals[len(m.modals)-1].Update(msg)
    }

    return m, m.activeChild().Update(msg)
```

### Mensagens de domínio

Mensagens de domínio (ex: `vaultChangedMsg{}`, `tickMsg`) são transmitidas para **todos** os modelos vivos:

```go
func (m *rootModel) broadcast(msg tea.Msg) []tea.Cmd {
    var cmds []tea.Cmd
    for _, c := range m.liveWorkChildren() {
        if cmd := c.Update(msg); cmd != nil { cmds = append(cmds, cmd) }
    }
    for _, modal := range m.modals {
        if cmd := modal.Update(msg); cmd != nil { cmds = append(cmds, cmd) }
    }
    return cmds
}

func (m *rootModel) liveWorkChildren() []childModel {
    var live []childModel
    if m.welcome != nil         { live = append(live, m.welcome) }
    if m.vaultTree != nil       { live = append(live, m.vaultTree) }
    if m.secretDetail != nil    { live = append(live, m.secretDetail) }
    if m.templateList != nil    { live = append(live, m.templateList) }
    if m.templateDetail != nil  { live = append(live, m.templateDetail) }
    if m.settings != nil        { live = append(live, m.settings) }
    return live
}
```

### Mensagens de resultado de modal (`modalResult`)

Mensagens que implementam `modalResult` são roteadas **somente** para `activeFlow`:

```go
case modalResult:
    if m.activeFlow != nil {
        return m, m.activeFlow.Update(msg)
    }
```

Dados sensíveis (bytes de senha) nunca entram em broadcast.

---

## Stack de Modais

Modais são gerenciados como uma pilha LIFO em `rootModel.modals []modalView`:

- **Push:** via `pushModalMsg{}` — `modals = append(modals, msg.modal)`. Nenhum child ou flow acessa a stack diretamente.
- **Pop por usuário:** via ESC ou seleção — o modal retorna `popModalMsg{}` como Cmd.
- **Pop programático:** o flow emite `popModalMsg{}` quando uma operação async conclui — fecha o modal sem ação do usuário.
- **Segurança do pop:** `ctrl+Q` é `ScopeLocal` — não dispara durante flows/modais. O único `ScopeGlobal` (`?`) empurra `helpModal` que é passivo (dismiss via ESC). Portanto, não há risco de um push externo intercalar com um pop pendente.
- **Invariante de callbacks:** callbacks `onYes`/`onNo` de `confirmModal` não devem ser `pushModalMsg` instantâneos. Se `onYes` precisa abrir outro modal, deve fazê-lo via `startFlowMsg` ou Cmd assíncrono — garantindo que o `popModalMsg` do confirm seja processado primeiro.
- O modal do topo recebe input de teclado/mouse (via passo 3 do despacho).
- Modais abaixo continuam vivos e recebem mensagens de domínio via `broadcast()`.
- Modais podem abrir outros modais (ex: confirmação abrindo outro modal de confirmação).

### Categorias de modal

| Modal | Retorno | Mecanismo |
|---|---|---|
| `confirmModal` | Decisão binária | Callbacks (`onYes`, `onNo` tea.Cmd) — contexto embutido |
| `messageModal` | Nenhum | Dismiss via ESC ou Enter |
| `passwordEntryModal` | `[]byte` | `modalResult` — roteado somente ao flow |
| `filePickerModal` | `string` (caminho) | `modalResult` — roteado somente ao flow |

**Feedback de progresso:** não existe `progressModal`. Operações assíncronas usam `MessageManager.Show(MsgBusy, ...)` na barra de mensagens — spinner animado a 1fps. O `activeFlow` já bloqueia input local (teclas caem no flow, que as ignora). Modal de progresso seria redundante nos três eixos: feedback visual, bloqueio de input, e animação.

**Tipos de modal (stubs no Phase 5):** entrada de senha, criação de senha, confirmação (sim/não), help.
**File picker modal:** adiado — implementado na fase que introduz seu primeiro caso de uso.

---

## Fluxos (`flowHandler`)

Fluxos encapsulam orquestração multi-passo que seria verbosa inline no `rootModel`. Exemplos: abrir cofre, criar cofre, salvar como, trocar senha, bloquear, sair com confirmação.

### Dois níveis de operações

O critério de divisão é simples: **a operação precisa de modal ou goroutine assíncrona?**

| Nível | Mecanismo | Exemplos |
|---|---|---|
| **Operação simples** (sem modal, sem async) | **Cmd factory** em `mutations.go` | favoritar, marcar exclusão, reordenar, renomear pasta |
| **Fluxo orquestrado** (modal e/ou async) | `flowHandler` via `Action.Handler` + `startFlowMsg` | abrir cofre, salvar como, alterar senha, sair com confirmação |

**Cmd factory — padrão para operações simples:**
```go
// mutations.go
func cmdToggleFavorite(mgr *vault.Manager, s *vault.Segredo) tea.Cmd {
    return func() tea.Msg {
        if err := mgr.ToggleFavorite(s); err != nil {
            return operationFailedMsg{err}
        }
        return secretModifiedMsg{s}
    }
}
```

**Como um fluxo orquestrado é acionado e funciona:**

```
Usuário pressiona tecla
        ↓
actions.Dispatch(key, inFlowOrModal)
        → Enabled() == true → Handler() retorna tea.Cmd
        ↓
Bubble Tea executa Cmd → startFlowMsg{flow: openVaultFlow{...}}
        ↓
rootModel.Update(startFlowMsg) → activeFlow = msg.flow
        ↓
rootModel.Update() delega input → activeFlow.Update()
        ↓
flow empurra modais via pushModalMsg{}
        ↓
modal coleta valor → emite modalResult → roteado ao flow
        ↓
flow executa operação async → emite popModalMsg{} + mensagem de domínio
        ↓
rootModel: activeFlow = nil | transição de estado
```

Cada fluxo vive em arquivo próprio (`flow_open_vault.go`, `flow_create_vault.go`, etc.). O `rootModel` não conhece os passos internos de nenhum fluxo.

### Encadeamento de fluxos

Em casos excepcionais, um fluxo que conclui pode solicitar a execução imediata de outro emitindo `startFlowMsg` diretamente:

```go
// Dentro do flowHandler, ao concluir:
return tea.Batch(
    func() tea.Msg { return vaultOpenedMsg{} },
    func() tea.Msg { return startFlowMsg{flow: newAutoSaveFlow(mgr)} },
)
```

O Bubble Tea processa **uma mensagem por `Update()`**. A mensagem de domínio chega primeiro — estado completamente atualizado — e só então `startFlowMsg` é processado.

**Nota:** encadeamento direto via `startFlowMsg` bypassa o `ActionManager` — não verifica `Enabled()`. Isso é aceitável porque o flow que encadeia já validou o estado.

---

## Serviços Compartilhados

Três objetos são instanciados em `main.go` e passados ao `rootModel`, que os repassa a cada filho no construtor.

### `vault.Manager`

API para todas as operações sobre o cofre (domínio). Fonte primária de dados para os filhos.

### `ActionManager`

> **Analogia:** assim como `vault.Manager` é a API para operações sobre o cofre, `ActionManager` é a API para definir quais ações estão disponíveis em cada momento **e o ponto único de despacho de input**.

- Objeto Go puro — sem `tea.Cmd`, sem mensagens, sem Bubble Tea.
- **Registro:** cada child registra suas ações no construtor via `actions.Register(owner, ...Action)`. `rootModel` registra ações de startup (`ctrl+Q` com `ScopeLocal`, `?` com `ScopeGlobal`).
- **Descarte:** `actions.ClearOwned(owner)` — chamado **antes** de setar o child para `nil` (invariante de ciclo de vida).
- **Dono ativo:** `actions.SetActiveOwner(owner)` — quando dois children estão vivos (`vaultTree` + `secretDetail`), `Dispatch` prioriza ações do dono ativo. Ações do `rootModel` (globais) são sempre elegíveis.
- **Despacho:** `actions.Dispatch(key string, inFlowOrModal bool) tea.Cmd` — verifica `Scope`, `Enabled()`, e executa `Handler()`.
- **Command bar:** `ActionManager.Visible()` — ações onde `Enabled() == true`, subconjunto priorizado para o espaço disponível.
- **Help modal:** `ActionManager.All()` — lista completa de todas as ações registradas, agrupadas por `Group`.

### `MessageManager`

> **Analogia:** assim como `ActionManager` é a API para ações disponíveis, `MessageManager` é a API para definir qual mensagem aparece na barra de mensagens — com tipo, duração e comportamento de descarte.

Objeto Go puro — sem `tea.Cmd`, sem mensagens Bubble Tea.

```go
type MsgKind int
const (
    MsgInfo  MsgKind = iota  // ✅ operação concluída com sucesso
    MsgWarn                   // ⚠️  atenção — bloqueio iminente, conflito externo
    MsgError                  // ❌ falha — salvamento, corrupção
    MsgBusy                   // ⏳ operação em andamento — salvando, exportando (spinner animado)
    MsgHint                   // 💡 explicação contextual — descrição de campo
)
```

```go
type MessageManager struct {
    current *activeMessage
}

type activeMessage struct {
    text         string
    kind         MsgKind
    startedAt    time.Time   // para calcular frame do spinner (MsgBusy)
    expiresAt    time.Time   // zero = permanente até substituição
    clearOnInput bool        // true = some ao próximo KeyPress/Mouse
}
```

**API:**

```go
// Escrita — children, flows e rootModel (dentro de Update, nunca em Cmd factories)
func (mm *MessageManager) Show(kind MsgKind, text string, ttlSeconds int, clearOnInput bool)
func (mm *MessageManager) Clear()

// Leitura — só rootModel.View()
type DisplayMessage struct {
    Text  string
    Kind  MsgKind
    Frame int      // índice de animação para MsgBusy (incrementa a cada segundo)
}
func (mm *MessageManager) Current() *DisplayMessage  // nil = sem mensagem

// Manutenção — só rootModel.Update()
func (mm *MessageManager) Tick()         // expira mensagens com TTL vencido
func (mm *MessageManager) HandleInput()  // limpa se clearOnInput == true
```

**Regras de uso:**

- **`Show()`/`Clear()` são chamados exclusivamente dentro de `Update()`** — nunca dentro de Cmd factories (`func() tea.Msg`). Cmd factories executam em goroutine separada no Bubble Tea; chamar `Show()` de lá causaria race condition.
- **Children e flows** chamam `Show()` ou `Clear()` de dentro do seu `Update()` — síncrono, seguro.
- O **`rootModel`** chama `Show()` ao processar mensagens de domínio retornadas por Cmd factories (ex: `secretModifiedMsg` → `Show(MsgInfo, "Favoritado", 2, false)`).
- **Prioridade:** last-write-wins — sem stack, sem fila. Se um `MsgInfo` ("Copiado", TTL=3s) sobrescrever um `MsgWarn` de lock, no próximo tick `IsLockWarning` re-emite o warning automaticamente.
- Filhos **não leem** do `MessageManager` — é write-only para eles.
- **Invariante de `MsgBusy`:** fluxos que emitem `Show(MsgBusy, ...)` devem emitir `Show()` ou `Clear()` em **todo** caminho de saída (sucesso, erro, cancelamento). `MsgBusy` não tem TTL — permanece até ser substituído.

**Renderização (responsabilidade do `rootModel.View()`):**

```go
if msg := m.messages.Current(); msg != nil {
    emoji := messageEmoji[msg.Kind]
    if msg.Kind == MsgBusy {
        frames := []string{"◐", "◓", "◑", "◒"}
        emoji = frames[msg.Frame % len(frames)]
    }
    messageBar = messageStyles[msg.Kind].Render(emoji + " " + msg.Text)
}
```

Estilos por `Kind` (cor + formatação) vivem na camada de renderização, não no manager.

**Exemplos de uso:**

```go
// rootModel.Update() — feedback de operação simples
case secretModifiedMsg:
    m.messages.Show(MsgInfo, "Favoritado", 2, false)
    return m, tea.Batch(m.broadcast(msg)...)

// flowHandler.Update() — progresso e resultado
case startSaving:
    f.msgs.Show(MsgBusy, "Salvando cofre...", 0, false)
    return cmdSaveVault(f.mgr)
case vaultSavedMsg:
    f.msgs.Show(MsgInfo, "Cofre salvo", 3, false)
    return endFlow()
case operationFailedMsg:
    f.msgs.Show(MsgError, msg.err.Error(), 5, false)
    return endFlow()

// child.Update() — hint contextual
m.messages.Show(MsgHint, campo.Description, 0, false)

// rootModel.Update(tickMsg) — aviso de bloqueio iminente
if m.mgr.IsLockWarning(m.lastActionAt) {
    m.messages.Show(MsgWarn, "Cofre será bloqueado em breve", 0, true)
}
```

---

## `dialogs` — Factory de Diálogos Pré-definidos

Diferente dos managers acima, `dialogs` não é estado compartilhado — são **funções puras** que produzem `tea.Cmd`:

```go
// informativo — dismiss via ESC ou Enter
dialogs.Message(title, text string) tea.Cmd

// pergunta sim/não — dispara onYes ou onNo conforme seleção
dialogs.Confirm(question string, onYes, onNo tea.Cmd) tea.Cmd
```

O Cmd emitido é um `pushModalMsg{}`. `rootModel.Update()` intercepta essa mensagem e empurra o modal na stack. Nenhum filho acessa a stack diretamente.

Callbacks (`onYes`, `onNo`) são adequados quando a decisão é binária e o contexto já é conhecido no momento da abertura. Para coleta de valores (senha, caminho), o modal emite `modalResult` em vez de usar callbacks.

---

## Timers e Timeouts

`rootModel` é o único dono de **todas** as decisões de timeout — lock, clipboard e ocultação de campo sensível.

### Justificativa: centralização no rootModel

Embora a ocultação de campo sensível (F16) seja visualmente local ao `secretDetailModel`, todos os três timers compartilham o mesmo padrão estrutural: um timestamp de reset, uma verificação por tick, e uma ação resultante. Centralizar no `rootModel` traz três benefícios:

1. **Localidade de raciocínio:** toda lógica temporal vive em um único `case tickMsg:` com 10-15 linhas. Distribuir entre `rootModel` e children criaria dois locais de verificação com o mesmo tick, sem eliminar complexidade.
2. **Coordenação com lock:** quando o lock dispara, o `rootModel` precisa garantir que clipboard e campo visível são limpos como parte do wipe de memória. Se o child controlasse o field hide, o `rootModel` precisaria de um mecanismo extra para forçar a limpeza — duplicando responsabilidade.
3. **Consistência:** os três timers usam a mesma infraestrutura (`vault.Manager.IsXxxExpired()`, timestamp no `rootModel`, mensagem tipada via broadcast). Patterns diferentes para o mesmo problema tornam o código mais difícil de manter.

A alternativa considerada — mover field hide para o child — resolvia melhor o caso de múltiplos campos revelados simultaneamente (um map `revealedAt` per-field). Porém, esse cenário é improvável na prática: o campo se oculta automaticamente após poucos segundos, e revelar um novo campo antes seria o caso normal. O `rootModel` pode emitir `fieldHideMsg{}` e o child decidir internamente quais campos ocultar, mantendo a decisão de *quando* centralizada e a decisão de *quais* encapsulada.

### Comportamento

- Rastreia `lastActionAt`, `lastCopyAt`, `lastRevealAt` — cada um resetado por evento diferente.
- No `tickMsg`, consulta o Manager: `mgr.IsLockExpired(lastActionAt)`, `mgr.IsClipboardExpired(lastCopyAt)`, `mgr.IsFieldHideExpired(lastRevealAt)`.
- A lógica de duração e habilitação fica encapsulada no Manager — `rootModel` recebe apenas `bool`.
- Se um timeout disparou, `rootModel` emite uma **mensagem tipada** para todos os filhos (`lockTimeoutMsg{}`, `clipboardTimeoutMsg{}`, `fieldHideMsg{}`).
- Aviso de bloqueio iminente usa `MessageManager.Show(MsgWarn, ..., 0, true)` — permanente até interação do usuário (ver seção MessageManager).
- Filhos recebem `tickMsg` apenas para **atualizar UI periódica** (ex: relógio no header). Nunca para implementar lógica de timeout.
- O tick global (1 segundo) **não começa em `Init()`**. É iniciado como `tea.Cmd` ao entrar em `workAreaVault`.

**Ordem de processamento no `tickMsg`:**

```go
case tickMsg:
    m.messages.Tick()  // 1. expira mensagens com TTL vencido

    if m.mgr.IsLockExpired(m.lastActionAt) {
        return m, startLockFlow(...)  // 2. lock tem prioridade absoluta
    }
    if m.mgr.IsLockWarning(m.lastActionAt) {
        m.messages.Show(MsgWarn, "Cofre será bloqueado em breve", 0, true)
    }
    // 3. clipboard
    if !m.lastCopyAt.IsZero() && m.mgr.IsClipboardExpired(m.lastCopyAt) {
        m.lastCopyAt = time.Time{}
        cmds = append(cmds, func() tea.Msg { return clipboardTimeoutMsg{} })
    }
    // 4. field hide
    if !m.lastRevealAt.IsZero() && m.mgr.IsFieldHideExpired(m.lastRevealAt) {
        m.lastRevealAt = time.Time{}
        cmds = append(cmds, func() tea.Msg { return fieldHideMsg{} })
    }
    cmds = append(cmds, m.broadcast(msg)...)  // 5. broadcast para children (UI periódica)
    cmds = append(cmds, tea.Tick(time.Second, ...))  // 6. re-agenda
    return m, tea.Batch(cmds...)
```

---

## Comunicação Entre Modelos

| Situação | Mecanismo |
|---|---|
| Filho notifica mutação de domínio | Retorna `tea.Cmd` emitindo mensagem de domínio tipada (ver tabela abaixo) |
| Filho lê dados do cofre | Chama `vault.Manager` diretamente |
| Filho lê estado do app | A definir (accessor read-only no `rootModel` ou valores passados no construtor) |
| Filho registra ações disponíveis | Via `ActionManager.Register(...)` no construtor |
| Filho define mensagem da barra | Chama `MessageManager.Show(kind, text, ttl, clearOnInput)` dentro de `Update()` |
| Filho abre diálogo | Retorna `dialogs.Confirm(...)` como Cmd |
| Action inicia fluxo multi-passo | `Handler` retorna `startFlowMsg{flow: ...}` |
| Fluxo empurra modal | Retorna Cmd emitindo `pushModalMsg{}` |
| Fluxo fecha modal programaticamente | Retorna Cmd emitindo `popModalMsg{}` |
| Modal devolve valor ao flow | Emite `modalResult` — roteado somente ao `activeFlow` |
| Fluxo conclui | Retorna Cmd emitindo mensagem de conclusão (ex: `vaultOpenedMsg{}`) |

**Regra absoluta:** nenhum filho acessa campos de outro filho. Toda comunicação passa pelo `rootModel` via mensagens de domínio.

### Taxonomia de mensagens de domínio

O Bubble Tea re-renderiza a tela inteira após todo `Update()`, então mensagens granulares não reduzem trabalho de renderização. Seu valor é permitir que children tomem decisões locais precisas — `secretDetailModel` ignora `secretReorderedMsg` sem consultar o Manager.

| Mensagem | Signficado |
|---|---|
| `secretAddedMsg{id}` | Segredo criado ou duplicado |
| `secretDeletedMsg{id}` | Segredo marcado para exclusão |
| `secretRestoredMsg{id}` | Marcação de exclusão removida |
| `secretModifiedMsg{id}` | Valores ou estrutura do segredo alterados |
| `secretMovedMsg{id, fromFolder, toFolder}` | Segredo movido entre pastas |
| `secretReorderedMsg{}` | Segredo reordenado dentro de uma pasta |
| `folderStructureChangedMsg{}` | Qualquer create/rename/move/reorder/delete de pasta |
| `vaultSavedMsg{}` | Cofre gravado em disco (segredos excluídos removidos da memória) |
| `vaultReloadedMsg{}` | Recarga completa do disco — todos os children resetam estado |
| `vaultClosedMsg{}` | Cofre bloqueado ou fechado — todos os children limpam memória sensível |
| `vaultChangedMsg{}` | Fallback genérico — usado por fluxos quando o tipo de mutação não é relevante para broadcast |
| `startFlowMsg{flow}` | Inicia um flowHandler — interceptado pelo `rootModel` |
| `pushModalMsg{modal}` | Empurra modal na stack — interceptado pelo `rootModel` |
| `popModalMsg{}` | Remove topmost modal — interceptado pelo `rootModel` |

Children que não necessitam de uma mensagem simplesmente a ignoram.

---

## Atalhos Globais

Atalhos registrados pelo `rootModel` no startup. Passam pelo `ActionManager.Dispatch()` como qualquer outra ação — sem interceptação hardcoded.

| Tecla | Comportamento | Scope | Justificativa |
|---|---|---|---|
| `ctrl+Q` | Quit — confirmação modal se há alterações não salvas | `ScopeLocal` | Durante flow/modal ativo, quit causaria conflitos: sobrescrita de `activeFlow`, modais órfãos na stack, Cmds assíncronos retornando para o flow errado. O caminho seguro é ESC (fecha modal/flow) → `ctrl+Q` |
| `?` | Abre `helpModal` com todas as ações registradas no `ActionManager` | `ScopeGlobal` | Help é passivo — overlay informacional sem estado, dismiss via ESC, sem conflito de flow |
| `ctrl+C` | **Não é quit** | — | |
| `q` | **Não é quit global** | — | |

---

## Ciclo de Vida dos Modelos

- Filhos são alocados ao entrar na area correspondente; `nil` ao sair.
- `nil` = inativo, sem memória retida. O GC recolhe o modelo antigo.
- A transição cria o novo filho via construtor, passando `mgr`, `actions`, `messages` e demais dependências.
- O construtor do child registra suas ações via `actions.Register(m, ...)` — ações vivem enquanto o child vive.
- **Invariante de desativação:** ao trocar de workArea, **sempre** chamar `actions.ClearOwned(child)` ANTES de setar o child para `nil`. Closures nas Actions seguram referência ao child — `ClearOwned` remove as ações antes que o ponteiro seja descartado.
- Modelos sensíveis (que retêm dados do cofre) têm sua memória zerada explicitamente antes de `nil`.
