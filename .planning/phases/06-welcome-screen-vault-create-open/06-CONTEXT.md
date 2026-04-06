# Phase 6: Welcome Screen + Vault Create/Open - Context

**Gathered:** 2026-04-06
**Status:** Ready for planning

<domain>
## Phase Boundary

This phase delivers the core user flows for vault setup: creating a new vault and opening an existing vault. It encompasses the initial welcome screen experience, file selection using a custom file picker, secure master password input (entry and creation), comprehensive error handling (authentication, integrity, external modification), and successful vault loading. All user-facing messages and dialogs adhere strictly to the TUI Design System and Specification.

This phase implements:
- **`headerModel`**: A new 2-line component rendering the app name, vault name (truncated with dirty indicator), and mode tabs (Cofre, Modelos, Config). It handles two structural states: no vault (welcome) and vault open.
- **`welcomeModel`**: The initial screen that renders the theme-aware ASCII art logo and provides action hints for "Open" and "Create" vault flows. It remains display-only.
- **Theme system**: The `Theme` struct, `ThemeTokyoNight` and `ThemeCyberpunk` instances (with all token hex values from `tui-design-system-novo.md`), and a mechanism for `rootModel` to toggle the active theme (`F12` shortcut) and propagate it to all live children/modals.
- **`filePickerModal`**: A new modal for selecting/saving vault files. Implements the two-panel layout (tree + file list), path navigation, `*.abditum` filtering, auto-selection, scroll, and error handling for inaccessible directories or existing files (in save mode).
- **`passwordEntryModal`**: A new modal for master password input when opening a vault. Features masked input, an attempt counter, and specific error messages for incorrect passwords.
- **`passwordCreateModal`**: A new modal for master password creation (and confirmation) when creating a new vault or changing it. Includes two masked input fields, Tab navigation, and a real-time password strength meter.
- **`openVaultFlow`**: Full implementation of the "Open Vault" flow, including checking for unsaved changes, file selection, password entry, external modification detection, error classification, and loading the vault. It supports CLI path fast-path.
- **`createVaultFlow`**: Full implementation of the "Create Vault" flow, including checking for unsaved changes, file path selection, overwrite confirmation, password creation, strength evaluation, and saving the new vault.
- **Error classification and dialogs**: Mapping of `internal/storage` and `internal/vault` errors to generic user-friendly messages and appropriate decision dialogs (Recognition × Error, Confirmation × Destructive).
- **Global shortcuts**: `F12` for theme toggle, `Ctrl+Q` for quit (with unsaved changes flow).

This phase does NOT implement:
- Vault tree display or secret navigation (Phase 7).
- Save / Save As / Discard / Change Password / Export / Import operations as standalone features (Phase 9). These are only implemented as part of the initial Open/Create flows as required by `fluxos.md`.
- Security timers, clipboard, manual lock/exit flows as independent actions (Phase 10).
- Any new `workArea` transitions beyond `workAreaPreVault → workAreaVault` (and back on error).
- Theme persistence to disk (Phase 9). Phase 6 defaults to Tokyo Night; F12 toggles in-memory only.
- Tab navigation between mode tabs (the tab strip is rendered but mode-switching F-keys are Phase 7+ scope).

</domain>

<decisions>
## Implementation Decisions

### Header Component (`header.go`)

- **D-01: `headerModel` — 2-line component with two structural states**
    - The header is a standalone struct (not a `childModel`) with a `Render(width int, vaultName string, isDirty bool, activeArea workArea, theme *Theme) string` method.
    - **Inputs at render time:** `width`, `vaultName` (radical only, no `.abditum`), `isDirty`, `activeArea` (determines active tab), `theme`.
    - **State 1 — No vault (welcome):** App name (`accent.primary` bold) on line 1; full-width `─` separator (`border.default`) on line 2. No vault name, dirty indicator, or tabs.
    - **State 2 — Vault open:**
        - Line 1: `  Abditum · {vaultName} {•}`. `•` (semantic.warning) for `isDirty`. `·` separator (`border.default`).
        - Tabs (Cofre, Modelos, Config): inactive tabs (`text.secondary`) show `╭ Text ╮`; active tab (`accent.primary` bold on `special.highlight`) shows `╭──────╮` on line 1, and `╯ Text ╰` on line 2, visually connecting to the work area.
        - Tabs use `border.default` for borders (`╭╮╯╰─`).
    - **Vault name truncation:** Truncated with `…` from the right if it exceeds available space, prioritizing tabs. Formula: `available = terminal_width - prefix_width - dirty_indicator_width - tabs_width - min_padding`. If `name + …` doesn't fit, only `…` is shown.
    - **Phase 6 tabs:** Only `workAreaVault` is functionally reachable. `workAreaTemplates` and `workAreaSettings` render in inactive style but are non-functional until Phase 7+.

