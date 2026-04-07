# Keybindings Specification Analysis

## Summary

**Scope:** Keybindings and action registration compliance with `@tui-specification-novo.md`

**Status:** 5 actions currently registered; 13 F-key actions defined in spec

**Gaps Identified:**
1. Welcome screen uses `O` and `N` instead of spec-defined `F5` and `F6`
2. F2, F3, F4 (workspace tabs) not registered — depends on post-Phase-06 vault area implementation
3. F7-F9 (save/export/import flows) not registered — depends on Phase 7+ flows
4. F12 (Toggle Theme) is POC feature, **not in spec** — requires decision
5. Help modal gaps detection missing from current registration flow

---

## I. Keybindings: Spec vs. Implementation

### Global Actions (Workspace-Independent)

| Spec `@tui-spec:46-79` | Current Code | Status | Notes |
|---|---|---|---|
| **F1** | `root.go:135-140` | ✅ SPEC-COMPLIANT | Help modal (Group 1, Priority 0, HideFromBar false) |
| **F12** | `root.go:131-134` | ⚠ **NOT IN SPEC** | Theme toggle (Group 0, Priority 100, HideFromBar true) — POC feature |
| **Ctrl+Q** | `root.go:111-114` + root.go:243-257 | ⚠ **PARTIAL** | Base action registered; special logic in Update() pre-dispatch |
| **Ctrl+Alt+Shift+Q** | *Not implemented* | ❌ NOT IMPLEMENTED | Lock vault (Fluxo 6) — Phase 7+ |

### Workspace Actions (Welcome Screen / Open Vault)

| Spec `@tui-spec:59-79` | Current Code | Status | Scope | Notes |
|---|---|---|---|---|
| **F2** (Cofre tab) | *Not registered* | ❌ NOT IMPLEMENTED | F2-F4: workspace tabs | Depends on vaultTree/templateList/settings models |
| **F3** (Modelos tab) | *Not registered* | ❌ NOT IMPLEMENTED | F2-F4: workspace tabs | Depends on vaultTree/templateList/settings models |
| **F4** (Config tab) | *Not registered* | ❌ NOT IMPLEMENTED | F2-F4: workspace tabs | Depends on vaultTree/templateList/settings models + Fluxo 14 |
| **F5** (Criar Novo) | `root.go:123-129` as **`N`** | ⚠ **WRONG KEY** | Fluxo 2 | Currently bound to `N` not `F5` |
| **F6** (Abrir Cofre) | `root.go:116-122` as **`O`** | ⚠ **WRONG KEY** | Fluxo 1 | Currently bound to `O` not `F6` |
| **Shift+F6** (Reload) | *Not registered* | ❌ NOT IMPLEMENTED | Fluxo 10 | Reload/discard changes — Phase 7+ |
| **F7** (Save current) | *Not registered* | ❌ NOT IMPLEMENTED | Fluxo 8 | Save to current file — Phase 7+ |
| **Shift+F7** (Save as) | *Not registered* | ❌ NOT IMPLEMENTED | Fluxo 9 | Save to new file — Phase 7+ |
| **Ctrl+F7** (Change pwd) | *Not registered* | ❌ NOT IMPLEMENTED | Fluxo 11 | Master password change — Phase 7+ |
| **F8** | — | — | Reserved | (No action assigned) |
| **F9** (Export) | *Not registered* | ❌ NOT IMPLEMENTED | Fluxo 12 | Export vault — Phase 7+ |
| **Shift+F9** (Import) | *Not registered* | ❌ NOT IMPLEMENTED | Fluxo 13 | Import vault — Phase 7+ |
| **F10** | — | — | Reserved | (No action assigned) |
| **F11** | — | — | Reserved | (No action assigned) |

---

## II. Keybindings: Implementation Architecture

### Current Registration Pattern

`root.go:109-142` registers 5 actions in `newRootModel()`:

```go
// Pattern: Action{Keys: []string{"key"}, Label, Description, Group, Scope, Priority, HideFromBar, Enabled, Handler}
Action{Keys: []string{"ctrl+q"}, ...}      // (1)
Action{Keys: []string{"o"}, ...}            // (2)
Action{Keys: []string{"n"}, ...}            // (3)
Action{Keys: []string{"f12"}, ...}          // (4)
Action{Keys: []string{"f1"}, ...}           // (5)
```

