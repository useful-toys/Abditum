package vault_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/useful-toys/abditum/internal/vault"
)

func TestErrModifiedExternally_IsComparavel(t *testing.T) {
	wrapped := fmt.Errorf("vault.Salvar: %w", vault.ErrModifiedExternally)
	if !errors.Is(wrapped, vault.ErrModifiedExternally) {
		t.Error("errors.Is deveria reconhecer ErrModifiedExternally em erro encadeado")
	}
}
