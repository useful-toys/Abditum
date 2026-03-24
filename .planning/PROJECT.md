# Abditum

## What This Is

Abditum is a portable, secure, offline password vault with a modern TUI interface. It stores and manages secrets (passwords, API keys, credit card data, and user-defined structured secrets) in an AES-256-GCM encrypted file, requiring no installation, no cloud account, and leaving no traces on the host system. It targets privacy-conscious users, developers, and sysadmins who want full ownership of their credentials data.

## Core Value

A single portable binary that any user can carry on a USB drive and use on any machine — opening, managing, and saving an encrypted vault without installing anything or touching the cloud.

## Requirements

### Validated

<!-- Shipped and confirmed valuable. -->

(None yet — ship to validate)

### Active

**Vault Lifecycle**
- [ ] Create new vault with path and master password
- [ ] Open existing vault with master password
- [ ] Save vault atomically (write to .abditum.tmp, rename on success, delete tmp on failure)
- [ ] Keep .abditum.bak backup of previous vault on every save
- [ ] Save vault to a new path (Save As)
- [ ] Discard unsaved changes and reload vault from disk
- [ ] Change master password
- [ ] Export vault to plain-text JSON (with security warning + confirmation)
- [ ] Import vault from plain-text JSON (with incremental suffix on name collision: "Secret" → "Secret (1)", "Secret (2)")
- [ ] Configure vault settings: auto-lock timeout, field reveal timeout, clipboard clear timeout

**Authentication**
- [ ] Create vault with master password
- [ ] Unlock vault with master password
- [ ] Manual lock (return to auth screen, clear sensitive data from memory)
- [ ] Auto-lock after configurable inactivity timeout
- [ ] Brute-force protection via Argon2id (high memory and time cost)

**Vault Navigation (read-only)**
- [ ] Display vault hierarchy (folders and secrets)
- [ ] Display secret details
- [ ] Temporarily reveal a sensitive field value; auto-hide after configurable timeout (default 15s)
- [ ] Virtual folder: Favorites
- [ ] Virtual folder: Trash (Lixeira)

**Secret Management**
- [ ] Create secret from predefined or custom template
  - Template is a snapshot at creation time; changing/deleting template does not affect existing secrets
- [ ] Create secret at vault root or inside any folder
- [ ] Duplicate a secret
- [ ] Favorite / unfavorite a secret
- [ ] Edit secret (field values, add/remove/reorder/retype fields — advanced mode)
- [ ] Delete secret (moves to Trash; purged permanently on save)
- [ ] Restore secret from Trash
- [ ] Move secret to another folder or vault root
- [ ] Reorder secret within its parent
- [ ] Search secrets by name, non-sensitive field value, note, or containing folder name

**Folder Management**
- [ ] Create folder at vault root or inside another folder
- [ ] Rename folder
- [ ] Move folder to another folder or vault root
- [ ] Reorder folder within its parent
- [ ] Delete folder (moves children to parent folder or vault root; folder itself goes to Trash)
- [ ] Pre-defined folders on vault creation: "Sites", "Financeiro", "Serviços" (user-editable/removable)

**Secret Template Management**
- [ ] Create custom template with named, typed fields
- [ ] Edit template (affects only future secrets; existing secrets unchanged)
- [ ] Delete template
- [ ] Create template from existing secret (copy field names and types, not values)
- [ ] Pre-defined templates on vault creation: "Login" (URL, Username, Password), "Cartão de Crédito" (Número do Cartão, Nome no Cartão, Data de Validade, CVV), "API Key" (Nome da API, Chave de API)
- [ ] Templates stored inside the vault file (each vault has its own template set)

**Clipboard**
- [ ] Copy any field value to clipboard
- [ ] Auto-clear clipboard after configurable timeout (default 30s)
- [ ] Clear clipboard on lock or vault close

**Security**
- [ ] Zero sensitive data from memory on lock or close
- [ ] Shoulder-surfing protection: hotkey to instantly hide the entire TUI
- [ ] Export warning: show security risk message and require explicit confirmation before exporting plain-text JSON

**Format & Compatibility**
- [ ] Vault file format: .abditum extension, JSON payload, AES-256-GCM encrypted, Argon2id key derivation
- [ ] Version field in encrypted payload header for N-1 backward compatibility
- [ ] No application logs containing vault paths, secret names, or field values

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

## Context

- **Greenfield project** — no existing codebase; starting from scratch
- **Language:** Go — compiles to a single static binary, strong stdlib crypto, excellent cross-platform support
- **TUI:** Bubble Tea + Lip Gloss — modern component-based TUI framework, production-proven
- **Vault format:** AES-256-GCM for data encryption, Argon2id for key derivation from master password
- **Target users:** privacy-conscious individuals, developers, sysadmins who distrust cloud password managers
- **Portability requirement** rules out external config files, registry writes, or temp files

## Constraints

- **Tech Stack:** Go + Bubble Tea + Lip Gloss — single static binary target, no runtime dependencies
- **Crypto:** AES-256-GCM + Argon2id — non-negotiable; must use high-cost Argon2id parameters
- **Portability:** No files written outside the vault file path; no OS registry or config directory usage
- **Compatibility:** Windows, macOS, Linux — all three must work from day one
- **Privacy:** Zero application logs containing sensitive data (paths, names, values)
- **Reliability:** Atomic save (tmp → rename) + .bak backup on every save
- **Backward Compatibility:** App version N must open vaults created with version N-1

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Go as implementation language | Single static binary requirement; excellent cross-platform builds; strong stdlib crypto | — Pending |
| Bubble Tea + Lip Gloss for TUI | Modern, component-based; widely adopted; pairs naturally with Go | — Pending |
| AES-256-GCM + Argon2id | Industry-standard symmetric encryption + memory-hard KDF; strong against offline brute-force | — Pending |
| Vault format: encrypted JSON | Human-readable internals (after decryption), easy migration/export, flexible schema | — Pending |
| Version field inside encrypted payload | Simpler format; version is part of the trusted payload, not exposed metadata | — Pending |
| Trash = virtual folder, purge on save | No separate "empty trash" UX; save intent is clear enough; keeps model simple | — Pending |
| Templates stored inside vault | Each vault is self-contained; portability guaranteed; no shared config files | — Pending |
| Templates are snapshots at creation time | Avoids retroactive mutations of existing secrets when templates evolve | — Pending |
| Favorites + Trash as only virtual folders | Covers the primary access patterns (quick access + undo delete) without over-engineering | — Pending |
| Search is in-memory scan | Vault is fully decrypted in memory; no index needed; simpler and avoids data leakage | — Pending |

---
*Last updated: 2026-03-24 after initial project specification*
