# Design — ActionLineView

**Data:** 2026-04-15  
**Contexto:** `internal/tui/screen/action_view.go`  
**Spec de referência:** `golden/tui-spec-barras.md#barra-de-comandos`

---

## Objetivo

Implementar a renderização visual da barra de comandos (`ActionLineView`) conforme a especificação do design system, com testes golden file no mesmo padrão de `message_view_test.go`.

---

## Prerequisito: Reorganização de tipos

### Problema

`tui/action.go` define `Action` e `ActionGroup` no package `tui`. O package `tui` importa `screen`. Se `screen` importar `tui` para acessar `Action`, há import cycle.

O alias `type Action = interface{}` em `screen/types.go` é o workaround atual — impede tipagem real em `ActionLineView.Render`.

### Solução

Mover `Action`, `ActionGroup` e `AppState` de `tui/action.go` para `tui/actions/action.go`.

**Novo grafo de dependência:**

```
design ← actions ← tui (root, view)
design ← actions ← screen (action_view)
```

**Arquivos afetados:**

| Arquivo | Mudança |
|---|---|
| `tui/actions/action.go` | **Novo** — recebe Action, ActionGroup, AppState movidos de tui/action.go |
| `tui/action.go` | **Deletado** — conteúdo movido para actions/action.go |
| `tui/root.go` | Importa `actions`; usa `actions.Action`, `actions.ActionGroup` |
| `tui/view.go` | `Actions() []actions.Action` (antes era `[]interface{}`) |
| `screen/types.go` | Remove alias `type Action = interface{}`; atualiza `ChildView.Actions()` |
| `screen/action_view.go` | `Render` recebe `[]actions.Action` |

### Interface `ChildView` em `actions/action.go`

`Action.AvailableWhen` tem assinatura `func(app AppState, view ChildView) bool`. Para que `actions` possa definir esse tipo sem importar `tui` ou `screen`, define-se uma interface mínima local:

```go
// ChildView é o subconjunto da interface tui.ChildView necessário para AvailableWhen.
// Definida aqui para evitar import cycle. Nenhuma action atual inspeciona a view —
// a interface está vazia agora, mas nomear o tipo documenta a intenção e permite
// adicionar métodos futuramente sem alterar a assinatura de AvailableWhen.
type ChildView interface{}
```

Funções `AvailableWhen` que precisarem inspecionar a view fazem type assertion para o tipo concreto.

### `Visible bool` — semântica adotada

O campo `Visible bool` cobre dois casos do design system:

- **`HideFromBar`:** ação funciona mas não aparece na barra (ex: `F12` tema — `Visible: false`)
- **Indisponível:** ação não satisfaz `AvailableWhen` — removida da barra (não fica dim)

Não existe campo `Enabled` separado neste escopo. Ações indisponíveis desaparecem da barra; o modal de Ajuda as exibe através de `ActiveViewActions()` que não filtra. Esta decisão é deliberada e documentada aqui — uma revisão futura pode introduzir `Enabled` quando o modal de Ajuda tiver implementação real.

---

## ActionLineView — Contrato

### Tipo

```go
type ActionLineView struct{} // stateless — zero value é válido
```

Não há estado interno. Toda informação necessária para renderizar é passada como argumento para `Render`.

### Assinatura de Render

```go
func (v *ActionLineView) Render(width int, theme *design.Theme, actions []actions.Action) string
```

- `width` — largura total disponível (100% do terminal)
- `theme` — tema ativo para tokens de cor
- `actions` — lista pré-filtrada e ordenada fornecida por `ActiveViewActionsForBar()`

O parâmetro `height` é removido em relação ao stub atual — a barra é sempre 1 linha.

### Identificação da âncora F1

`ActionLineView` extrai a âncora identificando a action cujo `Keys[0].Code` e `Keys[0].Mod` coincidem com `design.Shortcuts.Help` (F1). Essa action é separada da lista e sempre renderizada fixada à direita. As demais ações são renderizadas à esquerda.

Se a âncora não estiver na lista, os 8 colunas da direita são preenchidos com espaços — a barra sempre ocupa exatamente `width` colunas e nunca é corrompida.

---

## Algoritmo de Layout

Conforme `tui-spec-barras.md — Anatomia e Dimensionamento`:

```
[2 espaços][ação₁][ · ][ação₂][ · ]…[N espaços][F1 Ajuda]
```

**Constantes:**

| Elemento | Colunas |
|---|---|
| Prefixo | 2 |
| Âncora `F1 Ajuda` | 8 |
| Preenchimento mínimo | 1 |
| Separador ` · ` | 3 |
| Espaço disponível para ações | `width - 2 - 8 - 1 = width - 11` |

**Largura de uma ação:** `lipgloss.Width(Keys[0].Label) + 1 + lipgloss.Width(Label)` em colunas. Usar `lipgloss.Width` (não `len`) para suportar corretamente Unicode multi-coluna como `⌃` e `⇧`.

**Algoritmo de seleção:**

1. Separar âncora (F1) da lista de ações normais.
2. Calcular espaço disponível: `width - 11`.
3. Iterar ações normais na ordem recebida (prioridade crescente):
   - Calcular largura da ação + separador (se não é a primeira).
   - Se cabe, incluir; senão, parar (ações seguintes de menor prioridade são descartadas).
4. Calcular preenchimento: `width - 2 - larguraTotalAções - 8`.
5. Montar string final.

---

## Identidade Visual

Conforme `tui-spec-barras.md — Identidade Visual`:

