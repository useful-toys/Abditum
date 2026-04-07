# Phase 06.1 — Nothing Works as Expected

## Specification Compliance Analysis: Fluxos 1, 2, and 3-5

**Date:** 2026-04-06  
**Scope:** Detailed gap analysis comparing implementation against `@fluxos` specification  
**Status:** All three flows confronted; discrepancies documented with file/line references

---

## Executive Summary

Phase 06 (Welcome Screen + Vault Create/Open) passed all unit tests but implementation **deviates significantly from specification** in three critical flows:

1. **Fluxo 1** (Abrir Cofre) — Minor deviations in error handling paths
2. **Fluxo 2** (Criar Cofre) — **HIGH:** Missing password strength validation gate
3. **Fluxos 3-5** (Exit flows) — **CRITICAL:** Multiple core features missing; "Salvar e sair" is completely broken

**Total gaps identified:** 10 discrepancies across 3 flows  
**Severity breakdown:** 2 CRITICAL, 2 HIGH, 6 MEDIUM  
**Estimated remediation effort:** 350 lines of new code + fixes

---

## I. Fluxo 1 — Abrir Cofre Existente (Open Existing Vault)

### Specification Reference

`fluxos.md` lines 209-259: Fluxo 1 defines 5 steps for opening an existing vault:

1. Check unsaved changes in current vault
2. Request vault file path
3. Request master password (with retry on failure)
4. Deserialize JSON and validate model
5. Close current vault, load new vault, confirm success

### Implementation Reference

`internal/tui/flow_open_vault.go` lines 1-151: State machine with 5 states

---

### Gap 1.1: File Validation Errors Don't Return to File Picker

**Spec Requirement (Step 2):**
> "Se o `magic` for inválido ou o `versão_formato` for superior ao suportado → o sistema comunica o erro. O usuário pode corrigir o caminho e tentar novamente. Volta ao passo 2."

**Current Implementation:**

```go
// flow_open_vault.go:120-128
if errors.Is(err, storage.ErrInvalidMagic) {
    errMsg = "Arquivo não é um cofre válido."
} else if errors.Is(err, storage.ErrVersionTooNew) {
    errMsg = "Versão do cofre não é suportada por esta versão do Abditum."
} else if errors.Is(err, storage.ErrCorrupted) {
    errMsg = "Cofre corrompido e não pode ser recuperado."
}
f.messages.Show(MsgError, errMsg, 5, false)
return endFlow()  // ← WRONG: Exits flow instead of returning to file picker
```

**Problem:** Calls `endFlow()` which terminates the entire flow, instead of returning to `statePickFile` to let user select a different file.

**Impact:** User cannot retry with a different path; must restart the entire "Open Vault" flow.

**Spec vs. Implementation:**

| Aspect | Spec | Current | Compliance |
|--------|------|---------|-----------|
| Show error | Yes | Yes | ✓ |
| Action | "Volta ao passo 2" (retry with different file) | `endFlow()` (exit) | **✗ DEVIATION** |

**Remediation:**
- Replace `endFlow()` with `f.state = statePickFile; return pushModalMsg{modal: &filePickerModal{}}`
- Scope: ~3 lines of code

---

### Gap 1.2: Corrupted Vault Payload Doesn't Return to File Picker

**Spec Requirement (Step 4):**
> "Se o payload for corrompido, o JSON for inválido ou a Pasta Geral estiver ausente → o sistema comunica o erro genérico (categoria: integridade). Volta ao passo 2."

**Current Implementation (line 124-128):** Same as Gap 1.1 — calls `endFlow()` instead of returning to file picker.

**Impact:** User cannot retry with a different file; must restart entire flow.

**Remediation:** Same fix as Gap 1.1 (same code location).

---

## II. Fluxo 2 — Criar Novo Cofre (Create New Vault)

### Specification Reference

`fluxos.md` lines 263-336: Fluxo 2 defines 6 steps for creating a new vault:

1. Check unsaved changes in current vault
2. Request destination path and filename
3. Request master password (2x for confirmation)
4. **Evaluate password strength; if weak, offer choice to proceed or revise**
5. Create vault and save to disk (atomic)
6. Close current vault, load new vault, confirm success

### Implementation Reference

