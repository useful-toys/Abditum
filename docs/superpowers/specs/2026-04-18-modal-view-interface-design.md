# Design: Refatoração da interface ModalView

**Data:** 2026-04-18

## Contexto

A interface `ModalView` em `internal/tui/modal.go` possui atualmente um método `Update(msg tea.Msg) tea.Cmd` que foi adicionado com intenção de generalidade. Na prática, todos os 5 implementadores usam `Update` apenas como um wrapper que despacha `tea.KeyMsg` para `HandleKey` e retorna `nil` para qualquer outro tipo de mensagem.

Além disso, a interface já possui `HandleKey(msg tea.KeyMsg) tea.Cmd`, tornando `Update` redundante para o único caso de uso real atual.

## Decisão

Remover `Update` da interface `ModalView` e adicionar `HandleMouse(msg tea.MouseMsg) tea.Cmd` diretamente na interface principal, pois todos os modais responderão a teclado e mouse.

## Interface resultante

```go
// ModalView é a interface principal implementada por todos os modais.
type ModalView interface {
    Render(maxHeight, maxWidth int, theme *design.Theme) string
    HandleKey(msg tea.KeyMsg) tea.Cmd
    HandleMouse(msg tea.MouseMsg) tea.Cmd
    Cursor(topY, leftX int) *tea.Cursor
}
```

## Mudanças necessárias

### `internal/tui/modal.go`
- Remover `Update(msg tea.Msg) tea.Cmd` da interface `ModalView`
- Adicionar `HandleMouse(msg tea.MouseMsg) tea.Cmd` à interface `ModalView`

### Implementadores (5 modais em `internal/tui/modal/`)
- `confirm_modal.go` — remover `Update`, adicionar `HandleMouse` retornando `nil`
- `help_modal.go` — remover `Update`, adicionar `HandleMouse` retornando `nil`
- `password_create_modal.go` — remover `Update`, adicionar `HandleMouse` retornando `nil`
- `password_entry_modal.go` — remover `Update`, adicionar `HandleMouse` retornando `nil`
- `file_picker.go` — remover `Update`, adicionar `HandleMouse` retornando `nil`

### `internal/tui/root.go`
- Substituir todas as chamadas a `.Update(msg)` nos modais por `.HandleKey(msg.(tea.KeyMsg))` onde a mensagem for `tea.KeyMsg`
- Para os pontos onde `Update` era chamado com mensagens não-key (ex: `ModalReadyMsg`), remover o dispatch para o modal — modais não precisam receber esse tipo de mensagem

## Motivação

- **YAGNI:** `Update` genérico não agrega valor enquanto modais só reagem a key e mouse
- **Clareza semântica:** `HandleKey` e `HandleMouse` tornam explícito o contrato da interface
- **Sem custo futuro:** quando mouse for implementado, o método já existe na interface — basta substituir o `return nil` pela lógica real
