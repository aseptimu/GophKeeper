package utils

import (
	"testing"
)

func TestValidateBankCard(t *testing.T) {
	tests := []struct {
		name        string
		cardNumber  string
		expiryDate  string
		cvv         string
		expectError bool
	}{
		{
			name:        "valid visa card",
			cardNumber:  "4111111111111111",
			expiryDate:  "12/25",
			cvv:         "123",
			expectError: false,
		},
		{
			name:        "valid mastercard",
			cardNumber:  "5555555555554444",
			expiryDate:  "01/26",
			cvv:         "456",
			expectError: false,
		},
		{
			name:        "valid amex card",
			cardNumber:  "378282246310005",
			expiryDate:  "06/27",
			cvv:         "1234",
			expectError: false,
		},
		{
			name:        "invalid card number",
			cardNumber:  "1234567890123456",
			expiryDate:  "12/25",
			cvv:         "123",
			expectError: true,
		},
		{
			name:        "invalid expiry date format",
			cardNumber:  "4111111111111111",
			expiryDate:  "12-25",
			cvv:         "123",
			expectError: true,
		},
		{
			name:        "invalid expiry month",
			cardNumber:  "4111111111111111",
			expiryDate:  "13/25",
			cvv:         "123",
			expectError: true,
		},
		{
			name:        "invalid cvv",
			cardNumber:  "4111111111111111",
			expiryDate:  "12/25",
			cvv:         "12",
			expectError: true,
		},
		{
			name:        "card number with spaces",
			cardNumber:  "4111 1111 1111 1111",
			expiryDate:  "12/25",
			cvv:         "123",
			expectError: false,
		},
		{
			name:        "card number with dashes",
			cardNumber:  "4111-1111-1111-1111",
			expiryDate:  "12/25",
			cvv:         "123",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBankCard(tt.cardNumber, tt.expiryDate, tt.cvv)
			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateCardNumber(t *testing.T) {
	tests := []struct {
		name        string
		cardNumber  string
		expectError bool
	}{
		{
			name:        "valid visa",
			cardNumber:  "4111111111111111",
			expectError: false,
		},
		{
			name:        "valid mastercard",
			cardNumber:  "5555555555554444",
			expectError: false,
		},
		{
			name:        "valid amex",
			cardNumber:  "378282246310005",
			expectError: false,
		},
		{
			name:        "invalid luhn",
			cardNumber:  "1234567890123456",
			expectError: true,
		},
		{
			name:        "too short",
			cardNumber:  "123456789012",
			expectError: true,
		},
		{
			name:        "too long",
			cardNumber:  "12345678901234567890",
			expectError: true,
		},
		{
			name:        "non-numeric",
			cardNumber:  "411111111111111a",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCardNumber(tt.cardNumber)
			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestLuhnCheck(t *testing.T) {
	tests := []struct {
		name     string
		number   string
		expected bool
	}{
		{
			name:     "valid visa",
			number:   "4111111111111111",
			expected: true,
		},
		{
			name:     "valid mastercard",
			number:   "5555555555554444",
			expected: true,
		},
		{
			name:     "invalid number",
			number:   "1234567890123456",
			expected: false,
		},
		{
			name:     "empty string",
			number:   "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LuhnCheck(tt.number)
			if result != tt.expected {
				t.Errorf("LuhnCheck(%s) = %v, expected %v", tt.number, result, tt.expected)
			}
		})
	}
}

func TestValidateExpiryDate(t *testing.T) {
	tests := []struct {
		name        string
		expiryDate  string
		expectError bool
	}{
		{
			name:        "valid date",
			expiryDate:  "12/25",
			expectError: false,
		},
		{
			name:        "valid date with leading zero",
			expiryDate:  "01/26",
			expectError: false,
		},
		{
			name:        "invalid format",
			expiryDate:  "12-25",
			expectError: true,
		},
		{
			name:        "invalid month",
			expiryDate:  "13/25",
			expectError: true,
		},
		{
			name:        "zero month",
			expiryDate:  "00/25",
			expectError: true,
		},
		{
			name:        "invalid year",
			expiryDate:  "12/100",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateExpiryDate(tt.expiryDate)
			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateCVV(t *testing.T) {
	tests := []struct {
		name        string
		cvv         string
		expectError bool
	}{
		{
			name:        "valid 3-digit cvv",
			cvv:         "123",
			expectError: false,
		},
		{
			name:        "valid 4-digit cvv",
			cvv:         "1234",
			expectError: false,
		},
		{
			name:        "too short",
			cvv:         "12",
			expectError: true,
		},
		{
			name:        "too long",
			cvv:         "12345",
			expectError: true,
		},
		{
			name:        "non-numeric",
			cvv:         "12a",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCVV(tt.cvv)
			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
