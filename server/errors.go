package server

import (
	"errors"
	"net"

	multierror "github.com/hashicorp/go-multierror"
)

var (
	// ErrInvalidConfiguration is the error returned when a `Config` fails
	// validation.
	ErrInvalidConfiguration = errors.New("invalid configuration")
	// ErrPacketTooLarge is the error returned when a TCP packet from a read
	// is too large.
	ErrPacketTooLarge = errors.New("packet too large")
)

func appendErrs(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}
	combined := multierror.Append(errs[0], errs[1:]...)
	if len(combined.Errors) == 0 {
		return nil
	}
	return combined
}

func isTimeout(err error) bool {
	noe, ok := err.(*net.OpError)
	if !ok {
		return false
	}

	return noe.Timeout()
}
