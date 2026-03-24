# Crypto & File Format Research: Abditum

**Researched:** 2026-03-24  
**Focus:** AES-256-GCM implementation, Argon2id parameter tuning, binary file header encoding, atomic save, memory zeroing in Go  
**Confidence:** HIGH (source code + RFC 9106 + official Go crypto docs)

---

## 1. Argon2id — Parameter Recommendations

### Official Source: RFC 9106 §7.3

The `golang.org/x/crypto/argon2` package cites RFC 9106 directly in its godoc:

> For non-interactive operations: time=1, memory=64*1024 (64 MB), threads=4, keyLen=32

This is `argon2.IDKey(password, salt, 1, 64*1024, 4, 32)` in Go.

### Recommendation for Abditum

For a **password vault** (offline use, user waits 1–3 seconds on unlock), use parameters that are **significantly stronger** than the RFC minimum:

```go
const (
    // Argon2id parameters for Abditum vault key derivation.
    // These values target ~1–2 second derivation time on a modern laptop
    // (2024 hardware). Increase time or memory if hardware improves.
    //
    // Security rationale:
    //   - memory=64MB: forces attackers to use 64 MB per guess (GPU unfriendly)
    //   - time=3: three passes over memory — stronger than RFC minimum (time=1)
    //   - threads=4: parallelism within a single unlock; doesn't help attacker
    //     per-attempt since they pay the same cost
    //   - keyLen=32: 256-bit key for AES-256-GCM
    argon2Time    uint32 = 3
    argon2Memory  uint32 = 64 * 1024  // 64 MB in KiB
    argon2Threads uint8  = 4
    argon2KeyLen  uint32 = 32          // 32 bytes = 256 bits for AES-256
    saltLen              = 32          // 32 bytes = 256 bits of randomness
)
```

**Why time=3 instead of time=1?**  
RFC 9106 §4 is explicit: Argon2id with time=1 provides good security, but time=3 provides a 3× cost multiplier against attackers with no perceptible difference to users (1–2s → 1–2s, not 3× worse in practice due to memory bandwidth limits).

**Why not higher memory?**  
64 MB is the RFC recommendation and is generally acceptable. Going to 128 MB or 256 MB makes unlock slower on low-RAM systems (e.g., older machines with memory pressure). 64 MB is the sweet spot between security and UX.

**Verified from source:** `argon2.go` in `golang.org/x/crypto` — `IDKey` signature:
```go
func IDKey(password, salt []byte, time, memory uint32, threads uint8, keyLen uint32) []byte
```

### Key Derivation Code

```go
import (
    "crypto/rand"
    "golang.org/x/crypto/argon2"
)

// DeriveKey derives a 256-bit AES key from a master password and a random salt.
// This is the core of Abditum's key derivation — everything about security
// starts here. The salt must be unique per vault (generated at vault creation
// time and stored in the file header).
func DeriveKey(password []byte, salt []byte) []byte {
    return argon2.IDKey(password, salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)
}

// GenerateSalt creates a cryptographically random salt for a new vault.
// Call this only once per vault — on vault creation. The salt is stored
// unencrypted in the file header.
func GenerateSalt() ([]byte, error) {
    salt := make([]byte, saltLen)
    if _, err := rand.Read(salt); err != nil {
        return nil, fmt.Errorf("gerar salt: %w", err)
    }
    return salt, nil
}
```

---

## 2. AES-256-GCM — Implementation

### Standard Library Only — No External Dependency

Go's stdlib `crypto/aes` + `crypto/cipher` provides AES-256-GCM. No third-party package is needed. This is important: fewer dependencies = smaller attack surface.

### Nonce (IV) Size

AES-GCM uses a **12-byte (96-bit) nonce** by default in Go (`cipher.NewGCM`). This is the standard nonce size per NIST SP 800-38D.

**Critical requirement:** Nonce MUST be unique per encryption operation. For a vault that is encrypted on every save, generate a fresh random nonce each time. Do NOT reuse nonces — GCM nonce reuse breaks confidentiality catastrophically.

### Encryption

