package vault

import (
	"testing"
	"time"
)

// TestRenomearSegredo tests secret renaming with estadoSessao tracking (UAT-09).
func TestRenomearSegredo(t *testing.T) {
	// Setup: Create vault with folder and secret
	cofre := NovoCofre()
	err := cofre.InicializarConteudoPadrao()
	if err != nil {
		t.Fatalf("Failed to initialize default content: %v", err)
	}

	manager := NewManager(cofre, nil)

	pasta, err2 := manager.CriarPasta(cofre.PastaGeral(), "Trabalho", -1)
	if err2 != nil {
		t.Fatalf("Failed to create folder: %v", err2)
	}

	// Create a secret with Original state
	segredo := &Segredo{
		nome: "GitHub",
		campos: []CampoSegredo{
			{nome: "Usuário", tipo: TipoCampoComum, valor: []byte("alice")},
			{nome: "Senha", tipo: TipoCampoSensivel, valor: []byte("secret123")},
		},
		observacao: CampoSegredo{
			nome:  "Observação",
			tipo:  TipoCampoComum,
			valor: []byte("Corporate account"),
		},
		pasta:                 pasta,
		favorito:              false,
		estadoSessao:          EstadoOriginal,
		dataCriacao:           time.Now().UTC().Add(-24 * time.Hour),
		dataUltimaModificacao: time.Now().UTC().Add(-24 * time.Hour),
	}
	pasta.segredos = append(pasta.segredos, segredo)

	// Record original modification time
	oldModTime := segredo.dataUltimaModificacao
	time.Sleep(2 * time.Millisecond) // Ensure time difference

	t.Run("ValidRename", func(t *testing.T) {
		err := manager.RenomearSegredo(segredo, "GitHub-Work")
		if err != nil {
			t.Errorf("RenomearSegredo failed: %v", err)
		}

		// Verify name changed
		if segredo.Nome() != "GitHub-Work" {
			t.Errorf("Expected name 'GitHub-Work', got %q", segredo.Nome())
		}

		// Verify estadoSessao = Modificado (D-11)
		if segredo.EstadoSessao() != EstadoModificado {
			t.Errorf("Expected EstadoModificado, got %v", segredo.EstadoSessao())
		}

		// Verify timestamp updated
		if !segredo.DataUltimaModificacao().After(oldModTime) {
			t.Errorf("Expected timestamp to be updated")
		}

		// Verify cofre.modificado flag set
		if !cofre.modificado {
			t.Errorf("Expected cofre.modificado = true")
		}
	})

	t.Run("NoOpRename", func(t *testing.T) {
		// Reset modificado flag
		cofre.modificado = false
		oldModTime := segredo.dataUltimaModificacao
		time.Sleep(2 * time.Millisecond)

		// Rename to same name (D-12: no-op should not mark modified)
		err := manager.RenomearSegredo(segredo, "GitHub-Work")
		if err != nil {
			t.Errorf("RenomearSegredo failed: %v", err)
		}

		// Verify cofre NOT marked modified
		if cofre.modificado {
			t.Errorf("Expected cofre.modificado = false for no-op rename")
		}

		// Verify timestamp NOT updated
		if segredo.DataUltimaModificacao() != oldModTime {
			t.Errorf("Expected timestamp unchanged for no-op rename")
		}
	})

	t.Run("ConflictDetection", func(t *testing.T) {
		// Add another secret with conflicting name
		conflito := &Segredo{
			nome: "Bitbucket",
			campos: []CampoSegredo{
				{nome: "Usuário", tipo: TipoCampoComum, valor: []byte("bob")},
			},
			observacao: CampoSegredo{
				nome:  "Observação",
				tipo:  TipoCampoComum,
				valor: []byte(""),
			},
			pasta:                 pasta,
			estadoSessao:          EstadoOriginal,
			dataCriacao:           time.Now().UTC(),
			dataUltimaModificacao: time.Now().UTC(),
		}
		pasta.segredos = append(pasta.segredos, conflito)

		// Try to rename to existing name
		err := manager.RenomearSegredo(segredo, "Bitbucket")
		if err != ErrNameConflict {
			t.Errorf("Expected ErrNameConflict, got %v", err)
		}

		// Verify name unchanged
		if segredo.Nome() != "GitHub-Work" {
			t.Errorf("Expected name unchanged after failed rename")
		}
	})

	t.Run("EmptyName", func(t *testing.T) {
		err := manager.RenomearSegredo(segredo, "")
		if err != ErrNomeVazio {
			t.Errorf("Expected ErrNomeVazio, got %v", err)
		}
	})

	t.Run("NameTooLong", func(t *testing.T) {
		longName := string(make([]byte, 256))
		err := manager.RenomearSegredo(segredo, longName)
		if err != ErrNomeMuitoLongo {
			t.Errorf("Expected ErrNomeMuitoLongo, got %v", err)
		}
	})

	t.Run("PreservesEstadoIncluido", func(t *testing.T) {
		// Create new secret with EstadoIncluido
		novo := &Segredo{
			nome: "NewSecret",
			campos: []CampoSegredo{
				{nome: "Senha", tipo: TipoCampoSensivel, valor: []byte("test")},
			},
			observacao: CampoSegredo{
				nome:  "Observação",
				tipo:  TipoCampoComum,
				valor: []byte(""),
			},
			pasta:                 pasta,
			estadoSessao:          EstadoIncluido,
			dataCriacao:           time.Now().UTC(),
			dataUltimaModificacao: time.Now().UTC(),
		}
		pasta.segredos = append(pasta.segredos, novo)

		// Rename it
		err := manager.RenomearSegredo(novo, "RenamedNew")
		if err != nil {
			t.Errorf("RenomearSegredo failed: %v", err)
		}

		// Verify estado remains Incluido (D-11: only Original -> Modificado)
		if novo.EstadoSessao() != EstadoIncluido {
			t.Errorf("Expected EstadoIncluido preserved, got %v", novo.EstadoSessao())
		}
	})

	t.Run("LockedVault", func(t *testing.T) {
		manager.Lock()
		err := manager.RenomearSegredo(segredo, "ShouldFail")
		if err != ErrCofreBloqueado {
			t.Errorf("Expected ErrCofreBloqueado, got %v", err)
		}
	})
}
