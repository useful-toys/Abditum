# Abditum — Project Roadmap

**Last updated:** 2026-03-24  
**Status:** Research complete. Ready for Phase 1.

---

## Overview

Abditum is a portable, offline TUI password vault: a single Go binary that encrypts secrets with AES-256-GCM + Argon2id and stores them in a `.abditum` file. No installation, no cloud, no traces.

The roadmap is split into **5 phases**, each independently shippable:

| Phase | Name | Goal | Key deliverable |
|-------|------|------|-----------------|
| 1 | Foundation | Project scaffold + crypto + storage | Can create and open a vault file |
| 2 | Domain | Domain model + business logic | Complete vault manager API, all CRUD, search |
| 3 | Auth TUI | Authentication screen | Unlock / create vault via TUI |
| 4 | Main TUI | Full interactive interface | Navigate, view, edit secrets in the terminal |
| 5 | Polish & Release | Build pipeline, tests, UX hardening | Cross-platform binaries, CI, full test coverage |

---

## Phase 1 — Foundation

**Goal:** Set up the project skeleton, crypto layer, and storage layer. No TUI. Fully tested in isolation.

**Deliverables:**
- `go.mod` / `go.sum` with all dependencies
- `cmd/abditum/main.go` — minimal entry point (no TUI yet, just smoke test)
- `internal/crypto/` — Argon2id KDF + AES-256-GCM encrypt/decrypt
- `internal/storage/` — binary file format (header + ciphertext), atomic save, `.bak` backup
- Unit tests for crypto (success + failure cases) and storage (success + failure cases)

**Why first:** Everything else depends on being able to encrypt and persist a vault. This is the bedrock. Can be built and tested with zero UI.

**Done when:** `go test ./...` passes and a vault file can be written and read back from disk with the correct binary format.

---

## Phase 2 — Domain Model

**Goal:** Implement the full vault domain: entities, VaultManager API, virtual folders, search, import/export. No I/O, no TUI.

**Deliverables:**
- `domain/vault/` package — all entities (Vault, Folder, Secret, Field, SecretTemplate)
- `VaultManager` — full CRUD API: secrets, folders, templates, settings
- Virtual folder collectors: `CollectFavorites()`, `CollectTrashed()`
- `Search()` — in-memory scan (name, non-sensitive fields, note, folder name)
- `InitializeNewVault()` — seeds default templates + folders
- `ExportVault()` / `ImportVault()` — plain-text JSON round-trip
- DDD pattern enforced: all mutations through VaultManager, entities read-only
- Unit tests for every manager operation, search, virtual folder collection, import/export collision resolution

**Why second:** The domain is the core business logic. It can be developed and fully tested without crypto, storage, or TUI. Keeping it pure makes it easy to verify correctness.

**Done when:** All VaultManager operations are tested, including edge cases (collision import, restore-from-missing-folder, trash purge on save).

---

## Phase 3 — Authentication TUI

**Goal:** A working terminal UI for the authentication flow. User can create a new vault or open an existing one.

**Deliverables:**
- `tui/auth/` — authentication screen model (create vault / open vault)
- File picker integration for choosing vault path
- Password input with confirmation (create) and single-entry (open)
- Zero-knowledge warning on vault creation
- Argon2id progress indicator (non-blocking spinner during key derivation)
- Wires Phase 1 (storage) and Phase 2 (domain) into a running binary: `abditum` can create and open a `.abditum` file
- Golden file tests for auth screen at 80×24
- Keyboard command tests for the auth flow

**Why third:** Authentication is the entry point — it's the smallest useful slice of TUI that exercises the crypto and storage layers end-to-end. A user can create and open a vault after this phase.

**Done when:** The binary starts, shows the auth screen, and successfully creates or opens a vault. Ctrl+Q quits cleanly.

---

## Phase 4 — Main TUI

**Goal:** Full interactive vault management interface. All functional requirements from `PROJECT.MD` are implemented.

