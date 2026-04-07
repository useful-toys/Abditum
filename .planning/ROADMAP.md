# Roadmap: Abditum

**Milestone:** v1.0
**Goal:** Ship a security-auditable, offline, single-binary Go TUI password manager with AES-256-GCM encryption, atomic persistence, and a keyboard-driven Bubble Tea interface.
**Created:** 2012-03-27

---

## Phases

- [x] **Phase 1: Project Scaffold + CI Foundation** ï¿½ Go module, directory tree, static binary build, GitHub Actions CI green on Linux
 (completed 2012-03-29)
- [x] **Phase 2: Crypto Package** ï¿½ Argon2id key derivation, AES-256-GCM AEAD, memory wipe primitives, mlock, password strength evaluator
 (completed 2012-03-29)
- [x] **Phase 3: Vault Domain + Manager** ï¿½ All entity types, full Manager API, business rules, invariant enforcement ï¿½ verified by unit tests
 (completed 2012-03-29)
- [x] **Phase 4: Storage Package** ï¿½ Binary `.abditum` format, atomic save with `.bak`/`.bak2` chain, Windows MoveFileEx, migration scaffold, startup recovery
 (completed 2012-03-30)
- [x] **Phase 5: TUI Scaffold + Root Model**
 (completed 2026-04-06) â€” Session state machine, root model, global tick, timer fields, modal overlay â€” no screens yet
- [x] **Phase 05.1: 05-tui-scaffold-root-model-fix**
 (completed 2026-04-06)
- [x] **Phase 05.2: tui-scaffold-message-arch**
 (completed 2026-04-06)
- [x] **Phase 05.2.1: tui-scaffold-message-arch-fixes**
 (completed 2026-04-06)
- [x] **Phase 05.2.1.1: tui-scaffold-message-arch-fixes**
 (completed 2026-04-06)
- [x] **Phase 05.2.2: tui-scaffold-message-arch-fixes**
 (completed 2026-04-06)
- [x] **Phase 05.3: merge-poc-to-app**
 (completed 2026-04-06) â€” Transformar abditum.exe na PoC standalone (sem cofre, sem domÃ­nio, apenas demonstraÃ§Ã£o de componentes TUI) e remover cmd/poc-mensagens
- [x] **Phase 05.4: 05 tui-scaffold-action-arch**
 (completed 2026-04-06)
- [x] **Phase 05.6: tui-scaffold-dialogs**
 (completed 2026-04-06)
- [x] **Phase 05.7: Golden test architecture for TUI modals**
 (completed 2026-04-06)
- [x] **Phase 6: Welcome Screen + Vault Create/Open** (completed 2026-04-06)
---

## Phase Details

### Phase 1: Project Scaffold + CI Foundation

**Goal:** The repo compiles as a static binary (`CGO_ENABLED=0`), all Charm import paths are the canonical `charm.land/*` paths, and CI executes build + lint + tests on every push.

**Requirements:** COMPAT-01, CI-01

**Plans:** 3 plans in 1 wave

Plans:
- [x] 01-01-PLAN.md ï¿½ Initialize Go module with dependencies and create directory structure with package stubs
- [x] 01-02-PLAN.md ï¿½ Configure CI workflow and Makefile for automated building, linting, and testing
- [x] 01-03-PLAN.md ï¿½ Configure golangci-lint with security rules and verify forbidden pattern detection

**UAT:**
- [ ] `CGO_ENABLED=0 go build ./cmd/abditum` succeeds and produces an executable binary
- [ ] On Linux: `file ./abditum` reports "statically linked" (no libc dependency)
- [ ] CI workflow triggers on push to `main` and reports all jobs green
- [ ] All `.go` files use `charm.land/*` import paths ï¿½ zero `github.com/charmbracelet/*` imports present
- [ ] `go vet ./...` passes with zero findings on the empty stubs

**Pitfall watch:** Import paths migrated in late 2025 ï¿½ use `charm.land/bubbletea/v2` (NOT `github.com/charmbracelet/bubbletea`). Set `CGO_ENABLED=0` as an environment variable in CI, not only at build time ï¿½ it must be enforced globally so all subsequent `go test` runs also use static linking.

---

### Phase 2: Crypto Package

**Goal:** `internal/crypto` delivers production-ready Argon2id key derivation, AES-256-GCM authenticated encryption, secure memory primitives, and password strength evaluation ï¿½ all verified by tests that will catch any future security regression.

**Requirements:** CRYPTO-01, CRYPTO-02, CRYPTO-03, CRYPTO-04, CRYPTO-05, CRYPTO-06, PWD-01

**Plans:** 1/1 plans complete

Plans:
- [x] 02-01-PLAN.md ï¿½ Complete crypto package: Argon2id KDF, AES-256-GCM AEAD, memory security primitives, password strength evaluation, comprehensive tests

**UAT:**
- [ ] `Encrypt(key, p)` + `Decrypt(key, c)` roundtrip returns original plaintext byte-for-byte
- [ ] Encrypting identical plaintext twice with the same key produces two distinct ciphertexts (nonce not reused)
- [ ] `Decrypt` with wrong key returns `crypto.ErrAuthFailed` ï¿½ no panic, no internal error string exposed
- [ ] `ZeroBytes(b)` fills every byte of `b` with `0x00`
- [ ] `EvaluatePasswordStrength([]byte("Abc1!Abc1!12"))` returns `StrengthStrong`; `"abc123"` returns `StrengthWeak`
- [ ] `CGO_ENABLED=0 go test ./internal/crypto/... -race` passes with zero race conditions

**Pitfall watch:** Argon2id `m` parameter is in **KiB** ï¿½ 256 MiB = `262144`, not `256`. `string` for sensitive data burns you: every API in this package accepts `[]byte` for passwords and keys ï¿½ establish this at the function signature level. Nonce must be generated via `io.ReadFull(rand.Reader, nonce)` immediately before each `gcm.Seal()` call ï¿½ never reuse a nonce slice across calls. mlock/VirtualLock failure is best-effort ï¿½ always check the error and continue, never fatal.

---

### Phase 3: Vault Domain + Manager

**Goal:** `internal/vault` delivers a complete, fully-tested in-memory domain layer ï¿½ all entity types, the full Manager API, and every business rule enforced and verified via unit tests ï¿½ before any file I/O or TUI code depends on it.

**Requirements:** VAULT-02, SEC-05, FOLDER-01, FOLDER-02, FOLDER-03, FOLDER-04, FOLDER-05, TPL-01, TPL-02, TPL-03, TPL-04, TPL-05, TPL-06

**Plans:** 7/7 plans complete

