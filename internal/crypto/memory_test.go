package crypto_test

import (
	"bytes"
	"testing"

	"github.com/useful-toys/abditum/internal/crypto"
)

func TestWipeSlice(t *testing.T) {
	data := []byte("sensitive data that must be cleared")
	original := make([]byte, len(data))
	copy(original, data)

	crypto.Wipe(data)

	// Verify all bytes are zeroed
	for i, b := range data {
		if b != 0 {
			t.Errorf("byte at index %d not zeroed: got %d, want 0", i, b)
		}
	}

	// Verify we actually had data before
	if bytes.Equal(original, data) {
		t.Error("data was not modified by Wipe()")
	}
}

func TestWipeNil(t *testing.T) {
	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Wipe() panicked on nil slice: %v", r)
		}
	}()

	crypto.Wipe(nil)
}

func TestWipeEmpty(t *testing.T) {
	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Wipe() panicked on empty slice: %v", r)
		}
	}()

	crypto.Wipe([]byte{})
}

func TestSecureAllocate(t *testing.T) {
	size := 1024
	buf, cleanup, err := crypto.SecureAllocate(size)
	// Memory locking may fail on platforms without support - buffer still usable
	if err != nil {
		t.Logf("SecureAllocate() returned error (non-fatal): %v", err)
	}
	defer cleanup()

	// Verify size
	if len(buf) != size {
		t.Errorf("buffer size = %d, want %d", len(buf), size)
	}

	// Verify initialized to zero
	for i, b := range buf {
		if b != 0 {
			t.Errorf("byte at index %d not zero: got %d", i, b)
		}
	}

	// Verify writable
	buf[0] = 42
	if buf[0] != 42 {
		t.Error("buffer not writable")
	}
}

func TestSecureAllocateZeroSize(t *testing.T) {
	buf, cleanup, err := crypto.SecureAllocate(0)
	if err != nil {
		t.Logf("SecureAllocate(0) returned error (non-fatal): %v", err)
	}
	defer cleanup()

	if len(buf) != 0 {
		t.Errorf("buffer size = %d, want 0", len(buf))
	}
}

func TestSecureAllocateCleanup(t *testing.T) {
	buf, cleanup, err := crypto.SecureAllocate(256)
	if err != nil {
		t.Logf("SecureAllocate() returned error (non-fatal): %v", err)
	}

	// Write data
	for i := range buf {
		buf[i] = byte(i % 256)
	}

	// Cleanup should wipe
	cleanup()

	// Verify wiped (best effort - compiler might optimize away)
	allZero := true
	for _, b := range buf {
		if b != 0 {
			allZero = false
			break
		}
	}
	if !allZero {
		t.Error("cleanup() did not wipe buffer")
	}
}

func TestSecureAllocateMultipleCalls(t *testing.T) {
	// Verify multiple allocations work
	cleanups := make([]func(), 0, 5)
	defer func() {
		for _, cleanup := range cleanups {
			cleanup()
		}
	}()

	for i := 0; i < 5; i++ {
		buf, cleanup, err := crypto.SecureAllocate(128)
		if err != nil {
			t.Logf("SecureAllocate() call %d returned error (non-fatal): %v", i, err)
		}
		cleanups = append(cleanups, cleanup)

		if len(buf) != 128 {
			t.Errorf("call %d: buffer size = %d, want 128", i, len(buf))
		}
	}
}
