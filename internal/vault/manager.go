package vault

import (
	"time"

	"github.com/useful-toys/abditum/internal/crypto"
)

// Manager orchestrates vault operations and maintains session state.
// All vault mutations go through Manager methods.
// Manager knows WHAT operations exist (high-level workflows),
// entities know HOW to execute (validation + mutation logic).
type Manager struct {
	cofre       *Cofre
	repositorio RepositorioCofre
	senha       []byte
	caminho     string
	bloqueado   bool
}

// NewManager creates a new Manager for the given vault and repository.
// The vault is initially unlocked.
func NewManager(cofre *Cofre, repositorio RepositorioCofre) *Manager {
	return &Manager{
		cofre:       cofre,
		repositorio: repositorio,
		senha:       nil,
		caminho:     "",
		bloqueado:   false,
	}
}

// Vault returns a pointer to the managed vault.
// Per D-08: returns live pointer, safety via package encapsulation.
// TUI cannot mutate private fields even with this pointer.
func (m *Manager) Vault() *Cofre {
	if m.bloqueado {
		return nil
	}
	return m.cofre
}

// CriarSegredo creates a new secret in a folder using a template.
// The secret is initialized with campos from the template (empty values),
// estadoSessao set to EstadoModificado, and favorito set to false.
// Returns the created secret or an error if validation fails.
func (m *Manager) CriarSegredo(pasta *Pasta, nome string, modelo *ModeloSegredo) (*Segredo, error) {
	// Validate locked state
	if m.bloqueado {
		return nil, ErrCofreBloqueado
	}

	// Validate creation parameters (pasta not nil, nome unique, modelo not nil)
	if err := pasta.validarCriacaoSegredo(nome, modelo); err != nil {
		return nil, err
	}

	// Create secret with campos from template
	segredo := pasta.criarSegredo(nome, modelo)

	// Set estadoSessao to Modificado (new secret is a modification)
	segredo.estadoSessao = EstadoModificado

	// Update vault state
	now := time.Now().UTC()
	segredo.dataUltimaModificacao = now
	m.cofre.modificado = true
	m.cofre.dataUltimaModificacao = now

	return segredo, nil
}

// ExcluirSegredo marks a secret as deleted (soft delete).
// Per D-14: Delete is reversible until Salvar.
// State transitions: Original→Excluido, Incluido→removed, Modificado→Excluido.
func (m *Manager) ExcluirSegredo(segredo *Segredo) error {
	// Validate locked state
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Validate deletion parameters
	if err := segredo.validarExclusaoSegredo(); err != nil {
		return err
	}

	// Perform soft delete
	segredo.excluirSegredo()

	// Update vault state (only if still attached - EstadoIncluido removes from parent)
	if segredo.estadoSessao == EstadoExcluido {
		now := time.Now().UTC()
		segredo.dataUltimaModificacao = now
		m.cofre.modificado = true
		m.cofre.dataUltimaModificacao = now
	} else {
		// EstadoIncluido was removed from parent - still mark vault as modified
		m.cofre.modificado = true
		m.cofre.dataUltimaModificacao = time.Now().UTC()
	}

	return nil
}

// RestaurarSegredo restores a soft-deleted secret.
// Per D-14: Restore reverses deletion (Excluido → Modificado).
func (m *Manager) RestaurarSegredo(segredo *Segredo) error {
	// Validate locked state
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Validate restoration parameters
	if err := segredo.validarRestauracaoSegredo(); err != nil {
		return err
	}

	// Perform restoration
	segredo.restaurarSegredo()

	// Update vault state
	now := time.Now().UTC()
	segredo.dataUltimaModificacao = now
	m.cofre.modificado = true
	m.cofre.dataUltimaModificacao = now

	return nil
}

// AlternarFavoritoSegredo toggles the favorito flag on a secret.
// Per D-11: favorito is independent of estadoSessao (only updates cofre.modificado).
func (m *Manager) AlternarFavoritoSegredo(segredo *Segredo) error {
	// Validate locked state
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Validate toggle parameters
	if err := segredo.validarAlternarFavorito(); err != nil {
		return err
	}

	// Toggle favorite flag
	segredo.alternarFavorito()

	// Update vault state (cofre.modificado = true, but estadoSessao unchanged per D-11)
	now := time.Now().UTC()
	segredo.dataUltimaModificacao = now
	m.cofre.modificado = true
	m.cofre.dataUltimaModificacao = now

	return nil
}

