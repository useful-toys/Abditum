# Design Package Completion

## Problema e Abordagem

O pacote `internal/tui/design` define os tokens de cor do `Theme` e os dois temas (TokyoNight, Cyberpunk), mas não cobre todos os elementos que o design system golden (`golden/tui-design-system.md`) especifica. Views que precisarem de símbolos, tipos de teclas ou constantes de layout teriam de hardcodar esses valores.

A solução é completar o pacote `design` com três adições:

1. Extensão de `design.go` — gradiente do logo + constantes de layout
2. Novo `symbols.go` — inventário de símbolos Unicode
3. Novo `keys.go` — tipo `Key` com label e código tea, teclas pré-definidas, funções de composição

## design.go — Adições

### Campo LogoGradient no Theme

O design system define 5 cores por tema para as linhas do logo ASCII art. São cores tema-dependentes que não se encaixam nos tokens funcionais existentes.

```go
type Theme struct {
    // ...campos existentes...

    // LogoGradient são as 5 cores para as linhas do logo ASCII, de cima para baixo.
    LogoGradient [5]string
}
```

Valores:

| Linha | Tokyo Night | Cyberpunk |
|---|---|---|
| 0 | `#9d7cd8` | `#ff2975` |
| 1 | `#89ddff` | `#b026ff` |
| 2 | `#7aa2f7` | `#00fff5` |
| 3 | `#7dcfff` | `#05ffa1` |
| 4 | `#bb9af7` | `#ff2975` |

### Constantes de Layout

```go
const (
    MinWidth       = 80   // largura mínima do terminal em colunas (padrão POSIX)
    MinHeight      = 24   // altura mínima do terminal em linhas (padrão POSIX)
    PanelTreeRatio = 0.35 // proporção do painel de navegação esquerdo (árvore/lista)
)
```

`PanelTreeRatio` define ~35% para o painel esquerdo; o painel direito (detalhe) ocupa o restante (~65%). A implementação pode ajustar ±5%.

## symbols.go — Inventário de Símbolos

Arquivo novo no pacote `design`. Centraliza todos os caracteres Unicode usados pela interface como constantes nomeadas, evitando strings mágicas nas views.

Todos os símbolos são BMP (U+0000–U+FFFF), largura 1 coluna, exceto `SymTreeConnector` (2 colunas, composto).

```go
// Navegação em árvore
const (
    SymFolderCollapsed = "▶" // pasta recolhida
    SymFolderExpanded  = "▼" // pasta expandida
    SymFolderEmpty     = "▷" // pasta vazia
    SymLeaf            = "●" // item folha
)

// Estados de item
const (
    SymFavorite = "★" // item favorito
    SymDeleted  = "✗" // marcado para exclusão
    SymCreated  = "✦" // recém-criado, não salvo
    SymModified = "✎" // modificado, não salvo
)

// Semânticos
const (
    SymSuccess = "✓" // operação concluída
    SymInfo    = "ℹ" // informação contextual
    SymWarning = "⚠" // alerta
    SymError   = "✕" // erro (distinto de SymDeleted ✗)
)

// Elementos de UI
const (
    SymRevealable   = "◉"  // indicador de campo revelável
    SymMask         = "•"  // máscara de conteúdo sensível (mesmo glifo de SymBullet, papel distinto)
    SymCursor       = "▌"  // cursor em campo de texto
    SymScrollUp     = "↑"  // indicador de scroll para cima
    SymScrollDown   = "↓"  // indicador de scroll para baixo
    SymScrollThumb  = "■"  // posição do scroll (thumb)
    SymEllipsis     = "…"  // truncamento de texto
    SymBullet       = "•"  // indicador contextual e dirty marker no cabeçalho
    SymHeaderSep    = "·"  // separador no cabeçalho
    SymTreeConnector = "<╡" // conector árvore → detalhe (2 colunas)
)

// Estrutura (Box Drawing)
const (
    SymBorderH   = "─"
    SymBorderV   = "│"
    SymCornerTL  = "╭"
    SymCornerTR  = "╮"
    SymCornerBL  = "╰"
    SymCornerBR  = "╯"
    SymJunctionL = "├"
    SymJunctionT = "┬"
    SymJunctionB = "┴"
    SymJunctionR = "┤"
)

// Spinner — iterar SpinnerFrames[frame % 4] para animar
var SpinnerFrames = [4]string{"◐", "◓", "◑", "◒"}

// SpinnerFrame retorna o frame do spinner para o índice de animação fornecido.
func SpinnerFrame(frame int) string { return SpinnerFrames[frame%4] }
```

**Nota:** `SymBullet` e `SymMask` usam o mesmo caractere `•` (U+2022), mas constantes distintas documentam os dois papéis semânticos separados.

## keys.go — Tipo Key e Teclas Pré-definidas

Arquivo novo no pacote `design`. Centraliza rótulos de UI e definições de teclas para uso tipado em `HandleKey`, barra de comandos e diálogo de ajuda.

### Tipo Key

