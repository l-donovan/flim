package expressions

import (
	"github.com/l-donovan/flim/common"
	"fmt"
)

type ExpandingExpression struct {
	expr common.Expression
}

func NewExpandingExpression(expr common.Expression) (ExpandingExpression, error) {
	return ExpandingExpression{expr: expr}, nil
}

func (e ExpandingExpression) ToString() string {
	return fmt.Sprintf("ExpandingExpression<%s>", e.expr.ToString())
}

func (e ExpandingExpression) GetTags() map[string]common.Expression {
	return e.expr.GetTags()
}

func (e ExpandingExpression) ReplaceReferences(tags map[string]common.Expression) (common.Expression, error) {
	newExpr, err := e.expr.ReplaceReferences(tags)

	if err != nil {
		return nil, err
	}

	e.expr = newExpr

	return e, nil
}

func (e ExpandingExpression) Evaluate(handlers map[string]common.HandlerFunc) (interface{}, error) {
	return e.expr.Evaluate(handlers)
}

func (e ExpandingExpression) Serialize(config *common.SerializerConfig, indentLevel int) (string, error) {
	exprStr, err := e.expr.Serialize(config, indentLevel)

	if err != nil {
		return "", err
	}

	return "*" + exprStr, nil
}
