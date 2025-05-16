package crack

import (
	"crypto/sha256"
	"math"
	"strings"
)

// EstimateCrackingTime estimates the time to crack the password using SHA256 brute-force simulation.
func EstimateCrackingTime(password string) (float64, error) {
	// Calculate character set size
	charSetSize := 0
	hasUpper, hasLower, hasNumber, hasSpecial := false, false, false, false
	for _, r := range password {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= '0' && r <= '9':
			hasNumber = true
		case strings.ContainsRune("@#$", r): // Limited to @, #, $ as specified
			hasSpecial = true
		}
	}
	if hasUpper {
		charSetSize += 26
	}
	if hasLower {
		charSetSize += 26
	}
	if hasNumber {
		charSetSize += 10
	}
	if hasSpecial {
		charSetSize += 3 // Only @, #, $
	}
	if charSetSize == 0 {
		charSetSize = 1 // Avoid division by zero
	}

	// Calculate total combinations
	combinations := math.Pow(float64(charSetSize), float64(len(password)))

	// Assume 10 billion hashes per second (modern GPU estimate)
	hashesPerSecond := 1e10
	secondsToCrack := combinations / hashesPerSecond
	yearsToCrack := secondsToCrack / (365 * 24 * 60 * 60)

	// Simulate SHA256 hashing (single hash for demo, not actual brute-force)
	hash := sha256.Sum256([]byte(password))
	_ = hash // Simulate hashing work

	return yearsToCrack, nil
}
