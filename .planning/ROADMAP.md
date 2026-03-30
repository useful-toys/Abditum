# Roadmap: Abditum

**Milestone:** v1.0
**Goal:** Ship a security-auditable, offline, single-binary Go TUI password manager with AES-256-GCM encryption, atomic persistence, and a keyboard-driven Bubble Tea interface.
**Created:** 2026-03-27

---

## Phases

- [x] **Phase 1: Project Scaffold + CI Foundation** — Go module, directory tree, static binary build, GitHub Actions CI green on Linux
 (completed 2026-03-29)
- [x] **Phase 2: Crypto Package** — Argon2id key derivation, AES-256-GCM AEAD, memory wipe primitives, mlock, password strength evaluator
 (completed 2026-03-29)
- [x] **Phase 3: Vault Domain + Manager** — All entity types, full Manager API, business rules, invariant enforcement — verified by unit tests
 (completed 2026-03-29)
- [x] **Phase 4: Storage Package** — Binary `.abditum` format, atomic save with `.bak`/`.bak2` chain, Windows MoveFileEx, migration scaffold, startup recovery
 (completed 2026-03-30)
- [ ] **Phase 5: TUI Scaffold + Root Model** — Session state machine, root model, global tick, timer fields, modal overlay — no screens yet
- [ ] **Phase 6: Welcome Screen + Vault Create/Open** — First end-to-end flow: create vault, open vault, error classification, master password strength UI
- [ ] **Phase 7: Vault Tree + Search** — Custom nested tree renderer, keyboard navigation, fold/expand, search overlay scoped to non-sensitive fields
- [ ] **Phase 8: Secret Detail + Edit Flows** — View/create/edit-values/edit-structure modes, template management, duplicate/move/reorder/favorite/soft-delete
- [ ] **Phase 9: Vault Lifecycle Operations** — Save, Save As, Discard/Reload, Change Master Password, Export, Import, Settings configuration
- [ ] **Phase 10: Security Timers + Clipboard + Lock/Exit** — Auto-lock, manual lock, clipboard auto-clear, field reveal timer, screen+scrollback wipe, clean exit
- [ ] **Phase 11: Cross-Platform CI Matrix + Integration Tests** — Full Windows/macOS/Linux matrix, race detector, golden files, secret-pattern scanner, govulncheck

---

## Phase Details

### Phase 1: Project Scaffold + CI Foundation

**Goal:** The repo compiles as a static binary (`CGO_ENABLED=0`), all Charm import paths are the canonical `charm.land/*` paths, and CI executes build + lint + tests on every push.

**Requirements:** COMPAT-01, CI-01

**Plans:** 3 plans in 1 wave

Plans:
- [x] 01-01-PLAN.md — Initialize Go module with dependencies and create directory structure with package stubs
- [x] 01-02-PLAN.md — Configure CI workflow and Makefile for automated building, linting, and testing
- [x] 01-03-PLAN.md — Configure golangci-lint with security rules and verify forbidden pattern detection

**UAT:**
- [ ] `CGO_ENABLED=0 go build ./cmd/abditum` succeeds and produces an executable binary
- [ ] On Linux: `file ./abditum` reports "statically linked" (no libc dependency)
- [ ] CI workflow triggers on push to `main` and reports all jobs green
- [ ] All `.go` files use `charm.land/*` import paths — zero `github.com/charmbracelet/*` imports present
- [ ] `go vet ./...` passes with zero findings on the empty stubs

**Pitfall watch:** Import paths migrated in late 2025 — use `charm.land/bubbletea/v2` (NOT `github.com/charmbracelet/bubbletea`). Set `CGO_ENABLED=0` as an environment variable in CI, not only at build time — it must be enforced globally so all subsequent `go test` runs also use static linking.

---

### Phase 2: Crypto Package

**Goal:** `internal/crypto` delivers production-ready Argon2id key derivation, AES-256-GCM authenticated encryption, secure memory primitives, and password strength evaluation — all verified by tests that will catch any future security regression.

**Requirements:** CRYPTO-01, CRYPTO-02, CRYPTO-03, CRYPTO-04, CRYPTO-05, CRYPTO-06, PWD-01

**Plans:** 1/1 plans complete

Plans:
- [x] 02-01-PLAN.md — Complete crypto package: Argon2id KDF, AES-256-GCM AEAD, memory security primitives, password strength evaluation, comprehensive tests

**UAT:**
- [ ] `Encrypt(key, p)` + `Decrypt(key, c)` roundtrip returns original plaintext byte-for-byte
- [ ] Encrypting identical plaintext twice with the same key produces two distinct ciphertexts (nonce not reused)
- [ ] `Decrypt` with wrong key returns `crypto.ErrAuthFailed` — no panic, no internal error string exposed
- [ ] `ZeroBytes(b)` fills every byte of `b` with `0x00`
- [ ] `EvaluatePasswordStrength([]byte("Abc1!Abc1!12"))` returns `StrengthStrong`; `"abc123"` returns `StrengthWeak`
- [ ] `CGO_ENABLED=0 go test ./internal/crypto/... -race` passes with zero race conditions

**Pitfall watch:** Argon2id `m` parameter is in **KiB** — 256 MiB = `262144`, not `256`. `string` for sensitive data burns you: every API in this package accepts `[]byte` for passwords and keys — establish this at the function signature level. Nonce must be generated via `io.ReadFull(rand.Reader, nonce)` immediately before each `gcm.Seal()` call — never reuse a nonce slice across calls. mlock/VirtualLock failure is best-effort — always check the error and continue, never fatal.

---

### Phase 3: Vault Domain + Manager

**Goal:** `internal/vault` delivers a complete, fully-tested in-memory domain layer — all entity types, the full Manager API, and every business rule enforced and verified via unit tests — before any file I/O or TUI code depends on it.

**Requirements:** VAULT-02, SEC-05, FOLDER-01, FOLDER-02, FOLDER-03, FOLDER-04, FOLDER-05, TPL-01, TPL-02, TPL-03, TPL-04, TPL-05, TPL-06

**Plans:** 7/7 plans complete

