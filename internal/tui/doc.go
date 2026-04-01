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
package tui
