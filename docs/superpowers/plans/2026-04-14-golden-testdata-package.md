# Golden Testdata Package Implementation Plan

**Document ID:** 2026-04-14-golden-testdata-package  
**Date:** April 14, 2026  
**Specification:** See `docs/superpowers/specs/2026-04-14-golden-testdata-package-design.md`

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

## Goal

Implement the `internal/tui/testdata` package (subpackage of `internal/tui`) with reusable helpers for golden file testing of TUI components. The package provides:

1. **ANSI parser** (`ansiparser.go`): Extract visual style transitions from ANSI output
2. **Testing helpers** (`helpers.go`): Golden file management and component rendering utilities
3. **Unit tests** (`ansiparser_test.go`): Validate parser correctness

## Architecture

Three files, no circular dependencies:
- `internal/tui/testdata/ansiparser.go` — ANSI parsing + `StyleTransition` type
- `internal/tui/testdata/ansiparser_test.go` — Unit tests for parser
- `internal/tui/testdata/helpers.go` — Golden file helpers + `RenderFn` type

## Tech Stack

- **Language:** Go 1.21+
- **Dependencies:** 
  - `internal/tui/design` (for `*design.Theme`)
  - Go standard library only (encoding/json, flag, os, path/filepath, strconv, strings, testing, regexp)

---

## Task 1: ansiparser.go

**Files:**
- Create: `internal/tui/testdata/ansiparser.go`

**Objective:** Implement ANSI escape sequence parser that extracts visual style transitions.

### Implementation Steps

- [ ] **Step 1: Define StyleTransition type**

Create struct with fields: `Line int`, `Col int`, `FG *string`, `BG *string`, `Style []string`.  
Add JSON struct tags for serialization.

```go
type StyleTransition struct {
	Line  int      `json:"line"`
	Col   int      `json:"col"`
	FG    *string  `json:"fg"`
	BG    *string  `json:"bg"`
	Style []string `json:"style"`
}
```

- [ ] **Step 2: Implement StyleTransition.MarshalJSON()**

Serialize as 5-tuple: `[line, col, fg, bg, styles]`.  
Ensure `Style` is always `[]` (not `null`) when empty.

```go
func (st StyleTransition) MarshalJSON() ([]byte, error) {
	type tuple [5]any
	var fg, bg any
	if st.FG != nil {
		fg = *st.FG
	}
	if st.BG != nil {
		bg = *st.BG
	}
	style := st.Style
	if style == nil {
		style = []string{}
	}
	return json.Marshal(tuple{st.Line, st.Col, fg, bg, style})
}
```

- [ ] **Step 3: Implement marshalStyleChanges()**

Serialize slice of `StyleTransition` to indented JSON with each tuple on one line.

```go
func marshalStyleChanges(transitions []StyleTransition) ([]byte, error) {
	if transitions == nil {
		return []byte("null"), nil
	}
	if len(transitions) == 0 {
		return []byte("[]"), nil
	}
	var buf strings.Builder
	buf.WriteString("[\n")
	for i, t := range transitions {
		b, err := json.Marshal(t)
		if err != nil {
			return nil, err
		}
		buf.WriteString("  ")
		buf.Write(b)
		if i < len(transitions)-1 {
			buf.WriteByte(',')
		}
		buf.WriteByte('\n')
	}
	buf.WriteString("]")
	return []byte(buf.String()), nil
}
```

- [ ] **Step 4: Define ansiState struct and helper functions**

```go
type ansiState struct {
	fg    *string
	bg    *string
	style map[string]bool
}
```

Implement:
- `stateKey(s ansiState) string` — Returns unique key for state comparison
- `applyCode(state *ansiState, code int)` — Applies SGR code to state (handle codes 0-107)
- `colorCode16(code int) string` — Converts color code to hex (16-color palette)
- `styleMapToArray(styleMap map[string]bool) []string` — Converts style set to sorted array

- [ ] **Step 5: Implement ansiToStyleChanges()**

Parse ANSI sequences using regex `\x1b\[([0-9;]*)m`.  
Maintain line/col position as text is consumed.  
Track state changes and record transitions only when state changes.  
Deduplicate: if multiple SGR codes at same position, replace previous transition with combined state.

```go
func ansiToStyleChanges(output string) []StyleTransition {
	// Initialize result, currentState, lastStateKey
	// Regex to find sequences
	// Loop through matches:
	//   - Advance line/col through text before sequence
	//   - Parse SGR codes and apply to state
	//   - If state changed, record transition (or replace if same position)
	// Return result
}
```

- [ ] **Step 6: Verify build**

Run: `go build ./internal/tui/testdata`  
Expected: No errors

- [ ] **Step 7: Commit**

```bash
git add internal/tui/testdata/ansiparser.go
git commit -m "feat: add ansiparser.go with style transition parsing"
```

Expected: Commit succeeds, no build errors

---

## Task 2: ansiparser_test.go

**Files:**
- Create: `internal/tui/testdata/ansiparser_test.go`

**Objective:** Implement unit tests for ANSI parser. Validate parser correctness and edge cases.

