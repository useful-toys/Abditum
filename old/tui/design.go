// Package tui provides the design system foundation for the Abditum TUI.
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
// Prescribed typography combinations
// ---------------------------------------------------------------------------
//
// Documented in golden/tui-design-system.md §Tipografia — Combinações previstas.
// These are not enforced by the type system; they are the canonical recipes
// that components should follow when styling specific element types.
//
//   TypoAlertTitle   — Bold + semantic colour  → modal alert / info title
//   TypoDeletedItem  — Dim + Strikethrough      → deleted item (pair with SymDeleted + Special.Muted)
//   TypoHintText     — Italic + text.secondary  → hints, auxiliary text, virtual folders
//
// Usage example:
//
//	style := lipgloss.NewStyle().Bold(theme.Typography.Bold).Foreground(lipgloss.Color(theme.Semantic.Warning))
//	// produces TypoAlertTitle for a warning modal.

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
	// MessageSuccess indicates a successful operation.
	MessageSuccess MessageKind = iota
	// MessageInfo indicates contextual information.
	MessageInfo
	// MessageWarning indicates a warning or alert.
	MessageWarning
	// MessageError indicates an operation error.
	MessageError
	// MessageBusy indicates an ongoing background operation.
	MessageBusy
	// MessageHint indicates a helpful hint or suggestion.
	MessageHint
)

