package tui

import (
	"fmt"
	"sort"
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
	Group       int // grouping in Help modal — groups shown in ascending order
	Scope       ActionScope
	Priority    int  // higher = further left in command bar; governs truncation order
	HideFromBar bool // hidden from command bar; listed in Help modal
	Enabled     func() bool
	Handler     func() tea.Cmd
}

// ActionManager is the owner-tracked, dispatch-capable action registry (D-06).
type ActionManager struct {
	owners      []any
	byOwner     map[any][]Action
	activeOwner any
	groupLabels map[int]string
}

// NewActionManager creates a new, empty ActionManager.
func NewActionManager() *ActionManager {
	return &ActionManager{
		byOwner:     make(map[any][]Action),
		groupLabels: make(map[int]string),
	}
}

// RegisterGroupLabel associates a display name with a group int for the Help modal.
func (a *ActionManager) RegisterGroupLabel(group int, label string) {
	a.groupLabels[group] = label
}

// GroupLabel returns the display name for a group, or a fallback string if unregistered.
func (a *ActionManager) GroupLabel(group int) string {
	if label, ok := a.groupLabels[group]; ok {
		return label
	}
	return fmt.Sprintf("Group %d", group)
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

// Visible returns actions where Enabled() is true and HideFromBar is false, sorted by Priority descending.
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
			if act.HideFromBar {
				continue
			}
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

	// Sort by Priority descending (higher priority = further left in bar)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Priority > result[j].Priority
	})
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

// RenderCommandBar renders the command bar from a pre-computed slice of visible actions.
// Tokens per DS spec: key = accent.primary bold, label = text.primary, separator = text.secondary.
// The command bar uses the application's default background (surface.base) — no explicit bg set.
// F1 action is right-anchored; all other actions are left-padded with 2 spaces.
// Actions are sorted by Priority descending before rendering — callers may pass unsorted slices.
// Callers should pass am.Visible() to obtain the visible action slice.
func RenderCommandBar(actions []Action, width int, theme *Theme) string {
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Accent.Primary)).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Primary))
	sepStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Secondary))

	// renderPart builds a single "KEY Label" token pair.
	renderPart := func(key, label string) string {
		return keyStyle.Render(key) + " " + labelStyle.Render(label)
	}

	// Sort a copy of actions by Priority descending so callers may pass unsorted slices.
	sorted := make([]Action, len(actions))
	copy(sorted, actions)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Priority > sorted[j].Priority
	})
	visible := sorted

	// Separate F1 anchor from body actions.
	var bodyActions []Action
	var anchorAction *Action
	for i := range visible {
		if len(visible[i].Keys) > 0 && strings.EqualFold(visible[i].Keys[0], "f1") {
			act := visible[i]
			anchorAction = &act
		} else {
			bodyActions = append(bodyActions, visible[i])
		}
	}

	// Build body parts (visible widths stored separately for truncation loop).
	type part struct {
		rendered string
		visW     int
	}
	var bodyParts []part
	for _, act := range bodyActions {
		if len(act.Keys) == 0 {
			continue
		}
		r := renderPart(formatKeyForHelp(act.Keys[0]), act.Label)
		bodyParts = append(bodyParts, part{rendered: r, visW: lipgloss.Width(r)})
	}

	sep := sepStyle.Render(" · ")
	sepW := lipgloss.Width(sep)

	// joinBody assembles left padding + parts + separators.
	// Returns (assembled string, visible width).
	joinBody := func(parts []part) (string, int) {
		if len(parts) == 0 {
			return "", 0
		}
		out := "  "
		w := 2
		for i, p := range parts {
			if i > 0 {
				out += sep
				w += sepW
			}
			out += p.rendered
			w += p.visW
		}
		return out, w
	}

	// No anchor: return body only.
	if anchorAction == nil || len(anchorAction.Keys) == 0 {
		body, _ := joinBody(bodyParts)
		return body
	}

	// Build anchor string.
	anchor := renderPart(formatKeyForHelp(anchorAction.Keys[0]), anchorAction.Label)
	anchorW := lipgloss.Width(anchor)
	minGap := 1

	// Truncate lowest-priority actions until anchor fits.
	// bodyParts is already sorted by Priority desc (from Visible()).
	for len(bodyParts) > 0 {
		_, bodyW := joinBody(bodyParts)
		if width-bodyW-anchorW >= minGap {
			break
		}
		bodyParts = bodyParts[:len(bodyParts)-1]
	}

	body, bodyW := joinBody(bodyParts)
	if len(bodyParts) == 0 {
		body = "  "
		bodyW = 2
	}

	gap := width - bodyW - anchorW
	if gap < minGap {
		gap = minGap
	}
	return body + strings.Repeat(" ", gap) + anchor
}
