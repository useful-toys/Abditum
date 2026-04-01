---
phase: 05-tui-scaffold-root-model
plan: "01"
subsystem: ui

tags: [bubbletea-v2, lipgloss-v2, bubbles-v2, tui, go, interfaces]

requires:
  - phase: 03-vault-domain-manager
    provides: "Pasta, Segredo, CampoSegredo, ModeloSegredo entity types used in FlowContext"
  - phase: 04-storage
    provides: "vault.Manager with Vault() accessor used by rootModel to populate FlowContext"

provides:
  - "childModel interface (5-method contract all child TUI models must satisfy)"
  - "FlowContext struct (8 fields capturing navigation + vault state at flow dispatch time)"
  - "flowDescriptor, flowHandler interfaces + FlowRegistry for flow routing"
  - "chainFlowMsg for flow chaining after state transitions"
  - "workArea enum with 4 constants (preVault, vault, templates, settings)"
  - "tickMsg and 10 domain message types for broadcasting vault events"
  - "pushModalMsg for rootModel modal stack protocol"
  - "4 Cmd factory stubs in mutations.go (TODO phase-8)"
  - "All charm.land v2 dependencies in go.mod"
  - "Supporting stub models (modal, prevault, vaulttree, secretdetail, templatelist, templatedetail)"

affects:
  - 05-02-tui-scaffold-prevault
  - 05-03-tui-scaffold-vaulttree
  - 05-04-tui-scaffold-rootmodel
  - 05-05-tui-scaffold-root-model
  - all subsequent tui plans

tech-stack:
  added:
    - "charm.land/bubbletea/v2 v2.0.2 (direct)"
    - "charm.land/bubbles/v2 v2.0.0 (indirect)"
    - "charm.land/lipgloss/v2 v2.0.2 (indirect)"
    - "github.com/atotto/clipboard v0.1.4 (indirect)"
    - "github.com/charmbracelet/x/exp/teatest/v2 (indirect)"
  patterns:
    - "childModel does NOT implement tea.Model — only rootModel does"
    - "Update uses pointer receivers and mutates in place — no self-replacement"
    - "Domain messages broadcast to all live models; children ignore irrelevant types via default case"
    - "Flows push modals onto rootModel stack (pushModalMsg) — they have no View of their own"
    - "chainFlowMsg dispatches next flow after state transition completes"

key-files:
  created:
    - "internal/tui/flows.go — childModel, FlowContext, flowDescriptor, flowHandler, FlowRegistry, chainFlowMsg"
    - "internal/tui/state.go — workArea enum, tickMsg, 10 domain message types, pushModalMsg"
    - "internal/tui/mutations.go — 4 Cmd factory stubs"
    - "internal/tui/modal.go — modalModel implementing childModel"
    - "internal/tui/ascii.go — AsciiArt const + RenderLogo()"
    - "internal/tui/actions.go — ActionManager + Action struct (imports bubbles/v2)"
    - "internal/tui/messages.go — MessageManager + MessageSeverity (imports lipgloss/v2)"
    - "internal/tui/dialogs.go — Message() and Confirm() dialog factories"
    - "internal/tui/prevault.go — preVaultModel stub"
    - "internal/tui/vaulttree.go — vaultTreeModel stub"
    - "internal/tui/secretdetail.go — secretDetailModel stub"
    - "internal/tui/templatelist.go — templateListModel stub"
    - "internal/tui/templatedetail.go — templateDetailModel stub"
  modified:
    - "go.mod — added 5 charm.land/atotto/teatest dependencies"
    - "go.sum — updated with new dependency checksums"
    - "internal/tui/doc.go — updated package comment with full architecture description"

key-decisions:
  - "childModel does NOT implement tea.Model — only rootModel does; View() returns string not tea.View"
  - "bubbles/v2 and lipgloss/v2 kept as indirect deps until concrete models import them in later plans"
  - "Supporting stub files created beyond plan's listed files_modified to keep package buildable and all 5 deps in go.mod"
  - "FocusedField uses *vault.CampoSegredo (correct type); plan context D-20 had wrong type name"

patterns-established:
  - "Interface-first: type contracts defined in plan 01 so all subsequent plans implement against them"
  - "childModel pattern: pointer receiver Update, string View, no self-replacement in Elm loop"
  - "Domain message broadcast: rootModel fans out to all active children, children ignore unknowns"

requirements-completed: []

duration: ~45min
completed: 2026-04-01
---

# Phase 5 Plan 01: Core TUI Contracts Summary

**charm.land v2 dependency stack added to go.mod and all foundational TUI type contracts defined — childModel, FlowContext, flowDescriptor/flowHandler, FlowRegistry, workArea enum, 10 domain messages, and Cmd factory stubs**

## Performance

- **Duration:** ~45 min
- **Started:** 2026-04-01T00:00:00Z
- **Completed:** 2026-04-01T00:45:00Z
- **Tasks:** 2/2
- **Files modified:** 16 (2 modified, 14 created)

