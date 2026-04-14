package secret

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/tui"
)

type TreeView struct{}

func NewTreeView() *TreeView {
	return &TreeView{}
}

func (v *TreeView) Render(height, width int, theme tui.Theme) string {
	content := "Vault Tree"
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Base))
	return style.Render(content)
}

func (v *TreeView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }
func (v *TreeView) HandleEvent(event any)            {}
func (v *TreeView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }
func (v *TreeView) Update(msg tea.Msg) tea.Cmd       { return nil }
