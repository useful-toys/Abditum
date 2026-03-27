# Feature Landscape: TUI/CLI Password Managers

**Domain:** Offline, portable, TUI/CLI password manager
**Researched:** 2026-03-27
**Scope:** Informs Abditum v1 roadmap — what's table stakes, what differentiates, what to avoid

---

## Reference Landscape

Competitive products surveyed to calibrate expectations:

| Product | Model | Storage | Encryption | Generator |
|---------|-------|---------|------------|-----------|
| KeePass / KeePassXC | Desktop GUI + KeePass-CLI | Local file (.kdbx) | AES-256 + Argon2 | ✓ Built-in |
| pass (Unix password store) | CLI | Local files + Git | GPG | ✓ Via `pwgen` / `openssl` |
| gopass | CLI (enhanced pass) | Local files + Git | GPG | ✓ Built-in |
| Bitwarden CLI (`bw`) | CLI wrapper | Cloud vault | AES-256-CBC | ✓ Built-in |
| 1Password CLI (`op`) | CLI wrapper | Cloud vault | AES-256 | ✓ Built-in |
| passpy / tpm | TUI | Local file | varies | ✓ Usually present |
| age / rage | Encryption tool | Local file | ChaCha20-Poly1305 | ✗ (not a PM) |

---

## Table Stakes

Features users expect from any password manager. Missing = product feels incomplete or untrustworthy.

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Encrypted vault storage | Core contract — data must be protected at rest | Low | AES-256 is the baseline expectation |
| Master password authentication | The only door to the vault | Low | Single credential, no alternatives |
| Create / Edit / Delete secrets | Core CRUD | Low | No manager skips this |
| Multiple field types per secret | Username + password minimum; users need structured storage | Low | Plain text fields are enough at start |
| Mask sensitive fields by default | Passwords should not be visible on screen at all times | Low | Essential for shoulder-surfing environments |
| Temporary field reveal (show/hide) | Users need to see a password to type it elsewhere | Low | Auto-hide timer strongly expected |
| Copy field to clipboard | Primary access pattern — most users never type passwords | Low | Universal in every PM |
| Clipboard auto-clear | Clipboard persists across apps; must not linger | Low | 30s is widely accepted default |
| Search / find secrets | Any vault with >20 entries is unusable without search | Medium | Substring, case-insensitive minimum |
| Organize / group secrets | Users mentally categorize (work, personal, banking…) | Low | Flat categories or folders both acceptable |
| Auto-lock on inactivity | Unattended session must close | Low | 5–15 min is typical default |
| Manual lock | User must be able to lock immediately | Low | Single keystroke expected in TUI |
| Keyboard-only navigation | Non-negotiable for any TUI | Medium | Mouse = optional; keyboard = mandatory |
| Export vault | Users must be able to migrate away | Medium | JSON or CSV minimum; no lock-in |
| Import vault | Adoption blocker if missing | Medium | At least the app's own format minimum |
| Memory protection (no sensitive logs) | Passwords must never appear in stdout/stderr/log files | Low | Security hygiene baseline |
| Master password strength feedback | Inform user of weak choices; advisory, not blocking | Low | Visual indicator at minimum |
| **Password generator** | **Every CLI and TUI PM has this. Users expect to generate, not just store.** | **Medium** | **Abditum DEFERS to v2 — notable gap** |

### On Password Generator

Every competitor in the table above ships a password generator. `pass` uses external tools (`pwgen`, `openssl rand`). KeePassXC has a full generator UI with length, charset, and entropy controls. gopass ships one built-in.

Users who start using a new password manager and can't generate a new password directly inside it will immediately reach for their browser's generator or a third-party tool — which undermines the "single tool" experience. This is the **single most significant v1 gap** against the class of TUI password managers.

**Recommendation:** Treat the password generator as P0 for v2; design the field-edit flow in v1 so the generator can be inserted without structural changes (e.g., a "Generate" action on any sensitive field during edit).

---

## Differentiators

