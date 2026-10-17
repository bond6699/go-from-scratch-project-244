package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Parser interface {
	Parse([]byte) (map[string]any, error)
}

type JSONParser struct{}

func (jp JSONParser) Parse(bytes []byte) (map[string]any, error) {
	jsonData := make(map[string]any)

	err := json.Unmarshal(bytes, &jsonData)
	if err != nil {
		return jsonData, err
	}

	return jsonData, nil
}

type YamlParser struct{}

func (yp YamlParser) Parse(bytes []byte) (map[string]any, error) {
	yamlData := make(map[string]any)

	err := yaml.Unmarshal(bytes, &yamlData)
	if err != nil {
		return yamlData, err
	}

	return yamlData, nil
}

func NewParser(path string) (Parser, error) {
	if strings.HasSuffix(path, ".json") {
		return JSONParser{}, nil
	}

	if strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml") {
		return YamlParser{}, nil
	}

	return nil, fmt.Errorf("not supported file '%s'", path)
}

func Parse(filepath string) (map[string]any, error) {
	result := make(map[string]any)

	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return result, err
	}

	parser, err := NewParser(filepath)
	if err != nil {
		return result, err
	}

	result, err = parser.Parse(bytes)
	if err != nil {
		return result, err
	}

	return result, nil
}