// DuplicarSegredo creates a copy of a secret with automatic name conflict resolution.
// Per D-27: Uses "(N)" progression for name conflicts: "Name" → "Name (2)" → "Name (3)".
func (m *Manager) DuplicarSegredo(segredo *Segredo) (*Segredo, error) {
	// Validate locked state
	if m.bloqueado {
		return nil, ErrCofreBloqueado
	}

	// Validate duplication parameters
	if err := segredo.validarDuplicacaoSegredo(); err != nil {
		return nil, err
	}

	// Get parent pasta
	pasta := segredo.pasta
	if pasta == nil {
		return nil, ErrPastaInvalida
	}

	// Create duplicate with name conflict resolution
	duplicate := pasta.duplicarSegredo(segredo)

	// Set estadoSessao to Modificado (new content)
	duplicate.estadoSessao = EstadoModificado

	// Update vault state
	now := time.Now().UTC()
	duplicate.dataUltimaModificacao = now
	m.cofre.modificado = true
	m.cofre.dataUltimaModificacao = now

	return duplicate, nil
}

// IsLocked returns true if the vault is currently locked.
func (m *Manager) IsLocked() bool {
	return m.bloqueado
}

// IsModified returns true if the vault has unsaved changes.
func (m *Manager) IsModified() bool {
	if m.cofre == nil {
		return false
	}
	return m.cofre.modificado
}

// AlterarConfiguracoes updates vault timer settings.
// Per D-20: All timers are mandatory (must be > 0).
// Returns ErrConfigInvalida if any timer is <= 0.
// Marks vault as modified and updates timestamp.
func (m *Manager) AlterarConfiguracoes(novasConfig Configuracoes) error {
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Validate all timers > 0 (VAULT-17: all mandatory)
	if novasConfig.tempoBloqueioInatividadeMinutos <= 0 {
		return ErrConfigInvalida
	}
	if novasConfig.tempoOcultarSegredoSegundos <= 0 {
		return ErrConfigInvalida
	}
	if novasConfig.tempoLimparAreaTransferenciaSegundos <= 0 {
		return ErrConfigInvalida
	}

	// Update configuration
	m.cofre.configuracoes = novasConfig
	m.cofre.modificado = true
	m.cofre.dataUltimaModificacao = time.Now().UTC()

	return nil
}

// Lock securely locks the vault by wiping sensitive data from memory.
// Per CRYPTO-04: password and all sensitive field values are overwritten with zeros.
// After locking, vault reference is cleared and IsLocked() returns true.
func (m *Manager) Lock() {
	if m.bloqueado {
		return // Already locked
	}

	// Wipe master password
	if m.senha != nil {
		crypto.Wipe(m.senha)
		m.senha = nil
	}

	// Wipe all sensitive field values recursively
	if m.cofre != nil {
		m.limparCamposSensiveis(m.cofre.pastaGeral)
	}

	// Clear vault reference
	m.cofre = nil
	m.bloqueado = true
}

// limparCamposSensiveis recursively wipes sensitive field values in all secrets.
func (m *Manager) limparCamposSensiveis(pasta *Pasta) {
	if pasta == nil {
		return
	}

	// Wipe sensitive fields in all secrets
	for _, segredo := range pasta.segredos {
		for i := range segredo.campos {
			if segredo.campos[i].tipo == TipoCampoSensivel {
				crypto.Wipe(segredo.campos[i].valor)
				segredo.campos[i].valor = nil
			}
		}
		// Wipe observation (always common type but still sensitive content)
		crypto.Wipe(segredo.observacao.valor)
		segredo.observacao.valor = nil
	}

	// Recurse into subfolders
	for _, subpasta := range pasta.subpastas {
		m.limparCamposSensiveis(subpasta)
	}
}

