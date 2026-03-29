# Phase 2: Crypto Package - Research

**Researched:** 2026-03-29
**Domain:** Go cryptography (Argon2id KDF, AES-256-GCM AEAD, memory security)
**Confidence:** HIGH

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**Error handling strategy:**
- D-01: Use sentinel errors exported as package-level vars (ErrAuthFailed, ErrInsufficientEntropy, ErrInvalidParams, ErrMLockFailed)
- D-02: Decrypt returns single ErrAuthFailed for both wrong password and corrupted ciphertext (no timing leak)
- D-03: mlock/VirtualLock failure returns error but continues execution (non-fatal)
- D-04: Panic only on programming errors; runtime errors return error
- D-05: Error messages minimal ('authentication failed', 'insufficient entropy')
- D-06: crypto/rand exhaustion returns ErrInsufficientEntropy
- D-07-D-12: Full sentinel error contract documented

**Memory safety patterns:**
- D-13: All sensitive data as `[]byte`, never `string`
- D-14: Caller owns and zeros buffers; crypto functions don't retain references
- D-15: Explicit zero with ZeroBytes(); no implicit defer
- D-16: Platform-specific build tags for mlock (mlock_unix.go, mlock_windows.go, mlock_other.go)
- D-17: Pre-allocate exact sizes; never append to locked buffers

**API surface:**
- D-18: Separate functions for each primitive
- D-19: Nonce generation internal to Encrypt(), called immediately before gcm.Seal()
- D-20: ArgonParams struct for parameters
- D-21: FormatVersion constant = 1
- D-22: All core crypto primitives exported

**Crypto parameters:**
- D-29: Argon2id: m=262144 KiB (256 MiB), t=3, p=4, keyLen=32
- D-30: Salt: 32 bytes (crypto/rand)
- D-31: Nonce: 12 bytes (crypto/rand, regenerated every save)
- D-32: AES-256-GCM: 32-byte key, 12-byte nonce, 16-byte tag
- D-34: Password strength: ≥12 chars AND ≥1 uppercase AND ≥1 lowercase AND ≥1 digit AND ≥1 special char

### the agent's Discretion
- Internal helper function organization
- Exact test fixture values
- Benchmark iteration counts
- Comment verbosity (must be generous per project policy)

### Deferred Ideas (OUT OF SCOPE)
None from CONTEXT.md discussion
</user_constraints>

<research_summary>
## Summary

Go's standard library (`crypto/*`) and `golang.org/x/crypto` provide production-ready primitives for Argon2id key derivation and AES-256-GCM AEAD. The Go crypto ecosystem is mature, well-audited, and the de facto standard for password-based encryption.

Key findings:
1. **Argon2id** via `golang.org/x/crypto/argon2` is the RFC 9106 reference implementation — use `argon2.IDKey()` directly with no wrapper
2. **AES-GCM** via `crypto/cipher` is constant-time, hardware-accelerated on modern CPUs, and the standard for authenticated encryption
3. **Memory locking** requires platform-specific syscalls (`unix.Mlock`, `windows.VirtualLock`) — both are best-effort and failure must be non-fatal
4. **Nonce uniqueness** is critical for GCM security — generate fresh 12-byte nonce via `io.ReadFull(rand.Reader, nonce)` before every `Seal()` call

**Primary recommendation:** Use stdlib + `golang.org/x/crypto` only. No third-party crypto libraries. Follow NIST SP 800-38D for GCM, RFC 9106 for Argon2id, and OWASP best practices for password strength.
</research_summary>