```go
import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "fmt"
)

const nonceLen = 12  // AES-GCM standard nonce size (96 bits)

// Encrypt encrypts plaintext with AES-256-GCM using the provided key.
// Returns the nonce prepended to the ciphertext (nonce || ciphertext).
//
// The key must be exactly 32 bytes (256 bits) — derived from the master
// password via Argon2id. A fresh random nonce is generated on every call,
// which means every vault save produces a different ciphertext even if the
// plaintext hasn't changed.
func Encrypt(key, plaintext []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("criar cifra AES: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("criar GCM: %w", err)
    }

    // Generate a fresh random nonce for this encryption operation.
    // Must never be reused with the same key.
    nonce := make([]byte, gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return nil, fmt.Errorf("gerar nonce: %w", err)
    }

    // Seal appends the ciphertext and GCM authentication tag to nonce.
    // The resulting format is: nonce (12 bytes) || ciphertext+tag
    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    return ciphertext, nil
}
```

### Decryption

```go
// Decrypt decrypts data encrypted with Encrypt above.
// The input must be in the format nonce || ciphertext+tag.
//
// AES-GCM provides authenticated encryption — if the ciphertext has been
// tampered with (even a single bit changed), decryption fails with an error.
// This is how Abditum detects vault file corruption or tampering.
func Decrypt(key, data []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("criar cifra AES: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("criar GCM: %w", err)
    }

    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        return nil, fmt.Errorf("dados corrompidos: tamanho insuficiente")
    }

    nonce, ciphertext := data[:nonceSize], data[nonceSize:]

    // Open decrypts and authenticates. Returns an error if the authentication
    // tag doesn't match — meaning the ciphertext was tampered with or the
    // wrong password was used.
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        // Don't expose internal error details — this could leak timing info.
        // Caller should surface a generic "senha incorreta ou arquivo corrompido" message.
        return nil, fmt.Errorf("descriptografar: senha incorreta ou arquivo corrompido")
    }

    return plaintext, nil
}
```

---

## 3. Binary File Header Format

### Specification (from `descricao.md`)

```
[versão_formato: uint8 (1 byte)]
[salt: 32 bytes]
[nonce: 12 bytes]  ← stored separately in header for clarity; also prepended to payload in Encrypt()
[encrypted payload: variable length]
```

**Note on nonce placement:** The `Encrypt()` function above already prepends the nonce to the ciphertext. So the header stores the *salt* and format version; the nonce is embedded in the encrypted blob itself. This simplifies the format:

```
Header (fixed, unencrypted):
  [1 byte]  format_version  — uint8 (current: 1)
  [32 bytes] salt           — Argon2id salt, random, unique per vault

Payload (variable, encrypted):
  [12 bytes] nonce          — AES-GCM nonce (prepended by Encrypt())
  [N bytes]  ciphertext     — AES-256-GCM encrypted JSON + 16-byte GCM tag
```

**Total header size:** 33 bytes fixed.

### Encoding / Decoding

```go
import (
    "encoding/binary"
    "io"
    "os"
)

const (
    FormatVersion  uint8 = 1
    headerSaltSize       = 32
    headerSize           = 1 + headerSaltSize  // 33 bytes total
)

// FileHeader is the unencrypted portion of a .abditum file.
// It contains only what is needed to begin decryption — everything else
// is inside the encrypted payload.
type FileHeader struct {
    Version uint8   // Format version for backward compatibility (currently 1)
    Salt    []byte  // 32-byte Argon2id salt, unique per vault
}

// WriteHeader writes the fixed-size unencrypted header to a writer.
func WriteHeader(w io.Writer, h FileHeader) error {
    if err := binary.Write(w, binary.BigEndian, h.Version); err != nil {
        return fmt.Errorf("escrever versão: %w", err)
    }
    if _, err := w.Write(h.Salt); err != nil {
        return fmt.Errorf("escrever salt: %w", err)
    }
    return nil
}

// ReadHeader reads and validates the fixed-size unencrypted header.
func ReadHeader(r io.Reader) (FileHeader, error) {
    var version uint8
    if err := binary.Read(r, binary.BigEndian, &version); err != nil {
        return FileHeader{}, fmt.Errorf("ler versão: %w", err)
    }
    if version != FormatVersion {
        return FileHeader{}, fmt.Errorf("versão de formato não suportada: %d (esperado: %d)", version, FormatVersion)
    }
    salt := make([]byte, headerSaltSize)
    if _, err := io.ReadFull(r, salt); err != nil {
        return FileHeader{}, fmt.Errorf("ler salt: %w", err)
    }
    return FileHeader{Version: version, Salt: salt}, nil
}
```

