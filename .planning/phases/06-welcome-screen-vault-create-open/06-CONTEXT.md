# Phase 6: Welcome Screen + Vault Create/Open - Context

**Gathered:** 2026-04-02
**Status:** Ready for planning

<domain>
## Phase Boundary

This phase delivers the first end-to-end user flows: creating a new vault and opening an existing vault — from the welcome screen through file selection, master password input, error handling, and successful vault load. Every error case surfaces the correct generic user message with no technical detail leaked.

This phase implements:
- `preVaultModel` upgraded with action key hints (keys already registered in FlowRegistry from Phase 5)
- **`headerModel`** — new component implementing the 2-line header; two structural states: `stateNoVault` (welcome) and `stateVaultOpen` (vault open with 3 mode tabs); renders app name, `·` separator, vault name (truncated), `•` dirty indicator, and tab strip with the active tab suspended on the separator line; `rootModel` passes it the current vault name, dirty state, and active workArea on every render
- `filePickerModal` — new modal struct for selecting or saving vault files (tree + file panel, Open and Save modes), used by both flows
- `passwordEntryModal` — new modal struct for vault-open password input (single field, attempt counter up to 5)
- `passwordCreateModal` — new modal struct for vault-create password input (two fields, Tab navigation, inline strength meter)
- `openVaultFlow.Update()` fully implemented: FilePicker(Open) → PasswordEntry → RecoverOrphans → storage.Load → error classification → emit domain message
- `createVaultFlow.Update()` fully implemented: FilePicker(Save) → PasswordCreate → Manager.Create → storage.Save → emit domain message
- CLI path fast-path: if `initialPath` non-empty at startup, `openVaultFlow` skips FilePicker and goes straight to PasswordEntry
- `dialogs.go` extended with new factory helpers for the Phase 6 error dialogs
- **Theme infrastructure**: `Theme` struct with all token values, two theme instances (`ThemeTokyoNight`, `ThemeCyberpunk`), `theme.go` new file; `rootModel` holds active theme and broadcasts on toggle
- **F12 global shortcut**: toggles active theme in-memory; immediate re-render; registered in `ActionManager` alongside `ctrl+q` and `?`
- `RenderLogo()` in `ascii.go` updated to accept a `Theme` and use per-theme gradient colors
- All new modals (`filePickerModal`, `passwordEntryModal`, `passwordCreateModal`) receive `*Theme` at construction and use it for all `lipgloss` color values

This phase does NOT implement:
- Vault tree display or secret navigation (Phase 7)
- Save / Save As / Discard / Change Password / Export / Import (Phase 9)
- Security timers, clipboard, lock/exit flows (Phase 10)
- Any new workArea transition beyond `workAreaPreVault → workAreaVault` (and back on error)
- Theme persistence to disk (saved in settings — Phase 9); Phase 6 default is Tokyo Night, F12 toggles in-memory only
- Tab navigation between modes (the tab strip is rendered but mode-switching F-keys are Phase 7+ scope)

</domain>

<decisions>
## Implementation Decisions

### Header Component

**D-12: `headerModel` — 2-line component with two structural states**

The header is implemented as a standalone struct (not a `childModel` — it never receives raw `tea.Msg`). `rootModel` calls `headerModel.Render(width int) string` from its `renderFrame()` method, passing current state.

**Inputs at render time (passed by rootModel):**
- `width int` — current terminal width
- `vaultName string` — radical of vault filename (no `.abditum` extension); empty string = no vault
- `isDirty bool` — unsaved changes present
- `activeArea workArea` — current work area (determines active tab)

**State 1 — No vault (welcome):**
```
  Abditum
────────────────────────────────────────────────────────────────────────────────
```
- Line 1: `  Abditum` — app name in `accent.primary` bold
- Line 2: full-width `─` separator in `border.default`; no tab connectors

**State 2 — Vault open:**
```
  Abditum · cofre •                         ╭───────╮  ╭ Modelos ╮  ╭ Config ╮
─────────────────────────────────────────╯ Cofre ╰────────────────
```
- `·` separator: `border.default`
- Vault name (radical only, truncated with `…` per spec algorithm): `text.secondary`
- `•` dirty indicator: `semantic.warning`; hidden when `!isDirty`
- Tab text inactive: `text.secondary`; active: `accent.primary` bold on `special.highlight` background
- Tab borders (`╭╮╯╰─`): `border.default` in both states
- Separator line: `border.default`

