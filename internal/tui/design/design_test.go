package design

import "testing"

func TestTokyoNight_LogoGradient(t *testing.T) {
	want := [5]string{"#9d7cd8", "#89ddff", "#7aa2f7", "#7dcfff", "#bb9af7"}
	if TokyoNight.LogoGradient != want {
		t.Errorf("TokyoNight.LogoGradient = %v, want %v", TokyoNight.LogoGradient, want)
	}
}

func TestCyberpunk_LogoGradient(t *testing.T) {
	want := [5]string{"#ff2975", "#b026ff", "#00fff5", "#05ffa1", "#ff2975"}
	if Cyberpunk.LogoGradient != want {
		t.Errorf("Cyberpunk.LogoGradient = %v, want %v", Cyberpunk.LogoGradient, want)
	}
}

func TestLayoutConstants(t *testing.T) {
	if MinWidth != 80 {
		t.Errorf("MinWidth = %d, want 80", MinWidth)
	}
	if MinHeight != 24 {
		t.Errorf("MinHeight = %d, want 24", MinHeight)
	}
	const wantRatio = 0.35
	if PanelTreeRatio != wantRatio {
		t.Errorf("PanelTreeRatio = %v, want %v", PanelTreeRatio, wantRatio)
	}
}
