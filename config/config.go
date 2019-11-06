package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	ErrConfigNotLoaded = errors.New("configuration has not been loaded")
	configuration      *Config
)

type Config struct {
	Variables map[string]string `yaml:"variables"`
}

func NewConfigFromFile(configFile string) (err error) {
	configuration, err = parseConfigFile(configFile)
	return
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
	return &cfg, nil
}

func Get() (cfg *Config, err error) {
	if configuration == nil {
		return nil, ErrConfigNotLoaded
	}
	return configuration, nil
}
