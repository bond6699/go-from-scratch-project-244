package formatter

import "code/internal/analyzer"

func simpleDiffInput() analyzer.DiffNode {
	return analyzer.DiffNode{
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
	}
}

func withArrayDiffInput() analyzer.DiffNode {
	return analyzer.DiffNode{
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
	}
}

func complexDiffInput() analyzer.DiffNode {
	return analyzer.DiffNode{
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
	}
}
