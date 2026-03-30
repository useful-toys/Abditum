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