### Welcome Screen (`welcome.go`)

- **D-02: `welcomeModel` remains display-only; action hints integrated**
    - `welcomeModel.Update()` handles no keyboard input; keys fall through to `ActionManager` dispatch.
    - `welcomeModel.View()` renders the theme-aware logo and action hints: `n Novo cofre    o Abrir cofre` using `text.secondary` color.
    - `ctrl+q` quit flow is a global shortcut. No `q` shortcut on welcome screen.

### File Picker Modal (`filepicker.go`)

- **D-03: `filePickerModal` implements two-panel design with full navigation and feedback**
    - **Modality:** A new modal struct (`filePickerModal`) satisfying `modalView`, not using `charm.land/bubbles/v2/filepicker` directly due to layout mismatch (custom implementation using `os.ReadDir`).
    - **Layout:** Max DS dialog width (`min(70 cols, 80% of terminal)`), height up to `80% of terminal`. Two panels: `Estrutura` (~40%) and `Arquivos` (~60%). **0 vertical padding** (exception to DS rule, justified by vertical space scarcity).
    - **Borders:** `border.focused` for the main modal border. Internal separators and junctions (`├┬┴┤`) use `border.default`. Scroll indicators (`↑`/`↓`/`■`) on panel borders use `text.secondary`.
    - **Initial directory:** CWD. If inaccessible, fallback to user home (`~`), with `MsgWarn` message.
    - **Filtering:** Only `*.abditum` files shown. Hidden files (starting with `.`) are excluded. `.abditum` extension omitted from display names.
    - **File metadata:** Displays size (human-readable KB/MB) and relative date/time (e.g., `1h`, `3d`).
    - **Navigation:**
        - `Tab` cycles focus: Tree → Files → Filename input (Save mode only).
        - Keyboard navigation (`↑↓←→`, `PgUp`/`PgDn`, `Home`/`End`) for both panels.
        - `Enter` expands/collapses folders in tree, selects file in file list, or confirms input in filename field.
        - Mouse scroll affects the focused panel.
    - **Errors:** Inaccessible directories trigger `MsgError` in message bar; folder remains collapsed.
    - **Mode Save specifics:** Includes a filename input field. Selecting an existing file in the list copies its name (without extension) to the input field. Overwrite confirmation is handled by the calling flow, not the modal.
    - Emits `filePickedMsg{path string}` on confirm; `flowCancelledMsg{}` on ESC.
    - **Sensitive data:** File paths are treated as sensitive if they are part of a `pwdEnteredMsg` or `pwdCreatedMsg` for `modalResult` routing.

### Password Input Modals (`passwordentry.go`, `passwordcreate.go`)

- **D-04: `passwordEntryModal` — secure single-field password input**
    - **Modality:** New modal struct satisfying `modalView`.
    - **Layout:** Fixed width 50 columns. `border.focused` border.
    - **Input:** Single `textinput` (masked `••••••••`, `EchoCharacter = '•'`). Mask is fixed 8 `•` characters, not reflecting actual length.
    - **Attempt counter:** `Tentativa N de 5` (`text.secondary`), hidden on attempt 1, visible from attempt 2 onward (max 5 attempts).
    - **Actions:** `Enter Confirmar` is `text.disabled` when field is empty; `accent.primary` bold when non-empty. `Esc Cancelar` always active.
    - **Behavior:** On wrong password: field cleared, counter incremented, default action locked. On 5th attempt exhausted: dialog closes automatically (emits `flowCancelledMsg{}`).
    - **Messages (MessageManager):**
        - On open / field empty or valid: `• Digite a senha para desbloquear o cofre` (Hint).
        - On wrong password: `✕ Senha incorreta` (Error, 5s TTL).
    - Emits `pwdEnteredMsg{password []byte}` on confirm (buffer immediately zeroed after extraction). `Tab` does nothing.

