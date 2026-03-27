# Architecture Patterns — Abditum

**Domain:** Go TUI password manager (offline, single-binary)
**Researched:** 2026-03-27
**Confidence:** HIGH — Bubble Tea v2 verified against official docs and source; patterns from official examples

---

## System Overview

```
┌──────────────────────────────────────────────────────────────────┐
│                         cmd/abditum/                             │
│  main.go — wires Manager + TUI, calls tea.NewProgram().Run()     │
└────────────────────────┬─────────────────────────────────────────┘
                         │ passes *vault.Manager to root model
                         ▼
┌──────────────────────────────────────────────────────────────────┐
│                       internal/tui/                              │
│                                                                  │
│  rootModel ─── sessionState enum ─── active screen/overlay      │
│  ┌──────────────┐ ┌─────────────────┐ ┌───────────────────────┐ │
│  │ welcomeModel │ │   vaultModel    │ │ secretDetailModel     │ │
│  │ (no vault)   │ │ (tree + search) │ │ (view/create/edit)    │ │
│  └──────────────┘ └─────────────────┘ └───────────────────────┘ │
│  Overlays (composed on top of active screen):                    │
│  ┌──────────────┐ ┌─────────────────┐ ┌───────────────────────┐ │
│  │  modalModel  │ │ filePickerModel │ │     helpModel         │ │
│  │ (warn/confirm│ │ (open/save/exp) │ │  (key bindings)       │ │
│  └──────────────┘ └─────────────────┘ └───────────────────────┘ │
│  Timers (in rootModel):                                          │
│  ┌──────────────┐ ┌─────────────────┐ ┌───────────────────────┐ │
│  │ lockTimer    │ │ clipboardTimer  │ │ revealedFields map    │ │
│  │ (auto-lock)  │ │ (clear after N) │ │ + globalTick (1s)     │ │
│  └──────────────┘ └─────────────────┘ └───────────────────────┘ │
└────────────────────────┬─────────────────────────────────────────┘
                         │ reads via Manager.Vault() (snapshot)
                         │ mutates via Manager methods ONLY
                         ▼
┌──────────────────────────────────────────────────────────────────┐
│                      internal/vault/                             │
│  Manager — single entry point for ALL mutations                  │
│  Vault, Folder, Secret, Field, Template — domain entities        │
│  Business rules, lifecycle states, invariant enforcement         │
└──────────┬───────────────────────────────┬───────────────────────┘
           │                               │
           ▼                               ▼
┌──────────────────────┐     ┌─────────────────────────────────────┐
│   internal/crypto/   │     │        internal/storage/            │
│  Argon2id key deriv. │     │  .abditum binary format             │
│  AES-256-GCM AEAD    │     │  atomic save (tmp → rename)         │
│  crypto/rand salts   │     │  version migration on load          │
│  Zero stdlib deps    │     │  external change detection          │
└──────────────────────┘     └─────────────────────────────────────┘
```

**Key invariant:** `internal/tui` never calls `internal/crypto` or `internal/storage` directly. Everything goes through `vault.Manager`. The TUI is a pure consumer of the Manager's API surface.

---

## Component Responsibilities

### `cmd/abditum/`
- `main.go`: reads optional CLI path arg, instantiates `vault.Manager`, creates root TUI model, calls `tea.NewProgram(rootModel).Run()`
- Zero business logic. Return code: 0 on clean exit, 1 on unrecoverable error.

### `internal/vault/`

**Entities (read-only externally):**
| Type | Responsibility |
|------|---------------|
| `Vault` | Top-level container: folders, templates, settings |
| `Folder` | Recursive subtree: ID, name, parent, ordered children (folders + secrets) |
| `Secret` | Credential record: ID, name, templateID, fields, lifecycle state |
| `Field` | Named value: ID, name, type (common / sensitive), value |
| `Template` | Field schema: ID, name, ordered field definitions |
| `Settings` | Configurable timeouts: auto-lock, reveal, clipboard |

