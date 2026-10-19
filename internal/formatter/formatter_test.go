package formatter

import (
	"testing"

	"code/internal/analyzer"
)

func TestGetFormatter(t *testing.T) {
	tests := []struct {
		format string
		want   Formatter
	}{
		{
			"stylish",
			FormatterStylish{},
		},
		{
			"plain",
			FormatterPlain{},
		},
		{
			"json",
			FormatterJSON{},
		},
		{
			"",
			FormatterStylish{},
		},
	}

	for _, td := range tests {
		t.Run(td.format, func(t *testing.T) {
			if got := getFormatter(td.format); got != td.want {
				t.Fatalf("got %T, want %T", got, td.want)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		format  string
		diff    analyzer.DiffNode
		want    string
		wantErr bool
	}{
		{
			"stylish",
			analyzer.DiffNode{
				Key:    "root",
				Status: analyzer.Container,
				Children: map[string]*analyzer.DiffNode{
					"timeout": {
						Key:      "timeout",
						Status:   analyzer.Changed,
						OldValue: 50.0,
						NewValue: 100.0,
					},
				},
			},
			`{
  - timeout: 50
  + timeout: 100
}`,
			false,
		},
		{
			"plain",
			analyzer.DiffNode{
				Key:    "root",
				Status: analyzer.Container,
				Children: map[string]*analyzer.DiffNode{
					"timeout": {
						Key:      "timeout",
						Status:   analyzer.Changed,
						OldValue: 50.0,
						NewValue: 100.0,
					},
				},
			},
			"Property 'timeout' was updated. From 50 to 100\n",
			false,
		},
		{
			"json",
			analyzer.DiffNode{
				Key:    "root",
				Status: analyzer.Container,
				Children: map[string]*analyzer.DiffNode{
					"timeout": {
						Key:      "timeout",
						Status:   analyzer.Changed,
						OldValue: 50.0,
						NewValue: 100.0,
					},
				},
			},
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
			false,
		},
		{
			"",
			analyzer.DiffNode{
				Key:    "root",
				Status: analyzer.Container,
				Children: map[string]*analyzer.DiffNode{
					"timeout": {
						Key:      "timeout",
						Status:   analyzer.Changed,
						OldValue: 50.0,
						NewValue: 100.0,
					},
				},
			},
			`{
  - timeout: 50
  + timeout: 100
}`,
			false,
		},
	}

	for _, td := range tests {
		name := td.format
		if name == "" {
			name = "default"
		}

		t.Run(name, func(t *testing.T) {
			got, err := Format(td.format, &td.diff)

			if (err != nil) != td.wantErr {
				t.Fatalf("Format() error = %v, wantErr %v", err, td.wantErr)
			}

			if got != td.want {
				t.Errorf("Format() got:\n%q\nwant:\n%q", got, td.want)
			}
		})
	}
}
