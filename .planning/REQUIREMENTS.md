# Requirements: Abditum

**Defined:** 2026-03-24
**Revised:** 2026-03-24 (revision 2 — see change summary below); 2026-03-24 (revision 3 — see change summary below); 2026-03-24 (revision 4 — see change summary below)
**Core Value:** A single portable binary that any user can carry on a USB drive and use on any machine — opening, managing, and saving an encrypted vault without installing anything or touching the cloud.

---

## Change Summary (revision 2026-03-24)

Key changes from the original specification (revision 2):

- **Lock behavior** — refined from "zero all sensitive data" to "minimize retention and clear app-controlled buffers where possible" (honest about Go GC limitations)
- **Import conflict rules** — expanded: folders merge silently; secrets get incremental suffix + user notification; templates are replaced silently by imported version
- **Folder delete** — folders are permanently removed (not soft-deleted); secrets AND subfolders promoted to parent; promoted items appended at end of parent lists
- **Search scope** — added "field name" as a search criterion; removed "containing folder name" as a search criterion; `texto sensível` fields explicitly excluded
- **Master password confirmation** — double-entry required on create and change (now explicit)
- **Security framing** — "zero memory" replaced with "minimize retention, clear controlled buffers"
- **Portability wording** — tightened: `.tmp` and `.bak` are explicitly named exceptions to the "no data outside vault" rule
- **Backward compatibility** — expanded: N opens any previously supported version (not just N-1); saves always in current format; migration happens in-memory
- **Name uniqueness** — new explicit NFR: duplicate names are permitted for secrets, folders, and templates
- **Data model** — formal entity definitions added with NanoID (6 chars), field types (`texto` / `texto sensível`), implicit observation field, recursive folder hierarchy
- **File format** — header revised: `magic (ABDT) + versão_formato + salt + nonce`; Argon2id params removed from header; entire header used as GCM AAD
- **Trash** — applies to secrets only; folders are permanently removed when deleted

Key changes from revision 2 (revision 3):

- **Secret creation flow** — removed "predefined vs custom template" distinction; user now creates from any existing template OR starts with a blank secret (no initial fields); SEC-01 updated accordingly
- **Secret edit modes renamed** — "basic mode" and "advanced edit mode" eliminated; replaced with "edit secret data" (name, field values, observation) and "alter secret structure" (add, remove, rename, reorder fields); secret name editing now explicit in SEC-05
- **Secret field type change → Out of Scope** — changing a field's type on an existing secret is not supported; workaround is delete + re-add; SEC-06 updated; new Out of Scope entry added
- **Template field type change remains supported** — TMPL-02 retains "retype" and now explicitly notes the asymmetry: field type change is allowed for templates but not for existing secrets

Key changes from revision 3 (revision 4):

- **Import secret conflict — two distinct rules** — identity (ID) collision and name collision are now separate cases with different resolutions: ID collision → assigned new ID (data preserved, no suffix); name collision in same destination folder → name suffixed + user notified; VAULT-09 updated accordingly
- **Backup rotation (.bak2)** — precise protocol added: existing `.bak` renamed to `.bak2` before new backup; deleted on success; restored on failure; VAULT-04 updated; Reliability constraint updated; new Key Decision added
- **Re-authentication to save → Out of Scope** — vault is already authenticated in session; master password not required again to save; added to Out of Scope in both docs
- **Root-as-folder (MODEL-08)** — vault root is an implicit nameless folder; the data model is fully recursive with no structural special-casing; new requirement added
- **JSON order = display order (MODEL-09)** — array position in JSON is the canonical display order for secrets, folders, and fields; no separate position field; new requirement added
- **Settings defaults explicit** — auto-lock 2 min, field reveal 15 s, clipboard 30 s now named in MODEL-01

---

## v1 Requirements

### Vault Lifecycle

- [ ] **VAULT-01**: User can create a new vault by specifying a file path and master password; master password requires double-entry confirmation
- [ ] **VAULT-02**: User can open an existing `.abditum` vault file with its master password
- [ ] **VAULT-03**: User can save the vault atomically (write to `.abditum.tmp` → rename on success; delete `.abditum.tmp` on failure to avoid leaving transitory data)
- [ ] **VAULT-04**: On every save, the previous vault file is preserved as `.abditum.bak`; before writing the new backup, any existing `.abditum.bak` is renamed to `.abditum.bak2`; on successful save `.abditum.bak2` is deleted; if the operation fails, `.abditum.bak2` is restored to `.abditum.bak` where possible
- [ ] **VAULT-05**: User can save the vault to a new path (Save As)
- [ ] **VAULT-06**: User can discard unsaved changes and reload the vault from disk
- [ ] **VAULT-07**: User can change the master password of an open vault; new password requires double-entry confirmation
- [ ] **VAULT-08**: User can export the vault to a plain-text JSON file (with security risk warning and explicit confirmation before proceeding)
- [ ] **VAULT-09**: User can import from a plain-text JSON file with the following conflict rules:
  - Folders with the same identity as an existing folder → hierarchy merged silently; if folder conflict occurs, merge happens silently
  - Secrets with the same identity (ID) as an existing secret → the imported secret is assigned a new identity (new ID); its data is fully preserved
  - Secrets with the same name as an existing secret in the same destination folder → the imported secret's name is suffixed with an incremental number (e.g., "Secret (1)", "Secret (2)"); user is notified of all name-suffixed imports
  - Templates with the same identity as an existing template → imported template replaces the existing one silently
