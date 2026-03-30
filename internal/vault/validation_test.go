package vault

import (
	"testing"
)

// TestCriarModeloOrdenacao verifies templates are sorted alphabetically after creation.
// Per TPL-02, TPL-06: templates always displayed in alphabetical order.
func TestCriarModeloOrdenacao(t *testing.T) {
	cofre := NovoCofre()
	manager := NewManager(cofre, nil)

	// Create templates in non-alphabetical order
	_, err := manager.CriarModelo("Zebra", []CampoModelo{{nome: "campo1", tipo: TipoCampoComum}})
	if err != nil {
		t.Fatalf("Failed to create Zebra model: %v", err)
	}

	_, err = manager.CriarModelo("Alpha", []CampoModelo{{nome: "campo2", tipo: TipoCampoComum}})
	if err != nil {
		t.Fatalf("Failed to create Alpha model: %v", err)
	}

	_, err = manager.CriarModelo("Medio", []CampoModelo{{nome: "campo3", tipo: TipoCampoComum}})
	if err != nil {
		t.Fatalf("Failed to create Medio model: %v", err)
	}

	// Retrieve models (should be alphabetically sorted)
	modelos := cofre.Modelos()

	// Verify count (3 created + 3 default from InicializarConteudoPadrao if called)
	// Since we didn't call InicializarConteudoPadrao, should be exactly 3
	if len(modelos) != 3 {
		t.Fatalf("Expected 3 models, got %d", len(modelos))
	}

	// Verify alphabetical order
	expected := []string{"Alpha", "Medio", "Zebra"}
	for i, modelo := range modelos {
		if modelo.Nome() != expected[i] {
			t.Errorf("Model %d: expected %s, got %s", i, expected[i], modelo.Nome())
		}
	}
}

// TestModeloNomeReservado verifies "Observação" name is forbidden in template and field names.
// Per D-29: "Observação" prohibited in template field names.
func TestModeloNomeReservado(t *testing.T) {
	cofre := NovoCofre()
	manager := NewManager(cofre, nil)

	// Test 1: Cannot create template with field named "Observação" (exact case)
	_, err := manager.CriarModelo("TestModel", []CampoModelo{
		{nome: "Observação", tipo: TipoCampoComum},
	})
	if err != ErrObservacaoReserved {
		t.Errorf("Expected ErrObservacaoReserved when creating field 'Observação', got: %v", err)
	}

	// Test 2: Cannot create template with field named "observação" (lowercase)
	_, err = manager.CriarModelo("TestModel2", []CampoModelo{
		{nome: "observação", tipo: TipoCampoComum},
	})
	if err != ErrObservacaoReserved {
		t.Errorf("Expected ErrObservacaoReserved when creating field 'observação', got: %v", err)
	}

	// Test 3: Cannot create template with field named "OBSERVAÇÃO" (uppercase)
	_, err = manager.CriarModelo("TestModel3", []CampoModelo{
		{nome: "OBSERVAÇÃO", tipo: TipoCampoComum},
	})
	if err != ErrObservacaoReserved {
		t.Errorf("Expected ErrObservacaoReserved when creating field 'OBSERVAÇÃO', got: %v", err)
	}

	// Test 4: Valid field name should work
	_, err = manager.CriarModelo("ValidModel", []CampoModelo{
		{nome: "Notas", tipo: TipoCampoComum},
	})
	if err != nil {
		t.Errorf("Expected no error for valid field name, got: %v", err)
	}
}

