# Golden Testdata Package Specification

**Document ID:** 2026-04-14-golden-testdata-package-design  
**Date:** April 14, 2026  
**Status:** Specification

## Overview

The `internal/tui/testdata` package provides reusable helpers for golden file testing of TUI (terminal user interface) components. Golden file testing is a strategy where component output is captured and stored as a "golden" reference, allowing future renders to be validated against these references.

## Goals

1. **Centralize ANSI parsing logic** for extracting and recording style transitions (color, font styling changes)
2. **Provide helper functions** for golden file management (read, write, compare)
3. **Enable parameterized testing** of TUI components across multiple screen sizes
4. **Support clean text extraction** for comparing visual content independent of ANSI sequences

## Architecture

The package is organized into three core modules:

### 1. ANSI Parser (`ansiparser.go`)

**Purpose:** Parse ANSI escape sequences and extract visual style transitions.

**Key Types:**
- `StyleTransition`: Represents a single point where visual styling (color, font) changes
  - `Line`: 0-indexed line number
  - `Col`: 0-indexed column where change occurs
  - `FG`: Foreground color in hex (nullable)
  - `BG`: Background color in hex (nullable)
  - `Style`: Array of active styles (bold, italic, underline, etc.)

**Key Functions:**
- `ansiToStyleChanges(output string) []StyleTransition`: Extracts all style changes from ANSI output
  - Normalizes SGR (Select Graphic Rendition) codes to canonical state
  - Deduplicates redundant transitions
  - Handles multiple SGR codes at the same position by recording the final combined state
- `marshalStyleChanges(transitions []StyleTransition) ([]byte, error)`: Serializes transitions to JSON
  - Format: Array of 5-tuples `[line, col, fg_hex_or_null, bg_hex_or_null, [styles]]`
  - Indented for readability but each tuple on a single line (diff-friendly)

**Internal Functions:**
- `applyCode(state *ansiState, code int)`: Applies a single SGR code to the current state
- `colorCode16(code int) string`: Converts ANSI 16-color codes to hex values
- `styleMapToArray(styleMap map[string]bool) []string`: Converts style set to ordered array

### 2. Testing Helpers (`helpers.go`)

**Purpose:** Provide utilities for golden file-based testing of TUI components.

**Key Types:**
- `RenderFn`: Function signature for rendering a component
  - Input: width (int), height (int), theme (`*design.Theme`)
  - Output: rendered string (typically with ANSI sequences)

**Key Variables:**
- `update`: Command-line flag (`-update-golden`) to create/update golden files instead of comparing

**Key Functions:**
- `goldenPath(component, variant, size, ext string) string`: Constructs filesystem path
  - Example: `goldenPath("vault_tree", "expanded", "80x24", "txt")` → `"testdata/golden/vault_tree_expanded_80x24.txt"`
- `ansiToCleanText(s string) string`: Removes all ANSI sequences
  - Returns visible text only
  - Useful for text-level validation independent of coloring
- `checkOrUpdateGolden(t *testing.T, path, got string)`: Compare or update golden file
  - If `-update-golden` flag set: writes/updates the file
  - Otherwise: compares `got` against file contents, fails test if different
- `parseSize(size string) (int, int, error)`: Parses size strings
  - Accepts "30" (square 30x30) or "80x24" (specific width×height)
- `TestRenderManaged(t *testing.T, component, variant string, sizes []string, render RenderFn)`: 
  - Runs parameterized golden tests for multiple sizes
  - Each size becomes a sub-test
  - Uses `design.TokyoNight` theme by default
  - Handles golden file creation/comparison via `checkOrUpdateGolden`
- `TestRenderManual(t *testing.T, component, variant, size, output string)`:
  - Placeholder for manual golden file validation
  - Intended for quick checks without persistent golden files

### 3. Unit Tests (`ansiparser_test.go`)

**Purpose:** Validate ANSI parser correctness.

**Test Cases:**
- `TestAnsiToStyleChanges_EmptyInput`: Empty input produces no transitions
- `TestAnsiToStyleChanges_PlainText`: Plain text with no sequences produces no transitions
- `TestAnsiToStyleChanges_SingleColor`: Bold sequence and reset are correctly tracked
- `TestAnsiToStyleChanges_LineTracking`: Line and column numbers advance correctly through newlines
- `TestAnsiToStyleChanges_MultipleTransitions`: Multiple style changes are all recorded
- `TestAnsiToStyleChanges_CombinedStateMultipleSGR`: Multiple SGR codes at same position result in single combined transition

## No Circular Dependencies

- The package imports only `design` (for the `Theme` type)
- `design` does not import `testdata`
- All other imports are from the Go standard library

## Module Path

`github.com/useful-toys/abditum/internal/tui/testdata`

## Testing Strategy

All golden file tests use the `-update-golden` flag workflow:

1. **First run (generate):** `go test -update-golden`
   - Renders component output
   - Writes to golden file
   - Stores for future validation

2. **Subsequent runs (validate):** `go test`
   - Renders component output
   - Compares against stored golden file
   - Fails if output differs

This allows:
- Visual regression detection (golden file changes when behavior changes)
- Documentation via example output (golden files show what the UI looks like)
- Diff-based review of rendering changes (visual diffs in PRs)

## Design Decisions

1. **State-based ANSI parsing**: Maintains current state (FG color, BG color, styles) and records transitions when state changes, not when sequences occur. This deduplicates redundant SGR codes and normalizes order-independence.

2. **Nullable colors**: Use `*string` for FG/BG to distinguish between "no color specified" (nil) and "reset to default" (still nil after explicit code 39/49). This enables proper serialization and comparison.

3. **Style set normalization**: Styles are stored in a map during parsing for O(1) lookup, then converted to sorted array for deterministic JSON output.

4. **Golden file paths**: Use simple naming convention (`component_variant_size.ext`) for predictable, readable paths. Sizes support both square and rectangular formats.

5. **Configurable golden updates**: Use flag-based approach (`-update-golden`) to avoid accidentally committing goldens. This follows Go test conventions.

## Future Extensibility

- Color space support (256-color, truecolor via codes 38:2, 48:2)
- Multiple theme parameterization in `TestRenderManaged`
- Snapshot-style comparison with inline diff reporting
- Integration with UI component libraries (once they exist)

## References

- ANSI Escape Code Reference: SGR (Select Graphic Rendition) — codes 0-97
- Go testing best practices: flag-based test configuration
- Golden file strategy: Industry-standard approach for regression testing
