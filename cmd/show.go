package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [PROJECT] [VERSION]",
	Short: "Show a hidden project version.",
	Long: `Show a hidden project version. 
	This will make the project version reappear in the version select as well as the search index.

Show a project version:
	show myproject 1.0.0
	`,
	Args: cobra.ExactArgs(2),
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureHost()
	},
	Run: func(cmd *cobra.Command, args []string) {
		project, version := args[0], args[1]

		err := docat.HideOrShowVersion(project, version, false)

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Successfully undid hiding version %s of project %s", version, project)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
