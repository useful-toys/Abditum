# Criar/Abrir Cofre Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implementar os Fluxos 1 (Abrir Cofre) e 2 (Criar Cofre) conforme `golden/fluxos.md`, incluindo entrada via CLI `--vault`, reuso via `guardCofreAlterado`, refatoração da `QuitOperation`, e atalhos `Ctrl+N`/`Ctrl+O`.

**Architecture:** Cada fluxo é uma `Operation` (interface `Init/Update`) com máquina de estados interna. O passo inicial de verificar cofre alterado é extraído como helper `guardCofreAlterado` (interno ao pacote `operation/`) reutilizado por ambas as operations e pela `QuitOperation` refatorada. A interface `vaultSaver` é movida de `quit_operation.go` para `vault_saver.go` para ser compartilhada. Novas funções `storage.ValidateHeader` e `storage.NewFileRepositoryForOpen` são adicionadas. A entrada via CLI usa o padrão `RootModelOption` existente.

**Tech Stack:** Go, Bubble Tea v2, pacotes internos `storage`, `vault`, `crypto`, `tui/modal`, `tui/design`.

**Spec de referência:** `docs/superpowers/specs/2026-04-18-criar-abrir-cofre-design.md`

---

## Estrutura de Arquivos

| Arquivo | Ação | Responsabilidade |
|---------|------|-----------------|
| `internal/storage/storage.go` | Modificar | Adicionar `ValidateHeader` |
| `internal/storage/storage_test.go` | Modificar | Testes de `ValidateHeader` |
| `internal/storage/repository.go` | Modificar | Adicionar `NewFileRepositoryForOpen` |
| `internal/storage/repository_test.go` | Modificar | Testes de `NewFileRepositoryForOpen` |
| `internal/tui/operation/vault_saver.go` | Criar | Interface `vaultSaver` compartilhada |
| `internal/tui/operation/quit_operation.go` | Modificar | Remover `vaultSaver` (movida); refatorar para usar guard |
| `internal/tui/operation/guard_cofre_alterado.go` | Criar | Helper passo 1 (verificar cofre modificado) |
| `internal/tui/operation/guard_cofre_alterado_test.go` | Criar | Testes do guard |
| `internal/tui/operation/criar_cofre.go` | Criar | Fluxo 2 completo |
| `internal/tui/operation/criar_cofre_test.go` | Criar | Testes do Fluxo 2 |
| `internal/tui/operation/abrir_cofre.go` | Criar | Fluxo 1 completo |
| `internal/tui/operation/abrir_cofre_test.go` | Criar | Testes do Fluxo 1 |
| `internal/tui/design/keys.go` | Modificar | Adicionar `Shortcuts.NewVault` e `Shortcuts.OpenVault` |
| `internal/tui/root.go` | Modificar | Adicionar `WithInitialVault` option |
| `cmd/abditum/setup.go` | Modificar | Registrar ações Ctrl+N e Ctrl+O |
| `cmd/abditum/main.go` | Modificar | Decidir fluxo via `--vault` usando `WithInitialVault` |

---

## Task 1: Adicionar `storage.ValidateHeader` e `storage.NewFileRepositoryForOpen`

**Files:**
- Modify: `internal/storage/storage.go`
- Modify: `internal/storage/repository.go`
- Modify: `internal/storage/storage_test.go` (ou criar se não existir)
- Modify: `internal/storage/repository_test.go` (ou criar se não existir)

- [ ] **Step 1: Escrever o teste de `ValidateHeader`**

Em `internal/storage/storage_test.go`, adicionar:

```go
func TestValidateHeader_ArquivoValido(t *testing.T) {
	// Criar um cofre válido para testar
	dir := t.TempDir()
	path := filepath.Join(dir, "test.abditum")
	cofre := vault.NovoCofre()
	if err := cofre.InicializarConteudoPadrao(); err != nil {
		t.Fatal(err)
	}
	if err := SaveNew(path, cofre, []byte("SenhaForte123!")); err != nil {
		t.Fatal(err)
	}

	if err := ValidateHeader(path); err != nil {
		t.Errorf("ValidateHeader de cofre válido: erro inesperado %v", err)
	}
}

func TestValidateHeader_ArquivoPequenoDemais(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "small.bin")
	os.WriteFile(path, []byte("ABC"), 0600) // menos que HeaderSize

	err := ValidateHeader(path)
	if !errors.Is(err, ErrInvalidMagic) {
		t.Errorf("arquivo pequeno: esperado ErrInvalidMagic, obteve %v", err)
	}
}

func TestValidateHeader_MagicInvalida(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad_magic.bin")
	data := make([]byte, HeaderSize)
	copy(data[0:4], []byte("XXXX"))
	data[4] = CurrentFormatVersion
	os.WriteFile(path, data, 0600)

	err := ValidateHeader(path)
	if !errors.Is(err, ErrInvalidMagic) {
		t.Errorf("magic inválida: esperado ErrInvalidMagic, obteve %v", err)
	}
}

func TestValidateHeader_VersaoIncompativel(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "future.bin")
	data := make([]byte, HeaderSize)
	copy(data[0:4], Magic[:])
	data[4] = 255 // versão futura
	os.WriteFile(path, data, 0600)

	err := ValidateHeader(path)
	if !errors.Is(err, ErrVersionTooNew) {
		t.Errorf("versão futura: esperado ErrVersionTooNew, obteve %v", err)
	}
}

func TestValidateHeader_ArquivoInexistente(t *testing.T) {
	err := ValidateHeader("/caminho/que/nao/existe.abditum")
	if err == nil {
		t.Error("arquivo inexistente: esperado erro, obteve nil")
	}
}
```

Imports necessários: `"errors"`, `"os"`, `"path/filepath"`, `"testing"`, e os pacotes internos `vault`.

- [ ] **Step 2: Rodar os testes para confirmar que falham**

```
go test ./internal/storage/... -run TestValidateHeader
```
Expected: FAIL (`ValidateHeader` não existe)

- [ ] **Step 3: Implementar `ValidateHeader` em `storage.go`**

Adicionar ao final de `internal/storage/storage.go`:

```go
// ValidateHeader lê o header do arquivo de cofre e valida magic bytes e versão
// do formato. Não faz derivação de chave nem descriptografia — executa em
// microssegundos.
//
// Retorna nil se o header é válido, ErrInvalidMagic se o arquivo não tem os
// magic bytes "ABDT", ErrVersionTooNew se a versão do formato não é suportada,
// ou erro de IO se o arquivo não pode ser lido.
func ValidateHeader(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("storage.ValidateHeader: %w", err)
	}
	defer f.Close()

	header := make([]byte, HeaderSize)
	n, err := io.ReadFull(f, header)
	if err != nil || n < HeaderSize {
		return ErrInvalidMagic
	}

	if header[0] != Magic[0] || header[1] != Magic[1] || header[2] != Magic[2] || header[3] != Magic[3] {
		return ErrInvalidMagic
	}

	version := header[MagicSize]
	if _, err := ProfileForVersion(version); err != nil {
		return err
	}

	return nil
}
```

- [ ] **Step 4: Rodar os testes para confirmar que passam**

```
go test ./internal/storage/... -run TestValidateHeader
```
Expected: PASS

- [ ] **Step 5: Escrever o teste de `NewFileRepositoryForOpen`**

Em `internal/storage/repository_test.go`, adicionar:

```go
func TestNewFileRepositoryForOpen_Carregar_E_Salvar_Atomico(t *testing.T) {
	// Criar um cofre primeiro
	dir := t.TempDir()
	path := filepath.Join(dir, "test.abditum")
	cofre := vault.NovoCofre()
	if err := cofre.InicializarConteudoPadrao(); err != nil {
		t.Fatal(err)
	}
	password := []byte("SenhaForte123!")
	if err := SaveNew(path, cofre, password); err != nil {
		t.Fatal(err)
	}

	// Abrir com NewFileRepositoryForOpen
	repo := NewFileRepositoryForOpen(path, password)

	// Carregar deve funcionar
	loaded, err := repo.Carregar()
	if err != nil {
		t.Fatalf("Carregar: %v", err)
	}
	if loaded == nil {
		t.Fatal("Carregar retornou cofre nil")
	}

	// Salvar deve usar protocolo atômico (Save, não SaveNew)
	// Verificamos que .bak é criado como prova do protocolo atômico
	if err := repo.Salvar(loaded); err != nil {
		t.Fatalf("Salvar: %v", err)
	}
	bakPath := path + ".bak"
	if _, err := os.Stat(bakPath); os.IsNotExist(err) {
		t.Error("Salvar após Carregar via ForOpen: .bak não foi criado — protocolo atômico não usado")
	}
}

func TestNewFileRepositoryForOpen_SenhaErrada_ErroDeAutenticacao(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.abditum")
	cofre := vault.NovoCofre()
	if err := cofre.InicializarConteudoPadrao(); err != nil {
		t.Fatal(err)
	}
	if err := SaveNew(path, cofre, []byte("SenhaCorreta1!")); err != nil {
		t.Fatal(err)
	}

	repo := NewFileRepositoryForOpen(path, []byte("SenhaErrada1!"))
	_, err := repo.Carregar()
	if !errors.Is(err, crypto.ErrAuthFailed) {
		t.Errorf("senha errada: esperado ErrAuthFailed, obteve %v", err)
	}
}
```

Imports: `"errors"`, `"os"`, `"path/filepath"`, `"testing"`, `"github.com/useful-toys/abditum/internal/crypto"`, `"github.com/useful-toys/abditum/internal/vault"`.

- [ ] **Step 6: Rodar os testes para confirmar que falham**

```
go test ./internal/storage/... -run TestNewFileRepositoryForOpen
```
Expected: FAIL (`NewFileRepositoryForOpen` não existe)

- [ ] **Step 7: Implementar `NewFileRepositoryForOpen` em `repository.go`**

Adicionar em `internal/storage/repository.go`, após `NewFileRepositoryForCreate`:

