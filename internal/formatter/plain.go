package formatter

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"code/internal/analyzer"
)

type Plain struct{}

func NewFormatterPlain() Plain {
	return Plain{}
}

func (fp Plain) Format(node *analyzer.ContainerNode) (string, error) {
	var builder strings.Builder

	for _, key := range fp.getSortedChildKeys(node) {
		if err := fp.formatNode(&builder, []string{}, node.Children()[key]); err != nil {
			return "", err
		}
	}

	return strings.TrimSuffix(builder.String(), "\n"), nil
}

func (fp Plain) formatNode(builder *strings.Builder, path []string, node analyzer.DiffNode) error {
	path = append(path, node.Key())

	switch v := node.(type) {
	case *analyzer.ContainerNode:
		return fp.formatContainer(builder, path, v)
	case *analyzer.AddedNode:
		return fp.formatAdded(builder, path, v)
	case *analyzer.DeletedNode:
		return fp.formatDeleted(builder, path)
	case *analyzer.ChangedNode:
		return fp.formatChanged(builder, path, v)
	}

	return nil
}

func (fp Plain) formatContainer(builder *strings.Builder, path []string, node *analyzer.ContainerNode) error {
	for _, key := range fp.getSortedChildKeys(node) {
		if err := fp.formatNode(builder, path, node.Children()[key]); err != nil {
			return err
		}
	}

	return nil
}

func (fp Plain) formatAdded(
	builder *strings.Builder,
	path []string,
	node *analyzer.AddedNode,
) error {
	fmt.Fprintf(
		builder, "Property '%s' was added with value: %v\n",
		strings.Join(path, "."),
		fp.formatValue(node.Value()),
	)

	return nil
}

func (fp Plain) formatDeleted(builder *strings.Builder, path []string) error {
	fmt.Fprintf(builder, "Property '%s' was removed\n", strings.Join(path, "."))

	return nil
}

func (fp Plain) formatChanged(
	builder *strings.Builder,
	path []string,
	node *analyzer.ChangedNode,
) error {
	fmt.Fprintf(
		builder, "Property '%s' was updated. From %v to %v\n",
		strings.Join(path, "."),
		fp.formatValue(node.OldValue()),
		fp.formatValue(node.NewValue()),
	)

	return nil
}

func (fp Plain) formatValue(value any) any {
	if value == nil {
		return "null"
	}

	_, isMap := value.(map[string]any)
	_, isArray := value.([]any)

	if isMap || isArray {
		return "[complex value]"
	}

	if v, ok := value.(string); ok {
		escaped := strings.ReplaceAll(v, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "\n", "\\n")

		return fmt.Sprintf("'%s'", escaped)
	}

	return value
}

func (fp Plain) getSortedChildKeys(node *analyzer.ContainerNode) []string {
	keys := slices.Collect(maps.Keys(node.Children()))
	slices.Sort(keys)

	return keys
}
