package storage

import "os"

// RecoverOrphans cleans up stale temporary files left by a previous process
// that was interrupted during a Save operation.
//
// RecoverOrphans should be called once at startup, before any Load or Save call.
// It is a no-op when the vault is in a clean state.
//
// Current recovery actions:
//   - If vaultPath + ".tmp" exists, it is removed. A stale .tmp means a previous
//     Save was interrupted after the tmp write but before the atomic rename. It is
//     always safe to delete because the original vault file was not yet replaced.
//
// RecoverOrphans does NOT automatically restore from .bak on a corrupt vault.
// That decision is left to the user/operator to avoid silent data loss from
// unrelated corruption. If the vault file is unreadable, Load will return an
// appropriate error and the user can manually restore from .bak.
//
// Parameters:
//   - vaultPath: Absolute path to the vault file (e.g. "/home/user/myvault.abditum").
//
// Returns nil on success (including when no orphans were found).
// Returns an error if a stale .tmp file could not be removed.
func RecoverOrphans(vaultPath string) error {
	tmpPath := vaultPath + ".tmp"
	if _, err := os.Stat(tmpPath); err == nil {
		// .tmp exists — remove it
		if err := os.Remove(tmpPath); err != nil {
			return err
		}
	}
	return nil
}
