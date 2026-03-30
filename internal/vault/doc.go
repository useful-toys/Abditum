// Package vault implements the domain layer for Abditum password manager.
//
// Domain Model:
//   - Cofre (vault): Root aggregate containing folders, templates, and configuration
//   - Pasta (folder): Hierarchical containers for subfolders and secrets
//   - Segredo (secret): Password entries with fields and observation
//   - ModeloSegredo (template): Reusable field definitions for secrets
//   - Configuracoes (settings): Timer values for auto-lock, reveal, clipboard
//
// Encapsulation:
// All entity fields are package-private (lowercase). External access via exported
// getters returning defensive copies. TUI interacts only through Manager public API.
//
// Identity:
// Entities use Go pointer identity during session (no synthetic IDs needed).
// Uniqueness validated via composite keys: (parent, nome) for folders/secrets,
// global nome for templates.
//
// State Tracking:
// Two independent flags:
//   - cofre.modificado: ANY mutation (including favoriting)
//   - segredo.estadoSessao: Content changes only (NOT favoriting)
//
// Manager Pattern:
// Manager orchestrates operations (knows WHAT and WHY), entities own validation
// and mutation logic (know HOW). Two-phase pattern: validate (can fail) →
// mutate (cannot fail after validation passes).
package vault
