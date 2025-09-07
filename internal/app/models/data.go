// Package models contains data structures used throughout the application.
package models

import "time"

// DataType represents the type of data stored in the system.
type DataType string

const (
	// DataTypeLoginPassword represents login/password pairs.
	DataTypeLoginPassword DataType = "login_password"
	// DataTypeText represents arbitrary text data.
	DataTypeText DataType = "text"
	// DataTypeBinary represents binary file data.
	DataTypeBinary DataType = "binary"
	// DataTypeBankCard represents bank card information.
	DataTypeBankCard DataType = "bank_card"
)

// DataItem represents a single piece of user data stored in the system.
type DataItem struct {
	ID        string    `json:"id"`         // Unique identifier for the data item
	UserID    string    `json:"user_id"`    // ID of the user who owns this data
	Type      DataType  `json:"type"`       // Type of data (login_password, text, binary, bank_card)
	Data      string    `json:"data"`       // The actual data content (JSON for structured data)
	Metadata  string    `json:"metadata"`   // Optional metadata describing the data
	CreatedAt time.Time `json:"created_at"` // When the item was created
	UpdatedAt time.Time `json:"updated_at"` // When the item was last updated
}

// LoginPasswordData represents login and password information.
type LoginPasswordData struct {
	Login    string `json:"login"`    // Username or login identifier
	Password string `json:"password"` // Password for the login
}

// BankCardData represents bank card information.
type BankCardData struct {
	CardNumber string `json:"card_number"` // Bank card number
	ExpiryDate string `json:"expiry_date"` // Card expiry date in MM/YY format
	CVV        string `json:"cvv"`         // Card verification value
	Cardholder string `json:"cardholder"`  // Name of the cardholder
}

// BinaryData represents information about a binary file.
type BinaryData struct {
	FilePath    string `json:"file_path"`    // Path to the stored file
	ContentType string `json:"content_type"` // MIME type of the file
	FileName    string `json:"file_name"`    // Original filename
	FileSize    int64  `json:"file_size"`    // Size of the file in bytes
}
