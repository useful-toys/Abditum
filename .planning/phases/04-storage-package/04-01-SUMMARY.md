# Plan 04-01 Summary

**Phase:** 04-storage-package
**Plan:** 01 — Foundation Layer
**Status:** COMPLETE
**Completed:** 2026-03-30

## What Was Done

Built the foundation layer required by Plans 04-02 through 04-04.

### Task 1: AAD-aware Crypto Functions (completed in prior session)
- Added `EncryptWithAAD(key, plaintext, aad []byte) (nonce, ciphertext []byte, err error)` to `internal/crypto/aead.go`
- Added `DecryptWithAAD(key, ciphertext, nonce, aad []byte) ([]byte, error)` to `internal/crypto/aead.go`
- Returns nonce and ciphertext **separately** (nonce goes to header, ciphertext is payload)
- Commit: `feat(04-01): add EncryptWithAAD and DecryptWithAAD with AAD parameter` (49ca52d)

### Task 2: Vault JSON Serialization
- Created `internal/vault/serialization.go` with:
  - `SerializarCofre(cofre *Cofre) ([]byte, error)` — serializes to JSON, omits `EstadoExcluido` secrets, writes `pasta_geral.nome` as "Geral", encodes `CampoSegredo.valor` as UTF-8 string
  - `DeserializarCofre(data []byte) (*Cofre, error)` — deserializes JSON, validates `pasta_geral.nome == "Geral"`, sets all secrets to `EstadoOriginal`, populates parent-child refs
  - `popularReferencias(pasta *Pasta, pai *Pasta)` — recursive reference population
- Created `internal/vault/serialization_test.go` with 16 tests covering all behaviors
- Commit: `feat(04-01): add SerializarCofre and DeserializarCofre with roundtrip tests` (781481f)

### Task 3: Storage Package Foundation
- Updated `internal/storage/doc.go` with full package documentation
- Created `internal/storage/errors.go` with `ErrInvalidMagic`, `ErrVersionTooNew`, `ErrCorrupted`
- Created `internal/storage/format.go` with:
  - `var Magic = [4]byte{'A', 'B', 'D', 'T'}`
  - `const HeaderSize = 49` (MagicSize + VersionSize + SaltSize + NonceSize)
  - `const CurrentFormatVersion uint8 = 1`
  - `type FileMetadata struct { Size int64; Hash [32]byte }`
  - `type FormatProfile struct { Time, Memory uint32; Threads uint8; KeyLen uint32 }`
  - `func ProfileForVersion(version uint8) (FormatProfile, error)`
  - `func (p FormatProfile) ToArgonParams() crypto.ArgonParams`
- Commit: `feat(04-01): add storage package foundation with format constants, errors, and version profiles` (a191024)

## Verification Results

```
go test ./internal/crypto/... -count=1  → PASS
go test ./internal/vault/... -count=1   → PASS (16 new serialization tests + all existing)
go build ./internal/storage/...         → OK
go vet ./...                            → OK
```

## Key Decisions

- `NovoCofre()` creates `pastaGeral.nome = "Pasta Geral"` (internal name). Serialization writes "Geral" unconditionally for `pasta_geral.nome`; deserialization validates exactly "Geral". The internal nome field is reconstructed from the JSON `nome` value.
- Serialization intermediate structs are private to `serialization.go` (no exported DTOs crossing package boundaries).
- `ErrVersionTooNew` wraps with version number for diagnosability: `fmt.Errorf("%w: version %d", ErrVersionTooNew, version)`.
