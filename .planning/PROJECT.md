# Abditum

## What This Is

Abditum is a portable, secure, offline password vault with a modern TUI interface. It stores and manages secrets (passwords, API keys, credit card data, and user-defined structured secrets) in an AES-256-GCM encrypted file, requiring no installation, no cloud account, and persisting no data outside the vault file (except transitory artifacts and backups explicitly provided by the application). It targets privacy-conscious users, developers, and sysadmins who want full ownership of their credentials data.

## Core Value

A single portable binary that any user can carry on a USB drive and use on any machine — opening, managing, and saving an encrypted vault without installing anything or touching the cloud.

## Requirements

### Validated

<!-- Shipped and confirmed valuable. -->

(None yet — ship to validate)

### Active

**Vault Lifecycle**
- [ ] Create new vault with path and master password (double-entry confirmation for master password)
- [ ] Open existing vault with master password
- [ ] Save vault atomically (write to .abditum.tmp → rename on success; delete tmp on failure)
- [ ] Keep `.abditum.bak` backup of previous vault on every save: before writing the new backup, rename any existing `.abditum.bak` to `.abditum.bak2`; on success delete `.abditum.bak2`; on failure restore `.abditum.bak2` → `.abditum.bak` where possible
- [ ] Save vault to a new path (Save As)
- [ ] Discard unsaved changes and reload vault from disk
- [ ] Change master password (double-entry confirmation required)
- [ ] Export vault to plain-text JSON (with security warning + confirmation)
- [ ] Import vault from plain-text JSON with conflict handling:
  - Folders with same identity → hierarchy merged silently
  - Secrets with same identity (ID collision) → imported secret receives a new identity; its data is preserved
  - Secrets with same name in the same destination folder → name suffixed incrementally (e.g. "Secret (1)"); user is notified of all name-suffixed imports
  - Templates with same identity → imported template replaces existing silently
- [ ] Configure vault settings (stored inside vault file): auto-lock timeout, field reveal timeout, clipboard clear timeout

**Authentication**
- [ ] Unlock vault with master password
- [ ] Manual lock: return to vault open screen, clear app-controlled sensitive buffers where possible
- [ ] Auto-lock after configurable inactivity timeout (same behavior as manual lock)
- [ ] Brute-force protection via Argon2id (high memory and time cost)
- [ ] Spinner/progress indicator during Argon2id key derivation

**Vault Navigation (read-only)**
- [ ] Display vault hierarchy (folders and secrets) in a sidebar tree
- [ ] Display secret details in a detail panel
- [ ] Temporarily reveal a sensitive field value; auto-hide after configurable timeout (suggested default 15 s)
- [ ] Virtual folder: Favorites (all favorited secrets, regardless of location)
- [ ] Virtual folder: Trash / Exclusão Reversível (all soft-deleted secrets, restorable until next save)

**Secret Management**
- [ ] Create secret from any existing template, OR start with a blank secret (no initial fields); template is a snapshot at creation time — changes to the template do not affect existing secrets; template name stored as historical record only
- [ ] Create secret at vault root or inside any folder
- [ ] Duplicate an existing secret
- [ ] Favorite / unfavorite a secret
- [ ] Edit secret data: change secret name, field values, and observation
- [ ] Alter secret structure: add, remove, rename, reorder fields on an existing secret; changing a field's type is not supported — delete the field and add a new one with the desired type instead
- [ ] Soft-delete a secret (reversible until next save; deleted secret moves to Trash virtual folder)
- [ ] Restore a soft-deleted secret from Trash before the next save
- [ ] Move secret to another folder or vault root
- [ ] Reorder secret relative to other items in same parent
- [ ] Search secrets by: secret name, field name, `texto`-type field value, or note (observation); `texto sensível` fields are never searched; all matching secrets shown

**Folder Management**
- [ ] Create folder at vault root or inside another folder
- [ ] Rename folder
- [ ] Move folder to another folder or vault root
- [ ] Reorder folder relative to other items in same parent
- [ ] Delete folder: its secrets and subfolders are promoted to the parent (or vault root); promoted secrets appended to end of parent's secret list; promoted subfolders appended to end of parent's folder list; folder itself is permanently removed (not soft-deleted)
- [ ] Pre-defined folders on vault creation: "Sites", "Financeiro", "Serviços" (user-editable/removable)

**Secret Template Management**
- [ ] Create custom template with named, typed fields
- [ ] Edit template (add, remove, rename, retype, reorder fields); changes only affect future secrets
- [ ] Delete template; existing secrets are unaffected
- [ ] Create template from existing secret (copies field names and types as initial structure; values are not copied)
- [ ] Pre-defined templates on vault creation: "Login" (URL, Username, Password), "Cartão de Crédito" (Número do Cartão, Nome no Cartão, Data de Validade, CVV), "API Key" (Nome da API, Chave de API)
- [ ] Templates stored inside the vault file; each vault has its own independent template set

