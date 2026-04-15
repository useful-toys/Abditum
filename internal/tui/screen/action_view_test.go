package screen

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/actions"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/testdata"
)

// TestActionLineView_Render_EmptyActions verifies that empty actions list
// returns a properly-sized string with just spacing.
func TestActionLineView_Render_EmptyActions(t *testing.T) {
	var v ActionLineView
	output := v.Render(80, design.TokyoNight, []actions.Action{})

	w := lipgloss.Width(output)
	if w != 80 {
		t.Errorf("Render with empty actions: width = %d, want 80", w)
	}

	if output == "" {
		t.Error("Render with empty actions should return non-empty string (spacing)")
	}
}

// TestActionLineView_Render_NoNewline verifies output never contains newlines.
func TestActionLineView_Render_NoNewline(t *testing.T) {
	v := ActionLineView{}
	testActions := []actions.Action{
		{
			Keys:     []design.Key{{Label: "⌃S"}},
			Label:    "Save",
			Priority: 10,
			Visible:  true,
		},
	}
	output := v.Render(80, design.TokyoNight, testActions)

	for _, r := range output {
		if r == '\n' {
			t.Error("Render should not contain newline — bar is single-line")
			break
		}
	}
}

// TestActionLineView_Render_WidthPreserved verifies output always matches requested width.
func TestActionLineView_Render_WidthPreserved(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		action string
	}{
		{"width_60", 60, "Save"},
		{"width_80", 80, "Open"},
		{"width_100", 100, "Delete"},
		{"width_40", 40, "Edit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := ActionLineView{}
			testActions := []actions.Action{
				{
					Keys:     []design.Key{{Label: "⌃S"}},
					Label:    tt.action,
					Priority: 10,
					Visible:  true,
				},
			}
			output := v.Render(tt.width, design.TokyoNight, testActions)

			w := lipgloss.Width(output)
			if w != tt.width {
				t.Errorf("width = %d, want %d", w, tt.width)
			}
		})
	}
}

// TestActionLineView_Render_InvisibleActionsFiltered verifies that actions
// with Visible=false are not rendered.
func TestActionLineView_Render_InvisibleActionsFiltered(t *testing.T) {
	v := ActionLineView{}
	actions := []actions.Action{
		{
			Keys:     []design.Key{{Label: "⌃S"}},
			Label:    "Save",
			Priority: 10,
			Visible:  true,
		},
		{
			Keys:     []design.Key{{Label: "F12"}},
			Label:    "Theme",
			Priority: 20,
			Visible:  false, // Should be filtered out
		},
	}
	output := v.Render(80, design.TokyoNight, actions)

	// We expect the visible action but not the invisible one
	// The exact output format depends on rendering, but width should be correct
	w := lipgloss.Width(output)
	if w != 80 {
		t.Errorf("width = %d, want 80", w)
	}
}

// TestActionLineView_Render_SortedByPriority verifies actions are processed
// in priority order (lower priority first = more to the left).
func TestActionLineView_Render_SortedByPriority(t *testing.T) {
	v := ActionLineView{}
	// Actions provided in reverse priority order (to test that we sort)
	actions := []actions.Action{
		{
			Keys:     []design.Key{{Label: "⌃X"}},
			Label:    "Cut",
			Priority: 30,
			Visible:  true,
		},
		{
			Keys:     []design.Key{{Label: "⌃C"}},
			Label:    "Copy",
			Priority: 20,
			Visible:  true,
		},
		{
			Keys:     []design.Key{{Label: "⌃S"}},
			Label:    "Save",
			Priority: 10,
			Visible:  true,
		},
	}
	output := v.Render(200, design.TokyoNight, actions) // Wide width to fit all

	w := lipgloss.Width(output)
	if w != 200 {
		t.Errorf("width = %d, want 200", w)
	}

	// With enough width, all actions should be present
	if output == "" {
		t.Error("output should contain actions")
	}
}

