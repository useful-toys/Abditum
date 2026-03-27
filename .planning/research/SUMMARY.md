# Research Summary — Abditum

**Date:** 2026-03-27
**Project:** Abditum — Go TUI Password Manager

---

## Key Findings

### Stack

- **Go 1.26.1 + charm.land import paths are confirmed.** Charm migrated vanity domains in late 2025; use `charm.land/bubbletea/v2`, `charm.land/bubbles/v2`, `charm.land/lipgloss/v2` — the old `github.com/charmbracelet/*` paths redirect but are not canonical.
- **Bubble Tea v2 has breaking API differences from v1.** `View()` now returns `tea.View` (not `string`); key messages are `tea.KeyPressMsg` (not `tea.KeyMsg`); space is the string `"space"` not `" "`. These must be learned before writing any TUI code or the entire layer needs a rewrite.
- **Argon2id parameters: t=3, m=64 MiB, p=4 (RFC 9106 second set) as default.** The 2 GiB first set is too aggressive for low-end hardware and will freeze the UI. The chosen set gives ~200–500ms unlock on modern CPUs — acceptable.
- **Clipboard must use `github.com/atotto/clipboard`** (subprocess-based, CGO-free), not any CGO-linked library. `CGO_ENABLED=0` is a hard binary constraint (PORT-01). `golang.x/sys` handles mlock/VirtualLock on all three platforms without CGO.
- **teatest/v2** is confirmed at `github.com/charmbracelet/x/exp/teatest/v2` with `charm.land` import paths (updated 2026-03-14). Use `WithInitialTermSize(80, 24)` on all golden-file tests for stability.

### Table-Stakes Features

- **Encrypt + CRUD + master password + search** — baseline expectation; fully covered in v1 requirements.
- **Copy field → clipboard auto-clear and temporary field reveal** — universal in every competing tool; covered in QUERY-04/05.
- **Keyboard-only navigation** — non-negotiable for TUI; not explicitly called out in requisitos.md as a named requirement, but implied everywhere.
- **In-app keyboard help overlay** — not in requisitos.md but research flags it as a MEDIUM gap that belongs in v1, not v2. TUI discoverability collapses without a `?` help modal or footer key-hint bar.
- **Password generator** — the single largest table-stakes gap (every competing TUI/CLI PM ships one). Deferred to v2 by design, but v1 field-edit flow must reserve a "Generate →" slot on sensitive fields so v2 wires in cleanly with zero refactor.

### Architecture

- **Strict layering is the primary structural invariant.** `internal/tui` never imports `internal/crypto` or `internal/storage` — all mutations go through `vault.Manager`. Violating this at any point makes the system untestable and the Manager pattern pointless.
- **Build bottom-up: crypto → vault → storage → tui → cmd.** Each layer is independently testable and has no upward dependencies. This ordering must be respected — writing TUI code before the domain layer exists produces untestable UI.
- **Root TUI model uses a flat session-state enum, not a model stack.** Every screen transition explicitly resets the previous child model to zero value. Modals are a `nil`-able field in the root, not a separate push/pop stack. This prevents ghost-state and unbounded growth.
- **All timers (auto-lock, clipboard clear, field reveal) live on the root model, driven by a 1-second global tick.** Using `tea.Tick` inside the Bubble Tea event loop eliminates goroutine leaks and race conditions that `time.AfterFunc` would introduce.
- **Atomic save writes `.tmp` in the same directory as the vault file** (`filepath.Dir(vaultPath)`) — never in `os.TempDir()`. Cross-device rename (`EXDEV`) is a real failure mode on encrypted home directories and network drives.

### Watch Out For

- **`string` for sensitive data will burn you.** Bubble Tea's `textinput.Value()` returns `string`. Convert to `[]byte` inside the `Update` handler, zero the textinput buffer, and never pass passwords or keys as `string` deeper into the stack. Establish this as a compile-time convention before writing any vault code.
- **GCM nonce reuse is catastrophic.** Always call `io.ReadFull(rand.Reader, nonce)` immediately before each `gcm.Seal()`. Write a unit test that encrypts the same plaintext twice and asserts the ciphertexts differ.
- **Argon2id parameters must be stored in the file header, not hardcoded at decryption time.** Hardcoding makes parameter upgrades permanently break existing vaults. Store salt (32 bytes), time, memory, threads, and an `ArgonVersion` constant in the plaintext header before any vault files are created — format changes are breaking changes.
- **Terminal "clear screen" does not clear scrollback.** Only `\033[3J\033[2J\033[H` clears both. `\033[2J` alone or `exec.Command("clear")` leaves secrets visible on scroll-up. Test manually in xterm, iTerm2, and Windows Terminal before shipping.
- **Cross-platform clipboard clearing has silent failure modes.** `xclip` is X11-only; Wayland clipboard persists separately. Use `github.com/atotto/clipboard`, clear synchronously before process exit (not in a goroutine), and handle headless/SSH sessions gracefully (skip with warning, no crash).
- **Search must never touch sensitive field values.** The `QUERY-02` requirement explicitly excludes sensitive fields from search. Any loop over secret fields that reads `.Value` without a `field.IsSensitive` guard is a bug. Write a negative test from day one: searching for a string that appears only in a sensitive field must return zero results.
- **Atomic save's `.bak`/`.bak2` state machine needs a startup recovery scan.** If the process dies between `.bak` rename and `.bak2` cleanup, subsequent opens find an inconsistent state. Implement a recovery function that runs at startup and repairs orphaned `.bak2` files.
- **The mlock/VirtualLock quota is not guaranteed.** Always check the return value; treat failure as "best-effort" and continue. Log generically (no content). Never `append` to a locked buffer — pre-allocate exact sizes upfront. Build-tag Unix and Windows code paths separately.

