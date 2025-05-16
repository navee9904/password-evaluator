package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"password-evaluator/crack"
	"password-evaluator/eval"
	"time"
)

// EvaluationResult holds the password evaluation results.
type EvaluationResult struct {
	LengthValid       bool    `json:"lengthValid"`
	HasUppercase      bool    `json:"hasUppercase"`
	HasLowercase      bool    `json:"hasLowercase"`
	HasNumber         bool    `json:"hasNumber"`
	HasSpecial        bool    `json:"hasSpecial"`
	HasCommonPattern  bool    `json:"hasCommonPattern"`
	CrackingTimeYears float64 `json:"crackingTimeYears"`
	Strength          string  `json:"strength"`
	SuggestedPassword string  `json:"suggestedPassword,omitempty"`
}

// evaluateHandler processes password evaluation requests.
func evaluateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Password == "" {
		http.Error(w, "Password cannot be empty", http.StatusBadRequest)
		return
	}

	// Evaluate password
	lengthValid, hasUpper, hasLower, hasNumber, hasSpecial, err := eval.CheckLengthAndVariety(req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error evaluating password: %v", err), http.StatusInternalServerError)
		return
	}

	hasCommonPattern, err := eval.DetectCommonPatterns(req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error detecting patterns: %v", err), http.StatusInternalServerError)
		return
	}

	crackingTime, err := crack.EstimateCrackingTime(req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error estimating cracking time: %v", err), http.StatusInternalServerError)
		return
	}

	// Determine password strength
	strength := calculateStrength(lengthValid, hasUpper, hasLower, hasNumber, hasSpecial, hasCommonPattern, crackingTime)

	// Generate suggested password if strength is Weak
	var suggestedPassword string
	if strength == "Weak" {
		suggestedPassword = generateSuggestedPassword()
	}

	// Prepare response
	result := EvaluationResult{
		LengthValid:       lengthValid,
		HasUppercase:      hasUpper,
		HasLowercase:      hasLower,
		HasNumber:         hasNumber,
		HasSpecial:        hasSpecial,
		HasCommonPattern:  hasCommonPattern,
		CrackingTimeYears: crackingTime,
		Strength:          strength,
		SuggestedPassword: suggestedPassword,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// suggestHandler generates a suggested strong password.
func suggestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	suggestedPassword := generateSuggestedPassword()
	response := struct {
		SuggestedPassword string `json:"suggestedPassword"`
	}{
		SuggestedPassword: suggestedPassword,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve index.html at root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	// API endpoints
	http.HandleFunc("/api/evaluate", evaluateHandler)
	http.HandleFunc("/api/suggest", suggestHandler)

	// Start server
	addr := ":8080"
	fmt.Println("Server starting on", addr)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// calculateStrength determines the password strength based on evaluation criteria.
func calculateStrength(lengthValid, hasUpper, hasLower, hasNumber, hasSpecial, hasCommonPattern bool, crackingTime float64) string {
	// Count fulfilled criteria
	fulfilledCriteria := 0
	if lengthValid {
		fulfilledCriteria++
	}
	if hasUpper {
		fulfilledCriteria++
	}
	if hasLower {
		fulfilledCriteria++
	}
	if hasNumber {
		fulfilledCriteria++
	}
	if hasSpecial {
		fulfilledCriteria++
	}
	if !hasCommonPattern {
		fulfilledCriteria++
	}

	// Strength rules
	if fulfilledCriteria == 6 && crackingTime > 1000 { // All criteria met and cracking time is very high
		return "Strong"
	} else if fulfilledCriteria >= 4 && crackingTime > 1 { // Most criteria met and cracking time is decent
		return "Medium"
	} else { // Fails multiple criteria or cracking time is low
		return "Weak"
	}
}

// generateSuggestedPassword creates a strong password that meets all criteria.
func generateSuggestedPassword() string {
	rand.Seed(time.Now().UnixNano())
	length := 16 // Generate a 16-character password

	// Character sets
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lower := "abcdefghijklmnopqrstuvwxyz"
	numbers := "0123456789"
	special := "@#$"

	// Ensure at least one character from each set
	var password []byte
	password = append(password, upper[rand.Intn(len(upper))])
	password = append(password, lower[rand.Intn(len(lower))])
	password = append(password, numbers[rand.Intn(len(numbers))])
	password = append(password, special[rand.Intn(len(special))])

	// Fill the remaining characters
	allChars := upper + lower + numbers + special
	for i := 4; i < length; i++ {
		password = append(password, allChars[rand.Intn(len(allChars))])
	}

	// Shuffle the password
	for i := range password {
		j := rand.Intn(len(password))
		password[i], password[j] = password[j], password[i]
	}

	return string(password)
}