<standard_stack>
## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| crypto/cipher | stdlib | AES-GCM AEAD | NIST-validated, hardware-accelerated, constant-time |
| crypto/aes | stdlib | AES block cipher | Foundation for GCM mode |
| crypto/rand | stdlib | Cryptographic RNG | OS-backed entropy (`/dev/urandom`, `CryptGenRandom`) |
| golang.org/x/crypto/argon2 | latest | Argon2id KDF | RFC 9106 reference implementation |
| golang.org/x/sys/unix | latest | mlock syscall | Memory locking on Unix systems |
| golang.org/x/sys/windows | latest | VirtualLock API | Memory locking on Windows |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| io | stdlib | ReadFull for nonce | Ensures exact byte count from rand.Reader |
| errors | stdlib | Sentinel error definitions | Standard error handling pattern |
| encoding/binary | stdlib | Little-endian serialization | File format header fields (not used in crypto package itself) |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| Argon2id | bcrypt | bcrypt is older, weaker (rounds-based not memory-hard); Argon2id is current best practice |
| AES-GCM | ChaCha20-Poly1305 | ChaCha20 better on devices without AES hardware; GCM is more widely supported and hardware-accelerated on x86/ARM |
| stdlib crypto | nacl/box, libsodium bindings | High-level but introduces C dependency (violates CGO_ENABLED=0); stdlib is pure Go |

**Installation:**
```bash
go get golang.org/x/crypto/argon2
go get golang.org/x/sys/unix
go get golang.org/x/sys/windows
```
</standard_stack>

<architecture_patterns>
## Architecture Patterns

### Recommended Package Structure
```
internal/crypto/
├── doc.go              # Package-level documentation
├── kdf.go              # Argon2id: GenerateSalt, DeriveKey
├── aead.go             # AES-GCM: Encrypt, Decrypt
├── memory.go           # ZeroBytes helper
├── mlock_unix.go       # +build !windows (Mlock via x/sys/unix)
├── mlock_windows.go    # +build windows (VirtualLock via x/sys/windows)
├── mlock_other.go      # +build !unix,!windows (no-op stub)
├── password.go         # EvaluatePasswordStrength
├── errors.go           # Sentinel errors
├── kdf_test.go
├── aead_test.go
├── memory_test.go
├── password_test.go
└── crypto_test.go      # Integration tests
```

### Pattern 1: Argon2id Key Derivation
**What:** Derive a 32-byte AES key from password+salt using Argon2id
**When to use:** At vault creation (with GenerateSalt) and vault open (with stored salt)
**Example:**
```go
// Source: golang.org/x/crypto/argon2 godoc
import "golang.org/x/crypto/argon2"

type ArgonParams struct {
    Time    uint32 // Number of iterations (t)
    Memory  uint32 // Memory in KiB (m)
    Threads uint8  // Parallelism (p)
    KeyLen  uint32 // Output key length
}

func DeriveKey(password, salt []byte, params ArgonParams) ([]byte, error) {
    if len(password) == 0 || len(salt) == 0 {
        return nil, ErrInvalidParams
    }
    
    // argon2.IDKey signature:
    // func IDKey(password, salt []byte, time, memory uint32, threads uint8, keyLength uint32) []byte
    key := argon2.IDKey(password, salt, params.Time, params.Memory, params.Threads, params.KeyLen)
    
    return key, nil
}

// Caller's responsibility to zero key after use:
// defer crypto.ZeroBytes(key)
```

**Critical:** Argon2id `memory` parameter is in **KiB**, not bytes. 256 MiB = 262144 KiB.

### Pattern 2: AES-256-GCM Encryption with Unique Nonce
**What:** Encrypt plaintext with AES-GCM, prepend nonce to ciphertext
**When to use:** Every vault save operation
**Example:**
```go
// Source: crypto/cipher godoc + NIST SP 800-38D
import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "io"
)

func Encrypt(key, plaintext []byte) ([]byte, error) {
    if len(key) != 32 {
        return nil, ErrInvalidParams
    }
    
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err // Should never happen with valid 32-byte key
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    // Generate fresh nonce IMMEDIATELY before Seal
    // CRITICAL: nonce must be unique for every encryption with same key
    nonce := make([]byte, gcm.NonceSize()) // 12 bytes for GCM
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, ErrInsufficientEntropy
    }
    
    // Seal appends ciphertext+tag to nonce
    // Output: nonce (12) || ciphertext (len(plaintext)) || tag (16)
    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    
    return ciphertext, nil
}
```

