package urlx

// Code generated by /var/folders/3p/215s80gx7rx2qs2g9v5r601c0000gp/T/go-build788968830/b001/exe/peg -switch -inline parse.peg DO NOT EDIT.

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleURL
	ruleURI
	ruleScheme
	ruleHost
	ruleIPPort
	ruleHostNamePort
	ruleBracketsPort
	ruleIP
	ruleIPV4
	ruleHostName
	ruleOnlyPort
	rulePort
	rulePath
	ruleRelPath
	ruleAbsPath
	ruleBrackets
	ruleEnd
	rulePegText
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
)

var rul3s = [...]string{
	"Unknown",
	"URL",
	"URI",
	"Scheme",
	"Host",
	"IPPort",
	"HostNamePort",
	"BracketsPort",
	"IP",
	"IPV4",
	"HostName",
	"OnlyPort",
	"Port",
	"Path",
	"RelPath",
	"AbsPath",
	"Brackets",
	"End",
	"PegText",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(w io.Writer, pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Fprintf(w, " ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Fprintf(w, "%v %v\n", rule, quote)
			} else {
				fmt.Fprintf(w, "\x1B[36m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(w io.Writer, buffer string) {
	node.print(w, false, buffer)
}

func (node *node32) PrettyPrint(w io.Writer, buffer string) {
	node.print(w, true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(os.Stdout, buffer)
}

func (t *tokens32) WriteSyntaxTree(w io.Writer, buffer string) {
	t.AST().Print(w, buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(os.Stdout, buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	tree, i := t.tree, int(index)
	if i >= len(tree) {
		t.tree = append(tree, token32{pegRule: rule, begin: begin, end: end})
		return
	}
	tree[i] = token32{pegRule: rule, begin: begin, end: end}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type parser struct {
	url uri

	Buffer string
	buffer []rune
	rules  [27]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *parser) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *parser) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *parser
	max token32
}

func (e *parseError) Error() string {
	tokens, err := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		err += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return err
}

func (p *parser) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *parser) WriteSyntaxTree(w io.Writer) {
	p.tokens32.WriteSyntaxTree(w, p.Buffer)
}

func (p *parser) SprintSyntaxTree() string {
	var bldr strings.Builder
	p.WriteSyntaxTree(&bldr)
	return bldr.String()
}

func (p *parser) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for _, token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:

			p.url.uri = text

		case ruleAction1:

			p.url.scheme = text[:len(text)-1]

		case ruleAction2:

			p.url.host = text

		case ruleAction3:

			p.url.host = text

		case ruleAction4:

			p.url.port = text

		case ruleAction5:

			p.url.path = text

		case ruleAction6:

			p.url.path = text

		case ruleAction7:

			p.url.host = "[::]"

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func Pretty(pretty bool) func(*parser) error {
	return func(p *parser) error {
		p.Pretty = pretty
		return nil
	}
}

func Size(size int) func(*parser) error {
	return func(p *parser) error {
		p.tokens32 = tokens32{tree: make([]token32, 0, size)}
		return nil
	}
}
func (p *parser) Init(options ...func(*parser) error) error {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	for _, option := range options {
		err := option(p)
		if err != nil {
			return err
		}
	}
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := p.tokens32
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 URL <- <(URI / Path / Scheme / Host / (OnlyPort End))> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
				{
					position2, tokenIndex2 := position, tokenIndex
					{
						position4 := position
						{
							position5 := position
							if !_rules[ruleScheme]() {
								goto l3
							}
							if buffer[position] != rune('/') {
								goto l3
							}
							position++
							if buffer[position] != rune('/') {
								goto l3
							}
							position++
							if !_rules[ruleHost]() {
								goto l3
							}
							{
								position6, tokenIndex6 := position, tokenIndex
								if !_rules[rulePath]() {
									goto l6
								}
								goto l7
							l6:
								position, tokenIndex = position6, tokenIndex6
							}
						l7:
							add(rulePegText, position5)
						}
						{
							add(ruleAction0, position)
						}
						add(ruleURI, position4)
					}
					goto l2
				l3:
					position, tokenIndex = position2, tokenIndex2
					if !_rules[rulePath]() {
						goto l9
					}
					goto l2
				l9:
					position, tokenIndex = position2, tokenIndex2
					if !_rules[ruleScheme]() {
						goto l10
					}
					goto l2
				l10:
					position, tokenIndex = position2, tokenIndex2
					if !_rules[ruleHost]() {
						goto l11
					}
					goto l2
				l11:
					position, tokenIndex = position2, tokenIndex2
					{
						position12 := position
						{
							position13, tokenIndex13 := position, tokenIndex
							if buffer[position] != rune(':') {
								goto l14
							}
							position++
							if !_rules[rulePort]() {
								goto l14
							}
							goto l13
						l14:
							position, tokenIndex = position13, tokenIndex13
							if !_rules[rulePort]() {
								goto l0
							}
						}
					l13:
						add(ruleOnlyPort, position12)
					}
					{
						position15 := position
						{
							position16, tokenIndex16 := position, tokenIndex
							if !matchDot() {
								goto l16
							}
							goto l0
						l16:
							position, tokenIndex = position16, tokenIndex16
						}
						add(ruleEnd, position15)
					}
				}
			l2:
				add(ruleURL, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 URI <- <(<(Scheme ('/' '/') Host Path?)> Action0)> */
		nil,
		/* 2 Scheme <- <(<(([a-z] / [A-Z]) ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('+') '+') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))* ':')> Action1)> */
		func() bool {
			position18, tokenIndex18 := position, tokenIndex
			{
				position19 := position
				{
					position20 := position
					{
						position21, tokenIndex21 := position, tokenIndex
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l22
						}
						position++
						goto l21
					l22:
						position, tokenIndex = position21, tokenIndex21
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l18
						}
						position++
					}
				l21:
				l23:
					{
						position24, tokenIndex24 := position, tokenIndex
						{
							switch buffer[position] {
							case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l24
								}
								position++
							case '+':
								if buffer[position] != rune('+') {
									goto l24
								}
								position++
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l24
								}
								position++
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l24
								}
								position++
							}
						}

						goto l23
					l24:
						position, tokenIndex = position24, tokenIndex24
					}
					if buffer[position] != rune(':') {
						goto l18
					}
					position++
					add(rulePegText, position20)
				}
				{
					add(ruleAction1, position)
				}
				add(ruleScheme, position19)
			}
			return true
		l18:
			position, tokenIndex = position18, tokenIndex18
			return false
		},
		/* 3 Host <- <(IPPort / HostNamePort / BracketsPort / ((&('[') Brackets) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') IPV4) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') HostName)))> */
		func() bool {
			position27, tokenIndex27 := position, tokenIndex
			{
				position28 := position
				{
					position29, tokenIndex29 := position, tokenIndex
					{
						position31 := position
						{
							position32 := position
							if !_rules[ruleIPV4]() {
								goto l30
							}
							add(ruleIP, position32)
						}
						if buffer[position] != rune(':') {
							goto l30
						}
						position++
						if !_rules[rulePort]() {
							goto l30
						}
						add(ruleIPPort, position31)
					}
					goto l29
				l30:
					position, tokenIndex = position29, tokenIndex29
					{
						position34 := position
						if !_rules[ruleHostName]() {
							goto l33
						}
						if buffer[position] != rune(':') {
							goto l33
						}
						position++
						if !_rules[rulePort]() {
							goto l33
						}
						add(ruleHostNamePort, position34)
					}
					goto l29
				l33:
					position, tokenIndex = position29, tokenIndex29
					{
						position36 := position
						if !_rules[ruleBrackets]() {
							goto l35
						}
						if buffer[position] != rune(':') {
							goto l35
						}
						position++
						if !_rules[rulePort]() {
							goto l35
						}
						add(ruleBracketsPort, position36)
					}
					goto l29
				l35:
					position, tokenIndex = position29, tokenIndex29
					{
						switch buffer[position] {
						case '[':
							if !_rules[ruleBrackets]() {
								goto l27
							}
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if !_rules[ruleIPV4]() {
								goto l27
							}
						default:
							if !_rules[ruleHostName]() {
								goto l27
							}
						}
					}

				}
			l29:
				add(ruleHost, position28)
			}
			return true
		l27:
			position, tokenIndex = position27, tokenIndex27
			return false
		},
		/* 4 IPPort <- <(IP ':' Port)> */
		nil,
		/* 5 HostNamePort <- <(HostName ':' Port)> */
		nil,
		/* 6 BracketsPort <- <(Brackets ':' Port)> */
		nil,
		/* 7 IP <- <IPV4> */
		nil,
		/* 8 IPV4 <- <(<([0-9]+ '.' [0-9]+ '.' [0-9]+ '.' [0-9]+)> Action2)> */
		func() bool {
			position42, tokenIndex42 := position, tokenIndex
			{
				position43 := position
				{
					position44 := position
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l42
					}
					position++
				l45:
					{
						position46, tokenIndex46 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l46
						}
						position++
						goto l45
					l46:
						position, tokenIndex = position46, tokenIndex46
					}
					if buffer[position] != rune('.') {
						goto l42
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l42
					}
					position++
				l47:
					{
						position48, tokenIndex48 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l48
						}
						position++
						goto l47
					l48:
						position, tokenIndex = position48, tokenIndex48
					}
					if buffer[position] != rune('.') {
						goto l42
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l42
					}
					position++
				l49:
					{
						position50, tokenIndex50 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l50
						}
						position++
						goto l49
					l50:
						position, tokenIndex = position50, tokenIndex50
					}
					if buffer[position] != rune('.') {
						goto l42
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l42
					}
					position++
				l51:
					{
						position52, tokenIndex52 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l52
						}
						position++
						goto l51
					l52:
						position, tokenIndex = position52, tokenIndex52
					}
					add(rulePegText, position44)
				}
				{
					add(ruleAction2, position)
				}
				add(ruleIPV4, position43)
			}
			return true
		l42:
			position, tokenIndex = position42, tokenIndex42
			return false
		},
		/* 9 HostName <- <(<(([a-z] / [A-Z]) ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))*)> Action3)> */
		func() bool {
			position54, tokenIndex54 := position, tokenIndex
			{
				position55 := position
				{
					position56 := position
					{
						position57, tokenIndex57 := position, tokenIndex
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l58
						}
						position++
						goto l57
					l58:
						position, tokenIndex = position57, tokenIndex57
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l54
						}
						position++
					}
				l57:
				l59:
					{
						position60, tokenIndex60 := position, tokenIndex
						{
							switch buffer[position] {
							case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l60
								}
								position++
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l60
								}
								position++
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l60
								}
								position++
							}
						}

						goto l59
					l60:
						position, tokenIndex = position60, tokenIndex60
					}
					add(rulePegText, position56)
				}
				{
					add(ruleAction3, position)
				}
				add(ruleHostName, position55)
			}
			return true
		l54:
			position, tokenIndex = position54, tokenIndex54
			return false
		},
		/* 10 OnlyPort <- <((':' Port) / Port)> */
		nil,
		/* 11 Port <- <(<[0-9]+> Action4)> */
		func() bool {
			position64, tokenIndex64 := position, tokenIndex
			{
				position65 := position
				{
					position66 := position
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l64
					}
					position++
				l67:
					{
						position68, tokenIndex68 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l68
						}
						position++
						goto l67
					l68:
						position, tokenIndex = position68, tokenIndex68
					}
					add(rulePegText, position66)
				}
				{
					add(ruleAction4, position)
				}
				add(rulePort, position65)
			}
			return true
		l64:
			position, tokenIndex = position64, tokenIndex64
			return false
		},
		/* 12 Path <- <(RelPath / AbsPath)> */
		func() bool {
			position70, tokenIndex70 := position, tokenIndex
			{
				position71 := position
				{
					position72, tokenIndex72 := position, tokenIndex
					{
						position74 := position
						{
							position75 := position
							if buffer[position] != rune('.') {
								goto l73
							}
							position++
							if buffer[position] != rune('/') {
								goto l73
							}
							position++
						l76:
							{
								position77, tokenIndex77 := position, tokenIndex
								if !matchDot() {
									goto l77
								}
								goto l76
							l77:
								position, tokenIndex = position77, tokenIndex77
							}
							add(rulePegText, position75)
						}
						{
							add(ruleAction5, position)
						}
						add(ruleRelPath, position74)
					}
					goto l72
				l73:
					position, tokenIndex = position72, tokenIndex72
					{
						position79 := position
						{
							position80 := position
							if buffer[position] != rune('/') {
								goto l70
							}
							position++
						l81:
							{
								position82, tokenIndex82 := position, tokenIndex
								if !matchDot() {
									goto l82
								}
								goto l81
							l82:
								position, tokenIndex = position82, tokenIndex82
							}
							add(rulePegText, position80)
						}
						{
							add(ruleAction6, position)
						}
						add(ruleAbsPath, position79)
					}
				}
			l72:
				add(rulePath, position71)
			}
			return true
		l70:
			position, tokenIndex = position70, tokenIndex70
			return false
		},
		/* 13 RelPath <- <(<('.' '/' .*)> Action5)> */
		nil,
		/* 14 AbsPath <- <(<('/' .*)> Action6)> */
		nil,
		/* 15 Brackets <- <('[' ':' ':' ']' Action7)> */
		func() bool {
			position86, tokenIndex86 := position, tokenIndex
			{
				position87 := position
				if buffer[position] != rune('[') {
					goto l86
				}
				position++
				if buffer[position] != rune(':') {
					goto l86
				}
				position++
				if buffer[position] != rune(':') {
					goto l86
				}
				position++
				if buffer[position] != rune(']') {
					goto l86
				}
				position++
				{
					add(ruleAction7, position)
				}
				add(ruleBrackets, position87)
			}
			return true
		l86:
			position, tokenIndex = position86, tokenIndex86
			return false
		},
		/* 16 End <- <!.> */
		nil,
		nil,
		/* 19 Action0 <- <{
		  p.url.uri = text
		}> */
		nil,
		/* 20 Action1 <- <{
		  p.url.scheme = text[:len(text)-1]
		}> */
		nil,
		/* 21 Action2 <- <{
		  p.url.host = text
		}> */
		nil,
		/* 22 Action3 <- <{
		  p.url.host = text
		}> */
		nil,
		/* 23 Action4 <- <{
		  p.url.port = text
		}> */
		nil,
		/* 24 Action5 <- <{
		  p.url.path = text
		}> */
		nil,
		/* 25 Action6 <- <{
		  p.url.path = text
		}> */
		nil,
		/* 26 Action7 <- <{
		  p.url.host = "[::]"
		}> */
		nil,
	}
	p.rules = _rules
	return nil
}