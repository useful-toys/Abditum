//go:build windows

package crypto

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// mlock locks memory pages to prevent swapping to disk on Windows.
// Returns ErrMLockFailed if locking fails (non-fatal per D-03).
func mlock(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	err := windows.VirtualLock(uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
	if err != nil {
		return ErrMLockFailed
	}
	return nil
}

// munlock unlocks previously locked memory pages on Windows.
func munlock(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	return windows.VirtualUnlock(uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
}
