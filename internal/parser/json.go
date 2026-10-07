package parser

import (
	"encoding/json"
)

type JSONParser struct{}

func (jp JSONParser) Parse(bytes []byte) (*Model, error) {
	var jsonData map[string]any

	err := json.Unmarshal(bytes, &jsonData)
	if err != nil {
		return &Model{}, err
	}

	data := New("root", jsonData)

	return data, nil
}
