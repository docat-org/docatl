package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var updateIndexCmd = &cobra.Command{
	Use:   "update-index",
	Short: "Force re-creation of the search index.",
	Long: `Force re-creation of the search index.
Update the search index:
	update-index
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureHost()
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := docat.UpdateIndex()

		if err != nil {
			log.Fatal(err)
		}

		log.Print("Successfully updated search index")
	},
}

func init() {
	rootCmd.AddCommand(updateIndexCmd)
}
