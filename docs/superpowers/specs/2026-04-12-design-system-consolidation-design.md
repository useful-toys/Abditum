# Design — Consolidação do Design System em `design.go`

**Data:** 2026-04-12
**Pacote:** `internal/tui`
**Arquivo destino:** `internal/tui/design.go`

---

## Contexto e Motivação

O pacote `internal/tui` acumulou quatro arquivos concorrentes de design system:

| Arquivo | Pacote | Problema |
|---|---|---|
| `tui/theme.go` | `tui` | `Theme` com `color.Color` (interface), falta `TextLink`, `BorderFocused`, `SemanticOff`, `Special.*` |
| `tui/tokens.go` | `tui` | Tokens globais fixos (Tokyo Night hardcoded), `MsgKind` definido aqui também |
| `tui/theme/theme.go` | `theme` | Segunda `Theme` com `string`, falta `SemanticOff`, `Logo`, `SurfaceInput`; valores Cyberpunk divergentes |
| `tui/tokens/tokens.go` | `tokens` | Terceiro conjunto de tokens, cores por tema, falta `BorderFocused` e `SurfaceInput` |

Adicionalmente, `MsgKind` está definido **três vezes** com iota incompatível:
- `tui/messages.go`: `MsgSuccess=0`, usa `MsgWarn`
- `tui/types/types.go`: `MsgInfo=0`, usa `MsgWarning`
- `tui/tokens.go`: referencia o tipo do mesmo pacote

Isso causa bugs silenciosos quando valores de `types.MsgKind` são passados onde se espera `tui.MsgKind`.

A fonte de verdade é `golden/tui-design-system.md`. Os valores nos arquivos atuais divergem do documento em múltiplos tokens Cyberpunk.

---

## Decisões de Design

| Decisão | Escolha | Razão |
|---|---|---|
| Localização | `internal/tui/design.go`, pacote `tui` | Zero imports adicionais para componentes existentes |
| Estrutura do `Theme` | Hierárquica (structs aninhadas por categoria) | Espelha taxonomia do design doc; autocompletar por categoria |
| Tipos dos campos | `string` (hex lipgloss) | Compatível com `lipgloss.Color()`; sem dependência de `image/color` |
| `MessageKind` | Definido em `design.go`, único no projeto | Remove triplicação; é fundamentalmente um conceito visual |
| Constantes `Message*` | Prefixo longo (`MessageSuccess`, `MessageWarning`...) | Consistência com o tipo `MessageKind` |
| Acesso a cores | Via `*Theme` passado como parâmetro | Sem globals mutáveis; tipagem completa; tema-aware por construção |
| Tipografia | `Typography` struct em `Theme` | Contratos transversais centralizados; permite desabilitar italic globalmente |
| Notação de teclado | Constantes `KeyCtrl`, `KeyShift`, `KeyAlt` | Barra de comandos e Help renderizam essas strings; centralizar evita divergência |

---

## Estrutura de `design.go`

### 1. Tokens de Cor — `Theme` hierárquica

```go
type SurfaceTokens struct {
    Base   string // fundo da tela inteira
    Raised string // painéis laterais e janelas sobrepostas
    Input  string // campos de texto dentro de diálogos
}

type TextTokens struct {
    Primary   string // nomes, títulos, conteúdo legível
    Secondary string // apoio, hints, placeholders
    Disabled  string // opções inativas
    Link      string // URLs e referências externas
}

type BorderTokens struct {
    Default string // divisores de painel, bordas de modais neutros
    Focused string // painel ativo, campos de entrada, modais com foco
}

type AccentTokens struct {
    Primary   string // barra de seleção, cursor, botão principal
    Secondary string // favorito ★, nomes de pastas
}

type SemanticTokens struct {
    Success string // operação concluída, config ON
    Warning string // alerta, estado dirty (✦ ✎ ✗)
    Error   string // erro, senha incorreta, borda destrutiva
    Info    string // informação contextual
    Off     string // config OFF
}

type SpecialTokens struct {
    Muted     string // texto apagado sem conotação semântica
    Highlight string // fundo atrás do item selecionado na lista
    Match     string // trecho que corresponde ao termo de busca
}

type Typography struct {
    Bold          bool // universal — títulos, item selecionado, ação default
    Dim           bool // itens desabilitados, conteúdo secundário
    Italic        bool // hints, pastas virtuais, placeholders
    Underline     bool // uso pontual
    Strikethrough bool // itens marcados para exclusão (par: SymDeleted + Special.Muted)
    // Blink: ausente por definição — não usar
}

type Theme struct {
    Surface    SurfaceTokens
    Text       TextTokens
    Border     BorderTokens
    Accent     AccentTokens
    Semantic   SemanticTokens
    Special    SpecialTokens
    Logo       [5]string  // gradiente do logo, linha 0 (topo) a linha 4 (base)
    Typography Typography
}
```

