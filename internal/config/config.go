package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Telemetry struct {
		Interfaces []string `yaml:"interfaces"`
		Protocols  struct {
			HTTP bool `yaml:"http"`
			DNS  bool `yaml:"dns"`
			ICMP bool `yaml:"icmp"`
		} `yaml:"protocols"`
	} `yaml:"telemetry"`

	Logging struct {
		Level string `yaml:"level"`
		Path  string `yaml:"path"`
	} `yaml:"logging"`

	Backend struct {
		Host   string `yaml:"host"`
		Port   int    `yaml:"port"`
		UseTLS bool   `yaml:"use_tls"`
	} `yaml:"backend"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
