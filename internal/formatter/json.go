package formatter

import (
	"encoding/json"
	"errors"
	"fmt"

	"code/internal/analyzer"
)

var errFormatToJSON = errors.New("failed to format to JSON")

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
