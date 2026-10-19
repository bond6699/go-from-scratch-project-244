package formatter

import (
	"testing"

	"code/internal/analyzer"
)

func TestStylishFormat(t *testing.T) {
	tests := []struct {
		name  string
		input analyzer.DiffNode
		want  string
	}{
		{
			"simple values",
			simpleDiffInput(),
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
			"value with arrays",
			withArrayDiffInput(),
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
			"complex value",
			complexDiffInput(),
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
			"empty value",
			analyzer.DiffNode{
				Key:      "root",
				Status:   analyzer.Container,
				Children: map[string]*analyzer.DiffNode{},
			},
			`{
}`,
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			got, err := FormatterStylish{}.Format(&td.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != td.want {
				t.Errorf("got:\n%s\nwant:\n%s", got, td.want)
			}
		})
	}
}
