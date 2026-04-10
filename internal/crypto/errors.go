package crypto

import "errors"

// Sentinel errors for the crypto package.
//
// These errors are returned by crypto functions to indicate specific failure
// conditions. They are designed to be checked with errors.Is() and provide
// minimal information to avoid leaking internal implementation details.
//
// Example usage:
//
//	plaintext, err := crypto.Decrypt(key, ciphertext)
//	if errors.Is(err, crypto.ErrAuthFailed) {
//	    // Wrong password or corrupted data - allow retry
//	    return promptForPassword()
//	}

// ErrAuthFailed is returned by Decrypt when authentication fails.
//
// This can occur for two reasons:
//  1. The provided key is incorrect (wrong password)
//  2. The ciphertext has been corrupted or tampered with
//
// Returning a single error for both cases prevents timing attacks that
// could distinguish between "wrong password" and "corrupted data".
//
// This error indicates that decryption should be retried with a different
// password, but not indefinitely. After several failures, the application
// should assume the vault is corrupted.
var ErrAuthFailed = errors.New("authentication failed")

// ErrInsufficientEntropy is returned when crypto/rand.Reader cannot provide
// enough random bytes.
//
// This is an extremely rare condition that typically indicates:
//   - The operating system's entropy pool is exhausted
//   - The system is under severe resource pressure
//   - A kernel bug or misconfiguration
//
// When this error occurs, the operation should be retried after a short delay.
// If it persists, the application should fail gracefully rather than proceeding
// with weak cryptographic material.
var ErrInsufficientEntropy = errors.New("insufficient entropy")

// ErrInvalidParams is returned when a function receives invalid parameters.
//
// Common causes:
//   - nil or empty password slice
//   - nil or empty salt slice
//   - Invalid Argon2id parameters (e.g., keyLen=0)
//   - Wrong key size for AES (not 32 bytes)
//
// This error indicates a programming error in the caller, not a runtime
// condition. The caller should validate inputs before calling crypto functions.
var ErrInvalidParams = errors.New("invalid parameters")

// ErrMLockFailed is returned when memory locking (mlock/VirtualLock) fails.
//
// Memory locking prevents sensitive data from being swapped to disk. However,
// it is not always available:
//   - Containers often disable mlock
//   - Unprivileged users may not have permission
//   - Some operating systems don't support it
//   - On "other" platforms (not Unix/Windows), it is never available
//
// This error is NON-FATAL. The caller should log a warning and continue
// execution. Abditum can operate securely without memory locking, though
// sensitive data may be swapped to disk in low-memory situations.
//
// Example handling:
//
//	key := make([]byte, 32)
//	if err := crypto.Mlock(key); err != nil {
//	    if errors.Is(err, crypto.ErrMLockFailed) {
//	        log.Warn("Memory locking unavailable - sensitive data may be swapped")
//	        // Continue execution normally
//	    }
//	}
var ErrMLockFailed = errors.New("memory lock failed")
