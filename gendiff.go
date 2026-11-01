package code

import (
	"fmt"

	"code/internal/analyzer"
	"code/internal/formatter"
	"code/internal/parser"
)

func GenDiff(filepath1, filepath2, format string) (string, error) {
	data1, err := parser.Parse(filepath1)
	if err != nil {
		return "", fmt.Errorf("failed to process first file %q: %w", filepath1, err)
	}

	data2, err := parser.Parse(filepath2)
	if err != nil {
		return "", fmt.Errorf("failed to process second file %q: %w", filepath2, err)
	}

	data1 = parser.Normalize(data1)
	data2 = parser.Normalize(data2)
	diff := analyzer.Analyze(data1, data2)

	result, err := formatter.Format(format, diff)
	if err != nil {
		return "", fmt.Errorf("failed to format %q, %q in %q: %w", filepath1, filepath2, format, err)
	}

	return result, nil
}