### 2. Instâncias dos Temas

Valores extraídos diretamente de `golden/tui-design-system.md` (fonte de verdade).

```go
var DefaultTypography = Typography{
    Bold: true, Dim: true, Italic: true,
    Underline: true, Strikethrough: true,
}

var TokyoNight = &Theme{
    Surface:  SurfaceTokens{"#1a1b26", "#24283b", "#1e1f2e"},
    Text:     TextTokens{"#a9b1d6", "#565f89", "#3b4261", "#7aa2f7"},
    Border:   BorderTokens{"#414868", "#7aa2f7"},
    Accent:   AccentTokens{"#7aa2f7", "#bb9af7"},
    Semantic: SemanticTokens{"#9ece6a", "#e0af68", "#f7768e", "#7dcfff", "#737aa2"},
    Special:  SpecialTokens{"#8690b5", "#283457", "#f7c67a"},
    Logo:     [5]string{"#9d7cd8", "#89ddff", "#7aa2f7", "#7dcfff", "#bb9af7"},
    Typography: DefaultTypography,
}

var Cyberpunk = &Theme{
    Surface:  SurfaceTokens{"#0a0a1a", "#1a1a2e", "#0e0e22"},
    Text:     TextTokens{"#e0e0ff", "#8888aa", "#444466", "#ff2975"},
    Border:   BorderTokens{"#3a3a5c", "#ff2975"},
    Accent:   AccentTokens{"#ff2975", "#00fff5"},
    Semantic: SemanticTokens{"#05ffa1", "#ffe900", "#ff3860", "#00b4d8", "#9999cc"},
    Special:  SpecialTokens{"#666688", "#2a1533", "#ffc107"},
    Logo:     [5]string{"#ff2975", "#b026ff", "#00fff5", "#05ffa1", "#ff2975"},
    Typography: DefaultTypography,
}
```

### 3. `MessageKind` — definição única

```go
type MessageKind int

const (
    MessageSuccess MessageKind = iota
    MessageInfo
    MessageWarning
    MessageError
    MessageBusy
    MessageHint
)

func SymbolFor(kind MessageKind) string  // retorna símbolo canônico
func ColorFor(t *Theme, kind MessageKind) string  // retorna hex da cor semântica
```

### 4. Símbolos — inventário completo (32 constantes + spinner)

Cobertura 1:1 com o inventário de `golden/tui-design-system.md`. Grupos:

- **Árvore:** `SymFolderCollapsed`, `SymFolderExpanded`, `SymFolderEmpty`, `SymLeaf`
- **Estado de item:** `SymFavorite`, `SymDeleted`, `SymCreated`, `SymModified`
- **Mensagens:** `SymSuccess`, `SymInfo`, `SymWarning`, `SymError`
- **UI:** `SymSensitiveField`, `SymSensMask` (mesmo glifo de `SymBullet`, semântica distinta), `SymCursor`, `SymScrollUp`, `SymScrollDown`, `SymScrollThumb`, `SymEllipsis`, `SymBullet`, `SymHeaderSep`, `SymTreeConnector`
- **Bordas:** `SymBorderH`, `SymBorderV`, `SymCornerTL`, `SymCornerTR`, `SymCornerBL`, `SymCornerBR`, `SymJunctionL`, `SymJunctionT`, `SymJunctionB`, `SymJunctionR`
- **Spinner:** `var SpinnerFrames = [4]string{"◐", "◓", "◑", "◒"}`

