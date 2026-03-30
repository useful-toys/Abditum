# Phase 4: Storage Package - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-30
**Phase:** 04-storage-package
**Areas discussed:** File format spec resolution, AAD + crypto API gap, JSON serialization strategy, External change detection API

---

## File Format Spec Resolution

| Option | Description | Selected |
|--------|-------------|----------|
| formato-arquivo-abditum.md is canonical | ABDT magic (4b), uint8 version, params derived from version, 49-byte fixed header, full header as GCM AAD. Ignore ROADMAP Plan 1 divergences. | ✓ |
| ROADMAP Plan 1 is canonical | ABDITUM\x00 magic (8b), uint16 version, Argon2id params in header, variable header size. | |
| Hybrid: formato-arquivo spec + stored params | Follow formato-arquivo mostly, but store Argon2id params in header for forward-compatibility. | |

**User's choice:** formato-arquivo-abditum.md is canonical
**Notes:** The ROADMAP Plan 1 predates the detailed spec and contains outdated format decisions. The version-to-profile mapping addresses the ROADMAP's pitfall about param changes.

---

## AAD + Crypto API Gap

### Main question: How to handle AAD

| Option | Description | Selected |
|--------|-------------|----------|
| Extend crypto package with AAD variants | Add EncryptWithAAD/DecryptWithAAD to internal/crypto. Keeps crypto logic centralized. Existing Encrypt/Decrypt unchanged. | ✓ |
| Storage handles GCM directly | Storage uses crypto/aes + cipher.NewGCM directly, bypassing internal/crypto for encrypt/decrypt. | |
| Refactor Encrypt/Decrypt to accept external nonce | More invasive change to existing API. | |

### Follow-up: Nonce management

| Option | Description | Selected |
|--------|-------------|----------|
| Crypto generates nonce, returns it separately | EncryptWithAAD generates nonce internally and returns (nonce, ciphertext, error). Consistent with current Encrypt behavior. | ✓ |
| Caller provides nonce explicitly | Storage generates nonce, passes it to EncryptWithAAD. Full control for storage layer. | |
| Low-level Seal/Open API in crypto | Storage splits steps: generate nonce → write header → use low-level GCM Seal/Open. | |

**User's choice:** Extend crypto + crypto generates nonce
**Notes:** Keeps all crypto logic centralized. Nonce returned separately so storage can write it to the 49-byte header.

---

## JSON Serialization Strategy

### Initial discussion: Where does serialization live?

| Option | Description | Selected |
|--------|-------------|----------|
| Serializer in storage (via getters + factories) | Storage serializes using public getters, deserializes using vault factory functions. | |
| Package-level functions in vault (serialization.go) | SerializarCofre/DeserializarCofre in vault package. Accesses private fields. Entities untouched. | ✓ |

**User's initial thought:** Could serialization use public getters from storage?
**Resolution:** Serialization (write) could use getters, but deserialization (read) needs to populate private fields without triggering session state tracking (Manager methods mark state as Modificado/Incluido). Therefore both directions must live in the vault package.

### Follow-up: DTOs vs no DTOs

User explicitly rejected intermediate DTOs ("não quero DTOs intermediários").

### Follow-up: MarshalJSON vs package-level functions

| Option | Description | Selected |
|--------|-------------|----------|
| MarshalJSON/UnmarshalJSON on each entity | encoding/json calls automatically. Distributes serialization across entities. | |
| Package-level SerializarCofre/DeserializarCofre | Separate functions in serialization.go. Explicit, auditable. | ✓ |
| MarshalJSON only on Cofre | Single method traverses entire tree. Centralized but large method. | |

**User's choice:** Package-level functions in serialization.go, no DTOs
**Notes:** This is a consequence of the encapsulation architecture — private fields protect against TUI mutation but also prevent external serialization. The user requested that `arquitetura-dominio.md` be updated with a new section explaining this decision and its rationale. Section 9 was added, responsibility table updated, and principles updated.

---

## External Change Detection API

### Main question: Detection mechanism

| Option | Description | Selected |
|--------|-------------|----------|
| mtime + size | Load returns FileMetadata (mtime + size). DetectExternalChange compares with disk state. | |
| mtime + size + hash SHA-256 | Same as above but adds SHA-256 hash for infallible detection. | (initially selected) |
| Detection embedded in Save (automatic) | Manager.Salvar() checks automatically and returns specific error if file changed. | |

### Follow-up: mtime reliability

User questioned mtime reliability. After discussion, mtime was dropped entirely.

**Final decision:** Size + SHA-256 only (no mtime)
- Size as fast-path: if different, change is certain without reading file
- SHA-256 for confirmation when size matches
- Vault file is small (few KB), hash cost is negligible
- mtime discarded: variable resolution across filesystems (FAT32=2s, NTFS=100ns), cloud sync inconsistencies, cross-device copy truncation

**Notes:** `arquitetura.md` §7 and `historico-decisoes.md` §17 updated to reflect this decision.

---

## Agent's Discretion

- Exact sentinel error names and messages for storage errors
- Migration scaffold internal structure
- RecoverOrphans implementation details
- Build tag file organization for atomic rename
- Test fixture format and embedding strategy

## Deferred Ideas

None — discussion stayed within phase scope.