```go
import tea "charm.land/bubbletea/v2"

// Key associa o rótulo de UI da tecla com sua representação tipada no bubbletea.
// Label é exibido na barra de comandos e no diálogo de ajuda.
// Code e Mod são usados para comparação tipada com tea.KeyMsg — sem strings mágicas.
type Key struct {
    Label string     // rótulo exibido na UI: "Enter", "⌃Q", "F1"
    Code  rune       // código da tecla: tea.KeyEnter, tea.KeyF1, 'q', etc.
    Mod   tea.KeyMod // modificadores: tea.ModCtrl, tea.ModShift, etc.
}

// Matches reporta se o evento de teclado corresponde a esta Key.
// Compara Code e Mod exatamente — sem parsing de string.
func (k Key) Matches(msg tea.KeyMsg) bool {
    key := msg.Key()
    return key.Code == k.Code && key.Mod == k.Mod
}
```

### Constantes de Modificadores (para montar labels)

```go
const (
    ModLabelCtrl  = "⌃" // U+2303
    ModLabelShift = "⇧" // U+21E7
    ModLabelAlt   = "!" // sem glifo Unicode dedicado
)
```

### Var Keys — Teclas simples pré-definidas

```go
var Keys = struct {
    Enter, Esc, Tab, Del, Ins, Home, End, PgUp, PgDn Key
    F1, F2, F3, F4, F5, F6, F7, F8, F9, F10, F11, F12 Key
}{
    Enter: Key{Label: "Enter", Code: tea.KeyEnter},
    Esc:   Key{Label: "Esc",   Code: tea.KeyEscape},
    Tab:   Key{Label: "Tab",   Code: tea.KeyTab},
    Del:   Key{Label: "Del",   Code: tea.KeyDelete},
    Ins:   Key{Label: "Ins",   Code: tea.KeyInsert},
    Home:  Key{Label: "Home",  Code: tea.KeyHome},
    End:   Key{Label: "End",   Code: tea.KeyEnd},
    PgUp:  Key{Label: "PgUp",  Code: tea.KeyPgUp},
    PgDn:  Key{Label: "PgDn",  Code: tea.KeyPgDown},
    F1:    Key{Label: "F1",    Code: tea.KeyF1},
    // ... F2–F12 seguem o mesmo padrão
}
```

### Funções de Composição

Para teclas com modificadores — combinam o Code base com o(s) Mod(s) e atualizam o Label:

```go
// Letter cria uma Key para uma tecla de letra, sem modificadores.
// O Label usa a letra maiúscula por convenção de notação.
func Letter(r rune) Key

// WithCtrl adiciona o modificador Ctrl à tecla base.
func WithCtrl(base Key) Key

// WithShift adiciona o modificador Shift à tecla base.
func WithShift(base Key) Key

// WithAlt adiciona o modificador Alt à tecla base.
func WithAlt(base Key) Key
```

Exemplos:
```go
WithCtrl(Letter('q'))                      // Key{Label: "⌃Q",    Code: 'q', Mod: ModCtrl}
WithShift(Keys.F6)                         // Key{Label: "⇧F6",   Code: KeyF6, Mod: ModShift}
WithCtrl(WithAlt(WithShift(Letter('q')))) // Key{Label: "⌃!⇧Q", Code: 'q', Mod: ModCtrl|ModAlt|ModShift}
```

### Atalhos Globais Pré-definidos

```go
var Shortcuts = struct {
    Help, ThemeToggle, Quit, LockVault Key
}{
    Help:        Keys.F1,
    ThemeToggle: Keys.F12,
    Quit:        WithCtrl(Letter('q')),
    LockVault:   WithCtrl(WithAlt(WithShift(Letter('q')))),
}
```

Uso em `HandleKey`:
```go
case tea.KeyPressMsg:
    switch {
    case design.Shortcuts.Quit.Matches(msg):
        return tea.Quit
    case design.Keys.F1.Matches(msg):
        return openHelp()
    }
```

## Estrutura Final do Pacote

```
internal/tui/design/
├── design.go    # Theme + tokens de cor + LogoGradient + MinWidth/MinHeight/PanelTreeRatio
├── symbols.go   # constantes de símbolos Unicode + SpinnerFrames + SpinnerFrame()
└── keys.go      # type Key + Matches() + ModLabel* + var Keys + funções de composição + var Shortcuts
```

## Decisões Tomadas

- `LogoGradient` é campo do `Theme` — as cores do gradiente são tema-dependentes e variam ao trocar de tema
- Símbolos como constantes nomeadas (não struct) — consistente com padrão já adotado no `old/tui/design.go`
- `SymBullet` e `SymMask` são constantes separadas apesar do mesmo valor — documentam papéis semânticos distintos
- `type Key` encapsula `Label + Code + Mod` — elimina string mágicas em `HandleKey` e provê label correto para a barra de comandos
- Funções de composição (`WithCtrl`, `WithShift`, `WithAlt`, `Letter`) em vez de campos pré-definidos para todas as combinações possíveis — evita enumeração exaustiva
- `var Shortcuts` pré-define os 4 atalhos globais do design system (F1, F12, ⌃Q, ⌃!⇧Q)
- `keys.go` importa `charm.land/bubbletea/v2` — dependência adequada dado que o pacote é especificamente para o TUI
