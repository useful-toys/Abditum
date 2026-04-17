# Operation Pattern Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implementar o padrão Operation — uma mini-máquina de estados autônoma que encapsula operações multi-etapa no Bubble Tea, validado com uma FakeOperation de ponta a ponta.

**Architecture:** A interface `Operation` vive no pacote `tui` (arquivo `internal/tui/operation.go`). O `RootModel` ganha o campo `activeOperation` e roteia todas as mensagens para ela antes de bifurcar para modal/view. Implementações concretas ficam em `internal/tui/operation/` (subpacote), seguindo o mesmo padrão que `tui/modal` já usa em relação a `tui`.

**Tech Stack:** Go 1.26, Bubble Tea v2 (`charm.land/bubbletea/v2`), pacotes internos `tui`, `tui/modal`, `tui/design`.

---

## Arquivos

| Arquivo | Ação | Conteúdo |
|---|---|---|
| `internal/tui/operation.go` | **Criar** | Interface `Operation`, mensagens, funções helper |
| `internal/tui/root.go` | **Modificar** | Campo `activeOperation`, 4 novos cases, rewrite do bloco final de roteamento |
| `internal/tui/root_update_test.go` | **Criar** | Testes de integração do RootModel com Operation |
| `internal/tui/operation/fake_operation.go` | **Criar** | `FakeOperation` — implementação de validação do padrão |
| `internal/tui/operation/fake_operation_test.go` | **Criar** | Testes white-box da FakeOperation |
| `cmd/abditum/setup.go` | **Modificar** | Import do subpacote `operation`, ação F2 |

---

## Task 1: Definir a interface Operation

**Files:**
- Create: `internal/tui/operation.go`

- [ ] **Step 1: Criar o arquivo com interface e mensagens**

Crie `internal/tui/operation.go` com o seguinte conteúdo:

```go
package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/vault"
)

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
func StartOperation(op Operation) tea.Cmd {
	return func() tea.Msg { return StartOperationMsg{Op: op} }
}

// OperationCompleted cria um Cmd que emite OperationCompletedMsg.
func OperationCompleted() tea.Cmd {
	return func() tea.Msg { return OperationCompletedMsg{} }
}
```

- [ ] **Step 2: Verificar compilação**

```
go build ./internal/tui/...
```

Resultado esperado: sem erros.

- [ ] **Step 3: Commit**

```
git add internal/tui/operation.go
git commit -m "feat(tui): define Operation interface and messages

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

---

## Task 2: Atualizar RootModel (TDD)

**Files:**
- Create: `internal/tui/root_update_test.go`
- Modify: `internal/tui/root.go`

### 2a — Escrever os testes primeiro

- [ ] **Step 1: Criar `internal/tui/root_update_test.go` com os testes**

```go
package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// stubOperation é uma Operation mínima para uso nos testes do RootModel.
type stubOperation struct {
	initCalled    bool
	updateCalled  bool
	lastMsg       tea.Msg
	updateReturns tea.Cmd
}

func (s *stubOperation) Init() tea.Cmd {
	s.initCalled = true
	return nil
}

func (s *stubOperation) Update(msg tea.Msg) tea.Cmd {
	s.updateCalled = true
	s.lastMsg = msg
	return s.updateReturns
}

func TestRootModel_StartOperationMsg_SetsActiveOperationAndCallsInit(t *testing.T) {
	r := NewRootModel()
	op := &stubOperation{}

	_, _ = r.Update(StartOperationMsg{Op: op})

	if r.activeOperation != op {
		t.Error("StartOperationMsg: activeOperation não foi setado")
	}
	if !op.initCalled {
		t.Error("StartOperationMsg: Init() não foi chamado")
	}
}

func TestRootModel_StartOperationMsg_ReplacesExistingOperation(t *testing.T) {
	r := NewRootModel()
	op1 := &stubOperation{}
	op2 := &stubOperation{}

	_, _ = r.Update(StartOperationMsg{Op: op1})
	_, _ = r.Update(StartOperationMsg{Op: op2})

	if r.activeOperation != op2 {
		t.Error("StartOperationMsg: deve substituir a operação anterior sem encerrar")
	}
}

func TestRootModel_OperationCompletedMsg_ClearsActiveOperation(t *testing.T) {
	r := NewRootModel()
	op := &stubOperation{}
	r.activeOperation = op

	_, _ = r.Update(OperationCompletedMsg{})

	if r.activeOperation != nil {
		t.Error("OperationCompletedMsg: activeOperation deveria ser nil")
	}
}

type unknownMsg struct{}

