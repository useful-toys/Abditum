// Package vault implements the domain layer for Abditum password manager.
// All entities use package-private fields (lowercase) for encapsulation.
// External access via exported getters returning defensive copies.
package vault

import (
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
