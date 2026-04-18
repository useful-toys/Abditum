# Design — Fluxo 1 (Abrir Cofre) e Fluxo 2 (Criar Cofre)

**Data:** 2026-04-18  
**Referência:** `golden/fluxos.md` — Fluxo 1 e Fluxo 2  
**Estado do código base:** `QuitOperation` já implementada (Fluxos 3, 4, 5); `Manager.Salvar(forcarSobrescrita bool)` já existe.

---

## Contexto

Os Fluxos 1 (Abrir Cofre Existente) e 2 (Criar Novo Cofre) compartilham um passo inicial idêntico: verificar se há um cofre carregado com modificações não salvas antes de prosseguir. A `QuitOperation` já implementou lógica similar para o Fluxo 5, estabelecendo padrões reutilizáveis.

---

## Decisões de Design

### 1. Reuso via `guardCofreAlterado`

O passo 1 dos Fluxos 1 e 2 ("verificar cofre alterado antes de prosseguir") é extraído como um helper interno ao pacote `operation/`. Não é uma `Operation` pública — é um tipo auxiliar com ciclo `Init/Update` próprio.

**Diferença em relação ao Fluxo 5 (QuitOperation):** ambos chamam `Salvar(false)` que detecta modificação externa. A diferença é a **resposta ao conflito**: o passo 1 dos Fluxos 1/2 trata `ErrModifiedExternally` como qualquer outra falha (interrompe o fluxo, sem oferecer sobrescrever). O Fluxo 5 abre um modal de conflito com opção de sobrescrever. Por isso a `QuitOperation` mantém sua lógica própria para o sub-fluxo pós-conflito.

### 2. Interface `vaultSaver` movida para arquivo compartilhado

A interface `vaultSaver` (já em `quit_operation.go`) é movida para `vault_saver.go` no mesmo pacote, tornando-a acessível a `guardCofreAlterado`, `CriarCofreOperation` e `AbrirCofreOperation` sem duplicação.

```go
// vaultSaver é a interface mínima que as operations precisam do vault.Manager.
type vaultSaver interface {
    IsModified() bool
    Salvar(forcarSobrescrita bool) error
}
```

### 3. Sinalização de progresso (SetBusy/Clear)

Toda operação de IO segue o padrão estabelecido na `QuitOperation`:
1. `SetBusy(mensagem)` imediatamente antes de despachar o goroutine
2. `Clear()` — ou `SetSuccess` / `SetError` — ao receber o resultado

### 4. Estrutura de arquivos

```
internal/tui/operation/
├── fake_operation.go          (existente — sem alteração)
├── quit_operation.go          (existente — remover vaultSaver daqui)
├── vault_saver.go             (novo — interface vaultSaver compartilhada)
├── guard_cofre_alterado.go    (novo — helper para passo 1)
├── criar_cofre.go             (novo — Fluxo 2)
└── abrir_cofre.go             (novo — Fluxo 1)
```

---

## Componentes

### `vault_saver.go`

Move a interface `vaultSaver` de `quit_operation.go` para cá. Sem outras mudanças.

---

### `guard_cofre_alterado.go`

Helper interno que encapsula o passo 1 dos Fluxos 1 e 2.

**Assinatura:**

```go
type guardCofreAlterado struct {
    saver     vaultSaver
    notifier  tui.MessageController
    onProceder func() tea.Cmd
    onAbortado func() tea.Cmd
}

func novoGuardCofreAlterado(
    saver vaultSaver,
    notifier tui.MessageController,
    onProceder func() tea.Cmd,
    onAbortado func() tea.Cmd,
) *guardCofreAlterado
```

**Máquina de estados:**

```
Init():
  saver == nil || !IsModified()  →  emite onProceder() diretamente
  IsModified()                   →  OpenModal(buildModifiedModal())

Update(guardSaveMsg):
  SetBusy("Salvando...")
  → Salvar(false) assíncrono
  → sucesso: Clear() → onProceder()
  → falha: SetError(err) → onAbortado()
```

**Modal de decisão** (título "Alterações não salvas"):
- `Enter` → Salvar e prosseguir
- `D` → Descartar e prosseguir → onProceder() direto
- `Esc` → Voltar → onAbortado()

**Nota sobre modificação externa:** o guard chama `Salvar(false)`, que já detecta modificação externa via `ErrModifiedExternally`. Diferente do Fluxo 5 (QuitOperation), o passo 1 dos Fluxos 1 e 2 **não** oferece opção de sobrescrever — qualquer falha de save (incluindo conflito externo) resulta em `SetError` + `onAbortado`. O fluxo é interrompido e o cofre permanece carregado e `alterado`, conforme o mermaid do Fluxo 1/2.

