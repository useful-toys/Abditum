# Phase 3: Vault Domain + Manager - Discussion Log

**Date:** 2025-03-29
**Participants:** User, OpenCode Assistant
**Purpose:** Context discussion to capture implementation decisions for Phase 3

---

## Session Overview

This session conducted a thorough Phase 3 context discussion following the `discuss-phase.md` workflow. The goal was to capture all implementation decisions needed for the vault domain layer and Manager pattern, so downstream agents (researcher, planner) can work autonomously without needing to ask the user again.

**Total decisions captured:** 30 (D-01 through D-30)
**Key documents provided by user:** `arquitetura-camada-dominio.md` (architectural decisions already made)
**Gray areas explored:** 4 major areas with 6 follow-up questions

---

## Major Decision Areas

### Area 1: Identity and Reference Management

**Question:** How should entities be identified? Synthetic IDs (NanoID/UUID) or Go pointers?

**Options discussed:**
- **A**: Generate NanoIDs for all entities (persistent identity across sessions)
- **B**: Use Go pointers during session, JSON hierarchical structure for persistence
- **C**: Hybrid approach (IDs for some entities, pointers for others)

**Decision:** **Option B** (D-01, D-02, D-03)
- Go pointers are sufficient identifiers during session (globally unique, O(1) lookup)
- JSON hierarchical structure (nested objects/arrays) expresses relationships without synthetic IDs
- Uniqueness validated by composite key `(parent, nome)` at mutation time
- Eliminates ID generation, collision detection, storage overhead

**Rationale:** IDs exist to reference entities across boundaries that erase structure (network, database, serialization between processes). Abditum doesn't cross these boundaries during session - vault lives in memory as graph of Go pointers.

**User confirmation:** Explicit agreement that no synthetic IDs are needed.

---

### Area 2: Invariant Validation Strategy

**Question:** How to enforce invariants - make impossible by design vs validate explicitly?

**User provided key insight:** Two-tier classification
- **Category A** (structurally impossible): Observação immutability via encapsulation - operations only manipulate user fields slice
- **Category B** (explicit validation): Uniqueness, Pasta Geral protection, cycle detection

**Decision:** D-06 (Make impossible > validate)
- Design data structures to make invalid states unrepresentable
- Only validate what cannot be prevented by structure

**Key example:** Observação always last field in internal slice, but `Campos()` returns only user fields - manipulation methods never see Observação in their index space.

**Follow-up:** Error representation strategy
- **Question:** Sentinel errors vs custom error structs?
- **Decision:** D-07 (Hybrid approach)
  - Sentinels for simple validations (type check only)
  - Custom structs when TUI needs structured data from error

---

### Area 3: Manager.Vault() Snapshot Strategy

**Initial question:** How should Manager expose vault state to TUI without allowing accidental mutation?

**Options discussed:**
- **A**: Deep copy everything (safe but expensive)
- **B**: Live pointer with documented read-only contract (zero cost but relies on discipline)
- **C**: Read-only wrapper with interface (type-safe but complex boilerplate)
- **D**: Shallow copy with copy-on-write (partial protection)
- **E**: Snapshot struct optimized for read (safe, on-demand decryption)

**User provided:** `arquitetura-camada-dominio.md` document clarifying architectural decision already made

**Resolution:** **Option B with package encapsulation** (D-08, D-09, D-10)
- Manager.Vault() returns live `*Cofre` pointer
- Safety via Go package-level encapsulation (fields lowercase/private)
- Getters return defensive copies of mutable collections
- Secret values via explicit `ValorComoString()` method

**Key clarification:** ROADMAP pitfall "Manager.Vault() must return snapshot" satisfied by defensive copies in **getters**, not by copying entire Cofre structure.

---

### Area 4: Session State Tracking Mechanics

**Questions explored:**
1. When does `original` → `modificado` transition fire?
2. How to detect deep mutations (field value changes)?
3. Where to filter deleted secrets?
4. Do deleted secrets persist across save?

**User provided:** `arquitetura-camada-dominio.md` section 7 with complete session state rules

