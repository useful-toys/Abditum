# Phase 5: TUI Scaffold + Root Model - Context

**Gathered:** 2026-03-31
**Status:** Ready for planning

<domain>
## Phase Boundary

This phase delivers the foundational TUI infrastructure: the `rootModel` (the only `tea.Model` passed to `tea.NewProgram`), the session state machine, the child model interface and stubs, the persistent frame layout (header + message bar + work area + command bar + footer), the modal stack, and the message-passing architecture. No real screen content is implemented — each zone renders a placeholder. All subsequent TUI phases fill zones with real content.

This phase implements:
- All charm.land v2 dependencies added to `go.mod`
- `sessionState` enum: `stateWelcome`, `stateOpenVault`, `stateCreateVault`, `stateVaultOpen`
- `childModel` custom interface (NOT `tea.Model`)
- `rootModel` struct: concrete child pointer fields, modal stack, vault path, terminal size, `*vault.Manager`
- Frame compositor in `rootModel.View()`: constant layout zones, placeholder content
- `modalModel` with push/pop stack mechanics
- Message dispatch rules: domain messages → all live models; keyboard/mouse → focused model only
- Model lifecycle: concrete pointer fields, nil = inactive, no cross-references
- ASCII art logo (`RenderLogo()`) ported from reference project
- Global tick: starts only after vault is opened
- Global `ctrl+Q` quit shortcut wired in `rootModel.Update`
- `main.go` wiring: parse optional vault path arg, instantiate Manager, run program

This phase does NOT implement:
- Any real screen UI (welcome form, password fields, vault tree, secret detail, template editor)
- Real timer logic beyond the tick infrastructure
- Vault open/create/lock operations (Phase 6+)
- Scroll behavior (defined per-component in later phases)

</domain>

<decisions>
## Implementation Decisions

### Dependencies

**D-01: All charm.land v2 packages added in Phase 5**
- `charm.land/bubbletea/v2` — only `rootModel` implements `tea.Model`
- `charm.land/bubbles/v2` — input components for later phases, added now
- `charm.land/lipgloss/v2` — used immediately for ASCII art and frame layout
- `github.com/charmbracelet/x/exp/teatest/v2` — golden file tests
- Pin exact latest versions: `go get charm.land/bubbletea/v2@latest` etc.
- **CRITICAL:** No v1 packages. `View()` returns `tea.View` (not `string`), key events are `tea.KeyPressMsg` (not `tea.KeyMsg`), space key is `"space"` (not `" "`).

### Session State Machine

**D-02: Four session states, no `stateLocked`**
```go
type sessionState int
const (
    stateWelcome    sessionState = iota // pick vault path, show ASCII art
    stateOpenVault                      // enter password to open existing vault
    stateCreateVault                    // enter password to create new vault
    stateVaultOpen                      // vault open — tree + detail side by side
)
```
- Lock operation = wipe sensitive memory + transition to `stateOpenVault` (re-enter password for same path). There is no `stateLocked` state.
- `rootModel` stores `vaultPath string` — populated from child via domain message, persists across lock/unlock cycles.

### Child Model Interface

**D-03: Custom `childModel` interface — does NOT implement `tea.Model`**
```go
type childModel interface {
    Update(tea.Msg) tea.Cmd   // mutates in place, returns only Cmd (no self-replacement)
    View() string              // returns string, NOT tea.View
    SetSize(w, h int)          // receives allocated size from rootModel compositor
}
```
- `View()` returns `string`. Only `rootModel.View()` returns `tea.View` (satisfying `tea.Model`).
- `Update` uses pointer receivers and mutates in place — no self-replacement return.
- Whether child interface includes `Init() tea.Cmd` is **left to researcher/planner** — needs investigation of Bubble Tea v2 initialization patterns.
- `SetSize(w, h int)` receives the child's **allocated** size (not terminal size). `rootModel` computes each child's share on `tea.WindowSizeMsg` and calls `SetSize` on all live children.

### Model Lifecycle

**D-04: Concrete pointer fields — nil = inactive/dead**
```go
type rootModel struct {
    state        sessionState
    mgr          *vault.Manager
    vaultPath    string
    width, height int

    // Child models — nil means inactive (dead, no memory held)
    welcome      *welcomeModel
    openVault    *openVaultModel
    createVault  *createVaultModel
    vaultTree    *vaultTreeModel
    secretDetail *secretDetailModel
    templateList *templateListModel
    templateDetail *templateDetailModel

    // Modal stack — LIFO, last element = topmost/active
    modals       []*modalModel
}
```
- When transitioning to a new state: allocate new child via constructor, set old child field to `nil`. Go GC reclaims the old model.
- Modals: allocated on `pushModal(...)`, popped on dismiss via `popModal()`.
- **nil-pointer safety:** store children as concrete pointer types, never as `childModel` interface. A typed nil stored in an interface is NOT nil in Go — this is a compile trap. Interface is used only transiently (e.g., in `liveModels()` helper).

