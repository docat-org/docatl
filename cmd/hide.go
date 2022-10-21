package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var hideCmd = &cobra.Command{
	Use:   "hide [PROJECT] [VERSION]",
	Short: "Hide a project version.",
	Long: `Hide a project version. 
	This will remove the project version from the version select as well as the search index.

Hide a project version:
	hide myproject 1.0.0
	`,
	Args: cobra.ExactArgs(2),
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureHost()
	},
	Run: func(cmd *cobra.Command, args []string) {
		project, version := args[0], args[1]

		err := docat.HideOrShowVersion(project, version, true)

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Successfully hid version %s of project %s", version, project)
	},
}

func init() {
	rootCmd.AddCommand(hideCmd)
}
