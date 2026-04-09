# Roadmap: Abditum

**Milestone:** v1.0
**Phases:** 10
**Requirements mapped:** 76 / 76

## Phases

| # | Phase | Goal | Requirements | Success Criteria |
|---|-------|------|--------------|------------------|
| 1 | Foundation & Distribution | Project compiles; cross-platform binaries produced; zero external state | DIST-01, DIST-02 | 3 |
| 2 | TUI Design System | Shell renders with all design-system components; theming; message bar; dialogs | TUI-01–10 | 5 |
| 3 | Cryptography & File Format | `.abditum` files can be created, encrypted, decrypted, backed up, and validated | FILE-01–09 | 5 |
| 4 | Vault Lifecycle | Users can create, open, save, configure, and exit the vault | VAULT-01–08, VAULT-12, VAULT-15–16, SEC-01 | 5 |
| 5 | Session Security & Locking | Vault locks safely; memory zeroed; clipboard and screen cleared | VAULT-09–11, SEC-02–06 | 4 |
| 6 | Vault Navigation & Search | Users can navigate the vault tree, view favorites, see session indicators, and search | NAV-01–06 | 4 |
| 7 | Secret Management | Full secret CRUD; sensitive field reveal; clipboard copy with auto-clear | SECRET-01–11, FIELD-01–02 | 5 |
| 8 | Folder Management | Full folder CRUD; content promotion on delete; cycle and duplicate prevention | FOLDER-01–07 | 4 |
| 9 | Template Management | Users can create, edit, and manage secret templates | TMPL-01–07 | 4 |
| 10 | Export & Import | Users can transfer vault data via interchange files | VAULT-13–14 | 3 |

## Phase Details

### Phase 1: Foundation & Distribution

**Goal:** The project compiles and produces a single cross-platform binary; no external state files are created outside the `.abditum` vault file.
**Depends on:** —

**Requirements:**
- **DIST-01**: Aplicação distribuída como binário único executável cross-platform (Windows, macOS, Linux 64-bit); sem instalação; sem dependências externas em runtime
- **DIST-02**: Zero arquivos de estado, config ou dados fora do arquivo do cofre (exceto artefatos transitórios `.abditum.tmp`, `.abditum.bak`, `.abditum.bak2`)

**Success Criteria:**
1. Running `go build` (or the project build script) produces a single self-contained executable for Windows, macOS, and Linux 64-bit without any runtime dependencies
2. Copying the binary to a machine with only a terminal and running it produces no segfault, no missing-library errors, and no extra files in the working directory
3. No configuration files, registry entries, or application-state directories are created anywhere on the system when the binary exits

**Plans:** TBD
**UI hint**: no

---

### Phase 2: TUI Design System

**Goal:** The terminal renders the full design-system shell — 4-zone layout, split panels, themed colours, message bar, dialogs, scrollbars, keymaps, and NO_COLOR compliance — before any vault logic is wired in.
**Depends on:** Phase 1

**Requirements:**
- **TUI-01**: Interface TUI em Go com Bubbletea/Lipgloss baseada em design system: 4 zonas verticais (Cabeçalho 2L | Área de trabalho | Barra de mensagens 1L | Barra de comandos 1L)
- **TUI-02**: Painel esquerdo (≈35%) e painel direito (≈65%); separador `│`; conector `<╡` no item selecionado
- **TUI-03**: Terminal mínimo 80×24; degradação sem crash abaixo do mínimo; truncamento com `…`
- **TUI-04**: Todos os símbolos do inventário BMP-only (U+0000–U+FFFF); sem emojis, sem Nerd Fonts
- **TUI-05**: NO_COLOR compliant: todo estado crítico usa ≥2 camadas de comunicação (cor + símbolo ou atributo tipográfico)
- **TUI-06**: Temas visuais discretos e reativos; identificador do tema gravado no payload criptografado do cofre
- **TUI-07**: Mapa de teclas conforme design system: `Enter`, `Esc`, `Tab`/`⇧Tab`; atalhos globais `F1`, `F12`, `⌃Q`, `⌃!⇧Q`
- **TUI-08**: Barra de mensagens com 7 tipos: Sucesso ✓, Informação ℹ, Alerta ⚠, Erro ✕ (bold), Ocupado spinner `◐◓◑◒`, Dica de campo `•` (italic), Dica de uso `•` (italic)
- **TUI-09**: Diálogos e modais centralizados com bordas arredondadas `╭╮╰╯`; ações na borda inferior; severidade governa visual
- **TUI-10**: Scrollbar em diálogos com conteúdo que excede viewport: setas `↑↓` e thumb `■` na borda direita

