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
	OldV     *any                 `json:"old,omitempty"`
	NewV     *any                 `json:"new,omitempty"`
	Children map[string]*DiffNode `json:"child,omitempty"`
	Key      string               `json:"key,omitempty"`
	Status   DiffStatus           `json:"status,omitempty"`
}

func NewDiffNode(key string, status DiffStatus) *DiffNode {
	return &DiffNode{
		Key:    key,
		Status: status,
	}
}

func (node *DiffNode) SetOldValue(oldValue any) {
	node.OldV = &oldValue
}

func (node *DiffNode) SetNewValue(newValue any) {
	node.NewV = &newValue
}

func (node *DiffNode) OldValue() any {
	return *node.OldV
}

func (node *DiffNode) NewValue() any {
	return *node.NewV
}

func Analyze(data1, data2 map[string]any) *DiffNode {
	root := NewDiffNode("root", Container)
	root.Children = map[string]*DiffNode{}

	keys := getSortedKeysUnion(data1, data2)

	for _, key := range keys {
		oldValue, oldExist := data1[key]
		newValue, newExist := data2[key]

		root.Children[key] = AnalyzeValue(key, oldValue, oldExist, newValue, newExist)
	}

	return root
}

func AnalyzeValue(key string, oldValue any, oldExist bool, newValue any, newExist bool) *DiffNode {
	if !oldExist {
		node := NewDiffNode(key, Added)
		node.SetNewValue(newValue)
		return node
	}

	if !newExist {
		node := NewDiffNode(key, Deleted)
		node.SetOldValue(oldValue)
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
			node := NewDiffNode(key, Unchanged)
			node.SetOldValue(oldArray)
			return node
		}

		node := NewDiffNode(key, Changed)
		node.SetOldValue(oldArray)
		node.SetNewValue(newArray)
		return node
	}

	if oldValue == newValue {
		node := NewDiffNode(key, Unchanged)
		node.SetOldValue(oldValue)
		return node
	}

	node := NewDiffNode(key, Changed)
	node.SetOldValue(oldValue)
	node.SetNewValue(newValue)
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
