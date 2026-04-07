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
	// LogoGradient: 5-color violet→cyan progression per Design System § Gradiente do logo
	LogoGradient: []color.Color{
		lipgloss.Color("#9d7cd8"), // Line 1: Violet
		lipgloss.Color("#89ddff"), // Line 2: Cyan
		lipgloss.Color("#7aa2f7"), // Line 3: Blue
		lipgloss.Color("#7dcfff"), // Line 4: Light cyan
		lipgloss.Color("#bb9af7"), // Line 5: Purple
	},
}

// ThemeCyberpunk defines the Cyberpunk theme colors.
var ThemeCyberpunk = &Theme{
	SurfaceBase:     lipgloss.Color("#0a0a1a"),
	SurfaceRaised:   lipgloss.Color("#1a1a2e"),
	AccentPrimary:   lipgloss.Color("#ff2975"),
	AccentSecondary: lipgloss.Color("#00fff5"),
	TextPrimary:     lipgloss.Color("#e0e0ff"),
	TextSecondary:   lipgloss.Color("#8888aa"),
	SemanticSuccess: lipgloss.Color("#05ffa1"),
	SemanticWarning: lipgloss.Color("#ffe900"),
	SemanticError:   lipgloss.Color("#ff3860"),
	SemanticInfo:    lipgloss.Color("#00b4d8"),
	SemanticOff:     lipgloss.Color("#9999cc"),
	// LogoGradient: 5-color magenta→cyan→green progression per Design System § Gradiente do logo
	LogoGradient: []color.Color{
		lipgloss.Color("#ff2975"), // Line 1: Magenta/Pink
		lipgloss.Color("#b026ff"), // Line 2: Purple
		lipgloss.Color("#00fff5"), // Line 3: Cyan
		lipgloss.Color("#05ffa1"), // Line 4: Green
		lipgloss.Color("#ff2975"), // Line 5: Magenta/Pink
	},
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