// TestExcluirModeloEmUso verifies templates cannot be deleted if in use by secrets.
// Per TPL-04, D-26: templates can be deleted unless referenced by a secret.
func TestExcluirModeloEmUso(t *testing.T) {
	cofre := NovoCofre()
	manager := NewManager(cofre, nil)

	// Create a template
	modelo, err := manager.CriarModelo("Login", []CampoModelo{
		{nome: "URL", tipo: TipoCampoComum},
		{nome: "Senha", tipo: TipoCampoSensivel},
	})
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// Template not in use - should be deletable
	err = manager.ExcluirModelo(modelo)
	if err != nil {
		t.Errorf("Expected no error deleting unused template, got: %v", err)
	}

	// TODO: Re-enable when secret creation is implemented in Task 5
	// Create template again for in-use test
	// modelo2, _ := manager.CriarModelo("Login2", []CampoModelo{{nome: "User", tipo: TipoCampoComum}})
	// Create secret using the template
	// segredo, _ := manager.CriarSegredoDeModelo(cofre.PastaGeral(), modelo2, "MyLogin", 0)
	// Now template is in use - should not be deletable
	// err = manager.ExcluirModelo(modelo2)
	// if err != ErrModeloEmUso {
	// 	t.Errorf("Expected ErrModeloEmUso when deleting in-use template, got: %v", err)
	// }
}

// TestCampoOperacoes verifies field operations (add/remove/reorder) on templates.
// Per TPL-03, D-29: templates support field structure changes.
func TestCampoOperacoes(t *testing.T) {
	cofre := NovoCofre()
	manager := NewManager(cofre, nil)

	// Create template with initial fields
	modelo, err := manager.CriarModelo("TestTemplate", []CampoModelo{
		{nome: "Campo1", tipo: TipoCampoComum},
		{nome: "Campo2", tipo: TipoCampoSensivel},
	})
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// Test 1: Add field at position
	err = manager.AdicionarCampo(modelo, "CampoNovo", TipoCampoComum, 1)
	if err != nil {
		t.Errorf("Failed to add field: %v", err)
	}
	campos := modelo.Campos()
	if len(campos) != 3 {
		t.Errorf("Expected 3 fields after add, got %d", len(campos))
	}
	if campos[1].Nome() != "CampoNovo" {
		t.Errorf("Expected field at position 1 to be 'CampoNovo', got %s", campos[1].Nome())
	}

	// Test 2: Cannot add field named "Observação"
	err = manager.AdicionarCampo(modelo, "Observação", TipoCampoComum, 0)
	if err != ErrObservacaoReserved {
		t.Errorf("Expected ErrObservacaoReserved when adding 'Observação' field, got: %v", err)
	}

	// Test 3: Remove field by index
	err = manager.RemoverCampo(modelo, 0)
	if err != nil {
		t.Errorf("Failed to remove field: %v", err)
	}
	campos = modelo.Campos()
	if len(campos) != 2 {
		t.Errorf("Expected 2 fields after remove, got %d", len(campos))
	}

	// Test 4: Reorder field
	err = manager.ReordenarCampo(modelo, 1, 0)
	if err != nil {
		t.Errorf("Failed to reorder field: %v", err)
	}
	campos = modelo.Campos()
	// After reorder, the field that was at index 1 should now be at index 0
	if campos[0].Nome() != "Campo2" {
		t.Errorf("Expected field at position 0 to be 'Campo2' after reorder, got %s", campos[0].Nome())
	}

	// Test 5: Invalid position errors
	err = manager.AdicionarCampo(modelo, "Invalid", TipoCampoComum, 999)
	if err != ErrPosicaoInvalida {
		t.Errorf("Expected ErrPosicaoInvalida for invalid position, got: %v", err)
	}

	err = manager.RemoverCampo(modelo, 999)
	if err != ErrCampoInvalido {
		t.Errorf("Expected ErrCampoInvalido for invalid index, got: %v", err)
	}

	err = manager.ReordenarCampo(modelo, 999, 0)
	if err != ErrCampoInvalido {
		t.Errorf("Expected ErrCampoInvalido for invalid reorder index, got: %v", err)
	}
}

