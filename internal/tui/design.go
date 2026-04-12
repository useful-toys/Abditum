// Package tui — design system foundation.
//
// design.go is the single source of truth for all visual tokens, symbols,
// typography attributes and theme instances used by the TUI. It is aligned
// 1:1 with golden/tui-design-system.md.
//
// To change a colour value, edit only the theme instance (TokyoNight or
// Cyberpunk) below. To add a token category, add a new sub-struct and a
// field to Theme. Never hardcode hex strings in component files.
package tui

// ---------------------------------------------------------------------------
// Colour token groups
// ---------------------------------------------------------------------------

// SurfaceTokens holds background colours for the three surface levels.
type SurfaceTokens struct {
	Base   string // full-screen background
	Raised string // side panels and overlay windows
	Input  string // text fields inside dialogs (recessed tone)
}

// TextTokens holds colours for text roles.
type TextTokens struct {
	Primary   string // names, titles, readable content
	Secondary string // support text, hints, placeholders
	Disabled  string // inactive options
	Link      string // URLs and external references
}

// BorderTokens holds colours for dividing lines.
type BorderTokens struct {
	Default string // panel dividers, neutral modal borders
	Focused string // active panel, input fields, focused modals
}

// AccentTokens holds colours for interactive and highlight elements.
type AccentTokens struct {
	Primary   string // selection bar, navigation cursor, primary action button
	Secondary string // favourite star ★, folder names
}

// SemanticTokens holds colours that communicate application state.
// Never use semantic colours for decoration.
type SemanticTokens struct {
	Success string // completed operation, config ON
	Warning string // alert before permanent action, dirty-state prefixes (✦ ✎ ✗)
	Error   string // operation error, wrong password, destructive dialog border
	Info    string // contextual information
	Off     string // config OFF
}

// SpecialTokens holds colours for punctual uses without a semantic category.
type SpecialTokens struct {
	Muted     string // faded text without semantic connotation (intentionally low contrast)
	Highlight string // background behind the selected list item
	Match     string // text fragment matching the current search term
}

// ---------------------------------------------------------------------------
// Typography
// ---------------------------------------------------------------------------

// Typography holds ANSI attribute flags. Centralised here so a theme can
// disable an attribute globally for terminals that render it poorly.
// Blink is intentionally absent — do not use it.
type Typography struct {
	Bold          bool // universal — titles, selected item, default action
	Dim           bool // disabled items, secondary content
	Italic        bool // hints, virtual folders, auxiliary text
	Underline     bool // punctual use
	Strikethrough bool // items marked for deletion (pair with SymDeleted + Special.Muted)
}

// DefaultTypography enables all supported ANSI attributes.
var DefaultTypography = Typography{
	Bold:          true,
	Dim:           true,
	Italic:        true,
	Underline:     true,
	Strikethrough: true,
}

// ---------------------------------------------------------------------------
// Theme
// ---------------------------------------------------------------------------

// Theme groups all visual design tokens for a single theme.
// All string fields are lipgloss-compatible hex colours (e.g. "#1a1b26").
// Logo[0] is the topmost logo line; Logo[4] is the bottommost.
type Theme struct {
	Surface    SurfaceTokens
	Text       TextTokens
	Border     BorderTokens
	Accent     AccentTokens
	Semantic   SemanticTokens
	Special    SpecialTokens
	Logo       [5]string
	Typography Typography
}

// TokyoNight is the default theme — dark, cool-blue palette.
// Values sourced from golden/tui-design-system.md.
var TokyoNight = &Theme{
	Surface: SurfaceTokens{
		Base:   "#1a1b26",
		Raised: "#24283b",
		Input:  "#1e1f2e",
	},
	Text: TextTokens{
		Primary:   "#a9b1d6",
		Secondary: "#565f89",
		Disabled:  "#3b4261",
		Link:      "#7aa2f7",
	},
	Border: BorderTokens{
		Default: "#414868",
		Focused: "#7aa2f7",
	},
	Accent: AccentTokens{
		Primary:   "#7aa2f7",
		Secondary: "#bb9af7",
	},
	Semantic: SemanticTokens{
		Success: "#9ece6a",
		Warning: "#e0af68",
		Error:   "#f7768e",
		Info:    "#7dcfff",
		Off:     "#737aa2",
	},
	Special: SpecialTokens{
		Muted:     "#8690b5",
		Highlight: "#283457",
		Match:     "#f7c67a",
	},
	Logo: [5]string{
		"#9d7cd8",
		"#89ddff",
		"#7aa2f7",
		"#7dcfff",
		"#bb9af7",
	},
	Typography: DefaultTypography,
}

// Cyberpunk is the alternate theme — high-contrast neon palette.
// Values sourced from golden/tui-design-system.md.
var Cyberpunk = &Theme{
	Surface: SurfaceTokens{
		Base:   "#0a0a1a",
		Raised: "#1a1a2e",
		Input:  "#0e0e22",
	},
	Text: TextTokens{
		Primary:   "#e0e0ff",
		Secondary: "#8888aa",
		Disabled:  "#444466",
		Link:      "#ff2975",
	},
	Border: BorderTokens{
		Default: "#3a3a5c",
		Focused: "#ff2975",
	},
	Accent: AccentTokens{
		Primary:   "#ff2975",
		Secondary: "#00fff5",
	},
	Semantic: SemanticTokens{
		Success: "#05ffa1",
		Warning: "#ffe900",
		Error:   "#ff3860",
		Info:    "#00b4d8",
		Off:     "#9999cc",
	},
	Special: SpecialTokens{
		Muted:     "#666688",
		Highlight: "#2a1533",
		Match:     "#ffc107",
	},
	Logo: [5]string{
		"#ff2975",
		"#b026ff",
		"#00fff5",
		"#05ffa1",
		"#ff2975",
	},
	Typography: DefaultTypography,
}

