// Package tui implements the terminal user interface for Abditum using the Bubble Tea v2
// framework (charm.land/bubbletea/v2). It follows the Elm architecture: a single rootModel
// implements tea.Model and orchestrates all child models, modal overlays, flow handlers,
// and the persistent frame layout.
//
// Architecture overview:
//
//   - rootModel   — the only tea.Model; owns the work area, modal stack, active flow, and shared services
//   - childModel  — interface implemented by all child TUI models (preVaultModel, vaultTreeModel, etc.)
//   - flowHandler — interface for multi-step modal orchestration flows (open vault, create vault, etc.)
//   - FlowRegistry — repository of globally registered flowDescriptors
//   - ActionManager — centralized registry of available keyboard actions for the command bar and help overlay
//   - MessageManager — centralized API for setting the message bar content
//
// Message dispatch rules (D-06):
//  1. Global shortcuts (ctrl+Q, ?) — intercepted by rootModel first
//  2. Active flow — if any, receives all input events
//  3. Topmost modal — receives input when modal stack is non-empty
//  4. Base child — receives input when no flow or modal is active
//
// Domain messages (vault events, tick) are broadcast to all live models via liveModels().
//
// SetSize-Before-View Contract
// =============================
// The rootModel enforces a strict ordering: SetSize(w, h) ALWAYS precedes View().
// This eliminates the need for defensive checks inside View() implementations.
//
// All childModels and modals implement a panic guard as the first line of View():
//
//	func (m *myModel) View() string {
//		if m.width == 0 || m.height == 0 {
//			panic(fmt.Sprintf("myModel.View() called without SetSize: w=%d h=%d", m.width, m.height))
//		}
//		// ... rendering logic (no size checks needed)
//	}
//
// If View() panics, it's a bug in rootModel's render flow, not in the component.
// The panic message will immediately pinpoint the violation.
//
// This approach:
// - Eliminates ~30 defensive checks spread across the codebase
// - Makes the contract explicit and testable
// - Provides clear error messages when violations occur
// - Simplifies component implementation
package tui
