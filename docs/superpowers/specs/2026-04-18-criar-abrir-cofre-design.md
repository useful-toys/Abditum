# Design — Fluxo 1 (Abrir Cofre) e Fluxo 2 (Criar Cofre)

**Data:** 2026-04-18  
**Referência:** `golden/fluxos.md` — Fluxo 1, Fluxo 2 e Entrada via Linha de Comando  
**Estado do código base:** `QuitOperation` já implementada (Fluxos 3, 4, 5); `Manager.Salvar(forcarSobrescrita bool)` já existe.

---

## Contexto

Os Fluxos 1 (Abrir Cofre Existente) e 2 (Criar Novo Cofre) compartilham um passo inicial idêntico: verificar se há um cofre carregado com modificações não salvas antes de prosseguir. A `QuitOperation` já implementou lógica idêntica para o Fluxo 5. O guard extrai essa lógica como helper reutilizável, e a `QuitOperation` é refatorada para usá-lo.

A aplicação também pode ser iniciada via CLI com `--vault <caminho>`, o que dispara automaticamente o fluxo correto (Abrir ou Criar) pulando o guard e o picker.

---

## Decisões de Design

### 1. `guardCofreAlterado` como helper reutilizável

O passo 1 dos Fluxos 1, 2 e 5 ("verificar cofre alterado antes de prosseguir") é extraído como helper interno ao pacote `operation/`. Não é uma `Operation` pública — é um tipo auxiliar com ciclo `Init/Update` próprio, parametrizado por `onProceder` e `onAbortado`.

### 2. Refatoração da `QuitOperation` para usar o guard

A `QuitOperation` é refatorada para delegar ao `guardCofreAlterado` quando o cofre está alterado (Fluxo 5). A continuação é `onProceder = tea.Quit`, `onAbortado = OperationCompleted()`. Os Fluxos 3 e 4 (sem alterações) continuam com modal de confirmação simples, sem passar pelo guard. O comportamento externo é idêntico — os testes existentes devem continuar passando.

### 3. Interface `vaultSaver` movida para arquivo compartilhado

A interface `vaultSaver` (já em `quit_operation.go`) é movida para `vault_saver.go` no mesmo pacote, tornando-a acessível a `guardCofreAlterado`, `QuitOperation`, `CriarCofreOperation` e `AbrirCofreOperation` sem duplicação.

```go
// vaultSaver é a interface mínima que as operations precisam do vault.Manager.
type vaultSaver interface {
    IsModified() bool
    Salvar(forcarSobrescrita bool) error
}
```

### 4. `storage.ValidateHeader(path)` para validação rápida

Função pública em `storage` que lê apenas os primeiros 49 bytes (header) do arquivo, valida magic bytes e versão do formato. Retorna `nil`, `ErrInvalidMagic` ou `ErrVersionTooNew`. **Sem derivação de chave** — executa em microssegundos.

Usada pelo `AbrirCofreOperation` tanto no fluxo GUI (após o picker retornar caminho) quanto na entrada via CLI (antes de pedir senha).

```go
func ValidateHeader(path string) error
```

### 5. `storage.NewFileRepositoryForOpen(path, password)` para abrir cofre existente

Novo construtor em `storage` que cria um `FileRepository` com `isNew=false`, sem exigir salt ou metadata antecipados. O `Carregar()` subsequente popula ambos internamente.

**Motivação:** usar `NewFileRepositoryForCreate` (que seta `isNew=true`) para abrir cofre existente faria o primeiro `Salvar` usar `SaveNew` (escrita direta) em vez de `Save` (protocolo atômico com .tmp/.bak), violando o Princípio do Salvamento Atômico.

```go
func NewFileRepositoryForOpen(path string, password []byte) *FileRepository {
    return &FileRepository{
        path:     path,
        password: password,
        salt:     nil,
        isNew:    false,
    }
}
```

### 6. Decisão do fluxo via CLI em `main.go`

Conforme `golden/fluxos.md` (Entrada via Linha de Comando), o `main.go` decide qual fluxo disparar:

| Condição do argumento `--vault` | Comportamento |
|---|---|
| Arquivo existe | Fluxo 1 (Abrir) a partir do passo 3 — guard e picker ignorados |
| Arquivo não existe, diretório pai existe | Fluxo 2 (Criar) a partir do passo 3 — guard e picker ignorados |
| Arquivo não existe, diretório pai não existe | Nenhum fluxo disparado — tela normal |

A decisão é feita em `main.go` via `os.Stat`/`filepath.Dir`, passando `caminhoInicial` para a operation correspondente.

