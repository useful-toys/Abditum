package storage

import "fmt"

// MigrationFunc transforms a decrypted JSON payload from format version N to version N+1.
//
// Migration functions receive the raw JSON bytes after decryption and return
// transformed JSON bytes. They must not modify the caller's slice.
type MigrationFunc func(data []byte) ([]byte, error)

// migrations maps source format version to the migration function that upgrades
// that version to the next one (version N → version N+1).
//
// Version 1 is the initial format — there is no migration into it.
// Future migrations are registered here as new format versions are introduced.
// Example:
//
//	1: migrateV1toV2,
var migrations = map[uint8]MigrationFunc{
	// No migrations yet. V1 is the baseline format.
}

// Migrate applies the chain of registered migration functions to transform data
// from fromVersion to toVersion.
//
// If fromVersion == toVersion, data is returned unchanged (no-op).
// If fromVersion > toVersion, an error is returned (downgrade not supported).
// If a migration step is missing from the registry, an error is returned.
//
// Parameters:
//   - data: Decrypted JSON payload at fromVersion.
//   - fromVersion: The format version of the input data.
//   - toVersion: The target format version (usually CurrentFormatVersion).
//
// Returns the migrated JSON bytes and nil, or nil and an error if migration fails.
func Migrate(data []byte, fromVersion, toVersion uint8) ([]byte, error) {
	if fromVersion > toVersion {
		return nil, fmt.Errorf("cannot downgrade from version %d to %d", fromVersion, toVersion)
	}
	if fromVersion == toVersion {
		return data, nil
	}

	current := data
	for v := fromVersion; v < toVersion; v++ {
		fn, ok := migrations[v]
		if !ok {
			return nil, fmt.Errorf("no migration path from version %d to %d", v, v+1)
		}
		var err error
		current, err = fn(current)
		if err != nil {
			return nil, fmt.Errorf("migration v%d to v%d failed: %w", v, v+1, err)
		}
	}

	return current, nil
}
