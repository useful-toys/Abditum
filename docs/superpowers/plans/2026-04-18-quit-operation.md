# QuitOperation (ctrl Q) — Plano de Implementação

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Substituir o `tea.Quit` direto da action `ctrl Q` por uma `QuitOperation` que implementa os Fluxos 3, 4 e 5 do `fluxos.md`, incluindo detecção de modificação externa antes de salvar.

**Architecture:** Três camadas em ordem de dependência: (1) `internal/vault` — `RepositorioCofre` ganha `DetectarAlteracaoExterna()` e `Manager.Salvar` ganha flag `forcarSobrescrita`; (2) `internal/storage` — `FileRepository` implementa o novo método; (3) `internal/tui/operation` — `QuitOperation` orquestra os modais e chama o `Manager`.

**Tech Stack:** Go 1.26, Bubble Tea v2, pacotes internos `vault`, `storage`, `tui`.

**Spec:** `docs/superpowers/specs/2026-04-18-quit-operation-design.md`

---

## Mapa de arquivos

| Arquivo | Ação | Responsabilidade |
|---|---|---|
| `internal/vault/repository.go` | Modificar | Adicionar `DetectarAlteracaoExterna()` à interface |
| `internal/vault/errors.go` | Criar | Declarar `ErrModifiedExternally` |
| `internal/vault/manager.go` | Modificar | `Salvar(forcarSobrescrita bool) error` |
| `internal/vault/manager_test.go` | Modificar | Atualizar mock + testes de `Salvar` |
| `internal/storage/repository.go` | Modificar | Implementar `DetectarAlteracaoExterna()` |
| `internal/storage/storage_test.go` | Modificar | Testar `DetectarAlteracaoExterna()` via `FileRepository` |
| `internal/tui/operation/quit_operation.go` | Criar | `QuitOperation` — modais e lógica de saída |
| `internal/tui/operation/quit_operation_test.go` | Criar | Testes da `QuitOperation` |
| `cmd/abditum/setup.go` | Modificar | Trocar `tea.Quit` por `StartOperation(NewQuitOperation(...))` |

---

## Tarefa 1 — Declarar `ErrModifiedExternally`

**Arquivos:**
- Criar: `internal/vault/errors.go`
- Modificar: `internal/vault/manager_test.go` (linha ~8 — adicionar import `errors`)

- [ ] **Escrever o teste que verifica que o erro é comparável via `errors.Is`**

  Arquivo: `internal/vault/errors_test.go` (criar)

  ```go
  package vault_test

  import (
      "errors"
      "testing"

      "github.com/useful-toys/abditum/internal/vault"
  )

  func TestErrModifiedExternally_IsComparavel(t *testing.T) {
      wrapped := fmt.Errorf("vault.Salvar: %w", vault.ErrModifiedExternally)
      if !errors.Is(wrapped, vault.ErrModifiedExternally) {
          t.Error("errors.Is deveria reconhecer ErrModifiedExternally em erro encadeado")
      }
  }
  ```

  Adicionar `"fmt"` ao import do arquivo de teste.

- [ ] **Rodar o teste para confirmar que falha**

  ```
  CGO_ENABLED=0 go test ./internal/vault/... -run TestErrModifiedExternally -v
  ```

  Esperado: FAIL — `vault.ErrModifiedExternally undefined`

- [ ] **Criar `internal/vault/errors.go`**

  ```go
  package vault

  import "errors"

  // ErrModifiedExternally é retornado por Manager.Salvar quando o arquivo de cofre
  // foi modificado externamente desde o último Load ou Save, e forcarSobrescrita é false.
  var ErrModifiedExternally = errors.New("arquivo modificado externamente")
  ```

- [ ] **Rodar o teste para confirmar que passa**

  ```
  CGO_ENABLED=0 go test ./internal/vault/... -run TestErrModifiedExternally -v
  ```

  Esperado: PASS

- [ ] **Commit**

  ```
  git add internal/vault/errors.go internal/vault/errors_test.go
  git commit -m "feat(vault): declara ErrModifiedExternally"
  ```

---

## Tarefa 2 — Adicionar `DetectarAlteracaoExterna` à interface `RepositorioCofre`

**Arquivos:**
- Modificar: `internal/vault/repository.go`
- Modificar: `internal/vault/manager_test.go` (mock deve implementar o novo método)

