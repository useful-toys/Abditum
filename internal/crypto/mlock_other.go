//go:build !unix && !windows

package crypto

// mlock is not available on this platform.
// Returns ErrMLockFailed to signal unavailability (non-fatal per D-03).
func mlock(b []byte) error {
	return ErrMLockFailed
}

// munlock is a no-op on platforms without mlock support.
func munlock(b []byte) error {
	return nil
}
