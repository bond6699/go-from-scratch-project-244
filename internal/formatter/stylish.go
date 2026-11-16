package formatter

import (
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
		if err := fs.formatNode(builder, 1, node.Children[key]); err != nil {
			return "", err
		}
	}

	builder.WriteString("\n}")

	return builder.String(), nil
}

func (fs FormatterStylish) formatNode(builder *strings.Builder, indent int, node *analyzer.DiffNode) error {
	switch node.Status {
	case analyzer.Container:
		return fs.formatContainer(builder, indent, node)
	case analyzer.Added:
		return fs.formatAdded(builder, indent, node)
	case analyzer.Deleted:
		return fs.formatDeleted(builder, indent, node)
	case analyzer.Changed:
		return fs.formatChanged(builder, indent, node)
	case analyzer.Unchanged:
		return fs.formatUnchanged(builder, indent, node)
	}

	return nil
}

func (fs FormatterStylish) formatContainer(builder *strings.Builder, indent int, node *analyzer.DiffNode) error {
	fmt.Fprintf(builder, "\n%s%s: {", getIndent(indent), node.Key)

	for _, key := range fs.getSortedChildKeys(node) {
		if err := fs.formatNode(builder, indent+1, node.Children[key]); err != nil {
			return err
		}
	}

	fmt.Fprintf(builder, "\n%s}", getIndent(indent))

	return nil
}

func (fs FormatterStylish) formatAdded(builder *strings.Builder, indent int, node *analyzer.DiffNode) error {
	ok, newValue := node.GetNewValue()
	if !ok {
		return fmt.Errorf("%w: key %q", errGetNodeValue, node.Key)
	}

	fmt.Fprintf(
		builder,
		"\n%s + %s: %v",
		getMarkedIndent(indent),
		node.Key,
		fs.formatValue(newValue, indent+1),
	)

	return nil
}

func (fs FormatterStylish) formatDeleted(builder *strings.Builder, indent int, node *analyzer.DiffNode) error {
	ok, oldValue := node.GetOldValue()
	if !ok {
		return fmt.Errorf("%w: key %q", errGetNodeValue, node.Key)
	}

	fmt.Fprintf(
		builder,
		"\n%s - %s: %v",
		getMarkedIndent(indent),
		node.Key,
		fs.formatValue(oldValue, indent+1),
	)

	return nil
}

func (fs FormatterStylish) formatChanged(builder *strings.Builder, indent int, node *analyzer.DiffNode) error {
	okOld, oldValue := node.GetOldValue()
	if !okOld {
		return fmt.Errorf("%w: key %q", errGetNodeValue, node.Key)
	}

	okNew, newValue := node.GetNewValue()
	if !okNew {
		return fmt.Errorf("%w: key %q", errGetNodeValue, node.Key)
	}

	fmt.Fprintf(
		builder,
		"\n%s - %s: %v",
		getMarkedIndent(indent),
		node.Key,
		fs.formatValue(oldValue, indent+1),
	)
	fmt.Fprintf(
		builder,
		"\n%s + %s: %v",
		getMarkedIndent(indent),
		node.Key,
		fs.formatValue(newValue, indent+1),
	)

	return nil
}

func (fs FormatterStylish) formatUnchanged(builder *strings.Builder, indent int, node *analyzer.DiffNode) error {
	ok, oldValue := node.GetOldValue()
	if !ok {
		return fmt.Errorf("%w: key %q", errGetNodeValue, node.Key)
	}

	fmt.Fprintf(builder, "\n%s%s: %v", getIndent(indent), node.Key, fs.formatValue(oldValue, indent+1))

	return nil
}

func (fs FormatterStylish) formatValue(value any, indent int) any {
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

func (fs FormatterStylish) formatNestedMap(value map[string]any, indent int) string {
	builder := &strings.Builder{}

	fmt.Fprintf(builder, "{")

	for _, key := range getSortedKeys(value) {
		fmt.Fprintf(builder, "\n%s%s: %v", getIndent(indent), key, fs.formatValue(value[key], indent+1))
	}

	fmt.Fprintf(builder, "\n%s}", getIndent(indent-1))

	return builder.String()
}

func (fs FormatterStylish) formatNestedArray(value []any, indent int) string {
	builder := &strings.Builder{}

	fmt.Fprintf(builder, "[")

	for _, v := range value {
		fmt.Fprintf(builder, "\n%s%v", getIndent(indent), fs.formatValue(v, indent+1))
	}

	fmt.Fprintf(builder, "\n%s]", getIndent(indent-1))

	return builder.String()
}

func (fs FormatterStylish) getSortedChildKeys(node *analyzer.DiffNode) []string {
	keys := slices.Collect(maps.Keys(node.Children))
	slices.Sort(keys)

	return keys
}

func getSortedKeys(valueMap map[string]any) []string {
	keys := slices.Collect(maps.Keys(valueMap))
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