// Symbol returns the canonical symbol string for the message kind.
func (k MessageKind) Symbol() string {
	switch k {
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

// Color returns the theme hex colour that corresponds to the message kind.
func (k MessageKind) Color(t *Theme) string {
	if t == nil {
		return ""
	}
	switch k {
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
	// SymFolderCollapsed is the symbol for a collapsed folder.
	SymFolderCollapsed = "▶" // U+25B6 — collapsed folder
	// SymFolderExpanded is the symbol for an expanded folder.
	SymFolderExpanded = "▼" // U+25BC — expanded folder
	// SymFolderEmpty is the symbol for an empty folder.
	SymFolderEmpty = "▷" // U+25B7 — empty folder
	// SymLeaf is the symbol for a leaf item.
	SymLeaf = "●" // U+25CF — leaf item
)

// Item state
const (
	// SymFavorite is the symbol for a favourite item.
	SymFavorite = "★" // U+2605 — favourite item
	// SymDeleted is the symbol for an item marked for deletion.
	SymDeleted = "✗" // U+2717 — marked for deletion
	// SymCreated is the symbol for a newly created, unsaved item.
	SymCreated = "✦" // U+2726 — newly created, unsaved
	// SymModified is the symbol for a modified, unsaved item.
	SymModified = "✎" // U+270E — modified, unsaved
)

// Message bar symbols
const (
	// SymSuccess is the symbol for a success message.
	SymSuccess = "✓" // U+2713 — success
	// SymInfo is the symbol for an information message.
	SymInfo = "ℹ" // U+2139 — information
	// SymWarning is the symbol for a warning message.
	SymWarning = "⚠" // U+26A0 — warning / alert
	// SymError is the symbol for an error message.
	SymError = "✕" // U+2715 — error (distinct from SymDeleted ✗)
)

// UI elements
const (
	// SymSensitiveField is the symbol for a revealable field indicator.
	SymSensitiveField = "◉" // U+25C9 — revealable field indicator
	// SymSensMask is the character used to mask sensitive content.
	SymSensMask = "•" // U+2022 — sensitive content mask character (same glyph as SymBullet, distinct semantic)
	// SymCursor is the symbol for a text field cursor.
	SymCursor = "▌" // U+258C — text field cursor
	// SymScrollUp is the indicator for scrolling up.
	SymScrollUp = "↑" // U+2191 — scroll direction indicator (up)
	// SymScrollDown is the indicator for scrolling down.
	SymScrollDown = "↓" // U+2193 — scroll direction indicator (down)
	// SymScrollThumb is the symbol for the scroll position thumb.
	SymScrollThumb = "■" // U+25A0 — scroll position thumb
	// SymEllipsis is the symbol for truncation.
	SymEllipsis = "…" // U+2026 — truncation
	// SymBullet is a contextual indicator or hint marker.
	SymBullet = "•" // U+2022 — contextual indicator, hint marker, dirty marker
	// SymHeaderSep is the symbol for a header separator.
	SymHeaderSep = "·" // U+00B7 — header separator
	// SymTreeConnector is the connector from tree to detail.
	SymTreeConnector = "<╡" // Basic Latin + U+2561 — tree → detail connector (2 columns)
)

// Border characters (Box Drawing block)
const (
	// SymBorderH is the horizontal border separator.
	SymBorderH = "─" // U+2500 — horizontal separator
	// SymBorderV is the vertical border separator.
	SymBorderV = "│" // U+2502 — vertical separator
	// SymCornerTL is the top-left rounded corner.
	SymCornerTL = "╭" // U+256D — top-left rounded corner
	// SymCornerTR is the top-right rounded corner.
	SymCornerTR = "╮" // U+256E — top-right rounded corner
	// SymCornerBL is the bottom-left rounded corner.
	SymCornerBL = "╰" // U+2570 — bottom-left rounded corner
	// SymCornerBR is the bottom-right rounded corner.
	SymCornerBR = "╯" // U+256F — bottom-right rounded corner
	// SymJunctionL is the left-pointing T junction.
	SymJunctionL = "├" // U+251C — T junction (left-pointing)
	// SymJunctionT is the top-pointing T junction.
	SymJunctionT = "┬" // U+252C — T junction (top-pointing)
	// SymJunctionB is the bottom-pointing T junction.
	SymJunctionB = "┴" // U+2534 — T junction (bottom-pointing)
	// SymJunctionR is the right-pointing T junction.
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
// Layout constants
// ---------------------------------------------------------------------------
//
// Minimum terminal dimensions and fixed zone heights defined in
// golden/tui-design-system.md §Dimensionamento e Layout.

const (
	// LayoutMinWidth is the minimum supported terminal width in columns.
	LayoutMinWidth = 80
	// LayoutMinHeight is the minimum supported terminal height in lines.
	LayoutMinHeight = 24

	// LayoutHeaderHeight is the fixed height of the header zone (lines).
	LayoutHeaderHeight = 2
	// LayoutMessageBarHeight is the fixed height of the message bar zone (lines).
	LayoutMessageBarHeight = 1
	// LayoutCommandBarHeight is the fixed height of the command bar zone (lines).
	LayoutCommandBarHeight = 1

	// LayoutTreePanelPct is the nominal percentage of screen width for the
	// left (tree / list) panel in two-panel modes (Vault, Templates).
	// The implementation may adjust ±5% for aesthetic alignment.
	LayoutTreePanelPct = 35
	// LayoutDetailPanelPct is the nominal percentage of screen width for the
	// right (detail) panel in two-panel modes.
	LayoutDetailPanelPct = 65
)

// ---------------------------------------------------------------------------
// Keyboard notation constants
// ---------------------------------------------------------------------------
//
// Used by the command bar and Help dialog to render key bindings.
// These are the canonical representations defined in golden/tui-design-system.md §Teclado.

// Modifier symbols.
const (
	// KeyCtrl is the notation for the Ctrl modifier.
	KeyCtrl = "⌃" // U+2303 — Ctrl modifier
	// KeyShift is the notation for the Shift modifier.
	KeyShift = "⇧" // U+21E7 — Shift modifier
	// KeyAlt is the notation for the Alt modifier.
	KeyAlt = "!" // no dedicated Unicode glyph — rendered as literal "!"
)

// Special key labels — written in full as shown in the command bar and Help dialog.
// Combinations are concatenated without spaces: "⌃Q", "⇧F6", "⌃!⇧Q".
const (
	// KeyEnter is the label for the Enter key.
	KeyEnter = "Enter"
	// KeyEsc is the label for the Escape key.
	KeyEsc = "Esc"
	// KeyTab is the label for the Tab key.
	KeyTab = "Tab"
	// KeyDel is the label for the Delete key.
	KeyDel = "Del"
	// KeyIns is the label for the Insert key.
	KeyIns = "Ins"
	// KeyHome is the label for the Home key.
	KeyHome = "Home"
	// KeyEnd is the label for the End key.
	KeyEnd = "End"
	// KeyPgUp is the label for the Page Up key.
	KeyPgUp = "PgUp"
	// KeyPgDn is the label for the Page Down key.
	KeyPgDn = "PgDn"
)

// Function key labels — used in the command bar.
const (
	KeyF1  = "F1"
	KeyF2  = "F2"
	KeyF3  = "F3"
	KeyF4  = "F4"
	KeyF5  = "F5"
	KeyF6  = "F6"
	KeyF7  = "F7"
	KeyF8  = "F8"
	KeyF9  = "F9"
	KeyF10 = "F10"
	KeyF11 = "F11"
	KeyF12 = "F12"
)

// Global shortcut strings — canonical representations used in the Help dialog.
// These four shortcuts work in every application context (global scope).
const (
	// ShortcutHelp opens / closes the Help dialog.
	ShortcutHelp = KeyF1
	// ShortcutThemeToggle switches between Tokyo Night and Cyberpunk themes.
	// Not displayed in the command bar; listed here for completeness.
	ShortcutThemeToggle = KeyF12
	// ShortcutQuit exits the application (with confirmation when there are
	// unsaved changes).
	ShortcutQuit = KeyCtrl + "Q"
	// ShortcutLockVault locks the vault immediately, discarding unsaved changes
	// without confirmation. The deliberately complex combination prevents
	// accidental activation.
	ShortcutLockVault = KeyCtrl + KeyAlt + KeyShift + "Q"
)
