package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	errReadFile            = errors.New("failed read file")
	errCreateParserFile    = errors.New("failed create parser")
	errUnsupportedFileType = errors.New("not supported file type")
	errParse               = errors.New("failed parse")

	errDecodeJSON = errors.New("failed to decode JSON")
	errDecodeYAML = errors.New("failed to decode YAML")
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
	ext := filepath.Ext(path)

	switch ext {
	case ".json":
		return JSONParser{}, nil
	case ".yml", ".yaml":
		return YamlParser{}, nil
	default:
		return nil, fmt.Errorf("%w %s", errUnsupportedFileType, path)
	}
}

func Parse(path string) (map[string]any, error) {
	result := make(map[string]any)

	bytes, err := os.ReadFile(path)
	if err != nil {
		return result, fmt.Errorf("%w %q: %w", errReadFile, path, err)
	}

	parser, err := NewParser(path)
	if err != nil {
		return result, fmt.Errorf("%w %q: %w", errCreateParserFile, path, err)
	}

	result, err = parser.Parse(bytes)
	if err != nil {
		return result, fmt.Errorf("%w %q: %w", errParse, path, err)
	}

	return result, nil
}
