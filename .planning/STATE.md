---
gsd_state_version: 1.0
milestone: v1.0
milestone_name: milestone
current_phase: 06
status: planning
last_updated: "2026-04-06T12:50:34.168Z"
progress:
  total_phases: 22
  completed_phases: 13
  total_plans: 51
  completed_plans: 48
---

# Project State — Abditum

**Last updated:** 2026-04-01T03:03:16Z
**Current phase:** 06
**Milestone:** v1.0

## Status

| Artifact | Status |
|----------|--------|
| PROJECT.md | ✓ Created |
| config.json | ✓ Created |
| REQUIREMENTS.md | ✓ 55 requirements defined, traceability mapped |
| research/ | ✓ STACK, FEATURES, ARCHITECTURE, PITFALLS, SUMMARY |
| ROADMAP.md | ✓ 11 phases, all requirements mapped |
| modelo-dominio.md | ✓ Revisado — ordem corrigida (subpastas antes de segredos) |
| arquitetura.md | ✓ Revisado — versões Bubble Tea v2 explicitadas |
| formato-arquivo-abditum.md | ✓ Revisado — Argon2id m=256 MiB incorporado ao REQUIREMENTS |
| descricao.md | ✓ Revisado — ordem corrigida em fluxo "Visualizar hierarquia" |

## Current Phase

**Phase 05.6: TUI Scaffold Dialogs** — Planning ✓

Plans to execute:

- [ ] 05.6-01: Fix compilation errors in decision.go (unused titleStyle, missing BorderTitle) + rewrite View() with manual top-border construction
- [ ] 05.6-02: 10 matrix fixture constructors + 3 rendering tests (symbol presence, border chars, matrix render)
- [ ] 05.6-03: 4 interaction tests (Enter/Esc/arrows/unknown) + 4 edge case tests
- [ ] 05.6-04: 1 integration test (Acknowledge → modal stack push/Enter → pop) + final build+vet verification

**Pre-conditions (current state):**

- `go test ./internal/tui/...` FAILS (build error in decision.go)
- Error 1: `titleStyle` declared and not used (line 179)
- Error 2: `borderStyle.BorderTitle` undefined (lipgloss v2 has no such method)
- All other packages compile and test clean

**Next:** Execute Plan 01 (fix compilation errors)

## Accumulated Context

### Roadmap Evolution

- Phase 05.7 inserted after Phase 05: golden-test-architecture-for-tui-modals (INSERTED URGENT)
  - **Purpose:** Golden test architecture for TUI modals — detect visual regressions (spacing, colors, fonts) that conventional tests miss
  - **Scope:** SGR parser, golden test runners, 48 DecisionDialog golden files, 8 HelpDialog golden files

- Phase 12 added: 05 tui-scaffold-action-arch
- Phase 05.6 added: tui-scaffold-dialogs

- Phase 05.2.2 inserted after Phase 05.2: tui-scaffold-message-arch-fixes (INSERTED REFINEMENT)
  - **Purpose:** Test message truncation with >100 char long messages + F10 action in PoC
  - **Scope:** Expand poc-mensagens to validate RenderMessageBar behavior at terminal limits

- Phase 05.2.1 inserted after Phase 05.2: tui-scaffold-message-arch-fixes (INSERTED URGENT)

- Phase 05.2 inserted after Phase 5: tui-scaffold-message-arch (INSERTED)

- Phase 04.1 inserted after Phase 04: Refinamento da camada de domínio — encapsulamento e versioning (INSERTED)

### Phase 04.1 — Itens identificados em revisão de código (2026-03-31)

1. **`DeserializarCofre` receber `version uint8`** — passar versão do formato para o deserializador; substituir cadeia de migração JSON→JSON (`migrate.go`) por compat fields nas structs de serialização
2. **`segredo.marcarModificacao()`** — método privado que seta `dataUltimaModificacao`; elimina acesso direto a campo interno espalhado em 5 métodos do Manager
3. **`cofre.marcarModificado()`** — método privado que seta `modificado = true` e `dataUltimaModificacao`; elimina acesso direto a campos internos do Cofre no Manager
4. **Factory define `estadoSessao` inicial** — `pasta.criarSegredo` e `pasta.duplicarSegredo` devem retornar segredo já com estado correto; Manager não deve setar `estadoSessao` após factory
5. **`AlternarFavoritoSegredo` não deve atualizar `segredo.dataUltimaModificacao`** — favoritar não é edição de conteúdo; bug atual identificado

---

## Phase History

