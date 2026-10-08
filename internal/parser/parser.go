package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var errInvalidPath = errors.New("invalid filepath")

type Parser interface {
	Parse([]byte) (map[string]any, error)
}

type JSONParser struct{}

func (jp JSONParser) Parse(bytes []byte) (map[string]any, error) {
	var jsonData map[string]any

	err := json.Unmarshal(bytes, &jsonData)
	if err != nil {
		return jsonData, err
	}

	return jsonData, nil
}

type YamlParser struct{}

func (yp YamlParser) Parse(bytes []byte) (map[string]any, error) {
	var yamlData map[string]any

	err := yaml.Unmarshal(bytes, &yamlData)
	if err != nil {
		return yamlData, err
	}

	return yamlData, nil
}

func isValidPath(path string) bool {
	path = strings.TrimSpace(path)
	if !strings.HasSuffix(path, ".json") && !strings.HasSuffix(path, ".yml") && !strings.HasSuffix(path, ".yaml") {
		return false
	}

	_, err := os.Stat(path)
	fmt.Print(err)
	return err == nil
}

func getParser(path string) Parser {
	if strings.HasSuffix(path, "json") {
		return JSONParser{}
	}
	return YamlParser{}
}

func Parse(filepath string) (map[string]any, error) {
	var result map[string]any

	if !isValidPath(filepath) {
		return result, fmt.Errorf("%w: %s", errInvalidPath, filepath)
	}

	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return result, err
	}

	parser := getParser(filepath)
	result, err = parser.Parse(bytes)
	if err != nil {
		return result, err
	}

	return result, nil
}
