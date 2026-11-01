package parser

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsers(t *testing.T) {
	tests := []struct {
		name     string
		parser   Parser
		filepath string
		want     map[string]any
		wantErr  error
	}{
		{
			"JSONParser with simple value",
			JSONParser{},
			"testdata/json/simple.json",
			simpleJSONResult(),
			nil,
		},
		{
			"JSONParser with complex value",
			JSONParser{},
			"testdata/json/complex.json",
			complexJSONResult(),
			nil,
		},
		{
			"JSONParser with empty value",
			JSONParser{},
			"testdata/json/empty.json",
			emptyResult(),
			nil,
		},
		{
			"JSONParser with invalid value",
			JSONParser{},
			"testdata/json/invalid.json",
			emptyResult(),
			errDecodeJSON,
		},
		{
			"YAMLParser with simple value",
			YamlParser{},
			"testdata/yaml/simple.yaml",
			simpleYamlResult(),
			nil,
		},
		{
			"YAMLParser with complex value",
			YamlParser{},
			"testdata/yaml/complex.yaml",
			complexYamlResult(),
			nil,
		},
		{
			"YAMLParser with empty value",
			YamlParser{},
			"testdata/yaml/empty.yaml",
			emptyResult(),
			nil,
		},
		{
			"YAMLParser with invalid value",
			YamlParser{},
			"testdata/yaml/invalid.yaml",
			emptyResult(),
			errDecodeYAML,
		},
	}
	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			bytes, err := ReadFile(td.filepath)
			require.NoError(t, err)

			got, err := td.parser.Parse(bytes)
			require.ErrorIs(t, err, td.wantErr)
			assert.Equal(t, got, td.want)
		})
	}
}

func TestNewParser(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Parser
		wantErr error
	}{
		{"NewParser with json file", "file/path/data.json", JSONParser{}, nil},
		{"NewParser with yaml file (with .yaml ext)", "file/path/data.yaml", YamlParser{}, nil},
		{"NewParser with yaml file (with .yml ext)", "file/path/data.yml", YamlParser{}, nil},
		{"NewParser with unsupported file type", "file/path/data.txt", nil, errUnsupportedFileType},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			got, err := NewParser(td.input)
			assert.ErrorIs(t, err, td.wantErr)
			require.Equal(t, got, td.want)
		})
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]any
		want  map[string]any
	}{
		{
			"Normalize int, float types",
			map[string]any{
				"dataInt":    234,
				"dataInt64":  int64(100),
				"dataUint64": uint64(233),
				"dataFloat":  334.32,
			},
			map[string]any{
				"dataInt":    234.0,
				"dataInt64":  100.0,
				"dataUint64": 233.0,
				"dataFloat":  334.32,
			},
		},
		{
			"Normalize slice []any",
			map[string]any{
				"numbers": []any{1, 6, 9, 22.6},
				"others":  []any{"hello", "world", nil},
			},
			map[string]any{
				"numbers": []any{1.0, 6.0, 9.0, 22.6},
				"others":  []any{"hello", "world", nil},
			},
		},
		{
			"Normalize null value (must be dont changed)",
			map[string]any{
				"dataNull": nil,
			},
			map[string]any{
				"dataNull": nil,
			},
		},
		{
			"Normalize complex value",
			map[string]any{
				"common": map[string]any{
					"setting1": "Value 1",
					"setting2": 200,
					"setting3": true,
					"setting6": map[string]any{
						"key": "value",
						"doge": map[string]any{
							"wow": "",
						},
					},
				},
				"group1": map[string]any{
					"baz": "bas",
					"foo": "bar",
					"nest": map[string]any{
						"key": "value",
					},
					"secret": nil,
				},
				"group2": map[string]any{
					"abc": 12345,
					"deep": map[string]any{
						"id": 45,
					},
				},
			},
			map[string]any{
				"common": map[string]any{
					"setting1": "Value 1",
					"setting2": 200.0,
					"setting3": true,
					"setting6": map[string]any{
						"key": "value",
						"doge": map[string]any{
							"wow": "",
						},
					},
				},
				"group1": map[string]any{
					"baz": "bas",
					"foo": "bar",
					"nest": map[string]any{
						"key": "value",
					},
					"secret": nil,
				},
				"group2": map[string]any{
					"abc": 12345.0,
					"deep": map[string]any{
						"id": 45.0,
					},
				},
			},
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			require.Equal(t, td.want, Normalize(td.input))
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		want     map[string]any
		wantErr  error
	}{
		{
			"Parse simple JSON file",
			"testdata/json/simple.json",
			simpleJSONResult(),
			nil,
		},
		{
			"Parse complex JSON file",
			"testdata/json/complex.json",
			complexJSONResult(),
			nil,
		},
		{
			"Parse empty JSON file",
			"testdata/json/empty.json",
			emptyResult(),
			nil,
		},
		{
			"Parse invalid JSON file",
			"testdata/json/invalid.json",
			emptyResult(),
			errDecodeJSON,
		},
		{
			"Parse valid simple YAML file",
			"testdata/yaml/simple.yaml",
			simpleYamlResult(),
			nil,
		},
		{
			"Parse valid complex YAML file",
			"testdata/yaml/complex.yaml",
			complexYamlResult(),
			nil,
		},
		{
			"Parse empty YAML file",
			"testdata/yaml/empty.yaml",
			emptyResult(),
			nil,
		},
		{
			"Parse invalid YAML file",
			"testdata/yaml/invalid.yaml",
			emptyResult(),
			errDecodeYAML,
		},
		{
			"Parse unsupported file",
			"testdata/unsupported.txt",
			emptyResult(),
			errUnsupportedFileType,
		},
		{
			"Parse invalid filepath",
			"testdata/unsupp",
			emptyResult(),
			errReadFile,
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			got, err := Parse(td.filepath)
			assert.ErrorIs(t, err, td.wantErr)
			require.Equal(t, td.want, got)
		})
	}
}