**Tabs:** Cofre (`workAreaVault`), Modelos (`workAreaTemplates`), Config (`workAreaSettings`). In Phase 6, only `workAreaVault` is reachable — the other tabs render in inactive style but are non-functional until Phase 7+.

**File layout:** `header.go` — new file in `internal/tui/`

### Welcome Screen — Action Model

**D-01: `preVaultModel` stays display-only; flows are dispatched via FlowRegistry**

Phase 5 CONTEXT (D-02, D-09) is authoritative: `preVaultModel` has no sub-states and manages no sub-flows. The ROADMAP's description of a `welcomeModel` with `j/k/Enter` menu is superseded by the Phase 5 architecture decisions.

Concretely:
- `preVaultModel.Update()` handles no keyboard input — keys fall through to FlowRegistry dispatch (D-06 priority order, step 4)
- `openVaultDescriptor` (key `"o"`) and `createVaultDescriptor` (key `"n"`) were already registered in Phase 5 — they now receive real implementations
- `preVaultModel.View()` is upgraded to render logo + two action hints beneath it:
  ```
  n  Novo cofre    o  Abrir cofre
  ```
  Hints use `text.secondary` color via `renderHints()` helper (already defined in `prevault.go`)
- `ctrl+q` quit flow remains a global shortcut (D-12 from Phase 5)
- No `q` shortcut on the welcome screen — `ctrl+q` is the global quit

### Flow Architecture — Modal Sequencing

**D-02: Flows orchestrate their own modal sequences; each flow is a state machine**

Both `openVaultFlow` and `createVaultFlow` are internal state machines. They push modals via `pushModalMsg` and react to domain messages returned from those modals. The flow's `Update()` method is the single dispatch point.

**Open vault flow state machine:**
```
statePickFile → (FilePicker modal pushed, wait for filePickedMsg or flowCancelledMsg)
statePwdEntry → (PasswordEntry modal pushed, wait for pwdEnteredMsg or flowCancelledMsg)
statePwdEntry handles: retry (wrong password), hard errors (corrupt/invalid/too new → error modal → statePickFile or close)
stateLoading  → (RecoverOrphans + storage.Load called via tea.Cmd goroutine)
→ success: emit vaultOpenedMsg{path}, rootModel transitions to workAreaVault
→ failure: see error classification (D-04)
```

**Create vault flow state machine:**
```
statePickFile → (FilePicker(Save mode) pushed, wait for filePickedMsg or flowCancelledMsg)
statePwdCreate → (PasswordCreate modal pushed, wait for pwdCreatedMsg or flowCancelledMsg)
statePickFile handles: overwrite confirmation if target exists (D-05)
stateSaving   → (Manager.Create + storage.SaveNew called via tea.Cmd goroutine)
→ success: emit vaultOpenedMsg{path}, rootModel transitions to workAreaVault
→ failure: error modal, return to statePwdCreate
```

**Flow cancellation at any step:** flow sets `activeFlow = nil` in rootModel; preVault resumes keyboard focus.

### New Modal Types

**D-03: Three new modal structs in Phase 6; all satisfy `childModel`**

Unlike the generic `modalModel` from Phase 5 (which handles simple push/pop with text options), Phase 6 introduces three dedicated modal structs with their own state:

1. **`filePickerModal`** (`filepicker.go` — new file)
   - Shared by open and create vault flows; mode is a constructor parameter (`FilePickerModeOpen`, `FilePickerModeSave`)
   - State: directory tree (expanded dirs, selected dir, scroll), file list (selected file, scroll), filename input field (Save mode only), current path string
   - Implements DS/spec exactly: two-panel layout (Estrutura ~40% / Arquivos ~60%), `border.focused` border, path header, file filter `*.abditum`, auto-select first file on directory navigation
   - Emits `filePickedMsg{path string}` on confirm; `flowCancelledMsg{}` on ESC
   - `Tab` cycles focus between tree panel → file panel → filename field (Save mode only)
   - **`charm.land/bubbles/v2/filepicker` exists** (`fp.AllowedTypes`, `fp.CurrentDirectory`) but provides a **single-panel** layout that does not match the two-panel spec. Researcher must assess whether it can be wrapped to produce the required layout or whether a custom implementation using `os.ReadDir` is needed. If the bubbles component cannot satisfy the spec's two-panel structure, build custom.
   - Directory tree traversal (if custom): `os.ReadDir`; no symlink recursion; shows all directories, filters files to `*.abditum` only
   - Displays relative dates (e.g., `1h`, `3d`) and sizes in human-readable form
   - **Spec reference:** `tui-specification-novo.md` §FilePicker (both modes)

