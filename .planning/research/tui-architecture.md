# TUI Architecture Research: Abditum

**Researched:** 2026-03-24  
**Focus:** Bubble Tea v2 architecture patterns, multi-panel layout, password UX, clipboard, security features  
**Confidence:** HIGH (Context7 + official sources + verified examples)

---

## 1. Library Versions & Import Paths (CRITICAL — v2 Breaking Change)

All Charmbracelet libraries have migrated to v2 with **new import paths**. Use these canonical paths:

| Library | Import Path | Purpose |
|---------|------------|---------|
| Bubble Tea | `charm.land/bubbletea/v2` | TUI framework |
| Lip Gloss | `charm.land/lipgloss/v2` | Styling & layout |
| Bubbles | `charm.land/bubbles/v2` | UI components |
| Huh | `charm.land/huh/v2` | Form dialogs |

> **Do NOT use** `github.com/charmbracelet/...` paths — those are v0/v1 and the API is different.

---

## 2. Bubble Tea v2: Key API Changes from v1

### `View()` return type changed

```go
// v1 (OLD — do not use)
func (m Model) View() string { ... }

// v2 (CORRECT)
func (m Model) View() tea.View { ... }

// Return a view:
return tea.NewView(m.renderContent())
```

### Keyboard messages changed

```go
// v1 (OLD)
case tea.KeyMsg:
    switch msg.Type { ... }
    switch msg.String() { ... }

// v2 (CORRECT)
case tea.KeyPressMsg:
    switch msg.Code {
    case tea.KeyEnter:  // named key
    }
    switch msg.Text {  // printed text ("a", "B", "!", etc.)
    case "q":
    }
    // Modifiers: msg.Mod (tea.ModCtrl, tea.ModAlt, tea.ModShift)
    if msg.Mod.Contains(tea.ModCtrl) { ... }
```

### Space bar key

```go
// v1: case " ":
// v2: case "space": — or use msg.Code == tea.KeySpace
```

### Alt screen, mouse, bracketed paste → declarative View fields

```go
// v1 (OLD): tea.WithAltScreen() program option, tea.EnterAltScreen command
// v2 (CORRECT): set in View() return value

func (m Model) View() tea.View {
    v := tea.NewView(m.render())
    v.AltScreen = true           // enable alternate screen
    v.MouseAllMotion = true      // enable full mouse tracking
    v.BracketedPaste = true      // enable bracketed paste
    return v
}
```

### Program creation

```go
// v2
p := tea.NewProgram(initialModel)
finalModel, err := p.Run()
```

---

## 3. Multi-Panel Layout Architecture

### Recommended: Composable Views with `sessionState` enum

Track which panel is active using an integer enum. Delegate `Update()` to the active sub-model.

```go
type sessionState uint

const (
    sidebarView sessionState = iota
    detailView
    modalView
)

type Model struct {
    state   sessionState
    sidebar SidebarModel
    detail  DetailModel
    modal   ModalModel
    width   int
    height  int
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        // ... distribute to sub-models

    case tea.KeyPressMsg:
        switch {
        case msg.Code == tea.KeyTab:
            // cycle panels
            if m.state == sidebarView {
                m.state = detailView
            } else {
                m.state = sidebarView
            }
        case msg.Code == tea.KeyEscape && m.state == modalView:
            m.state = sidebarView
        default:
            // delegate to active panel
            return m.delegateUpdate(msg)
        }
    }
    return m, nil
}

func (m Model) delegateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    switch m.state {
    case sidebarView:
        m.sidebar, cmd = m.sidebar.Update(msg)
    case detailView:
        m.detail, cmd = m.detail.Update(msg)
    case modalView:
        m.modal, cmd = m.modal.Update(msg)
    }
    return m, cmd
}
```

Source: Official `composable-views` example in BubbleTea repo.

### Split Panel Rendering with Lip Gloss

