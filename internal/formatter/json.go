package formatter

import (
	"encoding/json"

	"code/internal/analyzer"
)

type FormatterJSON struct{}

func NewFormatterJSON() FormatterJSON {
	return FormatterJSON{}
}

func (FormatterJSON) Format(node *analyzer.DiffNode) (string, error) {
	resultBytes, err := json.MarshalIndent(node, "", "   ")
	if err != nil {
		return "", err
	}

	return string(resultBytes), nil
}