// ---------------------------------------------------------------------------
// MessageKind
// ---------------------------------------------------------------------------

// MessageKind classifies the type of message shown in the message bar.
// It governs colour, symbol and TTL behaviour.
type MessageKind int

const (
	MessageSuccess MessageKind = iota
	MessageInfo
	MessageWarning
	MessageError
	MessageBusy
	MessageHint
)

// SymbolFor returns the canonical symbol string for a MessageKind.
func SymbolFor(kind MessageKind) string {
	switch kind {
	case MessageSuccess:
		return SymSuccess
	case MessageInfo:
		return SymInfo
	case MessageWarning:
		return SymWarning
	case MessageError:
		return SymError
	case MessageBusy:
		return SpinnerFrames[0]
	default: // MessageHint
		return SymBullet
	}
}

// ColorFor returns the theme hex colour that corresponds to a MessageKind.
func ColorFor(t *Theme, kind MessageKind) string {
	switch kind {
	case MessageSuccess:
		return t.Semantic.Success
	case MessageInfo:
		return t.Semantic.Info
	case MessageWarning:
		return t.Semantic.Warning
	case MessageError:
		return t.Semantic.Error
	case MessageBusy:
		return t.Accent.Primary
	default: // MessageHint
		return t.Text.Secondary
	}
}

// ---------------------------------------------------------------------------
// Symbols
// ---------------------------------------------------------------------------
//
// Complete inventory aligned with golden/tui-design-system.md §Ícones e Símbolos.
// Constraints: BMP only (U+0000–U+FFFF), no Nerd Fonts, no emojis.
// All symbols are 1 column wide except SymTreeConnector (2 columns, composed).

// Tree navigation
const (
	SymFolderCollapsed = "▶" // U+25B6 — collapsed folder
	SymFolderExpanded  = "▼" // U+25BC — expanded folder
	SymFolderEmpty     = "▷" // U+25B7 — empty folder
	SymLeaf            = "●" // U+25CF — leaf item
)

// Item state
const (
	SymFavorite = "★" // U+2605 — favourite item
	SymDeleted  = "✗" // U+2717 — marked for deletion
	SymCreated  = "✦" // U+2726 — newly created, unsaved
	SymModified = "✎" // U+270E — modified, unsaved
)

// Message bar symbols
const (
	SymSuccess = "✓" // U+2713 — success
	SymInfo    = "ℹ" // U+2139 — information
	SymWarning = "⚠" // U+26A0 — warning / alert
	SymError   = "✕" // U+2715 — error (distinct from SymDeleted ✗)
)

// UI elements
const (
	SymSensitiveField = "◉"  // U+25C9 — revealable field indicator
	SymSensMask       = "•"  // U+2022 — sensitive content mask character (same glyph as SymBullet, distinct semantic)
	SymCursor         = "▌"  // U+258C — text field cursor
	SymScrollUp       = "↑"  // U+2191 — scroll direction indicator (up)
	SymScrollDown     = "↓"  // U+2193 — scroll direction indicator (down)
	SymScrollThumb    = "■"  // U+25A0 — scroll position thumb
	SymEllipsis       = "…"  // U+2026 — truncation
	SymBullet         = "•"  // U+2022 — contextual indicator, hint marker, dirty marker
	SymHeaderSep      = "·"  // U+00B7 — header separator
	SymTreeConnector  = "<╡" // Basic Latin + U+2561 — tree → detail connector (2 columns)
)

// Border characters (Box Drawing block)
const (
	SymBorderH   = "─" // U+2500 — horizontal separator
	SymBorderV   = "│" // U+2502 — vertical separator
	SymCornerTL  = "╭" // U+256D — top-left rounded corner
	SymCornerTR  = "╮" // U+256E — top-right rounded corner
	SymCornerBL  = "╰" // U+2570 — bottom-left rounded corner
	SymCornerBR  = "╯" // U+256F — bottom-right rounded corner
	SymJunctionL = "├" // U+251C — T junction (left-pointing)
	SymJunctionT = "┬" // U+252C — T junction (top-pointing)
	SymJunctionB = "┴" // U+2534 — T junction (bottom-pointing)
	SymJunctionR = "┤" // U+2524 — T junction (right-pointing)
)

// SpinnerFrames is the four-frame activity spinner sequence.
// Note: these Geometric Shapes characters may render as 2 columns in some
// locales. Reserve 2 columns in the surrounding layout to prevent jitter.
var SpinnerFrames = [4]string{"◐", "◓", "◑", "◒"}

// SpinnerFrame returns the spinner character for a given animation frame index.
func SpinnerFrame(frame int) string {
	return SpinnerFrames[frame%4]
}

// ---------------------------------------------------------------------------
// Keyboard notation constants
// ---------------------------------------------------------------------------
//
// Used by the command bar and Help dialog to render key bindings.
// These are the canonical representations defined in golden/tui-design-system.md §Teclado.

const (
	KeyCtrl  = "⌃" // U+2303 — Ctrl modifier
	KeyShift = "⇧" // U+21E7 — Shift modifier
	KeyAlt   = "!" // no dedicated Unicode glyph — rendered as literal "!"
)
