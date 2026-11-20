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
				"../../testdata/simple.json",
				"../../testdata/simple.json",
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
				"../../testdata/simple.json",
				"../../testdata/simple.json",
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
				"../../testdata/simple.json",
				"../../testdata/simple.json",
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
				"../../testdata/simple.json",
				"../../testdata/simple.json",
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
				"../../testdata/simple.json",
				"../../testdata/simple.json",
			},
			"",
			errUnsupportedFormat,
		},
		{
			"run with missing second file",
			[]string{
				"gendiff",
				"../../testdata/simple.json",
			},
			"",
			errInvalidArguments,
		},
		{
			"run with too many arguments",
			[]string{
				"gendiff",
				"../../testdata/simple.json",
				"../../testdata/simple.json",
				"../../testdata/simple.json",
			},
			"",
			errInvalidArguments,
		},
		{
			"run with invalid first file",
			[]string{
				"gendiff",
				"data/json/simple.json",
				"../../testdata/simple.json",
			},
			"",
			errValidatePath,
		},
		{
			"run with invalid second file",
			[]string{
				"gendiff",
				"../../testdata/simple.json",
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

			if td.wantErr == nil {
				require.NoError(t, err)
				assert.Equal(t, td.want, output.String())
			} else {
				require.ErrorIs(t, err, td.wantErr)
			}
		})
	}
}