2. **`passwordEntryModal`** (`passwordentry.go` — new file)
   - Fixed width: 50 columns; border token: `border.focused`
   - State: password `textinput` (masked, `EchoMode = textinput.EchoPassword`, `EchoCharacter = '•'`), attempt counter `int`, max attempts: **5**
   - Counter line (`Tentativa N de 5`) is **hidden on attempt 1**, shown from attempt 2 onward
   - Action default (`Enter Confirmar`) blocked when field empty; unlocked when non-empty
   - On wrong password: field cleared, counter incremented, action default locked again
   - On attempt 5 exhausted: dialog closes automatically (emits `flowCancelledMsg{}`)
   - MessageManager hints:
     - On open / field empty or valid: `• Digite a senha para desbloquear o cofre`
     - On wrong password: `✗ Senha incorreta` (5s error)
   - Emits `pwdEnteredMsg{password []byte}` on confirm; immediately zeros the textinput buffer after extracting bytes
   - `Tab` does nothing (single field)
   - **Spec reference:** `tui-specification-novo.md` §PasswordEntry

3. **`passwordCreateModal`** (`passwordcreate.go` — new file)
   - Fixed width: 50 columns; border token: `border.focused`
   - State: two `textinput` fields (both `EchoMode = textinput.EchoPassword`, `EchoCharacter = '•'`), focused field index (`0` = Nova senha, `1` = Confirmação)
   - `Tab` toggles between fields
   - Strength meter row is **hidden when field 1 empty**, appears when field 1 has content; renders `Força: ████████░░ Boa` or `Força: ████░░░░░░ Fraca` using `semantic.success`/`semantic.warning`
   - Password strength calls `crypto.EvaluatePasswordStrength` on every keystroke in field 1
   - Action default blocked when either field empty; unlocked when both non-empty (weakness does NOT block submit — non-blocking per ROADMAP plan)
   - On Enter with passwords divergent: confirm field cleared, focus returns to confirm field, action default locked, error message `✗ As senhas não conferem — digite novamente` (5s)
   - On confirmed submit: extracts password bytes from field 1, zeros both textinput buffers, emits `pwdCreatedMsg{password []byte}`
   - MessageManager hints:
     - Focus on field 1 (empty or valid): `• A senha mestra protege todo o cofre — use 12+ caracteres`
     - Focus on field 2 (empty or valid): `• Redigite a senha para confirmar`
     - After divergence error, re-focus on field 2: `✗ As senhas não conferem — digite novamente` (5s)
   - **Spec reference:** `tui-specification-novo.md` §PasswordCreate

### Error Classification (VAULT-04)

**D-04: Error sentinel mapping — no Go error strings, no internal detail ever exposed**

| Error sentinel | User-visible message | UX action |
|---|---|---|
| `storage.ErrAuthFailed` | `✗ Senha incorreta` (in MessageManager, 5s) | Clear input, increment attempt counter, allow retry |
| `storage.ErrInvalidMagic` | `O arquivo selecionado não é um cofre Abditum` | Recognition × Error modal; on dismiss → flow returns to `statePickFile` |
| `storage.ErrVersionTooNew` | `Este cofre foi criado por uma versão mais recente do Abditum` | Recognition × Error modal; on dismiss → flow returns to `statePickFile` |
| `storage.ErrCorrupted` | `O arquivo está corrompido e não pode ser aberto` | Recognition × Error modal; on dismiss → flow returns to `statePickFile` (no retry at pwd step) |
| save failure on create | `Não foi possível criar o cofre — verifique o caminho e as permissões` | Recognition × Error modal; on dismiss → flow returns to `statePwdCreate` (path already selected) |

Recognition × Error dialogs follow the spec: `✗` symbol, `semantic.error` border, `Enter OK` action.

### Overwrite Confirmation (Create Flow)

