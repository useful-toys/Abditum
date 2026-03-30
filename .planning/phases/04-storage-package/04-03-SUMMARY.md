# Plan 04-03 Summary

**Phase:** 04-storage-package
**Plan:** 03 — Orphan Recovery, Change Detection, Migration Scaffold
**Status:** COMPLETE
**Completed:** 2026-03-30

## What Was Done

Added three supporting lifecycle functions to the storage package: startup orphan recovery, external change detection, and a migration scaffold for future format upgrades.

### Task 1: RecoverOrphans + DetectExternalChange

**`internal/storage/recover.go`** — `RecoverOrphans(vaultPath string) error`
- Removes `vaultPath + ".tmp"` if it exists (stale from an interrupted Save)
- No-op and returns nil when the vault is in a clean state
- Does NOT auto-restore from `.bak` — that decision is left to the user to avoid silent data loss
- Called once at application startup before any Load or Save

**`internal/storage/detect.go`** — two exported functions:
- `DetectExternalChange(vaultPath string, metadata FileMetadata) (bool, error)`:
  - Fast path: `os.Stat` size comparison (O(1), no file read)
  - Slow path: full `sha256.Sum256` comparison when sizes match
  - Returns `(true, nil)` on change, `(false, nil)` on no change, `(false, error)` on I/O failure
- `ComputeFileMetadata(vaultPath string) (FileMetadata, error)`:
  - Convenience function to snapshot current file state without a full Load

### Task 2: Migration Scaffold

**`internal/storage/migrate.go`** — `Migrate(data []byte, fromVersion, toVersion uint8) ([]byte, error)`
- Applies a chain of registered `MigrationFunc` values from `fromVersion` to `toVersion`
- Returns data unchanged for same-version (no-op)
- Returns error for downgrade (`fromVersion > toVersion`)
- Returns error for missing migration step (version gap)
- `var migrations = map[uint8]MigrationFunc{}` — currently empty; future migrations registered here

### New Tests (added to storage_test.go) — 10 new tests, all pass

| Test | Verifies |
|------|----------|
| `TestRecoverOrphans_RemovesStaleTmp` | .tmp removed on RecoverOrphans call |
| `TestRecoverOrphans_NoOpWhenClean` | nil returned when no .tmp exists |
| `TestDetectExternalChange_NoChange` | false for unchanged file |
| `TestDetectExternalChange_SizeDiffers` | true when bytes appended |
| `TestDetectExternalChange_ContentDiffers` | true for same-size in-place modification |
| `TestDetectExternalChange_FileNotFound` | error for missing file |
| `TestMigrate_SameVersion_NoOp` | v1→v1 returns same bytes |
| `TestMigrate_FutureVersion_Error` | v1→v2 returns error (no path) |
| `TestMigrate_Downgrade_Error` | v2→v1 returns error |
| `TestMigrate_V1FixtureRoundtrip` | full SaveNew + Load pipeline with v1 format |

## Verification Results

```
CGO_ENABLED=0 go test ./internal/storage/... -count=1 -v  → 22/22 PASS
CGO_ENABLED=0 go test ./...                               → all packages PASS
CGO_ENABLED=0 go vet ./...                                → OK
```

## Key Decisions

- `RecoverOrphans` is intentionally conservative: it only removes `.tmp` (always safe to delete) and does not auto-restore from `.bak` (risky without knowing why the file is unreadable).
- `DetectExternalChange` uses size + SHA-256 instead of mtime per `arquitetura.md §7` — mtime is unreliable on cloud-synced filesystems.
- Migration scaffold uses an empty map for v1 (no migrations needed into the baseline format). The chain loop handles multi-version upgrades automatically once migration functions are registered.
