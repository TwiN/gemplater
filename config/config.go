package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var (
	ErrConfigNotLoaded = errors.New("configuration has not been loaded")
	configuration      *Config
)

type Config struct {
	Variables map[string]string `yaml:"variables"`
}

// Creates a new Config struct
// If the configuration file passed doesn't exist, an empty Config struct will be created instead.
func NewConfig(configFile string) (*Config, error) {
	configuration, err := parseConfigFile(configFile)
	// Check if the file doesn't exist. If it doesn't, then return en empty Config
	if os.IsNotExist(err) {
		configuration = &Config{Variables: make(map[string]string)}
		return configuration, nil
	}
	return configuration, err
}

func parseConfigFile(configFile string) (*Config, error) {
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}
	if cfg.Variables == nil {
		cfg.Variables = make(map[string]string)
	}
	return &cfg, nil
}

func Get() (cfg *Config, err error) {
	if configuration == nil {
		return nil, ErrConfigNotLoaded
	}
	return configuration, nil
}
