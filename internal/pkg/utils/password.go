package utils

import (
	"fmt"
	"strings"
	"unicode"
)

// ValidateAdminPassword checks if a password meets the minimum security requirements.
// Requirements: at least 8 characters, one uppercase, one lowercase, one digit, one special character.
func ValidateAdminPassword(password string) error {
	var errs []string

	if len(password) < 8 {
		errs = append(errs, "must be at least 8 characters")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case strings.ContainsAny(string(ch), "!@#$%^&*"):
			hasSpecial = true
		}
	}

	if !hasUpper {
		errs = append(errs, "must contain at least one uppercase letter")
	}
	if !hasLower {
		errs = append(errs, "must contain at least one lowercase letter")
	}
	if !hasDigit {
		errs = append(errs, "must contain at least one digit")
	}
	if !hasSpecial {
		errs = append(errs, "must contain at least one special character (!@#$%^&*)")
	}

	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "; "))
	}

	return nil
}
