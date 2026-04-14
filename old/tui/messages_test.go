package tui

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	testdatapkg "github.com/useful-toys/abditum/internal/tui/testdata"
)

// update flag: when set, golden files are regenerated instead of compared.
// Usage: go test ./internal/tui/... -run TestRenderMessageBar_Golden -update
var update = flag.Bool("update", false, "regenerate golden files")

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
	mm.Show(MessageSuccess, "hello", 5, false)

	curr := mm.Current()
	if curr == nil {
		t.Fatal("Current() must return a message after Show")
	}
	if curr.Text != "hello" {
		t.Errorf("expected text %q, got %q", "hello", curr.Text)
	}
	if curr.Kind != MessageSuccess {
		t.Errorf("expected kind MessageSuccess, got %v", curr.Kind)
	}
}

// TestMessageManager_Clear removes the current message immediately.
func TestMessageManager_Clear(t *testing.T) {
	mm := NewMessageManager()
	mm.Show(MessageSuccess, "x", 0, false)
	mm.Clear()
	if mm.Current() != nil {
		t.Error("Current() must be nil after Clear()")
	}
}

// TestMessageManager_Tick_ExpiresByTTL verifies that Tick decrements TTL and
// clears the message when it reaches zero.
func TestMessageManager_Tick_ExpiresByTTL(t *testing.T) {
	mm := NewMessageManager()
	mm.Show(MessageSuccess, "expire-me", 2, false) // TTL = 2 ticks

	mm.Tick()
	if mm.Current() == nil {
		t.Fatal("message must still be present after 1 tick (TTL=2)")
	}

	mm.Tick()
	if mm.Current() != nil {
		t.Error("message must have expired after 2 ticks (TTL=2)")
	}
}

// TestMessageManager_Tick_BusyNeverExpires verifies that MessageBusy messages
// never expire via Tick, regardless of tattempted TTL.
func TestMessageManager_Tick_BusyNeverExpires(t *testing.T) {
	mm := NewMessageManager()
	mm.Show(MessageBusy, "loading", 5, false) // ttlSeconds ignored for MessageBusy

	for i := 0; i < 10; i++ {
		mm.Tick()
	}
	if mm.Current() == nil {
		t.Error("MessageBusy must never expire via Tick()")
	}
}

// TestMessageManager_Tick_BusyAdvancesFrame verifies that each Tick advances
// the animation frame index in a 0→1→2→3→0 cycle.
func TestMessageManager_Tick_BusyAdvancesFrame(t *testing.T) {
	mm := NewMessageManager()
	mm.Show(MessageBusy, "loading", 0, false)

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
	mm.Show(MessageHint, "press key to dismiss", 0, true)

	mm.HandleInput()
	if mm.Current() != nil {
		t.Error("message with clearOnInput=true must be cleared by HandleInput()")
	}
}

// TestMessageManager_MsgKindDistinct verifies MessageSuccess and MessageInfo are distinct values.
func TestMessageManager_MsgKindDistinct(t *testing.T) {
	if MessageSuccess == MessageInfo {
		t.Error("MessageSuccess and MessageInfo must be distinct MessageKind values")
	}
	mm := NewMessageManager()
	mm.Show(MessageInfo, "neutral", 3, false)
	curr := mm.Current()
	if curr == nil || curr.Kind != MessageInfo {
		t.Error("Show(MessageInfo, ...) must store MessageInfo kind")
	}
}