- **D-05: `passwordCreateModal` — two-field password creation with strength meter**
    - **Modality:** New modal struct satisfying `modalView`.
    - **Layout:** Fixed width 50 columns. `border.focused` border.
    - **Input:** Two `textinput` fields (`Nova senha`, `Confirmação`), both masked (`••••••••`).
    - **Navigation:** `Tab` cycles focus between the two fields.
    - **Strength meter:** `Força: ████████░░ Boa` (or `Fraca`). Hidden when `Nova senha` field is empty. `██` uses `semantic.success`/`semantic.warning`; `░░` uses `text.disabled`. Label (`Boa`/`Forte`/`Fraca`) uses the same semantic token. Calls `crypto.EvaluatePasswordStrength` on every keystroke in `Nova senha`.
    - **Actions:** `Enter Confirmar` is `text.disabled` if either field is empty or passwords diverge; `accent.primary` bold when both non-empty and matching. `Esc Cancelar` always active. Weak password does NOT block submission (per `requisitos.md`).
    - **Behavior:** On passwords diverge: confirmation field cleared, focus returns to confirmation, default action locked, `MsgError` ("As senhas não conferem").
    - **Messages (MessageManager):**
        - Focus on `Nova senha` (empty/valid): `• A senha mestra protege todo o cofre — use 12+ caracteres` (Hint).
        - Focus on `Confirmação` (empty/valid): `• Redigite a senha para confirmar` (Hint).
        - On divergence: `✕ As senhas não conferem — digite novamente` (Error, 5s TTL).
    - Emits `pwdCreatedMsg{password []byte}` on confirm (buffers immediately zeroed after extraction).

### Vault Lifecycle Flows (`flow_open_vault.go`, `flow_create_vault.go`)

- **D-06: `openVaultFlow` orchestrates open sequence with pre-check for unsaved changes**
    - **Flow state machine:** `stateCheckDirty` → `statePickFile` (FilePicker) → `statePwdEntry` (PasswordEntry) → `statePreload` (RecoverOrphans + Load) → `stateDone`.
    - **Unsaved changes check (`stateCheckDirty`):** Before `FilePicker`, if `vault.Manager.IsDirty()`, push a `Confirmation × Neutro` dialog ("Alterações não salvas: Deseja salvar antes de sair?", with options: Salvar / Descartar / Voltar).
        - "Salvar": Trigger `saveVaultFlow` (Phase 9) or equivalent direct save. If save fails, flow cancels.
        - "Descartar": Proceed.
        - "Voltar": Flow cancels.
    - **Error classification (`D-09`):** `storage.Load` errors mapped to generic user messages and `Recognition × Error` modals. `storage.ErrAuthFailed` triggers retry in `passwordEntryModal`. Other errors lead to `statePickFile`.
    - Emits `vaultOpenedMsg{path}` on success.

- **D-07: `createVaultFlow` orchestrates create sequence with pre-check for unsaved changes**
    - **Flow state machine:** `stateCheckDirty` → `statePickFile` (FilePicker-Save) → `stateCheckOverwrite` → `statePwdCreate` (PasswordCreate) → `stateSaveNew` (Create + SaveNew) → `stateDone`.
    - **Unsaved changes check (`stateCheckDirty`):** Same logic as `openVaultFlow`.
    - **Overwrite confirmation (`stateCheckOverwrite`):** If target file exists after FilePicker, push `Confirmation × Destrutivo` dialog ("Sobrescrever arquivo?" with options: Sobrescrever / Voltar).
        - "Sobrescrever": Proceed.
        - "Voltar": Return to `statePickFile`.
    - **Weak password:** `passwordCreateModal` allows weak password submit. `createVaultFlow` then pushes `Confirmation × Alerta` dialog ("Senha fraca: Deseja prosseguir mesmo assim?" with options: Prosseguir / Revisar) if `crypto.EvaluatePasswordStrength` is weak.
        - "Prosseguir": Proceed.
        - "Revisar": Return to `statePwdCreate`.
    - Emits `vaultOpenedMsg{path}` on success.