Features that set Abditum apart. Not universally expected, but clearly valued by the target user.

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| **User-defined secret templates** | KeePass has preset entry types; nobody has user-defined field schemas with type control (common vs sensitive) | High | Core design choice; enables any credential kind |
| **Hierarchical folder tree (unlimited depth)** | Most CLI/TUI PMs use flat tags or one-level groups | Medium | Full tree with move/cycle-detection is non-trivial |
| **Soft-delete (mark-for-exclusion)** | Normal PMs delete immediately; soft-delete gives an "undo window" before commit | Low | Unique behavior for a CLI tool; reduces panic deletes |
| **Atomic save with .tmp/.bak protocol** | Most personal tools have no crash-safe write; Abditum behaves like a database | Medium | Differentiates trust and reliability |
| **Single portable binary — truly zero install** | pass needs GPG + git; KeePassXC needs Qt runtime; Abditum needs nothing | Low (design) | Go binary design; huge DX advantage |
| **Offline-only with contractual guarantee** | No cloud = zero attack surface outside the file; no account, no breach notification emails | Low (policy) | Trust through simplicity |
| **Template creation from existing secret** | Workflow shortcut: retroactively formalize a schema from live data | Low | No competitor has this |
| **File format backward-compatibility guarantee** | Professional guarantee; most TUI tools don't document this | Medium | Requires migration layer in code |
| **External file change detection before save** | Prevents silent data loss if vault file is on a network share or synced folder | Low | Smart safety net |
| **Clear screen on lock/exit** | Prevents residual data in terminal scrollback | Low | Thoughtful privacy detail; most tools skip this |
| **mlock / VirtualLock best-effort** | Reduces swap exposure of sensitive buffers | Medium | Platform-sensitive; honest about limits |
| **Observation field always present** | Universal notes field; users always have a place to add free-form context | Low | Simple but prevents "where do I put notes?" friction |

---

## Anti-Features

Features to explicitly NOT build. Each has been evaluated and rejected for principled reasons — not laziness.

| Anti-Feature | Why Avoid | What to Do Instead |
|--------------|-----------|-------------------|
| **TOTP code generation** | Changes the app's trust model: now it's also a 2FA authenticator, requires time-sync, complicates key material. Out of scope permanently. | Recommend Ente Auth, Raivo, or Aegis for TOTP |
| **Cloud sync / remote storage** | Contradicts the offline contract. Any network access expands attack surface and requires a server. Out permanently. | Users may optionally sync the `.abditum` file via their own Syncthing/Dropbox — app is agnostic |
| **Browser extension / auto-fill** | Requires browser-to-app IPC, extension privileges, browser-specific builds. App is a portable binary by design. | Clipboard copy is the bridge |
| **Hardware key authentication (YubiKey)** | Requires USB stack, changes key-derivation model, makes file non-portable without the token. Out permanently. | Single master password + strong Argon2id params |
| **Multi-vault (simultaneous open)** | Architectural invariant — one vault per session. Adding multi-vault requires multiplexed unlock, session isolation, UI chrome. | `save as` to switch vaults; restart app |
| **Mobile / web app** | TUI portability is the product. Mobile UX requires a completely different stack (Flutter, React Native, etc.) | Clipboard and QR code (v2) are the transfer surface |
| **Team sharing / collaboration** | Personal vault design. Sharing would require ACL, E2E key exchange, server infrastructure. | Use a team PM (Bitwarden, 1Password Teams) |
| **Automatic cloud backup management** | Adds dependency on cloud APIs, credentials for backups, complicates trust model. | Users own their backups; `.abditum.bak` is the safety net |
| **Partial data recovery from corrupt vault** | AES-256-GCM is all-or-nothing by design; partial decryption is not possible. Would create false confidence. | Honest error message; user keeps `.abditum.bak` |
| **Configurable Argon2 parameters UI** | Exposing KDF tuning to users leads to accidentally weak configurations. | Ship sane defaults; update with new binary versions |

---

## Feature Dependencies

