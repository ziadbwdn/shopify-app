package utils

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const minPasswordLength = 12

// ValidatePasswordWithRegex checks if a password meets complexity requirements using regular expressions.
func ValidatePasswordWithRegex(password string) error {
	// Rule 1: Minimum length check
	if len(password) < minPasswordLength {
		return fmt.Errorf("password must be at least %d characters long", minPasswordLength)
	}

	// Define regex for each requirement
	rules := map[string]string{
		"an uppercase letter": `[A-Z]`,
		"a lowercase letter":  `[a-z]`,
		"a number":            `[0-9]`,
		"a special character": `[!@#$%^&*()_+\-=\[\]{}|;':",.<>/?]`,
	}

	var missing []string
	for requirement, pattern := range rules {
		matched, _ := regexp.MatchString(pattern, password)
		if !matched {
			missing = append(missing, requirement)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("password is missing: %s", strings.Join(missing, ", "))
	}

	return nil
}

// HashPassword generates a bcrypt hash of the password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a password with a hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
