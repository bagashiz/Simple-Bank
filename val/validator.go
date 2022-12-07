package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	// Valid username must contain only lowercase alphanumeric characters and underscores
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	// valid full name must contain only alphabetical characters and spaces
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

// ValidateString validates that a string is between a minimum and maximum length
func ValidateString(value string, minLength, maxLength int) error {
	valLength := len(value)
	if valLength < minLength || valLength > maxLength {
		return fmt.Errorf("must contain between %d-%d characters", minLength, maxLength)
	}

	return nil
}

// ValidateUsername validates that a username is valid based on the following rules:
// - must contain between 3-100 characters
// - must contain only lowercase alphanumeric characters and underscores
func ValidateUsername(username string) error {
	if err := ValidateString(username, 3, 100); err != nil {
		return err
	}

	if !isValidUsername(username) {
		return fmt.Errorf("must contain only lowercase alphanumeric characters and underscores")
	}

	return nil
}

// ValidateFullName validates that a full name is valid based on the following rules:
// - must contain between 3-100 characters
// - must contain only alphabetical characters and spaces
func ValidateFullName(fullname string) error {
	if err := ValidateString(fullname, 3, 100); err != nil {
		return err
	}

	if !isValidFullName(fullname) {
		return fmt.Errorf("must contain only alphabetical characters and spaces")
	}

	return nil
}

// ValidatePassword validates that a password is valid based on the following rules:
// - must contain between 6-100 characters
func ValidatePassword(password string) error {
	return ValidateString(password, 6, 100)
}

// ValidateEmail validates that an email is valid based on the following rules:
// - must contain between 3-100 characters
// - must be a valid email address
func ValidateEmail(email string) error {
	if err := ValidateString(email, 3, 100); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("is not a valid email address")
	}

	return nil
}
