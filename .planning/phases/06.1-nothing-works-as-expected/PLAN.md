# Phase 06.1 PLAN — Remediation Strategy

**Goal:** Identify and fix all specification discrepancies discovered in Phase 06 implementation (Fluxos 1, 2, and 3-5)

**Scope:** 10 gaps across 3 flows; ~360 lines of code changes

**Dependencies:** Phase 06 complete; all 141 tests passing

---

## I. Phase Goal

Remediate specification compliance gaps discovered in Phase 06 so that:

1. **Fluxo 1** (Open Vault) — All error paths loop back to file picker instead of exiting flow
2. **Fluxo 2** (Create Vault) — Password strength validation gate implemented; user can choose to revise weak passwords
3. **Fluxos 3-5** (Exit flows) — All three exit scenarios work as specified:
   - **Fluxo 3:** Confirmation required when exiting with no vault
   - **Fluxo 4:** Confirmation required when exiting with clean vault
   - **Fluxo 5:** "Salvar e sair" fully implemented with external file check + atomic save

---

## II. Execution Strategy

### Phase Structure: 2 Sub-Phases

Since Phase 06.1 has interdependent components (Fluxos 3-5 require new flow architecture), execution is organized as:

- **06.1.A:** Fluxos 1 & 2 fixes (simpler, independent)
- **06.1.B:** Fluxos 3-5 fixes (requires new `saveAndExitFlow` component)

This allows 06.1.A to complete and test independently while 06.1.B is being built.

---

## III. Detailed Task Breakdown

### TIER 1.A: Fluxo 1 Error Path Fixes (Blocking: None | Effort: 0.5 hours)

#### Task 1.A.1: Fix Invalid Magic/Version Error Path

**File:** `internal/tui/flow_open_vault.go`  
**Lines:** 120-128  
**Current Code:**
```go
if errors.Is(err, storage.ErrInvalidMagic) {
    errMsg = "Arquivo não é um cofre válido."
} else if errors.Is(err, storage.ErrVersionTooNew) {
    errMsg = "Versão do cofre não é suportada por esta versão do Abditum."
} else if errors.Is(err, storage.ErrCorrupted) {
    errMsg = "Cofre corrompido e não pode ser recuperado."
}
f.messages.Show(MsgError, errMsg, 5, false)
return endFlow()  // ← WRONG
```

**Change:**
```go
if errors.Is(err, storage.ErrInvalidMagic) {
    errMsg = "Arquivo não é um cofre válido."
    f.messages.Show(MsgError, errMsg, 5, false)
    f.state = statePickFile  // ← volta ao passo 2
    return func() tea.Msg {
        return pushModalMsg{modal: &filePickerModal{}}
    }
} else if errors.Is(err, storage.ErrVersionTooNew) {
    errMsg = "Versão do cofre não é suportada por esta versão do Abditum."
    f.messages.Show(MsgError, errMsg, 5, false)
    f.state = statePickFile  // ← volta ao passo 2
    return func() tea.Msg {
        return pushModalMsg{modal: &filePickerModal{}}
    }
} else if errors.Is(err, storage.ErrCorrupted) {
    errMsg = "Cofre corrompido e não pode ser recuperado."
    f.messages.Show(MsgError, errMsg, 5, false)
    f.state = statePickFile  // ← volta ao passo 2
    return func() tea.Msg {
        return pushModalMsg{modal: &filePickerModal{}}
    }
}
```

**Spec Reference:** `fluxos.md` lines 226, 230 (Step 2 & Step 4)

**Testing:**
- Unit: Mock invalid magic file → should return to file picker
- Unit: Mock version-too-new file → should return to file picker
- Unit: Mock corrupted file → should return to file picker

**Definition of Done:**
- ✓ Code compiles
- ✓ Unit tests pass
- ✓ Error message shown + file picker reopens
- ✓ No `endFlow()` called for these errors

---

#### Task 1.A.2: Verify File Picker Reopens Correctly

**File:** `internal/tui/dialogs.go` (filePickerModal)  
**Purpose:** Ensure filePickerModal can be instantiated multiple times in same flow

