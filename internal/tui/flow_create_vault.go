package tui

import tea "charm.land/bubbletea/v2"

// createVaultDescriptor self-describes the "create vault" flow.
// The flow starts when the user presses "n" from workAreaPreVault.
// Phase 5: stub — no real modal orchestration yet. Implemented in Phase 6.
type createVaultDescriptor struct{}

func (d createVaultDescriptor) Key() string   { return "n" }
func (d createVaultDescriptor) Label() string { return "New vault" }

// IsApplicable: create-vault is only available when no vault is currently open.
func (d createVaultDescriptor) IsApplicable(ctx FlowContext) bool {
	return !ctx.VaultOpen
}

// New creates a fresh createVaultFlow from the current context.
func (d createVaultDescriptor) New(ctx FlowContext) flowHandler {
	return &createVaultFlow{}
}

// createVaultFlow orchestrates the multi-step new-vault modal sequence.
// Phase 5: stub — real implementation in Phase 6.
type createVaultFlow struct{}

// Update receives messages forwarded by rootModel during flow execution.
// Phase 5 stub: does nothing.
func (f *createVaultFlow) Update(msg tea.Msg) tea.Cmd {
	// TODO(phase-6): push folder-picker modal, then password-create modal, then emit vaultCreatedMsg
	return nil
}