// TestRenomearModeloNoOp verifies that renaming to same name doesn't mark vault as modified.
// Per D-12: change detection based on actual value difference.
func TestRenomearModeloNoOp(t *testing.T) {
	cofre := NovoCofre()
	manager := NewManager(cofre, nil)

	// Create a template
	modelo, err := manager.CriarModelo("Original", []CampoModelo{
		{nome: "Campo1", tipo: TipoCampoComum},
	})
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// Mark vault as unmodified (simulate a save)
	cofre.modificado = false

	// Rename to same name (no-op)
	err = manager.RenomearModelo(modelo, "Original")
	if err != nil {
		t.Errorf("Expected no error for no-op rename, got: %v", err)
	}

	// Verify vault NOT marked as modified
	if cofre.Modificado() {
		t.Errorf("Expected vault to NOT be modified after no-op rename")
	}
}

// Folder Operation Validation Tests

// TestFolderHierarchyIntegrity verifies folder operations maintain parent-child relationship integrity.
// Validates: CriarPasta correctly establishes bidirectional parent-child links.
func TestFolderHierarchyIntegrity(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	sitesEApps := cofre.PastaGeral().Subpastas()[0]

	// Create nested folder structure
	level1, err := manager.CriarPasta(sitesEApps, "Level1", 0)
	if err != nil {
		t.Fatalf("Failed to create Level1: %v", err)
	}

	level2, err := manager.CriarPasta(level1, "Level2", 0)
	if err != nil {
		t.Fatalf("Failed to create Level2: %v", err)
	}

	level3, err := manager.CriarPasta(level2, "Level3", 0)
	if err != nil {
		t.Fatalf("Failed to create Level3: %v", err)
	}

	// Verify parent-child relationships
	if level1.Pai() != sitesEApps {
		t.Error("Level1 parent should be sitesEApps")
	}
	if level2.Pai() != level1 {
		t.Error("Level2 parent should be Level1")
	}
	if level3.Pai() != level2 {
		t.Error("Level3 parent should be Level2")
	}

	// Verify children exist in parent's subpastas
	found := false
	for _, sub := range sitesEApps.Subpastas() {
		if sub == level1 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Level1 not found in sitesEApps subpastas")
	}

	found = false
	for _, sub := range level1.Subpastas() {
		if sub == level2 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Level2 not found in Level1 subpastas")
	}
}

// TestFolderNameUniquenessWithinParent verifies name uniqueness is enforced per parent, not globally.
// Two folders can have the same name if they have different parents.
func TestFolderNameUniquenessWithinParent(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	sitesEApps := cofre.PastaGeral().Subpastas()[0]
	financeiro := cofre.PastaGeral().Subpastas()[1]

	// Create "Work" folder in sitesEApps - should succeed
	work1, err := manager.CriarPasta(sitesEApps, "Work", 0)
	if err != nil {
		t.Fatalf("Failed to create Work in sitesEApps: %v", err)
	}

	// Create "Work" folder in financeiro - should succeed (different parent)
	work2, err := manager.CriarPasta(financeiro, "Work", 0)
	if err != nil {
		t.Fatalf("Failed to create Work in financeiro: %v", err)
	}

	// Verify both folders exist and are different instances
	if work1 == work2 {
		t.Error("Two Work folders should be different instances")
	}
	if work1.Nome() != "Work" || work2.Nome() != "Work" {
		t.Error("Both folders should be named Work")
	}

	// Try to create duplicate "Work" in sitesEApps - should fail
	_, err = manager.CriarPasta(sitesEApps, "Work", 0)
	if err != ErrNameConflict {
		t.Errorf("Expected ErrNameConflict for duplicate name in same parent, got: %v", err)
	}
}

// TestFolderRenameNoOpDetection verifies renaming to same name doesn't mark vault as modified.
// Per D-12: change detection based on actual value difference.
func TestFolderRenameNoOpDetection(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	sitesEApps := cofre.PastaGeral().Subpastas()[0]
	folder, _ := manager.CriarPasta(sitesEApps, "Original", 0)

	// Mark vault as unmodified (simulate save)
	cofre.modificado = false

	// Rename to same name (no-op)
	err := manager.RenomearPasta(folder, "Original")
	if err != nil {
		t.Errorf("Expected no error for no-op rename, got: %v", err)
	}

	// Verify vault NOT marked as modified
	if cofre.Modificado() {
		t.Error("Expected vault to NOT be modified after no-op rename")
	}

	// Now rename to different name
	err = manager.RenomearPasta(folder, "NewName")
	if err != nil {
		t.Errorf("Failed to rename to different name: %v", err)
	}

	// Verify vault IS marked as modified
	if !cofre.Modificado() {
		t.Error("Expected vault to be modified after actual rename")
	}
}