// Salvar persists the vault to storage using two-phase atomic commit.
// Phase 1: Create deep copy with StateExcluido filtered (live vault untouched)
// Phase 2: Persist via repository (if fails, live vault unchanged)
// Phase 3: Finalize deletions in memory only after successful save
// Per D-17: Guarantees atomicity — save failure doesn't cause data loss.
func (m *Manager) Salvar() error {
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Phase 1: Prepare immutable snapshot (filters excluido)
	snapshot := m.prepararSnapshot()

	// Phase 2: Persist snapshot
	if err := m.repositorio.Salvar(snapshot); err != nil {
		return err // Live vault unchanged on failure
	}

	// Phase 3: Finalize deletions only after successful save
	m.finalizarExclusoes()
	m.cofre.modificado = false

	return nil
}

// prepararSnapshot creates a deep copy of the vault with StateExcluido secrets filtered.
// Returns immutable snapshot safe for serialization.
func (m *Manager) prepararSnapshot() *Cofre {
	snapshot := &Cofre{
		dataCriacao:           m.cofre.dataCriacao,
		dataUltimaModificacao: m.cofre.dataUltimaModificacao,
		configuracoes:         m.cofre.configuracoes,
		modelos:               make([]*ModeloSegredo, len(m.cofre.modelos)),
		modificado:            false,
	}

	// Deep copy models
	for i, modelo := range m.cofre.modelos {
		snapshot.modelos[i] = m.copiarModelo(modelo)
	}

	// Deep copy pasta hierarchy (filters excluido)
	snapshot.pastaGeral = m.copiarPastaRecursivamente(m.cofre.pastaGeral, true)

	return snapshot
}

// copiarPastaRecursivamente creates deep copy of folder tree.
// If filtrarExcluidos=true, skips secrets with EstadoExcluido.
func (m *Manager) copiarPastaRecursivamente(pasta *Pasta, filtrarExcluidos bool) *Pasta {
	copia := &Pasta{
		nome:      pasta.nome,
		pai:       nil, // Will be set during tree reconstruction
		subpastas: make([]*Pasta, 0, len(pasta.subpastas)),
		segredos:  make([]*Segredo, 0),
	}

	// Copy subfolders recursively
	for _, sub := range pasta.subpastas {
		subCopia := m.copiarPastaRecursivamente(sub, filtrarExcluidos)
		subCopia.pai = copia
		copia.subpastas = append(copia.subpastas, subCopia)
	}

	// Copy secrets (filter excluido if requested)
	for _, segredo := range pasta.segredos {
		if filtrarExcluidos && segredo.estadoSessao == EstadoExcluido {
			continue // Skip deleted secrets
		}
		segredoCopia := m.copiarSegredo(segredo)
		segredoCopia.pasta = copia
		copia.segredos = append(copia.segredos, segredoCopia)
	}

	return copia
}

// copiarSegredo creates deep copy of secret with all fields.
func (m *Manager) copiarSegredo(s *Segredo) *Segredo {
	copia := &Segredo{
		nome:                  s.nome,
		favorito:              s.favorito,
		estadoSessao:          s.estadoSessao,
		dataCriacao:           s.dataCriacao,
		dataUltimaModificacao: s.dataUltimaModificacao,
		campos:                make([]CampoSegredo, len(s.campos)),
		observacao:            m.copiarCampo(s.observacao),
	}

	for i, campo := range s.campos {
		copia.campos[i] = m.copiarCampo(campo)
	}

	return copia
}

// copiarCampo creates deep copy of field with []byte value copy.
func (m *Manager) copiarCampo(c CampoSegredo) CampoSegredo {
	valorCopia := make([]byte, len(c.valor))
	copy(valorCopia, c.valor)
	return CampoSegredo{
		nome:  c.nome,
		tipo:  c.tipo,
		valor: valorCopia,
	}
}

// copiarModelo creates deep copy of template.
func (m *Manager) copiarModelo(modelo *ModeloSegredo) *ModeloSegredo {
	copia := &ModeloSegredo{
		nome:   modelo.nome,
		campos: make([]CampoModelo, len(modelo.campos)),
	}
	copy(copia.campos, modelo.campos)
	return copia
}

