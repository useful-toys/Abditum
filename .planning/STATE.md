---
gsd_state_version: 1.0
milestone: v1.0
milestone_name: milestone
current_phase: 05
status: pending
last_updated: "2026-03-30T16:30:00.000Z"
progress:
  total_phases: 11
  completed_phases: 4
  total_plans: 15
  completed_plans: 15
---

# Project State — Abditum

**Last updated:** 2026-03-30T16:30:00Z
**Current phase:** 05
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

**Phase 04: Storage Package** — Complete ✅

All 4 plans executed successfully:

- ✓ 04-01: Foundation (AAD-aware crypto, vault serialization, format constants)
- ✓ 04-02: Core I/O (SaveNew, Save, Load, atomic rename, backup chain)
- ✓ 04-03: Recovery & Migration (RecoverOrphans, DetectExternalChange, Migrate scaffold)
- ✓ 04-04: Integration (FileRepository adapter, end-to-end tests)

**Package Status:**

- 27 tests in `internal/storage`, all passing
- Full binary `.abditum` format with 49-byte AAD header
- Atomic writes via `.tmp` → rename, `.bak`/`.bak2` backup chain
- Windows-specific `MoveFileEx` via build tag
- `FileRepository` implements `vault.RepositorioCofre`
- Ready for Phase 5 (TUI Scaffold + Root Model)

**Next:** Phase 05 - TUI Scaffold + Root Model

## Accumulated Context

### Roadmap Evolution

- Phase 04.1 inserted after Phase 04: Refinamento da camada de domínio — encapsulamento e versioning (INSERTED)

### Phase 04.1 — Itens identificados em revisão de código (2026-03-31)

1. **`DeserializarCofre` receber `version uint8`** — passar versão do formato para o deserializador; substituir cadeia de migração JSON→JSON (`migrate.go`) por compat fields nas structs de serialização
2. **`segredo.marcarModificacao()`** — método privado que seta `dataUltimaModificacao`; elimina acesso direto a campo interno espalhado em 5 métodos do Manager
3. **`cofre.marcarModificado()`** — método privado que seta `modificado = true` e `dataUltimaModificacao`; elimina acesso direto a campos internos do Cofre no Manager
4. **Factory define `estadoSessao` inicial** — `pasta.criarSegredo` e `pasta.duplicarSegredo` devem retornar segredo já com estado correto; Manager não deve setar `estadoSessao` após factory
5. **`AlternarFavoritoSegredo` não deve atualizar `segredo.dataUltimaModificacao`** — favoritar não é edição de conteúdo; bug atual identificado

---

## Phase History

### Phase 04: Storage Package (Completed 2026-03-30)

**Status:** Complete ✅

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

### Phase 02 Decisions

1. **Phase 02-01:** Memory locking failures are non-fatal (D-03) — `SecureAllocate()` returns usable buffer even if mlock fails
2. **Phase 02-01:** Nonce generation is internal to `Encrypt()` function (D-19) — callers never handle nonces directly
3. **Phase 02-01:** Password strength evaluation operates on `[]byte` without string conversion (Pitfall 3) — prevents unzeroable copies
4. **Phase 02-01:** Platform-specific mlock/munlock kept internal to package — exposed only through `SecureAllocate()` API

## Open Decisions

- Help overlay de teclado (`?` key / footer hints) — especificar antes da Phase 5 (TUI scaffold)

## Notes

- Bubble Tea v2 import path: `charm.land/bubbletea/v2` (não v1)
- Argon2id: m=256 MiB, t=3, p=4 (conforme formato-arquivo-abditum.md)
- Lixeira: lista in-memory no Manager, sem campo no Segredo
- Exibição dentro de pasta: subpastas primeiro, segredos depois
- NanoID: implementar internamente com `crypto/rand`
- Datetime: RFC 3339 UTC
- Favoritos: DFS na ordem do JSON
