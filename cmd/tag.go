package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag PROJECT VERSION TAG..",
	Short: "Tag an existing documentation on a docat server",
	Long: `Tag an existing documentation on a docat server.

Tag documentation:

	docatl tag --host https://localhost:8000 myproject 1.0.0 latest
`,
	Args: cobra.MinimumNArgs(3),
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureHost()
	},
	Run: func(cmd *cobra.Command, args []string) {
		project, version, tags := args[0], args[1], args[2:]

		for _, tag := range tags {
			err := docat.Tag(project, version, tag)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Successfully tagged version %s of project %s as %s", version, project, tag)
		}
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
}
