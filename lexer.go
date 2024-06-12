package flim

import (
	"fmt"
	"regexp"
)

type TokenDefinition struct {
	Name    string
	Pattern regexp.Regexp
}

var tokenDefintions []TokenDefinition

type LexerToken struct {
	Name     string
	Contents string
}

func (t LexerToken) IsOfType(names ...string) bool {
	for _, name := range names {
		if t.Name == name {
			return true
		}
	}

	return false
}

func Lex(text string) ([]LexerToken, error) {
	tokens := []LexerToken{}
	var found bool
	offset := 0

	for len(text) > 0 {
		found = false

		for _, tokenDefinition := range tokenDefintions {
			match := tokenDefinition.Pattern.FindStringSubmatchIndex(text)

			if match == nil {
				continue
			}

			found = true
			contents := text[match[0]:match[1]]
			token := LexerToken{Name: tokenDefinition.Name, Contents: contents}
			offset += match[1]
			text = text[match[1]:]

			if tokenDefinition.Name != "Whitespace" && tokenDefinition.Name != "LineComment" && tokenDefinition.Name != "Newline" {
				tokens = append(tokens, token)
			}
		}

		if !found {
			return tokens, fmt.Errorf("could not find token matching %s", text)
		}
	}

	return tokens, nil
}

func init() {
	tokenDefintions = []TokenDefinition{
		{"Newline", *regexp.MustCompile(`^\n`)},
		{"LineComment", *regexp.MustCompile(`^//[^\n]+`)},
		{"Float", *regexp.MustCompile(`^-?\d*\.\d+`)},
		{"Integer", *regexp.MustCompile(`^-?\d+`)},
		{"Boolean", *regexp.MustCompile(`^(true|false)`)},
		{"Null", *regexp.MustCompile(`^null`)},
		{"Keyword", *regexp.MustCompile(`^[\w_]+`)},
		{"String", *regexp.MustCompile(`^\".+?\"`)},
		{"LeftCurlyBrace", *regexp.MustCompile(`^\{`)},
		{"RightCurlyBrace", *regexp.MustCompile(`^\}`)},
		{"LeftSquareBracket", *regexp.MustCompile(`^\[`)},
		{"RightSquareBracket", *regexp.MustCompile(`^\]`)},
		{"Whitespace", *regexp.MustCompile(`^\s+`)},
	}
}