```go
// NewFileRepositoryForOpen creates a FileRepository for opening an existing vault.
//
// Unlike NewFileRepositoryForCreate, this sets isNew=false so that Salvar uses
// the atomic Save protocol (with .tmp/.bak rotation) from the first call.
// Salt and metadata are populated by the subsequent Carregar() call.
//
// Parameters:
//   - path: Absolute path to the existing vault file.
//   - password: The user's master password (UTF-8 bytes). Stored by reference.
func NewFileRepositoryForOpen(path string, password []byte) *FileRepository {
	return &FileRepository{
		path:     path,
		password: password,
		salt:     nil,
		isNew:    false,
	}
}
```

- [ ] **Step 8: Rodar os testes para confirmar que passam**

```
go test ./internal/storage/... -run TestNewFileRepositoryForOpen
```
Expected: PASS

- [ ] **Step 9: Rodar todos os testes do pacote storage**

```
go test ./internal/storage/...
```
Expected: PASS

- [ ] **Step 10: Commit**

```
git add internal/storage/storage.go internal/storage/repository.go internal/storage/storage_test.go internal/storage/repository_test.go
git commit -m "feat: adicionar ValidateHeader e NewFileRepositoryForOpen ao pacote storage"
```

---

## Task 2: Mover `vaultSaver` para arquivo compartilhado

**Files:**
- Create: `internal/tui/operation/vault_saver.go`
- Modify: `internal/tui/operation/quit_operation.go`

- [ ] **Step 1: Criar `vault_saver.go`**

```go
package operation

// vaultSaver é a interface mínima que as operations precisam do vault.Manager.
// Usar uma interface aqui (em vez do tipo concreto) facilita os testes.
type vaultSaver interface {
	IsModified() bool
	Salvar(forcarSobrescrita bool) error
}
```

- [ ] **Step 2: Remover `vaultSaver` de `quit_operation.go`**

Em `quit_operation.go`, remover as linhas 13-18:
```go
// vaultSaver é a interface mínima que a QuitOperation precisa do vault.Manager.
// Usar uma interface aqui (em vez do tipo concreto) facilita os testes.
type vaultSaver interface {
	IsModified() bool
	Salvar(forcarSobrescrita bool) error
}
```

- [ ] **Step 3: Verificar que os testes existentes ainda passam**

```
go test ./internal/tui/operation/...
```
Expected: PASS (todos os testes de `quit_operation_test.go`)

- [ ] **Step 4: Commit**

```
git add internal/tui/operation/vault_saver.go internal/tui/operation/quit_operation.go
git commit -m "refactor: mover interface vaultSaver para vault_saver.go compartilhado"
```

---

## Task 3: Implementar `guardCofreAlterado`

**Files:**
- Create: `internal/tui/operation/guard_cofre_alterado.go`
- Create: `internal/tui/operation/guard_cofre_alterado_test.go`

- [ ] **Step 1: Escrever os testes**

Criar `internal/tui/operation/guard_cofre_alterado_test.go`:

```go
package operation

import (
	"errors"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/vault"
)

// --- Init ---

func TestGuard_Init_SemCofre_ChamaOnProceder(t *testing.T) {
	var chamado bool
	g := novoGuardCofreAlterado(
		&stubNotifier{}, nil,
		func() tea.Cmd { chamado = true; return nil },
		func() tea.Cmd { t.Error("onAbortado não deveria ser chamado"); return nil },
	)
	execCmd(g.Init())
	if !chamado {
		t.Error("Init sem cofre: onProceder não foi chamado")
	}
}

func TestGuard_Init_CofreInalterado_ChamaOnProceder(t *testing.T) {
	var chamado bool
	g := novoGuardCofreAlterado(
		&stubNotifier{},
		&stubManager{isModified: false},
		func() tea.Cmd { chamado = true; return nil },
		func() tea.Cmd { t.Error("onAbortado não deveria ser chamado"); return nil },
	)
	execCmd(g.Init())
	if !chamado {
		t.Error("Init cofre inalterado: onProceder não foi chamado")
	}
}

func TestGuard_Init_CofreAlterado_AbreModal(t *testing.T) {
	g := novoGuardCofreAlterado(
		&stubNotifier{},
		&stubManager{isModified: true},
		func() tea.Cmd { return nil },
		func() tea.Cmd { return nil },
	)
	msg := execCmd(g.Init())
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Init cofre alterado: esperado OpenModalMsg, obteve %T", msg)
	}
}

// --- Update: guardSaveMsg ---

func TestGuard_Update_SalvarSucesso_ChamaOnProceder(t *testing.T) {
	n := &stubNotifier{}
	var chamado bool
	g := novoGuardCofreAlterado(
		n,
		&stubManager{isModified: true},
		func() tea.Cmd { chamado = true; return nil },
		func() tea.Cmd { t.Error("onAbortado inesperado"); return nil },
	)
	// Disparar salvamento
	cmd := g.Update(guardSaveMsg{forced: false})
	if n.lastMethod != "SetBusy" {
		t.Errorf("esperado SetBusy, obteve %q", n.lastMethod)
	}
	// Executar goroutine e processar resultado
	resultMsg := execCmd(cmd)
	execCmd(g.Update(resultMsg))
	if !chamado {
		t.Error("após salvar OK: onProceder não foi chamado")
	}
	if n.lastMethod != "Clear" {
		t.Errorf("após salvar OK: esperado Clear, obteve %q", n.lastMethod)
	}
}

func TestGuard_Update_SalvarErroGenerico_ChamaOnAbortado(t *testing.T) {
	n := &stubNotifier{}
	var abortado bool
	g := novoGuardCofreAlterado(
		n,
		&stubManager{isModified: true, salvarErr: errors.New("disco cheio")},
		func() tea.Cmd { t.Error("onProceder inesperado"); return nil },
		func() tea.Cmd { abortado = true; return nil },
	)
	cmd := g.Update(guardSaveMsg{forced: false})
	resultMsg := execCmd(cmd)
	execCmd(g.Update(resultMsg))
	if !abortado {
		t.Error("após erro genérico: onAbortado não foi chamado")
	}
	if n.lastMethod != "SetError" {
		t.Errorf("após erro genérico: esperado SetError, obteve %q", n.lastMethod)
	}
}

func TestGuard_Update_ModificadoExternamente_AbreModalConflito(t *testing.T) {
	n := &stubNotifier{}
	g := novoGuardCofreAlterado(
		n,
		&stubManager{isModified: true, salvarErr: vault.ErrModifiedExternally},
		func() tea.Cmd { return nil },
		func() tea.Cmd { return nil },
	)
	cmd := g.Update(guardSaveMsg{forced: false})
	resultMsg := execCmd(cmd)
	resultCmd := g.Update(resultMsg)
	msg := execCmd(resultCmd)
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("ErrModifiedExternally: esperado OpenModalMsg (conflito), obteve %T", msg)
	}
	if n.lastMethod != "Clear" {
		t.Errorf("ErrModifiedExternally: esperado Clear, obteve %q", n.lastMethod)
	}
}

func TestGuard_Update_SalvarForcado_Sucesso_ChamaOnProceder(t *testing.T) {
	n := &stubNotifier{}
	var chamado bool
	g := novoGuardCofreAlterado(
		n,
		&stubManager{isModified: true},
		func() tea.Cmd { chamado = true; return nil },
		func() tea.Cmd { t.Error("onAbortado inesperado"); return nil },
	)
	cmd := g.Update(guardSaveMsg{forced: true})
	resultMsg := execCmd(cmd)
	execCmd(g.Update(resultMsg))
	if !chamado {
		t.Error("após salvar forçado OK: onProceder não foi chamado")
	}
}

func TestGuard_Update_SalvarForcado_Erro_ChamaOnAbortado(t *testing.T) {
	n := &stubNotifier{}
	var abortado bool
	m := &stubManager{isModified: true, salvarErr: errors.New("falha")}
	g := novoGuardCofreAlterado(
		n, m,
		func() tea.Cmd { t.Error("onProceder inesperado"); return nil },
		func() tea.Cmd { abortado = true; return nil },
	)
	cmd := g.Update(guardSaveMsg{forced: true})
	resultMsg := execCmd(cmd)
	execCmd(g.Update(resultMsg))
	if !abortado {
		t.Error("após forçado com erro: onAbortado não foi chamado")
	}
}

// --- Update: descarte ---

func TestGuard_Update_DescartarMsg_ChamaOnProceder(t *testing.T) {
	var chamado bool
	g := novoGuardCofreAlterado(
		&stubNotifier{},
		&stubManager{isModified: true},
		func() tea.Cmd { chamado = true; return nil },
		func() tea.Cmd { t.Error("onAbortado inesperado"); return nil },
	)
	execCmd(g.Update(guardDiscardMsg{}))
	if !chamado {
		t.Error("descartar: onProceder não foi chamado")
	}
}

// --- Update: cancelar ---

func TestGuard_Update_CancelarMsg_ChamaOnAbortado(t *testing.T) {
	var abortado bool
	g := novoGuardCofreAlterado(
		&stubNotifier{},
		&stubManager{isModified: true},
		func() tea.Cmd { t.Error("onProceder inesperado"); return nil },
		func() tea.Cmd { abortado = true; return nil },
	)
	execCmd(g.Update(guardCancelMsg{}))
	if !abortado {
		t.Error("cancelar: onAbortado não foi chamado")
	}
}
```

- [ ] **Step 2: Rodar os testes para confirmar que falham**

```
go test ./internal/tui/operation/... -run TestGuard
```
Expected: FAIL

- [ ] **Step 3: Criar `guard_cofre_alterado.go`**