// finalizarExclusoes permanently removes StateExcluido secrets from live vault.
// Called only after successful save.
func (m *Manager) finalizarExclusoes() {
	m.removerExcluidosRecursivamente(m.cofre.pastaGeral)
}

// removerExcluidosRecursivamente removes deleted secrets from folder tree.
func (m *Manager) removerExcluidosRecursivamente(pasta *Pasta) {
	if pasta == nil {
		return
	}

	// Filter out excluido secrets
	segredosAtivos := make([]*Segredo, 0, len(pasta.segredos))
	for _, segredo := range pasta.segredos {
		if segredo.estadoSessao != EstadoExcluido {
			segredosAtivos = append(segredosAtivos, segredo)
		}
	}
	pasta.segredos = segredosAtivos

	// Recurse into subfolders
	for _, subpasta := range pasta.subpastas {
		m.removerExcluidosRecursivamente(subpasta)
	}
}

// Template Operations

// CriarModelo creates a new template with the given name and fields.
// Templates are automatically sorted alphabetically after insertion (TPL-02, TPL-06).
// Per D-29: "Observação" is a reserved field name and is prohibited.
// Returns the created template or an error.
func (m *Manager) CriarModelo(nome string, campos []CampoModelo) (*ModeloSegredo, error) {
	if m.bloqueado {
		return nil, ErrCofreBloqueado
	}

	// Validation phase
	if err := m.cofre.validarCriacaoModelo(nome, campos); err != nil {
		return nil, err
	}

	// Mutation phase
	modelo := m.cofre.criarModelo(nome, campos)
	m.cofre.modificado = true
	m.cofre.dataUltimaModificacao = time.Now().UTC()

	return modelo, nil
}

// RenomearModelo renames a template and re-sorts the template list alphabetically.
// Per TPL-02, TPL-06, D-23: templates always displayed in alphabetical order.
// Per D-12: only marks modified if name actually changes.
// Returns error if name conflicts or validation fails.
func (m *Manager) RenomearModelo(modelo *ModeloSegredo, novoNome string) error {
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Validation phase
	if err := modelo.validarRenomear(m.cofre, novoNome); err != nil {
		return err
	}

	// Mutation phase (returns true if actually changed)
	if modelo.renomear(novoNome) {
		m.cofre.modificado = true
		m.cofre.dataUltimaModificacao = time.Now().UTC()
	}

	return nil
}

// ExcluirModelo deletes a template from the vault.
// Per TPL-04, D-26: templates can be deleted unless referenced by a secret.
// Returns ErrModeloEmUso if any secret references the template.
func (m *Manager) ExcluirModelo(modelo *ModeloSegredo) error {
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Validation phase
	if err := modelo.validarExclusao(m.cofre); err != nil {
		return err
	}

	// Mutation phase
	modelo.excluir(m.cofre)
	m.cofre.modificado = true
	m.cofre.dataUltimaModificacao = time.Now().UTC()

	return nil
}

// AdicionarCampo adds a field to a template at the specified position.
// Per D-29: "Observação" is a reserved field name and is prohibited.
// Position is 0-indexed. Position == len(campos) means append.
// Returns error if position is invalid or field name is reserved.
func (m *Manager) AdicionarCampo(modelo *ModeloSegredo, nome string, tipo TipoCampo, posicao int) error {
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Validation phase
	if err := modelo.validarAdicionarCampo(nome, posicao); err != nil {
		return err
	}

	// Mutation phase
	modelo.adicionarCampo(nome, tipo, posicao)
	m.cofre.modificado = true
	m.cofre.dataUltimaModificacao = time.Now().UTC()

	return nil
}

// RemoverCampo removes a field from a template by index.
// Returns error if index is out of bounds.
func (m *Manager) RemoverCampo(modelo *ModeloSegredo, indice int) error {
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Validation phase
	if err := modelo.validarRemoverCampo(indice); err != nil {
		return err
	}

	// Mutation phase
	modelo.removerCampo(indice)
	m.cofre.modificado = true
	m.cofre.dataUltimaModificacao = time.Now().UTC()

	return nil
}

