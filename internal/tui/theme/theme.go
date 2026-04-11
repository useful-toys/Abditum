package theme

import (
	abditum_tokens "github.com/useful-toys/abditum/internal/tui/tokens"
)

// Theme defines the visual styling of the application. (D-12)
type Theme struct {
	// Base colors (as hex strings, applied via lipgloss.Color function)
	SurfaceBase   string
	SurfaceRaised string

	// Text colors
	TextPrimary   string
	TextSecondary string
	TextDisabled  string
	TextLink      string

	// Accent colors
	AccentPrimary   string
	AccentSecondary string

	// Message bar colors
	MsgInfo    string
	MsgSuccess string
	MsgWarning string
	MsgError   string
	MsgBusy    string
	MsgHint    string

	// Border color
	Border string
}

// ThemeTokyoNight is a dark theme inspired by the Tokyo Night color palette. (D-13)
var ThemeTokyoNight = &Theme{
	SurfaceBase:   abditum_tokens.ColorSurfaceBaseTokyoNight,
	SurfaceRaised: abditum_tokens.ColorSurfaceRaisedTokyoNight,

	TextPrimary:   abditum_tokens.ColorTextPrimaryTokyoNight,
	TextSecondary: abditum_tokens.ColorTextSecondaryTokyoNight,
	TextDisabled:  abditum_tokens.ColorTextDisabledTokyoNight,
	TextLink:      abditum_tokens.ColorTextLinkTokyoNight,

	AccentPrimary:   abditum_tokens.ColorAccentPrimary,
	AccentSecondary: abditum_tokens.ColorAccentSecondary,

	MsgInfo:    abditum_tokens.ColorInfo,
	MsgSuccess: abditum_tokens.ColorSuccess,
	MsgWarning: abditum_tokens.ColorWarn,
	MsgError:   abditum_tokens.ColorError,
	MsgBusy:    abditum_tokens.ColorBusy,
	MsgHint:    abditum_tokens.ColorHint,

	Border: abditum_tokens.ColorBorder,
}

// ThemeCyberpunk is a vibrant, high-contrast theme. (D-14)
var ThemeCyberpunk = &Theme{
	SurfaceBase:   abditum_tokens.ColorSurfaceBaseCyberpunk,
	SurfaceRaised: abditum_tokens.ColorSurfaceRaisedCyberpunk,

	TextPrimary:   abditum_tokens.ColorTextPrimaryCyberpunk,
	TextSecondary: abditum_tokens.ColorTextSecondaryCyberpunk,
	TextDisabled:  abditum_tokens.ColorTextDisabledCyberpunk,
	TextLink:      abditum_tokens.ColorTextLinkCyberpunk,

	AccentPrimary:   abditum_tokens.ColorAccentPrimary,
	AccentSecondary: abditum_tokens.ColorAccentSecondary,

	MsgInfo:    abditum_tokens.ColorInfo,
	MsgSuccess: abditum_tokens.ColorSuccess,
	MsgWarning: abditum_tokens.ColorWarn,
	MsgError:   abditum_tokens.ColorError,
	MsgBusy:    abditum_tokens.ColorBusy,
	MsgHint:    abditum_tokens.ColorHint,

	Border: abditum_tokens.ColorBorder,
}
