package code

import (
	"code/internal/analyzer"
	"code/internal/formatter"
	"code/internal/parser"
)

func GenDiff(filepath1, filepath2, format string) (result string, err error) {
	data1, err := parser.Parse(filepath1)
	if err != nil {
		return "", err
	}

	data1 = parser.Normalize(data1)

	data2, err := parser.Parse(filepath2)
	if err != nil {
		return "", err
	}

	data2 = parser.Normalize(data2)

	diff := analyzer.Analyze(data1, data2)

	result, err = formatter.Format(format, diff)
	if err != nil {
		return "", err
	}

	return result, nil
}
