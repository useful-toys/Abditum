# Flow Implementation Mapping - Abditum

**Generated:** April 6, 2026

This document maps the implementation of flows 1, 2, and 3-5 (exit flows) against the @fluxos specification (fluxos.md).

---

## FLUXO 1: Abrir Cofre Existente (Open Existing Vault)

### Implementation: flow_open_vault.go

**File Location:** C:\git\Abditum\internal\tui\flow_open_vault.go

**Lines:** 1-151 (total file)

**Main Logic Line Numbers:**
- Flow initialization: Lines 52-78
- Message handling/state transitions: Lines 81-146
- State constants: Lines 30-36
- Flow struct definition: Lines 15-27

---

### Current Step Sequence

The implementation follows a state machine pattern with 5 states:

\\\
State: stateCheckDirty (0)
  -> Init checks if vault.IsModified()
      ├- YES: Show Acknowledge dialog (lines 57-59)
      └- NO or CLI fast-path: Proceed to statePickFile

State: statePickFile (2)
  -> Init or filePickerResult handler
      ├- CLI fast-path: Skip to statePwdEntry (lines 62-71)
      └- Normal path: Push filePickerModal (lines 73-77)

State: statePwdEntry (3)
  -> On filePickerResult (lines 83-93)
      └-> Push passwordEntryModal

State: statePreload (4)
  -> On pwdEnteredMsg (lines 95-138)
      └-> Background command: storage.Load()
          ├- Success: Transition to stateDone, emit vaultOpenedMsg
          └- Failure: Handle auth or structural errors
              ├- Wrong password: Retry (max 5 attempts, lines 108-116)
              └- File errors: Show error, endFlow()

State: stateDone (5) [implicit completion]
  -> vaultOpenedMsg emitted, flow completes
\\\

**Line References for Each State:**
- **stateCheckDirty:** Lines 52-60 (Init method handling)
- **statePickFile:** Lines 62-77 (CLI fast-path or file picker)
- **statePwdEntry:** Lines 88-93 (transition on filePickerResult)
- **statePreload:** Lines 96-138 (password validation and vault loading)

---

### Decisions/Branches

| Line Range | Condition | Branch A | Branch B | Branch C |
|---|---|---|---|---|
| 52-60 | mgr.IsModified()? | Show dirty dialog | Skip to file picker | N/A |
| 62-71 | cliPath != ""? | Skip to password | Show file picker | N/A |
| 84-85 | filePickerResult.Cancelled? | endFlow() | Continue | N/A |
| 107-129 | Load result | Success -> vaultOpenedMsg | Wrong pwd (max 5) | File error -> endFlow |
| 108-116 | Auth failure & attempts < 5 | Retry pwd entry (line 116) | Exceeded -> error | N/A |
| 120-128 | File validation errors | Magic invalid | Version too new | Corrupted |

**Specification Alignment (Fluxo 1, Step 1):**
- ✓ Checks for unsaved changes (line 52-60) -> matches spec step 1
- ✓ Offers 3-way decision dialog (implied via Acknowledge on line 57)
- ✗ **Deviation:** Implementation does NOT handle 'Salvar' option in flow init; this is delegated to rootModel (root.go line 245-256)

---

### Modals and Dialogs

| Modal | Function | Triggered Line | Dialog Type | User Actions |
|---|---|---|---|---|
| **Decision Dialog** (unsaved changes) | Show when vault modified | 57 | Acknowledge() with SeverityNeutral | N/A in flow (handled by root.go) |
| **filePickerModal** | Select vault file | 76 | pushModalMsg | Up/Down navigate, Tab switch panels, Enter select, Esc cancel |
| **passwordEntryModal** | Enter password | 92 | pushModalMsg | Type password, Enter confirm (with retry counter after 2 attempts), Esc cancel |
| **Message Bar** | Show status | 110, 115, 127 | MessageManager.Show() | None (user-initiated) |

---

### Error Handling

| Error Type | Detected at Line | Error Message | Recovery Path |
|---|---|---|---|
| **Unsaved changes** | 53 | 'Alterações não salvas' | Show dialog (spec: salvar/descartar/voltar) |
| **Wrong password** | 107 | 'Senha incorreta. Tente novamente' | Retry password entry (max 5 times, line 109) |
| **Attempt limit exceeded** | 109 | 'Limite de tentativas excedido' | endFlow() |
| **Invalid magic** | 120 | 'Arquivo não é um cofre válido.' | Return to file picker (implied) |
| **Unsupported version** | 122 | 'Versão do cofre não é suportada...' | Return to file picker (implied) |
| **Corrupted vault** | 124 | 'Cofre corrompido e não pode ser recuperado.' | endFlow() |
| **Other file errors** | 119 | Generic 'Não foi possível abrir o cofre...' | Return to file picker (implied) |

