package formatter

import (
	"code/internal/analyzer"
)

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