```go
package operation

import (
	"errors"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/vault"
)

// guardSaveMsg dispara a ação de salvar no guard.
type guardSaveMsg struct {
	forced bool // true = forçar sobrescrita de modificação externa
}

// guardSaveResultMsg carrega o resultado da tentativa de salvar no guard.
type guardSaveResultMsg struct {
	err    error
	forced bool
}

// guardDiscardMsg sinaliza que o usuário optou por descartar alterações.
type guardDiscardMsg struct{}

// guardCancelMsg sinaliza que o usuário optou por voltar (cancelar).
type guardCancelMsg struct{}

// guardCofreAlterado verifica se há um cofre com alterações não salvas antes
// de prosseguir para outra operação. É um helper interno do pacote operation/,
// compartilhado pelos Fluxos 1, 2 e 5.
//
// Se não há cofre ou o cofre está inalterado, chama onProceder() diretamente.
// Se há alterações, abre modal de decisão (Salvar / Descartar / Voltar).
// Se salvar falha com ErrModifiedExternally, abre modal de conflito.
type guardCofreAlterado struct {
	saver      vaultSaver
	notifier   tui.MessageController
	onProceder func() tea.Cmd
	onAbortado func() tea.Cmd
}

// novoGuardCofreAlterado cria um guardCofreAlterado.
// saver pode ser nil quando nenhum cofre está carregado.
func novoGuardCofreAlterado(
	notifier tui.MessageController,
	saver vaultSaver,
	onProceder func() tea.Cmd,
	onAbortado func() tea.Cmd,
) *guardCofreAlterado {
	return &guardCofreAlterado{
		saver:      saver,
		notifier:   notifier,
		onProceder: onProceder,
		onAbortado: onAbortado,
	}
}

// Init inicia o guard. Se não há cofre ou não há alterações, dispara onProceder
// imediatamente. Caso contrário, abre o modal de decisão.
func (g *guardCofreAlterado) Init() tea.Cmd {
	if g.saver == nil || !g.saver.IsModified() {
		return g.onProceder()
	}
	return tui.OpenModal(g.buildModifiedModal())
}

// Update trata as mensagens internas do guard.
func (g *guardCofreAlterado) Update(msg tea.Msg) tea.Cmd {
	switch m := msg.(type) {
	case guardSaveMsg:
		g.notifier.SetBusy("Salvando...")
		forced := m.forced
		return func() tea.Msg {
			return guardSaveResultMsg{err: g.saver.Salvar(forced), forced: forced}
		}

	case guardSaveResultMsg:
		if m.err == nil {
			g.notifier.Clear()
			return g.onProceder()
		}
		if !m.forced && errors.Is(m.err, vault.ErrModifiedExternally) {
			g.notifier.Clear()
			return tui.OpenModal(g.buildConflictModal())
		}
		g.notifier.SetError(m.err.Error())
		return g.onAbortado()

	case guardDiscardMsg:
		return g.onProceder()

	case guardCancelMsg:
		return g.onAbortado()
	}
	return nil
}

// buildModifiedModal cria o modal de decisão quando há alterações não salvas.
func (g *guardCofreAlterado) buildModifiedModal() *modal.ConfirmModal {
	return modal.NewConfirmModal(
		"Alterações não salvas",
		"Há alterações não salvas. O que deseja fazer?",
		[]modal.ModalOption{
			{
				Keys:  []design.Key{design.Keys.Enter},
				Label: "Salvar e prosseguir",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return guardSaveMsg{forced: false}
					})
				},
			},
			{
				Keys:  []design.Key{design.Letter('d')},
				Label: "Descartar e prosseguir",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return guardDiscardMsg{}
					})
				},
			},
			{
				Keys:  []design.Key{design.Keys.Esc},
				Label: "Voltar",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return guardCancelMsg{}
					})
				},
			},
		},
	)
}

// buildConflictModal cria o modal de resolução de conflito quando o arquivo foi
// modificado externamente enquanto o cofre estava aberto.
func (g *guardCofreAlterado) buildConflictModal() *modal.ConfirmModal {
	return modal.NewConfirmModal(
		"Conflito",
		"O arquivo foi modificado externamente. Deseja sobrescrever?",
		[]modal.ModalOption{
			{
				Keys:  []design.Key{design.Keys.Enter},
				Label: "Sobrescrever e prosseguir",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return guardSaveMsg{forced: true}
					})
				},
			},
			{
				Keys:  []design.Key{design.Keys.Esc},
				Label: "Voltar",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return guardCancelMsg{}
					})
				},
			},
		},
	)
}
```

- [ ] **Step 4: Rodar os testes para confirmar que passam**

```
go test ./internal/tui/operation/... -run TestGuard
```
Expected: PASS

- [ ] **Step 5: Commit**

```
git add internal/tui/operation/guard_cofre_alterado.go internal/tui/operation/guard_cofre_alterado_test.go
git commit -m "feat: adicionar guardCofreAlterado helper para passo 1 dos Fluxos 1, 2 e 5"
```

---

## Task 4: Refatorar `QuitOperation` para usar `guardCofreAlterado`

**Files:**
- Modify: `internal/tui/operation/quit_operation.go`

- [ ] **Step 1: Refatorar `QuitOperation`**

Substituir o conteúdo de `internal/tui/operation/quit_operation.go` por:

```go
package operation

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/vault"
)

// QuitOperation implementa o fluxo de saída (ctrl+Q) do gerenciador.
//
// Três fluxos possíveis:
//   - Fluxo 3: sem cofre aberto → confirmação simples → sair
//   - Fluxo 4: cofre aberto mas sem alterações → confirmação simples → sair
//   - Fluxo 5: cofre com alterações → guardCofreAlterado → sair
type QuitOperation struct {
	notifier tui.MessageController
	manager  vaultSaver
	guard    *guardCofreAlterado // não-nil apenas no Fluxo 5
}

// NewQuitOperation cria uma QuitOperation com o vault.Manager concreto.
// manager pode ser nil quando nenhum cofre está aberto (Fluxo 3).
func NewQuitOperation(notifier tui.MessageController, manager *vault.Manager) *QuitOperation {
	var saver vaultSaver
	if manager != nil {
		saver = manager
	}
	return &QuitOperation{notifier: notifier, manager: saver}
}

// newQuitOperationFromSaver é usada internamente e nos testes (mesmo pacote).
// Permite injetar qualquer implementação de vaultSaver sem depender do tipo concreto.
func newQuitOperationFromSaver(notifier tui.MessageController, saver vaultSaver) *QuitOperation {
	return &QuitOperation{notifier: notifier, manager: saver}
}

// Init inicia o fluxo de saída exibindo o modal adequado conforme o estado do cofre.
func (q *QuitOperation) Init() tea.Cmd {
	if q.manager != nil && q.manager.IsModified() {
		// Fluxo 5: cofre com alterações não salvas — delegar ao guard
		q.guard = novoGuardCofreAlterado(
			q.notifier,
			q.manager,
			func() tea.Cmd {
				q.notifier.SetSuccess("Cofre salvo.")
				return tea.Quit
			},
			func() tea.Cmd { return tui.OperationCompleted() },
		)
		return q.guard.Init()
	}
	// Fluxo 3 (sem cofre) ou Fluxo 4 (cofre inalterado): confirmação simples
	return tui.OpenModal(q.buildConfirmModal())
}

// Update trata as mensagens internas da QuitOperation.
func (q *QuitOperation) Update(msg tea.Msg) tea.Cmd {
	if q.guard != nil {
		return q.guard.Update(msg)
	}
	// Fluxos 3/4: nenhuma mensagem interna — tudo tratado pelo modal de confirmação
	return nil
}

// buildConfirmModal cria o modal de confirmação simples para os Fluxos 3 e 4.
func (q *QuitOperation) buildConfirmModal() *modal.ConfirmModal {
	return modal.NewConfirmModal(
		"Sair",
		"Deseja encerrar a aplicação?",
		[]modal.ModalOption{
			{
				Keys:  []design.Key{design.Keys.Enter},
				Label: "Confirmar",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), tea.Quit)
				},
			},
			{
				Keys:  []design.Key{design.Keys.Esc},
				Label: "Voltar",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), tui.OperationCompleted())
				},
			},
		},
	)
}
```

**Nota:** Os imports de `"errors"` e `"github.com/useful-toys/abditum/internal/vault"` são removidos se não usados diretamente (o guard trata tudo). O import de `vault` é necessário apenas em `NewQuitOperation` para o tipo `*vault.Manager`. Verificar se compila.

- [ ] **Step 2: Rodar os testes existentes**

```
go test ./internal/tui/operation/... -run TestQuitOperation
```
Expected: PASS — todos os testes de `quit_operation_test.go` devem continuar passando.

**Análise de compatibilidade dos testes:**
- `TestQuitOperation_Init_SemCofre_AbreModalConfirmacao` → OK (Fluxo 3, guard não criado)
- `TestQuitOperation_Init_CofreInalterado_AbreModalConfirmacao` → OK (Fluxo 4, guard não criado)
- `TestQuitOperation_Init_CofreAlterado_AbreModalDecisao` → OK (Fluxo 5, guard.Init() abre modal)
- `TestQuitOperation_Update_Saving_Sucesso_EmiteQuit` → Precisa verificar: este teste envia `quitMsg{state: quitStateSaving}` que é um tipo da implementação antiga. Agora o guard usa `guardSaveMsg{forced: false}`. **Os testes de Update que enviam `quitMsg` vão falhar.**

**IMPORTANTE:** Os testes existentes usam `quitMsg` e `quitSaveResultMsg` que são tipos da implementação antiga. Após a refatoração, a `QuitOperation` delega ao guard, que usa `guardSaveMsg` e `guardSaveResultMsg`. Os testes devem ser atualizados para usar os novos tipos.

- [ ] **Step 3: Atualizar os testes de `QuitOperation` para a nova implementação**

Em `internal/tui/operation/quit_operation_test.go`, substituir os testes de Update:

