package testdata

import (
	"testing"
)

func TestAnsiToStyleChanges_EmptyInput(t *testing.T) {
	result := ansiToStyleChanges("")
	if len(result) != 0 {
		t.Errorf("expected empty result for empty input, got %d transitions", len(result))
	}
}

func TestAnsiToStyleChanges_PlainText(t *testing.T) {
	result := ansiToStyleChanges("Hello, world!")
	if len(result) != 0 {
		t.Errorf("expected no transitions for plain text, got %d", len(result))
	}
}

func TestAnsiToStyleChanges_SingleColor(t *testing.T) {
	// \x1b[1m = bold
	input := "\x1b[1mBold text\x1b[0m"
	result := ansiToStyleChanges(input)

	if len(result) != 2 {
		t.Fatalf("expected 2 transitions (bold on, reset), got %d", len(result))
	}

	if len(result[0].Style) != 1 || result[0].Style[0] != "bold" {
		t.Errorf("expected first transition to have bold, got %v", result[0].Style)
	}

	if len(result[1].Style) != 0 {
		t.Errorf("expected second transition to reset style, got %v", result[1].Style)
	}
}

func TestAnsiToStyleChanges_LineTracking(t *testing.T) {
	// Two lines separated by newline
	input := "Line 1\n\x1b[1mLine 2\x1b[0m"
	result := ansiToStyleChanges(input)

	if len(result) < 1 {
		t.Fatalf("expected at least 1 transition, got %d", len(result))
	}

	// Should have transition on line 1 (second line)
	foundLine1 := false
	for _, tr := range result {
		if tr.Line == 1 {
			foundLine1 = true
			break
		}
	}
	if !foundLine1 {
		t.Errorf("expected transition on line 1, got lines: %v", func() []int {
			var lines []int
			for _, tr := range result {
				lines = append(lines, tr.Line)
			}
			return lines
		}())
	}
}

func TestAnsiToStyleChanges_MultipleTransitions(t *testing.T) {
	// Bold on, bold off, bold on again
	input := "\x1b[1mA\x1b[0mB\x1b[1mC"
	result := ansiToStyleChanges(input)

	if len(result) < 3 {
		t.Fatalf("expected at least 3 transitions, got %d", len(result))
	}
}

// TestAnsiToStyleChanges_CombinedStateMultipleSGR validates correct handling when
// multiple SGR codes appear at the same position. The final state should prevail,
// with only one transition recorded at that position (substitution, not addition).
func TestAnsiToStyleChanges_CombinedStateMultipleSGR(t *testing.T) {
	// \x1b[1;31m = bold + red foreground at the same position
	// This tests that multiple SGR codes at the same column result in
	// a single transition with the combined final state.
	input := "\x1b[1;31mText\x1b[0m"
	result := ansiToStyleChanges(input)

	if len(result) < 1 {
		t.Fatalf("expected at least 1 transition, got %d", len(result))
	}

	// First transition should have bold style and red foreground
	firstTransition := result[0]

	// Check bold is present
	if len(firstTransition.Style) == 0 {
		t.Errorf("expected style to be set, got empty style")
	}
	hasBold := false
	for _, s := range firstTransition.Style {
		if s == "bold" {
			hasBold = true
			break
		}
	}
	if !hasBold {
		t.Errorf("expected bold in style, got %v", firstTransition.Style)
	}

	// Check red foreground color
	if firstTransition.FG == nil {
		t.Errorf("expected FG color to be set, got nil")
	} else if *firstTransition.FG != "#800000" {
		t.Errorf("expected red (#800000), got %s", *firstTransition.FG)
	}

	// Should have second transition for reset
	if len(result) < 2 {
		t.Fatalf("expected at least 2 transitions (combined + reset), got %d", len(result))
	}

	secondTransition := result[1]
	if len(secondTransition.Style) != 0 {
		t.Errorf("expected reset (empty style), got %v", secondTransition.Style)
	}
	if secondTransition.FG != nil {
		t.Errorf("expected FG to be nil after reset, got %s", *secondTransition.FG)
	}
}