**Success Criteria:**
1. An 80×24 terminal displays the 4-zone skeleton: a 2-line header, the split workspace (≈35% / ≈65%, `│` separator, `<╡` connector), a 1-line message bar, and a 1-line command bar
2. Shrinking the terminal below 80×24 truncates content with `…` and does not crash
3. All 7 message-bar types render with correct symbol and text attribute; running `NO_COLOR=1` shows that every critical state remains distinguishable without colour
4. A confirmation dialog appears centred with rounded borders `╭╮╰╯`, its default action highlighted at the bottom edge, and a scrollbar (`↑↓` / `■`) appears when content overflows
5. Pressing `F12` cycles through available themes; each theme changes the colour palette reactively without a restart

**Plans:** TBD
**UI hint**: yes

---

### Phase 3: Cryptography & File Format

**Goal:** The `.abditum` file format engine is complete — files can be created, encrypted, decrypted, validated, atomically written, backed up, and backward-compat migrated.
**Depends on:** Phase 1

**Requirements:**
- **FILE-01**: Arquivo `.abditum` segue formato binário: cabeçalho fixo 49 bytes (`magic ABDT` + versão uint8 + salt 32B + nonce 12B) + payload JSON UTF-8 criptografado + tag GCM 16 bytes
- **FILE-02**: Criptografia: AES-256-GCM com chave derivada via Argon2id (m=256MiB, t=3, p=4); salt e nonce gerados via CSPRNG (`crypto/rand`)
- **FILE-03**: Cabeçalho completo autenticado como AAD do AES-256-GCM; salt substituído apenas ao alterar senha mestra; nonce regenerado a cada salvamento
- **FILE-04**: Gravação atômica: dados escritos em `.abditum.tmp` → rename substituindo arquivo original; em falha, `.abditum.tmp` apagado imediatamente
- **FILE-05**: Backup automático a cada sobrescrita: arquivo anterior preservado como `.abditum.bak` com rotação via `.abditum.bak2`; falha durante escrita exibe mensagem com caminho do backup
- **FILE-06**: Criação de novo cofre em caminho vazio e "salvar como" em caminho vazio gravam diretamente no destino (sem `.abditum.tmp`)
- **FILE-07**: Aplicação detecta modificação externa do arquivo antes de sobrescrever e oferece: Sobrescrever / Salvar como novo arquivo / Cancelar
- **FILE-08**: Aplicação detecta acesso concorrente ao arquivo (outro processo) e falha com mensagem informativa
- **FILE-09**: Compatibilidade retroativa: abre qualquer versão suportada; migra payload em memória; salva no formato atual; versão superior ao suportado rejeitada com erro claro

**Success Criteria:**
1. A newly-created `.abditum` file opened in a hex editor shows the 4-byte `ABDT` magic, a version byte, a 32-byte salt, and a 12-byte nonce in the first 49 bytes, followed by ciphertext
2. Attempting to open the vault with an incorrect password returns an authentication error and allows retry; a bit-flipped file returns an integrity error with no retry option
3. Saving the same vault twice produces two files with different nonces (visible in hex); the `.abditum.bak` file contains the previous version
4. Simulating an external file modification between open and save triggers the overwrite-or-save-as dialog; simulating a concurrent process lock triggers the concurrent-access error
5. Opening a vault file with a lower format version migrates it in memory and saves in the current format; a file with a higher-than-supported version is rejected with a clear version mismatch error