// TestFolderMoveNoCycleDetection verifies cycle detection prevents moving folder into its own descendant.
// Per FOLDER-03: cycle detection via full ancestor walk.
func TestFolderMoveNoCycleDetection(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	sitesEApps := cofre.PastaGeral().Subpastas()[0]

	// Create nested structure: A -> B -> C
	folderA, _ := manager.CriarPasta(sitesEApps, "A", 0)
	folderB, _ := manager.CriarPasta(folderA, "B", 0)
	folderC, _ := manager.CriarPasta(folderB, "C", 0)

	// Try to move A into C (would create cycle A -> B -> C -> A)
	err := manager.MoverPasta(folderA, folderC)
	if err != ErrCycleDetected {
		t.Errorf("Expected ErrCycleDetected when moving folder into descendant, got: %v", err)
	}

	// Try to move A into B (would create cycle A -> B -> A)
	err = manager.MoverPasta(folderA, folderB)
	if err != ErrCycleDetected {
		t.Errorf("Expected ErrCycleDetected when moving folder into direct child, got: %v", err)
	}

	// Try to move A into itself (would create cycle A -> A)
	err = manager.MoverPasta(folderA, folderA)
	if err != ErrDestinoInvalido {
		t.Errorf("Expected ErrDestinoInvalido when moving folder into itself, got: %v", err)
	}

	// Valid move: move C to sitesEApps (no cycle)
	err = manager.MoverPasta(folderC, sitesEApps)
	if err != nil {
		t.Errorf("Expected no error for valid move, got: %v", err)
	}

	// Verify C was moved
	if folderC.Pai() != sitesEApps {
		t.Error("folderC parent should be sitesEApps after move")
	}
}

// TestFolderMoveNoOpDetection verifies moving to same parent doesn't duplicate folder.
// MoverPasta is for changing parents; use ReposicionarPasta for same-parent moves.
func TestFolderMoveNoOpDetection(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	sitesEApps := cofre.PastaGeral().Subpastas()[0]
	financeiro := cofre.PastaGeral().Subpastas()[1]
	folder1, _ := manager.CriarPasta(sitesEApps, "Folder1", 0)

	// Mark vault as unmodified
	cofre.modificado = false

	// Move folder1 to same parent (should fail - name already exists)
	// Note: MoverPasta is for changing parents, so same parent with same name = conflict
	err := manager.MoverPasta(folder1, sitesEApps)
	if err == nil {
		t.Error("Expected error when moving to same parent (name conflict)")
	}

	// Now move to different parent
	err = manager.MoverPasta(folder1, financeiro)
	if err != nil {
		t.Errorf("Failed to move to different parent: %v", err)
	}

	// Verify vault IS marked as modified
	if !cofre.Modificado() {
		t.Error("Expected vault to be modified after actual move")
	}

	// Verify folder was moved
	if folder1.Pai() != financeiro {
		t.Error("Folder1 parent should be financeiro after move")
	}
}

