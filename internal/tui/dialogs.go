package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// DialogType is the semantic type of a question or confirm dialog.
type DialogType int

const (
	DialogQuestion DialogType = iota // neutral decision
	DialogAlert                      // destructive or irreversible action
	DialogInfo                       // information requiring explicit acknowledgement
)

// DialogOption is a selectable choice in an Ask or Confirm dialog.
type DialogOption struct {
	Label   string
	Keys    []string
	Cmd     tea.Cmd
	Default bool
	Cancel  bool
}

// FilePickerMode controls which items the file picker can select.
type FilePickerMode int

const (
	FilePickerFile      FilePickerMode = iota // files only
	FilePickerDirectory                       // directories only
	FilePickerAny                             // files and directories
)

// Message creates an informational dialog dismissed via ENTER or ESC.
func Message(dtype DialogType, title, message string) tea.Cmd {
	return func() tea.Msg {
		m := newDialogModal(dtype, title, message, nil)
		return pushModalMsg{modal: m}
	}
}

// Confirm creates a yes/no confirmation dialog.
func Confirm(dtype DialogType, title, message string, onYes, onNo tea.Cmd) tea.Cmd {
	options := []DialogOption{
		{Label: "Sim", Keys: []string{"s", "y"}, Cmd: onYes, Default: true},
		{Label: "Nao", Keys: []string{"n"}, Cmd: onNo, Cancel: true},
	}
	return func() tea.Msg {
		m := newDialogModal(dtype, title, message, options)
		return pushModalMsg{modal: m}
	}
}

// PasswordEntry creates a password-entry modal.
func PasswordEntry(title string) tea.Cmd {
	return func() tea.Msg {
		m := &passwordEntryModal{title: title}
		m.theme = ThemeTokyoNight
		m.messages = NewMessageManager()
		cmd := m.Init()
		return tea.Batch(
			func() tea.Msg { return pushModalMsg{modal: m} },
			cmd,
		)
	}
}

// PasswordCreate creates a password-creation modal.
func PasswordCreate(title string) tea.Cmd {
	return func() tea.Msg {
		m := &passwordCreateModal{title: title}
		m.theme = ThemeTokyoNight
		m.messages = NewMessageManager()
		cmd := m.Init()
		return tea.Batch(
			func() tea.Msg { return pushModalMsg{modal: m} },
			cmd,
		)
	}
}

// NewRecognitionError creates an error recognition dialog (acknowledgement-only).
// Used for unrecoverable errors like invalid vault files.
func NewRecognitionError(title, text string) tea.Cmd {
	return Acknowledge(SeverityError, title, text, nil)
}

// FilePicker creates a file-picker modal.
func FilePicker(title string, mode FilePickerMode, extension string) tea.Cmd {
	return func() tea.Msg {
		fpk := &filePickerModal{title: title, mode: mode, ext: extension}
		fpk.Init()
		return pushModalMsg{modal: fpk}
	}
}

type filePickerModal struct {
	title       string
	mode        FilePickerMode
	ext         string
	currentPath string
	files       []string      // filtered files without extension
	fileInfos   []os.FileInfo // metadata for files (size, mod time)
	directories []string      // subdirectories
	fileCursor  int           // cursor position in files list
	treeCursor  int           // cursor position in directories list
	fileScroll  int           // scroll offset for files
	treeScroll  int           // scroll offset for tree
	width       int
	height      int
	theme       *Theme
	focusPanel  int // 0 = tree, 1 = files
}

func (m *filePickerModal) Init() tea.Cmd {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = os.Getenv("HOME")
		if cwd == "" {
			cwd = "/"
		}
	}
	m.currentPath = cwd
	m.theme = ThemeTokyoNight
	m.focusPanel = 1 // start in files panel
	m.loadDirectory()
	return nil
}

// loadDirectory reads the current directory and filters for .abditum files.
func (m *filePickerModal) loadDirectory() {
	m.files = []string{}
	m.fileInfos = []os.FileInfo{}
	m.directories = []string{}
	m.fileCursor = 0
	m.treeCursor = 0

	entries, err := os.ReadDir(m.currentPath)
	if err != nil {
		return // silently skip inaccessible directories
	}

	for _, entry := range entries {
		// Skip hidden files
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		if entry.IsDir() {
			m.directories = append(m.directories, entry.Name())
		} else if strings.HasSuffix(entry.Name(), ".abditum") {
			// Store without extension
			name := strings.TrimSuffix(entry.Name(), ".abditum")
			m.files = append(m.files, name)
			// Store file info for metadata display
			info, _ := entry.Info()
			m.fileInfos = append(m.fileInfos, info)
		}
	}
}

func (m *filePickerModal) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Code {
		case tea.KeyEsc:
			// Cancel flow
			return tea.Batch(
				func() tea.Msg { return popModalMsg{} },
				func() tea.Msg { return filePickerResult{Cancelled: true} },
			)
		case tea.KeyEnter:
			// Select file or directory
			if m.focusPanel == 0 && len(m.directories) > 0 && m.treeCursor < len(m.directories) {
				// Navigate into directory
				m.currentPath = filepath.Join(m.currentPath, m.directories[m.treeCursor])
				m.loadDirectory()
				return nil
			} else if m.focusPanel == 1 && len(m.files) > 0 && m.fileCursor < len(m.files) {
				// Select file
				fileName := m.files[m.fileCursor] + ".abditum"
				fullPath := filepath.Join(m.currentPath, fileName)
				return tea.Batch(
					func() tea.Msg { return popModalMsg{} },
					func() tea.Msg { return filePickerResult{Path: fullPath} },
				)
			}
		case tea.KeyDown:
			if m.focusPanel == 0 {
				if len(m.directories) > 0 {
					m.treeCursor = (m.treeCursor + 1) % len(m.directories)
				}
			} else {
				if len(m.files) > 0 {
					m.fileCursor = (m.fileCursor + 1) % len(m.files)
				}
			}
		case tea.KeyUp:
			if m.focusPanel == 0 {
				if len(m.directories) > 0 {
					m.treeCursor = (m.treeCursor - 1 + len(m.directories)) % len(m.directories)
				}
			} else {
				if len(m.files) > 0 {
					m.fileCursor = (m.fileCursor - 1 + len(m.files)) % len(m.files)
				}
			}
		case tea.KeyTab:
			// Cycle focus
			m.focusPanel = (m.focusPanel + 1) % 2
		case tea.KeyLeft:
			if m.focusPanel == 1 {
				m.focusPanel = 0
			}
		case tea.KeyRight:
			if m.focusPanel == 0 {
				m.focusPanel = 1
			}
		}
	}
	return nil
}