- **D-08: CLI path fast-path via `rootModel.Init()` for `openVaultFlow`**
    - `rootModel.Init()`: If `m.initialPath` is non-empty, immediately dispatch `openVaultDescriptor.New(ctx)` as a `tea.Cmd`.
    - `openVaultFlow` with pre-filled path starts in `statePreVerify`. It calls `storage.Probe(path)` (or equivalent header-only read) to validate `*.abditum` magic/version.
        - If valid: push `passwordEntryModal` directly (skip `FilePicker`).
        - If invalid: show `Recognition × Error` modal; on dismiss, transition to `statePickFile` (open `FilePicker` normally).
    - This applies only to `openVaultFlow`.

- **D-09: `RecoverOrphans` called silently before `storage.Load` in `openVaultFlow`**
    - In the `statePreload` step of `openVaultFlow`, `storage.RecoverOrphans(path)` is called via a `tea.Cmd` goroutine before `storage.Load`.
    - `RecoverOrphans` errors are ignored silently (no user message, no logs). `storage.Load` proceeds regardless.

### Error Classification and Dialogs

- **D-10: Standardized error messages and dialogs for vault operations**
    - All errors during `openVaultFlow` or `createVaultFlow` are caught and mapped to generic user-friendly messages. No Go error strings or internal details exposed.
    - **`storage.ErrAuthFailed`**: Handled internally by `passwordEntryModal` as `✗ Senha incorreta` (Error, 5s TTL) in the MessageManager. Allows retry.
    - **`storage.ErrInvalidMagic`**: `O arquivo selecionado não é um cofre Abditum`. `Recognition × Error` modal.
    - **`storage.ErrVersionTooNew`**: `Este cofre foi criado por uma versão mais recente do Abditum`. `Recognition × Error` modal.
    - **`storage.ErrCorrupted`**: `O arquivo está corrompido e não pode ser aberto`. `Recognition × Error` modal.
    - **Save/Create failure (general)**: `Não foi possível [criar/salvar] o cofre — verifique o caminho e as permissões`. `Recognition × Error` modal.
    - All `Recognition × Error` dialogs use `✗` symbol, `semantic.error` border, `Enter OK` action.
    - `dialogs.go` needs new factory helpers for these `Recognition × Error` modals: `NewRecognitionError(title, text string) tea.Cmd`.

### Theme System (`theme.go`)

- **D-11: `Theme` struct holds all design system token values; `F12` toggles**
    - `Theme` struct: Contains `lipgloss.Color` fields for all tokens defined in `tui-design-system-novo.md` (§Paleta de Cores, §Gradiente do logo).
    - Two package-level instances: `ThemeTokyoNight` and `ThemeCyberpunk`, populated with exact hex values.
    - `rootModel` has a `theme *Theme` field (default `ThemeTokyoNight`).
    - **Propagation:** `rootModel` passes `theme` to children/modals at construction. When `rootModel.theme` changes, `applyTheme(child childModel, t *Theme)` is called on all live children and modals (type-switch based).
    - **`F12` global shortcut:** Registered in `ActionManager` as `ScopeGlobal` with `HideFromBar: true`. Intercepted in `rootModel.Update()`, toggles `m.theme` pointer, triggers `applyTheme` calls.
    - `RenderLogo(t *Theme) string` uses `t.LogoGradient` for colors.
    - Theme persistence to disk is Phase 9.

### Flow Control and Messaging

- **D-12: New domain message types for Phase 6 flows**
    - `filePickedMsg{path string}`: Absolute path chosen in FilePicker.
    - `pwdEnteredMsg{password []byte}`: Raw password bytes from PasswordEntry (zeroed in modal immediately after copy).
    - `pwdCreatedMsg{password []byte}`: Confirmed new password from PasswordCreate.
    - `vaultOpenedMsg{path string}`: Signals successful vault open/create to rootModel (triggers `workAreaVault` transition).
    - `flowCancelledMsg{}`: Signals active flow cancellation (sets `activeFlow = nil`).
    - These messages are handled by the orchestrating flow handler's `Update()` method, or by `rootModel.Update()` for global effects (`vaultOpenedMsg`, `flowCancelledMsg`).

