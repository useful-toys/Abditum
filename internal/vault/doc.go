package vault

// Manager is a stub for the vault manager.
// Full implementation is in phases 02-04.
type Manager struct{}

// IsModified returns whether the vault has unsaved changes.
func (m *Manager) IsModified() bool {
	return false
}
