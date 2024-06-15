package flim

import (
	"fmt"
	"os"
	"strconv"
)

type Parser struct {
	tokens []LexerToken
}

func (p *Parser) popToken() LexerToken {
	token := p.tokens[0]
	p.tokens = p.tokens[1:]
	return token
}

func (p Parser) peekToken() LexerToken {
	return p.tokens[0]
}

func (p *Parser) parseMapPairExpression() (ExpressionMeta, error) {
	if p.peekToken().IsOfType("Star") {
		expr, err := p.parseTaggedExpression()

		if err != nil {
			return undefinedExpressionMeta, err
		}

		return expr, nil
	}

	leftToken := p.popToken()

	if !leftToken.IsOfType("Keyword") {
		return undefinedExpressionMeta, fmt.Errorf("map pair cannot start with token of type %s", leftToken.Name)
	}

	left := leftToken.Contents

	right, err := p.parseTaggedExpression()

	if err != nil {
		return undefinedExpressionMeta, err
	}

	expr := PairExpression{key: left, val: right}
	pair := ExpressionMeta{expr, false, ""}

	return pair, nil
}

func (p *Parser) parseMapExpression() (Expression, error) {
	pairs := []ExpressionMeta{}

	for !p.peekToken().IsOfType("RightCurlyBrace") {
		pair, err := p.parseMapPairExpression()

		if err != nil {
			return nil, err
		}

		pairs = append(pairs, pair)
	}

	// Throw away the right curly brace
	p.popToken()

	return MapExpression{pairs}, nil
}

func (p *Parser) parseListExpression() (Expression, error) {
	listItems := []ExpressionMeta{}

	for !p.peekToken().IsOfType("RightSquareBracket") {
		listItem, err := p.parseTaggedExpression()

		if err != nil {
			return nil, err
		}

		listItems = append(listItems, listItem)
	}

	// Throw away the right square bracket
	p.popToken()

	return ListExpression{listItems}, nil
}

func (p *Parser) parseExpression() (ExpressionMeta, error) {
	token := p.popToken()
	shouldExpand := false

	if token.IsOfType("Star") {
		token = p.popToken()
		shouldExpand = true

		if !token.IsOfType("Keyword", "LeftSquareBracket", "LeftCurlyBrace", "Ampersand") {
			return undefinedExpressionMeta, fmt.Errorf("token of type %s cannot be expanded", token.Name)
		}
	}

	if token.IsOfType("Boolean") {
		val := token.Contents == "true"
		expr := BooleanLiteralExpression{val}
		return ExpressionMeta{expr, shouldExpand, ""}, nil
	}

	if token.IsOfType("String") {
		val := token.Contents[1 : len(token.Contents)-1]
		expr := StringLiteralExpression{val}
		return ExpressionMeta{expr, shouldExpand, ""}, nil
	}

	if token.IsOfType("Float") {
		val, err := strconv.ParseFloat(token.Contents, 64)

		if err != nil {
			return undefinedExpressionMeta, err
		}

		expr := FloatLiteralExpression{val}
		return ExpressionMeta{expr, shouldExpand, ""}, nil
	}

	if token.IsOfType("Integer") {
		val, err := strconv.ParseInt(token.Contents, 10, 64)

		if err != nil {
			return undefinedExpressionMeta, err
		}

		expr := IntegerLiteralExpression{val}
		return ExpressionMeta{expr, shouldExpand, ""}, nil
	}

	if token.IsOfType("Null") {
		expr := NullLiteralExpression{}
		return ExpressionMeta{expr, shouldExpand, ""}, nil
	}

	if token.IsOfType("LeftCurlyBrace") {
		expr, err := p.parseMapExpression()
		return ExpressionMeta{expr, shouldExpand, ""}, err
	}

	if token.IsOfType("LeftSquareBracket") {
		expr, err := p.parseListExpression()
		return ExpressionMeta{expr, shouldExpand, ""}, err
	}

	if token.IsOfType("Keyword") {
		baseExpr, err := p.parseTaggedExpression()

		if err != nil {
			return undefinedExpressionMeta, err
		}

		expr := NamedExpression{name: token.Contents, expr: baseExpr}
		return ExpressionMeta{expr, shouldExpand, ""}, nil
	}

	if token.IsOfType("Ampersand") {
		nameToken := p.popToken()

		if !nameToken.IsOfType("Keyword") {
			return undefinedExpressionMeta, fmt.Errorf("tag references must be keywords")
		}

		expr := ReferenceExpression{nameToken.Contents}
		return ExpressionMeta{expr, shouldExpand, ""}, nil
	}

	return undefinedExpressionMeta, fmt.Errorf("unknown token type %s", token.Name)
}

func (p *Parser) parseTaggedExpression() (ExpressionMeta, error) {
	tagName := ""

	if p.peekToken().IsOfType("Pound") {
		p.popToken()

		token := p.popToken()

		if !token.IsOfType("Keyword") {
			return undefinedExpressionMeta, fmt.Errorf("expected keyword")
		}

		tagName = token.Contents
	}

	exprMeta, err := p.parseExpression()

	if err != nil {
		return undefinedExpressionMeta, err
	}

	exprMeta.Tag = tagName

	return exprMeta, nil
}

func (p *Parser) parseFileExpression() (Expression, error) {
	expressions := []ExpressionMeta{}

	for len(p.tokens) > 0 {
		expr, err := p.parseTaggedExpression()

		if err != nil {
			return nil, err
		}

		expressions = append(expressions, expr)
	}

	return FileExpression{expressions}, nil
}

func (p *Parser) Parse(tokens []LexerToken) (Expression, error) {
	p.tokens = tokens

	fileExpr, err := p.parseFileExpression()

	if err != nil {
		return nil, err
	}

	return fileExpr, nil
}

func ParseFile(filename string) (Expression, error) {
	fileContents, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	tokens, err := Lex(string(fileContents))

	if err != nil {
		return nil, err
	}

	parser := Parser{}
	expr, err := parser.Parse(tokens)

	if err != nil {
		return nil, err
	}

	return expr, nil
}
