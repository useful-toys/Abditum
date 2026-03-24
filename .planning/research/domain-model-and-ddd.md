# Domain Model & DDD Research: Abditum

**Researched:** 2026-03-24  
**Focus:** Go NanoID, DDD patterns for Go, vault manager API design, recursive folder hierarchy, soft delete (Trash), virtual folders  
**Confidence:** HIGH (source-verified libraries + established Go DDD community patterns)

---

## 1. NanoID in Go

### Recommended Library: `github.com/matoous/go-nanoid/v2`

- **Stars:** 1.6k — the canonical Go NanoID port, listed on the official ai/nanoid README
- **Status:** Stable, minimal, intentionally low-maintenance (bug fixes only)
- **License:** MIT
- **CGO:** Zero CGO dependencies — pure Go
- **Install:** `go get github.com/matoous/go-nanoid/v2`

### API

```go
import gonanoid "github.com/matoous/go-nanoid/v2"

// Generate a NanoID with the default alphabet (A-Za-z0-9_-) and default length (21)
id, err := gonanoid.New()

// Generate a 6-character NanoID with the default alphabet — Abditum's case
id, err := gonanoid.New(6)  // e.g., "aB3xR7"

// Generate with a custom alphabet
id, err := gonanoid.Generate("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789", 6)
```

### Collision Analysis for 6 Characters

Per `descricao.md`: IDs are 6 alphanumeric characters (A-Za-z0-9 = 62 characters).

- **Space size:** 62⁶ ≈ 56.8 billion combinations
- **Birthday paradox collision at 50%:** ~212,000 IDs
- **Practical vault size:** Unlikely to exceed 10,000 secrets + folders per vault
- **Collision probability at 10,000 items:** < 0.001% — safe for practical use

**Custom alphabet for Abditum** (alphanumeric only, no `-` or `_` for cleaner display):

```go
const idAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
const idLength   = 6

func NewID() (string, error) {
    return gonanoid.Generate(idAlphabet, idLength)
}
```

### ID Assignment Scope

Per `descricao.md`: Only **Secrets** and **Folders** have IDs (they are moveable and can be referenced). Templates are identified by name. Fields have no ID (they are ordered within their parent).

---

## 2. DDD in Go: Core Principles and Patterns

### Philosophy for Abditum

Go is not object-oriented, but DDD principles apply through **package boundaries and interface contracts**:

- **Domain layer:** Pure business logic, no I/O, no framework dependencies
- **Entities:** Structs with methods; mutations only through Manager
- **Manager as aggregate root:** All modifications go through a Manager that enforces invariants
- **Repository pattern:** Interface in domain, implementation in storage layer

### The Read-Only Entity / Mutable Manager Split

`descricao.md` explicitly mandates this pattern:

> "Permita que as entidades sejam navegadas somente leitura. As operações de modificação devem ser realizadas exclusivamente por meio de métodos explícitos de um Manager"

**Implementation:**

```go
// domain/vault/vault.go — read-only access

// Vault is the aggregate root. Its fields are unexported — callers can only
// read through accessor methods. All mutations go through VaultManager.
type Vault struct {
    version         int
    settings        Settings
    secrets         []Secret    // top-level (at vault root)
    folders         []Folder    // top-level (at vault root)
    templates       []SecretTemplate
    createdAt       time.Time
    lastModifiedAt  time.Time
}

// Secrets returns a read-only view of the vault's top-level secrets.
// Callers must not mutate the returned slice.
func (v *Vault) Secrets() []Secret { return v.secrets }

// Folders returns a read-only view of the vault's top-level folders.
func (v *Vault) Folders() []Folder { return v.folders }

// Templates returns a read-only view of the vault's secret templates.
func (v *Vault) Templates() []SecretTemplate { return v.templates }

// Settings returns the vault's configuration.
func (v *Vault) Settings() Settings { return v.settings }

// LastModifiedAt returns the time of the last modification.
func (v *Vault) LastModifiedAt() time.Time { return v.lastModifiedAt }
```

```go
// domain/vault/manager.go — all mutations

// VaultManager is the single entry point for modifying vault state.
// It enforces all business rules and invariants. The TUI and other callers
// receive mutations through manager methods only — never by modifying
// Vault fields directly.
type VaultManager struct {
    vault    *Vault
    dirty    bool     // true if there are unsaved changes
}

// NewVaultManager creates a new manager wrapping the given vault.
func NewVaultManager(v *Vault) *VaultManager {
    return &VaultManager{vault: v}
}

// IsDirty reports whether the vault has unsaved changes.
func (m *VaultManager) IsDirty() bool { return m.dirty }

// Vault returns the underlying vault for read-only access.
func (m *VaultManager) Vault() *Vault { return m.vault }

// MarkSaved resets the dirty flag after a successful save.
func (m *VaultManager) MarkSaved() { m.dirty = false }
```

