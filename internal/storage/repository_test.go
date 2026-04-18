package storage_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/useful-toys/abditum/internal/crypto"
	"github.com/useful-toys/abditum/internal/storage"
	"github.com/useful-toys/abditum/internal/vault"
)

func TestNewFileRepositoryForOpen_Carregar_E_Salvar_Atomico(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.abditum")
	cofre := vault.NovoCofre()
	if err := cofre.InicializarConteudoPadrao(); err != nil {
		t.Fatal(err)
	}
	password := []byte("SenhaForte123!")
	if err := storage.SaveNew(path, cofre, password); err != nil {
		t.Fatal(err)
	}

	repo := storage.NewFileRepositoryForOpen(path, password)

	loaded, err := repo.Carregar()
	if err != nil {
		t.Fatalf("Carregar: %v", err)
	}
	if loaded == nil {
		t.Fatal("Carregar retornou cofre nil")
	}

	if err := repo.Salvar(loaded); err != nil {
		t.Fatalf("Salvar: %v", err)
	}
	bakPath := path + ".bak"
	if _, err := os.Stat(bakPath); os.IsNotExist(err) {
		t.Error("Salvar após Carregar via ForOpen: .bak não foi criado — protocolo atômico não usado")
	}
}

func TestNewFileRepositoryForOpen_SenhaErrada_ErroDeAutenticacao(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.abditum")
	cofre := vault.NovoCofre()
	if err := cofre.InicializarConteudoPadrao(); err != nil {
		t.Fatal(err)
	}
	if err := storage.SaveNew(path, cofre, []byte("SenhaCorreta1!")); err != nil {
		t.Fatal(err)
	}

	repo := storage.NewFileRepositoryForOpen(path, []byte("SenhaErrada1!"))
	_, err := repo.Carregar()
	if !errors.Is(err, crypto.ErrAuthFailed) {
		t.Errorf("senha errada: esperado ErrAuthFailed, obteve %v", err)
	}
}
