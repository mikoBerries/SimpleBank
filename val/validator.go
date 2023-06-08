// Package val populate func() to validate incoming request
package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

// Constant function of regex used
var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

// ValidateString validate string min&max behaviour
func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must between %d - %d character", minLength, maxLength)
	}
	return nil
}

// ValidateUsername to validate username request request(min max length & regex check)
func ValidateUsername(value string) error {
	// username string length berween 3-100
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	//regex check
	if !isValidUsername(value) {
		return fmt.Errorf("must contain only lowercase letters, digits, or underscore")
	}
	return nil
}

// ValidateFullName to validate full_name field request(min max length & regex check)
func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	//regex for full name
	if !isValidFullName(value) {
		return fmt.Errorf("must contain only letters or spaces")
	}
	return nil
}

// ValidatePassword to validate password request (min-max length)
func ValidatePassword(value string) error {
	return ValidateString(value, 6, 100)
}

// ValidateEmail to validate email reqeust (email structure and email min max length )
func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 200); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not a valid email address")
	}
	return nil
}