## Accomplishments
- Added all 5 required charm.land v2 dependencies to go.mod (bubbletea/v2, bubbles/v2, lipgloss/v2, clipboard, teatest/v2)
- Defined `childModel` interface with exactly 5 methods — the contract all child TUI models must satisfy
- Defined `FlowContext` struct with all 8 fields, `flowDescriptor`/`flowHandler` interfaces, `FlowRegistry`, `chainFlowMsg`
- Defined `workArea` enum (4 constants) and all 10 domain message types in state.go
- Created 4 Cmd factory stubs in mutations.go with `// TODO(phase-8):` markers
- Package builds and passes `go vet` with zero errors

## Task Commits

Each task was committed atomically:

1. **Task 1: Add charm.land v2 dependencies to go.mod** - `6deb61c` (chore)
2. **Task 2: Define core TUI type contracts** - `1016ca6` (feat)

**Plan metadata:** _(docs commit pending)_

## Files Created/Modified

- `go.mod` — Added 5 charm.land/atotto/teatest indirect and direct dependencies
- `go.sum` — Auto-updated with new dependency checksums
- `internal/tui/doc.go` — Updated package comment with full architecture description
- `internal/tui/flows.go` — childModel (5 methods), FlowContext (8 fields), flowDescriptor, flowHandler, FlowRegistry, chainFlowMsg
- `internal/tui/state.go` — workArea (4 consts), tickMsg, secretAddedMsg, secretDeletedMsg, secretRestoredMsg, secretModifiedMsg, secretMovedMsg, secretReorderedMsg, folderStructureChangedMsg, vaultSavedMsg, vaultReloadedMsg, vaultClosedMsg, vaultChangedMsg, pushModalMsg
- `internal/tui/mutations.go` — favoriteSecretCmd, softDeleteSecretCmd, restoreSecretCmd, reorderSecretCmd (all TODO phase-8 stubs)
- `internal/tui/modal.go` — modalModel implementing childModel
- `internal/tui/ascii.go` — AsciiArt const + RenderLogo() using lipgloss/v2
- `internal/tui/actions.go` — ActionManager + Action struct (imports bubbles/v2 textinput)
- `internal/tui/messages.go` — MessageManager + MessageSeverity (imports lipgloss/v2)
- `internal/tui/dialogs.go` — Message() and Confirm() factories, confirmModalModel
- `internal/tui/prevault.go` — preVaultModel stub (implements childModel)
- `internal/tui/vaulttree.go` — vaultTreeModel stub (implements childModel)
- `internal/tui/secretdetail.go` — secretDetailModel stub (implements childModel)
- `internal/tui/templatelist.go` — templateListModel stub (implements childModel)
- `internal/tui/templatedetail.go` — templateDetailModel stub (implements childModel)

## Decisions Made

- **childModel does NOT implement tea.Model** — only rootModel does; `View()` returns `string` not `tea.View`. This keeps child models from being accidentally passed to `tea.NewProgram`.
- **Extra stub files created** — The plan listed only `flows.go`, `state.go`, `mutations.go` in `files_modified`, but to keep the package buildable and all 5 deps in go.mod, 10 additional files were created (modal, ascii, actions, messages, dialogs, and 5 child-model stubs). This is necessary for `go mod tidy` not to drop unused deps.
- **`FocusedField *vault.CampoSegredo`** — plan CONTEXT.md D-20 mentioned `*vault.Campo` which doesn't exist; the correct type is `*vault.CampoSegredo` from entities.go.
- **bubbles/v2 and lipgloss/v2 remain indirect** — They are imported by actions.go and ascii.go/messages.go respectively, but `go mod tidy` keeps them as indirect because the module doesn't directly expose them in the `require` block as direct. Concrete child models in later plans will import them directly, promoting them.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 2 - Missing Critical] Created additional stub files beyond plan's listed files_modified**
- **Found during:** Task 2 (Define core TUI type contracts)
- **Issue:** Plan only listed flows.go, state.go, mutations.go, doc.go — but `go mod tidy` would drop bubbles/v2 and lipgloss/v2 if nothing imported them, and `modal.go` was referenced by pushModalMsg in state.go
- **Fix:** Created modal.go (referenced by state.go's pushModalMsg), ascii.go/actions.go/messages.go/dialogs.go (to keep all 5 deps in go.mod), and 5 child model stubs (prevault, vaulttree, secretdetail, templatelist, templatedetail) to scaffold the full package
- **Files modified:** 10 additional files created
- **Verification:** `go build ./internal/tui/...` and `go vet ./internal/tui/...` both pass
- **Committed in:** `1016ca6` (Task 2 commit)

---

**Total deviations:** 1 auto-fixed (Rule 2 — missing critical supporting files)
**Impact on plan:** Necessary for package correctness (modal.go referenced by pushModalMsg) and dependency management (keep all 5 deps in go.mod). No scope creep — all files are stubs implementing childModel.

## Issues Encountered

- `CGO_ENABLED=0` environment variable syntax differs on PowerShell — must use `$env:CGO_ENABLED="0"` instead of the Unix prefix syntax
- `go mod tidy` removes unused dependencies, so stub files importing bubbles/v2 and lipgloss/v2 were necessary to keep them in go.mod

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- All type contracts defined — Plans 02–05 can implement against childModel, FlowContext, workArea, and domain messages without discovering signatures
- All 5 charm.land/v2 dependencies available
- Supporting stub models ready to be fleshed out in subsequent plans
- No blockers

---
*Phase: 05-tui-scaffold-root-model*
*Completed: 2026-04-01*
