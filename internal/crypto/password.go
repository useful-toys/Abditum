package crypto

// StrengthLevel represents the evaluated strength of a password.
type StrengthLevel int

const (
	// StrengthWeak indicates a password that does not meet strength requirements:
	// - Less than 12 characters, OR
	// - Missing any character category (uppercase, lowercase, digit, special)
	StrengthWeak StrengthLevel = iota

	// StrengthStrong indicates a password that meets all strength requirements:
	// - At least 12 characters, AND
	// - Contains uppercase, lowercase, digit, and special characters
	StrengthStrong
)

// EvaluatePasswordStrength evaluates password strength based on length and character diversity.
//
// Requirements for StrengthStrong (PWD-01, D-34):
//   - Minimum length: 12 characters
//   - Must contain all four character categories:
//   - Uppercase letters (A-Z)
//   - Lowercase letters (a-z)
//   - Digits (0-9)
//   - Special characters (anything not in above categories)
//
// CRITICAL: Operates directly on []byte without string conversion (Pitfall 3).
// String conversion creates unzeroable copies in memory.
func EvaluatePasswordStrength(password []byte) StrengthLevel {
	// Length requirement
	if len(password) < 12 {
		return StrengthWeak
	}

	// Category flags
	var hasUpper, hasLower, hasDigit, hasSpecial bool

	// Check character categories
	for _, b := range password {
		switch {
		case b >= 'A' && b <= 'Z':
			hasUpper = true
		case b >= 'a' && b <= 'z':
			hasLower = true
		case b >= '0' && b <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	// All four categories required for strong password
	if hasUpper && hasLower && hasDigit && hasSpecial {
		return StrengthStrong
	}

	return StrengthWeak
}