| Elemento | Token | Atributo |
|---|---|---|
| Tecla da ação (ex: `⌃S`) | `theme.Accent.Primary` | **bold** |
| Label da ação (ex: `Salvar`) | `theme.Text.Primary` | — |
| Separador ` · ` | `theme.Text.Secondary` | — |
| Âncora `F1` — tecla | `theme.Accent.Primary` | **bold** |
| Âncora `F1` — label `Ajuda` | `theme.Text.Primary` | — |

**Regra — sem hard-code:** nenhum valor de cor, atributo tipográfico ou símbolo pode ser literal no código. Usar exclusivamente:
- Tokens do `*design.Theme` passado como argumento (ex: `theme.Accent.Primary`, `theme.Text.Secondary`)
- Constantes de `design.go` e `symbols.go` (ex: `design.SymHeaderSep`)

O separador ` · ` é composto por `" " + design.SymHeaderSep + " "` — nunca pela string literal `" · "`.

**Regra — sem background:** a barra não define cor de fundo. Background é sempre herdado do terminal, como em `MessageLineView`. Não usar `lipgloss.Style.Background()` na barra.

**Acessibilidade NO_COLOR:** sem cor, a tecla fica em **bold** e o label em texto normal. O espaço entre eles (`⌃S Salvar`) é separador suficiente. Padrão consistente com o restante da aplicação — a barra não tem estado crítico que dependa exclusivamente de cor.

---

## Ajustes no Root Model

### `ActiveViewActions()` — mudança de comportamento

O método **já existe** e é chamado em dois lugares:
- `root.go` (dentro de `View()`) — será substituído por `ActiveViewActionsForBar()`
- `actions/setup.go` — dentro do `OnExecute` da action F1, para montar o modal de Ajuda

**Mudanças:**

1. Tipo de retorno: `[]Action` → `[]actions.Action` (consequência da reorganização de tipos)
2. Passa a incluir também `activeView.Actions()`, combinando as três fontes:
   - `systemActions + applicationActions + activeView.Actions()`

O modal de Ajuda é atualmente um stub (`Render` retorna apenas `"Ajuda"`). Não há comportamento real a quebrar — a mudança é segura. Quando o modal for implementado, ele já receberá as ações completas.

### `ActiveViewActionsForBar()` — novo método

Novo método que filtra e ordena a lista para uso exclusivo da barra de comandos:

```go
func (r *RootModel) ActiveViewActionsForBar() []actions.Action {
    // combinar: ActiveViewActions() já retorna as três fontes
    // filtrar: Visible == true E (AvailableWhen == nil OU AvailableWhen satisfeita)
    // ordenar: por Priority crescente
    // retornar
}
```

### `View()` — chamada de Render

```go
r.actionLineView.Render(r.width, r.theme, r.ActiveViewActionsForBar())
```

O bloco de conversão `[]interface{}` existente é removido.

### `ChildView.Actions()` — tipagem

```go
Actions() []actions.Action
```

Todas as implementações de `ChildView` precisam atualizar a assinatura.

---

## Testes

### Testes unitários

- `TestActionLineView_ZeroValue` — Render de lista vazia produz string com largura correta (sem pânico)
- `TestActionLineView_Render_Width` — largura exata (`lipgloss.Width == width`) com diversas combinações de ações
- `TestActionLineView_Render_NoNewline` — resultado nunca contém `\n`
- `TestActionLineView_Render_F1IsAnchor` — F1 aparece quando presente na lista
- `TestActionLineView_Render_NoF1_WidthPreserved` — sem âncora F1, largura ainda é `width`
- `TestActionLineView_Render_TruncatesLowPriority` — ações de menor prioridade caem primeiro quando falta espaço

### Testes golden

Tamanho: `"80x1"` (consistente com barra de 1 linha).

| Variante | Descrição |
|---|---|
| `empty` | Lista vazia — sem ações, sem âncora F1 |
| `single-action` | Uma ação + F1 |
| `multiple-actions` | Três ações + F1 |
| `overflow` | Ações que excedem a largura — truncamento por prioridade |
| `no-f1` | Lista sem ação F1 — âncora ausente, largura preservada |

Golden files em `screen/testdata/golden/actions-*.golden.{txt,json}`.

---

## design_action.go — Helper de renderização

Criar `internal/tui/design/design_action.go` no package `design`, análogo a `design_message.go`.

### Responsabilidade

Encapsula a renderização de uma única ação em texto estilizado, desacoplando a lógica visual de `ActionLineView`. Segue as mesmas regras de identidade visual: sem hard-code, tokens do tema, constantes de symbols.go.

### API

```go
// RenderedAction encapsula o texto ANSI estilizado de uma ação e sua largura em colunas.
type RenderedAction struct {
    Text  string // texto com sequências ANSI
    Width int    // largura visual em colunas (lipgloss.Width — nunca len)
}

// RenderAction renderiza uma ação: tecla (theme.Accent.Primary + bold) + espaço + rótulo (theme.Text.Primary).
func RenderAction(key, label string, theme *Theme) RenderedAction

// ActionSeparator retorna " · " (espaço + SymHeaderSep + espaço) em theme.Text.Secondary.
// Sempre tem Width == 3.
func ActionSeparator(theme *Theme) RenderedAction
```

`ActionLineView.Render` usa `RenderAction` e `ActionSeparator` para montar o layout.

### Alinhamento com design_message.go

- Mesmo package `design`, mesmo padrão de arquivo isolado por domínio
- `RenderedAction` análogo ao par `(string, int)` retornado por `Message.Render`
- Testes unitários em `design_action_test.go`

---

## Não-escopo

- Não implementar interação por mouse (barra é passiva)
- Não implementar estado de hover ou foco na barra
- Não implementar suporte a grupos de ações na barra (grupos são apenas para o modal de Ajuda)
- Não adicionar animação ou transição
- Não implementar o campo `Enabled` separado de `Visible` (decisão registrada acima)
