# Phase 04: Storage Package — Research

**Researched:** 2026-03-30
**Status:** Complete
**Discovery Level:** 1 (Quick Verification — well-understood domain, established codebase patterns)

## Standard Stack

All dependencies already present. No new external libraries needed.

| Component | Library | Version | Already in go.mod |
|-----------|---------|---------|-------------------|
| AES-256-GCM | `crypto/aes`, `crypto/cipher` | stdlib | ✓ |
| Argon2id | `golang.org/x/crypto/argon2` | v0.49.0 | ✓ |
| Windows MoveFileEx | `golang.org/x/sys/windows` | v0.42.0 | ✓ |
| SHA-256 | `crypto/sha256` | stdlib | ✓ |
| JSON | `encoding/json` | stdlib | ✓ |
| Binary I/O | `encoding/binary` | stdlib | ✓ |

## Architecture Patterns

### File Format (from `formato-arquivo-abditum.md`)

49-byte fixed header + encrypted payload:
- Bytes 0–3: magic `ABDT` (4 bytes ASCII)
- Byte 4: `versão_formato` (uint8, currently `1`)
- Bytes 5–36: salt (32 bytes)
- Bytes 37–48: nonce (12 bytes)
- Bytes 49+: AES-256-GCM ciphertext (JSON UTF-8) + 16-byte GCM tag

Full header is authenticated as GCM AAD — tamper with any byte and decryption fails.

### Crypto API Extension (per D-02 in CONTEXT.md)

Current `Encrypt`/`Decrypt` in `internal/crypto/aead.go`:
- `Encrypt(key, plaintext)` → nonce||ciphertext||tag (nil AAD)
- `Decrypt(key, ciphertext)` → plaintext (nil AAD)

New functions needed:
- `EncryptWithAAD(key, plaintext, aad)` → `(nonce, ciphertext, error)` — nonce and ciphertext returned separately (storage writes nonce to header, ciphertext to payload)
- `DecryptWithAAD(key, ciphertext, nonce, aad)` → `(plaintext, error)` — nonce passed explicitly (read from header)

### Serialization (per D-03 in CONTEXT.md, `arquitetura-dominio.md` §9)

Lives in `internal/vault/serialization.go`:
- `SerializarCofre(cofre *Cofre) ([]byte, error)` — accesses private fields, omits `estadoSessao == excluido` secrets
- `DeserializarCofre(data []byte) (*Cofre, error)` — rebuilds graph, populates parent-child refs via `popularReferencias()`, sets all secrets to `estadoSessao = original`
- `CampoSegredo.valor` ([]byte) serialized as UTF-8 string in JSON, not Base64

### Storage Package Structure

```
internal/storage/
  doc.go                    — package documentation
  errors.go                 — sentinel errors (ErrInvalidMagic, ErrVersionTooNew, ErrCorrupted)
  format.go                 — constants (magic, header size, format version), FileMetadata struct
  storage.go                — Save, SaveNew, Load, DetectExternalChange, RecoverOrphans
  migrate.go                — Migrate function with version chain scaffold
  atomic_rename_unix.go     — os.Rename (build tag: !windows)
  atomic_rename_windows.go  — MoveFileEx (build tag: windows)
  storage_test.go           — comprehensive tests
  testdata/                 — v1 fixture file for migration tests
```

### Repository Interface Adaptation

Current `RepositorioCofre` interface:
```go
type RepositorioCofre interface {
    Salvar(cofre *Cofre) error
    Carregar() (*Cofre, error)
}
```

Needs extension: `Carregar` must also return `FileMetadata` for external change detection. Options:
1. Change to `Carregar() (*Cofre, FileMetadata, error)` — breaking change to interface
2. Add separate `MetadataAtual() FileMetadata` method

Option 1 is cleaner since Phase 4 is the first implementation. The interface was a placeholder.

### Atomic Save Protocol (ATOMIC-01, ATOMIC-02)

For existing file:
1. If `.bak` exists → rename to `.bak2`
2. Rename current file → `.bak`
3. Write payload to `.abditum.tmp` in `filepath.Dir(vaultPath)` (NOT `os.TempDir()`)
4. Atomic rename `.tmp` → target (platform-specific)
5. Delete `.bak2` on success
6. On any failure → delete `.tmp` immediately

For new file (ATOMIC-03): Write directly to `destPath` without `.tmp`.

