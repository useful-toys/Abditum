# Arquitetura TUI do Abditum

> Documento de referência arquitetural para o pacote `internal/tui`.  
> Baseado nas decisões capturadas em `.planning/phases/05-tui-scaffold-root-model/05-CONTEXT.md`.

---

## Visão Geral

A TUI do Abditum segue o modelo **ELM** (Model-Update-View) do Bubble Tea, mas com uma hierarquia de modelos customizada. Apenas o `rootModel` implementa a interface `tea.Model` do Bubble Tea. Todos os demais modelos implementam uma interface interna `childModel` mais simples.

```
tea.Program
    └── rootModel          (único tea.Model)
            ├── child models   (interface childModel)
            ├── modal stack    ([]*modalModel)
            ├── activeFlow     (flowHandler)
            └── shared services
                    ├── *vault.Manager
                    ├── *ActionManager
                    └── *MessageManager
```

---

## Interfaces Centrais

### `childModel`

Interface implementada por todos os modelos filhos de área de trabalho e pelos modais.

```go
type childModel interface {
    Update(tea.Msg) tea.Cmd       // muta in-place; retorna apenas Cmd
    View() string                  // retorna string, NÃO tea.View
    SetSize(w, h int)              // recebe o tamanho alocado pelo compositor
    Context() FlowContext          // expõe estado de navegação/seleção para despacho de fluxos
    ChildFlows() []flowDescriptor  // descritores de fluxos específicos do filho (nil se nenhum)
}
```

- `View()` retorna `string` (não `tea.View`). Somente `rootModel.View()` retorna `tea.View`.
- `Update()` muta o próprio estado via pointer receiver — sem retornar `(Model, Cmd)`.
- `rootModel` calcula o tamanho de cada filho e chama `SetSize()` ao receber `tea.WindowSizeMsg`.
- `Context()` preenche os campos de navegação/seleção do `FlowContext`; `rootModel` enriquece com campos de nível vault antes de consultar os candidatos.
- `ChildFlows()` é o **escape hatch para casos raros** em que um fluxo não pode ser completamente parametrizado a partir do `FlowContext` — quando `IsApplicable` ou `New` precisariam de estado interno do filho que o `FlowContext` não carrega e não deve carregar. Em todos os outros casos, o fluxo pertence ao `FlowRegistry` global. Os descritores seguem o mesmo contrato `IsApplicable(FlowContext)` — sem closures. São verificados **antes** dos fluxos globais no despacho por tecla.

### `FlowContext`

Estado completo de contexto montado por `rootModel` no momento do despacho. Alimentado de duas fontes:

```go
type FlowContext struct {
    // Preenchido por rootModel a partir de vault.Manager
    VaultOpen  bool
    VaultDirty bool
    // Preenchido pelo filho ativo via Context()
    FocusedFolder   *vault.Pasta
    FocusedSecret   *vault.Segredo
    SecretOpen      bool
    FocusedField    *vault.Campo
    FocusedTemplate *vault.ModeloSegredo
    Mode            int  // filho define: ex. modo view vs edit, painel esquerdo vs direito
}
```

`rootModel` monta o `FlowContext` final assim:
```go
ctx := m.activeChild.Context()           // filho preenche campos de navegação
ctx.VaultOpen = m.mgr.IsOpen()           // rootModel adiciona estado do cofre
ctx.VaultDirty = m.mgr.HasUnsavedChanges()
// ctx está completo — única fonte de verdade para fluxos globais
// Despacho: activeChild.ChildFlows() (escape hatch) → flows.ForKey(key, ctx) (regra geral)
```

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

## `rootModel` — Estrutura

```go
type rootModel struct {
    area          workArea
    mgr           *vault.Manager
    vaultPath     string
    width, height int
    lastActionAt  time.Time

    // Modelos filhos — nil = inativo (GC recolhe)
    preVault       *preVaultModel
    vaultTree      *vaultTreeModel
    secretDetail   *secretDetailModel
    templateList   *templateListModel
    templateDetail *templateDetailModel
    settings       *settingsModel

    // Stack de modais — LIFO
    modals         []*modalModel

    // Fluxo ativo — nil = nenhum fluxo em andamento
    activeFlow     flowHandler

    // Serviços compartilhados
    actions        *ActionManager
    messages       *MessageManager
}
```

**Regra de nil-safety:** filhos são armazenados sempre como ponteiros concretos, nunca como interface `childModel`. Um `*preVaultModel` nil guardado numa interface `childModel` **não é nil** em Go — isso é uma armadilha de compilação. A interface é usada apenas transitoriamente no helper `liveModels()`.

---

## Área de Trabalho (`workArea`)

O `rootModel` rastreia `area workArea` — descreve o que está montado na zona central da tela:

```go
type workArea int
const (
    workAreaPreVault  workArea = iota // tela de boas-vindas (ASCII art)
    workAreaVault                     // cofre aberto — árvore + detalhe
    workAreaTemplates                 // editor de modelos — lista + detalhe
    workAreaSettings                  // tela de configurações
)
```

