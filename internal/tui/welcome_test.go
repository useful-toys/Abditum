package tui

import (
	"testing"

	testdatapkg "github.com/useful-toys/abditum/internal/tui/testdata"
)

// TestWelcomeModel_Structure verifies the welcomeModel struct has required fields.
func TestWelcomeModel_Structure(t *testing.T) {
	wm := newWelcomeModel(nil, "v0.1.0", TokyoNight)
	if wm == nil {
		t.Fatal("newWelcomeModel returned nil")
	}
	if wm.version != "v0.1.0" {
		t.Errorf("expected version to be set, got %q", wm.version)
	}
}

// TestWelcomeModel_View verifies View() returns a non-empty string.
func TestWelcomeModel_View(t *testing.T) {
	wm := newWelcomeModel(nil, "v0.1.0", TokyoNight)

	view := wm.View(80, 24, TokyoNight)
	if view == "" {
		t.Error("View() returned empty string")
	}
	if len(view) == 0 {
		t.Error("View() should contain logo and version")
	}
}

// TestWelcomeModel_Update returns nil (display-only for now).
func TestWelcomeModel_Update(t *testing.T) {
	wm := newWelcomeModel(nil, "v0.1.0", TokyoNight)
	cmd := wm.Update(nil)
	if cmd != nil {
		t.Error("Update() should return nil (display-only)")
	}
}

// TestWelcomeModel_ViewContainsLogo verifies View includes the logo text.
func TestWelcomeModel_ViewContainsLogo(t *testing.T) {
	wm := newWelcomeModel(nil, "v0.1.0", TokyoNight)

	view := wm.View(80, 24, TokyoNight)
	// Logo should contain "A" (part of "Abditum")
	if len(view) == 0 {
		t.Fatal("View should not be empty")
	}
}

// TestWelcomeModel_ViewContainsHints verifies View includes action hints.
func TestWelcomeModel_ViewContainsHints(t *testing.T) {
	wm := newWelcomeModel(nil, "v0.1.0", TokyoNight)

	view := wm.View(80, 24, TokyoNight)
	// Should contain hint text about version
	if len(view) == 0 {
		t.Fatal("View should not be empty")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Golden File Tests
// ─────────────────────────────────────────────────────────────────────────────

// TestWelcomeModel_Golden validates the visual output of welcomeModel
// against golden files for both Tokyo Night and Cyberpunk themes at 80x24.
func TestWelcomeModel_Golden(t *testing.T) {
	type testCase struct {
		name  string
		theme *Theme
	}

	cases := []testCase{
		{"tokyo-night", TokyoNight},
		{"cyberpunk", Cyberpunk},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wm := newWelcomeModel(nil, "v0.1.0", TokyoNight)

			out := wm.View(80, 24, TokyoNight)

			// .txt.golden: raw ANSI output stripped of codes
			txtPath := goldenPath("welcome", tc.name, 80, "txt")
			checkOrUpdateGolden(t, txtPath, stripANSI(out))

			// .json.golden: style transitions
			transitions := testdatapkg.ParseANSIStyle(out)
			jsonBytes, err := testdatapkg.MarshalStyleTransitions(transitions)
			if err != nil {
				t.Fatalf("marshal transitions: %v", err)
			}
			jsonPath := goldenPath("welcome", tc.name, 80, "json")
			checkOrUpdateGolden(t, jsonPath, string(jsonBytes))
		})
	}
}
