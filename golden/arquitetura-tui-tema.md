# Arquitetura de Gerenciamento de Tema — `internal/tui`

> Descreve os padrões, convenções e decisões arquiteturais que governam o sistema de tema visual da TUI do Abditum.

---

## Sumário

- [Visão Geral](#visão-geral)
- [Estrutura de Arquivos](#estrutura-de-arquivos)
- [Modelo de Dados do Tema](#modelo-de-dados-do-tema)
  - [`Theme` (struct principal)](#theme-struct-principal)
  - [Subpacotes paralelos](#subpacotes-paralelos)
- [Tokens Estáticos (`tokens.go`)](#tokens-estáticos-tokensgo)
- [Temas Disponíveis](#temas-disponíveis)
- [Ciclo de Vida do Tema no `rootModel`](#ciclo-de-vida-do-tema-no-rootmodel)
  - [Inicialização](#inicialização)
  - [Alternância em Runtime (`F12`)](#alternância-em-runtime-f12)
  - [Propagação via `applyTheme()`](#propagação-via-applytheme)
- [Contrato da Interface `childModel`](#contrato-da-interface-childmodel)
- [Modais e o Protocolo Opcional de Tema](#modais-e-o-protocolo-opcional-de-tema)
- [Passagem de Tema por Injeção de Dependência](#passagem-de-tema-por-injeção-de-dependência)
  - [Construtores de filhos](#construtores-de-filhos)
  - [Construtores de flows](#construtores-de-flows)
  - [Fábricas de diálogos](#fábricas-de-diálogos)
- [Uso do Tema na Renderização](#uso-do-tema-na-renderização)
  - [Componentes stateless (header, message bar, command bar)](#componentes-stateless-header-message-bar-command-bar)
  - [Modelos com estado próprio (welcome, vaultTree, settings…)](#modelos-com-estado-próprio-welcome-vaulttree-settings)
  - [Background global do terminal](#background-global-do-terminal)
- [Divisão entre Tokens Estáticos e Tokens Dinâmicos](#divisão-entre-tokens-estáticos-e-tokens-dinâmicos)
- [Zona de Tensão: Duplicação de Tokens](#zona-de-tensão-duplicação-de-tokens)
- [Gradiente do Logo](#gradiente-do-logo)
- [Design System como Fonte de Verdade](#design-system-como-fonte-de-verdade)
- [Decisões Arquiteturais Relevantes](#decisões-arquiteturais-relevantes)

---

## Visão Geral

O sistema de tema da TUI do Abditum segue uma arquitetura de **tema como dado compartilhado**: um ponteiro `*Theme` é criado pelo `rootModel` na inicialização, propagado para todos os filhos vivos via injeção de dependência nos construtores e, quando o tema muda em runtime, redistribuído via chamada imperativa `applyTheme()`.

O tema não é um objeto global: ele nunca usa variáveis de pacote mutáveis nem singletons. O `rootModel` é o único proprietário do tema ativo.

---

## Estrutura de Arquivos

| Arquivo | Responsabilidade |
|---|---|
| `theme.go` | Define `Theme` (struct), instâncias `ThemeTokyoNight` e `ThemeCyberpunk`, método stub `ApplyTheme(childModel)` |
| `tokens.go` | Constantes de cor e símbolo, funções `Style*()` para estilos lipgloss pré-compostos — usam valores **hardcoded** do Tokyo Night |
| `tokens/tokens.go` | Subpacote com as mesmas constantes + variantes Cyberpunk separadas por nome, referenciadas por `theme/theme.go` |
| `theme/theme.go` | Subpacote com `Theme` alternativa (campos `string` em vez de `color.Color`) + temas pré-construídos; **atualmente não é o struct usado pelo restante do pacote** |
| `types/types.go` | Subpacote com tipos compartilhados (`MsgKind`, `WorkArea`, `FilePickerMode`)  |
| `ascii.go` | `RenderLogo(*Theme)` — renderiza o wordmark em gradiente consumindo `Theme.LogoGradient` |
| `messages.go` | `RenderMessageBar(*DisplayMessage, int, *Theme)` — renderiza a barra de status consumindo tokens do `*Theme` |
| `header.go` | `headerModel.Render(..., *Theme)` — recebe o tema como parâmetro a cada frame |
| `root.go` | `rootModel` — proprietário do `*Theme`, responsável por inicializar, alternar e propagar |

---

## Modelo de Dados do Tema

### `Theme` (struct principal)

Definida em `theme.go`, é o contrato visual central do pacote:

```go
type Theme struct {
    // Superfícies
    SurfaceBase   color.Color
    SurfaceRaised color.Color

    // Acentos
    AccentPrimary   color.Color
    AccentSecondary color.Color

    // Texto
    TextPrimary   color.Color
    TextSecondary color.Color
    TextDisabled  color.Color

    // Semânticas
    SemanticSuccess color.Color
    SemanticWarning color.Color
    SemanticError   color.Color
    SemanticInfo    color.Color
    SemanticOff     color.Color

    // Gradiente do logo (exatamente 5 cores)
    LogoGradient []color.Color
}
```

**Decisão de tipo:** os campos usam `color.Color` (interface padrão da stdlib `image/color`), preenchidos com `lipgloss.Color(hexString)`. Isso mantém compatibilidade direta com a API do lipgloss v2, que aceita `color.Color` como argumento de `Foreground()` e `Background()`.

**Decidido não incluir:**
- `border.default` — está hardcoded em `tokens.go` como `ColorBorderDefault`  
- `surface.input` — está hardcoded em `tokens.go` como `ColorSurfaceInput`  
- `text.link`, `special.*`, `border.focused` — presentes no Design System mas ainda não migrados para o struct

### Subpacotes paralelos

Existem dois subpacotes (`tokens/` e `theme/`) com definições redundantes de cores e temas. São artefatos de refatoração inacabada — o pacote principal `tui` ainda consome o `Theme` definido em `theme.go` (raiz do pacote), não os subpacotes. Os subpacotes estão isolados e não são usados por nenhuma importação ativa do código de produção.

---

## Tokens Estáticos (`tokens.go`)

`tokens.go` define constantes `const` de cor como strings hexadecimais brutas, funções `Style*()` que retornam `lipgloss.Style` pré-compostos, e constantes de símbolos Unicode:

```go
const (
    ColorSuccess = "#9ece6a"
    ColorError   = "#f7768e"
    ColorBorderFocused = "#7aa2f7"
    ColorSurfaceInput  = "#1e1f2e"
    // ...
)

func StyleSymbol(kind MsgKind) lipgloss.Style { ... }
func StyleCommandKey() lipgloss.Style { ... }
```

**Propósito:** centralizar todos os valores absolutos de cor para zero hardcoding nos consumidores. Por convenção, consumidores que **não recebem `*Theme`** (e.g. `passwordEntryModal.View()`) usam estas constantes diretamente.

**Limitação:** as constantes estão fixadas no Tokyo Night. Isso significa que `passwordEntryModal`, `passwordCreateModal` e `modalModel` **não respondem à alternância de tema**, pois não recebem `*Theme` e usam tokens fixos. Esta inconsistência está documentada implicitamente pela separação arquitetural entre modais que recebem `*Theme` e modais que não recebem.

---

## Temas Disponíveis

| Identidade | Estética | Superfície base | Acento primário |
|---|---|---|---|
| `ThemeTokyoNight` | Dark azul-roxo, contido | `#1a1b26` | `#7aa2f7` (azul) |
| `ThemeCyberpunk` | Dark negro, neon saturado | `#0a0a1a` | `#ff2975` (magenta) |

Ambos são variáveis `var` de pacote com valor `*Theme`. São imutáveis por convenção — nunca são modificados após criação. A alternância de tema consiste em mudar qual ponteiro o `rootModel` armazena em `.theme`.

---

## Ciclo de Vida do Tema no `rootModel`

### Inicialização

```go
// newRootModel()
m := &rootModel{
    theme: ThemeTokyoNight, // tema padrão
}
m.welcome = newWelcomeModel(actions, m.theme, m.version)
```

O Tokyo Night é o tema padrão. O `welcome` (primeiro filho criado) recebe o `*Theme` já no construtor.

### Alternância em Runtime (`F12`)

A alternância é tratada diretamente no método `Update()` do `rootModel`, **antes** do dispatch normal de ações:

```go
case tea.KeyPressMsg:
    if key == "f12" {
        if m.theme == ThemeTokyoNight {
            m.theme = ThemeCyberpunk
        } else {
            m.theme = ThemeTokyoNight
        }
        m.applyTheme()
        return m, nil
    }
```

**Decisão:** F12 interceptado antes do `ActionManager.Dispatch()`. Embora também registrado como `Action` no `ActionManager` (via `toggleThemeMsg`), o handler real opera diretamente no Update. O tratamento via mensagem (`toggleThemeMsg`) é emitido pelo handler da ação mas também capturado aqui — na prática a lógica real de troca está inline no Update, não no handler da ação. A `toggleThemeMsg` é usada como mecanismo de ação registrado para exibição na barra de ajuda.

A alternância é binária (toggle entre dois temas), sem suporte a temas arbitrários em runtime.

### Propagação via `applyTheme()`

```go
func (m *rootModel) applyTheme() {
    for _, child := range m.liveWorkChildren() {
        child.ApplyTheme(m.theme)
    }
    for _, modal := range m.modals {
        if themeableModal, ok := modal.(interface{ ApplyTheme(*Theme) }); ok {
            themeableModal.ApplyTheme(m.theme)
        }
    }
}
```

**Modelo pull vs. push:** o rootModel *empurra* o tema para todos os filhos vivos imediatamente após a troca. Os filhos não consultam o tema — recebem uma referência direta via campo `theme *Theme`.

**Filhos de área de trabalho:** todos implementam o método `ApplyTheme(*Theme)` como parte do contrato `childModel`. A chamada é uniforme via interface.

**Modais:** tratados separadamente. O protocolo é opcional — modais que não precisam responder ao tema (e.g. `modalModel` genérico, `helpModal`) simplesmente não implementam `ApplyTheme`. A asserção de tipo `interface{ ApplyTheme(*Theme) }` garante ausência de pânico.

**Componentes stateless:** `headerModel`, `RenderMessageBar` e `RenderCommandBar` **não armazenam tema**. Recebem `*Theme` como parâmetro a cada `View()`. Quando o tema muda, o próximo frame já usa o novo tema automaticamente — sem necessidade de propagação explícita.

---

## Contrato da Interface `childModel`

```go
type childModel interface {
    Update(tea.Msg) tea.Cmd
    View() string
    SetSize(w, h int)
    ApplyTheme(*Theme)  // obrigatório para todos os filhos de área de trabalho
}
```

**`ApplyTheme` é obrigatório** no contrato `childModel`. Isso garante que nenhum filho de área de trabalho pode esquecer de lidar com alternância de tema — o compilador força a implementação.

Cada filho implementa o método com corpo mínimo:

```go
func (m *welcomeModel) ApplyTheme(t *Theme) {
    m.theme = t
}
```

A referência local `m.theme` é atualizada; o próximo `View()` renderiza com o novo tema.

---

## Modais e o Protocolo Opcional de Tema

A interface `modalView` **não inclui** `ApplyTheme`:

```go
type modalView interface {
    Update(tea.Msg) tea.Cmd
    View() string
    Shortcuts() []Shortcut
    SetSize(w, h int)
}
```

Modais que precisam de tema (e.g. `filePickerModal`) recebem o `*Theme` no construtor e armazenam internamente. Para responder à alternância em runtime, precisam opcionalmente implementar `ApplyTheme(*Theme)`, detectado via type assertion em `applyTheme()`.

Modais que **não** implementam o protocolo (e.g. `modalModel`, `passwordEntryModal`, `passwordCreateModal`) continuam usando as constantes hardcoded de `tokens.go` e não respondem à troca de tema.

Essa assimetria é uma inconsistência arquitetural conhecida: `filePickerModal` é tema-aware, os modais de senha não são.

---

## Passagem de Tema por Injeção de Dependência

O tema nunca é acessado via global. Ele flui sempre como parâmetro explícito.

### Construtores de filhos

```go
func newWelcomeModel(actions *ActionManager, theme *Theme, version string) *welcomeModel
func newVaultTreeModel(mgr *vault.Manager, actions *ActionManager, msgs *MessageManager, theme *Theme) *vaultTreeModel
func newSecretDetailModel(mgr *vault.Manager, actions *ActionManager, msgs *MessageManager, theme *Theme) *secretDetailModel
func newSettingsModel(mgr *vault.Manager, actions *ActionManager, msgs *MessageManager, theme *Theme) *settingsModel
```

Padrão uniforme: `*Theme` é sempre o último parâmetro antes dos opcionais, após os serviços.

### Construtores de flows

```go
func newOpenVaultFlow(..., theme *Theme) *openVaultFlow
func newCreateVaultFlow(..., theme *Theme) *createVaultFlow
func newSaveAndExitFlow(..., theme *Theme) *saveAndExitFlow
```

Flows também recebem `*Theme` para repassar às fábricas de diálogos que criam durante sua execução.

### Fábricas de diálogos

```go
func FilePicker(title string, mode FilePickerMode, ext string, messages *MessageManager, theme *Theme) tea.Cmd
```

`PasswordEntry` e `PasswordCreate` **não** recebem `*Theme` — usam `ThemeTokyoNight` hardcoded internamente:

```go
func PasswordEntry(title string) tea.Cmd {
    m := &passwordEntryModal{title: title}
    m.theme = ThemeTokyoNight // hardcoded
    ...
}
```

Esta é a principal inconsistência de injeção no sistema atual.

---

## Uso do Tema na Renderização

### Componentes stateless (header, message bar, command bar)

Recebem `*Theme` como parâmetro de renderização, sem armazenamento:

```go
// header.go
func (h *headerModel) Render(width int, vaultName string, isDirty bool, area workArea, theme *Theme) string {
    appNameStyle := lipgloss.NewStyle().Foreground(theme.AccentPrimary).Bold(true)
    vaultNameStyle := lipgloss.NewStyle().Foreground(theme.TextPrimary).Bold(true)
    // ...
}

// messages.go
func RenderMessageBar(msg *DisplayMessage, width int, theme *Theme) string {
    borderStyle := lipgloss.NewStyle().Foreground(theme.SurfaceRaised)
    // colors per msg.Kind: theme.SemanticSuccess, theme.SemanticError, etc.
}

// actions.go
func RenderCommandBar(actions []Action, width int, theme *Theme) string {
    keyStyle := lipgloss.NewStyle().Foreground(theme.AccentPrimary).Bold(true)
    labelStyle := lipgloss.NewStyle().Foreground(theme.TextPrimary)
}
```

**Vantagem:** estes componentes são automaticamente tema-corretos no próximo frame após a alternância, sem necessidade de `ApplyTheme`.

### Modelos com estado próprio (welcome, vaultTree, settings…)

Armazenam `m.theme *Theme` e usam nos seus `View()`:

```go
// welcome.go
func (m *welcomeModel) View() string {
    versionStyle := lipgloss.NewStyle().Foreground(m.theme.TextSecondary)
    logoBlock := lipgloss.NewStyle().Width(43).Render(RenderLogo(m.theme))
    // ...
}
```

### Background global do terminal

O `rootModel.View()` define a cor de fundo do terminal inteiro via `tea.View`:

```go
v := tea.NewView(content)
v.AltScreen = true
v.BackgroundColor = m.theme.SurfaceBase
```

Isso garante que o fundo do Bubble Tea v2 coincide com `surface.base` do tema ativo, sem janelas pretas nas extremidades.

---

## Divisão entre Tokens Estáticos e Tokens Dinâmicos

O sistema mantém duas camadas de definição de cores:

| Camada | Arquivo | Tipo | Mutável em runtime |
|---|---|---|---|
| **Tokens estáticos** | `tokens.go` | `const string` hex | Não |
| **Tokens dinâmicos** | `theme.go` (struct) | `color.Color` por tema | Via `m.theme = ...` |

**Por que existem duas camadas?**

Historicamente, os tokens estáticos foram criados primeiro para centralizar valores e evitar hardcoding espalhado. O struct `Theme` foi introduzido depois para suportar alternância de temas. As constantes de `tokens.go` persistem para componentes que não recebem `*Theme` e para casos de borda onde hardcoding é deliberado (e.g. cores da modal de ajuda usam índices de paleta 256 — `"62"`, `"11"` — em vez de true color, por compatibilidade de terminal).

**Regra de uso:**
- Se o componente recebe `*Theme`: usar campos do struct
- Se o componente não recebe `*Theme` (modal autônoma, utilitário): usar constantes de `tokens.go`

---

## Gradiente do Logo

O wordmark ASCII é renderizado com gradiente de 5 cores, uma por linha:

```go
// ascii.go
func RenderLogo(t *Theme) string {
    colors := t.LogoGradient
    if len(colors) != 5 {
        panic("LogoGradient must contain exactly 5 colors")
    }
    for i, line := range lines {
        style := lipgloss.NewStyle().Foreground(colors[i])
        rendered = append(rendered, style.Render(line))
    }
}
```

**Decisão arquitetural:** o gradiente é parte do `Theme`, não de `tokens.go`. Cada tema tem progressão cromática própria (violeta→ciano para Tokyo Night, magenta→verde para Cyberpunk). A validação do tamanho é feita com `panic` — é um invariante de construção, não um erro de runtime.

---

## Design System como Fonte de Verdade

O documento `tui-design-system-novo.md` define os papéis de cor (e.g. `surface.base`, `accent.primary`) e seus valores por tema. O struct `Theme` mapeia diretamente esses papéis — os nomes dos campos seguem a nomenclatura do design system.

A regra de governança do design system é: **princípio prevalece sobre especificação de tela em conflito**. Isso se reflete na arquitetura: o tema é a única fonte dos valores de cor para componentes tema-aware; nenhum componente armazena alternativas paralelas.

---

## Decisões Arquiteturais Relevantes

| Código | Decisão | Consequência |
|---|---|---|
| D-17 | Todos os hex centralizados — zero hardcoding em consumidores | `tokens.go` e struct `Theme` como únicas fontes de valores de cor |
| D-18 | Paleta organizada por **papel funcional**, não por cor concreta | Trocar tema = trocar ponteiro `*Theme`, sem alterar lógica de renderização |
| — | `ApplyTheme` obrigatório em `childModel` | Compilador garante que filhos de área de trabalho sempre respondem à troca |
| — | Protocolo opcional em `modalView` | Modais legados funcionam sem modificação; novos modais opt-in com type assertion |
| — | Tema como injeção de dependência, não global | Testabilidade: testes instanciam com `ThemeCyberpunk` sem afetar outros testes |
| — | Toggle binário via F12, inline no `Update()` | Simplicidade; sem estado intermediário; sem suporte a temas externos em runtime |
| — | `passwordEntryModal` e `passwordCreateModal` hardcoded em `ThemeTokyoNight` | Inconsistência conhecida: modais de senha não respondem à troca de tema |
| — | Subpacotes `tokens/` e `theme/` existem mas não são consumidos pelo pacote principal | Artefatos de refatoração em andamento; código morto no estado atual |