**Testing:**
- Integration: Open Vault → enter invalid file path → error shown → file picker shown again → user can select different file

**Definition of Done:**
- ✓ File picker modal shows on retry
- ✓ Previous file path not cached/pre-selected
- ✓ User can navigate and select different file

---

### TIER 1.B: Fluxo 2 Fix #1 — Overwrite Dialog (Blocking: None | Effort: 1 hour)

#### Task 1.B.1: Handle Overwrite "Voltar" Decision

**File:** `internal/tui/flow_create_vault.go`  
**Lines:** 62-115 (Update method)  
**Current Issue:** When user clicks "Voltar" on overwrite dialog, flow exits instead of returning to file picker

**Change:** Add explicit handling for Decision result from overwrite dialog:

```go
// After case filePickerResult at line 64:
// ... existing code until line 80 (file exists check) ...

case decisionMsg:  // ← NEW: Handle decision modal results
    if msg.Title == "Arquivo existe" {  // ← Identify overwrite dialog
        if msg.Confirmed {  // ← User chose "Sobrescrever"
            f.state = statePwdCreate
            return func() tea.Msg {
                return pushModalMsg{modal: &passwordCreateModal{}}
            }
        } else {  // ← User chose "Voltar"
            f.state = statePickFile
            return func() tea.Msg {
                return pushModalMsg{modal: &filePickerModal{mode: FilePickerFile}}
            }
        }
    }
```

**Problem:** Need to identify which Decision dialog result is being processed. Two approaches:

**Approach A (Cleaner):** Create custom message type for overwrite decision
```go
type overwriteDecisionMsg struct {
    Confirmed bool  // true = "Sobrescrever", false = "Voltar"
}
```

Then modify decision.go or create new modal type for overwrite confirmation.

**Approach B (Quick Fix):** Check msg.Title to identify dialog (fragile but works now)

**Recommendation:** Use Approach A for robustness. Modify filePickerResult or create `filePickerOverwriteConflictMsg`.

**Testing:**
- Unit: Trigger overwrite dialog → click "Sobrescrever" → continue to password
- Unit: Trigger overwrite dialog → click "Voltar" → return to file picker
- Integration: Create Vault → select existing file → overwrite dialog → "Voltar" → select different file → proceed

**Definition of Done:**
- ✓ Code compiles
- ✓ Both decision paths tested
- ✓ File picker reopens after "Voltar"
- ✓ No flow exit on "Voltar"

---

### TIER 2.A: Fluxo 2 Fix #2 — Password Strength Gate (Blocking: None | Effort: 2-3 hours)

#### Task 2.A.1: Add stateStrengthCheck to createVaultFlow

**File:** `internal/tui/flow_create_vault.go`  
**Current States:**
```go
const (
    stateCheckDirty = iota       // 0 (reused)
    statePickFile                // 1 (reused)
    stateCheckOverwrite = iota + 10  // 10
    statePwdCreate               // 11
    stateSaveNew                 // 12
)
```

**New State:**
```go
const (
    // ... existing ...
    stateStrengthCheck = 13  // NEW: Between password creation and save
)
```

**Change in Init/Update:**

```go
// In Update method, modify pwdCreatedMsg case:
case pwdCreatedMsg:
    f.state = stateStrengthCheck  // ← NEW: Don't go directly to save
    password := msg.Password
    
    // Check password strength (delegated to modal)
    // For now, emit stateStrengthCheck signal
    // Modal will determine if weak and emit response
    
    return func() tea.Msg {
        // Check strength synchronously
        strength := evaluatePasswordStrength(password)
        if strength < ThresholdStrong {  // Assume constant defined
            return pwdWeakMsg{
                Password: password,
                Strength: strength,
            }
        } else {
            // Proceed directly to save
            return pwdAcceptableMsg{Password: password}
        }
    }
```

**Simplification:** Instead of modal handling, handle strength check in flow itself:

