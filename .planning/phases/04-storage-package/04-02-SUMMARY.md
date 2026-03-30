# Plan 04-02 Summary

**Phase:** 04-storage-package
**Plan:** 02 — Save/Load Functions + Atomic Rename
**Status:** COMPLETE
**Completed:** 2026-03-30

## What Was Done

Implemented the core I/O layer: `SaveNew`, `Save`, `Load`, and platform-specific atomic rename.

### Task 1: SaveNew, Save, and Load (storage.go)

Created `internal/storage/storage.go` with three exported functions:

**`SaveNew(destPath string, cofre *vault.Cofre, password []byte) error`**
- Serializes via `vault.SerializarCofre`, generates fresh salt + nonce
- Builds the full 49-byte header (Magic + Version + Salt + Nonce) and uses it as GCM AAD
- Calls `crypto.SealWithAAD` (caller-provided nonce) to encrypt
- Writes `header + ciphertext` to `destPath` with mode 0600
- Wipes JSON bytes and key via `crypto.Wipe` after use

**`Save(vaultPath string, cofre *vault.Cofre, password, salt []byte) error`**
- Same encrypt flow as `SaveNew` but reuses the provided salt (preserves password-derived key)
- Atomic rotation protocol:
  1. Write encrypted content to `vaultPath + ".tmp"`
  2. Rotate existing `.bak` → `.bak2` (if `.bak` exists)
  3. Rename `vaultPath` → `.bak`
  4. `atomicRename(".tmp" → vaultPath)` via platform-specific function
  5. Cleanup `.tmp` on any failure (best-effort `os.Remove`)

**`Load(vaultPath string, password []byte) (*vault.Cofre, FileMetadata, error)`**
- Reads entire file; computes `FileMetadata{Size, Hash: sha256.Sum256(data)}`
- Validates header: length >= 49, magic bytes == `ABDT`, version known via `ProfileForVersion`
- Derives key, decrypts via `crypto.DecryptWithAAD` with full 49-byte header as AAD
- Propagates `crypto.ErrAuthFailed` directly; wraps deserialization errors as `ErrCorrupted`

### Task 2: Platform-specific Atomic Rename

- `internal/storage/atomic_rename_unix.go` (`//go:build !windows`): thin wrapper over `os.Rename`
- `internal/storage/atomic_rename_windows.go` (`//go:build windows`): uses `windows.MoveFileEx` with `MOVEFILE_REPLACE_EXISTING` for true NTFS atomic replace

### Tests (storage_test.go) — 12 tests, all pass

| Test | Verifies |
|------|----------|
| `TestSaveNew_RoundTrip` | SaveNew + Load returns equivalent Cofre |
| `TestSaveNew_HeaderMagic` | Header bytes 0-3 = "ABDT", byte 4 = version 1 |
| `TestLoad_WrongMagic` | → ErrInvalidMagic |
| `TestLoad_FileTooShort` | → ErrInvalidMagic |
| `TestLoad_VersionTooNew` | version=255 → ErrVersionTooNew |
| `TestLoad_WrongPassword` | → crypto.ErrAuthFailed |
| `TestLoad_TamperedHeader` | flip salt byte → ErrAuthFailed (AAD auth) |
| `TestLoad_TamperedPayload` | flip payload byte → ErrAuthFailed |
| `TestSave_CreatesBackup` | .bak file exists after Save |
| `TestSave_CreatesBak2WhenBakExists` | second Save creates .bak2 |
| `TestSave_NoTmpAfterSuccess` | .tmp cleaned up on success |
| `TestLoad_FileMetadata` | Size and SHA-256 Hash match file-on-disk |

## Verification Results

```
CGO_ENABLED=0 go test ./internal/storage/... -count=1 -v  → 12/12 PASS (5.6s)
CGO_ENABLED=0 go test ./...                               → all packages PASS
CGO_ENABLED=0 go vet ./...                                → OK
```

## Key Decisions

- Used `SealWithAAD` (caller-provided nonce) rather than `EncryptWithAAD` (generates nonce internally) to resolve the chicken-and-egg problem: nonce must be in the header for AAD, but AAD is needed at encryption time. Storage generates nonce via `crypto/rand`, builds the full 49-byte header, uses it as AAD, then seals.
- `.bak2` is **not** removed on successful save — it is kept as a secondary backup for user recovery. Only `.tmp` is cleaned up.
- Salt is a caller parameter to `Save` (not re-generated) so the password-derived key is stable across saves. Salt only changes on explicit password change.
- Argon2id params for v1: Time=3, Memory=256 MiB, Threads=4 (from `ProfileForVersion`). Tests accept the ~0.3s per-call cost (no fast-param override needed for 12 tests).