### Implementation Steps

- [ ] **Step 1: Write test for empty input**

```go
func TestAnsiToStyleChanges_EmptyInput(t *testing.T) {
	result := ansiToStyleChanges("")
	if len(result) != 0 {
		t.Errorf("expected empty result for empty input, got %d transitions", len(result))
	}
}
```

- [ ] **Step 2: Write test for plain text**

```go
func TestAnsiToStyleChanges_PlainText(t *testing.T) {
	result := ansiToStyleChanges("Hello, world!")
	if len(result) != 0 {
		t.Errorf("expected no transitions for plain text, got %d", len(result))
	}
}
```

- [ ] **Step 3: Write test for single color (bold)**

Test sequence `\x1b[1m` (bold on) and `\x1b[0m` (reset).  
Verify first transition has `bold` in style, second has empty style.

```go
func TestAnsiToStyleChanges_SingleColor(t *testing.T) {
	input := "\x1b[1mBold text\x1b[0m"
	result := ansiToStyleChanges(input)
	// Assert len(result) == 2
	// Assert result[0].Style contains "bold"
	// Assert result[1].Style is empty
}
```

- [ ] **Step 4: Write test for line tracking**

Test input with newline: `"Line 1\n\x1b[1mLine 2"`.  
Verify transition occurs on line 1 (second line, 0-indexed).

- [ ] **Step 5: Write test for multiple transitions**

Test: `"\x1b[1mA\x1b[0mB\x1b[1mC"`.  
Verify at least 3 transitions are recorded (bold on, reset, bold on again).

- [ ] **Step 6: Write test for combined state with multiple SGR codes**

**Key test case:** Multiple SGR codes at same position result in single combined transition.

Test input: `"\x1b[1;31mText\x1b[0m"` (bold + red foreground at same position).

```go
func TestAnsiToStyleChanges_CombinedStateMultipleSGR(t *testing.T) {
	input := "\x1b[1;31mText\x1b[0m"
	result := ansiToStyleChanges(input)
	// Assert len(result) >= 1
	// Assert result[0].Style contains "bold"
	// Assert result[0].FG == "#800000" (red color)
	// Verify second transition is reset (empty style, nil FG)
}
```

- [ ] **Step 7: Run tests to verify they pass**

Run: `go test ./internal/tui/testdata -v`  
Expected: All tests pass (PASS)

- [ ] **Step 8: Verify build still works**

Run: `go build ./internal/tui/testdata`  
Expected: No errors

- [ ] **Step 9: Commit**

```bash
git add internal/tui/testdata/ansiparser_test.go
git commit -m "feat: add ansiparser tests with SGR code validation"
```

Expected: Commit succeeds, all tests pass

---

## Task 3: helpers.go

**Files:**
- Create: `internal/tui/testdata/helpers.go`

**Objective:** Implement golden file testing helpers. Provide utilities for comparing rendered output against saved golden files.

### Implementation Steps

- [ ] **Step 1: Define update flag and RenderFn type**

```go
var update = flag.Bool("update-golden", false, "update golden files")

type RenderFn func(w, h int, theme *design.Theme) string
```

- [ ] **Step 2: Implement goldenPath()**

Constructs path like `testdata/golden/component_variant_size.ext`.

```go
func goldenPath(component, variant, size, ext string) string {
	filename := component + "_" + variant + "_" + size + "." + ext
	return filepath.Join("testdata", "golden", filename)
}
```

- [ ] **Step 3: Implement ansiToCleanText()**

Remove all ANSI escape sequences from string.  
Track ESC character (`\x1b`) and skip everything until letter that ends sequence.

```go
func ansiToCleanText(s string) string {
	cleaned := ""
	inEscape := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= 'a' && s[i] <= 'z') {
				inEscape = false
			}
			continue
		}
		cleaned += string(s[i])
	}
	return cleaned
}
```

- [ ] **Step 4: Implement checkOrUpdateGolden()**

If `-update-golden` set: write/update file.  
Otherwise: read file and compare, fail test if different.

```go
func checkOrUpdateGolden(t *testing.T, path, got string) {
	t.Helper()
	if *update {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("failed to create golden directory: %v", err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatalf("failed to write golden file: %v", err)
		}
		return
	}
	expected, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("golden file not found: %s (run with -update-golden to create it)", path)
		}
		t.Fatalf("failed to read golden file: %v", err)
	}
	if string(expected) != got {
		t.Errorf("output does not match golden file %s", path)
	}
}
```

- [ ] **Step 5: Implement parseSize()**

Parse size strings like "30" (30x30) or "80x24" (80 width, 24 height).

```go
func parseSize(size string) (int, int, error) {
	if strings.Contains(size, "x") || strings.Contains(size, "X") {
		parts := strings.FieldsFunc(size, func(r rune) bool { return r == 'x' || r == 'X' })
		if len(parts) != 2 {
			return 0, 0, invalid("invalid size format")
		}
		w, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		h, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err1 != nil || err2 != nil {
			return 0, 0, invalid("invalid size values")
		}
		return w, h, nil
	}
	dim, err := strconv.Atoi(strings.TrimSpace(size))
	if err != nil {
		return 0, 0, invalid("invalid size")
	}
	return dim, dim, nil
}
```

