package flim

import (
	"flim/common"
	flimexpr "flim/expressions"
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

func (p *Parser) peekToken() LexerToken {
	return p.tokens[0]
}

func (p *Parser) parseMapPairExpression() (common.Expression, error) {
	if p.peekToken().IsOfType("Star") {
		expr, err := p.parseExpression()

		if err != nil {
			return nil, err
		}

		return expr, nil
	}

	leftToken := p.popToken()

	if !leftToken.IsOfType("Keyword") {
		return nil, fmt.Errorf("map pair cannot start with token of type %s", leftToken.Name)
	}

	left := leftToken.Contents

	right, err := p.parseExpression()

	if err != nil {
		return nil, err
	}

	return flimexpr.NewPairExpression(left, right)
}

func (p *Parser) parseMapExpression() (common.Expression, error) {
	pairs := []common.Expression{}

	for !p.peekToken().IsOfType("RightCurlyBrace") {
		pair, err := p.parseMapPairExpression()

		if err != nil {
			return nil, err
		}

		pairs = append(pairs, pair)
	}

	// Throw away the right curly brace
	p.popToken()

	return flimexpr.NewMapExpression(pairs)
}

func (p *Parser) parseListExpression() (common.Expression, error) {
	listItems := []common.Expression{}

	for !p.peekToken().IsOfType("RightSquareBracket") {
		listItem, err := p.parseExpression()

		if err != nil {
			return nil, err
		}

		listItems = append(listItems, listItem)
	}

	// Throw away the right square bracket
	p.popToken()

	return flimexpr.NewListExpression(listItems)
}

func (p *Parser) parseExpression() (common.Expression, error) {
	token := p.popToken()

	if token.IsOfType("Star") {
		baseExpr, err := p.parseExpression()

		if err != nil {
			return nil, err
		}

		return flimexpr.NewExpandingExpression(baseExpr)
	}

	if token.IsOfType("Pound") {
		token := p.popToken()

		if !token.IsOfType("Keyword") {
			return nil, fmt.Errorf("expected keyword")
		}

		tagName := token.Contents
		baseExpr, err := p.parseExpression()

		if err != nil {
			return nil, err
		}

		return flimexpr.NewTaggedExpression(tagName, baseExpr)
	}

	if token.IsOfType("AtSign") {
		token := p.popToken()

		if !token.IsOfType("Keyword") {
			return nil, fmt.Errorf("expected keyword")
		}

		transformerName := token.Contents
		baseExpr, err := p.parseExpression()

		if err != nil {
			return nil, err
		}

		return flimexpr.NewMappedTransformerExpression(transformerName, baseExpr)
	}

	if token.IsOfType("Boolean") {
		val := token.Contents == "true"
		return flimexpr.NewBooleanLiteralExpression(val)
	}

	if token.IsOfType("String") {
		val := token.Contents[1 : len(token.Contents)-1]
		return flimexpr.NewStringLiteralExpression(val)
	}

	if token.IsOfType("Float") {
		val, err := strconv.ParseFloat(token.Contents, 64)

		if err != nil {
			return nil, err
		}

		return flimexpr.NewFloatLiteralExpression(val)
	}

	if token.IsOfType("Integer") {
		val, err := strconv.ParseInt(token.Contents, 10, 64)

		if err != nil {
			return nil, err
		}

		return flimexpr.NewIntegerLiteralExpression(val)
	}

	if token.IsOfType("Null") {
		return flimexpr.NewNullLiteralExpression()
	}

	if token.IsOfType("LeftCurlyBrace") {
		return p.parseMapExpression()
	}

	if token.IsOfType("LeftSquareBracket") {
		return p.parseListExpression()
	}

	if token.IsOfType("Keyword") {
		baseExpr, err := p.parseExpression()

		if err != nil {
			return nil, err
		}

		return flimexpr.NewTransformerExpression(token.Contents, baseExpr)
	}

	if token.IsOfType("Ampersand") {
		nameToken := p.popToken()

		if !nameToken.IsOfType("Keyword") {
			return nil, fmt.Errorf("tag references must be keywords")
		}

		return flimexpr.NewReferenceExpression(nameToken.Contents)
	}

	return nil, fmt.Errorf("unknown token type %s", token.Name)
}

func (p *Parser) parseFileExpression() (common.Expression, error) {
	expressions := []common.Expression{}

	for len(p.tokens) > 0 {
		expr, err := p.parseExpression()

		if err != nil {
			return nil, err
		}

		expressions = append(expressions, expr)
	}

	return flimexpr.NewFileExpression(expressions)
}

func (p *Parser) Parse(tokens []LexerToken) (common.Expression, error) {
	p.tokens = tokens

	fileExpr, err := p.parseFileExpression()

	if err != nil {
		return nil, err
	}

	return fileExpr, nil
}

func ParseFile(filename string) (common.Expression, error) {
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

	tags := expr.GetTags()
	newExpr, err := expr.ReplaceReferences(tags)

	if err != nil {
		return nil, err
	}

	return newExpr, nil
}