**Manager API surface (representative, not exhaustive):**
```go
// --- Vault lifecycle ---
func (m *Manager) Create(path string, password []byte) error
func (m *Manager) Open(path string, password []byte) error
func (m *Manager) Save() error
func (m *Manager) SaveAs(path string) error
func (m *Manager) Lock()
func (m *Manager) Discard() error           // reload from disk
func (m *Manager) ChangeMasterPassword(current, next []byte) error
func (m *Manager) Export(path string) error
func (m *Manager) Import(path string) error

// --- Read-only state ---
func (m *Manager) Vault() *Vault            // snapshot for TUI reads
func (m *Manager) IsModified() bool
func (m *Manager) IsLocked() bool
func (m *Manager) CurrentPath() string
func (m *Manager) Settings() Settings

// --- Secrets ---
func (m *Manager) CreateSecret(folderID string, templateID string) (*Secret, error)
func (m *Manager) UpdateSecret(id string, changes SecretChanges) error
func (m *Manager) SoftDeleteSecret(id string) error
func (m *Manager) RestoreSecret(id string) error
func (m *Manager) MoveSecret(id, destFolderID string) error
func (m *Manager) ReorderSecret(id string, newIndex int) error
func (m *Manager) FavoriteSecret(id string, fav bool) error
func (m *Manager) DuplicateSecret(id string) (*Secret, error)

// --- Folders, Templates (similar shapes) ---
func (m *Manager) CreateFolder(parentID, name string) (*Folder, error)
// ... etc
```

**Invariant enforcement inside Manager:**
- Only one vault active at a time (Manager is not a pool)
- Lifecycle state transitions of Secret match the state machine in `descricao.md`
- Folder/template uniqueness rules
- Pasta Geral cannot be renamed, moved, or deleted

### `internal/crypto/`
| Function | Purpose |
|----------|---------|
| `DeriveKey(password, salt []byte, params Argon2Params) ([]byte, error)` | Argon2id key derivation |
| `Encrypt(key, plaintext []byte) (ciphertext, nonce []byte, err error)` | AES-256-GCM authenticated encryption |
| `Decrypt(key, ciphertext, nonce []byte) ([]byte, error)` | AES-256-GCM decryption + auth |
| `GenerateSalt() ([]byte, error)` | 32-byte random salt via `crypto/rand` |
| `GenerateNonce() ([]byte, error)` | 12-byte random nonce via `crypto/rand` |

No third-party crypto. Key must be zeroed by caller after use. Package does not retain references to keys or plaintexts.

### `internal/storage/`
| Function | Purpose |
|----------|---------|
| `Load(path string) (RawFile, error)` | Read file, validate magic + version |
| `Decrypt(raw RawFile, key []byte) (VaultPayload, error)` | Decrypt and deserialize payload |
| `Save(path string, payload VaultPayload, key []byte) error` | Encrypt and write atomically |
| `DetectExternalChange(path string, knownStat FileStat) bool` | Compare mtime+size |
| `Migrate(payload VaultPayload) (VaultPayload, error)` | In-memory version migration |

**Atomic save protocol:**
```
1. Encrypt payload → ciphertext
2. Write to <path>.abditum.tmp
3. If <path>.abditum.bak exists: rename to <path>.abditum.bak2
4. Rename <path>.abditum to <path>.abditum.bak
5. Rename <path>.abditum.tmp to <path>.abditum  ← OS atomic on most platforms
6. Delete <path>.abditum.bak2
```

**Binary format layout:**
```
[4 bytes  magic]  [2 bytes version] [argon2_params (fixed size)]
[12 bytes nonce]  [N bytes ciphertext+GCM tag]
```

### `internal/tui/`

Root model holds all session state. Child models are embedded by value (not pointer) — standard Bubble Tea v2 pattern. The root model is the only `tea.Model` passed to `tea.NewProgram`.

---

## Data Flow

