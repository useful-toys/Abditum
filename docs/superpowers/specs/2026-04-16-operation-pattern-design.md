# Design: Padrão Operation no Bubble Tea

**Contexto:** Abditum — gerenciador de cofre TUI em Go (Bubble Tea v2, bubbletea/v2)

## Problema

Em TEA (The Elm Architecture), o `Update` deve ser puro e síncrono. Operações reais como abrir cofre, criar cofre e exportar segredos exigem:

- Múltiplas etapas sequenciais com estado entre elas
- Coleta de input do usuário no meio do fluxo (modais)
- Trabalho assíncrono pesado (IO, criptografia) com feedback de progresso
- Recuperação de erros e retry (ex: senha incorreta)

Colocar essa lógica no `RootModel` resulta em flags espalhadas e métodos sem coesão. O padrão *Operation* encapsula cada operação como uma mini-máquina de estados autônoma.

## Decisões de Design

**A operação é uma máquina de estados com `Update`.**
Análogo a como as views filhas já funcionam — cada operação tem `Init()` e `Update(msg)`, processa apenas as mensagens que reconhece, e retorna `nil` para as demais.

**A operação é iniciada por mensagem, não por chamada direta.**
Simétrico ao padrão de modais (`OpenModalMsg` / `CloseModalMsg`). A ação dispara `StartOperation(op)` que emite `StartOperationMsg`. Root trata essa mensagem e chama `op.Init()`.

**O root roteia mensagens genericamente.**
Root não conhece as mensagens internas de nenhuma operação. Ele roteia *tudo* para `activeOperation.Update(msg)`. A operação ignora o que não reconhece.

**Modais são empurrados pela própria operação.**
A operação usa `tui.OpenModal()` diretamente (mesmo padrão de `tui/modal` que já importa `tui`). Root não precisa de lógica específica para "abrir modal quando operação pede".

**Mensagens terminais ficam em `tui`.**
`VaultOpenedMsg`, `SecretExportedMsg` etc. são definidas em `tui` porque root precisa tratá-las. A operação emite essas mensagens por conhecer o pacote `tui`.

**Encadeamento via `StartOperationMsg`.**
Uma operação pode encadear outra emitindo `StartOperationMsg{Op: nextOp}` em vez de `OperationCompleted()`. Root substitui `activeOperation` sem passar por idle. Isso permite `CreateVaultOperation` → `OpenVaultOperation` naturalmente.

## Interface e Mensagens

```go
// Em tui/operation.go

// Operation encapsula uma operação multi-etapa com estado próprio.
// Cada operação é uma mini-máquina de estados que processa mensagens
// e avança seu fluxo de forma autônoma via comandos Tea.
type Operation interface {
    // Init retorna o primeiro comando da operação — normalmente abre um modal
    // para coletar input inicial. Análogo a tea.Model.Init().
    Init() tea.Cmd
    // Update processa mensagens Tea e avança o estado interno.
    // Retorna um comando ou nil. A operação ignora mensagens que não reconhece.
    Update(msg tea.Msg) tea.Cmd
}

// StartOperationMsg inicia uma operação (e encerra a atual, se houver).
// Emitida por ações no setup.go ou por operações encadeando outra.
type StartOperationMsg struct{ Op Operation }

// OperationCompletedMsg sinaliza conclusão da operação ativa sem continuação.
// Root limpa activeOperation ao receber esta mensagem.
type OperationCompletedMsg struct{}

// VaultOpenedMsg sinaliza que um cofre foi aberto ou criado com sucesso.
// Emitida pela operação; tratada pelo root para configurar vaultManager.
type VaultOpenedMsg struct{ Manager *vault.Manager }

// SecretExportedMsg sinaliza exportação bem-sucedida de um segredo.
type SecretExportedMsg struct{}

// StartOperation cria um Cmd que emite StartOperationMsg.
func StartOperation(op Operation) tea.Cmd { ... }

// OperationCompleted cria um Cmd que emite OperationCompletedMsg.
func OperationCompleted() tea.Cmd { ... }
```

## Topologia de Imports

