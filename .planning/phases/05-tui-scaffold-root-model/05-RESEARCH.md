# Phase 5 Research: TUI Scaffold + Root Model

**Researched:** 2026-04-01
**Phase:** 05-tui-scaffold-root-model

---

## Standard Stack

### Dependency Versions (from reference project `c:\git\Abditum2\go.mod`)

| Package | Version | Notes |
|---------|---------|-------|
| `charm.land/bubbletea/v2` | `v2.0.2` | `View()` returns `tea.View`, not `string` |
| `charm.land/bubbles/v2` | `v2.0.0` | textinput, spinner components |
| `charm.land/lipgloss/v2` | `v2.0.2` | `lipgloss.Place()` for overlay centering |
| `github.com/atotto/clipboard` | `v0.1.4` | Clipboard clear on lock/exit |
| `github.com/charmbracelet/x/exp/teatest/v2` | `v2.0.0-20260316093931-f2fb44ab3145` | Golden file testing |

**Add all via:** `go get charm.land/bubbletea/v2@v2.0.2 charm.land/bubbles/v2@v2.0.0 charm.land/lipgloss/v2@v2.0.2 github.com/atotto/clipboard@v0.1.4 github.com/charmbracelet/x/exp/teatest/v2@v2.0.0-20260316093931-f2fb44ab3145`

---

## Bubble Tea v2 API — Critical Differences from v1

### `View()` Return Type

```go
// v1 (WRONG for this project):
func (m AppModel) View() string { ... }

// v2 (CORRECT):
func (m AppModel) View() tea.View {
    v := tea.NewView(content)  // content is string
    v.AltScreen = true         // enables alternate screen buffer
    return v
}
```

Only `rootModel` implements `tea.Model` (returns `tea.View`). All child models return plain `string` from their `View()`.

### Key Events

```go
// v1 (WRONG):
case tea.KeyMsg:
    switch msg.Type { case tea.KeyCtrlQ: ... }

// v2 (CORRECT):
case tea.KeyPressMsg:
    switch msg.String() {
    case "ctrl+q": ...
    case "space":  ...   // NOT " " — literal "space"
    case "enter":  ...
    }
```

### Tick Pattern

```go
// Start tick (deferred — NOT in Init()):
tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })

// Renew tick in Update after each tickMsg:
return m, tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
```

### WindowSizeMsg

```go
case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
    // Call SetSize on ALL live children (not just active one)
    for _, child := range m.liveModels() {
        child.SetSize(msg.Width, msg.Height) // rootModel computes allocated share
    }
```

### Program Bootstrap (v2)

```go
p := tea.NewProgram(rootModel, tea.WithAltScreen())
// tea.WithAltScreen() is separate option, NOT set on tea.View directly in main
// The reference project uses v.AltScreen = true on the View itself — both approaches work
_, err := p.Run()
```

---

## Lipgloss v2 — Overlay Mechanics

### `lipgloss.Place()` for Modal Overlay

```go
// Center a modal box over full-terminal background:
overlayBox := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("62")).
    Padding(0, 1).
    Width(boxW).
    Render(content)

// Place at center of terminal:
return lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, overlayBox)
```

The `View()` method returns the placed overlay string directly — it replaces the base frame entirely in the reference project when help/modal is active. For the Abditum architecture (modal stack over persistent frame), `rootModel.View()` renders the base frame first, then uses `lipgloss.Place()` to overlay the topmost modal.

### Width Behavior (v2 Change)

In lipgloss v2, `Width()` sets the TOTAL width including borders. Interior content width = `Width - 2` (for single border on each side).

### Frame Composition Pattern (from reference)

```go
// From vault_shell.go (adapted):
innerW := m.width - 2   // interior after frame borders
frameH := m.height - 3  // account for footer rows outside frame
panelH := frameH - 2    // meta line + separator
```

---

## Architecture Patterns (from Reference Implementation)

### Value vs Pointer Receivers

The reference project uses **value receivers** (`AppModel`) with self-replacement return `(AppModel, tea.Cmd)`. The CONTEXT.md D-04 specifies pointer receivers for child models (mutate in place, return only `tea.Cmd`) — this diverges from the reference. **Honor CONTEXT.md decision: children use pointer receivers.**

### `tea.View` Construction

```go
func (m *rootModel) View() tea.View {
    frame := m.renderFrame()      // renders full TUI as string
    v := tea.NewView(frame)
    v.AltScreen = true
    return v
}
```

### Modal Rendering Pattern

From `help_overlay.go` (reference uses full-replacement, but same `lipgloss.Place` call):
```go
func (m *rootModel) renderModalOverlay() string {
    modal := m.modals[len(m.modals)-1]
    box := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(0, 1).
        Width(40).
        Render(modal.View())
    return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}
```

### Init() Pattern in v2

