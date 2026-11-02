package code

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
