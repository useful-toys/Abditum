package vault

import (
	"encoding/json"
	"errors"
	"time"
)

// Sentinel errors for serialization failures.
var (
	ErrCofreNulo              = errors.New("cofre nao pode ser nulo")
	ErrPastaGeralAusente      = errors.New("pasta_geral ausente ou invalida no JSON")
	ErrPastaGeralNomeInvalido = errors.New("pasta_geral.nome deve ser 'Geral'")
)

// --------------------------------------------------------------------------
// Internal serialization structs
// These are private to this file and used only for JSON marshaling/unmarshaling.
// They mirror the domain entities but have exported field names (required by encoding/json).
// --------------------------------------------------------------------------

type configuracaoJSON struct {
	TempoBloqueioInatividadeMinutos      int `json:"tempo_bloqueio_inatividade_minutos"`
	TempoOcultarSegredoSegundos          int `json:"tempo_ocultar_segredo_segundos"`
	TempoLimparAreaTransferenciaSegundos int `json:"tempo_limpar_area_transferencia_segundos"`
}

type campoModeloJSON struct {
	Nome string    `json:"nome"`
	Tipo TipoCampo `json:"tipo"`
}

type modeloJSON struct {
	Nome   string            `json:"nome"`
	Campos []campoModeloJSON `json:"campos"`
}

type campoJSON struct {
	Nome  string    `json:"nome"`
	Tipo  TipoCampo `json:"tipo"`
	Valor string    `json:"valor"`
}

type segredoJSON struct {
	Nome                  string      `json:"nome"`
	Campos                []campoJSON `json:"campos"`
	Observacao            string      `json:"observacao"`
	Favorito              bool        `json:"favorito"`
	DataCriacao           time.Time   `json:"data_criacao"`
	DataUltimaModificacao time.Time   `json:"data_ultima_modificacao"`
}

type pastaJSON struct {
	Nome      string        `json:"nome"`
	Subpastas []pastaJSON   `json:"subpastas"`
	Segredos  []segredoJSON `json:"segredos"`
}

type cofreJSON struct {
	DataCriacao           time.Time        `json:"data_criacao"`
	DataUltimaModificacao time.Time        `json:"data_ultima_modificacao"`
	Configuracoes         configuracaoJSON `json:"configuracoes"`
	Modelos               []modeloJSON     `json:"modelos"`
	PastaGeral            *pastaJSON       `json:"pasta_geral"`
}

// --------------------------------------------------------------------------
// SerializarCofre converts a Cofre to JSON bytes.
//
// Secrets with estadoSessao == EstadoExcluido are omitted.
// CampoSegredo.valor is serialized as a UTF-8 string (not Base64).
// pasta_geral.nome is always written as "Geral" (canonical serialized name).
// --------------------------------------------------------------------------
func SerializarCofre(cofre *Cofre) ([]byte, error) {
	if cofre == nil {
		return nil, ErrCofreNulo
	}

	dto := cofreJSON{
		DataCriacao:           cofre.dataCriacao,
		DataUltimaModificacao: cofre.dataUltimaModificacao,
		Configuracoes: configuracaoJSON{
			TempoBloqueioInatividadeMinutos:      cofre.configuracoes.tempoBloqueioInatividadeMinutos,
			TempoOcultarSegredoSegundos:          cofre.configuracoes.tempoOcultarSegredoSegundos,
			TempoLimparAreaTransferenciaSegundos: cofre.configuracoes.tempoLimparAreaTransferenciaSegundos,
		},
		Modelos:    serializarModelos(cofre.modelos),
		PastaGeral: serializarPasta(cofre.pastaGeral, true),
	}

	return json.Marshal(dto)
}

// serializarModelos converts a slice of ModeloSegredo to serialization structs.
func serializarModelos(modelos []*ModeloSegredo) []modeloJSON {
	result := make([]modeloJSON, len(modelos))
	for i, m := range modelos {
		campos := make([]campoModeloJSON, len(m.campos))
		for j, c := range m.campos {
			campos[j] = campoModeloJSON{Nome: c.nome, Tipo: c.tipo}
		}
		result[i] = modeloJSON{Nome: m.nome, Campos: campos}
	}
	return result
}

