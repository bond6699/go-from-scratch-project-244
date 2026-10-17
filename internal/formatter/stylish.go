package formatter

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"

	"code/internal/analyzer"
)

type FormatterStylish struct{}

func NewFormatterStylish() FormatterStylish {
	return FormatterStylish{}
}

func (fs FormatterStylish) Format(node *analyzer.DiffNode) (string, error) {
	builder := &strings.Builder{}

	builder.WriteString("{")

	for _, key := range fs.getSortedChildKeys(node) {
		fs.formatNode(builder, 1, node.Children[key])
	}

	builder.WriteString("\n}")

	return builder.String(), nil
}

func (fs FormatterStylish) formatNode(builder *strings.Builder, indent int, node *analyzer.DiffNode) {
	switch node.Status {
	case analyzer.Container:
		fmt.Fprintf(builder, "\n%s%s: {", getIndent(indent), node.Key)

		for _, key := range fs.getSortedChildKeys(node) {
			fs.formatNode(builder, indent+1, node.Children[key])
		}

		fmt.Fprintf(builder, "\n%s}", getIndent(indent))

	case analyzer.Added:
		fmt.Fprintf(
			builder,
			"\n%s + %s: %v",
			getMarkedIndent(indent),
			node.Key,
			fs.formatValue(node.NewValue, indent+1),
		)
	case analyzer.Deleted:
		fmt.Fprintf(
			builder,
			"\n%s - %s: %v",
			getMarkedIndent(indent),
			node.Key,
			fs.formatValue(node.OldValue, indent+1),
		)
	case analyzer.Changed:
		fmt.Fprintf(
			builder,
			"\n%s - %s: %v",
			getMarkedIndent(indent),
			node.Key,
			fs.formatValue(node.OldValue, indent+1),
		)
		fmt.Fprintf(
			builder,
			"\n%s + %s: %v",
			getMarkedIndent(indent),
			node.Key,
			fs.formatValue(node.NewValue, indent+1),
		)
	case analyzer.Unchanged:
		fmt.Fprintf(builder, "\n%s%s: %v", getIndent(indent), node.Key, fs.formatValue(node.OldValue, indent+1))
	}
}

func (fs FormatterStylish) formatValue(value any, indent int) any {
	switch value.(type) {
	case map[string]any:
		return fs.formatNestedValue(value, indent)
	case []any:
		return fs.formatNestedValue(value, indent)
	default:
		return value
	}
}

func (fs FormatterStylish) formatNestedValue(value any, indent int) string {
	bytes, _ := json.MarshalIndent(value, "", getIndent(1))

	strSlice := strings.Split(string(bytes), "\n")
	for i, row := range strSlice {
		if i != 0 {
			strSlice[i] = getIndent(indent-1) + row
		}
	}

	result := strings.Join(strSlice, "\n")
	result = strings.Replace(result, "\"", "", len(strSlice)*5)
	result = strings.Replace(result, ",", "", len(strSlice)*5)

	return result
}

func (fs FormatterStylish) getSortedChildKeys(node *analyzer.DiffNode) []string {
	keys := slices.Collect(maps.Keys(node.Children))
	slices.Sort(keys)

	return keys
}

func getIndent(indent int) string {
	return strings.Repeat(" ", 4*indent)
}

func getMarkedIndent(indent int) string {
	indentStr := strings.Repeat(" ", 4*(indent-1))
	indentStr += " "

	return indentStr
}
