---
phase: 04
slug: storage-package
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-30
---

# Phase 04 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test |
| **Config file** | none — standard Go test tooling |
| **Quick run command** | `CGO_ENABLED=0 go test ./internal/storage/... -count=1` |
| **Full suite command** | `CGO_ENABLED=0 go test ./internal/... -count=1 -race` |
| **Estimated runtime** | ~5 seconds (excludes Argon2id with full params — tests use fast params) |

---

## Sampling Rate

- **After every task commit:** Run `CGO_ENABLED=0 go test ./internal/storage/... -count=1`
- **After every plan wave:** Run `CGO_ENABLED=0 go test ./internal/... -count=1 -race`
- **Before `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** 10 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 04-01-01 | 01 | 1 | CRYPTO ext | unit | `go test ./internal/crypto/... -run AAD -count=1` | ❌ W0 | ⬜ pending |
| 04-01-02 | 01 | 1 | vault serial | unit | `go test ./internal/vault/... -run Serializ -count=1` | ❌ W0 | ⬜ pending |
| 04-02-01 | 02 | 2 | ATOMIC-01,02,03 | unit | `go test ./internal/storage/... -run Save -count=1` | ❌ W0 | ⬜ pending |
| 04-02-02 | 02 | 2 | ATOMIC-04 | unit | `go test ./internal/storage/... -run AtomicRename -count=1` | ❌ W0 | ⬜ pending |
| 04-03-01 | 03 | 2 | ATOMIC-01,02 | unit | `go test ./internal/storage/... -run Backup -count=1` | ❌ W0 | ⬜ pending |
| 04-03-02 | 03 | 2 | COMPAT-03 | unit | `go test ./internal/storage/... -run Migra -count=1` | ❌ W0 | ⬜ pending |
| 04-04-01 | 04 | 3 | all | integration | `go test ./internal/storage/... -run Integration -count=1` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] Test infrastructure created within each plan's first task (TDD approach)

*Existing Go test infrastructure covers framework needs. Test files created alongside implementation.*

---

## Manual-Only Verifications

*All phase behaviors have automated verification.*

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 10s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