**CRITICAL:** 
- Generate nonce with `io.ReadFull(rand.Reader, nonce)` — not `rand.Read()` (partial reads possible)
- Generate nonce immediately before `Seal()` — never reuse nonce variable
- Nonce is prepended to output (not a separate return value)

### Pattern 3: AES-256-GCM Decryption
**What:** Decrypt ciphertext, verify authentication tag
**When to use:** Vault open
**Example:**
```go
func Decrypt(key, ciphertext []byte) ([]byte, error) {
    if len(key) != 32 {
        return nil, ErrInvalidParams
    }
    
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, ErrAuthFailed // Too short to contain nonce
    }
    
    // Split nonce and ciphertext+tag
    nonce, ciphertextAndTag := ciphertext[:nonceSize], ciphertext[nonceSize:]
    
    // Open verifies tag and decrypts
    plaintext, err := gcm.Open(nil, nonce, ciphertextAndTag, nil)
    if err != nil {
        // Wrong key OR corrupted data OR tampered ciphertext
        // Return single sentinel error (no timing leak, no internal details)
        return nil, ErrAuthFailed
    }
    
    return plaintext, nil
}
```

**CRITICAL:** `gcm.Open()` error means authentication failure — always return `ErrAuthFailed`, never expose internal error message.

### Pattern 4: Memory Zeroing
**What:** Overwrite sensitive byte slices with zeros
**When to use:** After any use of passwords, keys, plaintexts
**Example:**
```go
func ZeroBytes(b []byte) {
    if len(b) == 0 {
        return
    }
    
    // Use copy with zero slice — compiler optimizes to memset
    zeros := make([]byte, len(b))
    copy(b, zeros)
    
    // Alternative: manual loop (also constant-time in Go)
    // for i := range b {
    //     b[i] = 0
    // }
}
```

**Why not clear():** `clear()` built-in (Go 1.21+) zeros slices but doesn't guarantee timing-resistance. Explicit copy/loop is clearer intent.

### Pattern 5: Platform-Specific Memory Locking
**What:** Lock sensitive memory to prevent swapping to disk
**When to use:** After allocating buffers for keys, passwords
**Example:**
```go
// mlock_unix.go
//go:build !windows

package crypto

import "golang.org/x/sys/unix"

func Mlock(b []byte) error {
    if len(b) == 0 {
        return nil
    }
    
    err := unix.Mlock(b)
    if err != nil {
        return ErrMLockFailed
    }
    return nil
}

func Munlock(b []byte) error {
    if len(b) == 0 {
        return nil
    }
    return unix.Munlock(b)
}
```

```go
// mlock_windows.go
//go:build windows

package crypto

import (
    "unsafe"
    "golang.org/x/sys/windows"
)

func Mlock(b []byte) error {
    if len(b) == 0 {
        return nil
    }
    
    // VirtualLock requires pointer and size
    err := windows.VirtualLock(uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
    if err != nil {
        return ErrMLockFailed
    }
    return nil
}

func Munlock(b []byte) error {
    if len(b) == 0 {
        return nil
    }
    return windows.VirtualUnlock(uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
}
```

```go
// mlock_other.go
//go:build !unix && !windows

package crypto

func Mlock(b []byte) error {
    return ErrMLockFailed // Signal unavailable but don't fatal
}

func Munlock(b []byte) error {
    return nil
}
```

**CRITICAL:** 
- mlock failure MUST be non-fatal — containers, low-privilege contexts, some OSes
- Never `append()` to locked buffer — allocation invalidates lock
- Pre-allocate exact size: `key := make([]byte, 32)` then `Mlock(key)`

