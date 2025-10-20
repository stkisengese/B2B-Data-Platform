package database

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateKey   = errors.New("duplicate key")
	ErrConnection     = errors.New("database connection error")
)
