package parser

import "slices"

type ModelType int

const (
	Primitive = iota
	Object
	Array
)

type Model struct {
	Type  ModelType `json:"-"`
	Key   string    `json:"key"`
	Value any       `json:"value,omitempty"`

	Status string `json:"-"`

	Children map[string]*Model `json:"children,omitempty"`
	Items    []*Model          `json:"items,omitempty"`
	Keys     []string          `json:"keys,omitempty"`
}

func (m *Model) Add(key string, value any) {
	child := New(key, value)

	m.Children[key] = child
	m.Keys = append(m.Keys, key)
}

func (m *Model) AddItem(value any) {
	item := New("", value)
	m.Items = append(m.Items, item)
}

func NewPrimitive(key string, value any) *Model {
	return &Model{
		Type:  Primitive,
		Key:   key,
		Value: value,
	}
}

func NewObject(key string) *Model {
	return &Model{
		Type:     Object,
		Key:      key,
		Children: make(map[string]*Model),
	}
}

func NewArray(key string) *Model {
	return &Model{
		Type:  Array,
		Key:   key,
		Items: make([]*Model, 0),
	}
}

func New(key string, value any) *Model {
	switch v := value.(type) {
	case map[string]any:
		model := NewObject(key)
		for key, value := range v {
			model.Add(key, value)
		}
		slices.Sort(model.Keys)
		return model

	case []any:
		model := NewArray(key)
		for _, itemValue := range v {
			model.AddItem(itemValue)
		}
		return model
	default:
		return NewPrimitive(key, value)
	}
}
