package cli

import (
	"fmt"
	"os"

	"github.com/devhindo/storage/pkg/auth"
	"github.com/devhindo/storage/pkg/core"
	"github.com/devhindo/storage/pkg/storage/gdrive"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "storage",
	Short: "A terminal-based Google Drive file browser",
	Long:  "Storage is a CLI tool for browsing and managing files in Google Drive from your terminal.",
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// newFileService authenticates and creates a FileService backed by Google Drive.
func newFileService() (*core.FileService, error) {
	client, err := auth.GetClient()
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	backend, err := gdrive.New(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create drive backend: %w", err)
	}

	return core.NewFileService(backend), nil
}
