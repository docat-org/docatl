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
	"os"
	"path"

	util "github.com/docat-org/docatl/internal"
	docatl "github.com/docat-org/docatl/pkg"
	"github.com/spf13/cobra"
)

var additionalTag string

var pushCmd = &cobra.Command{
	Use:   "push DOCS PROJECT VERSION",
	Short: "Push documentation to a docat server",
	Long: `Push documentation to a docat server.

Upload documentation:

	docatl push ./docs.zip myproject 1.0.0 -t latest

Build & Upload documentation:

	docatl push ./docs/ myproject 1.0.0 -t latest

Upload documentation to specific docat server:

	docatl push --host localhost:8000 ./docs.zip myproject 1.0.0 -t latest
`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		docsPath, project, version := args[0], args[1], args[2]
		currentDir, err := os.Getwd()
		cobra.CheckErr(err)
		docsPath = path.Join(currentDir, docsPath)

		if util.IsDirectory(docsPath) {
			docsPath, err = docatl.Build(util.ResolvePath(docsPath), docatl.BuildMetadata{
				Project: project,
				Version: version,
			})
			cobra.CheckErr(err)
		}

		err = docat.Post(project, version, docsPath)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Successfully pushed documentation version %s to project %s", version, project)

		if additionalTag != "" {
			err = docat.Tag(project, version, additionalTag)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Successfully tagged version %s of project %s as %s", version, project, additionalTag)
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	rootCmd.PersistentFlags().StringVar(&additionalTag, "tag", "", "Additional Tag for this version")
}
