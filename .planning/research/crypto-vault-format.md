# Cryptography & Vault File Format Research

**Project:** Abditum — Portable Go Password Vault  
**Researched:** 2026-03-24  
**Go version tested against:** go1.26.1 windows/amd64  
**Overall confidence:** HIGH (all claims verified against Go stdlib docs, official sources, or inspected source code)

---

## 1. AES-256-GCM in Go's Standard Library

**Confidence:** HIGH — verified directly against `go doc crypto/aes` and `go doc crypto/cipher`

### Correct Usage Pattern

AES-256-GCM in Go requires chaining three stdlib calls:

```go
import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
)

// Encryption
func encrypt(key, plaintext, additionalData []byte) (ciphertext []byte, err error) {
    block, err := aes.NewCipher(key)       // key must be exactly 32 bytes for AES-256
    if err != nil {
        return nil, err
    }

    // NewGCMWithRandomNonce is available in Go 1.22+ and is the cleanest approach:
    // it generates a random 96-bit nonce internally, prepends it to the ciphertext,
    // and uses Overhead() = 28 bytes (12-byte nonce + 16-byte GCM tag).
    gcm, err := cipher.NewGCMWithRandomNonce(block)
    if err != nil {
        return nil, err
    }

    // Seal appends the nonce (prepended) + encrypted ciphertext + GCM tag.
    // Pass nil as dst to let Seal allocate.
    // additionalData (AAD) is authenticated but NOT encrypted — use it for the
    // unencrypted header so tampering with the header is detected.
    return gcm.Seal(nil, nil, plaintext, additionalData), nil
}

// Decryption
func decrypt(key, ciphertext, additionalData []byte) (plaintext []byte, err error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCMWithRandomNonce(block)
    if err != nil {
        return nil, err
    }

    // Open strips the nonce, decrypts, and verifies the GCM tag.
    // Returns error if tag verification fails (tampered data or wrong key).
    return gcm.Open(nil, nil, ciphertext, additionalData)
}
```

### Alternative: Manual Nonce (Explicit Header Storage)

If storing the nonce separately in the file header (which is the design in `descricao.md`), use
the lower-level `NewGCM` with a hand-generated nonce instead:

```go
func encryptWithSeparateNonce(key, plaintext, additionalData []byte) (nonce, ciphertext []byte, err error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, nil, err
    }
    gcm, err := cipher.NewGCM(block)  // standard 96-bit nonce, 128-bit tag
    if err != nil {
        return nil, nil, err
    }

    nonce = make([]byte, gcm.NonceSize()) // 12 bytes
    if _, err = rand.Read(nonce); err != nil {   // crypto/rand — never fails on modern OS
        return nil, nil, err
    }
    // Seal writes: encrypted payload + 16-byte GCM authentication tag
    ciphertext = gcm.Seal(nil, nonce, plaintext, additionalData)
    return nonce, ciphertext, nil
}
```

### Key Findings

- **Key length:** Must be exactly 32 bytes for AES-256. Pass the 32-byte output of `argon2.IDKey(...)` directly.
- **Nonce size:** Standard GCM uses 96-bit (12-byte) nonces. This is the correct size — do not change it.
- **GCM tag:** 16 bytes appended to the ciphertext by `Seal`. `Open` verifies this tag before returning plaintext.
- **`NewGCMWithRandomNonce`:** Added in Go 1.22. Generates nonce internally and prepends it to ciphertext (`Overhead() = 28`). Cleaner than manual nonce, but mixes nonce into the ciphertext blob rather than keeping it in the header. Since Abditum stores the nonce in the header separately, use `NewGCM` instead.
- **Hardware acceleration:** On amd64 with AES-NI (all modern x86 CPUs), `aes.NewCipher` + `cipher.NewGCM` run in constant time and use hardware-accelerated AES/GCM — no third-party library needed.
- **`crypto/rand.Read`:** In Go 1.20+, this function never returns an error and always fills the buffer. Safe to ignore the error, but do check it anyway for clarity.
- **Additional Data (AAD):** Pass the raw binary header bytes as AAD. This cryptographically binds the header to the ciphertext — any tampering with magic bytes, version, salt, or nonce bytes will cause decryption to fail with an authentication error. This provides the "integrity check" the spec requires without any extra checksum field.

### What NOT to Do

