package tui

import (
	"charm.land/lipgloss/v2"
)

// Theme identifica o tema visual ativo.
type Theme int

const (
	ThemeTokyoNight Theme = iota
	ThemeCyberpunk
)

// palette define os tokens de cor de um tema.
type palette struct {
	// Superfícies
	surfaceBase   string
	surfaceRaised string
	surfaceInput  string

	// Texto
	textPrimary   string
	textSecondary string
	textDisabled  string
	textLink      string

	// Bordas
	borderDefault string
	borderFocused string

	// Interação
	accentPrimary   string
	accentSecondary string

	// Semânticas
	semanticSuccess string
	semanticWarning string
	semanticError   string
	semanticInfo    string
	semanticOff     string

	// Especiais
	specialMuted     string
	specialHighlight string
	specialMatch     string

	// Gradiente do logo (5 linhas)
	logoGrad [5]string
}

var tokyoNight = palette{
	surfaceBase:      "#1a1b26",
	surfaceRaised:    "#24283b",
	surfaceInput:     "#1e1f2e",
	textPrimary:      "#a9b1d6",
	textSecondary:    "#565f89",
	textDisabled:     "#3b4261",
	textLink:         "#7aa2f7",
	borderDefault:    "#414868",
	borderFocused:    "#7aa2f7",
	accentPrimary:    "#7aa2f7",
	accentSecondary:  "#bb9af7",
	semanticSuccess:  "#9ece6a",
	semanticWarning:  "#e0af68",
	semanticError:    "#f7768e",
	semanticInfo:     "#7dcfff",
	semanticOff:      "#737aa2",
	specialMuted:     "#8690b5",
	specialHighlight: "#283457",
	specialMatch:     "#f7c67a",
	logoGrad:         [5]string{"#9d7cd8", "#89ddff", "#7aa2f7", "#7dcfff", "#bb9af7"},
}

var cyberpunk = palette{
	surfaceBase:      "#0a0a1a",
	surfaceRaised:    "#1a1a2e",
	surfaceInput:     "#0e0e22",
	textPrimary:      "#e0e0ff",
	textSecondary:    "#8888aa",
	textDisabled:     "#444466",
	textLink:         "#ff2975",
	borderDefault:    "#3a3a5c",
	borderFocused:    "#ff2975",
	accentPrimary:    "#ff2975",
	accentSecondary:  "#00fff5",
	semanticSuccess:  "#05ffa1",
	semanticWarning:  "#ffe900",
	semanticError:    "#ff3860",
	semanticInfo:     "#00b4d8",
	semanticOff:      "#9999cc",
	specialMuted:     "#666688",
	specialHighlight: "#2a1533",
	specialMatch:     "#ffc107",
	logoGrad:         [5]string{"#ff2975", "#b026ff", "#00fff5", "#05ffa1", "#ff2975"},
}

// styles contém lipgloss styles pré-calculados para um tema.
type styles struct {
	pal palette

	// Texto
	TextPrimary   lipgloss.Style
	TextSecondary lipgloss.Style
	TextDisabled  lipgloss.Style

	// Interação
	AccentPrimary   lipgloss.Style
	AccentSecondary lipgloss.Style

	// Semânticas
	SemanticSuccess lipgloss.Style
	SemanticWarning lipgloss.Style
	SemanticError   lipgloss.Style
	SemanticInfo    lipgloss.Style

	// Especiais
	SpecialMatch     lipgloss.Style
	SpecialHighlight lipgloss.Style
	SpecialMuted     lipgloss.Style

	// Bordas
	BorderDefault lipgloss.Style
	BorderFocused lipgloss.Style

	// Compostos
	AppName     lipgloss.Style // "Abditum" — accent.primary + bold
	CoffeeName  lipgloss.Style // nome do cofre — text.secondary
	DirtyDot    lipgloss.Style // "•" — semantic.warning
	TabActive   lipgloss.Style // aba ativa — accent.primary + bold + special.highlight bg
	TabInactive lipgloss.Style // aba inativa — text.secondary
	TabBorder   lipgloss.Style // bordas das abas — border.default

	ActionKey   lipgloss.Style // tecla na barra — accent.primary + bold
	ActionLabel lipgloss.Style // label na barra — text.primary
	ActionSep   lipgloss.Style // "·" — text.secondary

	MsgSuccess lipgloss.Style
	MsgInfo    lipgloss.Style
	MsgWarning lipgloss.Style
	MsgError   lipgloss.Style
	MsgSpinner lipgloss.Style
	MsgHint    lipgloss.Style

	Selected lipgloss.Style // item selecionado — special.highlight bg + bold
}