### Flow 1: Opening a Vault
```
TUI filePickerModel → path selected
  → TUI sends customMsg{OpenVault, path}
  → User types password in TUI input field
  → TUI sends customMsg{ConfirmOpen, password}
  → Manager.Open(path, password)
    → storage.Load(path)      // reads raw bytes
    → crypto.DeriveKey(...)   // derives AES key from password + stored salt
    → storage.Decrypt(raw, key) // authenticated decryption
    → storage.Migrate(payload)  // upgrade format if needed
    → vault.hydrate(payload)    // populate in-memory entities
  → password zeroed in TUI after call
  → TUI transitions: stateVault
  → lockTimer starts (Manager.Settings().AutoLockDuration)
```

### Flow 2: Editing a Secret
```
User navigates to secret (vaultModel cursor)
  → User presses edit key
  → secretDetailModel enters edit mode (holds provisional copy)
  → User edits fields in local model state (not yet in Manager)
  → User confirms (Enter)
  → TUI calls Manager.UpdateSecret(id, changes)
    → Manager applies changes to domain entity
    → Manager marks Vault as Modified
    → Manager updates secret lifecycle state per rules in descricao.md
  → TUI calls Manager.Vault() → refreshed read-only snapshot
  → vaultModel re-renders
  → Activity timer reset (key press detected in rootModel)
```

### Flow 3: Saving a Vault
```
User presses Ctrl+S → TUI sends SaveCmd
  → Manager.Save()
    → vault.serialize() → VaultPayload
    → crypto.Encrypt(sessionKey, payload) → ciphertext
    → storage.Save(path, ciphertext)     → atomic write
    → Manager records new FileStat (mtime + size)
    → Manager marks Vault as Saved
  → TUI reads Manager.IsModified() → updates status bar indicator
```

### Flow 4: Auto-Lock
```
lockTimer (bubbles/timer.Model) fires timer.TimeoutMsg
  → rootModel.Update intercepts
  → Manager.Lock()
    → zero master password in memory
    → zero AES session key in memory
    → discard vault data
  → TUI clears clipboard (if pending)
  → TUI clears screen (tea.ClearScreen equivalent in View)
  → rootModel transitions → stateWelcome
  → lockTimer stopped
```

### Flow 5: Field Reveal / Clipboard Clear
```
User presses "reveal" on sensitive field
  → TUI records: revealedFields[fieldID] = time.Now().Add(revealDuration)
  → globalTick fires every 1s
  → rootModel checks: for each revealedFields entry, if deadline passed → hide

User presses "copy" on field
  → OS clipboard written
  → TUI records: clipboardFieldID = fieldID, clipboardDeadline = now+clipboardDuration
  → globalTick checks clipboardDeadline → if passed: clear OS clipboard, reset tracking
  → On Lock or Quit: immediately clear OS clipboard regardless of timer
```

---

## Build Order

Build bottom-up. Each layer has zero upward dependencies.

### Layer 1 — `internal/crypto` (no project deps)
**What:** Argon2id + AES-256-GCM primitives, salt/nonce generation
**Tests:** Round-trip encrypt/decrypt, key derivation stability, error paths (wrong key, tampered ciphertext)
**Why first:** Every upper layer depends on this. Simple, pure functions, fast to test.

### Layer 2 — `internal/vault` (depends on: nothing from project)
**What:** Domain entities, Manager, business rules, state machine
**Tests:** All Manager methods; lifecycle state invariants; folder/secret ordering; soft delete/restore; duplicate names; Pasta Geral protection
**Why second:** Pure in-memory logic. No I/O. Can be fully verified before storage layer exists.

### Layer 3 — `internal/storage` (depends on: crypto)
**What:** Binary file format, serialization, atomic save, external change detection, version migration
**Tests:** Save/load round-trips with real temp files; atomic save with simulated mid-write failure; migration from each prior format version
**Why third:** Storage is a consumer of crypto and a producer for Manager.Open/Save.

### Layer 4 — `internal/tui` (depends on: vault, uses storage via Manager)
**Substep 4a:** `welcomeModel` — no vault, welcomes user, offers create/open/quit
**Substep 4b:** `filePickerModel` — path selection for open/save/export/import
**Substep 4c:** `vaultModel` — tree navigation, search overlay, status bar
**Substep 4d:** `secretDetailModel` — view, create, edit, edit-advanced modes
**Substep 4e:** `modalModel` + `helpModel` — overlays
**Substep 4f:** Timer integration — lock, clipboard, reveal
**Tests:** teatest/v2 for each screen; golden files at 80×24

