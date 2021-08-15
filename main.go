package main

import (
	"fmt"
	"os"

	"github.com/dhermes/postgresql-schema-router/server"
)

func main() {
	err := server.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
