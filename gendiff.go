package code

import (
	"fmt"

	"code/internal/analyzer"
	"code/internal/formatter"
	"code/internal/parser"
)

var (
	errFailedProcessFirstFile  = fmt.Errorf("failed to process first file")
	errFailedProcessSecondFile = fmt.Errorf("failed to process second file")
	errFailedFormat            = fmt.Errorf("failed to format")
)

func GenDiff(filepath1, filepath2, format string) (string, error) {
	if format == "" {
		format = "stylish"
	}

	data1, err := parser.Parse(filepath1)
	if err != nil {
		return "", fmt.Errorf("%w %q: %w", errFailedProcessFirstFile, filepath1, err)
	}

	data2, err := parser.Parse(filepath2)
	if err != nil {
		return "", fmt.Errorf("%w %q: %w", errFailedProcessSecondFile, filepath2, err)
	}

	data1 = parser.Normalize(data1)
	data2 = parser.Normalize(data2)
	diff := analyzer.Analyze(data1, data2)

	result, err := formatter.Format(format, diff)
	if err != nil {
		return "", fmt.Errorf("%w %q, %q in %q: %w", errFailedFormat, filepath1, filepath2, format, err)
	}

	return result, nil
}
