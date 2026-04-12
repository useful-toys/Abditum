# Design System Consolidation — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace four conflicting design system files (`theme.go`, `tokens.go`, `theme/theme.go`, `tokens/tokens.go`) with a single `internal/tui/design.go` that is 100% aligned with `golden/tui-design-system.md`, eliminating the triple-definition of `MsgKind` and all field-access divergences.

**Architecture:** One new file (`design.go`) is created in `package tui` containing all types, constants and instances. The file is written before any deletion so the codebase always compiles. Then all consumers are updated. Then the four old files are deleted. The subpackages `tui/theme` and `tui/tokens` are removed entirely. `tui/types` keeps everything except `MsgKind`.

**Tech Stack:** Go 1.21+, `github.com/charmbracelet/lipgloss`, `go build`, `go vet`, `go test`.

**Spec:** `docs/superpowers/specs/2026-04-12-design-system-consolidation-design.md`

---

## File Map

| Action | File | Responsibility |
|---|---|---|
| **Create** | `internal/tui/design.go` | All types, constants, instances for the design system |
| **Modify** | `internal/tui/messages.go` | Remove `MsgKind` definition; rename `MsgWarn` → `MessageWarning` throughout; update `RenderMessageBar` field accesses |
| **Modify** | `internal/tui/types/types.go` | Remove `MsgKind` definition and its constants |
| **Modify** | `internal/tui/ascii.go` | `t.LogoGradient` → `t.Logo` |
| **Modify** | `internal/tui/welcome.go` | `m.theme.TextSecondary` → `m.theme.Text.Secondary` |
| **Modify** | `internal/tui/header.go` | All flat field accesses → nested |
| **Modify** | `internal/tui/actions.go` | All flat field accesses → nested |
| **Modify** | `internal/tui/root.go` | `ThemeTokyoNight`/`ThemeCyberpunk` → `TokyoNight`/`Cyberpunk`; flat fields → nested |
| **Modify** | `internal/tui/filepicker.go` | Flat fields → nested; hardcoded `#3d59a1` → `theme.Special.Highlight`; remove `tokens` import |
| **Modify** | `internal/tui/settings.go` | `m.theme.SemanticInfo` → `m.theme.Semantic.Info` |
| **Modify** | `internal/tui/templatedetail.go` | `m.theme.SemanticInfo` → `m.theme.Semantic.Info` |
| **Modify** | `internal/tui/templatelist.go` | `m.theme.SemanticInfo` → `m.theme.Semantic.Info` |
| **Modify** | `internal/tui/secretdetail.go` | `m.theme.SemanticInfo` → `m.theme.Semantic.Info` |
| **Modify** | `internal/tui/vaulttree.go` | `m.theme.SemanticInfo` → `m.theme.Semantic.Info` |
| **Modify** | `internal/tui/dialogs.go` | `ThemeTokyoNight` → `TokyoNight` |
| **Modify** | `internal/tui/passwordentry.go` | `MsgError`/`MsgHint` → `MessageError`/`MessageHint`; remove `tokens` import; use `Sym*` from `design.go` |
| **Modify** | `internal/tui/passwordcreate.go` | Same as `passwordentry.go` |
| **Modify** | `internal/tui/flow_create_vault.go` | `MsgBusy`/`MsgError`/`MsgSuccess` → `MessageBusy`/`MessageError`/`MessageSuccess` |
| **Modify** | `internal/tui/flow_open_vault.go` | Same as `flow_create_vault.go` |
| **Modify** | `internal/tui/flow_save_and_exit.go` | `MsgError` → `MessageError` |
| **Modify** | `internal/tui/messages_test.go` | `ThemeTokyoNight` → `TokyoNight`; `MsgKind` refs → `MessageKind`; `MsgWarn` → `MessageWarning` etc. |
| **Modify** | `internal/tui/welcome_test.go` | `ThemeTokyoNight`/`ThemeCyberpunk` → `TokyoNight`/`Cyberpunk` |
| **Modify** | `internal/tui/root_test.go` | `ThemeTokyoNight` → `TokyoNight`; `MsgHint` → `MessageHint` |
| **Modify** | `internal/tui/passwordentry_test.go` | `ThemeTokyoNight` → `TokyoNight`; `MsgError` → `MessageError` |
| **Modify** | `internal/tui/passwordcreate_test.go` | Same as `passwordentry_test.go` |
| **Modify** | `internal/tui/actions_test.go` | `ThemeTokyoNight` → `TokyoNight` |
| **Modify** | `internal/tui/filepicker_test.go` | `ThemeTokyoNight` → `TokyoNight` |
| **Modify** | `internal/tui/exit_flow_integration_test.go` | `ThemeTokyoNight` → `TokyoNight` |
| **Modify** | `internal/tui/flow_open_vault_test.go` | `ThemeTokyoNight` → `TokyoNight`; `MsgBusy`/`MsgError` → `MessageBusy`/`MessageError` |
| **Modify** | `internal/tui/flow_create_vault_test.go` | Same as `flow_open_vault_test.go` |
| **Modify** | `internal/tui/flow_save_and_exit_test.go` | `ThemeTokyoNight` → `TokyoNight` |
| **Delete** | `internal/tui/theme.go` | Replaced by `design.go` |
| **Delete** | `internal/tui/tokens.go` | Replaced by `design.go` |
| **Delete** | `internal/tui/theme/theme.go` | Replaced by `design.go` |
| **Delete** | `internal/tui/tokens/tokens.go` | Replaced by `design.go` |

