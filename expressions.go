package flim

import (
	"fmt"
	"strings"
)

type HandlerFunc func(data interface{}) (interface{}, error)

type Expression interface {
	ToString() string
	Evaluate(map[string]HandlerFunc) (interface{}, error)
}

type IntegerLiteralExpression struct {
	val int64
}

func (e IntegerLiteralExpression) ToString() string {
	return fmt.Sprintf("IntegerLiteralExpression<%d>", e.val)
}

func (e IntegerLiteralExpression) Evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
	return e.val, nil
}

type FloatLiteralExpression struct {
	val float64
}

func (e FloatLiteralExpression) ToString() string {
	return fmt.Sprintf("FloatLiteralExpression<%f>", e.val)
}

func (e FloatLiteralExpression) Evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
	return e.val, nil
}

type BooleanLiteralExpression struct {
	val bool
}

func (e BooleanLiteralExpression) ToString() string {
	return fmt.Sprintf("BooleanLiteralExpression<%t>", e.val)
}

func (e BooleanLiteralExpression) Evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
	return e.val, nil
}

type StringLiteralExpression struct {
	val string
}

func (e StringLiteralExpression) ToString() string {
	return fmt.Sprintf("StringLiteralExpression<%s>", e.val)
}

func (e StringLiteralExpression) Evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
	return e.val, nil
}

type NullLiteralExpression struct {
}

func (e NullLiteralExpression) ToString() string {
	return "NullLiteralExpression<>"
}

func (e NullLiteralExpression) Evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
	return nil, nil
}

type ListExpression struct {
	listItems []Expression
}

func (e ListExpression) ToString() string {
	listItemStrings := []string{}

	for _, listItem := range e.listItems {
		listItemStrings = append(listItemStrings, listItem.ToString())
	}

	return fmt.Sprintf("ListExpression<%s>", strings.Join(listItemStrings, ", "))
}

func (e ListExpression) Evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
	listItemResults := []interface{}{}

	for _, listItem := range e.listItems {
		listItemResult, err := listItem.Evaluate(handlers)

		if err != nil {
			return nil, err
		}

		listItemResults = append(listItemResults, listItemResult)
	}

	return listItemResults, nil
}

type Pair struct {
	Key string
	Val interface{}
}

type PairExpression struct {
	key string
	val Expression
}

func (e PairExpression) toString() string {
	return fmt.Sprintf("PairExpression<%s: %s>", e.key, e.val.ToString())
}

func (e PairExpression) evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
	result, err := e.val.Evaluate(handlers)

	if err != nil {
		return nil, err
	}

	return Pair{Key: e.key, Val: result}, nil
}

type MapExpression struct {
	pairs []PairExpression
}

func (e MapExpression) ToString() string {
	pairStrings := []string{}

	for _, pair := range e.pairs {
		pairStrings = append(pairStrings, pair.toString())
	}

	return fmt.Sprintf("MapExpression<%s>", strings.Join(pairStrings, ", "))
}

func (e MapExpression) Evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
	pairResults := map[string]interface{}{}

	for _, pairExpr := range e.pairs {
		pairResult, err := pairExpr.evaluate(handlers)

		if err != nil {
			return nil, err
		}

		pair := pairResult.(Pair)
		pairResults[pair.Key] = pair.Val
	}

	return pairResults, nil
}

type NamedExpression struct {
	name string
	expr Expression
}

func (e NamedExpression) ToString() string {
	return fmt.Sprintf("NamedExpression<%s, %s>", e.name, e.expr.ToString())
}

func (e NamedExpression) Evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
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