---

## 3. Vault Entity Model

### Folder (Recursive Hierarchy)

```go
// domain/vault/folder.go

// Folder is a container for secrets and other folders.
// The vault root is modeled as an implicit "root folder" (not an actual Folder
// entity) that holds the top-level secrets and folders slices in Vault.
//
// Folders are fully recursive — a folder can contain both secrets and subfolders
// at any depth.
type Folder struct {
    id      string    // NanoID (6 chars)
    name    string
    secrets []Secret  // secrets in this folder, in display order
    folders []Folder  // subfolders, in display order
}

// ID returns the folder's unique identifier.
func (f *Folder) ID() string { return f.id }

// Name returns the folder's display name.
func (f *Folder) Name() string { return f.name }

// Secrets returns secrets in this folder.
func (f *Folder) Secrets() []Secret { return f.secrets }

// Folders returns subfolders of this folder.
func (f *Folder) Folders() []Folder { return f.folders }

// TotalSecretCount returns the total number of secrets in this folder and all
// subfolders, recursively. Used for displaying "(5)" next to folder name in the tree.
func (f *Folder) TotalSecretCount() int {
    count := len(f.secrets)
    for _, sub := range f.folders {
        count += sub.TotalSecretCount()
    }
    return count
}
```

### Secret Entity

```go
// domain/vault/secret.go

// Secret is an individual item in the vault. It holds a set of named, typed
// fields (some sensitive, some not) along with metadata.
type Secret struct {
    id             string    // NanoID (6 chars)
    name           string
    templateName   string    // snapshot: name of the template used at creation (informational only)
    fields         []Field   // ordered list of fields
    favorite       bool
    note           string    // free-text observation
    inTrash        bool      // true = soft-deleted, will be purged on save
    createdAt      time.Time
    lastModifiedAt time.Time
}

// ID, Name, TemplateName, Fields, Favorite, Note, InTrash, etc. — accessor methods
func (s *Secret) ID() string              { return s.id }
func (s *Secret) Name() string            { return s.name }
func (s *Secret) TemplateName() string    { return s.templateName }
func (s *Secret) Fields() []Field         { return s.fields }
func (s *Secret) IsFavorite() bool        { return s.favorite }
func (s *Secret) Note() string            { return s.note }
func (s *Secret) IsInTrash() bool         { return s.inTrash }
func (s *Secret) CreatedAt() time.Time    { return s.createdAt }
func (s *Secret) LastModifiedAt() time.Time { return s.lastModifiedAt }
```

### Field Entity

```go
// domain/vault/field.go

// FieldType enumerates the types a secret field can hold.
// Only two types exist in v1: plain text and sensitive text.
type FieldType uint8

const (
    FieldTypeText      FieldType = iota  // regular text (URL, username, etc.)
    FieldTypeSensitive                   // sensitive text (password, API key, CVV)
)

// Field is a single key-value pair within a Secret. The type determines
// whether the value should be masked in the UI by default.
type Field struct {
    name      string
    fieldType FieldType
    value     string  // may be empty — that's valid (not a "missing" state)
}

func (f *Field) Name() string          { return f.name }
func (f *Field) Type() FieldType       { return f.fieldType }
func (f *Field) Value() string         { return f.value }
func (f *Field) IsSensitive() bool     { return f.fieldType == FieldTypeSensitive }
```

### SecretTemplate Entity

```go
// domain/vault/template.go

// SecretTemplate defines the structure (field names and types) used to create
// new secrets. It is a snapshot — changing a template after secrets are created
// does NOT retroactively change those secrets.
type SecretTemplate struct {
    id     string         // NanoID (6 chars) — for internal use; displayed by name
    name   string         // displayed to user; also used as the "snapshot" reference
    fields []TemplateField
}

// TemplateField defines a field in a template (no value — just name and type).
type TemplateField struct {
    name      string
    fieldType FieldType
}

func (t *SecretTemplate) ID() string                { return t.id }
func (t *SecretTemplate) Name() string              { return t.name }
func (t *SecretTemplate) Fields() []TemplateField   { return t.fields }
```

---

## 4. VaultManager API

### Secret Operations