### Layer 5 — `cmd/abditum` (wires all layers)
**What:** `main.go`, optional CLI path arg, entry point
**Tests:** Integration test: full round-trip from create vault to open vault to edit secret to save

---

## Bubble Tea v2 Patterns for This Project

### Import Paths (v2 uses charm.land vanity domain)
```go
import (
    tea      "charm.land/bubbletea/v2"
    "charm.land/bubbles/v2/key"
    "charm.land/bubbles/v2/timer"
    "charm.land/bubbles/v2/textinput"
    "charm.land/bubbles/v2/viewport"
    "charm.land/bubbles/v2/help"
    "charm.land/lipgloss/v2"
)
```

### Pattern 1: Session State Enum for Screen Routing

```go
// tui/root.go
type sessionState int
const (
    stateWelcome     sessionState = iota
    stateVault                      // tree + search
    stateSecretDetail               // view/edit secret
    stateModal                      // overlay: warning / confirmation
    stateFilePicker                 // overlay: path selection
    stateHelp                       // overlay: key bindings
    stateTooSmall                   // overlay: terminal too small
)

type rootModel struct {
    state     sessionState
    prevState sessionState  // restored after overlay dismissed

    welcome      welcomeModel
    vault        vaultModel
    secretDetail secretDetailModel
    modal        modalModel
    filePicker   filePickerModel
    help         helpModel

    manager *vault.Manager
    width   int
    height  int

    // Timers
    lockTimer      timer.Model
    clipboardState clipboardState  // fieldID + deadline
    revealedFields map[string]time.Time  // fieldID → expiry
}
```

### Pattern 2: Overlay vs. Base Screen Message Routing

Overlays (modal, file picker, help, too-small) consume all input before the base screen sees it. This is enforced in root's `Update()` priority order:

```go
func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    // 1. Reset activity timer on any user interaction
    switch msg.(type) {
    case tea.KeyPressMsg, tea.MouseClickMsg:
        if !m.manager.IsLocked() {
            var cmd tea.Cmd
            m.lockTimer, cmd = m.lockTimer.Update(timer.ResetMsg{})
            cmds = append(cmds, cmd)
        }
    }

    // 2. Global messages handled regardless of state
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        return m.handleWindowSize(msg), tea.Batch(cmds...)
    case timer.TimeoutMsg:
        if msg.ID == m.lockTimer.ID() {
            return m.doLock(), tea.Batch(cmds...)
        }
    case globalTickMsg:
        m = m.checkTimers(msg.t)
        cmds = append(cmds, globalTick())
    }

    // 3. Overlay states intercept all remaining input
    switch m.state {
    case stateModal:
        var cmd tea.Cmd
        m.modal, cmd = m.modal.Update(msg)
        cmds = append(cmds, cmd)
        if result, done := m.modal.Result(); done {
            m = m.handleModalResult(result)
        }
        return m, tea.Batch(cmds...)

    case stateFilePicker:
        // ... similar delegation

    case stateHelp, stateTooSmall:
        // dismiss on any key press
    }

    // 4. Base state delegation
    switch m.state {
    case stateWelcome:
        var cmd tea.Cmd
        m.welcome, cmd = m.welcome.Update(msg)
        cmds = append(cmds, cmd)
    case stateVault:
        var cmd tea.Cmd
        m.vault, cmd = m.vault.Update(msg)
        cmds = append(cmds, cmd)
    case stateSecretDetail:
        var cmd tea.Cmd
        m.secretDetail, cmd = m.secretDetail.Update(msg)
        cmds = append(cmds, cmd)
    }

    return m, tea.Batch(cmds...)
}
```

### Pattern 3: Custom Messages for Cross-Screen Communication

Child models MUST NOT directly mutate state or access the Manager. They return custom `tea.Msg` values that the root model processes. This enforces unidirectional data flow.