- ❌ **Never reuse a nonce with the same key.** Nonce reuse completely breaks GCM confidentiality and authenticity. Abditum regenerates the nonce on every save — correct.
- ❌ **Do not use `crypto/cipher.NewCFB*` or `NewCBC*`** — these are unauthenticated modes. GCM is the right choice.
- ❌ **Do not use a counter-based nonce for a vault** — counter requires persistent state. Random nonce is simpler and safe given < 2^32 saves per key.

---

## 2. Argon2id Parameters for a Password Vault

**Confidence:** HIGH — based on RFC 9106 (Argon2 standard), OWASP Password Storage Cheat Sheet, and Go `x/crypto/argon2` package analysis  
**Source:** RFC 9106 §4 "Parameter Choice", OWASP recommendations

### Recommended Parameters

```go
import "golang.org/x/crypto/argon2"

const (
    // Memory cost: 64 MiB — good balance on low-RAM devices (e.g., 512MB USB-boot machines)
    // vs. 256 MiB for high-security workstations. Start here; expose as a config option.
    Argon2Memory      = 64 * 1024  // 64 MiB in KiB units

    // Time cost: 3 iterations. RFC 9106 recommends time=1 as minimum; 3 gives ~3x slowdown
    // over t=1 with same memory, meaningfully raises brute-force cost.
    Argon2Time        = 3

    // Parallelism: 4 threads. More threads = faster on multi-core attacker hardware too.
    // 4 is the sweet spot — uses 4 cores, but attacker gains proportionally less than
    // cost from memory × time.
    Argon2Threads     = 4

    // Output key length: 32 bytes — exactly the AES-256 key size needed.
    Argon2KeyLen      = 32

    // Salt length: 32 bytes — RFC 9106 recommends ≥ 16 bytes; 32 is the safe choice.
    Argon2SaltLen     = 32
)

func deriveKey(masterPassword, salt []byte) []byte {
    // argon2.IDKey uses the Argon2id variant, which is the hybrid
    // (protects against both side-channel and GPU attacks).
    // The returned slice is exactly Argon2KeyLen bytes.
    return argon2.IDKey(
        masterPassword,  // password bytes
        salt,            // random salt stored in header
        Argon2Time,      // time cost
        Argon2Memory,    // memory cost in KiB
        Argon2Threads,   // parallelism
        Argon2KeyLen,    // output key bytes
    )
}
```

### Parameter Rationale

| Parameter | Value | Rationale |
|-----------|-------|-----------|
| Algorithm | Argon2**id** | Hybrid: resistant to both side-channel attacks (d-variant) and GPU/ASIC attacks (i-variant). RFC 9106 §4 recommends Argon2id for password hashing when there is no side-channel concern, which covers this app. |
| Memory | 64 MiB | RFC 9106 "first recommended" parameters use 64 MiB. Ballpark: ~0.5s on a modern laptop. Too little (< 16 MiB) leaves GPU attacks feasible. |
| Time | 3 | Multiplies wall-clock time linearly without reducing the memory requirement. 3 iterations ~= 1–2 seconds on average hardware with 64 MiB. |
| Parallelism | 4 | OWASP recommends ≥ 1 and generally ≤ available cores. 4 is suitable for cross-platform; higher parallelism does not proportionally help the defender on single-core scenarios. |
| Salt | 32 bytes | Random per-vault. Generated once at vault creation with `crypto/rand`, stored in plaintext header. Never reused. |
| Key output | 32 bytes | Exactly the AES-256 key. No truncation/padding needed. |

### Expected Unlock Time (approximate, medium-tier hardware 2024)

- ~0.5–1.5 seconds on modern laptop/desktop
- ~2–4 seconds on a Raspberry Pi or very old machine

**Decision:** Expose `memory` and `time` as stored parameters in the vault header so users on more capable machines can increase them, and future versions can update defaults while still reading old vaults.

### Critical: Store Parameters in the Header

The Argon2id parameters (memory, time, threads, salt) MUST be stored in the vault file header so that:
1. The vault can be opened on any machine regardless of defaults.
2. Future Abditum versions can change defaults without breaking existing vaults.
3. The format is self-describing.

---

## 3. Vault File Format Design

**Confidence:** HIGH — design choices based on Go stdlib capabilities, security principles, and the project spec