```go
case pwdCreatedMsg:
    password := msg.Password
    strength := evaluatePasswordStrength(password)  // New function
    
    if strength < ThresholdStrong {
        f.state = stateStrengthCheck
        return func() tea.Msg {
            return Decision(SeverityWarn, "Senha fraca",
                "A senha é fraca. Critérios não atendidos: [...]\n\nDeseja prosseguir mesmo assim ou revisar?",
                DecisionAction{Key: "Enter", Label: "Prosseguir", Default: false},
                nil,
                DecisionAction{Key: "R", Label: "Revisar"})  // ← R for Revisar
        }
    } else {
        // Strong password: proceed to save
        f.state = stateSaveNew
        password := msg.Password
        path := f.targetPath
        // ... existing save logic
    }
```

**New Message Type:**
```go
// In flows.go or flow_create_vault.go:
type weakPasswordDecisionMsg struct {
    Confirmed bool  // true = "Prosseguir", false = "Revisar"
}
```

**Handle Decision Result:**
```go
case decisionMsg:
    if msg.Title == "Senha fraca" {
        if msg.Confirmed {  // ← "Prosseguir"
            f.state = stateSaveNew
            // proceed to save
        } else {  // ← "Revisar"
            f.state = statePwdCreate
            return func() tea.Msg {
                return pushModalMsg{modal: &passwordCreateModal{}}
            }
        }
    }
```

**Helper Function (New):**
```go
func evaluatePasswordStrength(pwd string) int {
    // Implementation: check length, entropy, common patterns, etc.
    // Return 1-5 (1=very weak, 5=very strong)
    // For MVP: simple length check
    if len(pwd) < 8 {
        return 1  // Very weak
    } else if len(pwd) < 12 {
        return 2  // Weak
    } else if len(pwd) < 16 {
        return 3  // Acceptable
    } else {
        return 4  // Strong
    }
    // TODO: Add entropy check via password strength library
}

const ThresholdStrong = 3  // Acceptable or above
```

**Testing:**
- Unit: passwordCreateModal returns weak password → Decision dialog shown
- Unit: User chooses "Prosseguir" → continue to save
- Unit: User chooses "Revisar" → return to password modal
- Unit: passwordCreateModal returns strong password → skip dialog, go to save
- Integration: Create Vault → enter weak password → decision dialog → choose "Revisar" → password modal reopens

**Definition of Done:**
- ✓ Code compiles
- ✓ Weak password detection works
- ✓ Decision dialog shown for weak passwords
- ✓ "Prosseguir" path continues to save
- ✓ "Revisar" path returns to password modal
- ✓ Strong passwords skip dialog
- ✓ Integration tests pass

---

### TIER 3: Fluxos 3-5 Exit Flow Fixes (Blocking: 1.A, 1.B, 2.A | Effort: 4-5 hours)

#### Task 3.1: Create New Message Types and Modal

**Files:** `internal/tui/flows.go`, `internal/tui/decision.go`  
**Purpose:** Define messages for exit flow coordination

**New Message Types:**
```go
// In flows.go:

type exitConfirmationMsg struct {
    HasVault    bool  // true = vault loaded
    IsDirty     bool  // true = unsaved changes
}

type exitDecisionMsg struct {
    Choice string  // "confirm" (exit), "cancel" (stay), "save" (Fluxo 5 only)
}

type externalModificationConflictMsg struct {
    Path string  // Path of modified file
}

type externalModificationDecisionMsg struct {
    Choice string  // "overwrite" (proceed), "cancel" (stay)
}

type saveAndExitProgressMsg struct {
    Status string  // "checking", "saving", "complete"
    Error  error   // if applicable
}
```

**Definition of Done:**
- ✓ Message types defined in flows.go
- ✓ Exported properly
- ✓ Used consistently in exit flow

---

#### Task 3.2: Implement Decision Dialogs for Fluxos 3 & 4

**File:** `internal/tui/root.go`  
**Lines:** 243-257  
**Current Code:**
```go
if key == "ctrl+q" {
    if m.mgr != nil && m.mgr.IsModified() {
        return m, func() tea.Msg {
            return Decision(...)  // Fluxo 5 dialog
        }
    }
    return m, tea.Quit  // ← Fluxos 3 & 4: direct quit
}
```