func simpleJSONResult() map[string]any {
	return map[string]any{
		"host":    "hexlet.io",
		"timeout": 50.0,
		"proxy":   "123.234.53.22",
		"follow":  false,
		"secret":  nil,
	}
}

func complexJSONResult() map[string]any {
	return map[string]any{
		"common": map[string]any{
			"setting1": "Value 1",
			"setting2": 200.0,
			"setting3": true,
			"setting6": map[string]any{
				"key": "value",
				"doge": map[string]any{
					"wow": "",
				},
			},
		},
		"group1": map[string]any{
			"baz": "bas",
			"foo": "bar",
			"nest": map[string]any{
				"key": "value",
			},
			"secret": nil,
		},
		"group2": map[string]any{
			"abc": 12345.0,
			"deep": map[string]any{
				"id": 45.0,
			},
		},
	}
}

func simpleYamlResult() map[string]any {
	return map[string]any{
		"host":    "hexlet.io",
		"timeout": 50,
		"proxy":   "123.234.53.22",
		"follow":  false,
		"secret":  nil,
	}
}

func complexYamlResult() map[string]any {
	return map[string]any{
		"common": map[string]any{
			"setting1": "Value 1",
			"setting2": 200,
			"setting3": true,
			"setting6": map[string]any{
				"key": "value",
				"doge": map[string]any{
					"wow": "",
				},
			},
		},
		"group1": map[string]any{
			"baz": "bas",
			"foo": "bar",
			"nest": map[string]any{
				"key": "value",
			},
			"secret": nil,
		},
		"group2": map[string]any{
			"abc": 12345,
			"deep": map[string]any{
				"id": 45,
			},
		},
	}
}

func emptyResult() map[string]any {
	return map[string]any{}
}

func ReadFile(filepath string) ([]byte, error) {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		err = fmt.Errorf("%s %w: %w", filepath, errReadFile, err)
	}
	return bytes, err
}
