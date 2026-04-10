package crypto_test

import (
	"bytes"
	"testing"

	"github.com/useful-toys/abditum/internal/crypto"
)

func TestFullRoundtrip(t *testing.T) {
	// Generate salt
	salt, err := crypto.GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() failed: %v", err)
	}

	// Derive key from password
	password := []byte("TestVaultPassword123!")
	params := crypto.ArgonParams{
		Time:    3,
		Memory:  262144, // 256 MiB
		Threads: 4,
		KeyLen:  32,
	}

	key, err := crypto.DeriveKey(password, salt, params)
	if err != nil {
		t.Fatalf("DeriveKey() failed: %v", err)
	}
	defer crypto.Wipe(key)

	// Encrypt plaintext
	plaintext := []byte("secret vault data")
	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt() failed: %v", err)
	}

	// Decrypt ciphertext
	decrypted, err := crypto.Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt() failed: %v", err)
	}

	// Verify roundtrip
	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypt() = %q, want %q", decrypted, plaintext)
	}

	// Zero key and verify
	crypto.Wipe(key)
	for i, b := range key {
		if b != 0 {
			t.Errorf("key[%d] = %d after Wipe(), want 0", i, b)
		}
	}
}

func TestKeyReuseProducesDistinctCiphertexts(t *testing.T) {
	// Derive key
	salt, err := crypto.GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() failed: %v", err)
	}

	password := []byte("TestPassword123!")
	params := crypto.ArgonParams{
		Time:    3,
		Memory:  262144,
		Threads: 4,
		KeyLen:  32,
	}

	key, err := crypto.DeriveKey(password, salt, params)
	if err != nil {
		t.Fatalf("DeriveKey() failed: %v", err)
	}
	defer crypto.Wipe(key)

	// Encrypt same plaintext twice with same key
	plaintext := []byte("same plaintext every time")

	ciphertext1, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt() #1 failed: %v", err)
	}

	ciphertext2, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt() #2 failed: %v", err)
	}

	// Verify ciphertexts are different (proves nonce uniqueness)
	if bytes.Equal(ciphertext1, ciphertext2) {
		t.Error("Encrypt() produced identical ciphertexts for same plaintext - nonce not unique!")
	}
}

func TestMemorySafety(t *testing.T) {
	// Create sensitive buffer
	sensitive := []byte("sensitive data here")

	// Verify buffer contains data
	if bytes.Equal(sensitive, make([]byte, len(sensitive))) {
		t.Error("sensitive buffer is already zeroed before Wipe()")
	}

	// Zero the buffer
	crypto.Wipe(sensitive)

	// Verify every byte is 0x00
	for i, b := range sensitive {
		if b != 0x00 {
			t.Errorf("sensitive[%d] = %#x after Wipe(), want 0x00", i, b)
		}
	}
}