**Plans:** TBD
**UI hint**: no

---

### Phase 4: Vault Lifecycle

**Goal:** Users can create a vault with default structure, open an existing vault, save (and save-as), discard changes, change the master password, configure settings, and exit gracefully.
**Depends on:** Phase 2, Phase 3

**Requirements:**
- **VAULT-01**: Usuário pode criar novo cofre em um arquivo `.abditum` com senha mestra, com confirmação dupla e aviso de senha fraca (não bloqueante)
- **VAULT-02**: Ao criar o cofre, estrutura padrão é criada: Pasta Geral com subpastas "Sites e Apps" e "Financeiro", e modelos padrão (Login, Cartão de Crédito, Chave de API)
- **VAULT-03**: Usuário pode abrir cofre existente a partir de um arquivo `.abditum` com senha mestra
- **VAULT-04**: Abertura valida o arquivo contra corrupção e senha incorreta — erros de autenticação permitem nova tentativa; erros de integridade encerram sem recuperação
- **VAULT-05**: Usuário pode salvar cofre no arquivo atual sem re-informar senha
- **VAULT-06**: Usuário pode salvar cofre em outro arquivo (salvar como), que passa a ser o arquivo de trabalho ativo
- **VAULT-07**: Usuário pode descartar alterações não salvas e recarregar o cofre do arquivo
- **VAULT-08**: Usuário pode alterar a senha mestra — cofre salvo imediatamente com a nova senha; operação irrevogável
- **VAULT-12**: Ao sair com alterações pendentes, usuário é perguntado: Salvar e Sair / Descartar e Sair / Cancelar
- **VAULT-15**: Usuário pode configurar: tempo de bloqueio automático, tempo de ocultação de campo sensível, tempo de limpeza de clipboard, tema visual
- **VAULT-16**: Configurações são persistidas no payload criptografado do cofre (não em arquivo externo)
- **SEC-01**: Senha mestra solicitada uma única vez por sessão; re-uso para todas as operações de criptografia; sem re-solicitação em salvamento ou descarte

**Success Criteria:**
1. Creating a new vault prompts for password with double confirmation, warns (non-blocking) on weak password, and the resulting file contains Pasta Geral with subfolders "Sites e Apps" and "Financeiro" plus the three default templates
2. Opening a vault with the correct password unlocks the workspace; a wrong password shows an error with retry; a corrupted file shows an integrity error and exits without data
3. The master password is asked only once per session; saving, discarding, and save-as all complete without re-prompting for the password
4. Changing the master password immediately saves the file with the new password and the old password no longer opens the vault
5. Exiting with unsaved changes presents the three-option dialog (Save & Exit / Discard & Exit / Cancel); settings changes (theme, timers) persist after re-opening the vault

**Plans:** TBD
**UI hint**: yes

---

### Phase 5: Session Security & Locking

**Goal:** The vault locks safely on demand, on inactivity, and on emergency; sensitive data is zeroed in memory; the terminal and clipboard are cleared on every lock/exit path.
**Depends on:** Phase 4

**Requirements:**
- **VAULT-09**: Usuário pode bloquear manualmente; senha removida da memória, clipboard limpa, terminal limpo (clear screen)
- **VAULT-10**: Cofre bloqueia automaticamente após tempo de inatividade configurável (padrão 5 min); qualquer interação reseta o timer
- **VAULT-11**: Usuário pode acionar bloqueio emergencial com tela disfarçada; atalho `⌃!⇧Q`; descarta alterações; sem confirmação
- **SEC-02**: Senha mestra e buffers sensíveis removidos da memória (zeros) ao bloquear ou encerrar
- **SEC-03**: Terminal limpo (clear screen) ao bloquear ou encerrar, evitando dados visíveis no buffer do terminal
- **SEC-04**: Dados sensíveis em memória protegida (mlock/VirtualLock) quando disponível no SO; ausência não impede operação
- **SEC-05**: Valor da clipboard removido ao bloquear, encerrar ou após tempo configurável (padrão 30s); limpeza depende de suporte do SO
- **SEC-06**: Nenhum log (stdout/stderr) contém caminhos de cofre, nomes de segredos ou valores de campos

