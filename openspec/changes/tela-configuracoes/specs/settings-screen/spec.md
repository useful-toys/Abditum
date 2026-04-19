## ADDED Requirements

### Requirement: Settings screen SHALL render as a structured first-class work area
The system SHALL replace the `"Settings"` placeholder in `WorkAreaSettings` with a centered, minimalist layout organized in labeled groups with vertical centering, consistent with the design system tokens, workspace dimensions, and interaction conventions of the TUI.

#### Scenario: Settings mode is activated
- **WHEN** the application switches to `WorkAreaSettings`
- **THEN** the work area MUST render the settings screen with a "Configurações" title, at minimum one visible group and one selectable item, replacing the raw `"Settings"` text

#### Scenario: Settings layout respects workspace boundaries
- **WHEN** the Settings screen is rendered at any terminal size
- **THEN** it MUST use the available workspace height and width, apply tokens from the active theme, and preserve the outer layout formed by header, work area, message line, and command bar

#### Scenario: Minimum-height fallback remains centralized in root
- **WHEN** the terminal height is below the minimum supported height
- **THEN** `SettingsView` MUST NOT introduce its own fallback rendering rule, because the root-level guard remains the single authority for blocking normal screen rendering in that condition

#### Scenario: Content is vertically centered in the work area
- **WHEN** the Settings screen is rendered
- **THEN** the settings content (title + groups + items) MUST be vertically centered in the available work area using symmetric blank-line padding above and below

### Requirement: Settings screen SHALL organize preferences into labeled groups
The system SHALL present settings items in named groups — at minimum: Aparência (theme), Segurança (timers), and Sobre (read-only app info) — separated by blank lines. Groups SHALL NOT use horizontal divider lines; visual hierarchy is established by bold group heading and indentation of items.

#### Scenario: Groups are visually distinguished
- **WHEN** the Settings screen is rendered
- **THEN** each group MUST have a bold section heading, each item MUST be indented relative to the heading, and groups MUST be separated by at least one blank line

#### Scenario: Each item shows name, current value and description inline
- **WHEN** an item in the settings list has keyboard focus
- **THEN** the screen MUST display the item name and its current value on the same line, and a contextual description on the immediately following line — using at least two visual layers (position and/or typographic weight, not color alone) to distinguish them

### Requirement: Settings screen SHALL support keyboard navigation
The system SHALL allow the user to move focus between settings items using ↑ and ↓, activate edit mode on numeric items using Enter, and cancel editing using Esc. The theme item is focusable but never enters edit mode.

#### Scenario: Moving focus down
- **WHEN** the user presses ↓ while a settings item is focused
- **THEN** focus MUST move to the next item, wrapping to the first item when already at the last

#### Scenario: Moving focus up
- **WHEN** the user presses ↑ while a settings item is focused
- **THEN** focus MUST move to the previous item, wrapping to the last item when already at the first

#### Scenario: Entering edit mode on a numeric item
- **WHEN** the user presses Enter on a numeric-type item that is NOT in edit mode
- **THEN** the item MUST enter edit mode, showing an editable inline field pre-filled with the current numeric value and a real terminal cursor positioned at end of buffer

#### Scenario: Only digits accepted in edit mode
- **WHEN** a numeric item is in edit mode and the user presses a non-digit key (excluding Backspace, Enter, Esc)
- **THEN** the keystroke MUST be silently ignored and the buffer MUST remain unchanged

#### Scenario: Confirming a numeric edit
- **WHEN** the user presses Enter while a numeric item is in edit mode and the value is within the valid range
- **THEN** the new value MUST be applied, persisted in-memory, and the item MUST exit edit mode

#### Scenario: Cancelling a numeric edit
- **WHEN** the user presses Esc while a numeric item is in edit mode
- **THEN** the original value MUST be restored and the item MUST exit edit mode without changes

