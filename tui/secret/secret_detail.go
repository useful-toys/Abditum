package secret

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/tui"
)

type DetailView struct{}

func NewDetailView() *DetailView {
	return &DetailView{}
}

func (v *DetailView) Render(height, width int, theme tui.Theme) string {
	content := "Secret Detail"
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Base))
	return style.Render(content)
}

func (v *DetailView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }
func (v *DetailView) HandleEvent(event any)            {}
func (v *DetailView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }
func (v *DetailView) Update(msg tea.Msg) tea.Cmd       { return nil }
