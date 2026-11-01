package formatter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"code/internal/analyzer"
)

func TestGetFormatter(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		want    Formatter
		wantErr error
	}{
		{
			"detect stylish formatter",
			"stylish",
			FormatterStylish{},
			nil,
		},
		{
			"detect plain formatter",
			"plain",
			FormatterPlain{},
			nil,
		},
		{
			"detect json formatter",
			"json",
			FormatterJSON{},
			nil,
		},
		{
			"detect error with invalid format",
			"invalidFormat",
			nil,
			errUnsupportedFormat,
		},
		{
			"detect error with empty format",
			"",
			nil,
			errUnsupportedFormat,
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			got, err := getFormatter(td.format)
			require.Equal(t, td.want, got)
			require.ErrorIs(t, err, td.wantErr)
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		diff    analyzer.DiffNode
		want    string
		wantErr error
	}{
		{
			"detect correct formatting stylish",
			"stylish",
			*ContainerNode("root", map[string]*analyzer.DiffNode{
				"timeout": ChangedNode("timeout", 50.0, 100.0),
			}),
			`{
  - timeout: 50
  + timeout: 100
}`,
			nil,
		},
		{
			"detect correct formatting stylish",
			"plain",
			*ContainerNode("root", map[string]*analyzer.DiffNode{
				"timeout": ChangedNode("timeout", 50.0, 100.0),
			}),
			"Property 'timeout' was updated. From 50 to 100",
			nil,
		},
		{
			"detect correct formatting stylish",
			"json",
			*ContainerNode("root", map[string]*analyzer.DiffNode{
				"timeout": ChangedNode("timeout", 50.0, 100.0),
			}),
			`{
   "child": {
      "timeout": {
         "old": 50,
         "new": 100,
         "key": "timeout",
         "status": "Changed"
      }
   },
   "key": "root",
   "status": "Container"
}`,
			nil,
		},
		{
			"detect error errCreateFormatter",
			"unsupportedFormat",
			*ContainerNode("root", map[string]*analyzer.DiffNode{}),
			"",
			errCreateFormatter,
		},
		{
			"detect error errFormatting",
			"json",
			*ContainerNode("root", map[string]*analyzer.DiffNode{
				"func": UnchangedNode("func", any(func() {})),
			}),
			"",
			errFormatting,
		},
	}

	for _, td := range tests {
		name := td.format
		if name == "" {
			name = "default"
		}

		t.Run(name, func(t *testing.T) {
			got, err := Format(td.format, &td.diff)
			assert.Equal(t, td.want, got)
			require.ErrorIs(t, err, td.wantErr)
		})
	}
}

func ContainerNode(key string, value map[string]*analyzer.DiffNode) *analyzer.DiffNode {
	node := analyzer.NewDiffNode(key, analyzer.Container)
	node.Children = value

	return node
}

func UnchangedNode(key string, value any) *analyzer.DiffNode {
	node := analyzer.NewDiffNode(key, analyzer.Unchanged)
	node.SetOldValue(value)

	return node
}

func ChangedNode(key string, oldValue, newValue any) *analyzer.DiffNode {
	node := analyzer.NewDiffNode(key, analyzer.Changed)
	node.SetOldValue(oldValue)
	node.SetNewValue(newValue)

	return node
}

func DeletedNode(key string, value any) *analyzer.DiffNode {
	node := analyzer.NewDiffNode(key, analyzer.Deleted)
	node.SetOldValue(value)

	return node
}

func AddedNode(key string, value any) *analyzer.DiffNode {
	node := analyzer.NewDiffNode(key, analyzer.Added)
	node.SetNewValue(value)

	return node
}