### Recommended Binary Layout

```
+--------------------------------------------------+
|           UNENCRYPTED HEADER (fixed size)        |
+--------------------------------------------------+
| [0..7]   Magic bytes: 0x41 0x42 0x44 0x49 0x54  |
|          0x55 0x4D 0x0A ("ABDITUM\n")            |
| [8..9]   Format version: uint16 big-endian       |
| [10..11] Argon2id time cost: uint16 big-endian   |
| [12..15] Argon2id memory cost: uint32 big-endian |
| [16]     Argon2id parallelism: uint8             |
| [17]     Argon2id key length: uint8 (=32)        |
| [18..49] Salt: 32 bytes (random, per-vault)      |
| [50..61] Nonce: 12 bytes (random, per-save)      |
+--------------------------------------------------+  ← Header ends at byte 62
|           ENCRYPTED PAYLOAD                      |
+--------------------------------------------------+
| AES-256-GCM ciphertext (variable length)         |
| Last 16 bytes = GCM authentication tag           |
+--------------------------------------------------+
```

**Total unencrypted header: 62 bytes (fixed)**

### Header Field Details

| Bytes | Field | Type | Notes |
|-------|-------|------|-------|
| 0–7 | Magic | `[8]byte` | `"ABDITUM\n"` — detect .abditum files, fail fast on wrong file |
| 8–9 | Format version | `uint16` big-endian | Start at `1`. Increment when breaking changes to payload JSON schema require migration. |
| 10–11 | Argon2 time | `uint16` big-endian | Iterations. Stored so old vaults still open after defaults change. |
| 12–15 | Argon2 memory | `uint32` big-endian | KiB. Same reason. Max ~4 TiB, plenty. |
| 16 | Argon2 threads | `uint8` | Parallelism factor. 1–255 range is sufficient. |
| 17 | Argon2 key len | `uint8` | Always 32 for AES-256. Stored for self-documentation. |
| 18–49 | Salt | `[32]byte` | `crypto/rand` generated at vault creation, never changed. |
| 50–61 | Nonce | `[12]byte` | `crypto/rand` generated on EVERY save. Never reuse. |

### What Goes in the Encrypted Payload

Everything else. The payload is a UTF-8 JSON document encoding the full `Cofre` (vault) struct:

```json
{
  "version": 1,
  "created_at": "2026-03-24T10:00:00Z",
  "updated_at": "2026-03-24T10:00:00Z",
  "config": { ... },
  "folders": [ ... ],
  "secrets": [ ... ],
  "templates": [ ... ]
}
```

**Why version is also inside the payload:** The payload version number refers to the JSON schema version (what fields exist in secrets, templates, etc.), while the header format version refers to the binary file format. These can evolve independently.

### Using AAD for Header Integrity

Pass the raw header bytes (all 62 bytes) as GCM Additional Authenticated Data:

```go
header := buildHeader(...) // 62 bytes
ciphertext := gcm.Seal(nil, nonce, plaintextJSON, header)
// On decrypt:
plaintext, err := gcm.Open(nil, nonce, ciphertext, header)
// If header was tampered, err != nil — the vault is corrupted or attacked.
```

This means you get integrity verification of the header for free from GCM — no extra HMAC or checksum field needed.

### Version Detection for N-1 Backward Compatibility

```go
func openVault(data []byte) (*Vault, error) {
    if len(data) < headerSize {
        return nil, ErrCorrupted
    }
    if !bytes.Equal(data[0:8], magicBytes) {
        return nil, ErrNotAVaultFile
    }
    version := binary.BigEndian.Uint16(data[8:10])
    switch version {
    case 1:
        return openVaultV1(data)
    case 2:
        return openVaultV2(data) // future
    default:
        return nil, fmt.Errorf("unsupported vault format version %d", version)
    }
}
```

**N-1 rule implementation:** Keep `openVaultV{N-1}` when shipping version N. The current build of Abditum should always be able to open format version `currentVersion - 1`. Keep a migration function `migrateV1toV2(payload)` that upgrades on open-then-save.

---

## 4. Atomic File Save Pattern (Cross-Platform)

**Confidence:** HIGH for POSIX; MEDIUM for Windows (atomic rename is not guaranteed without special APIs)

### The Problem with Windows `os.Rename`

