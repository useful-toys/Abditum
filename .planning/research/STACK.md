# Stack Research

**Domain:** Go TUI Password Manager
**Researched:** 2026-03-27
**Confidence:** HIGH (all versions verified against official release pages and documentation)

---

## Recommended Stack

### Core Technologies

| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| Go | 1.26.1 | Language / runtime | Latest stable (released 2026-02-10/03-05); `CGO_ENABLED=0` gives a fully static binary on all three platforms; strong stdlib coverage means fewer third-party dependencies |
| Bubble Tea v2 | v2.0.2 | TUI event loop and model management | Industry standard for Go TUIs; Elm-architecture keeps state predictable; v2 is the stable GA release (`charm.land/bubbletea/v2`); already decided in arquitetura.md |
| `golang.org/x/crypto` | v0.49.0 (latest) | Argon2id key derivation | *Only* Go-maintained library allowed by the crypto policy; `argon2.IDKey` implements Argon2id exactly; stdlib-grade trust level |

### Supporting Libraries — Bubble Tea Ecosystem

| Library | Version | Import Path | Purpose | When to Use |
|---------|---------|-------------|---------|-------------|
| Bubbles v2 | v2.1.0 | `charm.land/bubbles/v2` | Pre-built TUI components | Use `textinput` for password entry, `list` for vault tree, `viewport` for secret detail, `help` for key-binding guide |
| Lip Gloss v2 | v2.0.2 | `charm.land/lipgloss/v2` | Terminal styling, layout, borders | Renders all styled text; use `lipgloss.NewStyle()` for colours/borders; use `LightDark()` helper for light/dark terminal awareness |
| teatest (v2) | pseudo-version (2026-03-23) | `github.com/charmbracelet/x/exp/teatest/v2` | TUI integration + golden-file testing | `NewTestModel` + `RequireEqualOutput` for all screen snapshots; already chosen in arquitetura.md |

> **Note on charm.land import paths:** Charm migrated their vanity domain in late 2025. Use `charm.land/*` for all three packages above — the old `github.com/charmbracelet/*` paths redirect correctly but the canonical form is `charm.land/bubbletea/v2`, `charm.land/bubbles/v2`, `charm.land/lipgloss/v2`.

### Supporting Libraries — Systems

| Library | Version | Import Path | Purpose | Notes |
|---------|---------|-------------|---------|-------|
| `golang.org/x/sys` | latest | `golang.org/x/sys` | `mlock` (Linux/macOS) and `VirtualLock` (Windows) for locking master-key memory pages | Also needed for low-level syscall access; CGO-free on all platforms |

### Standard Library Usage (no external dependency)

| Package | Usage |
|---------|-------|
| `crypto/aes` + `crypto/cipher` | AES-256-GCM encryption/decryption |
| `crypto/rand` | Salt (16 bytes), nonce (12 bytes), and NanoID-equivalent ID generation |
| `encoding/binary` / `encoding/json` | File format serialisation |
| `os` | Atomic save via temp file + rename |
| `time` | Inactivity timer, clipboard clear timer |

---

## Cryptography Details

### Argon2id Parameter Recommendations

