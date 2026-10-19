package parser

import (
	"os"
	"reflect"
	"testing"
)

func TestParsers(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		parser   Parser
		want     map[string]any
		wantErr  bool
	}{
		{
			"valid simple JSON value",
			"testdata/json/simple.json",
			JSONParser{},
			map[string]any{
				"host":    "hexlet.io",
				"timeout": 50.0,
				"proxy":   "123.234.53.22",
				"follow":  false,
				"secret":  nil,
			},
			false,
		},
		{
			"valid complex JSON value",
			"testdata/json/complex.json",
			JSONParser{},
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
			false,
		},
		{
			"invalid JSON value",
			"testdata/json/invalid.json",
			JSONParser{},
			map[string]any{},
			true,
		},
		{
			"empty",
			"testdata/json/empty.json",
			JSONParser{},
			map[string]any{},
			false,
		},
		{
			"simple YAML value",
			"testdata/yaml/simple.yaml",
			YamlParser{},
			map[string]any{
				"host":    "hexlet.io",
				"timeout": 50,
				"proxy":   "123.234.53.22",
				"follow":  false,
				"secret":  nil,
			},
			false,
		},
		{
			"complex YAML value",
			"testdata/yaml/complex.yaml",
			YamlParser{},
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
			false,
		},
		{
			"invalid YAML value",
			"testdata/yaml/invalid.yaml",
			YamlParser{},
			map[string]any{},
			true,
		},
		{
			"empty YAML",
			"testdata/yaml/empty.yaml",
			YamlParser{},
			map[string]any{},
			false,
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			bytes, err := os.ReadFile(td.filepath)
			if err != nil {
				assertError(t, false, err)
			}

			got, err := td.parser.Parse(bytes)
			assertError(t, td.wantErr, err)

			if !td.wantErr {
				assertParserValue(t, td.want, got)
			}
		})
	}
}

func TestNewParser(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Parser
		wantErr bool
	}{
		{"json file", "file/path/data.json", JSONParser{}, false},
		{"yaml file", "file/path/data.yaml", YamlParser{}, false},
		{"yml file", "file/path/data.yml", YamlParser{}, false},
		{"unsupported file", "file/path/data.txt", nil, true},
		{"empty path", "", nil, true},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			got, err := NewParser(td.input)
			assertError(t, td.wantErr, err)

			if got != td.want {
				t.Fatalf("want %T, got %T", td.want, got)
			}
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
			"int, float types",
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
			"slice []any",
			map[string]any{
				"numbers": []any{1, 6, 9, 22.6},
				"others":  []any{"hello", "world", nil},
			},
			map[string]any{
				"numbers": []any{1.0, 6.0, 9.0, 22.6},
				"others":  []any{"hello", "world", "null"},
			},
		},
		{
			"null value",
			map[string]any{
				"dataNull": nil,
			},
			map[string]any{
				"dataNull": "null",
			},
		},
		{
			"mixed value",
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
					"secret": "null",
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
			assertParserValue(t, td.want, Normalize(td.input))
		})
	}
}

func TestParseAndNormalize(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		want     map[string]any
		wantErr  bool
	}{
		{
			"valid simple JSON file",
			"testdata/json/simple.json",
			map[string]any{
				"host":    "hexlet.io",
				"timeout": 50.0,
				"proxy":   "123.234.53.22",
				"follow":  false,
				"secret":  "null",
			},
			false,
		},
		{
			"valid complex JSON file",
			"testdata/json/complex.json",
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
					"secret": "null",
				},
				"group2": map[string]any{
					"abc": 12345.0,
					"deep": map[string]any{
						"id": 45.0,
					},
				},
			},
			false,
		},
		{
			"invalid JSON file",
			"testdata/json/invalid.json",
			map[string]any{},
			true,
		},
		{
			"empty JSON file",
			"testdata/json/empty.json",
			map[string]any{},
			false,
		},
		{
			"valid simple YAML file",
			"testdata/yaml/simple.yaml",
			map[string]any{
				"host":    "hexlet.io",
				"timeout": 50.0,
				"proxy":   "123.234.53.22",
				"follow":  false,
				"secret":  "null",
			},
			false,
		},
		{
			"valid complex YAML file",
			"testdata/yaml/complex.yaml",
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
					"secret": "null",
				},
				"group2": map[string]any{
					"abc": 12345.0,
					"deep": map[string]any{
						"id": 45.0,
					},
				},
			},
			false,
		},
		{
			"invalid YAML file",
			"testdata/yaml/invalid.yaml",
			map[string]any{},
			true,
		},
		{
			"empty YAML file",
			"testdata/yaml/empty.yaml",
			map[string]any{},
			false,
		},
		{
			"unsupported file",
			"testdata/unsupported.txt",
			map[string]any{},
			true,
		},
		{
			"invalid filepath",
			"testdata/unsupp",
			map[string]any{},
			true,
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			got, err := Parse(td.filepath)
			assertError(t, td.wantErr, err)

			if !td.wantErr {
				assertParserValue(t, td.want, Normalize(got))
			}
		})
	}
}

func assertError(t *testing.T, wantErr bool, got error) {
	t.Helper()

	if wantErr && got == nil {
		t.Fatal("expected an error, got nil")
	}

	if !wantErr && got != nil {
		t.Fatalf("unexpected error: %v", got)
	}
}

func assertParserValue(t *testing.T, want, got map[string]any) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
