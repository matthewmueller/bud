package radix

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/matthewmueller/bud/internal/routerv2/internal/lex"
)

func New() *Tree {
	return &Tree{}
}

type matchFunc = func(path string) (index int, slots []*Slot)

type node struct {
	Tokens   lex.Tokens
	Route    string
	Match    matchFunc
	Children []*node
	Wilds    []*node
}

func (n *node) IsWild() bool {
	token := n.Tokens[0]
	switch token.Type {
	case lex.TSlot, lex.TQuestion, lex.TStar:
		return true
	default:
		return false
	}
}

// Priority of the node
func (n *node) Priority() (priority int) {
	for _, token := range n.Tokens {
		switch token.Type {
		case lex.TSlash, lex.TPath:
			priority++
		}
	}
	return priority
}

type Match struct {
	Route string
	Slots []*Slot
}

type Slot struct {
	Key   string
	Value string
}

type Route struct {
	Path    string
	Handler http.Handler
}

type Tree struct {
	root *node
}

func (t *Tree) Insert(route string) error {
	lexer := lex.New(route)
	var tokens lex.Tokens
	for tok := lexer.Next(); tok.Type != lex.TEnd; tok = lexer.Next() {
		switch tok.Type {
		case lex.TQuestion:
			// Each optional tokens insert two routes
			if err := t.insert(stripTokenTrail(tokens), route); err != nil {
				return err
			}
			// Make the optional token required
			tokens = append(tokens, lex.Token{
				Text: strings.TrimRight(tok.Text, "?"),
				Type: lex.TSlot,
			})
		case lex.TStar:
			// Each optional tokens insert two routes
			if err := t.insert(stripTokenTrail(tokens), route); err != nil {
				return err
			}
			tokens = append(tokens, tok)
		case lex.TError:
			// Error parsing the route
			return errors.New(tok.Text)
		default:
			tokens = append(tokens, tok)
		}
	}
	return t.insert(tokens, route)
}

func (t *Tree) insert(tokens lex.Tokens, route string) error {
	if t.root == nil {
		t.root = &node{
			Tokens: tokens,
			Match:  matcher(tokens),
			Route:  route,
		}
		return nil
	}
	return t.insertAt(t.root, tokens, route)
}

func (t *Tree) insertAt(parent *node, tokens lex.Tokens, route string) error {
	// Compute the longest common prefix between new path and the node's path
	// before slots.
	lcp := longestCommonPrefix(tokens, parent.Tokens)
	parts := tokens.Split(lcp)
	inTreeAlready := len(parts) == 1
	// If longest common prefix is not the same length as the node path, We need
	// to split the node path into parent and child to prepare for another child.
	if lcp < parent.Tokens.Size() {
		parent = splitAt(parent, lcp)
		// This set of tokens are already in the tree
		// E.g. We've inserted "/a", "/b", then "/". "/" will already be in the tree
		// but not have a handler
		if inTreeAlready {
			parent.Route = route
			return nil
		}
		err := insertChild(parent, &node{
			Tokens: parts[1],
			Match:  matcher(parts[1]),
			Route:  route,
		})
		if err != nil {
			return err
		}
		// Unset the parent data
		parent.Route = ""
		return nil
	}
	// This set of tokens are already in the tree. Override any prior handler if
	// there was one.
	if inTreeAlready {
		parent.Route = route
		return nil
	}
	// For the remaining, non-common part of the path, check if any of the
	// children also start with that non-common part. If so, traverse that child.
	for _, child := range parent.Children {
		if child.Tokens.At(0) == parts[1].At(0) {
			return t.insertAt(child, parts[1], route)
		}
	}
	// Recurse wild children if the wild child matches exactly
	for _, wild := range parent.Wilds {
		if wild != nil && wild.Tokens.At(0) == parts[1].At(0) {
			return t.insertAt(wild, parts[1], route)
		}
	}
	// Otherwise, insert a new child on the parent with the remaining non-common
	// part of the path.
	return insertChild(parent, &node{
		Tokens: parts[1],
		Match:  matcher(parts[1]),
		Route:  route,
	})
}

// Split the single node into a parent and child node
func splitAt(parent *node, at int) *node {
	parts := parent.Tokens.Split(at)
	if len(parts) == 1 {
		return parent
	}
	child := &node{
		Tokens:   parts[1],
		Match:    matcher(parts[1]),
		Route:    parent.Route,
		Children: parent.Children,
		Wilds:    parent.Wilds,
	}
	// Add the split child, moving all existing children into the child
	parent.Children = []*node{}
	parent.Wilds = []*node{}
	insertChild(parent, child)
	// Split the tokens up and recompile the match function
	parent.Tokens = parts[0]
	parent.Match = matcher(parts[0])
	return parent
}

// Insert a child
func insertChild(parent *node, child *node) error {
	if child.IsWild() {
		return insertWild(parent, child)
	}
	parent.Children = append(parent.Children, child)
	return nil
}

