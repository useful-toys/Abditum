# SetSize Semantics Refactor (Phase 1) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rename SetSize to SetAvailableSize for modals and use descriptive parameter names throughout, clarifying the semantic difference between exact child sizing and maximum modal sizing.

**Architecture:** Update interface definitions, then propagate changes through all implementations and call sites. No logic changes—only naming and method signatures. All tests should pass without modification.

**Tech Stack:** Go 1.x, standard library, existing TUI framework

---

## File Structure

### Interface Definition
- `internal/tui/flows.go` - Update `childModel` and `modalView` interfaces

### childModel Implementations (6 files)
- `internal/tui/welcome.go` - SetSize(width, height)
- `internal/tui/vaulttree.go` - SetSize(width, height)
- `internal/tui/secretdetail.go` - SetSize(width, height)
- `internal/tui/templatelist.go` - SetSize(width, height)
- `internal/tui/templatedetail.go` - SetSize(width, height)
- `internal/tui/settings.go` - SetSize(width, height)

### modalView Implementations (6 files)
- `internal/tui/modal.go` - SetAvailableSize(maxWidth, maxHeight)
- `internal/tui/help.go` - SetAvailableSize(maxWidth, maxHeight)
- `internal/tui/decision.go` - SetAvailableSize(maxWidth, maxHeight)
- `internal/tui/filepicker.go` - SetAvailableSize(maxWidth, maxHeight) [NEW - discovered]
- `internal/tui/passwordcreate.go` - SetAvailableSize(maxWidth, maxHeight) [NEW - discovered]
- `internal/tui/passwordentry.go` - SetAvailableSize(maxWidth, maxHeight) [NEW - discovered]

### Call Sites (rootModel)
- `internal/tui/root.go` - Update all SetSize() and new SetAvailableSize() calls

### Tests
- `internal/tui/help_test.go` - SetAvailableSize()
- `internal/tui/decision_test.go` - SetAvailableSize()
- `internal/tui/filepicker_test.go` - SetAvailableSize() [NEW - discovered]
- `internal/tui/passwordcreate_test.go` - SetAvailableSize() [NEW - discovered]
- `internal/tui/passwordentry_test.go` - SetAvailableSize() [NEW - discovered]

---

## Task 1: Update Interface Definitions

**Files:**
- Modify: `internal/tui/flows.go:8-53`

- [ ] **Step 1: Read current interface definitions**

Run: `cd C:\git\Abditum-T2 && type internal\tui\flows.go | findstr /A:32 "type childModel\|type modalView" -A 10`

Expected: See both interface definitions with current SetSize methods

- [ ] **Step 2: Update childModel interface**

In `internal/tui/flows.go`, update the `childModel` interface (around line 8-25):

```go
type childModel interface {
	Update(tea.Msg) tea.Cmd
	View() string
	// SetSize(width, height) stores exact allocated dimensions for this component.
	// Component MUST occupy exactly this space.
	SetSize(width, height int)
	ApplyTheme(*Theme)
}
```

Change:
- `SetSize(w, h int)` → `SetSize(width, height int)` (parameter names only)
- Update comment to mention "exact allocated dimensions"

- [ ] **Step 3: Update modalView interface**

In `internal/tui/flows.go`, update the `modalView` interface (around line 33-53):

```go
type modalView interface {
	Update(tea.Msg) tea.Cmd
	View() string
	Shortcuts() []Shortcut
	// SetAvailableSize(maxWidth, maxHeight) stores maximum available dimensions.
	// Modal MAY use less space (e.g., center and pad).
	SetAvailableSize(maxWidth, maxHeight int)
}
```

Changes:
- Rename method: `SetSize(w, h int)` → `SetAvailableSize(maxWidth, maxHeight int)`
- Update comment to clarify "maximum available" semantics

- [ ] **Step 4: Verify syntax is correct**

Run: `cd C:\git\Abditum-T2 && go build ./internal/tui`