Plans:
- [x] 03-01-PLAN.md — Domain Entities + Factory: Define Cofre, Pasta, Segredo, Modelo, Campo types, EstadoSessao enum, factory functions, error types, comprehensive entity tests
- [x] 03-02-PLAN.md — Manager + Cofre Lifecycle: Manager struct, Lock with memory wiping, Salvar with atomic two-phase commit, timestamp tracking, lifecycle tests
- [x] 03-03-PLAN.md — Folder Management: CriarPasta, RenomearPasta, MoverPasta with cycle detection, ExcluirPasta with promotion and conflict resolution, Pasta Geral protection
- [x] 03-04-PLAN.md — Template Management: CriarModelo, RenomearModelo, ExcluirModelo with in-use check, field operations (add/remove/reorder), alphabetical ordering, "Observação" prohibition
- [x] 03-05-PLAN.md — Secret Lifecycle + State Machine: CriarSegredo, ExcluirSegredo, RestaurarSegredo, FavoritarSegredo, DuplicarSegredo with name progression, estadoSessao transitions
- [x] 03-06-PLAN.md — Secret CRUD + Structure: RenomearSegredo, EditarCampoSegredo, EditarObservacao (separate field), MoverSegredo, ReposicionarSegredo with estadoSessao tracking
- [x] 03-07-PLAN.md — Search + Favorites + Comprehensive Validation: Buscar with sensitive field exclusion (QUERY-02), ListarFavoritos with DFS, integration tests, UAT validation, final package documentation

**UAT:**
- [ ] `Manager.Create` initializes vault with Pasta Geral + "Sites e Apps" + "Financeiro" subfolders + 3 default templates (Login, Cartão de Crédito, Chave de API)
- [ ] `Manager.CreateFolder` returns error when name already exists in the target parent
- [ ] `Manager.MoveFolder` returns error when moving a folder into one of its own descendants (cycle)
- [ ] `Manager.DeleteFolder` promotes children (including `StateDeleted` secrets, which retain their `StateDeleted` state) and resolves name conflicts with numeric suffix; Pasta Geral deletion returns explicit error
- [ ] `Manager.CreateSecret` always produces a secret with Observation as the last field; `UpdateSecretStructure` cannot rename, reorder, or delete the Observation field
- [ ] `DuplicateSecret("X")` produces "X (1)"; duplicating "X (1)" produces "X (2)"
- [ ] `CreateTemplateFromSecret` excludes ALL fields named 'Observação' (both auto-Observation and any user-named field with that name); a template produced from such a secret never contains a field named 'Observação'
- [ ] `UpdateTemplateStructure` returns explicit error when attempting to add or rename a field to 'Observação'
- [ ] Search function called with a string present only in a `FieldTypeSensitive` field **value** returns zero results; search called with the **name** of a sensitive field returns secrets containing that field
- [ ] `CreateSecret` produces secret with `sessionState = StateIncluded`; `UpdateSecret` on a `StateOriginal` secret → `StateModified`; `UpdateSecret` on `StateIncluded` → remains `StateIncluded`; `SoftDeleteSecret` → `StateDeleted` (stores `prevState`); `RestoreSecret` → restores previous state
- [ ] `StateDeleted` secret is excluded from search results; at Open/Discard all loaded secrets have `sessionState = StateOriginal`
- [ ] Templates are always returned in alphabetical order by name regardless of creation order; creating "Z-Template", "A-Template", "M-Template" in sequence returns them as A-Template, M-Template, Z-Template (TPL-06)

**Pitfall watch:** Pasta Geral protection requires the Manager to store the root folder's NanoID and check it in every destructive method. NanoID must use `crypto/rand` — verify `go-nanoid`'s `gonanoid.New()` uses the secure alphabet by default. Cycle detection must walk the FULL ancestor chain, not just immediate parent. `Manager.Vault()` must return a snapshot — exposing a live pointer to domain state allows TUI bugs to corrupt the domain silently.

---

### Phase 4: Storage Package

**Goal:** `internal/storage` delivers the binary `.abditum` format with atomic writes, `.bak`/`.bak2` backup chain, startup orphan recovery, external change detection, and a migration scaffold — including the Windows-specific `MoveFileEx` rename path — fully verified.

**Requirements:** ATOMIC-01, ATOMIC-02, ATOMIC-03, ATOMIC-04, COMPAT-03

**Plans:** 4 plans in 3 waves

Plans:
- [x] 04-01-PLAN.md — Foundation: AAD-aware crypto, vault JSON serialization, storage format constants
- [x] 04-02-PLAN.md — Core I/O: Save/SaveNew/Load with atomic .tmp protocol, platform-specific rename
- [x] 04-03-PLAN.md — Recovery & Migration: RecoverOrphans, DetectExternalChange, migration scaffold
- [x] 04-04-PLAN.md — Integration: FileRepository adapter implementing RepositorioCofre, end-to-end tests

**UAT:**
- [ ] `SaveNew` + `Load` roundtrip returns an identical `*vault.Vault` (deep equality on all fields)
- [ ] File with wrong magic returns `ErrInvalidMagic` — no panic or internal error string exposed to caller
- [ ] File with correct magic but version `> currentVersion` returns `ErrVersionTooNew`
- [ ] File with correct header but wrong password returns `ErrAuthFailed`
- [ ] After second `Save` on same path: `.abditum.bak` exists; after third `Save`: old `.bak` becomes `.bak2` and new `.bak` is created
- [ ] Simulated crash (`.tmp` present, target missing): `RecoverOrphans` removes stale `.tmp` on next startup
- [ ] `FileMetadata` `mtime`/`size` returned from `Load`; modifying the file externally causes `DetectExternalChange` to return `true`
- [ ] v0 test fixture file loads cleanly via `Migrate` and passes validation at v1

**Pitfall watch:** Windows `os.Rename` is NOT atomic when replacing a file in use — must use `MoveFileEx` with `MOVEFILE_REPLACE_EXISTING`. The `.tmp` file MUST be written to `filepath.Dir(vaultPath)`, never `os.TempDir()` — cross-device rename (`EXDEV`) is a real failure on encrypted home directories and network mounts. Argon2id parameters (t, m, p, salt) MUST be stored in the plaintext file header and read back at decrypt time — hardcoding them at decrypt breaks every existing vault on any future parameter change.

