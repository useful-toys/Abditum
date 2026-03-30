//go:build !windows

package storage

import "os"

// atomicRename performs an atomic file rename on Unix systems.
//
// POSIX rename(2) is atomic when source and destination are on the same filesystem.
// The .tmp file MUST be in the same directory as the target to guarantee
// same-filesystem operation -- never use os.TempDir() which may be a different mount.
func atomicRename(src, dst string) error {
	return os.Rename(src, dst)
}
