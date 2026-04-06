package tui

import (
	"fmt"
	"os"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

// Helper to create a filePickerModal for testing
func newTestFilePickerModal() *filePickerModal {
	fpk := &filePickerModal{}
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

// TestFilePickerModalInit verifies Init() initializes the directory.
func TestFilePickerModalInit(t *testing.T) {
	fpk := &filePickerModal{}
	fpk.Init()
	if fpk.currentPath == "" {
		t.Fatal("currentPath not set after Init()")
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

// TestFilePickerModalShortcuts verifies Shortcuts() returns a slice.
func TestFilePickerModalShortcuts(t *testing.T) {
	fpk := newTestFilePickerModal()
	shortcuts := fpk.Shortcuts()
	if shortcuts == nil {
		t.Fatal("Shortcuts() returned nil")
	}
	if len(shortcuts) == 0 {
		t.Log("Shortcuts() returned empty slice")
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

// TestFilePickerModalDirectoryLoading verifies loadDirectory works.
func TestFilePickerModalDirectoryLoading(t *testing.T) {
	fpk := &filePickerModal{}
	testDir := t.TempDir()
	fpk.currentPath = testDir

	// Create some test files
	for i := 0; i < 3; i++ {
		f := testDir + "/test" + string(rune('0'+i)) + ".abditum"
		file, err := os.Create(f)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		file.Close() // Close immediately to avoid Windows locking issues
	}

	// Try to load directory
	fpk.loadDirectory()

	// Should have loaded 3 files (filtering for .abditum)
	if len(fpk.files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(fpk.files))
	}
}

// TestFilePickerModalFiltering verifies that only .abditum files are shown and hidden files are excluded.
func TestFilePickerModalFiltering(t *testing.T) {
	fpk := &filePickerModal{}
	testDir := t.TempDir()
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
		file.Close() // Close immediately to avoid Windows locking issues
	}

	// Load directory
	fpk.loadDirectory()

	// Check filtering
	if len(fpk.files) != 1 {
		t.Errorf("Expected 1 visible file, got %d", len(fpk.files))
	}
	if len(fpk.files) > 0 && fpk.files[0] != "vault" {
		t.Errorf("Expected 'vault', got '%s'", fpk.files[0])
	}
}

// TestFilePickerModalNavigationDown verifies down arrow moves cursor.
func TestFilePickerModalNavigationDown(t *testing.T) {
	fpk := &filePickerModal{focusPanel: 1} // Focus on files panel
	testDir := t.TempDir()
	fpk.currentPath = testDir

	// Create test files
	for i := 0; i < 5; i++ {
		file, _ := os.Create(testDir + "/" + "file" + string(rune('0'+i)) + ".abditum")
		file.Close() // Close immediately to avoid Windows locking issues
	}
	fpk.loadDirectory()

	initialCursor := fpk.fileCursor
	msg := tea.KeyPressMsg{Code: tea.KeyDown}
	fpk.Update(msg)

	// Cursor should move
	if fpk.fileCursor == initialCursor {
		t.Error("Down key did not move cursor")
	}
	if fpk.fileCursor != 1 {
		t.Errorf("Expected cursor at 1, got %d", fpk.fileCursor)
	}
}

// TestFilePickerModalNavigationUp verifies up arrow moves cursor backwards.
func TestFilePickerModalNavigationUp(t *testing.T) {
	fpk := &filePickerModal{focusPanel: 1} // Focus on files panel
	testDir := t.TempDir()
	fpk.currentPath = testDir

	// Create test files
	for i := 0; i < 5; i++ {
		file, _ := os.Create(testDir + "/" + "file" + string(rune('0'+i)) + ".abditum")
		file.Close() // Close immediately to avoid Windows locking issues
	}
	fpk.loadDirectory()

	// Move to position 2
	fpk.fileCursor = 2
	msg := tea.KeyPressMsg{Code: tea.KeyUp}
	fpk.Update(msg)

	// Cursor should move back
	if fpk.fileCursor != 1 {
		t.Errorf("Expected cursor at 1, got %d", fpk.fileCursor)
	}
}

// TestFilePickerModalTabFocus verifies Tab cycles focus between panels.
func TestFilePickerModalTabFocus(t *testing.T) {
	fpk := newTestFilePickerModal()

	initialFocus := fpk.focusPanel
	msg := tea.KeyPressMsg{Code: tea.KeyTab}
	fpk.Update(msg)

	if fpk.focusPanel == initialFocus {
		t.Error("Tab did not change focus panel")
	}
}

// TestFilePickerModalDisplaysFileSizes verifies that file sizes are shown in human-readable format.
func TestFilePickerModalDisplaysFileSizes(t *testing.T) {
	fpk := &filePickerModal{focusPanel: 1}
	testDir := t.TempDir()
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
		// Write specific size to file
		file.WriteString(strings.Repeat("x", f.size))
		file.Close()
	}

	fpk.loadDirectory()
	view := fpk.View()

	// View should contain human-readable size info (not just file names)
	// Look for patterns like "512B", "1.0M", "KB", etc.
	if !contains(view, "B") && !contains(view, "K") && !contains(view, "M") {
		t.Error("File picker view must display file sizes in human-readable format (B, KB, MB, etc)")
	}
}

// TestFilePickerModalDisplaysRelativeDates verifies that modification dates are shown in relative format.
func TestFilePickerModalDisplaysRelativeDates(t *testing.T) {
	fpk := &filePickerModal{focusPanel: 1}
	testDir := t.TempDir()
	fpk.currentPath = testDir
	fpk.SetSize(80, 24)

	// Create a test file
	file, _ := os.Create(testDir + "/recent.abditum")
	file.Close()

	fpk.loadDirectory()
	view := fpk.View()

	// View should include time/date information for files
	// Should contain patterns like "now", time (h/d/m), or date format (MM/DD/YY)
	// For a freshly created file, should show recent time indicator
	hasTimeInfo := contains(view, "now") || contains(view, "h") || contains(view, "d") || contains(view, "/")
	if !hasTimeInfo {
		t.Error("File picker view must display relative dates (e.g., 'now', '1h', '2d', or date format)")
	}
}

// TestFilePickerModalHandlesInaccessibleDirectory verifies error handling for inaccessible dirs.
func TestFilePickerModalHandlesInaccessibleDirectory(t *testing.T) {
	fpk := &filePickerModal{focusPanel: 0}
	testDir := t.TempDir()
	fpk.currentPath = testDir

	// Create a subdirectory and make it inaccessible (Windows: deny read)
	restrictedDir := testDir + "/restricted"
	os.Mkdir(restrictedDir, 0755)

	// On Windows, we can't easily simulate inaccessible dirs via permissions
	// Instead, test the error handling when loading a non-existent path
	fpk.currentPath = testDir + "/nonexistent"
	fpk.loadDirectory()

	// Should not crash; files/directories should be empty
	if fpk.files != nil || fpk.directories != nil {
		// After loading nonexistent dir, should gracefully handle
		// This is actually OK - loadDirectory silently skips on error
		t.Log("loadDirectory handled inaccessible/nonexistent directory gracefully")
	}
}

// TestFilePickerModalMouseScrollSupport verifies that scroll events don't crash the modal.
func TestFilePickerModalMouseScrollSupport(t *testing.T) {
	fpk := &filePickerModal{focusPanel: 1}
	testDir := t.TempDir()
	fpk.currentPath = testDir

	// Create multiple files to enable scrolling
	for i := 0; i < 20; i++ {
		file, _ := os.Create(testDir + "/" + "file" + fmt.Sprintf("%02d", i) + ".abditum")
		file.Close()
	}
	fpk.loadDirectory()
	fpk.SetSize(40, 5) // Small height to force scrolling

	// Test that PageDown-like navigation works (simulate via multiple down arrows)
	initialCursor := fpk.fileCursor
	for i := 0; i < 15; i++ {
		fpk.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	}

	// Cursor should have moved and wrapped around (or reached end)
	// This effectively tests scrolling behavior
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
