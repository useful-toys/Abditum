# Phase 3: Vault Domain + Manager - Context

**Gathered:** 2025-03-29
**Status:** Ready for planning

<domain>
## Phase Boundary

Implement the domain layer (`internal/vault`) with all entities (Cofre, Pasta, Segredo, ModeloSegredo, CampoSegredo, Configuracoes), the Manager pattern for orchestration, and complete business logic for vault operations. This phase delivers the core domain model with in-memory operations, validation, state management, and the interface for persistence (implementation in Phase 4).

**Fixed from ROADMAP.md:**
- Domain entities with encapsulation via package-private fields
- Manager as single public API for TUI
- Session state tracking (original/incluido/modificado/excluido)
- Validation logic (uniqueness, cycles, invariants)
- Factory methods for entity creation
- Search and favorites functionality
- Ready for storage layer integration (Phase 4)

</domain>

<decisions>
## Implementation Decisions

### Identity and References

**D-01: No synthetic IDs - Go pointers as identity**
- Entities identified by Go pointers (`*Pasta`, `*Segredo`, `*ModeloSegredo`) during session
- JSON uses hierarchical structure (nested objects/arrays) without `id` or `parentId` fields
- Uniqueness validated by composite key `(parent, nome)` at mutation time
- Rename doesn't break references - pointers stable throughout session
- Eliminates: ID generation, collision detection, storage overhead, dual identity confusion

**D-02: Uniqueness validated via entity query methods**
- **Segredo**: `(pasta, nome)` unique within pasta - `Pasta.contemSegredoComNome()` validates
- **Pasta**: `(pai, nome)` unique among siblings - `Pasta.contemSubpastaComNome()` validates
- **ModeloSegredo**: `nome` globally unique - `Cofre.contemModeloComNome()` validates
- Entity owns query methods (private), uses them in validation phase
- Manager/Cofre orchestrate by calling entity methods, don't implement verification logic
- Conflicts return `ErrNameConflict` sentinel error

**D-03: Rename doesn't break references**
- During session: TUI maintains Go pointers (stable)
- In JSON: graph reconstructed deterministicly on deserialization
- Storage populates `pai`/`pasta` references via single O(n) traversal after load
- References (`segredo.pasta`, `pasta.pai`) are non-persisted - exist only in memory

### Validation and Invariants

**D-04: Entities validate and execute, Manager orchestrates**
- **Entities** contain all validation and mutation logic (private methods)
- **Manager** knows WHAT operations exist and their semantics (high-level knowledge)
- Manager does NOT know HOW to validate or execute (implementation details)
- Manager delegates to entity methods, updates global state (`cofre.modificado`, timestamps)
- Analogy: Manager is maestro (knows the music), entities are musicians (play instruments)

**D-05: Eager validation (fail-fast)**
- Operations return error immediately if invariants would be violated
- No partial mutations or invalid intermediate states
- Two-phase pattern: validate (read-only) → mutate (cannot fail)

**D-06: Make impossible > validate (structural enforcement)**
- **Category A** (structurally impossible):
  - Observação immutability via encapsulation: `Campos()` returns only user fields
  - Manipulation methods operate on user fields slice, Observação invisible to edit operations
  - Observação exists in internal slice but architecturally unreachable for mutation
- **Category B** (explicit validation):
  - Name uniqueness, Pasta Geral protection, cycle detection
  - "Observação" name forbidden in ModeloSegredo (explicit check when adding field)

**D-07: Hybrid error strategy**
- **Sentinel errors** for simple validations where TUI only needs type check
- **Custom error structs** when TUI needs structured data from error
- Examples: `ErrNameConflict`, `ErrPastaGeralProtected`, `ErrCycleDetected`, `ErrObservacaoReserved`, `ErrConfigInvalida`, `ErrPosicaoInvalida`, `ErrSegredoNaoEncontrado`, `ErrNomeVazio`, `ErrNomeMuitoLongo`

### Vault Exposure and Encapsulation

**D-08: Manager.Vault() returns live `*Cofre` pointer**
- TUI receives direct pointer to aggregate root
- Safety via package encapsulation: all fields lowercase (private to `internal/vault`)
- Zero copy overhead - optimal performance for large vaults
- Mutation impossible from TUI - requires Manager method calls
- TUI navigates via exported getters: `Cofre.PastaGeral()`, `Pasta.Subpastas()`, `Pasta.Segredos()`, `Segredo.Pasta()`, `Pasta.Pai()`

