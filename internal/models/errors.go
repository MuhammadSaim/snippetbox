package models

import "errors"

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("model: invalid credentials")
	ErrDuplicateEmail     = errors.New("model: duplicate email")
)