```go
// tui/messages.go — all inter-screen messages defined here

type openVaultRequestMsg struct{ path string }
type openVaultSuccessMsg struct{}
type openVaultFailedMsg  struct{ err error }

type lockVaultMsg struct{}

type showModalMsg struct {
    kind    modalKind  // warn / confirm
    message string
    onConfirm tea.Msg  // message to send if user confirms
}

type showFilePickerMsg struct {
    mode filePickerMode  // open / save / export / import
}

type filePickerDoneMsg struct{ path string }

type navigateToSecretMsg struct{ secretID string }
type returnToVaultMsg    struct{}

type clipboardCopiedMsg  struct{ fieldID, value string }
type fieldRevealedMsg    struct{ fieldID string }
```

Root intercepts these and drives transitions:
```go
case openVaultSuccessMsg:
    m.state = stateVault
    dur := m.manager.Settings().AutoLockDuration
    m.lockTimer = timer.New(dur)
    cmds = append(cmds, m.lockTimer.Init(), globalTick())
```

### Pattern 4: Key Bindings with `bubbles/key`

Each screen defines its own `keyMap` struct using `key.Binding`. The root model holds only global bindings (quit, lock). Screen models hold screen-specific bindings and render them via `bubbles/help`.

```go
// tui/vault/keybindings.go
type vaultKeyMap struct {
    Up      key.Binding
    Down    key.Binding
    Search  key.Binding
    New     key.Binding
    Edit    key.Binding
    Delete  key.Binding
    Move    key.Binding
    Lock    key.Binding
    Save    key.Binding
    SaveAs  key.Binding
    Help    key.Binding
    Quit    key.Binding
}

func defaultVaultKeyMap() vaultKeyMap {
    return vaultKeyMap{
        Search: key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
        New:    key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
        Edit:   key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
        Lock:   key.NewBinding(key.WithKeys("ctrl+l"), key.WithHelp("ctrl+l", "lock")),
        Save:   key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")),
        // ...
    }
}
```

Key binding delegation is explicit: root's `Update()` does NOT intercept keyboard for screens; each screen's `Update()` handles its own keys. Only global actions (quit, lock) are intercepted at root based on matching the `lockKey` or `quitKey` bindings before delegation.

### Pattern 5: View Composition with AltScreen

```go
// tui/root.go
func (m rootModel) View() tea.View {
    var content string

    // Render base screen
    switch m.state {
    case stateWelcome:
        content = m.welcome.View()
    case stateVault, stateModal, stateFilePicker, stateHelp:
        content = m.vault.View()  // base content always rendered
    case stateSecretDetail:
        content = m.secretDetail.View()
    case stateTooSmall:
        content = m.renderTooSmallWarning()
    }

    // Render overlay on top using lipgloss.Place
    switch m.state {
    case stateModal:
        content = renderOverlay(content, m.modal.View(), m.width, m.height)
    case stateFilePicker:
        content = renderOverlay(content, m.filePicker.View(), m.width, m.height)
    case stateHelp:
        content = renderOverlay(content, m.help.View(), m.width, m.height)
    }

    v := tea.NewView(content)
    v.AltScreen = true  // Always full-screen
    return v
}
```

**Important v2 note:** `AltScreen`, `WindowTitle`, and cursor settings are all declared on `tea.View` — not as `tea.NewProgram` options or commands. There is no `tea.WithAltScreen()` in v2.

### Pattern 6: Window Size Handling

```go
func (m rootModel) handleWindowSize(msg tea.WindowSizeMsg) rootModel {
    m.width, m.height = msg.Width, msg.Height

    inSmallTerminal := msg.Width < minWidth || msg.Height < minHeight
    if inSmallTerminal && m.state != stateTooSmall {
        m.prevState = m.state
        m.state = stateTooSmall
    } else if !inSmallTerminal && m.state == stateTooSmall {
        m.state = m.prevState
    }

    // Propagate size to children that need it (viewport-based models)
    m.vault.SetSize(msg.Width, msg.Height)
    m.secretDetail.SetSize(msg.Width, msg.Height)
    // Overlays compute their own size from m.width/m.height at render time
    return m
}
```

---

## Timer / Tick Architecture

