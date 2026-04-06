package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"
	testdatapkg "github.com/useful-toys/abditum/internal/tui/testdata"
)

// ─────────────────────────────────────────────────────────────────────────────
// Action fixtures for helpModal tests
// ─────────────────────────────────────────────────────────────────────────────

// help3actions returns 3 actions in group 0 (no group header rendered).
func help3actions() []Action {
	return []Action{
		{Keys: []string{"f2"}, Description: "Novo cofre", Group: 0},
		{Keys: []string{"f3"}, Description: "Abrir cofre", Group: 0},
		{Keys: []string{"f4"}, Description: "Salvar cofre", Group: 0},
	}
}

// help15actions returns 15 actions distributed across 3 named groups.
func help15actions() []Action {
	return []Action{
		// Grupo 1: "Arquivo"
		{Keys: []string{"f2"}, Description: "Novo cofre", Group: 1},
		{Keys: []string{"f3"}, Description: "Abrir cofre existente", Group: 1},
		{Keys: []string{"f4"}, Description: "Salvar cofre atual", Group: 1},
		{Keys: []string{"f5"}, Description: "Fechar cofre", Group: 1},
		{Keys: []string{"f6"}, Description: "Exportar cofre", Group: 1},
		// Grupo 2: "Edição"
		{Keys: []string{"ctrl+n"}, Description: "Novo segredo", Group: 2},
		{Keys: []string{"ctrl+d"}, Description: "Duplicar segredo", Group: 2},
		{Keys: []string{"ctrl+x"}, Description: "Cortar segredo", Group: 2},
		{Keys: []string{"ctrl+v"}, Description: "Colar segredo", Group: 2},
		{Keys: []string{"del"}, Description: "Excluir segredo", Group: 2},
		// Grupo 3: "Geral"
		{Keys: []string{"ctrl+c"}, Description: "Copiar valor", Group: 3},
		{Keys: []string{"ctrl+s"}, Description: "Salvar alterações", Group: 3},
		{Keys: []string{"ctrl+f"}, Description: "Buscar segredos", Group: 3},
		{Keys: []string{"ctrl+r"}, Description: "Recarregar cofre", Group: 3},
		{Keys: []string{"f1"}, Description: "Esta ajuda", Group: 3},
	}
}

// helpGroupLabel returns display name for numbered action groups.
func helpGroupLabel(grp int) string {
	labels := map[int]string{1: "Arquivo", 2: "Edição", 3: "Geral"}
	return labels[grp]
}

// ─────────────────────────────────────────────────────────────────────────────
// Golden path helpers (decision-style: variant already encodes width)
// ─────────────────────────────────────────────────────────────────────────────

// helpGoldenPath returns the golden file path for a help modal scenario.
func helpGoldenPath(variant, ext string) string {
	name := fmt.Sprintf("help-%s.%s.golden", variant, ext)
	return filepath.Join("testdata", "golden", name)
}

