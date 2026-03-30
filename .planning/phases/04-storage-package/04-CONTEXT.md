# Phase 4: Storage Package - Context

**Gathered:** 2026-03-30
**Status:** Ready for planning

<domain>
## Phase Boundary

`internal/storage` delivers the binary `.abditum` file format with atomic writes, `.bak`/`.bak2` backup chain, startup orphan recovery, external change detection, and a migration scaffold — including the Windows-specific `MoveFileEx` rename path — fully verified.

This phase implements:
- Binary file format reading/writing per `formato-arquivo-abditum.md`
- Atomic save protocol (`.tmp` → rename) with backup rotation
- Platform-specific atomic rename (Unix `os.Rename` / Windows `MoveFileEx`)
- Startup orphan recovery (`RecoverOrphans`)
- External change detection (size + SHA-256)
- Format version migration scaffold
- Extension of `internal/crypto` with AAD-aware encrypt/decrypt variants

This phase does NOT implement:
- TUI integration (Phase 5+)
- Vault creation/open flows (Phase 6)
- Any user-facing UI

</domain>

<decisions>
## Implementation Decisions

### File Format Specification

**D-01: `formato-arquivo-abditum.md` is the canonical spec**
- The ROADMAP Phase 4 Plan 1 describes a different format (magic `ABDITUM\x00` 8 bytes, `uint16` LE version, Argon2id params stored in header). This is outdated.
- Canonical format: magic `ABDT` (4 bytes ASCII), `versão_formato` (1 byte uint8), salt (32 bytes), nonce (12 bytes) = 49-byte fixed header.
- Argon2id parameters are NOT stored in the header — they are derived from `versão_formato` via a lookup table. Version 1 maps to: Argon2id m=256 MiB, t=3, p=4, keyLen=32.
- The ROADMAP pitfall about "breaking existing vaults on param changes" is addressed by the version-to-profile mapping: changing params means bumping the format version and adding a migration path.
- Full 49-byte header is authenticated as GCM AAD — tamper with any header byte and decryption fails.

### AAD + Crypto API Extension

**D-02: Extend `internal/crypto` with AAD-aware variants**
- Add `EncryptWithAAD(key, plaintext, aad []byte) (nonce, ciphertext []byte, err error)` — generates fresh nonce via `crypto/rand` internally, returns nonce and ciphertext separately.
- Add `DecryptWithAAD(key, ciphertext, nonce, aad []byte) ([]byte, error)` — accepts nonce explicitly (read from file header), decrypts with AAD verification.
- Existing `Encrypt`/`Decrypt` (without AAD) remain unchanged for backward compatibility.
- Nonce management: crypto generates nonce (consistent with existing `Encrypt` behavior). Storage writes the returned nonce into the header, then writes the ciphertext as payload.
- On load: storage reads nonce from header, builds full header as AAD, passes both to `DecryptWithAAD`.

### JSON Serialization Strategy

**D-03: Serialization lives in `vault/serialization.go`, no DTOs**
- `internal/storage` cannot access private fields of vault entities (lowercase fields, private to package).
- This is a direct consequence of the encapsulation architecture (fields private to prevent TUI mutation — see `arquitetura-dominio.md` §1).
- Solution: `SerializarCofre(cofre *Cofre) ([]byte, error)` and `DeserializarCofre(data []byte) (*Cofre, error)` live in `internal/vault/serialization.go`.
- These functions access private fields directly (same package) and handle the full domain graph.
- `SerializarCofre` omits secrets with `estadoSessao == excluido` (per D-16/D-17 from Phase 3).
- `DeserializarCofre` rebuilds the graph, sets all secrets to `estadoSessao = original`, and reconstitutes parent-child references via `popularReferencias()` in O(n).
- No intermediate DTOs — avoids structural duplication (6 entities × 2 structs), mechanical conversion code, and silent bugs when fields are added.
- No `MarshalJSON`/`UnmarshalJSON` on entities — avoids coupling domain to serialization and prevents accidental serialization in debug/logging contexts.
- `CampoSegredo.valor` (`[]byte`) must be serialized as UTF-8 string in JSON (not Base64) — custom handling in `SerializarCofre`/`DeserializarCofre` per `arquitetura.md` §5.
- `arquitetura-dominio.md` has been updated with a new §9 documenting this decision and its rationale.

**D-04: Storage calls vault serialization functions**
- Save flow: `Manager` → `vault.SerializarCofre(cofre)` → `[]byte` JSON → `crypto.EncryptWithAAD` → `storage.Write`
- Load flow: `storage.Read` → `crypto.DecryptWithAAD` → `[]byte` JSON → `vault.DeserializarCofre(data)` → `*Cofre`
- Storage handles: file I/O, header parsing/writing, encryption orchestration, atomic write protocol, backup chain.
- Vault handles: JSON marshaling/unmarshaling, domain graph reconstruction, reference population.

### External Change Detection

**D-05: Size + SHA-256, no mtime**
- `mtime` was discarded: resolution varies by filesystem (FAT32 = 2s, NTFS = 100ns, ext4 = 1ns), cloud sync (Dropbox, OneDrive) may or may not preserve it, cross-device copies can truncate it. False positives and negatives are real.
- Detection uses **file size** as fast-path (O(1), no file read needed) — if size differs, change is certain.
- If size matches, compute **SHA-256** of the full file and compare — confirms whether content actually changed.
- The vault file is small (few KB even with hundreds of secrets), so hashing cost is negligible.
- `FileMetadata` struct: `{ Size int64, Hash [32]byte }` — returned by `Load`, stored by Manager, checked before `Save`.
- `DetectExternalChange(path string, metadata FileMetadata) (bool, error)` — compares current file state against stored metadata.
- `arquitetura.md` §7 and `historico-decisoes.md` §17 have been updated to reflect this decision.