### Dispatch Flow

`root.go:237-292` keyboard input handling order (D-09):

1. **Pre-dispatch special cases** (lines 243-268)
   - `ctrl+q` → special logic for unsaved changes + Decision dialog
   - `f12` → theme toggle (hardcoded, bypasses ActionManager)

2. **Help modal intercept** (lines 270-275)
   - If help modal is open, it consumes all keys (including F1 and ESC)

3. **ActionManager.Dispatch()** (line 277)
   - Checks Scope (Global vs. Local)
   - Checks Enabled() function
   - Matches key against all registered actions
   - Executes Handler if match found

4. **Other modals** (lines 281-282)
   - Top-level modal gets remaining key events

5. **Active flow** (lines 285-286)
   - If flow is running and no modal consumed key, flow gets it

6. **Active work-area child** (lines 289-290)
   - Current workspace model (welcome, vaultTree, templateList, settings) gets key

### Key Registration Strategy

Each action's `Keys` field is an array:

```go
Keys: []string{"key1", "key2", ...}
```

- **Display:** `RenderCommandBar()` shows `Keys[0]` in command bar
- **Dispatch:** All elements in `Keys` trigger the action
- **Alias support:** Multiple keys can trigger same action (not currently used)

### Priority and Bar Rendering

`actions.go:175-272` command bar rendering:

- **F1 is right-anchored** (lines 193-203)
- **Other actions are left-padded** with 2 spaces (line 228)
- **Truncation order:** Lowest-priority actions removed first when width is constrained (lines 252-260)
- **Priority levels used currently:**
  - Ctrl+Q: Priority 10 (pre-dispatch, not in bar)
  - N (Novo): Priority 94
  - O (Abrir): Priority 95
  - F12 (Toggle): Priority 100 (hidden from bar)
  - F1 (Help): Priority 0 (right-anchored)

---

## III. Current State Details

### File: `internal/tui/root.go`

**Lines 109-142: Action registration in `newRootModel()`**

```go
actions.Register(m,
    // Ctrl+Q — Special: pre-dispatched in Update(), line 243-257
    Action{Keys: []string{"ctrl+q"}, Label: "Sair", Description: "Sair do Abditum",
        Group: 0, Scope: ScopeLocal, Priority: 10, HideFromBar: false,
        Enabled: func() bool { return true },
        Handler: func() tea.Cmd { return tea.Quit }},
    
    // O key — Welcome screen, should be F6 per spec
    Action{Keys: []string{"o"}, Label: "Abrir", Description: "Abrir cofre existente",
        Group: 4, Scope: ScopeLocal, Priority: 95, HideFromBar: false,
        Enabled: func() bool { return m.area == workAreaWelcome },
        Handler: func() tea.Cmd {
            flow := newOpenVaultFlow(m.mgr, m.messages, actions, m.theme)
            return func() tea.Msg { return startFlowMsg{flow: flow} }
        }},
    
    // N key — Welcome screen, should be F5 per spec
    Action{Keys: []string{"n"}, Label: "Novo", Description: "Criar novo cofre",
        Group: 4, Scope: ScopeLocal, Priority: 94, HideFromBar: false,
        Enabled: func() bool { return m.area == workAreaWelcome },
        Handler: func() tea.Cmd {
            flow := newCreateVaultFlow(m.mgr, m.messages, actions, m.theme)
            return func() tea.Msg { return startFlowMsg{flow: flow} }
        }},
    
    // F12 — Theme toggle (NOT IN SPEC, POC feature)
    Action{Keys: []string{"f12"}, Label: "Toggle Theme", Description: "Alternar tema visual (Tokyo Night / Cyberpunk)",
        Group: 0, Scope: ScopeGlobal, Priority: 100, HideFromBar: true,
        Enabled: func() bool { return true },
        Handler: func() tea.Cmd { return func() tea.Msg { return toggleThemeMsg{} } }},
    
    // F1 — Help modal
    Action{Keys: []string{"f1"}, Label: "Ajuda", Description: "Mostrar atalhos de teclado",
        Group: 1, Scope: ScopeGlobal, Priority: 0, HideFromBar: false,
        Enabled: func() bool { return true },
        Handler: func() tea.Cmd {
            return func() tea.Msg { return pushModalMsg{modal: newHelpModal(actions.All(), actions.GroupLabel)} }
        }},
)
```

