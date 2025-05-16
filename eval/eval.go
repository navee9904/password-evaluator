package eval

import (
	"regexp"
	"strings"
	"unicode"
)

// Common passwords and patterns (simplified for demo)
var commonPasswords = []string{"password", "123456", "qwerty", "admin123"}
var keyboardPattern = regexp.MustCompile(`qwerty|asdf|zxcv`)

// CheckLengthAndVariety checks if the password meets length and character variety requirements.
func CheckLengthAndVariety(password string) (lengthValid, hasUpper, hasLower, hasNumber, hasSpecial bool, err error) {
	lengthValid = len(password) >= 12

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsNumber(r):
			hasNumber = true
		case strings.ContainsRune("@#$", r): // Limited to @, #, $ as specified
			hasSpecial = true
		}
	}

	return lengthValid, hasUpper, hasLower, hasNumber, hasSpecial, nil
}

// DetectCommonPatterns checks for common passwords or patterns.
func DetectCommonPatterns(password string) (bool, error) {
	lowerPass := strings.ToLower(password)
	// Check common passwords
	for _, common := range commonPasswords {
		if strings.Contains(lowerPass, common) {
			return true, nil
		}
	}

	// Check for repetitive characters (e.g., "aaa" or "111")
	count := 1
	for i := 1; i < len(password); i++ {
		if password[i] == password[i-1] {
			count++
			if count >= 3 { // 3 or more consecutive identical characters
				return true, nil
			}
		} else {
			count = 1
		}
	}

	// Check keyboard patterns
	if keyboardPattern.MatchString(lowerPass) {
		return true, nil
	}
	return false, nil
}
