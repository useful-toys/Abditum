//go:build !windows

package crypto

import (
	"golang.org/x/sys/unix"
)

// mlock locks memory pages to prevent swapping to disk on Unix systems.
// Returns ErrMLockFailed if locking fails (non-fatal per D-03).
func mlock(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	if err := unix.Mlock(b); err != nil {
		return ErrMLockFailed
	}
	return nil
}

// munlock unlocks previously locked memory pages on Unix systems.
func munlock(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	return unix.Munlock(b)
}