- **D-13: `rootModel.Update()` handles `startFlowMsg`/`endFlowMsg` and `modalResult` routing**
    - Confirms existing `startFlowMsg` sets `activeFlow` and calls `flow.Init()`.
    - Confirms `endFlowMsg` sets `activeFlow = nil`.
    - `modalResult` messages (e.g., `pwdEnteredMsg`, `filePickedMsg`) are routed **exclusively** to `activeFlow.Update()` when a flow is active, not broadcast.

### Global Shortcuts and Exit Flows

- **D-14: `Ctrl+Q` global quit shortcut implements `fluxos.md` exit flows**
    - Intercepted in `rootModel.Update()` as a global action.
    - **Flow 3 (Quit without vault):** Simple `Confirmation × Neutro` dialog.
    - **Flow 4 (Quit with clean vault):** Simple `Confirmation × Neutro` dialog.
    - **Flow 5 (Quit with dirty vault):**
        - First, `Confirmation × Neutro` dialog ("Alterações não salvas: Salvar / Descartar / Voltar?").
        - If "Salvar", then `rootModel` orchestrates external modification check (via `storage.DetectExternalChange`) and potentially `Confirmation × Destrutivo` dialog for overwrite.
        - Finally, calls `storage.Save` (via Cmd goroutine).
        - Errors during save (failure or backup available) result in `Recognition × Error` dialog; app remains open.
    - Memory wiping and screen clearing (`\033[3J`) for all exit paths (per `requisitos.md` §VAULT-13).

### Agent's Discretion

- Exact Go struct field counts and constructor signatures for the new modal types.
- Scroll implementation details within FilePicker panels.
- `pwdEnteredMsg.password` and `pwdCreatedMsg.password` lifetime management: flow handlers must zero the slices after use.
- Whether `openVaultFlow` and `createVaultFlow` own their modal structs as fields (for state continuity across `Update` calls) or recreate them on each flow state transition.
- Use `charm.land/bubbles/v2/key` for key bindings in new modal types if it enhances readability without architectural overhead (existing codebase uses string matching).
- `storage.Probe()` API (or equivalent) for CLI fast-path: if `internal/storage` does not expose a header-only read, `storage.Load` with a dummy password will return `ErrAuthFailed` for non-vault files, which is sufficient.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### UI/UX Specification
- `tui-design-system-novo.md` — **All visual foundations:** principles, color palette (exact hex values for Tokyo Night/Cyberpunk themes, logo gradient), typography, borders, sizing/layout (min terminal, vertical zones, dialog sizing, scroll in dialogs), icons/symbols, visual states, overlay patterns (modals, dialogs, action bar rules, semantic severity), messages (types, symbols, TTLs), focus/navigation, keymap (global keys, F12 toggle).
- `tui-specification-novo.md` — **All screen/component specifics:**
    - §Componentes/Cabeçalho — Full header spec (2-line anatomy, states, tabs, truncation, tokens, events).
    - §PasswordEntry — Exact wireframes, token table, state table, message table, behavior rules for vault-open password dialog.
    - §PasswordCreate — Exact wireframes, token table, state table, message table, behavior rules for vault-create password dialog (two fields, Tab, strength meter).
    - §FilePicker — Both Open and Save mode wireframes, element token tables, state tables, message tables, keyboard navigation rules, 0 vertical padding exception.
    - §Diálogos de Decisão — Recognition × Error and Confirmation × Destructive patterns for error modals and overwrite confirmation.
    - §Telas/Boas-vindas — Minimal layout, ASCII art logo, action hints.

### Behavior / Flow Specification
- `fluxos.md` — **All relevant flows:**
    - §Entrada via Linha de Comando — CLI path fast-path for Open/Create.
    - §Fluxo 1 — Abrir Cofre Existente — Full open vault flow: check dirty, file selection, password entry, error branching, retry rules, loading.
    - §Fluxo 2 — Criar Novo Cofre — Full create vault flow: check dirty, file path selection, overwrite confirmation, password creation, strength evaluation, saving.
    - §Fluxo 3, 4, 5 — Sair da Aplicação — All quit flows (no vault, clean vault, dirty vault with external modification check).
    - §Fluxo 6 — Bloquear cofre — Discard unsaved changes silently, re-initiate Flow 1 from password.
    - §Fluxo 7 — Aviso de bloqueio iminente por inatividade — Warning before auto-lock.

