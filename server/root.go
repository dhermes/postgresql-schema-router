package server

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

// Run starts the PostgreSQL reverse proxy server.
func Run(c Config) error {
	err := c.Validate()
	if err != nil {
		return err
	}

	proxyAddr := fmt.Sprintf("localhost:%d", c.ProxyPort)
	addr, err := net.ResolveTCPAddr("tcp", proxyAddr)
	if err != nil {
		return err
	}

	// LOG-TODO: Setting up TCP proxy on %s\n", proxyAddr)
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	for {
		tc, err := listener.AcceptTCP()
		if err != nil {
			return err
		}

		// TODO: Use a channel here and a fixed set of goroutines to handle it
		go proxy(tc, c)
	}
}

// Execute runs the PostgreSQL reverse proxy server as a command line (CLI)
// application.
func Execute() error {
	c := Config{}
	cmd := &cobra.Command{
		Use:           "postgresql-schema-router",
		Short:         "PostgreSQL Reverse Proxy",
		Long:          "PostgreSQL Reverse Proxy\n\nForward Queries Based on Schema.",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(_ *cobra.Command, _ []string) error {
			return Run(c)
		},
	}

	cmd.PersistentFlags().IntVar(
		&c.ProxyPort,
		"port",
		DefaultProxyPort,
		"The port where the proxy should expose the server",
	)
	cmd.PersistentFlags().StringVar(
		&c.RemoteAddr,
		"remote",
		"",
		"The remote address  where the proxy should forward traffic (e.g. localhost:22089)",
	)

	return cmd.Execute()
}