// serializarPasta converts a Pasta (and its subtree) to a serialization struct.
// isRoot: when true, the nome is written as "Geral" (canonical name for pasta_geral).
func serializarPasta(pasta *Pasta, isRoot bool) *pastaJSON {
	if pasta == nil {
		return nil
	}

	nome := pasta.nome
	if isRoot {
		nome = "Geral"
	}

	// Serialize subpastas recursively
	subpastas := make([]pastaJSON, len(pasta.subpastas))
	for i, sub := range pasta.subpastas {
		p := serializarPasta(sub, false)
		subpastas[i] = *p
	}

	// Serialize segredos, omitting EstadoExcluido
	var segredos []segredoJSON
	for _, s := range pasta.segredos {
		if s.estadoSessao == EstadoExcluido {
			continue
		}
		campos := make([]campoJSON, len(s.campos))
		for i, c := range s.campos {
			campos[i] = campoJSON{
				Nome:  c.nome,
				Tipo:  c.tipo,
				Valor: string(c.valor), // UTF-8, not Base64
			}
		}
		segredos = append(segredos, segredoJSON{
			Nome:                  s.nome,
			Campos:                campos,
			Observacao:            string(s.observacao.valor), // plain string
			Favorito:              s.favorito,
			DataCriacao:           s.dataCriacao,
			DataUltimaModificacao: s.dataUltimaModificacao,
		})
	}
	if segredos == nil {
		segredos = []segredoJSON{}
	}

	return &pastaJSON{
		Nome:      nome,
		Subpastas: subpastas,
		Segredos:  segredos,
	}
}

// --------------------------------------------------------------------------
// DeserializarCofre converts JSON bytes back to a Cofre.
//
// Validates that pasta_geral exists and has nome == "Geral".
// Sets all secrets to estadoSessao = EstadoOriginal.
// Populates all parent-child references via popularReferencias.
// --------------------------------------------------------------------------
func DeserializarCofre(data []byte) (*Cofre, error) {
	var dto cofreJSON
	if err := json.Unmarshal(data, &dto); err != nil {
		return nil, err
	}

	if dto.PastaGeral == nil {
		return nil, ErrPastaGeralAusente
	}
	if dto.PastaGeral.Nome != "Geral" {
		return nil, ErrPastaGeralNomeInvalido
	}

	cofre := &Cofre{
		dataCriacao:           dto.DataCriacao,
		dataUltimaModificacao: dto.DataUltimaModificacao,
		configuracoes: Configuracoes{
			tempoBloqueioInatividadeMinutos:      dto.Configuracoes.TempoBloqueioInatividadeMinutos,
			tempoOcultarSegredoSegundos:          dto.Configuracoes.TempoOcultarSegredoSegundos,
			tempoLimparAreaTransferenciaSegundos: dto.Configuracoes.TempoLimparAreaTransferenciaSegundos,
		},
		modificado: false,
	}

	// Reconstruct modelos
	cofre.modelos = make([]*ModeloSegredo, len(dto.Modelos))
	for i, m := range dto.Modelos {
		campos := make([]CampoModelo, len(m.Campos))
		for j, c := range m.Campos {
			campos[j] = CampoModelo{nome: c.Nome, tipo: c.Tipo}
		}
		cofre.modelos[i] = &ModeloSegredo{nome: m.Nome, campos: campos}
	}

	// Reconstruct pasta tree
	cofre.pastaGeral = deserializarPasta(dto.PastaGeral)

	// Populate parent-child references
	popularReferencias(cofre.pastaGeral, nil)

	return cofre, nil
}

// deserializarPasta converts a pastaJSON to a Pasta entity (without parent refs).
func deserializarPasta(dto *pastaJSON) *Pasta {
	if dto == nil {
		return nil
	}

	pasta := &Pasta{
		nome:      dto.Nome,
		pai:       nil, // set by popularReferencias
		subpastas: make([]*Pasta, len(dto.Subpastas)),
		segredos:  make([]*Segredo, 0, len(dto.Segredos)),
	}

	for i, sub := range dto.Subpastas {
		subCopy := sub
		pasta.subpastas[i] = deserializarPasta(&subCopy)
	}

	for _, sdto := range dto.Segredos {
		campos := make([]CampoSegredo, len(sdto.Campos))
		for i, c := range sdto.Campos {
			campos[i] = CampoSegredo{
				nome:  c.Nome,
				tipo:  c.Tipo,
				valor: []byte(c.Valor), // UTF-8 string back to []byte
			}
		}
		observacao := CampoSegredo{
			nome:  "Observacao",
			tipo:  TipoCampoComum,
			valor: []byte(sdto.Observacao),
		}
		s := &Segredo{
			nome:                  sdto.Nome,
			campos:                campos,
			observacao:            observacao,
			pasta:                 nil, // set by popularReferencias
			favorito:              sdto.Favorito,
			estadoSessao:          EstadoOriginal, // always reset on load
			dataCriacao:           sdto.DataCriacao,
			dataUltimaModificacao: sdto.DataUltimaModificacao,
		}
		pasta.segredos = append(pasta.segredos, s)
	}

	return pasta
}

// popularReferencias recursively sets parent-child back-references.
// Sets pasta.pai = pai, and segredo.pasta = pasta for all secrets in the tree.
func popularReferencias(pasta *Pasta, pai *Pasta) {
	if pasta == nil {
		return
	}
	pasta.pai = pai
	for _, segredo := range pasta.segredos {
		segredo.pasta = pasta
	}
	for _, sub := range pasta.subpastas {
		popularReferencias(sub, pasta)
	}
}