**D-05: Target file already exists → Confirmation × Destructive modal before creating**

Per `fluxos.md` Fluxo 2 step 1: if the chosen save path already exists, show a confirmation dialog before proceeding.

- After FilePicker returns a path, `createVaultFlow` checks `os.Stat(path)` — if file exists, push confirmation dialog
- Dialog: Confirmação × Destrutivo — `⚠ Sobrescrever arquivo?` / `S Sobrescrever` (semantic.error) / `Esc Cancelar` (semantic.warning)
- On confirm → proceed to `statePwdCreate`
- On cancel → return to `statePickFile` (FilePicker reopened)

### CLI Path Fast-Path

**D-06: If `initialPath` is non-empty at startup, `openVaultFlow` skips FilePicker**

Per `fluxos.md` Fluxo 1 ("Entrada antecipada via argumento de linha de comando"):

- `rootModel.Init()` checks: if `m.vaultPath` is non-empty → immediately dispatch `openVaultDescriptor.New(ctx)` and set `activeFlow`
- `openVaultFlow` starts in `statePreVerify` when constructed with a pre-filled path: calls `storage.Probe(path)` (or equivalent header read) to verify the path is a valid `.abditum` file before pushing PasswordEntry
  - If valid → push PasswordEntry directly (no FilePicker pushed)
  - If invalid → show Recognition × Error modal; on dismiss → flow transitions to `statePickFile` (open FilePicker normally)
- This fast-path **only applies to the open flow** — create vault flow has no CLI path argument

### Timing of `activeFlow` Reference

**D-07: `rootModel.Init()` (not `newRootModel()`) triggers CLI fast-path**

`rootModel.Init()` is called once by Bubble Tea before the event loop starts. This is the correct place to return the `tea.Cmd` that starts the open flow when `m.vaultPath` is pre-filled. `newRootModel()` must not block on I/O.

### `RecoverOrphans` Call

**D-08: `RecoverOrphans` called in tea.Cmd goroutine before `storage.Load`**

Per ROADMAP plan: before calling `storage.Load`, call `storage.RecoverOrphans(path)`. Orphan recovery errors are **ignored silently** — the cmd logs nothing (SEC-PRIV-01: no paths or internal errors in stdout/stderr). Even if RecoverOrphans fails, Load proceeds.

### New Domain Messages

**D-09: Phase 6 introduces these new domain message types**

```go
// filePickedMsg carries the absolute path chosen in FilePicker
type filePickedMsg struct{ path string }

// pwdEnteredMsg carries the raw password bytes from PasswordEntry
// Bytes are zeroed in the modal immediately after being copied here
type pwdEnteredMsg struct{ password []byte }

// pwdCreatedMsg carries the confirmed new password from PasswordCreate
type pwdCreatedMsg struct{ password []byte }

// vaultOpenedMsg signals successful vault open/create to rootModel
// rootModel transitions workArea to workAreaVault upon receiving this
type vaultOpenedMsg struct{ path string }

// flowCancelledMsg signals the user cancelled the active flow at any step
// rootModel sets activeFlow = nil upon receiving this
type flowCancelledMsg struct{}
```

These are in addition to the existing message taxonomy from Phase 5 (D-07).

### Theme System

**D-11: `Theme` struct as the single source of all color token values; passed to every child**

```go
// Theme holds all design-system color token values for one palette.
type Theme struct {
    // Surfaces
    SurfaceBase   lipgloss.Color
    SurfaceRaised lipgloss.Color
    SurfaceInput  lipgloss.Color
    // Text
    TextPrimary   lipgloss.Color
    TextSecondary lipgloss.Color
    TextDisabled  lipgloss.Color
    TextLink      lipgloss.Color
    // Borders
    BorderDefault lipgloss.Color
    BorderFocused lipgloss.Color
    // Interaction
    AccentPrimary   lipgloss.Color
    AccentSecondary lipgloss.Color
    // Semantics
    SemanticSuccess lipgloss.Color
    SemanticWarning lipgloss.Color
    SemanticError   lipgloss.Color
    SemanticInfo    lipgloss.Color
    SemanticOff     lipgloss.Color
    // Specials
    SpecialMuted     lipgloss.Color
    SpecialHighlight lipgloss.Color
    SpecialMatch     lipgloss.Color
    // Logo gradient (5 lines)
    LogoGradient [5]lipgloss.Color
}
```