```
Master password auth
  └── Session state tracking
        ├── Auto-lock timer (requires "last interaction" timestamp)
        ├── Clipboard auto-clear timer
        └── Field reveal timer

Sensitive field type system
  ├── Field masking (default hidden)
  ├── Temporary reveal → Field reveal timer
  └── Copy to clipboard → Clipboard auto-clear

Template system
  ├── Secret creation (from template or bare)
  ├── Duplicate secret (copies template lineage)
  └── Create template from existing secret

Hierarchical folder tree
  ├── Move secret → any folder
  ├── Move folder → cycle validation
  ├── Search (must traverse full tree)
  └── Export / Import (must serialize entire tree)

Soft-delete
  └── Save / Save As (permanent removal on persist)

Atomic save protocol (.tmp / .bak)
  ├── Save (overwrite)
  ├── Save As (new file)
  └── Master password change (immediate, irrevocable)

External change detection
  └── Save → check file mtime before write

Export (JSON)
  ├── JSON schema stability (must match Import)
  └── Backward compatibility layer

Import (JSON)
  ├── Merge logic (folder path deduplication)
  ├── ID conflict resolution (secrets + templates)
  └── Backward compatibility layer

Memory wipe
  ├── Lock (manual and auto)
  └── Exit (clean + dirty)

Clear screen
  ├── Lock
  └── Exit

[v2] Password generator
  └── Field edit flow (sensitive field context action)
      └── No structural dependency on v1 — can be inserted cleanly

[v2] Duress password
  └── Master password system (alternative decryption path)
      └── Separate from all v1 features; isolated upgrade

[v2] Secret version history
  └── Soft-delete (similar "marked but retained" concept)

[v2] Vault health report
  └── Full vault traversal + sensitive field access (read-only audit)

[v2] Tags
  └── Search (must filter by tag)
```

---

## Abditum v1 Coverage Assessment

### Fully Covered

| Feature | Requirement IDs | Notes |
|---------|----------------|-------|
| Encrypted vault (AES-256-GCM + Argon2id) | SEC-CRYPTO-01 | Strong algorithm choices |
| Master password auth with strength feedback | VAULT-01, VAULT-09 | Advisory; non-blocking |
| Create / Edit / Delete secrets | SEC-01, SEC-03, SEC-07 | Full CRUD |
| Custom fields (common + sensitive) per secret | SEC-04, TPL-01–05 | Unique field system |
| Field masking by default | QUERY-04 | Sensitive fields hidden |
| Temporary field reveal with auto-hide | QUERY-04 | Configurable timer |
| Copy any field to clipboard + auto-clear | QUERY-05 | Configurable timer |
| Search (name, field name, common value, note) | QUERY-02 | Sensitive fields excluded from search |
| Hierarchical folder tree (unlimited depth) | FOLDER-01–05 | Cycle detection, move, reorder |
| Auto-lock on inactivity + manual lock | VAULT-10, VAULT-11 | Configurable timer |
| Memory wipe on lock/exit | VAULT-12, MEM-01 | mlock/VirtualLock best-effort |
| Clear screen on lock/exit | VAULT-12 | Visual privacy |
| Save / Save As | VAULT-05, VAULT-06 | External change detection |
| External file change detection | VAULT-07 | Overwrite / Save As / Cancel |
| Discard and reload | VAULT-08 | Session rollback |
| Master password change (immediate) | VAULT-09 | Irrevocable, saves immediately |
| Export to JSON + security warning | VAULT-14 | Soft-deleted secrets excluded |
| Import from JSON + merge + conflict handling | VAULT-15 | ID dedup, folder merge |
| Configurable timers (3x) | VAULT-16 | Lock, reveal, clipboard |
| Soft-delete (mark for exclusion) | SEC-07 | Permanent on save |
| Favorites | SEC-06 | Visual priority |
| Duplicate secret | SEC-02 | Same folder, adjusted name |
| Move / Reorder secret | SEC-08, SEC-09 | Manual ordering persisted |
| Create / rename / move / reorder / delete folders | FOLDER-01–05 | Full folder lifecycle |
| User-defined templates (create, edit, delete) | TPL-01–04 | Alteration does not affect existing secrets |
| Create template from existing secret | TPL-05 | Observation excluded |
| Zero sensitive data in logs | SEC-PRIV-01 | No paths, names, or values in stdout/stderr |
| Backward-compatible file format | COMPAT-01 | In-memory migration, save in current format |
| Atomic save (.tmp / .bak / .bak2) | ATOMIC-01 | Full rollback protocol |
| Single portable binary, cross-platform | PORT-01 | Windows, macOS, Linux |
| CI (build + lint + tests on every push) | CI-01 | Mandatory gate |
| Observation field (always present, last position) | SEC-05, QUERY-03 | Non-removable, non-renameable |

