package qianmo

import "errors"

var (
	// ErrNotFound is returned when the requested item is not found.
	ErrNotFound = errors.New("Not found")
)