### External Change Detection (D-05 in CONTEXT.md)

`FileMetadata` struct: `{ Size int64, Hash [32]byte }`
- Fast-path: compare file sizes (O(1), no read needed)
- If sizes match: compute SHA-256 of full file, compare hashes
- `DetectExternalChange(path string, metadata FileMetadata) (bool, error)`

### Platform-Specific Atomic Rename

Unix (`atomic_rename_unix.go`, build tag `//go:build !windows`):
- `os.Rename(src, dst)` — POSIX rename is atomic

Windows (`atomic_rename_windows.go`, build tag `//go:build windows`):
- `windows.MoveFileEx(src, dst, windows.MOVEFILE_REPLACE_EXISTING)` via `golang.org/x/sys/windows`
- Standard `os.Rename` on Windows is NOT atomic when replacing existing file

Established pattern: `mlock_unix.go`, `mlock_windows.go`, `mlock_other.go` in crypto package.

## Don't Hand Roll

- Nonce generation: use `crypto/rand` via `io.ReadFull`, same as existing `Encrypt`
- GCM: use `crypto/cipher.NewGCM`, same as existing AEAD code
- Argon2id: use `DeriveKey` from `internal/crypto`
- JSON: use `encoding/json` with custom serialization in vault package

## Common Pitfalls

1. **Windows `os.Rename`**: NOT atomic when replacing — must use `MoveFileEx` with `MOVEFILE_REPLACE_EXISTING`
2. **`.tmp` location**: MUST be `filepath.Dir(vaultPath)`, not `os.TempDir()` — cross-device rename (`EXDEV`) fails on encrypted home dirs and network mounts
3. **Argon2id params from format version**: Derived from `versão_formato` via lookup table, NOT stored in file header (differs from ROADMAP description — CONTEXT.md D-01 overrides)
4. **`CampoSegredo.valor` as UTF-8 string**: `encoding/json` would Base64 encode `[]byte` — custom serialization needed
5. **GCM AAD**: Full 49-byte header is AAD — any header tampering causes auth failure
6. **Nonce in file vs. in Encrypt output**: Current `Encrypt` prepends nonce to output. New `EncryptWithAAD` must return nonce separately since storage writes it to header position
7. **`crypto.ErrAuthFailed` reuse**: Storage-level `ErrAuthFailed` could cause confusion with crypto-level. Storage should re-export or wrap distinctly. Actually, per CONTEXT.md agent discretion, storage can define its own sentinel errors.

## Key Decisions from CONTEXT.md

| ID | Decision | Impact |
|----|----------|--------|
| D-01 | `formato-arquivo-abditum.md` is canonical spec | Ignore ROADMAP format description |
| D-02 | Extend crypto with AAD variants | New functions in `internal/crypto/aead.go` |
| D-03 | Serialization in `vault/serialization.go`, no DTOs | Storage calls vault functions |
| D-04 | Storage orchestrates file I/O + encryption | Vault handles JSON marshal/unmarshal |
| D-05 | Size + SHA-256 for change detection, no mtime | `FileMetadata { Size, Hash }` |

## Validation Architecture

### Key Invariants to Test

1. **Roundtrip integrity**: `SaveNew` + `Load` returns identical vault (deep equality)
2. **Error classification**: wrong magic → `ErrInvalidMagic`, version too new → `ErrVersionTooNew`, wrong password → `ErrAuthFailed`, missing Pasta Geral → `ErrCorrupted`
3. **Atomic save sequence**: `.bak`/`.bak2` rotation on successive saves
4. **Orphan recovery**: stale `.tmp` removed, `.bak2` restored when target unreadable
5. **External change detection**: modified file size/content detected correctly
6. **Platform rename**: Windows uses `MoveFileEx`, Unix uses `os.Rename`
7. **AAD authentication**: tampered header bytes cause `ErrAuthFailed`
8. **Migration scaffold**: v1 fixture file loads and validates cleanly

### Test Strategy

- Unit tests in `internal/storage/storage_test.go`
- Use `t.TempDir()` for all file operations (auto-cleanup)
- Embed v1 fixture in `testdata/` for migration tests
- Test error sentinels with `errors.Is()`
- Test backup chain by performing 3 consecutive saves and verifying `.bak`/`.bak2` existence
- Simulate crash by pre-creating `.tmp` file, then calling `RecoverOrphans`
