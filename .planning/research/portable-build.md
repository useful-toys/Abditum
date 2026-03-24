# Portable Single Binary Build & Distribution — Abditum

**Focus:** Go cross-compilation, CGO, static binary, build tooling, binary size, code signing, version injection, and `golang.org/x/crypto` constraints for a TUI password vault.
**Researched:** 2026-03-24
**Overall Confidence:** HIGH — verified directly from source code, official repos, and real-world goreleaser configs.

---

## 1. Cross-Compilation Targets

Go's cross-compilation is first-class. Set `GOOS` and `GOARCH` at build time; no separate toolchain is needed when CGO is disabled (see section 2).

### Required targets for Abditum

| Platform | `GOOS` | `GOARCH` | Notes |
|---|---|---|---|
| Windows (Intel) | `windows` | `amd64` | Primary desktop target |
| Windows (ARM) | `windows` | `arm64` | Surface Pro X, ARM laptops — important for portability |
| macOS Intel | `darwin` | `amd64` | Rosetta covers this on M-series too |
| macOS Apple Silicon | `darwin` | `arm64` | Native M-series — required for best performance |
| Linux (Intel) | `linux` | `amd64` | Most servers/desktops |
| Linux (ARM) | `linux` | `arm64` | Raspberry Pi, ARM cloud VMs, Termux |

**Recommendation:** Target all six. goreleaser handles this trivially in one config. Adding them costs nothing at build time.

### Manual cross-compile example (without goreleaser)

```bash
# macOS arm64 — from any host with Go installed
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build \
  -ldflags="-s -w -X main.version=${VERSION}" \
  -trimpath \
  -o dist/abditum-darwin-arm64 \
  ./cmd/abditum

# Windows amd64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build \
  -ldflags="-s -w -X main.version=${VERSION}" \
  -trimpath \
  -o dist/abditum-windows-amd64.exe \
  ./cmd/abditum
```

---

## 2. CGO Implications — The Critical Question

### Bubble Tea v2 (charm.land/bubbletea/v2): **CGO-FREE** ✅

Verified directly from the `go.mod` and source files. All dependencies of Bubble Tea v2 are pure Go:

```
charm.land/bubbletea/v2
├── github.com/charmbracelet/colorprofile
├── github.com/charmbracelet/ultraviolet
├── github.com/charmbracelet/x/ansi
├── github.com/charmbracelet/x/term
├── github.com/charmbracelet/x/termios       ← syscall-based, not CGO
├── github.com/charmbracelet/x/windows       ← pure Go Win32 syscalls
├── github.com/muesli/cancelreader           ← pure Go
└── golang.org/x/sys                         ← pure Go syscall wrappers
```

Bubble Tea uses `golang.org/x/sys` (pure Go syscall bindings) for raw terminal mode and OS-level I/O — not `libncurses`, not CGO. It has explicit `_windows.go` and `_unix.go` build-tagged files. **No C compiler required.**

### Clipboard: **The Tricky Dependency**

#### Option A: `atotto/clipboard` (common choice) — **PARTIALLY PROBLEMATIC**

Verified from source (`clipboard_unix.go`):

- **Windows**: Pure Go syscall via `user32.dll` / `kernel32.dll`. No CGO. ✅
- **macOS**: Uses `pbcopy`/`pbpaste` subprocess calls. No CGO. ✅
- **Linux X11/Wayland**: Shells out to `xclip`, `xsel`, or `wl-clipboard`. **These tools must be present on the host system.** This breaks the "pure portable binary" contract.

For a headless Linux terminal without a clipboard manager installed, atotto/clipboard will fail silently (sets `Unsupported = true`). This is acceptable for a terminal app — clipboard requires a display server anyway.

#### Option B: Bubble Tea v2 built-in OSC52 clipboard (`SetClipboard`)

Verified from `clipboard.go`: Bubble Tea v2 has **native clipboard support via OSC52 escape sequences**. This works in most modern terminals (iTerm2, Windows Terminal, Alacritty, Kitty, etc.) without any external dependency — the terminal itself handles clipboard access.

**Recommendation for Abditum:** Use **Bubble Tea's `SetClipboard` (OSC52)** as primary method. It is:
- CGO-free
- No external binary dependency
- Works identically on all three platforms
- Functionally simpler (write-only is sufficient for a password vault)