#### Scenario: Adjusting a numeric item with +/-
- **WHEN** a numeric item has focus and is NOT in edit mode
- **THEN** pressing `+` MUST increment the value by exactly 5 seconds and pressing `-` MUST decrement it by exactly 5 seconds, applying the same range validation as edit mode

#### Scenario: +/- ignored in edit mode
- **WHEN** a numeric item is in edit mode
- **THEN** pressing `+` or `-` MUST be treated as a non-digit keystroke and silently ignored

#### Scenario: Enter ignored on theme item
- **WHEN** the user presses Enter while the theme item is focused
- **THEN** the theme item MUST NOT enter edit mode and the screen MUST remain unchanged

### Requirement: Settings screen SHALL expose theme selection
The system SHALL present the active theme name as a focusable, non-editable item in the Aparência group, while preserving `F12` as the standard interaction for changing themes anywhere in the application.

#### Scenario: Active theme is visible
- **WHEN** the Settings screen is shown
- **THEN** the Aparência group MUST display the currently active theme name as the current value of the theme item

#### Scenario: Theme item receives focus
- **WHEN** focus moves to the theme item
- **THEN** the item MUST become visually focused without entering edit mode and MUST remain non-editable from Enter

#### Scenario: F12 theme toggle reflects in settings
- **WHEN** the active theme is changed via F12 while the Settings screen is visible
- **THEN** the theme item in the Aparência group MUST display the updated theme name without requiring the user to leave and re-enter the mode

#### Scenario: Theme change dirties the cofre state
- **WHEN** the active theme is changed while a cofre is open
- **THEN** the new theme identifier MUST be written into the in-memory cofre configuration and the cofre MUST become modified until saved

### Requirement: Settings screen SHALL expose security timer configuration
The system SHALL allow the user to configure three security timers, all stored and displayed in **seconds**, adjusted in steps of 5 seconds:
- Auto-lock timeout: > 60 s (default 300 s)
- Sensitive field hide timeout: > 2 s (default 15 s)
- Clipboard clear timeout: > 10 s (default 30 s)

#### Scenario: Timer current value is visible
- **WHEN** the Settings screen is shown
- **THEN** the Segurança group MUST display the current configured value for each timer in seconds alongside its label

#### Scenario: +/- adjusts timer in steps of 5
- **WHEN** a timer item has focus and is NOT in edit mode and the user presses `+` or `-`
- **THEN** the value MUST change by exactly 5 seconds in the respective direction, and the change MUST NOT be applied if the resulting value would violate the minimum constraint

#### Scenario: Timer input accepts digits only
- **WHEN** a timer item enters edit mode
- **THEN** the editable portion of the field MUST accept only numeric digits, while the unit `s` remains non-editable outside the input area

#### Scenario: Timer field shows input treatment during editing
- **WHEN** a timer item is in edit mode
- **THEN** the editable portion MUST use the input-field treatment defined by the design system, including differentiated `surface.input` background and visible cursor

#### Scenario: Entering an out-of-range timer value
- **WHEN** the user confirms a timer edit with a value that does not satisfy the minimum constraint
- **THEN** the system MUST reject the value, keep the field in edit mode, and display an error message in the message bar indicating the minimum valid value

#### Scenario: Timer change dirties the cofre state
- **WHEN** the user applies a valid timer value while a cofre is open
- **THEN** the new timer MUST be written into the in-memory cofre configuration and the cofre MUST become modified until saved

### Requirement: Settings screen SHALL show read-only app information
The system SHALL display a Sobre group with read-only application and cofre information, using items that can receive focus for inspection but cannot be activated or edited.

#### Scenario: Sobre group is rendered
- **WHEN** the Settings screen is shown
- **THEN** the Sobre group MUST contain at minimum the application version string and one additional read-only contextual value, both non-editable

#### Scenario: Read-only item receives focus
- **WHEN** focus moves to a read-only item in the Sobre group
- **THEN** the item MUST update the contextual description and hint text, but MUST ignore activation keys that would start editing