`internal/tui/flow_create_vault.go` lines 1-120: State machine with 4 states (missing strength check state)

---

### Gap 2.1: Missing Password Strength Validation Gate (HIGH SEVERITY)

**Spec Requirement (Step 4):**
> "O sistema avalia a força da senha.
> Se a senha for considerada fraca → o sistema comunica os critérios não atendidos e solicita uma decisão: prosseguir mesmo assim ou revisar a senha.
> - Se o usuário escolhe revisar → volta ao passo 3.
> - Se o usuário escolhe prosseguir → continua para o passo 5."

**Current Implementation:**

```go
// flow_create_vault.go:83-84
f.state = statePwdCreate
return func() tea.Msg {
    return pushModalMsg{modal: &passwordCreateModal{}}  // ← Shows strength meter but no decision point
}

// Then immediately at line 87:
case pwdCreatedMsg:  // ← No intermediate strength check state
    f.state = stateSaveNew
    // ... proceeds directly to save
```

**Problem:** 
1. `passwordCreateModal` shows password strength meter
2. But flow has **no** `stateStrengthCheck` to intercept weak passwords
3. No decision dialog offering "Prosseguir mesmo assim" or "Revisar" options
4. Flow skips directly from password entry to save

**Spec vs. Implementation:**

| Aspect | Spec | Current | Compliance |
|--------|------|---------|-----------|
| Show strength evaluation | Yes | Yes (modal shows meter) | ✓ |
| Communicate weakness criteria | Yes | No visible | ✗ |
| **Offer user choice for weak passwords** | **Yes** | **No** | **✗ CRITICAL** |
| "Revisar" → volta ao passo 3 | Yes | N/A (path missing) | ✗ |
| "Prosseguir" → passo 5 | Yes | N/A (path missing) | ✗ |

**Impact:** Specification requirement for user confirmation on weak passwords cannot be met. Flow always allows weak passwords without explicit user consent.

**Root Cause:** State machine incomplete; missing intermediate state for strength validation decision.

**Remediation:**
- Add new state: `stateStrengthCheck` (after `statePwdCreate`, value=20)
- Extend `passwordCreateModal` to emit new message type `pwdStrengthCheckMsg` when strength < threshold
- Add new state transition in `Update()` to handle strength check
- Show Decision modal offering "Prosseguir" or "Revisar"
- If "Revisar": return to `statePwdCreate`
- If "Prosseguir": continue to `stateSaveNew`
- Scope: ~60-70 lines across flow_create_vault.go and passwordcreate.go

**Code Location:**
- Primary: `internal/tui/flow_create_vault.go` lines 62-115 (Update method)
- Secondary: `internal/tui/passwordcreate.go` (modal modification)

---

### Gap 2.2: Overwrite Confirmation "Voltar" Path Unclear

**Spec Requirement (Step 2b):**
> "Se o arquivo de destino já existir → o sistema alerta que o arquivo já existe e solicita uma decisão: sobrescrever ou informar outro caminho.
> - Se o usuário escolhe informar outro caminho → volta ao passo 2."

**Current Implementation:**

```go
// flow_create_vault.go:71-79
if _, err := os.Stat(f.targetPath); err == nil {
    return func() tea.Msg {
        return Decision(SeverityDestructive, "Arquivo existe",
            "Um cofre já existe neste caminho. Deseja sobrescrever?",
            DecisionAction{Key: "Enter", Label: "Sobrescrever", Default: true},
            nil,  // ← No "outro caminho" action provided
            DecisionAction{Key: "Esc", Label: "Voltar"})
    }
}
```

**Problem:** 
1. Decision dialog has "Voltar" action but flow doesn't explicitly handle it
2. When user clicks ESC/"Voltar", modal is dismissed, likely emitting `flowCancelledMsg`
3. This calls `endFlow()` (line 110), terminating entire flow
4. But spec says "volta ao passo 2" (file picker), not exit flow

**Spec vs. Implementation:**

| Scenario | Spec | Current | Compliance |
|----------|------|---------|-----------|
| File exists, show dialog | Yes | Yes | ✓ |
| User confirms "Sobrescrever" | Yes | Continues to step 3 | ✓ |
| User chooses "outro caminho" | Return to step 2 (file picker) | `endFlow()` (exit) | **✗ DEVIATION** |

