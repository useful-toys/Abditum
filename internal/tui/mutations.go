package tui

import tea "charm.land/bubbletea/v2"

// Cmd factories for simple vault operations (no modal, no async).
// These wrap vault.Manager calls and return the appropriate domain message.
// Simple operations: favorite, mark for deletion, reorder, rename folder.
// Each factory is a placeholder stub — concrete implementations in Phase 6+.

// favoriteSecretCmd returns a Cmd that toggles the favorite flag on a secret.
// id is the secret's identifier.
func favoriteSecretCmd(id string) tea.Cmd {
	return func() tea.Msg {
		// TODO(phase-8): call vault.Manager.AlternarFavoritoSegredo(id)
		return secretModifiedMsg{id: id}
	}
}

// softDeleteSecretCmd returns a Cmd that marks a secret for deletion.
func softDeleteSecretCmd(id string) tea.Cmd {
	return func() tea.Msg {
		// TODO(phase-8): call vault.Manager.ExcluirSegredo(id)
		return secretDeletedMsg{id: id}
	}
}

// restoreSecretCmd returns a Cmd that removes a soft-deletion mark.
func restoreSecretCmd(id string) tea.Cmd {
	return func() tea.Msg {
		// TODO(phase-8): call vault.Manager.RestaurarSegredo(id)
		return secretRestoredMsg{id: id}
	}
}

// reorderSecretCmd returns a Cmd that repositions a secret within its folder.
func reorderSecretCmd(id string, newIndex int) tea.Cmd {
	return func() tea.Msg {
		// TODO(phase-8): call vault.Manager.ReposicionarSegredo(id, newIndex)
		return secretReorderedMsg{}
	}
}
