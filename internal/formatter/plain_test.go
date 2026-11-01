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
		diff analyzer.DiffNode
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
			`Property 'numbers' was updated. From [1 6 9] to [1 6 9 22.6]
Property 'others' was updated. From [hello world null] to [hello world]`,
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
			"detect correct PLAIN result with empty value",
			emptyDiff,
			"",
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			got, err := FormatterPlain{}.Format(&td.diff)
			assert.Equal(t, td.want, got)
			require.NoError(t, err)
		})
	}
}