- [ ] **Adicionar o método à interface**

  Arquivo: `internal/vault/repository.go`

  ```go
  package vault

  // RepositorioCofre defines the storage interface for vault persistence.
  type RepositorioCofre interface {
      // Salvar persists the vault to storage.
      Salvar(cofre *Cofre) error

      // Carregar loads a vault from storage.
      Carregar() (*Cofre, error)

      // DetectarAlteracaoExterna verifica se o arquivo de cofre foi modificado
      // por processo externo desde o último Salvar ou Carregar.
      // Retorna false sem erro se não houver baseline (cofre recém-criado).
      DetectarAlteracaoExterna() (bool, error)
  }
  ```

- [ ] **Rodar os testes para confirmar que o mock quebra a compilação**

  ```
  CGO_ENABLED=0 go test ./internal/vault/... -v
  ```

  Esperado: FAIL — `mockRepository does not implement RepositorioCofre`

- [ ] **Atualizar o mock em `internal/vault/manager_test.go`**

  Localizar o struct `mockRepository` (linhas ~8–21) e adicionar o novo método:

  ```go
  type mockRepository struct {
      salvarCalled                 bool
      salvarError                  error
      detectarAlteracaoExternaResp bool
      detectarAlteracaoExternaErr  error
  }

  func (m *mockRepository) Salvar(cofre *Cofre) error {
      m.salvarCalled = true
      return m.salvarError
  }

  func (m *mockRepository) Carregar() (*Cofre, error) {
      return nil, errors.New("not implemented")
  }

  func (m *mockRepository) DetectarAlteracaoExterna() (bool, error) {
      return m.detectarAlteracaoExternaResp, m.detectarAlteracaoExternaErr
  }
  ```

- [ ] **Rodar os testes para confirmar que compilam e passam**

  ```
  CGO_ENABLED=0 go test ./internal/vault/... -v
  ```

  Esperado: PASS (todos os testes existentes devem continuar passando)

- [ ] **Commit**

  ```
  git add internal/vault/repository.go internal/vault/manager_test.go
  git commit -m "feat(vault): adiciona DetectarAlteracaoExterna à interface RepositorioCofre"
  ```

---

## Tarefa 3 — Implementar `DetectarAlteracaoExterna` no `FileRepository`

**Arquivos:**
- Modificar: `internal/storage/repository.go`
- Modificar: `internal/storage/storage_test.go`

- [ ] **Escrever o teste**

  Em `internal/storage/storage_test.go`, adicionar ao final (antes do último `}`):

  ```go
  func TestFileRepository_DetectarAlteracaoExterna_SemAlteracao(t *testing.T) {
      dir := t.TempDir()
      path := filepath.Join(dir, "cofre.abditum")
      cofre := vault.NewCofre()
      repo := NewFileRepositoryForCreate(path, testPassword)
      if err := repo.Salvar(cofre); err != nil {
          t.Fatalf("Salvar: %v", err)
      }

      changed, err := repo.DetectarAlteracaoExterna()
      if err != nil {
          t.Fatalf("DetectarAlteracaoExterna: %v", err)
      }
      if changed {
          t.Error("esperado false (sem alteração externa), obteve true")
      }
  }

  func TestFileRepository_DetectarAlteracaoExterna_ComAlteracao(t *testing.T) {
      dir := t.TempDir()
      path := filepath.Join(dir, "cofre.abditum")
      cofre := vault.NewCofre()
      repo := NewFileRepositoryForCreate(path, testPassword)
      if err := repo.Salvar(cofre); err != nil {
          t.Fatalf("Salvar: %v", err)
      }

      // Modificar o arquivo externamente
      if err := os.WriteFile(path, []byte("conteudo diferente"), 0600); err != nil {
          t.Fatalf("WriteFile: %v", err)
      }

      changed, err := repo.DetectarAlteracaoExterna()
      if err != nil {
          t.Fatalf("DetectarAlteracaoExterna: %v", err)
      }
      if !changed {
          t.Error("esperado true (arquivo alterado externamente), obteve false")
      }
  }

  func TestFileRepository_DetectarAlteracaoExterna_CofreNovo(t *testing.T) {
      dir := t.TempDir()
      path := filepath.Join(dir, "novo.abditum")
      repo := NewFileRepositoryForCreate(path, testPassword)

      // Antes do primeiro Salvar, metadata é zero — deve retornar false sem erro
      changed, err := repo.DetectarAlteracaoExterna()
      if err != nil {
          t.Fatalf("DetectarAlteracaoExterna em cofre novo: %v", err)
      }
      if changed {
          t.Error("cofre novo: esperado false, obteve true")
      }
  }
  ```

