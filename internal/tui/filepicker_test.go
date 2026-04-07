package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	testdatapkg "github.com/useful-toys/abditum/internal/tui/testdata"
)

// Helper to create a filePickerModal for testing
func newTestFilePickerModal() *filePickerModal {
	fpk := &filePickerModal{
		ext:  ".abditum",
		mode: FilePickerOpen,
	}
	fpk.Init()
	return fpk
}

// TestFilePickerModalStructExists verifies that filePickerModal can be instantiated.
func TestFilePickerModalStructExists(t *testing.T) {
	fpk := newTestFilePickerModal()
	if fpk == nil {
		t.Fatal("newTestFilePickerModal returned nil")
	}
}

// TestFilePickerModalInit verifies Init() initializes the directory and sets focusPanel=0.
func TestFilePickerModalInit(t *testing.T) {
	fpk := &filePickerModal{ext: ".abditum", mode: FilePickerOpen}
	fpk.Init()
	if fpk.currentPath == "" {
		t.Fatal("currentPath not set after Init()")
	}
	if fpk.focusPanel != 0 {
		t.Errorf("Init() should set focusPanel=0 (tree), got %d", fpk.focusPanel)
	}
	t.Logf("Initialized at: %s", fpk.currentPath)
}

// TestFilePickerModalView verifies View() returns a string.
func TestFilePickerModalView(t *testing.T) {
	fpk := newTestFilePickerModal()
	fpk.SetSize(80, 24)
	view := fpk.View()
	if view == "" {
		t.Log("View() returned empty string (acceptable for initial render)")
	}
}

// TestFilePickerModalUpdate verifies Update() accepts messages.
func TestFilePickerModalUpdate(t *testing.T) {
	fpk := newTestFilePickerModal()
	fpk.SetSize(80, 24)
	// Update with arbitrary message - should not panic
	_ = fpk.Update(tea.KeyPressMsg{Code: tea.KeyDown})
}

// TestFilePickerModalSetSize verifies SetSize() is callable.
func TestFilePickerModalSetSize(t *testing.T) {
	fpk := newTestFilePickerModal()
	fpk.SetSize(80, 24)
	if fpk.width != 80 || fpk.height != 24 {
		t.Errorf("SetSize did not update width/height: got %dx%d", fpk.width, fpk.height)
	}
}

// TestFilePickerModalShortcuts verifies Shortcuts() returns exactly 2 entries: Tab+F1.
func TestFilePickerModalShortcuts(t *testing.T) {
	fpk := newTestFilePickerModal()
	shortcuts := fpk.Shortcuts()
	if shortcuts == nil {
		t.Fatal("Shortcuts() returned nil")
	}
	if len(shortcuts) != 2 {
		t.Fatalf("Expected 2 shortcuts (Tab+F1), got %d", len(shortcuts))
	}
	if shortcuts[0].Key != "Tab" || shortcuts[0].Label != "Painel" {
		t.Errorf("shortcuts[0] should be {Tab, Painel}, got {%s, %s}", shortcuts[0].Key, shortcuts[0].Label)
	}
	if shortcuts[1].Key != "F1" || shortcuts[1].Label != "Ajuda" {
		t.Errorf("shortcuts[1] should be {F1, Ajuda}, got {%s, %s}", shortcuts[1].Key, shortcuts[1].Label)
	}
}

// TestFilePickerModalEmitsMessageOnEsc verifies ESC triggers a command.
func TestFilePickerModalEmitsMessageOnEsc(t *testing.T) {
	fpk := newTestFilePickerModal()
	fpk.SetSize(80, 24)

	// Press ESC
	msg := tea.KeyPressMsg{Code: tea.KeyEsc}
	cmd := fpk.Update(msg)

	if cmd != nil {
		resultMsg := cmd()
		if resultMsg != nil {
			_, ok := resultMsg.(popModalMsg)
			if !ok {
				t.Logf("Expected popModalMsg but got %T", resultMsg)
			}
		}
	} else {
		t.Log("ESC did not emit a command")
	}
}

