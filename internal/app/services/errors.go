package services

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrEmptyCredentials  = errors.New("empty login or password")
)
