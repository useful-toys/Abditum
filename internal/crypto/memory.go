package crypto

import (
	"runtime"
)

// Wipe overwrites a byte slice with zeros to clear sensitive data from memory.
// This is a best-effort approach - the compiler may optimize away the zeroing
// if it determines the slice is no longer used. For critical secrets, use
// SecureAllocate which includes memory locking on supported platforms.
//
// Wipe is safe to call on nil or empty slices.
func Wipe(data []byte) {
	for i := range data {
		data[i] = 0
	}
	runtime.KeepAlive(data)
}

// SecureAllocate creates a zeroed byte slice of the given size with additional
// security measures:
//   - Memory is locked to prevent swapping to disk (on supported platforms)
//   - Returns a cleanup function that wipes and unlocks the memory
//   - Memory is initially zeroed
//
// The cleanup function should be called via defer to ensure proper cleanup:
//
//	buf, cleanup, err := crypto.SecureAllocate(32)
//	if err != nil {
//	    return err
//	}
//	defer cleanup()
//
// On platforms without memory locking support (or if locking fails), this
// function still returns a usable buffer but memory may be swapped to disk.
// Check the returned error to determine if locking succeeded.
//
// For zero-size allocations, returns an empty slice with no error.
func SecureAllocate(size int) ([]byte, func(), error) {
	if size == 0 {
		return []byte{}, func() {}, nil
	}

	buf := make([]byte, size)

	// Attempt to lock memory (platform-specific)
	unlockErr := mlock(buf)

	cleanup := func() {
		Wipe(buf)
		if unlockErr == nil {
			// Only unlock if lock succeeded
			_ = munlock(buf) // Best effort - ignore unlock errors
		}
	}

	// Return the buffer even if locking failed - caller can decide based on error
	return buf, cleanup, unlockErr
}
