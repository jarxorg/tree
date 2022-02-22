package tree

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Query is an interface that defines the methods to query a node.
type Query interface {
	Exec(n Node) ([]Node, error)
	String() string
}

// NopQuery is a query that implements no-op Exec method.
type NopQuery struct{}

// Exec returns the provided node.
func (q NopQuery) Exec(n Node) ([]Node, error) {
	return []Node{n}, nil
}

func (q NopQuery) String() string {
	return "."
}

// ValueQuery is a query that returns the constant value.
type ValueQuery struct {
	Node
}

// Exec returns the constant value.
func (q ValueQuery) Exec(n Node) ([]Node, error) {
	return []Node{q.Node}, nil
}

func (q ValueQuery) String() string {
	s, _ := MarshalJSON(q.Node)
	return string(s)
}

// MapQuery is a key of the Map that implements methods of the Query.
type MapQuery string

func (q MapQuery) Exec(n Node) ([]Node, error) {
	key := string(q)
	if m := n.Map(); m != nil {
		return []Node{m.Get(key)}, nil
	}
	if a := n.Array(); a != nil {
		return []Node{a.Get(key)}, nil
	}
	return nil, fmt.Errorf("Cannot index array with string %q", key)
}

func (q MapQuery) String() string {
	return "." + string(q)
}

// ArrayQuery is an index of the Array that implements methods of the Query.
type ArrayQuery int

func (q ArrayQuery) Exec(n Node) ([]Node, error) {
	if a := n.Array(); a != nil {
		return []Node{a.Get(int(q))}, nil
	}
	return nil, fmt.Errorf(`Cannot index array with index %d`, q)
}

func (q ArrayQuery) String() string {
	return fmt.Sprintf("[%d]", q)
}

// ArrayRangeQuery represents a range of the Array that implements methods of the Query.
type ArrayRangeQuery []int