| `workArea` | Conteúdo renderizado |
|---|---|
| `workAreaPreVault` | `preVaultModel` — ASCII art de boas-vindas, sem sub-estados |
| `workAreaVault` | `vaultTreeModel` (esquerda) + `secretDetailModel` (direita) |
| `workAreaTemplates` | `templateListModel` (esquerda) + `templateDetailModel` (direita) |
| `workAreaSettings` | `settingsModel` ocupa toda a área |

A área de trabalho **não muda durante fluxos** — o usuário vê a área atual com modais sobrepostos. A transição só ocorre após um fluxo concluir com sucesso.

---

## Layout do Frame

`rootModel.View()` compõe sempre as mesmas zonas via lipgloss:

```
┌─────────────────────────────────┐
│ Header                          │  ← nome do app, nome do cofre, indicador de alterações
├─────────────────────────────────┤
│ Message bar                     │  ← MessageManager.Current()
├─────────────────────────────────┤
│                                 │
│ Work area                       │  ← childModel ativo
│                                 │
├─────────────────────────────────┤
│ Command bar                     │  ← ActionManager.Visible()
└─────────────────────────────────┘
```

Modais são sobrepostos **acima** de todo o frame via `lipgloss.Place()`.

---

## Despacho de Mensagens

`rootModel.Update()` despacha na seguinte ordem de prioridade:

```
1. Atalhos globais (ctrl+Q, ?)          → sempre interceptados primeiro
        ↓ senão
2. activeFlow != nil                     → delega ao flowHandler ativo
        ↓ senão
3. stack de modais não vazia             → topmost modal recebe input
        ↓ senão
4. child ativo da área de trabalho       → recebe input
```

**Mensagens de domínio** (ex: `vaultChangedMsg{}`, `tickMsg`) são transmitidas para **todos** os modelos vivos via `liveModels()`:

```go
func (m *rootModel) liveModels() []childModel {
    // retorna todos os filhos não-nil + todos os modais na stack
}
```

---

## Stack de Modais

Modais são gerenciados como uma pilha LIFO em `rootModel.modals []*modalModel`:

- **Push:** `modals = append(modals, newModal(...))`
- **Pop:** `modals = modals[:len(modals)-1]`
- O modal do topo recebe input de teclado/mouse.
- Modais abaixo continuam vivos e recebem mensagens de domínio.
- Modais podem abrir outros modais (ex: file picker abrindo confirmação).

**Tipos de modal conhecidos:** file picker, entrada de senha, criação de senha, confirmação (sim/não), help, progresso/spinner.

---

## Fluxos (`flowHandler`)

Fluxos encapsulam orquestração multi-passo que seria verbosa inline no `rootModel`. Exemplos: abrir cofre, criar cofre, salvar como, trocar senha, bloquear, sair com confirmação.

### Dois níveis de operações

O critério de divisão é simples: **a operação precisa de modal ou goroutine assíncrona?**

| Nível | Mecanismo | Exemplos |
|---|---|---|
| **Operação simples** (sem modal, sem async) | **Cmd factory** em `mutations.go` | favoritar, marcar exclusão, reordenar, renomear pasta |
| **Fluxo orquestrado** (modal e/ou async) | `flowHandler` no `FlowRegistry` | abrir cofre, salvar como, alterar senha, sair com confirmação |

**Cmd factory — padrão para operações simples:**
```go
// mutations.go
func cmdMarkSecretDeleted(mgr *vault.Manager, id string) tea.Cmd {
    return func() tea.Msg {
        if err := mgr.MarkDeleted(id); err != nil {
            return operationFailedMsg{err}
        }
        return secretDeletedMsg{id: id}
    }
}
```
O child chama `return cmdMarkSecretDeleted(m.mgr, id)` no seu `Update()`. Nunca chama Manager e fabrica Cmd manualmente — a factory é o contrato que **amarra mutacão → mensagem**.

**Como um fluxo orquestrado funciona:**

```
tecla acionada → ctx = activeChild.Context() + rootModel enriches vault state
        → candidatos: activeChild.ChildFlows() DEPOIS flows.ForKey(key, ctx)
        → primeiro IsApplicable(ctx) que passa?
        ↓ sim
rootModel: activeFlow = descriptor.New(ctx)
        ↓
rootModel.Update() delega input → activeFlow.Update()
        ↓
flow empurra modais via pushModalMsg{}  (file picker → progress → password → progress)
        ↓
operação assíncrona conclui → flow emite vaultOpenedMsg{} (+ chainFlowMsg se encadeamento)
        ↓
rootModel: activeFlow = nil  |  transição de estado
        ↓ (se chainFlowMsg presente — processado no Update() seguinte)
rootModel: ctx reconstruído → ForKey(key, ctx) → novo activeFlow (ou ignorado se inaplicável)
```

Cada fluxo vive em arquivo próprio (`flow_open_vault.go`, `flow_create_vault.go`, etc.). O `rootModel` não conhece os passos internos de nenhum fluxo.

### Encadeamento de fluxos (`chainFlowMsg`)

