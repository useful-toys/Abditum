package tui

import "charm.land/lipgloss/v2"

// Action represents a single registered keyboard action shown in the command bar or help.
type Action struct {
	Key         string // keyboard shortcut string (e.g., "ctrl+q", "?")
	Label       string // short display label for command bar
	Description string // longer description for help overlay
	Group       string // grouping label (e.g., "Global", "Vault", "Navigation")
	Priority    int    // higher priority shown first in command bar; 0 = default
}

// ActionManager is the centralized registry of currently available actions.
// It is a shared mutable object: rootModel and all children call Register/Clear
// on it; only rootModel.View() reads from it via Visible() and All().
//
// ActionManager does NOT know about Bubble Tea internals — it holds no tea.Cmd
// or messages. It is queried synchronously from View() only.
//
// Phase 5: stub implementation — Visible() returns a flat slice, no display-width
// awareness yet. Full priority/grouping logic added in later phases.
type ActionManager struct {
	actions []Action
}

// NewActionManager creates a new, empty ActionManager.
func NewActionManager() *ActionManager {
	return &ActionManager{}
}

// Register adds an action to the registry. Duplicate keys are allowed
// (the latest registration wins display priority when groups are merged).
func (a *ActionManager) Register(action Action) {
	a.actions = append(a.actions, action)
}

// ClearGroup removes all actions belonging to the given group.
// Children call this when they deactivate, passing their group name.
func (a *ActionManager) ClearGroup(group string) {
	filtered := a.actions[:0]
	for _, act := range a.actions {
		if act.Group != group {
			filtered = append(filtered, act)
		}
	}
	a.actions = filtered
}

// Visible returns a prioritized subset of registered actions for the command bar.
// Phase 5: returns all actions sorted by Priority descending (flat list).
// Later phases: add display-width awareness and truncation.
func (a *ActionManager) Visible() []Action {
	// Return copy sorted by priority (higher first).
	result := make([]Action, len(a.actions))
	copy(result, a.actions)
	// Simple insertion sort by Priority descending (small slice, no import needed).
	for i := 1; i < len(result); i++ {
		for j := i; j > 0 && result[j].Priority > result[j-1].Priority; j-- {
			result[j], result[j-1] = result[j-1], result[j]
		}
	}
	return result
}

// All returns all registered actions, grouped. Used by the help overlay.
func (a *ActionManager) All() []Action {
	result := make([]Action, len(a.actions))
	copy(result, a.actions)
	return result
}

// RenderCommandBar renders the command bar line from visible actions.
// Placeholder style: "  key  label  |  key  label  ..."
func (a *ActionManager) RenderCommandBar(width int) string {
	actions := a.Visible()
	if len(actions) == 0 {
		return ""
	}
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11")) // yellow
	var parts []string
	for _, act := range actions {
		parts = append(parts, keyStyle.Render(act.Key)+"  "+act.Label)
	}
	bar := ""
	for i, p := range parts {
		if i > 0 {
			bar += "  │  "
		}
		bar += p
	}
	return lipgloss.NewStyle().Width(width).Render(bar)
}