Go's `os.Rename` documentation explicitly states:

> "Even within the same directory, on non-Unix platforms Rename is not an atomic operation."

The underlying Windows `MoveFileW` API will fail with `ERROR_ALREADY_EXISTS` if the destination already exists. The replacement requires `MoveFileExW` with the `MOVEFILE_REPLACE_EXISTING` flag — which IS documented to be atomic when source and destination are on the same volume.

**Go's `os.Rename` on Windows** (since Go 1.5) actually calls `MoveFileExW` with `MOVEFILE_REPLACE_EXISTING` internally, which means it CAN replace an existing file. However, it is still **not guaranteed to be atomic** from a crash-consistency standpoint (a power-loss between write and rename can leave the destination replaced but the source unsynced).

The `google/renameio` library explicitly documents:

> "It is not possible to reliably write files atomically on Windows" — and exports **no functions** on Windows.

### Recommended Pattern for Abditum

Given the project's requirements (atomicity + backup), implement the following:

```go
// saveVaultAtomically writes vault data safely.
// On success: .abditum.tmp → .abditum.bak → .abditum (new)
// On failure: .abditum.tmp is deleted; original .abditum is untouched.
func saveVaultAtomically(vaultPath string, data []byte) error {
    dir := filepath.Dir(vaultPath)
    tmpPath := vaultPath + ".tmp"
    bakPath := vaultPath + ".bak"

    // Step 1: Write encrypted vault to .abditum.tmp
    f, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        return fmt.Errorf("cannot create tmp file: %w", err)
    }

    _, writeErr := f.Write(data)
    syncErr := f.Sync()  // flush to OS page cache → persistent storage
    closeErr := f.Close()

    // Cleanup tmp on any write failure
    if writeErr != nil || syncErr != nil || closeErr != nil {
        os.Remove(tmpPath)
        if writeErr != nil { return writeErr }
        if syncErr != nil  { return syncErr }
        return closeErr
    }

    // Step 2: If vault exists, copy it to .bak before replacing
    if _, err := os.Stat(vaultPath); err == nil {
        if err := copyFile(vaultPath, bakPath); err != nil {
            os.Remove(tmpPath)
            return fmt.Errorf("cannot create backup: %w", err)
        }
    }

    // Step 3: Rename .abditum.tmp → .abditum
    // On POSIX: atomic. On Windows: not truly atomic but best available option.
    if err := os.Rename(tmpPath, vaultPath); err != nil {
        os.Remove(tmpPath)
        return fmt.Errorf("cannot finalize vault: %w", err)
    }

    // Step 4: fsync the directory (POSIX only; no-op on Windows, which is fine)
    if df, err := os.Open(dir); err == nil {
        df.Sync()
        df.Close()
    }

    return nil
}
```

### Windows-Specific Considerations

- **`os.Rename` on Windows:** Works for replacing existing files since Go 1.5 (uses `MoveFileExW` + `MOVEFILE_REPLACE_EXISTING`). The source and destination must be on the same filesystem/drive, which is guaranteed since `.tmp` is in the same directory as `.abditum`.
- **Power loss during rename on Windows:** True atomic rename is not achievable without Windows `TRANSACTIONAL NTFS` (TxF), which is deprecated. The practical risk for a password vault (used interactively, not by automated systems) is extremely low.
- **Antivirus interference:** Windows AV software may hold a file lock during scan, causing `os.Rename` to fail with `ACCESS_DENIED`. Handle this by retrying a few times with a short delay before surfacing an error to the user.
- **`google/renameio`:** Do NOT use this library. It explicitly excludes Windows (no exported functions on `GOOS=windows`).

### The `.bak` Backup

The spec requires keeping a `.abditum.bak`. The pattern above handles this by copying the current vault to `.bak` before replacing it. This gives users a one-level-deep undo. The `.bak` file is always the PREVIOUS successful save, not the current save.

---

## 5. Secure Memory Zeroing in Go

**Confidence:** HIGH for the problem diagnosis; MEDIUM for the memguard recommendation (API is still marked "experimental")

### The Core Problem

Go's GC makes secure memory zeroing non-trivial:

