package docatl

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host   string `yaml:"host"`
	ApiKey string `yaml:"api-key"`
}

func WriteConfig(configPath string, config Config) error {
	doc, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("unable to marshal config '%v' to YAML: %w", config, err)
	}

	err = ioutil.WriteFile(configPath, doc, 0)
	if err != nil {
		return fmt.Errorf("unable to write config to '%s': %w", configPath, err)
	}

	return nil
}
