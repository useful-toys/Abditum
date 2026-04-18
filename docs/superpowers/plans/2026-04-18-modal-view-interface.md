# Refatoração da interface ModalView — Plano de Implementação

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remover `Update(tea.Msg)` da interface `ModalView` e adicionar `HandleMouse(tea.MouseMsg)`, tornando o contrato explícito para teclado e mouse.

**Architecture:** A interface `ModalView` passa a ter quatro métodos: `Render`, `HandleKey`, `HandleMouse`, `Cursor`. O `RootModel` substitui chamadas a `.Update(msg)` nos modais por `.HandleKey()` ou `.HandleMouse()` conforme o tipo da mensagem. Para `ModalReadyMsg`, o modal não recebe mais a mensagem — apenas a view ativa recebe.

**Tech Stack:** Go, charm.land/bubbletea/v2

---

### Task 1: Atualizar a interface `ModalView`

**Files:**
- Modify: `internal/tui/modal.go:10-21`

- [ ] **Step 1: Substituir a definição da interface**

Em `internal/tui/modal.go`, substituir o bloco da interface:

```go
// ModalView define o contrato para componentes de modal da interface.
// Modais são exibidos sobrepostos à área de trabalho e gerenciados por RootModel.
type ModalView interface {
	// Render retorna a representação visual do modal dentro dos limites fornecidos.
	// theme é passado por ponteiro para evitar cópia — design.Theme tem 400 bytes.
	Render(maxHeight, maxWidth int, theme *design.Theme) string
	// HandleKey processa eventos de teclado e retorna um comando ou nil.
	HandleKey(msg tea.KeyMsg) tea.Cmd
	// HandleMouse processa eventos de mouse e retorna um comando ou nil.
	HandleMouse(msg tea.MouseMsg) tea.Cmd
	// Cursor retorna a posição do cursor real para o modal ativo, ou nil se não houver cursor.
	// topY e leftX são as coordenadas absolutas do canto superior esquerdo do modal na tela.
	Cursor(topY, leftX int) *tea.Cursor
}
```

- [ ] **Step 2: Verificar que o projeto não compila (confirmar quebra)**

```
go build ./...
```

Esperado: erros de compilação nos 5 modais e em `root.go` por `Update` não existir mais na interface e por falta de `HandleMouse`.

---

### Task 2: Atualizar `ConfirmModal`

**Files:**
- Modify: `internal/tui/modal/confirm_modal.go:68-74`

- [ ] **Step 1: Remover `Update` e adicionar `HandleMouse`**

Substituir o método `Update` pelo `HandleMouse`:

```go
// HandleMouse processa eventos de mouse. ConfirmModal não reage a mouse — retorna nil.
func (m *ConfirmModal) HandleMouse(_ tea.MouseMsg) tea.Cmd {
	return nil
}
```

- [ ] **Step 2: Verificar compilação parcial**

```
go build ./internal/tui/modal/...
```

Esperado: `confirm_modal.go` compila sem erros (os outros modais ainda podem falhar).

---

### Task 3: Atualizar `HelpModal`

**Files:**
- Modify: `internal/tui/modal/help_modal.go`

- [ ] **Step 1: Localizar e remover o método `Update`**

Localizar o método `Update` (próximo à linha 168) e substituí-lo por:

```go
// HandleMouse processa eventos de mouse. HelpModal não reage a mouse — retorna nil.
func (m *HelpModal) HandleMouse(_ tea.MouseMsg) tea.Cmd {
	return nil
}
```

---

### Task 4: Atualizar `PasswordCreateModal`

**Files:**
- Modify: `internal/tui/modal/password_create_modal.go`

- [ ] **Step 1: Localizar e remover o método `Update`**

Localizar o método `Update` (próximo à linha 231) e substituí-lo por:

```go
// HandleMouse processa eventos de mouse. PasswordCreateModal não reage a mouse — retorna nil.
func (m *PasswordCreateModal) HandleMouse(_ tea.MouseMsg) tea.Cmd {
	return nil
}
```

---

### Task 5: Atualizar `PasswordEntryModal`

**Files:**
- Modify: `internal/tui/modal/password_entry_modal.go`

- [ ] **Step 1: Localizar e remover o método `Update`**

Localizar o método `Update` (próximo à linha 129) e substituí-lo por:

```go
// HandleMouse processa eventos de mouse. PasswordEntryModal não reage a mouse — retorna nil.
func (m *PasswordEntryModal) HandleMouse(_ tea.MouseMsg) tea.Cmd {
	return nil
}
```

---

### Task 6: Atualizar `FilePickerModal`

**Files:**
- Modify: `internal/tui/modal/file_picker.go`

- [ ] **Step 1: Localizar e remover o método `Update`**

Localizar o método `Update` (próximo à linha 667) e substituí-lo por:

```go
// HandleMouse processa eventos de mouse. FilePickerModal não reage a mouse — retorna nil.
func (m *FilePickerModal) HandleMouse(_ tea.MouseMsg) tea.Cmd {
	return nil
}
```

---

### Task 7: Atualizar `RootModel` — substituir chamadas a `.Update()` nos modais

**Files:**
- Modify: `internal/tui/root.go:328-395`

Há quatro pontos em `root.go` onde modais recebem `.Update(msg)`. Cada um tem um tratamento diferente:

**Ponto 1 — linha 331:** `ModalReadyMsg` para o modal pai na pilha
```go
// Antes:
return r, parent.Update(msg)
// Depois: modais não recebem ModalReadyMsg — apenas views recebem.
// O modal pai não precisa ser notificado; a view ativa cuida do resultado.
return r, nil
```

**Ponto 2 — linha 333:** `ModalReadyMsg` para a view ativa (não é modal — não muda)
```go
// Não muda — activeView.Update(msg) permanece igual.
return r, r.activeView.Update(msg)
```

**Ponto 3 — linha 344:** `tea.KeyMsg` para o modal no topo da pilha
```go
// Antes:
return r, r.modals[top].Update(msg)
// Depois:
return r, r.modals[top].HandleKey(msg)
```

**Ponto 4 — linha 390:** mensagens genéricas para o modal no topo da pilha
```go
// Antes:
cmds = append(cmds, r.modals[top].Update(msg))
// Depois: modais só recebem mouse além de key.
// tea.KeyMsg já foi tratado no case acima (linha 335).
// Aqui chegam apenas mensagens não-key — verificar se é MouseMsg.
if mouseMsg, ok := msg.(tea.MouseMsg); ok {
    cmds = append(cmds, r.modals[top].HandleMouse(mouseMsg))
}
```

- [ ] **Step 1: Aplicar as quatro substituições acima em `root.go`**

- [ ] **Step 2: Verificar que o projeto compila**

```
go build ./...
```

Esperado: compilação sem erros.

- [ ] **Step 3: Executar testes**

```
go test ./...
```

Esperado: todos os testes passam.

- [ ] **Step 4: Commit**

```
git add internal/tui/modal.go internal/tui/modal/confirm_modal.go internal/tui/modal/help_modal.go internal/tui/modal/password_create_modal.go internal/tui/modal/password_entry_modal.go internal/tui/modal/file_picker.go internal/tui/root.go
git commit -m "refactor(tui): substituir Update por HandleKey e HandleMouse em ModalView"
```