// ReordenarCampo moves a field from one position to another in a template.
// Both indices must be valid (0 <= index < len(campos)).
// Returns error if indices are out of bounds.
func (m *Manager) ReordenarCampo(modelo *ModeloSegredo, indiceOrigem, indiceDestino int) error {
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Validation phase
	if err := modelo.validarReordenarCampo(indiceOrigem, indiceDestino); err != nil {
		return err
	}

	// Mutation phase
	modelo.reordenarCampo(indiceOrigem, indiceDestino)
	m.cofre.modificado = true
	m.cofre.dataUltimaModificacao = time.Now().UTC()

	return nil
}

// Folder Operations

// CriarPasta creates a new subfolder in the specified parent folder at the given position.
// Position semantics (D-22): 0-indexed, position == len means append at end.
// Validates: name non-empty, length <= 255, unique in parent, valid position [0, len].
// Marks vault as modified and updates timestamp per D-05.
func (m *Manager) CriarPasta(pai *Pasta, nome string, posicao int) (*Pasta, error) {
	if m.bloqueado {
		return nil, ErrCofreBloqueado
	}

	// Phase 1: Validate (can fail)
	if err := pai.validarCriacaoSubpasta(nome, posicao); err != nil {
		return nil, err
	}

	// Phase 2: Mutate (cannot fail after validation per D-05)
	novaPasta := pai.criarSubpasta(nome, posicao)

	// Update global state
	m.cofre.modificado = true
	m.cofre.dataUltimaModificacao = time.Now().UTC()

	return novaPasta, nil
}

// RenomearPasta renames a folder with Pasta Geral protection and change detection.
// Per D-12: Only marks vault as modified if name actually changes (no-op if same name).
// Returns ErrPastaGeralProtected if attempting to rename Pasta Geral.
// Validates: not Pasta Geral, name non-empty, length <= 255, unique among siblings.
func (m *Manager) RenomearPasta(pasta *Pasta, novoNome string) error {
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Check Pasta Geral protection (additional Manager-level check)
	if pasta == m.cofre.pastaGeral {
		return ErrPastaGeralProtected
	}

	// Phase 1: Validate (can fail)
	if err := pasta.validarRenomear(novoNome); err != nil {
		return err
	}

	// Phase 2: Mutate and check if actual change occurred (D-12)
	alterado, err := pasta.renomear(novoNome)
	if err != nil {
		return err // Should never happen after validation per D-05
	}

	// Only update global state if actual change (D-12)
	if alterado {
		m.cofre.modificado = true
		m.cofre.dataUltimaModificacao = time.Now().UTC()
	}

	return nil
}

// MoverPasta moves a folder to a new parent with cycle detection.
// Performs full ancestor walk to detect cycles (ROADMAP pitfall).
// Returns ErrCycleDetected if destino is a descendant of pasta.
// Validates: not Pasta Geral, not moving to self, no cycle, name unique in destination.
func (m *Manager) MoverPasta(pasta *Pasta, destino *Pasta) error {
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Check Pasta Geral protection (additional Manager-level check)
	if pasta == m.cofre.pastaGeral {
		return ErrPastaGeralProtected
	}

	// Phase 1: Validate (can fail)
	if err := pasta.validarMover(destino); err != nil {
		return err
	}

	// Phase 2: Mutate (cannot fail after validation per D-05)
	pasta.mover(destino)

	// Update global state
	m.cofre.modificado = true
	m.cofre.dataUltimaModificacao = time.Now().UTC()

	return nil
}

// ReposicionarPasta moves a folder to a new position within its parent.
// Position semantics (D-22): 0-indexed, valid range [0, len-1].
// Edge case (D-12, D-23): Moving to current position is a no-op - returns nil without marking modified.
// Validates: position in valid range.
func (m *Manager) ReposicionarPasta(pasta *Pasta, novaPosicao int) error {
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Phase 1: Validate (can fail)
	if err := pasta.validarReposicionar(novaPosicao); err != nil {
		return err
	}

	// Phase 2: Mutate and check if actual change occurred (D-12, D-23)
	alterado, err := pasta.reposicionar(novaPosicao)
	if err != nil {
		return err // Should never happen after validation per D-05
	}

	// Only update global state if actual change (D-12, D-23)
	if alterado {
		m.cofre.modificado = true
		m.cofre.dataUltimaModificacao = time.Now().UTC()
	}

	return nil
}