### Requirements
- `requisitos.md` — **Specific requirements driving Phase 6 features:**
    - §Ciclo de Vida do Cofre (VAULT-01, VAULT-03, VAULT-04, VAULT-05, VAULT-07, VAULT-08, VAULT-09, VAULT-10, VAULT-11, VAULT-12, VAULT-13, VAULT-14) — Create, Open, Save, Discard, Change Password, Lock, Exit flows, external modification, error classification, memory/screen wipe.
    - §Força da Senha Mestra (PWD-01) — Strength criteria and non-blocking warning.
    - §Como protegemos seus dados (Security Principles, `mlock`/`VirtualLock`, zero logs).

### Architecture (Phase 5 decisions — binding for Phase 6)
- `.planning/phases/05.1-tui-scaffold-root-model-fix/05.1-CONTEXT.md` — TUI interfaces and dispatch rules:
    - §D-01: `childModel` (Update, View, SetSize)
    - §D-02: `modalView` (Update, View, Shortcuts)
    - §D-03: `modalResult` (for sensitive data messages)
    - §D-04: `flowHandler` (Init, Update)
    - §D-05: `Action` struct (Keys, Label, Group, Scope, Priority, HideFromBar, Enabled, Handler)
    - §D-06: `ActionManager` API (Register, ClearOwned, SetActiveOwner, Dispatch, Visible, All, RenderCommandBar)
    - §D-07: `MessageManager` API (Show, Clear, Current, Tick, HandleInput) with `MsgKind` and TTL.
    - §D-08: `startFlowMsg`/`endFlowMsg` for flow lifecycle.
    - §D-09: Dispatch order (Global actions → active modal → active flow → active child).
    - §D-10: `dialogs.go` factories (Message, Confirm, PasswordEntry, PasswordCreate, FilePicker).
- `.planning/phases/05-tui-scaffold-root-model/05-CONTEXT.md` — Broader TUI architecture:
    - §D-02: `workArea` state machine and modal flow architecture.
    - §D-04: Concrete pointer fields for child models; `nil` = inactive.
    - §D-10: Modal stack mechanics (`pushModalMsg`/`popModalMsg`).
    - §D-15: File layout conventions for `internal/tui/`.