---

## Task 1: Create `design.go`

**Files:**
- Create: `internal/tui/design.go`

- [ ] **Step 1.1: Create `design.go` with full content**

Create `internal/tui/design.go` with the following complete content:

```go
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
```

- [ ] **Step 1.2: Verify the file compiles in isolation**

```
go build github.com/useful-toys/abditum/internal/tui
```

Expected: compile error because `messages.go` still defines `MsgKind` (duplicate). That is fine — the old files still exist. We are checking that `design.go` itself has no syntax errors by looking at the error message (it should be about redeclaration, not syntax).

- [ ] **Step 1.3: Commit**

```
git add internal/tui/design.go
git commit -m "feat: add design.go — single source of truth for design system"
```

---

## Task 2: Remove `MsgKind` from `messages.go` and rename to `MessageKind`

**Files:**
- Modify: `internal/tui/messages.go`

- [ ] **Step 2.1: Edit `messages.go`**

In `messages.go`, make these changes:

**a) Delete the `MsgKind` type block** (the `type MsgKind int` and the `const (MsgSuccess ... MsgHint)` block).

**b) In the `DisplayMessage` struct**, the field `Kind MsgKind` becomes `Kind MessageKind`.

**c) In the `activeMessage` struct**, the field `kind MsgKind` becomes `kind MessageKind`.

**d) In `Show` signature**: `func (m *MessageManager) Show(kind MsgKind, ...)` → `func (m *MessageManager) Show(kind MessageKind, ...)`.

**e) Rename all `MsgWarn` occurrences** to `MessageWarning`. Rename all `MsgSuccess` → `MessageSuccess`, `MsgInfo` → `MessageInfo`, `MsgError` → `MessageError`, `MsgBusy` → `MessageBusy`, `MsgHint` → `MessageHint`.

**f) In `RenderMessageBar`**, update the flat `Theme` field accesses:
- `theme.SurfaceRaised` → `theme.Surface.Raised`
- `theme.SemanticSuccess` → `theme.Semantic.Success`
- `theme.SemanticInfo` → `theme.Semantic.Info`
- `theme.SemanticWarning` → `theme.Semantic.Warning`
- `theme.SemanticError` → `theme.Semantic.Error`
- `theme.AccentPrimary` → `theme.Accent.Primary`
- `theme.AccentSecondary` → `theme.Accent.Secondary`

**g) In `RenderMessageBar` line ~160**, `symbol = SymHint` → `symbol = SymBullet`.
`SymHint` was defined in `tokens.go` as `"•"` — the equivalent in `design.go` is `SymBullet`.

Also remove any `TickMsg`, `toggleThemeMsg`, `pwdEnteredMsg`, `pwdCreatedMsg` type definitions from `messages.go` **only if** they are also defined in `types/types.go` — check before deleting. If they exist only in `messages.go`, leave them.

- [ ] **Step 2.2: Verify**

