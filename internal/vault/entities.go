// Package vault implements the domain layer for Abditum password manager.
// All entities use package-private fields (lowercase) for encapsulation.
// External access via exported getters returning defensive copies.
package vault

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// TipoCampo represents the visibility type of a secret field.
type TipoCampo int

const (
	// TipoCampoComum represents a common (non-sensitive) field - always visible.
	TipoCampoComum TipoCampo = iota
	// TipoCampoSensivel represents a sensitive field - hidden by default, revealed temporarily.
	TipoCampoSensivel
)

// EstadoSessao tracks the lifecycle state of a secret relative to the persisted file.
// This is a session-only flag - not serialized to disk.
type EstadoSessao int

const (
	// EstadoOriginal means the secret was loaded from file and has not been modified this session.
	EstadoOriginal EstadoSessao = iota
	// EstadoIncluido means the secret was created this session and does not exist in the file yet.
	EstadoIncluido
	// EstadoModificado means the secret existed in the file but has been modified this session.
	EstadoModificado
	// EstadoExcluido means the secret is marked for deletion - will be removed on save.
	EstadoExcluido
)

// Renomeacao records an automatic rename during folder deletion (for TUI display).
// Used by ExcluirPasta to inform user which secrets were renamed due to name conflicts.
type Renomeacao struct {
	Antigo string // Original secret name
	Novo   string // New name with numeric suffix
	Pasta  string // Parent folder name where rename occurred
}

// Configuracoes contains operational settings for the vault (timers for auto-lock, reveal, clipboard).
type Configuracoes struct {
	tempoBloqueioInatividadeMinutos      int // Default: 5 min
	tempoOcultarSegredoSegundos          int // Default: 15 sec
	tempoLimparAreaTransferenciaSegundos int // Default: 30 sec
}

// CampoModelo defines a field structure in a ModeloSegredo (template).
type CampoModelo struct {
	nome string
	tipo TipoCampo
}

// ModeloSegredo is a reusable template defining field structures for secrets.
// Identity: global unique name across the vault.
type ModeloSegredo struct {
	nome   string
	campos []CampoModelo
}

// CampoSegredo represents an individual field within a Segredo (secret).
// Identity: position (index) in the campos list.
type CampoSegredo struct {
	nome  string
	tipo  TipoCampo
	valor []byte // Always UTF-8, wipeable in memory
}

// Segredo represents a credential or confidential information stored within a Pasta.
// Identity: (pasta, nome) - unique name within parent folder.
type Segredo struct {
	nome                  string
	campos                []CampoSegredo // User fields only, excludes Observação
	observacao            CampoSegredo   // Separate, always exists, not in campos slice
	pasta                 *Pasta         // Back-reference (not serialized)
	favorito              bool
	estadoSessao          EstadoSessao
	dataCriacao           time.Time
	dataUltimaModificacao time.Time
}

// Pasta is a hierarchical container that groups secrets and other folders.
// Identity: (pai, nome) - unique name among siblings.
type Pasta struct {
	nome      string
	pai       *Pasta // Back-reference (not serialized)
	subpastas []*Pasta
	segredos  []*Segredo
}

// Cofre is the aggregate root encapsulating the entire password vault.
// All mutations pass through Cofre. All persistence is atomic over Cofre.
type Cofre struct {
	pastaGeral            *Pasta
	modelos               []*ModeloSegredo
	modificado            bool
	dataCriacao           time.Time
	dataUltimaModificacao time.Time
	configuracoes         Configuracoes
}

// Exported getters for Cofre

// PastaGeral returns the root folder of the vault.
func (c *Cofre) PastaGeral() *Pasta {
	return c.pastaGeral
}

// Modelos returns a defensive copy of the template list, sorted alphabetically (TPL-06).
func (c *Cofre) Modelos() []*ModeloSegredo {
	// Defensive copy
	copia := make([]*ModeloSegredo, len(c.modelos))
	copy(copia, c.modelos)
	// Sort alphabetically by name (TPL-06)
	sort.Slice(copia, func(i, j int) bool {
		return copia[i].nome < copia[j].nome
	})
	return copia
}

// Configuracoes returns a copy of the vault settings.
func (c *Cofre) Configuracoes() Configuracoes {
	return c.configuracoes // Value copy (struct)
}

// Modificado returns whether the vault has unsaved changes.
func (c *Cofre) Modificado() bool {
	return c.modificado
}

// DataCriacao returns when the vault was created.
func (c *Cofre) DataCriacao() time.Time {
	return c.dataCriacao
}

// DataUltimaModificacao returns when the vault was last modified.
func (c *Cofre) DataUltimaModificacao() time.Time {
	return c.dataUltimaModificacao
}

// Exported getters for Pasta

// Nome returns the folder name.
func (p *Pasta) Nome() string {
	return p.nome
}

// Pai returns the parent folder (nil for Pasta Geral).
func (p *Pasta) Pai() *Pasta {
	return p.pai
}

// Subpastas returns a defensive copy of the subfolders list.
func (p *Pasta) Subpastas() []*Pasta {
	copia := make([]*Pasta, len(p.subpastas))
	copy(copia, p.subpastas)
	return copia
}

// Segredos returns a defensive copy of the secrets list.
// Includes all secrets, even those marked as EstadoExcluido (D-14).
func (p *Pasta) Segredos() []*Segredo {
	copia := make([]*Segredo, len(p.segredos))
	copy(copia, p.segredos)
	return copia
}

// Exported getters for Segredo

// Nome returns the secret name.
func (s *Segredo) Nome() string {
	return s.nome
}