Plans:
- [x] 03-01-PLAN.md ï¿½ Domain Entities + Factory: Define Cofre, Pasta, Segredo, Modelo, Campo types, EstadoSessao enum, factory functions, error types, comprehensive entity tests
- [x] 03-02-PLAN.md ï¿½ Manager + Cofre Lifecycle: Manager struct, Lock with memory wiping, Salvar with atomic two-phase commit, timestamp tracking, lifecycle tests
- [x] 03-03-PLAN.md ï¿½ Folder Management: CriarPasta, RenomearPasta, MoverPasta with cycle detection, ExcluirPasta with promotion and conflict resolution, Pasta Geral protection
- [x] 03-04-PLAN.md ï¿½ Template Management: CriarModelo, RenomearModelo, ExcluirModelo with in-use check, field operations (add/remove/reorder), alphabetical ordering, "Observaï¿½ï¿½o" prohibition
- [x] 03-05-PLAN.md ï¿½ Secret Lifecycle + State Machine: CriarSegredo, ExcluirSegredo, RestaurarSegredo, FavoritarSegredo, DuplicarSegredo with name progression, estadoSessao transitions
- [x] 03-06-PLAN.md ï¿½ Secret CRUD + Structure: RenomearSegredo, EditarCampoSegredo, EditarObservacao (separate field), MoverSegredo, ReposicionarSegredo with estadoSessao tracking
- [x] 03-07-PLAN.md ï¿½ Search + Favorites + Comprehensive Validation: Buscar with sensitive field exclusion (QUERY-02), ListarFavoritos with DFS, integration tests, UAT validation, final package documentation

**UAT:**
- [ ] `Manager.Create` initializes vault with Pasta Geral + "Sites e Apps" + "Financeiro" subfolders + 3 default templates (Login, Cartï¿½o de Crï¿½dito, Chave de API)
- [ ] `Manager.CreateFolder` returns error when name already exists in the target parent
- [ ] `Manager.MoveFolder` returns error when moving a folder into one of its own descendants (cycle)
- [ ] `Manager.DeleteFolder` promotes children (including `StateDeleted` secrets, which retain their `StateDeleted` state) and resolves name conflicts with numeric suffix; Pasta Geral deletion returns explicit error
- [ ] `Manager.CreateSecret` always produces a secret with Observation as the last field; `UpdateSecretStructure` cannot rename, reorder, or delete the Observation field
- [ ] `DuplicateSecret("X")` produces "X (1)"; duplicating "X (1)" produces "X (2)"
- [ ] `CreateTemplateFromSecret` excludes ALL fields named 'Observaï¿½ï¿½o' (both auto-Observation and any user-named field with that name); a template produced from such a secret never contains a field named 'Observaï¿½ï¿½o'
- [ ] `UpdateTemplateStructure` returns explicit error when attempting to add or rename a field to 'Observaï¿½ï¿½o'
- [ ] Search function called with a string present only in a `FieldTypeSensitive` field **value** returns zero results; search called with the **name** of a sensitive field returns secrets containing that field
- [ ] `CreateSecret` produces secret with `sessionState = StateIncluded`; `UpdateSecret` on a `StateOriginal` secret ? `StateModified`; `UpdateSecret` on `StateIncluded` ? remains `StateIncluded`; `SoftDeleteSecret` ? `StateDeleted` (stores `prevState`); `RestoreSecret` ? restores previous state
- [ ] `StateDeleted` secret is excluded from search results; at Open/Discard all loaded secrets have `sessionState = StateOriginal`
- [ ] Templates are always returned in alphabetical order by name regardless of creation order; creating "Z-Template", "A-Template", "M-Template" in sequence returns them as A-Template, M-Template, Z-Template (TPL-06)

**Pitfall watch:** Pasta Geral protection requires the Manager to store the root folder's NanoID and check it in every destructive method. NanoID must use `crypto/rand` ï¿½ verify `go-nanoid`'s `gonanoid.New()` uses the secure alphabet by default. Cycle detection must walk the FULL ancestor chain, not just immediate parent. `Manager.Vault()` must return a snapshot ï¿½ exposing a live pointer to domain state allows TUI bugs to corrupt the domain silently.

---

### Phase 4: Storage Package

**Goal:** `internal/storage` delivers the binary `.abditum` format with atomic writes, `.bak`/`.bak2` backup chain, startup orphan recovery, external change detection, and a migration scaffold ï¿½ including the Windows-specific `MoveFileEx` rename path ï¿½ fully verified.

**Requirements:** ATOMIC-01, ATOMIC-02, ATOMIC-03, ATOMIC-04, COMPAT-03

**Plans:** 4 plans in 3 waves

Plans:
- [x] 04-01-PLAN.md ï¿½ Foundation: AAD-aware crypto, vault JSON serialization, storage format constants
- [x] 04-02-PLAN.md ï¿½ Core I/O: Save/SaveNew/Load with atomic .tmp protocol, platform-specific rename
- [x] 04-03-PLAN.md ï¿½ Recovery & Migration: RecoverOrphans, DetectExternalChange, migration scaffold
- [x] 04-04-PLAN.md ï¿½ Integration: FileRepository adapter implementing RepositorioCofre, end-to-end tests

**UAT:**
- [ ] `SaveNew` + `Load` roundtrip returns an identical `*vault.Vault` (deep equality on all fields)
- [ ] File with wrong magic returns `ErrInvalidMagic` ï¿½ no panic or internal error string exposed to caller
- [ ] File with correct magic but version `> currentVersion` returns `ErrVersionTooNew`
- [ ] File with correct header but wrong password returns `ErrAuthFailed`
- [ ] After second `Save` on same path: `.abditum.bak` exists; after third `Save`: old `.bak` becomes `.bak2` and new `.bak` is created
- [ ] Simulated crash (`.tmp` present, target missing): `RecoverOrphans` removes stale `.tmp` on next startup
- [ ] `FileMetadata` `mtime`/`size` returned from `Load`; modifying the file externally causes `DetectExternalChange` to return `true`
- [ ] v0 test fixture file loads cleanly via `Migrate` and passes validation at v1

**Pitfall watch:** Windows `os.Rename` is NOT atomic when replacing a file in use ï¿½ must use `MoveFileEx` with `MOVEFILE_REPLACE_EXISTING`. The `.tmp` file MUST be written to `filepath.Dir(vaultPath)`, never `os.TempDir()` ï¿½ cross-device rename (`EXDEV`) is a real failure on encrypted home directories and network mounts. Argon2id parameters (t, m, p, salt) MUST be stored in the plaintext file header and read back at decrypt time ï¿½ hardcoding them at decrypt breaks every existing vault on any future parameter change.

---

### Phase 04.1: Refinamento da camada de domï¿½nio ï¿½ encapsulamento e versioning (completed 2012-03-31)

**Goal:** Corrigir desvios de encapsulamento identificados em revisï¿½o pï¿½s-fase 4 e adicionar suporte a versï¿½es de formato no deserializador; nenhuma funcionalidade nova ï¿½ apenas alinhamento da implementaï¿½ï¿½o com os princï¿½pios arquiteturais documentados.

**Requirements:** ARCH-01, ARCH-02, BUG-01, COMPAT-03

**Depends on:** Phase 4

**Plans:** 3 plans in 2 waves

Plans:
- [x] 04.1-01-PLAN.md â€” Entity private methods: state management (marcarModificado/marcarModificacao), wipe (zerarValoresSensiveis), factory estadoSessao fix, deep copy (copiarProfundo)
- [x] 04.1-02-PLAN.md â€” Serialization versioning: DeserializarCofre(version), remove migrate.go and its tests
- [x] 04.1-03-PLAN.md â€” Manager refactoring: replace direct field access with entity methods, fix AlternarFavorito bug, add regression tests

