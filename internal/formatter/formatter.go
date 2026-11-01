package formatter

import (
	"fmt"

	"code/internal/analyzer"
)

var (
	errUnsupportedFormat = fmt.Errorf(
		"format is not supported; select one of the possible options (stylish, plain, json)",
	)
	errCreateFormatter = fmt.Errorf("failed to create formatter")
	errFormatting      = fmt.Errorf("failed formatting")
)

type Formatter interface {
	Format(*analyzer.DiffNode) (string, error)
}

func getFormatter(format string) (Formatter, error) {
	switch format {
	case "stylish":
		return NewFormatterStylish(), nil
	case "plain":
		return NewFormatterPlain(), nil
	case "json":
		return NewFormatterJSON(), nil
	default:
		return nil, fmt.Errorf("%q %w", format, errUnsupportedFormat)
	}
}

func Format(format string, diff *analyzer.DiffNode) (string, error) {
	formatter, err := getFormatter(format)
	if err != nil {
		return "", fmt.Errorf("%w %q: %w", errCreateFormatter, format, err)
	}

	result, err := formatter.Format(diff)
	if err != nil {
		err = fmt.Errorf("%w: %w", errFormatting, err)
	}

	return result, err
}
