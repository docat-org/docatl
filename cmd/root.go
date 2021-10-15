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
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	docatl "github.com/timofurrer/docat-cli/pkg"

	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "docatl",
	Short: "Manage docat documentation easily",
	Long: `docatl - manage docat documentation, easily.

Upload documentation:

	docatl push ./docs.zip myproject 1.0.0 -t latest

Upload documentation to specific docat server:

	docatl push --host localhost:8000 ./docs.zip myproject 1.0.0 -t latest
`,
}

var docat docatl.Docat

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.docat-cli.yaml)")
	rootCmd.PersistentFlags().StringVar(&docat.Host, "host", "", "docat hostname (e.g. https://docat.company.com:8000)")
	rootCmd.PersistentFlags().StringVar(&docat.ApiKey, "api-key", "", "docat Api Key")
	err := rootCmd.MarkFlagRequired("host")
	if err != nil {
		log.Fatal("unable to mark flag `host` as required. Please create an upstream issue!")
	}
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".docat-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".docat-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
