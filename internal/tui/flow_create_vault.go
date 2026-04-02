package tui

import tea "charm.land/bubbletea/v2"

// createVaultFlow orchestrates the multi-step vault-creation modal sequence.
// Initiated via Action.Handler that returns startFlowMsg{flow: newCreateVaultFlow()}.
// Phase 5.1: stub - real implementation in Phase 6.
type createVaultFlow struct{}

// Init satisfies flowHandler (D-04). Phase 6 will push the folder-picker modal here.
func (f *createVaultFlow) Init() tea.Cmd {
	return nil
}

// Update receives messages forwarded by rootModel during flow execution.
func (f *createVaultFlow) Update(msg tea.Msg) tea.Cmd {
	// TODO(phase-6): process filePickerResult -> push PasswordCreate modal
	return nil
}