**D-09: Getters return defensive copies of mutable collections**
- Slices returned by getters are copies to prevent indirect mutation
- Example: `func (s *Segredo) Campos() []CampoSegredo` returns `copy(copia, s.camposUsuario)`
- Prevents TUI from mutating internal state via `append()` or index assignment
- ROADMAP pitfall "Manager.Vault() must return snapshot" satisfied by defensive copies in getters

**D-10: Secret values via explicit `ValorComoString()` method**
- `CampoSegredo.valor` is `[]byte` private to package
- TUI calls `campo.ValorComoString()` to get display string
- Method name signals security boundary crossing
- Conversion irreversible (string immutable, cannot be wiped)
- Natural interception point for future audit/rate limiting

### Session State Management

**D-11: Two independent flags - favoriting does NOT change estadoSessao**
- **`cofre.modificado`**: "vault has unsaved changes" - changes on ANY mutation including favoriting
- **`segredo.estadoSessao`**: "secret CONTENT changed" - does NOT change on favoriting
- Favoriting is navigation preference (metadata), not content edit
- User does NOT see "modified" indicator when favoriting secret
- User DOES see unsaved changes indicator at vault level
- **Overrides modelo-dominio.md line 73** (document deviation)

**D-12: Change detection based on actual value difference**
- Flags only change if resulting value actually differs from current
- Entity private methods return `(alterado bool, err error)` to signal real change
- Cofre/Manager use feedback to decide whether to update flags and timestamps
- Example: user opens rename dialog, doesn't change anything, confirms → no modification indicator
- Applies to both `cofre.modificado` and `segredo.estadoSessao` transitions

**D-13: EstadoSessao transition rules**
- `original` → `modificado`: Any content mutation with actually different value
- `incluido` → `incluido`: Mutations on newly-created secret stay `incluido` (not yet in file)
- Any → `excluido`: When marked for deletion; previous state memorized
- `excluido` → previous state: When unmarking deletion (restore)

**D-14: Deleted secrets filtering policy**
- `Pasta.Segredos()` returns ALL secrets including `excluido` (no separate "including deleted" method)
- TUI renders deleted with strikethrough - visible in normal listing
- Filtering happens only in specific use cases:
  - **Search**: `Manager.BuscarSegredos()` filters `excluido` internally
  - **Export**: Export logic filters `excluido`
  - **Save**: Storage removes `excluido` permanently from JSON

**D-15: Search responsibility split**
- **`Segredo.AtendeCriterio(criterio string) bool`**: Entity knows if content matches criterion
  - Normalized substring search (no accents, case-insensitive)
  - Searches: secret name, all field names, values of common fields and observation
  - Sensitive field VALUES never participate - only sensitive field NAME
- **`Manager.BuscarSegredos(query string) []*Segredo`**: Manager applies search policy
  - Traverses entire tree
  - Filters out `excluido`
  - Delegates matching to `Segredo.AtendeCriterio()`
- TUI can call `segredo.AtendeCriterio()` directly for individual checks (e.g., highlighting while navigating)

**D-16: Deletion finalized on save**
- Secrets marked `excluido` do NOT persist in JSON
- Save operation permanently removes them
- Undo is session-only feature - once saved, deletion irreversible
- Aligns with "save is explicit commit" model

**D-17: Save is atomic - two-phase commit for deletions**
- **Phase 1**: `prepararSnapshot()` creates deep copy with `excluido` filtered (live vault untouched)
- **Phase 2**: `repositorio.Salvar(snapshot)` - if fails, live vault unchanged
- **Phase 3**: `finalizarExclusoes()` removes `excluido` from memory only if save succeeded
- Guarantees: save failure doesn't cause data loss in memory
- Tradeoff: deep copy has memory cost, but only during save (infrequent operation)

### Additional Specifications

**D-18: NanoID pitfall obsolete**
- ROADMAP warning "NanoID must use crypto/rand" obsolete
- D-01 eliminated need for synthetic identifiers entirely
- No random ID generation required in Phase 3

**D-19: Timestamps only on Segredo and Cofre**
- **`Segredo`**: `dataCriacao` and `dataUltimaModificacao`
- **`Cofre`**: `dataCriacao` and `dataUltimaModificacao`
- **`Pasta`**: NO timestamps - no need to audit hierarchy changes via folders
- Simpler model: only content entities (secrets) + aggregate root track temporal metadata