```go
package operation

import (
	"errors"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/vault"
)

// stubManager simula vault.Manager para os testes da QuitOperation.
type stubManager struct {
	isModified                    bool
	salvarErr                     error
	salvarForcarSobrescritaCalled bool
}

func (s *stubManager) IsModified() bool { return s.isModified }
func (s *stubManager) Salvar(forcarSobrescrita bool) error {
	if forcarSobrescrita {
		s.salvarForcarSobrescritaCalled = true
	}
	return s.salvarErr
}

// --- Init ---

func TestQuitOperation_Init_SemCofre_AbreModalConfirmacao(t *testing.T) {
	op := newQuitOperationFromSaver(&stubNotifier{}, nil)
	msg := execCmd(op.Init())
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Init sem cofre: esperado OpenModalMsg, obteve %T", msg)
	}
}

func TestQuitOperation_Init_CofreInalterado_AbreModalConfirmacao(t *testing.T) {
	op := newQuitOperationFromSaver(&stubNotifier{}, &stubManager{isModified: false})
	msg := execCmd(op.Init())
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Init cofre inalterado: esperado OpenModalMsg, obteve %T", msg)
	}
}

func TestQuitOperation_Init_CofreAlterado_AbreModalDecisao(t *testing.T) {
	op := newQuitOperationFromSaver(&stubNotifier{}, &stubManager{isModified: true})
	msg := execCmd(op.Init())
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Init cofre alterado: esperado OpenModalMsg, obteve %T", msg)
	}
}

// --- Update (delegado ao guard no Fluxo 5) ---

func TestQuitOperation_Update_IgnoraMensagemDesconhecida(t *testing.T) {
	op := newQuitOperationFromSaver(&stubNotifier{}, nil)
	type outraMsg struct{}
	if cmd := op.Update(outraMsg{}); cmd != nil {
		t.Error("Update(outraMsg): esperado nil")
	}
}

func TestQuitOperation_Update_Saving_Sucesso_EmiteQuit(t *testing.T) {
	n := &stubNotifier{}
	m := &stubManager{}
	op := newQuitOperationFromSaver(n, m)
	// Inicializar sem alterações para guard=nil, depois recriamos com alterações
	// Na verdade precisamos inicializar com alterações para ativar o guard
	m.isModified = true
	op2 := newQuitOperationFromSaver(n, m)
	op2.Init() // ativa o guard

	cmd := op2.Update(guardSaveMsg{forced: false})
	if n.lastMethod != "SetBusy" {
		t.Errorf("Update(saving): esperado SetBusy, obteve %q", n.lastMethod)
	}
	resultMsg := execCmd(cmd)
	resultCmd := op2.Update(resultMsg)
	msg := execCmd(resultCmd)
	if n.lastMethod != "SetSuccess" {
		t.Errorf("após salvar OK: esperado SetSuccess, obteve %q", n.lastMethod)
	}
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Errorf("após salvar OK: esperado tea.QuitMsg, obteve %T", msg)
	}
}

func TestQuitOperation_Update_Saving_ErroExterno_AbreModalConflito(t *testing.T) {
	n := &stubNotifier{}
	m := &stubManager{isModified: true, salvarErr: vault.ErrModifiedExternally}
	op := newQuitOperationFromSaver(n, m)
	op.Init()

	cmd := op.Update(guardSaveMsg{forced: false})
	resultMsg := execCmd(cmd)
	resultCmd := op.Update(resultMsg)
	msg := execCmd(resultCmd)
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("após ErrModifiedExternally: esperado OpenModalMsg, obteve %T", msg)
	}
	if n.lastMethod != "Clear" {
		t.Errorf("após ErrModifiedExternally: esperado Clear, obteve %q", n.lastMethod)
	}
}

func TestQuitOperation_Update_Saving_ErroGenerico_NaoSai(t *testing.T) {
	n := &stubNotifier{}
	m := &stubManager{isModified: true, salvarErr: errors.New("disco cheio")}
	op := newQuitOperationFromSaver(n, m)
	op.Init()

	cmd := op.Update(guardSaveMsg{forced: false})
	resultMsg := execCmd(cmd)
	resultCmd := op.Update(resultMsg)
	msg := execCmd(resultCmd)
	if n.lastMethod != "SetError" {
		t.Errorf("após erro genérico: esperado SetError, obteve %q", n.lastMethod)
	}
	if _, ok := msg.(tui.OperationCompletedMsg); !ok {
		t.Errorf("após erro genérico: esperado OperationCompletedMsg, obteve %T", msg)
	}
}

func TestQuitOperation_Update_SavingForced_Sucesso_EmiteQuit(t *testing.T) {
	n := &stubNotifier{}
	m := &stubManager{isModified: true}
	op := newQuitOperationFromSaver(n, m)
	op.Init()

	cmd := op.Update(guardSaveMsg{forced: true})
	resultMsg := execCmd(cmd)
	resultCmd := op.Update(resultMsg)
	msg := execCmd(resultCmd)
	if !m.salvarForcarSobrescritaCalled {
		t.Error("SavingForced: Salvar(true) não foi chamado")
	}
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Errorf("SavingForced OK: esperado tea.QuitMsg, obteve %T", msg)
	}
}

func TestQuitOperation_Update_SavingForced_Erro_NaoSai(t *testing.T) {
	n := &stubNotifier{}
	m := &stubManager{isModified: true, salvarErr: errors.New("falha")}
	op := newQuitOperationFromSaver(n, m)
	op.Init()

	cmd := op.Update(guardSaveMsg{forced: true})
	resultMsg := execCmd(cmd)
	resultCmd := op.Update(resultMsg)
	msg := execCmd(resultCmd)
	if n.lastMethod != "SetError" {
		t.Errorf("SavingForced erro: esperado SetError, obteve %q", n.lastMethod)
	}
	if _, ok := msg.(tui.OperationCompletedMsg); !ok {
		t.Errorf("SavingForced erro: esperado OperationCompletedMsg, obteve %T", msg)
	}
}

// --- Update: Fluxo 5 — descartar ---

func TestQuitOperation_Update_Descartar_EmiteQuit(t *testing.T) {
	n := &stubNotifier{}
	m := &stubManager{isModified: true}
	op := newQuitOperationFromSaver(n, m)
	op.Init()

	// guardDiscardMsg → guard.onProceder → SetSuccess + tea.Quit
	cmd := op.Update(guardDiscardMsg{})
	msg := execCmd(cmd)
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Errorf("Descartar: esperado tea.QuitMsg, obteve %T", msg)
	}
}
```

**Nota sobre o teste de Descartar:** No guard da QuitOperation, `onProceder` chama `SetSuccess("Cofre salvo.")` e `tea.Quit`. Mas quando o usuário descarta, o guard chama `onProceder` diretamente (sem salvar). Isso faz o `SetSuccess("Cofre salvo.")` ser emitido mesmo quando o cofre não foi salvo. **Isso é um bug no design!**

A solução: o `onProceder` do guard na QuitOperation deve ser apenas `tea.Quit`. O `SetSuccess("Cofre salvo.")` deve ser emitido apenas quando o guard realmente salva com sucesso. Mas o guard não notifica se salvou ou descartou — ele apenas chama `onProceder` em ambos os casos.

**Correção:** O guard já emite `guardSaveResultMsg` com `err == nil` quando salva com sucesso, e `guardDiscardMsg` quando descarta. O `onProceder` é chamado em ambos os casos. Para distinguir, o `onProceder` da QuitOperation deve ser simplesmente `tea.Quit`, sem `SetSuccess`. O `SetSuccess` deve ser adicionado *dentro* do guard apenas no caminho de save-success — mas isso quebra a separação de responsabilidades.

**Decisão pragmática:** Ao descartar e sair, não faz mal não mostrar `SetSuccess`. Ao salvar e sair, o guard já chama `notifier.Clear()` antes de `onProceder()`. O `onProceder` da QuitOperation faz:

```go
func() tea.Cmd {
	return tea.Quit
},
```

E o `SetSuccess("Cofre salvo.")` é removido da `QuitOperation` refatorada. Para manter o comportamento anterior (mostrar "Cofre salvo." antes de sair), o guard teria que distinguir save vs discard — mas isso complica a API. Como a aplicação sai imediatamente após, o `SetSuccess` é efêmero e o usuário mal o vê. Removê-lo é aceitável.

**Atualizar o teste `TestQuitOperation_Update_Saving_Sucesso_EmiteQuit`**: trocar a asserção de `SetSuccess` para `Clear` (que é o que o guard faz antes de chamar onProceder):

```go
// na verificação:
if n.lastMethod != "Clear" {
    t.Errorf("após salvar OK: esperado Clear, obteve %q", n.lastMethod)
}
```

- [ ] **Step 4: Rodar os testes**

```
go test ./internal/tui/operation/... -run TestQuitOperation
```
Expected: PASS

- [ ] **Step 5: Rodar todos os testes do pacote**

```
go test ./internal/tui/operation/...
```
Expected: PASS

- [ ] **Step 6: Commit**

```
git add internal/tui/operation/quit_operation.go internal/tui/operation/quit_operation_test.go
git commit -m "refactor: QuitOperation usa guardCofreAlterado no Fluxo 5"
```

---

## Task 5: Adicionar `Shortcuts.NewVault` e `Shortcuts.OpenVault` em `keys.go`

**Files:**
- Modify: `internal/tui/design/keys.go`

- [ ] **Step 1: Adicionar campos ao struct `Shortcuts`**

Em `internal/tui/design/keys.go`, substituir o bloco `var Shortcuts = struct {` (linhas 93-110) por:

```go
// Shortcuts contém os atalhos globais do design system, ativos em qualquer contexto da aplicação.
var Shortcuts = struct {
	// Help abre e fecha o diálogo de ajuda.
	Help Key
	// ThemeToggle alterna entre os temas Tokyo Night e Cyberpunk.
	// Não é exibido na barra de comandos.
	ThemeToggle Key
	// Quit sai da aplicação (com confirmação quando há alterações não salvas).
	Quit Key
	// LockVault bloqueia o cofre imediatamente, descartando alterações sem confirmação.
	// O atalho complexo (⌃!⇧Q) é intencional para evitar acionamento acidental.
	LockVault Key
	// NewVault inicia o fluxo de criação de novo cofre (Fluxo 2).
	NewVault Key
	// OpenVault inicia o fluxo de abertura de cofre existente (Fluxo 1).
	OpenVault Key
}{
	Help:        Keys.F1,
	ThemeToggle: Keys.F12,
	Quit:        WithCtrl(Letter('q')),
	LockVault:   WithCtrl(WithAlt(WithShift(Letter('q')))),
	NewVault:    WithCtrl(Letter('n')),
	OpenVault:   WithCtrl(Letter('o')),
}
```