**Change:**
```go
if key == "ctrl+q" {
    if m.mgr != nil && m.mgr.IsModified() {
        // Fluxo 5: Dirty vault
        return m, func() tea.Msg {
            return Decision(SeverityNeutral, "Alterações não salvas",
                "Deseja salvar as alterações antes de sair?",
                DecisionAction{Key: "Enter", Label: "Salvar", Default: true},
                []DecisionAction{{Key: "D", Label: "Descartar"}},
                DecisionAction{Key: "Esc", Label: "Voltar"})
        }
    } else if m.mgr != nil {
        // Fluxo 4: Clean vault (no unsaved changes)
        return m, func() tea.Msg {
            return Decision(SeverityNeutral, "Sair do Abditum",
                "Tem certeza que deseja sair?",
                DecisionAction{Key: "Enter", Label: "Sim", Default: true},
                nil,
                DecisionAction{Key: "Esc", Label: "Não"})
        }
    } else {
        // Fluxo 3: No vault loaded
        return m, func() tea.Msg {
            return Decision(SeverityNeutral, "Sair do Abditum",
                "Tem certeza que deseja sair?",
                DecisionAction{Key: "Enter", Label: "Sim", Default: true},
                nil,
                DecisionAction{Key: "Esc", Label: "Não"})
        }
    }
}
```

**Testing:**
- Unit: No vault loaded → Ctrl+Q → confirmation dialog shown
- Unit: Clean vault → Ctrl+Q → confirmation dialog shown
- Unit: Dirty vault → Ctrl+Q → "unsaved changes" dialog shown (Fluxo 5)
- Unit: Confirmation dialog → Enter → tea.Quit
- Unit: Confirmation dialog → Esc → dismiss, stay in app

**Definition of Done:**
- ✓ Code compiles
- ✓ All three dialogs shown correctly
- ✓ Confirmation required for all exit scenarios
- ✓ Spec compliance verified

---

#### Task 3.3: Wire Decision Dialog Results to Handler

**File:** `internal/tui/root.go`  
**Purpose:** Handle Decision modal results (currently unhandled)

**Current Issue:** Decision modal is shown but result is not processed.

**Implementation:** Add case for Decision result in root.Update():

```go
// In Update method, add new case:
case decisionMsg:
    if msg.Title == "Sair do Abditum" {
        // Fluxo 3 & 4: Simple exit confirmation
        if msg.Confirmed {
            return m, tea.Quit
        }
        return m, nil  // User chose "Não", stay in app
    } else if msg.Title == "Alterações não salvas" {
        // Fluxo 5: Unsaved changes decision
        if msg.Label == "Salvar" || msg.Confirmed {
            // Start save-and-exit flow
            flow := newSaveAndExitFlow(m.mgr, m.messages, m.actions, m.theme)
            return m, func() tea.Msg { return startFlowMsg{flow: flow} }
        } else if msg.Label == "Descartar" {
            // Discard and exit immediately
            return m, tea.Quit
        }
        // "Voltar" = dismissed, stay in app
        return m, nil
    }
    return m, nil
```

**Problem:** Need to identify which action was chosen from Decision dialog. Decision modal should emit message with action label or result status.

**Review:** Check `internal/tui/decision.go` to understand how Decision result is communicated.

**Definition of Done:**
- ✓ Code compiles
- ✓ Decision results properly identified
- ✓ Correct action taken for each path
- ✓ "Salvar" triggers saveAndExitFlow (next task)

---

#### Task 3.4: Create saveAndExitFlow (New Component)

**File:** `internal/tui/flow_save_and_exit.go` (NEW)  
**Purpose:** Implement Fluxo 8 logic (atomic save) + Fluxo 5 Step 3-4 (external check + save + exit)

**Structure:**