```
go build github.com/useful-toys/abditum/internal/tui
```

Expected: errors about `MsgKind` still defined in `types/types.go` and possibly usage of old names in other files. That is expected — we are migrating incrementally. What must NOT appear: errors inside `messages.go` itself.

- [ ] **Step 2.3: Commit**

```
git add internal/tui/messages.go
git commit -m "refactor: migrate messages.go to MessageKind and nested Theme fields"
```

---

## Task 3: Remove `MsgKind` from `types/types.go`

**Files:**
- Modify: `internal/tui/types/types.go`

- [ ] **Step 3.1: Edit `types/types.go`**

Delete the `MsgKind` type and its constants block from `types/types.go`:

```go
// DELETE these lines:
type MsgKind int

const (
    MsgInfo    MsgKind = iota
    MsgSuccess
    MsgWarning
    MsgError
    MsgBusy
    MsgHint
)
```

Keep everything else: `WorkArea`, `FilePickerMode`, `PwdEnteredMsg`, `PwdCreatedMsg`, `ErrorMsg`, `ModalResult`, `TickMsg`, `ToggleThemeMsg`, `Action`, `Scope`.

- [ ] **Step 3.2: Verify**

```
go build github.com/useful-toys/abditum/internal/tui/types
```

Expected: PASS (the `types` subpackage no longer defines `MsgKind` so it should compile cleanly).

- [ ] **Step 3.3: Commit**

```
git add internal/tui/types/types.go
git commit -m "refactor: remove duplicate MsgKind from types package"
```

---

## Task 4: Update `tokens/tokens.go` to use `MessageKind`

**Files:**
- Modify: `internal/tui/tokens/tokens.go`

> Note: this file will be deleted in Task 12, but it must compile until then.

- [ ] **Step 4.1: Edit `tokens/tokens.go`**

The file currently imports `types "github.com/useful-toys/abditum/internal/tui/types"` and uses `types.MsgKind`. Since `types.MsgKind` no longer exists, change the import to the parent `tui` package and use `tui.MessageKind`:

Change import:
```go
// Before:
types "github.com/useful-toys/abditum/internal/tui/types"

// After:
tui "github.com/useful-toys/abditum/internal/tui"
```

Update all `types.MsgKind` references to `tui.MessageKind`, `types.MsgInfo` → `tui.MessageInfo`, `types.MsgSuccess` → `tui.MessageSuccess`, `types.MsgWarning` → `tui.MessageWarning`, `types.MsgError` → `tui.MessageError`, `types.MsgBusy` → `tui.MessageBusy`, `types.MsgHint` → `tui.MessageHint`.

> **Important:** this creates an import cycle if `tui` imports `tokens`. Check: does any non-deleted file in `package tui` currently import `tokens`? Yes — `passwordentry.go` and `passwordcreate.go`. Those will be fixed in Task 8. Until then, if the cycle blocks compilation, an alternative is to temporarily define a local type alias in `tokens.go` — but the clean path is to fix all consumers in the next tasks.

- [ ] **Step 4.2: Verify**

```
go build github.com/useful-toys/abditum/internal/tui/tokens
```

Expected: may fail due to import cycle until Tasks 8 is done. That is acceptable — move on.

- [ ] **Step 4.3: Commit**

```
git add internal/tui/tokens/tokens.go
git commit -m "refactor: update tokens subpackage to reference tui.MessageKind"
```

---

## Task 5: Update `ascii.go`

**Files:**
- Modify: `internal/tui/ascii.go`

- [ ] **Step 5.1: Edit `ascii.go`**

Change:
```go
colors := t.LogoGradient
```
To:
```go
colors := t.Logo[:]
```

> `Logo` is `[5]string` in the new design. The slice `t.Logo[:]` gives a `[]string` if the existing code iterates a slice. If the code already indexes directly (e.g., `t.LogoGradient[0]`), change each index access to `t.Logo[0]` etc.

Read `ascii.go` first to see the exact usage before editing.

- [ ] **Step 5.2: Verify**

```
go build github.com/useful-toys/abditum/internal/tui
```

- [ ] **Step 5.3: Commit**

```
git add internal/tui/ascii.go
git commit -m "refactor: ascii.go uses theme.Logo instead of theme.LogoGradient"
```