**Impact:** User cannot change destination path after learning a file exists; must restart entire "Create Vault" flow.

**Remediation:**
- Explicitly handle Decision dialog result for overwrite confirmation
- If "Sobrescrever": continue to `statePwdCreate`
- If "Voltar": return to `statePickFile` instead of calling `endFlow()`
- Scope: ~15-20 lines

**Code Location:** `internal/tui/flow_create_vault.go` lines 62-85 (Update method)

---

## III. Fluxos 3, 4, 5 — Exit Flows (Ctrl+Q)

### Specification Reference

`fluxos.md` lines 340-416: Three separate exit flows based on vault state

- **Fluxo 3** (lines 340-354): No vault loaded → request confirmation → exit
- **Fluxo 4** (lines 357-371): Vault loaded, clean → request confirmation → exit  
- **Fluxo 5** (lines 374-416): Vault loaded, dirty → 3 options (save/discard/cancel) + external file check + atomic save

### Implementation Reference

`internal/tui/root.go` lines 111-114, 243-257: Ctrl+Q handling in root model

---

### Gap 3.1: No Confirmation Dialog for Fluxo 3 (No Vault Loaded)

**Spec Requirement (Fluxo 3, Step 2):**
> "O usuário solicita sair.
> O sistema solicita confirmação.
> - Se o usuário confirma → a aplicação encerra.
> - Se o usuário volta → o fluxo é interrompido e nada muda."

**Current Implementation (Line 256):**
```go
return m, tea.Quit  // ← Direct quit, NO confirmation
```

**Condition:** Executes when `m.mgr == nil` (no vault loaded).

**Problem:** Ctrl+Q immediately exits without confirmation, violating spec.

**Impact:** Accidental exit possible with single keystroke.

**Spec vs. Implementation:**

| Step | Spec | Current | Compliance |
|------|------|---------|-----------|
| 1 | User requests exit | Ctrl+Q detected | ✓ |
| 2 | Request confirmation | Direct `tea.Quit` | **✗ MISSING** |

**Remediation:**
- Show Decision dialog: "Tem certeza que deseja sair do Abditum?"
- Add confirmation before Fluxo 3 exit
- Scope: ~15 lines

---

### Gap 3.2: No Confirmation Dialog for Fluxo 4 (Clean Vault)

**Spec Requirement (Fluxo 4, Step 2):**
> "O usuário solicita sair.
> O sistema solicita confirmação.
> - Se o usuário confirma → a aplicação encerra.
> - Se o usuário volta → o fluxo é interrompido e nada muda."

**Current Implementation (Line 256):**
```go
return m, tea.Quit  // ← Direct quit, NO confirmation
```

**Condition:** Executes when `m.mgr != nil && !m.mgr.IsModified()` (vault loaded, clean).

**Problem:** Same as Gap 3.1 — no confirmation shown.

**Impact:** Accidental exit with clean vault.

**Remediation:** Same as Gap 3.1 — add confirmation dialog.

---

### Gap 3.3: "Salvar" Action Has No Handler (CRITICAL)

**Spec Requirement (Fluxo 5, Step 2-3):**
> "O usuário solicita sair.
> O sistema comunica que há alterações não salvas e solicita uma decisão: salvar e sair, descartar e sair, ou voltar.
> - Se o usuário escolhe salvar e sair → continua para o passo 3.
> [Step 3 involves checking external file modification and atomic save]"

**Current Implementation (Lines 248-253):**
```go
return m, func() tea.Msg {
    return Decision(SeverityNeutral, "Alterações não salvas",
        "Deseja salvar as alterações antes de sair?",
        DecisionAction{Key: "Enter", Label: "Salvar", Default: true},  // ← Default action
        []DecisionAction{{Key: "D", Label: "Descartar"}},
        DecisionAction{Key: "Esc", Label: "Voltar"})
}
```

**Problem:** 
1. Decision dialog is shown ✓
2. User presses Enter to "Salvar"
3. Modal is dismissed
4. **But there is NO handler for this Decision result** — the action has no command bound
5. Application returns to normal state without saving

**Spec vs. Implementation:**

