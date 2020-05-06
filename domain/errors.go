package domain

import "github.com/pkg/errors"

var (
	// ErrNotFound is not found error
	ErrNotFound = errors.New("not found")
	// ErrBadRequest is bad request error
	ErrBadRequest = errors.New("bad request")
)
