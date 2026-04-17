# Design — Saída com Confirmação (Fluxo 3)

## Problema

A action "Sair" (`⌃Q`) hoje chama `tea.Quit` diretamente, sem nenhuma confirmação. O `fluxos.md` (Fluxo 3) especifica que ao sair sem cofre carregado o sistema deve solicitar confirmação antes de encerrar.

## Escopo

Implementa **somente o Fluxo 3**: saída quando **nenhum cofre está carregado**.

Fluxos 4 (cofre inalterado) e 5 (cofre alterado com save) estão fora de escopo desta fase.

## Abordagem

Seguir o padrão `Operation` já estabelecido no projeto (`FakeOperation` como referência).

## Arquitetura

### Novo arquivo: `internal/tui/operation/exit_operation.go`

`ExitOperation` implementa `tui.Operation`. É uma máquina de estados com um único estado:

```
stateAwaitingConfirmation  (estado inicial e único)
```

**`Init()`** — abre um `modal.ConfirmModal` com:
- Título: `"Sair"`
- Mensagem: `"Deseja encerrar a aplicação?"`
- Opções:
  - `Enter` / `"Confirmar"` / `IntentConfirm` → emite `exitConfirmedMsg` (tipo privado) + fecha modal
  - `Esc` / `"Voltar"` / `IntentCancel` → fecha modal + `tui.OperationCompleted()`

**`Update(msg)`** — responde apenas a `exitConfirmedMsg`:
- Retorna `tea.Quit`

Sem goroutines, sem IO, sem estados adicionais.

### Modificação: `cmd/abditum/setup.go`

A action `Sair` passa de:
```go
OnExecute: func() tea.Cmd { return tea.Quit },
```
para:
```go
OnExecute: func() tea.Cmd {
    return tui.StartOperation(operation.NewExitOperation())
},
```

`AvailableWhen` permanece `nil` (sempre disponível) — filtragem por estado do cofre é responsabilidade de fases futuras.

## Testes

`exit_operation_test.go` cobre:

1. `Init()` retorna um `OpenModalMsg` (não nil, não `tea.Quit`)
2. `Update(exitConfirmedMsg{})` retorna `tea.Quit`
3. `Update(outros msgs)` retorna nil (ignora mensagens desconhecidas)
4. Opção "Voltar": a action do modal emite `CloseModal + OperationCompleted` (não `tea.Quit`)

## Arquivos afetados

| Arquivo | Operação |
|---------|----------|
| `internal/tui/operation/exit_operation.go` | Criar |
| `internal/tui/operation/exit_operation_test.go` | Criar |
| `cmd/abditum/setup.go` | Modificar (1 linha) |

Nenhum outro arquivo existente é modificado.