---

### `criar_cofre.go` — `CriarCofreOperation`

Implementa o Fluxo 2 completo.

**Construtor:**

```go
func NewCriarCofreOperation(
    notifier tui.MessageController,
    manager *vault.Manager, // nil se sem cofre carregado
) *CriarCofreOperation
```

**Estados internos:**

| Estado | Descrição |
|--------|-----------|
| `criandoGuardando` | Executa `guardCofreAlterado`; ao concluir → `criandoInformandoCaminho` |
| `criandoInformandoCaminho` | `FilePicker` modo Save, extensão `.abditum` |
| `criandoConfirmandoSobrescrita` | `ConfirmModal` "arquivo já existe" |
| `criandoInformandoSenha` | `PasswordCreateModal` (2 campos, force meter) |
| `criandoAvaliacaoSenhaFraca` | `ConfirmModal` "senha fraca, prosseguir?" |
| `criandoCriando` | Busy + IO assíncrono |

**Máquina de estados detalhada:**

```
Init():
  → instancia guardCofreAlterado
    onProceder = emitir msg para → criandoInformandoCaminho
    onAbortado = OperationCompleted()
  → guard.Init()

criandoInformandoCaminho:
  → OpenModal(FilePicker{Mode: Save, Ext: ".abditum"})
  → desistiu (Esc no picker): OperationCompleted()
  → caminho selecionado:
      arquivo existe  → criandoConfirmandoSobrescrita (guardar: destino=existente)
      não existe      → criandoInformandoSenha (guardar: destino=novo)

criandoConfirmandoSobrescrita:
  → ConfirmModal "O arquivo já existe. Deseja sobrescrever?"
  → "Outro caminho" (Esc): → criandoInformandoCaminho
  → "Sobrescrever" (Enter): → criandoInformandoSenha (destino=existente)

criandoInformandoSenha:
  → OpenModal(PasswordCreateModal)
    onConfirm(password, forte bool): guardar senha; forte → criandoCriando; fraco → criandoAvaliacaoSenhaFraca
    onCancel(): → criandoInformandoCaminho
  (senhas não coincidentes: o modal trata internamente, onConfirm não é chamado)

criandoAvaliacaoSenhaFraca:
  → ConfirmModal "Senha fraca. Prosseguir assim mesmo?"
  → "Revisar" (Esc): → criandoInformandoSenha
  → "Prosseguir" (Enter): → criandoCriando

criandoCriando:
  → SetBusy("Criando cofre...")
  → goroutine:
      cofre = NovoCofre()
      InicializarConteudoPadrao(cofre)
      repo = NewFileRepositoryForCreate(caminho, senha)
      manager = NewManager(cofre, repo)
      err = manager.Salvar(false)
      → retorna criarCofreResultMsg{manager, err}
  → sucesso:
      Clear()
      VaultOpenedMsg{Manager: manager}
      SetSuccess("Cofre criado.")
  → falha:
      SetError(err.Error())
      OperationCompleted()
      (cofre anterior preservado — root não troca o manager)
```

**Tratamento de erros na gravação (passo 5 do Fluxo 2):**
- Falha ao gravar novo arquivo → `SetError` + `OperationCompleted`
- Falha em sobrescrita sem backup → `SetError` + `OperationCompleted`
- Falha em sobrescrita após backup → `SetError("...backup disponível em ...")` + `OperationCompleted`

Os detalhes da mensagem de erro com backup dependem de `storage.ErrBackupDisponivel` ou similar — a ser verificado durante implementação.

---

### `abrir_cofre.go` — `AbrirCofreOperation`

Implementa o Fluxo 1 completo.

**Construtor:**

```go
func NewAbrirCofreOperation(
    notifier tui.MessageController,
    manager *vault.Manager, // nil se sem cofre carregado
) *AbrirCofreOperation
```

**Estados internos:**

| Estado | Descrição |
|--------|-----------|
| `abrindoGuardando` | Executa `guardCofreAlterado`; ao concluir → `abrindoInformandoCaminho` |
| `abrindoInformandoCaminho` | `FilePicker` modo Open |
| `abrindoInformandoSenha` | `PasswordEntryModal` (1 campo) |
| `abrindoAbrindo` | Busy + IO assíncrono |

**Máquina de estados detalhada:**