**Clipboard**
- [ ] Copy any field value to the system clipboard
- [ ] Auto-clear clipboard after configurable timeout (suggested default 30 s)
- [ ] Clear clipboard on lock or vault close

**Security**
- [ ] Minimize sensitive data retention in memory; clear app-controlled buffers on lock or close where possible
- [ ] Shoulder-surfing protection: hotkey to instantly hide the entire TUI; restore on same hotkey
- [ ] Export warning: show security risk message and require explicit confirmation before exporting plain-text JSON

**Data Model**
- [ ] Vault payload contains: settings (auto-lock timeout default 2 min, field reveal timeout default 15 s, clipboard clear timeout default 30 s), secrets (root level), folders (root level), secret templates, creation date, last-modified date
- [ ] Secret fields: id (NanoID 6 chars), name, template name (optional, historical), fields list, favorite flag, note/observation (implicit on every secret, not part of templates), creation date, last-modified date
- [ ] Folder fields: id (NanoID 6 chars), name, secrets list, subfolders list (recursive)
- [ ] Template fields: id (NanoID 6 chars), name, template-fields list
- [ ] Field types: `texto` (plain text) and `texto sensível` (sensitive text); field value may be empty string
- [ ] Name uniqueness not enforced: duplicate names are permitted for secrets, folders, and templates
- [ ] Every secret has an implicit observation field (free text, non-sensitive, not in templates, cannot be removed)

**Format & Compatibility**
- [ ] Vault file format: `.abditum` extension
- [ ] Binary file: fixed header (`magic=ABDT` + `versão_formato` + `salt` + `nonce`) followed by AES-256-GCM ciphertext of JSON payload; entire header is GCM AAD
- [ ] Argon2id parameters for key derivation (high-cost: ≥64 MiB memory, time=3)
- [ ] App version N can open vaults created with any previously supported format version; payload migrated in-memory; saved in current format version
- [ ] No application logs (stdout/stderr) containing vault paths, secret names, or field values
- [ ] No data persisted outside the vault file except `.abditum.tmp` (transitory) and `.abditum.bak` (backup)

### Out of Scope

- Cloud sync / cloud storage — zero-knowledge by design; no network dependency
- Multiple vaults open simultaneously — single vault at a time, by design
- Mobile or web app — portable TUI is the product
- Tags — folder hierarchy is sufficient for v1
- Secret version history — deferred to v2
- TOTP/2FA field type — deferred to v2
- Password generator — deferred to v2
- Duress password — deferred to v2
- QR Code sharing — deferred to v2
- Vault health / password audit report — deferred to v2
- Hardware token / keyfile second factor — deferred to v2
- Changing the type of an existing secret field — delete the field and add a new one with the desired type instead
- Re-authentication to save — the vault is already unlocked and authenticated in the session; the master password is not required again to save

## Context

- **Greenfield project** — no existing codebase; starting from scratch
- **Language:** Go — compiles to a single static binary, strong stdlib crypto, excellent cross-platform support
- **TUI:** Bubble Tea + Lip Gloss — modern component-based TUI framework, production-proven
- **Vault format:** AES-256-GCM for data encryption, Argon2id for key derivation from master password
- **Target users:** privacy-conscious individuals, developers, sysadmins who distrust cloud password managers
- **Portability requirement** rules out external config files, registry writes, or any persistent data outside the vault file path

## Constraints

