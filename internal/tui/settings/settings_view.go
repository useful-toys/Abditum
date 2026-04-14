package settings

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui"
)

type SettingsView struct{}

func NewSettingsView() *SettingsView {
	return &SettingsView{}
}

func (v *SettingsView) Render(height, width int, theme tui.Theme) string {
	content := "Settings"
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Base))
	return style.Render(content)
}

func (v *SettingsView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }
func (v *SettingsView) HandleEvent(event any)            {}
func (v *SettingsView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }
func (v *SettingsView) Update(msg tea.Msg) tea.Cmd       { return nil }
