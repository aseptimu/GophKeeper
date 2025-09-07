// Package utils provides utility functions for data validation and other common operations.
package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	// ErrInvalidCardNumber is returned when a bank card number is invalid.
	ErrInvalidCardNumber = errors.New("invalid card number")
	// ErrInvalidExpiryDate is returned when a card expiry date is invalid.
	ErrInvalidExpiryDate = errors.New("invalid expiry date")
	// ErrInvalidCVV is returned when a card CVV is invalid.
	ErrInvalidCVV = errors.New("invalid CVV")
)

// ValidateBankCard validates all components of a bank card.
// It checks the card number using Luhn algorithm, expiry date format, and CVV format.
func ValidateBankCard(cardNumber, expiryDate, cvv string) error {
	if err := ValidateCardNumber(cardNumber); err != nil {
		return err
	}

	if err := ValidateExpiryDate(expiryDate); err != nil {
		return err
	}

	if err := ValidateCVV(cvv); err != nil {
		return err
	}

	return nil
}

// ValidateCardNumber validates a bank card number using the Luhn algorithm.
// It accepts card numbers with spaces and dashes, which are automatically removed.
func ValidateCardNumber(cardNumber string) error {
	cardNumber = strings.ReplaceAll(cardNumber, " ", "")
	cardNumber = strings.ReplaceAll(cardNumber, "-", "")

	if len(cardNumber) < 13 || len(cardNumber) > 19 {
		return ErrInvalidCardNumber
	}

	if !regexp.MustCompile(`^\d+$`).MatchString(cardNumber) {
		return ErrInvalidCardNumber
	}

	if !LuhnCheck(cardNumber) {
		return ErrInvalidCardNumber
	}

	return nil
}

// LuhnCheck implements the Luhn algorithm to validate credit card numbers.
// Returns true if the card number is valid according to the Luhn algorithm.
func LuhnCheck(cardNumber string) bool {
	if len(cardNumber) == 0 {
		return false
	}

	sum := 0
	alternate := false

	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(cardNumber[i]))
		if err != nil {
			return false
		}

		if alternate {
			digit *= 2
			if digit > 9 {
				digit = digit%10 + digit/10
			}
		}

		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}

// ValidateExpiryDate validates a card expiry date in MM/YY format.
// Month must be between 01-12, year must be between 00-99.
func ValidateExpiryDate(expiryDate string) error {
	if !regexp.MustCompile(`^\d{2}/\d{2}$`).MatchString(expiryDate) {
		return ErrInvalidExpiryDate
	}

	parts := strings.Split(expiryDate, "/")
	month, err := strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return ErrInvalidExpiryDate
	}

	year, err := strconv.Atoi(parts[1])
	if err != nil || year < 0 || year > 99 {
		return ErrInvalidExpiryDate
	}

	return nil
}

// ValidateCVV validates a card CVV (Card Verification Value).
// CVV must be 3 or 4 digits.
func ValidateCVV(cvv string) error {
	if !regexp.MustCompile(`^\d{3,4}$`).MatchString(cvv) {
		return ErrInvalidCVV
	}
	return nil
}
