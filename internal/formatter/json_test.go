package formatter

import (
	"testing"

	"code/internal/analyzer"
)

func TestJSONFormat(t *testing.T) {
	tests := []struct {
		name    string
		input   analyzer.DiffNode
		want    string
		wantErr bool
	}{
		{
			"simple values",
			analyzer.DiffNode{
				Key:    "root",
				Status: analyzer.Container,
				Children: map[string]*analyzer.DiffNode{
					"follow": {
						Key:      "follow",
						Status:   analyzer.Changed,
						OldValue: false,
						NewValue: true,
					},
					"host": {
						Key:      "host",
						Status:   analyzer.Unchanged,
						OldValue: "hexlet.io",
						NewValue: nil,
					},
					"proxy": {
						Key:      "proxy",
						Status:   analyzer.Changed,
						OldValue: "123.234.53.22",
						NewValue: "123.235.53.22",
					},
					"secret": {
						Key:      "secret",
						Status:   analyzer.Unchanged,
						OldValue: "null",
						NewValue: nil,
					},
					"timeout": {
						Key:      "timeout",
						Status:   analyzer.Changed,
						OldValue: 50.0,
						NewValue: 100.0,
					},
				},
			},
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
			false,
		},
		{
			"value with arrays",
			analyzer.DiffNode{
				Key:    "root",
				Status: analyzer.Container,
				Children: map[string]*analyzer.DiffNode{
					"languages": {
						Key:      "languages",
						Status:   analyzer.Unchanged,
						OldValue: []any{"GO", "RUST", "PYTHON"},
						NewValue: nil,
					},
					"numbers": {
						Key:      "numbers",
						Status:   analyzer.Changed,
						OldValue: []any{1.0, 6.0, 9.0},
						NewValue: []any{1.0, 6.0, 9.0, 22.6},
					},
					"others": {
						Key:      "others",
						Status:   analyzer.Changed,
						OldValue: []any{"hello", "world", "null"},
						NewValue: []any{"hello", "world"},
					},
				},
			},
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
			false,
		},
		{
			"complex value",
			analyzer.DiffNode{
				Key:    "root",
				Status: analyzer.Container,
				Children: map[string]*analyzer.DiffNode{
					"common": {
						Key:    "common",
						Status: analyzer.Container,
						Children: map[string]*analyzer.DiffNode{
							"setting1": {
								Key:      "setting1",
								Status:   analyzer.Deleted,
								OldValue: "Value 1",
								NewValue: nil,
							},
							"setting2": {
								Key:      "setting2",
								Status:   analyzer.Unchanged,
								OldValue: 200.0,
								NewValue: nil,
							},
							"setting3": {
								Key:      "setting3",
								Status:   analyzer.Unchanged,
								OldValue: true,
								NewValue: nil,
							},
							"setting6": {
								Key:    "setting6",
								Status: analyzer.Container,
								Children: map[string]*analyzer.DiffNode{
									"doge": {
										Key:    "doge",
										Status: analyzer.Container,
										Children: map[string]*analyzer.DiffNode{
											"wow": {
												Key:      "wow",
												Status:   analyzer.Unchanged,
												OldValue: "",
												NewValue: nil,
											},
										},
									},
									"key": {
										Key:      "key",
										Status:   analyzer.Unchanged,
										OldValue: "value",
										NewValue: nil,
									},
								},
							},
						},
					},
					"group1": {
						Key:    "group1",
						Status: analyzer.Container,
						Children: map[string]*analyzer.DiffNode{
							"baz": {
								Key:      "baz",
								Status:   analyzer.Unchanged,
								OldValue: "bas",
								NewValue: nil,
							},
							"foo": {
								Key:      "foo",
								Status:   analyzer.Unchanged,
								OldValue: "bar",
								NewValue: nil,
							},
							"nest": {
								Key:    "nest",
								Status: analyzer.Added,
								NewValue: map[string]any{
									"key": "value",
								},
							},
							"secret": {
								Key:      "secret",
								Status:   analyzer.Added,
								NewValue: "null",
							},
						},
					},
					"group2": {
						Key:    "group2",
						Status: analyzer.Added,
						NewValue: map[string]any{
							"abc": 12345.0,
							"deep": map[string]any{
								"id": 45.0,
							},
						},
					},
				},
			},
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
			false,
		},
		{
			"empty value",
			analyzer.DiffNode{
				Key:      "root",
				Status:   analyzer.Container,
				Children: map[string]*analyzer.DiffNode{},
			},
			`{
   "key": "root",
   "status": "Container"
}`,
			false,
		},
		{
			"incorrect value",
			analyzer.DiffNode{
				Key:    "root",
				Status: analyzer.Container,
				Children: map[string]*analyzer.DiffNode{
					"follow": {
						Key:      "follow",
						Status:   analyzer.Changed,
						OldValue: false,
						NewValue: func() {},
					},
				},
			},
			"",
			true,
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			got, err := FormatterJSON{}.Format(&td.input)
			assertError(t, td.wantErr, err)
			if got != td.want {
				t.Errorf("got:\n%s\nwant:\n%s", got, td.want)
			}
		})
	}
}

func assertError(t *testing.T, wantErr bool, got error) {
	t.Helper()
	if wantErr && got == nil {
		t.Fatal("expected an error, got nil")
	}

	if !wantErr && got != nil {
		t.Fatalf("unexpected error: %v", got)
	}
}
