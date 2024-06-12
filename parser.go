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

func (p *Parser) parseMapExpression() (Expression, error) {
	pairs := []PairExpression{}

	for !p.peekToken().IsOfType("RightCurlyBrace") {
		leftToken := p.popToken()

		if !leftToken.IsOfType("Keyword") {
			return nil, fmt.Errorf("NO GOOD")
		}

		left := leftToken.Contents

		right, err := p.parseExpression()

		if err != nil {
			return nil, err
		}

		pair := PairExpression{key: left, val: right}
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
		expr, err := p.parseExpression()

		if err != nil {
			return nil, err
		}

		return NamedExpression{name: token.Contents, expr: expr}, nil
	}

	return nil, fmt.Errorf("unknown token type %s", token.Name)
}

func (p *Parser) parse(tokens []LexerToken) (Expression, error) {
	p.tokens = tokens
	return p.parseExpression()
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
	expr, err := parser.parse(tokens)

	if err != nil {
		return nil, err
	}

	return expr, nil
}
