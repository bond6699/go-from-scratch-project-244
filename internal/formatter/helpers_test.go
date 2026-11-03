package formatter

import "code/internal/analyzer"

var (
	simpleDiff = *ContainerNode("root", map[string]*analyzer.DiffNode{
		"follow":  ChangedNode("follow", false, true),
		"host":    UnchangedNode("host", "hexlet.io"),
		"proxy":   ChangedNode("proxy", "123.234.53.22", "123.235.53.22"),
		"secret":  UnchangedNode("secret", "null"),
		"timeout": ChangedNode("timeout", int64(50), int64(100)),
	})

	diffWithArrays = *ContainerNode("root", map[string]*analyzer.DiffNode{
		"languages": UnchangedNode("languages", []any{"GO", "RUST", "PYTHON"}),
		"numbers":   ChangedNode("numbers", []any{int64(1), int64(6), int64(9)}, []any{int64(1), int64(6), int64(9), 22.6}),
		"others":    ChangedNode("others", []any{"hello", "world", "null"}, []any{"hello", "world"}),
	})

	diffWithNil = *ContainerNode("root", map[string]*analyzer.DiffNode{
		"toNullValue":   ChangedNode("toNullValue", "null", nil),
		"toStringValue": ChangedNode("toStringValue", nil, "null"),
	})

	complexDiff = *ContainerNode("root", map[string]*analyzer.DiffNode{
		"common": ContainerNode("common", map[string]*analyzer.DiffNode{
			"setting1": DeletedNode("setting1", "Value 1"),
			"setting2": UnchangedNode("setting2", int64(200)),
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
			"abc": int64(12345),
			"deep": map[string]any{
				"id": int64(45),
			},
		}),
	})

	diffWithValueContainsPunctuationMarks = *ContainerNode("root", map[string]*analyzer.DiffNode{
		"docsIndex": ChangedNode("docsIndex", "124/365-44, 124/365-45", "124/365-45, 124/365-46"),
	})

	emptyDiff = *ContainerNode("root", map[string]*analyzer.DiffNode{})
)
