package design

import "testing"

func TestSymbols_TreeNavigation(t *testing.T) {
	tests := []struct{ name, got, want string }{
		{"FolderCollapsed", SymFolderCollapsed, "▶"},
		{"FolderExpanded", SymFolderExpanded, "▼"},
		{"FolderEmpty", SymFolderEmpty, "▷"},
		{"Leaf", SymLeaf, "●"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestSymbols_ItemStates(t *testing.T) {
	tests := []struct{ name, got, want string }{
		{"Favorite", SymFavorite, "★"},
		{"Deleted", SymDeleted, "✗"},
		{"Created", SymCreated, "✦"},
		{"Modified", SymModified, "✎"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestSymbols_Semantic(t *testing.T) {
	tests := []struct{ name, got, want string }{
		{"Success", SymSuccess, "✓"},
		{"Info", SymInfo, "ℹ"},
		{"Warning", SymWarning, "⚠"},
		{"Error", SymError, "✕"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestSymbols_UI(t *testing.T) {
	tests := []struct{ name, got, want string }{
		{"Revealable", SymRevealable, "◉"},
		{"Mask", SymMask, "•"},
		{"Cursor", SymCursor, "▌"},
		{"ScrollUp", SymScrollUp, "↑"},
		{"ScrollDown", SymScrollDown, "↓"},
		{"ScrollThumb", SymScrollThumb, "■"},
		{"Ellipsis", SymEllipsis, "…"},
		{"Bullet", SymBullet, "•"},
		{"HeaderSep", SymHeaderSep, "·"},
		{"TreeConnector", SymTreeConnector, "<╡"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestSymbols_BoxDrawing(t *testing.T) {
	tests := []struct{ name, got, want string }{
		{"BorderH", SymBorderH, "─"},
		{"BorderV", SymBorderV, "│"},
		{"CornerTL", SymCornerTL, "╭"},
		{"CornerTR", SymCornerTR, "╮"},
		{"CornerBL", SymCornerBL, "╰"},
		{"CornerBR", SymCornerBR, "╯"},
		{"JunctionL", SymJunctionL, "├"},
		{"JunctionT", SymJunctionT, "┬"},
		{"JunctionB", SymJunctionB, "┴"},
		{"JunctionR", SymJunctionR, "┤"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestSpinnerFrames_Values(t *testing.T) {
	want := [4]string{"◐", "◓", "◑", "◒"}
	if SpinnerFrames != want {
		t.Errorf("SpinnerFrames = %v, want %v", SpinnerFrames, want)
	}
}

func TestSpinnerFrame_Wraps(t *testing.T) {
	// frame%4 deve iterar pelos 4 frames e reiniciar
	for i := 0; i < 8; i++ {
		got := SpinnerFrame(i)
		want := SpinnerFrames[i%4]
		if got != want {
			t.Errorf("SpinnerFrame(%d) = %q, want %q", i, got, want)
		}
	}
}