Fallback note: OSC52 is not supported in all terminals (notably older tmux configurations, some SSH setups, and VS Code integrated terminal). Document this limitation in the README. For unsupported terminals, gracefully surface an error toast.

If broader compatibility is required, adding `atotto/clipboard` as a secondary fallback is reasonable — it is still CGO-free, just requires system tools on Linux.

### `golang.org/x/crypto` (Argon2id): **CGO-FREE, with Assembly Optimization** ✅

Verified directly from source (`argon2/blamka_amd64.go`, `blamka_ref.go`):

```go
// blamka_amd64.go — build constraint:
//go:build amd64 && gc && !purego

// blamka_ref.go — fallback for other architectures:
//go:build !amd64 || purego || !gc
```

The Argon2id implementation uses **Go assembly** (`.s` files) for amd64 SSE2/SSE4 optimization — **not CGO**. On arm64, it falls through to a pure Go generic implementation. This means:

- `CGO_ENABLED=0` works on all platforms. ✅
- amd64 gets native SIMD-accelerated Argon2id (fast). ✅
- arm64 (M-series Mac, ARM Linux, Windows ARM) gets the generic Go implementation (slower but fully correct). ✅
- No C compiler, no `libargon2` shared library needed. ✅

**Summary: The entire dependency tree for Abditum can be built with `CGO_ENABLED=0` on all target platforms.**

---

## 3. Static Binary Considerations

### What "static" means in Go

With `CGO_ENABLED=0`, Go produces a **fully self-contained binary** with no shared library dependencies. Verification:

```bash
# Linux — should show "not a dynamic executable"
file abditum-linux-amd64
ldd abditum-linux-amd64  # → "not a dynamic executable"

# macOS — dynamically linked to system libSystem only (unavoidable)
otool -L abditum-darwin-arm64
# Expected: /usr/lib/libSystem.B.dylib (always present on macOS)
```

### macOS caveat: `libSystem`

Even with `CGO_ENABLED=0`, Go binaries on macOS link to `/usr/lib/libSystem.B.dylib`. This is always present on any macOS installation and is not a portability concern. No action needed.

### Linux musl vs glibc

With `CGO_ENABLED=0`, the Go runtime uses pure Go syscalls — **no glibc dependency**. The binary runs on any Linux kernel version ≥ 2.6.23 (the Go runtime minimum). This makes the binary compatible with Alpine Linux (musl-based) without any special build flags.

### Windows: No DLL dependencies beyond OS-provided ones

`atotto/clipboard` (if used) calls `user32.dll` and `kernel32.dll` via Go's `syscall` package. These are guaranteed to be present on every Windows installation.

---

## 4. Build Tooling

### Recommendation: **goreleaser** (free tier, open source)