**D-20: Configuracoes in Phase 3 scope with mutation method**
```go
type Configuracoes struct {
    tempoBloqueioInatividadeMinutos      int  // default: 5
    tempoOcultarSegredoSegundos          int  // default: 15
    tempoLimparAreaTransferenciaSegundos int  // default: 30
}
```
- All timers mandatory (cannot be disabled/zero/negative)
- Manager exposes `AlterarConfiguracoes(novasConfig Configuracoes) error`
- Validation: all values must be > 0
- Marks `cofre.modificado = true` and updates `dataUltimaModificacao`
- New sentinel error: `ErrConfigInvalida`

**D-21: ListarFavoritos() with depth-first traversal**
- Manager exposes `Manager.ListarFavoritos() []*Segredo`
- Traversal uses **DFS (depth-first)** following JSON order (per modelo-dominio.md line 176)
- Filters out `excluido` (consistent with search D-15)
- `Segredo` exposes getter `func (s *Segredo) Favorito() bool` for TUI individual checks
- Order: processes secrets in current folder first, then recurses into subfolders
- **No cache in Phase 3**: Simple on-demand traversal O(n). Future optimization opportunity if profiling reveals bottleneck - cache can be added internally without API changes (transparent to TUI).

**D-22: Factory methods with position parameter**
- Pasta factory methods (private): `criarSubpasta(nome, posicao)`, `criarSegredo(nome, campos, posicao)`, `criarSegredoDeModelo(nome, modelo, posicao)`
- Manager public API: `CriarPasta(pai, nome, posicao)`, `CriarSegredo(pasta, nome, campos, posicao)`, `CriarSegredoDeModelo(pasta, nome, modelo, posicao)`
- **Position semantics**:
  - `posicao` is 0-indexed
  - `posicao == len(slice)` means "append at end"
  - `posicao < len(slice)` means "insert here, shift existing elements right (+1)"
  - Invalid position (negative or > len) returns error
- Factory validates position before mutation

**D-23: Explicit repositioning operations**
- Manager exposes:
  - `ReposicionarSegredo(segredo, novaPosicao int) error`
  - `ReposicionarPasta(pasta, novaPosicao int) error`
  - `SubirSegredoNaPosicao(segredo) error` (posicao - 1)
  - `DescerSegredoNaPosicao(segredo) error` (posicao + 1)
  - `SubirPastaNaPosicao(pasta) error`
  - `DescerPastaNaPosicao(pasta) error`
- **Nomenclature**: "Reposicionar" (move to position) not "Reordenar" (apply ordering criterion like alphabetic)
- TUI mapping: `Ctrl+Up` → `SubirXNaPosicao()`, `Ctrl+Down` → `DescerXNaPosicao()`
- Edge cases (consistent with D-12):
  - `Subir` at position 0 → no-op, returns `nil`, doesn't mark modified
  - `Descer` at last position → no-op, returns `nil`, doesn't mark modified
  - `Reposicionar` to current position → no-op, returns `nil`, doesn't mark modified

### Timestamp Update Policy

**D-24: Timestamps update only on structural changes**

**`Segredo.dataUltimaModificacao` updates ONLY on**:
- ✅ Rename
- ✅ Field value changes (actually different value)
- ✅ Add/remove/reorder fields
- ✅ Change field names or types
- ✅ Move to different pasta (parent reference changes - structural)
- ✅ Mark/unmark deletion
- ❌ **NOT on favoriting** (metadata, not content)
- ❌ **NOT on repositioning within same pasta** (position is container concern)

**`Cofre.dataUltimaModificacao` updates on**:
- ✅ Any secret content change (which updates Segredo timestamp)
- ✅ Any pasta structural change (create/rename/move/delete)
- ✅ Any modelo change
- ✅ **Favoriting** (updates Cofre timestamp but not Segredo timestamp)
- ✅ **Repositioning** (updates Cofre timestamp but not entity timestamps)
- ✅ Configuration changes

**Rationale**:
- Favoriting is user preference, not content change
- Repositioning is container property, not entity content
- Move changes parent reference (structural to entity)

### Architectural Pattern

**D-25: Manager as orchestrator with high-level business knowledge**