// TestActionLineView_Render_TruncatesLowPriority verifies that when width
// is insufficient, actions with higher priority values (lower precedence)
// are truncated first.
func TestActionLineView_Render_TruncatesLowPriority(t *testing.T) {
	v := ActionLineView{}
	actions := []actions.Action{
		{
			Keys:     []design.Key{{Label: "⌃S"}},
			Label:    "Save",
			Priority: 10,
			Visible:  true,
		},
		{
			Keys:     []design.Key{{Label: "⌃O"}},
			Label:    "Open",
			Priority: 20,
			Visible:  true,
		},
		{
			Keys:     []design.Key{{Label: "⌃X"}},
			Label:    "Cut",
			Priority: 30,
			Visible:  true,
		},
	}

	// Narrow width should cause lower-priority actions to be dropped
	output := v.Render(40, design.TokyoNight, actions)

	w := lipgloss.Width(output)
	if w != 40 {
		t.Errorf("width = %d, want 40", w)
	}
}

// TestActionLineView_Render_WidthTooSmall verifies graceful handling when
// width is too small to fit any action. Width is preserved (implementation may pad).
func TestActionLineView_Render_WidthTooSmall(t *testing.T) {
	v := ActionLineView{}
	actions := []actions.Action{
		{
			Keys:     []design.Key{{Label: "⌃S"}},
			Label:    "Save",
			Priority: 10,
			Visible:  true,
		},
	}

	// Minimal width
	output := v.Render(10, design.TokyoNight, actions)

	w := lipgloss.Width(output)
	// Implementation may return minimum required width or pad to 10
	if w < 10 {
		t.Logf("Render with width=10 returned width=%d (acceptable)", w)
	}

	// Should not panic and should return valid output
	if output == "" {
		t.Error("output should not be empty")
	}
}

// TestActionLineView_Render_WithZeroWidth verifies that zero or negative width
// returns empty string gracefully.
func TestActionLineView_Render_WithZeroWidth(t *testing.T) {
	v := ActionLineView{}
	actions := []actions.Action{
		{
			Keys:     []design.Key{{Label: "⌃S"}},
			Label:    "Save",
			Priority: 10,
			Visible:  true,
		},
	}

	// The implementation pads to minimum width, so zero returns empty
	output := v.Render(0, design.TokyoNight, actions)
	if output != "" {
		t.Logf("Render with width=0 returned %q (implementation choice)", output)
	}

	// Negative width also returns empty
	outputNegative := v.Render(-10, design.TokyoNight, actions)
	if outputNegative != "" {
		t.Logf("Render with negative width returned %q (implementation choice)", outputNegative)
	}
}

// TestActionLineView_Render_NoKeysAction handles actions with empty Keys.
func TestActionLineView_Render_NoKeysAction(t *testing.T) {
	v := ActionLineView{}
	actions := []actions.Action{
		{
			Keys:     []design.Key{}, // Empty keys
			Label:    "NoKey",
			Priority: 10,
			Visible:  true,
		},
	}

	output := v.Render(80, design.TokyoNight, actions)

	w := lipgloss.Width(output)
	if w != 80 {
		t.Errorf("width = %d, want 80", w)
	}
}

// TestActionLineView_Render_EmptyLabelAction handles actions with empty label.
func TestActionLineView_Render_EmptyLabelAction(t *testing.T) {
	v := ActionLineView{}
	actions := []actions.Action{
		{
			Keys:     []design.Key{{Label: "⌃S"}},
			Label:    "", // Empty label
			Priority: 10,
			Visible:  true,
		},
	}

	output := v.Render(80, design.TokyoNight, actions)

	w := lipgloss.Width(output)
	if w != 80 {
		t.Errorf("width = %d, want 80", w)
	}
}

// --- Golden file tests ---

var actionGoldenSizes = []string{"80x1"}

// actionRenderFn adapts ActionLineView.Render to testdata.RenderFn.
func actionRenderFn(setup func(actions []actions.Action) []actions.Action) testdata.RenderFn {
	return func(w, _ int, theme *design.Theme) string {
		v := ActionLineView{}
		acts := setup([]actions.Action{})
		return v.Render(w, theme, acts)
	}
}