For a small OSS project that needs to ship cross-platform binaries to GitHub Releases, goreleaser is the standard choice. Used by Charm (Bubble Tea's creators), jira-cli, and thousands of Go CLI/TUI projects.

**Why goreleaser over raw `go build` scripts:**
- Parallel cross-compilation of all targets in one command
- Automatic GitHub Releases with checksums (`checksums.txt` with SHA256)
- Archive generation (`.tar.gz` for Unix, `.zip` for Windows)
- Version injection from git tags
- Homebrew formula generation (optional, useful later)
- GitHub Actions integration via `goreleaser/goreleaser-action`

**Why not raw Makefile:**
A Makefile is fine for local dev builds (and should exist for `make build`, `make test`) but goreleaser provides the release automation. Use both.

### Minimal goreleaser config for Abditum (`.goreleaser.yml`)

```yaml
version: 2

project_name: abditum

before:
  hooks:
    - go mod tidy

builds:
  - <<: &build_defaults
      binary: abditum
      main: ./cmd/abditum
      ldflags:
        - -s -w
        - -X main.version={{.Version}}
        - -X main.commit={{.FullCommit}}
        - -X main.buildDate={{.Date}}
      env:
        - CGO_ENABLED=0
      flags:
        - -trimpath

    id: darwin
    goos: [darwin]
    goarch: [amd64, arm64]

  - <<: *build_defaults
    id: linux
    goos: [linux]
    goarch: [amd64, arm64]

  - <<: *build_defaults
    id: windows
    goos: [windows]
    goarch: [amd64, arm64]

archives:
  - id: unix
    ids: [darwin, linux]
    name_template: >-
      {{ .ProjectName }}_{{ .Version }}_
      {{- if eq .Os "darwin" }}macOS{{- else }}{{ .Os }}{{- end }}_
      {{- if eq .Arch "amd64" }}x86_64{{- else }}{{ .Arch }}{{- end }}
    formats: [tar.gz]
    files: [LICENSE, README.md]

  - id: windows
    ids: [windows]
    name_template: >-
      {{ .ProjectName }}_{{ .Version }}_windows_
      {{- if eq .Arch "amd64" }}x86_64{{- else }}{{ .Arch }}{{- end }}
    formats: [zip]
    files: [LICENSE, README.md]

checksum:
  name_template: 'checksums.txt'
  algorithm: sha256

changelog:
  use: github
  sort: desc
```

### GitHub Actions workflow (CI/Release)

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0   # goreleaser needs full git history for changelog

      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - uses: goreleaser/goreleaser-action@v6
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## 5. Version String Injection

The standard Go idiom for embedding version at build time uses `ldflags -X`:

```go
// cmd/abditum/main.go
package main

var (
    version   = "dev"      // overridden by ldflags at release build
    commit    = "unknown"  // overridden by ldflags at release build
    buildDate = "unknown"  // overridden by ldflags at release build
)
```

Build command:
```bash
VERSION=1.0.0
COMMIT=$(git rev-parse --short HEAD)
DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

go build \
  -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildDate=${DATE}" \
  -trimpath \
  -o abditum \
  ./cmd/abditum
```

goreleaser automates this using template variables (`{{.Version}}`, `{{.FullCommit}}`, `{{.Date}}`).

**Important for portability:** Use `-trimpath` to remove local build paths from the binary. This prevents leaking developer machine paths in stack traces.

---

## 6. Binary Size

### Expected size for Abditum

A Go TUI app with Bubble Tea, x/crypto, and JSON encoding:

| Configuration | Estimated Size |
|---|---|
| Default build (no flags) | 12–18 MB |
| With `-ldflags="-s -w"` | 8–12 MB |
| With `-ldflags="-s -w" -trimpath` | 8–12 MB |

`-s` strips the symbol table. `-w` strips DWARF debug info. Together they reduce binary size by 20–35% with zero runtime impact.

Actual sizes from real Charm CLI projects (glow, mods) released on GitHub are in the 8–15 MB range for the full stack.

### UPX: **NOT RECOMMENDED for Abditum**

UPX is a binary packer/compressor. For Abditum specifically, there are two strong reasons to avoid it:

**1. AV False Positives (Critical for a security tool)**
UPX-compressed executables are flagged by antivirus software as suspicious because malware commonly uses UPX to evade detection. For a password vault — a security-sensitive application — having Windows Defender or other AV software quarantine the binary at download or first launch would be catastrophic for user trust. This is not a theoretical risk; it is well-documented community experience.

**2. macOS Gatekeeper incompatibility**
UPX-compressed macOS binaries cannot be notarized by Apple. Notarization requires the binary to pass static analysis, which UPX packing defeats.

**3. Negligible benefit**
8–12 MB is already small. Saving 3–5 MB via UPX is not worth the tradeoffs.

**Verdict:** Skip UPX. The raw `-s -w` binary at 8–12 MB is the right artifact.

---

## 7. Code Signing

### Windows: SmartScreen Warning

Unsigned Windows executables downloaded from the internet trigger SmartScreen's "Windows protected your PC" popup on first run. This is expected behavior for unknown publishers.

**Options:**

| Option | Cost | Effort | Result |
|---|---|---|---|
| No signing | Free | None | SmartScreen warning on first run — users must click "More info → Run anyway" |
| EV Code Signing Certificate | $200–$500/year | Significant (DigiCert, Sectigo) | No SmartScreen warning immediately |
| Standard OV Certificate | $70–$200/year | Moderate | SmartScreen warning eventually clears via reputation |
| Microsoft Trusted Signing (Azure) | ~$10/month | Moderate | Newer Microsoft-managed option, clears SmartScreen faster |

**Recommendation for Abditum v1:** Ship unsigned. Document the SmartScreen behavior in README with the "More info → Run anyway" instruction. The target user of a portable security tool is technically sophisticated enough to handle this. Pursue a code signing certificate in a future milestone if user feedback indicates it's a blocker.

Note: goreleaser supports Windows code signing via the `sbom` and `sign` customizations if a certificate is added later.

### macOS: Gatekeeper Quarantine

Unsigned binaries downloaded from the internet are quarantined by macOS Gatekeeper. Users will see "Apple cannot verify this app is free from malware."

**The user workaround (document in README):**
```bash
# Remove the quarantine attribute after download
xattr -dr com.apple.quarantine abditum-darwin-arm64

# OR: right-click in Finder → Open → Open (first time only)
```

**Full signing + notarization (future milestone):**
- Requires Apple Developer Program ($99/year)
- Must codesign the binary: `codesign --sign "Developer ID Application: ..." abditum`
- Must notarize with Apple's servers: `xcrun notarytool submit ...`
- Goreleaser Pro supports notarization via the `notarize` customization

**Recommendation for Abditum v1:** Ship unsigned with clear README documentation of the `xattr` workaround. This is standard practice for small OSS CLI tools (mise, fd, ripgrep all do this initially).

**Linux:** No equivalent signing requirement. Distribute directly. SHA256 checksums in the release are sufficient for integrity verification.

---

## 8. `golang.org/x/crypto` — Platform Constraints Summary

Verified from source. All relevant packages are **fully cross-compilable with `CGO_ENABLED=0`**:

| Package | Constraint | Notes |
|---|---|---|
| `argon2` | None (CGO-free) | Uses Go assembly for amd64 optimization; pure Go fallback for arm64, others |
| `aes` (stdlib) | None | Go stdlib uses platform assembly (Go asm, not CGO) |
| `blake2b` | None | Used internally by argon2; pure Go |
| `chacha20poly1305` | None | AEAD alternative if ever needed |
| `pbkdf2` | None | Pure Go |

The amd64 assembly in `argon2/blamka_amd64.s` is **Go assembly** (compiled by the Go toolchain), not C assembled by a C compiler. It is included automatically when cross-compiling for amd64 regardless of CGO state.

**Concrete build constraint:** `//go:build amd64 && gc && !purego`
- `gc` means "use the standard Go compiler" (always true with `go build`)
- `!purego` is a build tag you can set if you want to force the generic implementation (useful for testing)
- The fallback `blamka_ref.go` covers `!amd64 || purego || !gc` — this is what arm64 builds use

**No platform-specific build constraints block any target.** All six targets (Windows/macOS/Linux × amd64/arm64) build cleanly.

---

## 9. Recommended Makefile (Developer Convenience)

```makefile
# Abditum — local development build targets
# For release builds, use goreleaser

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -s -w \
  -X main.version=$(VERSION) \
  -X main.commit=$(COMMIT) \
  -X main.buildDate=$(DATE)

BUILD_FLAGS := -ldflags="$(LDFLAGS)" -trimpath

.PHONY: build test lint clean dist

## build: Build for the current platform
build:
	CGO_ENABLED=0 go build $(BUILD_FLAGS) -o abditum ./cmd/abditum

## test: Run all tests
test:
	go test ./...

## test-race: Run tests with race detector
test-race:
	go test -race ./...

## lint: Run golangci-lint
lint:
	golangci-lint run

## dist: Cross-compile all targets (for local testing; use goreleaser for releases)
dist:
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build $(BUILD_FLAGS) -o dist/abditum-darwin-amd64 ./cmd/abditum
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build $(BUILD_FLAGS) -o dist/abditum-darwin-arm64 ./cmd/abditum
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build $(BUILD_FLAGS) -o dist/abditum-linux-amd64 ./cmd/abditum
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build $(BUILD_FLAGS) -o dist/abditum-linux-arm64 ./cmd/abditum
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o dist/abditum-windows-amd64.exe ./cmd/abditum
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build $(BUILD_FLAGS) -o dist/abditum-windows-arm64.exe ./cmd/abditum

## clean: Remove build artifacts
clean:
	rm -rf abditum dist/
```

---

## 10. Complete Dependency CGO Audit

All key dependencies verified as CGO-free:

| Dependency | CGO? | Mechanism | Notes |
|---|---|---|---|
| `charm.land/bubbletea/v2` | No | Pure Go + Go syscalls | Has `_windows.go` / `_unix.go` build tags |
| `github.com/charmbracelet/lipgloss` | No | Pure Go | Styling only |
| `github.com/charmbracelet/bubbles` | No | Pure Go | TUI components |
| `golang.org/x/crypto/argon2` | No | Go assembly (not CGO) | amd64 SSE optimized, arm64 pure Go |
| `golang.org/x/sys` | No | Pure Go syscall wrappers | Required by Bubble Tea |
| `atotto/clipboard` (if used) | No | Win32 syscalls / exec on Linux | Linux needs xclip/xsel installed |
| `encoding/json` (stdlib) | No | Pure Go | — |

**Conclusion:** `CGO_ENABLED=0` is safe and correct for all Abditum builds.

---

## 11. Pitfalls & Anti-Patterns

### Pitfall 1: Not setting `CGO_ENABLED=0` explicitly
**Problem:** Go may enable CGO by default if a C compiler is found on the build machine. This can produce binaries that dynamically link against glibc on Linux CI (e.g., `GLIBC_2.34` required), which won't run on older distros.
**Prevention:** Always set `CGO_ENABLED=0` explicitly in `ldflags`, CI environment, and goreleaser config.

### Pitfall 2: Forgetting `-trimpath`
**Problem:** Go by default embeds the absolute path of source files in the binary (visible in stack traces and `strings` output). This leaks developer machine paths and makes builds non-reproducible.
**Prevention:** Always use `-trimpath`.

### Pitfall 3: Using UPX in CI
**Problem:** Compressed binaries trigger AV false positives and break macOS notarization.
**Prevention:** Do not use UPX. Use `-s -w` instead.

### Pitfall 4: Relying on `atotto/clipboard` on Linux without fallback
**Problem:** On a Linux machine without `xclip`, `xsel`, or `wl-clipboard` installed, clipboard operations silently fail.
**Prevention:** Prefer Bubble Tea's OSC52 `SetClipboard`. If using `atotto/clipboard`, catch the error and surface a helpful toast: "Clipboard not available. Install xclip or wl-clipboard."

### Pitfall 5: Windows arm64 build ignored
**Problem:** Many projects skip `windows/arm64`. With the growth of Snapdragon X Elite and ARM laptops on Windows, this target is increasingly relevant for a "truly portable" tool.
**Prevention:** Include `windows/arm64` in goreleaser config. It costs nothing.

### Pitfall 6: Not including `go.sum` in version control
**Problem:** `go.sum` must be committed for reproducible builds. Its absence breaks `go mod verify` and CI.
**Prevention:** Commit both `go.mod` and `go.sum`. Run `go mod tidy` in CI as a check.

### Pitfall 7: macOS quarantine surprises users
**Problem:** Users double-click the binary in Finder, macOS says it can't be opened, they give up.
**Prevention:** Document the `xattr -dr com.apple.quarantine abditum` command prominently in README and GitHub Release notes.

---

## 12. Sources

| Source | Confidence | URL |
|---|---|---|
| Bubble Tea v2 go.mod (direct) | HIGH | https://github.com/charmbracelet/bubbletea/blob/main/go.mod |
| Bubble Tea clipboard.go (OSC52) | HIGH | https://github.com/charmbracelet/bubbletea/blob/main/clipboard.go |
| atotto/clipboard — Unix impl | HIGH | https://github.com/atotto/clipboard/blob/master/clipboard_unix.go |
| atotto/clipboard — Windows impl | HIGH | https://github.com/atotto/clipboard/blob/master/clipboard_windows.go |
| golang/crypto argon2 directory | HIGH | https://github.com/golang/crypto/tree/master/argon2 |
| argon2 blamka_amd64.go (Go assembly, not CGO) | HIGH | https://github.com/golang/crypto/blob/master/argon2/blamka_amd64.go |
| argon2 blamka_ref.go (pure Go fallback) | HIGH | https://github.com/golang/crypto/blob/master/argon2/blamka_ref.go |
| jira-cli goreleaser config (real-world example) | HIGH | https://github.com/ankitpokhrel/jira-cli/blob/main/.goreleaser.yml |
| Bubble Tea .goreleaser.yml | HIGH | https://github.com/charmbracelet/bubbletea/blob/main/.goreleaser.yml |
| Glow .goreleaser.yml | HIGH | https://github.com/charmbracelet/glow/blob/master/.goreleaser.yml |