**UAT:**
- [x] CriarSegredo returns secret with estadoSessao == EstadoIncluido (not EstadoModificado)
- [x] DuplicarSegredo returns duplicate with estadoSessao == EstadoIncluido
- [x] AlternarFavoritoSegredo does NOT update segredo.DataUltimaModificacao() â€” only cofre.DataUltimaModificacao() changes
- [x] DeserializarCofre(data, version) new signature compiles and all existing tests pass
- [x] internal/storage/migrate.go does not exist; go test ./... green

### Phase 5: TUI Scaffold + Root Model

**Goal:** The application launches, shows a placeholder, and the entire root model architecture ï¿½ session state machine, flat child-model composition, global tick, timer field declarations, and modal overlay plumbing ï¿½ is in place and race-condition-free before any screen is implemented.

**Requirements:** *(foundational infrastructure ï¿½ no single v1 requirement is independently verifiable at this phase; all subsequent TUI phases depend on this foundation)*

**Plans:** 5 plans

Plans:
- [x] 05-01-PLAN.md â€” Add charm.land v2 deps to go.mod; define core type contracts (childModel, FlowContext, workArea enum, domain messages, Cmd factory stubs)
- [x] 05-02-PLAN.md â€” Shared services: ASCII art logo, ActionManager, MessageManager, modalModel, dialog factory functions
- [x] 05-03-PLAN.md â€” All 7 child model stubs (preVaultModel, vaultTreeModel, secretDetailModel, templateListModel, templateDetailModel, settingsModel, helpModal)
- [x] 05-04-PLAN.md â€” Flow stubs (openVaultFlow, createVaultFlow) + rootModel (tea.Model, modal stack, dispatch, frame compositor)
- [x] 05-05-PLAN.md â€” Bootstrap cmd/abditum/main.go + rootModel unit tests (race-free, typed-nil safe)

**UAT:**
- [ ] `./abditum` launches without panic and renders a placeholder welcome message; `q`/`ctrl+c` exits with code 0
- [ ] `./abditum /path/to/vault.abditum` launches and passes path through to root model without crashing
- [ ] `rootModel.View()` return type compiles as `tea.View` ï¿½ not `string` (Bubble Tea v2 API)
- [ ] `CGO_ENABLED=0 go test -race ./internal/tui/...` passes with zero race detector alerts
- [ ] Sending a sequence of simulated `tickMsg` in a unit test fires appropriate timer callbacks at correct counts

**Pitfall watch:** In Bubble Tea v2, `View()` returns `tea.View` (not `string`), key messages are `tea.KeyPressMsg` (not `tea.KeyMsg`), and the space key is the string `"space"` (not `" "`). Verify these against `charm.land/bubbletea/v2` source before writing any key handler. Timer decrement must fire on `tickMsg`, NOT inside goroutines ï¿½ using `time.AfterFunc` introduces goroutine leaks and races.

---

- [x] **Phase 05.7: Golden test architecture for TUI modals**
 (completed 2026-04-06)

### Phase 05.7: Golden test architecture for TUI modals (INSERTED)

**Goal:** Implement golden test architecture for 4 TUI components: message bar, command bar, help modal, and decision dialog. Parser SGR reutilizÃ¡vel, 36 pares de golden files (72 arquivos), flag `-update` para regeneraÃ§Ã£o.
**Requirements**: TBD
**Depends on:** Phase 05
**Plans:** 9/9 plans complete
(completed 2026-04-06)

Plans:
- [x] 05.7-01-PLAN.md â€” Refactor: RenderCommandBar pure fn + helpModal decoupled from ActionManager (Wave 1)
- [x] 05.7-02-PLAN.md â€” MessageBar golden tests: 6 kinds Ã— 2 widths = 24 files (Wave 2)
- [x] 05.7-03-PLAN.md â€” CommandBar golden tests: 5 scenarios Ã— 2 widths = 20 files (Wave 2)
- [x] 05.7-04-PLAN.md â€” DecisionDialog golden tests: 10 scenarios = 20 files + 8 Update() tests (Wave 2)
- [x] 05.7-05-PLAN.md â€” HelpModal golden tests: 8 scenarios = 16 files + 8 Update() tests (Wave 2)

- [x] **Phase 05.2: tui-scaffold-message-arch**
 (completed 2026-04-06)

**Goal:** Refactor the message bar and command bar architecture â€” split MsgKind, export RenderMessageBar, extend Action with Priority/HideFromBar/Group-as-int, rewrite RenderCommandBar with spec-correct tokens and F1 right anchor â€” and validate the full system with a standalone `cmd/poc-mensagens` PoC binary.

**Requirements:** 05.2-MSG-01, 05.2-ACT-01, 05.2-INT-01, 05.2-POC-01

**Depends on:** Phase 5

**Plans:** 3/3 plans complete

Plans:
- [x] 05.2-01-PLAN.md â€” messages.go: MsgKind split (MsgInfoâ†’MsgSuccess, new MsgInfo), export TickMsg and RenderMessageBar; messages_test.go: rename MsgInfoâ†’MsgSuccess
- [x] 05.2-02-PLAN.md â€” actions.go: Action.Group int+Priority+HideFromBar, Visible() sort, RegisterGroupLabel, RenderCommandBar spec rewrite; root.go: f1 key+Group 1+Priority; help.go: int grouping with labels
- [x] 05.2-03-PLAN.md â€” cmd/poc-mensagens: standalone PoC binary with all 15 actions, live tick, RenderMessageBar+RenderCommandBar demonstration

- [x] **Phase 05.2.2: tui-scaffold-message-arch-fixes**
 (completed 2026-04-06)

**Goal:** Fix residual bugs in message bar rendering (root.go not using RenderMessageBar, uncolored text, no truncation) and help modal (F1 stacking, ESC not closing, missing bottom action bar).
**Requirements**: TBD
**Depends on:** Phase 5.2
**Plans:** 3/3 plans complete

Plans:
- [x] 05.2.2-01-PLAN.md â€” messages.go: semantic text coloring + truncation; root.go: wire RenderMessageBar()
- [x] 05.2.2-02-PLAN.md â€” root.go: fix F1/ESC dispatch order; help.go: add bottom action bar
- [x] 05.2.2-03-PLAN.md â€” Gap closure: verify error icon (âœ• U+2715) and spinner frame order (â— â—“ â—‘ â—’) match design system spec

- [x] **Phase 05.2.1: tui-scaffold-message-arch-fixes**
 (completed 2026-04-06)

**Goal:** Fix residual bugs discovered after Phase 05.2 â€” help modal full dialog rewrite, command bar truncation, spinner frame order, PoC modal pattern alignment.
**Requirements**: TBD
**Depends on:** Phase 5.2
**Plans:** 2 plans in 1 wave

Plans:
- [ ] 05.2.1-01-PLAN.md â€” Help modal: full DS dialog rewrite (Portuguese title, scroll support, DS tokens, action bar)
- [ ] 05.2.1-02-PLAN.md â€” Command bar truncation, spinner frame fix, PoC modal pattern

- [x] **Phase 05.2.1.1: tui-scaffold-message-arch-fixes**
 (completed 2026-04-06)

**Goal:** [Urgent work - to be planned]
**Requirements**: TBD
**Depends on:** Phase 5.2.1
**Plans:** 0/0 plans complete
(completed 2026-04-06)

- [x] **Phase 05.1: 05-tui-scaffold-root-model-fix**
 (completed 2026-04-06)