func newStyles(p palette) styles {
	s := styles{pal: p}

	s.TextPrimary = lipgloss.NewStyle().Foreground(lipgloss.Color(p.textPrimary))
	s.TextSecondary = lipgloss.NewStyle().Foreground(lipgloss.Color(p.textSecondary))
	s.TextDisabled = lipgloss.NewStyle().Foreground(lipgloss.Color(p.textDisabled))

	s.AccentPrimary = lipgloss.NewStyle().Foreground(lipgloss.Color(p.accentPrimary))
	s.AccentSecondary = lipgloss.NewStyle().Foreground(lipgloss.Color(p.accentSecondary))

	s.SemanticSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color(p.semanticSuccess))
	s.SemanticWarning = lipgloss.NewStyle().Foreground(lipgloss.Color(p.semanticWarning))
	s.SemanticError = lipgloss.NewStyle().Foreground(lipgloss.Color(p.semanticError))
	s.SemanticInfo = lipgloss.NewStyle().Foreground(lipgloss.Color(p.semanticInfo))

	s.SpecialMatch = lipgloss.NewStyle().Foreground(lipgloss.Color(p.specialMatch)).Bold(true)
	s.SpecialHighlight = lipgloss.NewStyle().Background(lipgloss.Color(p.specialHighlight))
	s.SpecialMuted = lipgloss.NewStyle().Foreground(lipgloss.Color(p.specialMuted))

	s.BorderDefault = lipgloss.NewStyle().Foreground(lipgloss.Color(p.borderDefault))
	s.BorderFocused = lipgloss.NewStyle().Foreground(lipgloss.Color(p.borderFocused))

	s.AppName = lipgloss.NewStyle().Foreground(lipgloss.Color(p.accentPrimary)).Bold(true)
	s.CoffeeName = lipgloss.NewStyle().Foreground(lipgloss.Color(p.textSecondary))
	s.DirtyDot = lipgloss.NewStyle().Foreground(lipgloss.Color(p.semanticWarning))
	s.TabActive = lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.accentPrimary)).
		Background(lipgloss.Color(p.specialHighlight)).
		Bold(true)
	s.TabInactive = lipgloss.NewStyle().Foreground(lipgloss.Color(p.textSecondary))
	s.TabBorder = lipgloss.NewStyle().Foreground(lipgloss.Color(p.borderDefault))

	s.ActionKey = lipgloss.NewStyle().Foreground(lipgloss.Color(p.accentPrimary)).Bold(true)
	s.ActionLabel = lipgloss.NewStyle().Foreground(lipgloss.Color(p.textPrimary))
	s.ActionSep = lipgloss.NewStyle().Foreground(lipgloss.Color(p.textSecondary))

	s.MsgSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color(p.semanticSuccess))
	s.MsgInfo = lipgloss.NewStyle().Foreground(lipgloss.Color(p.semanticInfo))
	s.MsgWarning = lipgloss.NewStyle().Foreground(lipgloss.Color(p.semanticWarning))
	s.MsgError = lipgloss.NewStyle().Foreground(lipgloss.Color(p.semanticError)).Bold(true)
	s.MsgSpinner = lipgloss.NewStyle().Foreground(lipgloss.Color(p.accentPrimary))
	s.MsgHint = lipgloss.NewStyle().Foreground(lipgloss.Color(p.textSecondary)).Italic(true)

	s.Selected = lipgloss.NewStyle().
		Background(lipgloss.Color(p.specialHighlight)).
		Foreground(lipgloss.Color(p.textPrimary)).
		Bold(true)

	return s
}