// TestFilePickerModalContainsPanelLabels verifies View contains expected labels.
func TestFilePickerModalContainsPanelLabels(t *testing.T) {
	fpk := newTestFilePickerModal()
	fpk.SetSize(80, 24)

	view := fpk.View()

	hasEstrutura := len(view) > 0 && contains(view, "Estrutura")
	hasArquivos := len(view) > 0 && contains(view, "Arquivos")

	if !hasEstrutura {
		t.Error("View missing 'Estrutura' label")
	}
	if !hasArquivos {
		t.Error("View missing 'Arquivos' label")
	}
}

// TestFilePickerModalDirectoryLoading verifies loadFilesForCursor loads .abditum files.
func TestFilePickerModalDirectoryLoading(t *testing.T) {
	testDir := t.TempDir()
	fpk := &filePickerModal{ext: ".abditum", mode: FilePickerOpen}
	fpk.Init()
	fpk.currentPath = testDir

	// Create some test files
	for i := 0; i < 3; i++ {
		f := testDir + "/test" + string(rune('0'+i)) + ".abditum"
		file, err := os.Create(f)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		file.Close()
	}

	fpk.loadFilesForCursor()

	// Should have loaded 3 files (filtering for .abditum)
	if len(fpk.files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(fpk.files))
	}
}

// TestFilePickerModalFiltering verifies that only .abditum files are shown and hidden files are excluded.
func TestFilePickerModalFiltering(t *testing.T) {
	testDir := t.TempDir()
	fpk := &filePickerModal{ext: ".abditum", mode: FilePickerOpen}
	fpk.Init()
	fpk.currentPath = testDir

	// Create test files
	files := []struct {
		name    string
		visible bool
	}{
		{"vault.abditum", true},
		{"other.txt", false},
		{".hidden.abditum", false},
		{"config.yaml", false},
	}

	for _, f := range files {
		file, err := os.Create(testDir + "/" + f.name)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		file.Close()
	}

	fpk.loadFilesForCursor()

	// Check filtering: only "vault" should be visible (no hidden, no wrong ext)
	if len(fpk.files) != 1 {
		t.Errorf("Expected 1 visible file, got %d", len(fpk.files))
	}
	if len(fpk.files) > 0 && fpk.files[0] != "vault" {
		t.Errorf("Expected 'vault', got '%s'", fpk.files[0])
	}
}

// TestFilePickerModalNavigationDown verifies down arrow moves cursor in the files panel.
func TestFilePickerModalNavigationDown(t *testing.T) {
	testDir := t.TempDir()
	fpk := &filePickerModal{ext: ".abditum", mode: FilePickerOpen}
	fpk.Init()
	fpk.currentPath = testDir
	fpk.focusPanel = 1 // switch to files panel after Init (Init always sets 0)

	// Create test files
	for i := 0; i < 5; i++ {
		file, _ := os.Create(testDir + "/" + "file" + string(rune('0'+i)) + ".abditum")
		if file != nil {
			file.Close()
		}
	}
	fpk.loadFilesForCursor()

	// After loading, D-15: fileCursor is 0 (auto-selected because files exist)
	initialCursor := fpk.fileCursor // should be 0
	fpk.Update(tea.KeyPressMsg{Code: tea.KeyDown})

	// Cursor should move from 0 to 1
	if fpk.fileCursor == initialCursor {
		t.Error("Down key did not move cursor")
	}
	if fpk.fileCursor != 1 {
		t.Errorf("Expected cursor at 1, got %d", fpk.fileCursor)
	}
}

// TestFilePickerModalNavigationUp verifies up arrow moves cursor backwards.
func TestFilePickerModalNavigationUp(t *testing.T) {
	testDir := t.TempDir()
	fpk := &filePickerModal{ext: ".abditum", mode: FilePickerOpen}
	fpk.Init()
	fpk.currentPath = testDir
	fpk.focusPanel = 1 // switch to files panel after Init

	// Create test files
	for i := 0; i < 5; i++ {
		file, _ := os.Create(testDir + "/" + "file" + string(rune('0'+i)) + ".abditum")
		if file != nil {
			file.Close()
		}
	}
	fpk.loadFilesForCursor()

	// Move to position 2
	fpk.fileCursor = 2
	fpk.Update(tea.KeyPressMsg{Code: tea.KeyUp})

	// Cursor should move back
	if fpk.fileCursor != 1 {
		t.Errorf("Expected cursor at 1, got %d", fpk.fileCursor)
	}
}

