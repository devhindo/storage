package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list [folder-id]",
	Short: "List files and folders in a Google Drive folder",
	Long:  "List files and folders in the specified Google Drive folder. Defaults to the root folder.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		folderID := "root"
		if len(args) > 0 {
			folderID = args[0]
		}

		svc, err := newFileService()
		if err != nil {
			return err
		}

		entries, err := svc.ListFolder(context.Background(), folderID)
		if err != nil {
			return fmt.Errorf("failed to list folder: %w", err)
		}

		for _, e := range entries {
			icon := "   "
			if e.IsFolder {
				icon = "📁 "
			}
			fmt.Printf("%s%s\n", icon, e.Name)
		}

		return nil
	},
}