1. **The GC can copy heap objects.** When the runtime moves a `[]byte` containing the master password during garbage collection, the old memory location may retain a copy of the bytes. Simply zeroing the slice does not zero the abandoned memory.
2. **The compiler can elide "unnecessary" writes.** A zero-loop followed by no reads may be optimized away by the compiler as a dead write, just as `memset` can be elided in C without `memset_s`.
3. **`string` values are immutable.** A Go `string` containing the master password cannot be zeroed at all — strings share underlying memory and are GC-managed. Always work with `[]byte` for sensitive values.

### The Manual Zeroing Pattern (Minimum Viable)

```go
// zeroBytes overwrites a byte slice with zeros.
// Uses a range-based loop which is harder for the optimizer to elide
// than a single assignment. Not perfect, but a reasonable baseline.
func zeroBytes(b []byte) {
    for i := range b {
        b[i] = 0
    }
    // runtime.KeepAlive prevents the compiler from treating the write as dead.
    // This is the idiomatic Go approach.
    runtime.KeepAlive(b)
}
```

**Limitation:** This only zeroes the slice's backing array at its current address. If the GC moved the slice's backing array before this call, the old copy is not zeroed.

### The `memguard` Library (Stronger Guarantee)

**Repository:** `github.com/awnumar/memguard` (2.7k stars, active, last release Aug 2025)  
**Pure Go:** Yes — no CGO required, works on Windows/macOS/Linux.

`memguard` provides a software enclave that:
- Allocates memory via syscalls (bypasses Go GC), so the GC cannot move or copy it.
- Locks the memory page (`mlock`/`VirtualLock`) to prevent it from being swapped to disk.
- Encrypts the buffer contents when not in use (XSalsa20-Poly1305 internally).
- Adds guard pages and canary values around the buffer to detect overflow.
- Zeroes and frees memory via `Destroy()`.

```go
import "github.com/awnumar/memguard"

// At startup:
memguard.CatchInterrupt()  // ensure cleanup on Ctrl+C / SIGTERM

// Store master password:
lockedBuf := memguard.NewBufferFromBytes(passwordBytes)
// passwordBytes slice should be zeroed immediately after this call

// Access safely:
key := lockedBuf.Bytes()  // temporary access; only available while buffer is open

// On lock or app exit:
lockedBuf.Destroy()       // zeroes, unlocks, and frees the memory

// Nuclear option — destroy everything:
memguard.Purge()
```

**Use `memguard` for:**
- The master password (entered by the user at the keyboard — must not outlive the unlock session)
- The derived AES-256 key (32-byte key material — must be zeroed on vault lock)
- The decrypted plaintext JSON (potentially large — most sensitive; must be zeroed on lock)

**Use manual zeroing (`for i := range b { b[i] = 0 }`) for:**
- Temporary intermediate buffers (e.g., the byte slice returned from `gcm.Open`)
- Short-lived sensitive values that don't need GC resistance

### Recommendation for Abditum

Use **`memguard`** for the three key lifecycle values:
1. `masterPassword []byte` — from entry to key derivation
2. `derivedKey [32]byte` — from Argon2id output to end of unlock session
3. `decryptedPayload []byte` — from `gcm.Open` return to end of unlock session

Use **manual zeroing** for all other temporary byte slices. Don't store the master password as a Go `string`.

### API Stability Warning

memguard's README states: *"API is experimental and may have unstable changes. You should pin a version."* Pin it in `go.mod` (`go get github.com/awnumar/memguard@v0.23.0`).

---

## 6. Cross-Platform Clipboard

**Confidence:** HIGH for atotto/clipboard; HIGH for golang-design/clipboard CGO requirement  
**Source:** Inspected source code of both libraries directly

### Candidates Compared

| Library | Stars | Windows | macOS | Linux | CGO | Clear clipboard |
|---------|-------|---------|-------|-------|-----|-----------------|
| `github.com/atotto/clipboard` | 1.4k | ✅ syscall | ✅ `pbcopy`/`pbpaste` | ✅ xsel/xclip/wl-clipboard | ❌ No CGO | ❌ No built-in clear |
| `golang.design/x/clipboard` | 770 | ✅ syscall | ⚠️ Requires CGO | ⚠️ Requires CGO + X11 | ✅/❌ Mixed | ✅ Watch API |

### Critical Difference: CGO Requirements