### Phase 04.1: Refinamento da Camada de Domínio (Completed 2026-03-31)

**Status:** Ready to plan

**Key Deliverables:**

- `segredo.marcarModificacao()`, `cofre.marcarModificado()` — encapsulated state mutation
- `Pasta.copiarProfundo()`, `Segredo.copiar()`, `CampoSegredo.copiar()` — entity-owned deep copy
- `segredo.zerarValoresSensiveis()` — encapsulated crypto.Wipe
- Factory `criarSegredo`/`duplicarSegredo` correctly set `EstadoIncluido` from birth
- `DeserializarCofre(data, version uint8)` — versioned deserializer ready for compat fields
- `migrate.go` deleted — compatibility via compat fields in JSON structs
- Manager fully refactored: 20+ direct field mutations replaced with entity method calls
- `AlternarFavoritoSegredo` bug fixed: does NOT update secret timestamp (favoriting ≠ content edit)
- Removed `copiarPastaRecursivamente`, `copiarSegredo`, `copiarCampo` from Manager
- 3 regression tests: D-08, D-05 (create), D-05 (duplicate)

**Commits:**

- 6856183 — feat(vault): add private entity methods for encapsulation (04.1-01)
- 10bdac2 — docs(04.1-01): add execution summary
- a9477b2 — feat(storage): versioned DeserializarCofre, remove migrate.go (04.1-02)
- c022fe3 — docs(04.1-02): add execution summary
- b864754 — refactor(vault): manager uses entity methods, fix AlternarFavorito bug (04.1-03)
- db3c511 — docs(04.1-03): add execution summary

### Phase 04: Storage Package (Completed 2026-03-30)

**Status:** Executing Phase 04.1

**Key Deliverables:**

- `internal/storage` package: binary `.abditum` format, atomic saves, backup chain
- 49-byte file header as GCM AAD: magic(4) + version(1) + salt(32) + nonce(12)
- `SaveNew`, `Save`, `Load` with `.tmp` → rename atomic protocol
- Windows: `MoveFileEx` with `MOVEFILE_REPLACE_EXISTING` (build-tagged)
- `RecoverOrphans`: cleans stale `.tmp` files on startup
- `DetectExternalChange`: compares mtime+size to detect external mutations
- `Migrate` scaffold with `MigrationFunc` registry for future format versions
- `FileRepository` implementing `vault.RepositorioCofre` — bridge to vault domain
- 27 tests total, all passing

**Commits:**

- a191024 — feat(04-01): storage foundation (format constants, errors)
- 781481f — feat(04-01): SerializarCofre / DeserializarCofre roundtrip
- 49ca52d — feat(04-01): EncryptWithAAD / DecryptWithAAD
- fca0253 — feat(04-02): SaveNew, Save, Load, atomic rename, backup rotation
- 08a3e93 — feat(04-03): RecoverOrphans, DetectExternalChange, Migrate scaffold
- 8b0644b — feat(04-04): FileRepository adapter with integration tests

### Phase 03: Vault Domain + Manager (Context Complete 2026-03-29)

**Status:** Executing Phase 04

**Context Artifacts:**

- ✓ 03-CONTEXT.md — 30 implementation decisions across 4 gray areas
- ✓ 03-DISCUSSION-LOG.md — Full audit trail of architectural discussions

**Key Architectural Decisions:**

- D-01: No synthetic IDs (Go pointers as identifiers)
- D-08/D-09: Package-level encapsulation with exported getters
- D-11: Two independent state flags (cofre.modificado vs segredo.estadoSessao)
- D-04/D-25: Manager as thin orchestrator, entities own validation
- D-17: Atomic save with two-phase commit pattern
- D-24: Timestamps only on structural changes
- D-27: Segredo soft delete, Pasta hard delete
- D-28: Factory + Initializer separation

**Requirements Mapped:** VAULT-02, SEC-05, FOLDER-01 through FOLDER-05, TPL-01 through TPL-06

### Phase 02: Crypto Package (Completed 2026-03-29)

**Plans:** 1/1 complete

- ✓ 02-01-PLAN.md — Cryptographic primitives package (7 tasks, 28 tests, 15 files)

**Key Deliverables:**

- `internal/crypto` package with production-ready cryptographic primitives
- Argon2id key derivation with secure parameters (m=256 MiB, t=3, p=4)
- AES-256-GCM authenticated encryption with automatic nonce generation
- Platform-specific memory locking (Unix/Windows/fallback)
- Password strength evaluation (12+ chars, 4 categories)
- Comprehensive test coverage with TDD methodology

