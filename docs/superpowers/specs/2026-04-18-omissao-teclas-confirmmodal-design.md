# Design: Omissão de declaração de teclas ENTER e ESC em ConfirmModal

**Data:** 2026-04-18
**Status:** Aprovado

## Contexto

Atualmente, toda `ModalOption` requer a declaração explícita do slice `Keys` contendo, no mínimo, `design.Keys.Enter` para ações de confirmação e `design.Keys.Esc` para ações de cancelamento. Em diálogos de confirmação padrão (ConfirmModal), esse padrão é tão recorrente que gera verbosidade desnecessária.

## Decisão

Modificar o comportamento do `KeyHandler` (usado por todos os modais que possuem ações de rodapé) e do `frame.go` para tratar implicitamente:
- A **primeira opção** no slice `Options` como ativada por **Enter** (quando `Keys` estiver vazio ou nil)
- A **última opção** no slice `Options` como ativada por **Esc** (quando `Keys` estiver vazio ou nil)

Isso permite que desenvolvedores omita a declaração de `Keys` nos casos comuns, enquanto ainda suporta declarações explícitas quando necessário.

### Comportamento de Teclas (Key Handling)
Para cada opção em `Modal.Options` (onde Modal é qualquer modal que use KeyHandler, incluindo ConfirmModal):
- **Se a opção for a primeira (índice 0)**:
  - Se `Keys` for **não vazio**: usa as teclas declaradas em `Keys` **além de** `design.Keys.Enter` (para matching)
  - Se `Keys` for **vazio ou nil**: trata como tendo apenas `design.Keys.Enter` (para matching)
- **Se a opção for a última (índice len-1)**:
  - Se `Keys` for **não vazio**: usa as teclas declaradas em `Keys` **além de** `design.Keys.Esc` (para matching)
  - Se `Keys` for **vazio ou nil**: trata como tendo apenas `design.Keys.Esc` (para matching)
- **Se a opção for do meio (nem primeira nem última)**:
  - **Deve ter `Keys` não vazio** (caso contrário, não terá nenhuma tecla ativadora)
  - Usa apenas as teclas declaradas em `Keys` (para matching)
- **Se houver apenas uma opção (que é tanto primeira quanto última)**:
  - Se `Keys` for **não vazio**: usa as teclas declaradas em `Keys` **além de** `design.Keys.Enter` e `design.Keys.Esc` (para matching)
  - Se `Keys` for **vazio ou nil**: trata como tendo **ambas** `design.Keys.Enter` e `design.Keys.Esc` (para matching)

### Comportamento de UI (Renderização)
O `frame.go` renderiza a primeira tecla declarada em `opt.Keys[0]`. Com nossas modificações:
- Quando `Keys` é **não vazio**: exibe `opt.Keys[0]` (primeira tecla declarada)
- Quando `Keys` for **vazio ou nil**: 
  - Para primeira opção: exibe `[Enter]`
  - Para última opção: exibe `[Esc]`
  - Para opção única (primeira e última): exibe `[Enter]` (usa Enter como padrão para display quando ambas as teclas estão disponíveis implicitamente)

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

#### Teclas personalizadas (mantém compatibilidade com comportamento aprimorado)
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
    {Keys: []design.Key{design.Keys.F2}, Label: "Ajuda", Action: onHelp}, // F2 obrigatório (nem primeira nem última)
    {Label: "Cancelar", Action: onCancel},             // Esc implícito
}
```

#### Única opção (Enter e Esc)
```go
opts := []ModalOption{
    {Label: "OK", Action: onOK}, // Enter e Esc implícitos, mostra [Enter] OK
}
```

## Arquivos Afetados

### Modificações de Comportamento
- `internal/tui/modal/key_handler.go` — modificar lógica de despacho para tratar primeiro/último opção com teclas implícitas
- `internal/tui/modal/frame.go` — modificar `renderBottomBorder` para exibir teclas apropriadamente quando `Keys` vazio

### Observação de Escopo
Esta mudança afeta **todos** os modais que utilizam `KeyHandler`, pois modifica o comportamento fundamental do despacho de teclas. No entanto, o benefício é mais evidente em `ConfirmModal` devido ao seu padrão comum de uso.

## Critério de Sucesso

- O projeto compila sem erros
- Todos os testes existentes passam (exceto aqueles que dependem exatamente do comportamento anterior de teclas em primeiro/último opção quando `Keys` vazio - estes serão atualizados)
- Novos padrões de uso (omissão de `Keys`) funcionam conforme especificado
- Comportamento aprimorado: mesmo quando `Keys` é declarado com teclas personalizadas, primeira opção responde a Enter e última opção responde a Esc
- Compatibilidade total com código existente: nenhuma quebra introdutiva na funcionalidade existente
