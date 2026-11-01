package formatter

import (
	"encoding/json"
	"fmt"

	"code/internal/analyzer"
)

var errFormatToJSON = fmt.Errorf("failed to format to JSON")

type FormatterJSON struct{}

func NewFormatterJSON() FormatterJSON {
	return FormatterJSON{}
}

func (FormatterJSON) Format(node *analyzer.DiffNode) (string, error) {
	resultBytes, err := json.MarshalIndent(node, "", "   ")
	if err != nil {
		return "", fmt.Errorf("%w: %w", errFormatToJSON, err)
	}

	return string(resultBytes), nil
}
