package tokens

import (
	lipgloss "charm.land/lipgloss/v2"
	types "github.com/useful-toys/abditum/internal/tui/types" // Renamed alias
)

// Semantic color constants (D-17)
const (
	ColorSuccess = "#9ece6a"
	ColorInfo    = "#7dcfff"
	ColorWarn    = "#e0af68"
	ColorError   = "#f7768e"
	ColorBusy    = "#7aa2f7"
	ColorHint    = "#565f89"

	// Base colors for Theme (from D-13, D-14)
	ColorSurfaceBaseTokyoNight   = "#1a1b26"
	ColorSurfaceRaisedTokyoNight = "#24283b"
	ColorTextPrimaryTokyoNight   = "#a9b1d6"
	ColorTextSecondaryTokyoNight = "#565f89"
	ColorTextDisabledTokyoNight  = "#3b4261"
	ColorTextLinkTokyoNight      = "#7aa2f7"

	ColorSurfaceBaseCyberpunk   = "#0f0017"
	ColorSurfaceRaisedCyberpunk = "#20012f"
	ColorTextPrimaryCyberpunk   = "#00ffff"
	ColorTextSecondaryCyberpunk = "#6a0dad"
	ColorTextDisabledCyberpunk  = "#36005c"
	ColorTextLinkCyberpunk      = "#ff00ff"
)

// Border and structural colors.
const (
	ColorBorderDefault = "#414868"
	ColorBorder        = ColorBorderDefault
)

// Accent colors.
const (
	ColorAccentPrimary   = "#7aa2f7"
	ColorAccentSecondary = "#bb9af7"
)

// Command bar colors.
const (
	ColorCommandKey   = "#7aa2f7"
	ColorCommandLabel = "#a9b1d6"
	ColorSeparator    = "#565f89"
)

// Help modal colors
const (
	ColorHelpTitle = "62"
	ColorHelpKey   = "11"
	ColorHelpSep   = "240"
	ColorHelpGroup = "14"
)

// Symbol constants
const (
	SymSuccess  = "✓"
	SymInfo     = "ℹ"
	SymWarn     = "⚠"
	SymError    = "✕"
	SymHint     = "•"
	SymBorder   = "─"
	SymEllipsis = "…"
	SymBullet   = "•"
)

// SpinnerFrames in display order: ◐ ◓ ◑ ◒
var SpinnerFrames = []string{"◐", "◓", "◑", "◒"}

// StyleSymbol returns a lipgloss.Style with the correct color and formatting
// for the given MsgKind.
func StyleSymbol(kind types.MsgKind) lipgloss.Style {
	switch kind {
	case types.MsgSuccess:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSuccess))
	case types.MsgInfo:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorInfo))
	case types.MsgWarning:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorWarn))
	case types.MsgError:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorError)).Bold(true)
	case types.MsgBusy:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorBusy))
	default: // types.MsgHint
		return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorHint)).Italic(true)
	}
}

// StyleBorder returns a style for the border character.
func StyleBorder() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorBorder))
}

// StyleCommandKey returns a bold style for command bar keys.
func StyleCommandKey() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorCommandKey)).Bold(true)
}

// StyleCommandLabel returns a style for command bar labels.
func StyleCommandLabel() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorCommandLabel))
}

// SymbolForKind returns the correct symbol character for a MsgKind.
func SymbolForKind(kind types.MsgKind) string {
	switch kind {
	case types.MsgSuccess:
		return SymSuccess
	case types.MsgInfo:
		return SymInfo
	case types.MsgWarning:
		return SymWarn
	case types.MsgError:
		return SymError
	case types.MsgBusy:
		return "" // spinner uses SpinnerFrames directly
	default: // types.MsgHint
		return SymHint
	}
}
