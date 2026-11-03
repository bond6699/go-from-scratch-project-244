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
				"timeout": int64(50),
				"proxy":   "123.234.53.22",
				"follow":  false,
				"secret":  "null",
			},
			map[string]any{
				"host":    "hexlet.io",
				"timeout": int64(100),
				"proxy":   "123.235.53.22",
				"follow":  true,
				"secret":  "null",
			},
			*ContainerNode("root", map[string]*DiffNode{
				"host":    UnchangedNode("host", "hexlet.io"),
				"timeout": ChangedNode("timeout", int64(50), int64(100)),
				"proxy":   ChangedNode("proxy", "123.234.53.22", "123.235.53.22"),
				"follow":  ChangedNode("follow", false, true),
				"secret":  UnchangedNode("secret", "null"),
			}),
		},
		{
			"detect equal simple maps",
			map[string]any{
				"host":    "hexlet.io",
				"timeout": int64(50),
				"proxy":   "123.234.53.22",
			},
			map[string]any{
				"host":    "hexlet.io",
				"timeout": int64(50),
				"proxy":   "123.234.53.22",
			},
			*ContainerNode("root", map[string]*DiffNode{
				"host":    UnchangedNode("host", "hexlet.io"),
				"timeout": UnchangedNode("timeout", int64(50)),
				"proxy":   UnchangedNode("proxy", "123.234.53.22"),
			}),
		},
		{
			"detect additions and removals in complex maps",
			map[string]any{
				"common": map[string]any{
					"setting1": "Value 1",
					"setting2": int64(200),
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
					"setting2": int64(200),
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
					"abc": int64(12345),
					"deep": map[string]any{
						"id": int64(45),
					},
				},
			},
			*ContainerNode("root", map[string]*DiffNode{
				"common": ContainerNode("common", map[string]*DiffNode{
					"setting1": DeletedNode("setting1", "Value 1"),
					"setting2": UnchangedNode("setting2", int64(200)),
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
					"abc": int64(12345),
					"deep": map[string]any{
						"id": int64(45),
					},
				}),
			}),
		},
		{
			"detect additions and removals in arrays",
			map[string]any{
				"numbers":   []any{int64(1), int64(6), int64(9)},
				"others":    []any{"hello", "world", "null"},
				"languages": []any{"GO", "RUST", "PYTHON"},
			},
			map[string]any{
				"numbers":   []any{int64(1), int64(6), int64(9), 22.6},
				"others":    []any{"hello", "world"},
				"languages": []any{"GO", "RUST", "PYTHON"},
			},
			*ContainerNode("root", map[string]*DiffNode{
				"languages": UnchangedNode("languages", []any{"GO", "RUST", "PYTHON"}),
				"numbers": ChangedNode(
					"numbers",
					[]any{int64(1), int64(6), int64(9)},
					[]any{int64(1), int64(6), int64(9), 22.6},
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
					"setting2": int64(200),
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
					"setting2": int64(200),
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
					"setting2": int64(200),
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
					"setting2": int64(200),
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
					"item1": int64(1),
					"item2": int64(2),
				},
			},
			map[string]any{
				"object": "hello, world",
			},
			*ContainerNode("root", map[string]*DiffNode{
				"object": ChangedNode(
					"object",
					map[string]any{
						"item1": int64(1),
						"item2": int64(2),
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
					"item1": int64(1),
					"item2": int64(2),
				},
			},
			*ContainerNode("root", map[string]*DiffNode{
				"object": ChangedNode(
					"object",
					"hello, world",
					map[string]any{
						"item1": int64(1),
						"item2": int64(2),
					},
				),
			}),
		},
		{
			"detect change between large integers beyond float64 precision",
			map[string]any{
				"id": int64(9007199254740993), // 2^53 + 1
			},
			map[string]any{
				"id": int64(9007199254740994), // 2^53 + 2
			},
			*ContainerNode("root", map[string]*DiffNode{
				"id": ChangedNode("id", int64(9007199254740993), int64(9007199254740994)),
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
