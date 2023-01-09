package cmd

import (
	"log"

	util "github.com/docat-org/docatl/internal"
	docatl "github.com/docat-org/docatl/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pushCmd = &cobra.Command{
	Use:   "push DOCS [PROJECT [VERSION]]",
	Short: "Push documentation to a docat server",
	Long: `Push documentation to a docat server.

Upload documentation:

	docatl push ./docs.zip myproject 1.0.0 -t latest

Build & Upload documentation:

	docatl push ./docs/ myproject 1.0.0 -t latest

Upload documentation to specific docat server:

	docatl push --host https://localhost:8000 ./docs.zip myproject 1.0.0 -t latest
`,
	Args: cobra.RangeArgs(1, 3),
	Run: func(cmd *cobra.Command, args []string) {
		docsPath := util.ResolvePath(args[0])
		project := viper.GetString("project")
		version := viper.GetString("version")

		unpackArgs := func() (string, string) {
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

			return project, version
		}

		if util.IsDirectory(docsPath) {

			project, version = unpackArgs()

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

			project, version = unpackArgs()
		}

		ensureHost()

		err := docat.Post(project, version, docsPath)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Successfully pushed documentation version %s to project %s", version, project)

		tags, err := cmd.Flags().GetStringSlice("tag")
		cobra.CheckErr(err)
		for _, tag := range tags {
			err = docat.Tag(project, version, tag)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Successfully tagged version %s of project %s as %s", version, project, tag)
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.PersistentFlags().StringSliceP("tag", "t", []string{}, "Additional Tag for this version (repeatable)")

	setupEnv(pushCmd)
}
