package flim

import (
	"fmt"
	"strings"
)

type SerializerConfig struct {
	useTabs    bool
	indentSize int
	minify     bool
}

func (c SerializerConfig) Indent(indentLevel int) string {
	if c.minify {
		return ""
	}

	if c.useTabs {
		return strings.Repeat("\t", c.indentSize*indentLevel)
	}

	return strings.Repeat(" ", c.indentSize*indentLevel)
}

func (c SerializerConfig) Sep(separator string, alt string) string {
	if c.minify {
		return alt
	}

	return separator
}

func Serialize(expr Expression, useTabs bool, indentSize int) (string, error) {
	config := SerializerConfig{useTabs: useTabs, indentSize: indentSize}
	return expr.Serialize(&config, 0)
}

func Minify(expr Expression) (string, error) {
	config := SerializerConfig{minify: true}
	return expr.Serialize(&config, 0)
}

type HandlerFunc func(data interface{}) (interface{}, error)

type Expression interface {
	ToString() string
	GetTags() map[string]Expression
	ReplaceReferences(map[string]Expression) (Expression, error)
	Evaluate(map[string]HandlerFunc) (interface{}, error)
	Serialize(config *SerializerConfig, indentLevel int) (string, error)
}

type ExpandingExpression struct {
	expr Expression
}

func (e ExpandingExpression) ToString() string {
	return fmt.Sprintf("*%s", e.expr.ToString())
}

func (e ExpandingExpression) GetTags() map[string]Expression {
	return e.expr.GetTags()
}

func (e ExpandingExpression) ReplaceReferences(tags map[string]Expression) (Expression, error) {
	newExpr, err := e.expr.ReplaceReferences(tags)

	if err != nil {
		return nil, err
	}

	e.expr = newExpr

	return e, nil
}

func (e ExpandingExpression) Evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
	return e.expr.Evaluate(handlers)
}

func (e ExpandingExpression) Serialize(config *SerializerConfig, indentLevel int) (string, error) {
	exprStr, err := e.expr.Serialize(config, indentLevel)

	if err != nil {
		return "", err
	}

	return "*" + exprStr, nil
}

type TaggedExpression struct {
	tag  string
	expr Expression
}

func (e TaggedExpression) ToString() string {
	return fmt.Sprintf("#%s %s", e.tag, e.expr.ToString())
}

func (e TaggedExpression) GetTags() map[string]Expression {
	tags := e.expr.GetTags()
	tags[e.tag] = e.expr

	return tags
}

func (e TaggedExpression) ReplaceReferences(tags map[string]Expression) (Expression, error) {
	newExpr, err := e.expr.ReplaceReferences(tags)

	if err != nil {
		return nil, err
	}

	e.expr = newExpr

	return e, nil
}

func (e TaggedExpression) Evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
	return e.expr.Evaluate(handlers)
}

func (e TaggedExpression) Serialize(config *SerializerConfig, indentLevel int) (string, error) {
	exprStr, err := e.expr.Serialize(config, indentLevel)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("#%s %s", e.tag, exprStr), nil
}

type IntegerLiteralExpression struct {
	val int64
}

func (e IntegerLiteralExpression) ToString() string {
	return fmt.Sprintf("IntegerLiteralExpression<%d>", e.val)
}

func (e IntegerLiteralExpression) GetTags() map[string]Expression {
	return map[string]Expression{}
}

func (e IntegerLiteralExpression) ReplaceReferences(tags map[string]Expression) (Expression, error) {
	return e, nil
}

func (e IntegerLiteralExpression) Evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
	return e.val, nil
}

func (e IntegerLiteralExpression) Serialize(config *SerializerConfig, indentLevel int) (string, error) {
	return fmt.Sprintf("%d", e.val), nil
}

type FloatLiteralExpression struct {
	val float64
}

func (e FloatLiteralExpression) ToString() string {
	return fmt.Sprintf("FloatLiteralExpression<%f>", e.val)
}

func (e FloatLiteralExpression) GetTags() map[string]Expression {
	return map[string]Expression{}
}

