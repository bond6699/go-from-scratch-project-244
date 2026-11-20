package formatter

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"code/internal/analyzer"
)

const countSpacePerIndent = 4

type Stylish struct{}

func NewFormatterStylish() Stylish {
	return Stylish{}
}

func (fs Stylish) Format(node *analyzer.ContainerNode) (string, error) {
	builder := &strings.Builder{}

	builder.WriteString("{")

	for _, key := range fs.getSortedChildKeys(node) {
		if err := fs.formatNode(builder, 1, node.Children()[key]); err != nil {
			return "", err
		}
	}

	builder.WriteString("\n}")

	return builder.String(), nil
}

func (fs Stylish) formatNode(builder *strings.Builder, indent int, node analyzer.DiffNode) error {
	switch v := node.(type) {
	case *analyzer.ContainerNode:
		return fs.formatContainer(builder, indent, v)
	case *analyzer.AddedNode:
		return fs.formatAdded(builder, indent, v)
	case *analyzer.DeletedNode:
		return fs.formatDeleted(builder, indent, v)
	case *analyzer.ChangedNode:
		return fs.formatChanged(builder, indent, v)
	case *analyzer.UnchangedNode:
		return fs.formatUnchanged(builder, indent, v)
	}

	return nil
}

func (fs Stylish) formatContainer(builder *strings.Builder, indent int, node *analyzer.ContainerNode) error {
	fmt.Fprintf(builder, "\n%s%s: {", getIndent(indent), node.Key())

	for _, key := range fs.getSortedChildKeys(node) {
		if err := fs.formatNode(builder, indent+1, node.Children()[key]); err != nil {
			return err
		}
	}

	fmt.Fprintf(builder, "\n%s}", getIndent(indent))

	return nil
}

func (fs Stylish) formatAdded(builder *strings.Builder, indent int, node *analyzer.AddedNode) error {
	fmt.Fprintf(
		builder,
		"\n%s + %s: %v",
		getMarkedIndent(indent),
		node.Key(),
		fs.formatValue(node.Value(), indent+1),
	)

	return nil
}

func (fs Stylish) formatDeleted(builder *strings.Builder, indent int, node *analyzer.DeletedNode) error {
	fmt.Fprintf(
		builder,
		"\n%s - %s: %v",
		getMarkedIndent(indent),
		node.Key(),
		fs.formatValue(node.Value(), indent+1),
	)

	return nil
}

func (fs Stylish) formatChanged(builder *strings.Builder, indent int, node *analyzer.ChangedNode) error {
	fmt.Fprintf(
		builder,
		"\n%s - %s: %v",
		getMarkedIndent(indent),
		node.Key(),
		fs.formatValue(node.OldValue(), indent+1),
	)
	fmt.Fprintf(
		builder,
		"\n%s + %s: %v",
		getMarkedIndent(indent),
		node.Key(),
		fs.formatValue(node.NewValue(), indent+1),
	)

	return nil
}

func (fs Stylish) formatUnchanged(builder *strings.Builder, indent int, node *analyzer.UnchangedNode) error {
	fmt.Fprintf(builder, "\n%s%s: %v", getIndent(indent), node.Key(), fs.formatValue(node.Value(), indent+1))

	return nil
}

func (fs Stylish) formatValue(value any, indent int) any {
	switch v := value.(type) {
	case map[string]any:
		return fs.formatNestedMap(v, indent)
	case []any:
		return fs.formatNestedArray(v, indent)
	case nil:
		return "null"
	default:
		return value
	}
}

func (fs Stylish) formatNestedMap(value map[string]any, indent int) string {
	builder := &strings.Builder{}

	fmt.Fprintf(builder, "{")

	for _, key := range getSortedKeys(value) {
		fmt.Fprintf(builder, "\n%s%s: %v", getIndent(indent), key, fs.formatValue(value[key], indent+1))
	}

	fmt.Fprintf(builder, "\n%s}", getIndent(indent-1))

	return builder.String()
}

func (fs Stylish) formatNestedArray(value []any, indent int) string {
	builder := &strings.Builder{}

	fmt.Fprintf(builder, "[")

	for _, v := range value {
		fmt.Fprintf(builder, "\n%s%v", getIndent(indent), fs.formatValue(v, indent+1))
	}

	fmt.Fprintf(builder, "\n%s]", getIndent(indent-1))

	return builder.String()
}

func (fs Stylish) getSortedChildKeys(node *analyzer.ContainerNode) []string {
	keys := slices.Collect(maps.Keys(node.Children()))
	slices.Sort(keys)

	return keys
}

func getSortedKeys(valueMap map[string]any) []string {
	keys := slices.Collect(maps.Keys(valueMap))
	slices.Sort(keys)

	return keys
}

func getIndent(indent int) string {
	return strings.Repeat(" ", countSpacePerIndent*indent)
}

func getMarkedIndent(indent int) string {
	indentStr := strings.Repeat(" ", countSpacePerIndent*(indent-1))
	indentStr += " "

	return indentStr
}