---

## FLUXO 2: Criar Novo Cofre (Create New Vault)

### Implementation: flow_create_vault.go

**File Location:** C:\git\Abditum\internal\tui\flow_create_vault.go

**Lines:** 1-120 (total file)

**Main Logic Line Numbers:**
- Flow initialization: Lines 45-59
- Message handling/state transitions: Lines 62-114
- State constants: Lines 26-30
- Flow struct definition: Lines 14-23

---

### Current Step Sequence

The implementation follows a state machine pattern with 5+ states:

\\\
State: stateCheckDirty (0)
  -> Init checks if vault.IsModified()
      ├- YES: Show Acknowledge dialog (lines 49-51)
      └- NO: Proceed to statePickFile

State: statePickFile (2)
  -> On filePickerResult (lines 64-85)
      └-> Check if file exists (line 71)
          ├- Exists: Show overwrite confirmation dialog (lines 73-79)
          └- Not exists: Proceed to statePwdCreate

State: stateCheckOverwrite (10)
  -> On Decision result (lines 73-79)
      └-> Show overwrite Decision dialog
          ├- 'Sobrescrever': Proceed to statePwdCreate
          └- 'Voltar': Return to file picker

State: statePwdCreate (11)
  -> Show passwordCreateModal (line 84)

State: stateSaveNew (12)
  -> On pwdCreatedMsg (lines 87-107)
      └-> Background command: storage.SaveNew()
          ├- Success: Emit vaultOpenedMsg
          └- Failure: Show error, endFlow()
\\\

---

### Decisions/Branches

| Line Range | Condition | Branch A | Branch B | Branch C |
|---|---|---|---|---|
| 45-53 | mgr.IsModified()? | Show dirty dialog | Skip to file picker | N/A |
| 64-66 | filePickerResult.Cancelled? | endFlow() | Continue | N/A |
| 71 | os.Stat(targetPath)? | File exists (line 71) | File not exists -> pwd create | N/A |
| 73-79 | File exists decision | Overwrite -> pwd create | Voltar -> file picker | N/A |
| 87-107 | Save result | Success -> vaultOpenedMsg | Failure -> error + endFlow | N/A |

---

### Modals and Dialogs

| Modal | Function | Triggered Line | Dialog Type | User Actions |
|---|---|---|---|---|
| **Decision Dialog** (unsaved changes) | Show if vault modified | 50 | Acknowledge() with SeverityNeutral | N/A in flow |
| **Decision Dialog** (overwrite confirmation) | File already exists | 74 | Decision() with SeverityDestructive | Enter 'Sobrescrever' or Esc 'Voltar' |
| **filePickerModal** | Select new vault path | 57 | pushModalMsg | Up/Down navigate, Tab switch, Enter select, Esc cancel |
| **passwordCreateModal** | Create password twice | 84 | pushModalMsg | Type pwd, Tab to confirm, Enter create, Esc cancel |
| **Message Bar** | Show save errors | 102 | MessageManager.Show() | None |

---

### Error Handling

| Error Type | Detected at Line | Error Message | Recovery Path |
|---|---|---|---|
| **File exists** | 71 | 'Um cofre já existe neste caminho. Deseja sobrescrever?' | User chooses sobrescrever or return to file picker |
| **Password mismatch** | 102 | 'As senhas nao conferem' | passwordcreate.go line 103; user can retry within modal |
| **Save failure** | 101 | 'Não foi possível salvar o cofre.' | endFlow() (no recovery path) |
| **Empty password** | 94, 98 | (Within passwordCreateModal) 'Digite uma senha' / 'Confirme a senha' | passwordcreate.go lines 95-100; user retries within modal |

---

## FLUXO 3, 4, 5: Exit Flows (Sair / Ctrl+Q)

### Implementation: root.go

**File Location:** C:\git\Abditum\internal\tui\root.go

**Lines:** 243-257 (exit logic), plus related lines

**Main Logic Line Numbers:**
- Exit/quit action registration: Line 111-114
- Ctrl+Q key handling: Lines 243-257
- Theme toggle (F12) comparison: Lines 260-268
- Context of dispatch order: Lines 236-292

---

### Current Step Sequence

The implementation handles 3 exit scenarios based on vault state:

