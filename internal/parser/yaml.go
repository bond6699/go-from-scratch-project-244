package parser

import (
	"gopkg.in/yaml.v3"
)

type YamlParser struct{}

func (yp YamlParser) Parse(bytes []byte) (*Model, error) {
	var yamlData map[string]any

	err := yaml.Unmarshal(bytes, &yamlData)
	if err != nil {
		return &Model{}, err
	}

	data := New("yaml", yamlData)

	return data, nil
}
