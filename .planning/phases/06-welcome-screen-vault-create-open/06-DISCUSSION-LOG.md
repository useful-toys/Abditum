# Phase 6: Welcome Screen + Vault Create/Open - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-04-06
**Phase:** 06-welcome-screen-vault-create-open
**Areas discussed:** Welcome Screen Interaction, File Picker User Experience, Password Dialogs Behavior, Vault Lifecycle Flow Control, Global Shortcuts and Theme System.

---

## Welcome Screen Interaction Pattern

| Option | Description | Selected |
|--------|-------------|----------|
| Dedicated actions in Command Bar | User initiates Open/Create via explicit keys (e.g., 'o', 'n') visible in the command bar. | ✓ |

**User's choice:** Initiated via explicit actions ('o', 'n') displayed as hints on the welcome screen, dispatched by `ActionManager`.

**Notes:** This aligns with the `tui-elm-architecture.md` where `welcomeModel` is display-only and `ActionManager` handles key dispatch. `welcomeModel` renders visual hints.

---

## File Picker User Experience

| Option | Description | Selected |
|--------|-------------|----------|
| Custom two-panel layout with detailed navigation | FilePicker implements a two-panel layout with directory tree and file list, supporting full keyboard/mouse navigation, filtering, metadata display, and error handling for inaccessible paths. | ✓ |

**User's choice:** Custom two-panel FilePicker implementation, adhering to `tui-specification-novo.md` for layout, navigation, metadata, and error handling (e.g., 0 vertical padding, `os.ReadDir`).

**Notes:** `charm.land/bubbles/v2/filepicker` was considered but deemed unsuitable due to its single-panel layout. Custom implementation ensures full compliance with the two-panel specification. Error messages are generic and user-friendly.

---

## Password Entry & Creation Flows

| Option | Description | Selected |
|--------|-------------|----------|
| Dedicated modals with masked input, strength meter, attempt counter | Separate modals (`passwordEntryModal`, `passwordCreateModal`) handle password input, masked with '•', provide strength feedback, and manage attempts. | ✓ |

**User's choice:** Dedicated `passwordEntryModal` and `passwordCreateModal` conforming to `tui-specification-novo.md` for fixed width, masking, real-time strength meter (for create), attempt counter (for entry), and precise message feedback.

**Notes:** `bubbles/textinput` is used for input fields. Password strength is evaluated using `crypto.EvaluatePasswordStrength`. Password buffers are zeroed immediately after use.

---

## Vault Lifecycle Flow Control

| Option | Description | Selected |
|--------|-------------|----------|
| State machine flows (`flowHandler`) with explicit `startFlowMsg`/`endFlowMsg` | `openVaultFlow` and `createVaultFlow` act as internal state machines, pushing dialogs onto the modal stack and reacting to results. Includes pre-checks for unsaved changes, overwrite confirmation, and error handling. | ✓ |

**User's choice:** Orchestrate multi-step vault operations using `flowHandler` implementations (`openVaultFlow`, `createVaultFlow`) that manage their internal state and interact with the `rootModel` via specific messages (`filePickedMsg`, `pwdEnteredMsg`, `pwdCreatedMsg`, `vaultOpenedMsg`, `flowCancelledMsg`).

**Notes:** Flows begin by checking for unsaved changes (requiring a `Confirmation × Neutro` dialog). Overwrite scenarios in create flow trigger a `Confirmation × Destrutivo` dialog. All `storage` and `vault` errors are classified and mapped to user-friendly messages and appropriate dialogs. `RecoverOrphans` is called silently.

---

## Global Shortcuts and Theme System

| Option | Description | Selected |
|--------|-------------|----------|
| `Ctrl+Q` for exit with unsaved changes flow, `F12` for theme toggle | Global shortcuts registered in `ActionManager`, `Ctrl+Q` triggers `fluxos.md` exit flows. `F12` toggles between `ThemeTokyoNight` and `ThemeCyberpunk` without persistence. | ✓ |

**User's choice:** `rootModel` manages global shortcuts. `Ctrl+Q` orchestrates the appropriate exit flow (`Fluxo 3, 4, or 5`). `F12` directly toggles the active `Theme` instance, propagating changes to all active TUI components immediately.

**Notes:** The theme system uses a `Theme` struct to encapsulate all color tokens, populated with exact hex values from `tui-design-system-novo.md`. Theme changes are applied via a type-switch helper (`applyTheme`) to all live children/modals. `F12` is registered as `HideFromBar: true`. Memory and screen clearing on exit/lock are ensured.

---

## the agent's Discretion

- Exact Go struct field counts and constructor signatures for the new modal types.
- Scroll implementation details within FilePicker panels.
- `pwdEnteredMsg.password` and `pwdCreatedMsg.password` lifetime management: flow handlers must zero the slices after use.
- Whether `openVaultFlow` and `createVaultFlow` own their modal structs as fields (for state continuity across `Update` calls) or recreate them on each flow state transition.
- Use `charm.land/bubbles/v2/key` for key bindings in new modal types if it enhances readability without architectural overhead (existing codebase uses string matching).
- `storage.Probe()` API (or equivalent) for CLI fast-path: if `internal/storage` does not expose a header-only read, `storage.Load` with a dummy password will return `ErrAuthFailed` for non-vault files, which is sufficient.

## Deferred Ideas

- **Vault tree display and secret navigation**: Phase 7.
- **Save / Save As / Discard / Change Password / Export / Import as standalone features**: Phase 9 (only implemented as part of initial Open/Create flows in Phase 6).
- **Security timers, clipboard, manual lock/exit flows**: Phase 10.
- **Theme persistence**: Saving the user's selected theme to settings file is Phase 9.
- **Tab navigation between modes (F2, F3, F4)**: These are visually rendered but become functional in Phase 7+.
- **File picker for `bubbles/filepicker`**: If the current custom implementation proves too complex or deviates too much, re-evaluate wrapping `bubbles/filepicker` for the two-panel spec layout in a future phase.
