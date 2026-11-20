package formatter

import "code/internal/analyzer"

var (
	simpleDiff = analyzer.NewContainerNode("root", map[string]analyzer.DiffNode{
		"follow":  analyzer.NewChangedNode("follow", false, true),
		"host":    analyzer.NewUnchangedNode("host", "hexlet.io"),
		"proxy":   analyzer.NewChangedNode("proxy", "123.234.53.22", "123.235.53.22"),
		"secret":  analyzer.NewUnchangedNode("secret", "null"),
		"timeout": analyzer.NewChangedNode("timeout", int64(50), int64(100)),
	})

	diffWithArrays = analyzer.NewContainerNode("root", map[string]analyzer.DiffNode{
		"languages": analyzer.NewUnchangedNode("languages", []any{"GO", "RUST", "PYTHON"}),
		"numbers": analyzer.NewChangedNode(
			"numbers",
			[]any{int64(1), int64(6), int64(9)},
			[]any{int64(1), int64(6), int64(9), 22.6},
		),
		"others": analyzer.NewChangedNode("others", []any{"hello", "world", "null"}, []any{"hello", "world"}),
	})

	diffWithNil = analyzer.NewContainerNode("root", map[string]analyzer.DiffNode{
		"toNullValue":   analyzer.NewChangedNode("toNullValue", "null", nil),
		"toStringValue": analyzer.NewChangedNode("toStringValue", nil, "null"),
	})

	complexDiff = analyzer.NewContainerNode("root", map[string]analyzer.DiffNode{
		"common": analyzer.NewContainerNode("common", map[string]analyzer.DiffNode{
			"setting1": analyzer.NewDeletedNode("setting1", "Value 1"),
			"setting2": analyzer.NewUnchangedNode("setting2", int64(200)),
			"setting3": analyzer.NewUnchangedNode("setting3", true),
			"setting6": analyzer.NewContainerNode("setting6", map[string]analyzer.DiffNode{
				"key": analyzer.NewUnchangedNode("key", "value"),
				"doge": analyzer.NewContainerNode("doge", map[string]analyzer.DiffNode{
					"wow": analyzer.NewUnchangedNode("wow", ""),
				}),
			}),
		}),
		"group1": analyzer.NewContainerNode("group1", map[string]analyzer.DiffNode{
			"baz": analyzer.NewUnchangedNode("baz", "bas"),
			"foo": analyzer.NewUnchangedNode("foo", "bar"),
			"nest": analyzer.NewAddedNode("nest", map[string]any{
				"key": "value",
			}),
			"secret":  analyzer.NewAddedNode("secret", "null"),
			"secret2": analyzer.NewAddedNode("secret2", nil),
		}),
		"group2": analyzer.NewAddedNode("group2", map[string]any{
			"abc": int64(12345),
			"deep": map[string]any{
				"id": int64(45),
			},
		}),
	})

	diffWithValueContainsPunctuationMarks = analyzer.NewContainerNode("root", map[string]analyzer.DiffNode{
		"docsIndex": analyzer.NewChangedNode("docsIndex", "124/365-44, 124/365-45", "124/365-45, 124/365-46"),
	})

	diffWithMultilineString = analyzer.NewContainerNode("root", map[string]analyzer.DiffNode{
		"description": analyzer.NewChangedNode("description", "first line\nsecond line", "single line"),
	})

	diffWithNewlineAndLiteralBackslash = analyzer.NewContainerNode("root", map[string]analyzer.DiffNode{
		"description": analyzer.NewChangedNode("description", "a\nb", "a\\nb"),
	})

	emptyDiff = analyzer.NewContainerNode("root", map[string]analyzer.DiffNode{})
)
