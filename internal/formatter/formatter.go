package formatter

import (
	"fmt"

	"code/internal/analyzer"
)

var ErrInvalidFormat = fmt.Errorf("invalid format")

type Formatter interface {
	Format(*analyzer.DiffNode) (string, error)
}

func getFormatter(format string) Formatter {
	switch format {
	case "stylish":
		return NewFormatterStylish()
	case "plain":
		return NewFormatterPlain()
	case "json":
		return NewFormatterJSON()
	default:
		return NewFormatterStylish()
	}
}

func Format(format string, diff *analyzer.DiffNode) (string, error) {
	return getFormatter(format).Format(diff)
}

func IsValidFormat(format string) bool {
	switch format {
	case "stylish":
		return true
	case "plain":
		return true
	case "json":
		return true
	default:
		return false
	}
}