**Goal:** Realign internal/tui contracts with 	ui-elm-architecture.md before Phase 6 builds real screens. Rewrites childModel(3), modalView, flowHandler+Init, ActionManager owner API, MessageManager with MsgKind/TTL, dialogs expansion, and root.go dispatch ï¿½ eliminating FlowRegistry/FlowContext/flowDescriptor which Phase 5 introduced but the canonical architecture abandoned.
**Requirements**: TBD
**Depends on:** Phase 5
**Plans:** 2 plans in 2 waves

Plans:
- [ ] 05.1-01-PLAN.md ï¿½ Interface contracts + service rewrites + modal/child migrations (flows.go, actions.go, messages.go, dialogs.go, modal.go, help.go, welcome.go, child stubs)
- [ ] 05.1-02-PLAN.md ï¿½ root.go complete rewrite + root_test.go (restores compilation, all tests pass)

### Phase 6: Welcome Screen + Vault Create/Open

**Goal:** A user can pick a vault file path, create a new vault with a master password (with real-time strength feedback), or open an existing vault ï¿½ and every error case surfaces the correct generic user message with no technical detail leaked.

**Requirements:** VAULT-01, VAULT-03, VAULT-04

**Plans:** 5/5 plans complete
- [x] 06-01-PLAN.md â€” Theme System, Header, Welcome Screen + Theme Toggle
- [x] 06-02-PLAN.md â€” File Picker Modal (two-panel navigation, filtering, metadata)
- [x] 06-03-PLAN.md â€” Password Entry/Creation Modals (masked input, strength meter, attempt counter)
- [x] 06-04-PLAN.md â€” Vault Lifecycle Flows (Open, Create, CLI fast-path, exit flow, error handling)
- [x] 06-05-PLAN.md â€” Dialog Factories & Golden Tests (error dialogs, comprehensive UI tests)

**UAT:**
- [ ] User enters path + matching passwords ? vault file created on disk; TUI transitions to vault tree screen
- [ ] Strength badge shows "Weak" for `"abc123"`; shows "Strong" for `"Abc1!Abc1!12"`; creation is not blocked on Weak (only warned)
- [ ] Mismatched confirm password shows inline error message; no vault is created until they match
- [ ] Opening with correct password transitions to vault tree; opening with wrong password shows "Incorrect password ï¿½ please try again" and clears the password field for retry
- [ ] Opening a non-`.abditum` binary file shows "Invalid file ï¿½ not an Abditum vault" ï¿½ no Go error string visible
- [ ] Opening a vault from a newer format version shows "This vault was created by a newer version"

**Pitfall watch:** `textinput.Value()` returns `string` ï¿½ convert to `[]byte` inside `Update()` the instant the user submits, zero the textinput buffer (`ti.SetValue("")`), and never pass the password as `string` further down the call stack. Do not show file paths, Go error messages (`err.Error()`), or internal error types in any user-facing error message.

---

### Phase 06.3: Reimplement file picker (INSERTED)

**Goal:** Deliver a spec-compliant two-panel file picker modal (`filePickerModal`) in `internal/tui/filepicker.go`, replacing the broken stub in `dialogs.go` with lazy tree navigation, Open/Save modes, correct metadata format, MessageManager wiring, and golden test coverage.
**Requirements**: TBD
**Depends on:** Phase 06
**Plans:** 4/5 plans executed

Plans:
- [x] 06.3-01-PLAN.md — Foundation: filepicker.go skeleton + dialogs.go cleanup + token constants
- [x] 06.3-02-PLAN.md — Core Logic: Init(), tree building, Update() full keyboard handling
- [x] 06.3-03-PLAN.md — View Rendering: borders, panels, scroll indicators, Save mode section
- [x] 06.3-04-PLAN.md — Integration: 3 flow call sites + existing test fixes + 18 behavioral tests
- [ ] 06.3-05-PLAN.md — Golden Tests: 8 golden pairs (timeFmt injection, deterministic fixtures)

### Phase 06.2: adequacao-design-system (INSERTED)

**Goal:** Corrigir 10 desvios de spec nos diÃ¡logos de decisÃ£o/reconhecimento dos fluxos jÃ¡ implementados (sair, criar cofre, abrir cofre, salvar e sair) para conformidade total com o CatÃ¡logo de DiÃ¡logos da spec; implementar Fluxo 6 (bloqueio emergencial Ctrl+Alt+Shift+Q); e corrigir a formataÃ§Ã£o de teclas no modal de ajuda.
**Requirements:** DS-DIALOG-01 through DS-DIALOG-09, DS-HELP-01, FLOW-6
**Depends on:** Phase 06.1
**Plans:** 4/4 plans complete

Plans:
- [x] 06.2-01-PLAN.md â€” root.go: corrigir diÃ¡logos Ctrl+Q (desvios 1 e 2) + implementar Fluxo 6 (Ctrl+Alt+Shift+Q)
- [ ] 06.2-02-PLAN.md â€” flow_open_vault.go: dirty-check Decision + erros senha/arquivo Acknowledge (desvios 3, 4, 5)
- [x] 06.2-03-PLAN.md â€” flow_create_vault.go + flow_save_and_exit.go: desvios 6, 7, 8, 9 (dirty-check, overwrite, senha fraca, conflito externo + N Salvar como novo)
- [x] 06.2-04-PLAN.md â€” help.go: formatKeyForHelp() + aplicar em buildContentLines() (desvio 10)

### Phase 7: Vault Tree + Search

**Goal:** Users can navigate the complete folder/secret hierarchy, expand/collapse any folder, instantly filter by name and common-field content, and see favorites and soft-deleted secrets with distinct visual treatments.

**Requirements:** QUERY-01, QUERY-02, QUERY-06, QUERY-07

**Plans:**
1. Implement custom tree renderer (`tree.go`) ï¿½ **do NOT use `bubbles/list`** (insufficient for recursive folder+secret hierarchy); recursive `renderNode(folder *vault.Folder, depth int, expandState map[string]bool)` producing indented lines with box-drawing characters; folders shown before secrets at each level; folder lines: fold indicator (`?`/`?`) + folder name + recursive active-secret count `(N)` (counts all non-`StateDeleted` secrets in folder and all its subfolders recursively ï¿½ QUERY-01); secret lines: name + optional favorite `?` prefix + optional session-state indicator (`+` for `StateIncluded`, `~` for `StateModified`, `?` for `StateDeleted` ï¿½ QUERY-06) + optional dim/strikethrough for `StateDeleted`; render virtual "Favoritos" node as a fixed sibling above Pasta Geral ï¿½ always visible, never expandable past its own contents, fold indicator follows same `?`/`?` convention; Favoritos always expanded by default (QUERY-07)
2. Implement keyboard navigation in `vaultModel`: `j`/`?` (next item), `k`/`?` (prev item), `l`/`?` or Enter-on-folder (expand folder), `h`/`?` (collapse folder), Enter-on-secret (open secret detail ? `stateSecretDetail`); maintain `cursor int` into a flattened renderable item list rebuilt on each tree mutation; maintain `expandState map[string]bool` (default: Favoritos expanded, top-level real folders expanded, subfolders collapsed); cursor can reach the Favoritos virtual node and its secret entries ï¿½ Enter on a Favoritos entry navigates to that secret in detail view; pressing `n`/`N`/`d`/`m` while cursor is on the Favoritos node or its entries is a no-op (read-only); `n` (new secret ï¿½ trigger create flow from Phase 8); `N` (new folder); `/` (open search overlay)
3. Implement session-state and favorite visual treatment: `lipgloss` style `Strikethrough(true).Faint(true)` for `StateDeleted` secrets; `?` prefix for favorites; `StateDeleted` items shown in tree but visually distinguished; `StateIncluded` and `StateModified` secrets show inline session-state indicator (QUERY-06); search must exclude `StateDeleted` secrets from results
4. Implement search overlay (`/` key): inline `textinput` at bottom of tree view; on every keystroke call `vault.Search(query string, v *vault.Vault) []*vault.Secret` ï¿½ match by: secret name, all field names (including names of sensitive fields ï¿½ field names are always searchable regardless of type), common field values, Observation value (case-insensitive, accentuation-normalized via `golang.org/x/text/unicode/norm` NFD + stripping); **never read `field.Value` when `field.Type == FieldTypeSensitive`** (QUERY-02 hard gate); ESC clears search and returns to tree; Enter on result selects secret
5. Write QUERY-02 negative test (must be written FIRST, before search implementation): construct a secret with one sensitive field whose value is `"hunter2"`; call `Search("hunter2", vault)` ? assert result is empty slice; this test must fail if search ever accidentally reads sensitive field values
6. Write golden file tests: full tree with 3-level folder nesting, search overlay with results, search overlay with zero results, tree with favorited + soft-deleted secrets visible