**Manager POSSESSES (high-level business knowledge)**:
- ✅ Knowledge of domain operations ("create secret", "move pasta", "favoritar")
- ✅ Semantics of operations ("move updates timestamp", "favoriting doesn't change content")
- ✅ Workflows: order of steps in complex operations (validate → mutate → update global state)
- ✅ Global rules: "any mutation marks cofre.modificado", "save is atomic with two-phase commit"
- ✅ Relationships: Segredo belongs to Pasta, Pasta has pai, operations can cross entity boundaries

**Manager does NOT POSSESS (implementation logic)**:
- ❌ Validation algorithms: how to verify name unique, detect cycles, generate suffixed names
- ❌ Manipulation algorithms: how to insert at position, remove maintaining order, reposition element
- ❌ Internal structures: doesn't access `pasta.segredos` directly, doesn't know memory layout
- ❌ Entity construction: doesn't instantiate `&Segredo{...}`, doesn't know Observação is last field

**Pattern**: Manager knows **WHAT** to do (workflows, semantics), entities know **HOW** to execute (algorithms, structures)

**Atomic operations** (refines D-05):
- Entities implement two-phase: validation (read-only) → mutation (cannot fail)
- Validation uses query methods: `contemNome()`, `encontrarPosicao()`, `detectarCiclo()`
- If validation fails, return error without touching state
- If validation passes, mutation executes (no error returns needed)
- Manager receives result, updates global state only if actual change occurred (D-12)

**Exception**: Save (D-17) has unavoidable I/O failure, but maintains atomicity (vault stays in valid pre-save state)

### Navigation and Access

**D-26: Navigation via entity getters only**
- TUI navigates using only exported getters from entities
- Manager does NOT expose navigation helper methods (avoid unless specific technical need)
- Examples:
  - ✅ `manager.Vault().PastaGeral()`
  - ✅ `pasta.Subpastas()`
  - ✅ `pasta.Segredos()`
  - ✅ `segredo.Pasta()`
  - ✅ `pasta.Pai()`
  - ❌ `manager.ObterPastaRaiz()` (unnecessary)
  - ❌ `manager.ObterTodosModelos()` (unnecessary)

### Deletion and Lifecycle

**D-27: Pasta has hard delete, Segredo has soft delete**

**Segredo: Soft delete**
- `ExcluirSegredo(segredo)` marks `estadoSessao = excluido`
- Remains in memory until save succeeds
- `RestaurarSegredo(segredo)` unmarks flag
- Permanently removed only on save commit (D-16, D-17)

**Pasta: Hard delete**
- `ExcluirPasta(pasta)` removes immediately from hierarchy
- Subpastas and segredos promoted to parent
- Name conflicts generate automatic renaming
- Returns `[]Renomeacao` for TUI to communicate to user
- No "restore pasta" - operation irreversible in session
- Requires explicit user confirmation (TUI responsibility)

**Rationale**:
- Segredo is leaf entity → undo trivial (flip flag)
- Pasta has dependents → deletion is complex restructuring
- Promoting children maintains invariant "every secret has parent pasta"

**D-28: Factory NovoCofre() + Initializer InicializarConteudoPadrao()**

**D-28a: NovoCofre() - empty vault factory**
```go
func NovoCofre() *Cofre
```
- **Responsibility**: Construct base aggregate structure
- Returns vault with Pasta Geral and default configurations
- NO initial content (folders/templates)
- Useful for tests and special cases

**D-28b: Cofre.InicializarConteudoPadrao() - bootstrap initial content**
```go
func (c *Cofre) InicializarConteudoPadrao() error
```
- **Responsibility**: Populate new vault with canonical initial structure
- Creates default folders: "Sites e Apps", "Financeiro"
- Creates default templates: "Login", "Cartão de Crédito", "Chave de API"
- Per functional requirement "Criar novo cofre" in requisitos.md
- Must be called ONCE after `NovoCofre()` for new user
- Does NOT mark `cofre.modificado = true` (initial content is part of base state)
- NOT a user operation via Manager (system bootstrap, not domain operation)

**Separation of responsibilities**:
- **NovoCofre()**: Construction (factory) - minimal structure
- **InicializarConteudoPadrao()**: Bootstrap (domain service) - initial content
- **Manager**: User operations during normal use

