package tui

import (
	"charm.land/lipgloss/v2"
)

// ─────────────────────────────────────────────────────────────────────────────
// Semantic color constants (D-17)
// All hex values centralized — zero hardcoding in consumers (D-18).
// ─────────────────────────────────────────────────────────────────────────────

// Message bar semantic colors (MsgKind palette).
const (
	ColorSuccess = "#9ece6a" // MsgSuccess — operação concluída
	ColorInfo    = "#7dcfff" // MsgInfo    — informação neutra
	ColorWarn    = "#e0af68" // MsgWarn    — atenção
	ColorError   = "#f7768e" // MsgError   — falha
	ColorBusy    = "#7aa2f7" // MsgBusy    — spinner
	ColorHint    = "#565f89" // MsgHint    — dica contextual
)

// Border and structural colors.
const (
	ColorBorder = "#414868" // default border line
)

// Command bar colors.
const (
	ColorCommandKey   = "#7aa2f7" // action key token (bold)
	ColorCommandLabel = "#a9b1d6" // action label text
	ColorSeparator    = "#565f89" // separator dots
)

// Help modal colors (lipgloss 256-color palette indices as strings).
const (
	ColorHelpTitle = "62"  // help title border
	ColorHelpKey   = "11"  // shortcut key
	ColorHelpSep   = "240" // separator lines
	ColorHelpGroup = "14"  // group label
)

// ─────────────────────────────────────────────────────────────────────────────
// Symbol constants
// ─────────────────────────────────────────────────────────────────────────────

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

// ─────────────────────────────────────────────────────────────────────────────
// Pre-built lipgloss style helpers (functions, not package-level vars)
// ─────────────────────────────────────────────────────────────────────────────

// StyleSymbol returns a lipgloss.Style with the correct color and formatting
// for the given MsgKind. Matches the switch in messages.go lines 131-150.
func StyleSymbol(kind MsgKind) lipgloss.Style {
	switch kind {
	case MsgSuccess:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSuccess))
	case MsgInfo:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorInfo))
	case MsgWarn:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorWarn))
	case MsgError:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorError)).Bold(true)
	case MsgBusy:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorBusy))
	default: // MsgHint
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

// ─────────────────────────────────────────────────────────────────────────────
// Helper functions
// ─────────────────────────────────────────────────────────────────────────────

// SymbolForKind returns the correct symbol character for a MsgKind.
func SymbolForKind(kind MsgKind) string {
	switch kind {
	case MsgSuccess:
		return SymSuccess
	case MsgInfo:
		return SymInfo
	case MsgWarn:
		return SymWarn
	case MsgError:
		return SymError
	case MsgBusy:
		return "" // spinner uses SpinnerFrames directly
	default: // MsgHint
		return SymHint
	}
}

// SpinnerFrame returns the spinner character for the given animation frame.
func SpinnerFrame(frame int) string {
	return SpinnerFrames[frame%4]
}