- [ ] **Rodar os testes para confirmar que falham**

  ```
  CGO_ENABLED=0 go test ./internal/storage/... -run TestFileRepository_DetectarAlteracaoExterna -v
  ```

  Esperado: FAIL — `repo.DetectarAlteracaoExterna undefined`

- [ ] **Implementar o método em `internal/storage/repository.go`**

  Adicionar após o método `Metadata()` (linha ~147):

  ```go
  // DetectarAlteracaoExterna verifica se o arquivo de cofre foi modificado por processo
  // externo desde o último Salvar ou Carregar.
  // Retorna false sem erro se o metadata ainda não foi capturado (cofre recém-criado).
  func (r *FileRepository) DetectarAlteracaoExterna() (bool, error) {
      // Metadata zero significa que nenhum Salvar ou Carregar ocorreu ainda.
      // Não há baseline para comparar — considerar sem alteração externa.
      if r.metadata == (FileMetadata{}) {
          return false, nil
      }
      return DetectExternalChange(r.path, r.metadata)
  }
  ```

- [ ] **Verificar conformidade com a interface em tempo de compilação**

  A linha existente `var _ vault.RepositorioCofre = (*FileRepository)(nil)` em `repository.go` garante isso automaticamente. Rodar:

  ```
  CGO_ENABLED=0 go build ./internal/storage/...
  ```

  Esperado: sem erros.

- [ ] **Rodar os testes para confirmar que passam**

  ```
  CGO_ENABLED=0 go test ./internal/storage/... -run TestFileRepository_DetectarAlteracaoExterna -v
  ```

  Esperado: PASS (3 testes)

- [ ] **Rodar a suite completa de storage para garantir que nada quebrou**

  ```
  CGO_ENABLED=0 go test ./internal/storage/... -race -count=1 -v
  ```

  Esperado: PASS

- [ ] **Commit**

  ```
  git add internal/storage/repository.go internal/storage/storage_test.go
  git commit -m "feat(storage): implementa DetectarAlteracaoExterna no FileRepository"
  ```

---

## Tarefa 4 — Alterar `Manager.Salvar` para aceitar `forcarSobrescrita bool`

**Arquivos:**
- Modificar: `internal/vault/manager.go`
- Modificar: `internal/vault/manager_test.go`

- [ ] **Escrever os testes novos**

  Em `internal/vault/manager_test.go`, adicionar:

  ```go
  func TestManager_Salvar_SemAlteracaoExterna_Sucesso(t *testing.T) {
      cofre := novoCofre()
      repo := &mockRepository{}
      m := NewManager(cofre, repo)
      cofre.modificado = true

      if err := m.Salvar(false); err != nil {
          t.Errorf("Salvar(false) sem alteração externa: esperado nil, obteve %v", err)
      }
      if !repo.salvarCalled {
          t.Error("Salvar(false): repositório não foi chamado")
      }
  }

  func TestManager_Salvar_ComAlteracaoExterna_RetornaErro(t *testing.T) {
      cofre := novoCofre()
      repo := &mockRepository{detectarAlteracaoExternaResp: true}
      m := NewManager(cofre, repo)
      cofre.modificado = true

      err := m.Salvar(false)
      if !errors.Is(err, ErrModifiedExternally) {
          t.Errorf("Salvar(false) com alteração externa: esperado ErrModifiedExternally, obteve %v", err)
      }
      if repo.salvarCalled {
          t.Error("Salvar(false) com alteração externa: repositório não deveria ter sido chamado")
      }
  }

  func TestManager_Salvar_ComAlteracaoExterna_ForcarSobrescrita(t *testing.T) {
      cofre := novoCofre()
      repo := &mockRepository{detectarAlteracaoExternaResp: true}
      m := NewManager(cofre, repo)
      cofre.modificado = true

      if err := m.Salvar(true); err != nil {
          t.Errorf("Salvar(true) com alteração externa: esperado nil, obteve %v", err)
      }
      if !repo.salvarCalled {
          t.Error("Salvar(true): repositório deveria ter sido chamado")
      }
  }

  func TestManager_Salvar_ErroNoRepositorio(t *testing.T) {
      cofre := novoCofre()
      repo := &mockRepository{salvarError: errors.New("disco cheio")}
      m := NewManager(cofre, repo)
      cofre.modificado = true

      err := m.Salvar(false)
      if err == nil {
          t.Error("Salvar: esperado erro do repositório, obteve nil")
      }
  }
  ```

  Verificar se existe função auxiliar `novoCofre()` nos testes — se não existir, usar `NewCofre()` diretamente.