---

### Phase 5: TUI Scaffold + Root Model

**Goal:** The application launches, shows a placeholder, and the entire root model architecture — session state machine, flat child-model composition, global tick, timer field declarations, and modal overlay plumbing — is in place and race-condition-free before any screen is implemented.

**Requirements:** *(foundational infrastructure — no single v1 requirement is independently verifiable at this phase; all subsequent TUI phases depend on this foundation)*

**Plans:**
1. Define `sessionState` enum: `stateWelcome`, `stateVaultTree`, `stateSecretDetail`, `stateLocked`; define `rootModel` struct fields: `state sessionState`, `mgr *vault.Manager`, `welcome welcomeModel`, `tree vaultModel`, `detail secretDetailModel`, `modal *modalModel`, `lockTimer int`, `clipboardTimer int`, `revealedFields map[string]time.Time`, `clipboardFieldID string`; initialize all child models to zero value on screen transition
2. Implement `rootModel.Init() tea.Cmd` — returns `tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg{t} })` to start the global 1-second tick; implement `rootModel.Update(msg tea.Msg) (tea.Model, tea.Cmd)` — dispatch: modal intercepts all messages first when non-nil; else route to active child model by state; reset all inactivity timers on any `tea.KeyPressMsg`; implement `rootModel.View() tea.View` — render active child; overlay modal on top using `lipgloss.Place` if `modal != nil`
3. Implement `modalModel` (reusable overlay): fields `title`, `body`, `options []string`, `selectedIndex int`, `onSelect func(int) tea.Cmd`; keyboard: `j`/`k` or arrows to move selection, Enter to confirm, `ESC` to close (calls `onSelect(-1)`); exported constructor `newModal(title, body string, options []string, cb func(int) tea.Cmd) *modalModel`; renders centered in terminal using `lipgloss` border + padding
4. Implement global tick handler in `rootModel.Update` on `tickMsg`: decrement `lockTimer` (if >0); decrement `clipboardTimer` (if >0, fire `clearClipboardCmd` at 0); expire entries in `revealedFields` map where `time.Now().After(expiry)` (set field back to masked); re-issue `tea.Tick(time.Second, ...)` cmd to keep tick loop alive
5. Wire `cmd/abditum/main.go`: parse `os.Args[1]` as optional vault file path; instantiate `vault.Manager`; build `rootModel{mgr: mgr, ...}`; call `tea.NewProgram(root, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run()`; print generic fatal error to stderr (no paths/secrets) and exit 1 on `Run` error; exit 0 on clean return

**UAT:**
- [ ] `./abditum` launches without panic and renders a placeholder welcome message; `q`/`ctrl+c` exits with code 0
- [ ] `./abditum /path/to/vault.abditum` launches and passes path through to root model without crashing
- [ ] `rootModel.View()` return type compiles as `tea.View` — not `string` (Bubble Tea v2 API)
- [ ] `CGO_ENABLED=0 go test -race ./internal/tui/...` passes with zero race detector alerts
- [ ] Sending a sequence of simulated `tickMsg` in a unit test fires appropriate timer callbacks at correct counts

**Pitfall watch:** In Bubble Tea v2, `View()` returns `tea.View` (not `string`), key messages are `tea.KeyPressMsg` (not `tea.KeyMsg`), and the space key is the string `"space"` (not `" "`). Verify these against `charm.land/bubbletea/v2` source before writing any key handler. Timer decrement must fire on `tickMsg`, NOT inside goroutines — using `time.AfterFunc` introduces goroutine leaks and races.

---

### Phase 6: Welcome Screen + Vault Create/Open

**Goal:** A user can pick a vault file path, create a new vault with a master password (with real-time strength feedback), or open an existing vault — and every error case surfaces the correct generic user message with no technical detail leaked.

**Requirements:** VAULT-01, VAULT-03, VAULT-04, VAULT-05

**Plans:**
1. Implement `welcomeModel`: two sub-states (`subStatePickPath`, `subStateCreatePassword`, `subStateOpenPassword`); file path `textinput`; three action choices (Create / Open / Quit) navigated with `j`/`k`/Enter; key-hint footer ("↑↓ navigate · Enter select · q quit") using `lipgloss`
2. Implement vault creation flow (VAULT-01): dual password input fields (Enter password / Confirm password); real-time strength badge (Weak/Strong) rendered inline beneath first field — calls `crypto.EvaluatePasswordStrength` on every keystroke; mismatch shown as inline error; `str` mismatch or empty → no submit; Weak password → non-blocking warning banner, submit allowed; on submit: convert `textinput.Value()` → `[]byte` immediately, zero textinput buffer, call `Manager.Create(path, passwordBytes)`, transition to `stateVaultTree`
3. Implement vault open flow (VAULT-03): single password input; on submit: convert to `[]byte`, zero buffer, call `Manager.Open(path, passwordBytes)`; map storage sentinel errors to user messages (NO Go error strings, NO internal details): `storage.ErrInvalidMagic` → "Invalid file — not an Abditum vault"; `storage.ErrVersionTooNew` → "This vault was created by a newer version of Abditum"; `storage.ErrAuthFailed` → "Incorrect password — please try again" (clear input, allow retry); `storage.ErrCorrupted` → "Vault integrity error — file cannot be opened" (return to path selection, no retry) (VAULT-04, VAULT-05)
4. Implement file path input validation: check file exists for Open (suggest creating if not found); check path is writable directory for Create; friendly messages for common errors (directory not found, permission denied) — generic, no OS error strings
5. Implement `RecoverOrphans` call on `Manager.Open` path: before `storage.Load`, call `storage.RecoverOrphans(path)` — ignore orphan recovery errors (log generically); ensures startup recovery from Phase 4 fires on every vault open
6. Write `teatest/v2` golden file tests for: welcome screen initial state (80×24), create-vault step 1 (path entered), create-vault step 2 (password + strength badge: Weak), create-vault step 3 (strength badge: Strong), open-vault password prompt, open-vault ErrAuthFailed error state, open-vault ErrCorrupted error state; use `teatest.WithInitialTermSize(80, 24)` throughout