**`golang-design/clipboard`:**
- macOS: Uses Objective-C (`clipboard_darwin.m`) — **requires CGO**
- Linux: Uses C (`clipboard_linux.c`) — **requires CGO + `libx11-dev`**
- Windows: Pure Go (no CGO)
- Has `clipboard_nocgo.go` which simply panics: `"clipboard: cannot use when CGO_ENABLED=0"`

**This library cannot be used in a CGO-disabled static build for macOS/Linux.** Since Abditum targets a single portable binary with no external dependencies, `golang-design/clipboard` is unsuitable.

**`atotto/clipboard`:**
- Windows: Pure Go — calls Win32 clipboard API via `syscall` (no CGO)
- macOS: Calls `pbcopy`/`pbpaste` shell commands — no CGO
- Linux: Calls `xclip`, `xsel`, or `wl-clipboard` (Wayland) shell commands — no CGO, but requires the external tool to be installed
- Text-only (UTF-8 strings) — sufficient for Abditum (field values are text)
- No native "clear clipboard" function — implement clear by writing an empty string or a benign placeholder

### Linux Caveat

`atotto/clipboard` on Linux shells out to `xclip`/`xsel`/`wl-copy`. These must be installed on the host system. For a TUI app targeting Linux desktop users this is acceptable — they will have X11/Wayland tooling. For headless server environments (where the clipboard is irrelevant anyway), this will fail gracefully with an error.

`atotto/clipboard` auto-detects at init time what's available (checks Wayland → xclip → xsel → termux). If nothing is found, it sets `Unsupported = true` and returns an error on clipboard operations. Abditum should handle this error gracefully (show a message: "Clipboard not available — copy your field value manually").

### Clear Clipboard Implementation

`atotto/clipboard` does not have a native clear function. Implement it as:

```go
import "github.com/atotto/clipboard"

func clearClipboard() error {
    // Writing an empty string effectively clears user-visible content.
    // Alternatively write a single space to avoid some apps treating empty as "no data".
    return clipboard.WriteAll("")
}
```

### Recommendation

**Use `github.com/atotto/clipboard`.** It is pure Go on all three target platforms and does not require CGO. This preserves the ability to build with `CGO_ENABLED=0` for maximum portability.

**Do NOT use `golang.design/x/clipboard`** — it requires CGO on macOS and Linux, breaking the single static binary goal.

---

## 7. Putting It All Together: Implementation Skeleton

This skeleton illustrates how all six areas integrate:

```go
// internal/crypto/vault.go

package crypto

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/binary"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "runtime"

    "github.com/awnumar/memguard"
    "golang.org/x/crypto/argon2"
)

// File format constants
const (
    headerMagic   = "ABDITUM\n"  // 8 bytes
    headerSize    = 62            // bytes
    currentVersion = uint16(1)

    // Argon2id default parameters — stored in header, so these are just defaults for new vaults.
    defaultArgon2Time    = uint16(3)
    defaultArgon2Memory  = uint32(64 * 1024)  // 64 MiB
    defaultArgon2Threads = uint8(4)
    argon2KeyLen         = uint8(32)

    saltLen  = 32  // bytes
    nonceLen = 12  // bytes (96-bit GCM nonce)
)

// Header holds the unencrypted vault file header.
type Header struct {
    Magic         [8]byte
    FormatVersion uint16
    Argon2Time    uint16
    Argon2Memory  uint32
    Argon2Threads uint8
    Argon2KeyLen  uint8
    Salt          [saltLen]byte
    Nonce         [nonceLen]byte
}

// encodeHeader serializes the header to exactly headerSize bytes (big-endian).
func encodeHeader(h Header) [headerSize]byte {
    var buf [headerSize]byte
    copy(buf[0:8], h.Magic[:])
    binary.BigEndian.PutUint16(buf[8:10], h.FormatVersion)
    binary.BigEndian.PutUint16(buf[10:12], h.Argon2Time)
    binary.BigEndian.PutUint32(buf[12:16], h.Argon2Memory)
    buf[16] = h.Argon2Threads
    buf[17] = h.Argon2KeyLen
    copy(buf[18:50], h.Salt[:])
    copy(buf[50:62], h.Nonce[:])
    return buf
}

// decodeHeader parses and validates the first headerSize bytes of a vault file.
func decodeHeader(data []byte) (Header, error) {
    if len(data) < headerSize {
        return Header{}, fmt.Errorf("file too short to be a vault")
    }
    var h Header
    copy(h.Magic[:], data[0:8])
    if string(h.Magic[:]) != headerMagic {
        return Header{}, fmt.Errorf("not a valid Abditum vault file")
    }
    h.FormatVersion = binary.BigEndian.Uint16(data[8:10])
    h.Argon2Time    = binary.BigEndian.Uint16(data[10:12])
    h.Argon2Memory  = binary.BigEndian.Uint32(data[12:16])
    h.Argon2Threads = data[16]
    h.Argon2KeyLen  = data[17]
    copy(h.Salt[:], data[18:50])
    copy(h.Nonce[:], data[50:62])
    return h, nil
}

// EncryptVault encrypts the vault payload (JSON bytes) using the master password.
// Returns the complete file bytes (header + ciphertext).
func EncryptVault(masterPassword []byte, payload []byte) ([]byte, error) {
    // 1. Generate fresh salt and nonce for this save
    var salt [saltLen]byte
    var nonce [nonceLen]byte
    rand.Read(salt[:])
    rand.Read(nonce[:])

    // 2. Derive AES-256 key from master password using Argon2id
    key := argon2.IDKey(
        masterPassword,
        salt[:],
        uint32(defaultArgon2Time),
        defaultArgon2Memory,
        defaultArgon2Threads,
        uint32(argon2KeyLen),
    )
    defer zeroBytes(key)  // erase key from memory when done

    // 3. Build header struct and serialize it
    header := Header{
        FormatVersion: currentVersion,
        Argon2Time:    defaultArgon2Time,
        Argon2Memory:  defaultArgon2Memory,
        Argon2Threads: defaultArgon2Threads,
        Argon2KeyLen:  argon2KeyLen,
        Salt:          salt,
        Nonce:         nonce,
    }
    copy(header.Magic[:], headerMagic)
    headerBytes := encodeHeader(header)

    // 4. Encrypt payload with AES-256-GCM, using header bytes as AAD
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    // Seal: nonce is in the header, not prepended to ciphertext
    ciphertext := gcm.Seal(nil, nonce[:], payload, headerBytes[:])

    // 5. Concatenate header + ciphertext
    result := make([]byte, headerSize+len(ciphertext))
    copy(result[:headerSize], headerBytes[:])
    copy(result[headerSize:], ciphertext)
    return result, nil
}

// DecryptVault decrypts a vault file using the master password.
// Returns the raw JSON payload.
func DecryptVault(masterPassword []byte, fileData []byte) ([]byte, error) {
    // 1. Parse and validate the header
    header, err := decodeHeader(fileData)
    if err != nil {
        return nil, err
    }

    // 2. Derive the key using parameters stored in the header
    key := argon2.IDKey(
        masterPassword,
        header.Salt[:],
        uint32(header.Argon2Time),
        header.Argon2Memory,
        uint32(header.Argon2Threads),
        uint32(header.Argon2KeyLen),
    )
    defer zeroBytes(key)

    // 3. Decrypt — GCM authentication verifies both the ciphertext AND the header (AAD)
    headerBytes := encodeHeader(header)  // re-encode to get the exact bytes used as AAD
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    ciphertext := fileData[headerSize:]
    plaintext, err := gcm.Open(nil, header.Nonce[:], ciphertext, headerBytes[:])
    if err != nil {
        // This error means either: wrong password, corrupted file, or tampered header.
        return nil, fmt.Errorf("decryption failed: incorrect password or corrupted vault")
    }

    return plaintext, nil
}

// zeroBytes overwrites a byte slice with zeros and prevents the compiler from eliding it.
func zeroBytes(b []byte) {
    for i := range b {
        b[i] = 0
    }
    runtime.KeepAlive(b)
}
```

---

## 8. Summary of Decisions

