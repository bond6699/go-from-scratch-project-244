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
	OldValue any                  `json:"old,omitempty"`
	NewValue any                  `json:"new,omitempty"`
	Children map[string]*DiffNode `json:"child,omitempty"`
	Key      string               `json:"key,omitempty"`
	Status   DiffStatus           `json:"status,omitempty"`
}

func NewDiffNode(key string, status DiffStatus, oldValue, newValue any) *DiffNode {
	return &DiffNode{
		Key:      key,
		Status:   status,
		OldValue: oldValue,
		NewValue: newValue,
	}
}

func Analyze(data1, data2 map[string]any) *DiffNode {
	root := NewDiffNode("root", Container, nil, nil)
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
		return NewDiffNode(key, Added, nil, newValue)
	}

	if !newExist {
		return NewDiffNode(key, Deleted, oldValue, nil)
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
			return NewDiffNode(key, Unchanged, oldArray, newArray)
		}

		return NewDiffNode(key, Changed, oldArray, newArray)
	}

	if oldValue == newValue {
		return NewDiffNode(key, Unchanged, oldValue, nil)
	}

	return NewDiffNode(key, Changed, oldValue, newValue)
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
