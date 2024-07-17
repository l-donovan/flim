package expressions

import (
	"github.com/l-donovan/flim/common"
	"fmt"
)

type TaggedExpression struct {
	tag  string
	expr common.Expression
}

func NewTaggedExpression(tag string, expr common.Expression) (TaggedExpression, error) {
	return TaggedExpression{tag: tag, expr: expr}, nil
}

func (e TaggedExpression) ToString() string {
	return fmt.Sprintf("TaggedExpression<#%s, %s>", e.tag, e.expr.ToString())
}

func (e TaggedExpression) GetTags() map[string]common.Expression {
	tags := e.expr.GetTags()
	tags[e.tag] = e.expr

	return tags
}

func (e TaggedExpression) ReplaceReferences(tags map[string]common.Expression) (common.Expression, error) {
	newExpr, err := e.expr.ReplaceReferences(tags)

	if err != nil {
		return nil, err
	}

	e.expr = newExpr

	return e, nil
}

func (e TaggedExpression) Evaluate(handlers map[string]common.HandlerFunc) (interface{}, error) {
	return e.expr.Evaluate(handlers)
}

func (e TaggedExpression) Serialize(config *common.SerializerConfig, indentLevel int) (string, error) {
	exprStr, err := e.expr.Serialize(config, indentLevel)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("#%s %s", e.tag, exprStr), nil
}