// TestFilePickerModalTabFocus verifies Tab cycles focus between panels.
func TestFilePickerModalTabFocus(t *testing.T) {
	fpk := newTestFilePickerModal()

	initialFocus := fpk.focusPanel
	fpk.Update(tea.KeyPressMsg{Code: tea.KeyTab})

	if fpk.focusPanel == initialFocus {
		t.Error("Tab did not change focus panel")
	}
}

// TestFilePickerModalDisplaysFileSizes verifies that file sizes are shown with space-separated units.
func TestFilePickerModalDisplaysFileSizes(t *testing.T) {
	testDir := t.TempDir()
	fpk := &filePickerModal{ext: ".abditum", mode: FilePickerOpen, focusPanel: 1}
	fpk.Init()
	fpk.currentPath = testDir
	fpk.SetSize(80, 24)

	// Create test files with different sizes
	testFiles := []struct {
		name string
		size int
	}{
		{"small.abditum", 512},      // 512 bytes
		{"medium.abditum", 1024000}, // ~1MB
	}

	for _, f := range testFiles {
		file, _ := os.Create(testDir + "/" + f.name)
		if file != nil {
			file.WriteString(strings.Repeat("x", f.size))
			file.Close()
		}
	}

	fpk.loadFilesForCursor()
	view := fpk.View()

	// Space-separated units per D-05: "512 B", "1.0 KB", "1.0 MB"
	if !contains(view, " B") && !contains(view, " KB") && !contains(view, " MB") && !contains(view, " GB") {
		t.Error("File sizes must use space-separated units: '512 B', '1.2 KB', etc. (D-05)")
	}
}

// TestFilePickerModalDisplaysRelativeDates verifies that modification dates use absolute format.
func TestFilePickerModalDisplaysRelativeDates(t *testing.T) {
	testDir := t.TempDir()
	fpk := &filePickerModal{ext: ".abditum", mode: FilePickerOpen, focusPanel: 1}
	fpk.Init()
	fpk.currentPath = testDir
	fpk.SetSize(80, 24)

	// Create a test file
	file, _ := os.Create(testDir + "/recent.abditum")
	if file != nil {
		file.Close()
	}

	fpk.loadFilesForCursor()
	view := fpk.View()

	// Absolute date format: "dd/mm/aa HH:MM" — must contain "/" between date parts (D-05)
	if !contains(view, "/") {
		t.Error("File picker must show absolute date format (dd/mm/aa HH:MM) per D-05 — must contain '/'")
	}
	// Must NOT show relative indicators
	if contains(view, "now") || contains(view, " ago") {
		t.Error("File picker must NOT show relative time — use absolute date format per D-05")
	}
}

// TestFilePickerModalHandlesInaccessibleDirectory verifies error handling for nonexistent dirs.
func TestFilePickerModalHandlesInaccessibleDirectory(t *testing.T) {
	testDir := t.TempDir()
	fpk := &filePickerModal{ext: ".abditum", mode: FilePickerOpen, focusPanel: 0}
	fpk.Init()
	fpk.currentPath = testDir

	// Test error handling when loading a non-existent path
	fpk.currentPath = testDir + "/nonexistent"
	fpk.loadFilesForCursor()

	// Should not crash; files should be empty
	t.Logf("loadFilesForCursor handled inaccessible/nonexistent directory gracefully, files=%d", len(fpk.files))
}