---

## Task 6: Update `header.go`, `welcome.go`, `actions.go`

**Files:**
- Modify: `internal/tui/header.go`
- Modify: `internal/tui/welcome.go`
- Modify: `internal/tui/actions.go`

- [ ] **Step 6.1: Edit `header.go`**

Apply field renames (all fields used in `header.go`):
- `theme.AccentPrimary` → `theme.Accent.Primary`
- `theme.SurfaceRaised` → `theme.Surface.Raised`
- `theme.TextPrimary` → `theme.Text.Primary`
- `theme.SemanticWarning` → `theme.Semantic.Warning`
- `theme.TextSecondary` → `theme.Text.Secondary`

- [ ] **Step 6.2: Edit `welcome.go`**

Apply field renames:
- `m.theme.TextSecondary` → `m.theme.Text.Secondary`

- [ ] **Step 6.3: Edit `actions.go`**

Apply field renames:
- `theme.AccentPrimary` → `theme.Accent.Primary`
- `theme.TextPrimary` → `theme.Text.Primary`
- `theme.TextSecondary` → `theme.Text.Secondary`

- [ ] **Step 6.4: Verify**

```
go build github.com/useful-toys/abditum/internal/tui
```

- [ ] **Step 6.5: Commit**

```
git add internal/tui/header.go internal/tui/welcome.go internal/tui/actions.go
git commit -m "refactor: migrate header, welcome, actions to nested Theme fields"
```

---

## Task 7: Update `settings.go`, `templatedetail.go`, `templatelist.go`, `secretdetail.go`, `vaulttree.go`

**Files:**
- Modify: `internal/tui/settings.go`
- Modify: `internal/tui/templatedetail.go`
- Modify: `internal/tui/templatelist.go`
- Modify: `internal/tui/secretdetail.go`
- Modify: `internal/tui/vaulttree.go`

All five files have the same single change: `m.theme.SemanticInfo` → `m.theme.Semantic.Info`.

- [ ] **Step 7.1: Edit all five files**

In each file, change every occurrence of `m.theme.SemanticInfo` to `m.theme.Semantic.Info`.

- [ ] **Step 7.2: Verify**

```
go build github.com/useful-toys/abditum/internal/tui
```

- [ ] **Step 7.3: Commit**

```
git add internal/tui/settings.go internal/tui/templatedetail.go internal/tui/templatelist.go internal/tui/secretdetail.go internal/tui/vaulttree.go
git commit -m "refactor: migrate stub view files to nested Theme fields"
```

---

## Task 8: Update `passwordentry.go` and `passwordcreate.go`

**Files:**
- Modify: `internal/tui/passwordentry.go`
- Modify: `internal/tui/passwordcreate.go`

Both files import `tokens` for `tokens.SymError` and `tokens.SymHint`. After this task they use the constants from `design.go` directly (same package).

- [ ] **Step 8.1: Edit `passwordentry.go`**

**a) Remove import:**
```go
"github.com/useful-toys/abditum/internal/tui/tokens"
```

**b) Replace symbol references:**
- `tokens.SymError` → `SymError`
- `tokens.SymHint` → `SymBullet`

> `SymHint` does not exist in `design.go`. The hint symbol is `SymBullet` (`•`). Verify by reading the file first — if it uses `tokens.SymHint` for a hint message, replace with `SymBullet`.

**c) Replace `MsgKind` references:**
- `MsgError` → `MessageError`
- `MsgHint` → `MessageHint`

- [ ] **Step 8.2: Edit `passwordcreate.go`**

Apply the same changes as Step 8.1.

- [ ] **Step 8.3: Verify**

```
go build github.com/useful-toys/abditum/internal/tui
```

Expected: PASS. The import cycle introduced in Task 4 is now resolved because `passwordentry.go` and `passwordcreate.go` no longer import `tokens`.

- [ ] **Step 8.4: Commit**

```
git add internal/tui/passwordentry.go internal/tui/passwordcreate.go
git commit -m "refactor: remove tokens import from password components; use design.go symbols"
```

---

## Task 9: Update `dialogs.go` and `root.go`

**Files:**
- Modify: `internal/tui/dialogs.go`
- Modify: `internal/tui/root.go`