### 7. Sinalização de progresso (SetBusy/Clear)

Toda operação de IO segue o padrão estabelecido na `QuitOperation`:
1. `SetBusy(mensagem)` imediatamente antes de despachar o goroutine
2. `Clear()` — ou `SetSuccess` / `SetError` — ao receber o resultado

### 8. Estrutura de arquivos

```
internal/storage/
├── repository.go              (modificar — adicionar NewFileRepositoryForOpen)
├── storage.go                 (modificar — adicionar ValidateHeader)

internal/tui/operation/
├── fake_operation.go          (existente — sem alteração)
├── quit_operation.go          (existente — refatorar para usar guard; remover vaultSaver)
├── vault_saver.go             (novo — interface vaultSaver compartilhada)
├── guard_cofre_alterado.go    (novo — helper para passo 1)
├── criar_cofre.go             (novo — Fluxo 2)
└── abrir_cofre.go             (novo — Fluxo 1)

internal/tui/design/
└── keys.go                    (modificar — adicionar Shortcuts.NewVault e Shortcuts.OpenVault)

cmd/abditum/
├── main.go                    (modificar — lógica de decisão --vault)
└── setup.go                   (modificar — registrar ações Ctrl+N e Ctrl+O)
```

---

## Componentes

### `storage.ValidateHeader`

Adicionada em `storage/storage.go`.

```go
// ValidateHeader lê o header do arquivo de cofre e valida magic bytes e versão
// do formato. Não faz derivação de chave nem descriptografia.
// Retorna nil se o header é válido, ErrInvalidMagic ou ErrVersionTooNew caso contrário.
func ValidateHeader(path string) error
```

Implementação: lê os primeiros `HeaderSize` bytes, verifica magic bytes (`ABDT`) e chama `ProfileForVersion(version)` para validar a versão.

---

### `storage.NewFileRepositoryForOpen`

Adicionada em `storage/repository.go`.

```go
// NewFileRepositoryForOpen creates a FileRepository for opening an existing vault.
//
// Unlike NewFileRepositoryForCreate, this sets isNew=false so that Salvar uses
// the atomic Save protocol (not SaveNew). Salt and metadata are populated by
// the subsequent Carregar() call.
func NewFileRepositoryForOpen(path string, password []byte) *FileRepository {
    return &FileRepository{
        path:     path,
        password: password,
        salt:     nil,
        isNew:    false,
    }
}
```

---

### `vault_saver.go`

Move a interface `vaultSaver` de `quit_operation.go` para cá. Sem outras mudanças.

---

### `guard_cofre_alterado.go`

Helper interno que encapsula o passo 1 dos Fluxos 1, 2 e 5.

**Assinatura:**

```go
type guardCofreAlterado struct {
    saver      vaultSaver
    notifier   tui.MessageController
    onProceder func() tea.Cmd
    onAbortado func() tea.Cmd
}

func novoGuardCofreAlterado(
    notifier tui.MessageController,
    saver vaultSaver,
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
  → ErrModifiedExternally: Clear() → OpenModal(buildConflictModal())
  → erro genérico: SetError(err) → onAbortado()

Update(guardSaveResultMsg após conflito, forçado=true):
  → sucesso: Clear() → onProceder()
  → falha: SetError(err) → onAbortado()
```

**Modal de decisão** (título "Alterações não salvas"):
- `Enter` → Salvar e prosseguir
- `D` → Descartar e prosseguir → `onProceder()` direto
- `Esc` → Voltar → `onAbortado()`

**Modal de conflito** (título "Conflito"):
- `Enter` → Sobrescrever e prosseguir → `Salvar(true)`
- `Esc` → Voltar → `onAbortado()`

---

### `quit_operation.go` — Refatoração

A `QuitOperation` é simplificada para delegar ao guard quando há cofre alterado.

```
Init():
  manager != nil && IsModified():
    → cria guardCofreAlterado com:
        onProceder = tea.Quit
        onAbortado = OperationCompleted()
    → guard.Init()
  senão:
    → buildConfirmModal() (Fluxos 3/4: confirmação simples)

Update(msg):
  se guard != nil:
    → delega ao guard.Update(msg)
  senão:
    → trata apenas msgs do modal de confirmação simples (como antes)
```

O guard é armazenado em campo `guard *guardCofreAlterado` na struct.

**Nota:** o modal de confirmação simples dos Fluxos 3/4 ("Deseja encerrar a aplicação?") permanece na `QuitOperation` — não faz parte do guard, pois não envolve verificação de cofre alterado.