**Success Criteria:**
1. After 5 minutes of inactivity (configurable), the vault locks automatically: the terminal is cleared, the password is no longer in memory, and re-entry requires the master password
2. Manual lock (`⌃Q` or menu item) immediately clears the terminal, zeros the master password in memory, and clears any vault-related clipboard value
3. Pressing `⌃!⇧Q` immediately replaces the screen with the disguised view — no vault name, no secret data, no pending-changes dialog — regardless of current state
4. Running the application with stdout/stderr redirected to a log file produces no output containing vault file paths, secret names, or field values at any point during normal operation

**Plans:** TBD
**UI hint**: yes

---

### Phase 6: Vault Navigation & Search

**Goal:** Users can navigate the full vault tree with folder counts and session-state indicators, view the Favorites panel, and search secrets by name, field, or observation.
**Depends on:** Phase 4

**Requirements:**
- **NAV-01**: Interface exibe árvore hierárquica do cofre com pastas e segredos; cada pasta exibe contagem total de segredos ativos (recursivo)
- **NAV-02**: Pasta virtual "Favoritos" exibida como nó irmão da Pasta Geral (acima dela); somente leitura; lista segredos favoritos em profundidade seguindo ordem do cofre
- **NAV-03**: Usuário pode buscar segredos por nome, nome de campo, valor de campo comum e observação (substring, case-insensitive, sem distinção de acentuação); campos sensíveis nunca participam da busca por valor; segredos marcados para exclusão excluídos dos resultados
- **NAV-04**: Segredo exibe nome, todos os campos (com indicação de sensível/comum) e observação
- **NAV-05**: Indicadores de estado de sessão na listagem: ✦ adicionado, ✎ modificado, ✗ excluído (segredos sem alteração não exibem indicador)
- **NAV-06**: Dirty state global do cofre sempre visível no cabeçalho da interface

**Success Criteria:**
1. The left panel renders the full folder tree; each folder label includes the total recursive count of active secrets; expanding/collapsing folders updates counts correctly
2. The Favorites virtual node appears above Pasta Geral, is read-only, and lists every secret with `favorito=true` from the entire vault in depth-first order
3. Searching for a substring returns secrets whose name, field name, common-field value, or observation contains it (case-insensitive, accent-insensitive); sensitive field values never appear in results; secrets marked for deletion are excluded
4. Newly added secrets show ✦, edited secrets show ✎, deletion-marked secrets show ✗ in the left panel; the header shows the dirty-state indicator whenever unsaved changes exist

**Plans:** TBD
**UI hint**: yes

---

### Phase 7: Secret Management

**Goal:** Users can create, view, edit, duplicate, move, reorder, favorite, and soft-delete secrets; sensitive fields can be temporarily revealed; any field value can be copied to the clipboard with auto-clear.
**Depends on:** Phase 6

**Requirements:**
- **SECRET-01**: Usuário pode criar segredo a partir de modelo existente ou sem campos estruturados (apenas com Observação automática)
- **SECRET-02**: Usuário pode escolher a pasta de destino no momento da criação do segredo
- **SECRET-03**: Usuário pode duplicar segredo; duplicata recebe nome único automático (ex: "Gmail (1)") e é posicionada imediatamente após o original; histórico de modelo preservado
- **SECRET-04**: Usuário pode editar nome, valores de campos e observação de um segredo
- **SECRET-05**: Usuário pode editar estrutura de segredo: adicionar campo (nome + tipo), renomear campo, reordenar campos, excluir campo; não é permitido alterar o tipo de um campo existente
- **SECRET-06**: Observação é campo automático em todo segredo: não pode ser renomeada, excluída ou movida; sempre na última posição; campo comum (sempre visível)
- **SECRET-07**: Usuário pode favoritar e desfavoritar segredo
- **SECRET-08**: Usuário pode marcar segredo para exclusão; segredo permanece na lista com indicador ✗ e é excluído permanentemente apenas ao salvar; restaurável antes do salvamento
- **SECRET-09**: Usuário pode desmarcar exclusão de segredo, restaurando o estado anterior
- **SECRET-10**: Usuário pode mover segredo para outra pasta
- **SECRET-11**: Usuário pode reordenar segredo dentro da mesma pasta; ordem final persistida ao salvar
- **FIELD-01**: Usuário pode exibir temporariamente o valor de um campo sensível; valor oculto automaticamente após tempo configurável (padrão 15s)
- **FIELD-02**: Usuário pode copiar valor de qualquer campo para a clipboard; valor removido automaticamente após tempo configurável (padrão 30s) ou ao bloquear/encerrar