- [ ] **Step 9.1: Edit `dialogs.go`**

Replace `ThemeTokyoNight` → `TokyoNight` wherever it appears as a default value.

- [ ] **Step 9.2: Edit `root.go`**

**a) Theme variable references:**
- `ThemeTokyoNight` → `TokyoNight`
- `ThemeCyberpunk` → `Cyberpunk`

**b) Flat field accesses:**
- `m.theme.SurfaceBase` → `m.theme.Surface.Base`
- `theme.AccentPrimary` → `theme.Accent.Primary`
- `theme.TextPrimary` → `theme.Text.Primary`
- `theme.TextSecondary` → `theme.Text.Secondary`

- [ ] **Step 9.3: Verify**

```
go build github.com/useful-toys/abditum/internal/tui
```

- [ ] **Step 9.4: Commit**

```
git add internal/tui/dialogs.go internal/tui/root.go
git commit -m "refactor: root and dialogs use TokyoNight/Cyberpunk and nested Theme fields"
```

---

## Task 10: Update `filepicker.go`

**Files:**
- Modify: `internal/tui/filepicker.go`

`filepicker.go` is the largest file (1095 lines) and has the most field accesses plus one hardcoded hex colour.

- [ ] **Step 10.1: Edit `filepicker.go`**

**a) Theme variable reference:**
- `ThemeTokyoNight` → `TokyoNight`

**b) Flat field accesses** (all occurrences):
- `theme.TextPrimary` → `theme.Text.Primary`
- `theme.TextSecondary` → `theme.Text.Secondary`
- `theme.AccentPrimary` → `theme.Accent.Primary`
- `theme.AccentSecondary` → `theme.Accent.Secondary`

**c) Fix hardcoded colour** (lines ~824 and ~891):
```go
// Before:
lipgloss.Color("#3d59a1")

// After:
lipgloss.Color(theme.Special.Highlight)
```

> Read the exact context of each hardcoded line before editing — make sure `theme` is in scope at that point.

- [ ] **Step 10.2: Verify**

```
go build github.com/useful-toys/abditum/internal/tui
```

- [ ] **Step 10.3: Commit**

```
git add internal/tui/filepicker.go
git commit -m "refactor: filepicker uses nested Theme fields; replace hardcoded highlight colour"
```

---

## Task 11: Update `flow_*.go` files

**Files:**
- Modify: `internal/tui/flow_create_vault.go`
- Modify: `internal/tui/flow_open_vault.go`
- Modify: `internal/tui/flow_save_and_exit.go`

All three files use `MsgKind` constants with `messages.Show(...)`.

- [ ] **Step 11.1: Edit `flow_create_vault.go`**

- `MsgBusy` → `MessageBusy`
- `MsgError` → `MessageError`
- `MsgSuccess` → `MessageSuccess`

- [ ] **Step 11.2: Edit `flow_open_vault.go`**

- `MsgBusy` → `MessageBusy`
- `MsgError` → `MessageError`
- `MsgSuccess` → `MessageSuccess`

- [ ] **Step 11.3: Edit `flow_save_and_exit.go`**

- `MsgError` → `MessageError`

- [ ] **Step 11.4: Verify full build and tests**

```
go build ./internal/tui/...
go test ./internal/tui/... -count=1
```

Expected: build passes; tests may still fail due to test files not yet updated (next task). Build must not error.

- [ ] **Step 11.5: Commit**

```
git add internal/tui/flow_create_vault.go internal/tui/flow_open_vault.go internal/tui/flow_save_and_exit.go
git commit -m "refactor: flow files use MessageKind constants"
```

---

## Task 12: Update all test files

**Files:**
- Modify: `internal/tui/messages_test.go`
- Modify: `internal/tui/welcome_test.go`
- Modify: `internal/tui/root_test.go`
- Modify: `internal/tui/passwordentry_test.go`
- Modify: `internal/tui/passwordcreate_test.go`
- Modify: `internal/tui/actions_test.go`
- Modify: `internal/tui/filepicker_test.go`
- Modify: `internal/tui/exit_flow_integration_test.go`
- Modify: `internal/tui/flow_open_vault_test.go`
- Modify: `internal/tui/flow_create_vault_test.go`
- Modify: `internal/tui/flow_save_and_exit_test.go`

