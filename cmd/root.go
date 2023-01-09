package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	docatl "github.com/docat-org/docatl/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/spf13/viper"
)

const (
	envPrefix      = "DOCATL"
	configFileName = ".docatl"
	configFileType = "yaml"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "docatl",
	Short: "Manage docat documentation easily",
	Long: `docatl - manage docat documentation, easily.

Upload documentation:

	docatl push ./docs.zip myproject 1.0.0 -t latest

Upload documentation to specific docat server:

	docatl push --host https://localhost:8000 ./docs.zip myproject 1.0.0 -t latest
`,
}

var docat docatl.Docat

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	cwd, err := os.Getwd()
	cobra.CheckErr(err)
	defaultConfigPath := filepath.Join(cwd, fmt.Sprintf("%s.%s", configFileName, configFileType))

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", defaultConfigPath, "config file")
	rootCmd.PersistentFlags().StringVar(&docat.Host, "host", "", "docat hostname (e.g. https://docat.company.com:8000)")
	rootCmd.PersistentFlags().StringVar(&docat.ApiKey, "api-key", "", "docat Api Key")
}

func ensureHost() {
	if docat.Host == "" {
		log.Fatal("host setting is missing. Either use `--host <host>` or `DOCATL_HOST=<host>` or a config file with the `host:` field.")
	}
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		cwd, err := os.Getwd()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(cwd)
		viper.SetConfigType(configFileType)
		viper.SetConfigName(configFileName)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	setupEnv(rootCmd)
}

func setupEnv(cmd *cobra.Command) {
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	bindFlags(cmd, viper.GetViper())
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			err := v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
			cobra.CheckErr(err)
		}

		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			cobra.CheckErr(err)
		}
	})
}
