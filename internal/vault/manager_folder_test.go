package vault

import (
	"errors"
	"testing"
)

// Task 1 Tests: CriarPasta

func TestCriarPasta_Success(t *testing.T) {
	cofre := NovoCofre()
	manager := NewManager(cofre, &mockRepository{})

	pasta, err := manager.CriarPasta(cofre.PastaGeral(), "Nova Pasta", 0)
	if err != nil {
		t.Fatalf("CriarPasta failed: %v", err)
	}

	if pasta == nil {
		t.Fatal("CriarPasta returned nil pasta")
	}

	if pasta.Nome() != "Nova Pasta" {
		t.Errorf("Expected nome='Nova Pasta', got %q", pasta.Nome())
	}

	if pasta.Pai() != cofre.PastaGeral() {
		t.Error("Pasta pai should be Pasta Geral")
	}

	if !manager.IsModified() {
		t.Error("Vault should be marked modified after creating folder")
	}

	// Verify pasta was added to parent
	subpastas := cofre.PastaGeral().Subpastas()
	if len(subpastas) == 0 {
		t.Fatal("Pasta Geral should have subfolders")
	}

	found := false
	for _, sub := range subpastas {
		if sub == pasta {
			found = true
			break
		}
	}
	if !found {
		t.Error("Created pasta not found in parent subpastas")
	}
}

func TestCriarPasta_AtPosition(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao() // Creates "Sites e Apps" and "Financeiro"
	manager := NewManager(cofre, &mockRepository{})

	// Insert at position 1 (between Sites e Apps and Financeiro)
	pasta, err := manager.CriarPasta(cofre.PastaGeral(), "Meio", 1)
	if err != nil {
		t.Fatalf("CriarPasta at position 1 failed: %v", err)
	}

	subpastas := cofre.PastaGeral().Subpastas()
	if len(subpastas) != 3 {
		t.Fatalf("Expected 3 subpastas, got %d", len(subpastas))
	}

	if subpastas[1] != pasta {
		t.Error("Pasta should be at position 1")
	}
}

func TestCriarPasta_AtEnd(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	subpastasBefore := cofre.PastaGeral().Subpastas()
	count := len(subpastasBefore)

	// Append at end (position == len)
	pasta, err := manager.CriarPasta(cofre.PastaGeral(), "Ultimo", count)
	if err != nil {
		t.Fatalf("CriarPasta at end failed: %v", err)
	}

	subpastas := cofre.PastaGeral().Subpastas()
	if len(subpastas) != count+1 {
		t.Fatalf("Expected %d subpastas, got %d", count+1, len(subpastas))
	}

	if subpastas[count] != pasta {
		t.Error("Pasta should be at last position")
	}
}

func TestCriarPasta_NomeVazio(t *testing.T) {
	cofre := NovoCofre()
	manager := NewManager(cofre, &mockRepository{})

	_, err := manager.CriarPasta(cofre.PastaGeral(), "", 0)
	if !errors.Is(err, ErrNomeVazio) {
		t.Errorf("Expected ErrNomeVazio, got %v", err)
	}

	if manager.IsModified() {
		t.Error("Vault should not be modified after validation failure")
	}
}

func TestCriarPasta_NomeMuitoLongo(t *testing.T) {
	cofre := NovoCofre()
	manager := NewManager(cofre, &mockRepository{})

	// 256 characters
	nomeLongo := string(make([]byte, 256))
	for i := range nomeLongo {
		nomeLongo = nomeLongo[:i] + "a"
	}

	_, err := manager.CriarPasta(cofre.PastaGeral(), nomeLongo, 0)
	if !errors.Is(err, ErrNomeMuitoLongo) {
		t.Errorf("Expected ErrNomeMuitoLongo, got %v", err)
	}

	if manager.IsModified() {
		t.Error("Vault should not be modified after validation failure")
	}
}

func TestCriarPasta_NomeConflict(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao() // Creates "Sites e Apps"
	manager := NewManager(cofre, &mockRepository{})

	_, err := manager.CriarPasta(cofre.PastaGeral(), "Sites e Apps", 0)
	if !errors.Is(err, ErrNameConflict) {
		t.Errorf("Expected ErrNameConflict, got %v", err)
	}

	if manager.IsModified() {
		t.Error("Vault should not be modified after validation failure")
	}
}

func TestCriarPasta_PosicaoInvalida(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao() // 2 folders
	manager := NewManager(cofre, &mockRepository{})

	tests := []struct {
		nome     string
		posicao  int
		expected error
	}{
		{"negative", -1, ErrPosicaoInvalida},
		{"beyond end", 3, ErrPosicaoInvalida}, // len=2, max valid=2 (append)
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			_, err := manager.CriarPasta(cofre.PastaGeral(), "Test", tt.posicao)
			if !errors.Is(err, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, err)
			}

			if manager.IsModified() {
				t.Error("Vault should not be modified after validation failure")
			}
		})
	}
}

