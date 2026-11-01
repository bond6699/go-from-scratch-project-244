package analyzer

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyze(t *testing.T) {
	tests := []struct {
		name string
		map1 map[string]any
		map2 map[string]any
		want DiffNode
	}{
		{
			"detect additions and removals in simple maps",
			map[string]any{
				"host":    "hexlet.io",
				"timeout": 50.0,
				"proxy":   "123.234.53.22",
				"follow":  false,
				"secret":  "null",
			},
			map[string]any{
				"host":    "hexlet.io",
				"timeout": 100.0,
				"proxy":   "123.235.53.22",
				"follow":  true,
				"secret":  "null",
			},
			*ContainerNode("root", map[string]*DiffNode{
				"host":    UnchangedNode("host", "hexlet.io"),
				"timeout": ChangedNode("timeout", 50.0, 100.0),
				"proxy":   ChangedNode("proxy", "123.234.53.22", "123.235.53.22"),
				"follow":  ChangedNode("follow", false, true),
				"secret":  UnchangedNode("secret", "null"),
			}),
		},
		{
			"detect equal simple maps",
			map[string]any{
				"host":    "hexlet.io",
				"timeout": 50.0,
				"proxy":   "123.234.53.22",
			},
			map[string]any{
				"host":    "hexlet.io",
				"timeout": 50.0,
				"proxy":   "123.234.53.22",
			},
			*ContainerNode("root", map[string]*DiffNode{
				"host":    UnchangedNode("host", "hexlet.io"),
				"timeout": UnchangedNode("timeout", 50.0),
				"proxy":   UnchangedNode("proxy", "123.234.53.22"),
			}),
		},
		{
			"detect additions and removals in complex maps",
			map[string]any{
				"common": map[string]any{
					"setting1": "Value 1",
					"setting2": 200.0,
					"setting3": true,
					"setting6": map[string]any{
						"key": "value",
						"doge": map[string]any{
							"wow": "",
						},
					},
				},
				"group1": map[string]any{
					"baz": "bas",
					"foo": "bar",
				},
			},
			map[string]any{
				"common": map[string]any{
					"setting2": 200.0,
					"setting3": true,
					"setting6": map[string]any{
						"key": "value",
						"doge": map[string]any{
							"wow": "",
						},
					},
				},
				"group1": map[string]any{
					"baz": "bas",
					"foo": "bar",
					"nest": map[string]any{
						"key": "value",
					},
					"secret": "null",
				},
				"group2": map[string]any{
					"abc": 12345.0,
					"deep": map[string]any{
						"id": 45.0,
					},
				},
			},
			*ContainerNode("root", map[string]*DiffNode{
				"common": ContainerNode("common", map[string]*DiffNode{
					"setting1": DeletedNode("setting1", "Value 1"),
					"setting2": UnchangedNode("setting2", 200.0),
					"setting3": UnchangedNode("setting3", true),
					"setting6": ContainerNode("setting6", map[string]*DiffNode{
						"key": UnchangedNode("key", "value"),
						"doge": ContainerNode("doge", map[string]*DiffNode{
							"wow": UnchangedNode("wow", ""),
						}),
					}),
				}),
				"group1": ContainerNode("group1", map[string]*DiffNode{
					"baz": UnchangedNode("baz", "bas"),
					"foo": UnchangedNode("foo", "bar"),
					"nest": AddedNode("nest", map[string]any{
						"key": "value",
					}),
					"secret": AddedNode("secret", "null"),
				}),
				"group2": AddedNode("group2", map[string]any{
					"abc": 12345.0,
					"deep": map[string]any{
						"id": 45.0,
					},
				}),
			}),
		},
		{
			"detect additions and removals in arrays",
			map[string]any{
				"numbers":   []any{1.0, 6.0, 9.0},
				"others":    []any{"hello", "world", "null"},
				"languages": []any{"GO", "RUST", "PYTHON"},
			},
			map[string]any{
				"numbers":   []any{1.0, 6.0, 9.0, 22.6},
				"others":    []any{"hello", "world"},
				"languages": []any{"GO", "RUST", "PYTHON"},
			},
			*ContainerNode("root", map[string]*DiffNode{
				"languages": UnchangedNode("languages", []any{"GO", "RUST", "PYTHON"}),
				"numbers": ChangedNode(
					"numbers",
					[]any{1.0, 6.0, 9.0},
					[]any{1.0, 6.0, 9.0, 22.6},
				),
				"others": ChangedNode(
					"others",
					[]any{"hello", "world", "null"},
					[]any{"hello", "world"},
				),
			}),
		},
		{
			"detect equal empty maps",
			map[string]any{},
			map[string]any{},
			*ContainerNode("root", map[string]*DiffNode{}),
		},
		{
			"detect difference between empty map and map with added all keys",
			map[string]any{},
			map[string]any{
				"common": map[string]any{
					"setting1": "Value 1",
					"setting2": 200.0,
					"setting3": true,
					"setting6": map[string]any{
						"key": "value",
						"doge": map[string]any{
							"wow": "",
						},
					},
				},
				"group1": map[string]any{
					"baz": "bas",
					"foo": "bar",
				},
			},
			*ContainerNode("root", map[string]*DiffNode{
				"common": AddedNode("common", map[string]any{
					"setting1": "Value 1",
					"setting2": 200.0,
					"setting3": true,
					"setting6": map[string]any{
						"key": "value",
						"doge": map[string]any{
							"wow": "",
						},
					},
				}),
				"group1": AddedNode("group1", map[string]any{
					"baz": "bas",
					"foo": "bar",
				}),
			}),
		},
		{
			"detect difference between map with deleted all keys and empty map",
			map[string]any{
				"common": map[string]any{
					"setting1": "Value 1",
					"setting2": 200.0,
					"setting3": true,
					"setting6": map[string]any{
						"key": "value",
						"doge": map[string]any{
							"wow": "",
						},
					},
				},
				"group1": map[string]any{
					"baz": "bas",
					"foo": "bar",
				},
			},
			map[string]any{},
			*ContainerNode("root", map[string]*DiffNode{
				"common": DeletedNode("common", map[string]any{
					"setting1": "Value 1",
					"setting2": 200.0,
					"setting3": true,
					"setting6": map[string]any{
						"key": "value",
						"doge": map[string]any{
							"wow": "",
						},
					},
				}),
				"group1": DeletedNode("group1", map[string]any{
					"baz": "bas",
					"foo": "bar",
				}),
			}),
		},
		{
			"detects object to scalar changes",
			map[string]any{
				"object": map[string]any{
					"item1": 1.0,
					"item2": 2.0,
				},
			},
			map[string]any{
				"object": "hello, world",
			},
			*ContainerNode("root", map[string]*DiffNode{
				"object": ChangedNode(
					"object",
					map[string]any{
						"item1": 1.0,
						"item2": 2.0,
					},
					"hello, world",
				),
			}),
		},
		{
			"detect scalar to object changes",
			map[string]any{
				"object": "hello, world",
			},
			map[string]any{
				"object": map[string]any{
					"item1": 1.0,
					"item2": 2.0,
				},
			},
			*ContainerNode("root", map[string]*DiffNode{
				"object": ChangedNode(
					"object",
					"hello, world",
					map[string]any{
						"item1": 1.0,
						"item2": 2.0,
					},
				),
			}),
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			got := *Analyze(td.map1, td.map2)
			assert.Equal(t, td.want, got)
			wantJSON, _ := json.MarshalIndent(td.want, "", "  ")
			gotJSON, _ := json.MarshalIndent(got, "", "  ")

			assert.JSONEq(t, string(wantJSON), string(gotJSON))
		})
	}
}

func ContainerNode(key string, value map[string]*DiffNode) *DiffNode {
	node := NewDiffNode(key, Container)
	node.Children = value

	return node
}

func UnchangedNode(key string, value any) *DiffNode {
	node := NewDiffNode(key, Unchanged)
	node.SetOldValue(value)

	return node
}

func ChangedNode(key string, oldValue, newValue any) *DiffNode {
	node := NewDiffNode(key, Changed)
	node.SetOldValue(oldValue)
	node.SetNewValue(newValue)

	return node
}

func DeletedNode(key string, value any) *DiffNode {
	node := NewDiffNode(key, Deleted)
	node.SetOldValue(value)

	return node
}

func AddedNode(key string, value any) *DiffNode {
	node := NewDiffNode(key, Added)
	node.SetNewValue(value)

	return node
}
