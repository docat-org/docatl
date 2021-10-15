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

var deleteCmd = &cobra.Command{
	Use:   "delete PROJECT VERSION",
	Short: "Delete documentation from a docat server",
	Long: `Delete documentation from a docat server.

Delete documentation:

	docatl delete myproject 1.0.0
`,
	Args: cobra.ExactArgs(2),
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
