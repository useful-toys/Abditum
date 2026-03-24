# UX Benchmarks & Data Model Research

**Project:** Abditum  
**Researched:** 2026-03-24  
**Scope:** Password manager UX patterns and data model decisions from Bitwarden, KeePassXC, gopass, and pass; inactivity detection patterns from Bubble Tea  

---

## 1. Data Model Structures

### KeePassXC

Source: `keepassxreboot/keepassxc` — `Entry.h`, `Group.h`, `EntryAttributes.h`, `Database.h`

- **Recursive Group tree:** `Group` contains `[]Group` and `[]Entry` — arbitrary nesting depth
- **Fixed "default attributes" on Entry:** Title, UserName, Password, URL, Notes — hardcoded in `EntryAttributes`; these keys receive special treatment (e.g. Password is always protected)
- **Custom attributes:** Stored in `EntryAttributes` as `QMap<QString, QString>` with a separate `QSet<QString>` of protected (hidden) keys — protection is a per-attribute binary flag, not a per-type concept
- **Any attribute can be marked protected** — protection means: rendered as `****` in UI by default, encrypted in the inner KDBX stream separately from the main payload
- **Recycle Bin:** A special Group whose UUID is tracked in Database metadata; `recycleEntry(entry)` and `recycleGroup(group)` re-parent the item into the Recycle Bin group; `emptyRecycleBin()` permanently deletes its contents
- **`previousParentGroupUuid`** stored on both Entry and Group — enables "Restore" (re-parent back to origin)
- **`deletedObjects`:** A list of `{UUID, deletionTime}` for items that were hard-deleted (bypassing Recycle Bin); used for KDBX sync/merge, not exposed as a UX concept
- **Entry history:** `m_history: QList<Entry*>` — full list of previous entry snapshots, oldest first; each history item is itself a full Entry object

### Bitwarden

Source: `bitwarden/clients` — `cipher.ts`, `field-type.enum.ts`

- **Cipher types (enum):** `Login=1`, `SecureNote=2`, `Card=3`, `Identity=4`, `SshKey=5` — 5 hard-coded item types
- **Per-type sub-object:** Each type has a dedicated strongly-typed sub-object (`login`, `card`, `identity`, `secureNote`, `sshKey`) with named fields specific to that type
- **Generic custom fields overlay:** ALL cipher types also have a `fields: Field[]` array for user-defined key-value pairs
- **`FieldType` enum:** `Text=0`, `Hidden=1`, `Boolean=2`, `Linked=3`
  - `Text` — non-sensitive plain text
  - `Hidden` — sensitive/masked; rendered as `****` by default, requires explicit reveal action
  - `Boolean` — checkbox (yes/no toggle); useful for flags
  - `Linked` — references another field value within the same cipher (autofill use case; not relevant to offline vaults)
- **Flat folder system:** Items belong to one `folderId`; no folder nesting; "collections" are an org-level concept
- **`deletedDate?: Date`** on Cipher — soft delete is a nullable timestamp on the item itself; items with `deletedDate != null` are in Trash and filtered from normal views; server auto-purges after 30 days
- **`reprompt: CipherRepromptType`** — optional master-password re-prompt before revealing a cipher's sensitive fields (extra security layer)
- **`favorite: boolean`** on Cipher
- **`passwordHistory?: Password[]`** — history of previous password values per cipher

### gopass

Source: `gopasspw/gopass` — `docs/features.md`

- **Filesystem-based:** Each secret = one GPG-encrypted file; directories = groups/folders
- **No strict schema:** Convention only — first line = password, rest = YAML key-value pairs; separator `---` for structured sections
- **TOTP support via convention:** `otpauth://` URL or `totp: BASE32SECRET` key in YAML — no dedicated field type
- **`safecontent` mode:** Hides the password (first line) by default when showing a secret; `unsafe-keys` config can hide arbitrary keys
- **Search:** Literal match first; automatic fuzzy fallback if no exact match found; no GUI filtering
- **No templates concept:** Completely free-form text files; field names are YAML keys

### pass (unix password-store)

- One GPG-encrypted file per secret; directories = folders
- Convention only: password on first line, arbitrary text below; no structured fields, no types, no GUI
- No built-in search — users `grep` directly or use `pass grep` extension

---

## 2. Field Type Systems

| Tool | Type System | Sensitivity Model |
|------|-------------|-------------------|
| KeePassXC | No types — all attributes are strings | Per-attribute binary `isProtected` flag |
| Bitwarden | 4 types: `Text`, `Hidden`, `Boolean`, `Linked` | Sensitivity is encoded in the TYPE (`Hidden` = sensitive) |
| gopass | No types — YAML values are all strings | Convention-based (`safecontent` hides first line; `unsafe-keys` for others) |
| pass | No types | None |

**Key takeaway for Abditum:**

Abditum's current `enum(texto, texto_sensível)` maps precisely to Bitwarden's `Text`/`Hidden` — this is validated by two independent tools using the same binary sensitive/non-sensitive distinction. The `Boolean` (checkbox) type from Bitwarden is worth considering for v2 (e.g. "Active: yes/no", "2FA enabled: yes/no"), but is not essential for v1.