**D-05: `liveModels()` helper for broadcast**
```go
func (m *rootModel) liveModels() []childModel {
    var live []childModel
    if m.welcome != nil       { live = append(live, m.welcome) }
    if m.openVault != nil     { live = append(live, m.openVault) }
    // ... all child fields
    for _, modal := range m.modals { live = append(live, modal) }
    return live
}
```

### Message Dispatch Rules

**D-06: Domain messages broadcast, input messages focused**
- **Domain messages** (vault changed, timer tick, etc.) → `rootModel` dispatches to ALL live models via `liveModels()` + all modals in stack.
- **Keyboard/mouse events** (`tea.KeyPressMsg`, `tea.MouseMsg`) → only the focused model. Focus rules:
  - If modal stack non-empty → topmost modal (`modals[len(modals)-1]`) has focus.
  - Else → currently active base child model.
- No child model ever calls another child model's methods or reads another child model's fields. All state changes flow through `vault.Manager` → domain message → `rootModel` broadcast.

### Inter-Model Communication Pattern

**D-07: All mutations via Manager + Cmd messages**
- Child calls `vault.Manager` method → on success, returns `tea.Cmd` emitting a domain message (e.g., `vaultChangedMsg{}`).
- `rootModel.Update()` receives domain message → broadcasts to all live models.
- No direct method calls between children. No shared flags or shared fields between models.
- `rootModel` passes `*vault.Manager` to each child at construction time (via constructor argument).

### Frame Layout

**D-08: Constant frame with pluggable work area**
`rootModel.View()` always renders:
```
┌─────────────────────────────────┐
│ Header                          │  ← app name, vault name, unsaved indicator
├─────────────────────────────────┤
│ Message bar                     │  ← contextual hint: what user should do now
├─────────────────────────────────┤
│                                 │
│ Work area                       │  ← changes by state (see D-09)
│                                 │
├─────────────────────────────────┤
│ Command bar                     │  ← global commands always visible
└─────────────────────────────────┘
```
- Frame zones are composed with lipgloss; header and bars have fixed heights, work area gets remaining height.
- Frame structure **may grow** in later phases — not fully locked. Planner should not over-engineer zone count.

**D-09: Work area content by state**
| State | Work area content |
|-------|------------------|
| `stateWelcome` | ASCII art logo (`RenderLogo()`) centered + placeholder stub |
| `stateOpenVault` | ASCII art logo + open-vault form stub |
| `stateCreateVault` | ASCII art logo + create-vault form stub |
| `stateVaultOpen` | `vaultTreeModel` (left) + `secretDetailModel` (right) side by side |

- Template editor (later phase): left = `templateListModel`, right = `templateDetailModel` — same split layout as vault tree/detail.
- Settings (later phase): fills work area.
- All list/tree/detail child models **must support vertical scroll** — exact mechanism left to researcher.

**D-10: Modal stack as overlay layer**
- Modals float above the frame as a layer rendered on top via `lipgloss.Place()`.
- Stack: `modals []*modalModel` — LIFO.
- Push: `modals = append(modals, newModal(...))`.
- Pop: `modals = modals[:len(modals)-1]` on dismiss (ESC or selection).
- Topmost modal receives keyboard input; all modals receive domain messages.
- Background models remain live during modal display; they do not receive keyboard/mouse input.
- On pop to empty stack, keyboard/mouse returns to active base child.

### Timers

**D-11: Timers zero until vault opens; tick starts on vault open**
- `lockTimer`, `clipboardTimer` initialized to `0` in `rootModel`.
- On vault open (`stateVaultOpen` transition), `rootModel` reads `mgr.Vault().Configuracoes()` and sets timer values.
- Global 1-second tick (`tickMsg`) is NOT started by `rootModel.Init()`. It is started as a `tea.Cmd` returned from the state transition handler when entering `stateVaultOpen`.
- Before vault is open, no tick fires — welcome/open/create screens are tick-free.

### Quit Shortcut

**D-12: `ctrl+Q` global quit**
- `rootModel.Update()` intercepts `ctrl+Q` in all states before routing to child or modal.
- Quit behavior follows `fluxos.md`: confirmation modal if unsaved changes, direct quit if no changes.
- `ctrl+C` is NOT used as quit. `q` is NOT a global quit.

### Vault Path Ownership

**D-13: `rootModel` owns vault path**
- `rootModel.vaultPath string` is the single source of truth.
- Child models (welcome, openVault) communicate the path back to `rootModel` via a `Cmd` returning a domain message (e.g., `vaultPathSelectedMsg{path: "..."}`) — never via direct field access.
- `main.go` may also provide an initial path via constructor arg (`newRootModel(mgr, initialPath)`).

### ASCII Art Logo