```go
package tui

import (
    "os"
    "time"
    tea "charm.land/bubbletea/v2"
    "github.com/useful-toys/abditum/internal/vault"
    "github.com/useful-toys/abditum/internal/storage"
    "github.com/useful-toys/abditum/internal/crypto"
)

type saveAndExitFlow struct {
    state           int
    mgr             *vault.Manager
    messages        *MessageManager
    actions         *ActionManager
    theme           *Theme
    lastModTime     time.Time  // Last known modification time of vault file
    checkingExternal bool
    backupCreated   bool
}

const (
    stateCheckExternal = iota
    stateSaveVault
    stateDone
)

func newSaveAndExitFlow(mgr *vault.Manager, messages *MessageManager, actions *ActionManager, theme *Theme) *saveAndExitFlow {
    return &saveAndExitFlow{
        state:    stateCheckExternal,
        mgr:      mgr,
        messages: messages,
        actions:  actions,
        theme:    theme,
    }
}

func (f *saveAndExitFlow) Init() tea.Cmd {
    // Step 3: Check for external modification
    return func() tea.Msg {
        vaultPath := f.mgr.VaultPath()  // Assume method exists
        fileInfo, err := os.Stat(vaultPath)
        if err != nil {
            f.messages.Show(MsgError, "Não foi possível acessar o arquivo do cofre.", 5, false)
            return endFlow()
        }
        
        // Compare modification time with last-known time
        // NOTE: Need to store lastModTime when vault is loaded
        // For now, just proceed (TODO: implement proper check)
        
        f.state = stateSaveVault
        return nil
    }
}

func (f *saveAndExitFlow) Update(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case externalModificationConflictMsg:
        // User responded to external modification dialog
        if msg.Choice == "overwrite" {
            f.state = stateSaveVault
            return f.saveVault()
        } else {
            // User chose to cancel exit
            return endFlow()
        }
    
    case saveAndExitProgressMsg:
        if msg.Error != nil {
            f.messages.Show(MsgError, msg.Error.Error(), 5, false)
            if f.backupCreated {
                f.messages.Show(MsgWarn, "Backup disponível em .abditum.bak", 5, false)
            }
            return endFlow()
        }
        
        if msg.Status == "complete" {
            f.state = stateDone
            return tea.Quit
        }
        return nil
    
    case flowCancelledMsg:
        return endFlow()
    
    default:
        return nil
    }
}

func (f *saveAndExitFlow) saveVault() tea.Cmd {
    return func() tea.Msg {
        vaultPath := f.mgr.VaultPath()
        password := f.mgr.Password()  // Assume method exists to get current session password
        
        // Step 4: Atomic save
        err := storage.Save(vaultPath, f.mgr.Vault(), password, true)  // true = atomic
        crypto.Wipe(password)
        
        if err != nil {
            // Check if backup was created
            backupPath := vaultPath + ".abditum.bak"
            if _, err := os.Stat(backupPath); err == nil {
                f.backupCreated = true
            }
            return saveAndExitProgressMsg{
                Status: "error",
                Error:  err,
            }
        }
        
        return saveAndExitProgressMsg{
            Status:   "complete",
            Error:    nil,
        }
    }
}

func (f *saveAndExitFlow) View(width, height int) string {
    return ""  // No modal display; progress shown in message bar
}
```

**Key Points:**
1. **State machine:** Check external → Save → Done → Exit
2. **External modification check:** Compare file mtime with stored last-read time
3. **Atomic save:** Use storage.Save() with atomic flag
4. **Backup handling:** Check for .abditum.bak after failure
5. **Exit condition:** Only call tea.Quit after successful save

**Testing:**
- Unit: Flow starts → checks external mod → none found → saves → exits
- Unit: Flow starts → external mod detected → shows dialog → user confirms → saves → exits
- Unit: Save fails without backup → error shown → flow ends without quit
- Unit: Save fails with backup → error + backup info shown → flow ends without quit
- Integration: Dirty vault → Ctrl+Q → "Salvar" → save succeeds → app exits

**Definition of Done:**
- ✓ Code compiles
- ✓ External check logic works
- ✓ Atomic save integrated
- ✓ Backup protocol respected
- ✓ Error handling prevents data loss
- ✓ Integration tests pass

---

#### Task 3.5: Implement External File Modification Check Modal

