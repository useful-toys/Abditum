# Phase 5: TUI Scaffold + Root Model - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-31
**Phase:** 05 — TUI Scaffold + Root Model

---

## Area: go.mod Dependency Scope

**Q:** Which charm.land packages should be added to go.mod in Phase 5?
**Options presented:** bubbletea/v2 only · bubbletea/v2 + lipgloss/v2 · All TUI deps now
**Selected:** Add all TUI deps now (bubbletea/v2 + bubbles/v2 + lipgloss/v2 + teatest/v2)

**Q:** Version pinning strategy?
**Selected:** Pin exact latest versions (`go get @latest` for each)

---

## Area: File Layout in internal/tui/

**Q:** How should files be organized inside internal/tui/?
**Options presented:** Separate files per concern · One file (split later) · One file per model
**Selected:** Separate files per concern (root.go, modal.go, state.go)

**Q:** Should the three child model stub types live in their own files already in Phase 5?
**Options presented:** Yes — stub files created now · No — stubs in root.go, split when they grow
**Selected:** Yes — stub files (welcome.go, vault.go, detail.go) created now

---

## Area: Child Model Interface

**Q:** Should child models share a common Go interface?
**Options presented:** Yes — private childModel interface · No interface — concrete structs + switch · Pre-interface pattern
**Selected:** Yes — define a private childModel interface

**Q:** What methods does the childModel interface include?
**Options presented:** Standard tea.Model only · tea.Model + SetSize · tea.Model + SetManager
**Selected:** tea.Model methods + SetSize(w, h int)

*[Later refined — see "Child interface does not implement tea.Model" below]*

---

## Area: Initial Timer Values

**Q:** How should lockTimer / clipboardTimer be initialized before a vault is open?
**Options presented:** Zero until vault opened · Hardcoded defaults · Package-level constants
**Selected:** Zero / empty until vault is opened (read from Configuracoes on transition)

**Q:** Should the global tick start immediately or only after vault is opened?
**Options presented:** Tick always runs · Tick only starts after vault opened
**Selected:** Tick only starts after vault is opened

---

## Area: Additional — Placeholder Screen Content

**Q:** What does the placeholder screen render when the app launches?
**Options presented:** Delegate to welcomeModel.View() · Hardcoded string · Centered lipgloss box
**Selected:** Delegate to welcomeModel.View() — stub renders ASCII art + quit hint

---

## Area: Additional — State Machine (from "stateLocked" discussion)

**User clarification:** "Não haverá uma tela de lock. A aplicação volta para a tela de abrir cofre."
→ No `stateLocked` state. Lock = wipe + go to `stateOpenVault`.

**User clarification:** "não retorna ao stateWelcome, mas ao stateOpenVault (ou de nome similar)"
→ Lock returns to a dedicated open-vault state, not to the welcome/path-selection state.

**User clarification (state machine):** "imaginei um state de welcome, um de abrir cofre, um de criar novo cofre, mostrando cofre"
→ 4 states: stateWelcome / stateOpenVault / stateCreateVault / stateVaultOpen

---

## Area: Additional — Vault Layout

**User clarification:** "split view, side by side, mas cada componente terá seu model independente. Além disso, teremos cabeçalho, rodapé e uma barra que mostra comando."
→ Constant frame: header + message bar + work area + command bar + footer (further zones possible).

**User clarification (pre-vault screens):** "usa o mesmo layout que a apresentação do cofre, mas sem mostrar a árvore e os detalhes. Ao invés disso, mostra o ascii art"
→ Same frame always; ASCII art fills work area when vault not open.

---

## Area: Additional — Terminal Resize Handling

**Q:** How does rootModel handle tea.WindowSizeMsg?
**User clarification:** "o filho precisa saber seu tamanho dentro da tela, não só o tamanho do terminal"
→ rootModel computes per-child allocated size, calls SetSize(w, h) on each child.

**Selected:** rootModel computes per-child size, calls SetSize on each

---

## Area: Additional — Quit Shortcut

**User clarification:** "ctrl + Q em qualquer momento da aplicação. não deve usar ctrl + c"
→ Global quit is ctrl+Q. Not ctrl+C, not q.

**User clarification:** "mostra confirmação conforme especificado nos requisitos (veja fluxos.md)"
→ Confirmation modal per fluxos.md when unsaved changes present.

---

## Area: Additional — Inter-Model Communication Rules

**User specification:**
- A model must not access another model's memory/fields.
- A model must not call another model's methods to notify it.
- Mutations go through Manager; Manager triggers Update messages on rootModel to broadcast changes.

**Conflict flagged by agent:** nil interface trap in Go — storing a nil concrete pointer in an interface variable is NOT nil. Resolution: rootModel stores children as concrete pointer types, interface used only transiently in liveModels() helper.

---

## Area: Additional — Model Lifecycle

**User specification:**
- "Live" = currently on screen. Dead models = nil, no memory held, no references.
- Modal models: allocated on open, set to nil on dismiss.
- Pre-vault models (welcome, openVault, createVault): zeroed when vault opens.

---

## Area: Additional — Modal Stack

**User specification:** "devemos tratar as janelas como uma stack, podendo uma modal abrir outra modal"
→ Changed from single `modal *modalModel` to `modals []*modalModel` LIFO stack.

---

## Area: Additional — Child Interface Does Not Implement tea.Model

**User specification:** "Os modelos filho implementam uma interface própria e não o Model do bubbletea"
→ Custom `childModel` interface. `Update` returns only `tea.Cmd` (no self-replacement). `View()` returns `string`.

**Q:** Does View() return string or tea.View?
**User:** "suponho que string, certo?" → Confirmed: string.

**Q:** Does child interface include Init() tea.Cmd?
**User:** "não sei, precisaria pesquisar" → Left to researcher/planner.

---

## Area: Additional — Work Area and Scroll

**User specification:**
- Work area zones: header + message bar (contextual hints) + work area + command bar (global commands) + footer.
- Frame structure not fully locked — more zones possible in later phases.
- Template editor fills work area: left = template list, right = template view/edit (same pattern as tree+detail).
- All list/tree/detail components must support vertical scroll.

---

## Area: Additional — v2 Constraint

**User specification:** "Usaremos obrigatoriamente as versões v2 do charm.land"
→ All charm.land packages must be v2. No v1 packages permitted anywhere.

---

## Area: Additional — ASCII Art Logo

**User provided reference implementation** from `c:\git\Abditum2\internal\tui\ascii.go`:
- `AsciiArt` constant — 5-line "Abditum" wordmark
- `RenderLogo()` — violet→cyan gradient, one lipgloss color per line
- Palette: `["#9d7cd8", "#89ddff", "#7aa2f7", "#7dcfff", "#bb9af7"]`