| Area | Decision | Confidence |
|------|----------|------------|
| AES-GCM | Use `crypto/aes` + `cipher.NewGCM` + separate nonce in header | HIGH |
| Nonce generation | `crypto/rand.Read` into a 12-byte array per save | HIGH |
| AAD | Pass raw 62-byte header as GCM additional authenticated data | HIGH |
| Argon2id | `golang.org/x/crypto/argon2.IDKey`, memory=64MiB, time=3, threads=4, salt=32B | HIGH |
| Argon2 params stored in header | Yes — all params (time, memory, threads, key len, salt) in unencrypted header | HIGH |
| File format | Fixed 62-byte binary header + variable GCM ciphertext | HIGH |
| Format version field | `uint16` at bytes 8–9; handled via switch-case in `openVault` | HIGH |
| Magic bytes | `"ABDITUM\n"` (8 bytes, non-printable last byte prevents shell double-click issues) | HIGH |
| Atomic save | Write to `.abditum.tmp`, `fsync`, `os.Rename` to `.abditum`; copy old to `.bak` first | HIGH |
| Windows rename | `os.Rename` on Windows calls `MoveFileExW` + `MOVEFILE_REPLACE_EXISTING` internally — works, but add retry on `ACCESS_DENIED` for AV conflicts | MEDIUM |
| DO NOT use `google/renameio` | Explicitly drops Windows support | HIGH |
| Secure memory | `github.com/awnumar/memguard` for master password + derived key + decrypted payload | MEDIUM (API experimental but stable in practice) |
| Manual zeroing | `for i := range b { b[i] = 0 }` + `runtime.KeepAlive(b)` for short-lived slices | HIGH |
| Clipboard | `github.com/atotto/clipboard` — pure Go, no CGO, works on Win/Mac/Linux | HIGH |
| Clear clipboard | `clipboard.WriteAll("")` — no native clear API; empty write is sufficient | HIGH |
| DO NOT use `golang.design/x/clipboard` | Requires CGO on macOS + Linux; panics when CGO_ENABLED=0 | HIGH |

---

## 9. Open Questions / Flags for Future Phases

1. **Argon2id parameter tuning:** The defaults above (64 MiB, time=3, 4 threads) target ~1–2s on a modern laptop. Consider benchmarking on a Raspberry Pi 4 and a 2015-era Intel Core i5 to validate acceptable UX. If too slow, reduce memory to 32 MiB and time to 2. **Expose these as user-configurable parameters in vault settings for v2.**

2. **Windows antivirus and `.tmp` rename failures:** Real-world Windows AV software (Windows Defender, Malwarebytes) is known to briefly lock files during scanning, causing `os.Rename` to return `ACCESS_DENIED`. The atomic save function should implement a small retry loop (3 attempts, 100ms apart) before surfacing an error to the user.

3. **Linux without xclip/xsel/wl-clipboard:** The clipboard will be unavailable in minimal Linux environments. The app must handle this gracefully — show an informational message rather than crash. Test in a minimal Docker container to verify behavior.

4. **memguard `Purge()` on signal handling:** `memguard.CatchInterrupt()` registers a handler for `SIGINT`/`SIGTERM` that calls `Purge()`. This may conflict with Bubble Tea's signal handling. Coordinate signal registration to avoid double-handling.

5. **Payload JSON encoding:** Use `encoding/json` from stdlib. The payload will be UTF-8 encoded (Go `encoding/json` always outputs UTF-8). Verify that all string values (secret names, field values) round-trip correctly with non-ASCII characters (the spec mentions Portuguese characters like ê, ã, ç).

6. **Key re-derivation on password change:** When the user changes the master password, a new salt AND new nonce must be generated, and the vault must be re-encrypted. Do NOT reuse the old salt with a new password.

---

## Sources

- Go stdlib: `go doc crypto/aes`, `go doc crypto/cipher`, `go doc crypto/rand`, `go doc os.Rename`, `go doc os.File.Sync` (verified against go1.26.1)
- RFC 9106 — Argon2 Memory-Hard Function (IETF, Sep 2021) — parameter recommendations §4
- `github.com/awnumar/memguard` — source + README inspected (2.7k stars, v0.23.0, Aug 2025)
- `github.com/atotto/clipboard` — source inspected: `clipboard_windows.go`, `clipboard_unix.go`, `clipboard_darwin.go` (1.4k stars)
- `github.com/golang-design/clipboard` — source inspected: `clipboard_nocgo.go` confirms CGO requirement (770 stars, v0.7.1, Jun 2025)
- `github.com/google/renameio` — README confirmed: no Windows export (672 stars)
- `github.com/golang/go/issues/22397` — Go team discussion on atomic file replacement confirming Windows limitations
- `descricao.md` + `.planning/PROJECT.md` — Abditum project specification (2026-03-24)