**D-14: `RenderLogo()` ported from reference project**
- `ascii.go` in `internal/tui/` contains `AsciiArt` constant and `RenderLogo()` function.
- 5-line wordmark, gradient coloring: one lipgloss color per line.
- Palette: `["#9d7cd8", "#89ddff", "#7aa2f7", "#7dcfff", "#bb9af7"]` (violet → blue/cyan).
- Reference implementation: `c:\git\Abditum2\internal\tui\ascii.go`.

### File Layout

**D-15: Separate files per concern in `internal/tui/`**
- `root.go` — `rootModel`, `sessionState`, `Init`/`Update`/`View`, `liveModels()`, dispatch logic
- `modal.go` — `modalModel`, push/pop helpers
- `state.go` — timer helpers, tick handler, domain message types
- `ascii.go` — `AsciiArt` constant, `RenderLogo()`
- `welcome.go` — `welcomeModel` stub
- `openvault.go` — `openVaultModel` stub
- `createvault.go` — `createVaultModel` stub
- `vaulttree.go` — `vaultTreeModel` stub
- `secretdetail.go` — `secretDetailModel` stub
- `templatelist.go` — `templateListModel` stub
- `templatedetail.go` — `templateDetailModel` stub

### Agent's Discretion
- Whether `childModel` interface includes `Init() tea.Cmd` — needs Bubble Tea v2 research
- Exact lipgloss styles for frame zones (colors, borders, padding) — Phase 5 uses minimal/placeholder styles
- Exact height allocation for header, message bar, command bar rows
- Constructor signatures for each child stub
- Exact domain message type names and fields
- `main.go` error handling details (generic fatal message format)

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Bubble Tea v2 Architecture
- `arquitetura.md` §1 — Technology choices: `charm.land/bubbletea/v2`, `charm.land/bubbles/v2`, `charm.land/lipgloss/v2`, `teatest/v2`; `View()` returns `tea.View`, key events `tea.KeyPressMsg`
- `arquitetura.md` §2 — Package structure (`internal/tui/`)
- `arquitetura.md` §3 — Manager pattern: TUI interacts with domain exclusively through Manager
- `arquitetura.md` §5 — Build conventions: CGO_ENABLED=0, no net imports, crypto/rand only
- `arquitetura.md` §6 — Security: mlock, memory wipe on lock/exit, clipboard clear, screen clear (`\033[3J\033[2J\033[H`)

### Domain Layer
- `arquitetura-dominio.md` §1 — Encapsulation: all entity fields lowercase, TUI reads via exported getters only
- `arquitetura-dominio.md` §8 — `Manager.Vault()` returns live `*Cofre` pointer; safety via package encapsulation

### Behavior Specification
- `fluxos.md` — User task flows including quit confirmation logic (VAULT-14), lock behavior (VAULT-13)
- `.planning/REQUIREMENTS.md` §VAULT-11 through §VAULT-14 — Lock, quit, clipboard, screen clear requirements

### Reference Implementation
- `c:\git\Abditum2\internal\tui\ascii.go` — `AsciiArt` constant and `RenderLogo()` to port verbatim

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/vault.Manager` — all vault operations; TUI must exclusively use this
- `internal/vault.Cofre`, `Pasta`, `Segredo` — domain entities with exported getters
- `internal/crypto.EvaluatePasswordStrength` — used in create-vault flow (Phase 6)
- `internal/storage.RecoverOrphans`, `DetectExternalChange` — called on vault open (Phase 6)
- `c:\git\Abditum2\internal\tui\ascii.go` — `RenderLogo()` reference implementation

### Established Patterns
- All sensitive data as `[]byte`, never `string`; zeroed on lock/exit
- Error messages always generic — no Go error strings, no file paths in user-facing output
- `CGO_ENABLED=0` enforced globally

### Integration Points
- `cmd/abditum/main.go` — currently `os.Exit(0)`; Phase 5 replaces with real TUI bootstrap
- `internal/vault.NewManager()` (or equivalent constructor) — instantiated in `main.go`, passed to `rootModel`
- `internal/tui` package — currently only `doc.go`; Phase 5 populates it

</code_context>

<specifics>
## Specific Ideas

- ASCII art logo: exact implementation from `c:\git\Abditum2\internal\tui\ascii.go` — 5-line wordmark with violet→cyan gradient, `RenderLogo()` returns colored string via lipgloss.
- Reference project at `c:\git\Abditum2` may contain additional TUI patterns worth reviewing during research.
- Modal stack (not single pointer) enables modals opening modals — confirmed use case in the app (e.g., file picker opening a confirmation overlay).
- Work area layout for `stateVaultOpen` and template editor is structurally identical (left list/tree + right detail panel) — the compositor logic can be shared.

</specifics>

<deferred>
## Deferred Ideas

- Exact scroll implementation for tree/list/detail — each component decides its scroll mechanics in its respective phase (7, 8, etc.)
- Real visual design of header, message bar, command bar — Phase 6 onwards when content exists
- Settings screen layout — Phase 9
- Template editor UI layout — Phase 8

</deferred>

---

*Phase: 05-tui-scaffold-root-model*
*Context gathered: 2026-03-31*
