package storage

import (
	"crypto/sha256"
	"os"
)

// DetectExternalChange reports whether the vault file at vaultPath has been
// modified since the FileMetadata was last recorded (e.g. immediately after Load).
//
// Detection strategy (per arquitetura.md §7 and historico-decisoes.md §17):
//  1. Fast path: compare file size. A size mismatch is an immediate change signal
//     with O(1) cost (no file read required beyond os.Stat).
//  2. Slow path: if size matches, read the full file and compare SHA-256 hashes.
//     This catches in-place modifications that preserve file size.
//
// Modification time (mtime) is deliberately NOT used because mtime is unreliable
// under cloud sync services (Dropbox, iCloud, Google Drive) and network filesystems.
//
// Parameters:
//   - vaultPath: Absolute path to the vault file.
//   - metadata: FileMetadata recorded at the last successful Load or Save.
//
// Returns:
//   - (false, nil) when the file matches the stored metadata (unchanged).
//   - (true, nil) when the file has been modified (size or content differs).
//   - (false, error) when the file cannot be read (permission error, not found, etc).
func DetectExternalChange(vaultPath string, metadata FileMetadata) (bool, error) {
	info, err := os.Stat(vaultPath)
	if err != nil {
		return false, err
	}

	// Fast path: size mismatch
	if info.Size() != metadata.Size {
		return true, nil
	}

	// Slow path: full content comparison via SHA-256
	data, err := os.ReadFile(vaultPath)
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256(data)
	if hash != metadata.Hash {
		return true, nil
	}

	return false, nil
}

// ComputeFileMetadata reads vaultPath and computes its FileMetadata.
//
// This is a convenience function for callers that need to snapshot the current
// file state without going through a full Load (e.g. after an external tool
// modifies the file and the caller needs to re-anchor its change detection baseline).
//
// Parameters:
//   - vaultPath: Absolute path to the vault file.
//
// Returns the computed FileMetadata and nil, or zero FileMetadata and an error
// if the file cannot be read.
func ComputeFileMetadata(vaultPath string) (FileMetadata, error) {
	data, err := os.ReadFile(vaultPath)
	if err != nil {
		return FileMetadata{}, err
	}
	return FileMetadata{
		Size: int64(len(data)),
		Hash: sha256.Sum256(data),
	}, nil
}