**Deliverables:**
- `tui/main/` — two-panel layout (sidebar tree + detail/edit panel)
- Sidebar: recursive tree, expand/collapse folders, virtual folders (Favorites at top, Trash at bottom), search/filter mode with highlight, keyboard navigation, type-to-jump
- Detail panel: field display with sensitive masking, reveal-with-timeout, clipboard copy with countdown
- Edit mode: create/edit/delete secrets, add/remove/reorder fields (advanced mode), template picker
- Folder management: create, rename, move, reorder, delete
- Template management: create, edit, delete, create-from-secret
- Status bar: file path, vault status (new / saved / modified), secret count
- Help bar: context-sensitive, showing only currently-applicable bindings
- Toast notifications: success (green), error (red), warning (yellow), info (blue)
- Blocking confirmation dialogs for destructive actions
- Shoulder-surfing hotkey (instant hide)
- Auto-lock timer with imminent-lock warning
- Responsive layout: resize warnings, single-panel mode under 40 columns
- Save/Save As/Discard, Change Master Password, Export/Import with security warning
- Exit with unsaved-changes guard
- Full golden file tests for every screen state at 80×24
- Keyboard command tests for all flows

**Why fourth:** The auth phase provides the foundation; the main TUI is the product. This is the largest phase and the most complex. All UX spec from `descricao.md` §5-§7 is realized here.

**Done when:** All items in `PROJECT.MD` Active requirements are checked off. The full user journey (create vault → add secrets → organize → lock → reopen) works end-to-end.

---

## Phase 5 — Polish & Release

**Goal:** Cross-platform release pipeline, full CI, UX hardening, and documentation.

**Deliverables:**
- `Makefile` — local dev targets: `build`, `test`, `lint`, `cross-compile`
- `.goreleaser.yml` — builds all 6 targets (windows/amd64, windows/arm64, darwin/amd64, darwin/arm64, linux/amd64, linux/arm64) with `-ldflags="-s -w"` + `-trimpath`
- `.github/workflows/ci.yml` — test matrix: ubuntu-latest + windows-latest + macos-latest
- `.github/workflows/release.yml` — goreleaser on tag push
- Cross-compile verification CI job (all 6 targets, no execution)
- `gosec` linter integration
- `golangci-lint` configuration
- No CGO confirmation (`CGO_ENABLED=0` enforced in all build paths)
- `README.md` — usage, download, security model, macOS/Windows bypass instructions for unsigned binaries
- Integration test: full end-to-end flow (create vault → create secret → edit → lock → reopen → verify)
- Memory zeroing audit (all sensitive data zeroed on lock/close)
- No-log audit (confirm zero sensitive data in any stdout/stderr output)
- Version `v1.0.0` tag and release

**Why last:** Polish requires a complete product to polish. CI and release tooling are validated against the real build; doing this earlier would be premature.

**Done when:** `goreleaser release --snapshot` produces 6 valid binaries. CI passes on all 3 OS runners. `v1.0.0` is tagged.

---

## Phase Ordering Rationale

```
Phase 1 (crypto + storage)
  └─► Phase 2 (domain — depends on nothing except stdlib)
        └─► Phase 3 (auth TUI — wires 1 + 2 into a running program)
              └─► Phase 4 (main TUI — the full product)
                    └─► Phase 5 (polish + release)
```

- Phases 1 and 2 are **independent** and could be parallelized, but Phase 1 is simpler and Phase 2 is pure Go — sequential is cleaner and avoids integration confusion.
- Phase 3 is the **first vertical slice**: crypto + storage + domain + one TUI screen = runnable binary.
- Phase 4 is the **bulk of the work**. It is intentionally one phase because all panels, flows, and UX behaviors are deeply interconnected.
- Phase 5 is explicitly **last** — no release pipeline until the product is done.

---

## Requirements Coverage

All Active requirements from `PROJECT.MD` are covered:

| Requirement Group | Phase |
|-------------------|-------|
| Vault Lifecycle (create, open, save, backup, save-as, discard, change password, export, import, settings) | 1 (persistence), 2 (logic), 4 (UX) |
| Authentication (create, unlock, lock, auto-lock, Argon2id) | 1 (crypto), 3 (TUI) |
| Vault Navigation (hierarchy, details, reveal, virtual folders) | 2 (logic), 4 (TUI) |
| Secret Management (create, duplicate, favorite, edit, delete, restore, move, reorder, search) | 2 (logic), 4 (TUI) |
| Folder Management (create, rename, move, reorder, delete, pre-defined) | 2 (logic), 4 (TUI) |
| Template Management (create, edit, delete, from-secret, pre-defined) | 2 (logic), 4 (TUI) |
| Clipboard (copy, auto-clear, clear-on-lock) | 4 (TUI) |
| Security (zero memory, shoulder-surfing, export warning) | 2 (zero/purge), 4 (TUI) |
| Format & Compatibility (`.abditum`, version field, no logs) | 1 (format), 5 (audit) |
| Build & Release (6 targets, CI, goreleaser) | 5 |