**File:** `internal/tui/decision.go` or new file `internal/tui/external_mod_modal.go`  
**Purpose:** Present conflict dialog when vault file modified externally

**Modal Content:**
```
Title: Conflito de Modificação
Message: O arquivo do cofre foi modificado externamente desde a última leitura.
         Sobrescrever com as alterações em memória?
Actions:
  - Enter: "Sobrescrever e sair"
  - Esc: "Voltar"
```

**Usage in saveAndExitFlow:**
```go
// In Init or Update:
if fileWasModifiedExternally {
    return func() tea.Msg {
        return Decision(SeverityWarn, "Conflito de Modificação",
            "O arquivo do cofre foi modificado externamente...",
            DecisionAction{Key: "Enter", Label: "Sobrescrever e sair", Default: false},
            nil,
            DecisionAction{Key: "Esc", Label: "Voltar"})
    }
}
```

**Definition of Done:**
- ✓ Modal displays correctly
- ✓ External modification detection works
- ✓ User can choose to overwrite or cancel
- ✓ Messages wired to saveAndExitFlow

---

### Task 3.6: Integration Testing

**Test Suite:** `internal/tui/exit_flow_integration_test.go` (NEW)

**Test Cases:**

1. **Fluxo 3: No vault loaded**
   - Setup: App starts, no vault
   - Action: Press Ctrl+Q
   - Expected: "Sair?" dialog shown
   - Confirm: App exits
   - Cancel: App stays running

2. **Fluxo 4: Clean vault**
   - Setup: Open vault, make no changes
   - Action: Press Ctrl+Q
   - Expected: "Sair?" dialog shown
   - Confirm: App exits
   - Cancel: App stays running

3. **Fluxo 5: Dirty vault → Salvar e sair (success)**
   - Setup: Open vault, make changes
   - Action: Press Ctrl+Q
   - Expected: "Alterações não salvas" dialog shown
   - Choose "Salvar": saveAndExitFlow starts
   - Expected: External check → Save → Exit
   - Result: App exits after successful save

4. **Fluxo 5: Dirty vault → Descartar e sair**
   - Setup: Open vault, make changes
   - Action: Press Ctrl+Q
   - Expected: "Alterações não salvas" dialog shown
   - Choose "Descartar": App exits immediately

5. **Fluxo 5: Dirty vault → Voltar**
   - Setup: Open vault, make changes
   - Action: Press Ctrl+Q
   - Expected: "Alterações não salvas" dialog shown
   - Choose "Voltar": Dialog dismisses, app continues

6. **Fluxo 5: External modification detected**
   - Setup: Open vault, make changes, externally modify file
   - Action: Press Ctrl+Q → "Salvar"
   - Expected: External conflict dialog shown
   - Choose "Sobrescrever": Save with overwrite → Exit
   - Choose "Voltar": Dialog dismisses, app continues

7. **Fluxo 5: Save failure without backup**
   - Setup: Open vault, make changes, simulate save failure
   - Action: Press Ctrl+Q → "Salvar"
   - Expected: Error shown, vault remains dirty, app continues

8. **Fluxo 5: Save failure with backup**
   - Setup: Open vault, make changes, simulate save failure with backup creation
   - Action: Press Ctrl+Q → "Salvar"
   - Expected: Error + backup info shown, vault remains dirty, app continues

**Definition of Done:**
- ✓ All 8 test cases pass
- ✓ Golden files generated for each scenario
- ✓ Coverage > 90% for exit code paths

---

## IV. Task Dependencies

```
1.A.1 (File validation error paths)
  ↓
1.A.2 (File picker retry)
  ↓
1.B.1 (Overwrite dialog handling)
  ↓
2.A.1 (Password strength gate) ← Can parallel with 1.A/1.B
  ↓
3.1 (Message types) ← Depends on all Tier 1 & 2
  ↓
3.2 (Confirmation dialogs) ← Parallel with 3.1
3.3 (Wire decision results)
  ↓
3.4 (saveAndExitFlow) ← Depends on 3.2 & 3.3
  ↓
3.5 (External mod modal)
  ↓
3.6 (Integration testing)
```

### Parallelization Opportunities

