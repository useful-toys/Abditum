package tui

import "time"

// workArea describes what is currently mounted in the work area.
type workArea int

const (
	workAreaPreVault  workArea = iota // 0: welcome screen (ASCII art background)
	workAreaVault                     // 1: vault open — tree + detail side by side
	workAreaTemplates                 // 2: template editor — list + detail side by side
	workAreaSettings                  // 3: settings screen
)

// tickMsg is sent every second to drive inactivity and clipboard timers.
// The global tick starts only after the vault is opened (workAreaVault).
type tickMsg time.Time

// Domain message types — broadcast to all live models on relevant vault events.
// Children that don't care about a type simply ignore it via default case.

// secretAddedMsg is sent when a secret is created or duplicated.
type secretAddedMsg struct{ id string }

// secretDeletedMsg is sent when a secret is marked for soft-deletion.
type secretDeletedMsg struct{ id string }

// secretRestoredMsg is sent when a soft-deletion mark is removed.
type secretRestoredMsg struct{ id string }

// secretModifiedMsg is sent when a secret's values or structure change.
type secretModifiedMsg struct{ id string }

// secretMovedMsg is sent when a secret is moved between folders.
type secretMovedMsg struct {
	id         string
	fromFolder string
	toFolder   string
}

// secretReorderedMsg is sent when a secret is reordered within a folder.
type secretReorderedMsg struct{}

// folderStructureChangedMsg is sent on any folder create/rename/move/reorder/delete.
type folderStructureChangedMsg struct{}

// vaultSavedMsg is sent when the vault is written to disk.
// After this, soft-deleted secrets are removed from memory.
type vaultSavedMsg struct{}

// vaultReloadedMsg is sent on full reload from disk (all children reset state).
type vaultReloadedMsg struct{}

// vaultClosedMsg is sent when the vault is locked or closed.
// All children must wipe sensitive memory on receiving this.
type vaultClosedMsg struct{}

// vaultChangedMsg is a generic fallback when the specific mutation type
// is not relevant to broadcast.
type vaultChangedMsg struct{}

// pushModalMsg is emitted by dialog factories; rootModel appends the modal
// to its stack on receiving this message.
type pushModalMsg struct {
	modal *modalModel
}
