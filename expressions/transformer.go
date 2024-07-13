package expressions

import (
	"flim/common"
	"fmt"
)

type TransformerExpression struct {
	name string
	expr common.Expression
}

func NewTransformerExpression(name string, expr common.Expression) (TransformerExpression, error) {
	return TransformerExpression{name: name, expr: expr}, nil
}

func (e TransformerExpression) ToString() string {
	return fmt.Sprintf("TransformerExpression<%s, %s>", e.name, e.expr.ToString())
}

func (e TransformerExpression) GetTags() map[string]common.Expression {
	return e.expr.GetTags()
}

func (e TransformerExpression) ReplaceReferences(tags map[string]common.Expression) (common.Expression, error) {
	newExpr, err := e.expr.ReplaceReferences(tags)

	if err != nil {
		return nil, err
	}

	e.expr = newExpr

	return e, nil
}

func (e TransformerExpression) Evaluate(handlers map[string]common.HandlerFunc) (interface{}, error) {
	handler, exists := handlers[e.name]

	if !exists {
		return nil, fmt.Errorf("no handler for `%s'", e.name)
	}

	exprResult, err := e.expr.Evaluate(handlers)

	if err != nil {
		return nil, err
	}

	handlerResult, err := handler(exprResult)

	if err != nil {
		return nil, err
	}

	return handlerResult, nil
}

func (e TransformerExpression) Serialize(config *common.SerializerConfig, indentLevel int) (string, error) {
	exprStr, err := e.expr.Serialize(config, indentLevel)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s", e.name, exprStr), nil
}

type MappedTransformerExpression struct {
	transformer string
	expr        common.Expression
}

func NewMappedTransformerExpression(transformer string, expr common.Expression) (MappedTransformerExpression, error) {
	return MappedTransformerExpression{transformer: transformer, expr: expr}, nil
}

func (e MappedTransformerExpression) ToString() string {
	return fmt.Sprintf("MappedTransformerExpression<%s, %s>", e.transformer, e.expr.ToString())
}

func (e MappedTransformerExpression) GetTags() map[string]common.Expression {
	return e.expr.GetTags()
}

func (e MappedTransformerExpression) ReplaceReferences(tags map[string]common.Expression) (common.Expression, error) {
	newExpr, err := e.expr.ReplaceReferences(tags)

	if err != nil {
		return nil, err
	}

	e.expr = newExpr

	return e, nil
}

func (e MappedTransformerExpression) Evaluate(handlers map[string]common.HandlerFunc) (interface{}, error) {
	listItemResults := []interface{}{}
	listExpr := e.expr.(ListExpression)

	for _, expr := range listExpr.listItems {
		transformedExpr := TransformerExpression{name: e.transformer, expr: expr}
		transformedResult, err := transformedExpr.Evaluate(handlers)

		if err != nil {
			return nil, err
		}

		listItemResults = append(listItemResults, transformedResult)
	}

	return listItemResults, nil
}

func (e MappedTransformerExpression) Serialize(config *common.SerializerConfig, indentLevel int) (string, error) {
	exprStr, err := e.expr.Serialize(config, indentLevel)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("@%s %s", e.transformer, exprStr), nil
}
