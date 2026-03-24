# Testing Strategy Research: Abditum

**Researched:** 2026-03-24  
**Focus:** teatest/v2 golden files, unit test patterns for crypto/storage/domain layers, integration tests, CI pipeline design  
**Confidence:** HIGH (official Charmbracelet docs + Go stdlib testing + direct source verification)

---

## 1. Testing Layers Overview

Abditum requires five testing categories per `descricao.md`:

| Layer | Type | Tool | What it validates |
|-------|------|------|-------------------|
| Crypto service | Unit | Go stdlib `testing` | Encrypt/Decrypt round-trips, wrong password, corrupted file |
| Storage service | Unit | Go stdlib `testing` + `os.TempDir()` | Atomic save, .bak creation, .tmp cleanup on failure |
| Domain/Manager | Unit (white-box) | Go stdlib `testing` | VaultManager state transitions, invariants |
| TUI screens | Visual golden | `teatest` | Each screen at 80×24, visual regression |
| Integration | E2E | `teatest` + real vault | Full workflow: create vault → create secret → edit → save → reopen |

---

## 2. `teatest` for Bubble Tea v2

### Import Path (v2 — CRITICAL)

```go
import "charm.land/bubbletea/v2/teatest"
```

> Do NOT use `github.com/charmbracelet/bubbletea/v2/teatest` — that is the old path.

### Key API

```go
// Create a test model with a fixed terminal size
tm := teatest.NewTestModel(t, initialModel,
    teatest.WithInitialTermSize(80, 24))

// Send key events
tm.Send(tea.KeyPressMsg{Code: tea.KeyDown})
tm.Send(tea.KeyPressMsg{Text: "q", Mod: tea.ModCtrl})

// Send a string as if the user typed it
tm.Type("my-vault-name")

// Wait until the output contains expected text
teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
    return strings.Contains(string(bts), "Abditum")
}, teatest.WithDuration(3*time.Second))

// Compare output to a golden file
tm.FinalOutput(t, teatest.WithGoldenFile("testdata/golden/unlock-screen.txt"))

// Quit cleanly
tm.Quit()
```

### Golden File Workflow

Golden files are stored at `testdata/golden/*.txt` and committed to the repository. They capture the exact terminal output (ANSI escape codes stripped for readability, or raw depending on config).

**Update golden files:**
```bash
go test ./... -update  # custom flag pattern — see below
```

```go
// In test file:
var update = flag.Bool("update", false, "update golden files")

func TestUnlockScreen(t *testing.T) {
    tm := teatest.NewTestModel(t, newUnlockModel(),
        teatest.WithInitialTermSize(80, 24))

    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return strings.Contains(string(bts), "Abditum")
    }, teatest.WithDuration(time.Second))

    if *update {
        tm.FinalOutput(t, teatest.WithGoldenFile("testdata/golden/unlock-screen.txt"))
    } else {
        out, _ := io.ReadAll(tm.Output())
        golden, _ := os.ReadFile("testdata/golden/unlock-screen.txt")
        if !bytes.Equal(out, golden) {
            t.Errorf("output does not match golden file")
        }
    }
}
```

---

## 3. Crypto Service Tests

