package server

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Run starts the PostgreSQL reverse proxy server.
func Run(port int) error {
	fmt.Printf("Hello world; port=%d\n", port)
	return nil
}

// Execute runs the PostgreSQL reverse proxy server as a command line (CLI)
// application.
func Execute() error {
	port := 5397
	cmd := &cobra.Command{
		Use:           "postgresql-schema-router",
		Short:         "PostgreSQL Reverse Proxy",
		Long:          "PostgreSQL Reverse Proxy\n\nForward Queries Based on Schema.",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(_ *cobra.Command, _ []string) error {
			return Run(port)
		},
	}

	cmd.PersistentFlags().IntVar(
		&port,
		"port",
		5397,
		"The port where the proxy should expose the server",
	)

	return cmd.Execute()
}