### Notable V1 Gaps

| Missing Feature | Gap Severity | Notes |
|----------------|-------------|-------|
| **Password generator** | **HIGH — table stakes** | Every TUI/CLI PM has this. Users creating a new account during onboarding cannot generate a password without leaving the app. Deferred to v2. Design field-edit UX now to accommodate a "Generate" action on sensitive fields with zero refactor cost later. |
| **Keyboard help overlay / cheatsheet** | MEDIUM — TUI usability | requisitos.md does not mention any in-app help. TUI discoverability depends entirely on visible key hints (footer bar, `?` overlay). Without this, onboarding is hostile. Not a v2 item — should be in v1 UX. |
| **Import from competitor formats** | MEDIUM — adoption blocker | Only JSON (Abditum's own format) is importable. Users migrating from KeePass (.kdbx XML export), Bitwarden (CSV), LastPass (CSV), or 1Password (1PUX) must manually re-enter all credentials. A KeePass XML or Bitwarden CSV importer would dramatically lower adoption friction. |
| **Password strength indicator on stored passwords** | LOW | Strength feedback only appears at master password creation/change. Users cannot see which stored passwords are weak without a health report (v2). |
| **Partial field reveal (e.g., last 4 digits)** | LOW | Already acknowledged in v2. Useful for card numbers. |

### V2 Roadmap Items (Already Planned)

| Feature | Priority Signal | Notes |
|---------|----------------|-------|
| Password generator | P0 — table stakes gap | Design v1 field-edit flow to accommodate this cleanly |
| Duress password | P1 — security differentiator | Complex design; storage + validation strategy TBD |
| Vault health report | P2 — trust builder | Needs full vault read + weak/reused analysis |
| Secret version history | P2 — safety net | "Mark for exclusion" is a conceptual predecessor |
| Tags | P2 — discoverability | Must integrate with search |
| QR code field sharing | P3 — niche but clever | ASCII block QR in TUI; offline transfer surface |
| Partial field reveal | P3 — UX refinement | Per-field masking rules |
| Orphan artifact recovery | P3 — edge case | .abditum.tmp / .bak2 detection at open time |

---

## MVP Recommendation

**Build all table-stakes features except password generator for v1.** The generator gap is real but manageable: users coming from browsers or external generators can paste values; the clipboard copy flow covers retrieval.

**Critical v1 additions not in requisitos.md:**
1. **In-app keyboard help** (footer hints + `?` modal) — without this, TUI discoverability is broken regardless of how good the feature set is.
2. **Design field-edit UX with a generator slot** — even if the generator is not in v1, leave a "Generate →" placeholder action on sensitive field inputs so v2 can wire it up without redesigning the form.

**Do not add to v1:** import from competitor formats is the highest-value post-v1 addition but adds significant scope (format parsers + mapping logic). Ship v1 with JSON-only import and plan a v2 "import wizard" phase.

---

## Sources

- KeePass documentation: https://keepass.info
- pass (Unix password store): https://www.passwordstore.org
- gopass: https://github.com/gopasspw/gopass
- Bitwarden CLI: https://bitwarden.com/help/cli/
- requisitos.md (Abditum requirements — primary source)
- .planning/PROJECT.md (Abditum project definition — primary source)
- Confidence: HIGH for table stakes (cross-validated across multiple products); HIGH for Abditum feature coverage (direct from requirements document); MEDIUM for gap severity assessments (based on observed patterns in similar tools)