Em casos excepcionais, um fluxo que conclui pode solicitar a execução imediata de outro:

```go
// Dentro do flowHandler, ao concluir:
return tea.Batch(
    func() tea.Msg { return vaultOpenedMsg{} },            // transição de estado
    func() tea.Msg { return chainFlowMsg{key: "..."} },    // solicita próximo flow
)
```

O Bubble Tea processa **uma mensagem por `Update()`**. A mensagem de domínio chega primeiro — estado completamente atualizado — e só então `chainFlowMsg` é processado. `rootModel` reconstrói o `FlowContext` do estado atual e despacha via `FlowRegistry.ForKey(key, ctx)`. Se o flow alvo não for encontrado ou `IsApplicable(ctx)` retornar `false`, a requisição é ignorada silenciosamente.

---

## Serviços Compartilhados

Três objetos são instanciados em `main.go` e passados ao `rootModel`, que os repassa a cada filho no construtor.

### `vault.Manager`

API para todas as operações sobre o cofre (domínio). Fonte primária de dados para os filhos.

### `ActionManager`

> **Analogia:** assim como `vault.Manager` é a API para operações sobre o cofre, `ActionManager` é a API para definir quais ações estão disponíveis em cada momento.

- Objeto Go puro — sem `tea.Cmd`, sem mensagens, sem Bubble Tea.
- **Escrita:** cada filho registra suas ações ao ficar ativo; limpa ao ser desativado.
- **Leitura (command bar):** `ActionManager.Visible()` — subconjunto priorizado para o espaço disponível.
- **Leitura (help modal):** `ActionManager.All()` — lista completa de todas as ações registradas, agrupadas.
- `rootModel` registra os atalhos globais (`ctrl+Q`, `?`) no startup.

### `MessageManager`

> **Analogia:** assim como `ActionManager` é a API para ações disponíveis, `MessageManager` é a API para definir qual mensagem/dica aparece na barra de mensagens.

- Objeto Go puro — sem `tea.Cmd`, sem mensagens.
- **Escrita:** qualquer filho chama `messages.Set(text)` de dentro do seu `Update()` — síncrono, sem Cmd.
- **Leitura (message bar):** `rootModel.View()` chama `messages.Current()` em cada frame.
- Como o Bubble Tea re-renderiza após todo `Update()`, o frame sempre reflete o estado atual sem nenhum mecanismo de notificação.
- Filhos **não leem** do `MessageManager` — é write-only para eles.

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

---

## Timers e Timeouts

`rootModel` é o único dono das decisões de timeout:

- Rastreia `lastActionAt time.Time` — atualizado a cada input significativo.
- No `tickMsg`, consulta o Manager: `mgr.IsLockExpired(lastActionAt)`, `mgr.IsClipboardExpired(lastActionAt)`.
- A lógica de duração e habilitação fica encapsulada no Manager — `rootModel` recebe apenas `bool`.
- Se um timeout disparou, `rootModel` emite uma **mensagem tipada** para todos os filhos (`lockTimeoutMsg{}`, `clipboardTimeoutMsg{}`).
- Filhos recebem `tickMsg` apenas para **atualizar UI periódica** (ex: relógio no header). Nunca para implementar lógica de timeout.
- O tick global (1 segundo) **não começa em `Init()`**. É iniciado como `tea.Cmd` ao entrar em `workAreaVault`.

---

## Comunicação Entre Modelos

| Situação | Mecanismo |
|---|---|
| Filho notifica mutação de domínio | Retorna `tea.Cmd` emitindo mensagem de domínio tipada (ver tabela abaixo) |
| Filho lê dados do cofre | Chama `vault.Manager` diretamente |
| Filho lê estado do app | A definir (accessor read-only no `rootModel` ou valores passados no construtor) |
| Filho registra ações disponíveis | Chama `ActionManager.Register(...)` |
| Filho define mensagem da barra | Chama `MessageManager.Set(text)` |
| Filho abre diálogo | Retorna `dialogs.Confirm(...)` como Cmd |
| Filho inicia fluxo multi-passo | Via tecla de atalho — `rootModel` consulta `FlowRegistry.ForKey(key, ctx)` e inicia o flow |
| Fluxo empurra modal | Retorna Cmd emitindo `pushModalMsg{}` |
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

Children que não necessitam de uma mensagem simplesmente a ignoram.

---

## Atalhos Globais

| Tecla | Comportamento |
|---|---|
| `ctrl+Q` | Quit global — confirmação modal se há alterações não salvas |
| `?` | Abre `helpModal` com todas as ações registradas no `ActionManager` |
| `ctrl+C` | **Não é quit** |
| `q` | **Não é quit global** |

---

## Ciclo de Vida dos Modelos

- Filhos são alocados ao entrar na area correspondente; `nil` ao sair.
- `nil` = inativo, sem memória retida. O GC recolhe o modelo antigo.
- A transição cria o novo filho via construtor, passando `mgr`, `actions`, `messages` e demais dependências.
- Modelos sensíveis (que retêm dados do cofre) têm sua memória zerada explicitamente antes de `nil`.
