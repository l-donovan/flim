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

func (p *Parser) parseMapPairExpression() (Expression, error) {
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

	return PairExpression{key: left, val: right}, nil
}

func (p *Parser) parseMapExpression() (Expression, error) {
	pairs := []Expression{}

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
	listItems := []Expression{}

	for !p.peekToken().IsOfType("RightSquareBracket") {
		listItem, err := p.parseExpression()

		if err != nil {
			return nil, err
		}

		listItems = append(listItems, listItem)
	}

	// Throw away the right square bracket
	p.popToken()

	return ListExpression{listItems}, nil
}

func (p *Parser) parseExpression() (Expression, error) {
	token := p.popToken()

	if token.IsOfType("Star") {
		baseExpr, err := p.parseExpression()

		if err != nil {
			return nil, err
		}

		return ExpandingExpression{expr: baseExpr}, nil
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

		return TaggedExpression{tag: tagName, expr: baseExpr}, nil
	}

	if token.IsOfType("Boolean") {
		val := token.Contents == "true"
		return BooleanLiteralExpression{val}, nil
	}

	if token.IsOfType("String") {
		val := token.Contents[1 : len(token.Contents)-1]
		return StringLiteralExpression{val}, nil
	}

	if token.IsOfType("Float") {
		val, err := strconv.ParseFloat(token.Contents, 64)

		if err != nil {
			return nil, err
		}

		return FloatLiteralExpression{val}, nil
	}

	if token.IsOfType("Integer") {
		val, err := strconv.ParseInt(token.Contents, 10, 64)

		if err != nil {
			return nil, err
		}

		return IntegerLiteralExpression{val}, nil
	}

	if token.IsOfType("Null") {
		return NullLiteralExpression{}, nil
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

		return NamedExpression{name: token.Contents, expr: baseExpr}, nil
	}

	if token.IsOfType("Ampersand") {
		nameToken := p.popToken()

		if !nameToken.IsOfType("Keyword") {
			return nil, fmt.Errorf("tag references must be keywords")
		}

		return ReferenceExpression{nameToken.Contents}, nil
	}

	return nil, fmt.Errorf("unknown token type %s", token.Name)
}

func (p *Parser) parseFileExpression() (Expression, error) {
	expressions := []Expression{}

	for len(p.tokens) > 0 {
		expr, err := p.parseExpression()

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

	tags := expr.GetTags()
	newExpr, err := expr.ReplaceReferences(tags)

	if err != nil {
		return nil, err
	}

	return newExpr, nil
}