```
Init():
  → instancia guardCofreAlterado
    onProceder = emitir msg para → abrindoInformandoCaminho
    onAbortado = OperationCompleted()
  → guard.Init()

abrindoInformandoCaminho:
  → OpenModal(FilePicker{Mode: Open})
  → desistiu: OperationCompleted()
  → caminho selecionado:
      validar magic + versão_formato (IO síncrono rápido, sem busy)
      inválido → SetError(err) + volta a abrindoInformandoCaminho
      ok → abrindoInformandoSenha (guardar caminho)

abrindoInformandoSenha:
  → OpenModal(PasswordEntryModal)
    onConfirm(password): → abrindoAbrindo
    onCancel(): → abrindoInformandoCaminho

abrindoAbrindo:
  → SetBusy("Abrindo cofre...")
  → goroutine:
      repo = NewFileRepository(caminho, senha)
      cofre, err = repo.Carregar()
      → retorna abrirCofreResultMsg{manager, err}
  → senha errada / autenticação falhou (ErrAutenticacao):
      Clear()
      SetError("Senha incorreta ou arquivo corrompido.")  ← erro genérico (não revela causa)
      → abrindoInformandoSenha
  → payload corrompido / PastaGeral ausente (ErrIntegridade):
      Clear()
      SetError("Arquivo corrompido ou inválido.")
      → abrindoInformandoCaminho
  → sucesso:
      Clear()
      VaultOpenedMsg{Manager: manager}
      SetSuccess("Cofre aberto.")
```

**Nota sobre erros genéricos:** a spec do Fluxo 1 exige mensagens de erro por *categoria* (autenticação, integridade) sem revelar a causa exata — por segurança. A implementação usa essas duas categorias.

---

## Integração em `setup.go`

Dois novos atalhos serão registrados:

| Atalho | Label | Condição `AvailableWhen` |
|--------|-------|--------------------------|
| `Ctrl+N` | Criar cofre | sem cofre carregado (restrição de UX atual) |
| `Ctrl+O` | Abrir cofre | sem cofre carregado (restrição de UX atual) |

```go
// Criar cofre (Ctrl+N)
{
    Keys:          []design.Key{design.Shortcuts.NewVault},
    Label:         "Criar cofre",
    Description:   "Cria um novo cofre protegido por senha.",
    GroupID:       "app",
    Priority:      30,
    Visible:       true,
    AvailableWhen: func(app AppState, _ ChildView) bool { return app.Manager() == nil },
    OnExecute: func() tea.Cmd {
        return tui.StartOperation(operation.NewCriarCofreOperation(r.MessageController(), r.Manager()))
    },
},
// Abrir cofre (Ctrl+O)
{
    Keys:          []design.Key{design.Shortcuts.OpenVault},
    Label:         "Abrir cofre",
    Description:   "Abre um cofre existente a partir de um arquivo.",
    GroupID:       "app",
    Priority:      31,
    Visible:       true,
    AvailableWhen: func(app AppState, _ ChildView) bool { return app.Manager() == nil },
    OnExecute: func() tea.Cmd {
        return tui.StartOperation(operation.NewAbrirCofreOperation(r.MessageController(), r.Manager()))
    },
},
```

Os atalhos `Ctrl+N` e `Ctrl+O` precisam ser adicionados a `design/keys.go` como `Shortcuts.NewVault` e `Shortcuts.OpenVault`.

---

## Testes

Cada novo arquivo terá seu `_test.go` correspondente, seguindo o padrão de `quit_operation_test.go`:

- `vault_saver.go` — sem teste (só interface)
- `guard_cofre_alterado_test.go` — testa os 3 caminhos: sem cofre, cofre inalterado, cofre alterado (salvar ok, salvar falha, descartar)
- `criar_cofre_test.go` — testa estados principais: guard integrado, picker, sobrescrita, senha fraca, criação ok, criação falha
- `abrir_cofre_test.go` — testa: guard integrado, picker, magic inválido, senha errada, integridade, sucesso

Todos os testes usam stubs das interfaces (`vaultSaver`, `tui.MessageController`) — sem IO real.

---

## Pontos a verificar durante implementação

1. **`PasswordCreateModal.onConfirm`** — verificar se o callback já recebe indicação de força da senha ou se a avaliação precisa ser feita na operation após o callback.
2. **Erros de backup** — verificar se `storage` já expõe `ErrBackupDisponivel` ou similar para mensagens diferenciadas no `criandoCriando`.
3. **`design/keys.go`** — verificar se `Shortcuts.NewVault` e `Shortcuts.OpenVault` já existem ou precisam ser adicionados.
4. **`FilePicker` modo Open** — verificar se já valida magic/versão ou se a validação precisa ocorrer na operation após o picker retornar o caminho.
5. **`NewFileRepositoryForCreate`** — verificar assinatura exata em `internal/storage/repository.go`.