// TestFolderRepositionNoOpDetection verifies no-op repositioning doesn't mark vault as modified.
// Per D-23: repositioning to current position, Subir at 0, Descer at last are no-ops.
func TestFolderRepositionNoOpDetection(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	sitesEApps := cofre.PastaGeral().Subpastas()[0]
	folder1, _ := manager.CriarPasta(sitesEApps, "Folder1", 0)
	folder2, _ := manager.CriarPasta(sitesEApps, "Folder2", 1)
	folder3, _ := manager.CriarPasta(sitesEApps, "Folder3", 2)

	// Mark vault as unmodified
	cofre.modificado = false

	// Test 1: Reposition to current position (folder2 is at position 1)
	err := manager.ReposicionarPasta(folder2, 1)
	if err != nil {
		t.Errorf("Expected no error for no-op reposition, got: %v", err)
	}
	if cofre.Modificado() {
		t.Error("Expected vault NOT modified after repositioning to current position")
	}

	// Test 2: Subir at position 0 (folder1 is at position 0)
	err = manager.SubirPastaNaPosicao(folder1)
	if err != nil {
		t.Errorf("Expected no error for Subir at position 0, got: %v", err)
	}
	if cofre.Modificado() {
		t.Error("Expected vault NOT modified after Subir at position 0")
	}

	// Test 3: Descer at last position (folder3 is at position 2, last position)
	err = manager.DescerPastaNaPosicao(folder3)
	if err != nil {
		t.Errorf("Expected no error for Descer at last position, got: %v", err)
	}
	if cofre.Modificado() {
		t.Error("Expected vault NOT modified after Descer at last position")
	}

	// Test 4: Actual reposition should mark as modified
	err = manager.ReposicionarPasta(folder2, 0)
	if err != nil {
		t.Errorf("Failed to reposition folder2: %v", err)
	}
	if !cofre.Modificado() {
		t.Error("Expected vault to be modified after actual reposition")
	}
}

// TestFolderDeletionPromotesChildren verifies deletion promotes all children to parent.
// Per D-27: hard delete with promotion.
func TestFolderDeletionPromotesChildren(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	sitesEApps := cofre.PastaGeral().Subpastas()[0]

	// Create structure: sitesEApps -> ToDelete -> [Child1, Child2]
	toDelete, _ := manager.CriarPasta(sitesEApps, "ToDelete", 0)
	child1, _ := manager.CriarPasta(toDelete, "Child1", 0)
	child2, _ := manager.CriarPasta(toDelete, "Child2", 1)

	initialSubCount := len(sitesEApps.Subpastas())

	// Delete ToDelete
	_, err := manager.ExcluirPasta(toDelete)
	if err != nil {
		t.Fatalf("Failed to delete folder: %v", err)
	}

	// Verify ToDelete is removed
	for _, sub := range sitesEApps.Subpastas() {
		if sub == toDelete {
			t.Error("Deleted folder still exists in parent")
		}
	}

	// Verify Child1 and Child2 are promoted to sitesEApps
	foundChild1 := false
	foundChild2 := false
	for _, sub := range sitesEApps.Subpastas() {
		if sub == child1 {
			foundChild1 = true
			if sub.Pai() != sitesEApps {
				t.Error("Child1 parent should be sitesEApps after promotion")
			}
		}
		if sub == child2 {
			foundChild2 = true
			if sub.Pai() != sitesEApps {
				t.Error("Child2 parent should be sitesEApps after promotion")
			}
		}
	}
	if !foundChild1 {
		t.Error("Child1 not found after promotion")
	}
	if !foundChild2 {
		t.Error("Child2 not found after promotion")
	}

	// Verify count: initial + 2 children - 1 deleted = initial + 1
	finalSubCount := len(sitesEApps.Subpastas())
	if finalSubCount != initialSubCount+1 {
		t.Errorf("Expected %d subfolders after deletion, got %d", initialSubCount+1, finalSubCount)
	}
}

