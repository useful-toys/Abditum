//go:build windows

package storage

import (
	"golang.org/x/sys/windows"
)

// atomicRename performs an atomic file rename on Windows.
//
// Standard os.Rename on Windows is NOT atomic when replacing an existing file --
// it can leave the destination deleted if the process is interrupted between
// the internal DeleteFile and MoveFile calls.
// MoveFileEx with MOVEFILE_REPLACE_EXISTING provides atomic replace semantics
// on NTFS (the rename is a single metadata operation in the MFT).
func atomicRename(src, dst string) error {
	srcPtr, err := windows.UTF16PtrFromString(src)
	if err != nil {
		return err
	}
	dstPtr, err := windows.UTF16PtrFromString(dst)
	if err != nil {
		return err
	}
	return windows.MoveFileEx(srcPtr, dstPtr, windows.MOVEFILE_REPLACE_EXISTING)
}
