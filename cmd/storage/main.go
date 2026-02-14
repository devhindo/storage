package main

import (
	"fmt"
	"os"

	"github.com/devhindo/storage/internal/auth"
	"github.com/devhindo/storage/internal/drive"
	"github.com/devhindo/storage/internal/tui"
)

func main() {
	client, err := auth.GetClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "authentication failed: %v\n", err)
		os.Exit(1)
	}

	srv, err := drive.NewService(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create drive service: %v\n", err)
		os.Exit(1)
	}

	if err := tui.Run(srv); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
