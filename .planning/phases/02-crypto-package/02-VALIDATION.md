---
phase: 2
slug: crypto-package
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-29
---

# Phase 2 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test (stdlib testing package) |
| **Config file** | none — stdlib only |
| **Quick run command** | `go test ./internal/crypto -run TestUnit -v` |
| **Full suite command** | `go test ./internal/crypto/... -race -count=1` |
| **Estimated runtime** | ~5 seconds (Argon2id derivation dominates) |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/crypto -run TestUnit -v`
- **After every plan wave:** Run `go test ./internal/crypto/... -race -count=1`
- **Before `/gsd-verify-work`:** Full suite must be green with race detector
- **Max feedback latency:** 5 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 02-01-01 | 01 | 1 | CRYPTO-01 | unit | `go test ./internal/crypto -run TestDeriveKey` | ⬜ W0 | ⬜ pending |
| 02-01-02 | 01 | 1 | CRYPTO-01 | unit | `go test ./internal/crypto -run TestGenerateSalt` | ⬜ W0 | ⬜ pending |
| 02-01-03 | 01 | 1 | CRYPTO-01 | unit | `go test ./internal/crypto -run TestEncrypt` | ⬜ W0 | ⬜ pending |
| 02-01-04 | 01 | 1 | CRYPTO-01 | unit | `go test ./internal/crypto -run TestDecrypt` | ⬜ W0 | ⬜ pending |
| 02-01-05 | 01 | 1 | CRYPTO-01, CRYPTO-04 | unit | `go test ./internal/crypto -run TestNonceUniqueness` | ⬜ W0 | ⬜ pending |
| 02-01-06 | 01 | 1 | CRYPTO-03, CRYPTO-04 | unit | `go test ./internal/crypto -run TestZeroBytes` | ⬜ W0 | ⬜ pending |
| 02-01-07 | 01 | 1 | CRYPTO-05 | unit | `go test ./internal/crypto -run TestMlock` | ⬜ W0 | ⬜ pending |
| 02-01-08 | 01 | 1 | PWD-01 | unit | `go test ./internal/crypto -run TestPasswordStrength` | ⬜ W0 | ⬜ pending |
| 02-01-09 | 01 | 1 | ALL | integration | `go test ./internal/crypto -run TestRoundtrip` | ⬜ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/crypto/kdf_test.go` — stubs for CRYPTO-01 (GenerateSalt, DeriveKey)
- [ ] `internal/crypto/aead_test.go` — stubs for CRYPTO-01 (Encrypt, Decrypt, nonce uniqueness)
- [ ] `internal/crypto/memory_test.go` — stubs for CRYPTO-03, CRYPTO-04, CRYPTO-05 (ZeroBytes, Mlock)
- [ ] `internal/crypto/password_test.go` — stubs for PWD-01 (password strength evaluation)
- [ ] `internal/crypto/crypto_test.go` — integration test stub (full roundtrip)

**Wave 0 task pattern:**
```go
// File: internal/crypto/kdf_test.go
package crypto_test

import "testing"

func TestGenerateSalt(t *testing.T) {
    t.Skip("STUB — Wave 0 scaffold")
}

func TestDeriveKey(t *testing.T) {
    t.Skip("STUB — Wave 0 scaffold")
}
```

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Argon2id derivation time | CRYPTO-01 | Performance baseline check | Run `go test -bench=BenchmarkDeriveKey`; verify 200-500ms per operation on modern CPU |
| Memory locking availability | CRYPTO-05 | Platform-dependent | Run `go test -run TestMlock -v` on Unix, Windows, and "other" build tags; verify Unix/Windows succeed, other returns `ErrMLockFailed` |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 5s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