**Typical usage**:
```go
// cmd/abditum - create new vault (first time)
cofre := vault.NovoCofre()
cofre.InicializarConteudoPadrao()
manager := vault.NewManager(cofre, repo)
manager.Salvar()

// Tests - empty vault
cofre := vault.NovoCofre()
manager := vault.NewManager(cofre, repo)
```

### Field Access

**D-29: Campos() excludes Observação - dedicated getter for Observação**

```go
// Returns only user fields (EXCEPT Observação)
func (s *Segredo) Campos() []CampoSegredo

// Returns observation value as string
func (s *Segredo) Observacao() string
```

**Benefits**:
- ✅ TUI doesn't need to know Observação is "the last field"
- ✅ Manipulation methods operate only on user fields
- ✅ Observação never appears in edit operation indices
- ✅ Consistent with D-06 (make structurally impossible)
- ✅ Clearer API: `segredo.Observacao()` vs `segredo.Campos()[len-1].ValorComoString()`

**Implication**: Indices in mutation methods always refer to user fields (0 to len-1 of what `Campos()` returns), never to Observação.

**Mutation methods**:
```go
func (m *Manager) AdicionarCampoSegredo(segredo, nome, tipo, valor, posicao)
func (m *Manager) AlterarCampoSegredo(segredo, indiceCampo, novoValor)
func (m *Manager) RemoverCampoSegredo(segredo, indiceCampo)
func (m *Manager) ReordenarCampoSegredo(segredo, indiceCampo, novaPosicao)
func (m *Manager) AlterarObservacao(segredo, novoValor)  // dedicated method
```

### Field Validation

**D-30: No sensitive field validation in Phase 3**
- ✅ Accepts any `[]byte`
- ✅ Any size (no maximum limit)
- ✅ Empty values allowed (`[]byte{}`)
- ⚠️ Spec says "always UTF-8", but no encoding validation
- TUI responsibility to send valid UTF-8

**Rationale**:
- Simplicity for Phase 3
- Validation can be added later if needed
- Size limits defined in future based on real requirements

**Document as deferred**: "UTF-8 encoding validation and size limits can be added in future phase if specification or testing reveals need."

### the agent's Discretion

- Exact internal field names (as long as package-private and follow Go conventions)
- Helper method names (as long as purpose clear)
- Error message wording (as long as communicates problem)
- Internal validation helper function organization
- Test organization and coverage strategy (must include race detector per Phase 2 pattern)

</decisions>

<specifics>
## Specific Ideas

### Architectural References
- **Manager pattern**: Thin orchestrator like a maestro - knows the music (workflows), entities play instruments (algorithms)
- **Encapsulation**: Go package-level, not class-level - all fields lowercase, only exported getters
- **Two-phase operations**: Validate (read-only) → Mutate (cannot fail after validation passes)

### Code patterns from Phase 2
- Sentinel errors for simple cases: `var ErrXYZ = errors.New("message")`
- `[]byte` for sensitive data with explicit wipe
- Fail-fast error handling
- Comprehensive testing with race detector (`go test -race`)

### Domain model specifics
- Observação always last field in internal slice, but invisible to external manipulation
- Pasta Geral protection via guard clauses (check `pasta == cofre.PastaGeral`)
- Favoritos is calculated view (DFS traversal), not persisted collection
- Exclusão de pasta promove filhos com renomeação automática em caso de conflito

</specifics>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Domain Model
- `modelo-dominio.md` — Complete domain model (244 lines): entities (Cofre, Pasta, Segredo, ModeloSegredo, CampoSegredo, Configuracoes), identity rules (composite keys), session state transitions (original/incluido/modificado/excluido), field immutability rules, ordering preservation rules
  - **Deviation noted**: Line 73 states favoriting changes `estadoSessao` - D-11 overrides this (favoriting does NOT change estadoSessao)

### Architecture
- `arquitetura.md` — Manager pattern (171 lines): Manager orchestrates + Cofre contains logic, transaction boundaries, package structure (cmd/abditum, internal/vault, internal/crypto, internal/storage, internal/tui), no network deps policy, minimal external dependencies, generous comments policy, sensitive data as []byte, CGO_ENABLED=0 static builds

- `arquitetura-camada-dominio.md` — Detailed domain layer architecture decisions: encapsulation via package-private fields, bidirectional navigation (pasta.pai, segredo.pasta), Cofre as coordinator not god object, Manager as thin orchestrator, session state tracking, search responsibility split, timestamp update policy, memory wiping limitations in Go