**UAT:**
- [ ] User enters path + matching passwords → vault file created on disk; TUI transitions to vault tree screen
- [ ] Strength badge shows "Weak" for `"abc123"`; shows "Strong" for `"Abc1!Abc1!12"`; creation is not blocked on Weak (only warned)
- [ ] Mismatched confirm password shows inline error message; no vault is created until they match
- [ ] Opening with correct password transitions to vault tree; opening with wrong password shows "Incorrect password — please try again" and clears the password field for retry
- [ ] Opening a non-`.abditum` binary file shows "Invalid file — not an Abditum vault" — no Go error string visible
- [ ] Opening a vault from a newer format version shows "This vault was created by a newer version"

**Pitfall watch:** `textinput.Value()` returns `string` — convert to `[]byte` inside `Update()` the instant the user submits, zero the textinput buffer (`ti.SetValue("")`), and never pass the password as `string` further down the call stack. Do not show file paths, Go error messages (`err.Error()`), or internal error types in any user-facing error message.

---

### Phase 7: Vault Tree + Search

**Goal:** Users can navigate the complete folder/secret hierarchy, expand/collapse any folder, instantly filter by name and common-field content, and see favorites and soft-deleted secrets with distinct visual treatments.

**Requirements:** QUERY-01, QUERY-02, QUERY-06, QUERY-07

**Plans:**
1. Implement custom tree renderer (`tree.go`) — **do NOT use `bubbles/list`** (insufficient for recursive folder+secret hierarchy); recursive `renderNode(folder *vault.Folder, depth int, expandState map[string]bool)` producing indented lines with box-drawing characters; folders shown before secrets at each level; folder lines: fold indicator (`▶`/`▼`) + folder name + recursive active-secret count `(N)` (counts all non-`StateDeleted` secrets in folder and all its subfolders recursively — QUERY-01); secret lines: name + optional favorite `★` prefix + optional session-state indicator (`+` for `StateIncluded`, `~` for `StateModified`, `✗` for `StateDeleted` — QUERY-06) + optional dim/strikethrough for `StateDeleted`; render virtual "Favoritos" node as a fixed sibling above Pasta Geral — always visible, never expandable past its own contents, fold indicator follows same `▶`/`▼` convention; Favoritos always expanded by default (QUERY-07)
2. Implement keyboard navigation in `vaultModel`: `j`/`↓` (next item), `k`/`↑` (prev item), `l`/`→` or Enter-on-folder (expand folder), `h`/`←` (collapse folder), Enter-on-secret (open secret detail → `stateSecretDetail`); maintain `cursor int` into a flattened renderable item list rebuilt on each tree mutation; maintain `expandState map[string]bool` (default: Favoritos expanded, top-level real folders expanded, subfolders collapsed); cursor can reach the Favoritos virtual node and its secret entries — Enter on a Favoritos entry navigates to that secret in detail view; pressing `n`/`N`/`d`/`m` while cursor is on the Favoritos node or its entries is a no-op (read-only); `n` (new secret — trigger create flow from Phase 8); `N` (new folder); `/` (open search overlay)
3. Implement session-state and favorite visual treatment: `lipgloss` style `Strikethrough(true).Faint(true)` for `StateDeleted` secrets; `★` prefix for favorites; `StateDeleted` items shown in tree but visually distinguished; `StateIncluded` and `StateModified` secrets show inline session-state indicator (QUERY-06); search must exclude `StateDeleted` secrets from results
4. Implement search overlay (`/` key): inline `textinput` at bottom of tree view; on every keystroke call `vault.Search(query string, v *vault.Vault) []*vault.Secret` — match by: secret name, all field names (including names of sensitive fields — field names are always searchable regardless of type), common field values, Observation value (case-insensitive, accentuation-normalized via `golang.org/x/text/unicode/norm` NFD + stripping); **never read `field.Value` when `field.Type == FieldTypeSensitive`** (QUERY-02 hard gate); ESC clears search and returns to tree; Enter on result selects secret
5. Write QUERY-02 negative test (must be written FIRST, before search implementation): construct a secret with one sensitive field whose value is `"hunter2"`; call `Search("hunter2", vault)` → assert result is empty slice; this test must fail if search ever accidentally reads sensitive field values
6. Write golden file tests: full tree with 3-level folder nesting, search overlay with results, search overlay with zero results, tree with favorited + soft-deleted secrets visible

**UAT:**
- [ ] Vault tree displays subfolders first, then secrets within each folder, matching domain order
- [ ] Each folder line shows the recursive count of active (non-`StateDeleted`) secrets across all its subfolders
- [ ] User can expand/collapse any non-root folder; expansion state persists for the session
- [ ] Searching with `/` filters the display in real time; results exclude `StateDeleted` secrets
- [ ] Searching for a string matching a sensitive field **name** (e.g. "Senha") returns secrets containing that field
- [ ] Searching for a string that exists only in a sensitive field **value** returns zero results (hard requirement)
- [ ] Favorite secrets show `★`; `StateDeleted` secrets show dim/strikethrough; `StateIncluded` and `StateModified` secrets show distinct session-state indicator
- [ ] Pressing Enter on a secret transitions to secret detail view (Phase 8) with that secret's ID
- [ ] "Favoritos" virtual node appears above Pasta Geral in the tree, expanded by default; entering it shows all `favorito = true` secrets; Enter on a Favoritos entry opens that secret in detail view; no write operations (create/move/delete) are available from within the Favoritos node

**Pitfall watch:** `bubbles/list` cannot render a recursive folder+secret tree with expand/collapse — implement a custom renderer. Accentuation normalization for QUERY-02 requires `golang.org/x/text/unicode/norm` NFD decomposition (not just `strings.ToLower`). **Write the sensitive-field negative test FIRST, before the search loop** — it is easy to accidentally iterate all field values without the `IsSensitive` guard.

