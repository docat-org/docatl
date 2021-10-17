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

	docatl "github.com/docat-org/docatl/pkg"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a documentation artifact to push to a docat server",
	Long: `Build a documentation artifact to push to a docat server.

The 'docatl build' command can be used to create a documentation artifact
which can then be used with 'docatl push' to upload to a docat server.

The documentation artifact is *just* a ZIP archive.

Example:

	docatl build docs/
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		docsPath := args[0]
		project, err := cmd.Flags().GetString("project")
		cobra.CheckErr(err)
		version, err := cmd.Flags().GetString("version")
		cobra.CheckErr(err)

		outputPath, err := docatl.Build(docsPath, docatl.BuildMetadata{
			Project: project,
			Version: version,
		})
		if err != nil {
			log.Fatalf("unable to build documentation: %s", err)
		}
		log.Printf("Successfully build documentation, stored at: %s", outputPath)
		log.Printf("Push documentation with: `docatl push %s`", outputPath)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringP("project", "p", "", "the name of the docat project")
	buildCmd.Flags().StringP("version", "v", "", "the version of this documentation")

	setupEnv(buildCmd)
}