\\\
User presses Ctrl+Q
  |
  -> root.go line 244: if key == 'ctrl+q'
      |
      ├- Case A: mgr == nil (NO VAULT LOADED)
      |          -> Line 256: tea.Quit (immediate exit)
      |
      ├- Case B: mgr != nil && mgr.IsModified() == false (VAULT LOADED, CLEAN)
      |          -> Line 256: tea.Quit (immediate exit, per spec Fluxo 4)
      |
      └- Case C: mgr != nil && mgr.IsModified() == true (VAULT LOADED, DIRTY)
                 -> Lines 247-253: Show Decision dialog
                     ├- 'Salvar' (default): [NOT IMPLEMENTED - would need additional flow]
                     ├- 'Descartar': [Implicit - just quit without executing tea.Quit]
                     └- 'Voltar' (Esc): Return to normal operation
\\\

---

### Decisions/Branches

| Line Range | Condition | Branch A | Branch B | Branch C |
|---|---|---|---|---|
| 246 | mgr != nil && mgr.IsModified()? | Show dialog (line 248) | Quit immediately (line 256) | N/A |
| 248-252 | Dialog shown; user presses... | 'Enter' = Salvar | 'D' = Descartar | 'Esc' = Voltar |

---

### Modals and Dialogs

| Dialog | Function | Triggered Line | Severity | Actions |
|---|---|---|---|---|
| **Decision Dialog** | Unsaved changes on exit | 248 | SeverityNeutral | Enter 'Salvar' + D 'Descartar' + Esc 'Voltar' |

**Dialog Details (Fluxo 5 - Line 248-252):**
- Title: 'Alterações não salvas'
- Message: 'Deseja salvar as alterações antes de sair?'
- Type: SeverityNeutral
- Actions: 
  - Key 'Enter', Label 'Salvar', Default: true (BUT Cmd: nil - NOT IMPLEMENTED)
  - Key 'D', Label 'Descartar' (Cmd: nil - quit happens implicitly)
  - Key 'Esc', Label 'Voltar' (Cancel: true - returns to operation)

---

### Error Handling

| Scenario | Current Behavior | Spec Requirement | Alignment |
|---|---|---|---|
| **No vault, Ctrl+Q** | Immediate quit | Show confirmation dialog | ✗ Deviation |
| **Vault clean, Ctrl+Q** | Immediate quit | Show confirmation dialog | ✗ Deviation |
| **Vault dirty, 'Salvar'** | Dialog closes, no action | Execute Fluxo 5 save flow | ✗ Not implemented |
| **Vault dirty, 'Descartar'** | Dialog closes, quit | Discard and quit | ✓ Works (implicit) |
| **Vault dirty, 'Voltar'** | Return to normal | Return to operation | ✓ Implemented |
| **Vault dirty, save fails** | N/A (no save impl) | Show error, keep open | ✗ N/A |

---

## Summary Table: Specification vs. Implementation

### Fluxo 1 (Open Vault)

| Spec Step | Description | Impl Location | Status |
|---|---|---|---|
| 1 | Check unsaved changes | flow_open_vault.go:52-60 | Partial (dialog shown but handled by root) |
| 2 | Request vault path | flow_open_vault.go:73-77 (file picker) | Full |
| 3 | Request password | flow_open_vault.go:88-93 (password modal) | Full |
| 4 | Deserialize & validate | flow_open_vault.go:102-106 | Full |
| 5 | Load vault | flow_open_vault.go:137 (vaultOpenedMsg) | Full |
| Error: Wrong password | flow_open_vault.go:107-116 | Full (max 5 attempts) |
| Error: File issues | flow_open_vault.go:118-128 | Full (3 error types) |

---

### Fluxo 2 (Create Vault)

| Spec Step | Description | Impl Location | Status |
|---|---|---|---|
| 1 | Check unsaved changes | flow_create_vault.go:45-53 | Partial |
| 2 | Request save path | flow_create_vault.go:64-68 | Full |
| 2a | Check file exists | flow_create_vault.go:71-79 | Full |
| 3 | Request password (2x) | flow_create_vault.go:84 (password modal) | Full (in modal) |
| 4 | Evaluate password strength | passwordcreate.go:169-179 | Shown (modal) |
| 4a | Weak password decision | N/A | Missing |
| 5 | Create and save | flow_create_vault.go:94-98 | Full |
| 6 | Load new vault | flow_create_vault.go:106 (vaultOpenedMsg) | Full |

---

### Fluxo 3-5 (Exit Flows)

