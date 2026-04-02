package tui

import tea "charm.land/bubbletea/v2"

// openVaultFlow orchestrates the multi-step open-vault modal sequence.
// Initiated via Action.Handler that returns startFlowMsg{flow: newOpenVaultFlow()}.
// Phase 5.1: stub - real implementation in Phase 6.
type openVaultFlow struct{}

// Init satisfies flowHandler (D-04). Phase 6 will push the file-picker modal here.
func (f *openVaultFlow) Init() tea.Cmd {
	return nil
}

// Update receives messages forwarded by rootModel during flow execution.
func (f *openVaultFlow) Update(msg tea.Msg) tea.Cmd {
	// TODO(phase-6): process filePickerResult -> push PasswordEntry modal
	return nil
}
