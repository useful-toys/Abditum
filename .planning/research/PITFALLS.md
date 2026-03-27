# Domain Pitfalls — Abditum

**Domain:** Offline TUI password manager (Go + Bubble Tea + AES-256-GCM + Argon2id)
**Researched:** 2026-03-27
**Scope:** Security-critical; every pitfall here has been observed in open-source Go password managers or cryptographic Go projects.

---

## Critical Pitfalls

Mistakes that cause silent data loss, key exposure, or require architectural rewrites.

---

### Pitfall 1: Using `string` for Sensitive Data Instead of `[]byte`

**What goes wrong:** Master password and derived AES key are stored as Go `string` values at any point in the code. Strings are immutable and interned — the runtime may keep copies in multiple heap locations, string literals land in read-only data segments, and the garbage collector never zeroes them. When you "zero" a string, you only affect one reference; unknown copies persist until GC decides to reclaim them (potentially never before a swap or crash dump).

**Why it happens:** Password prompts (Bubble Tea text inputs) produce `string`. It is tempting to pass that value directly to `crypto/argon2.IDKey`. Developers forget to convert immediately and propagate the string through function signatures.

**How to avoid:**
- Convert the raw input to `[]byte` at the earliest possible moment — inside the input handler that receives keystrokes, before returning from the Bubble Tea `Update` function.
- Never define function signatures that accept the master password or AES key as `string`. Use `[]byte` exclusively in `internal/crypto` and `internal/vault`.
- After conversion, immediately overwrite the Bubble Tea text-input buffer state (set it to empty string) so the TUI model no longer holds the plaintext.
- After use, call `crypto/subtle` or a manual loop to zero the `[]byte` before releasing it: `for i := range buf { buf[i] = 0 }`.
- Define a `SecureBytes` type alias or a small wrapper that zeros on `Close()` as a compile-time reminder.

**Warning signs:**
- Function signatures with `password string` anywhere in `internal/crypto` or `internal/vault`.
- Passing `model.PasswordInput.Value()` directly to `argon2.IDKey()` without intermediate `[]byte` conversion.
- Absence of zeroing loops in `vault.Manager.Lock()` and `vault.Manager.Close()`.

**Phase to address:** Crypto package foundation (Phase 1 / earliest implementation). Retrofit is extremely risky — establish the `[]byte`-only convention before writing any crypto or vault code.

---

### Pitfall 2: Go Garbage Collector Moving Sensitive Data Before Zeroing

**What goes wrong:** Go's GC is a moving/compacting collector (in principle; current GC does not compact, but this is not guaranteed by the spec). More practically: when a `[]byte` slice is passed to a function, Go may allocate a new backing array on the heap and copy the data. The original backing array is now an unreachable orphan — zeroing the slice you hold does NOT zero the orphan copy. This is especially subtle with `append`, slice operations, and interface boxing.

**Why it happens:** Developers assume that zeroing the slice they hold is sufficient. Go's escape analysis and GC make this false. The `mlock`/`VirtualLock` call only protects pages that back the currently live allocation, not historical copies.

**How to avoid:**
- Treat zeroing as best-effort and document it as such (the architecture document already states this — make sure the code comments repeat it).
- Minimize the number of operations that cause re-allocation: pre-allocate fixed-size buffers for the master password and key; never `append` to them in ways that could trigger reallocation.
- Use `mlock`/`VirtualLock` to make swapping impossible for the pages you DO control (see Pitfall 9).
- Never pass sensitive `[]byte` as `interface{}` or `any` — interface boxing copies the data to a new heap allocation you cannot zero.
- Avoid `fmt.Sprintf`, `strings.Builder`, or any string interpolation with sensitive values — they all produce heap-allocated strings you cannot zero.
- Do not copy sensitive values into error messages or log strings (zero-log policy).

**Warning signs:**
- `append(sensitiveSlice, ...)` usage in crypto or vault packages.
- Sensitive buffers passed as `interface{}`.
- Use of `fmt.Sprintf` or `string(sensitiveBytes)` anywhere in the hot path.

**Phase to address:** Crypto package foundation and vault Manager initialization. Design the allocation strategy upfront.

---

### Pitfall 3: GCM Nonce Reuse with a Fixed Key