```go
func (m Model) render() string {
    sidebarWidth := 30
    detailWidth := m.width - sidebarWidth - 1  // 1 for separator

    sidebarStyle := lipgloss.NewStyle().
        Width(sidebarWidth).
        Height(m.height - 2)  // 2 for status bar + help bar

    detailStyle := lipgloss.NewStyle().
        Width(detailWidth).
        Height(m.height - 2)

    // Apply focused/blurred border
    if m.state == sidebarView {
        sidebarStyle = sidebarStyle.BorderStyle(lipgloss.RoundedBorder()).
            BorderForeground(lipgloss.Color("#FF6E6E"))  // cyberpunk accent
    } else {
        sidebarStyle = sidebarStyle.BorderStyle(lipgloss.RoundedBorder()).
            BorderForeground(lipgloss.Color("#444444"))  // dim when inactive
    }

    // Join horizontally
    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        sidebarStyle.Render(m.sidebar.View()),
        detailStyle.Render(m.detail.View()),
    )
}
```

Source: Official `split-editors` example in BubbleTea repo.

### Responsive Layout: Width Breakpoints

Per `descricao.md` requirements:

```go
func (m Model) View() tea.View {
    if m.height < 5 || m.width < 20 {
        msg := "Terminal muito pequeno.\nRedimensione para continuar."
        return tea.NewView(lipgloss.Place(m.width, m.height,
            lipgloss.Center, lipgloss.Center, msg))
    }
    if m.width < 40 {
        // Detail panel only
        return tea.NewView(m.renderDetailOnly())
    }
    // Full layout
    return tea.NewView(m.renderFull())
}
```

---

## 4. Modal / Overlay Pattern with Huh

`huh.Form` implements `tea.Model` — embed it directly. Route all input to it when modal is active.

```go
import "charm.land/huh/v2"

type Model struct {
    state     sessionState
    huhForm   *huh.Form
    // ...
}

// Open a modal
func (m *Model) openNewSecretModal() {
    m.huhForm = huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Title("Nome do Segredo").
                Placeholder("ex: GitHub").
                Value(&m.newSecretName),
            huh.NewSelect[string]().
                Title("Modelo").
                Options(huh.NewOptions("Login", "Cartão de Crédito", "API Key", "Personalizado")...).
                Value(&m.newSecretTemplate),
        ),
    ).WithTheme(huh.ThemeCharm())
    m.state = modalView
}

// In Update:
case modalView:
    form, cmd := m.huhForm.Update(msg)
    if f, ok := form.(*huh.Form); ok {
        m.huhForm = f
    }
    if m.huhForm.State == huh.StateCompleted {
        // form was submitted — handle result
        m.handleNewSecret()
        m.state = sidebarView
    } else if m.huhForm.State == huh.StateAborted {
        m.state = sidebarView
    }
    return m, cmd
```

**Huh's built-in themes:** `ThemeCharm()`, `ThemeDracula()`, `ThemeCatppuccin()`, `ThemeBase16()`, `ThemeDefault()`

**Password fields in Huh:**
```go
huh.NewInput().
    Title("Senha Mestra").
    EchoMode(huh.EchoPassword)  // shows • characters
```

---

## 5. Tree Rendering (Sidebar Hierarchy)

Lip Gloss v2 includes a `tree` sub-package for rendering tree structures:

```go
import "charm.land/lipgloss/v2/tree"

t := tree.Root("Cofre").
    Child(
        tree.New().Root("📁 Sites").
            Child("🔑 GitHub").
            Child("🔑 Gmail"),
        tree.New().Root("📁 Financeiro").
            Child("💳 Nubank"),
    )

rendered := t.String()
```

**Important caveat:** `lipgloss/tree` is a **rendering** library, not interactive. It produces a static string. For interactive navigation with keyboard support, you need a custom implementation.

### Recommended Approach: Custom Interactive Tree

Build on top of Lip Gloss rendering, maintaining your own cursor state:

```go
type TreeNode struct {
    ID       string
    Name     string
    Type     NodeType  // folder | secret
    Children []TreeNode
    Expanded bool
    Modified bool
    Favorite bool
}

type SidebarModel struct {
    root     []TreeNode
    cursor   int          // index into flattened visible list
    flat     []FlatNode   // pre-computed flat view for rendering
}

// Flatten the tree for linear navigation
func (m *SidebarModel) flatten() {
    m.flat = flattenNodes(m.root, 0)
}

// On Up/Down: move cursor, re-render
// On Right: expand folder or move to first child
// On Left: collapse folder or jump to parent
```

