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

	util "github.com/docat-org/docatl/internal"
	docatl "github.com/docat-org/docatl/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var additionalTag string

var pushCmd = &cobra.Command{
	Use:   "push DOCS [PROJECT [VERSION]]",
	Short: "Push documentation to a docat server",
	Long: `Push documentation to a docat server.

Upload documentation:

	docatl push ./docs.zip myproject 1.0.0 -t latest

Build & Upload documentation:

	docatl push ./docs/ myproject 1.0.0 -t latest

Upload documentation to specific docat server:

	docatl push --host localhost:8000 ./docs.zip myproject 1.0.0 -t latest
`,
	Args: cobra.RangeArgs(1, 3),
	Run: func(cmd *cobra.Command, args []string) {
		docsPath := util.ResolvePath(args[0])
		project := viper.GetString("project")
		version := viper.GetString("version")

		if util.IsDirectory(docsPath) {
			if project == "" {
				if len(args) < 2 {
					log.Fatalf("when PROJECT is not given, the DOCATL_PROJECT variable must contain it")
				}
				project = args[1]
			}

			if version == "" {
				if len(args) < 3 {
					log.Fatalf("when VERSION is not given, the DOCATL_VERSION variable must contain it")
				} else {
					version = args[2]
				}
			}

			docsPathBuilt, err := docatl.Build(docsPath, docatl.BuildMetadata{
				Host:    docat.Host,
				Project: project,
				Version: version,
			})
			if err != nil {
				log.Fatal(err)
			}
			docsPath = docsPathBuilt
		} else {
			meta, err := docatl.ExtractMetadata(docsPath)
			if err != nil {
				log.Fatal(err)
			}

			if meta.Host != "" {
				docat.Host = meta.Host
			}

			project = meta.Project
			version = meta.Version

			if project == "" || version == "" {
				log.Fatal("when PROJECT and VERSION are not given, the `.docatl.meta.yaml` file in the archive must contain them. Use `docatl build` to make sure it exists")
			}
		}

		ensureHost()

		err := docat.Post(project, version, docsPath)
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
	pushCmd.PersistentFlags().StringVar(&additionalTag, "tag", "", "Additional Tag for this version")

	setupEnv(pushCmd)
}