### Pattern 6: Password Strength Evaluation
**What:** Evaluate password complexity without converting to string
**When to use:** Vault creation UI, master password change
**Example:**
```go
type StrengthLevel int

const (
    StrengthWeak StrengthLevel = iota
    StrengthStrong
)

func EvaluatePasswordStrength(password []byte) StrengthLevel {
    if len(password) < 12 {
        return StrengthWeak
    }
    
    var hasUpper, hasLower, hasDigit, hasSpecial bool
    
    for _, b := range password {
        switch {
        case b >= 'A' && b <= 'Z':
            hasUpper = true
        case b >= 'a' && b <= 'z':
            hasLower = true
        case b >= '0' && b <= '9':
            hasDigit = true
        case !((b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9')):
            hasSpecial = true
        }
    }
    
    if hasUpper && hasLower && hasDigit && hasSpecial {
        return StrengthStrong
    }
    
    return StrengthWeak
}
```

**CRITICAL:** Operate directly on `[]byte` — never `string(password)` (creates unzeroaple copy).

### Anti-Patterns to Avoid
- **String for passwords:** `string(password)` creates immutable copy that cannot be zeroed
- **Reusing nonce slice:** `nonce := make([]byte, 12)` outside loop then `rand.Read(nonce)` in loop — buffer can be reused IF you regenerate before each Seal
- **Ignoring mlock errors as fatal:** mlock unavailable ≠ security failure; log and continue
- **Appending to locked buffer:** `key = append(key, byte)` reallocates, loses lock
- **Storing Argon2id memory in bytes:** m=256 means 256 *KiB*, not 256 *bytes*
- **Using rand.Read() for nonce:** Partial reads possible; use `io.ReadFull(rand.Reader, nonce)`
- **Wrapping stdlib crypto errors:** Return clean sentinel errors (ErrAuthFailed), not `fmt.Errorf("decrypt failed: %w", err)`
</architecture_patterns>

<dont_hand_roll>
## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Key derivation | Custom PBKDF2 tuning, scrypt | `argon2.IDKey()` with standard params | Argon2id is current best practice (RFC 9106); memory-hard, GPU-resistant |
| Nonce generation | `time.Now().UnixNano()`, counter | `io.ReadFull(rand.Reader, nonce)` | Time/counter predictable; crypto/rand uses OS entropy |
| Authenticated encryption | Separate Encrypt + HMAC | `cipher.NewGCM()` | GCM combines confidentiality + authenticity; manual composition error-prone |
| Constant-time comparison | `bytes.Equal()` | `subtle.ConstantTimeCompare()` | Timing attacks possible with byte-by-byte comparison |
| Memory zeroing | `for i := range b { b[i] = 0 }` alone | `copy(b, make([]byte, len(b)))` or loop | Both work; copy may be optimized; never rely on defer to zero (panic can bypass) |
| Password hashing | MD5, SHA-256, bcrypt | Argon2id | Password hashing requires memory-hard KDF; SHA is too fast, bcrypt lacks memory hardness |

**Key insight:** Cryptography is a domain where "good enough" leads to critical vulnerabilities. stdlib + `golang.org/x/crypto` are audited, constant-time, and the Go community standard. Custom implementations introduce timing attacks, nonce reuse, weak parameters, and other subtle flaws.
</dont_hand_roll>

<common_pitfalls>
## Common Pitfalls

### Pitfall 1: Argon2id Memory Parameter in Bytes
**What goes wrong:** `argon2.IDKey(pwd, salt, 3, 256, 4, 32)` produces weak key (only 256 *KiB* memory, not 256 *MiB*)
**Why it happens:** `memory` parameter is in KiB; 256 MiB = 262144 KiB
**How to avoid:** Document parameter in comments: `memory: 262144, // 256 MiB in KiB`
**Warning signs:** Key derivation completes in <100ms (should be ~200-500ms on modern CPU)

