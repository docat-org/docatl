package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:   "rename [PROJECT] [NEW_NAME]",
	Short: "Rename a project",
	Long: `Rename a project.
Rename a project:
	rename myproject newname
	`,
	Args: cobra.ExactArgs(2),
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureHost()
	},
	Run: func(cmd *cobra.Command, args []string) {
		project, newName := args[0], args[1]

		err := docat.Rename(project, newName)

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Successfully renamed project %s to %s", project, newName)
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