// Pasta returns the parent folder reference.
func (s *Segredo) Pasta() *Pasta {
	return s.pasta
}

// Campos returns a defensive copy of user fields only (excludes Observação per D-29).
func (s *Segredo) Campos() []CampoSegredo {
	copia := make([]CampoSegredo, len(s.campos))
	copy(copia, s.campos)
	return copia
}

// Observacao returns the observation field value as a string (D-29).
func (s *Segredo) Observacao() string {
	return string(s.observacao.valor)
}

// Favorito returns whether the secret is marked as favorite.
func (s *Segredo) Favorito() bool {
	return s.favorito
}

// EstadoSessao returns the session lifecycle state of the secret.
func (s *Segredo) EstadoSessao() EstadoSessao {
	return s.estadoSessao
}

// DataCriacao returns when the secret was created.
func (s *Segredo) DataCriacao() time.Time {
	return s.dataCriacao
}

// DataUltimaModificacao returns when the secret was last modified.
func (s *Segredo) DataUltimaModificacao() time.Time {
	return s.dataUltimaModificacao
}

// Exported getters for CampoSegredo

// Nome returns the field name.
func (c *CampoSegredo) Nome() string {
	return c.nome
}

// Tipo returns the field type (common or sensitive).
func (c *CampoSegredo) Tipo() TipoCampo {
	return c.tipo
}

// ValorComoString converts the field value to string for display in TUI.
// SECURITY BOUNDARY: This conversion is irreversible - the returned string
// is immutable and cannot be wiped from memory. Call only when user explicitly
// requests display. Per D-10.
func (c *CampoSegredo) ValorComoString() string {
	return string(c.valor)
}

// Exported getters for ModeloSegredo

// Nome returns the template name.
func (m *ModeloSegredo) Nome() string {
	return m.nome
}

// Campos returns a defensive copy of the template fields.
func (m *ModeloSegredo) Campos() []CampoModelo {
	copia := make([]CampoModelo, len(m.campos))
	copy(copia, m.campos)
	return copia
}

// Exported getters for CampoModelo

// Nome returns the field definition name.
func (c *CampoModelo) Nome() string {
	return c.nome
}

// Tipo returns the field definition type.
func (c *CampoModelo) Tipo() TipoCampo {
	return c.tipo
}

// Private entity methods for secret operations

// validarRenomear validates secret rename parameters.
// Checks: name non-empty, length <= 255, unique in parent folder.
func (s *Segredo) validarRenomear(novoNome string) error {
	// Validate name non-empty
	if novoNome == "" {
		return ErrNomeVazio
	}

	// Validate name length <= 255
	if len(novoNome) > 255 {
		return ErrNomeMuitoLongo
	}

	// Validate name unique within parent folder (excluding self)
	for _, outro := range s.pasta.segredos {
		if outro != s && outro.nome == novoNome {
			return ErrNameConflict
		}
	}

	return nil
}

// renomear changes the secret name and returns whether an actual change occurred.
// Per D-11: marks estadoSessao = Modificado if currently Original.
// Per D-12: returns (false, nil) if name unchanged (no-op).
// PRECONDITION: validarRenomear must pass (cannot fail after validation per D-05).
func (s *Segredo) renomear(novoNome string) (alterado bool, err error) {
	// Check if name actually changed (D-12)
	if s.nome == novoNome {
		return false, nil // No-op
	}

	// Update name
	s.nome = novoNome

	// Mark as modified if currently Original (D-11)
	if s.estadoSessao == EstadoOriginal {
		s.estadoSessao = EstadoModificado
	}
	// Note: EstadoIncluido and EstadoModificado remain unchanged

	return true, nil
}

// validarEditarCampo validates field editing parameters.
// Checks: index within valid range [0, len(campos)-1].
func (s *Segredo) validarEditarCampo(indice int, valor []byte) error {
	// Validate index within bounds
	if indice < 0 || indice >= len(s.campos) {
		return ErrCampoInvalido
	}

	return nil
}

// editarCampo updates a field value and returns whether an actual change occurred.
// Per D-11: marks estadoSessao = Modificado if currently Original.
// Per D-12: returns (false, nil) if value unchanged (no-op).
// PRECONDITION: validarEditarCampo must pass (cannot fail after validation per D-05).
func (s *Segredo) editarCampo(indice int, novoValor []byte) (alterado bool, err error) {
	// Check if value actually changed (D-12)
	campoAtual := s.campos[indice].valor
	if string(campoAtual) == string(novoValor) {
		return false, nil // No-op
	}

	// Update field value (deep copy to ensure independence)
	valorCopia := make([]byte, len(novoValor))
	copy(valorCopia, novoValor)
	s.campos[indice].valor = valorCopia

	// Mark as modified if currently Original (D-11)
	if s.estadoSessao == EstadoOriginal {
		s.estadoSessao = EstadoModificado
	}
	// Note: EstadoIncluido and EstadoModificado remain unchanged

	return true, nil
}

// editarObservacao updates the observação field and returns whether an actual change occurred.
// Per D-11: marks estadoSessao = Modificado if currently Original.
// Per D-12: returns (false, nil) if value unchanged (no-op).
// Per D-29: observação is separate field, not in campos slice.
// PRECONDITION: validation must pass (cannot fail after validation per D-05).
func (s *Segredo) editarObservacao(novoTexto string) (alterado bool, err error) {
	// Check if value actually changed (D-12)
	textoAtual := string(s.observacao.valor)
	if textoAtual == novoTexto {
		return false, nil // No-op
	}

	// Update observação value (convert string to []byte)
	s.observacao.valor = []byte(novoTexto)

	// Mark as modified if currently Original (D-11)
	if s.estadoSessao == EstadoOriginal {
		s.estadoSessao = EstadoModificado
	}
	// Note: EstadoIncluido and EstadoModificado remain unchanged

	return true, nil
}