func TestCriarPasta_WhenLocked(t *testing.T) {
	cofre := NovoCofre()
	manager := NewManager(cofre, &mockRepository{})
	manager.Lock()

	_, err := manager.CriarPasta(cofre.PastaGeral(), "Test", 0)
	if !errors.Is(err, ErrCofreBloqueado) {
		t.Errorf("Expected ErrCofreBloqueado, got %v", err)
	}
}

// Task 2 Tests: RenomearPasta

func TestRenomearPasta_Success(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	pasta := cofre.PastaGeral().Subpastas()[0] // "Sites e Apps"
	err := manager.RenomearPasta(pasta, "Novo Nome")
	if err != nil {
		t.Fatalf("RenomearPasta failed: %v", err)
	}

	if pasta.Nome() != "Novo Nome" {
		t.Errorf("Expected nome='Novo Nome', got %q", pasta.Nome())
	}

	if !manager.IsModified() {
		t.Error("Vault should be marked modified after renaming folder")
	}
}

func TestRenomearPasta_PastaGeralProtection(t *testing.T) {
	cofre := NovoCofre()
	manager := NewManager(cofre, &mockRepository{})

	err := manager.RenomearPasta(cofre.PastaGeral(), "Tentativa")
	if !errors.Is(err, ErrPastaGeralProtected) {
		t.Errorf("Expected ErrPastaGeralProtected, got %v", err)
	}

	if manager.IsModified() {
		t.Error("Vault should not be modified after failed rename")
	}

	if cofre.PastaGeral().Nome() != "Pasta Geral" {
		t.Error("Pasta Geral name should remain unchanged")
	}
}

func TestRenomearPasta_NoOpSameName(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	pasta := cofre.PastaGeral().Subpastas()[0] // "Sites e Apps"
	nomeOriginal := pasta.Nome()

	// Rename to same name (no-op per D-12)
	err := manager.RenomearPasta(pasta, nomeOriginal)
	if err != nil {
		t.Fatalf("RenomearPasta with same name should not fail: %v", err)
	}

	// Should NOT mark vault as modified (D-12 change detection)
	if manager.IsModified() {
		t.Error("Vault should not be modified when renaming to same name (D-12)")
	}
}

func TestRenomearPasta_NameConflict(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao() // Creates "Sites e Apps" and "Financeiro"
	manager := NewManager(cofre, &mockRepository{})

	subpastas := cofre.PastaGeral().Subpastas()
	pasta1 := subpastas[0] // "Sites e Apps"

	// Try to rename to existing sibling name
	err := manager.RenomearPasta(pasta1, "Financeiro")
	if !errors.Is(err, ErrNameConflict) {
		t.Errorf("Expected ErrNameConflict, got %v", err)
	}

	if manager.IsModified() {
		t.Error("Vault should not be modified after failed rename")
	}

	if pasta1.Nome() != "Sites e Apps" {
		t.Error("Pasta name should remain unchanged after conflict")
	}
}

func TestRenomearPasta_NomeVazio(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	pasta := cofre.PastaGeral().Subpastas()[0]
	err := manager.RenomearPasta(pasta, "")
	if !errors.Is(err, ErrNomeVazio) {
		t.Errorf("Expected ErrNomeVazio, got %v", err)
	}

	if manager.IsModified() {
		t.Error("Vault should not be modified after validation failure")
	}
}

func TestRenomearPasta_NomeMuitoLongo(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	pasta := cofre.PastaGeral().Subpastas()[0]
	nomeLongo := string(make([]byte, 256))
	for i := range nomeLongo {
		nomeLongo = nomeLongo[:i] + "a"
	}

	err := manager.RenomearPasta(pasta, nomeLongo)
	if !errors.Is(err, ErrNomeMuitoLongo) {
		t.Errorf("Expected ErrNomeMuitoLongo, got %v", err)
	}

	if manager.IsModified() {
		t.Error("Vault should not be modified after validation failure")
	}
}

func TestRenomearPasta_WhenLocked(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})
	pasta := cofre.PastaGeral().Subpastas()[0]
	manager.Lock()

	err := manager.RenomearPasta(pasta, "Test")
	if !errors.Is(err, ErrCofreBloqueado) {
		t.Errorf("Expected ErrCofreBloqueado, got %v", err)
	}
}

// Task 3 Tests: MoverPasta

func TestMoverPasta_Success(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao() // Creates "Sites e Apps" and "Financeiro"
	manager := NewManager(cofre, &mockRepository{})

	// Create a subfolder in Sites e Apps
	sitesEApps := cofre.PastaGeral().Subpastas()[0]
	subpasta, _ := manager.CriarPasta(sitesEApps, "Subfolder", 0)

	// Move subfolder to Financeiro
	financeiro := cofre.PastaGeral().Subpastas()[1]
	err := manager.MoverPasta(subpasta, financeiro)
	if err != nil {
		t.Fatalf("MoverPasta failed: %v", err)
	}

	// Verify parent changed
	if subpasta.Pai() != financeiro {
		t.Error("Pasta pai should be Financeiro after move")
	}

	// Verify removed from old parent
	if len(sitesEApps.Subpastas()) != 0 {
		t.Error("Sites e Apps should have no subfolders after move")
	}

	// Verify added to new parent
	financeiroSubs := financeiro.Subpastas()
	if len(financeiroSubs) != 1 {
		t.Fatalf("Financeiro should have 1 subfolder, got %d", len(financeiroSubs))
	}
	if financeiroSubs[0] != subpasta {
		t.Error("Moved subfolder not found in Financeiro")
	}

	if !manager.IsModified() {
		t.Error("Vault should be marked modified after move")
	}
}