**Decisions:** D-11 through D-17

**D-11 - Key clarification on favoriting:**
- Initial spec (modelo-dominio.md line 73): "Favoriting changes estadoSessao"
- User correction: Favoriting does NOT change `estadoSessao`
- Rationale: Favoriting is metadata/navigation preference, not content change
- Two independent flags: `cofre.modificado` (any mutation) vs `segredo.estadoSessao` (content change)

**D-12 - Change detection:**
- Flags only update on actual value difference, not on method call
- Entity methods return `(alterado bool, err error)` to signal real change
- Example: User opens rename dialog, doesn't change, confirms → no modification indicator

**D-17 - Atomic save:**
- **User requirement:** "Ao salvar, precisaremos garantir que, os segredos só serão excluidos do cofre se o salvar for realizado com sucesso."
- **Solution:** Two-phase commit pattern
  - Phase 1: Prepare snapshot (filter excluido, deep copy)
  - Phase 2: Persist snapshot
  - Phase 3: Commit deletions from live vault only if Phase 2 succeeds
- Guarantees: Save failure doesn't cause data loss in memory

---

## Follow-up Questions Resolved

### Q1: Reordenação - Manager methods or TUI concern?

**Decision:** D-23 (Manager exposes explicit repositioning operations)
- Manager has `ReposicionarSegredo()`, `ReposicionarPasta()`
- Helper methods: `SubirXNaPosicao()`, `DescerXNaPosicao()`
- TUI maps: `Ctrl+Up` → Subir, `Ctrl+Down` → Descer

**User refinement:** "Reposicionar" not "Reordenar"
- "Reordenar" suggests applying criterion (alphabetic, chronologic)
- "Reposicionar" indicates moving to specific position
- More semantically precise

### Q2: Criar segredo a partir de modelo - Who builds fields?

**Decision:** D-22 (Manager has specific method)
- `Manager.CriarSegredoDeModelo(pasta, nome, modelo, posicao)`
- Manager delegates to `Pasta.criarSegredoDeModelo()` which materializes template
- Logic belongs in domain, not TUI

**User clarification:** Manager delegates to Pasta factory method which receives modelo as parameter

### Q3: Favoritos - Dedicated query or TUI filters locally?

**Decision:** D-21 (Manager exposes ListarFavoritos())
- Consistent with `BuscarSegredos()` pattern (D-15)
- DFS traversal (not BFS - corrected from initial assumption based on modelo-dominio.md)
- Filters excluido internally

**Follow-up question:** Should Manager cache favoritos?
- **User insight:** "Tendo um manager com boa API, poderemos introduzir a cache mais tarde sem impacto em código existente"
- **Decision:** No cache in Phase 3 - simple on-demand traversal
- Future optimization transparent to TUI (no API changes needed)

### Q4: Configurações - Phase 3 scope or deferred?

**Decision:** D-20 (Configurações in Phase 3)
- Three timer fields from modelo-dominio.md
- Manager has `AlterarConfiguracoes()` method
- All timers mandatory (cannot be disabled)

### Q5: Timestamps - Does Pasta have them?

**Decision:** D-19 (Only Segredo and Cofre)
- Pasta has NO timestamps
- No need to audit hierarchy changes via folders
- Simpler model

### Q6: NanoID pitfall - Still relevant?

**Decision:** D-18 (Obsolete)
- ROADMAP warning "NanoID must use crypto/rand" obsolete
- D-01 eliminated need for synthetic IDs entirely
- No ID generation in Phase 3

---

## Architectural Clarifications

### Manager Responsibility ("Lógica de Negócio")

**User question:** "O que você quer dizer com Manager não contém lógica de negócio?"

**Critical refinement of D-04 and D-25:**

**Manager HAS (high-level business knowledge):**
- ✅ Knows WHAT operations exist and their semantics
- ✅ Knows workflows (order of steps in complex operations)
- ✅ Knows global rules ("any mutation marks cofre.modificado")
- ✅ Knows relationships (Segredo belongs to Pasta, operations cross boundaries)

