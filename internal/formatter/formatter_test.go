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
			Stylish{},
			nil,
		},
		{
			"detect plain formatter",
			"plain",
			Plain{},
			nil,
		},
		{
			"detect json formatter",
			"json",
			JSON{},
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
			analyzer.NewContainerNode("root", map[string]analyzer.DiffNode{
				"timeout": analyzer.NewChangedNode("timeout", 50.0, 100.0),
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
			analyzer.NewContainerNode("root", map[string]analyzer.DiffNode{
				"timeout": analyzer.NewChangedNode("timeout", 50.0, 100.0),
			}),
			"Property 'timeout' was updated. From 50 to 100",
			nil,
		},
		{
			"detect correct formatting stylish",
			"json",
			analyzer.NewContainerNode("root", map[string]analyzer.DiffNode{
				"timeout": analyzer.NewChangedNode("timeout", 50.0, 100.0),
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
			analyzer.NewContainerNode("root", map[string]analyzer.DiffNode{}),
			"",
			errCreateFormatter,
		},
		{
			"detect error errFormatting",
			"json",
			analyzer.NewContainerNode("root", map[string]analyzer.DiffNode{
				"func": analyzer.NewUnchangedNode("func", any(func() {})),
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
			got, err := Format(td.format, td.diff)
			assert.Equal(t, td.want, got)
			require.ErrorIs(t, err, td.wantErr)
		})
	}
}