func TestRootModel_UnknownMsg_RoutedToActiveOperation(t *testing.T) {
	r := NewRootModel()
	op := &stubOperation{}
	r.activeOperation = op

	_, _ = r.Update(unknownMsg{})

	if !op.updateCalled {
		t.Error("mensagem desconhecida: op.Update não foi chamado")
	}
	if _, ok := op.lastMsg.(unknownMsg); !ok {
		t.Errorf("mensagem desconhecida: op.Update recebeu %T, esperado unknownMsg", op.lastMsg)
	}
}

func TestRootModel_UnknownMsg_RoutedToOperation_EvenWhenModalActive(t *testing.T) {
	r := NewRootModel()
	op := &stubOperation{}
	r.activeOperation = op
	r.modals = append(r.modals, &stubModal{})

	_, _ = r.Update(unknownMsg{})

	if !op.updateCalled {
		t.Errorf("mensagem desconhecida com modal ativo: op.Update não foi chamado")
	}
}

// stubModal implementa ModalView para os testes do RootModel.
type stubModal struct{}

func (s *stubModal) Render(_ int, _ int, _ *design.Theme) string { return "" }
func (s *stubModal) HandleKey(_ tea.KeyMsg) tea.Cmd              { return nil }
func (s *stubModal) Update(_ tea.Msg) tea.Cmd                    { return nil }
```

- [ ] **Step 2: Rodar os testes — verificar que FALHAM**

```
go test ./internal/tui/ -run "TestRootModel_" -v
```

Resultado esperado: erros de compilação (campos e cases ainda não existem) ou FAIL.

### 2b — Implementar as mudanças no RootModel

- [ ] **Step 3: Adicionar campo `activeOperation` em `internal/tui/root.go`**

Localize o bloco de campos após `modals []ModalView` (linha ~54) e adicione o campo após ele:

```go
	// modals é a pilha de modais abertos; o topo da pilha é o modal ativo.
	modals []ModalView

	// activeOperation é a operação em andamento, ou nil quando ocioso.
	// Recebe todas as mensagens não tratadas pelo switch do root.
	activeOperation Operation
```

- [ ] **Step 4: Adicionar os 4 novos cases no switch do Update**

Localize o case `screen.WorkAreaChangedMsg` e adicione os 4 novos cases logo após ele (antes de `case ModalReadyMsg`):

```go
	case StartOperationMsg:
		r.activeOperation = msg.Op
		return r, msg.Op.Init()

	case OperationCompletedMsg:
		r.activeOperation = nil
		return r, nil

	case VaultOpenedMsg:
		r.setVaultManager(msg.Manager)
		r.setWorkArea(design.WorkAreaVault)
		return r, nil

	case SecretExportedMsg:
		return r, nil
```

- [ ] **Step 5: Reescrever o bloco final de roteamento**

Localize o bloco final do `Update` — as 8 linhas após o `switch` que terminam o método (linhas ~335–343 no arquivo original):

```go
	if len(r.modals) > 0 {
		top := len(r.modals) - 1
		return r, r.modals[top].Update(msg)
	}

	var cmds []tea.Cmd
	cmds = append(cmds, r.activeView.Update(msg))
	cmds = append(cmds, r.headerView.Update(msg))
	return r, tea.Batch(cmds...)
```

Substitua por:

```go
	var cmds []tea.Cmd

	// Roteia para activeOperation ANTES da bifurcação modal/view.
	// tea.Batch não garante ordem: mensagens privadas da operação (ex: fakeConfirmedMsg)
	// podem chegar antes de CloseModalMsg no mesmo Batch. Se aguardarmos o modal sair
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

- [ ] **Step 6: Rodar os testes — verificar que PASSAM**

```
go test ./internal/tui/ -run "TestRootModel_" -v
```

Resultado esperado: 5 testes PASS.

- [ ] **Step 7: Rodar a suite completa do pacote tui para garantir ausência de regressões**

```
go test ./internal/tui/... -v
```

Resultado esperado: todos os testes passam.

- [ ] **Step 8: Commit**

```
git add internal/tui/root.go internal/tui/root_update_test.go
git commit -m "feat(tui): add Operation support to RootModel

- Add activeOperation field to RootModel
- Handle StartOperationMsg, OperationCompletedMsg, VaultOpenedMsg, SecretExportedMsg
- Rewrite generic routing block: always route to activeOperation before modal/view
- Add integration tests for Operation routing

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

---

## Task 3: Implementar FakeOperation (TDD)

**Files:**
- Create: `internal/tui/operation/fake_operation_test.go`
- Create: `internal/tui/operation/fake_operation.go`

### 3a — Escrever os testes primeiro

- [ ] **Step 1: Criar o diretório**

```
mkdir internal\tui\operation
```

- [ ] **Step 2: Criar `internal/tui/operation/fake_operation_test.go`**

Os testes são white-box (package `operation`) para acessar os tipos privados `fakeConfirmedMsg` e `fakeWorkDoneMsg`.

```go
package operation

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
)