```go
// domain/crypto/crypto_test.go

package crypto_test

import (
    "bytes"
    "testing"
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
    key := make([]byte, 32)  // all-zero key for testing
    plaintext := []byte(`{"versao":1,"segredos":[]}`)

    encrypted, err := Encrypt(key, plaintext)
    if err != nil {
        t.Fatalf("Encrypt failed: %v", err)
    }

    decrypted, err := Decrypt(key, encrypted)
    if err != nil {
        t.Fatalf("Decrypt failed: %v", err)
    }

    if !bytes.Equal(plaintext, decrypted) {
        t.Errorf("round-trip mismatch: got %q, want %q", decrypted, plaintext)
    }
}

func TestDecryptWrongKey(t *testing.T) {
    key1 := make([]byte, 32)
    key2 := make([]byte, 32)
    key2[0] = 1  // differ by one byte

    plaintext := []byte("secret data")
    encrypted, _ := Encrypt(key1, plaintext)

    _, err := Decrypt(key2, encrypted)
    if err == nil {
        t.Error("expected error with wrong key, got nil")
    }
}

func TestDecryptCorruptedData(t *testing.T) {
    key := make([]byte, 32)
    plaintext := []byte("secret data")

    encrypted, _ := Encrypt(key, plaintext)

    // Corrupt a byte in the middle of the ciphertext
    encrypted[len(encrypted)/2] ^= 0xFF

    _, err := Decrypt(key, encrypted)
    if err == nil {
        t.Error("expected error with corrupted data, got nil")
    }
}

func TestDecryptTruncatedData(t *testing.T) {
    key := make([]byte, 32)
    _, err := Decrypt(key, []byte{1, 2, 3})  // too short for nonce
    if err == nil {
        t.Error("expected error with truncated data, got nil")
    }
}

func TestArgon2idKeyDerivation(t *testing.T) {
    password := []byte("my-master-password")
    salt := make([]byte, 32)

    key1 := DeriveKey(password, salt)
    key2 := DeriveKey(password, salt)

    if !bytes.Equal(key1, key2) {
        t.Error("Argon2id must be deterministic: same inputs must produce same key")
    }
    if len(key1) != 32 {
        t.Errorf("expected 32-byte key, got %d bytes", len(key1))
    }

    // Different salt → different key
    differentSalt := make([]byte, 32)
    differentSalt[0] = 1
    key3 := DeriveKey(password, differentSalt)
    if bytes.Equal(key1, key3) {
        t.Error("different salt must produce different key")
    }
}

// TestArgon2idCost is a benchmark-style test that verifies derivation takes
// at least a minimum time (anti-regression for accidentally lowered parameters).
func TestArgon2idMinimumCost(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping Argon2id cost test in short mode")
    }

    password := []byte("test")
    salt := make([]byte, 32)

    start := time.Now()
    DeriveKey(password, salt)
    elapsed := time.Since(start)

    // On modern hardware, 64MB/3-pass should take at least 200ms.
    // If it takes less, parameters may have been accidentally lowered.
    if elapsed < 200*time.Millisecond {
        t.Errorf("Argon2id too fast (%v): parameters may be too weak", elapsed)
    }
}
```

---

## 4. Storage Service Tests

```go
// storage/storage_test.go

package storage_test

import (
    "os"
    "path/filepath"
    "testing"
)

func TestAtomicSave(t *testing.T) {
    dir := t.TempDir()  // t.TempDir() auto-cleans on test completion
    vaultPath := filepath.Join(dir, "test.abditum")

    data := []byte("encrypted vault data")
    if err := SaveVault(vaultPath, data); err != nil {
        t.Fatalf("SaveVault failed: %v", err)
    }

    // Verify the vault file exists and has correct content
    got, err := os.ReadFile(vaultPath)
    if err != nil {
        t.Fatalf("vault file not found after save: %v", err)
    }
    if string(got) != string(data) {
        t.Errorf("vault content mismatch")
    }

    // Verify no leftover .tmp file
    if _, err := os.Stat(vaultPath + ".tmp"); !os.IsNotExist(err) {
        t.Error("expected .tmp file to be cleaned up after successful save")
    }
}

func TestAtomicSaveCreatesBackup(t *testing.T) {
    dir := t.TempDir()
    vaultPath := filepath.Join(dir, "test.abditum")

    // First save — no backup yet
    if err := SaveVault(vaultPath, []byte("version 1")); err != nil {
        t.Fatalf("first save failed: %v", err)
    }

    // Second save — should create .bak of first version
    if err := SaveVault(vaultPath, []byte("version 2")); err != nil {
        t.Fatalf("second save failed: %v", err)
    }

    // Verify .bak contains the old content
    bak, err := os.ReadFile(vaultPath + ".bak")
    if err != nil {
        t.Fatalf(".bak file not created: %v", err)
    }
    if string(bak) != "version 1" {
        t.Errorf("backup content wrong: got %q, want %q", bak, "version 1")
    }

    // Verify current vault has new content
    current, _ := os.ReadFile(vaultPath)
    if string(current) != "version 2" {
        t.Errorf("current vault content wrong: got %q, want %q", current, "version 2")
    }
}

func TestLoadVault(t *testing.T) {
    dir := t.TempDir()
    vaultPath := filepath.Join(dir, "test.abditum")

    data := []byte("some encrypted bytes")
    os.WriteFile(vaultPath, data, 0600)

    loaded, err := LoadVault(vaultPath)
    if err != nil {
        t.Fatalf("LoadVault failed: %v", err)
    }
    if string(loaded) != string(data) {
        t.Errorf("loaded content mismatch")
    }
}

func TestLoadVaultNotFound(t *testing.T) {
    _, err := LoadVault("/nonexistent/path/test.abditum")
    if err == nil {
        t.Error("expected error for non-existent vault")
    }
}

func TestFilePermissions(t *testing.T) {
    dir := t.TempDir()
    vaultPath := filepath.Join(dir, "test.abditum")

    SaveVault(vaultPath, []byte("data"))

    info, err := os.Stat(vaultPath)
    if err != nil {
        t.Fatal(err)
    }

    mode := info.Mode().Perm()
    if mode != 0600 {
        t.Errorf("vault file permissions: got %o, want 0600", mode)
    }
}
```