### Requirement: Settings screen SHALL react immediately to theme changes
The system SHALL re-render the Settings screen with the correct tokens whenever the active theme changes, regardless of whether the change originated from settings navigation or from the global F12 shortcut.

#### Scenario: Theme changes while settings are visible
- **WHEN** the active theme changes while `WorkAreaSettings` is the active work area
- **THEN** the Settings screen MUST update all colors, highlights, and emphasis without requiring navigation away from the mode

### Requirement: Settings screen SHALL use message bar feedback for focus and validation
The system SHALL use the TUI message bar to guide the user while navigating or editing settings, reusing the existing hint and feedback model instead of introducing inline validation widgets.

#### Scenario: Theme field shows focus hint
- **WHEN** the theme item has focus
- **THEN** the message bar MUST show a hint equivalent to `F12 para alternar tema visual`

#### Scenario: Numeric field shows focus hint when not editing
- **WHEN** a timer item has focus and is NOT in edit mode
- **THEN** the message bar MUST show a hint equivalent to `Enter edita · +/- altera o valor`

#### Scenario: Numeric field shows edit hint when in edit mode
- **WHEN** a timer item enters edit mode
- **THEN** the message bar MUST update to a hint equivalent to `Enter confirma · Esc cancela`

#### Scenario: Rejected setting shows local feedback
- **WHEN** the cofre domain rejects a setting change requested from the Settings screen
- **THEN** the message bar MUST show an error message and the Settings screen MUST remain responsible for keeping the editing state consistent

### Requirement: Settings mutations SHALL be applied locally by the settings view through the cofre domain
The system SHALL treat the Settings screen as a self-contained feature screen for synchronous setting edits, with `SettingsView` validating local input and applying persisted configuration changes directly through `vault.Manager`, while `RootModel` remains responsible only for global layout, screen routing, and theme shortcut handling.

#### Scenario: Settings view applies a persisted change directly
- **WHEN** the user applies a persisted setting from the Settings screen
- **THEN** the view MUST call the cofre-domain service directly, the mutation MUST be applied to the in-memory configuration on success, and any domain rejection MUST be handled in the Settings screen without requiring a root-level orchestration message

#### Scenario: Direct application does not emit application event
- **WHEN** a synchronous setting change is applied successfully from the Settings screen
- **THEN** the system MUST NOT require a dedicated application-wide `tea.Msg` or equivalent event solely to propagate that field mutation

#### Scenario: Future consumers re-read canonical configuration
- **WHEN** another part of the application needs updated settings after a synchronous change
- **THEN** it SHOULD obtain the canonical values from the cofre domain (`vault.Manager` / `Configuracoes`) instead of depending on a generic broadcast mechanism introduced by the Settings screen

### Requirement: Visual specification for settings SHALL be documented in golden/tui-spec-telas.md
The project SHALL define the Settings screen in `golden/tui-spec-telas.md`, following the documentation pattern used for other screens and respecting rules from `golden/tui-design-system.md`.

#### Scenario: Golden spec includes Modo Configurações section
- **WHEN** this change is implemented
- **THEN** `golden/tui-spec-telas.md` MUST include the "Modo Configurações" section with wireframe, visual identity table, states table, messages table, events table, and behavior list

### Requirement: Settings screen SHALL be protected by golden-file rendering tests
The project SHALL include golden-file tests for the Settings screen so regressions in grouping, focus, editing treatment, and overall composition are detected automatically.

#### Scenario: Structured screen render has golden coverage
- **WHEN** the Settings screen is rendered in its default structured state
- **THEN** the test suite MUST compare the rendered output against a committed golden file

#### Scenario: Editing state has golden coverage
- **WHEN** the Settings screen renders a numeric field in edit mode
- **THEN** the test suite MUST compare the rendered output against a committed golden file that includes input treatment and cursor state
