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
		from, to := q[0], q[1]
		if from == -1 {
			return a[:to], nil
		} else if q[1] == -1 {
			return a[from:], nil
		}
		return a[from:to], nil
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

// WalkQuery is a key of each nodes that implements methods of the Query.
type WalkQuery string

// Exec walks the specified root node and collects matching nodes using itself as a key.
func (q WalkQuery) Exec(root Node) (Node, error) {
	key := string(q)
	c := Array{}
	err := Walk(root, func(n Node, keys []interface{}) error {
		if nn := n.Get(key); nn != nil {
			if aa := nn.Array(); aa != nil {
				c = append(c, aa...)
			} else {
				c = append(c, nn)
			}
			return SkipWalk
		}
		return nil
	})
	if err != nil && err != SkipWalk {
		return nil, err
	}
	return c, nil
}

// Selector checks if a node is eligible for selection.
type Selector interface {
	Matches(i int, n Node) (bool, error)
}

// And represents selectors that combines each selector with and.
type And []Selector

func (ss And) Matches(i int, n Node) (bool, error) {
	for _, s := range ss {
		ok, err := s.Matches(i, n)
		if err != nil || !ok {
			return false, err
		}
	}
	return true, nil
}

// Or represents selectors that combines each selector with or.
type Or []Selector

func (ss Or) Matches(i int, n Node) (bool, error) {
	for _, s := range ss {
		ok, err := s.Matches(i, n)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

// Comparator represents a comparable selector.
type Comparator struct {
	Left  Query
	Op    Operator
	Right Query
}

// Matches evaluates left and right using the operator. (eg. .id == 0)
func (c Comparator) Matches(i int, n Node) (bool, error) {
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

// SelectQuery returns nodes that matched by selectors.
type SelectQuery struct {
	Selector
}

func (q SelectQuery) Exec(n Node) (Node, error) {
	if n == nil || q.Selector == nil {
		return n, nil
	}
	if a := n.Array(); a != nil {
		c := Array{}
		for i, nn := range a {
			ok, err := q.Selector.Matches(i, nn)
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

var (
	_ Selector = (And)(nil)
	_ Selector = (Or)(nil)
	_ Selector = (*Comparator)(nil)
	_ Selector = (*SelectQuery)(nil)
)

var tokenRegexp = regexp.MustCompile(`"([^"]*)"|(and|or|==|<=|>=|\.\.|[\.\[\]\(\)<>:])|(\w+)`)

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

func (t *token) indexOfCmd(cmd string) int {
	for i, c := range t.children {
		if c.cmd == cmd {
			return i
		}
	}
	return -1
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
			if lastChild != nil && (lastChild.cmd == "." || lastChild.cmd == "..") {
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
		case "]", ")":
			if (m[2] == "]" && current.cmd != "[") || (m[2] == ")" && current.cmd != "(") {
				return nil, fmt.Errorf(`Syntax error: no left bracket: "%s"`, expr)
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
	case "..":
		if t.value != "" {
			return WalkQuery(t.value), nil
		}
		return NopQuery, nil
	case "[":
		if child == 0 {
			return SelectQuery{}, nil
		}
		if child == 1 {
			i, err := strconv.Atoi(t.children[0].value)
			if err != nil {
				return nil, fmt.Errorf(`Syntax error: invalid array index: "%s"`, expr)
			}
			return ArrayQuery(i), nil
		}
		if i := t.indexOfCmd(":"); i != -1 {
			return tokensToArrayRangeQuery(t.children, i, expr)
		}
		selector, err := tokensToSelector(t.children, expr)
		if err != nil {
			return nil, err
		}
		return SelectQuery{selector}, nil
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

func tokensToArrayRangeQuery(ts []*token, i int, expr string) (Query, error) {
	from := -1
	to := -1
	if j := i - 1; j >= 0 {
		var err error
		from, err = strconv.Atoi(ts[j].value)
		if err != nil {
			return nil, fmt.Errorf(`Syntax error: invalid array range: %q`, expr)
		}
	}
	if j := i + 1; j < len(ts) {
		var err error
		to, err = strconv.Atoi(ts[j].value)
		if err != nil {
			return nil, fmt.Errorf(`Syntax error: invalid array range: %q`, expr)
		}
	}
	return ArrayRangeQuery{from, to}, nil
}

func tokensToSelector(ts []*token, expr string) (Selector, error) {
	andOr := ""
	var groups [][]*token
	off := 0
	for i, t := range ts {
		switch t.cmd {
		case "and", "or":
			if andOr != "" && andOr != t.cmd {
				return nil, fmt.Errorf(`Syntax error: mixed and|or: %q`, expr)
			}
			andOr = t.cmd
			groups = append(groups, ts[off:i])
			off = i + 1
		case "(":
			groups = append(groups, ts[off:i])
			groups = append(groups, []*token{t})
			off = i + 1
		}
	}
	groups = append(groups, ts[off:])

	var ss []Selector
	for _, group := range groups {
		op := -1
		for i, t := range group {
			if t.cmd == "(" {
				sss, err := tokensToSelector(t.children, expr)
				if err != nil {
					return nil, err
				}
				ss = append(ss, sss)
				break
			}
			switch Operator(t.cmd) {
			case EQ, GT, GE, LT, LE:
				op = i
				break
			}
		}
		if op == -1 {
			continue
		}
		left, err := tokenToQuery(&token{children: group[0:op]}, expr)
		if err != nil {
			return nil, err
		}
		right, err := tokenToQuery(&token{children: group[op+1:]}, expr)
		if err != nil {
			return nil, err
		}
		ss = append(ss, Comparator{left, Operator(group[op].cmd), right})
	}
	if andOr == "or" {
		return Or(ss), nil
	}
	return And(ss), nil
}

// Find finds a node from n using the Query.
func Find(n Node, expr string) (Node, error) {
	q, err := ParseQuery(expr)
	if err != nil {
		return nil, err
	}
	return q.Exec(n)
}