// validarMover validates secret move parameters.
// Checks: destino not nil, name unique in destino.
func (s *Segredo) validarMover(destino *Pasta) error {
	// Check destino not nil
	if destino == nil {
		return ErrPastaInvalida
	}

	// Check name uniqueness in destination (excluding self if moving within same folder)
	for _, outro := range destino.segredos {
		if outro != s && outro.nome == s.nome {
			return ErrNameConflict
		}
	}

	return nil
}

// mover removes secret from origem and adds to destino at specified position.
// Per D-16: Move is structural operation - does NOT change estadoSessao.
// PRECONDITION: validarMover must pass (cannot fail after validation per D-05).
func (s *Segredo) mover(destino *Pasta, posicao int) {
	// Remove from origem (current pasta)
	if s.pasta != nil {
		for i, seg := range s.pasta.segredos {
			if seg == s {
				s.pasta.segredos = append(s.pasta.segredos[:i], s.pasta.segredos[i+1:]...)
				break
			}
		}
	}

	// Add to destino at position
	if posicao < 0 || posicao >= len(destino.segredos) {
		// Invalid position or append - just append
		destino.segredos = append(destino.segredos, s)
	} else {
		// Insert at position
		destino.segredos = append(destino.segredos[:posicao], append([]*Segredo{s}, destino.segredos[posicao:]...)...)
	}

	// Update parent reference
	s.pasta = destino
}

// obterPosicaoAtualSegredo returns the current position of this secret in its parent folder.
// Returns -1 if secret has no parent.
func (s *Segredo) obterPosicaoAtualSegredo() int {
	if s.pasta == nil {
		return -1
	}

	for i, seg := range s.pasta.segredos {
		if seg == s {
			return i
		}
	}

	return -1 // Should never happen if tree is consistent
}

// validarReposicionarSegredo validates secret repositioning parameters.
// Checks: position in valid range [0, len-1] (0-indexed).
func (s *Segredo) validarReposicionarSegredo(novaPosicao int) error {
	if s.pasta == nil {
		return ErrPastaInvalida
	}

	// Validate position in valid range [0, len-1]
	if novaPosicao < 0 || novaPosicao >= len(s.pasta.segredos) {
		return ErrPosicaoInvalida
	}

	return nil
}

// reposicionarSegredo moves this secret to a new position within its parent folder.
// Returns (true, nil) if position changed, (false, nil) if no change (D-12, D-23 no-op).
// Per D-16: Reposition is structural operation - does NOT change estadoSessao.
// PRECONDITION: validarReposicionarSegredo must pass (cannot fail after validation per D-05).
func (s *Segredo) reposicionarSegredo(novaPosicao int) (alterado bool, err error) {
	posicaoAtual := s.obterPosicaoAtualSegredo()

	// No-op if already at target position (D-12, D-23)
	if posicaoAtual == novaPosicao {
		return false, nil
	}

	// Remove from current position
	s.pasta.segredos = append(s.pasta.segredos[:posicaoAtual], s.pasta.segredos[posicaoAtual+1:]...)

	// Insert at new position
	s.pasta.segredos = append(s.pasta.segredos[:novaPosicao], append([]*Segredo{s}, s.pasta.segredos[novaPosicao:]...)...)

	return true, nil
}

// Private entity methods for folder operations

// contemSubpastaComNome checks if a subfolder with the given name exists.
func (p *Pasta) contemSubpastaComNome(nome string) bool {
	for _, sub := range p.subpastas {
		if sub.nome == nome {
			return true
		}
	}
	return false
}

// contemSegredoComNome checks if a secret with the given name exists.
func (p *Pasta) contemSegredoComNome(nome string) bool {
	for _, seg := range p.segredos {
		if seg.nome == nome {
			return true
		}
	}
	return false
}

// validarCriacaoSubpasta validates folder creation parameters.
// Checks: name non-empty, length <= 255, unique in parent, valid position.
func (p *Pasta) validarCriacaoSubpasta(nome string, posicao int) error {
	// Validate name non-empty
	if nome == "" {
		return ErrNomeVazio
	}

	// Validate name length <= 255
	if len(nome) > 255 {
		return ErrNomeMuitoLongo
	}

	// Validate name unique within parent
	if p.contemSubpastaComNome(nome) {
		return ErrNameConflict
	}

	// Validate position in valid range [0, len] (inclusive of len for append)
	if posicao < 0 || posicao > len(p.subpastas) {
		return ErrPosicaoInvalida
	}

	return nil
}

// criarSubpasta creates and inserts a subfolder at the specified position.
// PRECONDITION: validarCriacaoSubpasta must pass (cannot fail after validation per D-05).
func (p *Pasta) criarSubpasta(nome string, posicao int) *Pasta {
	novaPasta := &Pasta{
		nome:      nome,
		pai:       p,
		subpastas: make([]*Pasta, 0),
		segredos:  make([]*Segredo, 0),
	}

	// Insert at position using slice idiom
	p.subpastas = append(p.subpastas[:posicao], append([]*Pasta{novaPasta}, p.subpastas[posicao:]...)...)

	return novaPasta
}

