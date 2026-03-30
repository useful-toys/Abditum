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

// TestEncryptWithAAD verifies EncryptWithAAD returns separate nonce and ciphertext
// and that the nonce is unique across calls.
func TestEncryptWithAAD(t *testing.T) {
	key := make([]byte, 32)
	plaintext := []byte("test plaintext for AAD")
	aad := []byte("additional authenticated data")

	t.Run("returns nonce and ciphertext separately", func(t *testing.T) {
		nonce, ciphertext, err := crypto.EncryptWithAAD(key, plaintext, aad)
		if err != nil {
			t.Fatalf("EncryptWithAAD() returned unexpected error: %v", err)
		}
		if len(nonce) != 12 {
			t.Errorf("EncryptWithAAD() nonce length = %d, want 12", len(nonce))
		}
		// ciphertext = plaintext + GCM tag (16 bytes)
		if len(ciphertext) != len(plaintext)+16 {
			t.Errorf("EncryptWithAAD() ciphertext length = %d, want %d", len(ciphertext), len(plaintext)+16)
		}
	})

	t.Run("invalid key returns ErrInvalidParams", func(t *testing.T) {
		_, _, err := crypto.EncryptWithAAD(nil, plaintext, aad)
		if err != crypto.ErrInvalidParams {
			t.Errorf("EncryptWithAAD(nil key) = %v, want ErrInvalidParams", err)
		}
		_, _, err = crypto.EncryptWithAAD([]byte{1, 2, 3}, plaintext, aad)
		if err != crypto.ErrInvalidParams {
			t.Errorf("EncryptWithAAD(short key) = %v, want ErrInvalidParams", err)
		}
	})

	t.Run("nonce uniqueness across calls", func(t *testing.T) {
		nonce1, _, err1 := crypto.EncryptWithAAD(key, plaintext, aad)
		if err1 != nil {
			t.Fatalf("EncryptWithAAD() first call error: %v", err1)
		}
		nonce2, _, err2 := crypto.EncryptWithAAD(key, plaintext, aad)
		if err2 != nil {
			t.Fatalf("EncryptWithAAD() second call error: %v", err2)
		}
		if bytes.Equal(nonce1, nonce2) {
			t.Error("EncryptWithAAD() produced identical nonces - NONCE REUSE DETECTED")
		}
	})

	t.Run("empty plaintext succeeds", func(t *testing.T) {
		nonce, ciphertext, err := crypto.EncryptWithAAD(key, []byte{}, aad)
		if err != nil {
			t.Fatalf("EncryptWithAAD(empty) returned error: %v", err)
		}
		if len(nonce) != 12 {
			t.Errorf("nonce length = %d, want 12", len(nonce))
		}
		// empty plaintext: ciphertext = 0 + 16 bytes GCM tag
		if len(ciphertext) != 16 {
			t.Errorf("ciphertext length = %d, want 16", len(ciphertext))
		}
	})
}

// TestDecryptWithAAD verifies DecryptWithAAD authenticates AAD and returns original plaintext.
func TestDecryptWithAAD(t *testing.T) {
	key := make([]byte, 32)
	plaintext := []byte("secret vault contents")
	aad := []byte("49-byte file header as AAD")

	// Encrypt once for reuse in sub-tests
	nonce, ciphertext, err := crypto.EncryptWithAAD(key, plaintext, aad)
	if err != nil {
		t.Fatalf("EncryptWithAAD() setup error: %v", err)
	}

	t.Run("roundtrip succeeds", func(t *testing.T) {
		got, err := crypto.DecryptWithAAD(key, ciphertext, nonce, aad)
		if err != nil {
			t.Fatalf("DecryptWithAAD() returned error: %v", err)
		}
		if !bytes.Equal(got, plaintext) {
			t.Errorf("DecryptWithAAD() = %q, want %q", got, plaintext)
		}
	})

	t.Run("wrong key returns ErrAuthFailed", func(t *testing.T) {
		wrongKey := make([]byte, 32)
		wrongKey[0] = 0xFF
		_, err := crypto.DecryptWithAAD(wrongKey, ciphertext, nonce, aad)
		if err != crypto.ErrAuthFailed {
			t.Errorf("DecryptWithAAD(wrong key) = %v, want ErrAuthFailed", err)
		}
	})

	t.Run("tampered AAD returns ErrAuthFailed", func(t *testing.T) {
		tamperedAAD := []byte("tampered header bytes here!!")
		_, err := crypto.DecryptWithAAD(key, ciphertext, nonce, tamperedAAD)
		if err != crypto.ErrAuthFailed {
			t.Errorf("DecryptWithAAD(tampered AAD) = %v, want ErrAuthFailed", err)
		}
	})

	t.Run("tampered ciphertext returns ErrAuthFailed", func(t *testing.T) {
		tampered := make([]byte, len(ciphertext))
		copy(tampered, ciphertext)
		tampered[len(tampered)/2] ^= 0xFF
		_, err := crypto.DecryptWithAAD(key, tampered, nonce, aad)
		if err != crypto.ErrAuthFailed {
			t.Errorf("DecryptWithAAD(tampered ciphertext) = %v, want ErrAuthFailed", err)
		}
	})

	t.Run("invalid key returns ErrInvalidParams", func(t *testing.T) {
		_, err := crypto.DecryptWithAAD(nil, ciphertext, nonce, aad)
		if err != crypto.ErrInvalidParams {
			t.Errorf("DecryptWithAAD(nil key) = %v, want ErrInvalidParams", err)
		}
	})

	t.Run("empty plaintext roundtrip", func(t *testing.T) {
		emptyNonce, emptyCipher, err := crypto.EncryptWithAAD(key, []byte{}, aad)
		if err != nil {
			t.Fatalf("EncryptWithAAD(empty) error: %v", err)
		}
		got, err := crypto.DecryptWithAAD(key, emptyCipher, emptyNonce, aad)
		if err != nil {
			t.Fatalf("DecryptWithAAD(empty) error: %v", err)
		}
		if !bytes.Equal(got, []byte{}) {
			t.Errorf("empty roundtrip = %v, want []byte{}", got)
		}
	})
}