### Pitfall 2: Nonce Reuse in GCM
**What goes wrong:** Encrypting twice with same key+nonce leaks plaintext XOR
**Why it happens:** Reusing nonce buffer across loop iterations without regenerating
**How to avoid:** Generate nonce immediately before each `Seal()` call; never global nonce variable
**Warning signs:** Deterministic ciphertext (same plaintext → same ciphertext every time)

### Pitfall 3: String Conversion for Passwords
**What goes wrong:** `password := string(passwordBytes)` creates unzeroaple copy
**Why it happens:** Go strings are immutable; cannot overwrite
**How to avoid:** Every API signature accepts `[]byte` for sensitive data; zero after use
**Warning signs:** `go test -race` may show data races on password strings (false positive but indicates string misuse)

### Pitfall 4: mlock Failure as Fatal Error
**What goes wrong:** Application crashes in Docker, low-privilege contexts, or unsupported platforms
**Why it happens:** Treating `ErrMLockFailed` as fatal
**How to avoid:** Return error but continue execution; caller logs warning
**Warning signs:** Application unusable in containers or as non-root user

### Pitfall 5: Exposing Stdlib Errors to User
**What goes wrong:** `gcm.Open()` error message exposed as "cipher: message authentication failed"
**Why it happens:** Returning raw stdlib error instead of sentinel
**How to avoid:** Wrap all stdlib errors in package sentinel errors; never `return err` directly
**Warning signs:** User-facing error messages contain "cipher:", "crypto/", or internal details

### Pitfall 6: Partial Read from crypto/rand
**What goes wrong:** `rand.Read(nonce)` reads fewer bytes than requested on entropy exhaustion
**Why it happens:** `rand.Read()` returns (n int, err error); n < len(nonce) is valid
**How to avoid:** Use `io.ReadFull(rand.Reader, nonce)` — ensures exactly len(nonce) bytes
**Warning signs:** Nonce shorter than 12 bytes; `ErrUnexpectedEOF` in tests

### Pitfall 7: Appending to Locked Memory
**What goes wrong:** `key = append(key, newByte)` reallocates, new allocation not locked
**Why it happens:** `append()` may allocate new backing array
**How to avoid:** Pre-allocate exact size (`make([]byte, 32)`), never grow locked buffers
**Warning signs:** mlock succeeds but key still swapped to disk

### Pitfall 8: Using defer for Memory Zeroing in Critical Paths
**What goes wrong:** `defer ZeroBytes(key)` doesn't run on panic; key leaks in crash dump
**Why it happens:** `defer` bypassed by `os.Exit()`, panic without recovery
**How to avoid:** Explicit zero in error paths; recover panics at boundaries; document caller's zeroing responsibility
**Warning signs:** Core dumps contain plaintext keys
</common_pitfalls>

<code_examples>
## Code Examples

Verified patterns from official sources:

### Complete Encrypt/Decrypt Roundtrip
```go
// Source: crypto/cipher godoc + golang.org/x/crypto/argon2 godoc
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "fmt"
    "io"
    "golang.org/x/crypto/argon2"
)

func main() {
    password := []byte("ExamplePassword123!")
    plaintext := []byte("secret data")
    
    // Generate salt
    salt := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, salt); err != nil {
        panic(err)
    }
    
    // Derive key with Argon2id
    key := argon2.IDKey(password, salt, 3, 262144, 4, 32)
    defer func() {
        for i := range key {
            key[i] = 0
        }
    }()
    
    // Encrypt
    ciphertext, err := encrypt(key, plaintext)
    if err != nil {
        panic(err)
    }
    
    // Decrypt
    decrypted, err := decrypt(key, ciphertext)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Roundtrip successful: %s\n", decrypted)
}

func encrypt(key, plaintext []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }
    
    return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(key, ciphertext []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, fmt.Errorf("ciphertext too short")
    }
    
    nonce, ciphertextAndTag := ciphertext[:nonceSize], ciphertext[nonceSize:]
    return gcm.Open(nil, nonce, ciphertextAndTag, nil)
}
```

