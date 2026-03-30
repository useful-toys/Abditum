---
phase: 03-vault-domain-manager
plan: 06
subsystem: domain
tags: [vault, domain-model, tdd, secret-mutations, state-management]

# Dependency graph
requires:
  - phase: 03-vault-domain-manager
    provides: Secret lifecycle operations (create, delete, restore, duplicate, favorite)
provides:
  - Secret content mutation methods (rename, edit fields, edit observação)
  - Secret structural operations (move folder, reposition within folder)
  - EstadoSessao tracking for content vs structural changes (D-11, D-16)
  - Change detection preventing no-op mutations (D-12)
  - Observação architectural separation enforcement (D-29)
affects: [04-storage, 05-tui]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Content mutations set estadoSessao = Modificado"
    - "Structural operations (move, reposition) do NOT change estadoSessao"
    - "Change detection returns (alterado bool, err error) from entity methods"
    - "No-op operations at boundary positions (subir at 0, descer at last) return early"

key-files:
  created: []
  modified:
    - internal/vault/entities.go
    - internal/vault/manager.go
    - internal/vault/errors.go
    - internal/vault/validation_test.go

key-decisions:
  - "D-11: Content mutations (rename, edit) mark estadoSessao = Modificado"
  - "D-16: Structural operations (move, reposition) only update cofre.modificado, NOT estadoSessao"
  - "D-12: Change detection prevents no-op edits from marking modified"
  - "D-29: Observação is separate CampoSegredo field, architecturally excluded from campos slice"
  - "D-23: Repositioning to current position, Subir at position 0, Descer at last position are all no-ops"

patterns-established:
  - "TDD cycle: RED (failing test) → GREEN (minimal implementation) → REFACTOR (cleanup)"
  - "Entity methods return (alterado bool, err error) for change detection feedback"
  - "Manager uses feedback to conditionally update global state"
  - "Structural operations distinguished from content mutations via state tracking"

requirements-completed: []

# Metrics
duration: 60min
completed: 2026-03-30
---

# Phase 3 Plan 6: Secret Content Mutations Summary

**Implemented secret content mutations (rename, edit fields, edit observação) and structural operations (move, reposition) with TDD methodology, establishing clear distinction between content changes (estadoSessao = Modificado) and structural organization (cofre.modificado only)**

## Performance

- **Duration:** 60 min (1h)
- **Started:** 2026-03-30T09:14:21Z
- **Completed:** 2026-03-30T10:14:43Z
- **Tasks:** 6 (1 already complete + 5 new)
- **Files modified:** 4
- **Tests:** 90 tests pass (added 2 new comprehensive validation tests)

## Accomplishments
- Implemented complete secret content mutation API (rename, edit fields, edit observação)
- Implemented structural operations API (move folder, reposition within folder, subir/descer)
- Established architectural separation for Observação (D-29) via structural enforcement
- Validated estadoSessao tracking distinguishes content vs structural operations (D-16)
- Full TDD coverage with RED-GREEN cycles for all operations

## Task Commits

Each task was committed atomically using TDD methodology:

1. **Task 1: RenomearSegredo** - Already complete before session (commits: 9bf3d5e RED, d486454 GREEN)
2. **Task 2: EditarCampoSegredo** - `87c69fc` (test), `a556ac0` (feat)
   - RED: Added failing test for field editing with index validation
   - GREEN: Implemented EditarCampoSegredo with estadoSessao = Modificado
3. **Task 3: EditarObservacao** - `159c244` (test), `b9b1cd7` (feat)
   - RED: Added failing test for observação as separate field
   - GREEN: Implemented EditarObservacao with length validation (max 1000 chars)
4. **Task 4: MoverSegredo** - `abd0c6d` (test), `e93aeda` (feat)
   - RED: Added failing test for move without estadoSessao change
   - GREEN: Implemented MoverSegredo as structural operation per D-16