func (q ArrayRangeQuery) Exec(n Node) ([]Node, error) {
	if len(q) != 2 {
		return nil, fmt.Errorf(`Invalid array range %s`, q.String())
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

func (q ArrayRangeQuery) String() string {
	ss := make([]string, len(q))
	for i, r := range q {
		if r != -1 {
			ss[i] = strconv.Itoa(r)
		}
	}
	return "[" + strings.Join(ss, ":") + "]"
}

// SlurpQuery is a special query that works in FilterQuery.
type SlurpQuery struct{}

// Exec returns the provided node into a single node array.
// FilterQuery calls q.Exec(Array(results)), which has the effect of to slurp
// all the results into a single node array.
func (q SlurpQuery) Exec(n Node) ([]Node, error) {
	return []Node{n}, nil
}

func (q SlurpQuery) String() string {
	return " | "
}

// FilterQuery consists of multiple queries that filter the nodes in order.
type FilterQuery []Query

func (qs FilterQuery) Exec(n Node) ([]Node, error) {
	rs := []Node{n}
	for _, q := range qs {
		switch q.(type) {
		case SlurpQuery:
			nrs, err := q.Exec(Array(rs))
			if err != nil {
				return nil, err
			}
			rs = nrs
			continue
		}
		var nrs []Node
		for _, r := range rs {
			if r == nil {
				continue
			}
			nr, err := q.Exec(r)
			if err != nil {
				return nil, err
			}
			nrs = append(nrs, nr...)
		}
		rs = nrs
	}
	return rs, nil
}

func (qs FilterQuery) String() string {
	ss := make([]string, len(qs))
	for i, q := range qs {
		ss[i] = q.String()
	}
	return strings.Join(ss, "")
}

// WalkQuery is a key of each nodes that implements methods of the Query.
type WalkQuery string

// Exec walks the specified root node and collects matching nodes using itself as a key.
func (q WalkQuery) Exec(root Node) ([]Node, error) {
	key := string(q)
	var r []Node
	err := Walk(root, func(n Node, keys []interface{}) error {
		if nn := n.Get(key); nn != nil {
			r = append(r, nn)
			return SkipWalk
		}
		return nil
	})
	if err != nil && err != SkipWalk {
		return nil, err
	}
	return r, nil
}

func (q WalkQuery) String() string {
	return ".." + string(q)
}

// Selector checks if a node is eligible for selection.
type Selector interface {
	Matches(n Node) (bool, error)
	String() string
}

// And represents selectors that combines each selector with and.
type And []Selector

// Matches returns true if all selectors returns true.
func (a And) Matches(n Node) (bool, error) {
	for _, s := range a {
		ok, err := s.Matches(n)
		if err != nil || !ok {
			return false, err
		}
	}
	return true, nil
}

func (a And) String() string {
	ss := make([]string, len(a))
	for i, s := range a {
		ss[i] = s.String()
	}
	return "(" + strings.Join(ss, " and ") + ")"
}

// Or represents selectors that combines each selector with or.
type Or []Selector

// Matches returns true if anyone returns true.
func (o Or) Matches(n Node) (bool, error) {
	for _, s := range o {
		ok, err := s.Matches(n)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

func (o Or) String() string {
	ss := make([]string, len(o))
	for i, s := range o {
		ss[i] = s.String()
	}
	return "(" + strings.Join(ss, " or ") + ")"
}

// Comparator represents a comparable selector.
type Comparator struct {
	Left  Query
	Op    Operator
	Right Query
}

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
	var l0, r0 Node
	switch len(l) {
	case 0:
		l0 = nil
	case 1:
		l0 = l[0]
	default:
		return false, fmt.Errorf("%#v returns no single value %+v", c.Left, l)
	}
	switch len(r) {
	case 0:
		r0 = nil
	case 1:
		r0 = r[0]
	default:
		return false, fmt.Errorf("%#v returns no single value %+v", c.Right, r)
	}
	if l0 == nil || r0 == nil {
		return (l0 == nil && r0 == nil), nil
	}
	return l0.Value().Compare(c.Op, r0.Value()), nil
}

func (c Comparator) String() string {
	return fmt.Sprintf("%s %s %s", c.Left, c.Op, c.Right)
}

// SelectQuery returns nodes that matched by selectors.
type SelectQuery struct {
	Selector
}

func (q SelectQuery) Exec(n Node) ([]Node, error) {
	if a := n.Array(); a != nil {
		if q.Selector == nil {
			return a, nil
		}
		var rs []Node
		for _, nn := range a {
			ok, err := q.Selector.Matches(nn)
			if err != nil {
				return nil, err
			}
			if ok {
				rs = append(rs, nn)
			}
		}
		return rs, nil
	}
	if m := n.Map(); m != nil {
		if q.Selector == nil {
			return m.Values(), nil
		}
		var rs []Node
		for _, nn := range m.Values() {
			ok, err := q.Selector.Matches(nn)
			if err != nil {
				return nil, err
			}
			if ok {
				rs = append(rs, nn)
			}
		}
		return rs, nil
	}
	return nil, nil
}

func (q SelectQuery) String() string {
	return "[" + q.Selector.String() + "]"
}

var (
	_ Selector = (And)(nil)
	_ Selector = (Or)(nil)
	_ Selector = (*Comparator)(nil)
	_ Selector = (*SelectQuery)(nil)
)

var tokenRegexp = regexp.MustCompile(`"([^"]*)"|(and|or|==|<=|>=|\.\.|[\.\[\]\(\)\|<>:])|(\w+)`)

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
	case "|":
		return SlurpQuery{}, nil
	case ".":
		if t.value != "" {
			return MapQuery(t.value), nil
		}
		return NopQuery{}, nil
	case "..":
		if t.value != "" {
			return WalkQuery(t.value), nil
		}
		return NopQuery{}, nil
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
func Find(n Node, expr string) ([]Node, error) {
	q, err := ParseQuery(expr)
	if err != nil {
		return nil, err
	}
	return q.Exec(n)
}