**Lines 243-268: Pre-dispatch special handling for Ctrl+Q and F12**

```go
case tea.KeyPressMsg:
    // ...
    key := msg.String()
    
    // Ctrl+Q special logic — check unsaved changes before dispatching
    if key == "ctrl+q" {
        if m.mgr != nil && m.mgr.IsModified() {
            return m, func() tea.Msg {
                return Decision(SeverityNeutral, "Alterações não salvas",
                    "Deseja salvar as alterações antes de sair?",
                    DecisionAction{Key: "Enter", Label: "Salvar", Default: true},
                    []DecisionAction{{Key: "D", Label: "Descartar"}},
                    DecisionAction{Key: "Esc", Label: "Voltar"})
            }
        }
        return m, tea.Quit
    }
    
    // F12 theme toggle — hardcoded dispatch (bypasses ActionManager)
    if key == "f12" {
        if m.theme == ThemeTokyoNight {
            m.theme = ThemeCyberpunk
        } else {
            m.theme = ThemeTokyoNight
        }
        m.applyTheme()
        return m, nil
    }
```

### File: `internal/tui/actions.go`

**Lines 92-120: `Dispatch()` implementation**

The dispatch flow follows scope and enabled-function checks:

```go
func (a *ActionManager) Dispatch(key string, inFlowOrModal bool) tea.Cmd {
    // Priority order: active owner first, then others
    var ordered []any
    if a.activeOwner != nil {
        ordered = append(ordered, a.activeOwner)
    }
    for _, o := range a.owners {
        if o != a.activeOwner {
            ordered = append(ordered, o)
        }
    }

    for _, owner := range ordered {
        for _, act := range a.byOwner[owner] {
            // Skip ScopeLocal actions if in flow/modal
            if act.Scope == ScopeLocal && inFlowOrModal {
                continue
            }
            // Check enabled function
            if act.Enabled != nil && !act.Enabled() {
                continue
            }
            // Match key against all registered keys for this action
            for _, k := range act.Keys {
                if k == key {
                    return act.Handler()
                }
            }
        }
    }
    return nil
}
```

**Lines 175-272: `RenderCommandBar()` implementation**

- F1 is separated as "anchor" and right-aligned (lines 193-203)
- Body actions are left-padded with 2 spaces (line 228)
- Separators are ` · ` (line 219)
- Truncation removes lowest-priority body actions first (lines 252-260)

---

## IV. Gaps Identified

### Gap K.1: O and N keys should be F5 and F6

**Specification Location:** `tui-specification-novo.md` lines 68-69

```
| `F5` | Criar Novo Cofre (Fluxo 2) | |
| `F6` | Abrir Cofre Existente (Fluxo 1) | |
```

**Current Implementation:** `root.go:116-129`

- O key → Abrir (should be F6)
- N key → Novo (should be F5)

**Issue:** Welcome screen action bar shows `O Abrir · N Novo` instead of spec-compliant `F5 Novo · F6 Abrir`

**Impact:** Users unfamiliar with Abditum won't find expected F-key shortcuts; help documentation must account for alias.

**Remediation:**
- Option A: Change keys to F5/F6, update help modal
- Option B: Support both O/N and F5/F6 as aliases (backward compatibility)
- Scope: ~2 lines per option (just Keys arrays)

---

### Gap K.2: F12 (Toggle Theme) is not in specification

**Specification Location:** `tui-specification-novo.md` lines 46-79 (Atalhos da Aplicação section)

**Current Implementation:** `root.go:131-134` (hardcoded dispatch) + `root.go:260-268` (pre-dispatch)

- Registered as Group 0, Scope ScopeGlobal, Priority 100, HideFromBar true
- Bypasses ActionManager dispatch (handled directly in root.go Update())

**Issue:** F12 is a POC (proof-of-concept) feature not aligned with spec; appears in Help modal but not mentioned in spec.

**Impact:** 
- Theme selection is useful for accessibility/preference, but not part of current spec scope
- If kept, should be documented in spec update
- If removed, delete ~40 lines of theme-toggle code