| Step | Spec | Current | Status |
|------|------|---------|--------|
| 2 | Show 3 options | Dialog shown | ✓ |
| "Descartar" | Exit without saving | Implicit quit works | ✓ |
| "Voltar" | Cancel exit | Dialog dismiss works | ✓ |
| **"Salvar"** | **Start save flow** | **Dialog dismissed; nothing happens** | **✗ CRITICAL** |

**Root Cause:** The Decision dialog is instantiated and shown, but the result is not handled. No command is executed when user selects "Salvar".

**Impact:** **CRITICAL** — User cannot save and exit. The "Salvar e sair" path from Spec is completely non-functional.

**Remediation:**
- Create new flow type: `saveAndExitFlow` OR integrate save logic into root model's Decision handling
- When "Salvar" chosen, start flow that:
  - Executes Fluxo 8 (Save cofre no arquivo atual) logic
  - Checks for external file modification (Spec Step 3)
  - Executes atomic save with backup protocol
  - Handles error recovery
  - Exits only on success
- Scope: ~100-120 lines of new code

**Code Location:**
- `internal/tui/root.go` lines 248-253 (Decision handling)
- New file: `internal/tui/flow_save_and_exit.go` (proposed)

---

### Gap 3.4: No External File Modification Check Before Save

**Spec Requirement (Fluxo 5, Step 3):**
> "O sistema verifica se o arquivo foi modificado externamente desde a última leitura ou salvamento.
> - Se foi modificado externamente → o sistema comunica o conflito e solicita uma decisão: sobrescrever e sair, ou voltar.
> - Se não foi modificado externamente → continua para o passo 4."

**Current Implementation:** Not implemented at all.

**Problem:** 
1. No mechanism to detect external file modification
2. No modal to present conflict to user
3. No decision point for user to choose overwrite vs. cancel

**Impact:** Can silently overwrite externally-modified vault files without user knowledge.

**Remediation:**
- Implement file modification detection (compare disk metadata against last-read state)
- Create new modal: `ExternalModificationConflictModal`
- Add state in saveAndExitFlow for this check
- Scope: ~40-50 lines of code

---

### Gap 3.5: No Atomic Save Execution During Exit

**Spec Requirement (Fluxo 5, Step 4):**
> "O cofre é gravado no arquivo atual usando o salvamento atômico (sinalização). Segredos marcados para exclusão não são gravados.
> - Se o salvamento falhar sem ter gerado backup → o sistema comunica o erro. O cofre permanece carregado e `alterado`.
> - Se o salvamento falhar após ter gerado backup → o sistema comunica o erro e informa que existe um backup disponível para intervenção manual. O cofre permanece carregado e `alterado`."

**Current Implementation:** Not integrated into exit flow.

**Problem:** No save logic executed when user selects "Salvar" in Fluxo 5.

**Impact:** Cannot persist data before exit.

**Remediation:**
- Integrate Fluxo 8 (Save cofre no arquivo atual) logic into saveAndExitFlow
- Execute `storage.Save()` with atomic protocol
- Handle both success and failure cases
- For failures: communicate error + backup availability
- Only exit on success
- Scope: ~50-70 lines

---

### Gap 3.6: No Backup Protocol Communication on Save Failure

**Spec Requirement (Fluxo 5, Step 4):**
> "Se o salvamento falhar após ter gerado backup → o sistema comunica o erro e informa que existe um backup disponível para intervenção manual."

**Current Implementation:** Not implemented.

**Problem:** If save fails during exit and a backup exists, user is not informed.

**Impact:** User may not know backup file is available for manual recovery.

**Remediation:**
- In error handling path of saveAndExitFlow, check if backup was created
- If backup exists: show error message mentioning `.abditum.bak` availability
- Scope: ~15-20 lines

---

## IV. Summary Table: All Gaps

