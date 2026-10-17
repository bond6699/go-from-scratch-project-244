package parser

import (
	"os"
	"reflect"
	"testing"
)

func assertError(t *testing.T, wantErr bool, got error) {
	t.Helper()
	if wantErr && got == nil {
		t.Fatal("expected an error, got nil")
	}

	if !wantErr && got != nil {
		t.Fatalf("unexpected error: %v", got)
	}
}

func assertParserValue(t *testing.T, want map[string]any, got map[string]any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestJsonParser(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    map[string]any
		wantErr bool
	}{
		{
			"valid JSON value",
			`{"common": {"setting1": "Value 1", "setting2": 200, "setting3": true}}`,
			map[string]any{
				"common": map[string]any{
					"setting1": "Value 1",
					"setting2": 200.0,
					"setting3": true,
				},
			},
			false,
		},
		{
			"invalid JSON value",
			`"common": {"setting1": "Value 1", "sett`,
			map[string]any{},
			true,
		},
		{
			"empty", "{}", map[string]any{}, false,
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			got, err := JSONParser{}.Parse([]byte(td.input))
			assertError(t, td.wantErr, err)

			if !td.wantErr {
				assertParserValue(t, td.want, got)
			}
		})
	}
}

func TestYamlParser(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    map[string]any
		wantErr bool
	}{
		{
			"valid YAML value",
			`
common:
  setting1: "Value 1"
  setting2: 200
  setting3: true`,
			map[string]any{
				"common": map[string]any{
					"setting1": "Value 1",
					"setting2": 200,
					"setting3": true,
				},
			},
			false,
		},
		{
			"invalid YAML value",
			`
common
  setting1 - Value 1"
  setting2: 200
  setting3: true`,
			map[string]any{},
			true,
		},
		{
			"empty", "", map[string]any{}, false,
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			got, err := YamlParser{}.Parse([]byte(td.input))
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
				"name": "Jany",
				"age":  23,
				"hobby": []any{
					map[string]any{
						"name":     "Go Code",
						"expYears": 2,
					},
					map[string]any{
						"name":     "Read Book",
						"expYears": 10,
					},
				},
				"Secret": nil,
			},
			map[string]any{
				"name": "Jany",
				"age":  float64(23.0),
				"hobby": []any{
					map[string]any{
						"name":     "Go Code",
						"expYears": float64(2.0),
					},
					map[string]any{
						"name":     "Read Book",
						"expYears": float64(10.0),
					},
				},
				"Secret": "null",
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
		filename string
		content  string
		want     map[string]any
		wantErr  bool
	}{
		{
			"valid JSON file",
			"config-*.json",
			`{"common": {"setting1": "Value 1", "setting2": 200, "setting3": true, "setting4": null}}`,
			map[string]any{
				"common": map[string]any{
					"setting1": "Value 1",
					"setting2": 200.0,
					"setting3": true,
					"setting4": "null",
				},
			},
			false,
		},
		{
			"invalid JSON file",
			"config-*.json",
			`{"common": {"setting1": "Value 1", "setting2": 200, "se`,
			map[string]any{},
			true,
		},
		{
			"empty JSON file",
			"config-*.json",
			"{}",
			map[string]any{},
			false,
		},
		{
			"valid YAML file",
			"config-*.yml",
			`
common:
  setting1: "Value 1"
  setting2: 200
  setting3: true
  setting4: null`,
			map[string]any{
				"common": map[string]any{
					"setting1": "Value 1",
					"setting2": 200.0,
					"setting3": true,
					"setting4": "null",
				},
			},
			false,
		},
		{
			"invalid YAML file",
			"config-*.yaml",
			`
common
  setting1 - "Value 1"
  setting2: 200
  settin`,
			map[string]any{},
			true,
		},
		{
			"empty YAML file",
			"config-*.yml",
			"",
			map[string]any{},
			false,
		},
		{
			"invalid file format",
			"data.txt",
			"Hello, World!",
			map[string]any{},
			true,
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			file, err := os.CreateTemp(t.TempDir(), td.filename)
			if err != nil {
				t.Fatal(err)
			}

			_, err = file.Write([]byte(td.content))
			if err != nil {
				t.Fatal(err)
			}

			got, err := Parse(file.Name())
			assertError(t, td.wantErr, err)
			if !td.wantErr {
				assertParserValue(t, td.want, Normalize(got))
			}
		})
	}
}
