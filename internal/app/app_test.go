package app

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr error
	}{
		{
			"run with DEFAULT format",
			[]string{
				"gendiff",
				"../parser/testdata/json/simple.json",
				"../parser/testdata/json/simple.json",
			},
			`{
    follow: false
    host: hexlet.io
    proxy: 123.234.53.22
    secret: null
    timeout: 50
}` + "\n",
			nil,
		},
		{
			"run with STYLISH format",
			[]string{
				"gendiff",
				"../parser/testdata/json/simple.json",
				"../parser/testdata/json/simple.json",
			},
			`{
    follow: false
    host: hexlet.io
    proxy: 123.234.53.22
    secret: null
    timeout: 50
}` + "\n",
			nil,
		},
		{
			"run with PLAIN format",
			[]string{
				"gendiff",
				"--format",
				"plain",
				"../parser/testdata/json/simple.json",
				"../parser/testdata/json/simple.json",
			},
			``,
			nil,
		},
		{
			"run with JSON format",
			[]string{
				"gendiff",
				"--format",
				"json",
				"../parser/testdata/json/simple.json",
				"../parser/testdata/json/simple.json",
			},
			`{
   "child": {
      "follow": {
         "old": false,
         "key": "follow",
         "status": "Unchanged"
      },
      "host": {
         "old": "hexlet.io",
         "key": "host",
         "status": "Unchanged"
      },
      "proxy": {
         "old": "123.234.53.22",
         "key": "proxy",
         "status": "Unchanged"
      },
      "secret": {
         "old": null,
         "key": "secret",
         "status": "Unchanged"
      },
      "timeout": {
         "old": 50,
         "key": "timeout",
         "status": "Unchanged"
      }
   },
   "key": "root",
   "status": "Container"
}` + "\n",
			nil,
		},
		{
			"run with unsupported format",
			[]string{
				"gendiff",
				"--format",
				"unsupported",
				"../parser/testdata/json/simple.json",
				"../parser/testdata/json/simple.json",
			},
			"",
			errUnsupportedFormat,
		},
		{
			"run with missing first file",
			[]string{
				"gendiff",
				"../parser/testdata/json/simple.json",
			},
			"",
			errInvalidArguments,
		},
		{
			"run with missing second file",
			[]string{
				"gendiff",
				"../parser/testdata/json/simple.json",
			},
			"",
			errInvalidArguments,
		},
		{
			"run with invalid first file",
			[]string{
				"gendiff",
				"data/json/simple.json",
				"../parser/testdata/json/simple.json",
			},
			"",
			errValidatePath,
		},
		{
			"run with invalid second file",
			[]string{
				"gendiff",
				"../parser/testdata/json/simple.json",
				"data/json/simple.json",
			},
			"",
			errValidatePath,
		},
		{
			"run without arguments",
			[]string{
				"gendiff",
			},
			"",
			errInvalidArguments,
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			var output bytes.Buffer

			cmd := New()
			cmd.Writer = &output

			err := cmd.Run(context.Background(), td.args)
			if err == nil {
				assert.Equal(t, td.want, output.String())
			} else {
				require.ErrorIs(t, err, td.wantErr)
			}
		})
	}
}
