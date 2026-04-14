package welcome

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/tui"
)

type ViewState int

const (
	StateNormal ViewState = iota
	StateAwaitingModal
)

type WelcomeView struct {
	state ViewState
}

func NewWelcomeView() *WelcomeView {
	return &WelcomeView{state: StateNormal}
}

func (v *WelcomeView) ID() string {
	return "welcome"
}

func (v *WelcomeView) Render(height, width int, theme tui.Theme) string {
	content := "Welcome"
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Base))
	return style.Render(content)
}

func (v *WelcomeView) HandleKey(msg tea.KeyMsg) tea.Cmd {
	return nil
}

func (v *WelcomeView) HandleEvent(event any) {}

func (v *WelcomeView) HandleTeaMsg(msg tea.Msg) {}
