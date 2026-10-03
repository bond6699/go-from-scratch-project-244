package parser

import "encoding/json"

type JsonParser struct{}

func (jp JsonParser) Parse(bytes []byte) (*Model, error) {
	var jsonData map[string]any

	err := json.Unmarshal(bytes, &jsonData)
	if err != nil {
		return &Model{}, err
	}

	data := CreateModel("json", nil)

	for k, v := range jsonData {
		data.AddPair(k, v)
	}

	return &data, nil
}