- Two package-level instances: `ThemeTokyoNight` and `ThemeCyberpunk` — values taken verbatim from `tui-design-system-novo.md` §Paleta de Cores and §Gradiente do logo
- `rootModel` gets a `theme *Theme` field. Default at startup: `ThemeTokyoNight`.
- `rootModel` passes `theme` to every child at construction (same pattern as `actions *ActionManager` and `messages *MessageManager`). When theme toggles, `rootModel` updates `m.theme` and calls `SetTheme(t *Theme)` on all live children; modals in the stack also receive `SetTheme`.
- `childModel` interface **does NOT grow** a `SetTheme` method — avoid interface churn. Instead, `rootModel` calls a package-level helper `applyTheme(child childModel, t *Theme)` that does a type-switch over the known concrete types. Each concrete modal/child has its own `theme *Theme` field updated directly.

**F12 global shortcut — theme toggle:**
- Registered in `rootModel` via `ActionManager` at startup with **`HideFromBar: true`**:
  ```go
  Action{Key: "f12", Label: "Tema", Description: "Alternar tema (Tokyo Night ↔ Cyberpunk)", Group: "Global", Priority: 80, HideFromBar: true}
  ```
- `HideFromBar: true` means `ActionManager.Visible()` omits it from the command bar, but `ActionManager.All()` still includes it — so it will appear in the Help modal when that is implemented
- Intercepted in `rootModel.Update()` at the global shortcut step (D-06 priority 1), **before** flow or modal dispatch — toggles theme even while a modal is open
- On toggle: flip `m.theme` pointer between the two instances; call `applyTheme` on all live children and modals; `RenderLogo()` is called with the new theme on next `View()`; no Cmd needed (re-render is automatic)

**`RenderLogo()` becomes theme-aware:**
- Signature changes to `RenderLogo(t *Theme) string`
- Uses `t.LogoGradient[i]` for each of the 5 lines instead of hardcoded hex values
- All callers (`preVaultModel.View()`, any future screen showing the logo) pass the current theme

**Theme persistence:** not in Phase 6. Default is always Tokyo Night at startup. Persisting the user's choice to the vault settings file is Phase 9 scope.

### File Layout

**D-10: New files added in Phase 6**

New files in `internal/tui/`:
- `theme.go` — `Theme` struct, `ThemeTokyoNight` and `ThemeCyberpunk` instances, `applyTheme()` helper
- `header.go` — `headerModel` struct, `Render(width, vaultName, isDirty, activeArea)` method
- `filepicker.go` — `filePickerModal` struct + `FilePickerModeOpen`/`FilePickerModeSave` consts
- `passwordentry.go` — `passwordEntryModal` struct
- `passwordcreate.go` — `passwordCreateModal` struct

Modified files in `internal/tui/`:
- `flow_open_vault.go` — full implementation of `openVaultFlow` (replacing stub)
- `flow_create_vault.go` — full implementation of `createVaultFlow` (replacing stub)
- `prevault.go` — `View()` updated to render action hints beneath logo; receives `*Theme`
- `ascii.go` — `RenderLogo(t *Theme) string` signature update; uses `t.LogoGradient`
- `root.go` — `Init()` updated to trigger CLI fast-path; `theme *Theme` field added; `header *headerModel` field added; F12 global shortcut; `applyTheme` calls on toggle
- `state.go` — new message types from D-09 added
- `dialogs.go` — new factory helpers: `NewRecognitionError(title, text)`, `NewOverwriteConfirm(name, onConfirm, onCancel)`

### Agent's Discretion

