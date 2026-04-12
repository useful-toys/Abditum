# Design: SetSize Semantics Refactor (Phase 1)

**Date:** 2026-04-12  
**Author:** Team  
**Status:** Approved for Phase 1 Implementation

---

## Overview

This refactor improves semantic clarity in the TUI layer by distinguishing between two fundamentally different sizing concepts:

- **childModel.SetSize(width, height)** → Component receives **exact** space it will occupy
- **modalView.SetAvailableSize(maxWidth, maxHeight)** → Modal receives **maximum** space available

This is a **Phase 1 interface refactor only**. Implementation logic in View/Render remains unchanged. Phase 2 will refactor the calculation logic.

---

## Problem

Currently, both `childModel` and `modalView` use the same method name `SetSize(w, h)`, obscuring their different contracts:

1. **childModel** must use exactly the space given (e.g., `vaultTree.SetSize(halfW, workH)` means vaultTree occupies exactly `halfW` width)
2. **modalView** receives maximum available space but may use less (e.g., `DecisionDialog` uses ~50 chars within available 80)

This ambiguity makes it unclear:
- Whether a component must fill its space or can be smaller
- Whether dimensions are constraints or limits

---

## Solution: Phase 1 (Interface Semantics)

### 1. Update Interface Definitions

**flows.go:**

```go
// childModel represents a work-area component that occupies exact space.
type childModel interface {
	Update(tea.Msg) tea.Cmd
	View() string
	// SetSize(width, height) stores exact allocated dimensions.
	// Component MUST use exactly this space.
	SetSize(width, height int)
	ApplyTheme(*Theme)
}

// modalView represents an overlay modal with maximum available space.
type modalView interface {
	Update(tea.Msg) tea.Cmd
	View() string
	Shortcuts() []Shortcut
	// SetAvailableSize(maxWidth, maxHeight) stores maximum available dimensions.
	// Modal MAY use less (e.g., center and pad with empty space).
	SetAvailableSize(maxWidth, maxHeight int)
}
```

### 2. Rename Parameter Names (for clarity)

In all implementations, use descriptive parameter names matching intent:

**childModel implementations:**
```go
func (m *welcomeModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
```

**modalView implementations:**
```go
func (d *DecisionDialog) SetAvailableSize(maxWidth, maxHeight int) {
	d.width = maxWidth
	d.height = maxHeight
}
```

### 3. Update All Call Sites

**rootModel calls:**
- childModel: `child.SetSize(width, height)`
- modalView: `modal.SetAvailableSize(maxWidth, maxHeight)`

**Example - renderVaultArea():**
```go
halfW := m.width / 2
m.vaultTree.SetSize(halfW, workH)              // Exact space
m.secretDetail.SetSize(m.width-halfW, workH)   // Exact space
```

**Example - View() modal rendering:**
```go
workH := m.height - headerH - msgBarH - cmdBarH
top.SetAvailableSize(m.width, workH)  // Maximum available space
content = m.renderFrame(top)
```

---

## Affected Files

### Interface Definition
- `internal/tui/flows.go` - Update `childModel` and `modalView` interfaces

### Implementation Updates

**childModel implementations:**
- `internal/tui/welcome.go` - SetSize(width, height)
- `internal/tui/vaulttree.go` - SetSize(width, height)
- `internal/tui/secretdetail.go` - SetSize(width, height)
- `internal/tui/templatelist.go` - SetSize(width, height)
- `internal/tui/templatedetail.go` - SetSize(width, height)
- `internal/tui/settings.go` - SetSize(width, height)

**modalView implementations:**
- `internal/tui/modal.go` - SetAvailableSize(maxWidth, maxHeight)
- `internal/tui/help.go` - SetAvailableSize(maxWidth, maxHeight)
- `internal/tui/decision.go` - SetAvailableSize(maxWidth, maxHeight)

### Call Site Updates
- `internal/tui/root.go` - All SetSize() and SetAvailableSize() calls

### Tests
- `internal/tui/help_test.go` - SetAvailableSize()
- `internal/tui/decision_test.go` - SetAvailableSize()

---

## Implementation Details: NO LOGIC CHANGES

**Important:** This phase only changes method names and parameter names. Implementation logic stays identical.

### childModel Implementations
```go
// BEFORE
func (m *vaultTreeModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// AFTER (logic unchanged)
func (m *vaultTreeModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
```

### modalView Implementations
```go
// BEFORE
func (d *DecisionDialog) SetSize(w, h int) {
	d.width = w
	d.height = h
}

// AFTER (logic unchanged, method renamed)
func (d *DecisionDialog) SetAvailableSize(maxWidth, maxHeight int) {
	d.width = maxWidth
	d.height = maxHeight
}
```

### View/Render Methods
- **NO CHANGES** to View() or Render() methods
- They continue to make calculations as before
- They continue to use `m.width` and `m.height` fields

---

## Naming Convention Justification

Follows Go idioms (from Effective Go):
- ✅ Full names for interface methods: `SetSize`, `SetAvailableSize`
- ✅ Descriptive parameter names: `width`, `height`, `maxWidth`, `maxHeight`
- ✅ No snake_case: using mixedCaps for multi-word names
- ✅ Stored in fields as `m.width`, `m.height` (unchanged)

---

## Phase 2 (Future)

Once this refactor is complete and tested, Phase 2 will:
1. Move calculation logic from View/Render into SetSize/SetAvailableSize
2. Store calculated values (e.g., `contentHeight`, `calculatedBoxWidth`) instead of raw dimensions
3. View/Render will only read pre-calculated state and render

This maintains clear separation:
- **SetSize/SetAvailableSize** = "planning phase" (calculate layout)
- **View/Render** = "rendering phase" (generate string)

---

## Testing Strategy

1. **Compile check:** Ensure all code compiles after interface changes
2. **Contract verification:** 
   - childModel: Verify components fill exact space
   - modalView: Verify modals center and don't exceed max dimensions
3. **Existing tests:** All existing tests should pass without modification
4. **Integration test:** Full TUI rendering with different terminal sizes

---

## Rollout Plan

1. Update interface definitions
2. Rename all method implementations
3. Update all call sites in rootModel
4. Update tests
5. Verify full TUI works
6. Commit with clear message describing semantic changes

---

## Success Criteria

- ✅ All code compiles
- ✅ All existing tests pass
- ✅ TUI renders correctly with different terminal sizes
- ✅ No logic changes, only interface/naming refactoring
- ✅ Code review confirms semantic clarity improvement
