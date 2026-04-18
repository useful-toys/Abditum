package vault

import (
	"errors"
	"testing"
)

// Mock repository for testing
type mockRepository struct {
	salvarCalled                 bool
	salvarError                  error
	detectarAlteracaoExternaResp bool
	detectarAlteracaoExternaErr  error
}

func (m *mockRepository) Salvar(cofre *Cofre) error {
	m.salvarCalled = true
	return m.salvarError
}

func (m *mockRepository) Carregar() (*Cofre, error) {
	return nil, errors.New("not implemented")
}

func (m *mockRepository) DetectarAlteracaoExterna() (bool, error) {
	return m.detectarAlteracaoExternaResp, m.detectarAlteracaoExternaErr
}

func TestNewManager(t *testing.T) {
	cofre := NovoCofre()
	repo := &mockRepository{}

	manager := NewManager(cofre, repo)

	if manager == nil {
		t.Fatal("NewManager returned nil")
	}

	if manager.IsLocked() {
		t.Error("New manager should not be locked")
	}

	if manager.Vault() == nil {
		t.Error("New manager should return vault")
	}

	if manager.IsModified() {
		t.Error("New manager should not be modified")
	}
}

func TestLock(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()

	// Create secret with sensitive field
	pasta := cofre.PastaGeral()
	segredo := &Segredo{
		nome:  "Test",
		pasta: pasta,
		campos: []CampoSegredo{
			{nome: "Senha", tipo: TipoCampoSensivel, valor: []byte("secret123")},
		},
		observacao:   CampoSegredo{nome: "Observação", tipo: TipoCampoComum, valor: []byte("note")},
		estadoSessao: EstadoOriginal,
	}
	pasta.segredos = append(pasta.segredos, segredo)

	manager := NewManager(cofre, &mockRepository{})
	manager.senha = []byte("master password")

	// Lock vault
	manager.Lock()

	// Verify locked state
	if !manager.IsLocked() {
		t.Error("Manager should be locked after Lock()")
	}

	if manager.Vault() != nil {
		t.Error("Vault() should return nil when locked")
	}

	// Verify sensitive data wiped
	if segredo.campos[0].valor != nil {
		t.Error("Sensitive field value should be wiped (nil) after Lock()")
	}

	if segredo.observacao.valor != nil {
		t.Error("Observation value should be wiped (nil) after Lock()")
	}

	// Verify password wiped (checking manager internals via package access)
	if manager.senha != nil {
		t.Error("Master password should be wiped (nil) after Lock()")
	}
}

func TestAlterarConfiguracoes(t *testing.T) {
	cofre := NovoCofre()
	manager := NewManager(cofre, &mockRepository{})

	// Valid configuration
	novasConfig := Configuracoes{
		tempoBloqueioInatividadeMinutos:      10,
		tempoOcultarSegredoSegundos:          20,
		tempoLimparAreaTransferenciaSegundos: 40,
	}

	err := manager.AlterarConfiguracoes(novasConfig)
	if err != nil {
		t.Fatalf("AlterarConfiguracoes with valid config failed: %v", err)
	}

	if !manager.IsModified() {
		t.Error("Vault should be marked modified after config change")
	}

	// Verify configuration updated
	config := cofre.Configuracoes()
	if config.tempoBloqueioInatividadeMinutos != 10 {
		t.Errorf("Expected tempoBloqueio=10, got %d", config.tempoBloqueioInatividadeMinutos)
	}
}

func TestAlterarConfiguracoesInvalid(t *testing.T) {
	cofre := NovoCofre()
	manager := NewManager(cofre, &mockRepository{})

	tests := []struct {
		nome   string
		config Configuracoes
	}{
		{"zero tempoBloqueio", Configuracoes{0, 15, 30}},
		{"negative tempoBloqueio", Configuracoes{-1, 15, 30}},
		{"zero tempoOcultar", Configuracoes{5, 0, 30}},
		{"zero tempoLimpar", Configuracoes{5, 15, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			err := manager.AlterarConfiguracoes(tt.config)
			if !errors.Is(err, ErrConfigInvalida) {
				t.Errorf("Expected ErrConfigInvalida, got %v", err)
			}

			if manager.IsModified() {
				t.Error("Vault should not be modified after invalid config")
			}
		})
	}
}

