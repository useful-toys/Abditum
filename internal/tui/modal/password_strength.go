package modal

import (
	"bytes"
	"fmt"
	"unicode"

	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// StrengthLevel represents the strength classification of a password.
type StrengthLevel int

const (
	// StrengthWeak represents passwords with 0-1 criteria met.
	StrengthWeak StrengthLevel = iota
	// StrengthFair represents passwords with 2-3 criteria met.
	StrengthFair
	// StrengthStrong represents passwords with all 4 criteria met.
	StrengthStrong
)

// symbolChars contains all valid symbols for password strength evaluation.
const symbolChars = "!@#$%^&*()-_=+[]{}|;:,.<>?/~"

// strengthMeterBlocks is the total number of blocks in the strength meter.
const strengthMeterBlocks = 10

// EvaluateStrength evaluates password strength based on 4 criteria:
// 1. Length >= 12 characters (1 point)
// 2. Contains at least 1 uppercase letter (1 point)
// 3. Contains at least 1 digit (1 point)
// 4. Contains at least 1 symbol from the predefined set (1 point)
//
// Returns StrengthWeak (0-1 points), StrengthFair (2-3 points), or StrengthStrong (4 points).
func EvaluateStrength(password []byte) StrengthLevel {
	score := 0

	// Criterion 1: Length >= 12 characters
	if len(password) >= 12 {
		score++
	}

	// Criterion 2: Contains at least 1 uppercase letter
	hasUpper := false
	for _, b := range password {
		r := rune(b)
		if unicode.IsUpper(r) {
			hasUpper = true
			break
		}
	}
	if hasUpper {
		score++
	}

	// Criterion 3: Contains at least 1 digit
	hasDigit := false
	for _, b := range password {
		r := rune(b)
		if unicode.IsDigit(r) {
			hasDigit = true
			break
		}
	}
	if hasDigit {
		score++
	}

	// Criterion 4: Contains at least 1 symbol
	hasSymbol := false
	for _, b := range password {
		if bytes.ContainsRune([]byte(symbolChars), rune(b)) {
			hasSymbol = true
			break
		}
	}
	if hasSymbol {
		score++
	}

	// Classify based on score
	if score >= 4 {
		return StrengthStrong
	} else if score >= 2 {
		return StrengthFair
	}
	return StrengthWeak
}

// filledBlocks returns the number of filled blocks for the given strength level.
func filledBlocks(level StrengthLevel) int {
	switch level {
	case StrengthWeak:
		return 2
	case StrengthFair:
		return 8
	case StrengthStrong:
		return strengthMeterBlocks
	default:
		return 0
	}
}

// RenderStrengthMeter renders a visual representation of password strength.
// Format: "Força: ████████░░ Boa"
// Uses innerWidth to size the bar. Returns the formatted string with styling applied.
func RenderStrengthMeter(password []byte, innerWidth int, theme *design.Theme) string {
	level := EvaluateStrength(password)

	// Determine label and color based on strength level
	var label string
	var color string
	switch level {
	case StrengthWeak:
		label = "⚠ Fraca"
		color = theme.Semantic.Warning
	case StrengthFair:
		label = "Boa"
		color = theme.Semantic.Success
	case StrengthStrong:
		label = "✓ Forte"
		color = theme.Semantic.Success
	}

	// Calculate filled blocks
	filled := filledBlocks(level)
	empty := strengthMeterBlocks - filled

	// Build meter string
	filledStr := "█"
	emptyStr := "░"
	meterStr := ""
	for i := 0; i < filled; i++ {
		meterStr += filledStr
	}
	for i := 0; i < empty; i++ {
		meterStr += emptyStr
	}

	// Format output with styling
	meterStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))

	output := fmt.Sprintf("Força: %s %s",
		meterStyle.Render(meterStr),
		labelStyle.Render(label),
	)

	return output
}
