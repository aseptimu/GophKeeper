package models

import "time"

type DataType string

const (
	DataTypeLoginPassword DataType = "login_password"
	DataTypeText          DataType = "text"
	DataTypeBinary        DataType = "binary"
	DataTypeBankCard      DataType = "bank_card"
)

type DataItem struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      DataType  `json:"type"`
	Data      string    `json:"data"`
	Metadata  string    `json:"metadata"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginPasswordData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type BankCardData struct {
	CardNumber string `json:"card_number"`
	ExpiryDate string `json:"expiry_date"`
	CVV        string `json:"cvv"`
	Cardholder string `json:"cardholder"`
}

type BinaryData struct {
	FilePath    string `json:"file_path"`
	ContentType string `json:"content_type"`
	FileName    string `json:"file_name"`
	FileSize    int64  `json:"file_size"`
}