**Commits:**

- 2cb9774 — Package documentation and sentinel errors
- 00d2e87 — Argon2id key derivation
- 84b7d83 — AES-256-GCM encryption
- f43d2f7 — Memory security primitives
- 8700a7c — Platform-specific memory locking
- 702fd00 — Password strength evaluation
- 8e07f1b — Integration tests

### Phase 01: Project Scaffold + CI Foundation (Completed 2026-03-29)

**Plans:** 3/3 complete

- ✓ 01-01-PLAN.md — Go module initialization and directory structure
- ✓ 01-02-PLAN.md — CI workflow and Makefile configuration
- ✓ 01-03-PLAN.md — golangci-lint security configuration

## Decisions Made

### Phase 03 Context Decisions

1. **D-01:** No synthetic IDs — Go pointers sufficient for in-memory identity
2. **D-08:** Package-level encapsulation — all entity fields private to `internal/vault`
3. **D-09:** Safe pointer sharing — `Manager.Vault()` returns live `*Cofre` pointer, getters return defensive copies
4. **D-11:** Two independent state flags — `cofre.modificado` (any mutation) vs `segredo.estadoSessao` (content only)
5. **D-17:** Atomic save with two-phase commit — prepare snapshot, persist, finalize deletions only on success
6. **D-24:** Timestamps on structural changes only — favoriting doesn't update `segredo.dataUltimaModificacao`
7. **D-27:** Deletion semantics differ — Segredo soft delete, Pasta hard delete
8. **D-28:** Factory vs Initializer — `NovoCofre()` creates structure, `InicializarConteudoPadrao()` bootstraps content

See `.planning/phases/03-vault-domain-manager/03-CONTEXT.md` for complete list of 30 decisions.