```
tui/operation  →  tui, tui/modal, tui/design, vault
tui/modal      →  tui, tui/design, tui/actions   (sem mudança)
tui            →  NÃO importa tui/operation
cmd/main       →  tui, tui/operation, tui/modal   (conecta tudo)
```

Sem ciclos. `tui/operation` segue o mesmo padrão que `tui/modal` já usa.

## Mudanças no RootModel

**Campo novo:**
```go
activeOperation Operation  // nil quando ocioso
```

**Novos cases no Update:**
```go
case StartOperationMsg:
    r.activeOperation = msg.Op
    return r, msg.Op.Init()

case OperationCompletedMsg:
    r.activeOperation = nil
    return r, nil

case VaultOpenedMsg:
    r.setVaultManager(msg.Manager)
    r.setWorkArea(design.WorkAreaVault)  // chamada direta — setWorkArea já existe em root.go
    return r, nil

case SecretExportedMsg:
    return r, nil
```

**Roteamento genérico** — substitui o bloco final do `Update` (após o switch):

```go
// Antes (sem Operation):
if len(r.modals) > 0 {
    return r, r.modals[top].Update(msg)
}
var cmds []tea.Cmd
cmds = append(cmds, r.activeView.Update(msg))
cmds = append(cmds, r.headerView.Update(msg))
return r, tea.Batch(cmds...)

// Depois (com Operation):
var cmds []tea.Cmd

// Rota para activeOperation ANTES da bifurcação modal/view.
// tea.Batch não garante ordem: mensagens privadas da operação (ex: fakeConfirmedMsg)
// podem chegar antes de CloseModalMsg no mesmo Batch. Se esperarmos o modal sair
// da pilha para rotear, a operação nunca receberia essas mensagens.
if r.activeOperation != nil {
    if cmd := r.activeOperation.Update(msg); cmd != nil {
        cmds = append(cmds, cmd)
    }
}

if len(r.modals) > 0 {
    top := len(r.modals) - 1
    cmds = append(cmds, r.modals[top].Update(msg))
} else {
    cmds = append(cmds, r.activeView.Update(msg))
    cmds = append(cmds, r.headerView.Update(msg))
}
return r, tea.Batch(cmds...)
```

Mensagens tratadas no `switch` com `return r, nil` (como `OperationCompletedMsg`) nunca chegam ao roteamento genérico — o `return` antecipado as consome. Isso é intencional: a operação não precisa processar sua própria mensagem de conclusão.

`ModalReadyMsg` atualmente vai para `activeView`. Para FakeOperation (e qualquer operação que usa action closures para resultados de modal) isso não importa — a operação não usa `ModalReadyMsg`. Operações futuras que precisem desse mecanismo devem incluir roteamento adicional de `ModalReadyMsg` para `activeOperation`.

Root não contém nenhuma lógica específica de operação além de tratar as mensagens terminais que já lhe pertencem.

## Fluxo de Mensagens

```
Ação pressionada → OnExecute (em setup.go) retorna StartOperation(op)
    root recebe StartOperationMsg → r.activeOperation = op, retorna op.Init()
    op.Init() retorna OpenModal(ConfirmModal("Deseja executar?", [Sim, Não]))
        root recebe OpenModalMsg → empilha modal (genérico)
        usuário pressiona "Não" → Action closure retorna Batch(CloseModal(), OperationCompleted())
            root recebe CloseModalMsg → desempilha modal
            root recebe OperationCompletedMsg → r.activeOperation = nil
        usuário pressiona "Sim" → Action closure retorna Batch(CloseModal(), fakeConfirmedMsg{})
            root recebe CloseModalMsg → desempilha modal
            root recebe fakeConfirmedMsg → roteia para op.Update(msg)  ← ROOT NÃO CONHECE ESSE TIPO
                op.Update(fakeConfirmedMsg) → state=executing, SetBusy, retorna fakeWorkCmd()
                fakeWorkCmd() roda em goroutine por 5s → retorna fakeWorkDoneMsg{}
                root recebe fakeWorkDoneMsg → roteia para op.Update(msg)
                    op.Update(fakeWorkDoneMsg) → Clear, retorna OpenModal(ConfirmModal("Executado", [OK]))
                        root recebe OpenModalMsg → empilha modal
                        usuário pressiona "OK" → Action closure retorna Batch(CloseModal(), OperationCompleted())
                            root recebe CloseModalMsg → desempilha modal
                            root recebe OperationCompletedMsg → r.activeOperation = nil
```