- [ ] **Step 2: Compilar para verificar**

```
go build ./...
```
Expected: sem erros

- [ ] **Step 3: Commit**

```
git add internal/tui/design/keys.go
git commit -m "feat: adicionar atalhos Ctrl+N (NewVault) e Ctrl+O (OpenVault)"
```

---

## Task 6: Implementar `CriarCofreOperation` (Fluxo 2)

**Files:**
- Create: `internal/tui/operation/criar_cofre.go`
- Create: `internal/tui/operation/criar_cofre_test.go`

- [ ] **Step 1: Escrever os testes**

Criar `internal/tui/operation/criar_cofre_test.go`:

```go
package operation

import (
	"errors"
	"testing"

	"github.com/useful-toys/abditum/internal/tui"
)

// --- Init ---

func TestCriarCofre_Init_SemCofre_AbreFilePicker(t *testing.T) {
	op := newCriarCofreOperationFromSaver(&stubNotifier{}, nil, "")
	// guard com saver nil → onProceder() imediato → emite criarAvancaMsg
	cmd := op.Init()
	msg := execCmd(cmd)
	// Deve abrir o FilePicker via OpenModalMsg
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Init sem cofre: esperado OpenModalMsg (FilePicker), obteve %T", msg)
	}
}

func TestCriarCofre_Init_ComCaminhoInicial_AbrePasswordModal(t *testing.T) {
	op := newCriarCofreOperationFromSaver(&stubNotifier{}, nil, "/tmp/cofre.abditum")
	cmd := op.Init()
	msg := execCmd(cmd)
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Init com caminho: esperado OpenModalMsg (PasswordCreate), obteve %T", msg)
	}
}

func TestCriarCofre_Init_CofreAlterado_AbreGuardModal(t *testing.T) {
	op := newCriarCofreOperationFromSaver(&stubNotifier{}, &stubManager{isModified: true}, "")
	cmd := op.Init()
	msg := execCmd(cmd)
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Init cofre alterado: esperado OpenModalMsg (guard), obteve %T", msg)
	}
}

// --- Update: criação ---

func TestCriarCofre_Update_CriandoEstado_SetsBusy(t *testing.T) {
	n := &stubNotifier{}
	op := newCriarCofreOperationFromSaver(n, nil, "/tmp/cofre.abditum")
	op.caminho = "/tmp/cofre.abditum"
	op.senha = []byte("SenhaForte123!")

	cmd := op.Update(criarAvancaMsg{estado: criandoCriando})
	if n.lastMethod != "SetBusy" {
		t.Errorf("criandoCriando: esperado SetBusy, obteve %q", n.lastMethod)
	}
	if cmd == nil {
		t.Error("criandoCriando: esperado cmd não-nil")
	}
}

func TestCriarCofre_Update_ResultMsg_Falha_EmiteSetErrorECompleta(t *testing.T) {
	n := &stubNotifier{}
	op := newCriarCofreOperationFromSaver(n, nil, "")

	cmd := op.Update(criarCofreResultMsg{err: errors.New("falha ao criar")})
	if n.lastMethod != "SetError" {
		t.Errorf("falha: esperado SetError, obteve %q", n.lastMethod)
	}
	msg := execCmd(cmd)
	if _, ok := msg.(tui.OperationCompletedMsg); !ok {
		t.Errorf("falha: esperado OperationCompletedMsg, obteve %T", msg)
	}
}

func TestCriarCofre_Update_ResultMsg_Sucesso_EmiteVaultOpened(t *testing.T) {
	n := &stubNotifier{}
	op := newCriarCofreOperationFromSaver(n, nil, "")

	cmd := op.Update(criarCofreResultMsg{err: nil})
	if n.lastMethod != "Clear" {
		t.Errorf("sucesso: esperado Clear, obteve %q", n.lastMethod)
	}
	if cmd == nil {
		t.Error("sucesso: esperado cmd não-nil")
	}
}

func TestCriarCofre_Update_GuardMsg_DelegaAoGuard(t *testing.T) {
	n := &stubNotifier{}
	op := newCriarCofreOperationFromSaver(n, &stubManager{isModified: true}, "")
	op.Init() // ativa o guard

	// Enviar guardDiscardMsg — deve chamar onProceder do guard (→ abrirFilePicker)
	cmd := op.Update(guardDiscardMsg{})
	msg := execCmd(cmd)
	// onProceder emite criarAvancaMsg{estado: criandoInformandoCaminho},
	// que é processado no Update como abrirFilePicker → OpenModalMsg
	if _, ok := msg.(criarAvancaMsg); ok {
		// A mensagem é processada novamente pelo Update
		cmd2 := op.Update(msg)
		msg2 := execCmd(cmd2)
		if _, ok := msg2.(tui.OpenModalMsg); !ok {
			t.Errorf("guard discard → picker: esperado OpenModalMsg, obteve %T", msg2)
		}
	}
}

func TestCriarCofre_Update_MensagemDesconhecida_RetornaNil(t *testing.T) {
	op := newCriarCofreOperationFromSaver(&stubNotifier{}, nil, "")
	type outraMsg struct{}
	if cmd := op.Update(outraMsg{}); cmd != nil {
		t.Error("Update(outraMsg): esperado nil")
	}
}
```

- [ ] **Step 2: Rodar os testes para confirmar que falham**

```
go test ./internal/tui/operation/... -run TestCriarCofre
```
Expected: FAIL

- [ ] **Step 3: Criar `criar_cofre.go`**

```go
package operation

import (
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/crypto"
	"github.com/useful-toys/abditum/internal/storage"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/vault"
)

// criarEstado representa em qual etapa do fluxo de criação estamos.
type criarEstado int

const (
	criandoInformandoCaminho      criarEstado = iota // FilePicker modo Save
	criandoConfirmandoSobrescrita                    // ConfirmModal "arquivo já existe"
	criandoInformandoSenha                           // PasswordCreateModal
	criandoAvaliacaoSenhaFraca                       // ConfirmModal "senha fraca"
	criandoCriando                                   // IO assíncrono
)

// criarAvancaMsg é a mensagem interna de transição de estado.
type criarAvancaMsg struct {
	estado criarEstado
}

// criarCofreResultMsg carrega o resultado da criação do cofre.
type criarCofreResultMsg struct {
	manager *vault.Manager
	err     error
}

// CriarCofreOperation implementa o Fluxo 2 (Criar Novo Cofre).
//
// Se caminhoInicial != "", pula o guard e o picker e vai direto para a entrada de senha.
// Se caminhoInicial == "", executa o fluxo completo via GUI.
type CriarCofreOperation struct {
	notifier tui.MessageController
	saver    vaultSaver
	guard    *guardCofreAlterado // não-nil apenas no fluxo GUI com cofre alterado
	caminho  string             // caminho destino selecionado
	senha    []byte             // senha informada
}

// NewCriarCofreOperation cria a operation com o vault.Manager concreto.
// manager pode ser nil quando nenhum cofre está carregado.
// caminhoInicial pode ser "" para fluxo GUI completo.
func NewCriarCofreOperation(
	notifier tui.MessageController,
	manager *vault.Manager,
	caminhoInicial string,
) *CriarCofreOperation {
	var saver vaultSaver
	if manager != nil {
		saver = manager
	}
	return &CriarCofreOperation{
		notifier: notifier,
		saver:    saver,
		caminho:  caminhoInicial,
	}
}

// newCriarCofreOperationFromSaver é usada nos testes (mesmo pacote).
func newCriarCofreOperationFromSaver(
	notifier tui.MessageController,
	saver vaultSaver,
	caminhoInicial string,
) *CriarCofreOperation {
	return &CriarCofreOperation{
		notifier: notifier,
		saver:    saver,
		caminho:  caminhoInicial,
	}
}

// Init inicia o fluxo de criação.
func (c *CriarCofreOperation) Init() tea.Cmd {
	if c.caminho != "" {
		// Entrada via CLI: pular guard e picker, ir direto para senha
		return c.abrirModalSenha()
	}
	// Fluxo GUI completo: verificar cofre alterado primeiro
	c.guard = novoGuardCofreAlterado(
		c.notifier,
		c.saver,
		func() tea.Cmd {
			return func() tea.Msg { return criarAvancaMsg{estado: criandoInformandoCaminho} }
		},
		func() tea.Cmd { return tui.OperationCompleted() },
	)
	return c.guard.Init()
}

// Update trata as mensagens internas da CriarCofreOperation.
func (c *CriarCofreOperation) Update(msg tea.Msg) tea.Cmd {
	switch m := msg.(type) {
	// Mensagens do guard — delegar se guard ativo
	case guardSaveMsg, guardSaveResultMsg, guardDiscardMsg, guardCancelMsg:
		if c.guard != nil {
			return c.guard.Update(msg)
		}
		return nil

	case criarAvancaMsg:
		switch m.estado {
		case criandoInformandoCaminho:
			return c.abrirFilePicker()
		case criandoConfirmandoSobrescrita:
			return c.abrirModalSobrescrita()
		case criandoInformandoSenha:
			return c.abrirModalSenha()
		case criandoAvaliacaoSenhaFraca:
			return c.abrirModalSenhaFraca()
		case criandoCriando:
			c.notifier.SetBusy("Criando cofre...")
			caminho := c.caminho
			senha := c.senha
			return func() tea.Msg {
				cofre := vault.NovoCofre()
				if err := cofre.InicializarConteudoPadrao(); err != nil {
					return criarCofreResultMsg{err: err}
				}
				repo := storage.NewFileRepositoryForCreate(caminho, senha)
				manager := vault.NewManager(cofre, repo)
				if err := manager.Salvar(false); err != nil {
					return criarCofreResultMsg{err: err}
				}
				return criarCofreResultMsg{manager: manager}
			}
		}

	case criarCofreResultMsg:
		if m.err != nil {
			c.notifier.SetError(m.err.Error())
			return tui.OperationCompleted()
		}
		c.notifier.Clear()
		c.notifier.SetSuccess("Cofre criado.")
		return func() tea.Msg { return tui.VaultOpenedMsg{Manager: m.manager} }
	}

	return nil
}

// abrirFilePicker abre o FilePicker no modo Save com extensão .abditum.
func (c *CriarCofreOperation) abrirFilePicker() tea.Cmd {
	return tui.OpenModal(modal.NewFilePicker(modal.FilePickerOptions{
		Mode:      modal.FilePickerSave,
		Extension: ".abditum",
		Messages:  c.notifier,
		OnResult: func(path string) tea.Cmd {
			if path == "" {
				return tui.OperationCompleted()
			}
			c.caminho = path
			if fileExists(path) {
				return func() tea.Msg { return criarAvancaMsg{estado: criandoConfirmandoSobrescrita} }
			}
			return func() tea.Msg { return criarAvancaMsg{estado: criandoInformandoSenha} }
		},
	}))
}

// abrirModalSobrescrita abre o modal de confirmação de sobrescrita.
func (c *CriarCofreOperation) abrirModalSobrescrita() tea.Cmd {
	return tui.OpenModal(modal.NewConfirmModal(
		"Arquivo existente",
		"O arquivo já existe. Deseja sobrescrever?",
		[]modal.ModalOption{
			{
				Keys:  []design.Key{design.Keys.Enter},
				Label: "Sobrescrever",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return criarAvancaMsg{estado: criandoInformandoSenha}
					})
				},
			},
			{
				Keys:  []design.Key{design.Keys.Esc},
				Label: "Outro caminho",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return criarAvancaMsg{estado: criandoInformandoCaminho}
					})
				},
			},
		},
	))
}

// abrirModalSenha abre o PasswordCreateModal.
func (c *CriarCofreOperation) abrirModalSenha() tea.Cmd {
	return tui.OpenModal(modal.NewPasswordCreateModal(
		c.notifier,
		func(password []byte) tea.Cmd {
			c.senha = password
			if crypto.EvaluatePasswordStrength(password) == crypto.StrengthWeak {
				return tea.Batch(tui.CloseModal(), func() tea.Msg {
					return criarAvancaMsg{estado: criandoAvaliacaoSenhaFraca}
				})
			}
			return tea.Batch(tui.CloseModal(), func() tea.Msg {
				return criarAvancaMsg{estado: criandoCriando}
			})
		},
		func() tea.Cmd {
			return tea.Batch(tui.CloseModal(), func() tea.Msg {
				return criarAvancaMsg{estado: criandoInformandoCaminho}
			})
		},
	))
}

// abrirModalSenhaFraca abre o modal de aviso de senha fraca.
func (c *CriarCofreOperation) abrirModalSenhaFraca() tea.Cmd {
	return tui.OpenModal(modal.NewConfirmModal(
		"Senha fraca",
		"A senha informada é fraca. Deseja prosseguir assim mesmo?",
		[]modal.ModalOption{
			{
				Keys:  []design.Key{design.Keys.Enter},
				Label: "Prosseguir",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return criarAvancaMsg{estado: criandoCriando}
					})
				},
			},
			{
				Keys:  []design.Key{design.Keys.Esc},
				Label: "Revisar",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return criarAvancaMsg{estado: criandoInformandoSenha}
					})
				},
			},
		},
	))
}

// fileExists reporta se o caminho aponta para um arquivo existente.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
```

