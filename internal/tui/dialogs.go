package tui

import tea "charm.land/bubbletea/v2"

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

// PasswordEntry creates a stub password-entry modal.
func PasswordEntry(title string) tea.Cmd {
	return func() tea.Msg {
		return pushModalMsg{modal: &passwordEntryModal{title: title}}
	}
}

// PasswordCreate creates a stub password-creation modal.
func PasswordCreate(title string) tea.Cmd {
	return func() tea.Msg {
		return pushModalMsg{modal: &passwordCreateModal{title: title}}
	}
}

// FilePicker creates a stub file-picker modal.
func FilePicker(title string, mode FilePickerMode, extension string) tea.Cmd {
	return func() tea.Msg {
		return pushModalMsg{modal: &filePickerModal{title: title, mode: mode, ext: extension}}
	}
}

// --- Stub modal implementations for Phase 5.1 ---

type passwordEntryModal struct{ title string }

func (m *passwordEntryModal) Update(msg tea.Msg) tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return popModalMsg{} },
		func() tea.Msg { return passwordEntryResult{Cancelled: true} },
	)
}
func (m *passwordEntryModal) View() string          { return "[PasswordEntry stub - Phase 6]" }
func (m *passwordEntryModal) Shortcuts() []Shortcut { return nil }
func (m *passwordEntryModal) SetSize(w, h int)      {}

type passwordCreateModal struct{ title string }

func (m *passwordCreateModal) Update(msg tea.Msg) tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return popModalMsg{} },
		func() tea.Msg { return passwordCreateResult{Cancelled: true} },
	)
}
func (m *passwordCreateModal) View() string          { return "[PasswordCreate stub - Phase 6]" }
func (m *passwordCreateModal) Shortcuts() []Shortcut { return nil }
func (m *passwordCreateModal) SetSize(w, h int)      {}

type filePickerModal struct {
	title string
	mode  FilePickerMode
	ext   string
}

func (m *filePickerModal) Update(msg tea.Msg) tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return popModalMsg{} },
		func() tea.Msg { return filePickerResult{Cancelled: true} },
	)
}
func (m *filePickerModal) View() string          { return "[FilePicker stub - Phase 6]" }
func (m *filePickerModal) Shortcuts() []Shortcut { return nil }
func (m *filePickerModal) SetSize(w, h int)      {}

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