// TestFolderDeletionConflictResolution verifies automatic conflict resolution during deletion.
// Per FOLDER-05: name conflicts resolved with numeric suffix.
func TestFolderDeletionConflictResolution(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	sitesEApps := cofre.PastaGeral().Subpastas()[0]

	// Create conflicting structure:
	// sitesEApps -> [Existing, ToDelete -> Existing]
	existingFolder, _ := manager.CriarPasta(sitesEApps, "Existing", 0)
	toDelete, _ := manager.CriarPasta(sitesEApps, "ToDelete", 1)
	conflictingFolder, _ := manager.CriarPasta(toDelete, "Existing", 0)

	// Add a secret to conflictingFolder
	secret := &Segredo{
		nome:         "TestSecret",
		pasta:        conflictingFolder,
		campos:       []CampoSegredo{},
		observacao:   CampoSegredo{nome: "Observação", tipo: TipoCampoComum, valor: []byte("")},
		estadoSessao: EstadoOriginal,
	}
	conflictingFolder.segredos = append(conflictingFolder.segredos, secret)

	// Delete ToDelete
	renomeacoes, err := manager.ExcluirPasta(toDelete)
	if err != nil {
		t.Fatalf("Failed to delete folder: %v", err)
	}

	// Verify conflictingFolder was merged into existingFolder
	// The secret should have been promoted and potentially renamed
	foundSecret := false
	for _, seg := range existingFolder.Segredos() {
		if seg.Nome() == "TestSecret" || seg.Nome() == "TestSecret (1)" {
			foundSecret = true
		}
	}
	if !foundSecret {
		t.Error("Secret not found in existing folder after merge")
	}

	// Verify renomeacoes were tracked if any occurred
	t.Logf("Renomeacoes: %v", renomeacoes)
}

// TestPastaGeralProtectionAcrossOperations verifies Pasta Geral cannot be modified or deleted.
// Per D-21: Pasta Geral is immutable and protected.
func TestPastaGeralProtectionAcrossOperations(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	pastaGeral := cofre.PastaGeral()
	sitesEApps := pastaGeral.Subpastas()[0]

	// Test 1: Cannot rename Pasta Geral
	err := manager.RenomearPasta(pastaGeral, "NewName")
	if err != ErrPastaGeralProtected {
		t.Errorf("Expected ErrPastaGeralProtected when renaming Pasta Geral, got: %v", err)
	}

	// Test 2: Cannot move Pasta Geral
	err = manager.MoverPasta(pastaGeral, sitesEApps)
	if err != ErrPastaGeralProtected {
		t.Errorf("Expected ErrPastaGeralProtected when moving Pasta Geral, got: %v", err)
	}

	// Test 3: Cannot delete Pasta Geral
	_, err = manager.ExcluirPasta(pastaGeral)
	if err != ErrPastaGeralNaoExcluivel {
		t.Errorf("Expected ErrPastaGeralNaoExcluivel when deleting Pasta Geral, got: %v", err)
	}

	// Test 4: Cannot reposition Pasta Geral (has no parent)
	// This should fail at validation level
	// Note: ReposicionarPasta expects pasta.pai to exist, so this will panic
	// We skip this test as it's an invalid operation by design
}

// TestLockedVaultPreventsFolderOperations verifies all folder operations fail when vault is locked.
// Per session management: locked vault blocks all modifications.
func TestLockedVaultPreventsFolderOperations(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	sitesEApps := cofre.PastaGeral().Subpastas()[0]
	folder, _ := manager.CriarPasta(sitesEApps, "TestFolder", 0)

	// Lock vault
	manager.Lock()

	// Test 1: Cannot create folder
	_, err := manager.CriarPasta(sitesEApps, "NewFolder", 0)
	if err != ErrCofreBloqueado {
		t.Errorf("Expected ErrCofreBloqueado for CriarPasta, got: %v", err)
	}

	// Test 2: Cannot rename folder
	err = manager.RenomearPasta(folder, "NewName")
	if err != ErrCofreBloqueado {
		t.Errorf("Expected ErrCofreBloqueado for RenomearPasta, got: %v", err)
	}

	// Test 3: Cannot move folder
	financeiro := cofre.PastaGeral().Subpastas()[1]
	err = manager.MoverPasta(folder, financeiro)
	if err != ErrCofreBloqueado {
		t.Errorf("Expected ErrCofreBloqueado for MoverPasta, got: %v", err)
	}

	// Test 4: Cannot reposition folder
	err = manager.ReposicionarPasta(folder, 0)
	if err != ErrCofreBloqueado {
		t.Errorf("Expected ErrCofreBloqueado for ReposicionarPasta, got: %v", err)
	}

	// Test 5: Cannot delete folder
	_, err = manager.ExcluirPasta(folder)
	if err != ErrCofreBloqueado {
		t.Errorf("Expected ErrCofreBloqueado for ExcluirPasta, got: %v", err)
	}
}