| Spec Flow | Scenario | Impl Location | Status |
|---|---|---|---|
| **Fluxo 3** | No vault loaded, user quits | root.go:256 | Partial (no confirmation dialog per spec) |
| **Fluxo 4** | Vault loaded, clean, user quits | root.go:256 (when !IsModified) | Full |
| **Fluxo 5-Step 1** | User requests exit | Implicit (Ctrl+Q) | Implicit |
| **Fluxo 5-Step 2** | Show unsaved decision | root.go:248-252 | Dialog shown |
| **Fluxo 5-Step 2a** | User chooses 'Salvar' | root.go:250 | No action bound (Cmd: nil) |
| **Fluxo 5-Step 2b** | User chooses 'Descartar' | Implicit after dialog | Implicit |
| **Fluxo 5-Step 2c** | User chooses 'Voltar' | root.go:252 (Esc) | Full |
| **Fluxo 5-Step 3** | Check external modification | N/A | Missing |
| **Fluxo 5-Step 4** | Execute atomic save | N/A | Missing |
| **Fluxo 5-Step 5** | Exit application | N/A (would follow save) | Incomplete |

---

## Key Deviations and Missing Implementations

### Fluxo 1 (Open Vault)
1. **Root model handling:** Ctrl+Q with unsaved vault changes is intercepted in root.go, not in flow
   - Spec expects flow to handle 'Salvar cofre -> continue to file picker'
   - Implementation: Dialog shown by root; flow receives Acknowledge message

2. **File validation errors:** Spec says 'volta ao passo 2' for magic/version errors
   - Implementation: Error shown and flow ends (endFlow)
   - Should: Return to file picker for correction

### Fluxo 2 (Create Vault)
1. **Password strength gate:** Spec step 4 defines weak password handling
   - Implementation: Strength meter shown in modal but flow does NOT offer 'prosseguir mesmo assim ou revisar' option
   - Passwordcreate modal has strength evaluation but no rejection/retry path in flow

2. **3-way overwrite decision:** Actually implemented as 2-way (Sobrescrever/Voltar)
   - Spec doesn't define a 3rd path for this step, so this is correct

### Fluxo 3-5 (Exit Flows)
1. **Fluxo 3 (no vault):** Spec expects confirmation dialog
   - Implementation: Direct quit without dialog
   - Reason: UX simplification (users expect quit when no data)

2. **Fluxo 4 (clean vault):** Spec expects confirmation dialog
   - Implementation: Direct quit without dialog
   - Reason: UX simplification (no unsaved data risk)

3. **Fluxo 5 (dirty vault):** Critical gaps
   - 'Salvar' action has no command bound (Cmd: nil)
   - No external file modification check (spec step 3)
   - No atomic save execution (spec step 4)
   - No backup protocol handling
   - **Impact:** User cannot save vault from quit dialog; must use separate save action

---

## File References Summary

| File | Lines | Purpose | Flows |
|---|---|---|---|
| flow_open_vault.go | 1-151 | State machine for opening vault | Fluxo 1 |
| flow_create_vault.go | 1-120 | State machine for creating vault | Fluxo 2 |
| root.go | 111-114, 243-257, 260-268 | Exit handling + action registration | Fluxo 3-5 |
| dialogs.go | 41-102 | Dialog factories (Message, Confirm, PasswordEntry) | Helper |
| decision.go | 1-476 | DecisionDialog implementation + factories | Fluxo 3-5 |
| passwordentry.go | 1-159 | Password entry modal + flow message type | Fluxo 1 |
| passwordcreate.go | 1-235 | Password creation modal with strength meter | Fluxo 2 |
| flows.go | 1-84 | Flow interfaces and message types | All |

---

## Conclusion

**Overall Implementation Status:**

- **Fluxo 1 (Open Vault):** ~85% complete
  - Core flow works; minor issues with error return paths and unsaved changes handling
  
- **Fluxo 2 (Create Vault):** ~75% complete
  - Missing password strength validation gate; otherwise solid
  
- **Fluxo 3-5 (Exit Flows):** ~40% complete
  - Dialog framework exists; save action not bound; external modification checks missing; atomic save not integrated

**Recommended Priorities:**
1. Fix Fluxo 5 'Salvar' action binding in root.go line 250
2. Implement external file modification check (Fluxo 5 step 3)
3. Implement password strength validation gate in Fluxo 2 (step 4)
4. Restore error recovery paths to file picker (Fluxo 1 file validation errors)
5. Add confirmation dialogs for Fluxo 3 & 4 (spec compliance, though current UX is arguably better)