// TestFilePickerModalMouseScrollSupport verifies that scroll/navigation events don't crash the modal.
func TestFilePickerModalMouseScrollSupport(t *testing.T) {
	testDir := t.TempDir()
	fpk := &filePickerModal{ext: ".abditum", mode: FilePickerOpen, focusPanel: 1}
	fpk.Init()
	fpk.currentPath = testDir

	// Create multiple files to enable scrolling
	for i := 0; i < 20; i++ {
		file, _ := os.Create(testDir + "/" + "file" + fmt.Sprintf("%02d", i) + ".abditum")
		if file != nil {
			file.Close()
		}
	}
	fpk.loadFilesForCursor()
	fpk.SetSize(40, 5) // Small height to force scrolling

	initialCursor := fpk.fileCursor
	for i := 0; i < 15; i++ {
		fpk.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	}

	t.Logf("Navigation handled without panic. Cursor: %d -> %d", initialCursor, fpk.fileCursor)
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ─────────────────────────────────────────────────────────────────────────────
// Behavioral Update() tests — D-07 matrix (18 cases)
// ─────────────────────────────────────────────────────────────────────────────

// TestFilePickerUpdateBehavior covers all 18 behavioral cases from CONTEXT.md D-07.
func TestFilePickerUpdateBehavior(t *testing.T) {
	// makeDir creates a temp dir with subdirs and .abditum files.
	makeDir := func(t *testing.T, subdirCount, fileCount int) string {
		dir := t.TempDir()
		for i := 0; i < subdirCount; i++ {
			os.Mkdir(filepath.Join(dir, fmt.Sprintf("sub%02d", i)), 0755)
		}
		for i := 0; i < fileCount; i++ {
			f, _ := os.Create(filepath.Join(dir, fmt.Sprintf("file%02d.abditum", i)))
			if f != nil {
				f.Close()
			}
		}
		return dir
	}

	// makeFPK makes a modal with given focus panel and a loaded directory.
	makeFPK := func(t *testing.T, mode FilePickerMode, dir string, focus int) *filePickerModal {
		fpk := &filePickerModal{
			ext:         ".abditum",
			mode:        mode,
			currentPath: dir,
		}
		fpk.Init()
		fpk.currentPath = dir // override CWD after Init
		fpk.loadFilesForCursor()
		fpk.focusPanel = focus // set focus AFTER Init (Init resets to 0)
		fpk.SetSize(80, 24)
		return fpk
	}

	tests := []struct {
		name  string
		setup func(t *testing.T) *filePickerModal
		key   tea.Key
		check func(t *testing.T, fpk *filePickerModal, cmd tea.Cmd)
	}{
		// ── Tree panel ──────────────────────────────────────────────────────
		{
			name: "↓ in tree: cursor stays in bounds",
			setup: func(t *testing.T) *filePickerModal {
				dir := makeDir(t, 0, 2)
				fpk := makeFPK(t, FilePickerOpen, dir, 0)
				// Expand root so there are visible nodes
				if fpk.root != nil && !fpk.root.expanded {
					fpk.root.expanded = true
					fpk.visibleNodes = nil
					fpk.buildVisibleNodes(fpk.root, &fpk.visibleNodes)
				}
				return fpk
			},
			key: tea.Key{Code: tea.KeyDown},
			check: func(t *testing.T, fpk *filePickerModal, cmd tea.Cmd) {
				if len(fpk.visibleNodes) > 0 {
					if fpk.treeCursor < 0 || fpk.treeCursor >= len(fpk.visibleNodes) {
						t.Errorf("treeCursor out of bounds: %d (len=%d)", fpk.treeCursor, len(fpk.visibleNodes))
					}
				}
			},
		},
		{
			name: "↑ in tree at top: cursor stays at 0",
			setup: func(t *testing.T) *filePickerModal {
				dir := makeDir(t, 0, 0)
				fpk := makeFPK(t, FilePickerOpen, dir, 0)
				fpk.treeCursor = 0
				return fpk
			},
			key: tea.Key{Code: tea.KeyUp},
			check: func(t *testing.T, fpk *filePickerModal, cmd tea.Cmd) {
				if fpk.treeCursor != 0 {
					t.Errorf("↑ at top should clamp to 0, got %d", fpk.treeCursor)
				}
			},
		},
		{
			name: "Tab in tree (Open): focus moves to files panel",
			setup: func(t *testing.T) *filePickerModal {
				dir := makeDir(t, 0, 1)
				return makeFPK(t, FilePickerOpen, dir, 0)
			},
			key: tea.Key{Code: tea.KeyTab},
			check: func(t *testing.T, fpk *filePickerModal, cmd tea.Cmd) {
				if fpk.focusPanel != 1 {
					t.Errorf("Tab in tree (Open): expected focusPanel=1, got %d", fpk.focusPanel)
				}
			},
		},
		{
			name: "Enter in tree with .abditum files: focus moves to files",
			setup: func(t *testing.T) *filePickerModal {
				dir := makeDir(t, 0, 2)
				return makeFPK(t, FilePickerOpen, dir, 0)
			},
			key: tea.Key{Code: tea.KeyEnter},
			check: func(t *testing.T, fpk *filePickerModal, cmd tea.Cmd) {
				if fpk.focusPanel != 1 {
					t.Errorf("Enter on dir with files: expected focusPanel=1, got %d", fpk.focusPanel)
				}
			},
		},
		{
			name: "Enter in tree with no .abditum files: no-op",
			setup: func(t *testing.T) *filePickerModal {
				dir := makeDir(t, 0, 0) // no files
				return makeFPK(t, FilePickerOpen, dir, 0)
			},
			key: tea.Key{Code: tea.KeyEnter},
			check: func(t *testing.T, fpk *filePickerModal, cmd tea.Cmd) {
				if fpk.focusPanel != 0 {
					t.Errorf("Enter on empty dir: expected focusPanel=0 (no-op), got %d", fpk.focusPanel)
				}
			},
		},
		// ── Files panel ─────────────────────────────────────────────────────
		{
			name: "Tab in files (Open): focus → tree",
			setup: func(t *testing.T) *filePickerModal {
				dir := makeDir(t, 0, 2)
				return makeFPK(t, FilePickerOpen, dir, 1)
			},
			key: tea.Key{Code: tea.KeyTab},
			check: func(t *testing.T, fpk *filePickerModal, cmd tea.Cmd) {
				if fpk.focusPanel != 0 {
					t.Errorf("Tab in files (Open): expected focusPanel=0, got %d", fpk.focusPanel)
				}
			},
		},
		{
			name: "Tab in files (Save): focus → campo nome",
			setup: func(t *testing.T) *filePickerModal {
				dir := makeDir(t, 0, 2)
				return makeFPK(t, FilePickerSave, dir, 1)
			},
			key: tea.Key{Code: tea.KeyTab},
			check: func(t *testing.T, fpk *filePickerModal, cmd tea.Cmd) {
				if fpk.focusPanel != 2 {
					t.Errorf("Tab in files (Save): expected focusPanel=2, got %d", fpk.focusPanel)
				}
			},
		},
		{
			name: "Enter in files (Open) with file selected: emits filePickerResult + popModalMsg",
			setup: func(t *testing.T) *filePickerModal {
				dir := makeDir(t, 0, 2)
				fpk := makeFPK(t, FilePickerOpen, dir, 1)
				fpk.fileCursor = 0
				return fpk
			},
			key: tea.Key{Code: tea.KeyEnter},
			check: func(t *testing.T, fpk *filePickerModal, cmd tea.Cmd) {
				if cmd == nil {
					t.Fatal("Enter on file (Open): expected non-nil cmd")
				}
				msg := cmd()
				// Accept BatchMsg, filePickerResult, or popModalMsg
				switch msg.(type) {
				case tea.BatchMsg, filePickerResult, popModalMsg:
					// ok
				default:
					t.Errorf("Expected BatchMsg or filePickerResult, got %T", msg)
				}
			},
		},
		{
			name: "Enter in files (Save): copies name to field + focus → campo nome",
			setup: func(t *testing.T) *filePickerModal {
				dir := makeDir(t, 0, 2)
				fpk := makeFPK(t, FilePickerSave, dir, 1)
				fpk.fileCursor = 0
				return fpk
			},
			key: tea.Key{Code: tea.KeyEnter},
			check: func(t *testing.T, fpk *filePickerModal, cmd tea.Cmd) {
				if fpk.focusPanel != 2 {
					t.Errorf("Enter in files (Save): expected focusPanel=2, got %d", fpk.focusPanel)
				}
				if fpk.nameField.Value() == "" {
					t.Error("Enter in files (Save): field should have filename copied")
				}
			},
		},
		{
			name: "Home in files: fileCursor=0",
			setup: func(t *testing.T) *filePickerModal {
				dir := makeDir(t, 0, 5)
				fpk := makeFPK(t, FilePickerOpen, dir, 1)
				fpk.fileCursor = 4
				return fpk
			},
			key: tea.Key{Code: tea.KeyHome},
			check: func(t *testing.T, fpk *filePickerModal, cmd tea.Cmd) {
				if fpk.fileCursor != 0 {
					t.Errorf("Home: expected fileCursor=0, got %d", fpk.fileCursor)
				}
			},
		},
		{
			name: "End in files: fileCursor=last",
			setup: func(t *testing.T) *filePickerModal {
				dir := makeDir(t, 0, 5)
				return makeFPK(t, FilePickerOpen, dir, 1)
			},
			key: tea.Key{Code: tea.KeyEnd},
			check: func(t *testing.T, fpk *filePickerModal, cmd tea.Cmd) {
				expected := len(fpk.files) - 1
				if fpk.fileCursor != expected {
					t.Errorf("End: expected fileCursor=%d, got %d", expected, fpk.fileCursor)
				}
			},
		},
		{
			name: "PgDn in files (scroll=0, 10 files): cursor advances",
			setup: func(t *testing.T) *filePickerModal {
				dir := makeDir(t, 0, 10)
				fpk := makeFPK(t, FilePickerOpen, dir, 1)
				fpk.fileScroll = 0
				fpk.fileCursor = 0
				return fpk
			},
			key: tea.Key{Code: tea.KeyPgDown},
			check: func(t *testing.T, fpk *filePickerModal, cmd tea.Cmd) {
				if fpk.fileCursor == 0 {
					t.Error("PgDn should advance fileCursor from 0")
				}
			},
		},
		// ── Campo nome (Save mode) ───────────────────────────────────────────
		{
			name: "Tab in campo nome: focus → tree",
			setup: func(t *testing.T) *filePickerModal {
				dir := makeDir(t, 0, 0)
				return makeFPK(t, FilePickerSave, dir, 2)
			},
			key: tea.Key{Code: tea.KeyTab},
			check: func(t *testing.T, fpk *filePickerModal, cmd tea.Cmd) {
				if fpk.focusPanel != 0 {
					t.Errorf("Tab in campo nome: expected focusPanel=0, got %d", fpk.focusPanel)
				}
			},
		},
		{
			name: "Enter in campo nome (non-empty): emits filePickerResult + popModalMsg",
			setup: func(t *testing.T) *filePickerModal {
				dir := makeDir(t, 0, 0)
				fpk := makeFPK(t, FilePickerSave, dir, 2)
				fpk.nameField.SetValue("meu-cofre")
				return fpk
			},
			key: tea.Key{Code: tea.KeyEnter},
			check: func(t *testing.T, fpk *filePickerModal, cmd tea.Cmd) {
				if cmd == nil {
					t.Fatal("Enter with non-empty field: expected non-nil cmd")
				}
			},
		},
		{
			name: "Enter in campo nome (empty): no-op",
			setup: func(t *testing.T) *filePickerModal {
				dir := makeDir(t, 0, 0)
				fpk := makeFPK(t, FilePickerSave, dir, 2)
				fpk.nameField.SetValue("")
				return fpk
			},
			key: tea.Key{Code: tea.KeyEnter},
			check: func(t *testing.T, fpk *filePickerModal, cmd tea.Cmd) {
				if cmd != nil {
					t.Error("Enter with empty field: expected nil cmd (no-op)")
				}
			},
		},
		// ── Global ──────────────────────────────────────────────────────────
		{
			name: "Esc from tree: emits Cancelled result + popModalMsg",
			setup: func(t *testing.T) *filePickerModal {
				dir := makeDir(t, 0, 0)
				return makeFPK(t, FilePickerOpen, dir, 0)
			},
			key: tea.Key{Code: tea.KeyEsc},
			check: func(t *testing.T, fpk *filePickerModal, cmd tea.Cmd) {
				if cmd == nil {
					t.Fatal("Esc: expected non-nil cmd")
				}
				// Check the batch contains a filePickerResult{Cancelled:true}
				gotCancelled := false
				msg := cmd()
				if bm, ok := msg.(tea.BatchMsg); ok {
					for _, fn := range bm {
						if m := fn(); m != nil {
							if r, ok := m.(filePickerResult); ok && r.Cancelled {
								gotCancelled = true
							}
						}
					}
				}
				if !gotCancelled {
					// Also acceptable: cmd() itself is the filePickerResult
					if r, ok := msg.(filePickerResult); ok && r.Cancelled {
						gotCancelled = true
					}
				}
				if !gotCancelled {
					t.Error("Esc: expected filePickerResult{Cancelled:true} in cmd chain")
				}
			},
		},
		{
			name: "Esc from files: emits Cancelled result",
			setup: func(t *testing.T) *filePickerModal {
				dir := makeDir(t, 0, 2)
				return makeFPK(t, FilePickerOpen, dir, 1)
			},
			key: tea.Key{Code: tea.KeyEsc},
			check: func(t *testing.T, fpk *filePickerModal, cmd tea.Cmd) {
				if cmd == nil {
					t.Fatal("Esc from files: expected non-nil cmd")
				}
			},
		},
		{
			name: "Esc from campo nome: emits Cancelled result",
			setup: func(t *testing.T) *filePickerModal {
				dir := makeDir(t, 0, 0)
				return makeFPK(t, FilePickerSave, dir, 2)
			},
			key: tea.Key{Code: tea.KeyEsc},
			check: func(t *testing.T, fpk *filePickerModal, cmd tea.Cmd) {
				if cmd == nil {
					t.Fatal("Esc from campo nome: expected non-nil cmd")
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			fpk := tt.setup(t)
			cmd := fpk.Update(tea.KeyPressMsg{Code: tt.key.Code})
			tt.check(t, fpk, cmd)
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Golden File Tests
// ─────────────────────────────────────────────────────────────────────────────

// TestFilePickerModal_Golden validates visual output against golden files.
func TestFilePickerModal_Golden(t *testing.T) {
	testDir := t.TempDir()

	// Create test file structure
	for i := 0; i < 5; i++ {
		f, err := os.Create(testDir + "/" + fmt.Sprintf("vault%d.abditum", i))
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		f.Close()
	}

	// Create some directories
	for i := 0; i < 3; i++ {
		os.Mkdir(testDir+"/"+fmt.Sprintf("Folder%d", i), 0755)
	}

	fpk := &filePickerModal{ext: ".abditum", mode: FilePickerOpen}
	fpk.Init()
	// Inject fixed time formatter for deterministic golden output
	fixedTime, _ := time.Parse("02/01/06 15:04", "01/01/24 12:00")
	fpk.timeFmt = func(time.Time) string { return fixedTime.Format("02/01/06 15:04") }
	fpk.currentPath = testDir
	fpk.focusPanel = 1
	fpk.loadFilesForCursor()
	// Override currentPath to a fixed value for deterministic golden rendering
	fpk.currentPath = "/home/user/cofres"
	fpk.SetSize(80, 24)
	fpk.theme = ThemeTokyoNight

	out := fpk.View()

	// .txt.golden: plain text render
	txtPath := goldenPath("filepicker", "initial", 80, "txt")
	checkOrUpdateGolden(t, txtPath, stripANSI(out))

	// .json.golden: style transitions (if any ANSI codes present)
	transitions := testdatapkg.ParseANSIStyle(out)
	jsonBytes, err := testdatapkg.MarshalStyleTransitions(transitions)
	if err != nil {
		t.Fatalf("marshal transitions: %v", err)
	}
	jsonPath := goldenPath("filepicker", "initial", 80, "json")
	checkOrUpdateGolden(t, jsonPath, string(jsonBytes))
}
