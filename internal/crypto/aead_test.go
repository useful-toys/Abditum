package crypto_test

import (
	"bytes"
	"testing"

	"github.com/useful-toys/abditum/internal/crypto"
)

// TestEncrypt verifies that Encrypt produces ciphertext with correct format:
// nonce (12) + ciphertext (len(plaintext)) + tag (16).
func TestEncrypt(t *testing.T) {
	key := make([]byte, 32)
	plaintext := []byte("test plaintext")

	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt() returned unexpected error: %v", err)
	}

	// Output should be: nonce (12) + encrypted(plaintext) + tag (16)
	expectedLen := 12 + len(plaintext) + 16
	if len(ciphertext) != expectedLen {
		t.Errorf("Encrypt() returned ciphertext of length %d, want %d", len(ciphertext), expectedLen)
	}
}

// TestNonceUniqueness is CRITICAL for GCM security.
// Encrypting the same plaintext twice with the same key MUST produce
// distinct ciphertexts (proves nonce is regenerated for each call).
func TestNonceUniqueness(t *testing.T) {
	key := make([]byte, 32)
	plaintext := []byte("same plaintext")

	ciphertext1, err1 := crypto.Encrypt(key, plaintext)
	if err1 != nil {
		t.Fatalf("Encrypt() first call returned error: %v", err1)
	}

	ciphertext2, err2 := crypto.Encrypt(key, plaintext)
	if err2 != nil {
		t.Fatalf("Encrypt() second call returned error: %v", err2)
	}

	if bytes.Equal(ciphertext1, ciphertext2) {
		t.Error("Encrypt() produced identical ciphertexts for same plaintext - NONCE REUSE DETECTED")
	}
}

// TestRoundtrip verifies that Encrypt → Decrypt returns original plaintext byte-for-byte.
func TestRoundtrip(t *testing.T) {
	key := make([]byte, 32)
	plaintext := []byte("secret vault data")

	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt() returned error: %v", err)
	}

	decrypted, err := crypto.Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt() returned error: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Roundtrip failed: got %q, want %q", decrypted, plaintext)
	}
}

// TestDecryptWrongKey verifies that Decrypt with wrong key returns ErrAuthFailed.
func TestDecryptWrongKey(t *testing.T) {
	key1 := make([]byte, 32)
	key1[0] = 1
	key2 := make([]byte, 32)
	key2[0] = 2

	plaintext := []byte("secret data")

	ciphertext, err := crypto.Encrypt(key1, plaintext)
	if err != nil {
		t.Fatalf("Encrypt() returned error: %v", err)
	}

	_, err = crypto.Decrypt(key2, ciphertext)
	if err != crypto.ErrAuthFailed {
		t.Errorf("Decrypt(wrong key) returned %v, want ErrAuthFailed", err)
	}
}

// TestDecryptShortCiphertext verifies that Decrypt with ciphertext shorter
// than nonce size returns ErrAuthFailed.
func TestDecryptShortCiphertext(t *testing.T) {
	key := make([]byte, 32)
	shortCiphertext := []byte{1, 2, 3} // Less than 12 bytes

	_, err := crypto.Decrypt(key, shortCiphertext)
	if err != crypto.ErrAuthFailed {
		t.Errorf("Decrypt(short ciphertext) returned %v, want ErrAuthFailed", err)
	}
}

// TestDecryptCorrupted verifies that Decrypt with corrupted ciphertext
// returns ErrAuthFailed (GCM tag verification fails).
func TestDecryptCorrupted(t *testing.T) {
	key := make([]byte, 32)
	plaintext := []byte("original data")

	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt() returned error: %v", err)
	}

	// Corrupt a byte in the middle of the ciphertext
	ciphertext[len(ciphertext)/2] ^= 0xFF

	_, err = crypto.Decrypt(key, ciphertext)
	if err != crypto.ErrAuthFailed {
		t.Errorf("Decrypt(corrupted) returned %v, want ErrAuthFailed", err)
	}
}

// BenchmarkEncrypt measures AES-GCM encryption throughput.
func BenchmarkEncrypt(b *testing.B) {
	key := make([]byte, 32)
	plaintext := make([]byte, 1024) // 1 KB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := crypto.Encrypt(key, plaintext)
		if err != nil {
			b.Fatalf("Encrypt failed: %v", err)
		}
	}
}

// BenchmarkDecrypt measures AES-GCM decryption throughput.
func BenchmarkDecrypt(b *testing.B) {
	key := make([]byte, 32)
	plaintext := make([]byte, 1024) // 1 KB

	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		b.Fatalf("Encrypt failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := crypto.Decrypt(key, ciphertext)
		if err != nil {
			b.Fatalf("Decrypt failed: %v", err)
		}
	}
}