// stubNotifier implementa tui.MessageController para os testes da FakeOperation.
// Registra qual método foi chamado por último e com qual texto.
type stubNotifier struct {
	lastMethod string
	lastText   string
}

func (s *stubNotifier) SetBusy(text string)     { s.lastMethod = "SetBusy"; s.lastText = text }
func (s *stubNotifier) SetSuccess(text string)  { s.lastMethod = "SetSuccess"; s.lastText = text }
func (s *stubNotifier) SetError(text string)    { s.lastMethod = "SetError"; s.lastText = text }
func (s *stubNotifier) SetWarning(text string)  { s.lastMethod = "SetWarning"; s.lastText = text }
func (s *stubNotifier) SetInfo(text string)     { s.lastMethod = "SetInfo"; s.lastText = text }
func (s *stubNotifier) SetHintField(text string){ s.lastMethod = "SetHintField"; s.lastText = text }
func (s *stubNotifier) SetHintUsage(text string){ s.lastMethod = "SetHintUsage"; s.lastText = text }
func (s *stubNotifier) Clear()                  { s.lastMethod = "Clear"; s.lastText = "" }

// execCmd executa um tea.Cmd e retorna a mensagem produzida.
// Retorna nil se cmd for nil.
func execCmd(cmd tea.Cmd) tea.Msg {
	if cmd == nil {
		return nil
	}
	return cmd()
}

func TestFakeOperation_Init_EmitsOpenModalMsg(t *testing.T) {
	op := NewFakeOperation(&stubNotifier{})

	cmd := op.Init()

	msg := execCmd(cmd)
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Init: esperado OpenModalMsg, obteve %T", msg)
	}
}

func TestFakeOperation_Update_IgnoresUnknownMsg(t *testing.T) {
	op := NewFakeOperation(&stubNotifier{})

	type randomMsg struct{}
	cmd := op.Update(randomMsg{})

	if cmd != nil {
		t.Error("Update(randomMsg): deveria retornar nil para mensagem desconhecida")
	}
}

func TestFakeOperation_Update_Confirmed_SetsBusyAndStartsWork(t *testing.T) {
	n := &stubNotifier{}
	op := NewFakeOperation(n)

	cmd := op.Update(fakeConfirmedMsg{})

	if op.state != stateExecuting {
		t.Errorf("Update(fakeConfirmedMsg): state esperado stateExecuting, obteve %v", op.state)
	}
	if n.lastMethod != "SetBusy" {
		t.Errorf("Update(fakeConfirmedMsg): esperado SetBusy, notifier recebeu %q", n.lastMethod)
	}
	if cmd == nil {
		t.Error("Update(fakeConfirmedMsg): esperado cmd não-nil (fakeWorkCmd)")
	}
}

func TestFakeOperation_Update_WorkDone_ClearsNotifierAndOpensResultModal(t *testing.T) {
	n := &stubNotifier{}
	op := NewFakeOperation(n)
	op.state = stateExecuting // avança estado manualmente para o correto

	cmd := op.Update(fakeWorkDoneMsg{})

	if n.lastMethod != "Clear" {
		t.Errorf("Update(fakeWorkDoneMsg): esperado Clear no notifier, obteve %q", n.lastMethod)
	}
	msg := execCmd(cmd)
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Update(fakeWorkDoneMsg): esperado OpenModalMsg, obteve %T", msg)
	}
}

func TestFakeOperation_Update_ConfirmedInWrongState_Ignored(t *testing.T) {
	op := NewFakeOperation(&stubNotifier{})
	op.state = stateExecuting // já está executando

	cmd := op.Update(fakeConfirmedMsg{})

	if cmd != nil {
		t.Error("Update(fakeConfirmedMsg) em stateExecuting: deveria ser ignorado (retornar nil)")
	}
}
```

- [ ] **Step 3: Rodar os testes — verificar que FALHAM (não compilam)**

```
go test ./internal/tui/operation/... -v
```

Resultado esperado: erros de compilação — `NewFakeOperation`, `fakeConfirmedMsg`, `fakeWorkDoneMsg`, `stateExecuting` não existem ainda.

### 3b — Implementar a FakeOperation

- [ ] **Step 4: Criar `internal/tui/operation/fake_operation.go`**

```go
package operation

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
)

// fakeState representa o estado interno da FakeOperation.
type fakeState int

const (
	// stateAwaitingConfirmation é o estado inicial — aguarda resposta do modal de confirmação.
	stateAwaitingConfirmation fakeState = iota
	// stateExecuting indica que o trabalho assíncrono está em andamento.
	stateExecuting
)

// fakeConfirmedMsg é emitida pela action closure do modal de confirmação quando o
// usuário confirma. Tipo não exportado — invisível ao RootModel.
type fakeConfirmedMsg struct{}

