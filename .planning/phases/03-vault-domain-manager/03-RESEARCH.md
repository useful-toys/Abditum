# Phase 3 Research: Vault Domain + Manager

**Phase:** 03-vault-domain-manager
**Research Date:** 2026-03-29
**Status:** Complete

## Executive Summary

Phase 3 implements a rich domain layer with complex business rules, state tracking, and validation logic. Key research areas include Go's encapsulation patterns for domain modeling, pointer-based identity management, defensive copying strategies, state machine implementation, and cycle detection algorithms.

**Critical Insight:** The CONTEXT.md decisions (D-01 through D-30) eliminate several ROADMAP concerns (NanoID, snapshot copies) and establish a clear pointer-based identity model. Plans must honor these locked decisions.

## Domain-Driven Design in Go

### Package-Level Encapsulation (D-08, D-09)

**Pattern:** Go uses package-level visibility (lowercase = private to package, uppercase = exported). This is MORE restrictive than object-level encapsulation in languages like Java.

**Application:**
```go
// internal/vault/entities.go
type Cofre struct {
    pastaGeral          *Pasta
    modelos             []*ModeloSegredo
    modificado          bool
    dataCriacao         time.Time
    dataUltimaModificacao time.Time
    configuracoes       Configuracoes
}

// Exported getter returns pointer (read-only via encapsulation)
func (c *Cofre) PastaGeral() *Pasta {
    return c.pastaGeral
}

// Defensive copy for mutable collections
func (c *Cofre) Modelos() []*ModeloSegredo {
    copia := make([]*ModeloSegredo, len(c.modelos))
    copy(copia, c.modelos)
    return copia
}
```

**Benefits:**
- TUI cannot mutate internal state even with live pointers
- Zero-copy performance for navigation
- Prevents `append()` or index assignment corruption
- Type system enforces boundaries

**Reference:** Effective Go - Names section on package-level visibility

### Pointer Identity in Go (D-01, D-03)

**Pattern:** Go pointers are stable memory addresses that survive across function calls and mutations within a session. No need for synthetic IDs when references don't cross process boundaries.

**Application:**
```go
// No ID fields needed
type Segredo struct {
    nome                 string
    campos               []CampoSegredo
    observacao           CampoSegredo  // Last field, separate from campos
    pasta                *Pasta        // Back-reference (not serialized)
    favorito             bool
    estadoSessao         EstadoSessao
    dataCriacao          time.Time
    dataUltimaModificacao time.Time
}

// Manager maintains references
type Manager struct {
    cofre    *Cofre
    repositorio RepositorioCofre
    senha    []byte
    caminho  string
    bloqueado bool
}
```

**Uniqueness via query methods (D-02):**
```go
// Pasta validates uniqueness
func (p *Pasta) contemSubpastaComNome(nome string) bool {
    for _, sub := range p.subpastas {
        if sub.nome == nome {
            return true
        }
    }
    return false
}
```

**Trade-offs:**
- ✅ No ID generation, collision handling, or dual identity confusion
- ✅ Rename doesn't break references (pointer unchanged)
- ✅ Simpler serialization (hierarchical structure)
- ⚠️ JSON deserialization requires reference reconstruction pass
- ⚠️ Cross-process communication would need ID layer (not needed for v1)

**Reference:** Go FAQ - Why is there no pointer arithmetic?

## State Management Patterns

### Two-Phase State Tracking (D-11, D-12, D-13)

**Pattern:** Separate flags for vault-level and entity-level change tracking with change detection based on actual value difference.

**Implementation:**
```go
type Cofre struct {
    modificado bool  // ANY mutation (including favoriting)
    // ...
}

type Segredo struct {
    estadoSessao EstadoSessao  // Content changes only (NOT favoriting)
    // ...
}

type EstadoSessao int

const (
    EstadoOriginal EstadoSessao = iota  // Loaded from file, unmodified
    EstadoIncluido                       // Created this session
    EstadoModificado                     // Existed in file, changed this session
    EstadoExcluido                       // Marked for deletion
)

// Deleted secrets store previous state for restoration
type SegredoExcluido struct {
    segredo      *Segredo
    estadoAnterior EstadoSessao
}
```

**Change detection pattern:**
```go
// Entity methods return bool indicating actual change
func (s *Segredo) renomear(novoNome string) (alterado bool, err error) {
    if s.nome == novoNome {
        return false, nil  // No actual change
    }
    s.nome = novoNome
    return true, nil
}

// Manager uses feedback to update flags
func (m *Manager) RenomearSegredo(segredo *Segredo, novoNome string) error {
    alterado, err := segredo.renomear(novoNome)
    if err != nil {
        return err
    }
    if alterado {
        m.atualizarEstadoSessao(segredo)  // original -> modificado
        m.cofre.modificado = true
        m.cofre.dataUltimaModificacao = time.Now().UTC()
    }
    return nil
}
```

**State transitions (D-13):**
```
original --[content mutation with actual change]--> modificado
incluido --[content mutation]--> incluido  (stays incluido)
any --[mark deletion]--> excluido  (stores previous state)
excluido --[restore]--> previous state
```

**Favoriting exception (D-11):**
```go
func (m *Manager) FavoritarSegredo(segredo *Segredo, favorito bool) error {
    if segredo.favorito == favorito {
        return nil  // No change
    }
    segredo.favorito = favorito
    m.cofre.modificado = true  // Vault modified
    m.cofre.dataUltimaModificacao = time.Now().UTC()
    // Note: segredo.estadoSessao NOT changed (D-11 override)
    return nil
}
```

**Filtering policy (D-14):**
- `Pasta.Segredos()` returns ALL including `excluido`
- `Manager.BuscarSegredos()` filters out `excluido`
- Export filters out `excluido`
- Save removes `excluido` permanently

### Atomic Save with Two-Phase Commit (D-17)

**Pattern:** Prepare immutable snapshot, persist, then finalize in-memory deletions only on success.