## Estrutura de Pacotes

**Novos arquivos:**

| Arquivo | Pacote | Conteúdo |
|---|---|---|
| `internal/tui/operation.go` | `tui` | Interface `Operation`, mensagens, helpers |
| `internal/tui/operation/fake_operation.go` | `operation` | `FakeOperation` para validação do padrão |

> Nota: `internal/tui/operation.go` é um **arquivo** dentro do pacote `tui`, não o subpacote. O subpacote `tui/operation` (diretório) contém as implementações concretas.

**Arquivos modificados:**

| Arquivo | Mudança |
|---|---|
| `internal/tui/root.go` | Campo `activeOperation`, novos cases no Update, roteamento genérico |
| `cmd/abditum/setup.go` | Nova ação que dispara `FakeOperation` |

## FakeOperation (implementação de validação)

Propósito: validar o padrão Operation de ponta a ponta sem lógica de negócio real.

**Estados:** `awaitingConfirmation` → `executing`

**Modais utilizados:** apenas `modal.ConfirmModal` (já implementado). Nenhum modal novo necessário.

**Construtor:**

```go
// NewFakeOperation cria uma FakeOperation.
// notifier é usado para reportar progresso na barra de mensagem.
func NewFakeOperation(notifier tui.MessageController) *FakeOperation {
    return &FakeOperation{notifier: notifier}
}
```

**Como disparar no setup.go:**

```go
{
    Keys:        []design.Key{design.Keys.F2},
    Label:       "Operação Fake",
    Description: "Demonstração do padrão Operation — confirmação + trabalho assíncrono.",
    GroupID:     "app",
    Priority:    99,
    Visible:     true,
    OnExecute: func() tea.Cmd {
        return tui.StartOperation(operation.NewFakeOperation(r.MessageController()))
    },
},
```

`design.Keys.F2` foi escolhido por não estar em uso no projeto.

**Fluxo:**

1. `Init()` abre `ConfirmModal("Deseja executar?", [Sim, Não])`
2. "Não": `CloseModal()` + `OperationCompleted()` — encerra sem passar por `Update`
3. "Sim": `CloseModal()` + `fakeConfirmedMsg{}` — operação recebe, avança para `executing`
4. `Update(fakeConfirmedMsg)`: `SetBusy("Executando...")`, lança `fakeWorkCmd()` (dorme 5s)
5. `Update(fakeWorkDoneMsg)`: `Clear()`, abre `ConfirmModal("Executado", [OK])`
6. "OK": `CloseModal()` + `OperationCompleted()`

**Mensagens privadas** (tipos não exportados, invisíveis ao root):
- `fakeConfirmedMsg{}`
- `fakeWorkDoneMsg{}`

## Operações Futuras

As três operações reais serão implementadas em fases subsequentes:

| Operação | Estados | Modais necessários |
|---|---|---|
| `OpenVaultOperation` | `awaitingPassword` → `decryptingVault` | `PasswordModal` (a criar) |
| `CreateVaultOperation` | `awaitingName` → `awaitingPassword` → `awaitingPasswordConfirm` → `creatingVault` | `TextInputModal` + `PasswordModal` (a criar) |
| `ExportSecretOperation` | `awaitingExportPath` → `exporting` | `TextInputModal` (a criar) |

`CreateVaultOperation` ao concluir encadeia `OpenVaultOperation` via `StartOperationMsg`.

## Testes

- `FakeOperation` deve ser testável em isolamento: instanciar, chamar `Init()`, simular mensagens, verificar comandos retornados — sem depender do `RootModel`.
- Testes de integração no `RootModel` verificam que `StartOperationMsg` seta `activeOperation` corretamente e que `OperationCompletedMsg` a limpa.
