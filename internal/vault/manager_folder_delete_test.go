package vault

import (
	"errors"
	"testing"
)

// Task 5 Tests: ExcluirPasta with promotion and conflict resolution

func TestExcluirPasta_Success_NoConflicts(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	// Create folder to delete with a subfolder and a secret
	sitesEApps := cofre.PastaGeral().Subpastas()[0]
	folderToDelete, _ := manager.CriarPasta(sitesEApps, "ToDelete", 0)

	// Add subfolder
	_, _ = manager.CriarPasta(folderToDelete, "Child", 0)

	// Add secret (we'll need to implement secret creation later, but mock it for now)
	// For now, create a secret directly since we don't have CriarSegredo yet
	segredo := &Segredo{
		nome:         "TestSecret",
		pasta:        folderToDelete,
		campos:       []CampoSegredo{},
		observacao:   CampoSegredo{nome: "Observação", tipo: TipoCampoComum, valor: []byte("")},
		estadoSessao: EstadoOriginal,
	}
	folderToDelete.segredos = append(folderToDelete.segredos, segredo)

	// Clear modified flag
	cofre.modificado = false

	// Delete folder
	renomeacoes, err := manager.ExcluirPasta(folderToDelete)
	if err != nil {
		t.Fatalf("ExcluirPasta failed: %v", err)
	}

	// No conflicts, so no renomeacoes
	if len(renomeacoes) != 0 {
		t.Errorf("Expected no renomeacoes, got %d", len(renomeacoes))
	}

	// Verify folder removed from parent
	subpastas := sitesEApps.Subpastas()
	for _, sub := range subpastas {
		if sub == folderToDelete {
			t.Error("Deleted folder still in parent")
		}
	}

	// Verify child promoted to parent
	found := false
	for _, sub := range subpastas {
		if sub.Nome() == "Child" {
			found = true
			if sub.Pai() != sitesEApps {
				t.Error("Child parent should be Sites e Apps")
			}
		}
	}
	if !found {
		t.Error("Child folder not promoted to parent")
	}

	// Verify secret promoted to parent
	segredos := sitesEApps.Segredos()
	foundSecret := false
	for _, s := range segredos {
		if s.Nome() == "TestSecret" {
			foundSecret = true
			if s.Pasta() != sitesEApps {
				t.Error("Secret parent should be Sites e Apps")
			}
		}
	}
	if !foundSecret {
		t.Error("Secret not promoted to parent")
	}

	if !manager.IsModified() {
		t.Error("Vault should be marked modified after folder deletion")
	}
}

func TestExcluirPasta_SecretNameConflict_RenamedWithSuffix(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	sitesEApps := cofre.PastaGeral().Subpastas()[0]

	// Create secret in parent
	secretInParent := &Segredo{
		nome:         "Conflict",
		pasta:        sitesEApps,
		campos:       []CampoSegredo{},
		observacao:   CampoSegredo{nome: "Observação", tipo: TipoCampoComum, valor: []byte("")},
		estadoSessao: EstadoOriginal,
	}
	sitesEApps.segredos = append(sitesEApps.segredos, secretInParent)

	// Create folder to delete
	folderToDelete, _ := manager.CriarPasta(sitesEApps, "ToDelete", 0)

	// Create secret with same name in folder to delete
	secretInChild := &Segredo{
		nome:         "Conflict",
		pasta:        folderToDelete,
		campos:       []CampoSegredo{},
		observacao:   CampoSegredo{nome: "Observação", tipo: TipoCampoComum, valor: []byte("")},
		estadoSessao: EstadoOriginal,
	}
	folderToDelete.segredos = append(folderToDelete.segredos, secretInChild)

	// Delete folder
	renomeacoes, err := manager.ExcluirPasta(folderToDelete)
	if err != nil {
		t.Fatalf("ExcluirPasta failed: %v", err)
	}

	// Should have one renomeacao
	if len(renomeacoes) != 1 {
		t.Fatalf("Expected 1 renomeacao, got %d", len(renomeacoes))
	}

	if renomeacoes[0].Antigo != "Conflict" {
		t.Errorf("Expected Antigo='Conflict', got %q", renomeacoes[0].Antigo)
	}

	if renomeacoes[0].Novo != "Conflict (1)" {
		t.Errorf("Expected Novo='Conflict (1)', got %q", renomeacoes[0].Novo)
	}

	if renomeacoes[0].Pasta != "Sites e Apps" {
		t.Errorf("Expected Pasta='Sites e Apps', got %q", renomeacoes[0].Pasta)
	}

	// Verify secret was renamed
	segredos := sitesEApps.Segredos()
	foundRenamed := false
	for _, s := range segredos {
		if s.Nome() == "Conflict (1)" {
			foundRenamed = true
		}
	}
	if !foundRenamed {
		t.Error("Conflicting secret should have been renamed to 'Conflict (1)'")
	}
}