### Storage Format
- `formato-arquivo-abditum.md` — Binary format spec (referenced for crypto parameters but not critical for domain layer - relevant for Phase 4 storage)

### Requirements
- `.planning/REQUIREMENTS.md` — All v1 requirements with traceability: VAULT-02, SEC-05, FOLDER-01 through 05, TPL-01 through 06 mapped to Phase 3

### Phase Context
- `.planning/ROADMAP.md` — Phase 3 details including 7 implementation steps, UAT criteria, pitfall watch (Manager.Vault() must return snapshot - satisfied by D-08/D-09; cycle detection required; Pasta Geral protection; NanoID obsolete per D-18)

### Prior Phases
- `.planning/phases/01-project-scaffold-ci-foundation/01-CONTEXT.md` — Phase 1 decisions: module path github.com/useful-toys/abditum, Go 1.26.1+, CI on main branch
- `.planning/phases/02-crypto-package/02-CONTEXT.md` — Phase 2 decisions: sentinel errors pattern with 12 design decisions, []byte for sensitive data, fail-fast error handling, comprehensive testing with race detector - patterns to follow in Phase 3

### Existing Code
- `internal/vault/doc.go` — Package stub (2 lines, empty doc comment, ready for implementation)
- `internal/crypto/` — Completed crypto package from Phase 2 (reusable for password/key handling: aead.go, kdf.go, memory.go, password.go, errors.go with sentinels)

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets from Phase 2
- **Sentinel error pattern**: `var ErrXYZ = errors.New("message")` in `internal/crypto/errors.go` - replicate in `internal/vault/errors.go`
- **Memory wiping**: `crypto.WipeBytes([]byte)` - use for sensitive field values when vault locks
- **Testing with race detector**: All Phase 2 tests run with `-race` flag - continue pattern in Phase 3

### Established Patterns
- **Package-level encapsulation**: Fields lowercase (private to package), exported getters only
- **`[]byte` for sensitive data**: Never `string` for passwords, keys, or secret field values
- **Fail-fast errors**: Return early on validation failure, no partial state mutations
- **Generous comments**: Public API fully documented with examples and edge cases

### Integration Points
- **Storage layer (Phase 4)**: `RepositorioCofre` interface defined in Phase 3, implemented in Phase 4
  - Interface: `Salvar(cofre *Cofre) error`
  - Manager receives repository via dependency injection
- **TUI layer (Phase 5+)**: Interacts only via Manager public API
  - Navigates vault via entity getters
  - Never accesses package-private fields
  - Receives errors, displays to user with friendly messages
- **Crypto package**: `internal/crypto` used for password/key handling when integrating with storage

</code_context>

<deferred>
## Deferred Ideas

### Future Optimizations
- **ListarFavoritos() caching** (D-21): Current O(n) traversal on each call. If profiling reveals bottleneck, add internal cache invalidated on mutations - transparent to TUI (no API changes).
- **UTF-8 encoding validation** (D-30): Field values currently accept any `[]byte` without encoding validation. Add if specification or testing reveals need.
- **Field size limits** (D-30): No maximum size enforced currently. Define limits in future based on real requirements.

### Phase 4: Storage Layer
- Implementation of `RepositorioCofre` interface
- JSON serialization/deserialization with custom marshalers
- AES-256-GCM encryption using `internal/crypto`
- Atomic file writes
- Reference population after deserialization (`popularReferencias()`)

### Phase 5+: TUI Layer
- Bubble Tea v2 components for vault navigation
- Secret editing forms with field type awareness (sensitive vs common)
- Favoritos view (calls `Manager.ListarFavoritos()`)
- Search interface (calls `Manager.BuscarSegredos()`)
- Confirmation dialogs for irreversible operations (delete pasta)

### Future Enhancements (Post-v1)
- Undo/redo beyond session (requires event sourcing)
- Conflict resolution for concurrent edits
- Export/import functionality
- Backup/restore operations
- Audit log of all mutations

</deferred>

---

*Phase: 03-vault-domain-manager*
*Context gathered: 2025-03-29*
*Total decisions: 30 (D-01 through D-30)*
*Canonical documents: 8 specs + 2 prior phase contexts*