### Agent's Discretion
- Exact sentinel error names and messages for storage-level errors (`ErrInvalidMagic`, `ErrVersionTooNew`, `ErrAuthFailed`, `ErrCorrupted`)
- Internal structure of the migration scaffold (`Migrate` function, `MigrationFunc` chain)
- `RecoverOrphans` implementation details (what to check, what to delete, what to restore)
- Build tag structure for platform-specific atomic rename files
- Test fixture file format and embedding strategy for migration tests

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### File Format
- `formato-arquivo-abditum.md` — Complete binary format specification: header layout, AAD, key derivation profiles, opening sequence, error categories, write/compatibility rules

### Architecture
- `arquitetura.md` §2 — Package structure (`internal/storage`, `internal/vault`, `internal/crypto`)
- `arquitetura.md` §3 — Manager pattern, persistence flow, save atomicity
- `arquitetura.md` §5 — Build conventions (CGO_ENABLED=0, no net imports, crypto/rand only), `CampoSegredo.valor` serialized as UTF-8 string not Base64
- `arquitetura.md` §6 — Session security: mlock, wipe, clipboard, clear screen
- `arquitetura.md` §7 — External change detection: size + SHA-256 (no mtime)

### Domain Layer
- `arquitetura-dominio.md` §9 — Serialization architecture: why it lives in vault, not storage. SerializarCofre/DeserializarCofre API
- `arquitetura-dominio.md` §3 — Cofre as aggregate root, save atomicity (D-17: two-phase commit)
- `arquitetura-dominio.md` §4 — Manager pattern, persistence vs mutation operations
- `arquitetura-dominio.md` §7 — Session state tracking, estadoSessao transitions, deletion finalization on save

### Requirements
- `.planning/REQUIREMENTS.md` §Salvamento Atômico — ATOMIC-01 through ATOMIC-04
- `.planning/REQUIREMENTS.md` §Compatibilidade — COMPAT-03 (backward compatibility, migration)
- `.planning/REQUIREMENTS.md` §Ciclo de Vida — VAULT-07, VAULT-08 (external change detection on save/discard)

### Existing Code
- `internal/vault/repository.go` — `RepositorioCofre` interface (Salvar/Carregar) that storage must implement
- `internal/crypto/aead.go` — Current Encrypt/Decrypt (nil AAD) — new EncryptWithAAD/DecryptWithAAD extend this
- `internal/crypto/kdf.go` — `DeriveKey`, `GenerateSalt`, `ArgonParams` struct
- `internal/crypto/errors.go` — Existing sentinel errors: `ErrAuthFailed`, `ErrInvalidParams`, `ErrInsufficientEntropy`
- `internal/vault/entities.go` — All entity types with private fields
- `internal/vault/manager.go` — Manager struct, NewManager, existing method patterns

### Prior Phase Context
- `.planning/phases/02-crypto-package/02-CONTEXT.md` — Crypto API decisions (D-18: separate functions, D-19: nonce internal to Encrypt, D-13: all sensitive data as []byte)
- `.planning/phases/03-vault-domain-manager/03-CONTEXT.md` — Domain decisions (D-01: no synthetic IDs, D-14: deleted secrets filtering, D-17: atomic save two-phase commit)

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `crypto.DeriveKey(password, salt, params)` — Key derivation ready to use
- `crypto.GenerateSalt()` — CSPRNG salt generation (32 bytes)
- `crypto.Encrypt/Decrypt` — Base AES-256-GCM implementation to extend with AAD variants
- `crypto.Wipe([]byte)` — Secure memory zeroing for key material after use
- `crypto.ArgonParams` — Struct for Argon2id parameters (t, m, p, keyLen)
- `crypto.FormatVersion` — Already exported as `const FormatVersion = 1`
- `vault.RepositorioCofre` interface — Storage must implement `Salvar(*Cofre) error` and `Carregar() (*Cofre, error)`

### Established Patterns
- Platform-specific build tags: `mlock_unix.go`, `mlock_windows.go`, `mlock_other.go` — same pattern for atomic rename
- Sentinel errors: `var ErrX = errors.New("...")` — extend for storage errors
- TDD: all existing packages developed test-first with comprehensive coverage
- `CGO_ENABLED=0` throughout — no C dependencies, static binary

### Integration Points
- `RepositorioCofre` interface in `internal/vault/repository.go` — storage implementation plugs in here
- `Manager.repositorio` field — receives storage instance via dependency injection
- `vault.SerializarCofre` / `vault.DeserializarCofre` — new functions that storage will call for JSON
- `crypto.EncryptWithAAD` / `crypto.DecryptWithAAD` — new functions that storage will call for encryption

</code_context>

<specifics>
## Specific Ideas

- The ROADMAP mentions `FileMetadata` carrying `mtime` and `size` — adjust to carry `size` and `hash [32]byte` instead.
- The ROADMAP Plan 1 describes a different binary format — ignore it entirely, follow `formato-arquivo-abditum.md`.
- `RepositorioCofre` interface may need adjustment: `Carregar` currently returns `(*Cofre, error)` but needs to also return `FileMetadata`. Consider `Carregar() (*Cofre, FileMetadata, error)` or surface metadata separately.
- The `formato-arquivo-abditum.md` specifies validation step 8: "Pasta Geral deve existir com nome 'Geral'" — this is the integrity check after JSON deserialization, implemented via `DeserializarCofre` validation.

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 04-storage-package*
*Context gathered: 2026-03-30*
