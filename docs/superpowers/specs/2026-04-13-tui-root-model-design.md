# TUI Root Model Design

## Visão Geral

O TUI Abditum usa uma arquitetura onde o RootModel é o orquestrador central que coordena Views isoladas e independentes. Views não acessam estado de outras views ou do root — comunicação via métodos diretos ou eventos quando necessário.

---

## Conceitos Centrais

### Theme

Theme define cores, tipografia e símbolos. Cada tema é definido por tokens funcionais — não por valores hardcoded.

```go
type Theme struct {
    Name     string
    Surface  SurfaceTokens
    Text     TextTokens
    Accent   AccentTokens
    Border   BorderTokens
    Semantic SemanticTokens
    Special  SpecialTokens
}

var Themes = map[string]*Theme{
    "TokyoNight": TokyoNight,
    "Cyberpunk": Cyberpunk,
}

type SurfaceTokens struct {
    Base   string  // Fundo da tela inteira
    Raised string  // Painéis laterais e janelas sobrepostas
    Input  string  // Campos de texto dentro de diálogos
}

type TextTokens struct {
    Primary   string  // Texto normal
    Secondary string  // Texto de apoio, hints
    Disabled  string  // Opções inativas
    Link      string  // URLs e referências
}

type AccentTokens struct {
    Primary   string  // Barra de seleção, cursor, botão principal
    Secondary string  // Favorito ★, nomes de pastas
}

type BorderTokens struct {
    Default string  // Linhas divisórias, bordas de janelas informativas
    Focused string  // Painel ativo, campos de entrada, diálogos
}

type SemanticTokens struct {
    Success string  // Operação concluída, config ON
    Warning string  // Alerta, dirty state
    Error   string  // Erro, senha incorreta
    Info    string  // Informação contextual
    Off     string  // Config OFF
}

type SpecialTokens struct {
    Muted     string  // Texto esmaecido
    Highlight string  // Fundo do item selecionado
    Match     string  // Texto que corresponde à busca
}
```

#### Temas disponíveis

| Tema | surface.base | surface.raised | surface.input | text.primary | text.secondary | text.disabled | text.link | accent.primary | accent.secondary | border.default | border.focused | semantic.success | semantic.warning | semantic.error | semantic.info | semantic.off | special.muted | special.highlight | special.match |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| **Tokyo Night** | `#1a1b26` | `#24283b` | `#1e1f2e` | `#a9b1d6` | `#565f89` | `#3b4261` | `#7aa2f7` | `#7aa2f7` | `#bb9af7` | `#414868` | `#7aa2f7` | `#9ece6a` | `#e0af68` | `#f7768e` | `#7dcfff` | `#737aa2` | `#8690b5` | `#283457` | `#f7c67a` |
| **Cyberpunk** | `#0a0a1a` | `#1a1a2e` | `#0e0e22` | `#e0e0ff` | `#8888aa` | `#444466` | `#ff2975` | `#ff2975` | `#00fff5` | `#3a3a5c` | `#ff2975` | `#05ffa1` | `#ffe900` | `#ff3860` | `#00b4d8` | `#9999cc` | `#666688` | `#2a1533` | `#ffc107` |

```go
var TokyoNight = &Theme{
    Name: "Tokyo Night",
    Surface:  SurfaceTokens{"#1a1b26", "#24283b", "#1e1f2e"},
    Text:     TextTokens{"#a9b1d6", "#565f89", "#3b4261", "#7aa2f7"},
    Accent:   AccentTokens{"#7aa2f7", "#bb9af7"},
    Border:   BorderTokens{"#414868", "#7aa2f7"},
    Semantic: SemanticTokens{"#9ece6a", "#e0af68", "#f7768e", "#7dcfff", "#737aa2"},
    Special:  SpecialTokens{"#8690b5", "#283457", "#f7c67a"},
}

var Cyberpunk = &Theme{
    Name: "Cyberpunk",
    Surface:  SurfaceTokens{"#0a0a1a", "#1a1a2e", "#0e0e22"},
    Text:     TextTokens{"#e0e0ff", "#8888aa", "#444466", "#ff2975"},
    Accent:   AccentTokens{"#ff2975", "#00fff5"},
    Border:   BorderTokens{"#3a3a5c", "#ff2975"},
    Semantic: SemanticTokens{"#05ffa1", "#ffe900", "#ff3860", "#00b4d8", "#9999cc"},
    Special:  SpecialTokens{"#666688", "#2a1533", "#ffc107"},
}
```

#### Atributos tipográficos

| Atributo | Suporte | Fallback |
|---|---|---|
| **Bold** | Universal | — |
| Dim/Faint | Amplo | Cor comunica estado |
| *Italic* | Parcial | text.secondary |
| Underline | Amplo | — |
| ~~Strikethrough~~ | Parcial | `✗` + muted |

