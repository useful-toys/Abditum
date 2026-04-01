# Phase 05 Validation Architecture

**Phase:** 05-tui-scaffold-root-model
**Generated from:** 05-RESEARCH.md § Validation Architecture
**Status:** Pending execution

---

## Test Coverage Plan

All tests in `internal/tui/` run with `CGO_ENABLED=0 go test -race ./internal/tui/...`.

### Unit Tests — root_test.go (Plan 05)

| Test | Assertion | Plan |
|------|-----------|------|
| `TestRootModelInit` | `rootModel.Init()` returns nil; `area == workAreaPreVault`; `preVault != nil`; `len(modals) == 0` | 05-05 |
| `TestRootModelViewType` | `rootModel.View()` compiles as `tea.View` — enforced by `var _ tea.Model = &rootModel{}` assertion in root.go | 05-04 |
| `TestModalStack_PushPop` | Push two modals → `len == 2`; pop → `len == 1`; pop → `len == 0`; extra pop does not panic | 05-05 |
| `TestLiveModels_TypedNilSafety` | Nil concrete pointer fields do not appear in `liveModels()` as typed-nil interface values | 05-05 |
| `TestDispatchPriority_CtrlQ` | `ctrl+Q` intercepted before reaching any child; returns non-nil Cmd (tea.Quit or confirm dialog) | 05-05 |
| `TestWindowSizeMsg_PropagatesToChildren` | `tea.WindowSizeMsg` updates `m.width`, `m.height`, and propagates `SetSize` to live children | 05-05 |

### Compile-Time Assertions

| Assertion | File | Meaning |
|-----------|------|---------|
| `var _ tea.Model = &rootModel{}` | root.go | rootModel satisfies tea.Model (View returns tea.View) |
| `var _ childModel = &modalModel{}` | modal.go | modalModel satisfies childModel |
| `var _ childModel = &preVaultModel{}` | prevault.go | preVaultModel satisfies childModel |
| `var _ childModel = &vaultTreeModel{}` | vaulttree.go | vaultTreeModel satisfies childModel |
| `var _ childModel = &secretDetailModel{}` | secretdetail.go | secretDetailModel satisfies childModel |
| `var _ childModel = &templateListModel{}` | templatelist.go | templateListModel satisfies childModel |
| `var _ childModel = &templateDetailModel{}` | templatedetail.go | templateDetailModel satisfies childModel |
| `var _ childModel = &settingsModel{}` | settings.go | settingsModel satisfies childModel |
| `var _ childModel = &helpModal{}` | help.go | helpModal satisfies childModel |

### Golden File Tests (Phase 5 scope — deferred to Phase 6)

These tests require a running TUI via `teatest.NewTestModel`. They are listed here as acceptance criteria but their _implementation_ is deferred until Phase 6 when the welcome screen has real content. The test infrastructure (teatest/v2 dependency) is added in Plan 01.

| Test | Assertion |
|------|-----------|
| `TestGolden_PreVaultPlaceholder` | `./abditum` renders placeholder frame at 80×24 terminal; output matches `testdata/golden_prevault.txt` |

All golden tests must use `teatest.WithInitialTermSize(80, 24)`.

---

## Automated Verification Commands

```bash
# Full phase verification — run after all 5 plans complete:
CGO_ENABLED=0 go build ./...
CGO_ENABLED=0 go vet ./...
CGO_ENABLED=0 go test -race ./internal/tui/... -v
CGO_ENABLED=0 go test -race ./... -count=1
```

Expected outcome: 0 build errors, 0 vet warnings, 0 test failures, 0 race conditions.

---

## Critical Pitfall Checks

These are verified by the tests above but noted explicitly for executor awareness:

1. **`View()` return type:** Children return `string`; only `rootModel.View()` returns `tea.View`. Enforced by `var _ tea.Model = &rootModel{}` compile-time assertion.

2. **Typed nil in interface trap:** `liveModels()` must use explicit `!= nil` checks on concrete pointer fields. `TestLiveModels_TypedNilSafety` catches any regression.

3. **Bubble Tea v2 key strings:** Match via `msg.String()` — returns `"ctrl+q"`, `"space"`, `"enter"` (not `tea.KeyCtrlQ` constants). Verified by `TestDispatchPriority_CtrlQ`.

4. **Tick renewal:** Every `tickMsg` handler must return a new `tea.Tick` cmd. Missing renewal silently stops all ticking. (No test for this in Phase 5 — tick doesn't start until workAreaVault; Phase 6 will test.)

5. **`tea.WithAltScreen()` does not exist in v2:** AltScreen is set via `v.AltScreen = true` in `rootModel.View()`. Build failure if this option is used.

6. **CGO_ENABLED=0:** All build and test commands must include this env var. CGO compilation is forbidden per `arquitetura.md §5`.
