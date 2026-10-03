package parser

type Model struct {
	Key      string
	Value    any
	Status   int
	Children []Model
}

func CreateModel(key string, value any) Model {
	return Model{
		Key:   key,
		Value: value,
	}
}

func (m *Model) AddPair(key string, value any) {
	if mapV, ok := value.(map[string]any); !ok {
		pair := CreateModel(key, value)
		m.Children = append(m.Children, pair)
	} else {
		pair := CreateModel(key, nil)
		for k, v := range mapV {
			pair.AddPair(k, v)
		}
		m.Children = append(m.Children, pair)
	}
}
