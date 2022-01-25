package tree

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Query is an interface that defines the methods to query a node.
type Query interface {
	Exec(n Node) (Node, error)
}

type nopQuery struct{}

func (q nopQuery) Exec(n Node) (Node, error) {
	return n, nil
}

// NopQuery with a no-op Exec method wrapping the provided interface.
var NopQuery Query = nopQuery{}

type valueQuery struct {
	v interface{}
}

func (q valueQuery) Exec(n Node) (Node, error) {
	return ToValue(q.v), nil
}

func ValueQuery(v interface{}) Query {
	return valueQuery{v: v}
}

type MapQuery string

func (q MapQuery) Exec(n Node) (Node, error) {
	if n == nil {
		return nil, nil
	}
	key := string(q)
	if m := n.Map(); m != nil {
		return m[key], nil
	}
	if a := n.Array(); a != nil {
		c := Array{}
		for _, aa := range a {
			if m := aa.Map(); m != nil {
				if v, ok := m[key]; ok {
					c = append(c, v)
				}
			}
		}
		return c, nil
	}
	return nil, fmt.Errorf(`Cannot index array with string "%s"`, key)
}

type ArrayQuery int

func (q ArrayQuery) Exec(n Node) (Node, error) {
	if n == nil {
		return nil, nil
	}
	index := int(q)
	if a := n.Array(); a != nil {
		return a[index], nil
	}
	return nil, fmt.Errorf(`Cannot index array with index %d`, index)
}

type ArrayRangeQuery []int

func (q ArrayRangeQuery) Exec(n Node) (Node, error) {
	if n == nil {
		return nil, nil
	}
	if len(q) != 2 {
		return nil, fmt.Errorf(`Invalid array range %v`, q)
	}
	if a := n.Array(); a != nil {
		return a[q[0] : q[1]+1], nil
	}
	return nil, fmt.Errorf(`Cannot index array with range %d:%d`, q[0], q[1])
}

type CollectionQuery []Query

func (qs CollectionQuery) Exec(n Node) (Node, error) {
	c := Array{}
	for _, q := range qs {
		nn, err := q.Exec(n)
		if err != nil {
			return nil, err
		}
		c = append(c, nn)
	}
	return c, nil
}

type FilterQuery []Query

func (qs FilterQuery) Exec(n Node) (Node, error) {
	nn := n
	for _, q := range qs {
		var err error
		nn, err = q.Exec(nn)
		if err != nil {
			return nil, err
		}
	}
	return nn, nil
}

type Selector interface {
	Matches(n Node) (bool, error)
}

type SelectQuery struct {
	Selectors []Selector
	Or        bool
}

func (q SelectQuery) Eval(n Node) (bool, error) {
	if len(q.Selectors) == 0 {
		return true, nil
	}
	for _, s := range q.Selectors {
		ok, err := s.Matches(n)
		if err != nil {
			return false, err
		}
		if ok {
			if q.Or {
				break
			}
		} else if !q.Or {
			return false, nil
		}
	}
	return true, nil
}

func (q SelectQuery) Exec(n Node) (Node, error) {
	if a := n.Array(); a != nil {
		c := Array{}
		for _, nn := range a {
			ok, err := q.Eval(nn)
			if err != nil {
				return nil, err
			}
			if ok {
				c = append(c, nn)
			}
		}
		return c, nil
	}
	return nil, nil
}

type Comparator struct {
	Operator Operator
	Left     Query
	Right    Query
}

var _ Selector = (*Comparator)(nil)

func (c Comparator) Matches(n Node) (bool, error) {
	l, err := c.Left.Exec(n)
	if err != nil {
		return false, err
	}
	r, err := c.Right.Exec(n)
	if err != nil {
		return false, err
	}
	if l == nil || r == nil {
		return (l == nil && r == nil), nil
	}
	return l.Value().Compare(c.Operator, r.Value()), nil
}

var tokenRegexp = regexp.MustCompile(`"([^"]*)"|(and|==|<=|>=|[\.\|\(\)\[\]<>:])|(\w+)`)

func ParseQuery(expr string) (Query, error) {
	token, err := tokenizeQuery(expr)
	if err != nil {
		return nil, err
	}
	fmt.Println("---")
	printToken(os.Stdout, token, 0)
	fmt.Println("---")
	return tokenToQuery(token, expr)
}

type token struct {
	cmd      string
	word     string
	parent   *token
	children []*token
}

