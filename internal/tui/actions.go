package tui

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// ActionManager is the centralized registry of keyboard actions available at any moment.
// It is a shared mutable object instantiated in main.go and passed to rootModel,
// which in turn passes it to every child at construction time.
//
// Guiding analogy: just as vault.Manager is the API for vault operations,
// ActionManager is the API for defining which actions are available (D-16).
//
// Usage:
//   - Registration: each child calls Register when it becomes active
//   - Command bar: reads Visible() — a prioritized, display-width-aware subset
//   - Help modal: reads All() — the full list of registered actions
//   - rootModel: registers global shortcuts (ctrl+Q, ?) at startup
//
// ActionManager does NOT know about Bubble Tea internals — it is a plain Go struct
// with no tea.Cmd or messaging. It is queried synchronously from View() only.
type ActionManager struct {
	actions []Action
}

// Action represents a single keyboard action available in the current context.
type Action struct {
	Key         string // keyboard shortcut (e.g., "ctrl+q", "enter", "n")
	Label       string // short display label for command bar (e.g., "Quit", "New")
	Description string // longer description for help overlay
	Group       string // grouping key for help overlay (e.g., "Global", "Vault", "Secret")
}

// Register adds an action to the registry.
// Children call this when they become active to surface their shortcuts.
func (a *ActionManager) Register(action Action) {
	a.actions = append(a.actions, action)
}

// Unregister removes all actions belonging to a given group.
// Children call this when they become inactive to clean up their shortcuts.
func (a *ActionManager) Unregister(group string) {
	filtered := a.actions[:0]
	for _, action := range a.actions {
		if action.Group != group {
			filtered = append(filtered, action)
		}
	}
	a.actions = filtered
}

// All returns the complete list of registered actions, for the help overlay.
func (a *ActionManager) All() []Action {
	result := make([]Action, len(a.actions))
	copy(result, a.actions)
	return result
}

// Visible returns a subset of registered actions for the command bar.
// Phase 5 stub: returns all actions (full prioritization logic deferred to later phases).
func (a *ActionManager) Visible() []Action {
	return a.All()
}

// NewActionManager creates a new ActionManager with no registered actions.
func NewActionManager() *ActionManager {
	return &ActionManager{}
}

// Ensure textinput import is used — bubbles/v2 provides components for future phases.
// This blank assignment keeps the import active until real textinput usage is added.
var _ = textinput.New

// Ensure tea import is used.
var _ tea.Cmd
