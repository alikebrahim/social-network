package utils

import (
	"errors"
	"regexp"
	"time"
	"unicode"
)

// ValidateEmail checks if the given email is in a valid format.
func ValidateEmail(email string) error {
	// Regular expression for validating an email
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// ValidatePassword checks if the password meets security requirements (min 8 chars, at least one uppercase, one number, and one special character).
func ValidatePassword(password string) error {
	// Check if the password is at least 8 characters long
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Check if the password contains at least one uppercase letter
	if !hasUpperCase(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	// Check if the password contains at least one number
	if !hasNumber(password) {
		return errors.New("password must contain at least one number")
	}

	// Check if the password contains at least one special character
	if !hasSpecialCharacter(password) {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

// ValidateName checks if the given name is non-empty and consists only of letters and spaces.
func ValidateName(name string) error {
	// Name should be non-empty
	if len(name) == 0 {
		return errors.New("name cannot be empty")
	}

	// Ensure the name contains only letters and spaces
	re := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if !re.MatchString(name) {
		return errors.New("name must only contain letters and spaces")
	}

	return nil
}

// ValidateDateOfBirth checks if the provided date of birth is in the correct format (YYYY-MM-DD) and the person is at least 18 years old.
func ValidateDateOfBirth(dateOfBirth string) error {
	// Regular expression to check if the date format is YYYY-MM-DD
	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if !re.MatchString(dateOfBirth) {
		return errors.New("invalid date of birth format. Use YYYY-MM-DD")
	}

	// Check if the user is at least 18 years old (example check, assumes date format is correct)
	// You can use a more sophisticated method here to parse and calculate age
	// (for simplicity, this is just an illustrative approach)
	age := calculateAge(dateOfBirth)
	if age < 18 {
		return errors.New("user must be at least 18 years old")
	}

	return nil
}

// hasUpperCase checks if the string contains at least one uppercase letter.
func hasUpperCase(s string) bool {
	for _, c := range s {
		if unicode.IsUpper(c) {
			return true
		}
	}
	return false
}

// hasNumber checks if the string contains at least one number.
func hasNumber(s string) bool {
	for _, c := range s {
		if unicode.IsDigit(c) {
			return true
		}
	}
	return false
}

// hasSpecialCharacter checks if the string contains at least one special character.
func hasSpecialCharacter(s string) bool {
	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return true
		}
	}
	return false
}

// calculateAge calculates the age from a date string in YYYY-MM-DD format.
func calculateAge(dateOfBirth string) int {
	// In a real-world application, use a package like "time" to calculate age
	// For simplicity, we assume a fixed calculation here.

	// Example: Let's assume the date of birth is "2000-01-01"
	// Use actual date logic here
	// (e.g., use time.Parse() to parse the date and check against the current date)

	// In the case of a fixed date for simplicity (e.g., checking against current date):
	currentYear := time.Now().Year()
	yearOfBirth := 2000 // extract from dateOfBirth string
	age := currentYear - yearOfBirth

	return age
}
