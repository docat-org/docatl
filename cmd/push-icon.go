package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var pushIconCmd = &cobra.Command{
	Use:   "push-icon [PROJECT] [ICON_PATH]",
	Short: "Push an icon for a project",
	Long: `Push an icon for a project.
Push an icon for a project:
	push-icon myproject /path/to/icon.png
	`,
	Args: cobra.ExactArgs(2),
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureHost()
	},
	Run: func(cmd *cobra.Command, args []string) {
		project, iconPath := args[0], args[1]

		err := docat.PushIcon(project, iconPath)

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Successfully pushed icon %s for project %s", iconPath, project)
	},
}

func init() {
	rootCmd.AddCommand(pushIconCmd)
}
