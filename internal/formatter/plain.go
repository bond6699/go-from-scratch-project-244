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

	for _, key := range fp.getSortedChildKeys(node) {
		if err := fp.formatNode(&builder, []string{}, node.Children[key]); err != nil {
			return "", err
		}
	}

	return strings.TrimSuffix(builder.String(), "\n"), nil
}

func (fp FormatterPlain) formatNode(builder *strings.Builder, path []string, node *analyzer.DiffNode) error {
	path = append(path, node.Key)
	hasChild := node.Children != nil

	switch node.Status {
	case analyzer.Container:
		return fp.formatContainer(builder, path, node)
	case analyzer.Added:
		return fp.formatAdded(builder, path, node, hasChild)
	case analyzer.Deleted:
		return fp.formatDeleted(builder, path, node)
	case analyzer.Changed:
		return fp.formatChanged(builder, path, node, hasChild)
	}

	return nil
}

func (fp FormatterPlain) formatContainer(builder *strings.Builder, path []string, node *analyzer.DiffNode) error {
	for _, key := range fp.getSortedChildKeys(node) {
		if err := fp.formatNode(builder, path, node.Children[key]); err != nil {
			return err
		}
	}

	return nil
}

func (fp FormatterPlain) formatAdded(
	builder *strings.Builder,
	path []string,
	node *analyzer.DiffNode,
	hasChild bool,
) error {
	ok, value := node.GetNewValue()
	if !ok {
		return fmt.Errorf("%w: key %q", errGetNodeValue, node.Key)
	}

	fmt.Fprintf(
		builder, "Property '%s' was added with value: %v\n",
		strings.Join(path, "."),
		fp.formatValue(value, hasChild),
	)

	return nil
}

func (fp FormatterPlain) formatDeleted(builder *strings.Builder, path []string, node *analyzer.DiffNode) error {
	if ok, _ := node.GetOldValue(); !ok {
		return fmt.Errorf("%w: key %q", errGetNodeValue, node.Key)
	}

	fmt.Fprintf(builder, "Property '%s' was removed\n", strings.Join(path, "."))

	return nil
}

func (fp FormatterPlain) formatChanged(
	builder *strings.Builder,
	path []string,
	node *analyzer.DiffNode,
	hasChild bool,
) error {
	okOld, oldValue := node.GetOldValue()
	if !okOld {
		return fmt.Errorf("%w: key %q", errGetNodeValue, node.Key)
	}

	okNew, newValue := node.GetNewValue()
	if !okNew {
		return fmt.Errorf("%w: key %q", errGetNodeValue, node.Key)
	}

	fmt.Fprintf(
		builder, "Property '%s' was updated. From %v to %v\n",
		strings.Join(path, "."),
		fp.formatValue(oldValue, hasChild),
		fp.formatValue(newValue, hasChild),
	)

	return nil
}

func (fp FormatterPlain) formatValue(value any, hasChild bool) any {
	if value == nil {
		return "null"
	}

	_, isMap := value.(map[string]any)
	_, isArray := value.([]any)

	if hasChild || isMap || isArray {
		return "[complex value]"
	}

	if v, ok := value.(string); ok {
		return fmt.Sprintf("'%s'", v)
	}

	if v, isString := value.(string); isString {
		return strings.ReplaceAll(v, "\n", "\\n")
	}

	return value
}

func (fp FormatterPlain) getSortedChildKeys(node *analyzer.DiffNode) []string {
	keys := slices.Collect(maps.Keys(node.Children))
	slices.Sort(keys)

	return keys
}