- [Phase 03-01]: Pointer identity (no synthetic IDs) - D-01
- [Phase 03-01]: Package-private encapsulation with defensive copies - D-08/D-09
- [Phase 03-01]: Factory+bootstrap separation (NovoCofre/InicializarConteudoPadrao) - D-28
- [Phase 03]: D-12 change detection: RenomearModelo only marks modified if name actually changes
- [Phase 03-vault-domain-manager]: Moving folder into itself returns ErrDestinoInvalido before cycle check
- [Phase 03-vault-domain-manager]: Repositioning no-ops: to current position, Subir at 0, Descer at last (D-23)
- [Phase 03-vault-domain-manager]: Folder deletion: hard delete with promotion, recursive merge for subfolders, numeric suffix for secrets
- [Phase 03-05]: Favorito flag independent of estadoSessao (only updates cofre.modificado per D-11)
- [Phase 03-05]: Soft-delete: EstadoIncluido secrets removed from parent list, others marked Excluido
- [Phase 03-05]: Name conflict uses fmt.Sprintf for (N) progression with 9999 safety limit
- [Phase 03-06]: Content mutations (rename, edit fields, edit observação) mark estadoSessao = Modificado per D-11
- [Phase 03-06]: Structural operations (move, reposition) only update cofre.modificado, NOT estadoSessao per D-16
- [Phase 03-06]: Observação architecturally separated (CampoSegredo field, excluded from campos slice) per D-29
- [Phase 03-06]: Change detection returns (alterado bool, err error) - no-op edits don't mark modified per D-12
- [Phase 05]: childModel does NOT implement tea.Model — only rootModel does; View() returns string not tea.View
- [Phase 05]: Extra stub files created beyond plan's files_modified to keep package buildable and all 5 deps in go.mod
- [Phase 05]: FocusedField uses *vault.CampoSegredo — plan CONTEXT.md D-20 had wrong type name *vault.Campo
- [Phase 05-02]: ActionManager uses insertion sort for priority (small slice, no sort import needed)
- [Phase 05-02]: MessageManager simplified to single string slot — severity tiers deferred per D-17
- [Phase 05-02]: modalModel fully interactive (j/k navigation, enter/esc) not a passive content container
- [Phase 05-02]: NewMessage/NewConfirm naming (not Message/Confirm) per plan spec
- [Phase 05-02]: popModalMsg defined in modal.go alongside the type that emits it
- [Phase 05-03]: preVaultModel constructor takes *ActionManager (not zero-arg) for consistency and forward compatibility
- [Phase 05-03]: Work area stubs take (mgr *vault.Manager, actions *ActionManager, msgs *MessageManager) — matches Plan 04 rootModel call sites
- [Phase 05-03]: helpModal.buildContent() groups actions by Group field in insertion order from ActionManager.All()
- [Phase 05-03]: renderHints() helper placed in prevault.go alongside the preVaultModel
- [Phase 05]: modals field changed to []childModel to support heterogeneous modal types (modalModel + helpModal) without typed-nil trap
- [Phase 05-05]: NewRootModel exported wrapper added to root.go — main.go is package main and cannot access unexported newRootModel
- [Phase 05-05]: vault.NovoCofre() used in main.go bootstrap (not NewCofre which doesn't exist) — Portuguese naming convention
- [Phase 05-05]: CGO_ENABLED=0 go test -race fails on Windows — race detector requires CGO; use CI (Linux) for race detection; local tests verified without -race
- [Phase 05.3]: tokens.go created with centralized semantic colors, symbols, and style helpers — messages.go, actions.go, help.go refactored to consume tokens (zero hardcoding)
- [Phase 05.3]: NewRootModel() takes no parameters (PoC mode), all 15 actions registered, Init() starts global tick, header simplified to Abditum only
- [Phase 05.6]: decision.go View() must build the top border manually — lipgloss v2 has no BorderTitle method; formula: "╭── " + styledTitle + dashes(fill) + " ──╮"
- [Phase 05.6]: wrapBody() uses stdlib only (strings.Split + rune counting); no external wrapping dependency
- [Phase 05.6]: Test assertions on View() use plain-text content checks only — ANSI color code assertions are fragile and deferred to Phase 06+ visual integration tests
- [Phase 05.6]: View() builds top border manually (╭── styledTitle dashes ──╮) — lipgloss v2 has no BorderTitle method
- [Phase 05.6]: Plain-text test assertions on View() — ANSI escape code assertions are fragile and deferred to Phase 06+
- [Phase 05.7-golden-test-architecture-for-tui-modals]: RenderCommandBar receives []Action as parameter — callers pass am.Visible() explicitly, keeping the function pure and testable without an ActionManager
- [Phase 05.7-golden-test-architecture-for-tui-modals]: helpModal.groupLabel stored as func(int) string field — preserves GroupLabel fallback logic without ActionManager dependency, enabling plain []Action fixture injection in tests
- [Phase 05.7]: RenderCommandBar sorts actions by Priority desc internally — callers may pass unsorted slices
- [Phase 05.7]: helpGoldenPath uses variant-only (no width param) — variant already encodes size+scroll position
- [Phase 05.7]: Negative scroll sentinel (-1 → maxScroll) simplifies fixture setup for bottom-scroll scenarios in help_test.go
- [Phase 05.7-04]: Variant-embedded width in decision golden path helper — decisionGoldenPath(variant, ext) does not append width suffix since variant already encodes it (e.g. destructive-1action-short-30)
- [Phase 05.7-04]: Explicit-key Update tests use lowercase rune codes (e.g. 'm') to verify case-insensitive dispatch in DecisionDialog.Update()
- [Phase 05.7-06]: StyleTransition.MarshalJSON emits compact [line,col,fg,bg,[styles]] tuple — encoding/json import added; nil Style coerces to []string{} for null-safe JSON
- [Phase 05.7-06]: styleMapToArray uses active := []string{} literal (not var nil) so no-style transitions always produce [] in JSON output
- [Phase 05.7-06]: ParseANSIStyle deduplicates same-(line,col) transitions by replacing last slice element — multiple SGR before any text advance produce one transition

### Phase 02 Decisions

1. **Phase 02-01:** Memory locking failures are non-fatal (D-03) — `SecureAllocate()` returns usable buffer even if mlock fails
2. **Phase 02-01:** Nonce generation is internal to `Encrypt()` function (D-19) — callers never handle nonces directly
3. **Phase 02-01:** Password strength evaluation operates on `[]byte` without string conversion (Pitfall 3) — prevents unzeroable copies
4. **Phase 02-01:** Platform-specific mlock/munlock kept internal to package — exposed only through `SecureAllocate()` API

## Open Decisions

- (None — help overlay implemented in 05-03 via helpModal reading ActionManager.All())

## Notes

- Bubble Tea v2 import path: `charm.land/bubbletea/v2` (não v1)
- Argon2id: m=256 MiB, t=3, p=4 (conforme formato-arquivo-abditum.md)
- Lixeira: lista in-memory no Manager, sem campo no Segredo
- Exibição dentro de pasta: subpastas primeiro, segredos depois
- NanoID: implementar internamente com `crypto/rand`
- Datetime: RFC 3339 UTC
- Favoritos: DFS na ordem do JSON
