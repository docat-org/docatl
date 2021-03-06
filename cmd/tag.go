/*
Copyright © 2021 Timo Furrer <tuxtimo@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