// SubirPastaNaPosicao moves a folder one position up (position - 1) within its parent.
// Edge case (D-23): If already at position 0, this is a no-op - returns nil without marking modified.
func (m *Manager) SubirPastaNaPosicao(pasta *Pasta) error {
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Get current position
	posicaoAtual := pasta.obterPosicaoAtual()

	// Edge case: already at top (D-23 no-op)
	if posicaoAtual == 0 {
		return nil // No-op, don't mark modified
	}

	// Move to position - 1
	return m.ReposicionarPasta(pasta, posicaoAtual-1)
}

// DescerPastaNaPosicao moves a folder one position down (position + 1) within its parent.
// Edge case (D-23): If already at last position, this is a no-op - returns nil without marking modified.
func (m *Manager) DescerPastaNaPosicao(pasta *Pasta) error {
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Get current position and max position
	posicaoAtual := pasta.obterPosicaoAtual()
	maxPosicao := len(pasta.pai.subpastas) - 1

	// Edge case: already at bottom (D-23 no-op)
	if posicaoAtual == maxPosicao {
		return nil // No-op, don't mark modified
	}

	// Move to position + 1
	return m.ReposicionarPasta(pasta, posicaoAtual+1)
}

// ExcluirPasta deletes a folder and promotes its children (secrets and subfolders) to the parent folder.
// Handles name conflicts automatically:
// - Subfolders with conflicting names: merge contents recursively
// - Secrets with conflicting names: rename with numeric suffix "(N)" and track in returned slice
// Per FOLDER-05: Secrets with EstadoExcluido retain that state when promoted.
// Per D-27: Hard delete (immediate removal from hierarchy).
// Returns slice of Renomeacao to inform TUI which secrets were renamed.
// Validates: locked vault, Pasta Geral protection.
func (m *Manager) ExcluirPasta(pasta *Pasta) ([]Renomeacao, error) {
	// Validate locked state
	if m.bloqueado {
		return nil, ErrCofreBloqueado
	}

	// Validate not Pasta Geral
	if err := pasta.validarExclusao(); err != nil {
		return nil, err
	}

	// Perform deletion with promotion and conflict resolution
	renomeacoes := pasta.excluir(pasta.pai)

	// Mark vault as modified
	m.cofre.modificado = true

	return renomeacoes, nil
}

// Secret content mutation operations

// RenomearSegredo renames a secret within its folder.
// Per D-11: marks estadoSessao = Modificado if currently Original.
// Per D-12: no-op if name unchanged (doesn't mark modified).
// Validates: locked vault, name non-empty, length <= 255, unique in parent.
func (m *Manager) RenomearSegredo(segredo *Segredo, novoNome string) error {
	// Validate locked state
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Validate rename parameters
	if err := segredo.validarRenomear(novoNome); err != nil {
		return err
	}

	// Perform rename
	alterado, err := segredo.renomear(novoNome)
	if err != nil {
		return err
	}

	// Only mark vault modified if actual change occurred (D-12)
	if alterado {
		segredo.dataUltimaModificacao = time.Now().UTC()
		m.cofre.modificado = true
		m.cofre.dataUltimaModificacao = time.Now().UTC()
	}

	return nil
}

// EditarCampoSegredo edits a field value in a secret by index.
// Per D-11: marks estadoSessao = Modificado if currently Original.
// Per D-12: no-op if value unchanged (doesn't mark modified).
// Validates: locked vault, index within valid range [0, len(campos)-1].
func (m *Manager) EditarCampoSegredo(segredo *Segredo, indice int, novoValor []byte) error {
	// Validate locked state
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Validate edit parameters
	if err := segredo.validarEditarCampo(indice, novoValor); err != nil {
		return err
	}

	// Perform edit
	alterado, err := segredo.editarCampo(indice, novoValor)
	if err != nil {
		return err
	}

	// Only mark vault modified if actual change occurred (D-12)
	if alterado {
		segredo.dataUltimaModificacao = time.Now().UTC()
		m.cofre.modificado = true
		m.cofre.dataUltimaModificacao = time.Now().UTC()
	}

	return nil
}

