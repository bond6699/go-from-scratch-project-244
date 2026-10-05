package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

var errInvalidPath = errors.New("invalid filepath")

type Parser interface {
	Parse([]byte) (*Model, error)
}

func isValidPath(path string) bool {
	path = strings.TrimSpace(path)
	if !strings.HasSuffix(path, ".json") && !strings.HasSuffix(path, ".yml") && !strings.HasSuffix(path, ".yaml") {
		return false
	}

	_, err := os.Stat(path)
	return err == nil
}

func getParser(path string) Parser {
	if strings.HasSuffix(path, "json") {
		return JSONParser{}
	}
	return YamlParser{}
}

func Parse(filepath string) (*Model, error) {
	empty := &Model{}
	if !isValidPath(filepath) {
		return empty, fmt.Errorf("%w: %s", errInvalidPath, filepath)
	}

	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return empty, err
	}

	parser := getParser(filepath)
	data, err := parser.Parse(bytes)
	if err != nil {
		return empty, err
	}

	return data, nil
}

func ToJSON(model *Model) ([]byte, error) {
	data, err := json.MarshalIndent(model, "", "  ")
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}
