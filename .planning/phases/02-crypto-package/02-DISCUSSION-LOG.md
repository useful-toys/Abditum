# Phase 2: Crypto Package - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-29
**Phase:** 2-crypto-package
**Areas discussed:** Error handling strategy

---

## Error handling strategy

### Q1: How should crypto package errors be structured?

| Option | Description | Selected |
|--------|-------------|----------|
| Sentinel errors | Package exports typed sentinel errors (e.g., ErrAuthFailed, ErrInvalidKeySize) that callers can compare with errors.Is(). Clearest contract, easiest to document. Recommended for crypto where callers need to distinguish authentication failure from corruption. | ✓ |
| Typed error structs | Package exports error types implementing error interface with additional fields (Code, Details). Allows structured error info but requires type assertions. More flexible but harder to use correctly. | |
| Wrapped errors | Package uses error wrapping with fmt.Errorf. Easier to add context but harder for callers to detect specific errors. Less idiomatic for library code. | |

**User's choice:** Sentinel errors (Recommended)

---

### Q2: When GCM authentication fails, what should Decrypt return?

| Option | Description | Selected |
|--------|-------------|----------|
| Single sentinel | Return generic ErrAuthFailed for both wrong password and corrupted ciphertext. Prevents timing attacks by treating both cases identically. Recommended for password managers where an attacker shouldn't distinguish wrong password from file corruption. | ✓ |
| Separate errors | Return distinct errors: ErrAuthFailed (wrong password) vs ErrCorrupted (invalid ciphertext structure). More debuggable but may leak timing information. | |
| Contextual errors | Return error with additional context string (e.g., 'authentication failed: tag mismatch'). Most debuggable but highest information leak risk. | |

**User's choice:** Single sentinel (Recommended)

---

### Q3: How should memory locking (mlock/VirtualLock) failures be handled?

| Option | Description | Selected |
|--------|-------------|----------|
| Non-fatal | mlock failure returns error but continues execution. Caller logs a warning, operation proceeds. Recommended for cross-platform code where mlock availability varies (not available in containers, low-privilege contexts). | ✓ |
| Fatal error | mlock failure returns fatal error, operation aborts. Strongest security posture but breaks in many real-world environments (Docker, restricted users). | |
| Silent ignore | mlock failure is silently ignored. Operation continues, no error returned. Simplest but provides no visibility into security degradation. | |

**User's choice:** Non-fatal (Recommended)

---

### Q4: When should crypto functions panic vs return error?

| Option | Description | Selected |
|--------|-------------|----------|
| Programming errors only | Panic only on programming errors (wrong key size passed to AES constructor, nil pointer). Runtime errors (auth failure, entropy exhaustion) return error. Recommended — makes API contract clear and testable. | ✓ |
| Never panic | Never panic, always return error. Most defensive but makes some bugs harder to catch during development (e.g., passing 16-byte key to AES-256). | |
| Panic liberally | Panic on any unrecoverable condition including entropy failures. Simplest error handling but can crash production unexpectedly. | |

**User's choice:** Programming errors only (Recommended)

---

### Q5: How much detail should error messages contain?

| Option | Description | Selected |
|--------|-------------|----------|
| Minimal | Error messages contain no details about failure mode ('authentication failed' only). Hardest to debug but safest against information leaks. Recommended for security-critical packages. | ✓ |
| Operation context | Error messages include operation context but no internal details ('encryption failed', 'key derivation failed'). Balance between debuggability and security. | |
| Full details | Error messages include full technical details ('GCM tag mismatch at byte 16', 'Argon2id out of memory'). Easiest to debug but may leak sensitive information. | |

**User's choice:** Minimal (Recommended)

---

### Q6: What should happen if crypto/rand fails during nonce generation?

| Option | Description | Selected |
|--------|-------------|----------|
| Return error | crypto/rand exhaustion returns error (ErrInsufficientEntropy). Allows caller to retry or fail gracefully. Recommended for most applications. | ✓ |
| Panic | crypto/rand exhaustion panics. Unrecoverable state, simplest to handle but crashes program. Appropriate if entropy failure indicates system compromise. | |
| Fallback (DO NOT USE) | Fall back to deterministic nonce generation with counter. Dangerous — nonce reuse breaks GCM security completely. Never use this. | |