Three timers with different behavioral requirements.

### Timer 1: Auto-Lock (activity-reset countdown)

**Behavior:** countdown reset by any user interaction; lock on expiry or on manual Ctrl+L.

**Implementation:** `bubbles/timer.Model` — it has built-in start/stop/reset, fires `timer.TimeoutMsg` and periodic `timer.TickMsg`.

```go
// In Init():
cmds = append(cmds, m.lockTimer.Init())

// In Update(), before delegation:
case tea.KeyPressMsg, tea.MouseClickMsg:
    m.lockTimer, cmd = m.lockTimer.Update(timer.ResetMsg{})
    cmds = append(cmds, cmd)

// In Update(), timer expiry:
case timer.TimeoutMsg:
    if msg.ID == m.lockTimer.ID() {
        return m.doLock(), nil
    }

// doLock():
func (m rootModel) doLock() rootModel {
    m.clearClipboard()        // OS clipboard clear
    m.manager.Lock()          // zero keys in vault.Manager
    m.revealedFields = nil    // drop all reveal deadlines
    m.state = stateWelcome
    m.lockTimer = timer.Model{} // discard
    return m
}
```

**The `ID()` check is critical** when multiple `timer.Model` instances exist — `timer.TickMsg` is keyed by ID.

### Timer 2: Clipboard Clear (deadline, single in-flight)

**Behavior:** after user copies a field, clear OS clipboard after N seconds. Last copy wins. Clear immediately on lock/quit.

**Implementation:** Simple `tea.Tick` + custom message — no `timer.Model` needed. Track the deadline and field in root state.

```go
type clipboardState struct {
    copiedAt time.Time
    duration time.Duration
    active   bool
}

// On copy action:
case clipboardCopiedMsg:
    setOSClipboard(msg.value)
    m.clipboard = clipboardState{
        copiedAt: time.Now(),
        duration: m.manager.Settings().ClipboardClearDuration,
        active:   true,
    }
    // globalTick will catch expiry at next 1s check

// In checkTimers() (called on globalTickMsg):
if m.clipboard.active {
    if time.Now().After(m.clipboard.copiedAt.Add(m.clipboard.duration)) {
        clearOSClipboard()
        m.clipboard = clipboardState{}
    }
}
```

### Timer 3: Field Reveal (multiple concurrent deadlines)

**Behavior:** sensitive field shown temporarily; multiple fields can be revealed simultaneously; each hides independently.

**Implementation:** `map[string]time.Time` (fieldID → hide deadline) + global 1s tick.

```go
type globalTickMsg struct{ t time.Time }

func globalTick() tea.Cmd {
    return tea.Tick(time.Second, func(t time.Time) tea.Msg {
        return globalTickMsg{t: t}
    })
}

// On reveal action:
case fieldRevealedMsg:
    m.revealedFields[msg.fieldID] = time.Now().Add(
        m.manager.Settings().FieldRevealDuration,
    )

// In checkTimers():
now := time.Now()
for id, deadline := range m.revealedFields {
    if now.After(deadline) {
        delete(m.revealedFields, id)
    }
}
```

**Global tick starts when vault opens, stops on lock.** A single 1-second tick covers both clipboard clear and field reveal — no separate goroutines needed.

### Timer Lifecycle Summary

| Timer | Start | Stop | Reset |
|-------|-------|-------|-------|
| Auto-lock | vault opens | vault locks | any key/mouse event |
| Clipboard | field copied | vault locks or deadline | next copy (last-write-wins) |
| Field reveal | reveal action | vault locks or deadline per field | not resettable; new action updates deadline |

### Security: On Lock/Quit, Always Cancel All Timers First

```go
func (m rootModel) doLock() rootModel {
    // 1. Clear clipboard immediately (before zero-ing Manager)
    clearOSClipboard()
    m.clipboard = clipboardState{}
    // 2. Hide all revealed fields
    m.revealedFields = nil
    // 3. Zero memory in Manager (password, AES key, vault data)
    m.manager.Lock()
    // 4. Transition UI
    m.state = stateWelcome
    return m
}
```

---

