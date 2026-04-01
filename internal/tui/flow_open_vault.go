package tui

import tea "charm.land/bubbletea/v2"

// openVaultDescriptor self-describes the "open vault" flow.
// The flow starts when the user presses "o" from workAreaPreVault.
// Phase 5: stub — no real modal orchestration yet. Implemented in Phase 6.
type openVaultDescriptor struct{}

func (d openVaultDescriptor) Key() string   { return "o" }
func (d openVaultDescriptor) Label() string { return "Open vault" }

// IsApplicable: open-vault is only available when no vault is currently open.
func (d openVaultDescriptor) IsApplicable(ctx FlowContext) bool {
	return !ctx.VaultOpen
}

// New creates a fresh openVaultFlow from the current context.
func (d openVaultDescriptor) New(ctx FlowContext) flowHandler {
	return &openVaultFlow{}
}

// openVaultFlow orchestrates the multi-step open-vault modal sequence.
// Phase 5: stub — real implementation in Phase 6.
type openVaultFlow struct{}

// Update receives messages forwarded by rootModel during flow execution.
// Phase 5 stub: does nothing and signals immediate (no-op) completion.
func (f *openVaultFlow) Update(msg tea.Msg) tea.Cmd {
	// TODO(phase-6): push file-picker modal, then password modal, then emit vaultOpenedMsg
	return nil
}