- Exact Go struct field counts and constructor signatures for the 3 new modal types
- `filePickerModal` implementation strategy: `charm.land/bubbles/v2/filepicker` exists but is single-panel — researcher must determine if it can be wrapped to produce the two-panel spec layout or whether a custom `os.ReadDir`-based implementation is required
- Exact `storage.Probe()` API (or equivalent) needed for CLI fast-path validation — researcher to check if storage package exposes a header-only read; if not, `storage.Load` with a dummy password will always return `ErrAuthFailed` (which is sufficient to distinguish "recognizable vault" from "not a vault")
- Scroll implementation inside FilePicker panels (how many rows are visible vs. total)
- `pwdEnteredMsg.password` lifetime management — flow should zero the slice after storage.Load returns, regardless of success or failure
- Whether `openVaultFlow` and `createVaultFlow` own their modal structs as fields (for state continuity across Update calls) or recreate them on each flow state transition
- Whether to use `charm.land/bubbles/v2/key` key bindings inside the new modal types — existing codebase uses string matching directly; key.Binding adds self-documenting structure but is not required for correctness
- Exact F12 key constant in Bubble Tea v2 — researcher to verify (`tea.KeyF12`? or string `"f12"`?)

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### UI/UX Specification
- `tui-specification-novo.md` §Componentes/Cabeçalho — **full header spec**: 2-line anatomy, both structural states (no vault / vault open), all 4 mode-tab wireframes (Cofre, Modelos, Config, base), tab mechanics (line 1↔2 visual transformation), truncation algorithm, token table, event table
- `tui-specification-novo.md` §PasswordEntry — exact wireframes, token table, state table, message table, behavior rules for the vault-open password dialog
- `tui-specification-novo.md` §PasswordCreate — same for vault-create password dialog (two fields, Tab, strength meter)
- `tui-specification-novo.md` §FilePicker — both Open and Save mode wireframes, element token tables, state tables, message tables, keyboard navigation rules
- `tui-specification-novo.md` §Diálogos de Decisão — Recognition × Error and Confirmation × Destructive patterns used for error modals and overwrite confirmation
- `tui-design-system-novo.md` §Paleta de Cores — **all token values for both Tokyo Night and Cyberpunk themes**; hex values must be transcribed verbatim into `ThemeTokyoNight` and `ThemeCyberpunk`
- `tui-design-system-novo.md` §Gradiente do logo — 5-line logo gradient colors per theme; goes into `Theme.LogoGradient`
- `tui-design-system-novo.md` §Temas — F12 global toggle behaviour, Settings persistence (Phase 9)
- `tui-design-system-novo.md` §Mapa de Teclas — F12 is global scope (overrides all other scopes); `ActionManager` registration rules
- `tui-design-system-novo.md` §Sobreposição — modal anatomy and overlay rules
- `tui-design-system-novo.md` §Dimensionamento e Layout — dialog sizing rules (50-col for password dialogs, 70-col max for FilePicker)

### Behavior / Flow Specification
- `fluxos.md` §Fluxo 1 — Abrir Cofre Existente — full open vault flow including CLI fast-path, error branching, retry rules
- `fluxos.md` §Fluxo 2 — Criar Novo Cofre — full create vault flow including overwrite confirmation, weak password behavior, gravação

### Requirements
- `requisitos.md` §VAULT-01 — Criar cofre (dual confirmation, strength evaluation)
- `requisitos.md` §VAULT-03 — Abrir cofre (password, validation)
- `requisitos.md` §VAULT-04 — Error classification: auth errors allow retry; integrity errors block
- `requisitos.md` §Força da Senha Mestra — strength criteria (12+ chars, uppercase, lowercase, digit, special) and non-blocking warning rule

### Architecture (Phase 5 decisions — binding for Phase 6)
- `.planning/phases/05-tui-scaffold-root-model/05-CONTEXT.md` §D-02 — workArea state machine and modal flow architecture
- `.planning/phases/05-tui-scaffold-root-model/05-CONTEXT.md` §D-06 — message dispatch priority (flows intercept input before modals)
- `.planning/phases/05-tui-scaffold-root-model/05-CONTEXT.md` §D-09 — `preVaultModel` is display-only, no sub-states
- `.planning/phases/05-tui-scaffold-root-model/05-CONTEXT.md` §D-10 — modal stack mechanics (`pushModalMsg`/`popModalMsg`)
- `.planning/phases/05-tui-scaffold-root-model/05-CONTEXT.md` §D-15 — file layout conventions

