package formatter

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"code/internal/analyzer"
)

type FormatterPlain struct{}

func NewFormatterPlain() FormatterPlain {
	return FormatterPlain{}
}

func (fp FormatterPlain) Format(node *analyzer.DiffNode) (string, error) {
	var builder strings.Builder

	for _, key := range fp.sortedChildKeys(node) {
		fp.formatNode(&builder, []string{}, node.Children[key])
	}

	return builder.String(), nil
}

func (fp FormatterPlain) formatNode(builder *strings.Builder, path []string, node *analyzer.DiffNode) {
	childPath := append(path, node.Key)
	hasChild := node.Children != nil

	switch node.Status {
	case analyzer.Unchanged:
		if hasChild {
			for _, key := range fp.sortedChildKeys(node) {
				fp.formatNode(builder, childPath, node.Children[key])
			}
		}
	case analyzer.Added:
		fmt.Fprintf(
			builder, "Property '%s' was added with value: %v\n",
			strings.Join(childPath, "."),
			fp.formatValue(node.NewValue, hasChild),
		)
	case analyzer.Deleted:
		fmt.Fprintf(builder, "Property '%s' was removed\n", strings.Join(childPath, "."))

	case analyzer.Changed:
		fmt.Fprintf(
			builder, "Property '%s' was updated. From %v to %v\n",
			strings.Join(childPath, "."),
			fp.formatValue(node.OldValue, hasChild),
			fp.formatValue(node.NewValue, hasChild),
		)

	}
}

func (fp FormatterPlain) sortedChildKeys(node *analyzer.DiffNode) []string {
	keys := slices.Collect(maps.Keys(node.Children))
	slices.Sort(keys)
	return keys
}

func (fp FormatterPlain) formatValue(value any, hasChild bool) any {
	_, isMap := value.(map[string]any)
	if hasChild || isMap {
		return "[complex value]"
	}

	if v, ok := value.(string); ok && v != "null" {
		value = fmt.Sprintf("'%s'", v)
	}
	return value
}