---

## v1 Readiness Assessment

**The v1 scope is realistic and well-defined.** The requirements in PROJECT.md are specific, internally consistent, and mapped to a verified architecture. No major feature is architecturally underspecified.

**Two risks worth watching:**

1. **In-app help is absent from requisitos.md** but belongs in v1 on research evidence. A TUI with no discoverable key bindings will feel broken to new users. This is a small implementation surface (footer key-hint bar + `?` modal using `charm.land/bubbles/v2/help`) and should be folded into the TUI phase rather than deferred.

2. **Cross-platform clipboard and screen-clear are deceptively hard to test.** These must be verified on Linux X11, Linux Wayland, macOS, and Windows before v1 ships — not assumed to work. The CI matrix should cover all three platforms from the first storage/vault-lifecycle phase.

**The password generator deferral is acceptable** given that users can paste from external generators. The field-edit UX must be designed now with the generator slot in mind so v2 does not require a structural UI refactor.

---

## Recommended Phase Sequence

Ordered strictly by dependency and risk. Each phase builds only on verified foundations.

| # | Phase | Rationale |
|---|-------|-----------|
| 1 | **Crypto package** (`internal/crypto`) | Foundation for every other layer. Pure functions, zero I/O, fast to test. Establish `[]byte`-only conventions, nonce generation, and mlock wrappers here — before any sensitive data flows exist. |
| 2 | **vault domain + Manager** (`internal/vault`) | Pure in-memory business logic. All entity types, lifecycle state machine, Manager API, business rules (soft-delete, Pasta Geral protection, cycle detection). Fully testable without storage or TUI. |
| 3 | **Storage package** (`internal/storage`) | Binary file format, Argon2id parameter storage in header, atomic save with `.tmp`/`.bak`/`.bak2` state machine and startup recovery, hash-based external change detection, version migration scaffold. |
| 4 | **TUI scaffold + root model** (`internal/tui`) | Session state enum, root model structure, flat child-model composition, message routing and overlay precedence, global tick, timer fields. No screens yet — just the skeleton that all screens live in. |
| 5 | **Welcome screen + vault create/open** | `welcomeModel`, file path input, `vault.Manager.Create` / `Open` integration, error domain mapping, master password input with strength indicator. First end-to-end flow. |
| 6 | **Vault tree + search** (`vaultModel`) | Folder/secret hierarchy rendering, keyboard navigation, search overlay (common fields + name + note only, sensitive excluded), favorites and soft-delete visual treatment. |
| 7 | **Secret detail + edit** (`secretDetailModel`) | View mode, create mode, edit-basic mode (values), edit-advanced mode (field structure), template picker, duplicate, move, reorder. Design with generator slot on sensitive field inputs. |
| 8 | **Vault lifecycle operations** | Save, Save As, Discard/Reload, Master password change, Export, Import (JSON merge), external change detection dialog. All through Manager — no new domain logic. |
| 9 | **Timers, clipboard, screen clear, lock/exit** | Auto-lock (inactivity reset on ALL input types), manual lock, clipboard auto-clear (all platforms), field reveal timer, screen+scrollback clear on lock/exit, memory wipe, clean exit with unsaved-change prompt. |
| 10 | **Keyboard help overlay + settings screen** | `?` help modal using `bubbles/v2/help`, configurable timer settings (VAULT-16), polish pass on key-hint footer in all screens. |
| 11 | **Cross-platform CI + integration tests** | Full build matrix (Windows/macOS/Linux), race detector on all tests, golden-file suite, atomic save failure simulation, platform clipboard tests, mlock quota tests, golden-file secret-pattern scan. |

---

## Gaps vs Requirements

| Gap | Source | Severity | Action |
|-----|--------|----------|--------|
| **Keyboard help overlay** (`?` modal + footer hints) | Not in requisitos.md; flagged by FEATURES research as essential for TUI onboarding | **Add to v1** — fold into Phase 10; small scope. |
| **`QUERY-02` negative test (sensitive fields excluded from search)** | Requirement exists, but easy to implement incorrectly | **Day-1 test** — write the negative test in Phase 2 (vault domain), not Phase 6. |
| **In-app help for folder/template management** | requisitos.md describes operations but no UX for discoverability | **Covered by Phase 10** keyboard help overlay. |
| **Startup recovery for orphaned `.bak2` files** | Not in requisitos.md; critical for backup protocol correctness | **Add to Phase 3** as part of storage package initialization. |
| **`COMPAT-01` migration test** | Requirement states backward compat; no explicit test strategy in PROJECT.md | **Add to Phase 3** — format versioning scaffold must include a migration test harness from day one. |
| **Windows rename atomicity** | ATOMIC-01 assumes `os.Rename` behavior that is not guaranteed on Windows | **Add to Phase 3** — Windows-specific `MoveFileEx` code path with retry on transient lock. |
| **Wayland clipboard clearing** | PORT-01 implies all platforms; Wayland clipboard is a separate surface from X11 | **Verify in Phase 9 CI** — add Wayland test to CI matrix or document explicitly as best-effort. |
| **Password generator placeholder in field-edit UX** | Feature deferred to v2; no design note in requirements | **Design decision in Phase 7** — leave a "Generate →" action slot on sensitive field inputs. Zero v2 refactor cost if done now. |
| **Import competitor formats (KeePass XML, Bitwarden CSV)** | Only Abditum JSON import in v1 | **v2 item confirmed** — note in POST-V1.md; document as known adoption friction. |