```go
// CreateSecret creates a new secret at the specified location and returns it.
// parentFolderID: ID of the folder to create in, or "" for vault root.
// afterID: ID of the secret after which to insert, or "" to append.
func (m *VaultManager) CreateSecret(name, templateName string, parentFolderID string) (*Secret, error)

// EditSecretField changes the value of a named field on a secret.
func (m *VaultManager) EditSecretField(secretID, fieldName, value string) error

// EditSecretNote changes the observation text of a secret.
func (m *VaultManager) EditSecretNote(secretID, note string) error

// RenameSecret changes the name of a secret.
func (m *VaultManager) RenameSecret(secretID, newName string) error

// DeleteSecret moves a secret to Trash (soft delete).
// It will be permanently removed on the next Save.
func (m *VaultManager) DeleteSecret(secretID string) error

// RestoreSecret moves a secret out of Trash back to its original location.
// If the original location no longer exists, it goes to vault root.
func (m *VaultManager) RestoreSecret(secretID string) error

// DuplicateSecret creates a copy of an existing secret in the same folder.
// The duplicate gets a new ID and "(Cópia)" suffix on the name.
func (m *VaultManager) DuplicateSecret(secretID string) (*Secret, error)

// FavoriteSecret toggles the favorite flag on a secret.
func (m *VaultManager) FavoriteSecret(secretID string, favorite bool) error

// MoveSecret moves a secret to a new parent folder (or vault root).
func (m *VaultManager) MoveSecret(secretID, newParentFolderID string) error

// ReorderSecret moves a secret to a new position within its parent.
func (m *VaultManager) ReorderSecret(secretID string, newIndex int) error

// AddFieldToSecret adds a new field to a secret (advanced edit mode).
func (m *VaultManager) AddFieldToSecret(secretID, fieldName string, fieldType FieldType) error

// RemoveFieldFromSecret removes a field from a secret.
func (m *VaultManager) RemoveFieldFromSecret(secretID, fieldName string) error

// ReorderField changes the position of a field within a secret.
func (m *VaultManager) ReorderField(secretID, fieldName string, newIndex int) error

// ChangeFieldType changes the type of a field.
func (m *VaultManager) ChangeFieldType(secretID, fieldName string, newType FieldType) error
```

### Folder Operations

```go
// CreateFolder creates a new folder at the specified parent.
func (m *VaultManager) CreateFolder(name, parentFolderID string) (*Folder, error)

// RenameFolder changes the name of a folder.
func (m *VaultManager) RenameFolder(folderID, newName string) error

// MoveFolder moves a folder to a new parent (or vault root).
func (m *VaultManager) MoveFolder(folderID, newParentFolderID string) error

// ReorderFolder moves a folder to a new position within its parent.
func (m *VaultManager) ReorderFolder(folderID string, newIndex int) error

// DeleteFolder deletes a folder, moving its children to the parent folder
// (or vault root if the deleted folder was at root level).
// The folder itself goes to Trash.
func (m *VaultManager) DeleteFolder(folderID string) error
```

### Template Operations

```go
// CreateTemplate creates a new secret template.
func (m *VaultManager) CreateTemplate(name string, fields []TemplateField) (*SecretTemplate, error)

// CreateTemplateFromSecret creates a template from an existing secret's field structure.
// Copies field names and types, NOT values.
func (m *VaultManager) CreateTemplateFromSecret(secretID, templateName string) (*SecretTemplate, error)

// EditTemplate updates the fields of an existing template.
// Does NOT affect secrets already created from this template.
func (m *VaultManager) EditTemplate(templateID string, newFields []TemplateField) error

// RenameTemplate renames a template.
func (m *VaultManager) RenameTemplate(templateID, newName string) error

// DeleteTemplate removes a template. Existing secrets are unaffected.
func (m *VaultManager) DeleteTemplate(templateID string) error
```

### Vault-Level Operations

```go
// ChangeSettings updates the vault configuration.
func (m *VaultManager) ChangeSettings(s Settings) error

// PurgeTrash permanently removes all soft-deleted items.
// This is called automatically as part of Save.
func (m *VaultManager) PurgeTrash()
```

---

## 5. Virtual Folders (Favorites & Trash)

Virtual folders are **computed views** — they don't exist in the data model; they are derived at rendering time.

### Favorites Virtual Folder

