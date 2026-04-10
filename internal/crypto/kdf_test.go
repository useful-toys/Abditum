package crypto_test

import (
	"bytes"
	"testing"

	"github.com/useful-toys/abditum/internal/crypto"
)

// TestGenerateSalt verifies that GenerateSalt produces a 32-byte slice
// and that calling it twice produces distinct salts (uniqueness).
func TestGenerateSalt(t *testing.T) {
	// Generate first salt
	salt1, err1 := crypto.GenerateSalt()
	if err1 != nil {
		t.Fatalf("GenerateSalt() returned unexpected error: %v", err1)
	}

	if len(salt1) != 32 {
		t.Errorf("GenerateSalt() returned salt of length %d, want 32", len(salt1))
	}

	// Generate second salt
	salt2, err2 := crypto.GenerateSalt()
	if err2 != nil {
		t.Fatalf("GenerateSalt() second call returned unexpected error: %v", err2)
	}

	if len(salt2) != 32 {
		t.Errorf("GenerateSalt() second call returned salt of length %d, want 32", len(salt2))
	}

	// Verify uniqueness - two salts should not be identical
	if bytes.Equal(salt1, salt2) {
		t.Error("GenerateSalt() produced identical salts on two calls - nonces must be unique")
	}
}

// TestDeriveKey verifies that DeriveKey produces a 32-byte key with valid inputs.
func TestDeriveKey(t *testing.T) {
	password := []byte("TestPassword123!")
	salt := make([]byte, 32)
	// Use a fixed salt for deterministic testing
	for i := range salt {
		salt[i] = byte(i)
	}

	params := crypto.ArgonParams{
		Time:    3,
		Memory:  262144, // 256 MiB in KiB
		Threads: 4,
		KeyLen:  32,
	}

	key, err := crypto.DeriveKey(password, salt, params)
	if err != nil {
		t.Fatalf("DeriveKey() returned unexpected error: %v", err)
	}

	if len(key) != 32 {
		t.Errorf("DeriveKey() returned key of length %d, want 32", len(key))
	}
}

// TestDeriveKeyInvalidPassword verifies that DeriveKey returns ErrInvalidParams
// when called with a nil password.
func TestDeriveKeyInvalidPassword(t *testing.T) {
	salt := make([]byte, 32)
	params := crypto.ArgonParams{
		Time:    3,
		Memory:  262144,
		Threads: 4,
		KeyLen:  32,
	}

	_, err := crypto.DeriveKey(nil, salt, params)
	if err != crypto.ErrInvalidParams {
		t.Errorf("DeriveKey(nil password) returned %v, want ErrInvalidParams", err)
	}
}

// TestDeriveKeyInvalidSalt verifies that DeriveKey returns ErrInvalidParams
// when called with a zero-length salt.
func TestDeriveKeyInvalidSalt(t *testing.T) {
	password := []byte("TestPassword123!")
	salt := []byte{} // Empty slice
	params := crypto.ArgonParams{
		Time:    3,
		Memory:  262144,
		Threads: 4,
		KeyLen:  32,
	}

	_, err := crypto.DeriveKey(password, salt, params)
	if err != crypto.ErrInvalidParams {
		t.Errorf("DeriveKey(empty salt) returned %v, want ErrInvalidParams", err)
	}
}

// BenchmarkDeriveKey measures Argon2id performance with standard parameters.
// Expected: 200-500ms per operation on modern CPU.
func BenchmarkDeriveKey(b *testing.B) {
	password := []byte("BenchmarkPassword123!")
	salt := make([]byte, 32)
	params := crypto.ArgonParams{
		Time:    3,
		Memory:  262144, // 256 MiB in KiB - CRITICAL: not 256
		Threads: 4,
		KeyLen:  32,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := crypto.DeriveKey(password, salt, params)
		if err != nil {
			b.Fatalf("DeriveKey failed: %v", err)
		}
	}
}
