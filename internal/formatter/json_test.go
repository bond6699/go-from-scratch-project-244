package formatter

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"code/internal/analyzer"
)

func TestJSONFormat(t *testing.T) {
	tests := []struct {
		name    string
		diff    analyzer.DiffNode
		want    string
		wantErr error
	}{
		{
			"detect correct JSON result with simple values",
			simpleDiff,
			`{
   "child": {
      "follow": {
         "old": false,
         "new": true,
         "key": "follow",
         "status": "Changed"
      },
      "host": {
         "old": "hexlet.io",
         "key": "host",
         "status": "Unchanged"
      },
      "proxy": {
         "old": "123.234.53.22",
         "new": "123.235.53.22",
         "key": "proxy",
         "status": "Changed"
      },
      "secret": {
         "old": "null",
         "key": "secret",
         "status": "Unchanged"
      },
      "timeout": {
         "old": 50,
         "new": 100,
         "key": "timeout",
         "status": "Changed"
      }
   },
   "key": "root",
   "status": "Container"
}`,
			nil,
		},
		{
			"detect correct JSON result with array values",
			diffWithArrays,
			`{
   "child": {
      "languages": {
         "old": [
            "GO",
            "RUST",
            "PYTHON"
         ],
         "key": "languages",
         "status": "Unchanged"
      },
      "numbers": {
         "old": [
            1,
            6,
            9
         ],
         "new": [
            1,
            6,
            9,
            22.6
         ],
         "key": "numbers",
         "status": "Changed"
      },
      "others": {
         "old": [
            "hello",
            "world",
            "null"
         ],
         "new": [
            "hello",
            "world"
         ],
         "key": "others",
         "status": "Changed"
      }
   },
   "key": "root",
   "status": "Container"
}`,
			nil,
		},
		{
			"detect saved nil value in JSON result",
			diffWithNil,
			`{
   "child": {
      "toNullValue": {
         "old": "null",
         "new": null,
         "key": "toNullValue",
         "status": "Changed"
      },
      "toStringValue": {
         "old": null,
         "new": "null",
         "key": "toStringValue",
         "status": "Changed"
      }
   },
   "key": "root",
   "status": "Container"
}`,
			nil,
		},
		{
			"detect correct JSON result with complex values",
			complexDiff,
			`{
   "child": {
      "common": {
         "child": {
            "setting1": {
               "old": "Value 1",
               "key": "setting1",
               "status": "Deleted"
            },
            "setting2": {
               "old": 200,
               "key": "setting2",
               "status": "Unchanged"
            },
            "setting3": {
               "old": true,
               "key": "setting3",
               "status": "Unchanged"
            },
            "setting6": {
               "child": {
                  "doge": {
                     "child": {
                        "wow": {
                           "old": "",
                           "key": "wow",
                           "status": "Unchanged"
                        }
                     },
                     "key": "doge",
                     "status": "Container"
                  },
                  "key": {
                     "old": "value",
                     "key": "key",
                     "status": "Unchanged"
                  }
               },
               "key": "setting6",
               "status": "Container"
            }
         },
         "key": "common",
         "status": "Container"
      },
      "group1": {
         "child": {
            "baz": {
               "old": "bas",
               "key": "baz",
               "status": "Unchanged"
            },
            "foo": {
               "old": "bar",
               "key": "foo",
               "status": "Unchanged"
            },
            "nest": {
               "new": {
                  "key": "value"
               },
               "key": "nest",
               "status": "Added"
            },
            "secret": {
               "new": "null",
               "key": "secret",
               "status": "Added"
            },
            "secret2": {
               "new": null,
               "key": "secret2",
               "status": "Added"
            }
         },
         "key": "group1",
         "status": "Container"
      },
      "group2": {
         "new": {
            "abc": 12345,
            "deep": {
               "id": 45
            }
         },
         "key": "group2",
         "status": "Added"
      }
   },
   "key": "root",
   "status": "Container"
}`,
			nil,
		},
		{
			"detect correct JSON result with value contains punctuation marks",
			diffWithValueContainsPunctuationMarks,
			`{
   "child": {
      "docsIndex": {
         "old": "124/365-44, 124/365-45",
         "new": "124/365-45, 124/365-46",
         "key": "docsIndex",
         "status": "Changed"
      }
   },
   "key": "root",
   "status": "Container"
}`,
			nil,
		},
		{
			"detect correct JSON result with empty value",
			emptyDiff,
			`{
   "key": "root",
   "status": "Container"
}`,
			nil,
		},
		{
			"detect errFormatToJSON with incorrect value",
			*ContainerNode("root", map[string]*analyzer.DiffNode{
				"follow": UnchangedNode("follow", any(func() {})),
			}),
			"",
			errFormatToJSON,
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			got, err := FormatterJSON{}.Format(&td.diff)
			assert.Equal(t, td.want, got)
			require.ErrorIs(t, err, td.wantErr)
		})
	}
}

func TestJSONFormat_SemanticStructure(t *testing.T) {
	t.Run("semantic check of null vs string null and omitted keys", func(t *testing.T) {
		got, err := FormatterJSON{}.Format(&diffWithNil)
		require.NoError(t, err)

		var parsed map[string]any

		err = json.Unmarshal([]byte(got), &parsed)
		require.NoError(t, err)

		children, ok := parsed["child"].(map[string]any)
		require.True(t, ok)

		toNull, ok := children["toNullValue"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "Changed", toNull["status"])
		assert.Equal(t, "null", toNull["old"])
		assert.Nil(t, toNull["new"])
		_, hasNew := toNull["new"]
		assert.True(t, hasNew, "new key must be present with null value")

		toString, ok := children["toStringValue"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "Changed", toString["status"])
		assert.Nil(t, toString["old"])
		_, hasOld := toString["old"]
		assert.True(t, hasOld, "old key must be present with null value")
		assert.Equal(t, "null", toString["new"])
	})

	t.Run("semantic check of complex diff and deleted/added keys omission", func(t *testing.T) {
		got, err := FormatterJSON{}.Format(&complexDiff)
		require.NoError(t, err)

		var parsed map[string]any

		err = json.Unmarshal([]byte(got), &parsed)
		require.NoError(t, err)

		children := parsed["child"].(map[string]any)
		common := children["common"].(map[string]any)["child"].(map[string]any)

		// Deleted node: old must exist, new must be omitted
		setting1 := common["setting1"].(map[string]any)
		assert.Equal(t, "Deleted", setting1["status"])
		assert.Equal(t, "Value 1", setting1["old"])
		_, hasNew := setting1["new"]
		assert.False(t, hasNew, "new key must be omitted for Deleted node")

		// Added node with null: secret2
		group1 := children["group1"].(map[string]any)["child"].(map[string]any)
		secret2 := group1["secret2"].(map[string]any)
		assert.Equal(t, "Added", secret2["status"])
		_, hasOld := secret2["old"]
		assert.False(t, hasOld, "old key must be omitted for Added node")
		assert.Nil(t, secret2["new"])
		_, hasNewSecret2 := secret2["new"]
		assert.True(t, hasNewSecret2, "new key must be present with null value")
	})
}
