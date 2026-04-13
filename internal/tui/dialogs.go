package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
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
	m := &passwordEntryModal{title: title}
	m.messages = NewMessageManager()
	cmd := m.Init()
	return tea.Batch(
		func() tea.Msg { return pushModalMsg{modal: m} },
		cmd,
	)
}

// PasswordCreate creates a password-creation modal.
func PasswordCreate(title string) tea.Cmd {
	m := &passwordCreateModal{title: title}
	m.messages = NewMessageManager()
	cmd := m.Init()
	return tea.Batch(
		func() tea.Msg { return pushModalMsg{modal: m} },
		cmd,
	)
}

// NewRecognitionError creates an error recognition dialog (acknowledgement-only).
// Used for unrecoverable errors like invalid vault files.
func NewRecognitionError(title, text string) tea.Cmd {
	return Acknowledge(SeverityError, title, text, nil)
}

// FilePicker creates a file-picker modal.
// messages: the flow's *MessageManager for status hints (D-03).
func FilePicker(title string, mode FilePickerMode, ext string, messages *MessageManager) tea.Cmd {
	fpk := &filePickerModal{
		title:    title,
		mode:     mode,
		ext:      ext,
		messages: messages,
	}
	cmd := fpk.Init()
	return tea.Batch(
		func() tea.Msg { return pushModalMsg{modal: fpk} },
		cmd,
	)
}

// formatFileSize returns a human-readable file size with space before unit: "25.8 MB", "1.2 KB", "512 B".
func formatFileSize(sizeBytes int64) string {
	const unit = 1024
	if sizeBytes < unit {
		return fmt.Sprintf("%d B", sizeBytes)
	}
	div, exp := int64(unit), 0
	for n := sizeBytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"KB", "MB", "GB"}
	return fmt.Sprintf("%.1f %s", float64(sizeBytes)/float64(div), units[exp])
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