func tokenizeQuery(expr string) (*token, error) {
	current := &token{}
	ms := tokenRegexp.FindAllStringSubmatch(expr, -1)
	for _, m := range ms {
		if current == nil {
			return nil, fmt.Errorf(`Syntax error: too right brackets: "%s"`, expr)
		}
		if m[1] != "" || m[3] != "" {
			word := m[1]
			if word == "" {
				word = m[3]
			}
			var lastChild *token
			if len(current.children) > 0 {
				lastChild = current.children[len(current.children)-1]
			}
			if lastChild != nil && lastChild.cmd == "." {
				lastChild.word = word
				continue
			}
			t := &token{word: word}
			current.children = append(current.children, t)
			continue
		}
		t := &token{cmd: m[2], parent: current}
		switch m[2] {
		case "]":
			if current.cmd != "[" {
				return nil, fmt.Errorf(`Syntax error: no right brackets: "%s"`, expr)
			}
			current = current.parent
		case ")":
			if current.cmd != "(" {
				return nil, fmt.Errorf(`Syntax error: no right brackets: "%s"`, expr)
			}
			current = current.parent
		case "[", "(":
			current.children = append(current.children, t)
			current = t
		default:
			current.children = append(current.children, t)
		}
	}
	if current.parent != nil {
		return nil, fmt.Errorf(`Syntax error: no right brackets: "%s"`, expr)
	}
	return current, nil
}

// printToken prints token tree for debug.
func printToken(w io.Writer, t *token, depth int) {
	indent := strings.Repeat("\t", depth)
	fmt.Fprintf(w, "%s{%s} %s\n", indent, t.cmd, t.word)
	if len(t.children) > 0 {
		depth++
		for _, c := range t.children {
			printToken(w, c, depth)
		}
	}
}

func tokenToQuery(t *token, expr string) (Query, error) {
	child := len(t.children)

	switch t.cmd {
	case "":
		if child == 0 {
			return ValueQuery(t.word), nil
		}
	case ".":
		if t.word != "" {
			return MapQuery(t.word), nil
		}
		return NopQuery, nil
	case "[":
		child := len(t.children)
		if child == 0 {
			return CollectionQuery{NopQuery}, nil
		}
		if child == 1 {
			i, err := strconv.Atoi(t.children[0].word)
			if err != nil {
				return nil, fmt.Errorf(`Syntax error: invalid index: "%s"`, expr)
			}
			return ArrayQuery(i), nil
		}
		if child == 3 && t.children[1].cmd == ":" {
			from, err := strconv.Atoi(t.children[0].word)
			if err != nil {
				return nil, fmt.Errorf(`Syntax error: invalid index: "%s"`, expr)
			}
			to, err := strconv.Atoi(t.children[2].word)
			if err != nil {
				return nil, fmt.Errorf(`Syntax error: invalid index: "%s"`, expr)
			}
			return ArrayRangeQuery{from, to}, nil
		}
		selectors, err := tokensToSelectors(t.children, expr)
		if err != nil {
			return nil, err
		}
		return SelectQuery{Selectors: selectors}, nil
	}
	if child == 0 {
		return nil, fmt.Errorf(`Syntax error: "%s": %#v`, expr, t)
	}
	if child == 1 {
		return tokenToQuery(t.children[0], expr)
	}
	var fq FilterQuery
	for _, c := range t.children {
		q, err := tokenToQuery(c, expr)
		if err != nil {
			return nil, err
		}
		fq = append(fq, q)
	}
	return fq, nil
}

func tokensToSelectors(ts []*token, expr string) ([]Selector, error) {
	var groups [][]*token
	off := 0
	for i, t := range ts {
		switch t.cmd {
		case "and":
			groups = append(groups, ts[off:i])
			off = i + 1
		}
	}
	groups = append(groups, ts[off:])

	var selectors []Selector
	for _, group := range groups {
		op := -1
		for i, t := range group {
			switch Operator(t.cmd) {
			case EQ, GT, GE, LT, LE:
				op = i
				break
			}
		}
		if op >= 0 {
			left, err := tokenToQuery(&token{children: group[0:op]}, expr)
			if err != nil {
				return nil, err
			}
			right, err := tokenToQuery(&token{children: group[op+1:]}, expr)
			if err != nil {
				return nil, err
			}
			selectors = append(selectors, Comparator{
				Operator: Operator(group[op].cmd),
				Left:     left,
				Right:    right,
			})
		}
	}
	return selectors, nil
}
