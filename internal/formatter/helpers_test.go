package formatter

import "code/internal/analyzer"

var (
	simpleDiff = *ContainerNode("root", map[string]*analyzer.DiffNode{
		"follow":  ChangedNode("follow", false, true),
		"host":    UnchangedNode("host", "hexlet.io"),
		"proxy":   ChangedNode("proxy", "123.234.53.22", "123.235.53.22"),
		"secret":  UnchangedNode("secret", "null"),
		"timeout": ChangedNode("timeout", 50.0, 100.0),
	})

	diffWithArrays = *ContainerNode("root", map[string]*analyzer.DiffNode{
		"languages": UnchangedNode("languages", []any{"GO", "RUST", "PYTHON"}),
		"numbers":   ChangedNode("numbers", []any{1.0, 6.0, 9.0}, []any{1.0, 6.0, 9.0, 22.6}),
		"others":    ChangedNode("others", []any{"hello", "world", "null"}, []any{"hello", "world"}),
	})

	diffWithNil = *ContainerNode("root", map[string]*analyzer.DiffNode{
		"toNullValue":   ChangedNode("toNullValue", "null", nil),
		"toStringValue": ChangedNode("toStringValue", nil, "null"),
	})

	complexDiff = *ContainerNode("root", map[string]*analyzer.DiffNode{
		"common": ContainerNode("common", map[string]*analyzer.DiffNode{
			"setting1": DeletedNode("setting1", "Value 1"),
			"setting2": UnchangedNode("setting2", 200.0),
			"setting3": UnchangedNode("setting3", true),
			"setting6": ContainerNode("setting6", map[string]*analyzer.DiffNode{
				"key": UnchangedNode("key", "value"),
				"doge": ContainerNode("doge", map[string]*analyzer.DiffNode{
					"wow": UnchangedNode("wow", ""),
				}),
			}),
		}),
		"group1": ContainerNode("group1", map[string]*analyzer.DiffNode{
			"baz": UnchangedNode("baz", "bas"),
			"foo": UnchangedNode("foo", "bar"),
			"nest": AddedNode("nest", map[string]any{
				"key": "value",
			}),
			"secret":  AddedNode("secret", "null"),
			"secret2": AddedNode("secret2", nil),
		}),
		"group2": AddedNode("group2", map[string]any{
			"abc": 12345.0,
			"deep": map[string]any{
				"id": 45.0,
			},
		}),
	})

	emptyDiff = *ContainerNode("root", map[string]*analyzer.DiffNode{})
)
