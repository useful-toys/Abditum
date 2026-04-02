package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ActionScope controls when an action is eligible for dispatch.
type ActionScope int

const (
	ScopeLocal  ActionScope = iota // eligible only when no flow or modal is active
	ScopeGlobal                    // always eligible - even during flows and modals
)

// Action is the central unit of keyboard interaction (D-05).
// Keys[0] is shown in the command bar. All keys in Keys trigger the action.
// Enabled and Handler are closures over the registering child.
type Action struct {
	Keys        []string
	Label       string
	Description string
	Group       string
	Scope       ActionScope
	Enabled     func() bool
	Handler     func() tea.Cmd
}

// ActionManager is the owner-tracked, dispatch-capable action registry (D-06).
type ActionManager struct {
	owners      []any
	byOwner     map[any][]Action
	activeOwner any
}

// NewActionManager creates a new, empty ActionManager.
func NewActionManager() *ActionManager {
	return &ActionManager{byOwner: make(map[any][]Action)}
}

// Register adds actions for the given owner.
func (a *ActionManager) Register(owner any, actions ...Action) {
	if _, exists := a.byOwner[owner]; !exists {
		a.owners = append(a.owners, owner)
	}
	a.byOwner[owner] = append(a.byOwner[owner], actions...)
}

// ClearOwned removes all actions registered for owner.
func (a *ActionManager) ClearOwned(owner any) {
	delete(a.byOwner, owner)
	filtered := a.owners[:0]
	for _, o := range a.owners {
		if o != owner {
			filtered = append(filtered, o)
		}
	}
	a.owners = filtered
	if a.activeOwner == owner {
		a.activeOwner = nil
	}
}

// SetActiveOwner prioritizes actions from the given owner during Dispatch.
func (a *ActionManager) SetActiveOwner(owner any) {
	a.activeOwner = owner
}

// Dispatch finds the first eligible action matching key and executes its Handler.
func (a *ActionManager) Dispatch(key string, inFlowOrModal bool) tea.Cmd {
	var ordered []any
	if a.activeOwner != nil {
		ordered = append(ordered, a.activeOwner)
	}
	for _, o := range a.owners {
		if o != a.activeOwner {
			ordered = append(ordered, o)
		}
	}

	for _, owner := range ordered {
		for _, act := range a.byOwner[owner] {
			if act.Scope == ScopeLocal && inFlowOrModal {
				continue
			}
			if act.Enabled != nil && !act.Enabled() {
				continue
			}
			for _, k := range act.Keys {
				if k == key {
					return act.Handler()
				}
			}
		}
	}
	return nil
}

// Visible returns actions where Enabled() is true, for the command bar.
func (a *ActionManager) Visible() []Action {
	var result []Action
	seen := make(map[string]bool)

	var ordered []any
	if a.activeOwner != nil {
		ordered = append(ordered, a.activeOwner)
	}
	for _, o := range a.owners {
		if o != a.activeOwner {
			ordered = append(ordered, o)
		}
	}

	for _, owner := range ordered {
		for _, act := range a.byOwner[owner] {
			if act.Enabled != nil && !act.Enabled() {
				continue
			}
			key := strings.Join(act.Keys, ",")
			if !seen[key] {
				seen[key] = true
				result = append(result, act)
			}
		}
	}
	return result
}

// All returns all registered actions in registration order, for the help overlay.
func (a *ActionManager) All() []Action {
	var result []Action
	for _, owner := range a.owners {
		result = append(result, a.byOwner[owner]...)
	}
	return result
}

// RenderCommandBar renders the command bar from currently visible actions.
func (a *ActionManager) RenderCommandBar(width int) string {
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	sepStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	var parts []string
	for _, act := range a.Visible() {
		if len(act.Keys) == 0 {
			continue
		}
		parts = append(parts, keyStyle.Render(act.Keys[0])+" "+labelStyle.Render(act.Label))
	}
	if len(parts) == 0 {
		return ""
	}
	return "  " + strings.Join(parts, sepStyle.Render("  |  ")+"  ")
}