// checkOrUpdateHelpGolden compares output against golden file, or writes it if -update or missing.
func checkOrUpdateHelpGolden(t *testing.T, variant, ext, got string) {
	t.Helper()
	path := helpGoldenPath(variant, ext)
	if *update {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("mkdirall %s: %v", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, []byte(got), 0644); err != nil {
			t.Fatalf("write golden %s: %v", path, err)
		}
		return
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		if err2 := os.MkdirAll(filepath.Dir(path), 0755); err2 != nil {
			t.Fatalf("mkdirall %s: %v", filepath.Dir(path), err2)
		}
		if err2 := os.WriteFile(path, []byte(got), 0644); err2 != nil {
			t.Fatalf("write golden %s: %v", path, err2)
		}
		return
	}
	if err != nil {
		t.Fatalf("read golden %s: %v", path, err)
	}
	if string(data) != got {
		t.Errorf("golden mismatch for %s:\nwant:\n%s\ngot:\n%s", path, string(data), got)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// TestHelpModal_Golden
// ─────────────────────────────────────────────────────────────────────────────

func TestHelpModal_Golden(t *testing.T) {
	type testCase struct {
		variant string
		modal   *helpModal
	}

	// Helper: build 15-action modal and compute maxScroll for "bottom" scenarios.
	make15 := func(w, h, scroll int) *helpModal {
		m := newHelpModal(help15actions(), helpGroupLabel)
		m.SetSize(w, h)
		if scroll < 0 {
			// Negative sentinel: use actual maxScroll.
			maxScroll := m.totalLines() - m.contentHeight()
			if maxScroll < 0 {
				maxScroll = 0
			}
			m.scroll = maxScroll
		} else {
			m.scroll = scroll
		}
		return m
	}

	cases := []testCase{
		// Set 1: 3 actions, height=12, no scroll needed
		{
			variant: "3actions-30x12",
			modal: func() *helpModal {
				m := newHelpModal(help3actions(), helpGroupLabel)
				m.SetSize(30, 12)
				return m
			}(),
		},
		{
			variant: "3actions-60x12",
			modal: func() *helpModal {
				m := newHelpModal(help3actions(), helpGroupLabel)
				m.SetSize(60, 12)
				return m
			}(),
		},
		// Set 2: 15 actions, height=16, 3 scroll positions × 2 widths
		{variant: "15actions-top-30x16", modal: make15(30, 16, 0)},
		{variant: "15actions-top-60x16", modal: make15(60, 16, 0)},
		{variant: "15actions-mid-30x16", modal: make15(30, 16, 3)},
		{variant: "15actions-mid-60x16", modal: make15(60, 16, 3)},
		{variant: "15actions-bottom-30x16", modal: make15(30, 16, -1)}, // -1 = maxScroll
		{variant: "15actions-bottom-60x16", modal: make15(60, 16, -1)}, // -1 = maxScroll
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.variant, func(t *testing.T) {
			out := tc.modal.View()

			// .txt.golden: raw ANSI output
			checkOrUpdateHelpGolden(t, tc.variant, "txt", stripANSI(out))

			// .json.golden: style transitions
			transitions := testdatapkg.ParseANSIStyle(out)
			jsonBytes, err := testdatapkg.MarshalStyleTransitions(transitions)
			if err != nil {
				t.Fatalf("marshal transitions: %v", err)
			}
			checkOrUpdateHelpGolden(t, tc.variant, "json", string(jsonBytes))
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// TestHelpModal_Update_* — keyboard navigation and dismiss
// ─────────────────────────────────────────────────────────────────────────────

// TestHelpModal_Update_DownIncreasesScroll: "down" key increases scroll by 1.
func TestHelpModal_Update_DownIncreasesScroll(t *testing.T) {
	m := newHelpModal(help15actions(), helpGroupLabel)
	m.SetSize(60, 16)
	initial := m.scroll // 0
	m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if m.scroll != initial+1 {
		t.Errorf("down key: expected scroll %d, got %d", initial+1, m.scroll)
	}
}

// TestHelpModal_Update_UpDoesNotGoBelowZero: "up" at scroll=0 stays at 0.
func TestHelpModal_Update_UpDoesNotGoBelowZero(t *testing.T) {
	m := newHelpModal(help15actions(), helpGroupLabel)
	m.SetSize(60, 16)
	m.scroll = 0
	m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if m.scroll != 0 {
		t.Errorf("up at scroll=0: expected 0, got %d", m.scroll)
	}
}

// TestHelpModal_Update_UpDecreasesScroll: "up" at scroll>0 decreases by 1.
func TestHelpModal_Update_UpDecreasesScroll(t *testing.T) {
	m := newHelpModal(help15actions(), helpGroupLabel)
	m.SetSize(60, 16)
	m.scroll = 3
	m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if m.scroll != 2 {
		t.Errorf("up at scroll=3: expected 2, got %d", m.scroll)
	}
}

// TestHelpModal_Update_HomeResetsScroll: "home" key sets scroll to 0.
func TestHelpModal_Update_HomeResetsScroll(t *testing.T) {
	m := newHelpModal(help15actions(), helpGroupLabel)
	m.SetSize(60, 16)
	m.scroll = 5
	m.Update(tea.KeyPressMsg{Code: tea.KeyHome})
	if m.scroll != 0 {
		t.Errorf("home key: expected scroll 0, got %d", m.scroll)
	}
}

// TestHelpModal_Update_EndScrollsToBottom: "end" key clamps scroll to maxScroll.
func TestHelpModal_Update_EndScrollsToBottom(t *testing.T) {
	m := newHelpModal(help15actions(), helpGroupLabel)
	m.SetSize(60, 16)
	maxScroll := m.totalLines() - m.contentHeight()
	if maxScroll < 0 {
		maxScroll = 0
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnd})
	if m.scroll != maxScroll {
		t.Errorf("end key: expected scroll %d (maxScroll), got %d", maxScroll, m.scroll)
	}
}

// TestHelpModal_Update_ScrollClampedAtMax: down past maxScroll stays at maxScroll.
func TestHelpModal_Update_ScrollClampedAtMax(t *testing.T) {
	m := newHelpModal(help15actions(), helpGroupLabel)
	m.SetSize(60, 16)
	maxScroll := m.totalLines() - m.contentHeight()
	if maxScroll < 0 {
		maxScroll = 0
	}
	m.scroll = maxScroll
	m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if m.scroll != maxScroll {
		t.Errorf("down at maxScroll: expected %d (clamped), got %d", maxScroll, m.scroll)
	}
}

// TestHelpModal_Update_EscDismisses: "esc" key returns non-nil cmd (popModalMsg).
func TestHelpModal_Update_EscDismisses(t *testing.T) {
	m := newHelpModal(help3actions(), helpGroupLabel)
	m.SetSize(60, 12)
	cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if cmd == nil {
		t.Error("esc key must return non-nil cmd (pop modal)")
	}
}

// TestHelpModal_Update_F1Dismisses: "f1" key returns non-nil cmd (popModalMsg).
func TestHelpModal_Update_F1Dismisses(t *testing.T) {
	m := newHelpModal(help3actions(), helpGroupLabel)
	m.SetSize(60, 12)
	cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyF1})
	if cmd == nil {
		t.Error("f1 key must return non-nil cmd (pop modal)")
	}
}
