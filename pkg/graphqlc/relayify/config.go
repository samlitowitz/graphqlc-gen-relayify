package relayify

import (
	"encoding/json"
	"io/ioutil"
)

type StringList []string

type Config struct {
	CursorType *CursorType   `json:"cursorType,omitempty"`
	Connectify []*Connectify `json:"connectify,omitempty"`
	Nodeify    []string      `json:"nodeify,omitempty"`
}

type Connectify struct {
	Type   string            `json:"type,omitempty"`
	Fields []ConnectifyField `json:"fields,omitempty"`
}

type ConnectifyField struct {
	Type      string `json:"type,omitempty"`
	Field     string `json:"field,omitempty"`
	Overwrite bool   `json:"overwrite,omitempty"`
}

type CursorType struct {
	Type     string `json:"type,omitempty"`
	Nullable bool   `json:"nullable,omitempty"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{
		CursorType: &CursorType{},
		Connectify: []*Connectify{},
		Nodeify:    []string{},
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
