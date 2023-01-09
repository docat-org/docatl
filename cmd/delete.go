package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete PROJECT VERSION",
	Short: "Delete documentation from a docat server",
	Long: `Delete documentation from a docat server.

Delete documentation:

	docatl delete myproject 1.0.0
`,
	Args: cobra.ExactArgs(2),
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureHost()
	},
	Run: func(cmd *cobra.Command, args []string) {
		project, version := args[0], args[1]

		err := docat.Delete(project, version)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Successfully deleted version %s of project %s", version, project)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
