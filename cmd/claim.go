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
	"github.com/spf13/viper"
)

var claimCmd = &cobra.Command{
	Use:   "claim PROJECT",
	Short: "Claim a docat project",
	Long: `Claim a docat project.

Claim a project:

	docatl claim myproject
`,
	Args: cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		ensureHost()
	},
	Run: func(cmd *cobra.Command, args []string) {
		project := args[0]

		claim, err := docat.Claim(project)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Successfully claimed project %s. Store and use the following token: %s", project, claim.Token)

		writeToConfig, err := cmd.Flags().GetBool("write-to-config")
		cobra.CheckErr(err)

		if writeToConfig {
			configPath := viper.ConfigFileUsed()
			err = docatl.WriteConfig(configPath, docatl.Config{
				Host:   docat.Host,
				ApiKey: claim.Token,
			})
			if err != nil {
				log.Fatalf("unable to write claim to config: %s", err)
			}
			log.Printf("Updated config at '%s' with claim token", configPath)
		}
	},
}

func init() {
	rootCmd.AddCommand(claimCmd)

	claimCmd.Flags().BoolP("write-to-config", "w", false, "write claim token to config file")
}