**Implementation:**
```go
func (m *Manager) Salvar() error {
    // Phase 1: Create deep copy with excluido filtered (live vault untouched)
    snapshot := m.prepararSnapshot()
    
    // Phase 2: Persist snapshot (if fails, live vault unchanged)
    if err := m.repositorio.Salvar(snapshot); err != nil {
        return err
    }
    
    // Phase 3: Finalize deletions in memory only after successful save
    m.finalizarExclusoes()
    m.cofre.modificado = false
    return nil
}

func (m *Manager) prepararSnapshot() *Cofre {
    // Deep copy entire vault, recursively filtering excluido
    snapshot := &Cofre{
        dataCriacao: m.cofre.dataCriacao,
        dataUltimaModificacao: m.cofre.dataUltimaModificacao,
        configuracoes: m.cofre.configuracoes,
        // ...
    }
    snapshot.pastaGeral = m.copiarPastaRecursivamente(m.cofre.pastaGeral, true)
    snapshot.modelos = m.copiarModelos(m.cofre.modelos)
    return snapshot
}

func (m *Manager) copiarPastaRecursivamente(pasta *Pasta, filtrarExcluidos bool) *Pasta {
    copia := &Pasta{
        nome: pasta.nome,
        subpastas: make([]*Pasta, 0, len(pasta.subpastas)),
        segredos: make([]*Segredo, 0),
    }
    
    // Copy subfolders recursively
    for _, sub := range pasta.subpastas {
        copia.subpastas = append(copia.subpastas, m.copiarPastaRecursivamente(sub, filtrarExcluidos))
    }
    
    // Copy secrets, filtering excluido if requested
    for _, seg := range pasta.segredos {
        if filtrarExcluidos && seg.estadoSessao == EstadoExcluido {
            continue  // Skip deleted
        }
        copiaSegredo := m.copiarSegredo(seg)
        copia.segredos = append(copia.segredos, copiaSegredo)
    }
    
    return copia
}

func (m *Manager) finalizarExclusoes() {
    // Remove excluido secrets from live vault
    m.removerExcluidosRecursivamente(m.cofre.pastaGeral)
}
```

**Trade-offs:**
- ✅ Atomicity: save failure leaves vault in valid pre-save state
- ✅ Safety: no half-completed deletions
- ⚠️ Memory cost: deep copy during save (acceptable - infrequent operation)

## Validation and Invariants

### Two-Phase Validation Pattern (D-05, D-25)

**Pattern:** Validation phase (read-only, can fail) followed by mutation phase (cannot fail after validation passes).

**Implementation:**
```go
// Entity owns validation logic
func (p *Pasta) validarCriacaoSubpasta(nome string, posicao int) error {
    if nome == "" {
        return ErrNomeVazio
    }
    if len(nome) > 255 {
        return ErrNomeMuitoLongo
    }
    if p.contemSubpastaComNome(nome) {
        return ErrNameConflict
    }
    if posicao < 0 || posicao > len(p.subpastas) {
        return ErrPosicaoInvalida
    }
    return nil
}

// Entity owns mutation logic
func (p *Pasta) criarSubpasta(nome string, posicao int) *Pasta {
    // Validation already passed, no errors possible
    novaPasta := &Pasta{
        nome: nome,
        pai: p,
        subpastas: make([]*Pasta, 0),
        segredos: make([]*Segredo, 0),
    }
    
    // Insert at position
    p.subpastas = append(p.subpastas[:posicao], append([]*Pasta{novaPasta}, p.subpastas[posicao:]...)...)
    return novaPasta
}

// Manager orchestrates
func (m *Manager) CriarPasta(pai *Pasta, nome string, posicao int) (*Pasta, error) {
    // Additional Manager-level checks
    if m.bloqueado {
        return nil, ErrCofreBloqueado
    }
    
    // Delegate validation to entity
    if err := pai.validarCriacaoSubpasta(nome, posicao); err != nil {
        return nil, err
    }
    
    // Delegate mutation to entity
    novaPasta := pai.criarSubpasta(nome, posicao)
    
    // Update global state
    m.cofre.modificado = true
    m.cofre.dataUltimaModificacao = time.Now().UTC()
    
    return novaPasta, nil
}
```