**What goes wrong:** AES-256-GCM is catastrophically broken when the same (key, nonce) pair is used twice: an attacker who has two ciphertexts encrypted under the same (key, nonce) can XOR them to cancel the keystream and recover plaintext, and can also forge arbitrary messages. With a randomly-generated 96-bit nonce (standard for GCM), the birthday bound for collision is $\approx 2^{48}$ encryptions — safe for a personal vault saved thousands of times. The risk is: if a code path ever **reuses** a previously generated nonce, or if the nonce is derived from a counter stored incorrectly, the guarantee collapses.

**Why it happens:** Developers copy nonce generation from examples that use a counter, or re-use the same nonce slice across multiple seal calls, or forget to generate a fresh nonce before each `gcm.Seal()` call.

**How to avoid:**
- Always generate the nonce with `io.ReadFull(rand.Reader, nonce)` immediately before each `gcm.Seal()` call — never reuse a nonce slice across calls.
- Never derive the nonce from a counter, timestamp, or any deterministic value.
- Never store or re-read the nonce from disk and use it again for a new encryption.
- The nonce is prepended to or stored alongside the ciphertext; on decryption, it is read back from that stored location and is never used again for encryption.
- Write a unit test that encrypts the same plaintext twice and asserts the two ciphertexts are different (nonce uniqueness smoke test).

**Warning signs:**
- A nonce variable allocated outside the encrypt function and reused across calls.
- Any counter, clock, or deterministic source for nonce generation.
- Absence of `io.ReadFull(rand.Reader, nonce)` immediately before `gcm.Seal`.

**Phase to address:** Crypto package foundation. Validate with unit tests before any other component uses the crypto package.

---

### Pitfall 4: Argon2id Parameters Too Weak or Too Strong (and Not Stored in the File)

**What goes wrong (too weak):** Using default/example Argon2id parameters that are far too low — e.g., `time=1, memory=64MB, threads=1` — provides inadequate resistance to offline brute force. An attacker with the `.abditum` file and a modern GPU can try millions of passwords per second.

**What goes wrong (too strong):** Parameters so high that unlock takes 10+ seconds on low-end hardware (e.g., a Raspberry Pi or old laptop), making the app unusable, or causing the OS to OOM-kill the process on constrained systems.

**What goes wrong (not stored):** The Argon2id parameters (time, memory, threads) and the salt are not stored in the file header. When the user opens the file, the parameters cannot be recovered — re-derivation uses hardcoded defaults, which will be wrong after a future parameter upgrade, permanently locking the user out.

**How to avoid:**
- Use the OWASP-2023 recommended baseline: `time=3, memory=64MB (65536 KiB), threads=4`. Benchmark on your CI hardware to confirm < 1s unlock time. Re-evaluate for the current year (2026) — 128MB is becoming the new minimum.
- Store salt (16 bytes minimum, 32 bytes preferred), time, memory, and threads in the plaintext file header alongside the ciphertext. Never hardcode parameters at decryption time.
- Write a migration test: create a vault with version N parameters, open it with the current code, confirm it decrypts correctly.
- Add a constant `ArgonVersion = 1` to the file format so future parameter upgrades can be identified and re-keyed on open.

**Warning signs:**
- Hardcoded Argon2id parameters at the call site of `argon2.IDKey` without reading them from the file.
- Parameters copied from a blog post without benchmarking on target hardware.
- File format has no field for Argon2id parameters.
- Unlock takes > 3 seconds on a modern laptop.

**Phase to address:** Crypto + storage packages. File format design must include parameter storage before any vault files are created (format changes are breaking).

---

### Pitfall 5: Clipboard Not Cleared on All Platforms

**What goes wrong:** The clipboard-clearing code calls `xclip`, `xsel`, or `pbcopy` (Unix) but does nothing on Windows, or vice-versa. On Wayland-based Linux, `xclip` writes to the X11 clipboard, not the Wayland clipboard — the sensitive value persists in the Wayland clipboard indefinitely.

A subtler variant: the clipboard library used spawns a subprocess, but when the application exits, the subprocess is killed before it can write the empty string, leaving the sensitive value in the clipboard.

**Why it happens:** Developers test on one platform. Cross-platform clipboard libraries hide differences but may fail silently. On Linux, the distinction between X11, Wayland, and wl-clipboard is often overlooked.