func TestSalvarWhenLocked(t *testing.T) {
	cofre := NovoCofre()
	manager := NewManager(cofre, &mockRepository{})

	manager.Lock()

	err := manager.Salvar()
	if !errors.Is(err, ErrCofreBloqueado) {
		t.Errorf("Expected ErrCofreBloqueado when saving locked vault, got %v", err)
	}
}

func TestSalvarSuccess(t *testing.T) {
	cofre := NovoCofre()
	repo := &mockRepository{}
	manager := NewManager(cofre, repo)

	// Mark as modified
	cofre.modificado = true

	err := manager.Salvar()
	if err != nil {
		t.Fatalf("Salvar failed: %v", err)
	}

	if !repo.salvarCalled {
		t.Error("Repository Salvar should have been called")
	}

	if manager.IsModified() {
		t.Error("Vault should not be modified after successful save")
	}
}

func TestSalvarFailureKeepsModifiedFlag(t *testing.T) {
	cofre := NovoCofre()
	repo := &mockRepository{salvarError: errors.New("disk full")}
	manager := NewManager(cofre, repo)

	cofre.modificado = true

	err := manager.Salvar()
	if err == nil {
		t.Fatal("Expected Salvar to return error")
	}

	// Vault should still be marked modified after failed save
	if !manager.IsModified() {
		t.Error("Vault should remain modified after failed save")
	}
}

// TestFullWorkflow exercises all Manager methods in a complete workflow.
// Creates vault structure, performs all operations, validates state.
func TestFullWorkflow(t *testing.T) {
	// Create vault with default content
	cofre := NovoCofre()
	err := cofre.InicializarConteudoPadrao()
	if err != nil {
		t.Fatalf("Failed to initialize vault: %v", err)
	}

	manager := NewManager(cofre, &mockRepository{})

	// Verify default content
	if len(cofre.Modelos()) != 3 {
		t.Errorf("Expected 3 default models, got %d", len(cofre.Modelos()))
	}

	// Create custom template
	modelo, err := manager.CriarModelo("CustomTemplate", []CampoModelo{
		{nome: "Field1", tipo: TipoCampoComum},
		{nome: "Field2", tipo: TipoCampoSensivel},
	})
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// Create folder structure
	pastaGeral := cofre.PastaGeral()
	pasta1, err := manager.CriarPasta(pastaGeral, "TestFolder", 0)
	if err != nil {
		t.Fatalf("Failed to create folder: %v", err)
	}

	// Create secrets
	secret1, err := manager.CriarSegredo(pasta1, "Secret1", modelo)
	if err != nil {
		t.Fatalf("Failed to create secret1: %v", err)
	}

	secret2, err := manager.CriarSegredo(pasta1, "Secret2", modelo)
	if err != nil {
		t.Fatalf("Failed to create secret2: %v", err)
	}

	// Edit secret
	err = manager.RenomearSegredo(secret1, "RenamedSecret")
	if err != nil {
		t.Fatalf("Failed to rename secret: %v", err)
	}

	err = manager.EditarCampoSegredo(secret1, 0, []byte("test value"))
	if err != nil {
		t.Fatalf("Failed to edit field: %v", err)
	}

	// Toggle favorite
	err = manager.AlternarFavoritoSegredo(secret1)
	if err != nil {
		t.Fatalf("Failed to toggle favorite: %v", err)
	}

	// Verify favorites
	favoritos := manager.ListarFavoritos()
	if len(favoritos) != 1 {
		t.Errorf("Expected 1 favorite, got %d", len(favoritos))
	}

	// Search
	results := manager.Buscar("renamed")
	if len(results) != 1 {
		t.Errorf("Expected 1 search result, got %d", len(results))
	}

	// Duplicate secret
	_, err = manager.DuplicarSegredo(secret2)
	if err != nil {
		t.Fatalf("Failed to duplicate secret: %v", err)
	}

	// Delete secret
	err = manager.ExcluirSegredo(secret2)
	if err != nil {
		t.Fatalf("Failed to delete secret: %v", err)
	}

	// Verify deleted excluded from search
	results = manager.Buscar("Secret2")
	if len(results) != 1 { // Should find the duplicate, not the deleted original
		t.Errorf("Expected 1 result (duplicate only), got %d", len(results))
	}

	// Lock and verify
	manager.Lock()
	if !manager.IsLocked() {
		t.Error("Manager should be locked")
	}
	if manager.Vault() != nil {
		t.Error("Vault should be nil when locked")
	}
}