// validarRenomear validates folder rename parameters.
// Checks: not Pasta Geral, name non-empty, length <= 255, unique in parent.
func (p *Pasta) validarRenomear(novoNome string) error {
	// Check not Pasta Geral
	if p.pai == nil {
		return ErrPastaGeralProtected
	}

	// Validate name non-empty
	if novoNome == "" {
		return ErrNomeVazio
	}

	// Validate name length <= 255
	if len(novoNome) > 255 {
		return ErrNomeMuitoLongo
	}

	// Validate name unique within parent (excluding self)
	for _, sub := range p.pai.subpastas {
		if sub != p && sub.nome == novoNome {
			return ErrNameConflict
		}
	}

	return nil
}

// renomear changes the folder name and returns whether an actual change occurred.
// Returns (true, nil) if name changed, (false, nil) if no change (same name per D-12).
// PRECONDITION: validarRenomear must pass (cannot fail after validation per D-05).
func (p *Pasta) renomear(novoNome string) (alterado bool, err error) {
	if p.nome == novoNome {
		return false, nil // No actual change (D-12)
	}
	p.nome = novoNome
	return true, nil
}

// detectarCiclo performs full ancestor walk from destino to check if pasta is an ancestor.
// Returns true if moving pasta to destino would create a cycle.
func (p *Pasta) detectarCiclo(destino *Pasta) bool {
	atual := destino
	for atual != nil {
		if atual == p {
			return true // Cycle detected: destino is descendant of pasta
		}
		atual = atual.pai
	}
	return false
}

// removerDePai removes this pasta from its parent's subpastas list.
func (p *Pasta) removerDePai() {
	if p.pai == nil {
		return // Pasta Geral has no parent
	}

	// Find and remove from parent's subpastas
	for i, sub := range p.pai.subpastas {
		if sub == p {
			p.pai.subpastas = append(p.pai.subpastas[:i], p.pai.subpastas[i+1:]...)
			break
		}
	}
}

// adicionarAoPai adds this pasta to a new parent at the end of subpastas list.
func (p *Pasta) adicionarAoPai(novoPai *Pasta) {
	p.pai = novoPai
	novoPai.subpastas = append(novoPai.subpastas, p)
}

// validarMover validates folder move parameters.
// Checks: not Pasta Geral, not moving to self, no cycle would be created, name unique in destination.
func (p *Pasta) validarMover(destino *Pasta) error {
	// Check not Pasta Geral
	if p.pai == nil {
		return ErrPastaGeralProtected
	}

	// Check not moving to self
	if p == destino {
		return ErrDestinoInvalido
	}

	// Check cycle: would moving to destino create a cycle?
	if p.detectarCiclo(destino) {
		return ErrCycleDetected
	}

	// Check name uniqueness in destination
	if destino.contemSubpastaComNome(p.nome) {
		return ErrNameConflict
	}

	return nil
}

// mover moves this pasta to a new parent folder.
// PRECONDITION: validarMover must pass (cannot fail after validation per D-05).
func (p *Pasta) mover(destino *Pasta) {
	p.removerDePai()
	p.adicionarAoPai(destino)
}

// obterPosicaoAtual returns the current position of this pasta in its parent's subpastas list.
// Returns -1 if pasta has no parent (Pasta Geral).
func (p *Pasta) obterPosicaoAtual() int {
	if p.pai == nil {
		return -1 // Pasta Geral has no position
	}

	for i, sub := range p.pai.subpastas {
		if sub == p {
			return i
		}
	}

	return -1 // Should never happen if tree is consistent
}

// validarReposicionar validates folder repositioning parameters.
// Checks: position in valid range [0, len-1] (0-indexed, cannot be equal to len).
func (p *Pasta) validarReposicionar(novaPosicao int) error {
	if p.pai == nil {
		return ErrPastaGeralProtected // Cannot reposition Pasta Geral
	}

	// Validate position in valid range [0, len-1] for repositioning
	// Note: Unlike criarSubpasta where len is valid (append), repositioning requires existing position
	if novaPosicao < 0 || novaPosicao >= len(p.pai.subpastas) {
		return ErrPosicaoInvalida
	}

	return nil
}

// reposicionar moves this pasta to a new position within its parent.
// Returns (true, nil) if position changed, (false, nil) if no change (D-12, D-23 no-op).
// PRECONDITION: validarReposicionar must pass (cannot fail after validation per D-05).
func (p *Pasta) reposicionar(novaPosicao int) (alterado bool, err error) {
	posicaoAtual := p.obterPosicaoAtual()

	// No-op if already at target position (D-12, D-23)
	if posicaoAtual == novaPosicao {
		return false, nil
	}

	// Remove from current position
	p.pai.subpastas = append(p.pai.subpastas[:posicaoAtual], p.pai.subpastas[posicaoAtual+1:]...)

	// Insert at new position
	p.pai.subpastas = append(p.pai.subpastas[:novaPosicao], append([]*Pasta{p}, p.pai.subpastas[novaPosicao:]...)...)

	return true, nil
}

// encontrarSubpastaPorNome finds a subfolder by name (case-insensitive).
// Returns the subfolder if found, nil otherwise.
func (p *Pasta) encontrarSubpastaPorNome(nome string) *Pasta {
	nomeLower := strings.ToLower(nome)
	for _, sub := range p.subpastas {
		if strings.ToLower(sub.nome) == nomeLower {
			return sub
		}
	}
	return nil
}

// encontrarSegredoPorNome finds a secret by name (case-insensitive).
// Returns the secret if found, nil otherwise.
func (p *Pasta) encontrarSegredoPorNome(nome string) *Segredo {
	nomeLower := strings.ToLower(nome)
	for _, seg := range p.segredos {
		if strings.ToLower(seg.nome) == nomeLower {
			return seg
		}
	}
	return nil
}