**How to avoid:**
- Use `golang.design/x/clipboard` or `github.com/atotto/clipboard` — evaluate which provides better cross-platform coverage and Wayland support in 2026. Test on all three platforms in CI.
- Implement clipboard clear as: write an empty string, then verify by reading back. If the read-back is not empty and matches the sensitive value, log a generic "clipboard clear failed" warning (no value in log — zero-log policy).
- On lock/exit, clear the clipboard synchronously before the process exits. Never rely on a goroutine that may be interrupted by `os.Exit`.
- Handle the case where clipboard operations are not supported (headless/SSH sessions) gracefully — do not crash; just skip and warn.
- Test on: Linux X11, Linux Wayland, macOS, Windows.

**Warning signs:**
- Clipboard code uses `exec.Command("xclip", ...)` without a Windows/macOS branch.
- Clipboard clear is triggered via a goroutine that uses `time.AfterFunc` without a cancellation path on lock/exit.
- No CI test for clipboard clear (even a mock/integration test).

**Phase to address:** Vault lifecycle / lock-exit phase. Must be cross-platform from day one — retrofitting is painful.

---

### Pitfall 6: Terminal "Clear" Does Not Clear Scrollback Buffer

**What goes wrong:** `fmt.Print("\033[2J\033[H")` (ANSI clear screen) clears the visible terminal area but does **not** clear the scrollback buffer. A user can scroll up and see all the secrets that were displayed before locking. On Windows (ConHost / Windows Terminal), `\033[2J` behavior varies; the new Windows Terminal has a separate escape for clearing scrollback (`\033[3J`).

Calling `exec.Command("clear")` on Unix or `exec.Command("cmd", "/c", "cls")` on Windows is also unreliable — it only works if the program is attached to a real terminal, and it does not clear scrollback in most terminal emulators.

**Why it happens:** Developers test by visual inspection ("the screen looks blank") without scrolling up to verify.

**How to avoid:**
- Send the combined sequence `\033[3J\033[2J\033[H` — `3J` clears scrollback, `2J` clears visible screen, `H` moves cursor to top. This works in most modern terminal emulators (xterm, iTerm2, Windows Terminal, GNOME Terminal, Alacritty).
- On Windows, also use `windows.SetConsoleCursorPosition` and clear via the Windows Console API if ANSI sequences are not accepted (legacy `conhost.exe`).
- Bubble Tea's `tea.ClearScrollArea()` or equivalent commands should be used inside the TUI framework before the program exits, so the framework's own rendering pipeline is flushed correctly.
- Accept that 100% clearing of scrollback is not guaranteed across all terminals. Document this limitation. The design goal is best-effort.
- Add a manual test protocol: run app in xterm, reveal a secret, lock, scroll up — verify the secret is not visible.

**Warning signs:**
- Only `\033[2J` or only `exec.Command("clear")` used at lock/exit.
- No scrollback-clear test case in the test plan.
- Terminal-clear code only tested on the developer's preferred terminal emulator.

**Phase to address:** Vault lifecycle / lock-exit phase. Implement and test on all supported platforms early.

---

### Pitfall 7: Bubble Tea State Machine Growing Unbounded / Modal Stacking Bugs

**What goes wrong:** In Bubble Tea, every screen/modal is typically a separate model pushed onto a stack or composed into a parent model. Two common failure modes:
1. **Stack overflow / unbounded growth:** Every keypress during a multi-step flow pushes a new model without popping old ones. After many interactions, the model tree is enormous, causing allocation pressure and subtle state bugs.
2. **Ghost state:** Closing a modal does not reset its internal state. Reopening the modal shows stale data — e.g., a "confirm delete" dialog pre-filled with the previous item's name, or a field reveal timer still running for a previously-closed secret.

A third variant: **event propagation bugs** — a `tea.KeyMsg` intended for a child modal is also processed by the parent, causing duplicate state mutations.

**Why it happens:** The Elm architecture in Bubble Tea v2 routes all messages through a single `Update` function. Developers do not implement explicit "pop" transitions, or they forget to reset child model state on open/close transitions. Without a disciplined state machine, valid and invalid states become indistinguishable.

