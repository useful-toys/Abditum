package crypto

import (
	"errors"
	"testing"
)

type errReader struct{}

func (errReader) Read(_ []byte) (int, error) {
	return 0, errors.New("entropy unavailable")
}

func withEntropyFailure(t *testing.T) {
	t.Helper()

	original := entropyReader
	entropyReader = errReader{}
	t.Cleanup(func() {
		entropyReader = original
	})
}

func TestEncryptEntropyFailure(t *testing.T) {
	withEntropyFailure(t)

	_, err := Encrypt(make([]byte, 32), []byte("secret"))
	if err != ErrInsufficientEntropy {
		t.Fatalf("Encrypt() error = %v, want ErrInsufficientEntropy", err)
	}
}

func TestEncryptWithAADEntropyFailure(t *testing.T) {
	withEntropyFailure(t)

	_, _, err := EncryptWithAAD(make([]byte, 32), []byte("secret"), []byte("aad"))
	if err != ErrInsufficientEntropy {
		t.Fatalf("EncryptWithAAD() error = %v, want ErrInsufficientEntropy", err)
	}
}

func TestGenerateSaltEntropyFailure(t *testing.T) {
	withEntropyFailure(t)

	_, err := GenerateSalt()
	if err != ErrInsufficientEntropy {
		t.Fatalf("GenerateSalt() error = %v, want ErrInsufficientEntropy", err)
	}
}