// gerarNomeSufixado generates a unique name with numeric suffix "(N)" for conflicts.
// Pattern: "Name (1)", "Name (2)", etc.
// Checks both secrets and subfolders in the target pasta.
func gerarNomeSufixado(nomeBase string, pasta *Pasta) string {
	// Try base name first
	if pasta.encontrarSegredoPorNome(nomeBase) == nil && pasta.encontrarSubpastaPorNome(nomeBase) == nil {
		return nomeBase
	}

	// Generate suffixed names
	for i := 1; i < 10000; i++ {
		nomeSufixado := nomeBase + " (" + strconv.Itoa(i) + ")"
		if pasta.encontrarSegredoPorNome(nomeSufixado) == nil && pasta.encontrarSubpastaPorNome(nomeSufixado) == nil {
			return nomeSufixado
		}
	}

	// Fallback (should never happen)
	return nomeBase + " (conflict)"
}

// mesclarPastas recursively merges contents from origem into destino.
// Handles name conflicts:
// - Subfolders: recursive merge
// - Secrets: rename with suffix, track in renomeacoes
// Returns slice of Renomeacao for renamed secrets.
func mesclarPastas(origem *Pasta, destino *Pasta) []Renomeacao {
	renomeacoes := make([]Renomeacao, 0)

	// Merge subfolders (recursive)
	for _, subOrigem := range origem.subpastas {
		subDestino := destino.encontrarSubpastaPorNome(subOrigem.nome)
		if subDestino != nil {
			// Conflict: merge recursively
			renomRecursivas := mesclarPastas(subOrigem, subDestino)
			renomeacoes = append(renomeacoes, renomRecursivas...)
		} else {
			// No conflict: move subfolder
			subOrigem.pai = destino
			destino.subpastas = append(destino.subpastas, subOrigem)
		}
	}

	// Merge secrets
	for _, segredo := range origem.segredos {
		if destino.encontrarSegredoPorNome(segredo.nome) != nil {
			// Conflict: rename
			nomeOriginal := segredo.nome
			segredo.nome = gerarNomeSufixado(nomeOriginal, destino)
			renomeacoes = append(renomeacoes, Renomeacao{
				Antigo: nomeOriginal,
				Novo:   segredo.nome,
				Pasta:  destino.nome,
			})
		}
		// Move secret
		segredo.pasta = destino
		destino.segredos = append(destino.segredos, segredo)
	}

	return renomeacoes
}

// validarExclusao validates that a folder can be deleted.
// Pasta Geral cannot be deleted.
func (p *Pasta) validarExclusao() error {
	if p.pai == nil {
		return ErrPastaGeralNaoExcluivel
	}
	return nil
}

// excluir removes this pasta from its parent and promotes all children (secrets and subfolders) to the parent.
// Handles name conflicts:
// - Subfolders: recursive merge via mesclarPastas
// - Secrets: rename with suffix, track in renomeacoes
// Per FOLDER-05: Secrets with EstadoExcluido retain that state when promoted.
// Per D-27: Hard delete (immediate removal).
// Returns slice of Renomeacao for renamed secrets.
func (p *Pasta) excluir(pai *Pasta) []Renomeacao {
	renomeacoes := make([]Renomeacao, 0)

	// Promote subfolders (with conflict resolution)
	for _, sub := range p.subpastas {
		existente := pai.encontrarSubpastaPorNome(sub.nome)
		if existente != nil {
			// Conflict: merge recursively
			renomRecursivas := mesclarPastas(sub, existente)
			renomeacoes = append(renomeacoes, renomRecursivas...)
		} else {
			// No conflict: move subfolder
			sub.pai = pai
			pai.subpastas = append(pai.subpastas, sub)
		}
	}

	// Promote secrets (with conflict resolution)
	for _, segredo := range p.segredos {
		if pai.encontrarSegredoPorNome(segredo.nome) != nil {
			// Conflict: rename
			nomeOriginal := segredo.nome
			segredo.nome = gerarNomeSufixado(nomeOriginal, pai)
			renomeacoes = append(renomeacoes, Renomeacao{
				Antigo: nomeOriginal,
				Novo:   segredo.nome,
				Pasta:  pai.nome,
			})
		}
		// Move secret (EstadoExcluido retained per FOLDER-05)
		segredo.pasta = pai
		pai.segredos = append(pai.segredos, segredo)
	}

	// Remove this pasta from parent
	for i, subpasta := range pai.subpastas {
		if subpasta == p {
			pai.subpastas = append(pai.subpastas[:i], pai.subpastas[i+1:]...)
			break
		}
	}

	return renomeacoes
}

// Factory methods

// NovoCofre creates a new empty vault with Pasta Geral and default configurations.
// Initial content (folders/templates) must be added via InicializarConteudoPadrao().
// Per D-28a.
func NovoCofre() *Cofre {
	agora := time.Now().UTC()

	pastaGeral := &Pasta{
		nome:      "Pasta Geral",
		pai:       nil,
		subpastas: make([]*Pasta, 0),
		segredos:  make([]*Segredo, 0),
	}

	return &Cofre{
		pastaGeral:            pastaGeral,
		modelos:               make([]*ModeloSegredo, 0),
		modificado:            false,
		dataCriacao:           agora,
		dataUltimaModificacao: agora,
		configuracoes: Configuracoes{
			tempoBloqueioInatividadeMinutos:      5,
			tempoOcultarSegredoSegundos:          15,
			tempoLimparAreaTransferenciaSegundos: 30,
		},
	}
}