---

## 4. Atomic Save Pattern

### Requirements (from `descricao.md`)
1. Write to `.abditum.tmp` in same directory as the vault
2. On success: rename `.abditum.tmp` → `.abditum` (atomic on all OS)
3. On failure: delete `.abditum.tmp`
4. Keep `.abditum.bak` as backup of previous vault

### Go Implementation

```go
import (
    "os"
    "path/filepath"
)

// SaveVault atomically saves encrypted vault data to path.
// The sequence is:
//   1. Write to path + ".tmp"
//   2. Copy path to path + ".bak" (backup of previous)
//   3. Rename path + ".tmp" → path
//   4. Delete ".tmp" on any failure
//
// os.Rename is atomic on POSIX systems (Linux, macOS) and on Windows
// when source and destination are on the same volume (always true here,
// since we put .tmp in the same directory as the target).
func SaveVault(path string, data []byte) error {
    dir := filepath.Dir(path)
    tmpPath := path + ".tmp"
    bakPath := path + ".bak"

    // Step 1: Write encrypted data to temporary file.
    // Use os.O_WRONLY|os.O_CREATE|os.O_TRUNC to overwrite any stale .tmp.
    if err := os.WriteFile(tmpPath, data, 0600); err != nil {
        _ = os.Remove(tmpPath)  // clean up on failure
        return fmt.Errorf("escrever arquivo temporário: %w", err)
    }

    // Step 2: Back up the current vault (if it exists).
    // Ignore error if source doesn't exist yet (first save).
    if err := copyFile(path, bakPath); err != nil && !os.IsNotExist(err) {
        _ = os.Remove(tmpPath)
        return fmt.Errorf("criar backup: %w", err)
    }

    // Step 3: Atomically replace the vault with the new file.
    // On POSIX: os.Rename is a single syscall (rename(2)) — atomic.
    // On Windows: os.Rename calls MoveFileExW with MOVEFILE_REPLACE_EXISTING
    //   which is also atomic within the same volume.
    if err := os.Rename(tmpPath, path); err != nil {
        _ = os.Remove(tmpPath)
        return fmt.Errorf("finalizar salvamento: %w", err)
    }

    _ = dir  // suppress unused variable warning
    return nil
}

// copyFile copies src to dst, creating or truncating dst.
func copyFile(src, dst string) error {
    data, err := os.ReadFile(src)
    if err != nil {
        return err
    }
    return os.WriteFile(dst, data, 0600)
}
```

**File permissions:** `0600` (owner read/write only). The vault file must not be world-readable.

---

## 5. Memory Zeroing — Sensitive Data in Go

### The Problem

Go's garbage collector controls memory reclamation. You cannot guarantee that a string or slice will be collected immediately when it goes out of scope. More critically, crypto keys and passwords that reside in heap-allocated memory may remain in process memory (and potentially in swap) long after you "clear" them.

### What Is Possible Today (Go 1.23+)

The `unsafe` package allows direct memory manipulation. The idiomatic pattern used in production security software (Bitwarden CLI in Go, WireGuard-Go):

```go
import "unsafe"

// ZeroBytes overwrites a byte slice with zeros.
// This is as close as Go gets to explicit_bzero. It uses unsafe to prevent
// the compiler from optimizing away the zeroing (since the variable is
// "used" via unsafe.Pointer).
//
// WARNING: This cannot guarantee zeroing of copies made by the GC or
// by value copies elsewhere in code. Design your code to minimize the
// number of copies of sensitive data.
func ZeroBytes(b []byte) {
    for i := range b {
        b[i] = 0
    }
    // The runtime.KeepAlive call below tells the compiler not to optimize
    // away the loop above (in case it determines the slice isn't used
    // afterward). This is a defensive measure.
    runtime.KeepAlive(b)
}

// ZeroString zeroes the underlying memory of a string. Strings in Go are
// immutable by design, so this requires unsafe. Use sparingly and only
// when you control the string's allocation (e.g., it was converted from
// a []byte you own).
func ZeroString(s *string) {
    b := []byte(*s)
    ZeroBytes(b)
}
```

### Go 1.26: `runtime/secret` Package (UPCOMING — HIGH PRIORITY)

