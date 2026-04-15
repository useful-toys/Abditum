package design

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestKeys_Labels(t *testing.T) {
	tests := []struct{ name, got, want string }{
		{"Enter", Keys.Enter.Label, "Enter"},
		{"Esc", Keys.Esc.Label, "Esc"},
		{"Tab", Keys.Tab.Label, "Tab"},
		{"Del", Keys.Del.Label, "Del"},
		{"Ins", Keys.Ins.Label, "Ins"},
		{"Home", Keys.Home.Label, "Home"},
		{"End", Keys.End.Label, "End"},
		{"PgUp", Keys.PgUp.Label, "PgUp"},
		{"PgDn", Keys.PgDn.Label, "PgDn"},
		{"Up", Keys.Up.Label, "↑"},
		{"Down", Keys.Down.Label, "↓"},
		{"F1", Keys.F1.Label, "F1"},
		{"F6", Keys.F6.Label, "F6"},
		{"F12", Keys.F12.Label, "F12"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("Keys.%s.Label = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestKeys_Matches_SimpleKeys(t *testing.T) {
	tests := []struct {
		name string
		key  Key
		msg  tea.KeyPressMsg
	}{
		{"Enter", Keys.Enter, tea.KeyPressMsg{Code: tea.KeyEnter}},
		{"Esc", Keys.Esc, tea.KeyPressMsg{Code: tea.KeyEscape}},
		{"Tab", Keys.Tab, tea.KeyPressMsg{Code: tea.KeyTab}},
		{"Del", Keys.Del, tea.KeyPressMsg{Code: tea.KeyDelete}},
		{"Ins", Keys.Ins, tea.KeyPressMsg{Code: tea.KeyInsert}},
		{"Home", Keys.Home, tea.KeyPressMsg{Code: tea.KeyHome}},
		{"End", Keys.End, tea.KeyPressMsg{Code: tea.KeyEnd}},
		{"PgUp", Keys.PgUp, tea.KeyPressMsg{Code: tea.KeyPgUp}},
		{"PgDn", Keys.PgDn, tea.KeyPressMsg{Code: tea.KeyPgDown}},
		{"Up", Keys.Up, tea.KeyPressMsg{Code: tea.KeyUp}},
		{"Down", Keys.Down, tea.KeyPressMsg{Code: tea.KeyDown}},
		{"F1", Keys.F1, tea.KeyPressMsg{Code: tea.KeyF1}},
		{"F6", Keys.F6, tea.KeyPressMsg{Code: tea.KeyF6}},
		{"F12", Keys.F12, tea.KeyPressMsg{Code: tea.KeyF12}},
	}
	for _, tt := range tests {
		if !tt.key.Matches(tt.msg) {
			t.Errorf("Keys.%s.Matches(correct msg) = false, want true", tt.name)
		}
	}
}

func TestKeys_Matches_DoesNotMatchWrong(t *testing.T) {
	// Enter não deve casar com Esc
	if Keys.Enter.Matches(tea.KeyPressMsg{Code: tea.KeyEscape}) {
		t.Error("Keys.Enter.Matches(Esc) = true, want false")
	}
	// F1 sem modificador não deve casar com Ctrl+F1
	if Keys.F1.Matches(tea.KeyPressMsg{Code: tea.KeyF1, Mod: tea.ModCtrl}) {
		t.Error("Keys.F1.Matches(ctrl+F1) = true, want false")
	}
}

func TestLetter(t *testing.T) {
	k := Letter('q')
	if k.Label != "Q" {
		t.Errorf("Letter('q').Label = %q, want \"Q\"", k.Label)
	}
	if k.Code != 'q' {
		t.Errorf("Letter('q').Code = %v, want 'q'", k.Code)
	}
	if k.Mod != 0 {
		t.Errorf("Letter('q').Mod = %v, want 0", k.Mod)
	}
	if !k.Matches(tea.KeyPressMsg{Code: 'q'}) {
		t.Error("Letter('q').Matches({Code:'q'}) = false, want true")
	}
}

func TestWithCtrl(t *testing.T) {
	k := WithCtrl(Letter('q'))
	if k.Label != "⌃Q" {
		t.Errorf("WithCtrl(Letter('q')).Label = %q, want \"⌃Q\"", k.Label)
	}
	if k.Code != 'q' {
		t.Errorf("WithCtrl(Letter('q')).Code = %v, want 'q'", k.Code)
	}
	if !k.Mod.Contains(tea.ModCtrl) {
		t.Error("WithCtrl(Letter('q')).Mod deve conter ModCtrl")
	}
	if !k.Matches(tea.KeyPressMsg{Code: 'q', Mod: tea.ModCtrl}) {
		t.Error("WithCtrl(Letter('q')).Matches(ctrl+q) = false, want true")
	}
	// sem modificador não deve casar
	if k.Matches(tea.KeyPressMsg{Code: 'q'}) {
		t.Error("WithCtrl(Letter('q')).Matches(q) = true, want false")
	}
}

func TestWithShift(t *testing.T) {
	k := WithShift(Keys.F6)
	if k.Label != "⇧F6" {
		t.Errorf("WithShift(Keys.F6).Label = %q, want \"⇧F6\"", k.Label)
	}
	if !k.Matches(tea.KeyPressMsg{Code: tea.KeyF6, Mod: tea.ModShift}) {
		t.Error("WithShift(Keys.F6).Matches(shift+F6) = false, want true")
	}
	// sem modificador não deve casar
	if k.Matches(tea.KeyPressMsg{Code: tea.KeyF6}) {
		t.Error("WithShift(Keys.F6).Matches(F6) = true, want false")
	}
}

func TestWithAlt(t *testing.T) {
	k := WithAlt(Letter('x'))
	if k.Label != "!X" {
		t.Errorf("WithAlt(Letter('x')).Label = %q, want \"!X\"", k.Label)
	}
	if !k.Matches(tea.KeyPressMsg{Code: 'x', Mod: tea.ModAlt}) {
		t.Error("WithAlt(Letter('x')).Matches(alt+x) = false, want true")
	}
}

func TestComposition_CtrlAltShift(t *testing.T) {
	// Ordem de composição: WithCtrl(WithAlt(WithShift(base))) produz "⌃!⇧Q"
	k := WithCtrl(WithAlt(WithShift(Letter('q'))))
	if k.Label != "⌃!⇧Q" {
		t.Errorf("label = %q, want \"⌃!⇧Q\"", k.Label)
	}
	msg := tea.KeyPressMsg{Code: 'q', Mod: tea.ModCtrl | tea.ModAlt | tea.ModShift}
	if !k.Matches(msg) {
		t.Error("Matches(ctrl+alt+shift+q) = false, want true")
	}
}

func TestShortcuts_Labels(t *testing.T) {
	tests := []struct{ name, got, want string }{
		{"Help", Shortcuts.Help.Label, "F1"},
		{"ThemeToggle", Shortcuts.ThemeToggle.Label, "F12"},
		{"Quit", Shortcuts.Quit.Label, "⌃Q"},
		{"LockVault", Shortcuts.LockVault.Label, "⌃!⇧Q"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("Shortcuts.%s.Label = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestShortcuts_Matches(t *testing.T) {
	if !Shortcuts.Help.Matches(tea.KeyPressMsg{Code: tea.KeyF1}) {
		t.Error("Shortcuts.Help.Matches(F1) = false, want true")
	}
	if !Shortcuts.ThemeToggle.Matches(tea.KeyPressMsg{Code: tea.KeyF12}) {
		t.Error("Shortcuts.ThemeToggle.Matches(F12) = false, want true")
	}
	if !Shortcuts.Quit.Matches(tea.KeyPressMsg{Code: 'q', Mod: tea.ModCtrl}) {
		t.Error("Shortcuts.Quit.Matches(ctrl+q) = false, want true")
	}
	emergency := tea.KeyPressMsg{Code: 'q', Mod: tea.ModCtrl | tea.ModAlt | tea.ModShift}
	if !Shortcuts.LockVault.Matches(emergency) {
		t.Error("Shortcuts.LockVault.Matches(ctrl+alt+shift+q) = false, want true")
	}
}
