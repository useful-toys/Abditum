// Package crypto implements stubs for memory locking on unsupported platforms.
// These will be replaced by platform-specific implementations in Task 5.

package crypto

// mlock attempts to lock memory to prevent swapping.
// This stub always returns ErrMLockFailed.
func mlock(b []byte) error {
	return ErrMLockFailed
}

// munlock unlocks previously locked memory.
// This stub always returns nil.
func munlock(b []byte) error {
	return nil
}