- [ ] **Step 4: Rodar os testes**

```
go test ./internal/tui/operation/... -run TestCriarCofre
```
Expected: PASS

- [ ] **Step 5: Rodar todos os testes do pacote**

```
go test ./internal/tui/operation/...
```
Expected: PASS

- [ ] **Step 6: Commit**

```
git add internal/tui/operation/criar_cofre.go internal/tui/operation/criar_cofre_test.go
git commit -m "feat: implementar CriarCofreOperation (Fluxo 2)"
```

---

## Task 7: Implementar `AbrirCofreOperation` (Fluxo 1)

**Files:**
- Create: `internal/tui/operation/abrir_cofre.go`
- Create: `internal/tui/operation/abrir_cofre_test.go`

- [ ] **Step 1: Escrever os testes**

Criar `internal/tui/operation/abrir_cofre_test.go`:

```go
package operation

import (
	"errors"
	"testing"

	"github.com/useful-toys/abditum/internal/crypto"
	"github.com/useful-toys/abditum/internal/storage"
	"github.com/useful-toys/abditum/internal/tui"
)

// --- Init ---

func TestAbrirCofre_Init_SemCofre_AbreFilePicker(t *testing.T) {
	op := newAbrirCofreOperationFromSaver(&stubNotifier{}, nil, "")
	cmd := op.Init()
	msg := execCmd(cmd)
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Init sem cofre: esperado OpenModalMsg (FilePicker), obteve %T", msg)
	}
}

func TestAbrirCofre_Init_CofreAlterado_AbreGuardModal(t *testing.T) {
	op := newAbrirCofreOperationFromSaver(&stubNotifier{}, &stubManager{isModified: true}, "")
	cmd := op.Init()
	msg := execCmd(cmd)
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Init cofre alterado: esperado OpenModalMsg (guard), obteve %T", msg)
	}
}

func TestAbrirCofre_Init_ComCaminhoInicial_HeaderInvalido_Erro(t *testing.T) {
	n := &stubNotifier{}
	// Caminho inexistente — ValidateHeader retornará erro de IO
	op := newAbrirCofreOperationFromSaver(n, nil, "/caminho/inexistente.abditum")
	cmd := op.Init()
	msg := execCmd(cmd)
	if n.lastMethod != "SetError" {
		t.Errorf("header inválido: esperado SetError, obteve %q", n.lastMethod)
	}
	if _, ok := msg.(tui.OperationCompletedMsg); !ok {
		t.Errorf("header inválido: esperado OperationCompletedMsg, obteve %T", msg)
	}
}

// --- Update: abertura ---

func TestAbrirCofre_Update_AbrindoEstado_SetsBusy(t *testing.T) {
	n := &stubNotifier{}
	op := newAbrirCofreOperationFromSaver(n, nil, "")
	op.caminho = "/tmp/cofre.abditum"
	op.senha = []byte("senha")

	cmd := op.Update(abrirAvancaMsg{estado: abrindoAbrindo})
	if n.lastMethod != "SetBusy" {
		t.Errorf("abrindoAbrindo: esperado SetBusy, obteve %q", n.lastMethod)
	}
	if cmd == nil {
		t.Error("abrindoAbrindo: esperado cmd não-nil")
	}
}

func TestAbrirCofre_Update_ResultMsg_SenhaErrada_VoltaParaSenha(t *testing.T) {
	n := &stubNotifier{}
	op := newAbrirCofreOperationFromSaver(n, nil, "")

	cmd := op.Update(abrirCofreResultMsg{err: crypto.ErrAuthFailed})
	if n.lastMethod != "SetError" {
		t.Errorf("senha errada: esperado SetError, obteve %q", n.lastMethod)
	}
	// Deve reabrir o modal de senha
	msg := execCmd(cmd)
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("senha errada: esperado OpenModalMsg (senha), obteve %T", msg)
	}
}

func TestAbrirCofre_Update_ResultMsg_Corrompido_VoltaParaCaminho(t *testing.T) {
	n := &stubNotifier{}
	op := newAbrirCofreOperationFromSaver(n, nil, "")

	cmd := op.Update(abrirCofreResultMsg{err: storage.ErrCorrupted})
	if n.lastMethod != "SetError" {
		t.Errorf("corrompido: esperado SetError, obteve %q", n.lastMethod)
	}
	msg := execCmd(cmd)
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("corrompido: esperado OpenModalMsg (picker), obteve %T", msg)
	}
}

func TestAbrirCofre_Update_ResultMsg_Sucesso_EmiteVaultOpened(t *testing.T) {
	n := &stubNotifier{}
	op := newAbrirCofreOperationFromSaver(n, nil, "")

	cmd := op.Update(abrirCofreResultMsg{err: nil})
	if n.lastMethod != "Clear" {
		t.Errorf("sucesso: esperado Clear, obteve %q", n.lastMethod)
	}
	if cmd == nil {
		t.Error("sucesso: esperado cmd não-nil")
	}
}

func TestAbrirCofre_Update_MensagemDesconhecida_RetornaNil(t *testing.T) {
	op := newAbrirCofreOperationFromSaver(&stubNotifier{}, nil, "")
	type outraMsg struct{}
	if cmd := op.Update(outraMsg{}); cmd != nil {
		t.Error("Update(outraMsg): esperado nil")
	}
}
```

- [ ] **Step 2: Rodar os testes para confirmar que falham**

```
go test ./internal/tui/operation/... -run TestAbrirCofre
```
Expected: FAIL

- [ ] **Step 3: Criar `abrir_cofre.go`**