**Success Criteria:**
1. Creating a secret from a template pre-populates field names and types; creating a blank secret contains only the Observation field; the destination folder can be changed before saving
2. Duplicating "Gmail" produces "Gmail (1)" inserted directly after the original, with the same field structure and template-origin history
3. A sensitive field shows `••••` by default; the user can temporarily reveal it; after the configured time (default 15 s) the value is hidden again automatically
4. Copying a field to the clipboard shows a countdown in the message bar; after 30 seconds (or on lock/exit) the clipboard is cleared; the Observation field is always present, always last, and cannot be renamed, deleted, or moved
5. Marking a secret for deletion shows it with `✗` in the tree; re-opening the vault after saving confirms the secret is permanently gone; before saving, unmarking restores the secret to its previous state

**Plans:** TBD
**UI hint**: yes

---

### Phase 8: Folder Management

**Goal:** Users can create, rename, move, reorder, and delete folders; Pasta Geral is immovable; content is promoted on folder deletion; cycle and duplicate sibling-name rules are enforced.
**Depends on:** Phase 7

**Requirements:**
- **FOLDER-01**: Usuário pode criar pasta dentro de outra pasta; nome deve ser único entre irmãos na mesma pasta pai
- **FOLDER-02**: Usuário pode renomear pasta; novo nome deve ser único entre irmãos na mesma pasta pai
- **FOLDER-03**: Usuário pode mover pasta para outra; sistema impede movimentos que criariam ciclos; impede mover para pasta que já contém irmão com o mesmo nome
- **FOLDER-04**: Usuário pode reordenar pasta dentro da mesma pasta pai; ordem final persistida ao salvar
- **FOLDER-05**: Usuário pode excluir pasta; conteúdo (segredos e subpastas) promovido para pasta pai; segredos com nome duplicado na pasta pai renomeados automaticamente com sufixo numérico; usuário avisado sobre renomeações
- **FOLDER-06**: Pasta Geral não pode ser renomeada, movida ou excluída; pode estar vazia
- **FOLDER-07**: Regras de hierarquia: sem dois segredos com mesmo nome na mesma pasta; sem duas subpastas com mesmo nome na mesma pasta pai; ciclos proibidos; todas as pastas navegáveis a partir da Pasta Geral

**Success Criteria:**
1. Creating a folder with a duplicate sibling name is rejected with an error; creating with a unique name inserts it and collapses/expands correctly in the tree
2. Attempting to rename, move, or delete Pasta Geral produces an informative blocked-action message; no change is made
3. Moving a folder into one of its own descendants is blocked with a cycle-prevention error; moving it to a folder that already has a sibling with the same name is also blocked
4. Deleting a folder with contents promotes all children to the parent folder; any name collisions with existing siblings in the parent are auto-renamed with a numeric suffix, and the message bar lists the renamed items

**Plans:** TBD
**UI hint**: yes

---

### Phase 9: Template Management

**Goal:** Users can create, rename, edit, and delete secret templates; templates can be created from existing secrets; secrets created from a template are decoupled — template changes do not affect them.
**Depends on:** Phase 7

