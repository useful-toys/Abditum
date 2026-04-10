package vault

import (
	"testing"
)

func TestNovoCofre(t *testing.T) {
	cofre := NovoCofre()

	if cofre == nil {
		t.Fatal("NovoCofre() returned nil")
	}

	if cofre.PastaGeral() == nil {
		t.Error("Pasta Geral is nil")
	}

	if cofre.PastaGeral().Nome() != "Pasta Geral" {
		t.Errorf("Expected 'Pasta Geral', got '%s'", cofre.PastaGeral().Nome())
	}

	if len(cofre.PastaGeral().Subpastas()) != 0 {
		t.Errorf("Expected empty subpastas, got %d", len(cofre.PastaGeral().Subpastas()))
	}

	if len(cofre.Modelos()) != 0 {
		t.Errorf("Expected no templates, got %d", len(cofre.Modelos()))
	}

	if cofre.Modificado() {
		t.Error("New cofre should not be marked modified")
	}

	// Verify default Configuracoes
	config := cofre.Configuracoes()
	if config.tempoBloqueioInatividadeMinutos != 5 {
		t.Errorf("Expected tempoBloqueio=5, got %d", config.tempoBloqueioInatividadeMinutos)
	}
	if config.tempoOcultarSegredoSegundos != 15 {
		t.Errorf("Expected tempoOcultar=15, got %d", config.tempoOcultarSegredoSegundos)
	}
	if config.tempoLimparAreaTransferenciaSegundos != 30 {
		t.Errorf("Expected tempoLimpar=30, got %d", config.tempoLimparAreaTransferenciaSegundos)
	}
}

func TestInicializarConteudoPadrao(t *testing.T) {
	cofre := NovoCofre()
	err := cofre.InicializarConteudoPadrao()

	if err != nil {
		t.Fatalf("InicializarConteudoPadrao() returned error: %v", err)
	}

	// Verify folders
	subpastas := cofre.PastaGeral().Subpastas()
	if len(subpastas) != 2 {
		t.Fatalf("Expected 2 subfolders, got %d", len(subpastas))
	}

	if subpastas[0].Nome() != "Sites e Apps" {
		t.Errorf("Expected 'Sites e Apps', got '%s'", subpastas[0].Nome())
	}

	if subpastas[1].Nome() != "Financeiro" {
		t.Errorf("Expected 'Financeiro', got '%s'", subpastas[1].Nome())
	}

	// Verify templates
	modelos := cofre.Modelos()
	if len(modelos) != 3 {
		t.Fatalf("Expected 3 templates, got %d", len(modelos))
	}

	// Check template names (already sorted by Modelos())
	esperados := []string{"Cartão de Crédito", "Chave de API", "Login"}
	for i, nome := range esperados {
		if modelos[i].Nome() != nome {
			t.Errorf("Expected template[%d]='%s', got '%s'", i, nome, modelos[i].Nome())
		}
	}

	// Verify Login template structure
	login := encontrarModelo(modelos, "Login")
	if login == nil {
		t.Fatal("Login template not found")
	}

	campos := login.Campos()
	if len(campos) != 3 {
		t.Fatalf("Login should have 3 fields, got %d", len(campos))
	}

	if campos[0].Nome() != "URL" || campos[0].Tipo() != TipoCampoComum {
		t.Error("Login field 0 should be URL (comum)")
	}
	if campos[1].Nome() != "Usuário" || campos[1].Tipo() != TipoCampoComum {
		t.Error("Login field 1 should be Usuário (comum)")
	}
	if campos[2].Nome() != "Senha" || campos[2].Tipo() != TipoCampoSensivel {
		t.Error("Login field 2 should be Senha (sensivel)")
	}

	// Verify modificado flag NOT set (D-28b)
	if cofre.Modificado() {
		t.Error("InicializarConteudoPadrao() should NOT set modificado=true")
	}
}

func TestDefensiveCopySubpastas(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()

	// Get subpastas slice
	subpastas1 := cofre.PastaGeral().Subpastas()
	subpastas2 := cofre.PastaGeral().Subpastas()

	// Verify defensive copy: different slice addresses
	if &subpastas1[0] == &subpastas2[0] {
		t.Error("Subpastas() should return defensive copy, not same slice")
	}

	// Modify returned slice
	subpastas1 = append(subpastas1, &Pasta{nome: "hacker"})

	// Verify entity not affected
	if len(cofre.PastaGeral().Subpastas()) != 2 {
		t.Error("Modifying returned slice affected entity internal state")
	}
}

func TestDefensiveCopyModelos(t *testing.T) {
	cofre := NovoCofre()
	cofre.InicializarConteudoPadrao()

	modelos1 := cofre.Modelos()
	modelos2 := cofre.Modelos()

	// Verify different slices
	if &modelos1[0] == &modelos2[0] {
		t.Error("Modelos() should return defensive copy")
	}

	// Modify returned slice
	modelos1 = append(modelos1, &ModeloSegredo{nome: "hacker"})

	// Verify entity not affected
	if len(cofre.Modelos()) != 3 {
		t.Error("Modifying returned slice affected entity")
	}
}

func TestModelosAlphabeticalSort(t *testing.T) {
	// Create cofre with templates in non-alphabetical order
	cofre := NovoCofre()
	cofre.modelos = []*ModeloSegredo{
		{nome: "Z-Template"},
		{nome: "A-Template"},
		{nome: "M-Template"},
	}

	// Get via public getter
	modelos := cofre.Modelos()

	// Verify alphabetical order (TPL-06)
	if modelos[0].Nome() != "A-Template" {
		t.Errorf("Expected A-Template first, got %s", modelos[0].Nome())
	}
	if modelos[1].Nome() != "M-Template" {
		t.Errorf("Expected M-Template second, got %s", modelos[1].Nome())
	}
	if modelos[2].Nome() != "Z-Template" {
		t.Errorf("Expected Z-Template third, got %s", modelos[2].Nome())
	}
}

// Helper function
func encontrarModelo(modelos []*ModeloSegredo, nome string) *ModeloSegredo {
	for _, m := range modelos {
		if m.Nome() == nome {
			return m
		}
	}
	return nil
}
