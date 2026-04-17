package modal_test

import (
	"testing"

	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/tui/testdata"
)

// TestEvaluateStrength tests password strength scoring.
func TestEvaluateStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected modal.StrengthLevel
	}{
		// Empty and very short passwords (0 criteria)
		{
			name:     "empty password",
			password: "",
			expected: modal.StrengthWeak,
		},
		{
			name:     "single lowercase",
			password: "a",
			expected: modal.StrengthWeak,
		},
		{
			name:     "three lowercase chars",
			password: "abc",
			expected: modal.StrengthWeak,
		},

		// Single criterion (0-1 points = Weak)
		{
			name:     "uppercase only (1 criterion)",
			password: "Abc",
			expected: modal.StrengthWeak,
		},
		{
			name:     "symbol only (1 criterion)",
			password: "!@#$%",
			expected: modal.StrengthWeak,
		},

		// Two criteria (2-3 points = Fair)
		{
			name:     "11 chars + uppercase (only 1 criterion - length needs 12)",
			password: "Abcdefghijk",
			expected: modal.StrengthWeak,
		},
		{
			name:     "12 chars + digit (2 criteria: length + digit)",
			password: "abcdefghijk1",
			expected: modal.StrengthFair,
		},
		{
			name:     "uppercase + digit (2 criteria)",
			password: "Abc1",
			expected: modal.StrengthFair,
		},
		{
			name:     "uppercase + symbol (2 criteria)",
			password: "Abc!",
			expected: modal.StrengthFair,
		},
		{
			name:     "digit + symbol (2 criteria)",
			password: "abc1!",
			expected: modal.StrengthFair,
		},

		// Three criteria (2-3 points = Fair)
		{
			name:     "12 chars + uppercase + digit (3 criteria)",
			password: "Abcdefghijk1",
			expected: modal.StrengthFair,
		},
		{
			name:     "12 chars + uppercase + symbol (3 criteria)",
			password: "Abcdefghijk!",
			expected: modal.StrengthFair,
		},
		{
			name:     "uppercase + digit + symbol without length (3 criteria)",
			password: "Aa1!",
			expected: modal.StrengthFair,
		},

		// Four criteria (4 points = Strong)
		{
			name:     "all four criteria - 12 chars + upper + digit + symbol",
			password: "Abcdefghijk1!",
			expected: modal.StrengthStrong,
		},
		{
			name:     "complex password with multiple symbols",
			password: "MyP@ssw0rd123",
			expected: modal.StrengthStrong,
		},
		{
			name:     "exactly 12 chars with all criteria",
			password: "Abcdefghij1!",
			expected: modal.StrengthStrong,
		},
		{
			name:     "all symbol types valid",
			password: "Pass@123!#$%^&*()_+-=[]{}|;:,.<>?/~",
			expected: modal.StrengthStrong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := modal.EvaluateStrength([]byte(tt.password))
			if got != tt.expected {
				t.Errorf("EvaluateStrength(%q) = %v, want %v", tt.password, got, tt.expected)
			}
		})
	}
}

// TestStrengthMeter_Weak tests rendering of weak strength meter.
func TestStrengthMeter_Weak(t *testing.T) {
	// Weak passwords have 0-1 criteria (2 blocks filled)
	password := "Pass"
	testdata.TestRenderManaged(t, "strength_meter", "weak", []string{"44x1"},
		func(w, h int, theme *design.Theme) string {
			return modal.RenderStrengthMeter([]byte(password), w, theme)
		})
}

// TestStrengthMeter_Fair tests rendering of fair strength meter.
func TestStrengthMeter_Fair(t *testing.T) {
	// Fair passwords have 2-3 criteria (8 blocks filled)
	password := "Abcdefghijk1"
	testdata.TestRenderManaged(t, "strength_meter", "fair", []string{"44x1"},
		func(w, h int, theme *design.Theme) string {
			return modal.RenderStrengthMeter([]byte(password), w, theme)
		})
}

// TestStrengthMeter_Strong tests rendering of strong strength meter.
func TestStrengthMeter_Strong(t *testing.T) {
	// Strong passwords have 4 criteria (10 blocks filled)
	password := "Abcdefghijk1!"
	testdata.TestRenderManaged(t, "strength_meter", "strong", []string{"44x1"},
		func(w, h int, theme *design.Theme) string {
			return modal.RenderStrengthMeter([]byte(password), w, theme)
		})
}