// TestAtomicSave validates two-phase commit pattern per D-17.
// Save failure doesn't cause data loss in memory.
func TestAtomicSave(t *testing.T) {
	cofre := NovoCofre()
	repo := &mockRepository{salvarError: errors.New("simulated failure")}
	manager := NewManager(cofre, repo)

	// Create model and secret
	modelo, err := manager.CriarModelo("Test", []CampoModelo{
		{nome: "Field", tipo: TipoCampoComum},
	})
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	pasta := cofre.PastaGeral()
	secret1, err := manager.CriarSegredo(pasta, "ToDelete", modelo)
	if err != nil {
		t.Fatalf("Failed to create secret: %v", err)
	}

	// Simulate secret having been persisted (EstadoOriginal) so soft-delete path is used
	secret1.estadoSessao = EstadoOriginal

	// Mark for deletion (soft delete — EstadoOriginal → EstadoExcluido)
	err = manager.ExcluirSegredo(secret1)
	if err != nil {
		t.Fatalf("Failed to delete secret: %v", err)
	}

	// Verify deleted in memory (estadoSessao == Excluido)
	if secret1.estadoSessao != EstadoExcluido {
		t.Error("Secret should be marked Excluido")
	}

	// Attempt save (will fail)
	err = manager.Salvar()
	if err == nil {
		t.Fatal("Expected save to fail")
	}

	// Verify secret still exists in memory (not removed despite being Excluido)
	// Because save failed, finalizarExclusoes was never called
	found := false
	for _, s := range pasta.Segredos() {
		if s == secret1 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Secret should still exist in memory after failed save (atomic save)")
	}

	// Now succeed the save
	repo.salvarError = nil
	err = manager.Salvar()
	if err != nil {
		t.Fatalf("Second save failed: %v", err)
	}

	// Verify secret removed from memory after successful save
	found = false
	for _, s := range pasta.Segredos() {
		if s == secret1 {
			found = true
			break
		}
	}
	if found {
		t.Error("Secret should be removed from memory after successful save")
	}
}

// TestCycleDetection validates hierarchy protection.
// Moving folder into its own descendant returns ErrCycleDetected.
func TestCycleDetection(t *testing.T) {
	cofre := NovoCofre()
	manager := NewManager(cofre, &mockRepository{})

	pastaGeral := cofre.PastaGeral()

	// Create hierarchy: A → B → C
	pastaA, err := manager.CriarPasta(pastaGeral, "A", 0)
	if err != nil {
		t.Fatalf("Failed to create A: %v", err)
	}

	pastaB, err := manager.CriarPasta(pastaA, "B", 0)
	if err != nil {
		t.Fatalf("Failed to create B: %v", err)
	}

	pastaC, err := manager.CriarPasta(pastaB, "C", 0)
	if err != nil {
		t.Fatalf("Failed to create C: %v", err)
	}

	// Attempt to move A into C (would create cycle: C → A → B → C)
	err = manager.MoverPasta(pastaA, pastaC)
	if !errors.Is(err, ErrCycleDetected) {
		t.Errorf("Expected ErrCycleDetected, got %v", err)
	}

	// Verify A still in original location
	if pastaA.Pai() != pastaGeral {
		t.Error("Pasta A should still be under Pasta Geral after failed move")
	}
}