- [ ] **Step 12.1: Apply renames across all test files**

In every test file, perform these substitutions (use your editor's global find-and-replace — all files are in the same directory):

| Find | Replace |
|---|---|
| `ThemeTokyoNight` | `TokyoNight` |
| `ThemeCyberpunk` | `Cyberpunk` |
| `MsgSuccess` | `MessageSuccess` |
| `MsgInfo` | `MessageInfo` |
| `MsgWarn\b` | `MessageWarning` |
| `MsgWarning` | `MessageWarning` |
| `MsgError` | `MessageError` |
| `MsgBusy` | `MessageBusy` |
| `MsgHint` | `MessageHint` |
| `MsgKind` | `MessageKind` |

> Use `\b` word boundary when searching to avoid partial matches. In `messages_test.go` there are string literals like `"expected kind MsgSuccess..."` — update those too for consistency with the new names.

- [ ] **Step 12.2: Verify all tests pass**

```
go test ./internal/tui/... -count=1 -v 2>&1 | tail -40
```

Expected: all tests PASS. If any test fails, read the error and fix.

- [ ] **Step 12.3: Commit**

```
git add internal/tui/*_test.go
git commit -m "test: update all test files to MessageKind and TokyoNight/Cyberpunk names"
```

---

## Task 13: Delete the four old files

**Files:**
- Delete: `internal/tui/theme.go`
- Delete: `internal/tui/tokens.go`
- Delete: `internal/tui/theme/theme.go`
- Delete: `internal/tui/tokens/tokens.go`

- [ ] **Step 13.1: Delete old files**

```
git rm internal/tui/theme.go
git rm internal/tui/tokens.go
git rm internal/tui/theme/theme.go
git rm internal/tui/tokens/tokens.go
```

If the `theme/` and `tokens/` directories are now empty (only contained those files), git will remove them automatically.

- [ ] **Step 13.2: Verify build and tests**

```
go build ./internal/tui/...
go vet ./internal/tui/...
go test ./internal/tui/... -count=1
```

Expected: all three commands pass with no errors.

- [ ] **Step 13.3: Verify no old names remain**

```powershell
Select-String -Path "internal\tui\*.go" -Pattern "ThemeTokyoNight|ThemeCyberpunk" -Recurse
Select-String -Path "internal\tui\*.go" -Pattern "\bMsgWarn\b" -Recurse
Select-String -Path "internal\tui\*.go" -Pattern "tui/theme\"|tui/tokens\"" -Recurse
```

Expected: all three commands return no output (zero matches).

- [ ] **Step 13.4: Commit**

```
git add -A
git commit -m "refactor: delete obsolete theme.go, tokens.go and subpackages"
```

---

## Task 14: Final validation

- [ ] **Step 14.1: Full build**

```
go build ./...
```

Expected: PASS.

- [ ] **Step 14.2: Vet**

```
go vet ./internal/tui/...
```

Expected: no warnings.

- [ ] **Step 14.3: Full test suite**

```
go test ./... -count=1
```

Expected: all tests pass.

- [ ] **Step 14.4: Confirm token coverage**

Open `internal/tui/design.go` and verify the following are present (visual scan):
- `SurfaceTokens` with `Base`, `Raised`, `Input`
- `TextTokens` with `Primary`, `Secondary`, `Disabled`, `Link`
- `BorderTokens` with `Default`, `Focused`
- `AccentTokens` with `Primary`, `Secondary`
- `SemanticTokens` with `Success`, `Warning`, `Error`, `Info`, `Off`
- `SpecialTokens` with `Muted`, `Highlight`, `Match`
- `Logo [5]string`
- `Typography` with `Bold`, `Dim`, `Italic`, `Underline`, `Strikethrough`
- `MessageKind` with `MessageSuccess`, `MessageInfo`, `MessageWarning`, `MessageError`, `MessageBusy`, `MessageHint`
- `KeyCtrl`, `KeyShift`, `KeyAlt`
- All symbol groups present

- [ ] **Step 14.5: Final commit if any cleanup needed**

```
git add -A
git commit -m "chore: final cleanup after design system consolidation"
```