- [ ] **VAULT-10**: User can configure vault settings (stored inside vault file): auto-lock timeout, field reveal timeout, clipboard clear timeout

### Authentication

- [ ] **AUTH-01**: User can unlock a vault with its master password
- [ ] **AUTH-02**: The vault locks automatically after a configurable inactivity timeout; locking returns the user to the vault open screen and clears app-controlled sensitive buffers where possible
- [ ] **AUTH-03**: User can lock the vault manually; same behavior as auto-lock
- [ ] **AUTH-04**: Key derivation uses Argon2id with high-cost parameters (≥64 MiB memory, time=3) to resist brute-force and offline attacks
- [ ] **AUTH-05**: A spinner/progress indicator is shown during Argon2id key derivation

### Vault Navigation

- [ ] **NAV-01**: User can view the full vault hierarchy (folders, subfolders, secrets) in a sidebar tree
- [ ] **NAV-02**: User can view the details of a selected secret in a detail panel
- [ ] **NAV-03**: User can temporarily reveal a sensitive field value; the value auto-hides after a configurable timeout (suggested default 15 s)
- [ ] **NAV-04**: A "Favorites" virtual folder surfaces all favorited secrets regardless of their location in the hierarchy
- [ ] **NAV-05**: A "Trash" virtual folder (Exclusão Reversível) surfaces all soft-deleted secrets; they are restorable until the next save, at which point they are permanently purged

### Secret Management

- [ ] **SEC-01**: User can create a secret from any existing template, OR start with a blank secret with no initial fields; the template is a snapshot at creation time — changes to or deletion of the template do not affect existing secrets; the template name is stored as a historical record only
- [ ] **SEC-02**: User can create a secret at the vault root or inside any folder
- [ ] **SEC-03**: User can duplicate an existing secret
- [ ] **SEC-04**: User can favorite or unfavorite a secret
- [ ] **SEC-05**: User can edit a secret's data: change the secret name, field values, and observation
- [ ] **SEC-06**: User can alter a secret's structure: add, remove, rename, and reorder fields on an existing secret; changing a field's type is not supported — the user must delete the field and add a new one with the desired type instead
- [ ] **SEC-07**: User can soft-delete a secret; the secret moves to the Trash virtual folder and remains restorable until the next save
- [ ] **SEC-08**: User can restore a soft-deleted secret from Trash before the next save
- [ ] **SEC-09**: User can move a secret to another folder or to the vault root
- [ ] **SEC-10**: User can reorder a secret relative to other items in the same parent
- [ ] **SEC-11**: User can search secrets by: secret name, field name, value of `texto`-type fields, or observation (note); `texto sensível` field values are never included in search results; all secrets matching any criterion are shown

### Folder Management

- [ ] **FOLD-01**: User can create a folder at the vault root or inside another folder
- [ ] **FOLD-02**: User can rename a folder
- [ ] **FOLD-03**: User can move a folder to another folder or to the vault root
- [ ] **FOLD-04**: User can reorder a folder relative to other items in the same parent
- [ ] **FOLD-05**: User can delete a folder; its direct secrets and direct subfolders are promoted to the parent (or vault root if the deleted folder was at root); promoted secrets are appended to the end of the parent's secret list; promoted subfolders are appended to the end of the parent's folder list; the folder itself is permanently removed (not soft-deleted)
- [ ] **FOLD-06**: A new vault is pre-populated with three editable/removable folders: "Sites", "Financeiro", "Serviços"

### Template Management

- [ ] **TMPL-01**: User can create a custom secret template with named, typed fields
- [ ] **TMPL-02**: User can edit a template (add, remove, rename, retype, reorder fields); changes only affect future secrets created from that template; existing secrets are unchanged; note: changing a field's type is supported for templates but not for existing secrets
- [ ] **TMPL-03**: User can delete a template; existing secrets created from it are unaffected
- [ ] **TMPL-04**: User can create a template from an existing secret (copies field names and types as the initial structure; field values are not copied)
- [ ] **TMPL-05**: A new vault is pre-populated with three editable/removable predefined templates: "Login" (URL, Username, Password), "Cartão de Crédito" (Número do Cartão, Nome no Cartão, Data de Validade, CVV), "API Key" (Nome da API, Chave de API)
- [ ] **TMPL-06**: Templates are stored inside the vault file; each vault has its own independent set of templates