**Manager DOES NOT HAVE (implementation logic):**
- ❌ How to verify name is unique (algorithm)
- ❌ How to insert at position in slice (manipulation)
- ❌ How to detect cycles (validation)
- ❌ Structure of internal fields (`pasta.segredos`)

**Analogy refined:** Manager is maestro (knows the music - WHAT to play, when each instrument enters), entities are musicians (know HOW to play their instruments).

**Key insight:** Manager has **business knowledge of high level** (workflows, semantics) but not **implementation logic** (algorithms, structures).

---

### Entity Validation Responsibility

**User question on D-02:** "Entendo que os métodos de Pasta poderão responder se já existe um filho com um nome informado. Manager será responsável por orquestrar a validação, mas recorrerá a métodos de validação ou teste das entidades. Estou correto?"

**Confirmation:** Absolutely correct.
- Entity owns query methods: `contemSegredoComNome()`, `segredoComNome()`, etc.
- Entity uses these in its own validation (factory/mutation methods)
- Manager orchestrates by calling entity methods, doesn't implement verification logic

**D-02 clarified:** Unicidade validada via métodos de consulta das entidades.

---

### Separation of Responsibilities: Factory vs Initializer

**User requirement:** "Eu prefiro uma factory que cria um cofre vazio e um método de inicialização - entendo que são dois propósitos e responsabilidades separadas"

**Decision D-28 refined:**

**D-28a: NovoCofre()** - Factory (construction)
- Responsibility: Build minimal aggregate structure
- Returns vault with Pasta Geral + default configs
- NO initial content

**D-28b: Cofre.InicializarConteudoPadrao()** - Bootstrap (initial content)
- Responsibility: Populate new vault with canonical structure
- Creates default folders: "Sites e Apps", "Financeiro"
- Creates default templates: "Login", "Cartão de Crédito", "Chave de API"
- Per requisitos.md functional requirement
- NOT a Manager operation (system bootstrap, not user operation)

**Separation principle:** Construction (factory) vs Bootstrap (domain service) are distinct responsibilities.

---

## Timestamp Policy Clarification

**User specification:** "Timestamps devem ser atualizados somente se houver alteração estrutural nas pastas, alteração de fato no segredo. Favoritar não atualiza timestamp do segredo, apenas do cofre."

**Decision D-24:**
- **Segredo timestamp**: Only on content/structural changes (rename, field edits, move to different pasta)
- **Cofre timestamp**: On any mutation including metadata (favoriting, repositioning)
- **Favoriting**: Updates Cofre timestamp, NOT Segredo timestamp
- **Repositioning**: Updates Cofre timestamp, NOT entity timestamps

**Rationale:**
- Favoriting is user preference (metadata), not content change
- Repositioning is container property, not entity content
- Move changes parent reference (structural to entity)

---

## Deletion Semantics

**User confirmation:** "Pastas não tem soft delete."

**Decision D-27:**
- **Segredo**: Soft delete (`estadoSessao = excluido`), restored via flag unmark
- **Pasta**: Hard delete (immediate removal, promote children, auto-rename on conflict)

**Rationale:**
- Segredo is leaf → undo trivial (flip flag)
- Pasta has dependents → deletion is complex restructuring
- Undo pasta deletion would require reversing multiple promotions + renaming (infeasible)

---

## Interface Dependency Analysis

**Question:** Does Phase 3 need to define `RepositorioCofre` interface, or defer to Phase 4?

**Analysis:**
- **Option A**: Define minimal interface now (`Salvar(cofre) error`)
- **Option B**: Defer entirely to Phase 4
- **Option C**: Full interface with `Carregar()` + `Salvar()`

**Decision:** **Option A** (minimal interface in Phase 3)
- Allows Phase 3 testing with mock repository
- Manager.Salvar() can be fully implemented
- Desocoupling: vault doesn't depend on storage implementation
- Interface is stable, unlikely to change

**Interface defined:**
```go
type RepositorioCofre interface {
    Salvar(cofre *Cofre) error
}
```