### 5. Notação de teclado

```go
const (
    KeyCtrl  = "⌃" // U+2303
    KeyShift = "⇧" // U+21E7
    KeyAlt   = "!" // sem Unicode dedicado — decisão do design system
)
```

---

## Cobertura do Design System

| Categoria | Itens | Cobertura |
|---|---|---|
| Tokens de cor | 20/20 | 100% |
| Valores hex | 48/48 sem divergência | 100% |
| Símbolos | 32 constantes + spinner | 100% |
| Tipografia ANSI | 5 atributos (Bold/Dim/Italic/Underline/Strikethrough) | 100% |
| Notação de teclado | 3 modificadores (Ctrl/Shift/Alt) | 100% |

---

## Estratégia de Migração

### Arquivos removidos

| Arquivo | Motivo |
|---|---|
| `internal/tui/theme.go` | Substituído por `design.go` |
| `internal/tui/tokens.go` | Substituído por `design.go` |
| `internal/tui/theme/theme.go` | Pacote `theme` removido |
| `internal/tui/tokens/tokens.go` | Pacote `tokens` removido |

### Arquivos mantidos com cirurgia mínima

| Arquivo | Mudança |
|---|---|
| `internal/tui/types/types.go` | Remove `MsgKind` e suas constantes; mantém `WorkArea`, `Action`, `Scope`, `ToggleThemeMsg`, `FilePickerMode` |
| `internal/tui/messages.go` | Remove definição de `MsgKind`; atualiza referências para `MessageKind`, `MessageSuccess`, `MessageWarning` etc. |

### Renomeações em componentes

| De | Para |
|---|---|
| `tui.MsgKind` / `types.MsgKind` | `tui.MessageKind` |
| `MsgWarn` | `MessageWarning` |
| `MsgSuccess` | `MessageSuccess` |
| `MsgInfo` | `MessageInfo` |
| `MsgError` | `MessageError` |
| `MsgBusy` | `MessageBusy` |
| `MsgHint` | `MessageHint` |
| `tui.ThemeTokyoNight` | `tui.TokyoNight` |
| `tui.ThemeCyberpunk` | `tui.Cyberpunk` |
| `theme.SurfaceBase` | `theme.Surface.Base` |
| `theme.SurfaceRaised` | `theme.Surface.Raised` |
| `theme.TextPrimary` | `theme.Text.Primary` |
| `theme.TextSecondary` | `theme.Text.Secondary` |
| `theme.TextDisabled` | `theme.Text.Disabled` |
| `theme.BorderDefault` | `theme.Border.Default` |
| `theme.AccentPrimary` | `theme.Accent.Primary` |
| `theme.AccentSecondary` | `theme.Accent.Secondary` |
| `theme.SemanticSuccess` | `theme.Semantic.Success` |
| `theme.SemanticWarning` | `theme.Semantic.Warning` |
| `theme.SemanticError` | `theme.Semantic.Error` |
| `theme.SemanticInfo` | `theme.Semantic.Info` |
| `theme.SemanticOff` | `theme.Semantic.Off` |
| `theme.LogoGradient[n]` | `theme.Logo[n]` |

### Hardcoding residual a corrigir

`filepicker.go` linha ~824 e ~891: `lipgloss.Color("#3d59a1")` → `theme.Special.Highlight`

### Diretórios não tocados

`actions/`, `header/`, `interfaces/`, `keys/`, `messagemanager/`, `dialogs/`, `textinput/`, `welcome/` — placeholders para refatorações futuras, fora do escopo deste trabalho.

---

## Validação

Após implementação:

1. `go build ./internal/tui/...` deve compilar sem erros
2. `go vet ./internal/tui/...` deve passar sem warnings
3. Busca por `MsgWarn\b` deve retornar zero resultados
4. Busca por `ThemeTokyoNight\|ThemeCyberpunk` deve retornar zero resultados
5. Busca por `tui/theme"` e `tui/tokens"` nos imports deve retornar zero resultados
6. Todos os 20 tokens de cor presentes em `design.go` devem bater com `golden/tui-design-system.md`
