package formatter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"code/internal/analyzer"
)

func TestPlainFormat(t *testing.T) {
	tests := []struct {
		name string
		diff *analyzer.ContainerNode
		want string
	}{
		{
			"detect correct PLAIN result with simple values",
			simpleDiff,
			`Property 'follow' was updated. From false to true
Property 'proxy' was updated. From '123.234.53.22' to '123.235.53.22'
Property 'timeout' was updated. From 50 to 100`,
		},
		{
			"detect correct PLAIN result with array values",
			diffWithArrays,
			`Property 'numbers' was updated. From [complex value] to [complex value]
Property 'others' was updated. From [complex value] to [complex value]`,
		},
		{
			"detect saved nil value in PLAIN result",
			diffWithNil,
			`Property 'toNullValue' was updated. From 'null' to null
Property 'toStringValue' was updated. From null to 'null'`,
		},
		{
			"detect correct PLAIN result with complex values",
			complexDiff,
			`Property 'common.setting1' was removed
Property 'group1.nest' was added with value: [complex value]
Property 'group1.secret' was added with value: 'null'
Property 'group1.secret2' was added with value: null
Property 'group2' was added with value: [complex value]`,
		},
		{
			"detect correct PLAIN result with value contains punctuation marks",
			diffWithValueContainsPunctuationMarks,
			"Property 'docsIndex' was updated. From '124/365-44, 124/365-45' to '124/365-45, 124/365-46'",
		},
		{
			"detect correct PLAIN result with multiline string",
			diffWithMultilineString,
			`Property 'description' was updated. From 'first line\nsecond line' to 'single line'`,
		},
		{
			"detect distinction between newline and literal backslash in PLAIN result",
			diffWithNewlineAndLiteralBackslash,
			`Property 'description' was updated. From 'a\nb' to 'a\\nb'`,
		},
		{
			"detect correct PLAIN result with empty value",
			emptyDiff,
			"",
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			got, err := Plain{}.Format(td.diff)
			assert.Equal(t, td.want, got)
			require.NoError(t, err)
		})
	}
}