## Testing Architecture

### Unit Tests (fast, no TUI)

**`internal/crypto/`**
```go
// crypto_test.go
func TestEncryptDecryptRoundTrip(t *testing.T)
func TestDecryptFailsWithWrongKey(t *testing.T)
func TestDecryptFailsWithTamperedCiphertext(t *testing.T)
func TestArgon2idDeterministic(t *testing.T)  // same password+salt → same key
func TestGenerateSaltUnique(t *testing.T)
```

**`internal/vault/`**
```go
// manager_test.go
func TestCreateVaultHasDefaultStructure(t *testing.T)
func TestSoftDeleteSecretPreservesState(t *testing.T)
func TestRestoreSecretReturnsToOriginalFolder(t *testing.T)
func TestSecretLifecycleTransitions(t *testing.T)  // all state machine edges
func TestFolderUniquenessWithinParent(t *testing.T)
func TestPastaGeralProtection(t *testing.T)        // rename/move/delete blocked
func TestOnlyOneVaultActive(t *testing.T)
func TestModifiedFlagSetOnMutation(t *testing.T)
func TestFavoriteDoesNotChangeLifecycleState(t *testing.T)
```

**`internal/storage/`**
```go
// storage_test.go
func TestSaveAndLoadRoundTrip(t *testing.T)
func TestAtomicSaveCreatesBackup(t *testing.T)
func TestExternalChangeDetection(t *testing.T)
func TestVersionMigration(t *testing.T)  // for each format version increment
func TestCorruptedFileReturnsError(t *testing.T)
```

### TUI Tests with teatest/v2

Import: `github.com/charmbracelet/x/exp/teatest/v2`

```go
// tui/welcome_test.go
func TestWelcomeScreenRendersCorrectly(t *testing.T) {
    mgr := vault.NewManager()
    m := tui.NewRootModel(mgr)

    tm := teatest.NewTestModel(t, m,
        teatest.WithInitialTermSize(80, 24),
    )

    // Wait for stable render
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("Abditum"))
    }, teatest.WithDuration(3*time.Second))

    // Golden file comparison (visual regression)
    golden.RequireEqual(t, tm.FinalOutput(t))
    tm.Quit(t)
}
```

**For interaction tests:**
```go
func TestVaultSearchTriggeredBySlash(t *testing.T) {
    tm := teatest.NewTestModel(t, openedVaultModel(t),
        teatest.WithInitialTermSize(80, 24),
    )
    tm.Send(tea.KeyPressMsg{Code: '/'})
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("Search"))
    })
    tm.Quit(t)
}
```

**Golden file workflow:**
- Store `.golden` files under `testdata/` next to each `_test.go` file
- Generate/update: `UPDATE_GOLDEN=1 go test ./...`
- CI runs without flag — deviation from golden = test failure
- One golden file per screen × terminal size (use 80×24 throughout)
- Use `tea.WithColorProfile(termenv.TrueColor)` in test options for consistent color output

### Integration Tests (end-to-end, no TUI)

```go
// integration/vault_lifecycle_test.go
func TestFullVaultLifecycle(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "test.abditum")

    mgr := vault.NewManager()
    require.NoError(t, mgr.Create(path, []byte("password123!")))
    require.True(t, mgr.IsModified() == false)  // just saved

    _, err := mgr.CreateFolder(mgr.Vault().RootFolderID(), "Work")
    require.NoError(t, err)

    s, err := mgr.CreateSecret(workFolderID, loginTemplateID)
    require.NoError(t, err)

    require.NoError(t, mgr.Save())
    require.False(t, mgr.IsModified())

    // Reopen
    mgr2 := vault.NewManager()
    require.NoError(t, mgr2.Open(path, []byte("password123!")))
    require.Equal(t, 1, len(mgr2.Vault().Folder(workFolderID).Secrets()))
    _ = s
}
```

---

## Notes for Abditum

### N1: Manager Owns the Session Key — TUI Is Stateless on Crypto
`internal/tui` never sees the AES session key or the master password after handing them to `Manager.Open()` / `Manager.Create()`. The TUI zeroes its local password buffer immediately after the call returns. The Manager internally holds the session key for re-encryption on save.

