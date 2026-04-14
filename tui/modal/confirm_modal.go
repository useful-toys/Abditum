package modal

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/tui"
)

type ConfirmModal struct {
	title   string
	message string
	options []ModalOption
}

func NewConfirmModal(title, message string) *ConfirmModal {
	return &ConfirmModal{
		title:   title,
		message: message,
		options: []ModalOption{
			{
				Keys:   []string{"Enter"},
				Label:  "Confirmar",
				Intent: IntentConfirm,
				Action: func() tea.Cmd {
					return func() tea.Msg { return tui.ModalReadyMsg{} }
				},
			},
			{
				Keys:   []string{"Esc"},
				Label:  "Cancelar",
				Intent: IntentCancel,
				Action: func() tea.Cmd {
					return func() tea.Msg { return tui.CloseModalMsg{} }
				},
			},
		},
	}
}

func (m *ConfirmModal) Render(maxHeight, maxWidth int, theme tui.Theme) string {
	content := m.title + "\n\n" + m.message
	style := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(lipgloss.Color(theme.Border.Default)).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Raised))
	return style.Render(content)
}

func (m *ConfirmModal) HandleKey(msg tea.KeyMsg) tea.Cmd {
	for _, opt := range m.options {
		for _, key := range opt.Keys {
			if msg.String() == key {
				return opt.Action()
			}
		}
	}
	return nil
}
