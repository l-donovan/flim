package expressions

import (
	"flim/common"
	"fmt"
	"strings"
)

type Pair struct {
	Key string
	Val interface{}
}

type PairExpression struct {
	key string
	val common.Expression
}

func NewPairExpression(key string, val common.Expression) (PairExpression, error) {
	return PairExpression{key: key, val: val}, nil
}

func (e PairExpression) ToString() string {
	return fmt.Sprintf("PairExpression<%s: %s>", e.key, e.val.ToString())
}

func (e PairExpression) GetTags() map[string]common.Expression {
	return e.val.GetTags()
}

func (e PairExpression) ReplaceReferences(tags map[string]common.Expression) (common.Expression, error) {
	newExpr, err := e.val.ReplaceReferences(tags)

	if err != nil {
		return nil, err
	}

	e.val = newExpr

	return e, nil
}

func (e PairExpression) Evaluate(handlers map[string]common.HandlerFunc) (interface{}, error) {
	result, err := e.val.Evaluate(handlers)

	if err != nil {
		return nil, err
	}

	return Pair{Key: e.key, Val: result}, nil
}

func (e PairExpression) Serialize(config *common.SerializerConfig, indentLevel int) (string, error) {
	exprStr, err := e.val.Serialize(config, indentLevel)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s", e.key, exprStr), nil
}

type MapExpression struct {
	pairs []common.Expression
}

func NewMapExpression(pairs []common.Expression) (MapExpression, error) {
	return MapExpression{pairs: pairs}, nil
}

func (e MapExpression) ToString() string {
	pairStrings := []string{}

	for _, pair := range e.pairs {
		pairStrings = append(pairStrings, pair.ToString())
	}

	return fmt.Sprintf("MapExpression<%s>", strings.Join(pairStrings, ", "))
}

func (e MapExpression) GetTags() map[string]common.Expression {
	tags := map[string]common.Expression{}

	for _, pairExpr := range e.pairs {
		pairTags := pairExpr.GetTags()

		for key, val := range pairTags {
			tags[key] = val
		}
	}

	return tags
}

func (e MapExpression) ReplaceReferences(tags map[string]common.Expression) (common.Expression, error) {
	for i, pairExpr := range e.pairs {
		newExpr, err := pairExpr.ReplaceReferences(tags)

		if err != nil {
			return nil, err
		}

		e.pairs[i] = newExpr
	}

	return e, nil
}

func (e MapExpression) Evaluate(handlers map[string]common.HandlerFunc) (interface{}, error) {
	pairResults := map[string]interface{}{}

	for _, pairExpr := range e.pairs {
		pairResult, err := pairExpr.Evaluate(handlers)

		if err != nil {
			return nil, err
		}

		if _, ok := pairExpr.(ExpandingExpression); ok {
			pairExprExpanded, ok := pairResult.(map[string]interface{})

			if !ok {
				return nil, fmt.Errorf("could not expand map pair")
			}

			for key, val := range pairExprExpanded {
				pairResults[key] = val
			}
		} else {
			pair := pairResult.(Pair)
			pairResults[pair.Key] = pair.Val
		}
	}

	return pairResults, nil
}

func (e MapExpression) Serialize(config *common.SerializerConfig, indentLevel int) (string, error) {
	exprStrs := make([]string, len(e.pairs))

	for i, pair := range e.pairs {
		pairStr, err := pair.Serialize(config, indentLevel+1)

		if err != nil {
			return "", err
		}

		exprStrs[i] = config.Indent(indentLevel) + pairStr
	}

	if len(e.pairs) == 0 {
		return "{}", nil
	} else {
		out := fmt.Sprintf(
			"{%s%s%s%s}",
			config.Sep("\n", ""),
			strings.Join(exprStrs, config.Sep("\n", " ")),
			config.Sep("\n", ""),
			config.Indent(indentLevel-1),
		)
		return out, nil
	}
}