### Component Research
- `.planning/research/tui-architecture.md` §14 — `bubbles/filepicker` exists (`fp.AllowedTypes`, `fp.CurrentDirectory`); single-panel — evaluate for wrapping vs. custom
- `.planning/research/tui-architecture.md` §6 — `textinput.EchoPassword` + `EchoCharacter` confirmed for masked password fields
- `.planning/research/tui-architecture.md` §2 — Bubble Tea v2 API: `tea.KeyPressMsg`, `msg.Code`, `msg.Text`, `msg.Mod.Contains(tea.ModCtrl)` for key handling in new modals
- `.planning/research/tui-architecture.md` §9 — `lipgloss.Place()` confirmed for modal overlay centering (already used in Phase 5)

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/tui/prevault.go` — `preVaultModel`, `renderHints()` already usable for welcome screen key hints
- `internal/tui/dialogs.go` — `NewMessage`, `NewConfirm` factories; extend here for recognition-error and overwrite-confirm factories
- `internal/tui/flow_open_vault.go` — stub `openVaultFlow` and `openVaultDescriptor` (key `"o"`, IsApplicable when !VaultOpen)
- `internal/tui/flow_create_vault.go` — stub `createVaultFlow` and `createVaultDescriptor` (key `"n"`, IsApplicable when !VaultOpen)
- `internal/tui/modal.go` — `modalModel`, `pushModalMsg`, `popModalMsg` push/pop helpers
- `internal/tui/state.go` — tick machinery, existing domain message types
- `internal/tui/root.go` — `rootModel.Init()` currently returns nil; needs CLI fast-path cmd here
- `internal/tui/ascii.go` — `RenderLogo()` will become `RenderLogo(t *Theme)` in Phase 6
- `charm.land/bubbles/v2/textinput` — `EchoMode`, `EchoCharacter` for masked password fields (confirmed in tui-architecture.md §6)
- `charm.land/bubbles/v2/filepicker` — single-panel file picker component; evaluate suitability for two-panel spec

### Established Patterns
- Charm v2 APIs: `View()` returns `string` (childModel), `tea.KeyPressMsg` for key events, `tea.View` only from `rootModel.View()`
- `textinput` component from `charm.land/bubbles/v2` — used for all text input fields
- `pushModalMsg` / `popModalMsg` — modal lifecycle; flows push by returning `pushModalMsg` cmd from `Update()`
- Pointer-receiver mutation in `childModel` — no self-replacement
- `MessageManager.Set(msg, duration)` — used for hints and timed errors in the message bar

### Integration Points
- `internal/crypto` — `EvaluatePasswordStrength([]byte) PasswordStrength` consumed by `passwordCreateModal` on each keystroke
- `internal/storage` — `RecoverOrphans(path)`, `Load(path, password) (*vault.Vault, FileMetadata, error)`, `SaveNew(path, vault, password)` consumed by flow cmd goroutines; error sentinels `ErrAuthFailed`, `ErrInvalidMagic`, `ErrVersionTooNew`, `ErrCorrupted`
- `internal/vault.Manager` — `Create(password []byte)` + `Load(*vault.Vault)` consumed by `createVaultFlow`
- `rootModel.vaultPath` — set from `vaultOpenedMsg.path`; persists across lock cycles

</code_context>

<specifics>
## Specific Ideas

- Wireframes in `tui-specification-novo.md` use rounded box characters (`╭╮╰╯│─`) — these are the required border style per design system
- FilePicker shows file metadata in format `● filename.abditum   25.8 MB 1h` — size is human-readable (KB/MB), date is relative (e.g., `1h`, `3d`)
- FilePicker `Caminho:` header row is read-only and updates as user navigates the tree; it shows the full absolute path of the current selected directory
- Strength meter format: `Força: ████████░░ Boa` — filled blocks (`█`) use `semantic.success` or `semantic.warning`, empty blocks (`░`) use `text.disabled`, label (`Boa`/`Forte`/`Fraca`) follows same semantic token
- Password masks use fixed 8 `•` characters regardless of actual password length — does not leak length
- Attempt counter format: `Tentativa 2 de 5` — hidden on first attempt, visible from second onward; uses `text.secondary`
- On vault open success, `rootModel` transitions to `workAreaVault` — `vaultTree` and `secretDetail` child models are allocated with the current `*Theme`, `preVault` is set to nil
- Theme toggle is visible immediately at any screen — changing theme while FilePicker or PasswordEntry modal is open updates that modal's colors on the next render cycle

</specifics>

<deferred>
## Deferred Ideas

- None — discussion stayed within phase scope

</deferred>

---

*Phase: 06-welcome-screen-vault-create-open*
*Context gathered: 2026-04-02*
