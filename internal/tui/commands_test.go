package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

// TestQuitWithCleanup verifies that QuitWithCleanup returns a valid Cmd
// that produces a QuitMsg when executed.
func TestQuitWithCleanup(t *testing.T) {
	cmd := QuitWithCleanup()
	if cmd == nil {
		t.Fatal("QuitWithCleanup() returned nil, expected a command")
	}

	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Fatalf("Expected tea.QuitMsg, got %T", msg)
	}
}
