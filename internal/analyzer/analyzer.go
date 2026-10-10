package analyzer

import (
	"reflect"
	"slices"
	"strconv"
)

type DiffStatus string

var (
	Unchanged DiffStatus = "Unchanged"
	Changed   DiffStatus = "Changed"
	Added     DiffStatus = "Added"
	Deleted   DiffStatus = "Deleted"
)

type DiffNode struct {
	Key      string               `json:"key,omitempty"`
	Status   DiffStatus           `json:"status,omitempty"`
	OldValue any                  `json:"old,omitempty"`
	NewValue any                  `json:"new,omitempty"`
	Children map[string]*DiffNode `json:"child,omitempty"`
}

func getKeys(data1 map[string]any, data2 map[string]any) []string {
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

	return keys
}

func Analyze(data1, data2 map[string]any) *DiffNode {
	root := &DiffNode{
		Status:   Unchanged,
		Children: map[string]*DiffNode{},
	}

	keys := getKeys(data1, data2)
	slices.Sort(keys)
	for _, key := range keys {
		oldValue, oldExist := data1[key]
		newValue, newExist := data2[key]

		root.Children[key] = AnalyzeValue(key, oldValue, oldExist, newValue, newExist)
	}

	return root
}

func AnalyzeValue(key string, oldValue any, oldExist bool, newValue any, newExist bool) *DiffNode {
	if !oldExist {
		return &DiffNode{
			Key:      key,
			NewValue: newValue,
			Status:   Added,
		}
	}

	if !newExist {
		return &DiffNode{
			Key:      key,
			OldValue: oldValue,
			Status:   Deleted,
		}
	}

	if reflect.TypeOf(oldValue) == reflect.TypeOf(newValue) {
		switch oldValue.(type) {
		case map[string]any:
			oldV := oldValue.(map[string]any)
			newV := newValue.(map[string]any)
			node := Analyze(oldV, newV)
			node.Key = key
			return node
		case []any:
			oldV := oldValue.([]any)
			newV := newValue.([]any)
			maxLen := max(len(oldV), len(newV))
			node := &DiffNode{
				Status:   Unchanged,
				Key:      key,
				Children: make(map[string]*DiffNode, maxLen),
			}

			for i := 0; i < maxLen; i++ {
				existOldItem := i < len(oldV)
				existNewItem := i < len(newV)

				var oldItemValue, newItemValue any

				if existOldItem {
					oldItemValue = oldV[i]
				}

				if existNewItem {
					newItemValue = newV[i]
				}

				index := strconv.Itoa(i)

				node.Children[index] = AnalyzeValue(
					index,
					oldItemValue, existOldItem,
					newItemValue, existNewItem,
				)
			}
			return node
		default:
			if oldValue == newValue {
				return &DiffNode{
					Key:      key,
					Status:   Unchanged,
					OldValue: oldValue,
				}
			} else {
				return &DiffNode{
					Key:      key,
					Status:   Changed,
					OldValue: oldValue,
					NewValue: newValue,
				}
			}
		}
	}
	return &DiffNode{
		Key:      key,
		Status:   Changed,
		OldValue: oldValue,
		NewValue: newValue,
	}
}