| ID | Flow | Gap | Severity | Scope | Code Location |
|:--:|:----:|-----|----------|-------|----------------|
| 1.1 | Fluxo 1 | File validation errors don't return to picker | MEDIUM | ~3 lines | flow_open_vault.go:120-128 |
| 1.2 | Fluxo 1 | Corrupted payload doesn't return to picker | MEDIUM | ~3 lines | flow_open_vault.go:124-128 |
| 2.1 | Fluxo 2 | **Missing password strength validation gate** | **HIGH** | **~60 lines** | flow_create_vault.go + passwordcreate.go |
| 2.2 | Fluxo 2 | Overwrite "Voltar" exits flow instead of returning to picker | MEDIUM | ~15 lines | flow_create_vault.go:71-85 |
| 3.1 | Fluxo 3 | No confirmation dialog when exiting with no vault | MEDIUM | ~15 lines | root.go:243-257 |
| 3.2 | Fluxo 4 | No confirmation dialog when exiting with clean vault | MEDIUM | ~15 lines | root.go:243-257 |
| **3.3** | **Fluxo 5** | **"Salvar" action has no handler** | **CRITICAL** | **~110 lines** | **root.go:248-253 + new saveAndExitFlow** |
| 3.4 | Fluxo 5 | No external file modification check | HIGH | ~45 lines | (New in saveAndExitFlow) |
| 3.5 | Fluxo 5 | No atomic save execution | CRITICAL | ~60 lines | (New in saveAndExitFlow) |
| 3.6 | Fluxo 5 | No backup protocol communication | MEDIUM | ~15 lines | (New in saveAndExitFlow) |

---

## V. Estimated Remediation Effort

### By Severity

| Severity | Count | Total Scope |
|----------|-------|-------------|
| CRITICAL | 2 | ~170 lines |
| HIGH | 2 | ~105 lines |
| MEDIUM | 6 | ~80 lines |
| **TOTAL** | **10** | **~355 lines** |

### By Flow

| Flow | Gaps | Complexity | Estimated Work |
|------|------|-----------|-----------------|
| Fluxo 1 | 2 | Simple path returns | ~6 lines |
| Fluxo 2 | 3 | State machine + modal | ~75 lines |
| Fluxos 3-5 | 5 | New flow + modals | ~280 lines |
| **TOTAL** | **10** | — | **~360 lines** |

---

## VI. Recommended Prioritization

### Tier 1: CRITICAL (Must Fix Before Vault Works)
- **Gap 3.3:** Implement "Salvar e sair" handler — **blocks save-and-exit completely**
- **Gap 3.5:** Atomic save execution during exit — **blocks save-and-exit completely**
- **Gap 3.4:** External file modification check — **required for data safety**

→ Requires new `saveAndExitFlow` component (~170 lines)

### Tier 2: HIGH (Spec Compliance)
- **Gap 2.1:** Password strength validation gate — **spec requires user choice for weak passwords**
- Confirmations for Fluxos 3 & 4 — **spec requires confirmation for all exit scenarios**

→ Requires state machine updates + decision dialog

### Tier 3: MEDIUM (UX Polish)
- Gap 1.1 & 1.2: Return to file picker on validation errors
- Gap 2.2: Fix overwrite "Voltar" path
- Gap 3.6: Backup communication

---

## VII. Investigation Summary

### Files Reviewed
- ✓ `fluxos.md` — Flow specifications (Fluxos 1-5)
- ✓ `internal/tui/root.go` — Ctrl+Q handling (lines 243-257)
- ✓ `internal/tui/flow_open_vault.go` — Fluxo 1 implementation
- ✓ `internal/tui/flow_create_vault.go` — Fluxo 2 implementation
- ⊙ `internal/tui/passwordcreate.go` — Strength modal (needs extension for Gap 2.1)
- ⊙ `internal/tui/decision.go` — Decision modal (result handling missing)

### Key Discoveries

1. **Tests pass but flows don't work** — 141/141 tests passing but critical decision paths untested
2. **Architecture gap** — Decision modal result not wired to root model handler
3. **Fluxo 5 is a stub** — "Salvar e sair" shows dialog but has no implementation
4. **No save-and-exit flow exists** — Entire save-during-exit logic missing from codebase

---

## VIII. Keybindings Analysis (Extended Scope)

As part of Phase 06.1 specification compliance review, keybindings have been analyzed against `@tui-specification-novo.md`. See `KEYBINDINGS-SPEC.md` for full analysis; this section summarizes findings.

### Current Keybindings vs. Specification

**Welcome Screen (Phase 06 Scope):**