- [ ] **Rodar os testes para confirmar que falham**

  ```
  CGO_ENABLED=0 go test ./internal/vault/... -run "TestManager_Salvar" -v
  ```

  Esperado: FAIL — assinatura de `Salvar` incompatível

- [ ] **Atualizar `Manager.Salvar` em `internal/vault/manager.go`**

  Localizar o método `Salvar()` (linha ~261) e substituir pela nova versão:

  ```go
  // Salvar persiste o cofre no arquivo usando o repositório configurado.
  //
  // Se forcarSobrescrita for false e o arquivo tiver sido modificado externamente
  // desde o último Load ou Save, retorna ErrModifiedExternally sem salvar.
  // Se forcarSobrescrita for true, pula essa verificação e salva diretamente.
  func (m *Manager) Salvar(forcarSobrescrita bool) error {
      if !forcarSobrescrita {
          changed, err := m.repositorio.DetectarAlteracaoExterna()
          if err != nil {
              return fmt.Errorf("vault.Salvar: verificação de alteração externa: %w", err)
          }
          if changed {
              return ErrModifiedExternally
          }
      }

      snapshot, err := m.prepararSnapshot()
      if err != nil {
          return fmt.Errorf("vault.Salvar: %w", err)
      }

      if err := m.repositorio.Salvar(snapshot); err != nil {
          return fmt.Errorf("vault.Salvar: %w", err)
      }

      m.finalizarExclusoes()
      m.cofre.modificado = false
      return nil
  }
  ```

  Verificar se `fmt` já está nos imports de `manager.go` — se não estiver, adicioná-lo.

- [ ] **Rodar os testes para confirmar que passam**

  ```
  CGO_ENABLED=0 go test ./internal/vault/... -run "TestManager_Salvar" -v
  ```

  Esperado: PASS

- [ ] **Rodar a suite completa de vault**

  ```
  CGO_ENABLED=0 go test ./internal/vault/... -race -count=1 -v
  ```

  Esperado: PASS

- [ ] **Commit**

  ```
  git add internal/vault/manager.go internal/vault/manager_test.go
  git commit -m "feat(vault): Salvar aceita forcarSobrescrita para detecção de modificação externa"
  ```

---

## Tarefa 5 — Criar `QuitOperation`

**Arquivos:**
- Criar: `internal/tui/operation/quit_operation.go`
- Criar: `internal/tui/operation/quit_operation_test.go`