### Clipboard

- [ ] **CLIP-01**: User can copy any field value to the system clipboard
- [ ] **CLIP-02**: The clipboard is automatically cleared after a configurable timeout (suggested default 30 s)
- [ ] **CLIP-03**: The clipboard is cleared when the vault is locked or the app is closed

### Security

- [ ] **SEC-A-01**: On lock or close, the application clears app-controlled buffers holding sensitive data (master password input, derived key, decrypted payload) where possible; full zeroing is best-effort given Go's GC
- [ ] **SEC-A-02**: A hotkey instantly hides the entire TUI (shoulder-surfing protection); the same hotkey restores it
- [ ] **SEC-A-03**: Exporting to plain-text JSON shows a security risk warning and requires explicit user confirmation before proceeding
- [ ] **SEC-A-04**: The application produces no log output (stdout/stderr) containing vault file paths, secret names, or field values

### Data Model

- [ ] **MODEL-01**: Vault payload structure: settings (auto-lock timeout, suggested default 2 min; field reveal timeout, suggested default 15 s; clipboard clear timeout, suggested default 30 s), root secrets list, root folders list, secret templates list, creation date, last-modified date
- [ ] **MODEL-02**: Secret entity: id (NanoID 6 alphanumeric chars), name, template name (optional, historical), fields list, favorite flag, observation (implicit non-sensitive free-text field, always present, cannot be removed), creation date, last-modified date
- [ ] **MODEL-03**: Folder entity: id (NanoID 6 alphanumeric chars), name, secrets list, subfolders list (recursive — folder hierarchy has no depth limit)
- [ ] **MODEL-04**: Template entity: id (NanoID 6 alphanumeric chars), name, template-fields list
- [ ] **MODEL-05**: Field types: `texto` (plain text, included in search) and `texto sensível` (sensitive text, masked in UI, excluded from search); field value may be an empty string
- [ ] **MODEL-06**: Name uniqueness is not enforced: duplicate names are permitted for secrets, folders, and templates; identity is determined by ID, not name
- [ ] **MODEL-07**: The observation field is implicit on every secret, not declared in templates, non-removable, and treated as non-sensitive data
- [ ] **MODEL-08**: The vault root is a nameless implicit folder; the recursive data model is uniform — the root contains a secrets list and a folders list, identical in structure to any named folder
- [ ] **MODEL-09**: The order of elements in JSON arrays is the canonical display order; there is no separate position or order field for secrets, folders, or fields

### Format & Compatibility

- [ ] **FMT-01**: Vault files use the `.abditum` extension
- [ ] **FMT-02**: Binary file format: fixed header (`magic=ABDT` bytes + `versão_formato` integer + `salt` bytes + `nonce` bytes) followed by AES-256-GCM ciphertext of the JSON payload; the entire header is used as GCM Additional Authenticated Data (AAD)
- [ ] **FMT-03**: Key derivation uses Argon2id with high-cost parameters; parameters are not stored in the file (fixed by the application)
- [ ] **FMT-04**: App version N can open vaults created with any previously supported format version; the decrypted payload is migrated in-memory to the current model; on save, the vault is always written in the current format version
- [ ] **FMT-05**: No data is persisted outside the vault file path, except `.abditum.tmp` (transitory, deleted after successful save) and `.abditum.bak` (explicit backup)

### Build & Release

- [ ] **REL-01**: Binaries are produced for 6 targets: windows/amd64, windows/arm64, darwin/amd64, darwin/arm64, linux/amd64, linux/arm64
- [ ] **REL-02**: All binaries are built with `CGO_ENABLED=0` (no C dependencies)
- [ ] **REL-03**: CI runs tests on Ubuntu, Windows, and macOS runners
- [ ] **REL-04**: Release pipeline produces binaries via goreleaser on version tag push

---

## v2 Requirements

### TOTP / 2FA

- **V2-TOTP-01**: A field type of "TOTP seed" generates a live 6-digit code with countdown indicator
- **V2-TOTP-02**: TOTP codes are calculated locally, offline, with no network access

### Password Generator

- **V2-GEN-01**: User can generate a random password from within a field editor, with configurable length and character sets

### Duress Password

- **V2-DURE-01**: User can configure a secondary "duress" master password that opens a restricted or decoy vault

### QR Code Sharing

- **V2-QR-01**: User can render a field value as a QR code directly in the TUI using ASCII blocks