// TestPromotion validates automatic conflict resolution per D-27.
// Deleting folder promotes children with numeric suffix for conflicts.
func TestPromotion(t *testing.T) {
	cofre := NovoCofre()
	manager := NewManager(cofre, &mockRepository{})

	// Create model
	modelo, err := manager.CriarModelo("Item", []CampoModelo{
		{nome: "Value", tipo: TipoCampoComum},
	})
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	pastaGeral := cofre.PastaGeral()

	// Create structure:
	// Pasta Geral
	//   ├─ Folder A
	//   │    └─ Secret "Conflict" (will conflict)
	//   └─ Secret "Conflict" (existing)
	folderA, err := manager.CriarPasta(pastaGeral, "Folder A", 0)
	if err != nil {
		t.Fatalf("Failed to create Folder A: %v", err)
	}

	// Create conflicting secrets
	_, err = manager.CriarSegredo(pastaGeral, "Conflict", modelo)
	if err != nil {
		t.Fatalf("Failed to create secret in Pasta Geral: %v", err)
	}

	secretInA, err := manager.CriarSegredo(folderA, "Conflict", modelo)
	if err != nil {
		t.Fatalf("Failed to create secret in Folder A: %v", err)
	}

	// Delete Folder A (should promote secretInA with renamed to avoid conflict)
	renomeacoes, err := manager.ExcluirPasta(folderA)
	if err != nil {
		t.Fatalf("Failed to delete Folder A: %v", err)
	}

	// Verify renomeacao occurred
	if len(renomeacoes) != 1 {
		t.Fatalf("Expected 1 renomeacao, got %d", len(renomeacoes))
	}

	if renomeacoes[0].Antigo != "Conflict" {
		t.Errorf("Expected original name 'Conflict', got '%s'", renomeacoes[0].Antigo)
	}

	// Verify promoted secret has new name
	if secretInA.Nome() == "Conflict" {
		t.Error("Promoted secret should have been renamed to avoid conflict")
	}

	// Should be "Conflict (1)" or similar
	if secretInA.Nome() != "Conflict (1)" {
		t.Errorf("Expected renamed to 'Conflict (1)', got '%s'", secretInA.Nome())
	}

	// Verify secret now in Pasta Geral
	if secretInA.Pasta() != pastaGeral {
		t.Error("Promoted secret should now be in Pasta Geral")
	}
}

// TestDuplication validates independent state per D-15.
// Duplicated secret has independent state and "(N)" name progression.
func TestDuplication(t *testing.T) {
	cofre := NovoCofre()
	manager := NewManager(cofre, &mockRepository{})

	// Create model
	modelo, err := manager.CriarModelo("Note", []CampoModelo{
		{nome: "Content", tipo: TipoCampoComum},
	})
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	pasta := cofre.PastaGeral()

	// Create original secret
	original, err := manager.CriarSegredo(pasta, "Original", modelo)
	if err != nil {
		t.Fatalf("Failed to create original: %v", err)
	}

	// Edit original
	err = manager.EditarCampoSegredo(original, 0, []byte("original content"))
	if err != nil {
		t.Fatalf("Failed to edit original field: %v", err)
	}

	// Toggle favorite on original
	err = manager.AlternarFavoritoSegredo(original)
	if err != nil {
		t.Fatalf("Failed to toggle favorite: %v", err)
	}

	// Duplicate secret
	duplicate, err := manager.DuplicarSegredo(original)
	if err != nil {
		t.Fatalf("Failed to duplicate: %v", err)
	}

	// Verify name progression (Original → Original (1))
	if duplicate.Nome() != "Original (1)" {
		t.Errorf("Expected duplicate name 'Original (1)', got '%s'", duplicate.Nome())
	}

	// Verify independent state: duplicate NOT favorite (favorito state not copied)
	if duplicate.Favorito() {
		t.Error("Duplicate should not inherit favorito state (independent)")
	}

	// Verify same folder
	if duplicate.Pasta() != pasta {
		t.Error("Duplicate should be in same folder as original")
	}

	// Verify field value copied
	campo := duplicate.Campos()[0]
	if string(campo.valor) != "original content" {
		t.Errorf("Expected duplicate field value 'original content', got '%s'", string(campo.valor))
	}

	// Verify independent modification: change duplicate doesn't affect original
	err = manager.EditarCampoSegredo(duplicate, 0, []byte("modified duplicate"))
	if err != nil {
		t.Fatalf("Failed to edit duplicate: %v", err)
	}

	// Original should still have original content
	originalCampo := original.Campos()[0]
	if string(originalCampo.valor) != "original content" {
		t.Error("Original should not be affected by duplicate modification")
	}

	// Duplicate a second time (test name progression: Original (1) → Original (2))
	duplicate2, err := manager.DuplicarSegredo(original)
	if err != nil {
		t.Fatalf("Failed to duplicate second time: %v", err)
	}

	if duplicate2.Nome() != "Original (2)" {
		t.Errorf("Expected second duplicate name 'Original (2)', got '%s'", duplicate2.Nome())
	}
}