**UAT:**
- [ ] Vault tree displays subfolders first, then secrets within each folder, matching domain order
- [ ] Each folder line shows the recursive count of active (non-`StateDeleted`) secrets across all its subfolders
- [ ] User can expand/collapse any non-root folder; expansion state persists for the session
- [ ] Searching with `/` filters the display in real time; results exclude `StateDeleted` secrets
- [ ] Searching for a string matching a sensitive field **name** (e.g. "Senha") returns secrets containing that field
- [ ] Searching for a string that exists only in a sensitive field **value** returns zero results (hard requirement)
- [ ] Favorite secrets show `?`; `StateDeleted` secrets show dim/strikethrough; `StateIncluded` and `StateModified` secrets show distinct session-state indicator
- [ ] Pressing Enter on a secret transitions to secret detail view (Phase 8) with that secret's ID
- [ ] "Favoritos" virtual node appears above Pasta Geral in the tree, expanded by default; entering it shows all `favorito = true` secrets; Enter on a Favoritos entry opens that secret in detail view; no write operations (create/move/delete) are available from within the Favoritos node

**Pitfall watch:** `bubbles/list` cannot render a recursive folder+secret tree with expand/collapse ï¿½ implement a custom renderer. Accentuation normalization for QUERY-02 requires `golang.org/x/text/unicode/norm` NFD decomposition (not just `strings.ToLower`). **Write the sensitive-field negative test FIRST, before the search loop** ï¿½ it is easy to accidentally iterate all field values without the `IsSensitive` guard.

---

### Phase 8: Secret Detail View + Edit Flows

**Goal:** Users have complete control over individual secrets: viewing (with sensitive fields masked), creating from template, editing values, editing structure, duplicating, moving, marking for deletion, reordering, favoriting ï¿½ plus full template management (create/rename/edit/delete/create-from-secret).

**Requirements:** SEC-01, SEC-02, SEC-03, SEC-04, SEC-06, SEC-07, SEC-08, SEC-09, QUERY-03

