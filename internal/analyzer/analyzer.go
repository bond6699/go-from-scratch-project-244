package analyzer

import (
	"encoding/json"
	"reflect"
)

type DiffNode interface {
	Key() string
	diffNode()
}

type AddedNode struct {
	key   string
	value any
}

func NewAddedNode(key string, value any) *AddedNode {
	return &AddedNode{key: key, value: value}
}

func (node *AddedNode) Key() string { return node.key }
func (node *AddedNode) Value() any  { return node.value }
func (*AddedNode) diffNode()        {}

func (node *AddedNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Value  any    `json:"new"`
		Key    string `json:"key"`
		Status string `json:"status"`
	}{
		Value:  node.value,
		Key:    node.key,
		Status: "Added",
	})
}

type DeletedNode struct {
	key   string
	value any
}

func NewDeletedNode(key string, value any) *DeletedNode {
	return &DeletedNode{key: key, value: value}
}

func (node *DeletedNode) Key() string { return node.key }
func (node *DeletedNode) Value() any  { return node.value }
func (*DeletedNode) diffNode()        {}

func (node *DeletedNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Value  any    `json:"old"`
		Key    string `json:"key"`
		Status string `json:"status"`
	}{
		Value:  node.value,
		Key:    node.key,
		Status: "Deleted",
	})
}

type UnchangedNode struct {
	key   string
	value any
}

func NewUnchangedNode(key string, value any) *UnchangedNode {
	return &UnchangedNode{key: key, value: value}
}
func (node *UnchangedNode) Key() string { return node.key }
func (node *UnchangedNode) Value() any  { return node.value }
func (*UnchangedNode) diffNode()        {}

func (node *UnchangedNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Value  any    `json:"old"`
		Key    string `json:"key"`
		Status string `json:"status"`
	}{
		Value:  node.value,
		Key:    node.key,
		Status: "Unchanged",
	})
}

type ChangedNode struct {
	key      string
	oldvalue any
	newValue any
}

func NewChangedNode(key string, oldValue, newValue any) *ChangedNode {
	return &ChangedNode{
		key:      key,
		oldvalue: oldValue,
		newValue: newValue,
	}
}
func (node *ChangedNode) Key() string   { return node.key }
func (node *ChangedNode) OldValue() any { return node.oldvalue }
func (node *ChangedNode) NewValue() any { return node.newValue }
func (*ChangedNode) diffNode()          {}

func (node *ChangedNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		OldValue any    `json:"old"`
		NewValue any    `json:"new"`
		Key      string `json:"key"`
		Status   string `json:"status"`
	}{
		OldValue: node.oldvalue,
		NewValue: node.newValue,
		Key:      node.key,
		Status:   "Changed",
	})
}

type ContainerNode struct {
	key      string
	children map[string]DiffNode
}

func NewContainerNode(key string, children map[string]DiffNode) *ContainerNode {
	return &ContainerNode{key: key, children: children}
}

func (node *ContainerNode) Key() string                   { return node.key }
func (node *ContainerNode) Children() map[string]DiffNode { return node.children }
func (*ContainerNode) diffNode()                          {}

func (node *ContainerNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Children map[string]DiffNode `json:"child,omitempty"`
		Key      string              `json:"key"`
		Status   string              `json:"status"`
	}{
		Children: node.children,
		Key:      node.key,
		Status:   "Container",
	})
}

func Analyze(key string, data1, data2 map[string]any) DiffNode {
	children := map[string]DiffNode{}
	keys := getKeysUnion(data1, data2)

	for _, k := range keys {
		oldValue, oldExist := data1[k]
		newValue, newExist := data2[k]

		children[k] = AnalyzeValue(k, oldValue, oldExist, newValue, newExist)
	}

	return NewContainerNode(key, children)
}

func AnalyzeValue(key string, oldValue any, oldExist bool, newValue any, newExist bool) DiffNode {
	if !oldExist {
		return NewAddedNode(key, newValue)
	}

	if !newExist {
		return NewDeletedNode(key, oldValue)
	}

	oldMap, oldIsMap := oldValue.(map[string]any)
	newMap, newIsMap := newValue.(map[string]any)

	if oldIsMap && newIsMap {
		return Analyze(key, oldMap, newMap)
	}

	oldArray, oldIsArray := oldValue.([]any)
	newArray, newIsArray := newValue.([]any)

	if oldIsArray && newIsArray {
		if reflect.DeepEqual(oldArray, newArray) {
			return NewUnchangedNode(key, oldArray)
		}

		return NewChangedNode(key, oldArray, newArray)
	}

	if oldValue == newValue {
		return NewUnchangedNode(key, oldValue)
	}

	return NewChangedNode(key, oldValue, newValue)
}

func getKeysUnion(data1, data2 map[string]any) []string {
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