### N2: `vault.Manager.Vault()` Returns a Snapshot — Not a Live Reference
`Vault()` returns a value snapshot of the current vault tree. TUI reads from this snapshot; mutations happen via Manager methods, then TUI calls `Vault()` again for the next render. This prevents the TUI from accidentally holding a stale reference to a mutated subtree.

### N3: No Goroutines in TUI Layer — Use tea.Tick and tea.Cmd
All async work (timers, clipboard, reveals) uses `tea.Tick` / `tea.Cmd`. No `go func(){}` in TUI code. Operations like `Manager.Save()` can block briefly — if needed, wrap in `tea.Cmd` so Bubble Tea runs them off the main update loop:
```go
func doSave(mgr *vault.Manager) tea.Cmd {
    return func() tea.Msg {
        if err := mgr.Save(); err != nil {
            return saveFailedMsg{err: err}
        }
        return saveSuccessMsg{}
    }
}
```

### N4: Screen Minimum Size
Track `minWidth = 80, minHeight = 24` as constants. The `stateTooSmall` overlay is activated in `handleWindowSize` and reactivated on every `WindowSizeMsg`. It should render a plain message without any decorative UI that might itself require space.

### N5: `tea.Batch` for Multiple Commands
```go
// Correct: use tea.Batch when returning multiple cmds
return m, tea.Batch(lockTimerReset, globalTick())
// Inside Update when firing multiple actions:
cmds := []tea.Cmd{ ... }
return m, tea.Batch(cmds...)
```

### N6: teatest/v2 Import Path
As of 2026-03, teatest v2 has been migrated to use `charm.land` imports internally. The package import path is still:
```
github.com/charmbracelet/x/exp/teatest/v2
```
Verify `go mod tidy` picks up the latest version. The `golden` helper is at `github.com/charmbracelet/x/exp/golden`.

### N7: File Picker — Use Custom Rather Than Bubbles
The Bubbles `filepicker` component has dependencies that may conflict with `CGO_ENABLED=0` on some platforms (directory listing via OS APIs). Evaluate at build time. A custom path-input text field with tab-completion is always safe and sufficient for this use case (open, save-as, export path).

### N8: Memory Zeroing Limitation in Go
Go's GC may move slices, meaning `copy(slice, zeros)` on a sensitive `[]byte` is best-effort, not guaranteed. Accept this limitation and document it. Still do it — it reduces the exposure window even if not a cryptographic guarantee. Never use `string` for passwords; always `[]byte`.

### N9: Recommended Bubbles Components
| Component | Use Case |
|-----------|----------|
| `bubbles/textinput` | Password entry, search bar, field editing, name inputs |
| `bubbles/textarea` | Observation/note fields (multi-line) |
| `bubbles/viewport` | Scrollable secret detail pane |
| `bubbles/key` | Key bindings + help rendering |
| `bubbles/help` | Help bar at bottom of screens |
| `bubbles/timer` | Auto-lock countdown only |
| `bubbles/list` | Folder/secret tree (if using list bubble) OR custom tree renderer |

The tree/hierarchy view may need a custom renderer rather than `bubbles/list`, because:
- Nested indentation with folder collapsing is not built into `bubbles/list`
- Dual-type items (folders + secrets with different icons/states) need rendering control

---

## Sources

- Bubble Tea v2 official docs: https://pkg.go.dev/charm.land/bubbletea/v2
- v2 upgrade guide: https://github.com/charmbracelet/bubbletea/blob/main/UPGRADE_GUIDE_V2.md
- Composable views example: https://github.com/charmbracelet/bubbletea/blob/main/examples/composable-views/main.go (HIGH — official repo)
- Timer example: https://github.com/charmbracelet/bubbletea/blob/main/examples/timer/main.go (HIGH — official repo)
- teatest v2 source: https://github.com/charmbracelet/x/blob/main/exp/teatest/v2/teatest.go (HIGH — official repo)
- Abditum arquitetura.md + descricao.md: project source documents (HIGH — authoritative)
