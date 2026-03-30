# Plan 04-04 Summary

**Phase:** 04-storage-package
**Plan:** 04 — Repository Adapter + Integration Tests
**Status:** COMPLETE
**Completed:** 2026-03-30

## What Was Done

Connected the storage layer to the domain layer via `FileRepository`, which implements `vault.RepositorioCofre`. Added integration tests proving the full end-to-end pipeline.

### Task 1: FileRepository Adapter (internal/storage/repository.go)

**`type FileRepository struct`** — holds `path`, `password`, `salt`, `isNew bool`, `metadata FileMetadata`

**Constructors:**
- `NewFileRepository(path, password, salt, metadata)` — for existing vaults (after Load)
- `NewFileRepositoryForCreate(path, password)` — for new vaults (isNew = true)

**Methods:**
- `Salvar(cofre *vault.Cofre) error` — first call uses `SaveNew` (direct write); subsequent calls use `Save` (atomic .tmp + rename). Updates internal `metadata` and `salt` after each successful save.
- `Carregar() (*vault.Cofre, error)` — delegates to `Load`, updates internal `metadata` and `salt`.
- `Metadata() FileMetadata` — returns last snapshot for use with `DetectExternalChange`.
- `Path() string` — returns vault file path.
- `UpdatePassword(password []byte)` — replaces stored password for master password change flow.

**Private helper:** `readSaltFromFile(path) ([]byte, error)` — reads `data[SaltOffset:SaltOffset+SaltSize]` from the file header, returns a copy.

**Compile-time interface check:**
```go
var _ vault.RepositorioCofre = (*FileRepository)(nil)
```

The `vault.RepositorioCofre` interface is kept as-is (no changes needed — `Salvar` and `Carregar` signatures match exactly).

### Task 2: Integration Tests (storage_test.go) — 5 new tests, all pass

| Test | Verifies |
|------|----------|
| `TestIntegration_FullPipelineRoundtrip` | NovoCofre + InicializarConteudoPadrao → Salvar → Carregar → verify 2 folders, 3 templates |
| `TestIntegration_BackupChainRotation` | 3 consecutive saves → .bak created on 2nd, .bak2 created on 3rd |
| `TestIntegration_ExternalChangeDetection` | Save → DetectExternalChange false → append byte → DetectExternalChange true |
| `TestIntegration_ErrorClassification` | 4 subtests: WrongMagic, VersionTooNew, WrongPassword, TamperedHeader |
| `TestIntegration_ManagerWithFileRepository` | Manager.CriarSegredo + Manager.Salvar → Load directly → secret present |

## Verification Results

```
CGO_ENABLED=0 go test ./internal/storage/... -count=1 -v  → 27/27 PASS (7.1s)
CGO_ENABLED=0 go test ./...                               → all packages PASS
CGO_ENABLED=0 go vet ./...                                → OK
CGO_ENABLED=0 go build ./cmd/abditum                     → OK
```

## Key Decisions

- `salt` is stored in `FileRepository` and reused across saves. It is extracted from the file header after each `SaveNew` or `Carregar`. This keeps the password-derived key stable — the key only changes when `UpdatePassword` is called and a new salt is generated at the next `SaveNew`.
- `NewFileRepository` accepts `nil` for `salt` — `Carregar` reads the salt from the file anyway, so the caller does not need to pre-extract it.
- The `RepositorioCofre` interface is intentionally minimal (`Salvar` + `Carregar` only). FileMetadata for change detection is surfaced via `repo.Metadata()` and `DetectExternalChange`, not through the interface — avoiding a cross-package dependency cycle.
- `Manager.Salvar()` calls `prepararSnapshot()` (deep copy filtering `EstadoExcluido`) then `repositorio.Salvar(snapshot)` — `FileRepository.Salvar` receives the snapshot, not the live vault.