---

## 5. Domain Manager Tests (White-Box)

```go
// domain/vault/manager_test.go

package vault_test

func TestCreateSecret(t *testing.T) {
    vault := InitializeNewVault()
    manager := NewVaultManager(vault)

    secret, err := manager.CreateSecret("GitHub", "Login", "")  // "" = vault root
    if err != nil {
        t.Fatalf("CreateSecret failed: %v", err)
    }
    if secret.Name() != "GitHub" {
        t.Errorf("got name %q, want %q", secret.Name(), "GitHub")
    }
    if !manager.IsDirty() {
        t.Error("expected manager to be dirty after CreateSecret")
    }
}

func TestDeleteSecretGoesToTrash(t *testing.T) {
    vault := InitializeNewVault()
    manager := NewVaultManager(vault)

    secret, _ := manager.CreateSecret("My Secret", "", "")
    id := secret.ID()

    if err := manager.DeleteSecret(id); err != nil {
        t.Fatalf("DeleteSecret failed: %v", err)
    }

    // Secret should now be in trash
    trashed := CollectTrashed(vault)
    found := false
    for _, s := range trashed {
        if s.ID() == id {
            found = true
        }
    }
    if !found {
        t.Error("deleted secret not found in trash")
    }
}

func TestRestoreSecret(t *testing.T) {
    vault := InitializeNewVault()
    manager := NewVaultManager(vault)

    secret, _ := manager.CreateSecret("My Secret", "", "")
    id := secret.ID()

    manager.DeleteSecret(id)
    if err := manager.RestoreSecret(id); err != nil {
        t.Fatalf("RestoreSecret failed: %v", err)
    }

    trashed := CollectTrashed(vault)
    for _, s := range trashed {
        if s.ID() == id {
            t.Error("restored secret still in trash")
        }
    }
}

func TestFavoriteToggle(t *testing.T) {
    vault := InitializeNewVault()
    manager := NewVaultManager(vault)

    secret, _ := manager.CreateSecret("My Secret", "", "")
    id := secret.ID()

    manager.FavoriteSecret(id, true)
    favs := CollectFavorites(vault)
    if len(favs) != 1 || favs[0].ID() != id {
        t.Error("secret not in favorites after favoriting")
    }

    manager.FavoriteSecret(id, false)
    favs = CollectFavorites(vault)
    if len(favs) != 0 {
        t.Error("secret still in favorites after un-favoriting")
    }
}

func TestSearch(t *testing.T) {
    vault := InitializeNewVault()
    manager := NewVaultManager(vault)

    manager.CreateSecret("GitHub Login", "Login", "")
    manager.CreateSecret("Gmail", "Login", "")

    results := Search(vault, "git")
    if len(results) != 1 || results[0].Secret.Name() != "GitHub Login" {
        t.Errorf("unexpected search results: %v", results)
    }
}

func TestSearchSensitiveFieldsExcluded(t *testing.T) {
    vault := InitializeNewVault()
    manager := NewVaultManager(vault)

    secret, _ := manager.CreateSecret("My Bank", "Login", "")
    manager.EditSecretField(secret.ID(), "Password", "hunter2")

    // Searching for the password value should NOT find it
    results := Search(vault, "hunter2")
    if len(results) != 0 {
        t.Error("search should not return results for sensitive field values")
    }
}
```

---

## 6. TUI Screen Golden File Tests

### Directory Structure

```
tui/
├── testdata/
│   └── golden/
│       ├── auth-unlock.txt         # Lock screen / unlock prompt
│       ├── auth-create-vault.txt   # New vault creation screen
│       ├── main-empty-vault.txt    # Main screen with empty vault
│       ├── main-with-secrets.txt   # Main screen with populated tree
│       ├── detail-secret.txt       # Secret detail panel
│       ├── detail-edit-secret.txt  # Secret edit mode
│       ├── modal-confirm-delete.txt # Confirmation dialog
│       ├── modal-new-secret.txt    # New secret form
│       ├── sidebar-search.txt      # Tree filtered by search
│       └── ...
```

