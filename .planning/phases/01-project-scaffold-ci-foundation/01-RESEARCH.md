# Phase 1 Research: Project Scaffold + CI Foundation

**Phase:** 01-project-scaffold-ci-foundation
**Researched:** 2026-03-29
**Status:** Complete

## Research Goal

Answer: "What do I need to know to PLAN Phase 1 well?"

This phase establishes the Go project foundation: module initialization, directory structure, CI pipeline, build tooling, and linter configuration. No business logic — purely scaffolding.

## Key Findings

### 1. Go Module & Import Paths

**Critical Decision from CONTEXT:**
- Module path: `github.com/useful-toys/abditum` (NOT `github.com/user/abditum` from roadmap placeholder)
- Go version: `1.26.1+` in go.mod
- Charm library import paths changed in late 2025 to canonical `charm.land/*` domain

**Dependencies for Phase 1:**
```
charm.land/bubbletea/v2
charm.land/bubbles/v2
charm.land/lipgloss/v2
golang.org/x/crypto
golang.org/x/sys
golang.org/x/text
github.com/atotto/clipboard
```

**NOT NEEDED in Phase 1:** `github.com/matoous/go-nanoid/v2` — per CONTEXT.md, will implement NanoID internally using `crypto/rand` in later phases.

### 2. Static Linking Requirements

From REQUIREMENTS.md (COMPAT-01) and arquitetura.md:

**Build constraints:**
- `CGO_ENABLED=0` — MUST be set as environment variable globally in CI
- Sets flag at job level, not just in build command
- Ensures `go test` also uses static linking
- Eliminates libc dependency for true portability

**Verification:**
- On Linux: `file ./abditum` should report "statically linked"
- No dynamic library dependencies

### 3. Directory Structure

Per arquitetura.md §2:

```
cmd/
  abditum/           -- entry point (main.go exits 0, no logic yet)
internal/
  vault/             -- domain logic (doc.go stub only)
  crypto/            -- Argon2id + AES-256-GCM (doc.go stub only)
  storage/           -- file I/O and atomic saves (doc.go stub only)
  tui/               -- Bubble Tea interface (doc.go stub only)
```

**Each doc.go:**
- Package declaration
- One-sentence package doc comment
- No code yet — just establishes import paths

### 4. CI Configuration

**Target: `.github/workflows/ci.yml`**

**Jobs required:**
1. **build** — `GOOS=linux CGO_ENABLED=0 go build ./cmd/abditum`
2. **lint** — `golangci-lint run`
3. **test** — `CGO_ENABLED=0 go test ./... -race -count=1`

**Triggers:**
- `push` to `main` branch (NOT `master` — renamed per STATE.md)
- `pull_request` to `main`

**Go version in CI:**
- Use `go-version: '1.26.x'` in setup-go action
- Matches go.mod requirement of `1.26.1+`

**Critical:** Set `CGO_ENABLED=0` as environment variable at job level:
```yaml
env:
  CGO_ENABLED: 0
```

### 5. Makefile Targets

Standard targets required:
- `build` — uses `CGO_ENABLED=0`, respects `GOOS` from environment
- `test` — includes `-race` flag
- `lint` — runs golangci-lint
- `vet` — runs go vet
- `clean` — removes build artifacts

**Output binary name:** `abditum` (per project name)

### 6. Linter Configuration (golangci-lint)

**Target: `.golangci.yml`**

**Linters to enable:**
- `errcheck` — unchecked errors
- `govet` — standard Go vet
- `staticcheck` — static analysis
- `revive` — style and best practices
- `gosec` — security issues
- `gocritic` — performance and style
- `forbidigo` — forbidden imports/functions

**Forbidden patterns (via forbidigo):**
- `math/rand` — MUST use `crypto/rand`
- `net` — network calls prohibited
- `net/http` — HTTP calls prohibited

**Configuration needs:**
- Severity thresholds appropriate for Go 1.26
- Nolint exceptions for empty stub packages (will have no code yet)
- Allow doc.go files with only package comment

### 7. Testing Considerations

**For this phase:**
- No actual tests yet (no logic to test)
- `go test ./...` should pass (reports "no tests")
- `go vet ./...` should pass with zero findings on stubs

**Golden file testing (future):**
- Use `theatest/v2` with `WithInitialTermSize(80, 24)`
- Ensures stable terminal size for reproducible snapshots

### 8. Security & Build Architecture

From arquitetura.md §5:

**Build principles:**
- Static binary: `CGO_ENABLED=0`
- Network isolation: no `net` or `net/http` imports
- Crypto only from stdlib and `golang.org/x/crypto`
- Minimal dependencies