// fakeWorkDoneMsg é emitida pela goroutine de trabalho fake após 5 segundos.
// Tipo não exportado — invisível ao RootModel.
type fakeWorkDoneMsg struct{}

// FakeOperation valida o padrão Operation de ponta a ponta sem lógica de negócio real.
// Fluxo: modal de confirmação → 5s de trabalho fake → modal de resultado.
type FakeOperation struct {
	state    fakeState
	notifier tui.MessageController
}

// NewFakeOperation cria uma FakeOperation.
// notifier é usado para reportar progresso na barra de mensagem durante o trabalho fake.
func NewFakeOperation(notifier tui.MessageController) *FakeOperation {
	return &FakeOperation{notifier: notifier}
}

// Init abre o modal de confirmação. Implementa tui.Operation.
func (f *FakeOperation) Init() tea.Cmd {
	return tui.OpenModal(f.buildConfirmModal())
}

// Update processa mensagens e avança a máquina de estados. Implementa tui.Operation.
func (f *FakeOperation) Update(msg tea.Msg) tea.Cmd {
	switch msg.(type) {
	case fakeConfirmedMsg:
		if f.state != stateAwaitingConfirmation {
			return nil
		}
		f.state = stateExecuting
		f.notifier.SetBusy("Executando operação fake...")
		return fakeWorkCmd()

	case fakeWorkDoneMsg:
		f.notifier.Clear()
		return tui.OpenModal(f.buildResultModal())
	}
	return nil
}

// buildConfirmModal cria o modal de confirmação inicial.
func (f *FakeOperation) buildConfirmModal() *modal.ConfirmModal {
	return modal.NewConfirmModal(
		"Operação Fake",
		"Deseja executar a operação fake?\nIsso simulará 5 segundos de trabalho.",
		[]modal.ModalOption{
			{
				Keys:   []design.Key{design.Keys.Enter},
				Label:  "Executar",
				Intent: modal.IntentConfirm,
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg { return fakeConfirmedMsg{} })
				},
			},
			{
				Keys:   []design.Key{design.Keys.Esc},
				Label:  "Cancelar",
				Intent: modal.IntentCancel,
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), tui.OperationCompleted())
				},
			},
		},
	)
}

// buildResultModal cria o modal de resultado após o trabalho fake.
func (f *FakeOperation) buildResultModal() *modal.ConfirmModal {
	return modal.NewConfirmModal(
		"Operação Fake",
		"Operação concluída com sucesso!",
		[]modal.ModalOption{
			{
				Keys:   []design.Key{design.Keys.Enter},
				Label:  "OK",
				Intent: modal.IntentConfirm,
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), tui.OperationCompleted())
				},
			},
		},
	)
}

// fakeWorkCmd retorna um Cmd que aguarda 5 segundos em goroutine e emite fakeWorkDoneMsg.
func fakeWorkCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(5 * time.Second)
		return fakeWorkDoneMsg{}
	}
}
```

- [ ] **Step 5: Rodar os testes — verificar que PASSAM**

```
go test ./internal/tui/operation/... -v
```

Resultado esperado: 5 testes PASS.

- [ ] **Step 6: Rodar suite completa para verificar ausência de regressões**

```
go test ./internal/tui/... -v
```

Resultado esperado: todos os testes passam.

- [ ] **Step 7: Commit**

```
git add internal/tui/operation/fake_operation.go internal/tui/operation/fake_operation_test.go
git commit -m "feat(operation): implement FakeOperation

State machine: awaitingConfirmation -> executing
- Init opens confirm modal
- fakeConfirmedMsg: sets busy, starts 5s fake work goroutine
- fakeWorkDoneMsg: clears notifier, opens result modal
- Cancel and OK use CloseModal + OperationCompleted in action closures

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

---

## Task 4: Conectar FakeOperation no setup.go

**Files:**
- Modify: `cmd/abditum/setup.go`

- [ ] **Step 1: Adicionar o import do subpacote `operation` em `cmd/abditum/setup.go`**

Localize o bloco de imports e adicione a linha do subpacote:

```go
import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/actions"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/tui/operation"
)
```

- [ ] **Step 2: Adicionar a ação F2 na função `setupApplication`**

Dentro de `setupApplication`, no slice passado para `r.RegisterApplicationActions`, adicione a ação ao final (após a ação Quit):

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

- [ ] **Step 3: Verificar compilação**

```
go build ./...
```

Resultado esperado: sem erros.

- [ ] **Step 4: Rodar suite completa**

```
go test ./... -count=1
```

Resultado esperado: todos os testes passam.

- [ ] **Step 5: Commit**

```
git add cmd/abditum/setup.go
git commit -m "feat(setup): wire FakeOperation to F2 key

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```