**Plans:**
1. Implement `secretDetailModel` view mode (QUERY-03): render secret name (bold), all fields in declared order with name ? value layout; sensitive fields masked as `ï¿½ï¿½ï¿½ï¿½ï¿½ï¿½ï¿½ï¿½` by default; Observation rendered last always; key-hint footer: `e` (edit values), `E` (edit structure), `d` (toggle soft-delete), `f` (favorite), `m` (move), `r` (reorder), `ctrl+d` (duplicate), `t` (template management), `v` (reveal field ï¿½ Phase 10), `c` (copy field ï¿½ Phase 10), `ESC` (back to tree), `?` (help ï¿½ Phase 11)
2. Implement create-secret flow (SEC-01): from tree view `n` key ? template picker overlay (scrollable from `bubbles/v2` list or custom); templates shown by name; "No template (Observation only)" option at top; folder picker (flat list of all folders); on confirm ? `Manager.CreateSecret(folderID, templateID)` ? transition to edit-values mode with new secret loaded
3. Implement edit-values mode (SEC-03): one `textinput` per field; sensitive fields use `textinput.EchoPassword` mode; Observation at bottom, also editable; `ctrl+s` ? batch `Manager.UpdateSecret(id, SecretChanges{Name: ..., FieldValues: ...})`; `ESC` with changes ? confirm-discard modal; `ESC` without changes ? return to view mode; password inputs: convert `textinput.Value()` ? `[]byte` immediately after submit, zero textinput buffer; leave a comment stub `// TODO(v2): "Generate ?" action for sensitive fields` at the sensitive field input site (v2 generator slot)
4. Implement edit-structure mode (SEC-04): enter with capital `E`; shows field list; controls per field: rename (inline input), reorder (`shift+j`/`shift+k`), delete (confirm modal); add field button at bottom: name input + type picker (Common / Sensitive); type cannot be changed post-creation; Observation row is rendered in a read-only non-interactive style with label "(auto ï¿½ cannot be modified)" ï¿½ no edit controls; `ctrl+s` ? `Manager.UpdateSecretStructure`
5. Implement duplicate (SEC-02): `ctrl+d` ? `Manager.DuplicateSecret(id)` ? stay in tree with new item selected; implement move (SEC-08): `m` ? folder picker overlay ? `Manager.MoveSecret(id, destFolderID)`; implement reorder (SEC-09): `r` ? reorder mode ï¿½ highlight item, `j`/`k` move it within parent folder, Enter confirm ? `Manager.ReorderSecret(id, newIndex)`, `ESC` cancel
6. Implement favorite toggle (SEC-06): `f` ? `Manager.FavoriteSecret(id, !current)` ? re-render view; implement soft-delete toggle (SEC-07): `d` ? if Active ? confirm-delete modal ? `Manager.SoftDeleteSecret(id)` ? view shows "(marked for deletion)" badge; if SoftDeleted ? `Manager.RestoreSecret(id)` ? badge removed; no confirmation needed to restore
7. Implement template management overlay (SEC-05's inverse ï¿½ TPL surface): `t` from any vault tree or detail view ? template list overlay; from list: Create (name input + add field loop), Rename (inline input), Edit Structure (same add/rename/reorder/delete as secret structure), Delete (confirm modal), Create From Secret (picks a secret from tree ï¿½ calls `Manager.CreateTemplateFromSecret`); all operations call respective `Manager.Template*` methods

**UAT:**
- [ ] Creating a secret from "Login" template produces fields: URL, Usuï¿½rio, Senha (sensitive), Observaï¿½ï¿½o (auto, last)
- [ ] Creating a secret with no template produces only Observaï¿½ï¿½o field
- [ ] Edit-values mode: Senha field is masked during editing (EchoPassword); submitting calls Manager and sets `IsModified() == true`
- [ ] Edit-structure mode: Observation row has no edit controls; adding a field appears before Observation; deleting a non-Observation field removes it; Observation cannot be deleted or moved
- [ ] Duplicating "MySecret" produces "MySecret (1)" immediately after original in tree
- [ ] Marking for deletion shows visual badge; `IsModified()` becomes true; toggling restores and clears flag
- [ ] Template management: creating a template + creating a secret from it populates fields correctly; deleting template does not affect existing secrets built from it

**Pitfall watch:** Sensitive field values must NEVER be passed to `textinput.SetValue()` for pre-population in edit mode ï¿½ this would copy them into a `string` that cannot be zeroed. Instead, leave the field blank and only populate on explicit user action. Comment the v2 generator slot at every sensitive `textinput` construction site so it is discoverable in Phase v2 without a codebase-wide search.

---

### Phase 9: Vault Lifecycle Operations

**Goal:** Users can save, save-as, discard/reload, change master password, export, import, and configure timer settings ï¿½ all vault-level operations wired through Manager with complete external-change detection and error handling.

**Requirements:** VAULT-06, VAULT-07, VAULT-08, VAULT-09, VAULT-10, VAULT-15, VAULT-16, VAULT-17

**Plans:**
1. Add key bindings in vault tree and detail views: `ctrl+s` (Save), `ctrl+shift+s` (Save As), `ctrl+r` (Discard/Reload), `ctrl+p` (Change Master Password), `ctrl+e` (Export), `ctrl+i` (Import), `ctrl+,` (Settings); each binding dispatches a typed `tea.Msg` to `rootModel` which handles all lifecycle operations (keeps TUI screens decoupled from file operations)
2. Implement Save flow (VAULT-06): `rootModel` receives `SaveMsg` ? call `Manager.Save()` (which calls `storage.Save` internally); if `storage.DetectExternalChange` returns true before saving, display modal: "This vault was modified externally.\n\nOverwrite / Save as new file / Cancel"; mapping: 0 ? force save, 1 ? trigger SaveAs flow, -1 ? cancel (VAULT-08)
3. Implement Save As flow (VAULT-07): modal with path `textinput`; validate path is not identical to `Manager.CurrentPath()` (reject with error message); on confirm ? `Manager.SaveAs(destPath)` ? `CurrentPath()` now returns `destPath`; IsModified becomes false
4. Implement Discard/Reload flow (VAULT-09): confirm modal "Discard all unsaved changes and reload from disk?"; if file was externally modified, add secondary warning line; on confirm ? `Manager.Discard()` (reloads from disk using current session key stored in Manager); tree re-renders with refreshed snapshot
5. Implement Change Master Password flow (VAULT-10): three-step modal sequence: step 1 = current password input (re-derive key in-process to verify; reject if wrong), step 2 = new password + strength badge (same as create flow), step 3 = confirm new password (must match); on all three confirmed ? `Manager.ChangeMasterPassword(current, next []byte)` which immediately re-encrypts and saves; operation is irreversible ï¿½ show "This saves immediately and cannot be undone" warning before step 1
6. Implement Export flow (VAULT-15): path input modal + risk warning: "The exported file is unencrypted. Store it somewhere safe and delete it when done."; on confirm ? `Manager.Export(path)` writes JSON: all active (non-soft-deleted) secrets including their timestamps (`data_criacao`/`data_ultima_modificacao`), all folders, all templates; vault-level timer settings (`Configuracoes` block) are not exported
7. Implement Import flow (VAULT-16): path input modal; call `Manager.Import(path)` ï¿½ validate file first: if invalid JSON or Pasta Geral absent ? return generic error, abort import; merge rules: folders merged by full path (sequence of names from root) ï¿½ existing folder ? merge contents; new folder ? create; within each merged folder: incoming secret with same **name** ? **replaces** existing; incoming secret with unique name ? added; incoming template with same **name** ? **replaces** existing; incoming template with unique name ? added; show import summary modal: "Imported N secrets, M folders, K templates"
8. Implement Settings screen (VAULT-17): dedicated screen reachable from `ctrl+,`; three integer inputs: "Auto-lock after (minutes):", "Reveal sensitive field for (seconds):", "Clear clipboard after (seconds):"; all three required, positive integers only; `ctrl+s` ? `Manager.UpdateSettings(settings)` ? new values used by Phase 10 timers immediately; `ESC` ? discard unsaved settings changes (confirm modal if modified)

**UAT:**
- [ ] Save writes to current path; `.abditum.bak` created; `IsModified()` returns false after successful save
- [ ] External change detection: modifying vault file externally between Open and Save triggers the 3-option modal
- [ ] Save As to a different path changes `CurrentPath()`; subsequent Ctrl+S writes to the new path
- [ ] Discard without pending changes still confirms; reloads vault cleanly; all unsaved mutations are gone
- [ ] Change Master Password: wrong current password shows error at step 1; correct current + strong new password saves immediately; vault can be reopened with new password afterward
- [ ] Export: resulting JSON contains all active secrets including timestamps (`data_criacao`/`data_ultima_modificacao`); soft-deleted secrets absent; `Configuracoes` timer-settings block absent
- [ ] Import from invalid JSON or file without Pasta Geral ? fails with generic error message, no changes applied to vault
- [ ] Import from valid Abditum JSON: folders merged by path; incoming secret with same name as existing secret in that folder replaces it; incoming secret with unique name is added; import summary modal shows correct counts
- [ ] Settings: changing auto-lock to 1 minute and waiting 1 minute triggers auto-lock (verified in Phase 10)

---

### Phase 10: Security Timers + Clipboard + Lock/Exit

**Goal:** All time-driven security behaviors work correctly and synchronously on all three platforms: auto-lock, manual lock, clipboard auto-clear on timer and on lock/exit, sensitive field auto-hide on timer, memory wipe, and screen+scrollback clear ï¿½ with clean exit handling for all pending-change scenarios.

**Requirements:** VAULT-11, VAULT-12, VAULT-13, VAULT-14, QUERY-04, QUERY-05

**Plans:**
1. Wire auto-lock timer (VAULT-11): every `tea.KeyPressMsg` or `tea.MouseMsg` in `rootModel.Update` resets `lockTimer` to `Manager.Settings().LockTimeoutMinutes * 60`; global tick decrements; at 0 ? emit `LockMsg{}` internally; test that ALL input varieties (key, mouse, window resize) reset the timer ï¿½ missing any input type breaks the invariant
2. Implement `LockMsg` handler in `rootModel` (VAULT-12, VAULT-13): call `Manager.Lock()` (zeros in-memory key and all field `[]byte` buffers); set `rootModel.state = stateLocked`; zero all `revealedFields` map entries; clear clipboard synchronously (before any further rendering); emit `clearScreenCmd` which outputs `"\033[3J\033[2J\033[H"` to stdout AFTER Bubble Tea program stops its render loop (`p.ReleaseTerminal()` or equivalent); transition welcome model to re-enter-password sub-state for same path
3. Implement manual lock keybinding `ctrl+l` from any active vault screen ? dispatch same `LockMsg{}` as auto-lock (VAULT-12); lock must work from vault tree, secret detail, and any overlay state
4. Implement clean exit (VAULT-14): `q` from welcome / `ctrl+q` from vault screens ? check `Manager.IsModified()`; if no pending changes ? proceed directly to exit; if pending ? modal "Save and Quit / Discard and Quit / Cancel"; on any exit path: `Manager.Lock()` (zeros memory), clear clipboard synchronously, emit `clearScreenCmd`, `tea.Quit`; exit code 0 on clean exit, 1 on fatal startup error only
5. Implement clipboard copy (QUERY-05): `c` key in secret detail view ? cursor selects field ? call `clipboard.WriteAll(string(field.Value))` via `github.com/atotto/clipboard`; set `clipboardTimer = Manager.Settings().ClipboardClearSeconds` and `clipboardFieldID = fieldID`; on timer expiry ? `clipboard.WriteAll("")`; on lock or exit ? `clipboard.WriteAll("")` synchronously before process returns; handle headless/Wayland/SSH: if `clipboard.WriteAll` returns error ? show gentle warning in footer, continue without crash (QUERY-05 Wayland best-effort)
6. Implement field reveal (QUERY-04): `v` key in secret detail view ? cursor selects a sensitive field ? set `revealedFields[fieldID] = time.Now().Add(time.Duration(settings.RevealSeconds) * time.Second)`; `secretDetailModel.View()` reads `rootModel.revealedFields` (passed via render args or message) to determine masking; global tick clears expired entries from map; re-triggering `v` on already-revealed field resets the timer
7. Wire Phase 9 Settings into Phase 10 timers: when `UpdateSettings` is called, update in-flight `lockTimer` cap (reset to new value if new value is smaller than current remaining); `clipboardTimer` and `revealedFields` expiry use `Manager.Settings()` at timer-set time; verify that settings screen changes take effect on next timer action without restart
8. Write synchronized timer tests: unit-test `rootModel` with mock tick messages; assert lock fires after exactly N ticks; assert any key between ticks resets count; assert clipboard cleared synchronously in `LockMsg` handler (not in a goroutine); assert reveal field correctly expires; run `go test -race`

**UAT:**
- [ ] Auto-lock: after the configured inactivity period with no input, vault locks automatically and terminal shows re-enter-password prompt
- [ ] Any key press or mouse event resets the inactivity timer to full duration
- [ ] Manual lock (`ctrl+l`) locks immediately from vault tree or detail view
- [ ] After lock: sensitive data is not visible in terminal scrollback buffer (manual verification in xterm, Windows Terminal, iTerm2)
- [ ] Clipboard: copying a field value and waiting `clipboardClearSeconds` clears the clipboard; copying another field resets the timer
- [ ] Clipboard is cleared synchronously when vault locks or app exits (no residual data after process exits)
- [ ] Revealing a sensitive field (`v`) shows value for `revealSeconds` then auto-hides to `ï¿½ï¿½ï¿½ï¿½ï¿½ï¿½ï¿½ï¿½`
- [ ] On Wayland without `wl-clipboard` present: clipboard operation shows a gentle footer notice; application continues normally

**Pitfall watch:** Terminal clear-screen `\033[2J` does NOT clear scrollback ï¿½ only `\033[3J\033[2J\033[H` clears both. Clipboard clear on lock/exit MUST be synchronous ï¿½ `time.AfterFunc` goroutine may not execute before `tea.Quit` unwinds; call `clipboard.WriteAll("")` directly in the `LockMsg` handler. On Wayland, `atotto/clipboard` may silently fail if `wl-clipboard` is absent ï¿½ wrap the call, check the error, show a footer notice (no crash, no log with sensitive content).

---

### Phase 11: Cross-Platform CI Matrix + Integration Tests + Golden Files

**Goal:** Every test passes on Windows, macOS, and Linux under the race detector; golden files lock all visual regressions; no sensitive data can appear in test fixtures; and dependency vulnerabilities are checked automatically.

**Requirements:** CI-02, COMPAT-02

**Plans:**
1. Expand GitHub Actions CI matrix to three OS: `ubuntu-latest`, `macos-latest`, `windows-latest`; all jobs run `CGO_ENABLED=0 go test ./... -race -count=1`; confirm Windows-specific `MoveFileEx` code path is exercised in the windows CI job explicitly (call `atomicRenameWindows` in a Windows-only integration test)
2. Write end-to-end manager integration test (in-process, no subprocess): `TestFullVaultRoundtrip` ï¿½ `Manager.Create` ? `Manager.CreateFolder` ? `Manager.CreateSecret` ? `Manager.UpdateSecret` (values) ? `Manager.FavoriteSecret` ? `Manager.Save` ? new `Manager.Open` (same path, same password) ? assert all mutations present ? assert `IsModified() == false`; run under `-race`
3. Collect all `teatest/v2` golden file tests from Phases 6-10 in CI; add any missing golden files for: locked state (re-enter password prompt), export confirmation modal, import summary modal, settings screen, help overlay; document golden file update procedure with `UPDATE_GOLDEN=true go test ./...`
4. Implement golden file secret-pattern scanner (`cmd/scan-golden/main.go` or `TestScanGoldenFiles` test): walk all files matching `testdata/**/*_golden*.txt`; apply regex patterns for common sensitive data signatures (passwords matching complexity rules, 16-digit groups, API key patterns like `sk-`, `Bearer `); fail CI if any pattern matches; use clearly fictional fixture values: `"P@ssw0rd!Test42"` (labeled fake), `"4242424242424242"` (recognizably fictional test card)
5. Add `govulncheck` step to CI (all three OS): `go install golang.org/x/vuln/cmd/govulncheck@latest && govulncheck ./...`; fail CI on any HIGH or CRITICAL severity finding; add `go-deadcode` and `govulncheck` to golangci-lint config
6. Write `TESTING.md` documenting: how to run tests locally with race detector, how to update golden files, how to manually verify clipboard clear on each platform (X11, Wayland, macOS, Windows Terminal), how to manually verify scrollback clear, known Wayland limitation (best-effort), golden file naming convention

**UAT:**
- [ ] CI matrix green on `ubuntu-latest`, `macos-latest`, `windows-latest` for every push to `main`
- [ ] `go test ./... -race` passes with zero race detector alerts on all three platforms
- [ ] All golden files match expected renders at 80ï¿½24; running `UPDATE_GOLDEN=true go test` regenerates them; running without the flag fails if they drift
- [ ] Golden file scanner finds zero sensitive-pattern matches in any `_golden*.txt` file
- [ ] `govulncheck ./...` reports no HIGH or CRITICAL vulnerabilities in the dependency tree
- [ ] Full `TestFullVaultRoundtrip` integration test passes on all three platforms under race detector

**Pitfall watch:** Golden files generated at different terminal sizes produce different output ï¿½ always use `teatest.WithInitialTermSize(80, 24)` in every `teatest` test; consider adding a test helper that enforces this. Golden files must never contain real passwords, credit card numbers, or API keys ï¿½ use obviously fictional values and clear comments marking them as test fixtures in the golden file scanner allowlist.

### Phase 12: tui-scaffold-dialogs

**Goal:** [To be planned]
**Requirements**: TBD
**Depends on:** Phase 11
**Plans:** 0/0 plans complete
(completed 2026-04-06)

Plans:
- [ ] TBD (run /gsd:plan-phase 12 to break down)

---

## Traceability Matrix

| Requirement | Phase | Description |
|-------------|-------|-------------|
| COMPAT-01 | 1 | Single static binary, no runtime dependencies |
| CI-01 | 1 | CI: build + lint + tests on every push (Linux) |
| CRYPTO-01 | 2 | AES-256-GCM + Argon2id (t=3, m=256 MiB, p=4, keyLen=32) |
| CRYPTO-02 | 2 | Crypto deps: stdlib + `golang.org/x/crypto` only |
| CRYPTO-03 | 2 | Sensitive data as `[]byte` only ï¿½ never `string` |
| CRYPTO-04 | 2 | ZeroBytes primitive; memory wipe on lock/exit |
| CRYPTO-05 | 2 | mlock/VirtualLock with build-tagged platform implementations |
| CRYPTO-06 | 2 | Zero stdout/stderr logs with paths, names, or field values |
| PWD-01 | 2 | Password strength evaluator (=12 chars + complexity rules) |
| VAULT-02 | 3 | Default seed on Create: Pasta Geral + subfolders + templates |
| SEC-05 | 3 | Observation field: auto-created, always last, immutable |
| FOLDER-01 | 3 | Create folder ï¿½ domain rule: unique name within parent |
| FOLDER-02 | 3 | Rename folder ï¿½ domain rule: unique name; Pasta Geral protected |
| FOLDER-03 | 3 | Move folder ï¿½ cycle detection; Pasta Geral protected |
| FOLDER-04 | 3 | Reorder folder ï¿½ domain: order persisted |
| FOLDER-05 | 3 | Delete folder ï¿½ promote children; numeric suffix conflicts |
| TPL-01 | 3 | Create template with custom fields |
| TPL-02 | 3 | Rename template ï¿½ unique among templates |
| TPL-03 | 3 | Alter template structure ï¿½ no effect on existing secrets |
| TPL-04 | 3 | Delete template |
| TPL-05 | 3 | Create template from secret ï¿½ auto-Observation excluded |
| TPL-06 | 3 | Templates always displayed in alphabetical order ï¿½ not user-reorderable |
| ATOMIC-01 | 4 | Atomic write via `.abditum.tmp` ? rename; delete tmp on failure |
| ATOMIC-02 | 4 | `.bak` / `.bak2` backup chain with startup recovery |
| ATOMIC-03 | 4 | New vault: direct write to destination (no tmp) |
| ATOMIC-04 | 4 | Windows: `MoveFileEx` with `MOVEFILE_REPLACE_EXISTING` |
| COMPAT-03 | 4 | Backward compat: versioned header, migration scaffold + test harness |
| VAULT-01 | 6 | Create vault with master password + confirmation |
| VAULT-03 | 6 | Open vault from existing file with master password |
| VAULT-04 | 6 | Error classification: invalid magic / version / auth / integrity |
| VAULT-05 | 6 | Pasta Geral must exist on open ï¿½ reject with error if absent |
| QUERY-01 | 7 | View folder/secret hierarchy with folder tree + recursive active-secret count per folder |
| QUERY-02 | 7 | Search by name/field name/common value/note; sensitive excluded |
| QUERY-06 | 7 | Session state indicators in tree: included/modified/deleted |
| QUERY-07 | 7 | Virtual "Favoritos" node as sibling above Pasta Geral; read-only; depth-first favorites list |
| SEC-01 | 8 | Create secret from template or blank (Observation only) |
| SEC-02 | 8 | Duplicate secret ï¿½ name "X (1)" / "X (2)" ï¿½ template history preserved |
| SEC-03 | 8 | Edit secret values (name, field values, Observation) |
| SEC-04 | 8 | Edit secret structure (add/rename/reorder/delete fields; Observation immutable) |
| SEC-06 | 8 | Favorite / unfavorite secret |
| SEC-07 | 8 | Mark / unmark secret for deletion (soft-delete; removed on save) |
| SEC-08 | 8 | Move secret to another folder |
| SEC-09 | 8 | Reorder secret within same folder |
| QUERY-03 | 8 | View secret: name, all fields, Observation; sensitive fields masked |
| VAULT-06 | 9 | Save vault to current path; soft-deleted secrets removed permanently |
| VAULT-07 | 9 | Save As: new path becomes current; cannot be same as current |
| VAULT-08 | 9 | External change detection before save ? 3-option modal |
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

### Phase 05.3: merge-poc-to-app (INSERTED)

**Goal:** Transformar `abditum.exe` na PoC standalone (sem cofre, sem domÃ­nio, apenas demonstraÃ§Ã£o de componentes TUI) e remover `cmd/poc-mensagens`.
**Requirements**: *(internal refactor + PoC merge â€” no v1 requirements)*
**Depends on:** Phase 5.2.2
**Plans:** 2 plans in 2 waves

Plans:
- [ ] 05.3-01-PLAN.md â€” Centralize tokens (colors, symbols, styles), refactor messages.go/actions.go/help.go to consume tokens
- [ ] 05.3-02-PLAN.md â€” Rewrite NewRootModel for PoC mode (15 actions, no vault.Manager), update main.go, delete cmd/poc-mensagens/

### Phase 05.4: 05 tui-scaffold-action-arch (INSERTED)

**Goal:** Validate the action architecture through targeted configuration changes â€” multi-key dispatch, selective command bar visibility, priority-based truncation, and help modal grouping â€” plus comprehensive unit tests covering RenderCommandBar edge cases.
**Requirements**: *(configuration validation + test coverage â€” no v1 requirements)*
**Depends on:** Phase 05.3
**Plans:** 2/2 plans complete

Plans:
- [ ] 05.4-01-PLAN.md â€” Update root.go action registrations (multi-key, HideFromBar, priorities, groups) + help.go group 0 skip
- [ ] 05.4-02-PLAN.md â€” 6 new unit tests: multi-key dispatch, HideFromBar visibility, RenderCommandBar truncation/anchor/narrow

### Phase 05.6: tui-scaffold-dialogs (INSERTED)

**Goal:** Fix DecisionDialog compilation errors (decision.go), complete View() with manual top-border construction (lipgloss v2 has no BorderTitle), and validate with comprehensive unit + integration tests covering all 10 severity Ã— intention combinations.
**Requirements**: *(foundational TUI infrastructure â€” prerequisite for Phase 06 flows)*
**Depends on:** Phase 05.4
**Plans:** 1/1 plans complete

Plans:
- [x] 05.6-01-PLAN â€” Fix compilation errors, rewrite View() with manual border construction, add wrapBody() helper
- [ ] 05.6-02-PLAN â€” 10 matrix fixture constructors + 3 rendering tests (matrix, symbols, border chars)
- [ ] 05.6-03-PLAN â€” 4 interaction tests (Enter/Esc/arrows/unknown) + 4 edge case tests (long body, short body, Ack no-Esc, small size)
- [ ] 05.6-04-PLAN â€” 1 integration test in root_test.go + final verification (build + vet + all tests)

**UAT:**
- [ ] `go build ./...` compiles clean (zero errors)
- [ ] `go test ./internal/tui/...` â€” 44+ tests, all pass
- [ ] ConfirmaÃ§Ã£o Ã— Destrutivo renders `âš ` symbol + `semantic.warning` border + `semantic.error` default key
- [ ] Reconhecimento renders only `Enter OK` action (no Esc)
- [ ] Enter key on any DecisionDialog triggers popModalMsg (modal stack decrements)
- [ ] All existing 33+ tui tests still passing (zero regressions)

