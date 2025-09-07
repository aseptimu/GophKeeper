package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidCardNumber = errors.New("invalid card number")
	ErrInvalidExpiryDate = errors.New("invalid expiry date")
	ErrInvalidCVV        = errors.New("invalid CVV")
)

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

func LuhnCheck(cardNumber string) bool {
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

func ValidateCVV(cvv string) error {
	if !regexp.MustCompile(`^\d{3,4}$`).MatchString(cvv) {
		return ErrInvalidCVV
	}
	return nil
}
