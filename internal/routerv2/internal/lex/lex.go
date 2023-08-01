package lex

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		states: []state{initialState},
	}
	l.step()
	return l
}

func Lex(input string) []Token {
	l := New(input)
	var tokens []Token
	for tok := l.Next(); tok.Type != TEnd; tok = l.Next() {
		tokens = append(tokens, tok)
	}
	return tokens
}

// Print the input as tokens
func Print(input string) string {
	tokens := Lex(input)
	stoken := make([]string, len(tokens))
	for i, token := range tokens {
		stoken[i] = token.String()
	}
	return strings.Join(stoken, " ")
}

type state func(*Lexer) Type

type Lexer struct {
	input  string
	start  int    // Index to the start of the current token
	end    int    // Index to the end of the current token
	cp     rune   // Code point being considered
	next   int    // Index to the next rune to be considered
	err    string // Error message for an error token
	states []state
}

// Use -1 to indicate the end of the file
const eof = -1

func (l *Lexer) step() {
	codePoint, width := utf8.DecodeRuneInString(l.input[l.next:])
	if width == 0 {
		codePoint = eof
	}
	l.cp = codePoint
	l.end = l.next
	l.next += width
}

func (l *Lexer) errorf(msg string, args ...interface{}) Type {
	l.err = fmt.Sprintf(msg, args...)
	return TError
}

func (l *Lexer) Next() Token {
	l.start = l.end
	kind := l.nextType()
	if kind == TError {
		token := Token{
			Type: kind,
			Text: l.err,
		}
		l.err = ""
		return token
	}
	return Token{
		Type: kind,
		Text: l.input[l.start:l.end],
	}
}

func (l *Lexer) nextType() Type {
	return l.states[len(l.states)-1](l)
}

func (l *Lexer) eof() bool {
	return l.cp == eof
}

func (l *Lexer) pushState(s state) {
	l.states = append(l.states, s)
}

func (l *Lexer) popState() {
	l.states = l.states[:len(l.states)-1]
}

func (l *Lexer) stepUntil(delims ...string) bool {
	for !l.eof() {
		for _, delim := range delims {
			if strings.HasPrefix(l.input[l.end:l.end+len(delim)], delim) {
				return true
			}
		}
		l.step()
	}
	return false
}

func (l *Lexer) stepEach(fn func(rune) bool) bool {
	steps := 0
	for !l.eof() && fn(l.cp) {
		l.step()
		steps++
	}
	return steps > 0
}

func initialState(l *Lexer) Type {
	switch l.cp {
	case eof:
		return TEnd
	case '/':
		l.step()
		l.pushState(pathState)
		return TSlash
	default:
		l.stepUntil("/")
		return l.errorf(`path must start with a slash /`)
	}
}

func pathState(l *Lexer) Type {
	switch l.cp {
	case eof:
		return TEnd
	case '/':
		l.step()
		return TSlash
	case '{':
		l.step()
		l.pushState(slotNameState)
		return TOpenCurly
	}
	if l.stepEach(isPathChar) {
		return TPath
	}
	for {
		l.step()
		if l.eof() || l.cp == '/' || l.cp == '{' || isPathChar(l.cp) {
			break
		}
	}
	return l.errorf("unexpected character '%s' in path", l.input[l.start:l.end])
}

func slotNameState(l *Lexer) Type {
	if l.eof() {
		l.pushState(endState)
		return l.errorf("unclosed slot")
	}
	if l.stepEach(isSlotChar) {
		l.popState()
		l.pushState(slotModifierState)
		return TSlot
	}
	l.step()
	return l.errorf("slot can't start with '%s'", l.input[l.start:l.end])
}

func slotModifierState(l *Lexer) Type {
	switch l.cp {
	case eof:
		l.pushState(endState)
		return l.errorf("unclosed slot")
	case '}':
		l.step()
		l.popState()
		return TCloseCurly
	case '?':
		l.step()
		l.popState()
		l.pushState(slotCloseState)
		return TQuestion
	case '*':
		l.step()
		l.popState()
		l.pushState(slotCloseState)
		return TStar
	case '|':
		l.step()
		l.pushState(slotRegexpState)
		return TPipe
	default:
		for {
			l.step()
			if l.eof() || l.cp == '}' || l.cp == '?' || l.cp == '*' || l.cp == '|' {
				break
			}
		}
		return l.errorf("invalid character '%s' in slot", l.input[l.start:l.end])
	}
}

func slotCloseState(l *Lexer) Type {
	switch l.cp {
	case eof:
		l.pushState(endState)
		return l.errorf("unclosed slot")
	case '}':
		l.step()
		l.popState()
		return TCloseCurly
	}
	if l.stepUntil("}") {
		return l.errorf(`expected '}' but got '%s'`, l.input[l.start:l.end])
	}
	l.pushState(endState)
	return l.errorf("unclosed slot")
}

func slotRegexpState(l *Lexer) Type {
	depth := 0
loop:
	for !l.eof() {
		switch l.cp {
		case '{':
			l.step()
			depth++
		case '}':
			if depth > 0 {
				depth--
				l.step()
				continue loop
			}
			l.popState()
			return TRegexp
		default:
			l.step()
		}
	}
	return TEnd
}

func endState(l *Lexer) Type {
	l.step()
	l.popState()
	return TEnd
}

func isLowerLetter(r rune) bool {
	return unicode.IsLetter(r) && unicode.IsLower(r)
}

func isNumber(r rune) bool {
	return unicode.IsNumber(r)
}

func isDash(r rune) bool {
	return r == '-'
}

func isUnderscore(r rune) bool {
	return r == '_'
}

func isPeriod(r rune) bool {
	return r == '.'
}

func isLowerAlpha(r rune) bool {
	return 'a' <= r && r <= 'z'
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isPathChar(r rune) bool {
	return isLowerLetter(r) || isNumber(r) || isDash(r) || isUnderscore(r) || isPeriod(r)
}

func isSlotChar(r rune) bool {
	return isLowerAlpha(r) || isDigit(r) || isUnderscore(r)
}