### Test Pattern

```go
// tui/auth_test.go

func TestUnlockScreen_Golden(t *testing.T) {
    tm := teatest.NewTestModel(t,
        newAuthModel(),
        teatest.WithInitialTermSize(80, 24),
    )

    // Wait for initial render
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("Abditum"))
    }, teatest.WithDuration(time.Second))

    tm.FinalOutput(t, teatest.WithGoldenFile("testdata/golden/auth-unlock.txt"))
    tm.Quit()
}

func TestUnlockScreen_KeyboardNavigation(t *testing.T) {
    tm := teatest.NewTestModel(t,
        newAuthModel(),
        teatest.WithInitialTermSize(80, 24),
    )

    // Type a password
    tm.Type("my-test-password")

    // Verify the password field shows masked input
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("●"))  // masked char
    }, teatest.WithDuration(time.Second))

    tm.Quit()
}
```

---

## 7. Integration Tests (Full Workflow)

Integration tests use a real temporary vault file and exercise the full stack: TUI → Manager → Crypto → Storage.

```go
// integration/integration_test.go

package integration_test

import (
    "os"
    "path/filepath"
    "testing"
)

// TestFullWorkflow exercises the complete create → populate → save → reopen flow.
func TestFullWorkflow(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode (Argon2id is slow)")
    }

    dir := t.TempDir()
    vaultPath := filepath.Join(dir, "test.abditum")

    // Step 1: Create a new vault
    vault := InitializeNewVault()
    manager := NewVaultManager(vault)

    // Step 2: Add a secret
    secret, err := manager.CreateSecret("GitHub", "Login", "")
    if err != nil {
        t.Fatalf("CreateSecret: %v", err)
    }
    manager.EditSecretField(secret.ID(), "URL", "https://github.com")
    manager.EditSecretField(secret.ID(), "Username", "testuser")
    manager.EditSecretField(secret.ID(), "Password", "s3cr3t!")

    // Step 3: Save the vault with a master password
    password := []byte("integration-test-password")
    salt, _ := GenerateSalt()
    key := DeriveKey(password, salt)

    jsonData, _ := json.Marshal(vault)
    encrypted, _ := Encrypt(key, jsonData)

    header := FileHeader{Version: FormatVersion, Salt: salt}
    var buf bytes.Buffer
    WriteHeader(&buf, header)
    buf.Write(encrypted)

    if err := SaveVault(vaultPath, buf.Bytes()); err != nil {
        t.Fatalf("SaveVault: %v", err)
    }

    // Step 4: Reopen the vault
    fileData, _ := os.ReadFile(vaultPath)
    reader := bytes.NewReader(fileData)
    hdr, _ := ReadHeader(reader)

    reopenedKey := DeriveKey(password, hdr.Salt)
    encryptedPayload, _ := io.ReadAll(reader)
    decrypted, err := Decrypt(reopenedKey, encryptedPayload)
    if err != nil {
        t.Fatalf("Decrypt on reopen: %v", err)
    }

    var reopenedVault Vault
    if err := json.Unmarshal(decrypted, &reopenedVault); err != nil {
        t.Fatalf("json.Unmarshal on reopen: %v", err)
    }

    // Step 5: Verify secrets survived the round-trip
    if len(reopenedVault.Secrets()) != 1 {
        t.Errorf("expected 1 secret after reopen, got %d", len(reopenedVault.Secrets()))
    }
    if reopenedVault.Secrets()[0].Name() != "GitHub" {
        t.Errorf("secret name mismatch after reopen")
    }
}
```

---

## 8. CI Pipeline Design

### GitHub Actions (`.github/workflows/ci.yml`)

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    name: Test (${{ matrix.os }})
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
          cache: true

      - name: Verify go.mod is tidy
        run: |
          go mod tidy
          git diff --exit-code go.mod go.sum

      - name: Build (current platform)
        run: CGO_ENABLED=0 go build -trimpath ./...

      - name: Test (short — fast tests only)
        run: go test -short -count=1 ./...

      - name: Test (full — with Argon2id cost tests)
        run: go test -count=1 -timeout 300s ./...
        if: github.event_name == 'push'  # only on merge, not every PR

      - name: Test with race detector
        run: go test -race -short ./...

  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - uses: golangci/golangci-lint-action@v6
        with:
          version: latest

  build-matrix:
    name: Cross-compile verification
    runs-on: ubuntu-latest
    strategy:
      matrix:
        include:
          - goos: linux
            goarch: amd64
          - goos: linux
            goarch: arm64
          - goos: darwin
            goarch: amd64
          - goos: darwin
            goarch: arm64
          - goos: windows
            goarch: amd64
          - goos: windows
            goarch: arm64
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - name: Cross-compile ${{ matrix.goos }}/${{ matrix.goarch }}
        run: CGO_ENABLED=0 GOOS=${{ matrix.goos }} GOARCH=${{ matrix.goarch }} go build -trimpath ./cmd/abditum
        env:
          CGO_ENABLED: 0
          GOOS: ${{ matrix.goos }}
          GOARCH: ${{ matrix.goarch }}