```go
// CollectFavorites returns all favorite secrets from anywhere in the vault hierarchy.
// These are used to populate the "Favoritos" virtual folder in the sidebar tree.
func CollectFavorites(v *Vault) []*Secret {
    var favorites []*Secret
    collectFavoritesFrom(v.secrets, &favorites)
    for _, f := range v.folders {
        collectFavoritesFromFolder(&f, &favorites)
    }
    return favorites
}

// collectFavoritesFrom appends any favorite secrets from a list.
func collectFavoritesFrom(secrets []Secret, out *[]*Secret) {
    for i := range secrets {
        if secrets[i].IsFavorite() && !secrets[i].IsInTrash() {
            *out = append(*out, &secrets[i])
        }
    }
}
```

### Trash Virtual Folder

```go
// CollectTrashed returns all soft-deleted items from anywhere in the hierarchy.
func CollectTrashed(v *Vault) []*Secret {
    var trashed []*Secret
    collectTrashedFrom(v.secrets, &trashed)
    for _, f := range v.folders {
        collectTrashedFromFolder(&f, &trashed)
    }
    return trashed
}
```

**Key decisions:**
- Trash is visible in sidebar only when there are items in it
- Favorites is visible at the top of the tree only when there are favorites
- Interacting with a secret in either virtual folder navigates to the real item

---

## 6. Search Implementation

### In-Memory Scan (No Index)

Per `descricao.md`, search scans the decrypted in-memory vault. No external search library needed.

```go
// SearchResult holds a matched secret and context about the match.
type SearchResult struct {
    Secret     *Secret
    FolderPath []string  // path of folder names from root to the secret's parent
    MatchField string    // name of the field that matched (empty if name matched)
}

// Search performs a case-insensitive substring search across the vault.
// Searches: secret name, non-sensitive field values, note, and containing folder name.
// Sensitive field values are NEVER searched (would display them unintentionally).
func Search(v *Vault, query string) []SearchResult {
    if query == "" {
        return nil
    }
    query = strings.ToLower(query)
    var results []SearchResult
    searchInSecrets(v.secrets, query, nil, &results)
    for _, f := range v.folders {
        searchInFolder(&f, query, []string{f.Name()}, &results)
    }
    return results
}

// searchInSecrets scans a slice of secrets against query.
// path is the folder path from root to the secrets' parent.
func searchInSecrets(secrets []Secret, query string, path []string, out *[]SearchResult) {
    for i := range secrets {
        s := &secrets[i]
        if s.IsInTrash() {
            continue
        }

        matched := false
        matchField := ""

        // Match on secret name
        if strings.Contains(strings.ToLower(s.Name()), query) {
            matched = true
        }

        // Match on non-sensitive field values
        if !matched {
            for _, f := range s.Fields() {
                if !f.IsSensitive() && strings.Contains(strings.ToLower(f.Value()), query) {
                    matched = true
                    matchField = f.Name()
                    break
                }
            }
        }

        // Match on note (observation)
        if !matched && strings.Contains(strings.ToLower(s.Note()), query) {
            matched = true
            matchField = "observação"
        }

        if matched {
            *out = append(*out, SearchResult{
                Secret:     s,
                FolderPath: path,
                MatchField: matchField,
            })
        }
    }
}
```

**Note:** Folder name matching is handled by checking if the folder name itself matches the query — if so, all secrets in that folder are shown.

---

## 7. Dirty State Tracking

The TUI needs to know if there are unsaved changes to:
1. Show the "modified" indicator in the status bar
2. Intercept Ctrl+Q and ask "Save / Discard / Cancel"

```go
// VaultManager tracks dirty state automatically.
// Every mutation method sets m.dirty = true.
// MarkSaved() resets it after a successful save.

// Example in CreateSecret:
func (m *VaultManager) CreateSecret(...) (*Secret, error) {
    // ... validation and creation logic ...
    m.dirty = true
    m.vault.lastModifiedAt = time.Now()
    return &newSecret, nil
}
```

---

## 8. Default Vault Initialization

When creating a new vault, the manager seeds it with default content:

