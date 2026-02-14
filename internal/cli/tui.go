package cli

import (
	"github.com/devhindo/storage/internal/tui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(tuiCmd)
}

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive terminal file browser",
	Long:  "Launch a full-screen interactive TUI for browsing Google Drive files.",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, err := newFileService()
		if err != nil {
			return err
		}

		return tui.Run(svc)
	},
}
