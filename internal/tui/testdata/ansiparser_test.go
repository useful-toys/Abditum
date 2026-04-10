package testdata

import (
	"testing"
)

func TestParseANSIStyle_EmptyInput(t *testing.T) {
	result := ParseANSIStyle("")
	if len(result) != 0 {
		t.Errorf("expected empty result for empty input, got %d transitions", len(result))
	}
}

func TestParseANSIStyle_PlainText(t *testing.T) {
	result := ParseANSIStyle("Hello, world!")
	if len(result) != 0 {
		t.Errorf("expected no transitions for plain text, got %d", len(result))
	}
}

func TestParseANSIStyle_SingleColor(t *testing.T) {
	// \x1b[1m = bold
	input := "\x1b[1mBold text\x1b[0m"
	result := ParseANSIStyle(input)

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

func TestParseANSIStyle_LineTracking(t *testing.T) {
	// Two lines separated by newline
	input := "Line 1\n\x1b[1mLine 2\x1b[0m"
	result := ParseANSIStyle(input)

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

func TestParseANSIStyle_MultipleTransitions(t *testing.T) {
	// Bold on, bold off, bold on again
	input := "\x1b[1mA\x1b[0mB\x1b[1mC"
	result := ParseANSIStyle(input)

	if len(result) < 3 {
		t.Fatalf("expected at least 3 transitions, got %d", len(result))
	}
}
