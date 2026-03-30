// Package vault implements the domain layer for Abditum password manager.
// All entities use package-private fields (lowercase) for encapsulation.
// External access via exported getters returning defensive copies.
package vault

import (
	"sort"
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
