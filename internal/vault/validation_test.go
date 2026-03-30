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
