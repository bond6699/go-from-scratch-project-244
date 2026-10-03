package parser

type YamlParser struct{}

func (yp YamlParser) Parse(bytes []byte) (*Model, error) {
	return &Model{}, nil
}