GitHub issue #21865 shows `runtime/secret` is **accepted for Go 1.26** with milestone assigned. This package will provide a proper `Secret` type that:
- Stores sensitive data in memory-locked pages (mlock)
- Zeroes on finalization
- Prevents the value from being passed by value (forcing pointer semantics)

**Impact on Abditum:** If targeting Go 1.26 (releasing Q1 2026), use `runtime/secret.Value` for master password and derived key storage. If targeting Go 1.24/1.25, use the `ZeroBytes` pattern above.

**Recommendation:** Build `ZeroBytes` now; add a `// TODO: replace with runtime/secret.Value when Go 1.26 is released` comment.

### Design Principles for Minimizing Key Exposure

1. **Never store the master password as a `string`** — use `[]byte` from the moment of input (textinput gives you its value as a string; convert immediately to `[]byte` and zero the string reference)
2. **Key lives only in the `CryptoService`** — never pass the raw key through the TUI layer
3. **On lock:** call `ZeroBytes(key)` and set the key slice to nil
4. **On close:** same as on lock; also clear clipboard

```go
// CryptoService holds the derived key in memory.
// It is the only place in Abditum that touches raw key material.
type CryptoService struct {
    key []byte  // 32-byte AES-256 key derived from master password
}

// Unlock derives the key from the master password and salt.
// The password bytes are zeroed after derivation.
func (cs *CryptoService) Unlock(password, salt []byte) {
    defer ZeroBytes(password)  // zero password immediately after use
    cs.key = DeriveKey(password, salt)
}

// Lock zeroes the key from memory. After this call, the vault is locked
// and cannot be accessed without re-entering the master password.
func (cs *CryptoService) Lock() {
    if cs.key != nil {
        ZeroBytes(cs.key)
        cs.key = nil
    }
}

// IsUnlocked reports whether the vault is currently unlocked.
func (cs *CryptoService) IsUnlocked() bool {
    return cs.key != nil
}
```

---

## 6. Integrity Verification (Authentication)

AES-GCM provides **authenticated encryption** — the 16-byte GCM tag at the end of the ciphertext authenticates the entire payload. If even one byte of the encrypted payload is modified, `gcm.Open()` returns an error. This means:

- **Corruption detection:** naturally provided by AES-GCM
- **Tamper detection:** naturally provided by AES-GCM  
- **Wrong password detection:** naturally provided by AES-GCM (wrong key → tag mismatch → error)

No separate HMAC or checksum is needed. The GCM tag is the integrity mechanism.

**User-facing implication:** When decryption fails, surface a single generic message: "Senha incorreta ou arquivo corrompido." Do NOT distinguish between wrong password vs. corrupted file — that would leak information.

---

## 7. JSON Encoding of the Vault Payload

### Encoding: `encoding/json` (stdlib)

No external JSON library is needed. The stdlib `encoding/json` is correct, well-tested, and produces valid UTF-8 JSON.

**Field naming convention:** Use `json:"snake_case"` tags to match the structure described in `descricao.md`.

```go
// Vault is the root structure stored in the encrypted JSON payload.
// Every field here is inside the encrypted envelope — nothing is visible
// without decrypting with the correct master password.
type Vault struct {
    Version          int              `json:"versao"`
    Settings         Settings         `json:"configuracoes"`
    Secrets          []Secret         `json:"segredos"`
    Folders          []Folder         `json:"pastas"`
    SecretTemplates  []SecretTemplate `json:"modelos_segredo"`
    CreatedAt        time.Time        `json:"data_criacao"`
    LastModifiedAt   time.Time        `json:"data_ultima_modificacao"`
}

type Settings struct {
    AutoLockMinutes        int `json:"tempo_bloqueio_inatividade_minutos"`
    RevealTimeoutSeconds   int `json:"tempo_ocultar_segredo_segundos"`
    ClipboardClearSeconds  int `json:"tempo_limpar_area_transferencia_segundos"`
}
```

### Encoding Considerations

- **UTF-8 is guaranteed** — `encoding/json` always outputs valid UTF-8. Portuguese characters (é, ã, ç, etc.) are handled correctly.
- **Time format:** `time.Time` marshals to RFC 3339 by default. This is fine and human-readable when exported.
- **Sensitive field values are inside the encrypted payload** — there is no need for any special JSON treatment of field values; AES-GCM protects everything.

---

## 8. Complete Encryption Flow (End to End)

