# Theme Parameter in View() — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended)

**Goal:** Atualizar todas as assinaturas de View() para receber theme como parâmetro, eliminando campo theme dos modelos.

**Arquitetura:** O rootModel possui o tema e o passa como parâmetro em todas as chamadas View(). Apenas rootModel armazena theme.

**Tech Stack:** Go, BubbleTea, lipgloss

---

## File Map

| Ação | Arquivo | Responsabilidade |
|---|---|---|
| Modificar | `internal/tui/flows.go` | ✅ childModel e modalView com theme |
| Modificar | `internal/tui/welcome.go` | WorkArea welcome |
| Modificar | `internal/tui/vaulttree.go` | WorkArea vault tree |
| Modificar | `internal/tui/templatelist.go` | WorkArea template list |
| Modificar | `internal/tui/templatedetail.go` | WorkArea template detail |
| Modificar | `internal/tui/settings.go` | WorkArea settings |
| Modificar | `internal/tui/secretdetail.go` | WorkArea secret detail |
| Modificar | `internal/tui/decision.go` | Modal dialog |
| Modificar | `internal/tui/passwordentry.go` | Modal password entry |
| Modificar | `internal/tui/passwordcreate.go` | Modal password create |
| Modificar | `internal/tui/help.go` | Modal help |
| Modificar | `internal/tui/filepicker.go` | Modal file picker |
| Modificar | `internal/tui/modal.go` | Modal base |
| Modificar | `internal/tui/flow_create_vault.go` | Flow create vault |
| Modificar | `internal/tui/flow_open_vault.go` | Flow open vault |
| Modificar | `internal/tui/flow_save_and_exit.go` | Flow save and exit |
| Modificar | `internal/tui/root.go` | Root orchestration |
| Modificar | `internal/tui/welcome_test.go` | Testes |
| Modificar | `internal/tui/vaulttree_test.go` | Testes |
| Modificar | `internal/tui/templatelist_test.go` | Testes |
| Modificar | `internal/tui/templatedetail_test.go` | Testes |
| Modificar | `internal/tui/settings_test.go` | Testes |
| Modificar | `internal/tui/secretdetail_test.go` | Testes |
| Modificar | `internal/tui/decision_test.go` | Testes |
| Modificar | `internal/tui/passwordentry_test.go` | Testes |
| Modificar | `internal/tui/passwordcreate_test.go` | Testes |
| Modificar | `internal/tui/help_test.go` | Testes |
| Modificar | `internal/tui/filepicker_test.go` | Testes |
| Modificar | `internal/tui/flow_create_vault_test.go` | Testes |
| Modificar | `internal/tui/flow_open_vault_test.go` | Testes |
| Modificar | `internal/tui/flow_save_and_exit_test.go` | Testes |
| Modificar | `internal/tui/root_test.go` | Testes |

---

## Tarefa 1: Atualizar workarea models (welcome, vaulttree, templatelist, templatedetail, settings, secretdetail)

**Files:** Modify cada arquivo de workarea

Para cada workarea model (ex: welcome.go):
- Remover campo `theme *Theme` da struct
- Remover método `ApplyTheme(*Theme)`
- Atualizar `View(width, height int)` para `View(width, height int, theme *Theme)`
- Atualizar todas as referências internas de `m.theme` para `theme`

### Tarefa 1a: welcome.go

- Remover campo `theme` da struct `welcomeModel`
- Remover método `ApplyTheme`
- Atualizar View signature
- Atualizar corpo para usar parâmetro `theme`

### Tarefa 1b: vaulttree.go

Similar a welcome.go

### Tarefa 1c: templatelist.go, templatedetail.go, settings.go, secretdetail.go

Similar a welcome.go

### Verificação Tarefa 1
```bash
go build ./internal/tui/...
```
Expected: erros de compilação indicando quais modelos ainda precisam ser atualizados

---

## Tarefa 2: Atualizar modal models (decision, passwordentry, passwordcreate, help, filepicker, modal)

**Files:** Modify cada arquivo de modal

Para cada modal model:
- Se tiver campo `theme`, remover
- Se tiver método `ApplyTheme`, remover
- Atualizar `View(maxWidth, maxHeight int)` para `View(maxWidth, maxHeight int, theme *Theme)`
- Atualizar referências internas

### Tarefa 2a: decision.go

- Campo `theme` já adicionado como `*Theme` (não removido ainda)
- Remover getter `activeTheme()` 
- Remover método `WithTheme()`
- Atualizar View signature
- Atualizar corpo

### Tarefa 2b: passwordentry.go, passwordcreate.go, help.go, filepicker.go, modal.go

Similar a decision.go

### Verificação Tarefa 2
```bash
go build ./internal/tui/...
```

---

## Tarefa 3: Atualizar flow models (flow_create_vault, flow_open_vault, flow_save_and_exit)

**Files:** Modify flow_*.go

Para cada flow model:
- Se tiver campo `theme`, remover
- Atualizar View signature
- Atualizar referências internas

### Verificação Tarefa 3
```bash
go build ./internal/tui/...
```

---

## Tarefa 4: Atualizar root.go

**Files:** Modify root.go

- Atualizar todas as chamadas para passar `m.theme` como parâmetro
- Ex: `m.welcome.View(m.width, m.height, m.theme)`
- Não precisa de campo theme aqui (root já tem)

### Verificação Tarefa 4
```bash
go build ./internal/tui/...
go vet ./internal/tui/...
```

---

## Tarefa 5: Atualizar testes (todos os *_test.go)

**Files:** Modify todos os arquivos de teste

Para cada teste:
- Atualizar chamadas View() para passar `TokyoNight` como terceiro parâmetro
- Ex: `m.View(80, 24)` → `m.View(80, 24, TokyoNight)`

**Importante:** Não alterar golden files

### Verificação Tarefa 5
```bash
go test ./internal/tui/... -count=1
```

Expected: Todos os testes passam (golden files não precisam mudar)

---

## Validação Final

```bash
go build ./...
go vet ./...
go test ./... -count=1
git diff --stat main..HEAD
```

---

## Commit Sugerido

```
$ git commit -m "refactor: pass theme as parameter to all View() methods

- Updated childModel and modalView interfaces to include theme parameter
- Removed theme field from all models (except rootModel)
- Updated all View() calls in root.go to pass m.theme
- Updated tests to pass TokyoNight
- No golden files changed"
```

---

## Execução

Two options:

1. **Subagent-Driven** - recomended
2. **Inline Execution**

Which approach?