**Decision Required:**
- **Option A (KEEP):** Add F12 to spec's "Global" section; formalize it as feature
- **Option B (REMOVE):** Delete theme toggle; users use config file instead
- **Option C (DEFER):** Keep code, mark as "POC experimental", plan for Phase 8+ refinement

**Current recommendation (pending decision):** Option A — keep it, add to spec's global actions

**Remediation:** If removing: delete lines 131-134, 260-268, toggleThemeMsg type, m.applyTheme() method; adjust RenderCommandBar tests

---

### Gap K.3: F2, F3, F4 (Workspace Tabs) Not Registered

**Specification Location:** `tui-specification-novo.md` lines 65-67

```
| `F2` | Modo Cofre (aba) | Só com cofre aberto |
| `F3` | Modo Modelos (aba) | Só com cofre aberto |
| `F4` | Modo Configurações (aba) | Abrange o Fluxo 14: Configurar o Cofre |
```

**Current Implementation:** *Not implemented*

**Why:** These depend on vault-area child models:
- F2 → switch to vaultTree area (Phase 06.5+)
- F3 → switch to templateList area (Phase 06.6+)
- F4 → switch to settings area + Fluxo 14 (Phase 06.7+)

**Scope:** Phase 06.1 → Document as "deferred to Phase 06.5+" in PLAN.md

**Remediation (Phase 06.5+):** 
```go
Action{Keys: []string{"f2"}, Label: "Cofre", Description: "Modo Cofre",
    Group: 2, Scope: ScopeLocal, Priority: 90, 
    Enabled: func() bool { return m.area == workAreaVault && m.vaultTree != nil },
    Handler: func() tea.Cmd { m.area = workAreaVault; return nil }},
// ... similar for F3, F4
```

---

### Gap K.4: F7-F9 (Save/Export/Import) Not Registered

**Specification Location:** `tui-specification-novo.md` lines 71-76

```
| `F7` | Salvar Cofre no Arquivo Atual (Fluxo 8) | |
| `Shift+F7` | Salvar Cofre em Outro Arquivo (Fluxo 9) | |
| `Ctrl+F7` | Alterar Senha Mestra (Fluxo 11) | |
| `F9` | Exportar Cofre (Fluxo 12) | |
| `Shift+F9` | Importar Cofre (Fluxo 13) | |
```

**Current Implementation:** *Not implemented*

**Why:** These flows don't exist yet (Fluxos 8-13 are Phase 7+ work)

**Scope:** Phase 06.1 → Document as "deferred to Phase 7+" in PLAN.md

**Remediation (Phase 7+):** Each flow gets keybinding registration when flow is implemented

---

### Gap K.5: Shift+F6 (Reload/Discard) Not Registered

**Specification Location:** `tui-specification-novo.md` line 70

```
| `Shift+F6` | Descartar Alterações e Recarregar Cofre (Fluxo 10) | Similaridade semântica com F6 |
```

**Current Implementation:** *Not implemented*

**Why:** Fluxo 10 (reload/discard flow) not implemented yet (Phase 7+)

**Scope:** Phase 06.1 → Document as "deferred to Phase 7+" in PLAN.md

---

### Gap K.6: Ctrl+Alt+Shift+Q (Lock Vault) Not Registered

**Specification Location:** `tui-specification-novo.md` line 57

```
| `Ctrl+Alt+Shift+Q` | Bloquear Cofre (Fluxo 6) | Global | Bloqueio emergencial, descarta alterações, sem confirmação. Atalho "complicado" para evitar acidentes. |
```

**Current Implementation:** *Not implemented*

**Why:** Fluxo 6 (emergency vault lock) not implemented yet (Phase 6+ consideration)

**Scope:** Phase 06.1 → Document as "deferred to Phase 6.5+ or Phase 7+" in PLAN.md

---

## V. F-Key Grouping Architecture (Design System Reference)

**From `tui-design-system-novo.md` (Mapa de Teclas):**

The spec defines semantic groupings for F-keys:

- **F1, F12:** Global — help, theme
- **Ctrl+Q, Ctrl+Alt+Shift+Q:** Global — exit, lock
- **F2-F4:** Workspace tabs
- **F5-F6:** Workspace entry flows (create, open)
- **Shift+F6:** Reload variant of F6
- **F7, Shift+F7, Ctrl+F7:** Save variants
- **F8:** Reserved
- **F9, Shift+F9:** Export/Import variants
- **F10, F11:** Reserved