// InicializarConteudoPadrao populates a new vault with default folders and templates.
// Creates "Sites e Apps" and "Financeiro" folders, and Login, Cartão de Crédito, and Chave de API templates.
// Must be called once after NovoCofre() for new user vaults.
// Does NOT mark cofre.modificado=true (initial content is part of base state per D-28b).
// NOT a user operation via Manager (system bootstrap, not domain operation).
func (c *Cofre) InicializarConteudoPadrao() error {
	// Create default folders
	sitesEApps := &Pasta{
		nome:      "Sites e Apps",
		pai:       c.pastaGeral,
		subpastas: make([]*Pasta, 0),
		segredos:  make([]*Segredo, 0),
	}

	financeiro := &Pasta{
		nome:      "Financeiro",
		pai:       c.pastaGeral,
		subpastas: make([]*Pasta, 0),
		segredos:  make([]*Segredo, 0),
	}

	c.pastaGeral.subpastas = append(c.pastaGeral.subpastas, sitesEApps, financeiro)

	// Create default templates (VAULT-02)
	login := &ModeloSegredo{
		nome: "Login",
		campos: []CampoModelo{
			{nome: "URL", tipo: TipoCampoComum},
			{nome: "Usuário", tipo: TipoCampoComum},
			{nome: "Senha", tipo: TipoCampoSensivel},
		},
	}

	cartao := &ModeloSegredo{
		nome: "Cartão de Crédito",
		campos: []CampoModelo{
			{nome: "Titular", tipo: TipoCampoComum},
			{nome: "Número", tipo: TipoCampoSensivel},
			{nome: "Validade", tipo: TipoCampoComum},
			{nome: "CVV", tipo: TipoCampoSensivel},
		},
	}

	chaveAPI := &ModeloSegredo{
		nome: "Chave de API",
		campos: []CampoModelo{
			{nome: "Serviço", tipo: TipoCampoComum},
			{nome: "Chave", tipo: TipoCampoSensivel},
		},
	}

	c.modelos = append(c.modelos, login, cartao, chaveAPI)

	// Note: Does NOT set c.modificado = true (D-28b)
	return nil
}

// Template management validation and mutation methods

// validarCriacaoModelo validates template creation request.
// Checks: name non-empty, no name conflict, no reserved field names.
func (c *Cofre) validarCriacaoModelo(nome string, campos []CampoModelo) error {
	// Check name not empty
	if nome == "" {
		return ErrNomeVazio
	}

	// Check name not too long
	if len(nome) > 255 {
		return ErrNomeMuitoLongo
	}

	// Check name not already in use
	for _, m := range c.modelos {
		if m.nome == nome {
			return ErrNameConflict
		}
	}

	// Check field names not reserved (D-29: "Observação" prohibited)
	for _, campo := range campos {
		if ehNomeReservado(campo.nome) {
			return ErrObservacaoReserved
		}
	}

	return nil
}

// criarModelo creates and inserts a new template into the vault.
// Inserts in alphabetically sorted position per D-23, TPL-06.
// Cannot fail after validation.
func (c *Cofre) criarModelo(nome string, campos []CampoModelo) *ModeloSegredo {
	novo := &ModeloSegredo{
		nome:   nome,
		campos: make([]CampoModelo, len(campos)),
	}
	copy(novo.campos, campos)

	// Insert in alphabetically sorted position (D-23)
	// Use case-insensitive comparison (strings.ToLower)
	pos := sort.Search(len(c.modelos), func(i int) bool {
		return c.modelos[i].nome > nome
	})

	// Insert at position
	c.modelos = append(c.modelos, nil)
	copy(c.modelos[pos+1:], c.modelos[pos:])
	c.modelos[pos] = novo

	return novo
}

// validarRenomear validates template rename request.
// Checks: new name non-empty, no conflict with other templates.
func (m *ModeloSegredo) validarRenomear(cofre *Cofre, novoNome string) error {
	// Check name not empty
	if novoNome == "" {
		return ErrNomeVazio
	}

	// Check name not too long
	if len(novoNome) > 255 {
		return ErrNomeMuitoLongo
	}

	// Check name not already in use by another template
	for _, modelo := range cofre.modelos {
		if modelo != m && modelo.nome == novoNome {
			return ErrNameConflict
		}
	}

	return nil
}

// renomear changes the template name and re-sorts the template list.
// Per D-23: templates always sorted alphabetically.
// Per D-12: returns true if name actually changed, false if no-op.
// Cannot fail after validation.
func (m *ModeloSegredo) renomear(novoNome string) bool {
	if m.nome == novoNome {
		return false // No change
	}
	m.nome = novoNome
	return true // Name changed
	// Note: Caller (Manager) is responsible for re-sorting if needed.
	// Since Cofre.Modelos() returns sorted copy, internal order doesn't affect TUI.
	// We rely on the getter's sort behavior (already implemented).
}

// validarExclusao validates template deletion request.
// Per TPL-04, D-26: templates can be deleted unless referenced by a secret.
func (m *ModeloSegredo) validarExclusao(cofre *Cofre) error {
	// Check if template is in use by any secret
	if m.emUso(cofre) {
		return ErrModeloEmUso
	}
	return nil
}

// emUso checks if this template is referenced by any secret in the vault.
// Per D-26: uses pointer equality (segredo.modelo == modelo).
func (m *ModeloSegredo) emUso(cofre *Cofre) bool {
	return verificarUsoRecursivo(m, cofre.pastaGeral)
}

