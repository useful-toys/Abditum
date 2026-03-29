# Phase 2: Crypto Package - Context

**Gathered:** 2026-03-29
**Status:** Ready for planning

<domain>
## Phase Boundary

`internal/crypto` delivers production-ready Argon2id key derivation, AES-256-GCM authenticated encryption, secure memory primitives (zero, mlock/VirtualLock), and password strength evaluation — all verified by tests that will catch any future security regression.

This phase implements cryptographic primitives only — no vault logic, no file I/O, no TUI. Pure crypto functions with comprehensive test coverage.

</domain>

<decisions>
## Implementation Decisions

### Error handling strategy
- **D-01:** Use sentinel errors exported as package-level vars (ErrAuthFailed, ErrInsufficientEntropy, ErrInvalidParams, ErrMLockFailed). Clearest contract, easiest to document and test with errors.Is().
- **D-02:** Decrypt returns single ErrAuthFailed for both wrong password and corrupted ciphertext (no timing leak).
- **D-03:** mlock/VirtualLock failure returns error but continues execution (non-fatal). Caller logs warning, operation proceeds. Required for cross-platform compatibility (containers, low-privilege contexts).
- **D-04:** Panic only on programming errors (wrong key size, nil pointer). Runtime errors (auth failure, entropy exhaustion, invalid params) return error.
- **D-05:** Error messages contain minimal details ('authentication failed', 'insufficient entropy') with no internal failure modes exposed.
- **D-06:** crypto/rand exhaustion during nonce generation returns ErrInsufficientEntropy (not panic).
- **D-07:** Every exported function documents which sentinel errors it can return in godoc.
- **D-08:** Unit tests verify exact sentinel error identity with errors.Is() to catch contract breaking changes.
- **D-09:** Invalid parameters (nil password, zero-length salt, invalid Argon2id params) return ErrInvalidParams.
- **D-10:** Crypto package never wraps errors — returns only sentinels. Caller adds context if needed.
- **D-11:** Export all sentinels as package-level vars for maximum caller flexibility.
- **D-12:** When stdlib crypto functions return errors (gcm.Open, gcm.Seal), wrap in appropriate sentinel and return (hides stdlib implementation details).

### Memory safety patterns
- **D-13:** All sensitive data (keys, passwords, plaintext) handled as `[]byte` — never as `string`.
- **D-14:** Caller owns and zeros all `[]byte` buffers after use. Crypto functions do not retain references.
- **D-15:** Explicit zero with ZeroBytes() helper — no implicit defer patterns (caller controls timing).
- **D-16:** mlock/VirtualLock implemented with platform-specific build tags (mlock_unix.go, mlock_windows.go, mlock_other.go).
- **D-17:** Pre-allocate exact buffer sizes — never append to locked buffers (invalidates mlock).

### API surface decisions
- **D-18:** Separate functions for each primitive: GenerateSalt(), DeriveKey(), GenerateNonce(), Encrypt(), Decrypt(). No high-level combined operations (keeps API flexible).
- **D-19:** Nonce generation is internal to Encrypt() — called immediately before gcm.Seal() to ensure uniqueness.
- **D-20:** Argon2id parameters passed via ArgonParams struct (t, m, p, keyLen fields).
- **D-21:** FormatVersion constant exported (value: 1) to match formato-arquivo-abditum.md.
- **D-22:** All core crypto primitives are exported (public API). Helper functions (internal validation) are internal.

### Testing and verification
- **D-23:** Comprehensive unit test suite: nonce uniqueness (encrypt twice → distinct ciphertexts), roundtrip equality, wrong-key → ErrAuthFailed, short ciphertext → error, ZeroBytes fills entire slice.
- **D-24:** Password strength boundary tests: 11 chars = Weak, 12 chars with all categories = Strong, missing any category = Weak.
- **D-25:** All tests run with race detector (`go test -race`).
- **D-26:** Deterministic test fixtures where possible (fixed salts/nonces for decrypt tests). Randomized for uniqueness tests.
- **D-27:** No fuzzing or property-based tests in Phase 2 (can be added later if needed).
- **D-28:** Benchmarks included for DeriveKey (Argon2id performance baseline) and Encrypt/Decrypt (AES-GCM throughput).

### Crypto parameters (from formato-arquivo-abditum.md)
- **D-29:** Argon2id parameters for format version 1: m=262144 KiB (256 MiB), t=3, p=4, keyLen=32 bytes.
- **D-30:** Salt size: 32 bytes (generated with crypto/rand at vault creation, replaced only on master password change).
- **D-31:** Nonce size: 12 bytes (regenerated with crypto/rand on every save operation).
- **D-32:** AES-256-GCM: 32-byte key, 12-byte nonce (NIST standard), 16-byte authentication tag.
- **D-33:** Magic bytes: "ABDT" (4 bytes ASCII) — not handled in crypto package (storage layer responsibility).
- **D-34:** Password strength: Strong requires ≥12 chars AND ≥1 uppercase AND ≥1 lowercase AND ≥1 digit AND ≥1 special char (PWD-01).

### the agent's Discretion
- Internal helper function organization
- Exact test fixture values (as long as they demonstrate the tested property)
- Benchmark iteration counts
- Comment verbosity (must be generous per project policy)

</decisions>

<specifics>
## Specific Ideas

- "I can't answer crypto implementation questions — use recommended best practices and ask if you need decisions."
- Check formato-arquivo-abditum.md for all crypto parameters — they're fully specified there.

</specifics>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Crypto specifications
- `formato-arquivo-abditum.md` — Complete binary format spec with Argon2id params (m=262144 KiB, t=3, p=4), AES-256-GCM layout, nonce/salt sizes, error categories, and format versioning strategy
- `.planning/REQUIREMENTS.md` §CRYPTO-01 through §CRYPTO-06 — All crypto requirements: AES-256-GCM, Argon2id KDF, dependency restrictions (stdlib + golang.org/x/crypto only), []byte-only sensitive data, memory wipe on lock/exit, mlock/VirtualLock, zero logs with sensitive content
- `.planning/REQUIREMENTS.md` §PWD-01 — Password strength evaluation criteria (≥12 chars, character class requirements)
- `.planning/ROADMAP.md` §Phase 2 — Complete plan breakdown, UAT criteria, and pitfall watch (Argon2id m parameter in KiB, never string for sensitive data, nonce uniqueness, mlock best-effort)

### Architecture and conventions
- `arquitetura.md` — Import path restrictions, dependency policy (no net, no math/rand), generous comment policy for non-Go/crypto experts
- `.planning/phases/01-project-scaffold-ci-foundation/01-CONTEXT.md` — Phase 1 decisions (module path, Go version, CI configuration)

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/crypto/doc.go` — Package stub exists, ready to be populated with package-level documentation
- `go.mod` — All required dependencies already present: `golang.org/x/crypto` (Argon2id), `golang.org/x/sys` (mlock/VirtualLock)

### Established Patterns
- Static linking enforced (`CGO_ENABLED=0`) — verified by Phase 1 CI
- Canonical Charm import paths (`charm.land/*`) — established in Phase 1
- Generous comment policy — code must be accessible to readers less familiar with Go, Bubble Tea, and cryptography

### Integration Points
- This package will be consumed by `internal/storage` (Phase 4) for file encryption/decryption
- Password strength evaluator will be used by TUI welcome screen (Phase 6) for real-time feedback
- Memory primitives (ZeroBytes, mlock) will be used by `internal/vault` Manager (Phase 3) for sensitive data cleanup

</code_context>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 02-crypto-package*
*Context gathered: 2026-03-29*
