package formatter

import (
	"testing"

	"code/internal/analyzer"
)

func TestPlainFormat(t *testing.T) {
	tests := []struct {
		name  string
		input analyzer.DiffNode
		want  string
	}{
		{
			"simple values",
			simpleDiffInput(),
			`Property 'follow' was updated. From false to true
Property 'proxy' was updated. From '123.234.53.22' to '123.235.53.22'
Property 'timeout' was updated. From 50 to 100`,
		},
		{
			"value with arrays",
			withArrayDiffInput(),
			`Property 'numbers' was updated. From [1 6 9] to [1 6 9 22.6]
Property 'others' was updated. From [hello world null] to [hello world]`,
		},
		{
			"complex value",
			complexDiffInput(),
			`Property 'common.setting1' was removed
Property 'group1.nest' was added with value: [complex value]
Property 'group1.secret' was added with value: null
Property 'group2' was added with value: [complex value]`,
		},
		{
			"empty value",
			analyzer.DiffNode{
				Key:      "root",
				Status:   analyzer.Container,
				Children: map[string]*analyzer.DiffNode{},
			},
			"",
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			got, err := FormatterPlain{}.Format(&td.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != td.want {
				t.Errorf("got:\n%s\nwant:\n%s", got, td.want)
			}
		})
	}
}