---

### `criar_cofre.go` — `CriarCofreOperation`

Implementa o Fluxo 2 completo.

**Construtor:**

```go
func NewCriarCofreOperation(
    notifier tui.MessageController,
    manager *vault.Manager, // nil se sem cofre carregado
    caminhoInicial string,  // "" = fluxo completo; preenchido = entrada via CLI (pula guard + picker)
) *CriarCofreOperation
```

**Estado inicial depende do construtor:**
- `caminhoInicial == ""` → começa em `criandoGuardando` (fluxo completo via GUI).
- `caminhoInicial != ""` → `Init()` vai direto para `criandoInformandoSenha`, sem guard nem picker. (Para Criar não há validação de magic/versão pois o arquivo ainda não existe.)

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
  se caminhoInicial != "":
    → guardar caminho + emitir msg para → criandoInformandoSenha
  senão:
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
    onConfirm(password): guardar senha;
      EvaluatePasswordStrength == StrengthWeak → criandoAvaliacaoSenhaFraca
      StrengthStrong → criandoCriando
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
      cofre.InicializarConteudoPadrao()
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

**Nota sobre salvamento atômico na criação:**
- Destino é arquivo novo: `NewFileRepositoryForCreate` → `manager.Salvar(false)` → `SaveNew` (escrita direta). Correto: não há arquivo existente para proteger.
- Destino é arquivo existente (sobrescrita): o `FileRepository` detecta `isNew=true` e usa `SaveNew`, que escreve diretamente. O protocolo atômico (.tmp/.bak) **não é usado** neste caso pela implementação atual de `FileRepository` — o primeiro save de um `ForCreate` sempre usa `SaveNew`. Isso é aceitável porque: (a) o fluxo já confirmou a sobrescrita com o usuário, e (b) não há cofre anterior "nosso" para proteger — o arquivo pertencia a outra sessão.

---

### `abrir_cofre.go` — `AbrirCofreOperation`

Implementa o Fluxo 1 completo.

**Construtor:**

```go
func NewAbrirCofreOperation(
    notifier tui.MessageController,
    manager *vault.Manager, // nil se sem cofre carregado
    caminhoInicial string,  // "" = fluxo completo; preenchido = entrada via CLI (pula guard + picker)
) *AbrirCofreOperation
```

**Estado inicial depende do construtor:**
- `caminhoInicial == ""` → começa em `abrindoGuardando` (fluxo completo via GUI).
- `caminhoInicial != ""` → `Init()` valida header via `storage.ValidateHeader(caminho)` e vai direto para `abrindoInformandoSenha`, sem guard nem picker. Se a validação falhar, emite `SetError` + `OperationCompleted()`.

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
  se caminhoInicial != "":
    err = storage.ValidateHeader(caminho) — síncrono, sem busy
    inválido → SetError(erroDeAberturaCategoria(err)) + OperationCompleted()
    ok → guardar caminho + emitir msg para → abrindoInformandoSenha
  senão:
    → instancia guardCofreAlterado
      onProceder = emitir msg para → abrindoInformandoCaminho
      onAbortado = OperationCompleted()
    → guard.Init()

abrindoInformandoCaminho:
  → OpenModal(FilePicker{Mode: Open})
  → desistiu: OperationCompleted()
  → caminho selecionado:
      err = storage.ValidateHeader(caminho) — IO síncrono rápido, sem busy
      inválido → SetError(erroDeAberturaCategoria(err)) + volta a abrindoInformandoCaminho
      ok → abrindoInformandoSenha (guardar caminho)

abrindoInformandoSenha:
  → OpenModal(PasswordEntryModal)
    onConfirm(password): → abrindoAbrindo
    onCancel(): → abrindoInformandoCaminho

abrindoAbrindo:
  → SetBusy("Abrindo cofre...")
  → goroutine:
      repo = NewFileRepositoryForOpen(caminho, senha)
      cofre, err = repo.Carregar()
      se err != nil → retorna abrirCofreResultMsg{err: err}
      manager = NewManager(cofre, repo)
      → retorna abrirCofreResultMsg{manager: manager}
  → senha errada / autenticação falhou (crypto.ErrAuthFailed):
      Clear()
      SetError("Senha incorreta ou arquivo corrompido.")
      → abrindoInformandoSenha
  → payload corrompido / PastaGeral ausente (storage.ErrCorrupted):
      Clear()
      SetError("Arquivo corrompido ou inválido.")
      → abrindoInformandoCaminho
  → sucesso:
      Clear()
      VaultOpenedMsg{Manager: manager}
      SetSuccess("Cofre aberto.")
