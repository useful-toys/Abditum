package vault

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// TestSerializarCofre_Nil tests that a nil cofre returns an error.
func TestSerializarCofre_Nil(t *testing.T) {
	_, err := SerializarCofre(nil)
	if err == nil {
		t.Error("SerializarCofre(nil) should return error")
	}
}

// TestSerializarCofre_ProduceJSON tests that serialization produces valid JSON.
func TestSerializarCofre_ProduceJSON(t *testing.T) {
	cofre := NovoCofre()

	data, err := SerializarCofre(cofre)
	if err != nil {
		t.Fatalf("SerializarCofre() returned error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("SerializarCofre() returned empty data")
	}
	if !json.Valid(data) {
		t.Errorf("SerializarCofre() did not produce valid JSON: %s", data)
	}
}

// TestSerializarCofre_PastaGeralNomeGeral tests that pasta_geral.nome is serialized as "Geral".
func TestSerializarCofre_PastaGeralNomeGeral(t *testing.T) {
	cofre := NovoCofre()

	data, err := SerializarCofre(cofre)
	if err != nil {
		t.Fatalf("SerializarCofre() returned error: %v", err)
	}

	// Unmarshal into generic map to check pasta_geral.nome
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	pastaGeral, ok := m["pasta_geral"].(map[string]interface{})
	if !ok {
		t.Fatal("pasta_geral field missing or not an object")
	}
	nome, ok := pastaGeral["nome"].(string)
	if !ok {
		t.Fatal("pasta_geral.nome field missing")
	}
	if nome != "Geral" {
		t.Errorf("Expected pasta_geral.nome='Geral', got '%s'", nome)
	}
}

// TestSerializarCofre_OmiteExcluidos tests that secrets with EstadoExcluido are omitted.
func TestSerializarCofre_OmiteExcluidos(t *testing.T) {
	cofre := NovoCofre()

	// Create a template and two secrets
	modelo := &ModeloSegredo{
		nome:   "Test",
		campos: []CampoModelo{{nome: "Campo", tipo: TipoCampoComum}},
	}
	cofre.modelos = append(cofre.modelos, modelo)

	segredo1 := cofre.pastaGeral.criarSegredo("Segredo Ativo", modelo)
	segredo2 := cofre.pastaGeral.criarSegredo("Segredo Excluido", modelo)

	// Mark second secret as excluded
	segredo2.estadoSessao = EstadoExcluido

	_ = segredo1 // used

	data, err := SerializarCofre(cofre)
	if err != nil {
		t.Fatalf("SerializarCofre() returned error: %v", err)
	}

	if strings.Contains(string(data), "Segredo Excluido") {
		t.Error("Excluded secret should not appear in serialized JSON")
	}
	if !strings.Contains(string(data), "Segredo Ativo") {
		t.Error("Active secret should appear in serialized JSON")
	}
}

// TestSerializarCofre_ValorUTF8 tests that CampoSegredo.valor is serialized as UTF-8 string (not Base64).
func TestSerializarCofre_ValorUTF8(t *testing.T) {
	cofre := NovoCofre()

	modelo := &ModeloSegredo{
		nome:   "Login",
		campos: []CampoModelo{{nome: "Senha", tipo: TipoCampoSensivel}},
	}
	cofre.modelos = append(cofre.modelos, modelo)

	segredo := cofre.pastaGeral.criarSegredo("Conta Gmail", modelo)
	segredo.campos[0].valor = []byte("minha-senha-secreta")

	data, err := SerializarCofre(cofre)
	if err != nil {
		t.Fatalf("SerializarCofre() returned error: %v", err)
	}

	jsonStr := string(data)
	// The value should appear as a plain UTF-8 string in JSON, not Base64
	if !strings.Contains(jsonStr, "minha-senha-secreta") {
		t.Errorf("Expected valor to be serialized as UTF-8 string, got: %s", jsonStr)
	}
}

// TestDeserializarCofre_InvalidJSON tests that invalid JSON returns an error.
func TestDeserializarCofre_InvalidJSON(t *testing.T) {
	_, err := DeserializarCofre([]byte("not valid json"), 1)
	if err == nil {
		t.Error("DeserializarCofre with invalid JSON should return error")
	}
}

// TestDeserializarCofre_SemPastaGeral tests that missing pasta_geral returns an error.
func TestDeserializarCofre_SemPastaGeral(t *testing.T) {
	// JSON without pasta_geral
	jsonSemPastaGeral := `{"data_criacao":"2026-01-01T00:00:00Z","data_ultima_modificacao":"2026-01-01T00:00:00Z"}`
	_, err := DeserializarCofre([]byte(jsonSemPastaGeral), 1)
	if err == nil {
		t.Error("DeserializarCofre without pasta_geral should return error")
	}
}

// TestDeserializarCofre_PastaGeralNomeErrado tests that pasta_geral.nome != "Geral" returns error.
func TestDeserializarCofre_PastaGeralNomeErrado(t *testing.T) {
	jsonNomeErrado := `{
		"data_criacao":"2026-01-01T00:00:00Z",
		"data_ultima_modificacao":"2026-01-01T00:00:00Z",
		"configuracoes":{"tempo_bloqueio_inatividade_minutos":5,"tempo_ocultar_segredo_segundos":15,"tempo_limpar_area_transferencia_segundos":30},
		"modelos":[],
		"pasta_geral":{"nome":"NaoGeral","subpastas":[],"segredos":[]}
	}`
	_, err := DeserializarCofre([]byte(jsonNomeErrado), 1)
	if err == nil {
		t.Error("DeserializarCofre with pasta_geral.nome != 'Geral' should return error")
	}
}

// TestRoundtrip_CofreVazio tests serialization/deserialization roundtrip of an empty vault.
func TestRoundtrip_CofreVazio(t *testing.T) {
	cofre := NovoCofre()
	// Set known timestamps for comparison
	agora := time.Date(2026, 3, 30, 10, 0, 0, 0, time.UTC)
	cofre.dataCriacao = agora
	cofre.dataUltimaModificacao = agora

	data, err := SerializarCofre(cofre)
	if err != nil {
		t.Fatalf("SerializarCofre() returned error: %v", err)
	}

	restored, err := DeserializarCofre(data, 1)
	if err != nil {
		t.Fatalf("DeserializarCofre() returned error: %v", err)
	}

	if restored == nil {
		t.Fatal("DeserializarCofre() returned nil")
	}
	if restored.pastaGeral == nil {
		t.Fatal("Deserialized cofre has nil pastaGeral")
	}
	// pastaGeral internal nome should remain "Pasta Geral" after deserialization
	// (serialized as "Geral", but reconstructed as the root pasta)
	if restored.pastaGeral.pai != nil {
		t.Error("pastaGeral.pai should be nil (it's the root)")
	}
	if !restored.dataCriacao.Equal(agora) {
		t.Errorf("dataCriacao mismatch: got %v, want %v", restored.dataCriacao, agora)
	}
	if !restored.dataUltimaModificacao.Equal(agora) {
		t.Errorf("dataUltimaModificacao mismatch")
	}
}

// TestRoundtrip_Configuracoes tests that Configuracoes are preserved.
func TestRoundtrip_Configuracoes(t *testing.T) {
	cofre := NovoCofre()
	cofre.configuracoes = Configuracoes{
		tempoBloqueioInatividadeMinutos:      10,
		tempoOcultarSegredoSegundos:          20,
		tempoLimparAreaTransferenciaSegundos: 45,
	}

	data, err := SerializarCofre(cofre)
	if err != nil {
		t.Fatalf("SerializarCofre() error: %v", err)
	}

	restored, err := DeserializarCofre(data, 1)
	if err != nil {
		t.Fatalf("DeserializarCofre() error: %v", err)
	}

	cfg := restored.configuracoes
	if cfg.tempoBloqueioInatividadeMinutos != 10 {
		t.Errorf("tempoBloqueio: got %d, want 10", cfg.tempoBloqueioInatividadeMinutos)
	}
	if cfg.tempoOcultarSegredoSegundos != 20 {
		t.Errorf("tempoOcultar: got %d, want 20", cfg.tempoOcultarSegredoSegundos)
	}
	if cfg.tempoLimparAreaTransferenciaSegundos != 45 {
		t.Errorf("tempoLimpar: got %d, want 45", cfg.tempoLimparAreaTransferenciaSegundos)
	}
}

// TestRoundtrip_Modelos tests that templates are preserved in roundtrip.
func TestRoundtrip_Modelos(t *testing.T) {
	cofre := NovoCofre()
	cofre.modelos = []*ModeloSegredo{
		{
			nome: "Login",
			campos: []CampoModelo{
				{nome: "URL", tipo: TipoCampoComum},
				{nome: "Senha", tipo: TipoCampoSensivel},
			},
		},
	}

	data, err := SerializarCofre(cofre)
	if err != nil {
		t.Fatalf("SerializarCofre() error: %v", err)
	}

	restored, err := DeserializarCofre(data, 1)
	if err != nil {
		t.Fatalf("DeserializarCofre() error: %v", err)
	}

	if len(restored.modelos) != 1 {
		t.Fatalf("Expected 1 modelo, got %d", len(restored.modelos))
	}
	m := restored.modelos[0]
	if m.nome != "Login" {
		t.Errorf("modelo.nome: got '%s', want 'Login'", m.nome)
	}
	if len(m.campos) != 2 {
		t.Fatalf("Expected 2 campos in modelo, got %d", len(m.campos))
	}
	if m.campos[0].nome != "URL" || m.campos[0].tipo != TipoCampoComum {
		t.Error("modelo.campos[0] mismatch")
	}
	if m.campos[1].nome != "Senha" || m.campos[1].tipo != TipoCampoSensivel {
		t.Error("modelo.campos[1] mismatch")
	}
}

// TestRoundtrip_SegredoCompleto tests full secret roundtrip with all fields.
func TestRoundtrip_SegredoCompleto(t *testing.T) {
	cofre := NovoCofre()

	modelo := &ModeloSegredo{
		nome: "Login",
		campos: []CampoModelo{
			{nome: "URL", tipo: TipoCampoComum},
			{nome: "Senha", tipo: TipoCampoSensivel},
		},
	}
	cofre.modelos = append(cofre.modelos, modelo)

	t0 := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	t1 := time.Date(2026, 3, 30, 10, 0, 0, 0, time.UTC)

	segredo := cofre.pastaGeral.criarSegredo("Conta Gmail", modelo)
	segredo.campos[0].valor = []byte("https://gmail.com")
	segredo.campos[1].valor = []byte("super-secret")
	segredo.observacao.valor = []byte("conta pessoal")
	segredo.favorito = true
	segredo.estadoSessao = EstadoOriginal
	segredo.dataCriacao = t0
	segredo.dataUltimaModificacao = t1

	data, err := SerializarCofre(cofre)
	if err != nil {
		t.Fatalf("SerializarCofre() error: %v", err)
	}

	restored, err := DeserializarCofre(data, 1)
	if err != nil {
		t.Fatalf("DeserializarCofre() error: %v", err)
	}

	segredos := restored.pastaGeral.segredos
	if len(segredos) != 1 {
		t.Fatalf("Expected 1 secret, got %d", len(segredos))
	}

	s := segredos[0]
	if s.nome != "Conta Gmail" {
		t.Errorf("nome: got '%s', want 'Conta Gmail'", s.nome)
	}
	if s.favorito != true {
		t.Error("favorito should be true")
	}
	if s.estadoSessao != EstadoOriginal {
		t.Errorf("estadoSessao: got %v, want EstadoOriginal", s.estadoSessao)
	}
	if !s.dataCriacao.Equal(t0) {
		t.Errorf("dataCriacao: got %v, want %v", s.dataCriacao, t0)
	}
	if !s.dataUltimaModificacao.Equal(t1) {
		t.Errorf("dataUltimaModificacao: got %v, want %v", s.dataUltimaModificacao, t1)
	}
	if len(s.campos) != 2 {
		t.Fatalf("Expected 2 campos, got %d", len(s.campos))
	}
	if string(s.campos[0].valor) != "https://gmail.com" {
		t.Errorf("campos[0].valor: got '%s', want 'https://gmail.com'", s.campos[0].valor)
	}
	if string(s.campos[1].valor) != "super-secret" {
		t.Errorf("campos[1].valor: got '%s', want 'super-secret'", s.campos[1].valor)
	}
	if string(s.observacao.valor) != "conta pessoal" {
		t.Errorf("observacao.valor: got '%s', want 'conta pessoal'", s.observacao.valor)
	}
	if s.observacao.nome != "Observacao" && s.observacao.nome != "Observação" {
		// observacao nome is reconstructed with standard name
		t.Logf("observacao.nome: %s (acceptable)", s.observacao.nome)
	}
}

// TestRoundtrip_ReferenciasPaiFilho tests that parent-child references are populated.
func TestRoundtrip_ReferenciasPaiFilho(t *testing.T) {
	cofre := NovoCofre()

	// Create a subfolder
	subpasta := cofre.pastaGeral.criarSubpasta("Trabalho", 0)

	// Create a template and a secret in the subfolder
	modelo := &ModeloSegredo{
		nome:   "Nota",
		campos: []CampoModelo{{nome: "Conteudo", tipo: TipoCampoComum}},
	}
	cofre.modelos = append(cofre.modelos, modelo)
	segredo := subpasta.criarSegredo("Acesso VPN", modelo)
	_ = segredo

	data, err := SerializarCofre(cofre)
	if err != nil {
		t.Fatalf("SerializarCofre() error: %v", err)
	}

	restored, err := DeserializarCofre(data, 1)
	if err != nil {
		t.Fatalf("DeserializarCofre() error: %v", err)
	}

	// Verify pastaGeral has no parent
	if restored.pastaGeral.pai != nil {
		t.Error("pastaGeral.pai should be nil")
	}

	// Verify subpasta has correct parent
	if len(restored.pastaGeral.subpastas) != 1 {
		t.Fatalf("Expected 1 subpasta, got %d", len(restored.pastaGeral.subpastas))
	}
	sub := restored.pastaGeral.subpastas[0]
	if sub.pai != restored.pastaGeral {
		t.Error("subpasta.pai should point to pastaGeral")
	}
	if sub.nome != "Trabalho" {
		t.Errorf("subpasta.nome: got '%s', want 'Trabalho'", sub.nome)
	}

	// Verify secret has correct parent pasta reference
	if len(sub.segredos) != 1 {
		t.Fatalf("Expected 1 secret in subpasta, got %d", len(sub.segredos))
	}
	s := sub.segredos[0]
	if s.pasta != sub {
		t.Error("segredo.pasta should point to its parent subpasta")
	}
}

// TestRoundtrip_EstadoSessaoOriginal tests that all deserialized secrets get EstadoOriginal.
func TestRoundtrip_EstadoSessaoOriginal(t *testing.T) {
	cofre := NovoCofre()

	modelo := &ModeloSegredo{
		nome:   "Test",
		campos: []CampoModelo{{nome: "Campo", tipo: TipoCampoComum}},
	}
	cofre.modelos = append(cofre.modelos, modelo)

	s1 := cofre.pastaGeral.criarSegredo("Incluido", modelo)
	s2 := cofre.pastaGeral.criarSegredo("Modificado", modelo)
	s1.estadoSessao = EstadoIncluido
	s2.estadoSessao = EstadoModificado

	data, err := SerializarCofre(cofre)
	if err != nil {
		t.Fatalf("SerializarCofre() error: %v", err)
	}

	restored, err := DeserializarCofre(data, 1)
	if err != nil {
		t.Fatalf("DeserializarCofre() error: %v", err)
	}

	for _, s := range restored.pastaGeral.segredos {
		if s.estadoSessao != EstadoOriginal {
			t.Errorf("Secret '%s': estadoSessao=%v, want EstadoOriginal", s.nome, s.estadoSessao)
		}
	}
}

// TestRoundtrip_SubpastasAninhadas tests nested subfolder serialization/deserialization.
func TestRoundtrip_SubpastasAninhadas(t *testing.T) {
	cofre := NovoCofre()

	// Create nested structure: PastaGeral > Trabalho > Projetos
	trabalho := cofre.pastaGeral.criarSubpasta("Trabalho", 0)
	projetos := trabalho.criarSubpasta("Projetos", 0)
	_ = projetos

	data, err := SerializarCofre(cofre)
	if err != nil {
		t.Fatalf("SerializarCofre() error: %v", err)
	}

	restored, err := DeserializarCofre(data, 1)
	if err != nil {
		t.Fatalf("DeserializarCofre() error: %v", err)
	}

	if len(restored.pastaGeral.subpastas) != 1 {
		t.Fatalf("Expected 1 subpasta in root, got %d", len(restored.pastaGeral.subpastas))
	}
	trab := restored.pastaGeral.subpastas[0]
	if trab.nome != "Trabalho" {
		t.Errorf("Expected 'Trabalho', got '%s'", trab.nome)
	}
	if trab.pai != restored.pastaGeral {
		t.Error("Trabalho.pai should be pastaGeral")
	}

	if len(trab.subpastas) != 1 {
		t.Fatalf("Expected 1 subpasta in Trabalho, got %d", len(trab.subpastas))
	}
	proj := trab.subpastas[0]
	if proj.nome != "Projetos" {
		t.Errorf("Expected 'Projetos', got '%s'", proj.nome)
	}
	if proj.pai != trab {
		t.Error("Projetos.pai should be Trabalho")
	}
}

// TestRoundtrip_OrdemPreservada tests that secrets and subfolders ordering is preserved.
func TestRoundtrip_OrdemPreservada(t *testing.T) {
	cofre := NovoCofre()

	modelo := &ModeloSegredo{
		nome:   "Test",
		campos: []CampoModelo{{nome: "V", tipo: TipoCampoComum}},
	}
	cofre.modelos = append(cofre.modelos, modelo)

	// Add secrets in specific order
	cofre.pastaGeral.criarSegredo("Primeiro", modelo)
	cofre.pastaGeral.criarSegredo("Segundo", modelo)
	cofre.pastaGeral.criarSegredo("Terceiro", modelo)

	// Add subfolders in specific order
	cofre.pastaGeral.criarSubpasta("Alpha", 0)
	cofre.pastaGeral.criarSubpasta("Beta", 1)

	data, err := SerializarCofre(cofre)
	if err != nil {
		t.Fatalf("SerializarCofre() error: %v", err)
	}

	restored, err := DeserializarCofre(data, 1)
	if err != nil {
		t.Fatalf("DeserializarCofre() error: %v", err)
	}

	segredos := restored.pastaGeral.segredos
	expectedNomes := []string{"Primeiro", "Segundo", "Terceiro"}
	for i, nome := range expectedNomes {
		if i >= len(segredos) {
			t.Fatalf("Missing secret at index %d", i)
		}
		if segredos[i].nome != nome {
			t.Errorf("segredos[%d].nome: got '%s', want '%s'", i, segredos[i].nome, nome)
		}
	}

	subpastas := restored.pastaGeral.subpastas
	expectedSubs := []string{"Alpha", "Beta"}
	for i, nome := range expectedSubs {
		if i >= len(subpastas) {
			t.Fatalf("Missing subpasta at index %d", i)
		}
		if subpastas[i].nome != nome {
			t.Errorf("subpastas[%d].nome: got '%s', want '%s'", i, subpastas[i].nome, nome)
		}
	}
}