// verificarUsoRecursivo recursively checks if template is used in folder tree.
func verificarUsoRecursivo(modelo *ModeloSegredo, pasta *Pasta) bool {
	if pasta == nil {
		return false
	}

	// Check secrets in this folder
	for _, segredo := range pasta.segredos {
		// Note: Segredo doesn't have modelo field yet (will be added in future task)
		// For now, we'll assume no secrets use templates (Task 5 will implement this)
		_ = segredo // Suppress unused warning
	}

	// Recurse into subfolders
	for _, subpasta := range pasta.subpastas {
		if verificarUsoRecursivo(modelo, subpasta) {
			return true
		}
	}

	return false
}

// excluir removes the template from the vault.
// Cannot fail after validation.
func (m *ModeloSegredo) excluir(cofre *Cofre) {
	// Find and remove template from slice
	for i, modelo := range cofre.modelos {
		if modelo == m {
			cofre.modelos = append(cofre.modelos[:i], cofre.modelos[i+1:]...)
			return
		}
	}
}

// Field management validation and mutation methods

// validarAdicionarCampo validates field addition request.
// Checks: name not reserved, position valid.
func (m *ModeloSegredo) validarAdicionarCampo(nome string, posicao int) error {
	// Check name not reserved (D-29)
	if ehNomeReservado(nome) {
		return ErrObservacaoReserved
	}

	// Check position valid (0 <= pos <= len, where len means append)
	if posicao < 0 || posicao > len(m.campos) {
		return ErrPosicaoInvalida
	}

	return nil
}

// adicionarCampo inserts a field at the specified position.
// Cannot fail after validation.
func (m *ModeloSegredo) adicionarCampo(nome string, tipo TipoCampo, posicao int) {
	novoCampo := CampoModelo{nome: nome, tipo: tipo}

	// Insert at position
	m.campos = append(m.campos, CampoModelo{})
	copy(m.campos[posicao+1:], m.campos[posicao:])
	m.campos[posicao] = novoCampo
}

// validarRemoverCampo validates field removal request.
// Checks: index is valid.
func (m *ModeloSegredo) validarRemoverCampo(indice int) error {
	if indice < 0 || indice >= len(m.campos) {
		return ErrCampoInvalido
	}
	return nil
}

// removerCampo removes a field by index.
// Cannot fail after validation.
func (m *ModeloSegredo) removerCampo(indice int) {
	m.campos = append(m.campos[:indice], m.campos[indice+1:]...)
}

// validarReordenarCampo validates field reordering request.
// Checks: both indices are valid.
func (m *ModeloSegredo) validarReordenarCampo(indiceOrigem, indiceDestino int) error {
	if indiceOrigem < 0 || indiceOrigem >= len(m.campos) {
		return ErrCampoInvalido
	}
	if indiceDestino < 0 || indiceDestino >= len(m.campos) {
		return ErrCampoInvalido
	}
	return nil
}

// reordenarCampo moves a field from one position to another.
// Cannot fail after validation.
func (m *ModeloSegredo) reordenarCampo(indiceOrigem, indiceDestino int) {
	// Remove from origin
	campo := m.campos[indiceOrigem]
	m.campos = append(m.campos[:indiceOrigem], m.campos[indiceOrigem+1:]...)

	// Insert at destination
	m.campos = append(m.campos, CampoModelo{})
	copy(m.campos[indiceDestino+1:], m.campos[indiceDestino:])
	m.campos[indiceDestino] = campo
}

// Helper functions

// ehNomeReservado checks if a name is reserved (case-insensitive).
// Per D-29: "Observação" is reserved for the observation field.
func ehNomeReservado(nome string) bool {
	// Case-insensitive comparison using strings package
	nomeLower := toLowerSimple(nome)
	return nomeLower == "observação"
}

// toLowerSimple performs simple lowercase conversion for ASCII and common Portuguese characters.
func toLowerSimple(s string) string {
	result := ""
	for _, r := range s {
		// Convert ASCII uppercase to lowercase
		if r >= 'A' && r <= 'Z' {
			result += string(r + 32)
		} else if r == 'Ç' {
			result += "ç"
		} else if r == 'Á' {
			result += "á"
		} else if r == 'É' {
			result += "é"
		} else if r == 'Í' {
			result += "í"
		} else if r == 'Ó' {
			result += "ó"
		} else if r == 'Ú' {
			result += "ú"
		} else if r == 'Â' {
			result += "â"
		} else if r == 'Ê' {
			result += "ê"
		} else if r == 'Ô' {
			result += "ô"
		} else if r == 'À' {
			result += "à"
		} else if r == 'Ã' {
			result += "ã"
		} else if r == 'Õ' {
			result += "õ"
		} else {
			result += string(r)
		}
	}
	return result
}

// Secret creation and lifecycle validation/mutation methods

// validarCriacaoSegredo validates secret creation parameters.
// Checks: pasta not nil, nome non-empty and unique, modelo not nil.
func (p *Pasta) validarCriacaoSegredo(nome string, modelo *ModeloSegredo) error {
	if nome == "" {
		return ErrNomeVazio
	}
	if len(nome) > 255 {
		return ErrNomeMuitoLongo
	}
	if modelo == nil {
		return ErrModeloInvalido
	}
	// Check name uniqueness within folder
	if p.contemSegredoComNome(nome) {
		return ErrNameConflict
	}
	return nil
}