#### Símbolos

```go
var Symbols = struct {
    FolderCollapsed  string  // ▶
    FolderExpanded   string  // ▼
    FolderEmpty      string  // ▷
    ItemLeaf         string  // ●
    Favorite        string  // ★
    MarkedDeleted    string  // ✗
    NewlyCreated    string  // ✦
    Modified        string  // ✎
    ContextIndicator string  // •
    Revealable       string  // ◉
    Success          string  // ✓
    Info             string  // ℹ
    Warning         string  // ⚠
    Error           string  // ✕
    SeparatorV      string  // │
    SeparatorH      string  // ─
    Junction        string  // ├ ┬ ┴ ┤
    Corners         string  // ╭ ╮ ╰ ╯
    Truncation      string  // …
    Masked          string  // ••••
    ScrollThumb     string  // ■
    ScrollArrow    string  // ↑ ↓
}
```

#### Dimensões mínimos

```go
const (
    MinWidth  = 80  // Colunas - padrão POSIX
    MinHeight = 24  // Linhas - padrão POSIX
)

### Frame (Tela Principal)

A tela principal é dividida em:
- **Header** (2 linhas fixas)
- **Work Area** (área de trabalho)
- **Message Bar** (1 linha)
- **Action Bar** (1 linha)

### Work Area States

| WorkArea | childViews |
|---------|-----------|
| WorkAreaWelcome | Welcome |
| WorkAreaSettings | Settings |
| WorkAreaVault | Tree + Detail (split) |
| WorkAreaTemplates | List + Detail (split) |

---

## Interfaces

### ChildView

- Componente da tela principal
- Coexiste em paralelo com outras ChildViews na mesma Work Area
- Espaço fixo dedicado pelo root (não negociado)
- Referenciado diretamente via ponteiro

```go
type ChildView interface {
    Render(height, width int, theme Theme) string
    HandleKey(msg tea.KeyMsg) tea.Cmd
    HandleEvent(event any)
    HandleTeaMsg(msg tea.Msg)
}
```

### ModalView

- Stack sobre a tela principal
- Só interage com quem a apresentou (parent)
- Espaço máximo, centralizado pelo root
- Parametrizado via construtor
- Retorno via Cmd que produz ModalSubmit

```go
type ModalView interface {
    Render(maxHeight, maxWidth int, theme Theme) string
    HandleKey(msg tea.KeyMsg) tea.Cmd
}
```

### Modal Stack

```go
type RootModel struct {
    // ...
    modals []ModalView
}

// PushModal adiciona modal à stack
func (r *RootModel) PushModal(modal ModalView)

// PopModal remove modal do topo
func (r *RootModel) PopModal()

// TopModal retorna o modal do topo (se houver)
func (r *RootModel) TopModal() ModalView
```

### Mensagens de Modais

```go
// Empilha um novo modal no topo da pilha
type OpenModalMsg struct {
    Modal ModalView
}

// Desempilha o modal do topo da pilha
type CloseModalMsg struct{}

// Modal emitiu quando tem resultado disponível
// Não implica fechamento — quem decide fechar é o chamador
type ModalReadyMsg struct{}
```

### Funções auxiliares

```go
import tea "charm.land/bubbletea/v2"
import "charm.land/lipgloss/v2"

func OpenModal(modal ModalView) tea.Cmd {
    return func() tea.Msg { return OpenModalMsg{Modal: modal} }
}

func CloseModal() tea.Cmd {
    return func() tea.Msg { return CloseModalMsg{} }
}

func CloseApp() tea.Cmd {
    return tea.Quit
}

### Parametrização de Modal

Modais são parametrizados via construtor:

```go
// passwordModal
func NewPasswordModal(title string, id int) ModalView

// filePickerModal
func NewFilePickerModal(initialPath string, suggestions []string) ModalView

// confirmModal
func NewConfirmModal(title, message string) ModalView
```

---

## RootModel

### Estado

```go
type RootModel struct {
    // Dimensões
    width  int
    height int

    // Estado da aplicação
    workArea     WorkArea // qual WorkArea está ativa
    focusedChild ChildView // ponteiro para ChildView com foco (input)

    // Tema
    theme *Theme
}

// ToggleTheme alterna entre TokyoNight e Cyberpunk
func (r *RootModel) ToggleTheme()

    // Vault
    vaultManager *vault.Manager

    // Views (tipos dos subpackages)
    welcome      *welcome.Welcome    // de tui/welcome
    settings    *settings.Settings // de tui/settings
    vaultTree    *secret.Tree      // de tui/secret
    secretDetail *secret.Detail    // de tui/secret
    templateList *template.List     // de tui/template
    templateDetail *template.Detail // de tui/template

    // Modals stack
    modals []ModalView

    // Timers
    lastActionAt time.Time
}
```