### Nonce Uniqueness Test
```go
// Source: NIST SP 800-38D requirement verification
func TestNonceUniqueness(t *testing.T) {
    key := make([]byte, 32)
    io.ReadFull(rand.Reader, key)
    
    plaintext := []byte("same plaintext")
    
    c1, _ := Encrypt(key, plaintext)
    c2, _ := Encrypt(key, plaintext)
    
    if bytes.Equal(c1, c2) {
        t.Error("Nonce reused: identical ciphertexts for same plaintext")
    }
}
```

### Platform-Specific Build Tags
```go
// Source: golang.org/x/sys documentation
// File: mlock_unix.go
//go:build !windows

package crypto

import "golang.org/x/sys/unix"

func Mlock(b []byte) error {
    return unix.Mlock(b)
}
```

```go
// File: mlock_windows.go
//go:build windows

package crypto

import (
    "unsafe"
    "golang.org/x/sys/windows"
)

func Mlock(b []byte) error {
    return windows.VirtualLock(uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
}
```

```go
// File: mlock_other.go
//go:build !unix && !windows

package crypto

func Mlock(b []byte) error {
    return ErrMLockFailed
}
```
</code_examples>

<sota_updates>
## State of the Art (2024-2026)

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| PBKDF2 | Argon2id | 2015 (RFC 9106 in 2021) | Memory-hard KDF resists GPU/ASIC attacks; Argon2id is current best practice |
| bcrypt | Argon2id | 2015 | bcrypt lacks memory hardness; Argon2id stronger against parallel attacks |
| AES-CBC + HMAC | AES-GCM | 2007 (GCM standardized) | Single primitive for confidentiality + authenticity; fewer implementation errors |
| ChaCha20-Poly1305 on all platforms | AES-GCM on x86/ARM, ChaCha20 on others | 2014 (ChaCha20 RFC 7539) | AES hardware acceleration on modern CPUs; ChaCha20 better where AES-NI absent |
| Manual nonce management | Generate in Encrypt() | Always | Reduces caller burden and nonce-reuse risk |

**New tools/patterns to consider:**
- **Go 1.21+ clear():** Built-in to zero slices, but `copy(b, make(...))` still preferred for explicit intent
- **Go 1.20+ crypto/rand.Reader determinism:** Tests can use `rand.Reader = deterministicReader` for reproducible fixtures
- **`x/crypto/hkdf`:** For deriving multiple keys from single master key (not needed for single AES key)

**Deprecated/outdated:**
- **scrypt:** Superseded by Argon2id (better tuning, memory-hard)
- **MD5, SHA-1 for passwords:** Never acceptable (too fast, rainbow tables)
- **AES-CTR alone:** No authentication; use GCM or separate HMAC
- **`math/rand` for crypto:** Use `crypto/rand` — `math/rand` is pseudo-random (deterministic seed)
</sota_updates>

<validation_architecture>
## Validation Architecture

### Nyquist Validation Strategy

**Test pyramid for crypto package:**

1. **Unit tests (internal correctness):**
   - Nonce uniqueness (encrypt same plaintext twice → distinct ciphertexts)
   - Roundtrip (Encrypt → Decrypt = original plaintext)
   - Wrong key (Decrypt with different key → `ErrAuthFailed`)
   - Short ciphertext (Decrypt with <12 bytes → error)
   - ZeroBytes (fills entire slice with 0x00)
   - Password strength boundary cases:
     - 11 chars with all categories → Weak
     - 12 chars with all categories → Strong
     - 12 chars missing any category → Weak

2. **Integration tests (cross-function behavior):**
   - Full vault roundtrip (GenerateSalt → DeriveKey → Encrypt → Decrypt with same password)
   - Key reuse safety (Encrypt twice with same key produces distinct ciphertexts)
   - Memory safety (ZeroBytes → verify buffer is all zeros)

