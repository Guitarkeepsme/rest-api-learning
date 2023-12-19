package storage

import "errors"

var (
	ErrNotFound  = errors.New("not found")
	ErrURLExists = errors.New("url already exists")
)
