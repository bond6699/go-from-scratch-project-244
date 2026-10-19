package formatter

import (
	"testing"

	"code/internal/analyzer"
)

func TestPlainFormat(t *testing.T) {
	tests := []struct {
		name  string
		input analyzer.DiffNode
		want  string
	}{
		{
			"simple values",
			analyzer.DiffNode{
				Key:    "root",
				Status: analyzer.Container,
				Children: map[string]*analyzer.DiffNode{
					"follow": {
						Key:      "follow",
						Status:   analyzer.Changed,
						OldValue: false,
						NewValue: true,
					},
					"host": {
						Key:      "host",
						Status:   analyzer.Unchanged,
						OldValue: "hexlet.io",
					},
					"proxy": {
						Key:      "proxy",
						Status:   analyzer.Changed,
						OldValue: "123.234.53.22",
						NewValue: "123.235.53.22",
					},
					"secret": {
						Key:      "secret",
						Status:   analyzer.Unchanged,
						OldValue: "null",
					},
					"timeout": {
						Key:      "timeout",
						Status:   analyzer.Changed,
						OldValue: 50.0,
						NewValue: 100.0,
					},
				},
			},
			`Property 'follow' was updated. From false to true
Property 'proxy' was updated. From '123.234.53.22' to '123.235.53.22'
Property 'timeout' was updated. From 50 to 100
`,
		},
		{
			"value with arrays",
			analyzer.DiffNode{
				Key:    "root",
				Status: analyzer.Container,
				Children: map[string]*analyzer.DiffNode{
					"languages": {
						Key:      "languages",
						Status:   analyzer.Unchanged,
						OldValue: []any{"GO", "RUST", "PYTHON"},
					},
					"numbers": {
						Key:      "numbers",
						Status:   analyzer.Changed,
						OldValue: []any{1.0, 6.0, 9.0},
						NewValue: []any{1.0, 6.0, 9.0, 22.6},
					},
					"others": {
						Key:      "others",
						Status:   analyzer.Changed,
						OldValue: []any{"hello", "world", "null"},
						NewValue: []any{"hello", "world"},
					},
				},
			},
			`Property 'numbers' was updated. From [1 6 9] to [1 6 9 22.6]
Property 'others' was updated. From [hello world null] to [hello world]
`,
		},
		{
			"complex value",
			analyzer.DiffNode{
				Key:    "root",
				Status: analyzer.Container,
				Children: map[string]*analyzer.DiffNode{
					"common": {
						Key:    "common",
						Status: analyzer.Container,
						Children: map[string]*analyzer.DiffNode{
							"setting1": {
								Key:      "setting1",
								Status:   analyzer.Deleted,
								OldValue: "Value 1",
							},
							"setting2": {
								Key:      "setting2",
								Status:   analyzer.Unchanged,
								OldValue: 200.0,
							},
							"setting3": {
								Key:      "setting3",
								Status:   analyzer.Unchanged,
								OldValue: true,
							},
							"setting6": {
								Key:    "setting6",
								Status: analyzer.Container,
								Children: map[string]*analyzer.DiffNode{
									"doge": {
										Key:    "doge",
										Status: analyzer.Container,
										Children: map[string]*analyzer.DiffNode{
											"wow": {
												Key:      "wow",
												Status:   analyzer.Unchanged,
												OldValue: "",
											},
										},
									},
									"key": {
										Key:      "key",
										Status:   analyzer.Unchanged,
										OldValue: "value",
									},
								},
							},
						},
					},
					"group1": {
						Key:    "group1",
						Status: analyzer.Container,
						Children: map[string]*analyzer.DiffNode{
							"baz": {
								Key:      "baz",
								Status:   analyzer.Unchanged,
								OldValue: "bas",
							},
							"foo": {
								Key:      "foo",
								Status:   analyzer.Unchanged,
								OldValue: "bar",
							},
							"nest": {
								Key:    "nest",
								Status: analyzer.Added,
								NewValue: map[string]any{
									"key": "value",
								},
							},
							"secret": {
								Key:      "secret",
								Status:   analyzer.Added,
								NewValue: "null",
							},
						},
					},
					"group2": {
						Key:    "group2",
						Status: analyzer.Added,
						NewValue: map[string]any{
							"abc": 12345.0,
							"deep": map[string]any{
								"id": 45.0,
							},
						},
					},
				},
			},
			`Property 'common.setting1' was removed
Property 'group1.nest' was added with value: [complex value]
Property 'group1.secret' was added with value: null
Property 'group2' was added with value: [complex value]
`,
		},
		{
			"empty value",
			analyzer.DiffNode{
				Key:      "root",
				Status:   analyzer.Container,
				Children: map[string]*analyzer.DiffNode{},
			},
			"",
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			got, err := FormatterPlain{}.Format(&td.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != td.want {
				t.Errorf("got:\n%s\nwant:\n%s", got, td.want)
			}
		})
	}
}