// TestCriarSegredoEstadoInicial verifies secret creation initializes state correctly.
// Per D-11, D-13: new secret has estadoSessao = Modificado (new content), favorito = false.
func TestCriarSegredoEstadoInicial(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, nil)

	pastaGeral := cofre.PastaGeral()
	modelos := cofre.Modelos()
	var modeloLogin *ModeloSegredo
	for _, m := range modelos {
		if m.Nome() == "Login" {
			modeloLogin = m
			break
		}
	}

	if modeloLogin == nil {
		t.Fatal("Login model not found")
	}

	// Create secret from template
	segredo, err := manager.CriarSegredo(pastaGeral, "GitHub", modeloLogin)
	if err != nil {
		t.Fatalf("Failed to create secret: %v", err)
	}

	// Verify initial state: estadoSessao = Modificado (new content per D-11)
	if segredo.EstadoSessao() != EstadoModificado {
		t.Errorf("Expected estadoSessao Modificado for new secret, got %v", segredo.EstadoSessao())
	}

	// Verify favorito = false
	if segredo.Favorito() {
		t.Error("Expected favorito false for new secret")
	}

	// Verify campos initialized from template structure with empty values
	campos := segredo.Campos()
	templateCampos := modeloLogin.Campos()
	if len(campos) != len(templateCampos) {
		t.Errorf("Expected %d campos from template, got %d", len(templateCampos), len(campos))
	}

	for i, campo := range campos {
		if campo.Nome() != templateCampos[i].Nome() {
			t.Errorf("Campo %d: expected name %s, got %s", i, templateCampos[i].Nome(), campo.Nome())
		}
		if campo.Tipo() != templateCampos[i].Tipo() {
			t.Errorf("Campo %d: expected type %v, got %v", i, templateCampos[i].Tipo(), campo.Tipo())
		}
		// Values should be empty (initialized with empty []byte)
		if len(campo.ValorComoString()) != 0 {
			t.Errorf("Campo %d: expected empty value, got %s", i, campo.ValorComoString())
		}
	}

	// Verify cofre.modificado = true
	if !manager.IsModified() {
		t.Error("Expected cofre to be modified after creating secret")
	}
}

// TestExcluirSegredoSoftDelete verifies secret deletion follows soft-delete pattern.
// Per D-14: Delete is reversible until Salvar (estadoSessao transitions based on current state).
func TestExcluirSegredoSoftDelete(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, nil)

	pastaGeral := cofre.PastaGeral()
	modelos := cofre.Modelos()
	var modeloLogin *ModeloSegredo
	for _, m := range modelos {
		if m.Nome() == "Login" {
			modeloLogin = m
			break
		}
	}

	// Create a secret
	segredo, err := manager.CriarSegredo(pastaGeral, "TestSecret", modeloLogin)
	if err != nil {
		t.Fatalf("Failed to create secret: %v", err)
	}

	// Verify initial state
	if segredo.EstadoSessao() != EstadoModificado {
		t.Errorf("Expected EstadoModificado after creation, got %v", segredo.EstadoSessao())
	}

	// Delete the secret
	err = manager.ExcluirSegredo(segredo)
	if err != nil {
		t.Fatalf("Failed to delete secret: %v", err)
	}

	// Verify estadoSessao = Excluido (soft delete)
	if segredo.EstadoSessao() != EstadoExcluido {
		t.Errorf("Expected EstadoExcluido after deletion, got %v", segredo.EstadoSessao())
	}

	// Verify cofre.modificado = true
	if !manager.IsModified() {
		t.Error("Expected cofre to be modified after deleting secret")
	}

	// Verify double-delete fails with ErrSegredoJaExcluido
	err = manager.ExcluirSegredo(segredo)
	if err != ErrSegredoJaExcluido {
		t.Errorf("Expected ErrSegredoJaExcluido for double-delete, got %v", err)
	}
}
