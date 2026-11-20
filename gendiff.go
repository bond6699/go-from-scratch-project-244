package code

import (
	"errors"
	"fmt"

	"code/internal/analyzer"
	"code/internal/formatter"
	"code/internal/parser"
)

var (
	errFailedProcessFirstFile  = errors.New("failed to process first file")
	errFailedProcessSecondFile = errors.New("failed to process second file")
	errFailedFormat            = errors.New("failed to format")
	errFailedNormalize         = errors.New("failed to normalize")
)

const DefaultFormatter = "stylish"

func GenDiff(filepath1, filepath2, format string) (string, error) {
	if format == "" {
		format = DefaultFormatter
	}

	data1, err := parser.Parse(filepath1)
	if err != nil {
		return "", fmt.Errorf("%w %q: %w", errFailedProcessFirstFile, filepath1, err)
	}

	data2, err := parser.Parse(filepath2)
	if err != nil {
		return "", fmt.Errorf("%w %q: %w", errFailedProcessSecondFile, filepath2, err)
	}

	data1, err = parser.Normalize(data1)
	if err != nil {
		return "", fmt.Errorf("%w %q: %w", errFailedNormalize, filepath1, err)
	}

	data2, err = parser.Normalize(data2)
	if err != nil {
		return "", fmt.Errorf("%w %q: %w", errFailedNormalize, filepath2, err)
	}

	diff := analyzer.Analyze("root", data1, data2)

	result, err := formatter.Format(format, diff)
	if err != nil {
		return "", fmt.Errorf("%w %q, %q in %q: %w", errFailedFormat, filepath1, filepath2, format, err)
	}

	return result, nil
}
