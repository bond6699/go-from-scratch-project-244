package analyzer

import (
	"reflect"
	"testing"
)

func TestAnalyze(t *testing.T) {
	tests := []struct {
		name string
		map1 map[string]any
		map2 map[string]any
		want DiffNode
	}{
		{
			"simple maps",
			map[string]any{
				"host":    "hexlet.io",
				"timeout": 50.0,
				"proxy":   "123.234.53.22",
				"follow":  false,
				"secret":  "null",
			},
			map[string]any{
				"host":    "hexlet.io",
				"timeout": 100.0,
				"proxy":   "123.235.53.22",
				"follow":  true,
				"secret":  "null",
			},
			DiffNode{
				Key:    "root",
				Status: Container,
				Children: map[string]*DiffNode{
					"follow": {
						Key:      "follow",
						Status:   Changed,
						OldValue: false,
						NewValue: true,
					},
					"host": {
						Key:      "host",
						Status:   Unchanged,
						OldValue: "hexlet.io",
						NewValue: nil,
					},
					"proxy": {
						Key:      "proxy",
						Status:   Changed,
						OldValue: "123.234.53.22",
						NewValue: "123.235.53.22",
					},
					"secret": {
						Key:      "secret",
						Status:   Unchanged,
						OldValue: "null",
						NewValue: nil,
					},
					"timeout": {
						Key:      "timeout",
						Status:   Changed,
						OldValue: 50.0,
						NewValue: 100.0,
					},
				},
			},
		},
		{
			"simple equal maps",
			map[string]any{
				"host":    "hexlet.io",
				"timeout": 50.0,
				"proxy":   "123.234.53.22",
				"follow":  false,
				"secret":  "null",
			},
			map[string]any{
				"host":    "hexlet.io",
				"timeout": 50.0,
				"proxy":   "123.234.53.22",
				"follow":  false,
				"secret":  "null",
			},
			DiffNode{
				Key:    "root",
				Status: Container,
				Children: map[string]*DiffNode{
					"follow": {
						Key:      "follow",
						Status:   Unchanged,
						OldValue: false,
						NewValue: nil,
					},
					"host": {
						Key:      "host",
						Status:   Unchanged,
						OldValue: "hexlet.io",
						NewValue: nil,
					},
					"proxy": {
						Key:      "proxy",
						Status:   Unchanged,
						OldValue: "123.234.53.22",
						NewValue: nil,
					},
					"secret": {
						Key:      "secret",
						Status:   Unchanged,
						OldValue: "null",
						NewValue: nil,
					},
					"timeout": {
						Key:      "timeout",
						Status:   Unchanged,
						OldValue: 50.0,
						NewValue: nil,
					},
				},
			},
		},
		{
			"complex maps",
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
				},
			},
			map[string]any{
				"common": map[string]any{
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
			DiffNode{
				Key:    "root",
				Status: Container,
				Children: map[string]*DiffNode{
					"common": {
						Key:    "common",
						Status: Container,
						Children: map[string]*DiffNode{
							"setting1": {
								Key:      "setting1",
								Status:   Deleted,
								OldValue: "Value 1",
								NewValue: nil,
							},
							"setting2": {
								Key:      "setting2",
								Status:   Unchanged,
								OldValue: 200.0,
								NewValue: nil,
							},
							"setting3": {
								Key:      "setting3",
								Status:   Unchanged,
								OldValue: true,
								NewValue: nil,
							},
							"setting6": {
								Key:    "setting6",
								Status: Container,
								Children: map[string]*DiffNode{
									"doge": {
										Key:    "doge",
										Status: Container,
										Children: map[string]*DiffNode{
											"wow": {
												Key:      "wow",
												Status:   Unchanged,
												OldValue: "",
												NewValue: nil,
											},
										},
									},
									"key": {
										Key:      "key",
										Status:   Unchanged,
										OldValue: "value",
										NewValue: nil,
									},
								},
							},
						},
					},
					"group1": {
						Key:    "group1",
						Status: Container,
						Children: map[string]*DiffNode{
							"baz": {
								Key:      "baz",
								Status:   Unchanged,
								OldValue: "bas",
								NewValue: nil,
							},
							"foo": {
								Key:      "foo",
								Status:   Unchanged,
								OldValue: "bar",
								NewValue: nil,
							},
							"nest": {
								Key:    "nest",
								Status: Added,
								NewValue: map[string]any{
									"key": "value",
								},
							},
							"secret": {
								Key:      "secret",
								Status:   Added,
								NewValue: "null",
							},
						},
					},
					"group2": {
						Key:    "group2",
						Status: Added,
						NewValue: map[string]any{
							"abc": 12345.0,
							"deep": map[string]any{
								"id": 45.0,
							},
						},
					},
				},
			},
		},
		{
			"maps wint array",
			map[string]any{
				"numbers":   []any{1.0, 6.0, 9.0},
				"others":    []any{"hello", "world", "null"},
				"languages": []any{"GO", "RUST", "PYTHON"},
			},
			map[string]any{
				"numbers":   []any{1.0, 6.0, 9.0, 22.6},
				"others":    []any{"hello", "world"},
				"languages": []any{"GO", "RUST", "PYTHON"},
			},
			DiffNode{
				Key:    "root",
				Status: Container,
				Children: map[string]*DiffNode{
					"languages": {
						Key:      "languages",
						Status:   Unchanged,
						OldValue: []any{"GO", "RUST", "PYTHON"},
						NewValue: nil,
					},
					"numbers": {
						Key:      "numbers",
						Status:   Changed,
						OldValue: []any{1.0, 6.0, 9.0},
						NewValue: []any{1.0, 6.0, 9.0, 22.6},
					},
					"others": {
						Key:      "others",
						Status:   Changed,
						OldValue: []any{"hello", "world", "null"},
						NewValue: []any{"hello", "world"},
					},
				},
			},
		},
		{
			"empty maps",
			map[string]any{},
			map[string]any{},
			DiffNode{
				Key:      "root",
				Status:   Container,
				Children: map[string]*DiffNode{},
			},
		},
		{
			"all keys added",
			map[string]any{},
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
				},
			},
			DiffNode{
				Key:    "root",
				Status: Container,
				Children: map[string]*DiffNode{
					"common": {
						Key:    "common",
						Status: Added,
						NewValue: map[string]any{
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
					},
					"group1": {
						Key:    "group1",
						Status: Added,
						NewValue: map[string]any{
							"baz": "bas",
							"foo": "bar",
						},
					},
				},
			},
		},
		{
			"all keys deleted",
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
				},
			},
			map[string]any{},
			DiffNode{
				Key:    "root",
				Status: Container,
				Children: map[string]*DiffNode{
					"common": {
						Key:    "common",
						Status: Deleted,
						OldValue: map[string]any{
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
					},
					"group1": {
						Key:    "group1",
						Status: Deleted,
						OldValue: map[string]any{
							"baz": "bas",
							"foo": "bar",
						},
					},
				},
			},
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			got := *Analyze(td.map1, td.map2)
			assertAnalyzerValue(t, td.want, got)
		})
	}
}

func assertAnalyzerValue(t *testing.T, want DiffNode, got DiffNode) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got:\n%#v\nwant:\n%#v", got, want)
	}
}