Expected: Build should fail with compile errors about method not found (this is expected—we'll fix implementations next)

- [ ] **Step 5: Commit interface changes**

```bash
cd C:\git\Abditum-T2
git add internal/tui/flows.go
git commit -m "refactor: update SetSize interface definitions and parameter names

- childModel.SetSize(width, height): exact space semantics
- modalView.SetAvailableSize(maxWidth, maxHeight): maximum space semantics
- No logic changes, interface only"
```

---

## Task 2: Update welcome.go Implementation

**Files:**
- Modify: `internal/tui/welcome.go:66-69`

- [ ] **Step 1: Read current SetSize implementation**

Run: `cd C:\git\Abditum-T2 && type internal\tui\welcome.go | findstr "SetSize" -A 4`

Expected: See current `SetSize(w, h int)` method

- [ ] **Step 2: Update method signature and parameter names**

In `internal/tui/welcome.go` (around line 67), change:

```go
// Before
func (m *welcomeModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// After
func (m *welcomeModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
```

Only change parameter names: `w` → `width`, `h` → `height`

- [ ] **Step 3: Verify build**

Run: `cd C:\git\Abditum-T2 && go build ./internal/tui`

Expected: Build should progress further (welcome.go no longer has compile error)

- [ ] **Step 4: Commit**

```bash
cd C:\git\Abditum-T2
git add internal/tui/welcome.go
git commit -m "refactor: update welcomeModel.SetSize parameter names to width, height"
```

---

## Task 3: Update vaulttree.go Implementation

**Files:**
- Modify: `internal/tui/vaulttree.go:49-52`

- [ ] **Step 1: Update method signature**

In `internal/tui/vaulttree.go` (around line 50), change:

```go
// Before
func (m *vaultTreeModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// After
func (m *vaultTreeModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
```

- [ ] **Step 2: Verify build**

Run: `cd C:\git\Abditum-T2 && go build ./internal/tui`

Expected: Build progresses further

- [ ] **Step 3: Commit**

```bash
cd C:\git\Abditum-T2
git add internal/tui/vaulttree.go
git commit -m "refactor: update vaultTreeModel.SetSize parameter names to width, height"
```

---

## Task 4: Update secretdetail.go Implementation

**Files:**
- Modify: `internal/tui/secretdetail.go:49-52`

- [ ] **Step 1: Update method signature**

In `internal/tui/secretdetail.go` (around line 50), change:

```go
// Before
func (m *secretDetailModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// After
func (m *secretDetailModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
```

- [ ] **Step 2: Verify build**

Run: `cd C:\git\Abditum-T2 && go build ./internal/tui`

Expected: Build progresses further

- [ ] **Step 3: Commit**

```bash
cd C:\git\Abditum-T2
git add internal/tui/secretdetail.go
git commit -m "refactor: update secretDetailModel.SetSize parameter names to width, height"
```

---

## Task 5: Update templatelist.go Implementation

**Files:**
- Modify: `internal/tui/templatelist.go:49-52`

- [ ] **Step 1: Update method signature**

In `internal/tui/templatelist.go` (around line 50), change:

```go
// Before
func (m *templateListModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// After
func (m *templateListModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
```

- [ ] **Step 2: Verify build**

Run: `cd C:\git\Abditum-T2 && go build ./internal/tui`

Expected: Build progresses further

- [ ] **Step 3: Commit**

```bash
cd C:\git\Abditum-T2
git add internal/tui/templatelist.go
git commit -m "refactor: update templateListModel.SetSize parameter names to width, height"
```

---

## Task 6: Update templatedetail.go Implementation

**Files:**
- Modify: `internal/tui/templatedetail.go:49-52`

- [ ] **Step 1: Update method signature**

In `internal/tui/templatedetail.go` (around line 50), change:

```go
// Before
func (m *templateDetailModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// After
func (m *templateDetailModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
```

- [ ] **Step 2: Verify build**

Run: `cd C:\git\Abditum-T2 && go build ./internal/tui`

Expected: Build progresses further

- [ ] **Step 3: Commit**

```bash
cd C:\git\Abditum-T2
git add internal/tui/templatedetail.go
git commit -m "refactor: update templateDetailModel.SetSize parameter names to width, height"
```

---

## Task 7: Update settings.go Implementation

**Files:**
- Modify: `internal/tui/settings.go:48-51`

- [ ] **Step 1: Update method signature**

In `internal/tui/settings.go` (around line 49), change:

```go
// Before
func (m *settingsModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// After
func (m *settingsModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
```

- [ ] **Step 2: Verify build**

Run: `cd C:\git\Abditum-T2 && go build ./internal/tui`

Expected: Build progresses further

- [ ] **Step 3: Commit**

```bash
cd C:\git\Abditum-T2
git add internal/tui/settings.go
git commit -m "refactor: update settingsModel.SetSize parameter names to width, height"
```

---

## Task 8: Update modalModel Implementation

**Files:**
- Modify: `internal/tui/modal.go:101-105`

- [ ] **Step 1: Update method signature and name**

In `internal/tui/modal.go` (around line 105), change:

```go
// Before
func (m *modalModel) SetSize(w, h int) {}

// After
func (m *modalModel) SetAvailableSize(maxWidth, maxHeight int) {}
```

Changes:
- Rename: `SetSize` → `SetAvailableSize`
- Parameter names: `w, h` → `maxWidth, maxHeight`
- Body remains empty (no-op)

- [ ] **Step 2: Update comment**

Update comment above the method (line 101-104) to:

```go
// SetAvailableSize stores terminal dimensions for layout calculations.
// For modalModel, this is called by rootModel before View() per the SetAvailableSize-before-View contract,
// but modalModel uses fixed width and doesn't need the dimensions.
// Other modalView implementations may use these dimensions to constrain their layout.
```

- [ ] **Step 3: Verify build**

Run: `cd C:\git\Abditum-T2 && go build ./internal/tui`

Expected: Build progresses, but still has errors about missing implementations

- [ ] **Step 4: Commit**

```bash
cd C:\git\Abditum-T2
git add internal/tui/modal.go
git commit -m "refactor: rename modalModel.SetSize to SetAvailableSize with parameter names"
```

---

## Task 9: Update help.go Implementation

**Files:**
- Modify: `internal/tui/help.go:119-123`

- [ ] **Step 1: Update method signature and name**

In `internal/tui/help.go` (around line 120), change:

```go
// Before
func (m *helpModal) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// After
func (m *helpModal) SetAvailableSize(maxWidth, maxHeight int) {
	m.width = maxWidth
	m.height = maxHeight
}
```

Changes:
- Rename: `SetSize` → `SetAvailableSize`
- Parameter names: `w, h` → `maxWidth, maxHeight`
- Assignment: `m.width = maxWidth`, `m.height = maxHeight`

- [ ] **Step 2: Update comment**

Update comment above the method (line 119) to:

```go
// SetAvailableSize sets the maximum available dimensions for dynamic modal sizing.
```

- [ ] **Step 3: Verify build**

Run: `cd C:\git\Abditum-T2 && go build ./internal/tui`

Expected: Build progresses

- [ ] **Step 4: Commit**

```bash
cd C:\git\Abditum-T2
git add internal/tui/help.go
git commit -m "refactor: rename helpModal.SetSize to SetAvailableSize with parameter names"
```

---

## Task 10: Update decision.go Implementation

**Files:**
- Modify: `internal/tui/decision.go:120-124`

- [ ] **Step 1: Update method signature and name**

In `internal/tui/decision.go` (around line 121), change:

```go
// Before
func (d *DecisionDialog) SetSize(w, h int) {
	d.width = w
	d.height = h
}

// After
func (d *DecisionDialog) SetAvailableSize(maxWidth, maxHeight int) {
	d.width = maxWidth
	d.height = maxHeight
}
```

Changes:
- Rename: `SetSize` → `SetAvailableSize`
- Parameter names: `w, h` → `maxWidth, maxHeight`
- Assignment: `d.width = maxWidth`, `d.height = maxHeight`

- [ ] **Step 2: Update comment**

Update comment above the method (line 120) to:

```go
// SetAvailableSize stores maximum available dimensions for future rendering.
```

- [ ] **Step 3: Verify build**

Run: `cd C:\git\Abditum-T2 && go build ./internal/tui`

Expected: Build should now have only root.go errors (implementation methods are fixed)

- [ ] **Step 4: Commit**

```bash
cd C:\git\Abditum-T2
git add internal/tui/decision.go
git commit -m "refactor: rename DecisionDialog.SetSize to SetAvailableSize with parameter names"
```

---

## Task 11: Update filepicker.go Implementation

**Files:**
- Modify: `internal/tui/filepicker.go:82-90`

- [ ] **Step 1: Update method signature and name**

In `internal/tui/filepicker.go` (around line 85), change:

```go
// Before
func (m *filePickerModal) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.treeScroll = 0
	m.adjustTreeScroll()
}

// After
func (m *filePickerModal) SetAvailableSize(maxWidth, maxHeight int) {
	m.width = maxWidth
	m.height = maxHeight
	m.treeScroll = 0
	m.adjustTreeScroll()
}
```

Changes:
- Rename: `SetSize` → `SetAvailableSize`
- Parameter names: `w, h` → `maxWidth, maxHeight`
- Assignment: `m.width = maxWidth`, `m.height = maxHeight`
- Behavior unchanged: still calls `adjustTreeScroll()` after storing dimensions

- [ ] **Step 2: Update comment**

Update comment above the method (line 82-84) to:

```go
// SetAvailableSize sets the maximum available dimensions for the file picker modal.
// Also resets scroll state to ensure valid viewport positioning.
```

- [ ] **Step 3: Verify build**

Run: `cd C:\git\Abditum-T2 && go build ./internal/tui`

Expected: Build progresses

- [ ] **Step 4: Commit**

```bash
cd C:\git\Abditum-T2
git add internal/tui/filepicker.go
git commit -m "refactor: rename filePickerModal.SetSize to SetAvailableSize with parameter names"
```

---

## Task 12: Update passwordcreate.go Implementation

**Files:**
- Modify: `internal/tui/passwordcreate.go:316-319`

- [ ] **Step 1: Update method signature and name**

In `internal/tui/passwordcreate.go` (around line 318), change:

```go
// Before
func (m *passwordCreateModal) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// After
func (m *passwordCreateModal) SetAvailableSize(maxWidth, maxHeight int) {
	m.width = maxWidth
	m.height = maxHeight
}
```

Changes:
- Rename: `SetSize` → `SetAvailableSize`
- Parameter names: `w, h` → `maxWidth, maxHeight`
- Assignment: `m.width = maxWidth`, `m.height = maxHeight`

- [ ] **Step 2: Update comment**

Update comment above the method (line 316-317) to:

```go
// SetAvailableSize sets the maximum available dimensions for modal layout.
```

- [ ] **Step 3: Verify build**

Run: `cd C:\git\Abditum-T2 && go build ./internal/tui`

Expected: Build progresses

- [ ] **Step 4: Commit**

```bash
cd C:\git\Abditum-T2
git add internal/tui/passwordcreate.go
git commit -m "refactor: rename passwordCreateModal.SetSize to SetAvailableSize with parameter names"
```

---

## Task 13: Update passwordentry.go Implementation

**Files:**
- Modify: `internal/tui/passwordentry.go:218-221`

- [ ] **Step 1: Update method signature and name**

In `internal/tui/passwordentry.go` (around line 220), change:

```go
// Before
func (m *passwordEntryModal) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// After
func (m *passwordEntryModal) SetAvailableSize(maxWidth, maxHeight int) {
	m.width = maxWidth
	m.height = maxHeight
}
```

Changes:
- Rename: `SetSize` → `SetAvailableSize`
- Parameter names: `w, h` → `maxWidth, maxHeight`
- Assignment: `m.width = maxWidth`, `m.height = maxHeight`

- [ ] **Step 2: Update comment**

Update comment above the method (line 218-219) to:

```go
// SetAvailableSize sets the maximum available dimensions for modal layout.
```

- [ ] **Step 3: Verify build**

Run: `cd C:\git\Abditum-T2 && go build ./internal/tui`

Expected: Build progresses

- [ ] **Step 4: Commit**

```bash
cd C:\git\Abditum-T2
git add internal/tui/passwordentry.go
git commit -m "refactor: rename passwordEntryModal.SetSize to SetAvailableSize with parameter names"
```

---

## Task 14: Update rootModel Call Sites (Part 1: WindowSizeMsg)

**Files:**
- Modify: `internal/tui/root.go:192-198`

- [ ] **Step 1: Update WindowSizeMsg handling**

In `internal/tui/root.go` (around line 192-198), the code already uses full names (`msg.Width`, `msg.Height`), so we update the loop call:

Current code:
```go
case tea.WindowSizeMsg:
	m.width = msg.Width
	m.height = msg.Height
	for _, child := range m.liveWorkChildren() {
		child.SetSize(msg.Width, msg.Height)
	}
	return m, nil
```

This part is already using full parameter semantics (passing exact space to children). No code change needed here, but verify it's correct:

Run: `cd C:\git\Abditum-T2 && go build ./internal/tui`

Expected: This part should compile now (no changes to make)

- [ ] **Step 2: Commit (no changes)**

```bash
cd C:\git\Abditum-T2
git commit --allow-empty -m "refactor: windowSizeMsg handling already uses correct semantics"
```

---

## Task 15: Update rootModel Call Sites (Part 2: Modal rendering)

**Files:**
- Modify: `internal/tui/root.go:419-430`

- [ ] **Step 1: Update View method modal call**

In `internal/tui/root.go` (around line 419-430), change:

```go
// Before
if len(m.modals) > 0 {
	const headerH = 2
	const msgBarH = 1
	const cmdBarH = 1
	workH := m.height - headerH - msgBarH - cmdBarH
	if workH < 0 {
		workH = 0
	}
	top := m.modals[len(m.modals)-1]
	top.SetSize(m.width, workH) // pass workH, not total height
	content = m.renderFrame(top)
}

// After
if len(m.modals) > 0 {
	const headerH = 2
	const msgBarH = 1
	const cmdBarH = 1
	workH := m.height - headerH - msgBarH - cmdBarH
	if workH < 0 {
		workH = 0
	}
	top := m.modals[len(m.modals)-1]
	top.SetAvailableSize(m.width, workH) // pass workH as maxHeight
	content = m.renderFrame(top)
}
```

Only change: `SetSize` → `SetAvailableSize`

- [ ] **Step 2: Verify build**

Run: `cd C:\git\Abditum-T2 && go build ./internal/tui`

Expected: Build progresses further

- [ ] **Step 3: Commit**

```bash
cd C:\git\Abditum-T2
git add internal/tui/root.go
git commit -m "refactor: update modal SetSize calls to SetAvailableSize in View()"
```

---

## Task 16: Update rootModel Call Sites (Part 3: renderFrame welcome/settings)

**Files:**
- Modify: `internal/tui/root.go:480-495`

- [ ] **Step 1: Check current code**

In `internal/tui/root.go` (around line 480-495), the code should look like:

```go
case workAreaWelcome:
	if m.welcome != nil {
		m.welcome.SetSize(m.width, workH)
		workContent = m.welcome.View()
	}
case workAreaSettings:
	if m.settings != nil {
		m.settings.SetSize(m.width, workH)
		workContent = m.settings.View()
	}
```

These are already correct (SetSize for childModel, exact space). No changes needed.

- [ ] **Step 2: Verify build**

Run: `cd C:\git\Abditum-T2 && go build ./internal/tui`

Expected: Build proceeds

- [ ] **Step 3: Commit (documentation)**

```bash
cd C:\git\Abditum-T2
git commit --allow-empty -m "refactor: welcome and settings SetSize calls already use correct semantics"
```

---

## Task 17: Update rootModel Call Sites (Part 4: renderVaultArea)

**Files:**
- Modify: `internal/tui/root.go:520-541`

- [ ] **Step 1: Check current code**

In `internal/tui/root.go` (around line 520-541), the code should look like:

```go
func (m *rootModel) renderVaultArea(workH int) string {
	halfW := m.width / 2
	if m.vaultTree != nil {
		m.vaultTree.SetSize(halfW, workH)
	}
	if m.secretDetail != nil {
		m.secretDetail.SetSize(m.width-halfW, workH)
	}
	// ... rest of rendering
}
```

This is already correct (SetSize for childModel, exact space). No changes needed.

- [ ] **Step 2: Verify build**

Run: `cd C:\git\Abditum-T2 && go build ./internal/tui`

Expected: Build proceeds

- [ ] **Step 3: Commit (documentation)**

```bash
cd C:\git\Abditum-T2
git commit --allow-empty -m "refactor: renderVaultArea SetSize calls already use correct semantics"
```

---

## Task 18: Update rootModel Call Sites (Part 5: renderTemplatesArea)

**Files:**
- Modify: `internal/tui/root.go:544-565`

- [ ] **Step 1: Check current code**

In `internal/tui/root.go` (around line 544-565), the code should look like:

```go
func (m *rootModel) renderTemplatesArea(workH int) string {
	halfW := m.width / 2
	if m.templateList != nil {
		m.templateList.SetSize(halfW, workH)
	}
	if m.templateDetail != nil {
		m.templateDetail.SetSize(m.width-halfW, workH)
	}
	// ... rest of rendering
}
```

This is already correct (SetSize for childModel, exact space). No changes needed.

- [ ] **Step 2: Verify build**

Run: `cd C:\git\Abditum-T2 && go build ./internal/tui`

Expected: Build proceeds

- [ ] **Step 3: Commit (documentation)**

```bash
cd C:\git\Abditum-T2
git commit --allow-empty -m "refactor: renderTemplatesArea SetSize calls already use correct semantics"
```

---

## Task 19: Update rootModel Call Sites (Part 6: enterVault)

**Files:**
- Modify: `internal/tui/root.go:575-576`

- [ ] **Step 1: Check current code**

In `internal/tui/root.go` (around line 575-576), the code should look like:

```go
m.vaultTree.SetSize(m.width/2, m.height-4)
m.secretDetail.SetSize(m.width-m.width/2, m.height-4)
```

This is already correct (SetSize for childModel, exact space). No changes needed.

- [ ] **Step 2: Verify build**

Run: `cd C:\git\Abditum-T2 && go build ./internal/tui`

Expected: Build should now succeed! All code should compile.

- [ ] **Step 3: Commit (documentation)**

```bash
cd C:\git\Abditum-T2
git commit --allow-empty -m "refactor: enterVault SetSize calls already use correct semantics"
```

---

## Task 20: Update help_test.go

**Files:**
- Modify: `internal/tui/help_test.go` (multiple lines)

- [ ] **Step 1: Find all SetSize calls in tests**

Run: `cd C:\git\Abditum-T2 && findstr /N "SetSize" internal\tui\help_test.go`

Expected: Multiple matches (lines ~110, 130, 138, 177, 188, 199, 210, 221, 235, 250, 260)

- [ ] **Step 2: Replace SetSize with SetAvailableSize**

In `internal/tui/help_test.go`, replace all:
- `m.SetSize(` → `m.SetAvailableSize(`

All lines that call SetSize should now call SetAvailableSize. Example:

```go
// Before
m.SetSize(w, h)

// After
m.SetAvailableSize(w, h)
```

- [ ] **Step 3: Verify tests compile**

Run: `cd C:\git\Abditum-T2 && go test -v ./internal/tui -run TestHelpModal 2>&1 | head -30`

Expected: Tests should compile and run (should pass or show expected failures if any)

- [ ] **Step 4: Run specific help tests**

Run: `cd C:\git\Abditum-T2 && go test -v ./internal/tui -run TestHelpModal_UpdateScroll`

Expected: PASS

- [ ] **Step 5: Run all help tests**

Run: `cd C:\git\Abditum-T2 && go test -v ./internal/tui -run TestHelpModal`

Expected: All tests PASS

- [ ] **Step 6: Commit**

```bash
cd C:\git\Abditum-T2
git add internal/tui/help_test.go
git commit -m "refactor: update help_test.go SetSize calls to SetAvailableSize"
```

---

## Task 21: Update decision_test.go

**Files:**
- Modify: `internal/tui/decision_test.go` (multiple lines)

- [ ] **Step 1: Find all SetSize calls in tests**

Run: `cd C:\git\Abditum-T2 && findstr /N "SetSize" internal\tui\decision_test.go`

Expected: Multiple matches (lines ~28, 40, 53, 66, 78, 91, and comments)

- [ ] **Step 2: Replace SetSize with SetAvailableSize**

In `internal/tui/decision_test.go`, replace all:
- `d.SetSize(` → `d.SetAvailableSize(`
- Also update comment on line 17: `SetSize(80, 24)` → `SetAvailableSize(80, 24)`

Example:

```go
// Before
d.SetSize(80, 24)

// After
d.SetAvailableSize(80, 24)
```

- [ ] **Step 3: Verify tests compile**

Run: `cd C:\git\Abditum-T2 && go test -v ./internal/tui -run TestDecision 2>&1 | head -30`

Expected: Tests compile and run

- [ ] **Step 4: Run decision tests**

Run: `cd C:\git\Abditum-T2 && go test -v ./internal/tui -run TestDecision`

Expected: All tests PASS

- [ ] **Step 5: Commit**

```bash
cd C:\git\Abditum-T2
git add internal/tui/decision_test.go
git commit -m "refactor: update decision_test.go SetSize calls to SetAvailableSize"
```

---

## Task 22: Update filepicker_test.go

**Files:**
- Modify: `internal/tui/filepicker_test.go` (multiple lines)

- [ ] **Step 1: Find all SetSize calls in tests**

Run: `cd C:\git\Abditum-T2 && findstr /N "SetSize" internal\tui\filepicker_test.go`

Expected: Multiple matches (exact count depends on test file size)

- [ ] **Step 2: Replace SetSize with SetAvailableSize**

In `internal/tui/filepicker_test.go`, replace all:
- `m.SetSize(` → `m.SetAvailableSize(`

Example:

```go
// Before
m.SetSize(w, h)

// After
m.SetAvailableSize(w, h)
```

- [ ] **Step 3: Verify tests compile**

Run: `cd C:\git\Abditum-T2 && go test -v ./internal/tui -run TestFilePickerModal 2>&1 | head -30`

Expected: Tests compile and run

- [ ] **Step 4: Run filepicker tests**

Run: `cd C:\git\Abditum-T2 && go test -v ./internal/tui -run TestFilePickerModal`

Expected: All tests PASS

- [ ] **Step 5: Commit**

```bash
cd C:\git\Abditum-T2
git add internal/tui/filepicker_test.go
git commit -m "refactor: update filepicker_test.go SetSize calls to SetAvailableSize"
```

---

## Task 23: Update passwordcreate_test.go

**Files:**
- Modify: `internal/tui/passwordcreate_test.go` (multiple lines)

- [ ] **Step 1: Find all SetSize calls in tests**

Run: `cd C:\git\Abditum-T2 && findstr /N "SetSize" internal\tui\passwordcreate_test.go`

Expected: Multiple matches (exact count depends on test file size)

- [ ] **Step 2: Replace SetSize with SetAvailableSize**

In `internal/tui/passwordcreate_test.go`, replace all:
- `m.SetSize(` → `m.SetAvailableSize(`

Example:

```go
// Before
m.SetSize(w, h)

// After
m.SetAvailableSize(w, h)
```

- [ ] **Step 3: Verify tests compile**

Run: `cd C:\git\Abditum-T2 && go test -v ./internal/tui -run TestPasswordCreateModal 2>&1 | head -30`

Expected: Tests compile and run

- [ ] **Step 4: Run passwordcreate tests**

Run: `cd C:\git\Abditum-T2 && go test -v ./internal/tui -run TestPasswordCreateModal`

Expected: All tests PASS

- [ ] **Step 5: Commit**

```bash
cd C:\git\Abditum-T2
git add internal/tui/passwordcreate_test.go
git commit -m "refactor: update passwordcreate_test.go SetSize calls to SetAvailableSize"
```

---

## Task 24: Update passwordentry_test.go

**Files:**
- Modify: `internal/tui/passwordentry_test.go` (multiple lines)

- [ ] **Step 1: Find all SetSize calls in tests**

Run: `cd C:\git\Abditum-T2 && findstr /N "SetSize" internal\tui\passwordentry_test.go`

Expected: Multiple matches (exact count depends on test file size)

- [ ] **Step 2: Replace SetSize with SetAvailableSize**

In `internal/tui/passwordentry_test.go`, replace all:
- `m.SetSize(` → `m.SetAvailableSize(`

Example:

```go
// Before
m.SetSize(w, h)

// After
m.SetAvailableSize(w, h)
```

- [ ] **Step 3: Verify tests compile**

Run: `cd C:\git\Abditum-T2 && go test -v ./internal/tui -run TestPasswordEntryModal 2>&1 | head -30`

Expected: Tests compile and run

- [ ] **Step 4: Run passwordentry tests**

Run: `cd C:\git\Abditum-T2 && go test -v ./internal/tui -run TestPasswordEntryModal`

Expected: All tests PASS

- [ ] **Step 5: Commit**

```bash
cd C:\git\Abditum-T2
git add internal/tui/passwordentry_test.go
git commit -m "refactor: update passwordentry_test.go SetSize calls to SetAvailableSize"
```

---

## Task 25: Full Build and Test Verification

**Files:**
- All modified files (verification only)

- [ ] **Step 1: Clean build**

Run: `cd C:\git\Abditum-T2 && go clean ./internal/tui && go build ./internal/tui`

Expected: Build succeeds with no errors

- [ ] **Step 2: Run all TUI tests**

Run: `cd C:\git\Abditum-T2 && go test -v ./internal/tui`

Expected: All tests PASS

- [ ] **Step 3: Check for any remaining SetSize references (should only be in comments)**

Run: `cd C:\git\Abditum-T2 && findstr /N "\.SetSize(" internal\tui\*.go | findstr -v "^.*// "`

Expected: Empty output (all code references replaced, only comments remain)

- [ ] **Step 4: Verify git log shows semantic changes**

Run: `cd C:\git\Abditum-T2 && git log --oneline | head -20`

Expected: See commits for each interface/implementation change

- [ ] **Step 5: Final commit (summary)**

```bash
cd C:\git\Abditum-T2
git commit --allow-empty -m "chore: SetSize semantics refactor complete (Phase 1)

All SetSize methods renamed:
- childModel.SetSize(width, height): exact space semantics
- modalView.SetAvailableSize(maxWidth, maxHeight): maximum space semantics

Parameter names updated for clarity throughout codebase.
No logic changes—interface and naming refactor only.
All tests passing."
```

---

## Plan Verification

✅ **Spec coverage:**
- Interface definitions updated (Task 1)
- All childModel implementations updated (Tasks 2-7)
- All modalView implementations updated (Tasks 8-13) [+3 new discovered modalViews]
- All rootModel call sites verified (Tasks 14-19)
- All tests updated (Tasks 20-24) [+3 new test files]
- Full verification (Task 25)

✅ **No placeholders:** All steps have concrete code and commands

✅ **Type consistency:** 
- `SetSize(width, height int)` for childModel throughout
- `SetAvailableSize(maxWidth, maxHeight int)` for modalView throughout
- Parameter names consistent across all implementations

✅ **Commit frequency:** One commit per logical unit (~1 per file modified)

---

## Execution Options

Plan complete and saved to `docs/superpowers/plans/2026-04-12-setsize-semantics-refactor-phase1.md`.

Two execution options:

**1. Subagent-Driven (recommended)**
- Fresh subagent per task (or grouped tasks)
- Review between major checkpoints
- Best for catching issues early

**2. Inline Execution**
- Execute tasks sequentially in this session
- Batch execution with checkpoints for review

Which approach would you like to use?
