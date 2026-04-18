# Design: Omissão de declaração de teclas ENTER e ESC em ConfirmModal

**Data:** 2026-04-18
**Status:** Aprovado

## Contexto

Atualmente, toda `ModalOption` requer a declaração explícita do slice `Keys` contendo, no mínimo, `design.Keys.Enter` para ações de confirmação e `design.Keys.Esc` para ações de cancelamento. Em diálogos de confirmação padrão (ConfirmModal), esse padrão é tão recorrente que gera verbosidade desnecessária.

## Decisão

Modificar o comportamento do `ConfirmModal` para tratar implicitamente:
- A **primeira opção** no slice `Options` como ativada por **Enter** (independentemente da declaração de `Keys`)
- A **última opção** no slice `Options` como ativada por **Esc** (independentemente da declaração de `Keys`)

Para opções do meio (nem primeira nem última), a declaração explícita de `Keys` permanece necessária.

### Comportamento de Teclas (Key Handling)
- **Primeira opção**: responde a `Enter` **e** a quaisquer teclas declaradas em `Keys` (se presente)
- **Última opção**: responde a `Esc` **e** a quaisquer teclas declaradas em `Keys` (se presente)
- **Opções do meio**: respondem **apenas** às teclas declaradas em `Keys`

### Comportamento de UI (Renderização)
- Se `Keys` for **não vazio**: exibe a **primeira tecla declarada** (ex: `[AltC] Confirmar`)
- Se `Keys` for **vazio ou nil**: 
  - Para primeira opção: exibe `[Enter]` 
  - Para última opção: exibe `[Esc]`
  - Para opção única (primeira e última): exibe `[Enter]`

### Exemplos

#### Antes (verbose)
```go
opts := []ModalOption{
    {Keys: []design.Key{design.Keys.Enter}, Label: "Confirmar", Action: onConfirm},
    {Keys: []design.Key{design.Keys.Esc}, Label: "Cancelar", Action: onCancel},
}
```

#### Depois (conciso)
```go
opts := []ModalOption{
    {Label: "Confirmar", Action: onConfirm}, // Enter implícito
    {Label: "Cancelar", Action: onCancel},   // Esc implícito
}
```

#### Teclas personalizadas (mantém compatibilidade)
```go
opts := []ModalOption{
    {Keys: []design.Key{design.Keys.AltC}, Label: "Confirmar", Action: onConfirm},
    {Label: "Cancelar", Action: onCancel}, // Esc ainda implícito
}
```
→ Primeira opção responde a `AltC` **e** `Enter`, mostra `[AltC] Confirmar`  
→ Segunda opção responde a `Esc`, mostra `[Esc] Cancelar`

#### Opção do meio (requer declaração explícita)
```go
opts := []ModalOption{
    {Label: "Confirmar", Action: onConfirm},           // Enter implícito
    {Keys: []design.Key{design.Keys.F2}, Label: "Ajuda", Action: onHelp}, // F2 obrigatório
    {Label: "Cancelar", Action: onCancel},             // Esc implícito
}
```

## Arquivos Afetados

### Modificações de Comportamento
- `internal/tui/modal/key_handler.go` — modificar lógica de despacho para tratar primeiro/último opção de ConfirmModal specially
- `internal/tui/modal/frame.go` — modificar `renderBottomBorder` para exibir teclas apropriadamente quando `Keys` vazio

### Observação de Escopo
Esta mudança afeta especificamente o `ConfirmModal`. Outros modais (como `PasswordEntryModal`, `HelpModal`, etc.) continuarão funcionando como antes, já que eles não dependem dessa lógica implícita — a menos que sejam refatorados para usar `ConfirmModal` no futuro.

## Critério de Sucesso

- O projeto compila sem erros
- Todos os testes existentes passam (exceto aqueles que dependem explicitamente da declaração de `Keys` em primeiro/último opção — estes serão atualizados)
- Novos padrões de uso (omissão de `Keys`) funcionam conforme especificado
- Compatibilidade total com código existente: nenhuma quebra introdutiva
