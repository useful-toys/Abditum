// Package crypto provides cryptographic primitives for Abditum password manager.
//
// # Key Derivation Function (KDF)
//
// Argon2id is used for deriving encryption keys from user passwords. Argon2id
// is a Key Derivation Function (KDF) — NOT an encryption algorithm itself.
// Its purpose is to convert a human-memorable password into a strong cryptographic
// key suitable for AES-256 encryption.
//
// Argon2id combines the memory-hard properties of Argon2i with the side-channel
// resistance of Argon2d, making it resistant to both GPU-based brute-force attacks
// and timing attacks. It is the current best practice for password-based key derivation
// (RFC 9106).
//
// # Chosen Parameters
//
// For Abditum format version 1, we use the following Argon2id parameters:
//
//   - t=3: Number of iterations (time cost)
//   - m=262144: Memory cost in KiB (equivalent to 256 MiB)
//   - p=4: Parallelism (number of threads)
//   - keyLen=32: Output key length in bytes (256 bits for AES-256)
//
// These parameters are chosen to provide strong security while remaining
// practical on consumer hardware. Key derivation takes approximately 200-500ms
// on modern CPUs, which provides resistance against brute-force attacks while
// not significantly impacting user experience during vault open operations.
//
// CRITICAL: The m parameter is specified in KiB (kibibytes), not bytes.
// 256 MiB = 262144 KiB. This is a common source of errors when using Argon2id.
//
// # AES-256-GCM Authenticated Encryption
//
// AES-256-GCM (Galois/Counter Mode) is used for authenticated encryption with
// associated data (AEAD). This combines confidentiality (encryption) and
// authenticity (authentication) in a single primitive.
//
// GCM requires a unique nonce (number used once) for every encryption operation
// with the same key. Nonce reuse with the same key catastrophically breaks GCM
// security. This package generates a fresh 12-byte nonce using crypto/rand for
// every Encrypt() call.
//
// # Ciphertext Format
//
// The output of Encrypt() consists of:
//
//   - Nonce: 12 bytes (prepended)
//   - Ciphertext: Variable length (same as plaintext)
//   - Authentication Tag: 16 bytes (appended by GCM)
//
// Total overhead: 28 bytes (12-byte nonce + 16-byte tag)
//
// The nonce is prepended to the ciphertext so that Decrypt() can extract it.
// The authentication tag is automatically included by gcm.Seal() and verified
// by gcm.Open().
//
// # Memory Security
//
// All sensitive data (passwords, keys, plaintext) must be handled as []byte slices,
// never as strings. Go strings are immutable and cannot be securely zeroed.
//
// Callers are responsible for zeroing sensitive buffers after use with ZeroBytes().
// This package does not automatically zero buffers to give callers full control
// over memory management and timing.
//
// Platform-specific memory locking (mlock on Unix, VirtualLock on Windows) is
// available to prevent sensitive data from being swapped to disk. Memory locking
// is best-effort and failure is non-fatal — the package continues to operate
// normally even when memory locking is unavailable (e.g., in containers or
// low-privilege contexts).
package crypto

// ArgonParams holds the parameters for Argon2id key derivation.
//
// These parameters control the computational cost and memory requirements
// of the key derivation process. Higher values provide better security
// against brute-force attacks but increase the time required to derive a key.
type ArgonParams struct {
	// Time is the number of iterations (time cost).
	// For Abditum format version 1, this is 3.
	Time uint32

	// Memory is the amount of memory to use in KiB (NOT bytes).
	// For Abditum format version 1, this is 262144 (256 MiB).
	Memory uint32

	// Threads is the degree of parallelism (number of threads).
	// For Abditum format version 1, this is 4.
	Threads uint8

	// KeyLen is the length of the derived key in bytes.
	// For Abditum format version 1, this is 32 (256 bits for AES-256).
	KeyLen uint32
}

// FormatVersion is the current Abditum vault format version.
//
// This version number is stored in the plaintext header of .abditum files
// and determines which cryptographic parameters and file structure to use.
//
// Version 1 uses:
//   - Argon2id with t=3, m=262144 KiB, p=4, keyLen=32
//   - AES-256-GCM with 12-byte nonce and 16-byte tag
//   - 32-byte salt generated with crypto/rand
const FormatVersion = 1