// TestMessageManager_HandleInput_Persistent verifies that messages created
// with clearOnInput=false survive a HandleInput call.
func TestMessageManager_HandleInput_Persistent(t *testing.T) {
	mm := NewMessageManager()
	mm.Show(MessageWarning, "persistent warning", 0, false)

	mm.HandleInput()
	if mm.Current() == nil {
		t.Error("message with clearOnInput=false must NOT be cleared by HandleInput()")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Golden test helpers
// ─────────────────────────────────────────────────────────────────────────────

// goldenPath returns the path to a golden file for a test case.
// Files are stored under testdata/golden/ with the naming convention:
//
//	{component}-{variant}-{width}.{ext}.golden
//
// Example: testdata/golden/messages-success-30.txt.golden
func goldenPath(component, variant string, width int, ext string) string {
	name := fmt.Sprintf("%s-%s-%d.%s.golden", component, variant, width, ext)
	return filepath.Join("testdata", "golden", name)
}

// checkOrUpdateGolden compares output against a golden file, or writes it if
// -update is set or the file does not yet exist (first run auto-generation).
// On mismatch, the test is failed with a diff showing want vs got.
func checkOrUpdateGolden(t *testing.T, path, got string) {
	t.Helper()
	if *update {
		// Explicit regeneration: overwrite unconditionally.
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
		// First run: auto-create the golden file (treated as baseline).
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

// stripANSI removes all ANSI escape sequences from s, returning plain visible text.
// Used when writing .txt.golden files — these must contain no escape codes.
func stripANSI(s string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(s, "")
}

// ─────────────────────────────────────────────────────────────────────────────
// TestRenderMessageBar_Golden
// ─────────────────────────────────────────────────────────────────────────────

// TestRenderMessageBar_Golden validates the visual output of RenderMessageBar
// against golden files for all 6 message kinds × 2 terminal widths (= 12 sub-tests).
//
// Each sub-test produces two golden files:
//   - .txt.golden — raw ANSI output, validated byte-for-byte
//   - .json.golden — style transitions from ParseANSIStyle, validated byte-for-byte
//
// First run (no golden files present) auto-generates the baselines.
// Subsequent runs compare against the recorded baselines.
// Run with -update to intentionally regenerate all baselines.
func TestRenderMessageBar_Golden(t *testing.T) {
	type testCase struct {
		variant string
		msg     *DisplayMessage
	}

	cases := []testCase{
		{"success", &DisplayMessage{Text: "Cofre salvo com sucesso — 12 segredos sincronizados", Kind: MessageSuccess, Frame: 0}},
		{"error", &DisplayMessage{Text: "Falha ao salvar o cofre — verifique permissões do arquivo", Kind: MessageError, Frame: 0}},
		{"warn", &DisplayMessage{Text: "Cofre modificado externamente — revisar antes de salvar", Kind: MessageWarning, Frame: 0}},
		{"info", &DisplayMessage{Text: "Cofre aberto — 12 segredos, 3 pastas, 2 modelos", Kind: MessageInfo, Frame: 0}},
		{"busy", &DisplayMessage{Text: "Carregando cofre, por favor aguarde...", Kind: MessageBusy, Frame: 0}},
		{"hint", &DisplayMessage{Text: "Pressione F1 para ver todos os atalhos disponíveis", Kind: MessageHint, Frame: 0}},
	}
	widths := []int{30, 60}

	for _, tc := range cases {
		for _, w := range widths {
			tc := tc // capture loop vars
			w := w
			name := fmt.Sprintf("%s-%d", tc.variant, w)
			t.Run(name, func(t *testing.T) {
				out := RenderMessageBar(tc.msg, w, TokyoNight)

				// .txt.golden: raw ANSI output — validates layout, spacing, truncation, borders
				txtPath := goldenPath("messages", tc.variant, w, "txt")
				checkOrUpdateGolden(t, txtPath, stripANSI(out))

				// .json.golden: style transitions — validates colors and font attributes
				transitions := testdatapkg.ParseANSIStyle(out)
				jsonBytes, err := testdatapkg.MarshalStyleTransitions(transitions)
				if err != nil {
					t.Fatalf("marshal transitions: %v", err)
				}
				jsonPath := goldenPath("messages", tc.variant, w, "json")
				checkOrUpdateGolden(t, jsonPath, string(jsonBytes))
			})
		}
	}
}
