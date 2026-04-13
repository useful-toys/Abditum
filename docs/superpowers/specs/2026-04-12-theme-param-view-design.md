# Design — Passagem de Tema via Parâmetro em View()

**Data:** 2026-04-12
**Pacote:** `internal/tui`

---

## Contexto

O design system foi consolidado em `design.go` mas a maneira como o tema é passado para os modelos ainda é diversa. Cada modelo armazena o tema internamente e a mudança de tema requer mensagens de update.

Este spec endereça a simplificação: tema passado como parâmetro em todas as chamadas `View()`.

---

## Design Aprovado

### Arquitetura Proposta

1. **Nova assinatura unificada:**
```go
// WorkArea models
func (m *welcomeModel) View(width, height int, theme *Theme) string
func (m *vaultTreeModel) View(width, height int, theme *Theme) string
func (m *settingsModel) View(width, height int, theme *Theme) string
func (m *secretDetailModel) View(width, height int, theme *Theme) string
func (m *templateListModel) View(width, height int, theme *Theme) string
func (m *templateDetailModel) View(width, height int, theme *Theme) string

// Modal/Flow models
func (d *DecisionDialog) View(maxWidth, maxHeight int, theme *Theme) string
func (m *passwordEntryModal) View(maxWidth, maxHeight int, theme *Theme) string
func (m *passwordCreateModal) View(maxWidth, maxHeight int, theme *Theme) string
func (m *helpModal) View(maxWidth, maxHeight int, theme *Theme) string
func (m *filePickerModal) View(maxWidth, maxHeight int, theme *Theme) string
func (m *modalModel) View(maxWidth, maxHeight int, theme *Theme) string

// Flow models
func (f *openVaultFlow) View(width, height int, theme *Theme) string
func (f *createVaultFlow) View(width, height int, theme *Theme) string
func (f *saveAndExitFlow) View(width, height int, theme *Theme) string
```

2. **root.View() orchestracommand:**
```go
func (m *rootModel) View() string {
    theme := m.activeTheme() // retorna tema atual (TokyoNight ou Cyberpunk)
    return header...
    + m.welcome.View(m.width, m.height, theme)
    + m.vaultTree.View(m.width, m.height, theme)
    + ...etc para cada workarea e modal
}
```

3. **Toggle Tema (F12):**
```go
case tea.KeyF12:
    m.theme = Cyberpunk // ou TokyoNight
    return m, nil
// NÃO envia mensagem - próxima View() usa novo tema automaticamente
```

4. **Interface modalView atualizada:**
```go
type modalView interface {
    Update(tea.Msg) tea.Cmd
    View(maxWidth, maxHeight int, theme *Theme) string
    Shortcuts() []Shortcut
}
```

---

## Modelos Afetados (~16 arquivos)

| Categoria | Arquivos |
|---|---|
| WorkArea | welcome.go, vaulttree.go, templatelist.go, templatedetail.go, settings.go, secretdetail.go |
| Modal | decision.go, passwordentry.go, passwordcreate.go, help.go, filepicker.go, modal.go |
| Flow | flow_create_vault.go, flow_open_vault.go, flow_save_and_exit.go |
| Components | header.go, actions.go, messages.go |
| Root | root.go |

---

## Testes

- Atualizar todas as chamadas `.View(w, h)` para `.View(w, h, TokyoNight)`
- Golden files **NÃO** precisam ser alterados (testam output, não assinatura)

---

## Não Afetados

- Funções que já usam parâmetro theme: `header.Render()`, `RenderCommandBar()`, `RenderMessageBar()`