- **1.A** and **1.B** can be done in parallel (different flows)
- **2.A** can be done in parallel with **1.A/1.B**
- **3.1** must wait for 1.A/1.B/2.A to complete
- **3.2** and **3.3** can be done in parallel after 3.1
- **3.4** and **3.5** must wait for 3.2 & 3.3
- **3.6** is final validation

---

## V. Success Criteria

### Phase 06.1 Complete When:

1. **All 10 gaps fixed:**
   - ✓ Gap 1.1: Invalid magic error returns to file picker
   - ✓ Gap 1.2: Corrupted payload error returns to file picker
   - ✓ Gap 2.1: Password strength validation gate implemented
   - ✓ Gap 2.2: Overwrite "Voltar" returns to file picker
   - ✓ Gap 3.1: Fluxo 3 confirmation dialog shown
   - ✓ Gap 3.2: Fluxo 4 confirmation dialog shown
   - ✓ Gap 3.3: "Salvar e sair" fully implemented
   - ✓ Gap 3.4: External modification check implemented
   - ✓ Gap 3.5: Atomic save execution working
   - ✓ Gap 3.6: Backup communication working

2. **All tests pass:**
   - ✓ 141 original Phase 06 tests still passing
   - ✓ New Phase 06.1 tests added and passing (target: 20+ new tests)
   - ✓ Integration tests for exit flows passing

3. **Code quality:**
   - ✓ No compilation errors
   - ✓ No linting issues
   - ✓ All code changes documented with spec references
   - ✓ Commit messages reference spec sections

4. **Specification compliance:**
   - ✓ Fluxo 1 behavior matches spec Step 2 & 4
   - ✓ Fluxo 2 behavior matches spec Step 2, 3, 4, 5
   - ✓ Fluxo 3 behavior matches spec Step 1, 2
   - ✓ Fluxo 4 behavior matches spec Step 1, 2
   - ✓ Fluxo 5 behavior matches spec Step 1-5

---

## VI. Risk Assessment

### High-Risk Areas

1. **External File Modification Detection**
   - Risk: Comparing file timestamps may not work reliably on all systems
   - Mitigation: Store file hash/size in vault metadata; compare against current disk state
   - Alternative: Use file watcher library for robust implementation

2. **Atomic Save Integration**
   - Risk: storage.Save() may not have atomic flag or proper backup handling
   - Mitigation: Review storage package implementation before integration
   - Alternative: Implement atomic save wrapper if storage package lacks it

3. **Message Routing in root.Update()**
   - Risk: Decision modal results may be routed to wrong handler
   - Mitigation: Add explicit message type discrimination; test each path
   - Alternative: Use flow-based decision handling instead of root model

### Medium-Risk Areas

1. **File Picker Multiple Instantiation** — Ensure no state is cached between calls
2. **Password Strength Evaluation** — Implement proper strength calculation (not just length)
3. **Backup File Availability Communication** — Ensure message is clear and actionable

---

## VII. Estimated Timeline

Assuming 1 developer:

| Task | Hours | Effort |
|------|-------|--------|
| 1.A.1 + 1.A.2 | 0.5 | Small |
| 1.B.1 | 1 | Small |
| 2.A.1 | 2-3 | Medium |
| 3.1 | 0.5 | Small |
| 3.2 + 3.3 | 1.5 | Medium |
| 3.4 | 2-3 | Medium |
| 3.5 | 0.5 | Small |
| 3.6 + Testing | 1-2 | Medium |
| **TOTAL** | **9-12 hours** | — |

**Note:** With parallel work on different flows, can reduce to **7-8 hours** wall-clock time.

---

## VIII. Artifacts to Create

1. ✓ **RESEARCH.md** — Gap analysis (done)
2. ✓ **PLAN.md** — This document
3. **06.1-IMPLEMENTATION.md** — Task execution log (to be created during build)
4. **New files:**
   - `internal/tui/flow_save_and_exit.go` — saveAndExitFlow implementation
   - `internal/tui/exit_flow_integration_test.go` — Exit flow tests

