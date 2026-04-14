package template

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui"
)

type ListView struct{}

func NewListView() *ListView {
	return &ListView{}
}

func (v *ListView) Render(height, width int, theme tui.Theme) string {
	content := "Template List"
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Base))
	return style.Render(content)
}

func (v *ListView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }
func (v *ListView) HandleEvent(event any)            {}
func (v *ListView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }
func (v *ListView) Update(msg tea.Msg) tea.Cmd       { return nil }
