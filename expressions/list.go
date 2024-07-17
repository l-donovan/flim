package expressions

import (
	"github.com/l-donovan/flim/common"
	"fmt"
	"strings"
)

type ListExpression struct {
	listItems []common.Expression
}

func NewListExpression(listItems []common.Expression) (ListExpression, error) {
	return ListExpression{listItems: listItems}, nil
}

func (e ListExpression) ToString() string {
	listItemStrings := []string{}

	for _, listItem := range e.listItems {
		listItemStrings = append(listItemStrings, listItem.ToString())
	}

	return fmt.Sprintf("ListExpression<%s>", strings.Join(listItemStrings, ", "))
}

func (e ListExpression) GetTags() map[string]common.Expression {
	tags := map[string]common.Expression{}

	for _, listItem := range e.listItems {
		childTags := listItem.GetTags()

		for key, val := range childTags {
			tags[key] = val
		}
	}

	return tags
}

func (e ListExpression) ReplaceReferences(tags map[string]common.Expression) (common.Expression, error) {
	for i, listItem := range e.listItems {
		newExpr, err := listItem.ReplaceReferences(tags)

		if err != nil {
			return nil, err
		}

		e.listItems[i] = newExpr
	}

	return e, nil
}

func (e ListExpression) Evaluate(handlers map[string]common.HandlerFunc) (interface{}, error) {
	listItemResults := []interface{}{}

	for _, listItem := range e.listItems {
		listItemResult, err := listItem.Evaluate(handlers)

		if err != nil {
			return nil, err
		}

		if _, ok := listItem.(ExpandingExpression); ok {
			listItemExpanded, ok := listItemResult.([]interface{})

			if !ok {
				return nil, fmt.Errorf("could not expand list item")
			}

			listItemResults = append(listItemResults, listItemExpanded...)
		} else {
			listItemResults = append(listItemResults, listItemResult)
		}
	}

	return listItemResults, nil
}

func (e ListExpression) Serialize(config *common.SerializerConfig, indentLevel int) (string, error) {
	exprStrs := make([]string, len(e.listItems))

	for i, listItem := range e.listItems {
		listItemStr, err := listItem.Serialize(config, indentLevel+1)

		if err != nil {
			return "", err
		}

		exprStrs[i] = config.Indent(indentLevel) + listItemStr
	}

	if len(e.listItems) == 0 {
		return "[]", nil
	} else {
		out := fmt.Sprintf(
			"[%s%s%s%s]",
			config.Sep("\n", ""),
			strings.Join(exprStrs, config.Sep("\n", " ")),
			config.Sep("\n", ""),
			config.Indent(indentLevel-1),
		)
		return out, nil
	}
}