### Responsabilidades

1. **Dimensões**
   - Armazena width/height do terminal
   - Se dimensões desconhecidas, exibe "Aguarde..."

2. **Coordenação de Render**
   ```
   Render():
     se width == 0 || height == 0:
         → "Aguarde..."
     senão:
         → render all ChildViews → render modalStack
   ```

3. **Render (View)**
```go
func (r *RootModel) View() string {
    // Se dimensões desconhecidas → mostra mensagem
    if r.width == 0 || r.height == 0 {
        return "Aguarde..."
    }

    base := r.activeView.Render(r.width, r.height, r.theme)

    if len(r.modals) == 0 {
        return base
    }

    top := r.modals[len(r.modals)-1]
    modalView := top.Render(r.width, r.height, r.theme)

    // overlay() centraliza modal sobre a view base usando lipgloss.Place
    return overlay(base, modalView, r.width, r.height)
}

// overlay posiciona modal centralizado sobre a view base
func overlay(base, modal string, width, height int) string {
    // dimensions da work area (subtraindo header, msg bar, action bar)
    workH := height - 4  // 2 header + 1 msg + 1 action
    return lipgloss.Place(width, workH, lipgloss.Center, lipgloss.Center, modal)
}
```

4. **Dispatch de Teclas**
   ```
   HandleKey(tea.KeyMsg):
     se modalStack não vazia:
         → modal no topo .HandleKey()
     senão se focusedChild != nil:
         → focusedChild.HandleKey()
     senão:
         → root trata
   ```

5. **Update - Roteamento de Mensagens**
```go
func (r *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // WindowSizeMsg → root
    if size, ok := msg.(tea.WindowSizeMsg); ok {
        r.width = size.Width
        r.height = size.Height
        return r, nil
    }

    // Atalhos globais (ctrl+q, f12, etc) → root
    // ...

    // Mensagens de ciclo de vida dos modais
    switch msg := msg.(type) {
    case OpenModalMsg:
        r.modals = append(r.modals, msg.Modal)
        return r, nil

    case CloseModalMsg:
        if len(r.modals) > 0 {
            r.modals = r.modals[:len(r.modals)-1]
        }
        return r, nil

    case ModalReadyMsg:
        // Roteia para view (ou modal pai se aninhado)
        if len(r.modals) > 1 {
            parent := r.modals[len(r.modals)-2]
            return r, parent.Update(msg)
        }
        return r, r.activeView.Update(msg)
    }

    // Se há modal ativo → modal recebe mensagens
    if len(r.modals) > 0 {
        top := len(r.modals) - 1
        return r, r.modals[top].Update(msg)
    }

    // Sem modal → view ativa recebe mensagens
    return r, r.activeView.Update(msg)
}
```

4. **Event Routing**
   | Mensagem | Destino |
   |---------|--------|
   | WindowSizeMsg | root |
   | domain events (secretAdded, etc) | todas as ChildViews |
   | tickMsg | todas as ChildViews |
   | ModalReadyMsg | view (ou modal pai se aninhado) |

5. **Modal Submit (retorno)**
   - Modal retorna Cmd que produz `ModalReadyMsg`
   - Root recebe e roteia para view (ou modal pai se aninhado)
   - View consulta getters do modal diretamente (não via mensagem)
   - View pode rejeitar resultado e manter modal ativo
   - Fechamento explícito via `CloseModal()` Cmd

6. **Ações (futuro)**
   - ActionManager armazenará actions globais, reusáveis, específicas
   - Verifica `enabled()` antes de executar
   - Por enquanto: não implementar

---

## Estrutura de Arquivos

```
tui/
├── view.go              # interfaces ChildView, ModalView + Theme + WorkArea
├── root.go             # RootModel
├── welcome/
│   └── welcome_view.go   # type Welcome
├── settings/
│   └── settings_view.go # type Settings
├── secret/
│   ├── vault_tree.go    # type Tree - árvore de segredos
│   └── secret_detail.go # type Detail - detalhe de segredo
├── template/
│   ├── template_list.go   # type List - lista de modelos
│   └── template_detail.go # type Detail - detalhe de modelo
└── modal/
    ├── modal_base.go      # Intent + ModalOption + baseModal
    ├── password_modal.go  # passwordModal
    ├── confirm_modal.go   # confirmModal
    ├── filepicker_modal.go # filepickerModal
    └── help_modal.go      # helpModal
```

---

## Keyboard Flow

