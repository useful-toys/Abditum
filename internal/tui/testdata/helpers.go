package testdata

import (
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/useful-toys/abditum/internal/tui/design"
)

// update is a flag to update golden files instead of comparing against them.
var update = flag.Bool("update-golden", false, "update golden files")

// RenderFn is a function that renders a TUI component for golden file testing.
// It takes width and height (in characters) and a theme, and returns the rendered
// output as a string (typically with ANSI escape sequences).
type RenderFn func(w, h int, theme *design.Theme) string

// goldenPath constructs the filesystem path to a golden file.
// Parameters: component (e.g., "tree"), variant (e.g., "expanded"), size (e.g., "30x20" or "80"),
// ext (e.g., "txt" or "json").
// Example: goldenPath("messages", "success", "80", "txt") → "testdata/golden/messages-success-80.golden.txt"
func goldenPath(component, variant, size, ext string) string {
	filename := component + "-" + variant + "-" + size + ".golden." + ext
	return filepath.Join("testdata", "golden", filename)
}

// ansiEscapeRe captura todas as sequências de escape ANSI SGR.
var ansiEscapeRe = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// ansiToCleanText removes all ANSI escape sequences from a string,
// returning only the visible text content. Safe for multi-byte UTF-8.
func ansiToCleanText(s string) string {
	return ansiEscapeRe.ReplaceAllString(s, "")
}

// checkOrUpdateGolden compares the rendered output against a golden file.
// If the -update-golden flag is set, it writes/updates the golden file.
// If the flag is not set, it compares and fails the test if they differ.
func checkOrUpdateGolden(t *testing.T, path, got string) {
	t.Helper()

	if *update {
		// Ensure directory exists
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("failed to create golden directory %s: %v", dir, err)
		}
		// Write/update the golden file
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatalf("failed to write golden file %s: %v", path, err)
		}
		return
	}

	// Compare against golden file
	expected, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("golden file not found: %s (run with -update-golden to create it)", path)
		}
		t.Fatalf("failed to read golden file %s: %v", path, err)
	}

	if string(expected) != got {
		t.Errorf("output does not match golden file %s", path)
		// Optionally output diff information for debugging
		cleanExpected := ansiToCleanText(string(expected))
		cleanGot := ansiToCleanText(got)
		if cleanExpected != cleanGot {
			t.Logf("Clean text differs:\nExpected:\n%s\n\nGot:\n%s", cleanExpected, cleanGot)
		}
	}
}

// parseSize parses a size string into width and height.
// Accepts formats like "30" (width=30, height=30) or "80x24" (width=80, height=24).
func parseSize(size string) (int, int, error) {
	if strings.Contains(size, "x") || strings.Contains(size, "X") {
		parts := strings.FieldsFunc(size, func(r rune) bool { return r == 'x' || r == 'X' })
		if len(parts) != 2 {
			return 0, 0, invalid("invalid size format: " + size)
		}
		w, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		h, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err1 != nil || err2 != nil {
			return 0, 0, invalid("invalid size values: " + size)
		}
		return w, h, nil
	}

	// Single number: square dimension
	dim, err := strconv.Atoi(strings.TrimSpace(size))
	if err != nil {
		return 0, 0, invalid("invalid size: " + size)
	}
	return dim, dim, nil
}

// invalid is a helper to create an error without exposing implementation details.
func invalid(msg string) error {
	return parseError{msg}
}

type parseError struct {
	message string
}

func (e parseError) Error() string {
	return e.message
}

// TestRenderManaged runs golden file tests for a component with multiple size variants.
// It uses the managed golden file naming convention and the -update-golden flag.
// Parameters:
//   - t: testing.T instance
//   - component: component name (e.g., "vault_tree")
//   - variant: variant name (e.g., "expanded")
//   - sizes: slice of size strings (e.g., []string{"30", "80x24"})
//   - render: RenderFn that renders the component for a given size and theme
func TestRenderManaged(t *testing.T, component, variant string, sizes []string, render RenderFn) {
	t.Helper()

	// Use the default theme (TokyoNight)
	theme := design.TokyoNight

	for _, sizeStr := range sizes {
		t.Run(sizeStr, func(t *testing.T) {
			w, h, err := parseSize(sizeStr)
			if err != nil {
				t.Fatalf("failed to parse size %q: %v", sizeStr, err)
			}

			// Render the component
			output := render(w, h, theme)

			// Check/update clean text golden file (.txt.golden) — ANSI stripped
			txtPath := goldenPath(component, variant, sizeStr, "txt")
			checkOrUpdateGolden(t, txtPath, ansiToCleanText(output))

			// Check/update style transitions golden file (.json)
			transitions := ansiToStyleChanges(output)
			jsonBytes, err := marshalStyleChanges(transitions)
			if err != nil {
				t.Fatalf("failed to marshal style transitions: %v", err)
			}
			jsonPath := goldenPath(component, variant, sizeStr, "json")
			checkOrUpdateGolden(t, jsonPath, string(jsonBytes))
		})
	}
}

// TestRenderManual is a simpler variant for manually-specified golden files.
// It does not use the -update-golden flag; instead, it always compares against
// the provided golden output string. This is useful for quick validation tests
// that don't need persistent golden files.
// Parameters:
//   - t: testing.T instance
//   - component: component name (for logging)
//   - variant: variant name (for logging)
//   - size: size string (e.g., "30x20")
//   - output: the expected (golden) output string
func TestRenderManual(t *testing.T, component, variant, size, output string) {
	t.Helper()

	// This function is intended to be called within a test context
	// with an already-rendered output. It serves as a validation helper
	// that ensures the output matches a manually-specified golden string.
	// For now, this is a placeholder for documentation and future expansion.

	// In practice, you would call this like:
	// rendered := renderMyComponent(30, 24, theme)
	// testdata.TestRenderManual(t, "my_component", "normal", "30x24", rendered)
	// which would then compare rendered against a manually-created golden string.
}
