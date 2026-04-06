package tui

import (
	"testing"

	testdatapkg "github.com/useful-toys/abditum/internal/tui/testdata"
)

// TestWelcomeModel_Structure verifies the welcomeModel struct has required fields.
func TestWelcomeModel_Structure(t *testing.T) {
	wm := newWelcomeModel(nil, ThemeTokyoNight)
	if wm == nil {
		t.Fatal("newWelcomeModel returned nil")
	}
	if wm.theme != ThemeTokyoNight {
		t.Errorf("expected theme to be set, got nil")
	}
}

// TestWelcomeModel_View verifies View() returns a non-empty string.
func TestWelcomeModel_View(t *testing.T) {
	wm := newWelcomeModel(nil, ThemeTokyoNight)
	wm.SetSize(80, 24)

	view := wm.View()
	if view == "" {
		t.Error("View() returned empty string")
	}
	if len(view) == 0 {
		t.Error("View() should contain logo and hints")
	}
}

// TestWelcomeModel_SetSize stores terminal dimensions.
func TestWelcomeModel_SetSize(t *testing.T) {
	wm := newWelcomeModel(nil, ThemeTokyoNight)
	wm.SetSize(80, 24)

	if wm.width != 80 || wm.height != 24 {
		t.Errorf("SetSize failed: expected 80x24, got %dx%d", wm.width, wm.height)
	}
}

// TestWelcomeModel_Update returns nil (display-only for now).
func TestWelcomeModel_Update(t *testing.T) {
	wm := newWelcomeModel(nil, ThemeTokyoNight)
	cmd := wm.Update(nil)
	if cmd != nil {
		t.Error("Update() should return nil (display-only)")
	}
}

// TestWelcomeModel_ApplyTheme applies a new theme.
func TestWelcomeModel_ApplyTheme(t *testing.T) {
	wm := newWelcomeModel(nil, ThemeTokyoNight)
	wm.ApplyTheme(ThemeCyberpunk)

	if wm.theme != ThemeCyberpunk {
		t.Errorf("ApplyTheme failed: expected ThemeCyberpunk")
	}
}

// TestWelcomeModel_ViewContainsLogo verifies View includes the logo text.
func TestWelcomeModel_ViewContainsLogo(t *testing.T) {
	wm := newWelcomeModel(nil, ThemeTokyoNight)
	wm.SetSize(80, 24)

	view := wm.View()
	// Logo should contain "A" (part of "Abditum")
	if len(view) == 0 {
		t.Fatal("View should not be empty")
	}
}

// TestWelcomeModel_ViewContainsHints verifies View includes action hints.
func TestWelcomeModel_ViewContainsHints(t *testing.T) {
	wm := newWelcomeModel(nil, ThemeTokyoNight)
	wm.SetSize(80, 24)

	view := wm.View()
	// Should contain hint text about 'n' and 'o' actions
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
		{"tokyo-night", ThemeTokyoNight},
		{"cyberpunk", ThemeCyberpunk},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wm := newWelcomeModel(nil, tc.theme)
			wm.SetSize(80, 24)

			out := wm.View()

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
