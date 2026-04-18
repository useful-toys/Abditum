# Remover `Intent` de `ModalOption` — Plano de Implementação

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remover o campo `Intent`, o tipo `Intent int` e suas constantes de `ModalOption`, eliminando dead code sem alterar nenhum comportamento.

**Architecture:** Remoção pura — apagar declarações de tipo/constantes em `modal_base.go`, remover o campo de todos os literais de struct no código de produção e nos testes, e deletar o teste que verificava apenas o round-trip do campo.

**Tech Stack:** Go, `go build`, `go test`

---

### Task 1: Remover tipo, constantes e campo em `modal_base.go`

**Files:**
- Modify: `internal/tui/modal/modal_base.go`

- [ ] **Step 1: Abrir o arquivo e localizar as declarações**

  Arquivo atual (linhas 8–29):
  ```go
  // Intent classifica a intenção semântica de uma opção de modal.
  type Intent int

  const (
      IntentConfirm Intent = iota
      IntentCancel
      IntentOther
  )

  // ModalOption represents an action available to the user within a modal.
  type ModalOption struct {
      Keys   []design.Key
      Label  string
      Intent Intent
      Action func() tea.Cmd
  }
  ```

- [ ] **Step 2: Remover o bloco de tipo/constantes e o campo `Intent`**

  Resultado esperado após edição:
  ```go
  // ModalOption representa uma ação disponível ao usuário dentro de um modal.
  type ModalOption struct {
      // Keys lista as teclas que ativam esta opção.
      // Keys[0].Label é exibido no rodapé do diálogo.
      // Outras Keys são aliases funcionais (ex: Enter como alias para "S Sobrescrever").
      Keys []design.Key
      // Label é o texto exibido ao usuário descrevendo a ação.
      Label string
      // Action é a função executada quando a opção é escolhida.
      Action func() tea.Cmd
  }
  ```

- [ ] **Step 3: Verificar compilação**

  ```bash
  go build ./internal/tui/modal/...
  ```
  Esperado: erros de compilação nos arquivos que ainda referenciam `IntentConfirm`, `IntentCancel`, `IntentOther`. Isso é esperado — as próximas tasks corrigem isso.

---

### Task 2: Remover `Intent` dos literais em `password_entry_modal.go`

**Files:**
- Modify: `internal/tui/modal/password_entry_modal.go`

- [ ] **Step 1: Localizar os literais de `ModalOption`**

  Em torno das linhas 70–91 há dois literais com `Intent: IntentConfirm` e `Intent: IntentCancel`.

- [ ] **Step 2: Remover o campo `Intent:` de ambos os literais**

  Antes:
  ```go
  modal.ModalOption{
      Keys:   []design.Key{...},
      Label:  "...",
      Intent: IntentConfirm,
      Action: func() tea.Cmd { ... },
  },
  ```
  Depois:
  ```go
  modal.ModalOption{
      Keys:   []design.Key{...},
      Label:  "...",
      Action: func() tea.Cmd { ... },
  },
  ```
  Aplicar para os dois literais no arquivo.

- [ ] **Step 3: Verificar compilação parcial**

  ```bash
  go build ./internal/tui/modal/...
  ```

---

### Task 3: Remover `Intent` dos literais em `password_create_modal.go`

**Files:**
- Modify: `internal/tui/modal/password_create_modal.go`

- [ ] **Step 1: Localizar os literais** (em torno das linhas 148–174)

- [ ] **Step 2: Remover `Intent: IntentConfirm` e `Intent: IntentCancel`** dos dois literais

  Padrão idêntico ao da Task 2.

- [ ] **Step 3: Verificar compilação**

  ```bash
  go build ./internal/tui/modal/...
  ```

---

### Task 4: Remover `Intent` dos literais em `help_modal.go`

**Files:**
- Modify: `internal/tui/modal/help_modal.go`

- [ ] **Step 1: Localizar os dois locais** — `NewHelpModal` (linhas ~37–44) e `Render` (linhas ~77–84)

- [ ] **Step 2: Remover `Intent: IntentCancel`** de ambos os literais

- [ ] **Step 3: Verificar compilação**

  ```bash
  go build ./internal/tui/modal/...
  ```

---

### Task 5: Remover `Intent` dos literais em `fake_operation.go`

**Files:**
- Modify: `internal/tui/operation/fake_operation.go`

- [ ] **Step 1: Localizar os três literais** (linhas ~56–75 e ~83–92) com `Intent: modal.IntentConfirm` e `Intent: modal.IntentCancel`

- [ ] **Step 2: Remover o campo `Intent:` dos três literais**

- [ ] **Step 3: Verificar compilação**

  ```bash
  go build ./internal/tui/operation/...
  ```

---

### Task 6: Remover `Intent` dos literais em `cmd/test_calc/main.go`

**Files:**
- Modify: `cmd/test_calc/main.go`

- [ ] **Step 1: Localizar os literais** (linhas ~13–16) com `modal.IntentConfirm` e `modal.IntentCancel`

- [ ] **Step 2: Remover o campo `Intent:` dos dois literais**

- [ ] **Step 3: Verificar compilação completa do projeto**

  ```bash
  go build ./...
  ```
  Esperado: sem erros.

---

### Task 7: Remover `Intent` dos testes

**Files:**
- Modify: `internal/tui/modal/confirm_modal_test.go`
- Modify: `internal/tui/modal/frame_test.go`
- Modify: `internal/tui/modal/key_handler_test.go`

- [ ] **Step 1: `confirm_modal_test.go` — remover `TestConfirmModal_IntentTypes_Preserved`**

  Deletar a função de teste inteira (buscar por `TestConfirmModal_IntentTypes_Preserved`).

- [ ] **Step 2: `confirm_modal_test.go` — remover todos os campos `Intent:` dos literais**

  São aproximadamente 28 ocorrências. Remover cada linha do tipo:
  ```go
  Intent: modal.IntentConfirm,
  ```
  ou
  ```go
  Intent: modal.IntentCancel,
  ```
  ou
  ```go
  Intent: modal.IntentOther,
  ```

- [ ] **Step 3: `frame_test.go` — remover campos `Intent:` dos literais** (linhas ~38–39, ~95–96, ~107)

- [ ] **Step 4: `key_handler_test.go` — remover campo `Intent:` do literal** (linha ~20)

- [ ] **Step 5: Rodar todos os testes do pacote**

  ```bash
  go test ./internal/tui/modal/...
  ```
  Esperado: todos PASS.

---

### Task 8: Build e testes finais + commit

- [ ] **Step 1: Build completo**

  ```bash
  go build ./...
  ```
  Esperado: sem erros.

- [ ] **Step 2: Testes completos**

  ```bash
  go test ./...
  ```
  Esperado: todos PASS (menos o teste removido).

- [ ] **Step 3: Commit**

  ```bash
  git add internal/tui/modal/modal_base.go \
          internal/tui/modal/password_entry_modal.go \
          internal/tui/modal/password_create_modal.go \
          internal/tui/modal/help_modal.go \
          internal/tui/modal/confirm_modal_test.go \
          internal/tui/modal/frame_test.go \
          internal/tui/modal/key_handler_test.go \
          internal/tui/operation/fake_operation.go \
          cmd/test_calc/main.go
  git commit -m "refactor: remove campo Intent de ModalOption (dead code)"
  ```