**User's choice:** Return error (Recommended)

---

### Q7: Should error returns be explicitly documented in godoc?

| Option | Description | Selected |
|--------|-------------|----------|
| Document all | Every exported function documents which sentinel errors it can return in godoc. Callers know exactly what to check for. Recommended for library code. | ✓ |
| Document common only | Only document the 'expected' errors (auth failure, wrong key). Document rare errors (entropy failure) only in package doc. Less cluttered but less discoverable. | |
| No documentation | No explicit error documentation, rely on code inspection. Simplest but poorest developer experience. | |

**User's choice:** Document all (Recommended)

---

### Q8: How should error returns be tested?

| Option | Description | Selected |
|--------|-------------|----------|
| Test identity | Unit tests verify exact sentinel error identity with errors.Is(). Catches error contract breaking changes. Recommended for stable API. | ✓ |
| Test presence only | Unit tests verify error is non-nil only. Faster to write but doesn't validate error contract. | |
| Test messages | Unit tests also verify error message strings. Most brittle but catches message regressions. | |

**User's choice:** Test identity (Recommended)

---

### Q9: How should parameter validation failures be handled?

| Option | Description | Selected |
|--------|-------------|----------|
| Return error | Invalid parameters (nil password, zero-length salt, invalid Argon2id params) return ErrInvalidParams. Clear separation between programmer error and runtime error. Recommended for production robustness. | ✓ |
| Panic | Invalid parameters panic. Fails fast during development but can crash production. Appropriate if all parameters are constants or validated upstream. | |
| No validation | Assume caller validates, no checks. Fastest but unsafe — garbage-in-garbage-out. | |

**User's choice:** Return error (Recommended)

---

### Q10: Should crypto package wrap errors with context?

| Option | Description | Selected |
|--------|-------------|----------|
| Never wrap | Crypto package returns only sentinel errors, never wraps. Caller adds context if needed with fmt.Errorf. Cleanest separation, easiest to test. Recommended for low-level libraries. | ✓ |
| Wrap with context | Crypto functions wrap sentinels with operation context ('encrypt: %w'). More context in stack traces but harder to test. | |
| Wrap with params | Crypto functions wrap AND add parameters to error message ('encrypt(32-byte key): %w'). Most debuggable but may leak key size info. | |

**User's choice:** Never wrap (Recommended)

---

### Q11: Which errors should be exported from crypto package?

| Option | Description | Selected |
|--------|-------------|----------|
| Export all | Export all sentinels as package-level vars (ErrAuthFailed, ErrInsufficientEntropy, ErrInvalidParams, ErrMLockFailed). Allows callers to check any condition. Recommended for flexibility. | ✓ |
| Export critical only | Export only 'important' errors (ErrAuthFailed), keep others internal. Simpler public API but forces callers to treat all failures identically. | |
| Constructor funcs | Export error constructor functions instead of vars (NewAuthError()). More flexible (can add fields later) but unusual pattern in Go. | |

**User's choice:** Export all (Recommended)

---

### Q12: How should errors from stdlib crypto functions be handled?

| Option | Description | Selected |
|--------|-------------|----------|
| Wrap in sentinel | When stdlib crypto functions return errors (gcm.Open, gcm.Seal), wrap in sentinel and return. Hides stdlib details, stable API. Recommended for encapsulation. | ✓ |
| Return directly | Return stdlib errors directly. Simplest but exposes stdlib implementation details in API contract. | |
| Selective wrapping | Catch specific stdlib errors and convert, pass through others. Most flexible but requires maintenance when stdlib changes. | |

**User's choice:** Wrap in sentinel (Recommended)

---

## the agent's Discretion

**Memory safety patterns, API surface decisions, Testing and verification, Constants and configuration** — User deferred all remaining decisions to implementation agent with instruction: "Use recommended best practices. Ask if decision is needed. Check formato-arquivo-abditum.md for crypto parameters."

All crypto parameters (Argon2id m/t/p, salt/nonce sizes, password strength criteria) are fully specified in formato-arquivo-abditum.md and incorporated into CONTEXT.md as decisions D-29 through D-34.

## Deferred Ideas

None — discussion stayed within phase scope
