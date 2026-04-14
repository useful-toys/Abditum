package tui

import tea "charm.land/bubbletea/v2"

// WorkArea representa o estado da área de trabalho
type WorkArea int

const (
	WorkAreaWelcome WorkArea = iota
	WorkAreaSettings
	WorkAreaVault
	WorkAreaTemplates
)

// Theme define cores, tipografia e símbolos
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

// ChildView interface para componentes da tela principal
type ChildView interface {
	Render(height, width int, theme Theme) string
	HandleKey(msg tea.KeyMsg) tea.Cmd
	HandleEvent(event any)
	HandleTeaMsg(msg tea.Msg)
}

// ModalView interface para modais
type ModalView interface {
	Render(maxHeight, maxWidth int, theme Theme) string
	HandleKey(msg tea.KeyMsg) tea.Cmd
}

// OpenModalMsg empilha um novo modal
type OpenModalMsg struct {
	Modal ModalView
}

// CloseModalMsg desempilha o modal do topo
type CloseModalMsg struct{}

// ModalReadyMsg indica que modal tem resultado
type ModalReadyMsg struct{}

// OpenModal cria um comando para abrir um modal
func OpenModal(modal ModalView) tea.Cmd {
	return func() tea.Msg { return OpenModalMsg{Modal: modal} }
}

// CloseModal cria um comando para fechar o modal do topo
func CloseModal() tea.Cmd {
	return func() tea.Msg { return CloseModalMsg{} }
}

var TokyoNight = &Theme{
	Name:     "Tokyo Night",
	Surface:  SurfaceTokens{"#1a1b26", "#24283b", "#1e1f2e"},
	Text:     TextTokens{"#a9b1d6", "#565f89", "#3b4261", "#7aa2f7"},
	Accent:   AccentTokens{"#7aa2f7", "#bb9af7"},
	Border:   BorderTokens{"#414868", "#7aa2f7"},
	Semantic: SemanticTokens{"#9ece6a", "#e0af68", "#f7768e", "#7dcfff", "#737aa2"},
	Special:  SpecialTokens{"#8690b5", "#283457", "#f7c67a"},
}

var Cyberpunk = &Theme{
	Name:     "Cyberpunk",
	Surface:  SurfaceTokens{"#0a0a1a", "#1a1a2e", "#0e0e22"},
	Text:     TextTokens{"#e0e0ff", "#8888aa", "#444466", "#ff2975"},
	Accent:   AccentTokens{"#ff2975", "#00fff5"},
	Border:   BorderTokens{"#3a3a5c", "#ff2975"},
	Semantic: SemanticTokens{"#05ffa1", "#ffe900", "#ff3860", "#00b4d8", "#9999cc"},
	Special:  SpecialTokens{"#666688", "#2a1533", "#ffc107"},
}