```go
package operation

import (
	"errors"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/crypto"
	"github.com/useful-toys/abditum/internal/storage"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/vault"
)

// abrirEstado representa em qual etapa do fluxo de abertura estamos.
type abrirEstado int

const (
	abrindoInformandoCaminho abrirEstado = iota // FilePicker modo Open
	abrindoInformandoSenha                      // PasswordEntryModal
	abrindoAbrindo                              // IO assíncrono
)

// abrirAvancaMsg é a mensagem interna de transição de estado.
type abrirAvancaMsg struct {
	estado abrirEstado
}

// abrirCofreResultMsg carrega o resultado da abertura do cofre.
type abrirCofreResultMsg struct {
	manager *vault.Manager
	err     error
}

// AbrirCofreOperation implementa o Fluxo 1 (Abrir Cofre Existente).
//
// Se caminhoInicial != "", valida o arquivo e vai direto para a entrada de senha.
// Se caminhoInicial == "", executa o fluxo completo via GUI.
type AbrirCofreOperation struct {
	notifier tui.MessageController
	saver    vaultSaver
	guard    *guardCofreAlterado // não-nil apenas no fluxo GUI com cofre alterado
	caminho  string             // caminho do cofre a abrir
	senha    []byte             // senha informada
}

// NewAbrirCofreOperation cria a operation com o vault.Manager concreto.
// manager pode ser nil quando nenhum cofre está carregado.
// caminhoInicial pode ser "" para fluxo GUI completo.
func NewAbrirCofreOperation(
	notifier tui.MessageController,
	manager *vault.Manager,
	caminhoInicial string,
) *AbrirCofreOperation {
	var saver vaultSaver
	if manager != nil {
		saver = manager
	}
	return &AbrirCofreOperation{
		notifier: notifier,
		saver:    saver,
		caminho:  caminhoInicial,
	}
}

// newAbrirCofreOperationFromSaver é usada nos testes (mesmo pacote).
func newAbrirCofreOperationFromSaver(
	notifier tui.MessageController,
	saver vaultSaver,
	caminhoInicial string,
) *AbrirCofreOperation {
	return &AbrirCofreOperation{
		notifier: notifier,
		saver:    saver,
		caminho:  caminhoInicial,
	}
}

// Init inicia o fluxo de abertura.
func (a *AbrirCofreOperation) Init() tea.Cmd {
	if a.caminho != "" {
		// Entrada via CLI: validar magic/versão antes de pedir senha
		if err := storage.ValidateHeader(a.caminho); err != nil {
			a.notifier.SetError(erroDeAberturaCategoria(err))
			return tui.OperationCompleted()
		}
		return a.abrirModalSenha()
	}
	// Fluxo GUI completo: verificar cofre alterado primeiro
	a.guard = novoGuardCofreAlterado(
		a.notifier,
		a.saver,
		func() tea.Cmd {
			return func() tea.Msg { return abrirAvancaMsg{estado: abrindoInformandoCaminho} }
		},
		func() tea.Cmd { return tui.OperationCompleted() },
	)
	return a.guard.Init()
}

// Update trata as mensagens internas da AbrirCofreOperation.
func (a *AbrirCofreOperation) Update(msg tea.Msg) tea.Cmd {
	switch m := msg.(type) {
	// Mensagens do guard — delegar se guard ativo
	case guardSaveMsg, guardSaveResultMsg, guardDiscardMsg, guardCancelMsg:
		if a.guard != nil {
			return a.guard.Update(msg)
		}
		return nil

	case abrirAvancaMsg:
		switch m.estado {
		case abrindoInformandoCaminho:
			return a.abrirFilePicker()
		case abrindoInformandoSenha:
			return a.abrirModalSenha()
		case abrindoAbrindo:
			a.notifier.SetBusy("Abrindo cofre...")
			caminho := a.caminho
			senha := a.senha
			return func() tea.Msg {
				repo := storage.NewFileRepositoryForOpen(caminho, senha)
				cofre, err := repo.Carregar()
				if err != nil {
					return abrirCofreResultMsg{err: err}
				}
				manager := vault.NewManager(cofre, repo)
				return abrirCofreResultMsg{manager: manager}
			}
		}

	case abrirCofreResultMsg:
		if m.err != nil {
			a.notifier.SetError(erroDeAberturaCategoria(m.err))
			if isErrAutenticacao(m.err) {
				// Senha errada: voltar ao modal de senha
				return a.abrirModalSenha()
			}
			// Integridade/formato: voltar ao picker
			return a.abrirFilePicker()
		}
		a.notifier.Clear()
		a.notifier.SetSuccess("Cofre aberto.")
		return func() tea.Msg { return tui.VaultOpenedMsg{Manager: m.manager} }
	}

	_ = msg
	return nil
}

// abrirFilePicker abre o FilePicker no modo Open.
func (a *AbrirCofreOperation) abrirFilePicker() tea.Cmd {
	return tui.OpenModal(modal.NewFilePicker(modal.FilePickerOptions{
		Mode:      modal.FilePickerOpen,
		Extension: ".abditum",
		Messages:  a.notifier,
		OnResult: func(path string) tea.Cmd {
			if path == "" {
				return tui.OperationCompleted()
			}
			// Validar magic/versão antes de pedir senha
			if err := storage.ValidateHeader(path); err != nil {
				a.notifier.SetError(erroDeAberturaCategoria(err))
				return a.abrirFilePicker()
			}
			a.caminho = path
			return func() tea.Msg { return abrirAvancaMsg{estado: abrindoInformandoSenha} }
		},
	}))
}

// abrirModalSenha abre o PasswordEntryModal.
func (a *AbrirCofreOperation) abrirModalSenha() tea.Cmd {
	return tui.OpenModal(modal.NewPasswordEntryModal(
		a.notifier,
		func(password []byte) tea.Cmd {
			a.senha = password
			return tea.Batch(tui.CloseModal(), func() tea.Msg {
				return abrirAvancaMsg{estado: abrindoAbrindo}
			})
		},
		func() tea.Cmd {
			return tea.Batch(tui.CloseModal(), func() tea.Msg {
				return abrirAvancaMsg{estado: abrindoInformandoCaminho}
			})
		},
	))
}

// erroDeAberturaCategoria converte o erro técnico em mensagem de categoria
// sem revelar a causa exata (por segurança).
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

// isErrAutenticacao reporta se o erro é de autenticação (senha errada).
func isErrAutenticacao(err error) bool {
	return errors.Is(err, crypto.ErrAuthFailed)
}
```

- [ ] **Step 4: Rodar os testes**

```
go test ./internal/tui/operation/... -run TestAbrirCofre
```
Expected: PASS

- [ ] **Step 5: Rodar todos os testes do pacote**

```
go test ./internal/tui/operation/...
```
Expected: PASS

- [ ] **Step 6: Commit**

```
git add internal/tui/operation/abrir_cofre.go internal/tui/operation/abrir_cofre_test.go
git commit -m "feat: implementar AbrirCofreOperation (Fluxo 1)"
```

---

## Task 8: Integração — `keys.go`, `setup.go`, `root.go` e `main.go`

**Files:**
- Modify: `internal/tui/root.go`
- Modify: `cmd/abditum/setup.go`
- Modify: `cmd/abditum/main.go`

- [ ] **Step 1: Adicionar `WithInitialVault` option em `root.go`**

Em `internal/tui/root.go`, após `WithVersion` (linha 427-431), adicionar:

```go
// WithInitialVault define o caminho de cofre a ser aberto/criado na inicialização.
// Se o arquivo existir, dispara AbrirCofreOperation; se não existir mas o diretório
// pai existir, dispara CriarCofreOperation. Caso contrário, é ignorado.
func WithInitialVault(path string) RootModelOption {
	return func(m *RootModel) {
		m.initialVaultPath = path
	}
}
```

Adicionar o campo `initialVaultPath string` ao struct `RootModel` (após `version string`, linha 70):

```go
	// initialVaultPath é o caminho de cofre passado via --vault na CLI.
	// Se preenchido, Init() dispara a operation correspondente.
	initialVaultPath string
```

Modificar `Init()` (linhas 411-413) para disparar a operation:

```go
func (r *RootModel) Init() tea.Cmd {
	cmds := []tea.Cmd{tickCmd()}
	if r.initialVaultPath != "" {
		cmds = append(cmds, r.buildInitialVaultCmd())
	}
	return tea.Batch(cmds...)
}

// buildInitialVaultCmd cria o comando que dispara a operation de abrir ou criar
// cofre baseado no caminho passado via --vault.
func (r *RootModel) buildInitialVaultCmd() tea.Cmd {
	path := r.initialVaultPath
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		// Arquivo existe → Fluxo 1 (Abrir)
		return StartOperation(operation.NewAbrirCofreOperation(r.MessageController(), nil, path))
	}
	if os.IsNotExist(err) {
		dir := filepath.Dir(path)
		if dirInfo, dirErr := os.Stat(dir); dirErr == nil && dirInfo.IsDir() {
			// Arquivo não existe, dir pai existe → Fluxo 2 (Criar)
			return StartOperation(operation.NewCriarCofreOperation(r.MessageController(), nil, path))
		}
	}
	// Dir pai não existe ou outro erro → ignorar, tela normal
	return nil
}
```

Adicionar os imports necessários em `root.go`:

```go
import (
	"os"
	"path/filepath"
	// ... imports existentes ...
	"github.com/useful-toys/abditum/internal/tui/operation"
)
```

**Nota sobre import cycle:** `tui/root.go` importando `tui/operation` pode causar ciclo se `operation` já importa `tui`. Verificar: `operation` importa `tui` (para `tui.MessageController`, `tui.OpenModal`, etc.). Isso cria um import cycle: `tui → tui/operation → tui`.

**Solução:** Em vez de importar `operation` diretamente em `root.go`, usar o padrão de closure — `buildInitialVaultCmd` é uma função que retorna `tea.Cmd`, e a criação da operation é feita em `main.go` via `RootModelOption`. Alternativa: mover `buildInitialVaultCmd` para `main.go` passando uma `tea.Cmd` diretamente como option:

```go
// WithInitialCommand define um comando a ser emitido junto com Init().
func WithInitialCommand(cmd tea.Cmd) RootModelOption {
	return func(m *RootModel) {
		m.initialCmd = cmd
	}
}
```

Campo no struct:
```go
	// initialCmd é um comando opcional a ser emitido junto com Init().
	initialCmd tea.Cmd
```

`Init()`:
```go
func (r *RootModel) Init() tea.Cmd {
	cmds := []tea.Cmd{tickCmd()}
	if r.initialCmd != nil {
		cmds = append(cmds, r.initialCmd)
	}
	return tea.Batch(cmds...)
}
```

E a lógica de decisão fica em `main.go`. Essa abordagem evita o import cycle.

**Usar esta abordagem (WithInitialCommand).**