// EditarObservacao edits the observação field of a secret.
// Per D-11: marks estadoSessao = Modificado if currently Original.
// Per D-12: no-op if value unchanged (doesn't mark modified).
// Per D-29: observação is separate field, not in campos slice.
// Validates: locked vault, length <= 1000 chars.
func (m *Manager) EditarObservacao(segredo *Segredo, novoTexto string) error {
	// Validate locked state
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Validate length <= 1000 chars
	if len(novoTexto) > 1000 {
		return ErrObservacaoMuitoLonga
	}

	// Perform edit
	alterado, err := segredo.editarObservacao(novoTexto)
	if err != nil {
		return err
	}

	// Only mark vault modified if actual change occurred (D-12)
	if alterado {
		segredo.dataUltimaModificacao = time.Now().UTC()
		m.cofre.modificado = true
		m.cofre.dataUltimaModificacao = time.Now().UTC()
	}

	return nil
}

// MoverSegredo moves a secret to a different folder at the specified position.
// Per D-16: Move is structural operation - does NOT change estadoSessao.
// Only updates cofre.modificado.
// Validates: locked vault, destino not nil, name unique in destino.
func (m *Manager) MoverSegredo(segredo *Segredo, destino *Pasta, posicao int) error {
	// Validate locked state
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Validate move parameters
	if err := segredo.validarMover(destino); err != nil {
		return err
	}

	// Perform move
	segredo.mover(destino, posicao)

	// Update vault state (cofre.modificado = true, but NOT estadoSessao per D-16)
	m.cofre.modificado = true
	m.cofre.dataUltimaModificacao = time.Now().UTC()

	return nil
}

// ReposicionarSegredo moves a secret to a new position within its current folder.
// Per D-16: Reposition is structural operation - does NOT change estadoSessao.
// Only updates cofre.modificado if position actually changes.
// Per D-12, D-23: Moving to current position is a no-op (returns nil without marking modified).
// Validates: locked vault, position in valid range [0, len-1].
func (m *Manager) ReposicionarSegredo(segredo *Segredo, novaPosicao int) error {
	// Validate locked state
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Validate reposition parameters
	if err := segredo.validarReposicionarSegredo(novaPosicao); err != nil {
		return err
	}

	// Perform reposition (returns true if position actually changed)
	alterado, err := segredo.reposicionarSegredo(novaPosicao)
	if err != nil {
		return err // Should never happen after validation per D-05
	}

	// Only update vault state if position changed (D-12, D-23)
	if alterado {
		m.cofre.modificado = true
		m.cofre.dataUltimaModificacao = time.Now().UTC()
	}

	return nil
}

// SubirSegredoNaPosicao moves a secret one position up (position - 1) within its folder.
// Per D-16: Structural operation - does NOT change estadoSessao.
// Edge case (D-23): If already at position 0, this is a no-op (returns nil without marking modified).
func (m *Manager) SubirSegredoNaPosicao(segredo *Segredo) error {
	// Validate locked state
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Get current position
	posicaoAtual := segredo.obterPosicaoAtualSegredo()

	// Edge case: already at top (D-23 no-op)
	if posicaoAtual == 0 {
		return nil // No-op, don't mark modified
	}

	// Move to position - 1
	return m.ReposicionarSegredo(segredo, posicaoAtual-1)
}

// DescerSegredoNaPosicao moves a secret one position down (position + 1) within its folder.
// Per D-16: Structural operation - does NOT change estadoSessao.
// Edge case (D-23): If already at last position, this is a no-op (returns nil without marking modified).
func (m *Manager) DescerSegredoNaPosicao(segredo *Segredo) error {
	// Validate locked state
	if m.bloqueado {
		return ErrCofreBloqueado
	}

	// Get current position and max position
	posicaoAtual := segredo.obterPosicaoAtualSegredo()
	maxPosicao := len(segredo.pasta.segredos) - 1

	// Edge case: already at bottom (D-23 no-op)
	if posicaoAtual == maxPosicao {
		return nil // No-op, don't mark modified
	}

	// Move to position + 1
	return m.ReposicionarSegredo(segredo, posicaoAtual+1)
}