func (m *filePickerModal) View() string {
	// Simple two-panel layout
	panelWidth := (m.width - 1) / 2 // subtract 1 for separator

	// Left panel: Estrutura (tree)
	leftPanel := m.renderTreePanel(panelWidth)
	// Right panel: Arquivos (files)
	rightPanel := m.renderFilesPanel(panelWidth)

	// Combine with separator
	lines := []string{}
	leftLines := strings.Split(leftPanel, "\n")
	rightLines := strings.Split(rightPanel, "\n")

	maxLines := len(leftLines)
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}

	for i := 0; i < maxLines; i++ {
		left := ""
		if i < len(leftLines) {
			left = leftLines[i]
		}
		right := ""
		if i < len(rightLines) {
			right = rightLines[i]
		}

		// Pad left to width
		for lipgloss.Width(left) < panelWidth {
			left += " "
		}

		lines = append(lines, left+"│"+right)
	}

	return strings.Join(lines, "\n")
}

// formatFileSize returns a human-readable file size string (e.g., "1.2MB", "512B").
func formatFileSize(sizeBytes int64) string {
	const unit = 1024
	if sizeBytes < unit {
		return fmt.Sprintf("%dB", sizeBytes)
	}
	div, exp := int64(unit), 0
	for n := sizeBytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%s", float64(sizeBytes)/float64(div), []string{"", "K", "M", "G"}[exp]) + "B"
}

// formatRelativeDate returns a relative time string for a file's modification time.
// Examples: "now", "1h", "2d", or "01/01/24" for older files.
func formatRelativeDate(modTime time.Time) string {
	now := time.Now()
	diff := now.Sub(modTime)

	// If modified within the last minute, show "now"
	if diff < time.Minute {
		return "now"
	}
	// If within the last hour
	if diff < time.Hour {
		mins := int(diff.Minutes())
		return fmt.Sprintf("%dm", mins)
	}
	// If within the last day
	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%dh", hours)
	}
	// If within the last week
	if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%dd", days)
	}
	// For older files, show date in MM/DD/YY format
	return modTime.Format("01/02/06")
}

// renderTreePanel renders the left panel with directories.
func (m *filePickerModal) renderTreePanel(width int) string {
	content := []string{"  Estrutura"}

	for i, dir := range m.directories {
		prefix := "  "
		if i == m.treeCursor && m.focusPanel == 0 {
			prefix = "→ "
		}
		line := prefix + dir
		if lipgloss.Width(line) > width {
			line = line[:width]
		}
		content = append(content, line)
	}

	// Pad to height
	for len(content) < m.height {
		content = append(content, "")
	}

	return strings.Join(content[:m.height], "\n")
}

// renderFilesPanel renders the right panel with files.
func (m *filePickerModal) renderFilesPanel(width int) string {
	content := []string{" Arquivos"}

	for i, file := range m.files {
		prefix := "  "
		if i == m.fileCursor && m.focusPanel == 1 {
			prefix = "→ "
		}

		// Build file entry with metadata
		var line string
		if i < len(m.fileInfos) && m.fileInfos[i] != nil {
			// Include file size and relative date
			size := formatFileSize(m.fileInfos[i].Size())
			date := formatRelativeDate(m.fileInfos[i].ModTime())
			line = fmt.Sprintf("%s%-20s %8s  %s", prefix, file, size, date)
		} else {
			line = prefix + file
		}

		// Truncate if too wide
		if lipgloss.Width(line) > width {
			line = line[:width-1] + "…"
		}
		content = append(content, line)
	}

	// Pad to height
	for len(content) < m.height {
		content = append(content, "")
	}

	return strings.Join(content[:m.height], "\n")
}

func (m *filePickerModal) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m *filePickerModal) Shortcuts() []Shortcut {
	return []Shortcut{
		{Key: "↑↓", Label: "Navegar"},
		{Key: "Tab", Label: "Trocar"},
		{Key: "Enter", Label: "Selecionar"},
		{Key: "Esc", Label: "Cancelar"},
	}
}

// newDialogModal creates a modalModel-backed dialog for Message and Confirm.
func newDialogModal(dtype DialogType, title, message string, options []DialogOption) modalView {
	var optLabels []string
	var onSelect func(int) tea.Cmd

	if len(options) > 0 {
		for _, o := range options {
			optLabels = append(optLabels, o.Label)
		}
		onSelect = func(idx int) tea.Cmd {
			if idx < 0 {
				for _, o := range options {
					if o.Cancel {
						return o.Cmd
					}
				}
				return nil
			}
			if idx < len(options) {
				return options[idx].Cmd
			}
			return nil
		}
	}

	return newModal(title, message, optLabels, onSelect)
}
