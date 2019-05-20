package generator

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type StringList []string

type Config struct {
	Connectify []string `yaml:"connectify,omitempty"`
	Nodeify    []string `yaml:"nodeify,omitempty"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{
		Connectify: []string{},
		Nodeify:    []string{},
	}
	err = yaml.UnmarshalStrict(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
