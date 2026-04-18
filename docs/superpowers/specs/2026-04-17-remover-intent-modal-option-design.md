# Design: Remover `Intent` de `ModalOption`

**Data:** 2026-04-17
**Status:** Aprovado

## Contexto

`ModalOption` possui um campo `Intent Intent` que classifica semanticamente a intenção de uma opção de modal (`IntentConfirm`, `IntentCancel`, `IntentOther`). Após análise do codebase, constatou-se que esse campo **nunca é lido** em código de produção — `frame.go` e `key_handler.go` consomem apenas `Keys`, `Label` e `Action`. A intenção já está implícita pela posição da opção no slice.

## Decisão

Remover completamente o campo `Intent`, o tipo `Intent int` e todas as constantes associadas. Nenhum comportamento será alterado.

## Arquivos Afetados

### Remoções de tipo e campo

- `internal/tui/modal/modal_base.go`
  - Remove `type Intent int`
  - Remove constantes `IntentConfirm`, `IntentCancel`, `IntentOther`
  - Remove campo `Intent Intent` de `ModalOption`

### Remoção do campo em literais de struct

- `internal/tui/modal/password_entry_modal.go`
- `internal/tui/modal/password_create_modal.go`
- `internal/tui/modal/help_modal.go`
- `internal/tui/operation/fake_operation.go`
- `cmd/test_calc/main.go`

### Testes

- `internal/tui/modal/confirm_modal_test.go` — remove `TestConfirmModal_IntentTypes_Preserved` e todos os campos `Intent:` nos literais
- `internal/tui/modal/frame_test.go` — remove campos `Intent:` nos literais
- `internal/tui/modal/key_handler_test.go` — remove campos `Intent:` nos literais

## Critério de Sucesso

- O projeto compila sem erros após as remoções
- Todos os testes passam (exceto o teste removido)
- Nenhuma lógica de renderização ou despacho de teclas é alterada
