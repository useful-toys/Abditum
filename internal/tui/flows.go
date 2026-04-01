package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/vault"
)

// FlowContext captures the complete navigation/vault state at the moment of
// flow dispatch. The active child fills navigation fields via Context();
// rootModel adds vault-level fields (VaultOpen, VaultDirty).
type FlowContext struct {
	// Filled by rootModel from vault.Manager
	VaultOpen  bool
	VaultDirty bool
	// Filled by the active child's Context() method
	FocusedFolder   *vault.Pasta
	FocusedSecret   *vault.Segredo
	SecretOpen      bool
	FocusedField    *vault.CampoSegredo
	FocusedTemplate *vault.ModeloSegredo
	Mode            int // child-defined: e.g., view vs edit, left vs right pane focus
}

// flowDescriptor self-describes a flow's trigger key, label, applicability,
// and factory. Implementations live in individual flow files.
type flowDescriptor interface {
	Key() string                   // keyboard shortcut that triggers this flow
	Label() string                 // display label for ActionManager / help
	IsApplicable(FlowContext) bool // pure function — can this flow start given context?
	New(FlowContext) flowHandler   // factory — creates a fresh handler from current context
}

// flowHandler encapsulates multi-step modal orchestration for a single flow.
// It has no View — flows push modals onto rootModel's stack to show UI.
type flowHandler interface {
	Update(tea.Msg) tea.Cmd // mutates flow state, returns Cmds (push modals, async work, domain msgs)
}

// FlowRegistry is the repository of all globally registered flows.
// rootModel owns one FlowRegistry, populated at startup.
type FlowRegistry struct {
	descriptors []flowDescriptor
}

// Register adds a flow descriptor to the registry.
func (r *FlowRegistry) Register(d flowDescriptor) {
	r.descriptors = append(r.descriptors, d)
}

// ForKey returns the first applicable flow descriptor matching the given key
// and context, or nil if none matches.
func (r *FlowRegistry) ForKey(key string, ctx FlowContext) flowDescriptor {
	for _, d := range r.descriptors {
		if d.Key() == key && d.IsApplicable(ctx) {
			return d
		}
	}
	return nil
}

// Applicable returns all flow descriptors applicable in the given context.
// Used to populate ActionManager for the command bar.
func (r *FlowRegistry) Applicable(ctx FlowContext) []flowDescriptor {
	var result []flowDescriptor
	for _, d := range r.descriptors {
		if d.IsApplicable(ctx) {
			result = append(result, d)
		}
	}
	return result
}

// chainFlowMsg is emitted by a completing flow to request starting another
// flow immediately after state transition. rootModel handles it by rebuilding
// FlowContext and dispatching through FlowRegistry.ForKey.
type chainFlowMsg struct {
	key string
}

// childModel is the interface all child TUI models must satisfy.
// It deliberately does NOT implement tea.Model — only rootModel does.
// View() returns string (not tea.View); only rootModel.View() returns tea.View.
// Update uses pointer receivers and mutates in place — no self-replacement.
type childModel interface {
	Update(tea.Msg) tea.Cmd       // mutates in place, returns only Cmd (no self-replacement)
	View() string                 // returns string, NOT tea.View
	SetSize(w, h int)             // receives allocated size from rootModel compositor
	Context() FlowContext         // exposes navigation/selection state for flow dispatch
	ChildFlows() []flowDescriptor // child-specific flows not in global registry (escape hatch)
}
