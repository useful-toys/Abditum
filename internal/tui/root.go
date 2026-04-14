package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type RootModel struct {
	width        int
	height       int
	theme        *Theme
	workArea     WorkArea
	activeView   ChildView
	modals       []ModalView
	lastActionAt time.Time
	version      string
}

func (r *RootModel) View() tea.View {
	if r.width == 0 || r.height == 0 {
		return tea.NewView("Aguarde...")
	}

	base := r.activeView.Render(r.width, r.height, *r.theme)

	if len(r.modals) == 0 {
		return tea.NewView(base)
	}

	top := r.modals[len(r.modals)-1]
	modalView := top.Render(r.width, r.height, *r.theme)

	workH := r.height - 4 // 2 header + 1 msg + 1 action
	modalContent := lipgloss.Place(r.width, workH, lipgloss.Center, lipgloss.Center, modalView)
	v := tea.NewView(modalContent)
	v.AltScreen = true
	v.BackgroundColor = lipgloss.Color(r.theme.Surface.Base)
	return v
}

func (r *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.width = msg.Width
		r.height = msg.Height
		return r, nil

	case OpenModalMsg:
		r.modals = append(r.modals, msg.Modal)
		return r, nil

	case CloseModalMsg:
		if len(r.modals) > 0 {
			r.modals = r.modals[:len(r.modals)-1]
		}
		return r, nil

	case ModalReadyMsg:
		if len(r.modals) > 1 {
			parent := r.modals[len(r.modals)-2]
			return r, parent.Update(msg)
		}
		return r, r.activeView.Update(msg)
	}

	if len(r.modals) > 0 {
		top := len(r.modals) - 1
		return r, r.modals[top].Update(msg)
	}

	return r, r.activeView.Update(msg)
}

func (r *RootModel) Init() tea.Cmd {
	return nil
}

type RootModelOption func(*RootModel)

func WithVersion(version string) RootModelOption {
	return func(m *RootModel) {
		m.version = version
	}
}

func NewRootModel(opts ...RootModelOption) *RootModel {
	m := &RootModel{
		theme:      TokyoNight,
		workArea:   WorkAreaWelcome,
		activeView: nil,
		version:    "dev",
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}
