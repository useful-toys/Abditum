//go:build windows

package crypto

import "testing"

func TestMlockEmptyBuffer(t *testing.T) {
	if err := mlock([]byte{}); err != nil {
		t.Fatalf("mlock(empty) error = %v, want nil", err)
	}
}

func TestMunlockEmptyBuffer(t *testing.T) {
	if err := munlock([]byte{}); err != nil {
		t.Fatalf("munlock(empty) error = %v, want nil", err)
	}
}

func TestMlockNonEmptyBuffer(t *testing.T) {
	buf := make([]byte, 1)

	err := mlock(buf)
	if err != nil && err != ErrMLockFailed {
		t.Fatalf("mlock(non-empty) error = %v, want nil or ErrMLockFailed", err)
	}

	if err == nil {
		if unlockErr := munlock(buf); unlockErr != nil {
			t.Fatalf("munlock(non-empty) error = %v, want nil", unlockErr)
		}
	}
}