- **Tech Stack:** Go + Bubble Tea + Lip Gloss (TUI) + teatest/v2 (TUI testing) — single static binary, no runtime dependencies, CGO_ENABLED=0
- **Crypto:** AES-256-GCM + Argon2id — non-negotiable; Argon2id used exclusively for KDF (not content integrity); GCM handles integrity via authentication tag
- **Argon2id parameters (v1, hard-coded):** memory=256 MiB, time=3 iterations, parallelism=min(4, logical CPUs); security floor=128 MiB (app refuses to operate below); reference ceiling=512 MiB; parameters fixed per app version, not configurable per file
- **Portability:** No data persisted outside the vault file (except `.abditum.tmp` and `.abditum.bak`); no OS registry, config directory, or external files
- **Compatibility:** Windows, macOS, Linux 64-bit — all three from day one; same security policy on all platforms
- **Privacy:** Zero application logs containing sensitive data (paths, names, values)
- **Reliability:** Atomic save (`.abditum.tmp` → rename) + `.abditum.bak` backup rotation (existing `.bak` renamed to `.bak2` before new backup is written; `.bak2` deleted on success, restored on failure); nonce regenerated on every save; salt unique per vault (random at creation, stored in header)
- **Backward Compatibility:** App version N opens vaults from any previously supported format version; Argon2id profile selected by `versão_formato` from header; format version > app max → explicit incompatibility error; migration in-memory; saves always in current format
- **CI:** Mandatory; runs on Ubuntu, Windows, macOS runners

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Go as implementation language | Single static binary requirement; excellent cross-platform builds; strong stdlib crypto; CGO_ENABLED=0 | — Pending |
| Bubble Tea + Lip Gloss for TUI; teatest/v2 for TUI tests | Modern, component-based; widely adopted; pairs naturally with Go; teatest/v2 is the official testing companion | — Pending |
| AES-256-GCM + Argon2id | Industry-standard symmetric encryption + memory-hard KDF; strong against offline brute-force | — Pending |
| Argon2id used exclusively for KDF | Content integrity handled by GCM authentication tag; clean separation of concerns | — Pending |
| Argon2id v1 params: 256 MiB, time=3, parallelism=min(4,CPUs) | UX target 0.8–1.5 s unlock on modern hardware; strong offline resistance; floor=128 MiB, ceiling=512 MiB | — Pending |
| Argon2id params hard-coded per app version, not stored in file | Simplifies format; eliminates param downgrade attack via crafted file; all v1 vaults use the same strong params | — Pending |
| Argon2id profile selected by `versão_formato` when opening old vaults | Decouples app params from file params; allows future param evolution without breaking old vaults | — Pending |
| Vault format: encrypted JSON, UTF-8 | Human-readable internals (after decryption), easy migration/export, flexible schema; UTF-8 supports all characters | — Pending |
| Binary header: `magic(ABDT)` + `versão_formato` + `salt` + `nonce` | Enables format detection before decryption; magic bytes distinguish wrong-file-type from wrong-password | — Pending |
| Entire header as GCM AAD | Tamper-detection of magic/version/salt/nonce at no extra cost; no separate checksum needed | — Pending |
| Salt: random at creation, unique per vault, stored in header | Standard Argon2id practice; prevents rainbow table attacks across vaults | — Pending |
| Nonce: regenerated on every save | Prevents nonce reuse under the same key; nonce reuse with GCM is catastrophic | — Pending |
| Format version > app max → explicit incompatibility error | Prevents silent data corruption when a newer vault is opened with an older app | — Pending |
| Trash = virtual folder, purge on save | No separate "empty trash" UX; save intent is clear; keeps model simple | — Pending |
| Folders are NOT soft-deleted | Only secrets are soft-deleted; folders are permanently removed and their children promoted | — Pending |
| Templates stored inside vault | Each vault is self-contained; portability guaranteed; no shared config files | — Pending |
| Templates are snapshots at creation time | Avoids retroactive mutations of existing secrets when templates evolve | — Pending |
| Favorites + Trash as only virtual folders | Covers the primary access patterns (quick access + undo delete) without over-engineering | — Pending |
| Search is in-memory sequential scan | Vault is fully in memory after unlock; no index needed; avoids any persistent index data | — Pending |
| NanoID (6 chars, 62^6 ≈ 56B combinations) for entity IDs | Uniqueness for identity/conflict resolution in import; names are not identifiers; JSON-serializable | — Pending |
| Name uniqueness not enforced | Simplifies data model; UX responsibility, not data integrity concern | — Pending |
| Implicit observation field on every secret | Always available for notes without polluting templates; treated as non-sensitive | — Pending |
| Import: identity collision → new ID for imported secret | Prevents secret data loss without silently overwriting existing secrets; preserves both | — Pending |
| Import: name collision in same folder → suffix name | Avoids visual ambiguity; user is notified; distinct from identity collision resolution | — Pending |
| Import: folders merged, templates replaced | Merge preserves structure; template replace assumes newer = better | — Pending |
| Vault root is a nameless implicit folder | Recursive hierarchy: root contains secrets and folders just like any other folder; no special-casing in data model | — Pending |
| JSON array order = display order | No separate order field; position in the JSON array is the canonical display order for secrets, folders, and fields | — Pending |
| .bak2 rotation before generating new backup | Prevents losing the last good backup if a crash occurs mid-save; safe rollback to .bak2 on failure | — Pending |
| Lock = minimize retention, clear controlled buffers | Go GC limits hard zeroing guarantees; honest about what the app can and cannot control | — Pending |
| CI mandatory on Ubuntu + Windows + macOS | Cross-platform guarantee; catches platform-specific regressions early | — Pending |

---
*Last updated: 2026-03-24 after 4th requirement revision (import two-level conflict rules, .bak2 rotation, re-auth-to-save out of scope, root-as-folder and JSON-order modeling decisions, settings defaults)*
