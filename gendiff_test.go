package code

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenDiff(t *testing.T) {
	tests := []struct {
		name      string
		filepath1 string
		filepath2 string
		format    string
		want      string
		wantErr   error
	}{
		{
			"detect equal JSON files in stylish format",
			"testdata/simple.json",
			"testdata/simple.json",
			"stylish",
			`{
    follow: false
    host: hexlet.io
    proxy: 123.234.53.22
    secret: null
    timeout: 50
}`,
			nil,
		},
		{
			"detect equal JSON files with default empty format",
			"testdata/simple.json",
			"testdata/simple.json",
			"",
			`{
    follow: false
    host: hexlet.io
    proxy: 123.234.53.22
    secret: null
    timeout: 50
}`,
			nil,
		},
		{
			"detect equal JSON files in plain format",
			"testdata/simple.json",
			"testdata/simple.json",
			"plain",
			``,
			nil,
		},
		{
			"detect equal JSON and YAML files",
			"testdata/simple.json",
			"testdata/simple.yaml",
			"stylish",
			`{
    follow: false
    host: hexlet.io
    proxy: 123.234.53.22
    secret: null
    timeout: 50
}`,
			nil,
		},
		{
			"detect differences in stylish format",
			"testdata/complex3.json",
			"testdata/complex4.json",
			"stylish",
			`{
    additional: {
        one: 1
      - three: [
            1, 2, 4
            2
            hello from, GO
        ]
      + three: [
            1, 2
            2
            hello from, GO
        ]
        two: 2
    }
    follow: false
    host: hexlet.io
    proxy: 123.234.53.22
  - secret: null
  + secret: null
    timeout: 50
}`,
			nil,
		},
		{
			"detect differences in plain format",
			"testdata/complex3.json",
			"testdata/complex4.json",
			"plain",
			"Property 'additional.three' was updated. From [complex value] to [complex value]\nProperty 'secret' was updated. From null to 'null'",
			nil,
		},
		{
			"detect error with missing first file",
			"testdata/missing.json",
			"testdata/simple.json",
			"stylish",
			``,
			errFailedProcessFirstFile,
		},
		{
			"detect error with missing second file",
			"testdata/simple.json",
			"testdata/missing.json",
			"stylish",
			``,
			errFailedProcessSecondFile,
		},
		{
			"detect error with unsupported format",
			"testdata/simple.json",
			"testdata/simple.json",
			"unsupported",
			``,
			errFailedFormat,
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			got, err := GenDiff(td.filepath1, td.filepath2, td.format)

			assert.Equal(t, td.want, got)
			assert.ErrorIs(t, err, td.wantErr)
		})
	}
}

func TestGenDiff_JSONFormat_SemanticIntegration(t *testing.T) {
	t.Run("complex diff semantic check", func(t *testing.T) {
		result, err := GenDiff("testdata/complex1.json", "testdata/complex2.json", "json")
		require.NoError(t, err)

		var diff map[string]any

		err = json.Unmarshal([]byte(result), &diff)
		require.NoError(t, err)

		require.Equal(t, "root", diff["key"])
		require.Equal(t, "Container", diff["status"])

		children, ok := diff["child"].(map[string]any)
		require.True(t, ok)

		common, ok := children["common"].(map[string]any)
		require.True(t, ok)
		require.Equal(t, "Container", common["status"])

		commonChildren, ok := common["child"].(map[string]any)
		require.True(t, ok)

		setting1, ok := commonChildren["setting1"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "Deleted", setting1["status"])
		assert.Equal(t, "Value 1", setting1["old"])
		assert.Nil(t, setting1["new"])
	})

	t.Run("null vs string null preservation across full pipeline", func(t *testing.T) {
		result, err := GenDiff("testdata/complex3.json", "testdata/complex4.json", "json")
		require.NoError(t, err)

		var diff map[string]any

		err = json.Unmarshal([]byte(result), &diff)
		require.NoError(t, err)

		children, ok := diff["child"].(map[string]any)
		require.True(t, ok)

		secret, ok := children["secret"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "Changed", secret["status"])
		assert.Nil(t, secret["old"])
		_, hasOld := secret["old"]
		assert.True(t, hasOld, "old key must be present with null value")
		assert.Equal(t, "null", secret["new"])
	})
}