func (e FloatLiteralExpression) ReplaceReferences(tags map[string]Expression) (Expression, error) {
	return e, nil
}

func (e FloatLiteralExpression) Evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
	return e.val, nil
}

func (e FloatLiteralExpression) Serialize(config *SerializerConfig, indentLevel int) (string, error) {
	return fmt.Sprintf("%f", e.val), nil
}

type BooleanLiteralExpression struct {
	val bool
}

func (e BooleanLiteralExpression) ToString() string {
	return fmt.Sprintf("BooleanLiteralExpression<%t>", e.val)
}

func (e BooleanLiteralExpression) ReplaceReferences(tags map[string]Expression) (Expression, error) {
	return e, nil
}

func (e BooleanLiteralExpression) GetTags() map[string]Expression {
	return map[string]Expression{}
}

func (e BooleanLiteralExpression) Evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
	return e.val, nil
}

func (e BooleanLiteralExpression) Serialize(config *SerializerConfig, indentLevel int) (string, error) {
	return fmt.Sprintf("%t", e.val), nil
}

type StringLiteralExpression struct {
	val string
}

func (e StringLiteralExpression) ToString() string {
	return fmt.Sprintf("StringLiteralExpression<%s>", e.val)
}

func (e StringLiteralExpression) ReplaceReferences(tags map[string]Expression) (Expression, error) {
	return e, nil
}

func (e StringLiteralExpression) GetTags() map[string]Expression {
	return map[string]Expression{}
}

func (e StringLiteralExpression) Evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
	return e.val, nil
}

func (e StringLiteralExpression) Serialize(config *SerializerConfig, indentLevel int) (string, error) {
	// TODO Escape more things
	replacer := strings.NewReplacer(
		"\n", "\\n",
	)
	escapedVal := replacer.Replace(e.val)
	return fmt.Sprintf("\"%s\"", escapedVal), nil
}

type NullLiteralExpression struct {
}

func (e NullLiteralExpression) ToString() string {
	return "NullLiteralExpression<>"
}

func (e NullLiteralExpression) ReplaceReferences(tags map[string]Expression) (Expression, error) {
	return e, nil
}

func (e NullLiteralExpression) GetTags() map[string]Expression {
	return map[string]Expression{}
}

func (e NullLiteralExpression) Evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
	return nil, nil
}