```

### golangci-lint Configuration (`.golangci.yml`)

```yaml
linters:
  enable:
    - errcheck       # Check all error returns handled
    - gosimple       # Suggest simpler code
    - govet          # `go vet` checks
    - ineffassign    # Detect ineffectual assignments
    - staticcheck    # Comprehensive static analysis
    - unused         # Detect unused code
    - gosec          # Security-focused linter (important for a password vault)
    - gofmt          # Format check

linters-settings:
  gosec:
    excludes:
      - G304  # File path provided as taint input — acceptable for vault file
```

---

## 9. Test Run Commands

```makefile
# Add to Makefile:

## test: Run all tests (short mode — fast)
test:
	go test -short ./...

## test-full: Run all tests including Argon2id cost tests
test-full:
	go test -timeout 300s ./...

## test-race: Run tests with race detector
test-race:
	go test -race -short ./...

## test-update-golden: Update golden files
test-update-golden:
	go test ./tui/... -update

## test-integration: Run integration tests only
test-integration:
	go test ./integration/... -timeout 300s
```

---

## 10. Testing Pitfalls

### Pitfall 1: Argon2id Makes Tests Slow
**Problem:** Full Argon2id derivation (64MB, 3 passes) takes 1–3 seconds per test. If every test creates a vault, CI becomes very slow.  
**Prevention:** 
- Use `testing.Short()` guard for Argon2id tests
- In unit tests for storage/domain, use a pre-derived test key (all-zero bytes or a mock crypto service)
- Only integration tests do real Argon2id

### Pitfall 2: Golden File Flakiness from Terminal Width
**Problem:** Golden files capture ANSI-escaped output. If terminal width differs between test runs, the layout breaks and tests fail.  
**Prevention:** Always use `teatest.WithInitialTermSize(80, 24)` — fixed 80×24 for all golden file tests.

### Pitfall 3: teatest Race Conditions
**Problem:** `WaitFor` with too short a timeout causes flaky tests on slow CI machines.  
**Prevention:** Use `teatest.WithDuration(3*time.Second)` minimum. For integration tests, use 10 seconds.

### Pitfall 4: Windows Path Separators in Tests
**Problem:** Hardcoded `/` path separators fail on Windows.  
**Prevention:** Always use `filepath.Join()` and `filepath.Dir()` in tests. The CI matrix tests on `windows-latest` catches this.

### Pitfall 5: t.Parallel() with TempDir
**Problem:** `t.TempDir()` is safe with parallel tests, but file-based tests may interfere if they use hardcoded paths.  
**Prevention:** Always use `t.TempDir()` — it's unique per test and auto-cleaned.

---

## 11. Key Decisions Summary

| Decision | Choice | Rationale |
|----------|--------|-----------|
| TUI testing | `teatest` (official library) | Only viable option for Bubble Tea |
| Golden file size | 80×24 | Standard terminal size; specified in descricao.md |
| Slow tests | Guard with `testing.Short()` | Argon2id must not block every `go test` run |
| CI matrix | ubuntu + windows + macos | Portability must be validated on all three |
| Cross-compile verification | Separate CI job | Catches import/build errors early, before goreleaser |
| Lint | golangci-lint + gosec | gosec catches security-specific anti-patterns |

---

## 12. Sources

| Source | Confidence | Notes |
|--------|------------|-------|
| `charm.land/bubbletea/v2/teatest` (official package) | HIGH | Official Charmbracelet testing library for Bubble Tea |
| Go stdlib `testing` package documentation | HIGH | Standard library |
| `os.TempDir()` / `t.TempDir()` (Go 1.15+) | HIGH | Auto-cleanup per test, widely used |
| `golangci/golangci-lint-action` GitHub Action | HIGH | Official golangci-lint CI integration |
| `gosec` linter documentation | MEDIUM | Security linter for Go — well-established |
| `descricao.md` — Testes section | HIGH | Primary source for test requirements |