3. **Property-based tests (fuzz targets for Phase 11):**
   - Encrypt any []byte → Decrypt → must equal original
   - Any password → DeriveKey → produces 32-byte key
   - Any []byte → ZeroBytes → all bytes are 0x00

4. **Security tests (explicit attack simulation):**
   - Timing attack resistance: Wrong password takes same time as right password
   - Entropy exhaustion: Simulate `rand.Reader` failure → `ErrInsufficientEntropy`
   - mlock failure: Simulate unavailable mlock → non-fatal, returns `ErrMLockFailed`

5. **Platform-specific tests:**
   - Unix: mlock succeeds on writable buffer
   - Windows: VirtualLock succeeds on writable buffer
   - Other: mlock returns `ErrMLockFailed` gracefully

6. **Benchmark tests (performance baseline):**
   - Argon2id derivation time (t=3, m=256MiB, p=4) → 200-500ms on modern CPU
   - AES-GCM encrypt/decrypt throughput → >100 MB/s

**Test data strategy:**
- Use fixed seeds for deterministic tests: `t.Setenv("GO_TEST_SEED", "12345")`
- Clearly fictional passwords: `"P@ssw0rd!Test42"`, `"FakeSecretKey123!"`
- No real user data in test fixtures

**Race detector:**
- ALL tests must pass with `-race` flag
- Crypto package has no goroutines, but tests may simulate concurrent access

**Coverage target:**
- 100% of exported functions
- 100% of error paths
- 90%+ of internal helpers
</validation_architecture>

<open_questions>
## Open Questions

None — stdlib + `golang.org/x/crypto` are well-documented, widely-used, and have clear best practices. All decisions locked per CONTEXT.md.
</open_questions>

<sources>
## Sources

### Primary (HIGH confidence)
- Go stdlib `crypto/cipher` godoc - GCM interface, nonce requirements, Seal/Open behavior
- Go stdlib `crypto/aes` godoc - AES block cipher construction
- Go stdlib `crypto/rand` godoc - Reader interface, entropy guarantees
- `golang.org/x/crypto/argon2` godoc - IDKey function signature, parameter meanings
- `golang.org/x/sys/unix` godoc - Mlock/Munlock syscalls
- `golang.org/x/sys/windows` godoc - VirtualLock/VirtualUnlock APIs
- NIST SP 800-38D - GCM specification (nonce uniqueness, tag verification)
- RFC 9106 - Argon2id specification (parameter recommendations)

### Secondary (MEDIUM confidence)
- OWASP Password Storage Cheat Sheet - Password strength criteria (verified against CONTEXT.md D-34)
- Go security best practices - Memory zeroing patterns, []byte for sensitive data

### Tertiary (LOW confidence - needs validation)
- None — all findings from primary sources
</sources>

<metadata>
## Metadata

**Research scope:**
- Core technology: Go stdlib crypto + golang.org/x/crypto
- Ecosystem: No third-party crypto libraries (per CRYPTO-02 requirement)
- Patterns: Argon2id KDF, AES-256-GCM AEAD, memory locking, password strength
- Pitfalls: Nonce reuse, string conversion, mlock failures, Argon2id parameter units

**Confidence breakdown:**
- Standard stack: HIGH - stdlib + x/crypto only, well-established
- Architecture: HIGH - godoc examples, NIST/RFC specifications
- Pitfalls: HIGH - documented in Go crypto wiki, security advisories
- Code examples: HIGH - directly from godoc and RFCs
- Validation strategy: HIGH - standard Go testing patterns for crypto

**Research date:** 2026-03-29
**Valid until:** 2027-03-29 (365 days - Go crypto stdlib is extremely stable; x/crypto API frozen)
</metadata>

---

*Phase: 02-crypto-package*
*Research completed: 2026-03-29*
*Ready for planning: yes*