Source: [RFC 9106 §4 and §7.4](https://www.rfc-editor.org/rfc/rfc9106#section-7.4) + `pkg.go.dev/golang.org/x/crypto/argon2`.

RFC 9106 gives two canonical parameter sets:

| Option | `time` (t) | `memory` (m) | `threads` (p) | `keyLen` | Context |
|--------|-----------|-------------|--------------|---------|---------|
| **First recommended** | 1 | 2 GiB (2 097 152 KiB) | 4 | 32 bytes | Uniformly safe; preferred on high-memory machines |
| **Second recommended** | 3 | 64 MiB (65 536 KiB) | 4 | 32 bytes | Memory-constrained environments |

**Recommendation for Abditum:**

Use the **second recommended set** as the *default* (t=3, m=64 MiB, p=4, keyLen=32) with the **first set** as an upgrade path for power users or future versions.

**Rationale:** A desktop password manager may run on machines with as little as 256 MB of free RAM (CI runner, low-end PC). 2 GiB is too aggressive as a non-negotiable default and would freeze the UI for several seconds on HDD-backed swap. 64 MiB at t=3 takes ~200–500 ms on a modern CPU, which is acceptable for an interactive unlock screen. If the user's machine has ≥ 4 GB RAM, t=1 / 2 GiB can be offered in settings.

**Concrete call:**
```go
// salt: 16 bytes from crypto/rand (128-bit)
// key derivation for AES-256 (32-byte output)
key := argon2.IDKey(
    password,   // []byte  — master password
    salt,       // []byte  — 16 random bytes stored in vault header
    3,          // time    — 3 passes over memory
    64*1024,    // memory  — 64 MiB in KiB
    4,          // threads — 4 parallel lanes
    32,         // keyLen  — 256-bit AES key
)
```

**Encoding in vault file header:** Store `{argon2_time, argon2_memory, argon2_threads, salt}` alongside the ciphertext so the parameters can be upgraded in future versions without breaking backward compatibility (COMPAT-01 in requirements).

### AES-256-GCM Usage Notes

- **Key**: 32-byte output of `argon2.IDKey` above.
- **Nonce**: 12 bytes (standard GCM nonce = 96 bits), generated fresh from `crypto/rand` on *every* save. Never reuse a nonce with the same key.
- **Tag**: GCM provides a built-in 16-byte authentication tag (appended by `cipher.AEAD.Seal`); verify via `Open` before accepting the plaintext.
- **Nonce collision risk**: With a 96-bit random nonce and a single fixed key, the birthday bound is ~2⁴⁸ encryptions before a collision becomes probable. A personal vault saved thousands of times over its lifetime is safe. The salt changes on every master-password rotation, which refreshes the key and resets the birthday counter.
- **CGO**: `crypto/aes` and `crypto/cipher` are pure Go stdlib — no CGO, no third-party dependency.

---

## Development Tools

| Tool | Purpose | Notes |
|------|---------|-------|
| `go build -trimpath -ldflags="-s -w"` | Reproducible, stripped binary | Remove debug info + symbol table; combine with `CGO_ENABLED=0` for fully static output |
| `golangci-lint` (latest) | Linting | Charm ships their own `.golangci.yml`; mirrors the Charm project conventions; catches common TUI mistakes |
| `go test -v -race ./...` | Test runner with race detector | The Bubble Tea event loop is concurrent; race detector catches subtle goroutine bugs |
| `teatest.RequireEqualOutput` + `-update` flag | Golden-file update workflow | Re-generate `.testdata/*.golden` files when intentional UI changes are made |

---

## Alternatives Considered

| Recommended | Alternative | When to Use Alternative |
|-------------|-------------|------------------------|
| Bubble Tea v2 (`charm.land/bubbletea/v2`) | Bubble Tea v1 (`github.com/charmbracelet/bubbletea`) | **Never for new projects** — v1 is in maintenance mode; v2 has a declarative `View()` struct, better renderer, and is where all Charm ecosystem development happens |
| Bubbles v2 (`charm.land/bubbles/v2`) | Roll your own TUI components | If no Bubbles component fits, write a custom `tea.Model`; but prefer Bubbles because it handles edge cases (Unicode widths, real cursor support, focus/blur) that you'd have to solve manually |
| Lip Gloss v2 (`charm.land/lipgloss/v2`) | `termenv` / raw ANSI | Lip Gloss v2 is pure (no I/O contention with Bubble Tea); raw ANSI requires manual colour downsampling for 256-colour and dumb terminals |
| Argon2id (RFC 9106) | bcrypt / scrypt | Argon2id is the PHC winner; bcrypt is limited to 72-byte passwords and has no memory-hardness tuning; scrypt is supported in x/crypto but Argon2id is strictly recommended by the RFC for new applications |
| `golang.org/x/sys` for mlock | `memguard` (github.com/awnumar/memguard) | `memguard` is a third-party library with its own enclave/buffer model; only use if you want automatic zeroing on GC movement, but it adds a dependency and complexity inconsistent with Abditum's zero-dependencies-unless-justified policy |
| `crypto/rand` for IDs | `github.com/matoous/go-nanoid/v2` | NanoID is fine if you need human-readable IDs; for internal opaque UUIDs, `crypto/rand` + `hex.EncodeToString` is zero-dependency and just as secure |

---

## What NOT to Use

- **Bubble Tea v1** (`github.com/charmbracelet/bubbletea` without `/v2`): In maintenance mode; `View()` returns a `string` rather than the new `tea.View` struct; you'd immediately hit API incompatibilities with Bubbles v2 and Lip Gloss v2 which already target the new interface.

- **`github.com/charmbracelet/bubbles` (v1 import path)**: Incompatible with Bubble Tea v2. The new import path is `charm.land/bubbles/v2`.

- **`tview`** (rivo/tview): Uses a different architecture (immediate mode); no Elm model; harder to test; not compatible with the teatest approach already chosen.

- **`tcell`** directly: Low-level; you'd rebuild what Bubble Tea/Lip Gloss already provide; only needed for custom renderers.

- **Any CGO-linked library** (including `golang.design/x/clipboard` on some platforms): `CGO_ENABLED=0` is a hard constraint (PORT-01, safety surface reduction). For clipboard, use `github.com/atotto/clipboard` which shells out to `xclip`/`xsel`/`pbcopy`/`clip.exe` — subprocess-based, CGO-free.

- **`math/rand`** for any randomness: Prohibited; use `crypto/rand` exclusively (arquitetura.md constraint).

- **Third-party encryption libraries** (e.g., `github.com/miscreant/miscreant`): Architecture decision restricts crypto to stdlib + `golang.org/x/crypto`; any third-party crypto lib is a potential supply-chain attack surface.

- **`bcrypt` or `scrypt`** for KDF: Argon2id is already chosen; bcrypt caps at 72 bytes password length; scrypt has no PHC endorsement for interactive use.

- **`encoding/gob`** for vault file format: Gob is Go-specific and not human-debuggable; a custom binary format (with explicit version field for COMPAT-01) is more portable for future tooling.

---

## Notes for Abditum

### Bubble Tea v2 API Differences vs. v1 (high impact)

The `View()` method in v2 returns `tea.View`, not `string`:

```go
// v2 pattern — always wrap content in tea.NewView()
func (m model) View() tea.View {
    v := tea.NewView(m.renderContent())
    v.AltScreen = true   // full-screen mode declared here, not in Init()
    return v
}
```

Features that were `tea.Cmd` in v1 are now fields on `tea.View`: alt-screen, mouse mode, focus reporting, window title. This eliminates the race condition between startup options and the first render.

### Key Bindings in Bubble Tea v2

Key messages are now typed: `tea.KeyPressMsg` (for key-down) and `tea.KeyReleaseMsg`. The space bar is now the string `"space"` instead of `" "`. Match keys using `msg.String()`:

```go
case tea.KeyPressMsg:
    switch msg.String() {
    case "ctrl+c", "q":
        return m, tea.Quit
    case "space":
        // toggle something
    case "ctrl+s":
        return m, saveVaultCmd()
    }
```

### Model Composition Pattern

For a multi-screen TUI like Abditum, the recommended pattern is a **root model** that owns a `currentScreen` enum and delegates `Update`/`View` to sub-models. Each screen (unlock, vault tree, secret detail, edit, settings) is its own `tea.Model`-implementing struct. Use `tea.Sequence` to chain commands (e.g., derive key → decrypt vault → transition to vault tree screen).

### Background Color Detection (Lip Gloss v2)

Request background color in `Init()` to style the UI correctly for light/dark terminals:

```go
func (m Model) Init() tea.Cmd {
    return tea.RequestBackgroundColor // sends tea.BackgroundColorMsg
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.BackgroundColorMsg:
        m.styles = newStyles(msg.IsDark())
    }
    // ...
}
```

### inactivity Timer Pattern

Bubble Tea v2's `tea.Tick` is the standard approach. Send a tick every second; compare `time.Since(lastActivity)` in `Update`. Reset `lastActivity` on any `tea.KeyPressMsg` or `tea.MouseClickMsg`. On timeout, emit a `lockVaultMsg`. This is cleaner than a goroutine-based timer because it runs inside the Bubble Tea event loop (no race conditions on the model).

### teatest v2 Golden File Workflow

The v2 subdirectory at `github.com/charmbracelet/x/exp/teatest/v2` was updated (2026-03-14) to use the `charm.land` import paths, matching Bubble Tea v2. Use `WithInitialTermSize(80, 24)` consistently across all tests so golden files are stable. Update golden files with:

```
go test ./internal/tui/... -update
```

### Argon2id Timing on CI

On a typical CI runner (2 CPU cores, 4–8 GB RAM), `t=3, m=64 MiB, p=4` completes in 300–800 ms. This is acceptable for unlock tests if the test uses a pre-derived key; do not re-derive in every test case. Keep one `TestMain` that derives a test key once and reuses it across the test suite.

### Static Build Command

```
CGO_ENABLED=0 GOOS=<target> GOARCH=<arch> go build \
  -trimpath \
  -ldflags="-s -w" \
  -o abditum \
  ./cmd/abditum
```

No linker flags for version injection are strictly needed, but `-X main.version=$(git describe --tags)` is conventional.

---

## Sources

| Source | Content | Confidence |
|--------|---------|------------|
| [github.com/charmbracelet/bubbletea/releases](https://github.com/charmbracelet/bubbletea/releases) | Bubble Tea v2.0.2 (latest), API changes, import paths | HIGH |
| [github.com/charmbracelet/bubbles/releases](https://github.com/charmbracelet/bubbles/releases) | Bubbles v2.1.0 (latest), functional options, v2 migration | HIGH |
| [github.com/charmbracelet/lipgloss/releases](https://github.com/charmbracelet/lipgloss/releases) | Lip Gloss v2.0.2 (latest), color API, LightDark helper | HIGH |
| [pkg.go.dev/github.com/charmbracelet/x/exp/teatest](https://pkg.go.dev/github.com/charmbracelet/x/exp/teatest) | teatest API, v2 subdir (published 2026-03-23) | HIGH |
| [pkg.go.dev/golang.org/x/crypto/argon2](https://pkg.go.dev/golang.org/x/crypto/argon2) | IDKey API, parameter recommendations | HIGH |
| [RFC 9106 §4 and §7.4](https://www.rfc-editor.org/rfc/rfc9106#section-7.4) | Argon2id first/second recommended parameter sets | HIGH |
| [go.dev/doc/devel/release](https://go.dev/doc/devel/release) | Go 1.26.1 current stable, release dates | HIGH |
