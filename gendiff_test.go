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
			"internal/parser/testdata/json/simple.json",
			"internal/parser/testdata/json/simple.json",
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
			"internal/parser/testdata/json/simple.json",
			"internal/parser/testdata/json/simple.json",
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
			"internal/parser/testdata/json/simple.json",
			"internal/parser/testdata/json/simple.json",
			"plain",
			``,
			nil,
		},
		{
			"detect equal JSON and YAML files",
			"internal/parser/testdata/json/simple.json",
			"internal/parser/testdata/yaml/simple.yaml",
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
			"detect error with missing first file",
			"testdata/missing.json",
			"internal/parser/testdata/json/simple.json",
			"stylish",
			``,
			errFailedProcessFirstFile,
		},
		{
			"detect error with missing second file",
			"internal/parser/testdata/json/simple.json",
			"testdata/missing.json",
			"stylish",
			``,
			errFailedProcessSecondFile,
		},
		{
			"detect error with unsupported format",
			"internal/parser/testdata/json/simple.json",
			"internal/parser/testdata/json/simple.json",
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
}
