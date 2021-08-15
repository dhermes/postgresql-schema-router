package server

import (
	"fmt"
)

const (
	// DefaultProxyPort is the default value for `Config.Port`.
	DefaultProxyPort = 5397
)

// Config represents the values needed to configure a server.
type Config struct {
	// ProxyPort is the port where the proxy should expose the server
	ProxyPort int
	// RemoteAddr is the address where the proxy should forward traffic. For
	// example `localhost:22089`
	RemoteAddr string
}

func (c Config) Validate() error {
	if c.ProxyPort == 0 {
		return fmt.Errorf("%w, ProxyPort is required", ErrInvalidConfiguration)
	}
	if c.RemoteAddr == "" {
		// TODO: Validate it's of the form `host:port` as well
		return fmt.Errorf("%w, RemoteAddr is required", ErrInvalidConfiguration)
	}
	return nil
}
