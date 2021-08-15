package server

import (
	"errors"
)

var (
	// ErrInvalidConfiguration is the error returned when a `Config` fails
	// validation.
	ErrInvalidConfiguration = errors.New("invalid configuration")
)