- [ ] **Escrever os testes**

  Arquivo: `internal/tui/operation/quit_operation_test.go`

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
      isModified                   bool
      salvarErr                    error
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
      op := NewQuitOperation(&stubNotifier{}, nil)
      msg := execCmd(op.Init())
      if _, ok := msg.(tui.OpenModalMsg); !ok {
          t.Errorf("Init sem cofre: esperado OpenModalMsg, obteve %T", msg)
      }
  }

  func TestQuitOperation_Init_CoffeInalterado_AbreModalConfirmacao(t *testing.T) {
      op := NewQuitOperation(&stubNotifier{}, &stubManager{isModified: false})
      msg := execCmd(op.Init())
      if _, ok := msg.(tui.OpenModalMsg); !ok {
          t.Errorf("Init cofre inalterado: esperado OpenModalMsg, obteve %T", msg)
      }
  }

  func TestQuitOperation_Init_CofreAlterado_AbreModalDecisao(t *testing.T) {
      op := NewQuitOperation(&stubNotifier{}, &stubManager{isModified: true})
      msg := execCmd(op.Init())
      if _, ok := msg.(tui.OpenModalMsg); !ok {
          t.Errorf("Init cofre alterado: esperado OpenModalMsg, obteve %T", msg)
      }
  }

  // --- Update ---

  func TestQuitOperation_Update_IgnoraMensagemDesconhecida(t *testing.T) {
      op := NewQuitOperation(&stubNotifier{}, nil)
      type outraMsg struct{}
      if cmd := op.Update(outraMsg{}); cmd != nil {
          t.Error("Update(outraMsg): esperado nil")
      }
  }

  func TestQuitOperation_Update_Saving_Sucesso_EmiteQuit(t *testing.T) {
      n := &stubNotifier{}
      m := &stubManager{}
      op := NewQuitOperation(n, m)

      cmd := op.Update(quitMsg{state: quitStateSaving})
      if n.lastMethod != "SetBusy" {
          t.Errorf("Update(saving): esperado SetBusy, obteve %q", n.lastMethod)
      }
      // cmd é o saveCmd — executar para obter o resultado
      resultCmd := execCmd(cmd).(tea.Cmd)
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
      m := &stubManager{salvarErr: vault.ErrModifiedExternally}
      op := NewQuitOperation(n, m)

      cmd := op.Update(quitMsg{state: quitStateSaving})
      resultCmd := execCmd(cmd).(tea.Cmd)
      msg := execCmd(resultCmd)
      if _, ok := msg.(tui.OpenModalMsg); !ok {
          t.Errorf("após ErrModifiedExternally: esperado OpenModalMsg, obteve %T", msg)
      }
      if n.lastMethod != "Clear" {
          t.Errorf("após ErrModifiedExternally: esperado Clear no notifier, obteve %q", n.lastMethod)
      }
  }

  func TestQuitOperation_Update_Saving_ErroGenerico_NaoSai(t *testing.T) {
      n := &stubNotifier{}
      m := &stubManager{salvarErr: errors.New("disco cheio")}
      op := NewQuitOperation(n, m)

      cmd := op.Update(quitMsg{state: quitStateSaving})
      resultCmd := execCmd(cmd).(tea.Cmd)
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
      m := &stubManager{}
      op := NewQuitOperation(n, m)

      cmd := op.Update(quitMsg{state: quitStateSavingForced})
      resultCmd := execCmd(cmd).(tea.Cmd)
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
      m := &stubManager{salvarErr: errors.New("falha")}
      op := NewQuitOperation(n, m)

      cmd := op.Update(quitMsg{state: quitStateSavingForced})
      resultCmd := execCmd(cmd).(tea.Cmd)
      msg := execCmd(resultCmd)
      if n.lastMethod != "SetError" {
          t.Errorf("SavingForced erro: esperado SetError, obteve %q", n.lastMethod)
      }
      if _, ok := msg.(tui.OperationCompletedMsg); !ok {
          t.Errorf("SavingForced erro: esperado OperationCompletedMsg, obteve %T", msg)
      }
  }
  ```

- [ ] **Rodar os testes para confirmar que falham**

  ```
  CGO_ENABLED=0 go test ./internal/tui/operation/... -run "TestQuitOperation" -v
  ```

  Esperado: FAIL — `NewQuitOperation undefined`

- [ ] **Criar `internal/tui/operation/quit_operation.go`**

  ```go
  package operation

  // QuitOperation implementa o fluxo de saída da aplicação (Fluxos 3, 4 e 5 do fluxos.md).
  //
  //   - Fluxo 3 (sem cofre carregado): confirmação simples → tea.Quit.
  //   - Fluxo 4 (cofre inalterado): confirmação simples → tea.Quit.
  //   - Fluxo 5 (cofre alterado): salvar e sair / descartar e sair / voltar.
  //     Se o repositório detectar modificação externa, abre modal de conflito antes de forçar.

  import (
      tea "charm.land/bubbletea/v2"
      "github.com/useful-toys/abditum/internal/tui"
      "github.com/useful-toys/abditum/internal/tui/design"
      "github.com/useful-toys/abditum/internal/tui/modal"
      "github.com/useful-toys/abditum/internal/vault"
  )

  type quitState int

  const (
      quitStateSaving       quitState = iota // salvar com forcarSobrescrita=false
      quitStateSavingForced                  // salvar com forcarSobrescrita=true
  )

  type quitMsg struct {
      state quitState
  }

  // vaultSaver é a fatia da interface vault.Manager que QuitOperation precisa.
  // Facilita substituição por stub nos testes sem depender do tipo concreto.
  type vaultSaver interface {
      IsModified() bool
      Salvar(forcarSobrescrita bool) error
  }

  // QuitOperation gerencia o fluxo de encerramento da aplicação.
  type QuitOperation struct {
      notifier tui.MessageController
      manager  vaultSaver // nil se nenhum cofre estiver carregado
  }

  // NewQuitOperation cria uma QuitOperation com o estado de cofre atual.
  // manager pode ser nil (nenhum cofre carregado).
  func NewQuitOperation(notifier tui.MessageController, manager *vault.Manager) *QuitOperation {
      var saver vaultSaver
      if manager != nil {
          saver = manager
      }
      return &QuitOperation{notifier: notifier, manager: saver}
  }

  // Init abre o modal adequado de acordo com o estado do cofre.
  func (q *QuitOperation) Init() tea.Cmd {
      if q.manager != nil && q.manager.IsModified() {
          return tui.OpenModal(q.buildModifiedModal())
      }
      return tui.OpenModal(q.buildConfirmModal())
  }

  // Update processa mensagens internas da operação.
  func (q *QuitOperation) Update(msg tea.Msg) tea.Cmd {
      m, ok := msg.(quitMsg)
      if !ok {
          return nil
      }
      switch m.state {
      case quitStateSaving:
          q.notifier.SetBusy("Salvando...")
          return q.saveCmd(false)
      case quitStateSavingForced:
          q.notifier.SetBusy("Salvando...")
          return q.saveCmd(true)
      }
      return nil
  }

  // saveCmd executa o salvamento em goroutine e retorna o comando com o resultado.
  func (q *QuitOperation) saveCmd(forcar bool) tea.Cmd {
      return func() tea.Msg {
          err := q.manager.Salvar(forcar)
          if err == nil {
              return quitSaveResultMsg{err: nil}
          }
          return quitSaveResultMsg{err: err}
      }
  }

  // quitSaveResultMsg carrega o resultado do salvamento de volta ao Update via Bubble Tea.
  type quitSaveResultMsg struct {
      err     error
      forcado bool
  }

  // --- NOTA: Update também precisa tratar quitSaveResultMsg ---
  // Atualizar Update para tratar quitSaveResultMsg:

  func (q *QuitOperation) UpdateFull(msg tea.Msg) tea.Cmd {
      switch m := msg.(type) {
      case quitMsg:
          switch m.state {
          case quitStateSaving:
              q.notifier.SetBusy("Salvando...")
              return func() tea.Msg {
                  return quitSaveResultMsg{err: q.manager.Salvar(false), forcado: false}
              }
          case quitStateSavingForced:
              q.notifier.SetBusy("Salvando...")
              return func() tea.Msg {
                  return quitSaveResultMsg{err: q.manager.Salvar(true), forcado: true}
              }
          }
      case quitSaveResultMsg:
          if m.err == nil {
              q.notifier.SetSuccess("Cofre salvo.")
              return tea.Quit
          }
          if !m.forcado && isModifiedExternally(m.err) {
              q.notifier.Clear()
              return tui.OpenModal(q.buildConflictModal())
          }
          q.notifier.SetError(m.err.Error())
          return tui.OperationCompleted()
      }
      return nil
  }
  ```

  **Atenção:** o esboço acima tem `Update` e `UpdateFull` separados para clareza, mas a implementação final deve ter apenas um método `Update` que trata ambos os tipos de mensagem. Usar a versão de `UpdateFull` como `Update`:

  ```go
  func (q *QuitOperation) Update(msg tea.Msg) tea.Cmd {
      switch m := msg.(type) {
      case quitMsg:
          switch m.state {
          case quitStateSaving:
              q.notifier.SetBusy("Salvando...")
              return func() tea.Msg {
                  return quitSaveResultMsg{err: q.manager.Salvar(false), forcado: false}
              }
          case quitStateSavingForced:
              q.notifier.SetBusy("Salvando...")
              return func() tea.Msg {
                  return quitSaveResultMsg{err: q.manager.Salvar(true), forcado: true}
              }
          }
      case quitSaveResultMsg:
          if m.err == nil {
              q.notifier.SetSuccess("Cofre salvo.")
              return tea.Quit
          }
          if !m.forcado && isModifiedExternally(m.err) {
              q.notifier.Clear()
              return tui.OpenModal(q.buildConflictModal())
          }
          q.notifier.SetError(m.err.Error())
          return tui.OperationCompleted()
      }
      return nil
  }

  func isModifiedExternally(err error) bool {
      return errors.Is(err, vault.ErrModifiedExternally)
  }
  ```

  Adicionar `"errors"` ao import.

  **Modais:**

  ```go
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

  func (q *QuitOperation) buildModifiedModal() *modal.ConfirmModal {
      return modal.NewConfirmModal(
          "Sair",
          "Há alterações não salvas. O que deseja fazer?",
          []modal.ModalOption{
              {
                  Keys:  []design.Key{design.Keys.Enter},
                  Label: "Salvar e sair",
                  Action: func() tea.Cmd {
                      return tea.Batch(tui.CloseModal(), func() tea.Msg {
                          return quitMsg{state: quitStateSaving}
                      })
                  },
              },
              {
                  Keys:  []design.Key{design.Letter('d')},
                  Label: "Descartar e sair",
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

  func (q *QuitOperation) buildConflictModal() *modal.ConfirmModal {
      return modal.NewConfirmModal(
          "Conflito",
          "O arquivo foi modificado externamente. Deseja sobrescrever?",
          []modal.ModalOption{
              {
                  Keys:  []design.Key{design.Keys.Enter},
                  Label: "Sobrescrever e sair",
                  Action: func() tea.Cmd {
                      return tea.Batch(tui.CloseModal(), func() tea.Msg {
                          return quitMsg{state: quitStateSavingForced}
                      })
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

- [ ] **Rodar os testes para confirmar que passam**

  ```
  CGO_ENABLED=0 go test ./internal/tui/operation/... -run "TestQuitOperation" -v
  ```

  Esperado: PASS

- [ ] **Rodar a suite completa de operation**

  ```
  CGO_ENABLED=0 go test ./internal/tui/operation/... -race -count=1 -v
  ```

  Esperado: PASS

- [ ] **Commit**

  ```
  git add internal/tui/operation/quit_operation.go internal/tui/operation/quit_operation_test.go
  git commit -m "feat(tui): implementa QuitOperation com Fluxos 3, 4 e 5"
  ```

---

## Tarefa 6 — Conectar `QuitOperation` ao `setup.go`

**Arquivos:**
- Modificar: `cmd/abditum/setup.go`

- [ ] **Remover o import `tea` se não for mais usado diretamente**

  Após a mudança, `tea.Quit` não será mais chamado diretamente no `setup.go`. Verificar se `tea` ainda é usado (ex: no `ThemeToggle` retorna `nil`, não usa `tea`). Se não for mais usado, o compilador apontará — remover do import.

- [ ] **Substituir a action Quit**

  Localizar o bloco da action `Quit` (linhas ~58–66) e substituir:

  ```go
  {
      Keys:        []design.Key{design.Shortcuts.Quit},
      Label:       "Sair",
      Description: "Encerra a aplicação, com opção de salvar se houver alterações.",
      GroupID:     "app",
      Priority:    20,
      Visible:     true,
      OnExecute: func() tea.Cmd {
          return tui.StartOperation(operation.NewQuitOperation(r.MessageController(), r.Manager()))
      },
  },
  ```

- [ ] **Compilar o projeto inteiro**

  ```
  CGO_ENABLED=0 go build ./...
  ```

  Esperado: sem erros.

- [ ] **Rodar a suite completa**

  ```
  CGO_ENABLED=0 go test ./... -race -count=1 -v
  ```

  Esperado: PASS

- [ ] **Commit**

  ```
  git add cmd/abditum/setup.go
  git commit -m "feat(cmd): ctrl Q passa a usar QuitOperation"
  ```

---

## Verificação final

- [ ] Rodar `CGO_ENABLED=0 go test ./... -race -count=1` — todos os testes passam.
- [ ] Iniciar a aplicação e testar manualmente:
  - `ctrl Q` sem cofre carregado → modal de confirmação simples.
  - `ctrl Q` com cofre inalterado → modal de confirmação simples.
  - `ctrl Q` com cofre alterado → modal de decisão (salvar / descartar / voltar).
  - Escolher "Salvar e sair" → indicador de progresso → mensagem de sucesso → aplicação encerra.
  - Escolher "Descartar e sair" → aplicação encerra sem salvar.
  - Escolher "Voltar" → retorna ao estado anterior, operação encerrada.
