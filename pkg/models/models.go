package models
import (
"errors"
)
var (
	ErrNoRecord = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

type User struct {
	ID int
	Name string
	Email string
	Password []byte
}

type CustomGPT struct {
	ID int
	SystemName string
	SystemPrompt string
}