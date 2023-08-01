package lex

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

type Token struct {
	Type Type
	Text string
}

func (t *Token) String() string {
	s := new(strings.Builder)
	s.WriteString(string(t.Type))
	if t.Text != "" && t.Text != string(t.Type) {
		s.WriteString(":")
		s.WriteString(strconv.Quote(t.Text))
	}
	return s.String()
}

type Type string

const (
	TEnd        Type = "end"
	TError      Type = "error"
	TRegexp     Type = "regexp"
	TPath       Type = "path"
	TSlot       Type = "slot"
	TSlash      Type = "/"
	TOpenCurly  Type = "{"
	TCloseCurly Type = "}"
	TQuestion   Type = "?"
	TStar       Type = "*"
	TPipe       Type = "|"
)

// Tokens is a list of tokens
type Tokens []Token

// At returns the individual value at i. For paths, this will be a single
// character, for slots, this will be the whole slot name.
func (tokens Tokens) At(i int) string {
	for _, token := range tokens {
		switch token.Type {
		case TPath, TSlash:
			for _, char := range token.Text {
				if i == 0 {
					return string(char)
				}
				i--
			}
		case TSlot, TQuestion, TStar:
			if i == 0 {
				return token.Text
			}
			i--
		}
	}
	return ""
}

// Size returns the number individual values. For paths, this will count the
// number of characters, whereas slots will count as one.
func (tokens Tokens) Size() (n int) {
	for _, token := range tokens {
		switch token.Type {
		case TPath, TSlash:
			n += utf8.RuneCountInString(token.Text)
		case TSlot, TQuestion, TStar:
			n++
		}
	}
	return n
}

// Split the token list into two lists of tokens. If we land in the middle of a
// path token then split that token into two path tokens.
func (tokens Tokens) Split(at int) []Tokens {
	for i, token := range tokens {
		switch token.Type {
		case TPath, TSlash:
			for j := range token.Text {
				if at != 0 {
					at--
					continue
				}
				left, right := token.Text[:j], token.Text[j:]
				// At the edge
				if left == "" || right == "" {
					if i > 0 && i < len(tokens) {
						return []Tokens{tokens[:i], tokens[i:]}
					}
					return []Tokens{tokens}
				}
				newToken := Token{token.Type, left}
				leftTokens := append(append(Tokens{}, tokens[:i]...), newToken)
				rightTokens := append(Tokens{}, tokens[i:]...)
				rightTokens[0].Text = right
				return []Tokens{leftTokens, rightTokens}
			}
		case TSlot, TQuestion, TStar:
			if at != 0 {
				at--
				continue
			}
			if i > 0 && i < len(tokens) {
				return []Tokens{tokens[:i], tokens[i:]}
			}
			return []Tokens{tokens}
		}
	}
	return []Tokens{tokens}
}

// String returns a string of tokens for testing
func (tokens Tokens) String() string {
	var ts []string
	for _, t := range tokens {
		s := t.String()
		if s == "" {
			continue
		}
		ts = append(ts, s)
	}
	return strings.Join(ts, " ")
}
