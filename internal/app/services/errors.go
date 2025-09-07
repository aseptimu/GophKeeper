package services

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrEmptyCredentials  = errors.New("empty login or password")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidDataType   = errors.New("invalid data type")
	ErrInvalidDataFormat = errors.New("invalid data format")
)
