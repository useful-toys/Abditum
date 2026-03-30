// Package vault implements the domain layer for Abditum password manager.
//
// # Domain Model
//
// The vault domain consists of five core entity types:
//   - Cofre (vault): Root aggregate containing folders, templates, and configuration
//   - Pasta (folder): Hierarchical containers for subfolders and secrets
//   - Segredo (secret): Password entries with fields and observation
//   - ModeloSegredo (template): Reusable field definitions for secrets
//   - Configuracoes (settings): Timer values for auto-lock, reveal, clipboard
//
// # Encapsulation
//
// All entity fields are package-private (lowercase). External access via exported
// getters returning defensive copies. TUI interacts only through Manager public API.
// This prevents accidental corruption of domain state from presentation layer bugs.
//
// # Identity (D-01)
//
// Entities use Go pointer identity during session (no synthetic IDs needed).
// Uniqueness validated via composite keys: (parent, nome) for folders/secrets,
// global nome for templates. Pointers are stable for the session lifetime.
//
// # State Tracking (D-11)
//
// Two independent flags:
//   - cofre.modificado: ANY mutation (including favoriting)
//   - segredo.estadoSessao: Content changes only (NOT favoriting)
//
// EstadoSessao transitions (D-13):
//   - EstadoOriginal: Loaded from disk, unmodified
//   - EstadoModificado: Content changed (rename, edit field, etc.)
//   - EstadoExcluido: Soft deleted, pending finalization on save
//
// # Manager Pattern (D-05)
//
// Manager orchestrates operations (knows WHAT and WHY), entities own validation
// and mutation logic (know HOW). Two-phase pattern: validate (can fail) →
// mutate (cannot fail after validation passes).
//
// Manager responsibilities:
//   - Enforce business rules (e.g., lock state, unique names)
//   - Update timestamps and modified flags
//   - Coordinate cross-entity operations (e.g., folder deletion with promotion)
//
// Entity responsibilities:
//   - Validate operation preconditions
//   - Execute state mutations
//   - Maintain invariants (e.g., Observação always exists)
//
// # Atomic Save (D-17)
//
// Save uses two-phase commit:
//  1. prepararSnapshot(): Create deep copy with EstadoExcluido filtered out
//  2. repository.Salvar(snapshot): Write to disk
//
// If save fails, live vault remains untouched. If save succeeds, finalizarExclusoes()
// permanently removes EstadoExcluido secrets from memory.
//
// # Memory Security (D-29)
//
// Lock() wipes all sensitive field values (tipo == TipoCampoSensivel) and master
// password from memory using ZeroBytes(). Observation field always wiped on lock
// regardless of type.
//
// # Search Policy (D-19, QUERY-02)
//
// Buscar() searches:
//   - Secret name (segredo.nome)
//   - Field NAMES (all types, including sensitive)
//   - Field VALUES (common fields only, excludes sensitive)
//   - Observation VALUE
//
// Case-insensitive using strings.ToLower(). Excludes EstadoExcluido secrets.
//
// # Favorites Traversal (D-20)
//
// ListarFavoritos() uses depth-first search (DFS) starting from PastaGeral.
// Returns secrets where favorito == true AND estadoSessao != EstadoExcluido.
//
// # Usage Example
//
//	// Create vault
//	cofre := vault.NovoCofre()
//	cofre.InicializarConteudoPadrao() // Adds default folders and templates
//	manager := vault.NewManager(cofre, repository)
//
//	// Create secret
//	modelo := cofre.Modelos()[0] // "Login" template
//	pasta := cofre.PastaGeral()
//	secret, err := manager.CriarSegredo(pasta, "GitHub", modelo)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Edit field
//	err = manager.EditarCampoSegredo(secret, 0, []byte("mypassword"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Save atomically
//	err = manager.Salvar()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Lock (wipe sensitive data from memory)
//	manager.Lock()
//
// # Key Design Decisions
//
// D-01: Go pointers as identity (no synthetic IDs)
// D-05: Manager pattern (orchestration vs execution split)
// D-11: Two independent flags (modificado vs estadoSessao)
// D-17: Atomic save with two-phase commit
// D-19: Case-insensitive search with strings.ToLower()
// D-20: Favorites use DFS traversal
// D-27: Pasta hard delete (immediate removal), Segredo soft delete (session-scoped)
// D-29: Lock wipes sensitive data from memory
//
// See .planning/phases/03-vault-domain-manager/03-CONTEXT.md for all 30 decisions.
package vault
