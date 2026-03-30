package vault

import (
	"errors"
	"testing"
)

// Mock repository for testing
type mockRepository struct {
	salvarCalled bool
	salvarError  error
}

func (m *mockRepository) Salvar(cofre *Cofre) error {
	m.salvarCalled = true
	return m.salvarError
}

func (m *mockRepository) Carregar() (*Cofre, error) {
	return nil, errors.New("not implemented")
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