| Spec Requirement | Current Implementation | Status | Scope |
|---|---|---|---|
| F5 → Criar Novo Cofre (Fluxo 2) | N key (not F5) | ⚠️ WRONG KEY | Phase 06.1 (easy fix) |
| F6 → Abrir Cofre Existente (Fluxo 1) | O key (not F6) | ⚠️ WRONG KEY | Phase 06.1 (easy fix) |
| F1 → Help | F1 ✓ | ✅ SPEC-COMPLIANT | Phase 06 |
| Ctrl+Q → Exit | Ctrl+Q ✓ | ⚠️ PARTIAL (handler broken for Fluxo 5) | Phase 06.1 (fix via saveAndExitFlow) |
| F12 → Toggle Theme | F12 ✓ | ❌ NOT IN SPEC | Decision needed |

**Workspace Tabs (Phase 06.5+ Scope):**

| Spec Requirement | Current Implementation | Status | Scope |
|---|---|---|---|
| F2 → Modo Cofre | *Not implemented* | ❌ NOT IMPLEMENTED | Phase 06.5+ |
| F3 → Modo Modelos | *Not implemented* | ❌ NOT IMPLEMENTED | Phase 06.5+ |
| F4 → Modo Configurações | *Not implemented* | ❌ NOT IMPLEMENTED | Phase 06.5+ |

**Save/Export/Import (Phase 7+ Scope):**

| Spec Requirement | Current Implementation | Status | Scope |
|---|---|---|---|
| F7 → Salvar Atual | *Not implemented* | ❌ NOT IMPLEMENTED | Phase 7+ (Fluxo 8) |
| Shift+F7 → Salvar Como | *Not implemented* | ❌ NOT IMPLEMENTED | Phase 7+ (Fluxo 9) |
| Ctrl+F7 → Alterar Senha | *Not implemented* | ❌ NOT IMPLEMENTED | Phase 7+ (Fluxo 11) |
| F9 → Exportar | *Not implemented* | ❌ NOT IMPLEMENTED | Phase 7+ (Fluxo 12) |
| Shift+F9 → Importar | *Not implemented* | ❌ NOT IMPLEMENTED | Phase 7+ (Fluxo 13) |
| Shift+F6 → Recarregar | *Not implemented* | ❌ NOT IMPLEMENTED | Phase 7+ (Fluxo 10) |
| Ctrl+Alt+Shift+Q → Lock Vault | *Not implemented* | ❌ NOT IMPLEMENTED | Phase 6.5+ or 7+ (Fluxo 6) |

### Keybindings Gaps Summary

| ID | Gap | Severity | Impact | Remediation |
|:--:|-----|----------|--------|-------------|
| K.1 | O/N should be F5/F6 | MEDIUM | Users can't find expected F-keys | Change Keys arrays in root.go |
| K.2 | F12 (Toggle) not in spec | LOW | Feature exists but not documented | Decision: keep/remove/defer |
| K.3-K.6 | F2-F4, F7-F9, variants not implemented | LOW | No blocker; depends on future phases | Document in PLAN.md |

### Remediation for Phase 06.1

**Include in Phase 06.1 PLAN.md:**
- Task K.1: Migrate O→F6, N→F5 (~10 minutes)
- Task K.2: Document F12 decision (~5 minutes)
- Task K.3: Document deferred keybindings (~5 minutes)

**Total effort:** 1-1.5 hours (non-blocking, can be parallelized)

### Keybindings Architecture Notes

- Action registration in `root.go:109-142` uses `Keys: []string{...}` array
- Dispatch in `actions.go:92-120` matches all keys in the array
- Help modal groups actions by Group ID
- Command bar shows `Keys[0]` for display; all keys trigger the action
- F1 is right-anchored; other actions left-padded
- Spec defines semantic F-key groupings: F1/F12 (global), F2-F4 (tabs), F5-F6 (entry), F7-F9 (persist)

---

## Conclusion

Phase 06 implementation has **significant specification deviations** that prevent core workflows from functioning as documented. The most critical issue is **Fluxo 5 "Salvar e sair" being completely non-functional**. Phase 06.1 must address all CRITICAL and HIGH severity gaps before vault operations can be considered compliant with specification. Keybindings alignment (Phase 06.1 extension) improves UX and spec compliance without blocking core functionality.