```

**Nota sobre erros genéricos:** a spec do Fluxo 1 exige mensagens de erro por *categoria* (autenticação, integridade) sem revelar a causa exata — por segurança. A implementação usa essas duas categorias.

**Função auxiliar `erroDeAberturaCategoria`:**

```go
func erroDeAberturaCategoria(err error) string {
    if errors.Is(err, crypto.ErrAuthFailed) {
        return "Senha incorreta ou arquivo corrompido."
    }
    if errors.Is(err, storage.ErrCorrupted) {
        return "Arquivo corrompido ou inválido."
    }
    if errors.Is(err, storage.ErrInvalidMagic) {
        return "O arquivo selecionado não é um cofre Abditum."
    }
    if errors.Is(err, storage.ErrVersionTooNew) {
        return "O cofre foi criado com uma versão mais recente do Abditum. Atualize o aplicativo."
    }
    return "Não foi possível abrir o cofre."
}
```

---

## Integração em `setup.go`

Dois novos atalhos registrados:

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
        return tui.StartOperation(operation.NewCriarCofreOperation(r.MessageController(), r.Manager(), ""))
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
        return tui.StartOperation(operation.NewAbrirCofreOperation(r.MessageController(), r.Manager(), ""))
    },
},
```

Os atalhos `Ctrl+N` e `Ctrl+O` precisam ser adicionados a `design/keys.go` como `Shortcuts.NewVault` e `Shortcuts.OpenVault`.

---

## Integração em `main.go` — Entrada via CLI

```go
if vaultPath != "" {
    info, err := os.Stat(vaultPath)
    if err == nil && !info.IsDir() {
        // Arquivo existe → Fluxo 1 (Abrir) a partir do passo 3
        root.SetInitialOperation(
            operation.NewAbrirCofreOperation(root.MessageController(), nil, vaultPath),
        )
    } else if os.IsNotExist(err) {
        dir := filepath.Dir(vaultPath)
        if dirInfo, dirErr := os.Stat(dir); dirErr == nil && dirInfo.IsDir() {
            // Arquivo não existe, dir pai existe → Fluxo 2 (Criar) a partir do passo 3
            root.SetInitialOperation(
                operation.NewCriarCofreOperation(root.MessageController(), nil, vaultPath),
            )
        }
        // Dir pai não existe → nada, tela normal de abertura
    }
}
```

**Nota:** `SetInitialOperation` (ou `SetInitialCommand`) precisa ser verificado/adicionado ao `RootModel` — ver seção Pontos a verificar.

---

## Testes

Cada novo arquivo terá seu `_test.go` correspondente, seguindo o padrão de `quit_operation_test.go`:

- `vault_saver.go` — sem teste (só interface)
- `guard_cofre_alterado_test.go` — testa os 3 caminhos: sem cofre, cofre inalterado, cofre alterado (salvar ok, salvar falha, descartar, conflito externo, forçar ok, forçar falha)
- `quit_operation_test.go` — testes existentes devem continuar passando após refatoração
- `criar_cofre_test.go` — testa estados principais: guard integrado, picker, sobrescrita, senha fraca, criação ok, criação falha, caminhoInicial
- `abrir_cofre_test.go` — testa: guard integrado, picker, magic inválido, versão incompatível, senha errada, integridade, sucesso, caminhoInicial
- `storage/` — testes para `ValidateHeader` e `NewFileRepositoryForOpen`

Todos os testes de operations usam stubs das interfaces (`vaultSaver`, `tui.MessageController`) — sem IO real.

---

## Pontos a verificar durante implementação

1. **`PasswordCreateModal.onConfirm`** — o callback recebe apenas `password []byte` (sem indicação de força). A avaliação de força deve ser feita na operation via `crypto.EvaluatePasswordStrength(password)`.
2. **Erros de backup** — verificar se `storage` já expõe erros diferenciados para "falha após backup criado" para mensagens no `criandoCriando`.
3. **`FilePicker` modo Open** — não valida magic/versão; a validação ocorre na operation via `storage.ValidateHeader` após o picker retornar o caminho.
4. **`RootModel.SetInitialOperation`** — verificar se existe ou criar método para despachar uma operation na inicialização (para suportar `--vault`).
5. **Refatoração da `QuitOperation`** — o guard deve ser armazenado como campo na struct para que `Update` possa delegar. Mensagens do guard (`guardSaveMsg`, `guardSaveResultMsg`) devem ser processadas pelo `Update` da `QuitOperation` e delegadas ao guard.