// TestAlternarFavorito_NaoAlteraDataUltimaModificacao verifies that favoriting a secret
// does NOT update the secret's dataUltimaModificacao (D-08, BUG-01).
// The vault IS marked modified (cofre.modificado = true) but the secret timestamp stays.
func TestAlternarFavorito_NaoAlteraDataUltimaModificacao(t *testing.T) {
	cofre := NovoCofre()
	if err := cofre.InicializarConteudoPadrao(); err != nil {
		t.Fatal(err)
	}
	manager := NewManager(cofre, &mockRepository{})

	pasta := cofre.PastaGeral()
	modelo := cofre.Modelos()[0]
	seg, err := manager.CriarSegredo(pasta, "Test Secret", modelo)
	if err != nil {
		t.Fatal(err)
	}

	// Capture modification time after creation
	dataBefore := seg.DataUltimaModificacao()

	// Toggle favorite
	if err := manager.AlternarFavoritoSegredo(seg); err != nil {
		t.Fatalf("AlternarFavoritoSegredo() error: %v", err)
	}

	// Secret modification time MUST NOT change (D-08)
	if !seg.DataUltimaModificacao().Equal(dataBefore) {
		t.Errorf("AlternarFavoritoSegredo changed dataUltimaModificacao: got %v, want %v (unchanged)",
			seg.DataUltimaModificacao(), dataBefore)
	}

	// Vault MUST be marked modified (cofre.marcarModificado was called)
	if !manager.IsModified() {
		t.Error("Vault should be marked modified after favoriting")
	}

	// Secret MUST be toggled
	if !seg.Favorito() {
		t.Error("Secret should be marked as favorite after toggle")
	}
}

// TestCriarSegredo_EstadoIncluido verifies that CriarSegredo returns a secret
// with estadoSessao == EstadoIncluido, not EstadoModificado (D-05).
func TestCriarSegredo_EstadoIncluido(t *testing.T) {
	cofre := NovoCofre()
	if err := cofre.InicializarConteudoPadrao(); err != nil {
		t.Fatal(err)
	}
	manager := NewManager(cofre, &mockRepository{})

	pasta := cofre.PastaGeral()
	modelo := cofre.Modelos()[0]
	seg, err := manager.CriarSegredo(pasta, "New Secret", modelo)
	if err != nil {
		t.Fatal(err)
	}

	if seg.EstadoSessao() != EstadoIncluido {
		t.Errorf("CriarSegredo estadoSessao = %v, want EstadoIncluido (%v)",
			seg.EstadoSessao(), EstadoIncluido)
	}
}

// TestDuplicarSegredo_EstadoIncluido verifies that DuplicarSegredo returns a duplicate
// with estadoSessao == EstadoIncluido, not EstadoModificado (D-05).
func TestDuplicarSegredo_EstadoIncluido(t *testing.T) {
	cofre := NovoCofre()
	if err := cofre.InicializarConteudoPadrao(); err != nil {
		t.Fatal(err)
	}
	manager := NewManager(cofre, &mockRepository{})

	pasta := cofre.PastaGeral()
	modelo := cofre.Modelos()[0]
	original, err := manager.CriarSegredo(pasta, "Original Secret", modelo)
	if err != nil {
		t.Fatal(err)
	}

	dup, err := manager.DuplicarSegredo(original)
	if err != nil {
		t.Fatal(err)
	}

	if dup.EstadoSessao() != EstadoIncluido {
		t.Errorf("DuplicarSegredo estadoSessao = %v, want EstadoIncluido (%v)",
			dup.EstadoSessao(), EstadoIncluido)
	}
}
