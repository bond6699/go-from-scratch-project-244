package formatter

import (
	"encoding/json"
	"errors"
	"fmt"

	"code/internal/analyzer"
)

var errFormatToJSON = errors.New("failed to format to JSON")

type JSON struct{}

func NewFormatterJSON() JSON {
	return JSON{}
}

func (JSON) Format(node *analyzer.ContainerNode) (string, error) {
	resultBytes, err := json.MarshalIndent(node, "", "   ")
	if err != nil {
		return "", fmt.Errorf("%w: %w", errFormatToJSON, err)
	}

	return string(resultBytes), nil
}