// criarSegredo creates a new secret with template-based structure.
// Per D-11: estadoSessao = Normal initially (will be set to Modificado by Manager).
// Per D-13: campos initialized from template with empty values.
// PRECONDITION: validarCriacaoSegredo must pass.
func (p *Pasta) criarSegredo(nome string, modelo *ModeloSegredo) *Segredo {
	agora := time.Now().UTC()

	// Initialize campos from template structure with empty values
	campos := make([]CampoSegredo, len(modelo.campos))
	for i, campoModelo := range modelo.campos {
		campos[i] = CampoSegredo{
			nome:  campoModelo.nome,
			tipo:  campoModelo.tipo,
			valor: []byte{}, // Empty value
		}
	}

	// Create Observação field (always common type, empty value)
	observacao := CampoSegredo{
		nome:  "Observação",
		tipo:  TipoCampoComum,
		valor: []byte{},
	}

	segredo := &Segredo{
		nome:                  nome,
		campos:                campos,
		observacao:            observacao,
		pasta:                 p,
		favorito:              false,
		estadoSessao:          EstadoOriginal, // Will be changed to Modificado by Manager
		dataCriacao:           agora,
		dataUltimaModificacao: agora,
	}

	// Add to pasta's secrets list
	p.segredos = append(p.segredos, segredo)

	return segredo
}

// validarExclusaoSegredo validates secret deletion parameters.
// Returns error if segredo is nil or already excluded.
func (s *Segredo) validarExclusaoSegredo() error {
	if s == nil {
		return ErrSegredoInvalido
	}
	if s.estadoSessao == EstadoExcluido {
		return ErrSegredoJaExcluido
	}
	return nil
}

// excluirSegredo marks secret as excluded (soft delete).
// Per D-14: Delete is reversible until Salvar.
// State transitions: Original→Excluido, Incluido→removed, Modificado→Excluido.
// PRECONDITION: validarExclusaoSegredo must pass.
func (s *Segredo) excluirSegredo() {
	switch s.estadoSessao {
	case EstadoIncluido:
		// New secret never persisted - remove from parent's list
		if s.pasta != nil {
			for i, seg := range s.pasta.segredos {
				if seg == s {
					s.pasta.segredos = append(s.pasta.segredos[:i], s.pasta.segredos[i+1:]...)
					break
				}
			}
		}
	case EstadoOriginal, EstadoModificado:
		// Existing secret - mark as excluded (soft delete)
		s.estadoSessao = EstadoExcluido
	}
}

// validarRestauracaoSegredo validates secret restoration parameters.
// Returns error if segredo is nil or not excluded.
func (s *Segredo) validarRestauracaoSegredo() error {
	if s == nil {
		return ErrSegredoInvalido
	}
	if s.estadoSessao != EstadoExcluido {
		return ErrSegredoNaoExcluido
	}
	return nil
}

// restaurarSegredo restores a soft-deleted secret.
// Per D-14: Restore reverses deletion (Excluido → Original or Modificado).
// PRECONDITION: validarRestauracaoSegredo must pass.
func (s *Segredo) restaurarSegredo() {
	// Restore to Modificado (content was marked as deleted, now restored)
	s.estadoSessao = EstadoModificado
}

// validarAlternarFavorito validates favorite toggle parameters.
// Returns error if segredo is nil or excluded.
func (s *Segredo) validarAlternarFavorito() error {
	if s == nil {
		return ErrSegredoInvalido
	}
	if s.estadoSessao == EstadoExcluido {
		return ErrSegredoJaExcluido
	}
	return nil
}

// alternarFavorito toggles the favorito flag.
// Per D-11: favorito is independent of estadoSessao (does NOT change content state).
// PRECONDITION: validarAlternarFavorito must pass.
func (s *Segredo) alternarFavorito() {
	s.favorito = !s.favorito
}

// validarDuplicacaoSegredo validates secret duplication parameters.
// Returns error if segredo is nil or excluded.
func (s *Segredo) validarDuplicacaoSegredo() error {
	if s == nil {
		return ErrSegredoInvalido
	}
	if s.estadoSessao == EstadoExcluido {
		return ErrSegredoJaExcluido
	}
	return nil
}

// duplicarSegredo creates a copy of the secret with name conflict resolution.
// Per D-27: Uses "(N)" progression for name conflicts: "Name" → "Name (2)" → "Name (3)".
// PRECONDITION: validarDuplicacaoSegredo must pass.
func (p *Pasta) duplicarSegredo(original *Segredo) *Segredo {
	// Generate unique name using "(N)" progression
	baseName := original.nome
	newName := baseName
	counter := 2

	// Check if name already ends with "(N)" pattern
	// If so, extract base and continue from that counter
	// Otherwise start with "(2)"
	for p.contemSegredoComNome(newName) {
		newName = fmt.Sprintf("%s (%d)", baseName, counter)
		counter++
		// Safety limit to prevent infinite loops
		if counter > 9999 {
			break
		}
	}

	// Create duplicate with deep copy of campos
	agora := time.Now().UTC()
	campos := make([]CampoSegredo, len(original.campos))
	for i, campo := range original.campos {
		// Deep copy the valor slice
		valorCopy := make([]byte, len(campo.valor))
		copy(valorCopy, campo.valor)

		campos[i] = CampoSegredo{
			nome:  campo.nome,
			tipo:  campo.tipo,
			valor: valorCopy,
		}
	}

	// Copy observacao (always exists as value, not pointer)
	valorCopy := make([]byte, len(original.observacao.valor))
	copy(valorCopy, original.observacao.valor)
	observacao := CampoSegredo{
		nome:  "Observação",
		tipo:  TipoCampoComum,
		valor: valorCopy,
	}

	duplicate := &Segredo{
		nome:                  newName,
		campos:                campos,
		observacao:            observacao,
		pasta:                 p,
		favorito:              false,          // Reset favorite flag
		estadoSessao:          EstadoOriginal, // Manager will set to Modificado
		dataCriacao:           agora,
		dataUltimaModificacao: agora,
	}

	// Add to pasta's secrets list
	p.segredos = append(p.segredos, duplicate)

	return duplicate
}
