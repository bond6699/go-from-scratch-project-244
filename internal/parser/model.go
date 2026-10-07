package parser

import "slices"

type ModelType int

const (
	Primitive ModelType = iota
	Object
	Array
)

type ModelStatus string

const (
	Unchanged ModelStatus = "Unchanged"
	Added     ModelStatus = "Added"
	Deleted   ModelStatus = "Deleted"
	Changed   ModelStatus = "Changed"
)

type Model struct {
	Root *Model
	Type ModelType `json:"-"`
	Key  string    `json:"key"`

	Value    any `json:"value,omitempty"`
	NewValue any `json:"new_value,omitempty"`

	Status ModelStatus `json:"status"`

	Children map[string]*Model `json:"children,omitempty"` // если здесь DELETED то смотреть в AddedChildren
	// Если children удалён, то
	Keys []string `json:"keys,omitempty"`

	AddedChildren map[string]*Model //
	AddedKeys     []string

	Items []*Model `json:"items,omitempty"`
}

func (m *Model) Add(key string, value any) {
	child := New(key, value)
	child.Root = m

	m.Children[key] = child
	m.Keys = append(m.Keys, key)
}

func (m *Model) AddItem(value any) {
	item := New("", value)
	item.Root = m
	m.Items = append(m.Items, item)
}

func (m *Model) HasKey(key string) bool {
	_, ok := m.Children[key]
	return ok
}

func (m *Model) HasItem(index int) bool {
	return index < len(m.Items)
}

func (m *Model) SetStatus(status ModelStatus) {
	m.Status = status
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
