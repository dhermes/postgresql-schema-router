package server

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Run starts the PostgreSQL reverse proxy server.
func Run(c Config) error {
	fmt.Printf("Hello world; port=%d\n", c.Port)
	return nil
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
		&c.Port,
		"port",
		DefaultPort,
		"The port where the proxy should expose the server",
	)

	return cmd.Execute()
}