Also define `parseError` struct and `invalid()` helper.

- [ ] **Step 6: Implement TestRenderManaged()**

Run parameterized tests for multiple sizes.  
For each size:
  - Parse size string
  - Render component using `RenderFn`
  - Call `checkOrUpdateGolden()` to compare/update
  
Use `design.TokyoNight` as default theme.

```go
func TestRenderManaged(t *testing.T, component, variant string, sizes []string, render RenderFn) {
	t.Helper()
	theme := design.TokyoNight
	for _, sizeStr := range sizes {
		t.Run(sizeStr, func(t *testing.T) {
			w, h, err := parseSize(sizeStr)
			if err != nil {
				t.Fatalf("failed to parse size: %v", err)
			}
			output := render(w, h, theme)
			path := goldenPath(component, variant, sizeStr, "txt")
			checkOrUpdateGolden(t, path, output)
		})
	}
}
```

- [ ] **Step 7: Implement TestRenderManual()**

Placeholder for manual validation (no persistent golden file).  
Document intended usage.

```go
func TestRenderManual(t *testing.T, component, variant, size, output string) {
	t.Helper()
	// Placeholder for manual golden file validation
	// In practice: rendered := renderMyComponent(...); TestRenderManual(t, "comp", "var", "30x24", rendered)
}
```

- [ ] **Step 8: Verify build**

Run: `go build ./internal/tui/testdata`  
Expected: No errors

- [ ] **Step 9: Run existing tests to confirm no regressions**

Run: `go test ./internal/tui/testdata -v`  
Expected: All 6 tests pass (5 from ansiparser_test.go + 1 placeholder)

- [ ] **Step 10: Commit**

```bash
git add internal/tui/testdata/helpers.go
git commit -m "feat: add golden file testing helpers (ansiToCleanText, goldenPath, TestRenderManaged, etc)"
```

Expected: Commit succeeds, all tests pass, no build errors

---

## Task 4: Documentation

**Files:**
- Create: `docs/superpowers/specs/2026-04-14-golden-testdata-package-design.md`
- Create: `docs/superpowers/plans/2026-04-14-golden-testdata-package.md`

**Objective:** Document the specification and implementation plan for reference and future maintenance.

### Implementation Steps

- [ ] **Step 1: Create specification document**

File: `docs/superpowers/specs/2026-04-14-golden-testdata-package-design.md`

Include sections:
- Overview and goals
- Architecture (3 modules: ansiparser, helpers, tests)
- Key types and functions
- No circular dependencies
- Module path
- Testing strategy
- Design decisions
- Future extensibility
- References

- [ ] **Step 2: Create plan document**

File: `docs/superpowers/plans/2026-04-14-golden-testdata-package.md`

This file (copy spec header + all tasks with substeps).

- [ ] **Step 3: Verify files exist and are readable**

Run: `ls -la docs/superpowers/specs/ docs/superpowers/plans/`  
Expected: Both .md files present

- [ ] **Step 4: Commit both documentation files**

```bash
git add docs/superpowers/specs/2026-04-14-golden-testdata-package-design.md docs/superpowers/plans/2026-04-14-golden-testdata-package.md
git commit -m "docs: add spec and plan for golden testdata package"
```

Expected: Commit succeeds

---

## Verification at the End

After all 4 tasks complete:

- [ ] **Step 1: Build the testdata package**

Run: `go build ./internal/tui/testdata`  
Expected: Zero errors

- [ ] **Step 2: Run all testdata tests**

Run: `go test ./internal/tui/testdata -v`  
Expected: All tests pass (PASS, 0 failures)

- [ ] **Step 3: Run full test suite to confirm no regressions**

Run: `go test ./...`  
Expected: All tests pass

- [ ] **Step 4: Verify all commits were created**

Run: `git log --oneline -10`  
Expected: See commits:
  1. "feat: add ansiparser.go with style transition parsing"
  2. "feat: add ansiparser tests with SGR code validation"
  3. "feat: add golden file testing helpers (ansiToCleanText, goldenPath, TestRenderManaged, etc)"
  4. "docs: add spec and plan for golden testdata package"

- [ ] **Step 5: Check worktree status**

Run: `git status`  
Expected: "working tree clean"

- [ ] **Final: Worktree ready for merge or cleanup**

Location: `.worktrees/golden-testdata-package/`  
Branch: `feature/golden-testdata-package`  
Ready for: PR creation or merge to main

---

## Notes

- **Compilation before commit:** Each task must compile successfully before commit
- **Tests before commit:** All tests must pass before commit
- **No partial commits:** Never commit if build or tests fail
- **Worktree isolation:** All work happens in `.worktrees/golden-testdata-package/`
- **No force push:** This is a feature branch, never force push
- **Documentation files:** These are supporting documentation only, not part of the production codebase