5. **Modified files:**
   - `internal/tui/flow_open_vault.go` — Error path fixes
   - `internal/tui/flow_create_vault.go` — Overwrite dialog + strength gate
   - `internal/tui/root.go` — Confirmation dialogs + decision handling
   - `internal/tui/passwordcreate.go` — Strength evaluation function (if needed)

6. **Test artifacts:**
   - Golden files for all exit flow scenarios
   - Unit tests for new components

---

## IX. Keybindings Remediation (Scope Extension)

**Note:** Keybindings are NOT blocking for Phase 06.1 completion but are part of broader spec compliance effort. See `KEYBINDINGS-SPEC.md` for full analysis.

### Phase 06.1 Keybindings Tasks

#### Task K.1: Migrate O→F6, N→F5

**File:** `internal/tui/root.go`  
**Lines:** 116-129  
**Effort:** ~10 minutes

**Change:**
```go
// Current (lines 116-122):
Action{Keys: []string{"o"}, Label: "Abrir", ...},

// Change to:
Action{Keys: []string{"f6"}, Label: "Abrir", ...},

// Current (lines 123-129):
Action{Keys: []string{"n"}, Label: "Novo", ...},

// Change to:
Action{Keys: []string{"f5"}, Label: "Novo", ...},
```

**Testing:**
- Verify F5 triggers "Criar Novo Cofre" flow
- Verify F6 triggers "Abrir Cofre Existente" flow
- Verify help modal shows F5/F6 (automatic via ActionManager)
- Update any tests that hardcode "o" or "n" key presses

**Spec Reference:** `tui-specification-novo.md` lines 68-69

**Definition of Done:**
- ✓ F5/F6 keys work
- ✓ O/N keys no longer trigger actions
- ✓ Help modal displays correctly
- ✓ Tests pass

---

#### Task K.2: Formalize F12 (Toggle Theme) in Spec

**DECISION (2026-04-06): KEEP — Formalize as feature.**

**File:** `tui-specification-novo.md` (Global Actions table)  
**Effort:** ~5 minutes

F12 Toggle Theme is already implemented and working. Add it to the spec's Global Actions table as a documented feature.

**Change:** Update `tui-specification-novo.md` line 55 — F12 row already exists and is correct. No code changes needed.

**Definition of Done:**
- ✓ F12 is documented in spec as a formal feature
- ✓ No code changes required

---

#### Task K.3: Document Deferred Keybindings

**File:** This PLAN.md + KEYBINDINGS-SPEC.md

**Keybindings deferred to Phase 06.5+:**
- F2, F3, F4 (Workspace tabs) — depends on vaultTree/templateList/settings models

**Keybindings deferred to Phase 7+:**
- F7, Shift+F7, Ctrl+F7 (Save variants) — depends on Fluxo 8, 9, 11
- F9, Shift+F9 (Export/Import) — depends on Fluxo 12, 13
- Shift+F6 (Reload) — depends on Fluxo 10
- Ctrl+Alt+Shift+Q (Lock vault) — depends on Fluxo 6

**Effort:** ~5 minutes (already documented in KEYBINDINGS-SPEC.md)

**Definition of Done:**
- ✓ Deferred keybindings listed in PLAN.md
- ✓ Phase dependencies documented
- ✓ Roadmap clarity improved

---

### Keybindings Success Criteria

Phase 06.1 keybindings work:

1. ✓ F5/F6 migrate from O/N (spec compliance)
2. ✓ F12 decision documented (feature or removed)
3. ✓ Deferred keybindings identified (Phase 06.5+, Phase 7+)
4. ✓ Help modal reflects new keybindings (automatic)
5. ✓ Tests updated for F5/F6

**Estimated Effort:** 1-1.5 hours total

---

## Conclusion

Phase 06.1 is a **focused remediation phase** that fixes critical specification compliance gaps discovered in Phase 06. The work is divided into three tiers with clear dependencies, allowing parallel execution of independent tasks. A keybindings extension is included to improve spec alignment and user experience. Upon completion, vault operations will fully comply with the `@fluxos` specification, and keybindings will align with `@tui-specification-novo.md`.