func TestExcluirPasta_SubfolderNameConflict_ContentsMerged(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	sitesEApps := cofre.PastaGeral().Subpastas()[0]

	// Create "Common" folder in parent
	commonInParent, _ := manager.CriarPasta(sitesEApps, "Common", 0)

	// Add a secret to Common in parent
	secretInParentCommon := &Segredo{
		nome:         "ParentSecret",
		pasta:        commonInParent,
		campos:       []CampoSegredo{},
		observacao:   CampoSegredo{nome: "Observação", tipo: TipoCampoComum, valor: []byte("")},
		estadoSessao: EstadoOriginal,
	}
	commonInParent.segredos = append(commonInParent.segredos, secretInParentCommon)

	// Create folder to delete
	folderToDelete, _ := manager.CriarPasta(sitesEApps, "ToDelete", 1)

	// Create "Common" folder in folder to delete
	commonInChild, _ := manager.CriarPasta(folderToDelete, "Common", 0)

	// Add a secret to Common in child
	secretInChildCommon := &Segredo{
		nome:         "ChildSecret",
		pasta:        commonInChild,
		campos:       []CampoSegredo{},
		observacao:   CampoSegredo{nome: "Observação", tipo: TipoCampoComum, valor: []byte("")},
		estadoSessao: EstadoOriginal,
	}
	commonInChild.segredos = append(commonInChild.segredos, secretInChildCommon)

	// Delete folder
	_, err := manager.ExcluirPasta(folderToDelete)
	if err != nil {
		t.Fatalf("ExcluirPasta failed: %v", err)
	}

	// Verify only one "Common" folder exists in parent
	subpastas := sitesEApps.Subpastas()
	commonCount := 0
	var mergedCommon *Pasta
	for _, sub := range subpastas {
		if sub.Nome() == "Common" {
			commonCount++
			mergedCommon = sub
		}
	}
	if commonCount != 1 {
		t.Fatalf("Expected 1 'Common' folder after merge, found %d", commonCount)
	}

	// Verify both secrets are in merged folder
	segredos := mergedCommon.Segredos()
	foundParentSecret := false
	foundChildSecret := false
	for _, s := range segredos {
		if s.Nome() == "ParentSecret" {
			foundParentSecret = true
		}
		if s.Nome() == "ChildSecret" {
			foundChildSecret = true
		}
	}
	if !foundParentSecret {
		t.Error("ParentSecret should be in merged Common folder")
	}
	if !foundChildSecret {
		t.Error("ChildSecret should be in merged Common folder")
	}
}

