package generator

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type StringList []string

type Config struct {
	CursorType CursorType `yaml:"cursor_type,omitempty"`
	Connectify []string   `yaml:"connectify,omitempty"`
	Nodeify    []string   `yaml:"nodeify,omitempty"`
}

type CursorType struct {
	Type     string `yaml:"type,omitempty"`
	Nullable bool   `yaml:"nullable,omitempty"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{
		CursorType: CursorType{},
		Connectify: []string{},
		Nodeify:    []string{},
	}
	err = yaml.UnmarshalStrict(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