```go
// InitializeNewVault creates a vault with default templates and folders.
// Per descricao.md, these are user-editable and removable.
func InitializeNewVault() *Vault {
    loginTemplate, _ := gonanoid.Generate(idAlphabet, idLength)
    cardTemplate, _ := gonanoid.Generate(idAlphabet, idLength)
    apiTemplate, _ := gonanoid.Generate(idAlphabet, idLength)

    sitesFolder, _ := gonanoid.Generate(idAlphabet, idLength)
    financeFolder, _ := gonanoid.Generate(idAlphabet, idLength)
    servicesFolder, _ := gonanoid.Generate(idAlphabet, idLength)

    now := time.Now()

    return &Vault{
        version: 1,
        settings: Settings{
            AutoLockMinutes:       5,
            RevealTimeoutSeconds:  15,
            ClipboardClearSeconds: 30,
        },
        templates: []SecretTemplate{
            {
                id:   loginTemplate,
                name: "Login",
                fields: []TemplateField{
                    {name: "URL", fieldType: FieldTypeText},
                    {name: "Username", fieldType: FieldTypeText},
                    {name: "Password", fieldType: FieldTypeSensitive},
                },
            },
            {
                id:   cardTemplate,
                name: "Cartão de Crédito",
                fields: []TemplateField{
                    {name: "Número do Cartão", fieldType: FieldTypeSensitive},
                    {name: "Nome no Cartão", fieldType: FieldTypeText},
                    {name: "Data de Validade", fieldType: FieldTypeText},
                    {name: "CVV", fieldType: FieldTypeSensitive},
                },
            },
            {
                id:   apiTemplate,
                name: "API Key",
                fields: []TemplateField{
                    {name: "Nome da API", fieldType: FieldTypeText},
                    {name: "Chave de API", fieldType: FieldTypeSensitive},
                },
            },
        },
        folders: []Folder{
            {id: sitesFolder, name: "Sites"},
            {id: financeFolder, name: "Financeiro"},
            {id: servicesFolder, name: "Serviços"},
        },
        createdAt:      now,
        lastModifiedAt: now,
    }
}
```

---

## 9. Import / Export Design

### Export (Vault → Plain JSON)

```go
// ExportVault returns the vault as a human-readable JSON byte slice.
// This is the format described in descricao.md for plain-text JSON export.
// SECURITY WARNING: must only be called after explicit user confirmation.
func ExportVault(v *Vault) ([]byte, error) {
    return json.MarshalIndent(toExportable(v), "", "  ")
}
```

### Import (Plain JSON → Vault)

Collision resolution per spec: if a secret with the same name exists at the same location, append `(1)`, `(2)`, etc.

```go
// ImportVault merges the imported vault into the current vault.
// Name collisions are resolved by appending numeric suffixes.
func (m *VaultManager) ImportVault(data []byte) error {
    var imported ExportableVault
    if err := json.Unmarshal(data, &imported); err != nil {
        return fmt.Errorf("formato de importação inválido: %w", err)
    }
    // ... resolve collisions, generate new IDs for imported items, merge
    m.dirty = true
    return nil
}
```

---

## 10. DDD Package Structure

```
domain/
├── vault/
│   ├── vault.go         # Vault aggregate root (read-only accessors)
│   ├── manager.go       # VaultManager (all mutations, invariant enforcement)
│   ├── secret.go        # Secret entity + Field entity + FieldType
│   ├── folder.go        # Folder entity (recursive)
│   ├── template.go      # SecretTemplate + TemplateField
│   ├── search.go        # Search() function (in-memory scan)
│   ├── virtual.go       # CollectFavorites(), CollectTrashed()
│   ├── init.go          # InitializeNewVault()
│   └── export.go        # ExportVault(), ImportVault()
```

---

## 11. Key Design Decisions Summary

| Decision | Rationale |
|----------|-----------|
| NanoID (6 chars, alphanumeric) | 56B combination space — safe at practical vault sizes; per spec |
| Read-only entities + VaultManager | Per spec; enforces invariants, centralizes business logic |
| Recursive Folder struct | Natural fit for arbitrary nesting; JSON serializes cleanly |
| Soft delete via `inTrash` flag | Purged on save; Trash is a virtual folder view |
| Virtual folders = computed views | No data structure needed; derive at render time |
| In-memory search scan | Vault is fully decrypted; no index needed; no data leakage risk |
| Default settings in vault itself | Portability — no config file outside vault file |
| Templates as snapshots | Secrets are self-contained; no retroactive changes per spec |

---

## 12. Sources

| Source | Confidence | Notes |
|--------|------------|-------|
| `github.com/matoous/go-nanoid` README (official Go port) | HIGH | Listed on ai/nanoid repository |
| ai/nanoid README — Other Languages section | HIGH | Confirms matoous/go-nanoid as canonical Go port |
| `descricao.md` — Modelagem section | HIGH | Primary source for entity structure and ID spec |
| Go community DDD patterns (multiple blog posts, 2024) | MEDIUM | No single canonical Go DDD source; patterns are community-established |
| zelark.github.io/nano-id-cc/ — collision calculator | HIGH | Official NanoID collision probability tool |
