# Phase 6: Welcome Screen + Vault Create/Open - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the context gathered and alternatives considered.

**Date:** 2026-04-02
**Mode:** Specification-driven (user provided spec files instead of interactive Q&A)

---

## Context Sources

The user provided four specification documents as the basis for all decisions in this phase, bypassing the interactive gray-area discussion flow:

1. `tui-specification-novo.md` — Full visual spec for PasswordEntry, PasswordCreate, FilePicker (Open + Save modes), Help, and decision dialogs
2. `tui-design-system-novo.md` — Design system with color palette, typography, borders, layout rules, and modal anatomy
3. `fluxos.md` — Behavioral flow specs for Fluxo 1 (Abrir Cofre) and Fluxo 2 (Criar Cofre), including CLI fast-path, retry rules, and overwrite handling
4. `requisitos.md` — Requirements including VAULT-01, VAULT-03, VAULT-04, and password strength rules

---

## Key Decisions Derived from Specs

### Welcome screen action model

**Considered:** ROADMAP Phase 6 plan described a `welcomeModel` with j/k/Enter menu and sub-states (`subStatePickPath`, `subStateCreatePassword`, `subStateOpenPassword`).

**Resolved:** Phase 5 CONTEXT (D-02, D-09) is authoritative. `preVaultModel` stays display-only; open/create are FlowRegistry flows dispatched by keys `n` and `o` already registered in Phase 5. The ROADMAP plan description was written before Phase 5 locked the architecture. No menu, no sub-states.

### Max retries on vault open

**Spec source:** `tui-specification-novo.md` §PasswordEntry — "Tentativa 2 de 5" wireframe.

**Decision:** Max 5 attempts. Counter hidden on attempt 1, shown from attempt 2 onward. After 5 exhausted attempts, dialog closes automatically (emits `flowCancelledMsg{}`). No hard lock-out beyond that — user returns to welcome screen.

### Weak password UX

**Spec source:** `tui-specification-novo.md` §PasswordCreate — action default unlocked when both fields non-empty, regardless of strength.

**Resolved:** Strength meter (`Força: ████████░░ Boa/Fraca`) with semantic colors is the full weak-password UX. No separate confirmation modal for weak password. Submit is always allowed when both fields are non-empty. This aligns with the ROADMAP plan's "non-blocking warning banner, submit allowed" and with `requisitos.md` §Força da Senha Mestra: "O aviso de senha fraca é apenas informativo."

### FilePicker scope

**Spec source:** `tui-specification-novo.md` §FilePicker — full two-panel spec (tree + files, Open + Save modes).

**Decision:** Implement the full spec-compliant FilePicker as a dedicated `filePickerModal` struct. No simplified text-input shortcut. Phase 6 is the first use case for FilePicker per Phase 5 CONTEXT (D-10: "File picker modal is deferred — implemented in the phase that introduces its first use case").

### CLI path fast-path

**Spec source:** `fluxos.md` Fluxo 1 — "Entrada antecipada via argumento de linha de comando."

**Decision:** When `initialPath` is non-empty at startup, `rootModel.Init()` returns a cmd that starts `openVaultFlow` starting from path verification (not FilePicker). If path is invalid, show error modal and fall back to FilePicker.

### Error classification

**Spec source:** `fluxos.md` Fluxo 1 step 1/3, `requisitos.md` VAULT-04.

**Decision:** Auth errors → retry (up to 5); integrity/magic/version errors → Recognition × Error modal → return to FilePicker. Mapped to storage sentinel errors with no technical strings exposed.

---

## Alternatives Not Pursued

- **Text-input path entry** instead of FilePicker: rejected — full FilePicker is specced and Phase 6 is its designated introduction phase
- **Blocking weak-password modal**: rejected — spec and requirements both say informative/non-blocking
- **`welcomeModel` with sub-states**: rejected — superseded by Phase 5 architecture (D-09)
- **Single `passwordModal` for both entry and create**: rejected — the two dialogs have different field counts, Tab behavior, and strength meter; separate structs are cleaner