**Benefits:**
- Clear separation: validation can fail, mutation cannot
- No partial state or rollback logic needed
- Entity owns both phases (Manager orchestrates, doesn't implement)

### Structural Enforcement (D-06)

**Pattern:** Make impossible states unrepresentable via encapsulation.

**Observação immutability:**
```go
type Segredo struct {
    nome       string
    campos     []CampoSegredo  // User fields only
    observacao CampoSegredo    // Separate, not in campos slice
    // ...
}

// Public getter excludes Observação
func (s *Segredo) Campos() []CampoSegredo {
    copia := make([]CampoSegredo, len(s.campos))
    copy(copia, s.campos)
    return copia
}

// Dedicated getter for Observação
func (s *Segredo) Observacao() string {
    return string(s.observacao.valor)
}

// Manipulation methods only touch campos slice
func (s *Segredo) adicionarCampo(campo CampoSegredo, posicao int) {
    s.campos = append(s.campos[:posicao], append([]CampoSegredo{campo}, s.campos[posicao:]...)...)
}

// Observação mutation requires dedicated method
func (m *Manager) AlterarObservacao(segredo *Segredo, novoValor string) error {
    if string(segredo.observacao.valor) == novoValor {
        return nil  // No change
    }
    segredo.observacao.valor = []byte(novoValor)
    m.atualizarEstadoSessao(segredo)
    m.cofre.modificado = true
    m.cofre.dataUltimaModificacao = time.Now().UTC()
    segredo.dataUltimaModificacao = time.Now().UTC()
    return nil
}
```

**Benefits:**
- TUI cannot accidentally manipulate Observação via index operations
- No runtime checks needed for "is this Observação?"
- Compiler enforces invariant

**"Observação" name forbidden in templates (D-06 Category B):**
```go
func (c *Cofre) validarCampoModelo(nome string) error {
    if nome == "Observação" {
        return ErrObservacaoReserved
    }
    return nil
}
```

## Complex Algorithms

### Cycle Detection (D-03, ROADMAP Pitfall)

**Requirement:** `MoveFolder` must detect if destination is a descendant of the folder being moved.

**Algorithm:** Full ancestor walk from destination up to root, checking for match with source folder.

```go
func (m *Manager) MoverPasta(pasta *Pasta, destino *Pasta) error {
    // Validate not moving Pasta Geral
    if pasta == m.cofre.pastaGeral {
        return ErrPastaGeralProtected
    }
    
    // Validate not moving into self
    if pasta == destino {
        return ErrDestinoInvalido
    }
    
    // Validate cycle: walk ancestors of destination
    atual := destino
    for atual != nil {
        if atual == pasta {
            return ErrCycleDetected  // destino is descendant of pasta
        }
        atual = atual.pai
    }
    
    // Validate name uniqueness in destination
    if destino.contemSubpastaComNome(pasta.nome) {
        return ErrNameConflict
    }
    
    // Mutation: remove from old parent, add to new parent
    pasta.removerDePai()
    pasta.adicionarAoPai(destino)
    
    // Update global state
    m.cofre.modificado = true
    m.cofre.dataUltimaModificacao = time.Now().UTC()
    
    return nil
}
```

**Complexity:** O(depth) where depth is the depth of the destination folder. Worst case O(n) for deeply nested hierarchies, but acceptable for typical vault structures.

**Edge cases:**
- Moving to sibling: no cycle
- Moving to parent: no cycle
- Moving to grandparent: no cycle
- Moving to own child: cycle detected ✓
- Moving to grandchild: cycle detected ✓

### Folder Deletion with Promotion (D-27, FOLDER-05)

**Requirement:** When deleting folder (except Pasta Geral), promote all children and secrets to parent with automatic conflict resolution.

**Algorithm:**
```go
func (m *Manager) ExcluirPasta(pasta *Pasta) ([]Renomeacao, error) {
    // Validate not Pasta Geral
    if pasta == m.cofre.pastaGeral {
        return nil, ErrPastaGeralProtected
    }
    
    renomeacoes := []Renomeacao{}
    pai := pasta.pai
    
    // Promote subfolders
    for _, subpasta := range pasta.subpastas {
        // Check name conflict
        if pai.contemSubpastaComNome(subpasta.nome) {
            // Merge contents into existing folder
            pastaExistente := pai.encontrarSubpastaPorNome(subpasta.nome)
            renomeacoes = append(renomeacoes, m.mesclarPastas(subpasta, pastaExistente)...)
        } else {
            // No conflict, simple promotion
            subpasta.pai = pai
            pai.subpastas = append(pai.subpastas, subpasta)
        }
    }
    
    // Promote secrets
    for _, segredo := range pasta.segredos {
        // Check name conflict
        if pai.contemSegredoComNome(segredo.nome) {
            // Rename with numeric suffix
            novoNome := m.gerarNomeSufixado(segredo.nome, pai)
            renomeacoes = append(renomeacoes, Renomeacao{
                Antigo: segredo.nome,
                Novo: novoNome,
                Pasta: pai.nome,
            })
            segredo.nome = novoNome
        }
        segredo.pasta = pai
        pai.segredos = append(pai.segredos, segredo)
    }
    
    // Remove pasta from parent
    pasta.removerDePai()
    
    // Update global state
    m.cofre.modificado = true
    m.cofre.dataUltimaModificacao = time.Now().UTC()
    
    return renomeacoes, nil
}

func (m *Manager) gerarNomeSufixado(nomeBase string, pasta *Pasta) string {
    contador := 1
    for {
        candidato := fmt.Sprintf("%s (%d)", nomeBase, contador)
        if !pasta.contemSegredoComNome(candidato) {
            return candidato
        }
        contador++
    }
}
```

**Complexity:** O(n * m) where n = items being promoted, m = items in parent. Acceptable for typical folder sizes.

**Important (FOLDER-05):** `StateDeleted` secrets retain their state when promoted - they don't get restored by promotion.

### Secret Duplication with Name Progression (SEC-02, UAT)

**Requirement:** "Name (1)", "Name (2)" progression for successive copies.

**Algorithm:**
```go
func (m *Manager) DuplicarSegredo(segredo *Segredo) (*Segredo, error) {
    pasta := segredo.pasta
    nomeBase := segredo.nome
    
    // Extract existing suffix if present
    baseSemSufixo, sufixoAtual := m.extrairSufixoNumerico(nomeBase)
    
    // Find next available suffix
    proximoSufixo := sufixoAtual + 1
    novoNome := fmt.Sprintf("%s (%d)", baseSemSufixo, proximoSufixo)
    
    // Keep incrementing if name already exists
    for pasta.contemSegredoComNome(novoNome) {
        proximoSufixo++
        novoNome = fmt.Sprintf("%s (%d)", baseSemSufixo, proximoSufixo)
    }
    
    // Create duplicate
    duplicado := m.copiarSegredo(segredo)
    duplicado.nome = novoNome
    duplicado.estadoSessao = EstadoIncluido  // New secret
    duplicado.dataCriacao = time.Now().UTC()
    duplicado.dataUltimaModificacao = time.Now().UTC()
    
    // Insert immediately after original
    posicaoOriginal := pasta.encontrarPosicaoSegredo(segredo)
    pasta.inserirSegredo(duplicado, posicaoOriginal+1)
    
    // Update global state
    m.cofre.modificado = true
    m.cofre.dataUltimaModificacao = time.Now().UTC()
    
    return duplicado, nil
}

func (m *Manager) extrairSufixoNumerico(nome string) (base string, sufixo int) {
    // Match " (N)" at end
    re := regexp.MustCompile(`^(.+) \((\d+)\)$`)
    matches := re.FindStringSubmatch(nome)
    if matches != nil {
        base = matches[1]
        sufixo, _ = strconv.Atoi(matches[2])
        return base, sufixo
    }
    // No suffix, treat entire name as base
    return nome, 0
}
```

**Examples:**
- "X" → duplicate → "X (1)"
- "X (1)" → duplicate → "X (2)"
- "X (2)" → duplicate → "X (3)"
- "X (1)" exists, duplicating "X" → "X (2)" (skip conflict)

### Template Alphabetical Sorting (TPL-06)

**Requirement:** Templates always returned in alphabetical order, regardless of creation order.

**Implementation:**
```go
import "sort"

func (c *Cofre) Modelos() []*ModeloSegredo {
    // Create defensive copy
    copia := make([]*ModeloSegredo, len(c.modelos))
    copy(copia, c.modelos)
    
    // Sort by name
    sort.Slice(copia, func(i, j int) bool {
        return copia[i].nome < copia[j].nome
    })
    
    return copia
}

// After create or rename, re-sort internal slice
func (m *Manager) CriarModelo(nome string, campos []CampoModelo) (*ModeloSegredo, error) {
    // Validation...
    
    modelo := &ModeloSegredo{
        nome: nome,
        campos: campos,
    }
    
    m.cofre.modelos = append(m.cofre.modelos, modelo)
    
    // Re-sort internal slice
    sort.Slice(m.cofre.modelos, func(i, j int) bool {
        return m.cofre.modelos[i].nome < m.cofre.modelos[j].nome
    })
    
    m.cofre.modificado = true
    m.cofre.dataUltimaModificacao = time.Now().UTC()
    
    return modelo, nil
}
```

**Complexity:** O(n log n) on create/rename. Acceptable since template operations are infrequent.

### Search with Sensitive Field Exclusion (QUERY-02)

**Requirement:** Search matches secret name, field names (including sensitive field NAMES), common field VALUES, observation VALUES. Sensitive field VALUES excluded. Normalized substring matching (case-insensitive, accent-insensitive).

**Implementation:**
```go
import (
    "strings"
    "golang.org/x/text/transform"
    "golang.org/x/text/unicode/norm"
)

// Normalize: lowercase + remove accents
func normalizar(s string) string {
    t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
        return r >= 0x0300 && r <= 0x036F  // Combining diacritical marks
    }), norm.NFC)
    resultado, _, _ := transform.String(t, strings.ToLower(s))
    return resultado
}

// Entity method: does this secret match criterion?
func (s *Segredo) AtendeCriterio(criterio string) bool {
    critNorm := normalizar(criterio)
    
    // Match secret name
    if strings.Contains(normalizar(s.nome), critNorm) {
        return true
    }
    
    // Match all field NAMES (including sensitive)
    for _, campo := range s.campos {
        if strings.Contains(normalizar(campo.nome), critNorm) {
            return true
        }
    }
    
    // Match common field VALUES only
    for _, campo := range s.campos {
        if campo.tipo == TipoCampoComum {
            if strings.Contains(normalizar(string(campo.valor)), critNorm) {
                return true
            }
        }
        // Sensitive field VALUES: skip
    }
    
    // Match observation VALUE
    if strings.Contains(normalizar(string(s.observacao.valor)), critNorm) {
        return true
    }
    
    return false
}

// Manager method: apply search policy
func (m *Manager) BuscarSegredos(query string) []*Segredo {
    if query == "" {
        return nil
    }
    
    resultados := []*Segredo{}
    m.buscarRecursivamente(m.cofre.pastaGeral, query, &resultados)
    return resultados
}

func (m *Manager) buscarRecursivamente(pasta *Pasta, query string, resultados *[]*Segredo) {
    // Search secrets in this folder
    for _, segredo := range pasta.segredos {
        // Filter out excluido (D-14)
        if segredo.estadoSessao == EstadoExcluido {
            continue
        }
        
        // Delegate matching to entity
        if segredo.AtendeCriterio(query) {
            *resultados = append(*resultados, segredo)
        }
    }
    
    // Recurse into subfolders
    for _, subpasta := range pasta.subpastas {
        m.buscarRecursivamente(subpasta, query, resultados)
    }
}
```

**Test case (UAT):** "Search function called with a string present only in a FieldTypeSensitive field **value** returns zero results; search called with the **name** of a sensitive field returns secrets containing that field"

```go
// Test: search for value in sensitive field returns nothing
func TestSearchSensitiveFieldValue(t *testing.T) {
    manager := setupManager()
    segredo := manager.CriarSegredoDeModelo(pastaGeral, "Login", "GitHub", 0)
    manager.AlterarCampoSegredo(segredo, 2, []byte("my_secret_password_123"))  // Senha field
    
    // Search for password value - should return nothing
    resultados := manager.BuscarSegredos("my_secret_password_123")
    assert.Empty(t, resultados)
    
    // Search for field name "Senha" - should return secret
    resultados = manager.BuscarSegredos("Senha")
    assert.Len(t, resultados, 1)
    assert.Equal(t, segredo, resultados[0])
}
```

### Favoritos DFS Traversal (D-21)

**Requirement:** List favorites in depth-first order following JSON structure.

**Implementation:**
```go
func (m *Manager) ListarFavoritos() []*Segredo {
    favoritos := []*Segredo{}
    m.coletarFavoritosRecursivamente(m.cofre.pastaGeral, &favoritos)
    return favoritos
}

func (m *Manager) coletarFavoritosRecursivamente(pasta *Pasta, favoritos *[]*Segredo) {
    // Process secrets in current folder first (DFS)
    for _, segredo := range pasta.segredos {
        // Filter excluido (D-21)
        if segredo.estadoSessao == EstadoExcluido {
            continue
        }
        
        if segredo.favorito {
            *favoritos = append(*favoritos, segredo)
        }
    }
    
    // Then recurse into subfolders
    for _, subpasta := range pasta.subpastas {
        m.coletarFavoritosRecursivamente(subpasta, favoritos)
    }
}
```

**Complexity:** O(n) where n = total secrets. No caching in Phase 3 (D-21) - simple on-demand traversal. Future optimization possible if profiling reveals bottleneck.

## Testing Strategy

### Test Organization

**File structure:**
```
internal/vault/
  entities.go           # All entity types
  manager.go            # Manager implementation
  errors.go             # Sentinel errors
  doc.go               # Package documentation
  entities_test.go      # Entity unit tests
  manager_test.go       # Manager integration tests
  validation_test.go    # Business rule tests
  state_machine_test.go # State transition tests
```

### Critical Test Cases (from ROADMAP UAT)

1. **Pasta Geral protection:** All mutating operations return `ErrPastaGeralProtected`
2. **Cycle detection:** `MoveFolder` to descendant returns `ErrCycleDetected`
3. **Name uniqueness:** Create/Rename with existing name returns `ErrNameConflict`
4. **Observação invariant:**
   - Auto-created on `CreateSecret`
   - Absent from `UpdateSecretStructure` operations
   - Immutable (cannot rename/delete/move)
5. **State machine:**
   - `CreateSecret` → `EstadoIncluido`
   - `UpdateSecret(original)` → `EstadoModificado`
   - `UpdateSecret(incluido)` → stays `EstadoIncluido`
   - `SoftDeleteSecret` → `EstadoExcluido` (stores previous state)
   - `RestoreSecret` → restores previous state
6. **Duplication:** "X" → "X (1)" → "X (2)"
7. **CreateTemplateFromSecret:** Excludes ALL fields named 'Observação'
8. **UpdateTemplateStructure:** Returns error when adding/renaming field to 'Observação'
9. **Search QUERY-02:** Sensitive field VALUE excluded, sensitive field NAME included
10. **Template sorting (TPL-06):** Always alphabetical regardless of creation order

### Race Detector (from Phase 2 pattern)

**Command:** `go test -race ./internal/vault/...`

Must pass with no race conditions. Critical for Manager methods that might be called concurrently in future phases (timer-based auto-lock).

### Table-Driven Tests

**Pattern from Go community:**
```go
func TestRenomearSegredo(t *testing.T) {
    tests := []struct {
        nome        string
        estadoInicial EstadoSessao
        nomeNovo    string
        esperado    EstadoSessao
        erro        error
    }{
        {"original sem mudança", EstadoOriginal, "mesmo", EstadoOriginal, nil},
        {"original com mudança", EstadoOriginal, "novo", EstadoModificado, nil},
        {"incluido com mudança", EstadoIncluido, "novo", EstadoIncluido, nil},
        {"nome vazio", EstadoOriginal, "", EstadoOriginal, ErrNomeVazio},
    }
    
    for _, tt := range tests {
        t.Run(tt.nome, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Error Handling

### Sentinel Errors (D-07, from Phase 2 pattern)

**Pattern:** Simple sentinel errors for common validation failures.

```go
// internal/vault/errors.go
package vault

import "errors"

var (
    ErrNomeVazio            = errors.New("nome não pode ser vazio")
    ErrNomeMuitoLongo       = errors.New("nome excede 255 caracteres")
    ErrNameConflict         = errors.New("já existe item com este nome")
    ErrPastaGeralProtected  = errors.New("Pasta Geral não pode ser modificada")
    ErrCycleDetected        = errors.New("operação criaria ciclo na hierarquia")
    ErrObservacaoReserved   = errors.New("nome 'Observação' é reservado")
    ErrConfigInvalida       = errors.New("configuração inválida")
    ErrPosicaoInvalida      = errors.New("posição inválida")
    ErrSegredoNaoEncontrado = errors.New("segredo não encontrado")
    ErrCofreBloqueado       = errors.New("cofre está bloqueado")
)
```

### Custom Error Structs (D-07)

**When TUI needs structured data from error:**
```go
type RenomeacoesError struct {
    Renomeacoes []Renomeacao
}

func (e *RenomeacoesError) Error() string {
    return fmt.Sprintf("%d itens foram renomeados automaticamente", len(e.Renomeacoes))
}

type Renomeacao struct {
    Antigo string
    Novo   string
    Pasta  string
}
```

TUI can type-assert to extract details for user display.

## Integration Points

### Storage Interface (Phase 4)

**Definition in Phase 3, implementation in Phase 4:**
```go
// internal/vault/repository.go
type RepositorioCofre interface {
    Salvar(cofre *Cofre) error
    Carregar() (*Cofre, error)
}
```

Manager receives repository via dependency injection:
```go
func NewManager(cofre *Cofre, repositorio RepositorioCofre) *Manager {
    return &Manager{
        cofre: cofre,
        repositorio: repositorio,
        bloqueado: false,
    }
}
```

### Crypto Package (from Phase 2)

**Available for Phase 3:**
```go
import "github.com/useful-toys/abditum/internal/crypto"

// Memory wiping when locking
func (m *Manager) Lock() {
    if m.senha != nil {
        crypto.WipeBytes(m.senha)
        m.senha = nil
    }
    
    // Wipe all sensitive field values
    m.limparCamposSensiveis(m.cofre.pastaGeral)
    
    m.cofre = nil
    m.bloqueado = true
}

func (m *Manager) limparCamposSensiveis(pasta *Pasta) {
    for _, segredo := range pasta.segredos {
        for i := range segredo.campos {
            if segredo.campos[i].tipo == TipoCampoSensivel {
                crypto.WipeBytes(segredo.campos[i].valor)
            }
        }
        crypto.WipeBytes(segredo.observacao.valor)
    }
    
    for _, subpasta := range pasta.subpastas {
        m.limparCamposSensiveis(subpasta)
    }
}
```

## Pitfalls and Gotchas

### Pitfall 1: NanoID No Longer Needed (D-18)

**ROADMAP says:** "NanoID must use crypto/rand"

**CONTEXT.md D-01 eliminates:** No synthetic IDs needed. Go pointers provide identity within session. JSON uses hierarchical structure without ID fields.

**Action:** Plans should NOT include NanoID generation. Update ROADMAP understanding during planning.

### Pitfall 2: Manager.Vault() Snapshot (ROADMAP vs D-08/D-09)

**ROADMAP says:** "Manager.Vault() must return a snapshot — exposing a live pointer to domain state allows TUI bugs to corrupt the domain silently"

**CONTEXT.md D-08/D-09 refines:** Manager returns live `*Cofre` pointer. Safety via package encapsulation (all fields lowercase) + defensive copies in getters. Zero copy overhead.

**Rationale:** Package-level encapsulation is stronger than deep copying. TUI literally cannot mutate private fields even with live pointer. Getters return defensive slice copies to prevent indirect mutation via `append()`.

**Action:** Plans should implement package-level encapsulation, NOT full snapshot deep copy on every `Vault()` call.

### Pitfall 3: Cycle Detection Must Walk Full Ancestor Chain (ROADMAP)

**ROADMAP says:** "Cycle detection must walk the FULL ancestor chain, not just immediate parent"

**Implementation:** Walk from destination up to root, comparing each ancestor with source folder.

**Edge cases:**
- Immediate parent: no cycle
- Grandparent: no cycle
- Own child: cycle ✓
- Grandchild: cycle ✓
- Distant descendant: cycle ✓

### Pitfall 4: StateDeleted Secrets Retain State When Promoted (FOLDER-05)

**ROADMAP FOLDER-05:** "including segredos with state `StateDeleted`, which retain their `StateDeleted` state"

**Implication:** When deleting folder, promoted secrets marked `StateDeleted` stay marked. They're still visible with strikethrough in parent folder. Final removal happens only on save.

### Pitfall 5: Go Slice Mutation via Defensive Copy Bypass

**Problem:** Even with defensive copies, if slice contains pointers, caller can mutate pointed-to objects.

**Solution:** All entity getters return slices of pointers (`[]*Pasta`, `[]*Segredo`). Slices are copied, but pointers remain. However, entities are encapsulated (lowercase fields), so TUI cannot mutate even with pointer access. Type system prevents corruption.

**Example:**
```go
// TUI code
modelos := cofre.Modelos()  // Defensive copy of slice
modelo := modelos[0]        // Pointer from slice
modelo.nome = "hack"        // COMPILE ERROR: nome is unexported
```

### Pitfall 6: Timestamp Update on No-Op Operations (D-12, D-24)

**Problem:** User opens rename dialog, doesn't change anything, confirms. Should this mark as modified?

**Solution:** Entity methods return `(alterado bool, err error)`. Manager only updates flags/timestamps if `alterado == true`. Prevents false positives.

### Pitfall 7: Favoriting vs EstadoSessao Confusion (D-11)

**CONTEXT overrides modelo-dominio.md line 73:** Favoriting does NOT change `segredo.estadoSessao`. Only changes `cofre.modificado`.

**Rationale:** Favoriting is navigation preference (metadata), not content edit. User doesn't see "modified" indicator on secret when favoriting.

**Tests must verify:** Favorite/unfavorite a `EstadoOriginal` secret → `estadoSessao` stays `EstadoOriginal`, but `cofre.modificado` becomes true.

### Pitfall 8: Position Semantics in Insertion (D-22)

**Position parameter:** 0-indexed. `posicao == len(slice)` means append at end. `posicao < len(slice)` means insert here, shift right.

**Edge cases:**
- `posicao = 0, len = 0`: append (OK)
- `posicao = 0, len > 0`: insert at start (OK)
- `posicao = len`: append (OK)
- `posicao > len`: error
- `posicao < 0`: error

**Implementation:**
```go
// Insert at position
slice = append(slice[:pos], append([]*T{item}, slice[pos:]...)...)
```

This idiom handles both insert and append naturally.

### Pitfall 9: Observação Exclusion in CreateTemplateFromSecret (TPL-05, UAT)

**Requirement:** "Excludes ALL fields named 'Observação' — both the auto-Observation and any user-named field with that name"

**Edge case:** User manually added a common field named "Observação" to a secret (before it was forbidden in templates). When creating template from this secret, BOTH the auto-Observation AND the user field must be excluded.

**Implementation:**
```go
func (m *Manager) CriarModeloDeSegredo(segredo *Segredo, nomeModelo string) (*ModeloSegredo, error) {
    // ... validation ...
    
    campos := []CampoModelo{}
    for _, campo := range segredo.campos {
        // Exclude any field named "Observação"
        if campo.nome == "Observação" {
            continue
        }
        campos = append(campos, CampoModelo{
            nome: campo.nome,
            tipo: campo.tipo,
        })
    }
    
    // Note: observacao field is separate, never in campos slice,
    // so it's already excluded by not iterating over it
    
    return m.CriarModelo(nomeModelo, campos)
}
```

## External Dependencies

### Standard Library Only (for Phase 3)

**No external dependencies needed beyond stdlib:**
- `errors` - sentinel errors
- `fmt` - string formatting
- `sort` - template sorting
- `strings` - string manipulation
- `time` - timestamps
- `golang.org/x/text` - text normalization for search (already in go.mod from Phase 2)

**Future phases:**
- Phase 4: `encoding/json` for serialization
- Phase 2 (already complete): `golang.org/x/crypto` for Argon2id/AES-GCM

### No NanoID Library Needed (D-18)

ROADMAP mentioned `go-nanoid/v2` but D-01 eliminated synthetic ID requirements. Don't add this dependency.

## API Surface

### Manager Public Methods (orchestration layer)

**Lifecycle:**
- `func NewManager(cofre *Cofre, repo RepositorioCofre) *Manager`
- `func (m *Manager) Salvar() error`
- `func (m *Manager) Lock()`
- `func (m *Manager) IsLocked() bool`
- `func (m *Manager) IsModified() bool`
- `func (m *Manager) Vault() *Cofre`

**Folder operations:**
- `func (m *Manager) CriarPasta(pai *Pasta, nome string, pos int) (*Pasta, error)`
- `func (m *Manager) RenomearPasta(pasta *Pasta, nome string) error`
- `func (m *Manager) MoverPasta(pasta, destino *Pasta) error`
- `func (m *Manager) ReposicionarPasta(pasta *Pasta, pos int) error`
- `func (m *Manager) SubirPastaNaPosicao(pasta *Pasta) error`
- `func (m *Manager) DescerPastaNaPosicao(pasta *Pasta) error`
- `func (m *Manager) ExcluirPasta(pasta *Pasta) ([]Renomeacao, error)`

**Secret operations:**
- `func (m *Manager) CriarSegredo(pasta *Pasta, nome string, campos []CampoSegredo, pos int) (*Segredo, error)`
- `func (m *Manager) CriarSegredoDeModelo(pasta *Pasta, modelo *ModeloSegredo, nome string, pos int) (*Segredo, error)`
- `func (m *Manager) DuplicarSegredo(segredo *Segredo) (*Segredo, error)`
- `func (m *Manager) RenomearSegredo(segredo *Segredo, nome string) error`
- `func (m *Manager) AlterarCampoSegredo(segredo *Segredo, indice int, valor []byte) error`
- `func (m *Manager) AlterarObservacao(segredo *Segredo, valor string) error`
- `func (m *Manager) AdicionarCampoSegredo(segredo *Segredo, nome, tipo string, valor []byte, pos int) error`
- `func (m *Manager) RenomearCampoSegredo(segredo *Segredo, indice int, nome string) error`
- `func (m *Manager) RemoverCampoSegredo(segredo *Segredo, indice int) error`
- `func (m *Manager) ReordenarCampoSegredo(segredo *Segredo, indice, pos int) error`
- `func (m *Manager) MoverSegredo(segredo *Segredo, destino *Pasta) error`
- `func (m *Manager) ReposicionarSegredo(segredo *Segredo, pos int) error`
- `func (m *Manager) SubirSegredoNaPosicao(segredo *Segredo) error`
- `func (m *Manager) DescerSegredoNaPosicao(segredo *Segredo) error`
- `func (m *Manager) FavoritarSegredo(segredo *Segredo, fav bool) error`
- `func (m *Manager) ExcluirSegredo(segredo *Segredo) error`
- `func (m *Manager) RestaurarSegredo(segredo *Segredo) error`

**Template operations:**
- `func (m *Manager) CriarModelo(nome string, campos []CampoModelo) (*ModeloSegredo, error)`
- `func (m *Manager) CriarModeloDeSegredo(segredo *Segredo, nome string) (*ModeloSegredo, error)`
- `func (m *Manager) RenomearModelo(modelo *ModeloSegredo, nome string) error`
- `func (m *Manager) AdicionarCampoModelo(modelo *ModeloSegredo, nome, tipo string, pos int) error`
- `func (m *Manager) RenomearCampoModelo(modelo *ModeloSegredo, indice int, nome string) error`
- `func (m *Manager) AlterarTipoCampoModelo(modelo *ModeloSegredo, indice int, tipo string) error`
- `func (m *Manager) RemoverCampoModelo(modelo *ModeloSegredo, indice int) error`
- `func (m *Manager) ReordenarCampoModelo(modelo *ModeloSegredo, indice, pos int) error`
- `func (m *Manager) ExcluirModelo(modelo *ModeloSegredo) error`

**Query operations:**
- `func (m *Manager) BuscarSegredos(query string) []*Segredo`
- `func (m *Manager) ListarFavoritos() []*Segredo`

**Configuration:**
- `func (m *Manager) AlterarConfiguracoes(config Configuracoes) error`

### Entity Exported Getters (navigation layer)

**Cofre:**
- `func (c *Cofre) PastaGeral() *Pasta`
- `func (c *Cofre) Modelos() []*ModeloSegredo` (defensive copy, sorted)
- `func (c *Cofre) Configuracoes() Configuracoes`
- `func (c *Cofre) Modificado() bool`
- `func (c *Cofre) DataCriacao() time.Time`
- `func (c *Cofre) DataUltimaModificacao() time.Time`

**Pasta:**
- `func (p *Pasta) Nome() string`
- `func (p *Pasta) Pai() *Pasta`
- `func (p *Pasta) Subpastas() []*Pasta` (defensive copy)
- `func (p *Pasta) Segredos() []*Segredo` (defensive copy, includes excluido)

**Segredo:**
- `func (s *Segredo) Nome() string`
- `func (s *Segredo) Pasta() *Pasta`
- `func (s *Segredo) Campos() []CampoSegredo` (defensive copy, excludes Observação)
- `func (s *Segredo) Observacao() string`
- `func (s *Segredo) Favorito() bool`
- `func (s *Segredo) EstadoSessao() EstadoSessao`
- `func (s *Segredo) DataCriacao() time.Time`
- `func (s *Segredo) DataUltimaModificacao() time.Time`

**CampoSegredo:**
- `func (c *CampoSegredo) Nome() string`
- `func (c *CampoSegredo) Tipo() TipoCampo`
- `func (c *CampoSegredo) ValorComoString() string`

**ModeloSegredo:**
- `func (m *ModeloSegredo) Nome() string`
- `func (m *ModeloSegredo) Campos() []CampoModelo` (defensive copy)

**CampoModelo:**
- `func (c *CampoModelo) Nome() string`
- `func (c *CampoModelo) Tipo() TipoCampo`

## Validation Architecture

### Nyquist Rule Compliance

**Every Manager method that creates/modifies state will have corresponding test:**
- Unit tests for validation logic (entity-level)
- Integration tests for workflows (Manager-level)
- State machine tests for transitions
- Business rule tests for invariants

**Test files:**
- `manager_test.go` - Manager method integration tests
- `entities_test.go` - Entity validation and mutation tests
- `state_machine_test.go` - EstadoSessao transition tests
- `validation_test.go` - Business rule tests (cycles, uniqueness, Pasta Geral protection)

**All tests run with race detector:** `go test -race ./internal/vault/...`

### Automated Verification Commands

**For each Manager method, tests verify:**
```bash
# All tests pass
go test ./internal/vault/... -v

# No race conditions
go test ./internal/vault/... -race

# High coverage (aim for >90% on logic paths)
go test ./internal/vault/... -cover
```

## Planning Recommendations

### Plan Breakdown Suggestion

Based on complexity and dependencies:

**Plan 03-01: Domain Entities + Factory (Foundation)**
- Define all entity structs (`Cofre`, `Pasta`, `Segredo`, `ModeloSegredo`, `CampoSegredo`, `Configuracoes`)
- Define enums (`EstadoSessao`, `TipoCampo`)
- Implement all exported getters with defensive copies
- Implement `NovoCofre()` factory
- Implement `Cofre.InicializarConteudoPadrao()` bootstrap
- Implement sentinel errors in `errors.go`
- No Manager yet, no mutations (read-only foundation)
- Tests: Entity construction, getter defensive copies, initial content

**Plan 03-02: Manager + Cofre Lifecycle**
- Implement `Manager` struct
- Implement `NewManager()`, `Salvar()`, `Lock()`, `IsLocked()`, `IsModified()`, `Vault()`
- Define `RepositorioCofre` interface (no implementation yet)
- Implement `Configuracoes` mutation
- Tests: Manager creation, lock/unlock, configuration changes

**Plan 03-03: Folder Management**
- Implement all folder Manager methods (`CriarPasta`, `RenomearPasta`, `MoverPasta`, `ReposicionarPasta`, etc.)
- Implement folder entity private methods (validation + mutation)
- Implement cycle detection algorithm
- Implement folder deletion with promotion and conflict resolution
- Tests: Cycle detection, name uniqueness, Pasta Geral protection, deletion promotion, rename suffixing

**Plan 03-04: Template Management**
- Implement all template Manager methods
- Implement template sorting (TPL-06)
- Implement "Observação" name prohibition
- Tests: Template CRUD, alphabetical sorting, reserved name rejection, CreateTemplateFromSecret exclusion

**Plan 03-05: Secret Lifecycle + State Machine**
- Implement secret creation methods (`CriarSegredo`, `CriarSegredoDeModelo`, `DuplicarSegredo`)
- Implement secret deletion/restoration (`ExcluirSegredo`, `RestaurarSegredo`)
- Implement favoriting
- Implement state machine transitions (D-13)
- Tests: State transitions, duplication name progression, deletion/restoration, favoriting doesn't change estadoSessao

**Plan 03-06: Secret CRUD + Structure**
- Implement secret content mutation (`RenomearSegredo`, `AlterarCampoSegredo`, `AlterarObservacao`)
- Implement secret structure mutation (`AdicionarCampoSegredo`, `RenomearCampoSegredo`, `RemoverCampoSegredo`, `ReordenarCampoSegredo`)
- Implement move/reposition
- Tests: Content changes update estadoSessao, Observação immutability, field manipulation, move updates timestamp

**Plan 03-07: Search + Favorites + Comprehensive Validation**
- Implement `BuscarSegredos()` with normalization and sensitive field exclusion (QUERY-02)
- Implement `ListarFavoritos()` with DFS traversal
- Implement atomic save with two-phase commit (D-17)
- Comprehensive validation test suite covering all ROADMAP UAT criteria
- Race detector validation
- Tests: Search sensitivity, favorite traversal, save atomicity, full UAT coverage

### Key Integration Points for Planners

**Dependencies Phase 2 → Phase 3:**
- `internal/crypto` package available for memory wiping on Lock()
- Sentinel error pattern established
- Testing with race detector established

**Dependencies Phase 3 → Phase 4:**
- `RepositorioCofre` interface defined
- `Salvar()` expects JSON serialization implementation
- Reference reconstruction after deserialization (`popularReferencias()`)

**Dependencies Phase 3 → Phase 5+:**
- Manager public API is TUI's complete interface
- Entity getters provide navigation
- No direct field access from TUI

## Research Complete

**Key Takeaways:**
1. CONTEXT.md decisions (D-01 through D-30) override some ROADMAP assumptions (NanoID, snapshot copies)
2. Pointer-based identity with package-level encapsulation is sufficient for safety
3. Two independent state flags (cofre.modificado vs segredo.estadoSessao) with change detection
4. Structural enforcement (Observação separate from campos slice) eliminates runtime checks
5. Two-phase validation pattern (validate → mutate) ensures atomicity
6. Cycle detection requires full ancestor walk
7. Atomic save uses two-phase commit (snapshot → persist → finalize)
8. Search must exclude sensitive field VALUES but include sensitive field NAMES
9. Templates always sorted alphabetically (TPL-06)
10. Comprehensive test suite required covering all ROADMAP UAT criteria + race detector

**Files to create:**
- `internal/vault/entities.go`
- `internal/vault/manager.go`
- `internal/vault/errors.go`
- `internal/vault/repository.go` (interface only)
- `internal/vault/entities_test.go`
- `internal/vault/manager_test.go`
- `internal/vault/state_machine_test.go`
- `internal/vault/validation_test.go`

**Critical specs to reference:**
- `.planning/phases/03-vault-domain-manager/03-CONTEXT.md` (30 decisions)
- `modelo-dominio.md` (domain model details)
- `arquitetura-camada-dominio.md` (architecture patterns)
- `.planning/REQUIREMENTS.md` (VAULT-02, SEC-05, FOLDER-01 through 05, TPL-01 through 06)
- `.planning/ROADMAP.md` Phase 3 (UAT criteria)
