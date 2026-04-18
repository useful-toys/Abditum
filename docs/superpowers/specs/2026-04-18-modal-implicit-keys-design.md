# Design: Teclas implícitas Enter/Esc em `ModalOption`

**Data:** 2026-04-18  
**Status:** Aprovado

## Contexto

Diálogos de confirmação são criados via `modal.NewConfirmModal` com uma slice de `[]ModalOption`.
Hoje, cada `ModalOption` exige a declaração explícita de suas teclas:

```go
[]modal.ModalOption{
    {
        Keys:   []design.Key{design.Keys.Enter},
        Label:  "Confirmar",
        Action: func() tea.Cmd { return tea.Batch(tui.CloseModal(), doAction()) },
    },
    {
        Keys:   []design.Key{design.Keys.Esc},
        Label:  "Cancelar",
        Action: func() tea.Cmd { return tui.CloseModal() },
    },
}
```

As regras de Enter/Esc já estão hardcoded em `key_handler.go` para o caso de option única
(`confirm_modal.go` adiciona Esc manualmente). A ideia é generalizar essas regras e eliminar
a repetição na declaração das teclas quando o caller usa apenas Enter/Esc.

## Objetivo

Tornar `ModalOption.Keys` opcional. Quando omitido (nil ou vazio):

- A **primeira** option responde a `Enter`
- A **última** option responde a `Esc`
- Com **uma única** option, ela responde a ambos (`Enter` e `Esc`)

Quando `Keys` está preenchido, as teclas implícitas são **adicionadas como aliases** (não substituídas).
O footer exibe `Keys[0].Label` quando há keys declaradas, ou `Enter`/`Esc` quando `Keys` está vazio.

## Regras de negócio

| Situação | Teclas ativas | Label no footer |
|---|---|---|
| 1ª option, `Keys` vazio | `Enter` (implícito) | `Enter` |
| última option, `Keys` vazio | `Esc` (implícito) | `Esc` |
| option única, `Keys` vazio | `Enter` + `Esc` (implícitos) | `Enter` |
| 1ª option, `Keys: [letter('s')]` | `s` + `Enter` (alias implícito) | `S` (Keys[0]) |
| última option, `Keys: [letter('n')]` | `n` + `Esc` (alias implícito) | `N` (Keys[0]) |
| option intermediária, `Keys` vazio | sem tecla | omitida do footer |

## Arquivos afetados

### 1. `internal/tui/modal/modal_base.go`

Atualizar comentário de `ModalOption.Keys` para documentar que o campo é opcional e descrever
o comportamento implícito de Enter/Esc.

**Mudança:** apenas comentário — sem alteração de struct.

### 2. `internal/tui/modal/key_handler.go` — `KeyHandler.Handle()`

Após o loop de keys explícitas, adicionar despacho implícito:

```
se len(Options) > 0:
    se Enter pressionado:
        executar Options[0].Action()
    se Esc pressionado:
        executar Options[len-1].Action()
```

Essa lógica ocorre **apenas** se o loop de keys explícitas não consumiu o evento.
Resultado: Enter/Esc funcionam sempre como aliases da 1ª e última option,
independentemente de estarem declarados em `Keys`.

**Remoção:** a lógica manual em `confirm_modal.go` que adiciona Esc como alias para option única
pode ser removida — `KeyHandler` passa a cobrir esse caso centralizadamente.

### 3. `internal/tui/modal/frame.go` — `DialogFrame`

Dois locais fazem `if len(opt.Keys) == 0 { continue }`:

- `calculateBodyWidth` (linha 165): cálculo da largura das ações no rodapé
- `renderBottomBorder` (linha 227): renderização das ações no rodapé

Ambos precisam de uma função auxiliar privada:

```go
// implicitKey retorna a tecla implícita de uma option dado seu índice e o total de options.
// Índice 0 → Enter; último índice → Esc; outros → zero value (sem tecla).
func implicitKey(index, total int) (design.Key, bool)
```

Quando `Keys` está vazio, usar `implicitKey(i, total)` como fallback para obter o label e a largura.

## Compatibilidade

- **Retrocompatível:** callers existentes com `Keys: []design.Key{design.Keys.Enter}` continuam
  funcionando. Enter ficará registrado duas vezes (explícito + implícito) — inofensivo.
- **Sem quebra de API:** `ModalOption` mantém a mesma struct. `KeyHandler` e `DialogFrame`
  mantêm as mesmas assinaturas públicas.

## Testes a atualizar/adicionar

- `key_handler_test.go`: casos para option sem Keys (Enter → 1ª, Esc → última, ambos para única)
- `frame_test.go`: caso onde `Keys` está vazio — footer deve mostrar Enter/Esc automaticamente
- `confirm_modal_test.go`: verificar que a remoção da lógica manual de Esc não quebra o caso de option única