// Insert a wild child
func insertWild(parent *node, child *node) error {
	lwilds := len(parent.Wilds)
	childp := child.Priority()
	for i := 0; i < lwilds; i++ {
		wild := parent.Wilds[i]
		wildp := wild.Priority()
		// Prioritize more specific slots over less specific slots.
		if childp > wildp {
			parent.Wilds = append(parent.Wilds[:i], append([]*node{child}, parent.Wilds[i:]...)...)
			return nil
		}
		// Don't allow /:id and /:hi on the same level.
		if childp == wildp {
			return fmt.Errorf("radix: ambiguous routes %q and %q", child.Route, wild.Route)
		}
	}
	parent.Wilds = append(parent.Wilds, child)
	return nil
}

// Turn the tokens into a matcher
func matcher(tokens lex.Tokens) matchFunc {
	var matchers []matchFunc
	for _, token := range tokens {
		switch token.Type {
		case lex.TPath, lex.TSlash:
			matchers = append(matchers, matchExact(token))
		case lex.TSlot:
			matchers = append(matchers, matchSlot(token))
		case lex.TStar:
			matchers = append(matchers, matchStar(token))
		}
	}
	return compose(matchers)
}

// Compose the match functions into one function
func compose(matchers []matchFunc) matchFunc {
	return func(path string) (index int, slots []*Slot) {
		for _, match := range matchers {
			i, matchSlots := match(path)
			if i == -1 {
				return -1, matchSlots
			}
			path = path[i:]
			index += i
			slots = append(slots, matchSlots...)
		}
		return index, slots
	}
}

// Match a slot exactly (/users)
func matchExact(token lex.Token) matchFunc {
	route := token.Text
	rlen := len(route)
	return func(path string) (index int, slots []*Slot) {
		if len(path) < rlen {
			return -1, nil
		}
		if !strings.EqualFold(path[:rlen], route) {
			return -1, nil
		}
		return rlen, nil
	}
}

// Match a slot (/{id})
func matchSlot(token lex.Token) matchFunc {
	slotKey := token.Text
	return func(path string) (index int, slots []*Slot) {
		lpath := len(path)
		for i := 0; i < lpath; i++ {
			if path[i] == '.' || path[i] == '/' {
				break
			}
			index++
		}
		if index == 0 {
			return -1, nil
		}
		return index, []*Slot{
			{
				Key:   slotKey,
				Value: path[:index],
			},
		}
	}
}

// Match a star (e.g. /:path*)
func matchStar(token lex.Token) matchFunc {
	lvalue := len(token.Text)
	slotKey := token.Text[1 : lvalue-1]
	return func(path string) (index int, slots []*Slot) {
		return len(path), []*Slot{
			{
				Key:   slotKey,
				Value: path,
			},
		}
	}
}

// Match the node
func (t *Tree) match(node *node, path string, slots []*Slot) *Match {
	index, matchSlots := node.Match(path)
	if index < 0 {
		return nil
	}
	path = path[index:]
	if matchSlots != nil {
		slots = append(slots, matchSlots...)
	}
	// No more path, we're done!
	if path == "" {
		// At a junction node, but this node isn't a route, so it's not a match
		// TODO: double-check this
		if node.Route == "" {
			return nil
		}
		return &Match{
			// Handler: node.handler,
			Route: node.Route,
			Slots: slots,
		}
	}
	// First try matching the children
	for _, child := range node.Children {
		if match := t.match(child, path, slots); match != nil {
			return match
		}
	}
	// Next try matching the wild children
	for _, wild := range node.Wilds {
		if match := t.match(wild, path, slots); match != nil {
			return match
		}
	}
	return nil
}

func (t *Tree) Match(path string) (*Match, bool) {
	// A tree without any routes shouldn't panic
	if t.root == nil {
		return nil, false
	}
	match := t.match(t.root, path, []*Slot{})
	if match == nil {
		return nil, false
	}
	return match, true
}

func (t *Tree) String() string {
	return t.string(t.root, "")
}

func (t *Tree) string(n *node, indent string) string {
	if n == nil {
		return ""
	}
	route := ""
	for _, token := range n.Tokens {
		route += token.Text
	}
	kind := "c"
	if n.IsWild() {
		kind = "w"
	}
	out := fmt.Sprintf("%s%s[%d%s]\r\n", indent, route, len(n.Children)+len(n.Wilds), kind)
	for l := len(route); l > 0; l-- {
		indent += " "
	}
	for _, child := range n.Children {
		out += t.string(child, indent)
	}
	for _, wild := range n.Wilds {
		out += t.string(wild, indent)
	}
	return out
}

// strip token trail removes path tokens up to either a slot or a slash
// e.g. /:id. => /:id
//
//	/a/b => /a
func stripTokenTrail(tokens lex.Tokens) lex.Tokens {
	i := len(tokens) - 1
loop:
	for ; i >= 0; i-- {
		switch tokens[i].Type {
		case lex.TSlot:
			i++ // Include the slot
			break loop
		case lex.TSlash:
			break loop
		}
	}
	if i == 0 {
		return tokens[:1]
	}
	newTokens := make(lex.Tokens, i)
	copy(newTokens[:], tokens)
	return newTokens
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func longestCommonPrefix(a, b lex.Tokens) int {
	i := 0
	max := min(a.Size(), b.Size())
	for i < max && a.At(i) == b.At(i) {
		i++
	}
	return i
}