**Manager construction:** Receives `*Cofre` already loaded + repository via dependency injection.

---

## Deviation from Specification

### modelo-dominio.md Line 73 Override

**Spec states:** "Alterar nome, campos ou favorito de segredo `original`: → `modificado`"

**User decision (D-11):** Favoritar does NOT transition `estadoSessao`

**Documented as deviation** in CONTEXT.md with rationale.

---

## Deferred Topics

### No Validation in Phase 3

**D-30:** No UTF-8 encoding validation, no field size limits
- Simplicity for Phase 3
- TUI responsible for sending valid UTF-8
- Can be added later if needed

### Performance Optimizations

**ListarFavoritos() caching (D-21):**
- Current: O(n) traversal on each call
- Future: Internal cache invalidated on mutations
- Transparent to TUI (no API changes)
- Only if profiling reveals bottleneck

---

## Workflow Compliance

This discussion followed the `discuss-phase.md` workflow:

1. ✅ **Loaded prior context** - PROJECT.md, REQUIREMENTS.md, STATE.md, Phase 1 & 2 CONTEXT.md
2. ✅ **Identified canonical references** - modelo-dominio.md, arquitetura.md, arquitetura-camada-dominio.md, formato-arquivo-abditum.md
3. ✅ **Explored gray areas** - 4 major areas (identity, validation, exposure, state) + 6 follow-up questions
4. ✅ **Captured decisions** - 30 decisions with IDs D-01 through D-30
5. ✅ **Documented deviations** - modelo-dominio.md line 73 override
6. ✅ **Deferred optimizations** - Caching, validation (documented as future work)

**Next steps:**
- Researcher agent reads CONTEXT.md + canonical refs to focus investigation
- Planner agent creates executable tasks based on decisions
- No need to ask user again - all decisions captured

---

## Key Insights from Discussion

1. **Go package encapsulation is powerful:** Fields lowercase = inaccessible outside package. No need for complex snapshot patterns.

2. **Pointers as identity:** In-process graph doesn't need synthetic IDs. JSON hierarchical structure is sufficient for persistence.

3. **Make impossible > validate:** Design structures so invalid states cannot be represented (Observação encapsulation).

4. **Two independent concerns:** `cofre.modificado` (any change needing save) vs `segredo.estadoSessao` (content change indicator).

5. **Manager knows WHAT, entities know HOW:** Clear separation between high-level orchestration and low-level implementation.

6. **Atomic operations via two-phase:** Validate (read-only) → Mutate (cannot fail). No rollback complexity needed.

7. **Good API enables future optimization:** ListarFavoritos() can be internally cached later without TUI changes.

8. **Separation of responsibilities:** Factory (construction) vs Initializer (bootstrap) vs Manager (user operations).

---

## Questions Asked by User

1. ✅ "Temos mais perguntas?" - Led to 6 follow-up questions
2. ✅ "Explique melhor cada uma das perguntas" - Detailed explanation with scenarios
3. ✅ "O manager deverá delegar o máximo possível..." - Led to D-25 refinement
4. ✅ "O que você quer dizer com Manager não contém lógica de negócio?" - Led to critical clarification of D-04/D-25
5. ✅ "Entendo que os métodos de Pasta poderão responder..." - Confirmed D-02 understanding
6. ✅ "Na sua opinião, o manager poderia manter cache dos favoritos?" - Led to D-21 optimization discussion
7. ✅ "Eu preferiria que os métodos Reordenar se chamassem Reposicionar" - Led to D-23 nomenclature refinement
8. ✅ "Eu prefiro factory + método de inicialização separado" - Led to D-28 separation of responsibilities

All questions led to clarifications that improved decision quality.

---

## Validation Complete

User confirmed: "tudo ok!"

All 30 decisions validated and ready for CONTEXT.md creation.

---

*Discussion completed: 2025-03-29*
*Total duration: Comprehensive context discussion*
*Output: 03-CONTEXT.md with 30 decisions + this discussion log*