```
tea.KeyPressMsg
    ↓
Root.HandleKey()
    ├── modalStack não vazia?
    │   └── yes: modal.top.HandleKey() → executa cmd
    ├── focusedChild != nil?
    │   └── yes: child.HandleKey() → executa cmd
    └── root trata
```

---

## Exemplos de WorkArea

### WorkAreaWelcome

```
┌────────────────────────────────┐
│ Header (2 linhas)               │
├────────────────────────────────┤
│                                │
│ welcomeView.Render(h, w, theme)  │
│                                │
├────────────────────────────────┤
│ msg bar (1 linha)              │
├────────────────────────────────┤
│ action bar (1 linha)            │
└────────────────────────────────┘
```

### WorkAreaVault

```
┌────────────────────────────────┐
│ Header (2 linhas)               │
├──────────────┬─────────────────┤
│ vaultTree  │ secretDetail    │
│ Render    │ Render        │
│ (w/2, h) │ (w-w/2, h)   │
├──────────────┴─────────────────┤
│ msg bar                      │
├─────────────────────────── │
│ action bar                  │
└─────────────────────────── │
```

---

## Interação entre ChildView e Modal

### View cria e guarda referência

```go
type ViewState int

const (
    StateNormal ViewState = iota
    StateAwaitingModal
)

type WelcomeView struct {
    state         ViewState
    passwordModal *PasswordModal // nil quando inativo
}

// Em algum ponto do HandleKey da view:
func (v *WelcomeView) HandleKey(msg tea.KeyMsg) tea.Cmd {
    switch msg.String() {
    case "enter":
        v.passwordModal = NewPasswordModal("Digite sua senha:")
        v.state = StateAwaitingModal
        return OpenModal(v.passwordModal)
    }
    return nil
}
```

### HandleKey - exemplo completo

```go
func (v *WelcomeView) HandleKey(msg tea.KeyMsg) tea.Cmd {
    // 1. trata ModalReadyMsg quando aguardar modal
    if ready, ok := msg.(ModalReadyMsg); ok && v.state == StateAwaitingModal {
        senha := v.passwordModal.Password()
        if len(senha) < 8 {
            v.passwordModal.SetError("Senha muito curta")
            return nil // modal permanece ativo
        }
        v.password = senha
        v.passwordModal = nil
        v.state = StateNormal
        return CloseModal()
    }

    // 2. enquanto aguarda modal, ignora outras mensagens
    if v.state == StateAwaitingModal {
        return nil
    }

    // 3. lógica normal da view
    switch msg.String() {
    case "enter":
        // abrir modal...
    case "q":
        return CloseApp() // retorna Cmd para sair
    }
    return nil
}
```

### Render - preenche espaço dedicado

```go
func (v *WelcomeView) Render(width, height int, theme Theme) string {
    // View DEVE preencher exatamente width x height
    content := v.renderContent(theme)
    style := lipgloss.NewStyle().Width(width).Height(height)
    return style.Render(content)
}

func (v *WelcomeView) renderContent(theme Theme) string {
    // conteúdo visual da view
    return theme.Text.Primary + "Welcome to Abditum"
}
```

---

## Modais Aninhados

A pilha suporta naturalmente modais sobre modais. Um modal pode abrir outro emitindo `OpenModalMsg`. O root empilha o novo modal. Quando o modal filho emite `ModalReadyMsg`, o root o entrega ao modal pai (elemento imediatamente abaixo na pilha).

---

## Separação de Responsabilidades

| Camada | Responsabilidade |
|--------|-----------------|
| **RootModel** | Armazenar width, height e theme; tratar tea.WindowSizeMsg; interceptar OpenModalMsg, CloseModalMsg, ModalReadyMsg; rotear mensagens para modal ou view; renderizar overlay |
| **ChildView** | Criar modais com configuração correta; guardar referência tipada; tratar ModalReadyMsg; validar dados via getters; injetar feedback via setters; emitir CloseModal quando satisfeita |
| **ModalView** | Processar input do usuário; manter estado interno; emitir ModalReadyMsg quando há resultado; expor getters para dados e setters para feedback; pode abrir modais filhos via OpenModal |

---

## Pendentes

- [ ] ActionManager (não implementar agora)
- [ ] MessageManager (não implementar agora)
- [ ] Comunicação entre ChildViews (listeners)
- [ ] Event naming convention

---

## Decisões Tomadas

- ChildViews referenciadas diretamente via ponteiro (não usa ID)
- ModalView não guarda estado de root — comunicação via Cmd/msg
- Modal retorna Cmd opcional, root executa transparentemente
- ModalOption.Action() executado pelo HandleKey do modal para decidir cmd
- Subpackages organizados por domínio (welcome, settings, secret, template, modal)
- ActionManager e MessageManager: não implementar nesta fase