**Memory safety patterns (future phases):**
- Sensitive data as `[]byte` (not `string`)
- Custom JSON marshaling for password fields
- Explicit zeroing on lock/exit

## Validation Architecture

**Not applicable for Phase 1** — no business logic or user-facing behavior to validate. Verification is purely build/compile checks:
- Binary compiles successfully
- Static linking confirmed
- CI jobs all green
- Import paths all canonical `charm.land/*`

## Common Pitfalls

### Import Path Migration
**Issue:** Charm libraries migrated from `github.com/charmbracelet/*` to `charm.land/*` in late 2025.

**Impact:** Using old import paths will fail or pull wrong versions.

**Mitigation:**
- All imports MUST use `charm.land/bubbletea/v2`, etc.
- Check `go.mod` requires to ensure no `github.com/charmbracelet` entries
- CI UAT includes explicit check: zero `github.com/charmbracelet/*` imports in `.go` files

### CGO_ENABLED Scope
**Issue:** Setting `CGO_ENABLED=0` only on build command doesn't affect `go test`.

**Impact:** Tests may pass with dynamic linking, hiding portability issues.

**Mitigation:**
- Set as environment variable at CI job level:
  ```yaml
  env:
    CGO_ENABLED: 0
  ```
- NOT just in command: `CGO_ENABLED=0 go build ...`

### Branch Name
**Issue:** Roadmap examples may reference `master` branch.

**Impact:** CI triggers won't match actual branch name.

**Mitigation:**
- Per CONTEXT.md and STATE.md: use `main` (already renamed)
- CI workflow: `on: push: branches: [main]`

### Empty Package Linting
**Issue:** Linters may complain about packages with only doc.go and no functions.

**Impact:** CI lint job fails on valid stubs.

**Mitigation:**
- Configure `.golangci.yml` with nolint exceptions for stub packages
- OR accept that some linters will have zero issues to report (not a failure)

## Dependencies Summary

**Production dependencies (go.mod):**
```
charm.land/bubbletea/v2
charm.land/bubbles/v2
charm.land/lipgloss/v2
golang.org/x/crypto
golang.org/x/sys
golang.org/x/text
github.com/atotto/clipboard
```

**Dev/test dependencies (future):**
```
github.com/charmbracelet/x/exp/teatest/v2  (Phase 5+)
```

**NOT included in Phase 1:**
- `github.com/matoous/go-nanoid/v2` — will implement internally

## Planning Guidance

### Task Breakdown Strategy
**Recommend 5 atomic tasks** (matching roadmap structure):

1. **Initialize Go module + dependencies**
   - Create `go.mod` with correct module path and Go version
   - Add all required dependencies
   - Verify with `go mod tidy` and `go mod verify`

2. **Create directory structure with package stubs**
   - Create `cmd/abditum/main.go` (exits 0)
   - Create `internal/{crypto,vault,storage,tui}/doc.go` stubs
   - Each stub: package declaration + one-sentence doc

3. **Configure CI workflow**
   - Create `.github/workflows/ci.yml`
   - Three jobs: build, lint, test
   - `CGO_ENABLED=0` at job level
   - Triggers on push/PR to `main`

4. **Add Makefile**
   - Targets: build, test, lint, vet, clean
   - `CGO_ENABLED=0` in build target
   - `-race` in test target

5. **Configure golangci-lint**
   - Create `.golangci.yml`
   - Enable required linters
   - Configure forbidigo for `math/rand`, `net`, `net/http`
   - Severity thresholds and stub exceptions

### Verification Strategy
Each task should verify:
- Files created in expected locations
- Content matches requirements
- Commands succeed: `go build ./cmd/abditum`, `go vet ./...`, `golangci-lint run`

### Must-Have Artifacts
From goal-backward analysis:

**Observable truths:**
- Binary compiles successfully with `CGO_ENABLED=0`
- CI workflow exists and is triggered by push to `main`
- All imports use `charm.land/*` canonical paths
- Linter runs without configuration errors

**Required artifacts:**
- `go.mod` with module path `github.com/useful-toys/abditum`
- Directory tree: `cmd/abditum/`, `internal/{crypto,vault,storage,tui}/`
- `.github/workflows/ci.yml` with 3 jobs
- `Makefile` with 5 targets
- `.golangci.yml` with forbidigo rules

**Key connections:**
- `go.mod` → `cmd/abditum/main.go` (compiles)
- CI workflow → `main` branch (triggers)
- forbidigo config → imports in `.go` files (enforces restrictions)

## Questions for Planner

None — requirements are clear and specific. All decisions locked in CONTEXT.md.

---

*Research complete. Ready for planning.*
