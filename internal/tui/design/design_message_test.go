package design

import (
	"testing"
)

func TestMessageKind_DefaultTTL(t *testing.T) {
	tests := []struct {
		kind MessageKind
		want int
	}{
		{MsgSuccess, 5},
		{MsgInfo, 5},
		{MsgWarning, 5},
		{MsgError, 5},
		{MsgBusy, 0},
		{MsgHintField, 0},
		{MsgHintUsage, 0},
	}
	for _, tt := range tests {
		if got := tt.kind.DefaultTTL(); got != tt.want {
			t.Errorf("MessageKind(%d).DefaultTTL() = %d, want %d", tt.kind, got, tt.want)
		}
	}
}

func TestMessageKind_Symbol(t *testing.T) {
	tests := []struct {
		kind MessageKind
		want string
	}{
		{MsgSuccess, SymSuccess},    // "✓"
		{MsgInfo, SymInfo},          // "ℹ"
		{MsgWarning, SymWarning},    // "⚠"
		{MsgError, SymError},        // "✕"
		{MsgBusy, SpinnerFrames[0]}, // "◐"
		{MsgHintField, SymBullet},   // "•"
		{MsgHintUsage, SymBullet},   // "•"
	}
	for _, tt := range tests {
		if got := tt.kind.Symbol(); got != tt.want {
			t.Errorf("MessageKind(%d).Symbol() = %q, want %q", tt.kind, got, tt.want)
		}
	}
}

func TestMessageKind_Color_TokyoNight(t *testing.T) {
	theme := TokyoNight
	tests := []struct {
		kind MessageKind
		want string
	}{
		{MsgSuccess, theme.Semantic.Success},
		{MsgInfo, theme.Semantic.Info},
		{MsgWarning, theme.Semantic.Warning},
		{MsgError, theme.Semantic.Error},
		{MsgBusy, theme.Accent.Primary},
		{MsgHintField, theme.Text.Secondary},
		{MsgHintUsage, theme.Text.Secondary},
	}
	for _, tt := range tests {
		if got := tt.kind.Color(theme); got != tt.want {
			t.Errorf("MessageKind(%d).Color(TokyoNight) = %q, want %q", tt.kind, got, tt.want)
		}
	}
}

func TestMessageHelpers(t *testing.T) {
	tests := []struct {
		name string
		msg  Message
		kind MessageKind
		text string
	}{
		{"Success", Success("ok"), MsgSuccess, "ok"},
		{"Error", Error("fail"), MsgError, "fail"},
		{"Info", Info("note"), MsgInfo, "note"},
		{"Warning", Warning("warn"), MsgWarning, "warn"},
		{"HintField", HintField("hint"), MsgHintField, "hint"},
		{"HintUsage", HintUsage("usage"), MsgHintUsage, "usage"},
	}
	for _, tt := range tests {
		if tt.msg.Kind != tt.kind {
			t.Errorf("%s: Kind = %d, want %d", tt.name, tt.msg.Kind, tt.kind)
		}
		if tt.msg.Text != tt.text {
			t.Errorf("%s: Text = %q, want %q", tt.name, tt.msg.Text, tt.text)
		}
	}
}

func TestBusyHelper_WithText(t *testing.T) {
	msg := Busy("carregando...")
	if msg.Kind != MsgBusy {
		t.Errorf("Busy().Kind = %d, want %d", msg.Kind, MsgBusy)
	}
	if msg.Text != "carregando..." {
		t.Errorf("Busy().Text = %q, want %q", msg.Text, "carregando...")
	}
}

func TestBusyHelper_WithoutText(t *testing.T) {
	msg := Busy()
	if msg.Kind != MsgBusy {
		t.Errorf("Busy().Kind = %d, want %d", msg.Kind, MsgBusy)
	}
	if msg.Text != "" {
		t.Errorf("Busy() without text: Text = %q, want empty", msg.Text)
	}
}

func TestBusyHelper_SpinnerFrame_ZeroValue(t *testing.T) {
	msg := Busy("teste")
	if msg.SpinnerFrame != 0 {
		t.Errorf("Busy(): SpinnerFrame = %d, want 0 (zero value)", msg.SpinnerFrame)
	}
}