**How to avoid:**
- Design an explicit screen/mode enum before writing TUI code: `type AppScreen int; const (ScreenUnlocked AppScreen = iota; ScreenSecretList; ScreenSecretDetail; ...)`. Every transition must be enumerated.
- Use a flat model where only the **current** screen's model is active, not a stack of accumulated models. When transitioning away from a screen, explicitly reset its model to zero value.
- For modals (confirmation dialogs, field reveal), use a dedicated `Modal` field in the parent model that is set to `nil`  when closed and re-initialized on open — never re-use the previous modal instance.
- Ensure `KeyMsg` handling is guarded: if a modal is active, the parent must **not** process the key; return early.
- Write state machine transition tests using `teatest/v2`: assert that after closing a secret detail view and reopening a different one, no state from the first view bleeds through.

**Warning signs:**
- A slice/stack of models that grows with each navigation action.
- Modal struct fields that are initialized once at startup and mutated in place.
- `Update` function returning `m, cmd` without explicitly zeroing the old child model on transition.
- A field reveal timer (for a different secret) still ticking after navigating away.

**Phase to address:** TUI architecture design phase. Establish the flat model / enum-driven state machine pattern before building any screens.

---

### Pitfall 8: Atomic Save Edge Cases — `os.Rename` Atomicity Across Directory Boundaries

