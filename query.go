package tree

import (
	"fmt"
	"regexp"
	"strconv"
)

// Query is an interface that defines the methods to query a node.
type Query interface {
	Exec(n Node) (Node, error)
}

type nopQuery struct{}

// Exec returns the provided node.
func (q nopQuery) Exec(n Node) (Node, error) {
	return n, nil
}

// NopQuery is a query that implements no-op Exec method.
var NopQuery Query = nopQuery{}

// ValueQuery is a query that returns the constant value.
type ValueQuery struct {
	Node
}

// Exec returns the constant value.
func (q ValueQuery) Exec(n Node) (Node, error) {
	return q.Node, nil
}

// MapQuery is a key of the Map that implements methods of the Query.
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
		if v := a.Get(key); v != nil {
			return v, nil
		}
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

// ArrayQuery is an index of the Array that implements methods of the Query.
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

// ArrayRangeQuery represents a range of the Array that implements methods of the Query.
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

// FilterQuery consists of multiple queries that filter the nodes in order.
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

// Selector checks if a node is eligible for selection.
type Selector interface {
	Matches(n Node) (bool, error)
}

// SelectQuery returns nodes that matched by selectors.
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
	if n == nil {
		return nil, nil
	}
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

// Comparator represents a comparable selector.
type Comparator struct {
	Left  Query
	Op    Operator
	Right Query
}

var _ Selector = (*Comparator)(nil)

// Matches evaluates left and right using the operator. (eg. .id == 0)
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
	return l.Value().Compare(c.Op, r.Value()), nil
}

var tokenRegexp = regexp.MustCompile(`"([^"]*)"|(and|==|<=|>=|[\.\[\]<>:])|(\w+)`)

// ParseQuery parses the provided expr to a Query.
// See https://github.com/jarxorg/tree#Query
func ParseQuery(expr string) (Query, error) {
	token, err := tokenizeQuery(expr)
	if err != nil {
		return nil, err
	}
	return tokenToQuery(token, expr)
}

type token struct {
	cmd      string
	quoted   bool
	value    string
	parent   *token
	children []*token
}

func (t *token) toValue() Node {
	if !t.quoted {
		if t.value == "true" {
			return BoolValue(true)
		}
		if t.value == "false" {
			return BoolValue(false)
		}
		if n, err := strconv.ParseFloat(t.value, 64); err == nil {
			return NumberValue(n)
		}
	}
	return StringValue(t.value)
}

func tokenizeQuery(expr string) (*token, error) {
	current := &token{}
	ms := tokenRegexp.FindAllStringSubmatch(expr, -1)
	for _, m := range ms {
		if m[1] != "" || m[3] != "" {
			value := m[1]
			quoted := value != ""
			if !quoted {
				value = m[3]
			}
			var lastChild *token
			if len(current.children) > 0 {
				lastChild = current.children[len(current.children)-1]
			}
			if lastChild != nil && lastChild.cmd == "." {
				lastChild.value = value
				lastChild.quoted = quoted
				continue
			}
			t := &token{value: value, quoted: quoted}
			current.children = append(current.children, t)
			continue
		}
		t := &token{cmd: m[2], parent: current}
		switch m[2] {
		case "]":
			if current.cmd != "[" {
				return nil, fmt.Errorf(`Syntax error: no left bracket '[': "%s"`, expr)
			}
			current = current.parent
		case "[":
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

func tokenToQuery(t *token, expr string) (Query, error) {
	child := len(t.children)

	switch t.cmd {
	case "":
		if child == 0 {
			return ValueQuery{t.toValue()}, nil
		}
	case ".":
		if t.value != "" {
			return MapQuery(t.value), nil
		}
		return NopQuery, nil
	case "[":
		if child == 0 {
			return SelectQuery{}, nil
		}
		if child == 1 {
			i, err := strconv.Atoi(t.children[0].value)
			if err != nil {
				return nil, fmt.Errorf(`Syntax error: invalid index: "%s"`, expr)
			}
			return ArrayQuery(i), nil
		}
		if child == 3 && t.children[1].cmd == ":" {
			from, err := strconv.Atoi(t.children[0].value)
			if err != nil {
				return nil, fmt.Errorf(`Syntax error: invalid range: "%s"`, expr)
			}
			to, err := strconv.Atoi(t.children[2].value)
			if err != nil {
				return nil, fmt.Errorf(`Syntax error: invalid range: "%s"`, expr)
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
		return nil, fmt.Errorf(`Syntax error: invalid token %s: "%s"`, t.cmd, expr)
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
			selectors = append(selectors, Comparator{left, Operator(group[op].cmd), right})
		}
	}
	return selectors, nil
}

// Find finds a node from n using the Query.
func Find(n Node, expr string) (Node, error) {
	q, err := ParseQuery(expr)
	if err != nil {
		return nil, err
	}
	return q.Exec(n)
}