### Vault Health Audit

- **V2-AUDIT-01**: User can run a vault health report that identifies weak, reused, or old passwords

### Hardware 2FA

- **V2-HW2FA-01**: User can require a keyfile or hardware token (e.g., YubiKey) in addition to the master password

---

## Out of Scope

| Feature | Reason |
|---------|--------|
| Cloud sync / remote storage | Zero-knowledge by design; portability is the core value |
| Multiple vaults open simultaneously | Single-vault model simplifies UX and security surface |
| Mobile or web app | The portable TUI is the product; mobile is an architectural pivot |
| Tags on secrets | Folder hierarchy covers organization needs for v1 |
| Secret version history | Significant storage and UX complexity; deferred to v2 |
| Changing the type of an existing secret field | Delete the field and add a new one with the desired type instead |
| Re-authentication to save | The vault is already unlocked and authenticated in the session; the master password is not required again to save |

---

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| VAULT-01 | Phase 1 (storage), Phase 3 (TUI) | Pending |
| VAULT-02 | Phase 1 (storage), Phase 3 (TUI) | Pending |
| VAULT-03 | Phase 1 (storage) | Pending |
| VAULT-04 | Phase 1 (storage) | Pending |
| VAULT-05 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| VAULT-06 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| VAULT-07 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| VAULT-08 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| VAULT-09 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| VAULT-10 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| AUTH-01 | Phase 1 (crypto), Phase 3 (TUI) | Pending |
| AUTH-02 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| AUTH-03 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| AUTH-04 | Phase 1 (crypto) | Pending |
| AUTH-05 | Phase 3 (TUI) | Pending |
| NAV-01 | Phase 4 (TUI) | Pending |
| NAV-02 | Phase 4 (TUI) | Pending |
| NAV-03 | Phase 4 (TUI) | Pending |
| NAV-04 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| NAV-05 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| SEC-01 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| SEC-02 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| SEC-03 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| SEC-04 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| SEC-05 | Phase 4 (TUI) | Pending |
| SEC-06 | Phase 4 (TUI) | Pending |
| SEC-07 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| SEC-08 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| SEC-09 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| SEC-10 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| SEC-11 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| FOLD-01 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| FOLD-02 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| FOLD-03 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| FOLD-04 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| FOLD-05 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| FOLD-06 | Phase 2 (domain) | Pending |
| TMPL-01 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| TMPL-02 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| TMPL-03 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| TMPL-04 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| TMPL-05 | Phase 2 (domain) | Pending |
| TMPL-06 | Phase 1 (format), Phase 2 (domain) | Pending |
| CLIP-01 | Phase 4 (TUI) | Pending |
| CLIP-02 | Phase 4 (TUI) | Pending |
| CLIP-03 | Phase 4 (TUI) | Pending |
| SEC-A-01 | Phase 2 (domain), Phase 4 (TUI) | Pending |
| SEC-A-02 | Phase 4 (TUI) | Pending |
| SEC-A-03 | Phase 4 (TUI) | Pending |
| SEC-A-04 | Phase 5 (audit) | Pending |
| MODEL-01 | Phase 1 (format), Phase 2 (domain) | Pending |
| MODEL-02 | Phase 1 (format), Phase 2 (domain) | Pending |
| MODEL-03 | Phase 1 (format), Phase 2 (domain) | Pending |
| MODEL-04 | Phase 1 (format), Phase 2 (domain) | Pending |
| MODEL-05 | Phase 1 (format), Phase 2 (domain) | Pending |
| MODEL-06 | Phase 2 (domain) | Pending |
| MODEL-07 | Phase 2 (domain) | Pending |
| MODEL-08 | Phase 1 (format), Phase 2 (domain) | Pending |
| MODEL-09 | Phase 1 (format), Phase 2 (domain) | Pending |
| FMT-01 | Phase 1 (storage) | Pending |
| FMT-02 | Phase 1 (crypto + storage) | Pending |
| FMT-03 | Phase 1 (crypto) | Pending |
| FMT-04 | Phase 1 (format), Phase 2 (domain) | Pending |
| FMT-05 | Phase 1 (storage), Phase 5 (audit) | Pending |
| REL-01 | Phase 5 | Pending |
| REL-02 | Phase 5 | Pending |
| REL-03 | Phase 5 | Pending |
| REL-04 | Phase 5 | Pending |

**Coverage:**
- v1 requirements: 66 total
- Mapped to phases: 66
- Unmapped: 0 ✓

---
*Requirements defined: 2026-03-24*
*Last updated: 2026-03-24 after 4th requirement revision (import two-level conflict rules, .bak2 rotation, re-auth-to-save out of scope, root-as-folder, JSON-order, settings defaults)*
