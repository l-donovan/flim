package expressions

import (
	"github.com/l-donovan/flim/common"
	"fmt"
	"strings"
)

type FileExpression struct {
	expressions []common.Expression
}

func NewFileExpression(expressions []common.Expression) (FileExpression, error) {
	return FileExpression{expressions: expressions}, nil
}

func (e FileExpression) ToString() string {
	expressionStrings := []string{}

	for _, expr := range e.expressions {
		expressionStrings = append(expressionStrings, expr.ToString())
	}

	return fmt.Sprintf("FileExpression<%s>", strings.Join(expressionStrings, ", "))
}

func (e FileExpression) GetTags() map[string]common.Expression {
	tags := map[string]common.Expression{}

	for _, expr := range e.expressions {
		childTags := expr.GetTags()

		for key, val := range childTags {
			tags[key] = val
		}
	}

	return tags
}

func (e FileExpression) ReplaceReferences(tags map[string]common.Expression) (common.Expression, error) {
	for i, expr := range e.expressions {
		newExpr, err := expr.ReplaceReferences(tags)

		if err != nil {
			return nil, err
		}

		e.expressions[i] = newExpr
	}

	return e, nil
}

func (e FileExpression) Evaluate(handlers map[string]common.HandlerFunc) (interface{}, error) {
	lastListItemResult := interface{}(nil)

	for _, expr := range e.expressions {
		listItemResult, err := expr.Evaluate(handlers)

		if err != nil {
			return nil, err
		}

		if _, ok := expr.(ExpandingExpression); ok {
			listItemExpanded, ok := listItemResult.([]interface{})

			if !ok {
				return nil, fmt.Errorf("could not expand list item")
			}

			for _, listItemSingleResult := range listItemExpanded {
				lastListItemResult = listItemSingleResult
			}
		} else {
			lastListItemResult = listItemResult
		}
	}

	return lastListItemResult, nil
}

func (e FileExpression) Serialize(config *common.SerializerConfig, indentLevel int) (string, error) {
	exprStrs := make([]string, len(e.expressions))

	for i, expr := range e.expressions {
		exprStr, err := expr.Serialize(config, indentLevel+1)

		if err != nil {
			return "", err
		}

		exprStrs[i] = config.Indent(indentLevel) + exprStr
	}

	return strings.Join(exprStrs, config.Sep("\n\n", " ")), nil
}
