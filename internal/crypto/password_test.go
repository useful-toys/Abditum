package crypto_test

import (
	"testing"

	"github.com/useful-toys/abditum/internal/crypto"
)

func TestPasswordStrength12CharsAllCategories(t *testing.T) {
	password := []byte("Abc1!Abc1!12")
	strength := crypto.EvaluatePasswordStrength(password)
	if strength != crypto.StrengthStrong {
		t.Errorf("EvaluatePasswordStrength(%q) = %v, want StrengthStrong", password, strength)
	}
}

func TestPasswordStrength11Chars(t *testing.T) {
	password := []byte("Abc1!Abc1!1")
	strength := crypto.EvaluatePasswordStrength(password)
	if strength != crypto.StrengthWeak {
		t.Errorf("EvaluatePasswordStrength(%q) = %v, want StrengthWeak (11 chars)", password, strength)
	}
}

func TestPasswordStrengthNoUppercase(t *testing.T) {
	password := []byte("abc123!@#456")
	strength := crypto.EvaluatePasswordStrength(password)
	if strength != crypto.StrengthWeak {
		t.Errorf("EvaluatePasswordStrength(%q) = %v, want StrengthWeak (no uppercase)", password, strength)
	}
}

func TestPasswordStrengthNoLowercase(t *testing.T) {
	password := []byte("ABC123!@#456")
	strength := crypto.EvaluatePasswordStrength(password)
	if strength != crypto.StrengthWeak {
		t.Errorf("EvaluatePasswordStrength(%q) = %v, want StrengthWeak (no lowercase)", password, strength)
	}
}

func TestPasswordStrengthNoDigit(t *testing.T) {
	password := []byte("Abcdef!@#Xyz")
	strength := crypto.EvaluatePasswordStrength(password)
	if strength != crypto.StrengthWeak {
		t.Errorf("EvaluatePasswordStrength(%q) = %v, want StrengthWeak (no digit)", password, strength)
	}
}

func TestPasswordStrengthNoSpecial(t *testing.T) {
	password := []byte("Abcdef123456")
	strength := crypto.EvaluatePasswordStrength(password)
	if strength != crypto.StrengthWeak {
		t.Errorf("EvaluatePasswordStrength(%q) = %v, want StrengthWeak (no special)", password, strength)
	}
}

func TestPasswordStrengthWeak(t *testing.T) {
	password := []byte("abc123")
	strength := crypto.EvaluatePasswordStrength(password)
	if strength != crypto.StrengthWeak {
		t.Errorf("EvaluatePasswordStrength(%q) = %v, want StrengthWeak", password, strength)
	}
}