**Requirements:**
- **TMPL-01**: Usuário pode criar modelo de segredo com campos personalizados (nome + tipo)
- **TMPL-02**: Usuário pode renomear modelo; nome deve ser único globalmente
- **TMPL-03**: Usuário pode editar estrutura do modelo: adicionar campo, renomear campo, alterar tipo de campo, reordenar campos, excluir campo; campos do modelo permitem nomes duplicados entre si
- **TMPL-04**: Usuário pode excluir modelo de segredo
- **TMPL-05**: Usuário pode criar modelo a partir de segredo existente; campo Observação do segredo não é copiado para o modelo
- **TMPL-06**: Segredo criado a partir de modelo não mantém vínculo — alterações no modelo não afetam segredos existentes; nome do modelo registrado apenas como histórico
- **TMPL-07**: Modelos exibidos em ordem alfabética; não são reordenáveis manualmente

**Success Criteria:**
1. A custom template can be created with arbitrary field names and types; renaming it to an existing template name is rejected; deleting it removes it from the template list
2. Creating a template from an existing secret copies all fields except Observation; the resulting template appears in alphabetical order in the template list
3. Editing a template (adding, renaming, or reordering fields) has no effect on secrets that were previously created from that template; those secrets retain their original structure
4. Templates are always displayed in alphabetical order with no drag-to-reorder option available in the UI

**Plans:** TBD
**UI hint**: yes

---

### Phase 10: Export & Import

**Goal:** Users can export the vault as an unencrypted interchange file (with risk confirmation) and import an interchange file with folder merging and conflict overwriting.
**Depends on:** Phase 7

**Requirements:**
- **VAULT-13**: Usuário pode exportar cofre para arquivo de intercâmbio não criptografado, com aviso de risco e confirmação explícita
- **VAULT-14**: Usuário pode importar cofre de arquivo de intercâmbio, com mesclagem de pastas e sobrescrita de segredos/modelos conflitantes

**Success Criteria:**
1. Initiating an export presents a risk-warning dialog requiring explicit confirmation; cancelling leaves the vault unchanged; confirming writes a human-readable (unencrypted) file to the chosen path
2. Importing an interchange file merges its folders into the live vault, overwrites conflicting secrets and templates with the imported versions, and adds new ones; the tree updates immediately to reflect the merged state
3. Importing a file with no conflicts produces a vault identical to having manually created every element in the interchange file

**Plans:** TBD
**UI hint**: yes

---

## Progress Table

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Foundation & Distribution | 0/? | Not started | — |
| 2. TUI Design System | 0/? | Not started | — |
| 3. Cryptography & File Format | 0/? | Not started | — |
| 4. Vault Lifecycle | 0/? | Not started | — |
| 5. Session Security & Locking | 0/? | Not started | — |
| 6. Vault Navigation & Search | 0/? | Not started | — |
| 7. Secret Management | 0/? | Not started | — |
| 8. Folder Management | 0/? | Not started | — |
| 9. Template Management | 0/? | Not started | — |
| 10. Export & Import | 0/? | Not started | — |

---

## Coverage Validation

| Category | Requirements | Phase |
|----------|-------------|-------|
| DIST | DIST-01, DIST-02 | 1 |
| TUI | TUI-01–10 | 2 |
| FILE | FILE-01–09 | 3 |
| VAULT (lifecycle) | VAULT-01–08, VAULT-12, VAULT-15–16 | 4 |
| SEC (session password) | SEC-01 | 4 |
| VAULT (locking) | VAULT-09–11 | 5 |
| SEC (hardening) | SEC-02–06 | 5 |
| NAV | NAV-01–06 | 6 |
| SECRET | SECRET-01–11 | 7 |
| FIELD | FIELD-01–02 | 7 |
| FOLDER | FOLDER-01–07 | 8 |
| TMPL | TMPL-01–07 | 9 |
| VAULT (exchange) | VAULT-13–14 | 10 |

**Total mapped: 76 / 76 ✓** *(Note: REQUIREMENTS.md initially stated "63" — the actual count is 76; traceability updated accordingly)*

---
*Roadmap created: April 9, 2026*