---

### Phase 8: Secret Detail View + Edit Flows

**Goal:** Users have complete control over individual secrets: viewing (with sensitive fields masked), creating from template, editing values, editing structure, duplicating, moving, marking for deletion, reordering, favoriting — plus full template management (create/rename/edit/delete/create-from-secret).

**Requirements:** SEC-01, SEC-02, SEC-03, SEC-04, SEC-06, SEC-07, SEC-08, SEC-09, QUERY-03

**Plans:**
1. Implement `secretDetailModel` view mode (QUERY-03): render secret name (bold), all fields in declared order with name → value layout; sensitive fields masked as `••••••••` by default; Observation rendered last always; key-hint footer: `e` (edit values), `E` (edit structure), `d` (toggle soft-delete), `f` (favorite), `m` (move), `r` (reorder), `ctrl+d` (duplicate), `t` (template management), `v` (reveal field — Phase 10), `c` (copy field — Phase 10), `ESC` (back to tree), `?` (help — Phase 11)
2. Implement create-secret flow (SEC-01): from tree view `n` key → template picker overlay (scrollable from `bubbles/v2` list or custom); templates shown by name; "No template (Observation only)" option at top; folder picker (flat list of all folders); on confirm → `Manager.CreateSecret(folderID, templateID)` → transition to edit-values mode with new secret loaded
3. Implement edit-values mode (SEC-03): one `textinput` per field; sensitive fields use `textinput.EchoPassword` mode; Observation at bottom, also editable; `ctrl+s` → batch `Manager.UpdateSecret(id, SecretChanges{Name: ..., FieldValues: ...})`; `ESC` with changes → confirm-discard modal; `ESC` without changes → return to view mode; password inputs: convert `textinput.Value()` → `[]byte` immediately after submit, zero textinput buffer; leave a comment stub `// TODO(v2): "Generate →" action for sensitive fields` at the sensitive field input site (v2 generator slot)
4. Implement edit-structure mode (SEC-04): enter with capital `E`; shows field list; controls per field: rename (inline input), reorder (`shift+j`/`shift+k`), delete (confirm modal); add field button at bottom: name input + type picker (Common / Sensitive); type cannot be changed post-creation; Observation row is rendered in a read-only non-interactive style with label "(auto — cannot be modified)" — no edit controls; `ctrl+s` → `Manager.UpdateSecretStructure`
5. Implement duplicate (SEC-02): `ctrl+d` → `Manager.DuplicateSecret(id)` → stay in tree with new item selected; implement move (SEC-08): `m` → folder picker overlay → `Manager.MoveSecret(id, destFolderID)`; implement reorder (SEC-09): `r` → reorder mode — highlight item, `j`/`k` move it within parent folder, Enter confirm → `Manager.ReorderSecret(id, newIndex)`, `ESC` cancel
6. Implement favorite toggle (SEC-06): `f` → `Manager.FavoriteSecret(id, !current)` → re-render view; implement soft-delete toggle (SEC-07): `d` → if Active → confirm-delete modal → `Manager.SoftDeleteSecret(id)` → view shows "(marked for deletion)" badge; if SoftDeleted → `Manager.RestoreSecret(id)` → badge removed; no confirmation needed to restore
7. Implement template management overlay (SEC-05's inverse — TPL surface): `t` from any vault tree or detail view → template list overlay; from list: Create (name input + add field loop), Rename (inline input), Edit Structure (same add/rename/reorder/delete as secret structure), Delete (confirm modal), Create From Secret (picks a secret from tree — calls `Manager.CreateTemplateFromSecret`); all operations call respective `Manager.Template*` methods

**UAT:**
- [ ] Creating a secret from "Login" template produces fields: URL, Usuário, Senha (sensitive), Observação (auto, last)
- [ ] Creating a secret with no template produces only Observação field
- [ ] Edit-values mode: Senha field is masked during editing (EchoPassword); submitting calls Manager and sets `IsModified() == true`
- [ ] Edit-structure mode: Observation row has no edit controls; adding a field appears before Observation; deleting a non-Observation field removes it; Observation cannot be deleted or moved
- [ ] Duplicating "MySecret" produces "MySecret (1)" immediately after original in tree
- [ ] Marking for deletion shows visual badge; `IsModified()` becomes true; toggling restores and clears flag
- [ ] Template management: creating a template + creating a secret from it populates fields correctly; deleting template does not affect existing secrets built from it

**Pitfall watch:** Sensitive field values must NEVER be passed to `textinput.SetValue()` for pre-population in edit mode — this would copy them into a `string` that cannot be zeroed. Instead, leave the field blank and only populate on explicit user action. Comment the v2 generator slot at every sensitive `textinput` construction site so it is discoverable in Phase v2 without a codebase-wide search.

---

### Phase 9: Vault Lifecycle Operations

**Goal:** Users can save, save-as, discard/reload, change master password, export, import, and configure timer settings — all vault-level operations wired through Manager with complete external-change detection and error handling.

**Requirements:** VAULT-06, VAULT-07, VAULT-08, VAULT-09, VAULT-10, VAULT-15, VAULT-16, VAULT-17

**Plans:**
1. Add key bindings in vault tree and detail views: `ctrl+s` (Save), `ctrl+shift+s` (Save As), `ctrl+r` (Discard/Reload), `ctrl+p` (Change Master Password), `ctrl+e` (Export), `ctrl+i` (Import), `ctrl+,` (Settings); each binding dispatches a typed `tea.Msg` to `rootModel` which handles all lifecycle operations (keeps TUI screens decoupled from file operations)
2. Implement Save flow (VAULT-06): `rootModel` receives `SaveMsg` → call `Manager.Save()` (which calls `storage.Save` internally); if `storage.DetectExternalChange` returns true before saving, display modal: "This vault was modified externally.\n\nOverwrite / Save as new file / Cancel"; mapping: 0 → force save, 1 → trigger SaveAs flow, -1 → cancel (VAULT-08)
3. Implement Save As flow (VAULT-07): modal with path `textinput`; validate path is not identical to `Manager.CurrentPath()` (reject with error message); on confirm → `Manager.SaveAs(destPath)` → `CurrentPath()` now returns `destPath`; IsModified becomes false
4. Implement Discard/Reload flow (VAULT-09): confirm modal "Discard all unsaved changes and reload from disk?"; if file was externally modified, add secondary warning line; on confirm → `Manager.Discard()` (reloads from disk using current session key stored in Manager); tree re-renders with refreshed snapshot
5. Implement Change Master Password flow (VAULT-10): three-step modal sequence: step 1 = current password input (re-derive key in-process to verify; reject if wrong), step 2 = new password + strength badge (same as create flow), step 3 = confirm new password (must match); on all three confirmed → `Manager.ChangeMasterPassword(current, next []byte)` which immediately re-encrypts and saves; operation is irreversible — show "This saves immediately and cannot be undone" warning before step 1
6. Implement Export flow (VAULT-15): path input modal + risk warning: "The exported file is unencrypted. Store it somewhere safe and delete it when done."; on confirm → `Manager.Export(path)` writes JSON: all active (non-soft-deleted) secrets including their timestamps (`data_criacao`/`data_ultima_modificacao`), all folders, all templates; vault-level timer settings (`Configuracoes` block) are not exported
7. Implement Import flow (VAULT-16): path input modal; call `Manager.Import(path)` — validate file first: if invalid JSON or Pasta Geral absent → return generic error, abort import; merge rules: folders merged by full path (sequence of names from root) — existing folder → merge contents; new folder → create; within each merged folder: incoming secret with same **name** → **replaces** existing; incoming secret with unique name → added; incoming template with same **name** → **replaces** existing; incoming template with unique name → added; show import summary modal: "Imported N secrets, M folders, K templates"
8. Implement Settings screen (VAULT-17): dedicated screen reachable from `ctrl+,`; three integer inputs: "Auto-lock after (minutes):", "Reveal sensitive field for (seconds):", "Clear clipboard after (seconds):"; all three required, positive integers only; `ctrl+s` → `Manager.UpdateSettings(settings)` → new values used by Phase 10 timers immediately; `ESC` → discard unsaved settings changes (confirm modal if modified)

**UAT:**
- [ ] Save writes to current path; `.abditum.bak` created; `IsModified()` returns false after successful save
- [ ] External change detection: modifying vault file externally between Open and Save triggers the 3-option modal
- [ ] Save As to a different path changes `CurrentPath()`; subsequent Ctrl+S writes to the new path
- [ ] Discard without pending changes still confirms; reloads vault cleanly; all unsaved mutations are gone
- [ ] Change Master Password: wrong current password shows error at step 1; correct current + strong new password saves immediately; vault can be reopened with new password afterward
- [ ] Export: resulting JSON contains all active secrets including timestamps (`data_criacao`/`data_ultima_modificacao`); soft-deleted secrets absent; `Configuracoes` timer-settings block absent
- [ ] Import from invalid JSON or file without Pasta Geral → fails with generic error message, no changes applied to vault
- [ ] Import from valid Abditum JSON: folders merged by path; incoming secret with same name as existing secret in that folder replaces it; incoming secret with unique name is added; import summary modal shows correct counts
- [ ] Settings: changing auto-lock to 1 minute and waiting 1 minute triggers auto-lock (verified in Phase 10)

---

### Phase 10: Security Timers + Clipboard + Lock/Exit

**Goal:** All time-driven security behaviors work correctly and synchronously on all three platforms: auto-lock, manual lock, clipboard auto-clear on timer and on lock/exit, sensitive field auto-hide on timer, memory wipe, and screen+scrollback clear — with clean exit handling for all pending-change scenarios.

**Requirements:** VAULT-11, VAULT-12, VAULT-13, VAULT-14, QUERY-04, QUERY-05

**Plans:**
1. Wire auto-lock timer (VAULT-11): every `tea.KeyPressMsg` or `tea.MouseMsg` in `rootModel.Update` resets `lockTimer` to `Manager.Settings().LockTimeoutMinutes * 60`; global tick decrements; at 0 → emit `LockMsg{}` internally; test that ALL input varieties (key, mouse, window resize) reset the timer — missing any input type breaks the invariant
2. Implement `LockMsg` handler in `rootModel` (VAULT-12, VAULT-13): call `Manager.Lock()` (zeros in-memory key and all field `[]byte` buffers); set `rootModel.state = stateLocked`; zero all `revealedFields` map entries; clear clipboard synchronously (before any further rendering); emit `clearScreenCmd` which outputs `"\033[3J\033[2J\033[H"` to stdout AFTER Bubble Tea program stops its render loop (`p.ReleaseTerminal()` or equivalent); transition welcome model to re-enter-password sub-state for same path
3. Implement manual lock keybinding `ctrl+l` from any active vault screen → dispatch same `LockMsg{}` as auto-lock (VAULT-12); lock must work from vault tree, secret detail, and any overlay state
4. Implement clean exit (VAULT-14): `q` from welcome / `ctrl+q` from vault screens → check `Manager.IsModified()`; if no pending changes → proceed directly to exit; if pending → modal "Save and Quit / Discard and Quit / Cancel"; on any exit path: `Manager.Lock()` (zeros memory), clear clipboard synchronously, emit `clearScreenCmd`, `tea.Quit`; exit code 0 on clean exit, 1 on fatal startup error only
5. Implement clipboard copy (QUERY-05): `c` key in secret detail view → cursor selects field → call `clipboard.WriteAll(string(field.Value))` via `github.com/atotto/clipboard`; set `clipboardTimer = Manager.Settings().ClipboardClearSeconds` and `clipboardFieldID = fieldID`; on timer expiry → `clipboard.WriteAll("")`; on lock or exit → `clipboard.WriteAll("")` synchronously before process returns; handle headless/Wayland/SSH: if `clipboard.WriteAll` returns error → show gentle warning in footer, continue without crash (QUERY-05 Wayland best-effort)
6. Implement field reveal (QUERY-04): `v` key in secret detail view → cursor selects a sensitive field → set `revealedFields[fieldID] = time.Now().Add(time.Duration(settings.RevealSeconds) * time.Second)`; `secretDetailModel.View()` reads `rootModel.revealedFields` (passed via render args or message) to determine masking; global tick clears expired entries from map; re-triggering `v` on already-revealed field resets the timer
7. Wire Phase 9 Settings into Phase 10 timers: when `UpdateSettings` is called, update in-flight `lockTimer` cap (reset to new value if new value is smaller than current remaining); `clipboardTimer` and `revealedFields` expiry use `Manager.Settings()` at timer-set time; verify that settings screen changes take effect on next timer action without restart
8. Write synchronized timer tests: unit-test `rootModel` with mock tick messages; assert lock fires after exactly N ticks; assert any key between ticks resets count; assert clipboard cleared synchronously in `LockMsg` handler (not in a goroutine); assert reveal field correctly expires; run `go test -race`

**UAT:**
- [ ] Auto-lock: after the configured inactivity period with no input, vault locks automatically and terminal shows re-enter-password prompt
- [ ] Any key press or mouse event resets the inactivity timer to full duration
- [ ] Manual lock (`ctrl+l`) locks immediately from vault tree or detail view
- [ ] After lock: sensitive data is not visible in terminal scrollback buffer (manual verification in xterm, Windows Terminal, iTerm2)
- [ ] Clipboard: copying a field value and waiting `clipboardClearSeconds` clears the clipboard; copying another field resets the timer
- [ ] Clipboard is cleared synchronously when vault locks or app exits (no residual data after process exits)
- [ ] Revealing a sensitive field (`v`) shows value for `revealSeconds` then auto-hides to `••••••••`
- [ ] On Wayland without `wl-clipboard` present: clipboard operation shows a gentle footer notice; application continues normally

**Pitfall watch:** Terminal clear-screen `\033[2J` does NOT clear scrollback — only `\033[3J\033[2J\033[H` clears both. Clipboard clear on lock/exit MUST be synchronous — `time.AfterFunc` goroutine may not execute before `tea.Quit` unwinds; call `clipboard.WriteAll("")` directly in the `LockMsg` handler. On Wayland, `atotto/clipboard` may silently fail if `wl-clipboard` is absent — wrap the call, check the error, show a footer notice (no crash, no log with sensitive content).

---

### Phase 11: Cross-Platform CI Matrix + Integration Tests + Golden Files

**Goal:** Every test passes on Windows, macOS, and Linux under the race detector; golden files lock all visual regressions; no sensitive data can appear in test fixtures; and dependency vulnerabilities are checked automatically.

**Requirements:** CI-02, COMPAT-02

**Plans:**
1. Expand GitHub Actions CI matrix to three OS: `ubuntu-latest`, `macos-latest`, `windows-latest`; all jobs run `CGO_ENABLED=0 go test ./... -race -count=1`; confirm Windows-specific `MoveFileEx` code path is exercised in the windows CI job explicitly (call `atomicRenameWindows` in a Windows-only integration test)
2. Write end-to-end manager integration test (in-process, no subprocess): `TestFullVaultRoundtrip` — `Manager.Create` → `Manager.CreateFolder` → `Manager.CreateSecret` → `Manager.UpdateSecret` (values) → `Manager.FavoriteSecret` → `Manager.Save` → new `Manager.Open` (same path, same password) → assert all mutations present → assert `IsModified() == false`; run under `-race`
3. Collect all `teatest/v2` golden file tests from Phases 6-10 in CI; add any missing golden files for: locked state (re-enter password prompt), export confirmation modal, import summary modal, settings screen, help overlay; document golden file update procedure with `UPDATE_GOLDEN=true go test ./...`
4. Implement golden file secret-pattern scanner (`cmd/scan-golden/main.go` or `TestScanGoldenFiles` test): walk all files matching `testdata/**/*_golden*.txt`; apply regex patterns for common sensitive data signatures (passwords matching complexity rules, 16-digit groups, API key patterns like `sk-`, `Bearer `); fail CI if any pattern matches; use clearly fictional fixture values: `"P@ssw0rd!Test42"` (labeled fake), `"4242424242424242"` (recognizably fictional test card)
5. Add `govulncheck` step to CI (all three OS): `go install golang.org/x/vuln/cmd/govulncheck@latest && govulncheck ./...`; fail CI on any HIGH or CRITICAL severity finding; add `go-deadcode` and `govulncheck` to golangci-lint config
6. Write `TESTING.md` documenting: how to run tests locally with race detector, how to update golden files, how to manually verify clipboard clear on each platform (X11, Wayland, macOS, Windows Terminal), how to manually verify scrollback clear, known Wayland limitation (best-effort), golden file naming convention

**UAT:**
- [ ] CI matrix green on `ubuntu-latest`, `macos-latest`, `windows-latest` for every push to `main`
- [ ] `go test ./... -race` passes with zero race detector alerts on all three platforms
- [ ] All golden files match expected renders at 80×24; running `UPDATE_GOLDEN=true go test` regenerates them; running without the flag fails if they drift
- [ ] Golden file scanner finds zero sensitive-pattern matches in any `_golden*.txt` file
- [ ] `govulncheck ./...` reports no HIGH or CRITICAL vulnerabilities in the dependency tree
- [ ] Full `TestFullVaultRoundtrip` integration test passes on all three platforms under race detector

**Pitfall watch:** Golden files generated at different terminal sizes produce different output — always use `teatest.WithInitialTermSize(80, 24)` in every `teatest` test; consider adding a test helper that enforces this. Golden files must never contain real passwords, credit card numbers, or API keys — use obviously fictional values and clear comments marking them as test fixtures in the golden file scanner allowlist.

---

## Traceability Matrix

| Requirement | Phase | Description |
|-------------|-------|-------------|
| COMPAT-01 | 1 | Single static binary, no runtime dependencies |
| CI-01 | 1 | CI: build + lint + tests on every push (Linux) |
| CRYPTO-01 | 2 | AES-256-GCM + Argon2id (t=3, m=256 MiB, p=4, keyLen=32) |
| CRYPTO-02 | 2 | Crypto deps: stdlib + `golang.org/x/crypto` only |
| CRYPTO-03 | 2 | Sensitive data as `[]byte` only — never `string` |
| CRYPTO-04 | 2 | ZeroBytes primitive; memory wipe on lock/exit |
| CRYPTO-05 | 2 | mlock/VirtualLock with build-tagged platform implementations |
| CRYPTO-06 | 2 | Zero stdout/stderr logs with paths, names, or field values |
| PWD-01 | 2 | Password strength evaluator (≥12 chars + complexity rules) |
| VAULT-02 | 3 | Default seed on Create: Pasta Geral + subfolders + templates |
| SEC-05 | 3 | Observation field: auto-created, always last, immutable |
| FOLDER-01 | 3 | Create folder — domain rule: unique name within parent |
| FOLDER-02 | 3 | Rename folder — domain rule: unique name; Pasta Geral protected |
| FOLDER-03 | 3 | Move folder — cycle detection; Pasta Geral protected |
| FOLDER-04 | 3 | Reorder folder — domain: order persisted |
| FOLDER-05 | 3 | Delete folder — promote children; numeric suffix conflicts |
| TPL-01 | 3 | Create template with custom fields |
| TPL-02 | 3 | Rename template — unique among templates |
| TPL-03 | 3 | Alter template structure — no effect on existing secrets |
| TPL-04 | 3 | Delete template |
| TPL-05 | 3 | Create template from secret — auto-Observation excluded |
| TPL-06 | 3 | Templates always displayed in alphabetical order — not user-reorderable |
| ATOMIC-01 | 4 | Atomic write via `.abditum.tmp` → rename; delete tmp on failure |
| ATOMIC-02 | 4 | `.bak` / `.bak2` backup chain with startup recovery |
| ATOMIC-03 | 4 | New vault: direct write to destination (no tmp) |
| ATOMIC-04 | 4 | Windows: `MoveFileEx` with `MOVEFILE_REPLACE_EXISTING` |
| COMPAT-03 | 4 | Backward compat: versioned header, migration scaffold + test harness |
| VAULT-01 | 6 | Create vault with master password + confirmation |
| VAULT-03 | 6 | Open vault from existing file with master password |
| VAULT-04 | 6 | Error classification: invalid magic / version / auth / integrity |
| VAULT-05 | 6 | Pasta Geral must exist on open — reject with error if absent |
| QUERY-01 | 7 | View folder/secret hierarchy with folder tree + recursive active-secret count per folder |
| QUERY-02 | 7 | Search by name/field name/common value/note; sensitive excluded |
| QUERY-06 | 7 | Session state indicators in tree: included/modified/deleted |
| QUERY-07 | 7 | Virtual "Favoritos" node as sibling above Pasta Geral; read-only; depth-first favorites list |
| SEC-01 | 8 | Create secret from template or blank (Observation only) |
| SEC-02 | 8 | Duplicate secret — name "X (1)" / "X (2)" — template history preserved |
| SEC-03 | 8 | Edit secret values (name, field values, Observation) |
| SEC-04 | 8 | Edit secret structure (add/rename/reorder/delete fields; Observation immutable) |
| SEC-06 | 8 | Favorite / unfavorite secret |
| SEC-07 | 8 | Mark / unmark secret for deletion (soft-delete; removed on save) |
| SEC-08 | 8 | Move secret to another folder |
| SEC-09 | 8 | Reorder secret within same folder |
| QUERY-03 | 8 | View secret: name, all fields, Observation; sensitive fields masked |
| VAULT-06 | 9 | Save vault to current path; soft-deleted secrets removed permanently |
| VAULT-07 | 9 | Save As: new path becomes current; cannot be same as current |
| VAULT-08 | 9 | External change detection before save → 3-option modal |
| VAULT-09 | 9 | Discard/Reload: confirm modal; reload from disk using session key |
| VAULT-10 | 9 | Change master password: re-verify current; saves immediately |
| VAULT-15 | 9 | Export to unencrypted JSON with risk warning |
| VAULT-16 | 9 | Import from JSON: validate Pasta Geral; merge by folder path; name-based conflict rules |
| VAULT-17 | 9 | Configure lock/reveal/clipboard timer values |
| VAULT-11 | 10 | Auto-lock after configurable inactivity; any input resets timer |
| VAULT-12 | 10 | Manual lock (`ctrl+l`) |
| VAULT-13 | 10 | On lock: zero memory, clear screen + scrollback (`\033[3J\033[2J\033[H`) |
| VAULT-14 | 10 | Exit: confirm if unsaved changes; zero memory + clear screen on any exit |
| QUERY-04 | 10 | Reveal sensitive field temporarily; auto-hide after timer |
| QUERY-05 | 10 | Copy field to clipboard; auto-clear on timer/lock/exit; Wayland best-effort |
| CI-02 | 11 | CI matrix: Windows + macOS + Linux on every push |
| COMPAT-02 | 11 | Binary runs on Windows, macOS, Linux (verified in CI matrix) |

---

## Coverage Summary

- **v1 requirements:** 60
- **Phases:** 11
- **All requirements mapped:** ✓

| Phase | Requirements Covered | Count |
|-------|---------------------|-------|
| 1 | COMPAT-01, CI-01 | 2 |
| 2 | CRYPTO-01–06, PWD-01 | 7 |
| 3 | VAULT-02, SEC-05, FOLDER-01–05, TPL-01–06 | 13 |
| 4 | ATOMIC-01–04, COMPAT-03 | 5 |
| 5 | *(foundational — no direct user requirement)* | 0 |
| 6 | VAULT-01, VAULT-03–05 | 4 |
| 7 | QUERY-01–02, QUERY-06–07 | 4 |
| 8 | SEC-01–04, SEC-06–09, QUERY-03 | 9 |
| 9 | VAULT-06–10, VAULT-15–17 | 8 |
| 10 | VAULT-11–14, QUERY-04–05 | 6 |
| 11 | CI-02, COMPAT-02 | 2 |
| **Total** | | **60** |
