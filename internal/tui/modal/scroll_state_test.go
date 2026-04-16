package modal

import "testing"

func TestScrollState_CanScrollUp_Down(t *testing.T) {
	s := ScrollState{Offset: 0, Total: 20, Viewport: 10}
	if s.CanScrollUp() {
		t.Error("CanScrollUp should be false when Offset=0")
	}
	if !s.CanScrollDown() {
		t.Error("CanScrollDown should be true when Total > Viewport")
	}

	s.Offset = 5
	if !s.CanScrollUp() {
		t.Error("CanScrollUp should be true when Offset > 0")
	}
	if !s.CanScrollDown() {
		t.Error("CanScrollDown should be true when offset+viewport < total")
	}

	s.Offset = 10
	// offset + viewport == total → cannot scroll down
	if s.CanScrollDown() {
		t.Error("CanScrollDown should be false when offset+viewport == total")
	}
}

func TestScrollState_Up_Down(t *testing.T) {
	s := ScrollState{Offset: 5, Total: 20, Viewport: 10}
	s.Up()
	if s.Offset != 4 {
		t.Errorf("Up: Offset = %d, want 4", s.Offset)
	}
	s.Offset = 0
	s.Up()
	if s.Offset != 0 {
		t.Errorf("Up at 0: Offset = %d, want 0 (clamped)", s.Offset)
	}
	s.Offset = 5
	s.Down()
	if s.Offset != 6 {
		t.Errorf("Down: Offset = %d, want 6", s.Offset)
	}
	s.Offset = 10 // offset + viewport == total → at bottom
	s.Down()
	if s.Offset != 10 {
		t.Errorf("Down at bottom: Offset = %d, want 10 (clamped)", s.Offset)
	}
}

func TestScrollState_PageUp_PageDown(t *testing.T) {
	s := ScrollState{Offset: 15, Total: 40, Viewport: 10}
	s.PageUp()
	if s.Offset != 5 {
		t.Errorf("PageUp: Offset = %d, want 5", s.Offset)
	}
	s.PageUp()
	if s.Offset != 0 {
		t.Errorf("PageUp clamp: Offset = %d, want 0", s.Offset)
	}

	s.Offset = 15
	s.PageDown()
	if s.Offset != 25 {
		t.Errorf("PageDown: Offset = %d, want 25", s.Offset)
	}
	s.PageDown()
	if s.Offset != 30 {
		t.Errorf("PageDown clamp: Offset = %d, want 30 (total-viewport)", s.Offset)
	}
}

func TestScrollState_Home_End(t *testing.T) {
	s := ScrollState{Offset: 15, Total: 40, Viewport: 10}
	s.Home()
	if s.Offset != 0 {
		t.Errorf("Home: Offset = %d, want 0", s.Offset)
	}
	s.End()
	if s.Offset != 30 {
		t.Errorf("End: Offset = %d, want 30 (total-viewport)", s.Offset)
	}
}

func TestScrollState_ThumbLine_InactiveWhenNoScroll(t *testing.T) {
	// Total <= Viewport → no scroll → ThumbLine returns -1
	s := ScrollState{Offset: 0, Total: 5, Viewport: 10}
	if got := s.ThumbLine(); got != -1 {
		t.Errorf("ThumbLine (no scroll): got %d, want -1", got)
	}
	s.Total = 10
	if got := s.ThumbLine(); got != -1 {
		t.Errorf("ThumbLine (total==viewport): got %d, want -1", got)
	}
}

func TestScrollState_ThumbLine_ArrowsHavePriority(t *testing.T) {
	// At top: ↑ inactive, ↓ active on last line.
	// Thumb must NOT be on line viewport (last line — occupied by ↓).
	s := ScrollState{Offset: 0, Total: 30, Viewport: 10}
	thumb := s.ThumbLine()
	if thumb == s.Viewport {
		t.Errorf("ThumbLine at top: thumb on last line %d (occupied by ↓)", thumb)
	}
	if thumb < 1 || thumb > s.Viewport {
		t.Errorf("ThumbLine at top: got %d, want in [1..%d]", thumb, s.Viewport)
	}

	// At bottom: ↑ active on line 1, ↓ inactive.
	// Thumb must NOT be on line 1.
	s.Offset = 20 // offset + viewport == total
	thumb = s.ThumbLine()
	if thumb == 1 {
		t.Errorf("ThumbLine at bottom: thumb on line 1 (occupied by ↑)")
	}
	if thumb < 1 || thumb > s.Viewport {
		t.Errorf("ThumbLine at bottom: got %d, want in [1..%d]", thumb, s.Viewport)
	}

	// In middle: both ↑ and ↓ active.
	// Thumb must NOT be on line 1 or line viewport.
	s.Offset = 10
	thumb = s.ThumbLine()
	if thumb == 1 {
		t.Errorf("ThumbLine in middle: thumb on line 1 (occupied by ↑)")
	}
	if thumb == s.Viewport {
		t.Errorf("ThumbLine in middle: thumb on last line (occupied by ↓)")
	}
	if thumb < 2 || thumb > s.Viewport-1 {
		t.Errorf("ThumbLine in middle: got %d, want in [2..%d]", thumb, s.Viewport-1)
	}
}

func TestScrollState_ThumbLine_ReturnsMinus1WhenNoSpace(t *testing.T) {
	// Viewport = 2, both arrows active → no room for thumb.
	s := ScrollState{Offset: 5, Total: 20, Viewport: 2}
	if got := s.ThumbLine(); got != -1 {
		t.Errorf("ThumbLine tiny viewport: got %d, want -1", got)
	}
}
