package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	errReadFile            = fmt.Errorf("failed read file")
	errCreateParserFile    = fmt.Errorf("failed create parser")
	errUnsupportedFileType = fmt.Errorf("not supported file type")
	errParse               = fmt.Errorf("failed parse")

	errDecodeJSON = fmt.Errorf("failed to decode JSON")
	errDecodeYAML = fmt.Errorf("failed to decode YAML")
)

type Parser interface {
	Parse([]byte) (map[string]any, error)
}

type JSONParser struct{}

func (jp JSONParser) Parse(bytes []byte) (map[string]any, error) {
	jsonData := make(map[string]any)

	err := json.Unmarshal(bytes, &jsonData)
	if err != nil {
		return jsonData, fmt.Errorf("%w: %w", errDecodeJSON, err)
	}

	return jsonData, nil
}

type YamlParser struct{}

func (yp YamlParser) Parse(bytes []byte) (map[string]any, error) {
	yamlData := make(map[string]any)

	err := yaml.Unmarshal(bytes, &yamlData)
	if err != nil {
		return yamlData, fmt.Errorf("%w: %w", errDecodeYAML, err)
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

	return nil, fmt.Errorf("%w %s", errUnsupportedFileType, path)
}

func Parse(filepath string) (map[string]any, error) {
	result := make(map[string]any)

	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return result, fmt.Errorf("%w: %w", errReadFile, err)
	}

	parser, err := NewParser(filepath)
	if err != nil {
		return result, fmt.Errorf("%w: %w", errCreateParserFile, err)
	}

	result, err = parser.Parse(bytes)
	if err != nil {
		return result, fmt.Errorf("%w: %w", errParse, err)
	}

	return result, nil
}