func (e NullLiteralExpression) Serialize(config *SerializerConfig, indentLevel int) (string, error) {
	return "null", nil
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

func (e ListExpression) GetTags() map[string]Expression {
	tags := map[string]Expression{}

	for _, listItem := range e.listItems {
		childTags := listItem.GetTags()

		for key, val := range childTags {
			tags[key] = val
		}
	}

	return tags
}

func (e ListExpression) ReplaceReferences(tags map[string]Expression) (Expression, error) {
	for i, listItem := range e.listItems {
		newExpr, err := listItem.ReplaceReferences(tags)

		if err != nil {
			return nil, err
		}

		e.listItems[i] = newExpr
	}

	return e, nil
}

func (e ListExpression) Evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
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

func (e ListExpression) Serialize(config *SerializerConfig, indentLevel int) (string, error) {
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

type Pair struct {
	Key string
	Val interface{}
}

type PairExpression struct {
	key string
	val Expression
}

func (e PairExpression) ToString() string {
	return fmt.Sprintf("PairExpression<%s: %s>", e.key, e.val.ToString())
}

func (e PairExpression) GetTags() map[string]Expression {
	return e.val.GetTags()
}

func (e PairExpression) ReplaceReferences(tags map[string]Expression) (Expression, error) {
	newExpr, err := e.val.ReplaceReferences(tags)

	if err != nil {
		return nil, err
	}

	e.val = newExpr

	return e, nil
}

func (e PairExpression) Evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
	result, err := e.val.Evaluate(handlers)

	if err != nil {
		return nil, err
	}

	return Pair{Key: e.key, Val: result}, nil
}

func (e PairExpression) Serialize(config *SerializerConfig, indentLevel int) (string, error) {
	exprStr, err := e.val.Serialize(config, indentLevel)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s", e.key, exprStr), nil
}

type MapExpression struct {
	pairs []Expression
}

func (e MapExpression) ToString() string {
	pairStrings := []string{}

	for _, pair := range e.pairs {
		pairStrings = append(pairStrings, pair.ToString())
	}

	return fmt.Sprintf("MapExpression<%s>", strings.Join(pairStrings, ", "))
}

func (e MapExpression) GetTags() map[string]Expression {
	tags := map[string]Expression{}

	for _, pairExpr := range e.pairs {
		pairTags := pairExpr.GetTags()

		for key, val := range pairTags {
			tags[key] = val
		}
	}

	return tags
}

func (e MapExpression) ReplaceReferences(tags map[string]Expression) (Expression, error) {
	for i, pairExpr := range e.pairs {
		newExpr, err := pairExpr.ReplaceReferences(tags)

		if err != nil {
			return nil, err
		}

		e.pairs[i] = newExpr
	}

	return e, nil
}

func (e MapExpression) Evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
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

func (e MapExpression) Serialize(config *SerializerConfig, indentLevel int) (string, error) {
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

type NamedExpression struct {
	name string
	expr Expression
}

func (e NamedExpression) ToString() string {
	return fmt.Sprintf("NamedExpression<%s, %s>", e.name, e.expr.ToString())
}

func (e NamedExpression) GetTags() map[string]Expression {
	return e.expr.GetTags()
}

func (e NamedExpression) ReplaceReferences(tags map[string]Expression) (Expression, error) {
	newExpr, err := e.expr.ReplaceReferences(tags)

	if err != nil {
		return nil, err
	}

	e.expr = newExpr

	return e, nil
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

func (e NamedExpression) Serialize(config *SerializerConfig, indentLevel int) (string, error) {
	exprStr, err := e.expr.Serialize(config, indentLevel)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s", e.name, exprStr), nil
}

type ReferenceExpression struct {
	name string
}

func (e ReferenceExpression) ToString() string {
	return fmt.Sprintf("ReferenceExpression<%s>", e.name)
}

func (e ReferenceExpression) GetTags() map[string]Expression {
	return map[string]Expression{}
}

func (e ReferenceExpression) ReplaceReferences(tags map[string]Expression) (Expression, error) {
	if replacement, exists := tags[e.name]; exists {
		return replacement, nil
	} else {
		return nil, fmt.Errorf("could not find tag `%s'", e.name)
	}
}

func (e ReferenceExpression) Evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
	return nil, fmt.Errorf("attempted to Evaluate a ReferenceExpression (hint: call ReplaceReferences first)")
}

func (e ReferenceExpression) Serialize(config *SerializerConfig, indentLevel int) (string, error) {
	return fmt.Sprintf("&%s", e.name), nil
}

type FileExpression struct {
	expressions []Expression
}

func (e FileExpression) ToString() string {
	expressionStrings := []string{}

	for _, expr := range e.expressions {
		expressionStrings = append(expressionStrings, expr.ToString())
	}

	return fmt.Sprintf("FileExpression<%s>", strings.Join(expressionStrings, ", "))
}

func (e FileExpression) GetTags() map[string]Expression {
	tags := map[string]Expression{}

	for _, expr := range e.expressions {
		childTags := expr.GetTags()

		for key, val := range childTags {
			tags[key] = val
		}
	}

	return tags
}

func (e FileExpression) ReplaceReferences(tags map[string]Expression) (Expression, error) {
	for i, expr := range e.expressions {
		newExpr, err := expr.ReplaceReferences(tags)

		if err != nil {
			return nil, err
		}

		e.expressions[i] = newExpr
	}

	return e, nil
}

func (e FileExpression) Evaluate(handlers map[string]HandlerFunc) (interface{}, error) {
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

func (e FileExpression) Serialize(config *SerializerConfig, indentLevel int) (string, error) {
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