func TestMoverPasta_PastaGeralProtection(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	destino := cofre.PastaGeral().Subpastas()[0]
	err := manager.MoverPasta(cofre.PastaGeral(), destino)
	if !errors.Is(err, ErrPastaGeralProtected) {
		t.Errorf("Expected ErrPastaGeralProtected, got %v", err)
	}

	if manager.IsModified() {
		t.Error("Vault should not be modified after failed move")
	}
}

func TestMoverPasta_MoveToSelf(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	pasta := cofre.PastaGeral().Subpastas()[0]
	err := manager.MoverPasta(pasta, pasta)
	if !errors.Is(err, ErrDestinoInvalido) {
		t.Errorf("Expected ErrDestinoInvalido, got %v", err)
	}

	if manager.IsModified() {
		t.Error("Vault should not be modified after failed move")
	}
}

func TestMoverPasta_CycleDetectionDirectChild(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	// Create hierarchy: Sites e Apps -> SubA
	sitesEApps := cofre.PastaGeral().Subpastas()[0]
	subA, _ := manager.CriarPasta(sitesEApps, "SubA", 0)

	// Save to clear modified flag
	cofre.modificado = false

	// Try to move Sites e Apps into its own child SubA (creates cycle)
	err := manager.MoverPasta(sitesEApps, subA)
	if !errors.Is(err, ErrCycleDetected) {
		t.Errorf("Expected ErrCycleDetected, got %v", err)
	}

	if manager.IsModified() {
		t.Error("Vault should not be modified after cycle detection")
	}

	// Verify hierarchy unchanged
	if subA.Pai() != sitesEApps {
		t.Error("SubA parent should still be Sites e Apps")
	}
	if sitesEApps.Pai() != cofre.PastaGeral() {
		t.Error("Sites e Apps parent should still be Pasta Geral")
	}
}

func TestMoverPasta_CycleDetectionGrandchild(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	// Create hierarchy: Sites e Apps -> SubA -> SubB
	sitesEApps := cofre.PastaGeral().Subpastas()[0]
	subA, _ := manager.CriarPasta(sitesEApps, "SubA", 0)
	subB, _ := manager.CriarPasta(subA, "SubB", 0)

	// Save to clear modified flag
	cofre.modificado = false

	// Try to move Sites e Apps into grandchild SubB (creates cycle)
	err := manager.MoverPasta(sitesEApps, subB)
	if !errors.Is(err, ErrCycleDetected) {
		t.Errorf("Expected ErrCycleDetected, got %v", err)
	}

	if manager.IsModified() {
		t.Error("Vault should not be modified after cycle detection")
	}
}

func TestMoverPasta_NameConflict(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	// Create "X" in Sites e Apps and "X" in Financeiro
	sitesEApps := cofre.PastaGeral().Subpastas()[0]
	financeiro := cofre.PastaGeral().Subpastas()[1]

	pastaX1, _ := manager.CriarPasta(sitesEApps, "X", 0)
	manager.CriarPasta(financeiro, "X", 0)

	// Save to clear modified flag
	cofre.modificado = false

	// Try to move X from Sites e Apps to Financeiro (name conflict)
	err := manager.MoverPasta(pastaX1, financeiro)
	if !errors.Is(err, ErrNameConflict) {
		t.Errorf("Expected ErrNameConflict, got %v", err)
	}

	if manager.IsModified() {
		t.Error("Vault should not be modified after failed move")
	}

	// Verify pastaX1 still in Sites e Apps
	if pastaX1.Pai() != sitesEApps {
		t.Error("Pasta should still be in Sites e Apps after failed move")
	}
}

func TestMoverPasta_MoveToSibling(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	// Move Sites e Apps to Financeiro (siblings, no cycle)
	sitesEApps := cofre.PastaGeral().Subpastas()[0]
	financeiro := cofre.PastaGeral().Subpastas()[1]

	err := manager.MoverPasta(sitesEApps, financeiro)
	if err != nil {
		t.Fatalf("MoverPasta to sibling should succeed: %v", err)
	}

	if sitesEApps.Pai() != financeiro {
		t.Error("Sites e Apps should be child of Financeiro")
	}

	if !manager.IsModified() {
		t.Error("Vault should be marked modified")
	}
}

func TestMoverPasta_WhenLocked(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()
	manager := NewManager(cofre, &mockRepository{})

	sitesEApps := cofre.PastaGeral().Subpastas()[0]
	financeiro := cofre.PastaGeral().Subpastas()[1]
	manager.Lock()

	err := manager.MoverPasta(sitesEApps, financeiro)
	if !errors.Is(err, ErrCofreBloqueado) {
		t.Errorf("Expected ErrCofreBloqueado, got %v", err)
	}
}