**`Linked` type is irrelevant** to Abditum (it's an autofill-specific concept with no use in an offline TUI vault).

---

## 3. Search UX

| Tool | Search Trigger | Scope | Highlighting |
|------|---------------|-------|-------------|
| KeePassXC GUI | Real-time as-you-type | Title, Username, URL, Notes, custom attributes (configurable) | Yes |
| Bitwarden CLI | Explicit `--search` flag | All fields | No (JSON output) |
| Bitwarden Web/Desktop | Real-time as-you-type | All fields | Yes |
| gopass | Explicit CLI arg | Secret name (path); `grep` for content | No |
| pass | No built-in | N/A | N/A |

**Pattern for TUI:** Real-time filtering is the dominant pattern for GUI-class password managers. The list/tree filters to matching items as the user types, with match highlighting. Abditum's current spec already aligns with this pattern.

**Scope recommendation:** Search name, non-sensitive field values, notes, and containing folder name — exactly as Abditum's spec states. Deliberately exclude sensitive (`texto_sensível`) field values from search to avoid inadvertently displaying them in search result context snippets.

---

## 4. Import/Export Formats

### KeePass XML (KDBX XML export)

- De facto interchange standard; supported by virtually every major password manager as both import and export
- Structure:
  ```xml
  <KeePassFile>
    <Root>
      <Group>
        <Name>Sites</Name>
        <Entry>
          <String><Key>Title</Key><Value>GitHub</Value></String>
          <String><Key>Password</Key><Value ProtectInMemory="True">secret</Value></String>
          <String><Key>CustomField</Key><Value ProtectInMemory="True">value</Value></String>
        </Entry>
      </Group>
    </Root>
  </KeePassFile>
  ```
- `ProtectInMemory="True"` attribute maps cleanly to Abditum's `texto_sensível` field type
- Preserves arbitrary folder hierarchy naturally via nested `<Group>` elements

### Bitwarden JSON export

- Used for Bitwarden-to-Bitwarden or Bitwarden-to-other migrations
- Structure: `{"encrypted": false, "folders": [...], "items": [...]}` where items include `type` (int), `name`, `notes`, `fields[]`, and type-specific sub-objects
- Flat folder model (no nesting) is a significant structural mismatch with Abditum's hierarchy
- Less universally supported outside the Bitwarden ecosystem

### Recommendation for Abditum

Use a **custom Abditum JSON format** for native export/import (already decided in `PROJECT.md`). For future interoperability (v2), **KeePass XML is the better interop target** because:
1. Wider tool support than Bitwarden JSON
2. XML structure naturally represents nested folder hierarchy
3. `ProtectInMemory` attribute maps cleanly to `texto_sensível`
4. Bitwarden JSON's flat folder model would lose hierarchy information

---

## 5. Auto-lock / Inactivity Detection in Bubble Tea

### How Bubble Tea Receives Input

Every keypress and mouse event flows through `Update(msg tea.Msg)` as a `tea.KeyMsg` or `tea.MouseMsg`. There is no OS-level idle API accessible from a terminal application across Windows/macOS/Linux.

### Inactivity Timer Pattern

The only viable cross-platform pattern for TUI inactivity detection:

```go
// State in model
type model struct {
    lastActivity time.Time
    lockTimeout  time.Duration
    // ...
}

// Tick periodically (e.g. every 10 seconds)
func tickCmd() tea.Cmd {
    return tea.Tick(10*time.Second, func(t time.Time) tea.Msg {
        return inactivityTickMsg{t}
    })
}

// In Update:
case tea.KeyMsg, tea.MouseMsg:
    m.lastActivity = time.Now()  // reset on every input event
    return m, nil

case inactivityTickMsg:
    if time.Since(m.lastActivity) >= m.lockTimeout {
        return m, lockCmd()
    }
    if time.Since(m.lastActivity) >= m.lockTimeout-30*time.Second {
        // show "locking in 30s" warning
    }
    return m, tickCmd()  // re-arm
```

### Bubble Tea `timer` Bubble (charmbracelet/bubbles)

- Available at `charm.land/bubbles/v2/timer`
- Provides `timer.Model{Timeout, Interval}` — a countdown timer component
- Sends `timer.TickMsg` periodically and `timer.TimeoutMsg` when it expires
- Can be started/stopped via `timer.StartStopMsg`
- **Simpler alternative to manual tick management** — but requires resetting on every input event
- `timer.Model.Reset()` resets the countdown; call this on every `tea.KeyMsg`/`tea.MouseMsg`

### Warning-Before-Lock Pattern (Abditum spec)

Abditum's spec calls for a "lock imminent" warning before locking. Implementation:

1. Primary timer: fires at `timeout - warningDuration` → show warning overlay
2. Secondary/remaining timer: fires at `warningDuration` → lock
3. Any key event resets both timers and dismisses the warning

---

## 6. Secret Template / Item Type Concept

| Tool | Approach | User-Definable? | Snapshot at Creation? |
|------|----------|-----------------|----------------------|
| Bitwarden | Hard-coded enum (5 types); each has a typed sub-object | No | N/A (type persists) |
| KeePassXC | No templates — all entries have the same 5 default fields + free custom attributes | No (no templates) | N/A |
| gopass | No templates — YAML keys are the "fields" | N/A | N/A |
| **Abditum** | **User-defined templates stored in vault; template = field name + type pairs** | **Yes** | **Yes — snapshot at creation time** |

**Assessment:** Abditum's approach is uniquely well-suited:
- More flexible than Bitwarden (user-defined vs. hard-coded)
- More structured than KeePassXC (named typed fields vs. free-form attributes)
- Snapshot semantics prevent retroactive mutations of existing secrets when templates evolve — this is the correct decision

**Pre-defined templates (from `PROJECT.md`):**
- `Login`: URL (`texto`), Username (`texto`), Password (`texto_sensível`)
- `Cartão de Crédito`: Número do Cartão (`texto_sensível`), Nome no Cartão (`texto`), Data de Validade (`texto`), CVV (`texto_sensível`)
- `API Key`: Nome da API (`texto`), Chave de API (`texto_sensível`)

These cover the primary use cases without overcomplicating the initial template set.

---

## 7. Trash / Soft Delete Patterns

### KeePassXC

- Recycle Bin = a special `Group` with a UUID tracked in `Database` metadata
- `recycleEntry(entry)` / `recycleGroup(group)` — re-parent item into Recycle Bin group
- `emptyRecycleBin()` — permanently deletes all contents of the Recycle Bin group
- Restore = re-parent item from Recycle Bin back to `previousParentGroupUuid`
- Recycle Bin can be disabled in settings
- Trash items persist across sessions until explicitly emptied

### Bitwarden

- Soft delete = set `deletedDate: Date` timestamp on the Cipher; no folder/group move
- Items with `deletedDate != null` are filtered from normal views, shown only in Trash view
- Restore = clear `deletedDate`
- Server auto-purges Trash items after 30 days
- Hard delete = permanent immediate removal

### Abditum (current design — from `PROJECT.md`)

- Trash = virtual folder showing items flagged `deleted=true` (or equivalent)
- **Items in Trash are purged permanently on save** — no "empty trash" action needed; the save is the commit point
- Restore before saving = unflag the item
- Once saved, deleted items are permanently gone

**Assessment:** Abditum's "purge on save" is a valid and simpler third approach, well-suited to an offline vault where:
- There is no server-side retention period
- The save action is already the user's explicit "commit" intent
- "Undo delete" = restore from Trash before saving (within-session undo)
- The `.abditum.bak` backup provides recovery from accidental saves

**Open question:** Should Trash items survive across sessions (i.e., be persisted in the vault file between saves)? The current design implies purge-on-save, meaning items deleted in a previous session cannot be recovered from Trash. The `.bak` backup is the recovery mechanism for that scenario. This is consistent and intentional.

---

## 8. Cross-Cutting Findings & Validation of Abditum Decisions

| Decision | Industry Evidence | Verdict |
|----------|------------------|---------|
| `texto` / `texto_sensível` field types | Matches Bitwarden `Text`/`Hidden`; matches KeePassXC protected flag | ✅ Validated |
| User-defined snapshot templates | More flexible than any surveyed tool; snapshot prevents mutation bugs | ✅ Correct approach |
| Real-time search filtering | KeePassXC GUI, Bitwarden GUI both use this; dominant GUI pattern | ✅ Validated |
| Exclude sensitive fields from search | Best practice — avoids displaying secrets in search result snippets | ✅ Validated |
| Trash purge on save | Simpler than KeePass (explicit empty) and Bitwarden (server timeout); suits offline vault | ✅ Appropriate |
| Inactivity via tick + reset pattern | Only viable cross-platform TUI approach; Bubble Tea architecture supports it | ✅ Confirmed |
| KeePass XML for future interop | Widest tool support; preserves hierarchy; `ProtectInMemory` maps cleanly | ✅ Recommended |
| `Boolean` field type | Bitwarden has it; could be useful (flags, toggles) | ⏳ Defer to v2 |
| TOTP field type | gopass supports via convention; Bitwarden has no dedicated type either | ⏳ Defer to v2 (already in PROJECT.md) |

---

## Sources

| Source | What Was Examined | Confidence |
|--------|------------------|------------|
| `github.com/keepassxreboot/keepassxc` — `src/core/Entry.h`, `Group.h`, `EntryAttributes.h`, `Database.h` | Data model, protected fields, Recycle Bin, history | HIGH |
| `github.com/bitwarden/clients` — `libs/common/src/vault/domain/cipher.ts`, `field-type.enum.ts` | Cipher type enum, FieldType enum, deletedDate, reprompt | HIGH |
| `github.com/gopasspw/gopass` — `docs/features.md` | Data organization, YAML structure, search, OTP | HIGH |
| `github.com/charmbracelet/bubbles` — `timer/timer.go`, examples | Bubble Tea timer component API | HIGH |
| `github.com/charmbracelet/bubbletea` — `examples/timer/main.go`, `examples/realtime/main.go` | Tick patterns, message flow | HIGH |