func TestExcluirPasta_StateDeletedSecretsRetainState(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	sitesEApps := cofre.PastaGeral().Subpastas()[0]

	// Create folder to delete
	folderToDelete, _ := manager.CriarPasta(sitesEApps, "ToDelete", 0)

	// Create secret and mark as deleted
	deletedSecret := &Segredo{
		nome:         "DeletedSecret",
		pasta:        folderToDelete,
		campos:       []CampoSegredo{},
		observacao:   CampoSegredo{nome: "Observação", tipo: TipoCampoComum, valor: []byte("")},
		estadoSessao: EstadoExcluido, // Marked for deletion
	}
	folderToDelete.segredos = append(folderToDelete.segredos, deletedSecret)

	// Delete folder
	_, err := manager.ExcluirPasta(folderToDelete)
	if err != nil {
		t.Fatalf("ExcluirPasta failed: %v", err)
	}

	// Verify secret was promoted and still has EstadoExcluido (FOLDER-05)
	segredos := sitesEApps.Segredos()
	foundDeleted := false
	for _, s := range segredos {
		if s.Nome() == "DeletedSecret" {
			foundDeleted = true
			if s.EstadoSessao() != EstadoExcluido {
				t.Errorf("Expected EstadoExcluido, got %v", s.EstadoSessao())
			}
		}
	}
	if !foundDeleted {
		t.Error("Deleted secret should have been promoted to parent")
	}
}

func TestExcluirPasta_PastaGeralProtection(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	_, err := manager.ExcluirPasta(cofre.PastaGeral())
	if !errors.Is(err, ErrPastaGeralNaoExcluivel) {
		t.Errorf("Expected ErrPastaGeralNaoExcluivel, got %v", err)
	}

	if manager.IsModified() {
		t.Error("Vault should not be modified after failed deletion")
	}
}

func TestExcluirPasta_WhenLocked(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	sitesEApps := cofre.PastaGeral().Subpastas()[0]
	manager.Lock()

	_, err := manager.ExcluirPasta(sitesEApps)
	if !errors.Is(err, ErrCofreBloqueado) {
		t.Errorf("Expected ErrCofreBloqueado, got %v", err)
	}
}

func TestExcluirPasta_MultipleConflicts(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	sitesEApps := cofre.PastaGeral().Subpastas()[0]

	// Create three secrets in parent
	for i := 1; i <= 3; i++ {
		secret := &Segredo{
			nome:         "Secret",
			pasta:        sitesEApps,
			campos:       []CampoSegredo{},
			observacao:   CampoSegredo{nome: "Observação", tipo: TipoCampoComum, valor: []byte("")},
			estadoSessao: EstadoOriginal,
		}
		sitesEApps.segredos = append(sitesEApps.segredos, secret)
	}

	// Create folder to delete with three "Secret" named secrets
	folderToDelete, _ := manager.CriarPasta(sitesEApps, "ToDelete", 0)
	for i := 1; i <= 3; i++ {
		secret := &Segredo{
			nome:         "Secret",
			pasta:        folderToDelete,
			campos:       []CampoSegredo{},
			observacao:   CampoSegredo{nome: "Observação", tipo: TipoCampoComum, valor: []byte("")},
			estadoSessao: EstadoOriginal,
		}
		folderToDelete.segredos = append(folderToDelete.segredos, secret)
	}

	// Delete folder
	renomeacoes, err := manager.ExcluirPasta(folderToDelete)
	if err != nil {
		t.Fatalf("ExcluirPasta failed: %v", err)
	}

	// Should have 3 renomeacoes (all conflicted)
	if len(renomeacoes) != 3 {
		t.Fatalf("Expected 3 renomeacoes, got %d", len(renomeacoes))
	}

	// Verify secrets were renamed to Secret (1), Secret (2), Secret (3)
	segredos := sitesEApps.Segredos()
	found1 := false
	found2 := false
	found3 := false
	for _, s := range segredos {
		if s.Nome() == "Secret (1)" {
			found1 = true
		}
		if s.Nome() == "Secret (2)" {
			found2 = true
		}
		if s.Nome() == "Secret (3)" {
			found3 = true
		}
	}
	if !found1 || !found2 || !found3 {
		t.Error("All conflicting secrets should have been renamed with unique suffixes")
	}
}