This gives full control over indicators (favorites, modified, etc.) as required.

---

## 6. Input Masking for Password Fields

Use `bubbles/textinput` directly for custom forms, or `huh.Input.EchoMode()` for huh forms.

### With Bubbles textinput

```go
import "charm.land/bubbles/v2/textinput"

masterPasswordInput := textinput.New()
masterPasswordInput.Placeholder = "Senha Mestra"
masterPasswordInput.EchoMode = textinput.EchoPassword   // shows •
// Or: textinput.EchoNone (shows nothing at all — maximum privacy)

// Custom echo character:
masterPasswordInput.EchoCharacter = '●'
```

### Echo modes

| Mode | Display | Use Case |
|------|---------|---------|
| `EchoNormal` | plain text | regular fields |
| `EchoPassword` | `••••••••` | password entry (shows length) |
| `EchoNone` | ` ` (nothing) | maximum privacy, hides even length |

**Recommendation for Abditum:** Use `EchoPassword` for master password input (user needs feedback that they're typing). Use `EchoNone` as an option for extreme privacy scenarios.

---

## 7. Shoulder-Surfing Protection: Screen Blanking

Based on the official `vanish` example:

```go
type Model struct {
    hidden bool
    // ...
}

func (m Model) View() tea.View {
    if m.hidden {
        return tea.NewView("")  // completely blank screen
    }
    return tea.NewView(m.render())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyPressMsg:
        // Toggle hide (e.g., Ctrl+H)
        if msg.Mod.Contains(tea.ModCtrl) && msg.Text == "h" {
            m.hidden = !m.hidden
            return m, nil
        }
        // Unhide on any keypress if hidden
        if m.hidden {
            m.hidden = false
            return m, nil
        }
    }
    // ... normal handling
}
```

Combined with `v.AltScreen = true`, the alternate terminal buffer is cleared completely — nothing leaks to scrollback history.

---

## 8. Clipboard Operations

### Strategy: Try OSC52 (BubbleTea native), fall back to atotto/clipboard

**OSC52 via BubbleTea v2** (works in modern terminals, SSH, tmux):
```go
// Write to clipboard
cmd := tea.SetClipboard(fieldValue)
return m, cmd

// Read clipboard (returns tea.ClipboardMsg)
cmd := tea.ReadClipboard()
return m, cmd

// Handle read result:
case tea.ClipboardMsg:
    clipboardContent := string(msg)
```

**Fallback: `github.com/atotto/clipboard`** (cross-platform, requires xclip/xsel on Linux):
```go
import "github.com/atotto/clipboard"

err := clipboard.WriteAll(fieldValue)
content, err := clipboard.ReadAll()
```

### Recommended clipboard service for Abditum

```go
type ClipboardService struct {
    clearTimer *time.Timer
    clearDelay time.Duration  // configurable, default 30s
}

func (cs *ClipboardService) Copy(value string) tea.Cmd {
    return func() tea.Msg {
        // Try OSC52 first (handled by tea.SetClipboard)
        // Reset auto-clear timer
        if cs.clearTimer != nil {
            cs.clearTimer.Stop()
        }
        cs.clearTimer = time.AfterFunc(cs.clearDelay, cs.clear)
        return clipboardCopiedMsg{value: value}
    }
}

func (cs *ClipboardService) ClearOnLock() {
    if cs.clearTimer != nil {
        cs.clearTimer.Stop()
    }
    cs.clear()
}
```

**Per requirements:** Clear clipboard on lock, on close, and after 30s automatically.

---

## 9. Lip Gloss Styling Best Practices for Portability

### Color handling

Lip Gloss automatically downsamples colors to the terminal's capability (256-color, 16-color, no-color). Use hex colors freely:

```go
lipgloss.Color("#FF6E6E")   // automatically adapted
lipgloss.Color("#7B61FF")
```

### Adaptive theming

```go
var accentColor lipgloss.TerminalColor
if lipgloss.HasDarkBackground() {
    accentColor = lipgloss.Color("#FF6E6E")
} else {
    accentColor = lipgloss.Color("#CC0000")
}
```

### Centering overlays (for modals)

```go
overlay := lipgloss.Place(
    m.width, m.height,
    lipgloss.Center, lipgloss.Center,
    modalStyle.Render(content),
)
```

### True overlay via compositor (advanced)

```go
// Lip Gloss v2 has a compositor for Z-layered overlays
import "charm.land/lipgloss/v2"

// Layer a modal over existing content
baseLayer := lipgloss.NewLayer(baseContent)
modalLayer := lipgloss.NewLayer(modalContent).
    X(centerX).Y(centerY).Z(1)

result := lipgloss.Compose(baseLayer, modalLayer)
```

### Border styles

```go
lipgloss.NormalBorder()   // single line: ─ │ ┌ etc.
lipgloss.RoundedBorder()  // rounded: ╭ ╮ ╰ ╯
lipgloss.ThickBorder()    // double-width
lipgloss.DoubleBorder()   // ═ ║ ╔ etc.
lipgloss.HiddenBorder()   // invisible (spacing only)
```

---

## 10. Keyboard Bindings with `bubbles/key`

Use the `key` package for structured, self-documenting keymaps:

```go
import "charm.land/bubbles/v2/key"

type KeyMap struct {
    Up       key.Binding
    Down     key.Binding
    Expand   key.Binding
    Collapse key.Binding
    NewItem  key.Binding
    Copy     key.Binding
    Lock     key.Binding
    Hide     key.Binding
    Quit     key.Binding
}

var DefaultKeyMap = KeyMap{
    Up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "subir")),
    Down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "descer")),
    Expand:   key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→", "expandir")),
    Collapse: key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←", "colapsar")),
    NewItem:  key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("ctrl+n", "novo segredo")),
    Copy:     key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "copiar campo")),
    Lock:     key.NewBinding(key.WithKeys("ctrl+l"), key.WithHelp("ctrl+l", "bloquear")),
    Hide:     key.NewBinding(key.WithKeys("ctrl+h"), key.WithHelp("ctrl+h", "ocultar tela")),
    Quit:     key.NewBinding(key.WithKeys("ctrl+q"), key.WithHelp("ctrl+q", "sair")),
}

// In Update:
case tea.KeyPressMsg:
    switch {
    case key.Matches(msg, m.keys.Up):
        m.moveCursorUp()
    case key.Matches(msg, m.keys.Down):
        m.moveCursorDown()
    }
```

The `key.Binding` struct integrates with `bubbles/help` to auto-generate the context-sensitive help bar.

---

## 11. Help Bar with `bubbles/help`

```go
import "charm.land/bubbles/v2/help"

type Model struct {
    help help.Model
    keys KeyMap
}

func (m Model) View() tea.View {
    helpView := m.help.View(m.keys)  // KeyMap must implement help.KeyMap interface
    // ... render at bottom of screen
}

// KeyMap interface requires:
// ShortHelp() []key.Binding
// FullHelp() [][]key.Binding
```

---

## 12. Auto-Lock via Inactivity Timer

```go
type tickMsg time.Time

func tickCmd() tea.Cmd {
    return tea.Tick(time.Second, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

type Model struct {
    lastActivity time.Time
    lockTimeout  time.Duration
    locked       bool
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tickMsg:
        if !m.locked && time.Since(m.lastActivity) >= m.lockTimeout {
            m = m.lock()
        }
        return m, tickCmd()  // keep ticking

    case tea.KeyPressMsg, tea.MouseMsg:
        m.lastActivity = time.Now()
        // ... handle normally
    }
}
```

---

## 13. Spinner for Long Operations

```go
import "charm.land/bubbles/v2/spinner"

type Model struct {
    spinner  spinner.Model
    loading  bool
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if m.loading {
        var cmd tea.Cmd
        m.spinner, cmd = m.spinner.Update(msg)
        return m, cmd
    }
    // ...
}

func (m Model) render() string {
    if m.loading {
        return m.spinner.View() + " Salvando cofre..."
    }
    // ...
}
```

---

## 14. File Picker for Open/Save Dialogs

`bubbles/filepicker` provides TUI-native file navigation:

```go
import "charm.land/bubbles/v2/filepicker"

fp := filepicker.New()
fp.AllowedTypes = []string{".abditum"}
fp.CurrentDirectory, _ = os.UserHomeDir()
```

This satisfies the requirement for a TUI file picker instead of blind path typing.

---

## 15. Recommended Project Structure (DDD-aligned)

```
abditum/
├── main.go
├── domain/
│   ├── vault/
│   │   ├── vault.go          # Vault entity (read-only access)
│   │   ├── manager.go        # VaultManager (all mutations)
│   │   ├── secret.go         # Secret entity
│   │   ├── folder.go         # Folder entity
│   │   └── template.go       # SecretTemplate entity
│   └── crypto/
│       ├── crypto.go         # AES-256-GCM + Argon2id
│       └── crypto_test.go
├── storage/
│   ├── storage.go            # Atomic save, .bak, .tmp pattern
│   └── storage_test.go
├── tui/
│   ├── app.go                # Root model (WindowSizeMsg, global state)
│   ├── sidebar/
│   │   ├── model.go          # SidebarModel
│   │   └── tree.go           # Tree flattening and rendering
│   ├── detail/
│   │   └── model.go          # DetailModel
│   ├── modal/
│   │   ├── new_secret.go     # New/Edit secret modal (huh.Form)
│   │   ├── confirm.go        # Confirmation dialog
│   │   └── filepicker.go     # File picker wrapper
│   ├── auth/
│   │   └── model.go          # Unlock screen model
│   ├── statusbar/
│   │   └── model.go          # Status bar (path, modified, count)
│   ├── helpbar/
│   │   └── model.go          # Context-sensitive help bar
│   └── keys/
│       └── keys.go           # All key bindings
├── services/
│   ├── clipboard.go          # Copy + auto-clear service
│   └── autolock.go           # Inactivity timer service
└── testdata/
    └── golden/               # Visual golden files (80×24)
```

---

## 16. Testing Strategy

### Golden file tests with `teatest`

```go
import "charm.land/bubbletea/v2/teatest"

func TestUnlockScreen(t *testing.T) {
    tm := teatest.NewTestModel(t, initialModel(),
        teatest.WithInitialTermSize(80, 24))

    // Wait for stable render
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return strings.Contains(string(bts), "Abditum")
    })

    // Capture golden file
    tm.FinalOutput(t, teatest.WithGoldenFile("testdata/golden/unlock.txt"))
}
```

### Key event simulation

```go
tm.Send(tea.KeyPressMsg{Code: tea.KeyDown})
tm.Send(tea.KeyPressMsg{Text: "q", Mod: tea.ModCtrl})
```

---

## 17. Key Decisions Summary

| Decision | Choice | Rationale |
|----------|--------|-----------|
| TUI framework | Bubble Tea v2 | Official Charmbracelet, active maintenance, v2 is current |
| Layout | Composable views + `sessionState` enum | Official pattern, clean delegation |
| Forms/modals | `huh.Form` embedded as `tea.Model` | Native integration, no custom form logic |
| Tree sidebar | Custom flat-list cursor on top of lipgloss/tree rendering | `lipgloss/tree` is render-only; need custom navigation |
| Password echo | `textinput.EchoPassword` | Shows • (user gets length feedback) |
| Clipboard | `tea.SetClipboard` (OSC52) + `atotto/clipboard` fallback | Maximum compatibility |
| Screen blank | `tea.NewView("")` when `model.hidden == true` | Official `vanish` example pattern |
| File picker | `bubbles/filepicker` | TUI-native, no external dependencies |
| Key bindings | `bubbles/key.Binding` + `bubbles/help` | Auto-generates help bar |
| Auto-lock | `tea.Tick` + `time.Since(lastActivity)` | Simple, no goroutine leaks |
| Testing | `teatest` + golden files at 80×24 | Official testing library |

---

## Sources

- BubbleTea v2 upgrade guide: `charm.land/bubbletea/v2` (Context7, verified)
- BubbleTea examples: `split-editors`, `composable-views`, `vanish`, `clipboard` (official GitHub)
- Lip Gloss v2 README + tree sub-package (Context7, verified)
- Huh v2 README + form embedding docs (Context7, verified)
- Bubbles v2: `textinput`, `filepicker`, `key`, `help`, `spinner` (Context7, verified)
- `github.com/atotto/clipboard` README (WebFetch, official)