5. **Task 5: Repositioning methods** - `380fa2e` (test), `ded2368` (feat)
   - RED: Added failing test for ReposicionarSegredo, SubirSegredoNaPosicao, DescerSegredoNaPosicao
   - GREEN: Implemented repositioning with no-op detection for boundary positions
6. **Task 6: Comprehensive validation tests** - `daf4377` (test)
   - Added TestRenomearSegredoEstado: verifies rename marks estadoSessao = Modificado
   - Added TestObservacaoSeparada: verifies Observação architectural separation (D-29)

## Files Created/Modified
- `internal/vault/entities.go` - Added entity methods for secret mutations (renomear, editarCampo, editarObservacao, mover, reposicionar, obterPosicaoAtual, validation methods)
- `internal/vault/manager.go` - Added Manager methods: RenomearSegredo, EditarCampoSegredo, EditarObservacao, MoverSegredo, ReposicionarSegredo, SubirSegredoNaPosicao, DescerSegredoNaPosicao
- `internal/vault/errors.go` - Added ErrObservacaoMuitoLonga sentinel error for length validation
- `internal/vault/validation_test.go` - Added comprehensive tests covering all UAT criteria (SEC-03, SEC-04, SEC-05, SEC-08, SEC-09)

## Decisions Made

**Key architectural decisions enforced:**

1. **D-11 (Content vs Metadata)**: Content mutations (rename, edit fields, edit observação) mark `estadoSessao = Modificado`. Metadata changes (favorite toggle) do NOT change estadoSessao.

2. **D-16 (Structural vs Content)**: Structural operations (MoverSegredo, ReposicionarSegredo) only update `cofre.modificado`, NOT `estadoSessao`. This distinguishes organization from actual content changes.

3. **D-29 (Architectural Separation)**: Observação stored as separate `CampoSegredo` field (not in `campos` slice). `Campos()` getter excludes observação by design. This prevents manipulation by structural impossibility.

4. **D-12 (Change Detection)**: Entity methods return `(alterado bool, err error)`. Manager uses feedback to decide whether to update flags/timestamps. No-op operations (same value) don't mark modified.

5. **D-23 (No-op Boundary Cases)**: Repositioning to current position, Subir at position 0, Descer at last position are all no-ops (return early without marking modified).

## Deviations from Plan

None - plan executed exactly as written. TDD methodology followed for all tasks with RED-GREEN cycles. All specifications from 03-CONTEXT.md (D-11, D-12, D-16, D-29) were correctly implemented.

## Issues Encountered

None. All tasks completed successfully following TDD cycle. Tests passed on first GREEN implementation attempt for each task.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

**Ready for Phase 4 (Storage Layer):**
- All vault operations complete (CRUD, lifecycle, mutations, structural operations)
- Manager API stable and fully tested
- EstadoSessao tracking ready for persistence layer integration
- Validation logic comprehensive (90 tests covering all edge cases)

**Remaining work in Phase 3:**
- Additional plans may cover search, favorites, or other domain operations per ROADMAP.md

---
*Phase: 03-vault-domain-manager*
*Completed: 2026-03-30*

## Self-Check: PASSED

All claimed files exist:
- ✓ internal/vault/entities.go
- ✓ internal/vault/manager.go
- ✓ internal/vault/errors.go
- ✓ internal/vault/validation_test.go

All claimed commits exist:
- ✓ 87c69fc (RED Task 2: EditarCampoSegredo)
- ✓ a556ac0 (GREEN Task 2: EditarCampoSegredo)
- ✓ 159c244 (RED Task 3: EditarObservacao)
- ✓ b9b1cd7 (GREEN Task 3: EditarObservacao)
- ✓ abd0c6d (RED Task 4: MoverSegredo)
- ✓ e93aeda (GREEN Task 4: MoverSegredo)
- ✓ 380fa2e (RED Task 5: Repositioning)
- ✓ ded2368 (GREEN Task 5: Repositioning)
- ✓ daf4377 (Task 6: Comprehensive validation tests)