**What goes wrong:** `os.Rename` is only atomic within the same filesystem and same directory (on Linux/macOS POSIX rename semantics). Writing `.abditum.tmp` to `/tmp/` and renaming to `/home/user/vault.abditum` crosses a filesystem boundary and will fail with `EXDEV` (invalid cross-device link). More subtly: writing to a different directory than the target (e.g., the current working directory vs. the vault file's directory) causes the same problem.

On Windows, `os.Rename` is NOT atomic — it is internally a `MoveFileEx` which on the same volume can be atomic, but is NOT if source and target are on different volumes. Additionally, Windows does not allow renaming over an existing file in all scenarios — the rename may fail if the target is open by another process (e.g., a backup agent).

**Why it happens:** Developers test on their dev machine where `/tmp` and home are on the same filesystem. The bug surfaces only in production (encrypted home directories, network drives, USB sticks).

**How to avoid:**
- Always write `.abditum.tmp` in **the same directory** as the target vault file. Use `filepath.Dir(vaultPath)` to construct the temp path, not `os.TempDir()`.
- On Windows, use `windows.MoveFileEx` with `MOVEFILE_REPLACE_EXISTING | MOVEFILE_WRITE_THROUGH` for near-atomic behavior. Implement a Windows-specific code path if `os.Rename` is insufficient.
- Test the atomic save on: a FAT32-formatted USB stick (no hard links, rename may fail); a network drive; an encrypted home directory (eCryptfs / VeraCrypt container).
- Add an integration test that simulates a crash mid-write (truncate the `.tmp` file) and verifies the original vault is intact and the `.tmp` is cleaned up on next open.

**Warning signs:**
- `os.TempDir()` or `ioutil.TempFile("", ...)` used for the `.tmp` file.
- `filepath.Dir(vaultPath) != filepath.Dir(tmpPath)`.
- No integration test for interrupted save.
- No Windows-specific rename code path.

**Phase to address:** Storage package. Validate atomicity guarantees across platforms as part of the storage package test suite.

---

### Pitfall 9: File Watch Race Condition in External Modification Detection

**What goes wrong:** The architecture uses timestamp + size comparison at save time to detect external modification. Two races:
1. **TOCTOU:** Between the `os.Stat` call (to read timestamp/size) and the write of `.abditum.tmp`, another process modifies the file. The stat was taken before the write, so the comparison is against a stale baseline.
2. **Filesystem timestamp resolution:** FAT32 has 2-second timestamp resolution. A fast external modification within the same 2-second window goes undetected. Even ext4 on Linux truncates timestamps to 1-nanosecond precision, but some network filesystems (NFS, SMB) have 1-second resolution — an external modification in the same second is invisible.
3. **Size collision:** A different file with the same size as the original satisfies the check. This is rare but possible for sequential saves that happen to produce the same byte count.

**Why it happens:** The stat-based approach is simple and portable, but is fundamentally a heuristic, not a guarantee.

**How to avoid:**
- Record BOTH timestamp AND size AND a hash (e.g., BLAKE2b or SHA-256) of the file at open/save time. Compare all three before overwriting. The hash eliminates size collisions and reduces the effective TOCTOU window.
- Document that the detection is best-effort. The user confirmation dialog on conflict is the safety net.
- On supported platforms, use `inotify` (Linux) or `FSEvents` (macOS) or `ReadDirectoryChangesW` (Windows) for real-time detection, as opt-in enhanced monitoring. This is a v2 improvement.
- Test: open vault, modify it externally in the same second, attempt save — verify the conflict dialog appears.

**Warning signs:**
- Only `ModTime()` compared, no size or hash.
- The baseline stat is taken at open time and never refreshed after a successful save.
- No test for the sub-second modification case.

**Phase to address:** Storage package. The hash-based detection should be part of the initial implementation, not added later.

---

### Pitfall 10: mlock/VirtualLock Availability Assumptions

**What goes wrong:** The code calls `unix.Mlock(buf)` without checking the return value, or assumes the call succeeded, or calls it after the slice has already been accessed (too late — the data may already have been paged out). A subtler variant: `mlock` succeeds but the slice is reallocated by a later `append`, creating a non-locked copy.

On Linux, unprivileged processes have an `RLIMIT_MEMLOCK` quota (commonly 64KB on older systems, 8MB on newer). Trying to lock more memory than the quota silently fails or returns `EPERM`. On Windows, `VirtualLock` has similar per-process quotas managed via user rights policy.

On macOS (Apple Silicon), `mlock` is available but enforced differently — it may succeed but the OS reserves the right to page encrypted swap.

**Why it happens:** Developers test as root or on systems with high limits. The CI runner has different limits.

**How to avoid:**
- Always check the error return of `mlock`/`VirtualLock`. Log (generically, no sensitive content) that memory locking failed and continue — the application must operate correctly without it.
- Allocate the master password buffer and AES key buffer **before** any other allocation, ideally at program start, to minimize fragmentation that could force the OS to use non-lockable pages.
- Use `mlock` immediately after allocation, before writing sensitive data into the buffer.
- Do NOT `append` to a locked buffer — allocate the exact size needed upfront.
- Add a startup diagnostic (debug build only) that reports whether `mlock` succeeded.
- Build-tag the mlock code: `//go:build linux || darwin` and `//go:build windows`.

**Warning signs:**
- `unix.Mlock(buf)` with no error check.
- Sensitive buffer created as `make([]byte, 0, 256)` and then grown with `append`.
- mlock called after the first read of the buffer.
- No build tag separating Unix and Windows mlock implementations.

**Phase to address:** Crypto package foundation, memory safety layer. Implement mlock wrappers before any sensitive data handling.

---

### Pitfall 11: Error Messages Leaking Sensitive Information

**What goes wrong:** Go's standard library error types (`*os.PathError`, `*fs.PathError`, `cipher.ErrAuthFailed`) carry contextual information — including file paths and operation names — that may be passed directly to the user via `err.Error()`. A corrupted vault returning `"open /home/alice/.secrets/vault.abditum: permission denied"` leaks the file path (violates `SEC-PRIV-01`). A decryption error returning the raw `crypto/cipher` error could help an attacker distinguish between a wrong password and a corrupted file.

**Why it happens:** `fmt.Errorf("failed to open vault: %w", err)` wraps and exposes the underlying OS error. Developers use `err.Error()` in the UI layer without sanitizing.

**How to avoid:**
- Define a closed set of domain errors in `internal/vault`: `ErrAuthFailed`, `ErrCorrupted`, `ErrIOFailure` — opaque sentinel values with no embedded path or OS details.
- In `internal/storage` and `internal/crypto`, catch all OS/stdlib errors and translate them into domain errors before returning to the vault Manager.
- The TUI layer only ever receives domain errors and maps them to user-facing strings that contain no paths, field names, or technical details.
- Write a test that triggers each error path and asserts that neither the error value nor the display string contains the vault file path.

**Warning signs:**
- `%w` wrapping of `*os.PathError` returned from storage package to TUI.
- `err.Error()` called in TUI `Update` and included in a `model.ErrorMessage` field.
- No domain error type hierarchy in `internal/vault`.

**Phase to address:** Storage and crypto packages, alongside vault Manager. Error taxonomy must be defined before TUI error handling is written.

---

### Pitfall 12: JSON Marshaling of Sensitive Fields Creates Uncontrolled Copies

**What goes wrong:** The vault's in-memory representation (Go structs) is serialized to JSON before encryption. `encoding/json` (and `encoding/json/v2`) internally uses `reflect` and builds an intermediate `[]byte` buffer that is heap-allocated outside your control. After marshaling, you hold the ciphertext, but the plaintext JSON bytes that were marshaled sit as unreachable heap objects — you cannot zero them.

**Why it happens:** JSON marshaling is the obvious serialization choice. The hidden cost is that the marshaler may create multiple intermediate string and byte slice allocations that are impossible to trace or zero.

**How to avoid:**
- Accept this as an irreducible limitation of using `encoding/json` with Go's GC. Document it explicitly in comments.
- Minimize the window: marshal → encrypt → zero the plaintext JSON buffer immediately after `gcm.Seal`. Although you can't guarantee zeroing all copies, you reduce the window drastically.
- Prefer `[]byte` for the marshal output (use `json.Marshal` which returns `[]byte`, not `json.Encoder` with a `strings.Builder`).
- Consider using a custom serialization format (MessagePack, protocol buffers, or a custom binary format) that can be written directly into a pre-allocated `[]byte` without intermediate allocations, for a future security hardening phase.
- Zero the ciphertext plaintext buffer (`plaintext = nil; runtime.GC()` is NOT sufficient, but zeroing the slice you own is better than nothing).

**Warning signs:**
- `json.NewEncoder(strings.Builder{})` or similar string-based marshal path.
- The marshal buffer kept alive beyond the encrypt call.
- No comment documenting the limitation.

**Phase to address:** Storage package serialization design. Acknowledge the limitation upfront and design around minimizing exposure window.

---

### Pitfall 13: Search Indexing Sensitive Field Values

**What goes wrong:** The search implementation iterates over all fields and their values. A developer accidentally includes sensitive fields in the search index (substring match against `field.Value` for all fields, not just common ones). A user searches for a common word that happens to appear in a password — the search returns that secret, revealing that the password contains that word.

Beyond the UX violation, if search results are ranked or cached, sensitive substrings appear in non-secure memory regions outside the normal sensitive buffer lifecycle.

**Why it happens:** The field type (`common` vs `sensitive`) is not checked in the search loop. A late addition of search over the "observation" field accidentally starts searching all fields.

**How to avoid:**
- Define a `Field.IsSensitive bool` (or `FieldType` enum) that is checked at the top of every loop that reads field values. Make it impossible to access the raw value of a sensitive field without explicitly acknowledging the sensitivity.
- Write a unit test: create a secret with a sensitive field containing "searchable", run a search for "searchable", assert zero results.
- Code review checklist: any loop over secret fields must include a `if field.IsSensitive { continue }` guard for value access.

**Warning signs:**
- Search loop iterates `secret.Fields` and accesses `field.Value` without type check.
- `QUERY-02` acceptance test is not written as a negative test (sensitive fields must NOT appear).

**Phase to address:** Vault Manager search implementation, from the first commit.

---

### Pitfall 14: Auto-Lock Timer Not Reset by All Input Events

**What goes wrong:** The inactivity timer is reset only on specific Bubble Tea messages (e.g., `tea.KeyMsg`) but not on all input events (mouse events, window resize events, paste events on Windows). A user who is actively scrolling with the mouse but not pressing keys gets locked out unexpectedly. Alternatively, the timer fires while a write operation is in progress, causing a lock that corrupts the state machine.

**Why it happens:** Bubble Tea delivers multiple message types; resetting only on `tea.KeyMsg` is the obvious first implementation.

**How to avoid:**
- Reset the inactivity timer on any message that represents user interaction: `tea.KeyMsg`, `tea.MouseMsg`, potentially `tea.WindowSizeMsg` (debatable — resize may be automatic).
- Ensure that timer firing during an async operation (ongoing save) defers the lock until the operation completes: use a `lockPending bool` flag that is checked after the save goroutine finishes.
- Write a test: simulate 30 messages with only mouse events; assert the timer resets; assert no lock fires.
- Verify that the timer goroutine uses `time.NewTimer` (resettable) rather than `time.NewTicker` (non-resettable without recreating).

**Warning signs:**
- `time.AfterFunc` used for the auto-lock timer (not easily cancellable/resettable).
- Timer reset only in the `tea.KeyMsg` branch of `Update`.
- Goroutine leak: old timer goroutine not cancelled when timer is reset.

**Phase to address:** Vault lifecycle / auto-lock implementation phase.

---

### Pitfall 15: Goroutine Leaks from Timer-Based Clipboard and Field Reveal

**What goes wrong:** The clipboard auto-clear and field-reveal timers launch goroutines via `time.AfterFunc`. When the user locks or exits before the timer fires, those goroutines are still running and may:
1. Try to clear the clipboard after the application state has been torn down.
2. Send a Bubble Tea command on a channel that is no longer being read, causing a goroutine leak or panic.
3. Race with the main lock/exit cleanup.

**Why it happens:** `time.AfterFunc` is fire-and-forget. The goroutine holds a reference to the application state, preventing GC.

**How to avoid:**
- Use `time.NewTimer` with a dedicated `chan struct{}` cancel channel for each active timer. On lock/exit, close the cancel channel to stop all pending timers before proceeding with cleanup.
- Alternatively, use Bubble Tea's own command mechanism: return a `tea.Tick` command that sends a message back through the normal `Update` loop, so the timer is inherently cancelled when the program exits (the event loop stops).
- Write a test: start a clipboard timer with a 10-second delay, lock the vault after 100ms, assert that the goroutine count does not increase (use `runtime.NumGoroutine`).

**Warning signs:**
- `time.AfterFunc` with no corresponding cancel/stop path.
- Goroutine references the vault Manager model directly (not via a message channel).
- No lock/exit cleanup that cancels all outstanding timers.

**Phase to address:** Clipboard and field-reveal feature implementation.

---

### Pitfall 16: Backup Protocol State Machine Bugs (`.bak` / `.bak2` Restoration)

**What goes wrong:** The backup protocol (`.tmp` → rename → `.bak` rotation with `.bak2` intermediary) has several failure modes:
1. Between renaming `.bak` to `.bak2` and writing the new `.bak`, the process crashes. Result: `.bak2` exists, `.bak` does not — on next open, the backup rotation is wrong.
2. On Windows, renaming `.abditum` to `.abditum.bak` fails because a virus scanner or backup agent has the file open — the entire save operation fails with no clear error.
3. The cleanup of `.bak2` is forgotten after a successful save due to an early `return err` in an error path, leaving orphaned `.bak2` files.
4. On a filesystem with no rename atomicity (rarely, e.g., certain FUSE mounts), the old `.bak` is deleted before the new one is written, eliminating the backup.

**Why it happens:** Multi-step file operations with multiple failure points are inherently stateful. Each step has an error path, and the error handling matrix grows combinatorially.

**How to avoid:**
- Implement the backup protocol as an explicit state machine with named steps. Write a recovery function that runs at startup and repairs any inconsistent state (`.bak2` without `.bak` → rename `.bak2` to `.bak`).
- Write exhaustive integration tests for each failure point: fail after `.bak2` creation, fail during rename, fail after rename before cleanup.
- On Windows, retry the rename with exponential backoff (up to 500ms total) before failing — file lock contention from AV scanners is transient.
- Log (generically) any unexpected `.bak2` file found at startup and attempt automatic recovery.

**Warning signs:**
- No startup recovery scan for orphaned `.bak2` files.
- The backup protocol implementation is a single linear function without explicit failure recovery for each step.
- No integration test that simulates process death mid-save.

**Phase to address:** Storage package. Include state machine recovery in the initial storage implementation, not as a later addition.

---

### Pitfall 17: Sensitive Data in Golden File Snapshots

**What goes wrong:** The test suite uses `teatest/v2` golden files (`.golden` snapshot files of terminal output) to detect visual regressions. If a test opens a real or mock vault and navigates to a screen showing a revealed sensitive field, the snapshot captures the plain-text value. These snapshots are committed to version control and visible to anyone with repository access.

**Why it happens:** Golden file tests capture the full terminal buffer. If the test data includes realistic-looking (but fake) credentials, those appear in the committed snapshot.

**How to avoid:**
- Use semantically valid but obviously fake credentials in all test fixtures: `password: "TESTPASSWORD_NOT_REAL"`, API key: `"sk_test_0000000000000000"`. Never use real-looking passwords.
- Add `.golden` files to a review checklist: verify no sensitive-looking patterns before committing.
- Add a CI step that scans golden files for patterns matching common secret formats (regex for `[a-zA-Z0-9]{32,}` long tokens, PEM headers, etc.) and fails the build if found.
- Never use real vault files as test fixtures.

**Warning signs:**
- `.golden` files in version control without a review step.
- Test fixtures using `password123` or similar realistic-looking credentials.
- No CI scan for secrets in test files.

**Phase to address:** Test infrastructure setup, before writing any TUI tests.

---

### Pitfall 18: Cross-Platform Binary Build Inconsistency (CGO)

**What goes wrong:** `CGO_ENABLED=0` is required for a static binary. However, some clipboard or mlock libraries have optional CGO paths that are silently used when CGO is enabled locally (on a developer's machine) but disabled in CI. The result: the binary works locally but the CI artifact (the actual distributable) silently skips clipboard or mlock functionality without any error at runtime.

A related variant: `syscall.Mlock` is available on Linux/macOS but the Windows equivalent requires importing `golang.org/x/sys/windows` — if the build tag is missing, the Windows binary silently skips mlock without the developer realizing it.

**Why it happens:** Developers build with CGO enabled for development convenience. The mismatch is invisible until someone tests the release binary.

**How to avoid:**
- CI builds all platform targets with `CGO_ENABLED=0 GOOS=windows/linux/darwin GOARCH=amd64/arm64` explicitly set.
- Use build tags consistently (`//go:build linux || darwin` and `//go:build windows`) for platform-specific code (mlock, VirtualLock, clipboard).
- Write a build test that verifies no CGO dependencies: `go build -v 2>&1 | grep -i cgo` should be empty.
- Provide `no-op` stub implementations for mlock/clipboard for GOOS targets that don't support them, so the function always exists and can always be called safely.

**Warning signs:**
- No explicit `CGO_ENABLED=0` in CI build commands.
- mlock implementation in a single file without build tags.
- No cross-compilation CI job (builds only for the host platform).

**Phase to address:** CI/CD setup and build pipeline, first sprint.

---

## Phase-Specific Warning Map

| Phase / Feature | Likely Pitfall | Priority |
|-----------------|----------------|----------|
| Crypto package | Pitfall 1 (string vs []byte), Pitfall 2 (GC copies), Pitfall 3 (nonce reuse), Pitfall 4 (Argon2id params), Pitfall 10 (mlock) | **CRITICAL** |
| File format design | Pitfall 4 (params not stored), Pitfall 12 (JSON copies) | **CRITICAL** |
| Storage package | Pitfall 7 (atomic rename cross-device), Pitfall 8 (stat race), Pitfall 16 (bak state machine) | **HIGH** |
| Vault Manager | Pitfall 11 (error message leakage), Pitfall 13 (search sensitive fields) | **HIGH** |
| TUI architecture | Pitfall 6 (state machine unbounded), Pitfall 14 (auto-lock timer) | **HIGH** |
| Lock/exit lifecycle | Pitfall 5 (clipboard cross-platform), Pitfall 6 (terminal scrollback), Pitfall 15 (goroutine leaks) | **HIGH** |
| Test infrastructure | Pitfall 17 (golden file secrets) | **MEDIUM** |
| Build/CI | Pitfall 18 (CGO mismatch) | **MEDIUM** |

---

## Sources

- OWASP Cryptographic Storage Cheat Sheet (2025) — Argon2id parameters
- Go security advisories and `golang.org/x/crypto` changelogs
- Dave Cheney, "Practical cryptography with Go" — nonce reuse analysis
- Bubble Tea v2 official examples and `teatest/v2` documentation
- POSIX `rename(2)` man page — same-filesystem atomicity guarantee
- Windows `MoveFileEx` documentation — cross-volume limitations
- `golang.org/x/sys/unix` mlock documentation
- Analysis of open-source Go password managers: `gopass`, `passage`, `kpcli-go` post-mortems on GitHub Issues
- Go memory model specification (2022 revision) — GC copying behavior
- xterm, Windows Terminal, and GNOME Terminal documentation — scrollback buffer clear escape codes