Em `internal/tui/root.go`:

Adicionar campo ao struct RootModel (após `version string`):
```go
	// initialCmd é um comando opcional a ser emitido junto com Init().
	// Usado para disparar operations na inicialização (ex: --vault via CLI).
	initialCmd tea.Cmd
```

Adicionar option após `WithVersion`:
```go
// WithInitialCommand define um comando a ser emitido junto com Init().
// Usado para disparar uma operation na inicialização da aplicação.
func WithInitialCommand(cmd tea.Cmd) RootModelOption {
	return func(m *RootModel) {
		m.initialCmd = cmd
	}
}
```

Modificar `Init()`:
```go
func (r *RootModel) Init() tea.Cmd {
	cmds := []tea.Cmd{tickCmd()}
	if r.initialCmd != nil {
		cmds = append(cmds, r.initialCmd)
	}
	return tea.Batch(cmds...)
}
```

- [ ] **Step 2: Adicionar ações Ctrl+N e Ctrl+O em `setup.go`**

Em `cmd/abditum/setup.go`, na função `setupApplication`, após a ação de `Sair` (linhas 58-68), antes da `FakeOperation` (linha 69), adicionar:

```go
		{
			Keys:        []design.Key{design.Shortcuts.NewVault},
			Label:       "Criar cofre",
			Description: "Cria um novo cofre protegido por senha.",
			GroupID:     "app",
			Priority:    30,
			Visible:     true,
			AvailableWhen: func(app actions.AppState, _ actions.ChildView) bool {
				return app.Manager() == nil
			},
			OnExecute: func() tea.Cmd {
				return tui.StartOperation(operation.NewCriarCofreOperation(r.MessageController(), r.Manager(), ""))
			},
		},
		{
			Keys:        []design.Key{design.Shortcuts.OpenVault},
			Label:       "Abrir cofre",
			Description: "Abre um cofre existente a partir de um arquivo.",
			GroupID:     "app",
			Priority:    31,
			Visible:     true,
			AvailableWhen: func(app actions.AppState, _ actions.ChildView) bool {
				return app.Manager() == nil
			},
			OnExecute: func() tea.Cmd {
				return tui.StartOperation(operation.NewAbrirCofreOperation(r.MessageController(), r.Manager(), ""))
			},
		},
```

- [ ] **Step 3: Implementar lógica de `--vault` em `main.go`**

Em `cmd/abditum/main.go`, na função `run()`, substituir a criação do root (linhas 41-43) por:

```go
	opts := []tui.RootModelOption{tui.WithVersion(version)}

	// Se um caminho de cofre foi passado via --vault, decidir qual fluxo disparar
	if vaultPath != "" {
		if cmd := buildVaultCmd(vaultPath); cmd != nil {
			opts = append(opts, tui.WithInitialCommand(cmd))
		}
	}

	root := tui.NewRootModel(opts...)
```

Adicionar a função `buildVaultCmd` no mesmo arquivo:

```go
// buildVaultCmd decide qual operation disparar com base no caminho --vault.
//   - Arquivo existe → Fluxo 1 (Abrir)
//   - Arquivo não existe, dir pai existe → Fluxo 2 (Criar)
//   - Senão → nil (tela normal)
func buildVaultCmd(vaultPath string) tea.Cmd {
	info, err := os.Stat(vaultPath)
	if err == nil && !info.IsDir() {
		return tui.StartOperation(operation.NewAbrirCofreOperation(nil, nil, vaultPath))
	}
	if os.IsNotExist(err) {
		dir := filepath.Dir(vaultPath)
		if dirInfo, dirErr := os.Stat(dir); dirErr == nil && dirInfo.IsDir() {
			return tui.StartOperation(operation.NewCriarCofreOperation(nil, nil, vaultPath))
		}
	}
	return nil
}
```

Adicionar imports em `main.go`:

```go
import (
	"path/filepath"
	// ... imports existentes ...
	"github.com/useful-toys/abditum/internal/tui/operation"
)
```

**Nota:** O `notifier` é passado como `nil` no `buildVaultCmd`. Isso é um problema — `NewAbrirCofreOperation` e `NewCriarCofreOperation` usam o notifier. Mas no momento da criação em `main.go`, o root model ainda não foi criado (pois o command é passado como option).

**Solução:** Usar `WithInitialCommand` com uma closure que será executada por `tea.Batch` — nesse momento o root já existirá. Mas a operation precisa do notifier no construtor, não em `Init()`.

**Alternativa:** Passar `nil` como notifier no `buildVaultCmd`, e mudar `main.go` para criar o root primeiro, depois injetar o command:

```go
	root := tui.NewRootModel(tui.WithVersion(version))
	setupActions(root)

	// Se um caminho de cofre foi passado via --vault, disparar o fluxo correspondente
	if vaultPath != "" {
		if cmd := buildVaultCmd(vaultPath, root); cmd != nil {
			root.SetInitialCommand(cmd)
		}
	}
```

Adicionar método `SetInitialCommand` ao `RootModel` (em vez de `WithInitialCommand` option):

```go
// SetInitialCommand define um comando a ser emitido junto com Init().
// Deve ser chamado antes de iniciar o loop Tea.
func (r *RootModel) SetInitialCommand(cmd tea.Cmd) {
	r.initialCmd = cmd
}
```

E `buildVaultCmd` recebe o root:

```go
func buildVaultCmd(vaultPath string, root *tui.RootModel) tea.Cmd {
	info, err := os.Stat(vaultPath)
	if err == nil && !info.IsDir() {
		return tui.StartOperation(
			operation.NewAbrirCofreOperation(root.MessageController(), nil, vaultPath),
		)
	}
	if os.IsNotExist(err) {
		dir := filepath.Dir(vaultPath)
		if dirInfo, dirErr := os.Stat(dir); dirErr == nil && dirInfo.IsDir() {
			return tui.StartOperation(
				operation.NewCriarCofreOperation(root.MessageController(), nil, vaultPath),
			)
		}
	}
	return nil
}
```

**Usar esta abordagem (SetInitialCommand + buildVaultCmd com root).**

Atualizar o passo: em `internal/tui/root.go`, adicionar campo e método:

Campo (após `version string`):
```go
	// initialCmd é um comando opcional a ser emitido junto com Init().
	// Configurado via SetInitialCommand antes de iniciar o loop Tea.
	initialCmd tea.Cmd
```

Método:
```go
// SetInitialCommand define um comando a ser emitido junto com Init().
// Deve ser chamado antes de iniciar o loop Tea.
func (r *RootModel) SetInitialCommand(cmd tea.Cmd) {
	r.initialCmd = cmd
}
```

Modificar `Init()`:
```go
func (r *RootModel) Init() tea.Cmd {
	cmds := []tea.Cmd{tickCmd()}
	if r.initialCmd != nil {
		cmds = append(cmds, r.initialCmd)
	}
	return tea.Batch(cmds...)
}
```

E em `main.go`:
```go
func run() error {
	var vaultPath string
	flag.StringVar(&vaultPath, "vault", "", "Path to the Abditum vault file")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	defer clipboard.WriteAll("") //nolint:errcheck

	root := tui.NewRootModel(tui.WithVersion(version))
	setupActions(root)

	if vaultPath != "" {
		if cmd := buildVaultCmd(vaultPath, root); cmd != nil {
			root.SetInitialCommand(cmd)
		}
	}

	p := tea.NewProgram(root,
		tea.WithContext(ctx),
	)
	_, err := p.Run()
	return err
}

// buildVaultCmd decide qual operation disparar com base no caminho --vault.
func buildVaultCmd(vaultPath string, root *tui.RootModel) tea.Cmd {
	info, err := os.Stat(vaultPath)
	if err == nil && !info.IsDir() {
		return tui.StartOperation(
			operation.NewAbrirCofreOperation(root.MessageController(), nil, vaultPath),
		)
	}
	if os.IsNotExist(err) {
		dir := filepath.Dir(vaultPath)
		if dirInfo, dirErr := os.Stat(dir); dirErr == nil && dirInfo.IsDir() {
			return tui.StartOperation(
				operation.NewCriarCofreOperation(root.MessageController(), nil, vaultPath),
			)
		}
	}
	return nil
}
```

Adicionar imports: `"path/filepath"`, `"github.com/useful-toys/abditum/internal/tui/operation"`.

- [ ] **Step 4: Compilar para verificar**

```
go build ./...
```
Expected: sem erros

- [ ] **Step 5: Rodar todos os testes**

```
go test ./...
```
Expected: PASS

- [ ] **Step 6: Commit**

```
git add internal/tui/root.go internal/tui/design/keys.go cmd/abditum/setup.go cmd/abditum/main.go
git commit -m "feat: integrar Criar/Abrir Cofre via Ctrl+N/Ctrl+O e --vault na CLI"
```

---

## Checklist de Cobertura do Spec

| Requisito do Spec | Task |
|-------------------|------|
| `storage.ValidateHeader` | Task 1 |
| `storage.NewFileRepositoryForOpen` | Task 1 |
| `vaultSaver` movida para arquivo compartilhado | Task 2 |
| `guardCofreAlterado` com sub-fluxo de conflito externo | Task 3 |
| Refatoração da `QuitOperation` para usar guard | Task 4 |
| `Shortcuts.NewVault` e `Shortcuts.OpenVault` | Task 5 |
| `CriarCofreOperation` com `caminhoInicial` | Task 6 |
| `AbrirCofreOperation` com `caminhoInicial` e `ValidateHeader` | Task 7 |
| Ações Ctrl+N / Ctrl+O em `setup.go` | Task 8 |
| `--vault` CLI dispara operation correta | Task 8 |
| `RootModel.SetInitialCommand` + `Init()` | Task 8 |
| Avaliação de força de senha em `CriarCofreOperation` | Task 6 |
| Mensagens de erro por categoria em `AbrirCofreOperation` | Task 7 |
| `AbrirCofreOperation` usa `NewFileRepositoryForOpen` (não `ForCreate`) | Task 7 |