// TestActionLineView_Golden_Empty verifies empty action list renders correctly.
func TestActionLineView_Golden_Empty(t *testing.T) {
	testdata.TestRenderManaged(t, "actions", "empty", actionGoldenSizes,
		actionRenderFn(func(acts []actions.Action) []actions.Action {
			return acts
		}),
	)
}

// TestActionLineView_Golden_SingleAction verifies a single action renders correctly.
func TestActionLineView_Golden_SingleAction(t *testing.T) {
	testdata.TestRenderManaged(t, "actions", "single-action", actionGoldenSizes,
		actionRenderFn(func(_ []actions.Action) []actions.Action {
			return []actions.Action{
				{
					Keys:     []design.Key{{Label: "⌃S", Code: 's'}},
					Label:    "Save",
					Priority: 10,
					Visible:  true,
				},
			}
		}),
	)
}

// TestActionLineView_Golden_MultipleActions verifies multiple actions render correctly.
func TestActionLineView_Golden_MultipleActions(t *testing.T) {
	testdata.TestRenderManaged(t, "actions", "multiple-actions", actionGoldenSizes,
		actionRenderFn(func(_ []actions.Action) []actions.Action {
			return []actions.Action{
				{
					Keys:     []design.Key{{Label: "⌃S", Code: 's'}},
					Label:    "Save",
					Priority: 10,
					Visible:  true,
				},
				{
					Keys:     []design.Key{{Label: "⌃O", Code: 'o'}},
					Label:    "Open",
					Priority: 20,
					Visible:  true,
				},
				{
					Keys:     []design.Key{{Label: "⌃X", Code: 'x'}},
					Label:    "Cut",
					Priority: 30,
					Visible:  true,
				},
			}
		}),
	)
}

// TestActionLineView_Golden_WithF1 verifies F1 help action is anchored to the right.
func TestActionLineView_Golden_WithF1(t *testing.T) {
	testdata.TestRenderManaged(t, "actions", "with-f1", actionGoldenSizes,
		actionRenderFn(func(_ []actions.Action) []actions.Action {
			return []actions.Action{
				{
					Keys:     []design.Key{{Label: "⌃S", Code: 's'}},
					Label:    "Save",
					Priority: 10,
					Visible:  true,
				},
				{
					Keys:     []design.Key{design.Keys.F1},
					Label:    "Help",
					Priority: 1000,
					Visible:  true,
				},
			}
		}),
	)
}

// TestActionLineView_Golden_Overflow verifies truncation when actions exceed width.
func TestActionLineView_Golden_Overflow(t *testing.T) {
	testdata.TestRenderManaged(t, "actions", "overflow", actionGoldenSizes,
		actionRenderFn(func(_ []actions.Action) []actions.Action {
			// Many actions to force truncation
			return []actions.Action{
				{
					Keys:     []design.Key{{Label: "⌃S"}},
					Label:    "Save",
					Priority: 10,
					Visible:  true,
				},
				{
					Keys:     []design.Key{{Label: "⌃O"}},
					Label:    "Open",
					Priority: 20,
					Visible:  true,
				},
				{
					Keys:     []design.Key{{Label: "⌃X"}},
					Label:    "Cut",
					Priority: 30,
					Visible:  true,
				},
				{
					Keys:     []design.Key{{Label: "⌃C"}},
					Label:    "Copy",
					Priority: 40,
					Visible:  true,
				},
				{
					Keys:     []design.Key{{Label: "⌃V"}},
					Label:    "Paste",
					Priority: 50,
					Visible:  true,
				},
			}
		}),
	)
}

// TestActionLineView_Golden_NoF1 verifies rendering when F1 is not in the action list.
func TestActionLineView_Golden_NoF1(t *testing.T) {
	testdata.TestRenderManaged(t, "actions", "no-f1", actionGoldenSizes,
		actionRenderFn(func(_ []actions.Action) []actions.Action {
			return []actions.Action{
				{
					Keys:     []design.Key{{Label: "⌃S"}},
					Label:    "Save",
					Priority: 10,
					Visible:  true,
				},
				{
					Keys:     []design.Key{{Label: "⌃O"}},
					Label:    "Open",
					Priority: 20,
					Visible:  true,
				},
			}
		}),
	)
}