This grouping ensures:
- F1-F4 are "always visible" shortcuts
- F5-F6 for main entry flows
- F7-F9 for persistence operations
- Shift/Ctrl modifiers for variants

---

## VI. Decision Matrix: F5/F6 vs. O/N

| Factor | F5/F6 (Spec) | O/N (Current) |
|--------|---|---|
| **Spec Compliance** | ✅ Yes | ❌ No |
| **User Discoverability** | ✅ High (F-keys visible on keyboard) | ⚠️ Medium (must explore or press ?) |
| **Command Bar Space** | ❌ Uses more width (F5/F6 vs O/N) | ✅ Compact |
| **Backward Compatibility** | ❌ Breaks existing bindings | ✅ No change |
| **Help Modal** | ✅ Matches spec | ⚠️ Shows O/N, confuses users |
| **Terminal Portability** | ✅ Standard on all keyboards | ⚠️ Conflict risk with terminal emulators |

**Recommendation:** Migrate to F5/F6 in Phase 06.1 as part of keybinding remediation

---

## VII. Keybindings Remediation Plan (Draft for PLAN.md)

### Phase 06.1 Scope (Welcome Screen)

**Task K.1: Update O→F6, N→F5 keybindings**
- Update `root.go:116-129` Keys arrays: `"o"` → `"f6"`, `"n"` → `"f5"`
- Update help modal display labels
- Update tests (if any hardcode O/N)
- Scope: ~10 lines

**Task K.2: Decide on F12 (Toggle Theme)**
- Document decision in PLAN.md (keep/remove/defer)
- If keeping: update spec or note as experimental
- If removing: delete toggle code from root.go (~40 lines)

**Task K.3: Document F2-F4, F7-F9, Shift+F6, Ctrl+Alt+Shift+Q as deferred**
- Add section to PLAN.md: "Deferred Keybindings (Phase 6.5+, Phase 7+)"
- List which flows depend on which keybindings
- Ensures roadmap clarity

### Phase 06.5+ Scope (Workspace Tabs)

**Task K.4: Register F2, F3, F4 workspace tab switches**
- Depends on: vaultTree, templateList, settings models available
- Scope: ~25 lines in action registration

### Phase 7+ Scope (Vault Persistence)

**Task K.5: Register F7, Shift+F7, Ctrl+F7 (Save variants)**
- Depends on: Fluxo 8, 9, 11 implementation
- Scope: ~15 lines per keybinding

**Task K.6: Register F9, Shift+F9 (Export/Import)**
- Depends on: Fluxo 12, 13 implementation
- Scope: ~15 lines per keybinding

**Task K.7: Register Shift+F6 (Reload/Discard)**
- Depends on: Fluxo 10 implementation
- Scope: ~10 lines

**Task K.8: Register Ctrl+Alt+Shift+Q (Lock Vault)**
- Depends on: Fluxo 6 implementation
- Scope: ~10 lines

---

## Appendix: Keybindings in Help Modal

**File:** `internal/tui/help.go`

The Help modal (`newHelpModal()`) displays all actions via `ActionManager.All()`, grouped by Group ID using `GroupLabel()`.

**Current Groups:**
- Group 0: Global (Ctrl+Q, F12)
- Group 1: Help (F1)
- Group 4: Cofre (O, N) ← should be Fluxos 1-2 or "Vault Actions"

**Impact:** Help modal shows O/N but spec says F5/F6 → user confusion.

**Remediation:** Once O/N → F5/F6 change is made, Help modal automatically reflects it (no code changes needed).

---

## Conclusion

**5 of 13 spec-defined keybindings** currently work (F1, Ctrl+Q, O, N, F12). The remaining 8 are deferred to Phase 06.5+ and Phase 7+ due to flow and model dependencies.

**Phase 06.1 keybindings work:**
1. ✅ Migrate O → F6, N → F5 (spec compliance)
2. ⚠️ Decide on F12 (keep/remove/defer)
3. 📋 Document deferred keybindings (F2-F4, F7-F9, etc.) in PLAN.md

**No blocker:** Keybindings are not critical for Phase 06.1 compliance but improve UX and spec alignment.
