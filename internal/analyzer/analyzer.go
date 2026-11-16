package analyzer

import (
	"reflect"
	"slices"
)

type DiffStatus string

var (
	Unchanged DiffStatus = "Unchanged"
	Changed   DiffStatus = "Changed"
	Added     DiffStatus = "Added"
	Deleted   DiffStatus = "Deleted"
	Container DiffStatus = "Container"
)

type DiffNode struct {
	OldValue *any                 `json:"old,omitempty"`
	NewValue *any                 `json:"new,omitempty"`
	Children map[string]*DiffNode `json:"child,omitempty"`
	Key      string               `json:"key,omitempty"`
	Status   DiffStatus           `json:"status,omitempty"`
}

func (node *DiffNode) SetOldValue(oldValue any) {
	node.OldValue = &oldValue
}

func (node *DiffNode) SetNewValue(newValue any) {
	node.NewValue = &newValue
}

func (node *DiffNode) GetOldValue() (ok bool, value any) {
	if node.OldValue == nil {
		return false, nil
	}

	return true, *node.OldValue
}

func (node *DiffNode) GetNewValue() (ok bool, value any) {
	if node.NewValue == nil {
		return false, nil
	}

	return true, *node.NewValue
}

func NewDiffNode(key string, status DiffStatus) *DiffNode {
	return &DiffNode{
		Key:    key,
		Status: status,
	}
}

func NewAddedNode(key string, value any) *DiffNode {
	node := NewDiffNode(key, Added)
	node.SetNewValue(value)

	return node
}

func NewDeletedNode(key string, value any) *DiffNode {
	node := NewDiffNode(key, Deleted)
	node.SetOldValue(value)

	return node
}

func NewUnchangedNode(key string, value any) *DiffNode {
	node := NewDiffNode(key, Unchanged)
	node.SetOldValue(value)

	return node
}

func NewChangedNode(key string, oldValue, newValue any) *DiffNode {
	node := NewDiffNode(key, Changed)
	node.SetOldValue(oldValue)
	node.SetNewValue(newValue)

	return node
}

func NewContainerNode(key string, children map[string]*DiffNode) *DiffNode {
	return &DiffNode{
		Key:      key,
		Status:   Container,
		Children: children,
	}
}

func Analyze(data1, data2 map[string]any) *DiffNode {
	children := map[string]*DiffNode{}
	keys := getSortedKeysUnion(data1, data2)

	for _, key := range keys {
		oldValue, oldExist := data1[key]
		newValue, newExist := data2[key]

		children[key] = AnalyzeValue(key, oldValue, oldExist, newValue, newExist)
	}

	return NewContainerNode("root", children)
}

func AnalyzeValue(key string, oldValue any, oldExist bool, newValue any, newExist bool) *DiffNode {
	if !oldExist {
		node := NewAddedNode(key, newValue)

		return node
	}

	if !newExist {
		node := NewDeletedNode(key, oldValue)

		return node
	}

	oldMap, oldIsMap := oldValue.(map[string]any)
	newMap, newIsMap := newValue.(map[string]any)

	if oldIsMap && newIsMap {
		node := Analyze(oldMap, newMap)
		node.Key = key

		return node
	}

	oldArray, oldIsArray := oldValue.([]any)
	newArray, newIsArray := newValue.([]any)

	if oldIsArray && newIsArray {
		if reflect.DeepEqual(oldArray, newArray) {
			node := NewUnchangedNode(key, oldArray)

			return node
		}

		node := NewChangedNode(key, oldArray, newArray)

		return node
	}

	if oldValue == newValue {
		node := NewUnchangedNode(key, oldValue)

		return node
	}

	node := NewChangedNode(key, oldValue, newValue)

	return node
}

func getSortedKeysUnion(data1, data2 map[string]any) []string {
	keysMap := make(map[string]struct{})

	for key := range data1 {
		keysMap[key] = struct{}{}
	}

	for key := range data2 {
		keysMap[key] = struct{}{}
	}

	keys := []string{}
	for key := range keysMap {
		keys = append(keys, key)
	}

	slices.Sort(keys)

	return keys
}