### Component Research
- `.planning/research/tui-architecture.md`
    - §2: Bubble Tea v2 API (`tea.KeyPressMsg`, `msg.Code`, `msg.Text`, `msg.Mod.Contains(tea.ModCtrl)` for key handling).
    - §6: `textinput.EchoPassword` + `EchoCharacter` confirmed for masked password fields.
    - §9: `lipgloss.Place()` confirmed for modal overlay centering.
    - §14: `bubbles/filepicker` is single-panel — must be wrapped or custom implemented for two-panel spec.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/vault.Manager` — All vault operations. TUI must exclusively use this.
- `internal/vault.Cofre`, `Pasta`, `Segredo` — Domain entities with exported getters.
- `internal/crypto.EvaluatePasswordStrength([]byte) PasswordStrength` — Consumed by `passwordCreateModal` for real-time feedback.
- `internal/storage.RecoverOrphans(path string)`, `storage.Load(path, password) (*vault.Vault, FileMetadata, error)`, `storage.SaveNew(path, vault, password)` — Consumed by flow command goroutines.
- `internal/storage` error sentinels: `ErrAuthFailed`, `ErrInvalidMagic`, `ErrVersionTooNew`, `ErrCorrupted`.
- `internal/tui/ascii.go` — `RenderLogo()` will be updated to accept `*Theme`.
- `internal/tui/dialogs.go` — `NewMessage`, `NewConfirm` factories (extend for new error/overwrite dialogs).
- `internal/tui/flow_open_vault.go` — Stub `openVaultFlow` and `openVaultDescriptor` (key `"o"`, `IsApplicable` when `!VaultOpen`).
- `internal/tui/flow_create_vault.go` — Stub `createVaultFlow` and `createVaultDescriptor` (key `"n"`, `IsApplicable` when `!VaultOpen`).
- `internal/tui/modal.go` — `modalModel`, `pushModalMsg`, `popModalMsg` push/pop helpers.
- `internal/tui/messages.go` — `MessageManager`, `RenderMessageBar`.
- `internal/tui/actions.go` — `ActionManager`, `RenderCommandBar`.
- `internal/tui/prevault.go` (to be renamed to `welcome.go`) — `welcomeModel`, `renderHints()` helper.
- `internal/tui/root.go` — `rootModel.Init()` currently returns nil; needs CLI fast-path Cmd.
- `charm.land/bubbles/v2/textinput` — `EchoMode`, `EchoCharacter` for masked password fields.
- `charm.land/bubbles/v2/viewport` — Potentially useful for scrollable dialog content (FilePicker, Help).
- `charm.land/lipgloss/v2` — Extensive styling capabilities.

### Established Patterns
- Charm v2 APIs: `View()` returns `string` (childModel), `tea.KeyPressMsg` for key events, `tea.View` only from `rootModel.View()`.
- `textinput` component from `charm.land/bubbles/v2` for all text input fields.
- `pushModalMsg` / `popModalMsg` for modal lifecycle; flows push by returning `pushModalMsg` cmd from `Update()`.
- Pointer-receiver mutation in `childModel` (no self-replacement).
- `MessageManager.Show(msg, duration)` for hints and timed errors in the message bar.
- `CGO_ENABLED=0` enforced globally.
- `time.Now().UTC()` for timestamps.

### Integration Points
- `cmd/abditum/main.go` — Will be updated to bootstrap the TUI (`tea.NewProgram`), passing `*vault.Manager` and `initialPath` to `rootModel`.
- `internal/vault.NewManager()` — Instantiated in `main.go`, passed to `rootModel`.
- `rootModel.vaultPath` — Set from `vaultOpenedMsg.path`; persists across lock cycles.
- `internal/tui` package — Will be populated with all new components and flows.

</code_context>

<specifics>
## Specific Ideas

- **Header `headerModel`:** Implement the `Render` method to compose the two-line header with dynamic vault name truncation (radical only, no `.abditum`), dirty indicator (`•`), and the three mode tabs (`Cofre`, `Modelos`, `Config`), ensuring the active tab visually connects to the work area below.
- **FilePicker (`filePickerModal`):** Implement custom two-panel layout. Use `os.ReadDir` for directory traversal. Display file sizes as "KB", "MB", etc., and dates as "1h", "3d", "2y" or "DD/MM/YY". Ensure keyboard navigation (`↑↓←→`, `PgUp`/`PgDn`, `Home`/`End`) and mouse scroll work for both panels. Truncate long paths in the header with `…`.
- **Password Entry/Create (`passwordEntryModal`, `passwordCreateModal`):** Use `bubbles/textinput` for masked fields. Implement the 8-character fixed `•` mask. Password strength meter `████████░░ Boa` with correct token colors and dynamic label. Attempt counter (`Tentativa N de 5`) hidden on first attempt.
- **Theme toggle (`F12`):** `rootModel` needs to manage the active `*Theme` pointer and call `applyTheme` on all children/modals when `F12` is pressed.
- **CLI fast-path for Open Vault:** `rootModel.Init()` will call `storage.Probe(path)` (if available, or a `storage.Load` dummy call to check magic/version) before pushing `passwordEntryModal`.
- **Error messages:** Ensure all user-facing error messages are generic and consistent with `fluxos.md` and `tui-design-system-novo.md` (e.g., `✗ Senha incorreta`, `O arquivo selecionado não é um cofre Abditum`).

</specifics>

<deferred>
## Deferred Ideas

- **Vault tree display and secret navigation**: Phase 7.
- **Save / Save As / Discard / Change Password / Export / Import as standalone features**: Phase 9 (only implemented as part of initial Open/Create flows in Phase 6).
- **Security timers, clipboard, manual lock/exit flows**: Phase 10.
- **Theme persistence**: Saving the user's selected theme to settings file is Phase 9.
- **Tab navigation between modes (F2, F3, F4)**: These are visually rendered but become functional in Phase 7+.
- **File picker for `bubbles/filepicker`**: If the current custom implementation proves too complex or deviates too much, re-evaluate wrapping `bubbles/filepicker` for the two-panel spec layout in a future phase.

</deferred>

---

*Phase: 06-welcome-screen-vault-create-open*
*Context gathered: 2026-04-06*