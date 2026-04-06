package tui

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// Theme holds all the design system tokens for a specific theme.
type Theme struct {
	// Surface colors
	SurfaceBase   color.Color
	SurfaceRaised color.Color

	// Accent colors
	AccentPrimary   color.Color
	AccentSecondary color.Color

	// Text colors
	TextPrimary   color.Color
	TextSecondary color.Color

	// Semantic colors
	SemanticSuccess color.Color
	SemanticWarning color.Color
	SemanticError   color.Color
	SemanticInfo    color.Color
	SemanticOff     color.Color

	// Gradient colors for the logo
	LogoGradient []color.Color
}

// ThemeTokyoNight defines the Tokyo Night theme colors.
var ThemeTokyoNight = &Theme{
	SurfaceBase:     lipgloss.Color("#1a1b26"),
	SurfaceRaised:   lipgloss.Color("#24283b"),
	AccentPrimary:   lipgloss.Color("#7aa2f7"),
	AccentSecondary: lipgloss.Color("#bb9af7"),
	TextPrimary:     lipgloss.Color("#a9b1d6"),
	TextSecondary:   lipgloss.Color("#565f89"),
	SemanticSuccess: lipgloss.Color("#9ece6a"),
	SemanticWarning: lipgloss.Color("#e0af68"),
	SemanticError:   lipgloss.Color("#f7768e"),
	SemanticInfo:    lipgloss.Color("#7dcfff"),
	SemanticOff:     lipgloss.Color("#737aa2"),
	LogoGradient:    []color.Color{lipgloss.Color("#bb9af7"), lipgloss.Color("#7aa2f7")},
}

// ThemeCyberpunk defines the Cyberpunk theme colors (placeholder values for now).
var ThemeCyberpunk = &Theme{
	SurfaceBase:     lipgloss.Color("#0d0d0d"),
	SurfaceRaised:   lipgloss.Color("#1a1a1a"),
	AccentPrimary:   lipgloss.Color("#00ffff"),
	AccentSecondary: lipgloss.Color("#ff00ff"),
	TextPrimary:     lipgloss.Color("#ffffff"),
	TextSecondary:   lipgloss.Color("#cccccc"),
	SemanticSuccess: lipgloss.Color("#00ff00"),
	SemanticWarning: lipgloss.Color("#ffff00"),
	SemanticError:   lipgloss.Color("#ff0000"),
	SemanticInfo:    lipgloss.Color("#00aaff"),
	SemanticOff:     lipgloss.Color("#888888"),
	LogoGradient:    []color.Color{lipgloss.Color("#ff00ff"), lipgloss.Color("#00ffff")},
}

// ApplyTheme is a placeholder for propagating theme to a child model.
// Actual implementation will use type-switching in rootModel.
func (t *Theme) ApplyTheme(child childModel) {
	// This method is a stub for now. The rootModel will handle the actual
	// propagation by calling ApplyTheme on concrete child types.
}

// To implement the ApplyTheme method on child models,
// each child model will need to have a field like:
// `theme *Theme`
// and a method:
// `func (m *childModelType) ApplyTheme(t *Theme) { m.theme = t }`