The reference `AppModel.Init()` returns `nil`. Children that need initialization return their Cmds from their own `Init()` which the parent calls explicitly. Per CONTEXT.md D-11, global tick starts only when entering `workAreaVault`, NOT in `rootModel.Init()`.

```go
func (m *rootModel) Init() tea.Cmd {
    return nil  // tick starts later on vaultOpenedMsg
}
```

---

## File Layout (from CONTEXT.md D-15)

17 files to create in `internal/tui/`:

| File | Contents |
|------|----------|
| `root.go` | `rootModel`, `workArea` enum, Init/Update/View, `liveModels()` |
| `modal.go` | `modalModel`, push/pop helpers |
| `state.go` | timer helpers, tick handler, domain message types |
| `mutations.go` | Cmd factories for simple operations |
| `ascii.go` | `AsciiArt` constant, `RenderLogo()` (port verbatim) |
| `actions.go` | `ActionManager` |
| `messages.go` | `MessageManager` |
| `dialogs.go` | dialog factory functions |
| `flow_open_vault.go` | `openVaultFlow` + `openVaultDescriptor` stubs |
| `flow_create_vault.go` | `createVaultFlow` + `createVaultDescriptor` stubs |
| `flows.go` | `FlowRegistry`, `flowDescriptor`, `FlowContext`, `chainFlowMsg` |
| `prevault.go` | `preVaultModel` stub |
| `vaulttree.go` | `vaultTreeModel` stub |
| `secretdetail.go` | `secretDetailModel` stub |
| `templatelist.go` | `templateListModel` stub |
| `templatedetail.go` | `templateDetailModel` stub |
| `settings.go` | `settingsModel` stub |
| `help.go` | `helpModal` stub |

Plus: `cmd/abditum/main.go` updated.

---

## Common Pitfalls

1. **`View()` return type:** Children return `string`; only `rootModel.View()` returns `tea.View`. Returning `string` from `rootModel.View()` is a compile error.

2. **Typed nil in interface trap:** Never store children as `childModel` interface in struct fields. A typed nil (`(*preVaultModel)(nil)` stored as `childModel`) is NOT nil. Store as concrete pointer types; use interface only transiently in `liveModels()`.

3. **v2 key strings:** `"space"` not `" "`, `"ctrl+q"` not `tea.KeyCtrlQ`, `"enter"` not `tea.KeyEnter`. All key matching via `msg.String()`.

4. **Tick renewal:** Must return new `tea.Tick` cmd on every `tickMsg` to keep loop alive. Missing renewal silently stops all ticking.

5. **modalModel as childModel:** `modalModel` implements `childModel` interface (has `Update`, `View`, `SetSize`, `Context`, `ChildFlows`). Modals are included in `liveModels()` broadcast.

6. **CGO_ENABLED=0 enforcement:** teatest/v2 tests must also run with `CGO_ENABLED=0` — verify CI job-level env var is set.

7. **lipgloss.Place for overlay:** Overlay replaces the full frame string, NOT composited over it. The returned string from `lipgloss.Place()` IS the final view when a modal is active.

---

## Dependency Management Commands

```bash
# Add all Phase 5 TUI dependencies:
CGO_ENABLED=0 go get charm.land/bubbletea/v2@v2.0.2
CGO_ENABLED=0 go get charm.land/bubbles/v2@v2.0.0
CGO_ENABLED=0 go get charm.land/lipgloss/v2@v2.0.2
CGO_ENABLED=0 go get github.com/atotto/clipboard@v0.1.4
CGO_ENABLED=0 go get github.com/charmbracelet/x/exp/teatest/v2@v2.0.0-20260316093931-f2fb44ab3145

# Verify static build still works:
CGO_ENABLED=0 go build ./cmd/abditum

# Run tests with race detector:
CGO_ENABLED=0 go test -race ./internal/tui/...
```

---

## Validation Architecture

Tests for Phase 5 use `teatest/v2` golden files plus unit tests for `rootModel`:

- `TestRootModelInit` — `rootModel.Init()` returns nil (no tick before vault open)
- `TestRootModelViewType` — `rootModel.View()` return compiles as `tea.View`
- `TestTickMsg_NoFireBeforeVaultOpen` — tickMsg sent to rootModel before workAreaVault has no effect on timers
- `TestModalStack_PushPop` — push modal → topmost receives input; pop → base child receives input
- `TestLiveModels_TypedNilSafety` — nil concrete pointer fields do not appear in liveModels() result
- `TestDispatchPriority` — ctrl+Q intercepted before reaching any child; modal intercepts before base child
- `TestWindowSizeMsg` — SetSize called on all live children, not just active one

Golden file tests (after main.go wire-up works):
- `TestGolden_PreVaultPlaceholder` — `./abditum` renders placeholder at 80×24
- All golden tests use `teatest.WithInitialTermSize(80, 24)`
