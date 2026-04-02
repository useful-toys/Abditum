package tui

import (
	"testing"
)

// TestMessageManager_InitiallyEmpty verifies Current() is nil before any Show.
func TestMessageManager_InitiallyEmpty(t *testing.T) {
	mm := NewMessageManager()
	if mm.Current() != nil {
		t.Error("Current() must be nil on a fresh MessageManager")
	}
}

// TestMessageManager_ShowAndCurrent verifies Show stores the message and
// Current() returns it with the correct fields.
func TestMessageManager_ShowAndCurrent(t *testing.T) {
	mm := NewMessageManager()
	mm.Show(MsgInfo, "hello", 5, false)

	curr := mm.Current()
	if curr == nil {
		t.Fatal("Current() must return a message after Show")
	}
	if curr.Text != "hello" {
		t.Errorf("expected text %q, got %q", "hello", curr.Text)
	}
	if curr.Kind != MsgInfo {
		t.Errorf("expected kind MsgInfo, got %v", curr.Kind)
	}
}

// TestMessageManager_Clear removes the current message immediately.
func TestMessageManager_Clear(t *testing.T) {
	mm := NewMessageManager()
	mm.Show(MsgInfo, "x", 0, false)
	mm.Clear()
	if mm.Current() != nil {
		t.Error("Current() must be nil after Clear()")
	}
}

// TestMessageManager_Tick_ExpiresByTTL verifies that Tick decrements TTL and
// clears the message when it reaches zero.
func TestMessageManager_Tick_ExpiresByTTL(t *testing.T) {
	mm := NewMessageManager()
	mm.Show(MsgInfo, "expire-me", 2, false) // TTL = 2 ticks

	mm.Tick()
	if mm.Current() == nil {
		t.Fatal("message must still be present after 1 tick (TTL=2)")
	}

	mm.Tick()
	if mm.Current() != nil {
		t.Error("message must have expired after 2 ticks (TTL=2)")
	}
}

// TestMessageManager_Tick_BusyNeverExpires verifies that MsgBusy messages
// never expire via Tick, regardless of tattempted TTL.
func TestMessageManager_Tick_BusyNeverExpires(t *testing.T) {
	mm := NewMessageManager()
	mm.Show(MsgBusy, "loading", 5, false) // ttlSeconds ignored for MsgBusy

	for i := 0; i < 10; i++ {
		mm.Tick()
	}
	if mm.Current() == nil {
		t.Error("MsgBusy must never expire via Tick()")
	}
}

// TestMessageManager_Tick_BusyAdvancesFrame verifies that each Tick advances
// the animation frame index in a 0→1→2→3→0 cycle.
func TestMessageManager_Tick_BusyAdvancesFrame(t *testing.T) {
	mm := NewMessageManager()
	mm.Show(MsgBusy, "loading", 0, false)

	expected := []int{1, 2, 3, 0, 1}
	for i, want := range expected {
		mm.Tick()
		got := mm.Current().Frame
		if got != want {
			t.Errorf("after tick %d: expected frame %d, got %d", i+1, want, got)
		}
	}
}

// TestMessageManager_HandleInput_ClearOnInput verifies that messages created
// with clearOnInput=true are cleared on the next HandleInput call.
func TestMessageManager_HandleInput_ClearOnInput(t *testing.T) {
	mm := NewMessageManager()
	mm.Show(MsgHint, "press key to dismiss", 0, true)

	mm.HandleInput()
	if mm.Current() != nil {
		t.Error("message with clearOnInput=true must be cleared by HandleInput()")
	}
}

// TestMessageManager_HandleInput_Persistent verifies that messages created
// with clearOnInput=false survive a HandleInput call.
func TestMessageManager_HandleInput_Persistent(t *testing.T) {
	mm := NewMessageManager()
	mm.Show(MsgWarn, "persistent warning", 0, false)

	mm.HandleInput()
	if mm.Current() == nil {
		t.Error("message with clearOnInput=false must NOT be cleared by HandleInput()")
	}
}
