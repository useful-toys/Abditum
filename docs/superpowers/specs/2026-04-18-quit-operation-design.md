# Design — QuitOperation (ctrl Q)

**Data:** 2026-04-18  
**Contexto:** `fluxos.md` — Fluxos 3, 4 e 5

---

## Escopo

Substituir o `tea.Quit` direto da action `ctrl Q` por uma `QuitOperation` que implementa os três fluxos de saída do `fluxos.md`:

- **Fluxo 3** — sair sem cofre carregado: confirmação simples.
- **Fluxo 4** — sair com cofre inalterado: confirmação simples.
- **Fluxo 5** — sair com cofre alterado: salvar e sair / descartar e sair / voltar; com detecção de modificação externa antes do salvamento.

---

## Componentes e mudanças

### `internal/vault` — interface e Manager

**`RepositorioCofre`** ganha um terceiro método:

```go
DetectarAlteracaoExterna() (bool, error)
```

Responsabilidade: verificar se o arquivo em disco foi modificado por processo externo desde o último Load ou Save.

**`Manager.Salvar`** muda de assinatura:

```go
func (m *Manager) Salvar(forcarSobrescrita bool) error
```

Comportamento:
- Se `forcarSobrescrita == false`: chama `repo.DetectarAlteracaoExterna()` antes de salvar. Se detectar mudança externa, retorna `vault.ErrModifiedExternally` sem salvar.
- Se `forcarSobrescrita == true`: pula a verificação e salva diretamente via protocolo atômico.

**Novo erro sentinela** declarado em `internal/vault`:

```go
var ErrModifiedExternally = errors.New("arquivo modificado externamente")
```

Verificável via `errors.Is`.

### `internal/storage` — FileRepository

**`FileRepository`** implementa `DetectarAlteracaoExterna()` delegando para a função já existente:

```go
func (r *FileRepository) DetectarAlteracaoExterna() (bool, error) {
    return storage.DetectExternalChange(r.path, r.metadata)
}
```

`DetectExternalChange` usa tamanho + SHA-256 (mtime deliberadamente ignorado — ver `arquitetura-storage.md §Detecção de Modificação Externa`).

### `internal/tui/operation` — QuitOperation

Nova struct `QuitOperation` que implementa `tui.Operation`.

**Construtor:**

```go
func NewQuitOperation(notifier tui.MessageController, manager *vault.Manager) *QuitOperation
```

`manager` pode ser `nil` (nenhum cofre carregado).

**`Init()`** — decide qual modal abrir:

- `manager == nil` ou `!manager.IsModified()` → modal de confirmação simples (Fluxos 3 e 4).
- `manager.IsModified()` → modal de decisão para cofre alterado (Fluxo 5).

**`Update(msg)`** — processa mensagens internas (`quitMsg`):

- `quitStateSaving` → `notifier.SetBusy("Salvando...")` + chama `manager.Salvar(false)` em goroutine:
  - `ErrModifiedExternally` → `notifier.Clear()` + abre modal de conflito externo.
  - Outro erro → `notifier.SetError(...)` + `OperationCompleted()`.
  - Sucesso → `notifier.SetSuccess("Cofre salvo.")` + `tea.Quit`.
- `quitStateSavingForced` → `notifier.SetBusy("Salvando...")` + chama `manager.Salvar(true)` em goroutine:
  - Erro → `notifier.SetError(...)` + `OperationCompleted()`.
  - Sucesso → `notifier.SetSuccess("Cofre salvo.")` + `tea.Quit`.

**Modais:**

| Modal | Opções |
|---|---|
| Confirmação simples | Confirmar (`Enter`) → `CloseModal + tea.Quit` · Voltar (`Esc`) → `CloseModal + OperationCompleted` |
| Cofre alterado | Salvar e sair (`Enter`) → `CloseModal + quitMsg{saving}` · Descartar e sair (`d`) → `CloseModal + tea.Quit` · Voltar (`Esc`) → `CloseModal + OperationCompleted` |
| Conflito externo | Sobrescrever e sair (`Enter`) → `CloseModal + quitMsg{savingForced}` · Voltar (`Esc`) → `CloseModal + OperationCompleted` |

### `cmd/abditum/setup.go`

A action `Quit` passa de:

```go
OnExecute: func() tea.Cmd { return tea.Quit }
```

Para:

```go
OnExecute: func() tea.Cmd {
    return tui.StartOperation(operation.NewQuitOperation(r.MessageController(), r.Manager()))
}
```

---

## Fluxo de dados

```
ctrl Q
  → Action.OnExecute
  → StartOperation(QuitOperation)
  → QuitOperation.Init()
      → [sem cofre ou inalterado] → ConfirmModal
      → [cofre alterado]         → ModifiedModal
  → usuário escolhe
      → [confirmar / descartar]  → tea.Quit
      → [salvar e sair]          → quitMsg{saving}
          → Manager.Salvar(false)
              → [OK]                    → tea.Quit
              → [ErrModifiedExternally] → ConflictModal
                  → [sobrescrever]      → quitMsg{savingForced}
                      → Manager.Salvar(true) → tea.Quit / SetError
                  → [voltar]            → OperationCompleted
              → [outro erro]            → SetError + OperationCompleted
      → [voltar]                 → OperationCompleted
```

## Sinalização de progresso

Durante qualquer operação de salvamento, o sistema sinaliza o estado via `MessageController`:

| Momento | Chamada |
|---|---|
| Início do salvamento | `notifier.SetBusy("Salvando...")` |
| Salvamento bem-sucedido | `notifier.SetSuccess("Cofre salvo.")` |
| Falha no salvamento | `notifier.SetError(mensagem do erro)` |
| Conflito externo detectado | `notifier.Clear()` — o modal de conflito substitui a sinalização |

O `SetBusy` é emitido antes de disparar a goroutine de salvamento, garantindo que o indicador apareça imediatamente. O `tea.Quit` é emitido apenas após o `SetSuccess`, permitindo que a mensagem seja renderizada antes do encerramento.

---



| Situação | Comportamento |
|---|---|
| `Salvar(false)` retorna `ErrModifiedExternally` | Abre modal de conflito; usuário decide sobrescrever ou voltar |
| `Salvar(false)` retorna outro erro | `notifier.SetError(mensagem)` + `OperationCompleted()` — cofre permanece carregado e alterado |
| `Salvar(true)` retorna erro | `notifier.SetError(mensagem)` + `OperationCompleted()` — cofre permanece carregado e alterado |
| `DetectarAlteracaoExterna` retorna erro (ex: arquivo inacessível) | Propagado como erro de `Salvar(false)` — tratado como erro genérico |
| Cofre recém-criado (`isNew == true`, metadata zerada) | `DetectarAlteracaoExterna` retorna `false` por definição — arquivo não existia antes, não há baseline para comparar. `Salvar(false)` prossegue normalmente. |

---

## O que não muda

- `RepositorioCofre` continua com `Salvar` e `Carregar` — apenas ganha `DetectarAlteracaoExterna`.
- Nenhuma outra `Operation` é alterada.
- O padrão `Action → StartOperation → Operation` já estabelecido pela `FakeOperation` é seguido sem modificação.
- Testes existentes de `Manager` que passam `nil` ou `&mockRepository{}` precisarão implementar `DetectarAlteracaoExterna` no mock — mudança trivial de uma linha.
