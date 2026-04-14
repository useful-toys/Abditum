package tui

import tea "charm.land/bubbletea/v2"

type WorkArea int

const (
	WorkAreaWelcome WorkArea = iota
	WorkAreaSettings
	WorkAreaVault
	WorkAreaTemplates
)

type Theme struct {
	Name     string
	Surface  SurfaceTokens
	Text     TextTokens
	Accent   AccentTokens
	Border   BorderTokens
	Semantic SemanticTokens
	Special  SpecialTokens
}

type SurfaceTokens struct {
	Base   string
	Raised string
	Input  string
}

type TextTokens struct {
	Primary   string
	Secondary string
	Disabled  string
	Link      string
}

type AccentTokens struct {
	Primary   string
	Secondary string
}

type BorderTokens struct {
	Default string
	Focused string
}

type SemanticTokens struct {
	Success string
	Warning string
	Error   string
	Info    string
	Off     string
}

type SpecialTokens struct {
	Muted     string
	Highlight string
	Match     string
}

type ChildView interface {
	ID() string
	Render(height, width int, theme Theme) string
	HandleKey(msg tea.KeyMsg) tea.Cmd
	HandleEvent(event any)
	HandleTeaMsg(msg tea.Msg)
}

type ModalView interface {
	Render(maxHeight, maxWidth int, theme Theme) string
	HandleKey(msg tea.KeyMsg) tea.Cmd
	SetOnComplete(func(any))
}
