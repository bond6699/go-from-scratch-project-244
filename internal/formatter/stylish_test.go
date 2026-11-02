package formatter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"code/internal/analyzer"
)

func TestStylishFormat(t *testing.T) {
	tests := []struct {
		name  string
		input analyzer.DiffNode
		want  string
	}{
		{
			"detect correct STYLISH result with simple values",
			simpleDiff,
			`{
  - follow: false
  + follow: true
    host: hexlet.io
  - proxy: 123.234.53.22
  + proxy: 123.235.53.22
    secret: null
  - timeout: 50
  + timeout: 100
}`,
		},
		{
			"detect correct STYLISH result with array values",
			diffWithArrays,
			`{
    languages: [
        GO
        RUST
        PYTHON
    ]
  - numbers: [
        1
        6
        9
    ]
  + numbers: [
        1
        6
        9
        22.6
    ]
  - others: [
        hello
        world
        null
    ]
  + others: [
        hello
        world
    ]
}`,
		},
		{
			"detect saved nil value in STYLISH result",
			diffWithNil,
			`{
  - toNullValue: null
  + toNullValue: null
  - toStringValue: null
  + toStringValue: null
}`,
		},
		{
			"detect correct STYLISH result with complex values",
			complexDiff,
			`{
    common: {
      - setting1: Value 1
        setting2: 200
        setting3: true
        setting6: {
            doge: {
                wow: 
            }
            key: value
        }
    }
    group1: {
        baz: bas
        foo: bar
      + nest: {
            key: value
        }
      + secret: null
      + secret2: null
    }
  + group2: {
        abc: 12345
        deep: {
            id: 45
        }
    }
}`,
		},
		{
			"detect correct STYLISH result with value contains punctuation marks",
			diffWithValueContainsPunctuationMarks,
			`{
  - docsIndex: 124/365-44, 124/365-45
  + docsIndex: 124/365-45, 124/365-46
}`,
		},
		{
			"detect correct STYLISH result with empty value",
			emptyDiff,
			`{
}`,
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			got, err := FormatterStylish{}.Format(&td.input)
			assert.Equal(t, td.want, got)
			require.NoError(t, err)
		})
	}
}