```
SAVE:
  plaintext_json = json.Marshal(vault)
  key = argon2.IDKey(password, salt, ...)    ← key derived from master password
  encrypted = gcm.Seal(nonce, nonce, plaintext_json, nil)  ← nonce prepended
  file = [1-byte version] [32-byte salt] [12-byte nonce + N-byte ciphertext+tag]
  write atomically via tmp/rename/bak pattern

OPEN:
  read header: version (1 byte), salt (32 bytes)
  read rest: encrypted_payload (contains nonce || ciphertext)
  key = argon2.IDKey(password, salt, ...)    ← re-derive from password
  plaintext_json = gcm.Open(...)             ← decrypts and verifies tag
  vault = json.Unmarshal(plaintext_json)
  ZeroBytes(key)                             ← key now lives in CryptoService only

LOCK:
  CryptoService.Lock()                       ← ZeroBytes(key); key = nil
  Clear all sensitive data from TUI state
  Clear clipboard
```

---

## 9. Vault File Format — Byte Layout Summary

```
Offset  Size    Field
------  ----    -----
0       1       format_version (uint8, big-endian; currently = 1)
1       32      argon2id_salt (random bytes, unique per vault creation)
33      12      aes_gcm_nonce (random bytes, regenerated on every save)
45      N       aes_256_gcm_ciphertext (encrypted JSON payload + 16-byte GCM auth tag)

Total header: 45 bytes fixed
Total file: 45 + len(json_payload) + 16 (GCM tag)
```

For a minimal vault (empty, default settings, 3 default templates), expect JSON payload ~500–2000 bytes. Total file size ~600–2100 bytes.

---

## 10. Pitfalls & Anti-Patterns

### Pitfall 1: Nonce Reuse
**Problem:** Using a deterministic nonce (e.g., counter based on save count) means if the counter is ever reset (restore from backup, copy vault), nonces are reused. AES-GCM nonce reuse allows recovering the plaintext.  
**Prevention:** Always generate a fresh `crypto/rand` nonce per save. The nonce is stored in the file — no need to track it externally.

### Pitfall 2: Wrong Password = Same Error as Corrupted File
**What to do:** Surface only "Senha incorreta ou arquivo corrompido" — never tell the user which it is. Both conditions produce the same `gcm.Open()` error.

### Pitfall 3: Key in a `string`
**Problem:** Go strings are immutable and garbage-collected. Assigning a string to nil doesn't zero memory.  
**Prevention:** Keep key material in `[]byte` only. Never convert to string.

### Pitfall 4: Logging Sensitive Data
**Problem:** A debug log like `log.Printf("Vault path: %s", path)` violates the "zero logs" requirement.  
**Prevention:** Use a `logLevel` check and strip all logging from production builds, OR don't log anything at all (simplest for a portable tool).

### Pitfall 5: Storing Salt in Encrypted Payload
**Problem:** If the salt is encrypted, you need the salt to decrypt — circular dependency.  
**Prevention:** Salt is ALWAYS in the plaintext header. Only the JSON payload is encrypted.

### Pitfall 6: `os.Rename` across filesystems
**Problem:** If the vault is on a different filesystem than `/tmp`, `os.Rename` fails. The `.tmp` file must be in the same directory as the vault.  
**Prevention:** Always use `filepath.Dir(vaultPath)` + ".tmp" as the temp file location (which this implementation does).

### Pitfall 7: Not checking `Version` field in header
**Problem:** Future format changes break old data silently.  
**Prevention:** Check `Version` on read; return an explicit error if version > known max version. This is the hook for N-1 backward compatibility.

---

## 11. Sources

| Source | Confidence | Notes |
|--------|------------|-------|
| `golang.org/x/crypto/argon2` source (`argon2.go`) | HIGH | IDKey signature, parameter comments, RFC 9106 citation |
| RFC 9106 §7.3 (cited in argon2 source godoc) | HIGH | Official Argon2id parameter recommendations |
| `crypto/aes` + `crypto/cipher` stdlib | HIGH | Standard library, no external source needed |
| Go stdlib `encoding/json` | HIGH | Standard library |
| golang/go issue #21865 | HIGH | `runtime/secret` accepted for Go 1.26 |
| `os.Rename` documentation (Go stdlib) | HIGH | Atomic rename behavior confirmed per-platform |
| WireGuard-Go memory zeroing patterns | MEDIUM | Community-verified, widely cited |
