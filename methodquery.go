package tree

import (
	"fmt"
	"sync"
)

var (
	// newMethodQueryFuncs stores registered method query factory functions
	newMethodQueryFuncs = make(map[string]NewMethodQueryFunc)
	// methodQueryMux protects concurrent access to newMethodQueryFuncs
	methodQueryMux      sync.Mutex
)

// NewMethodQueryFunc is a factory function type for creating method queries.
// It takes string arguments and returns a Query implementation.
type NewMethodQueryFunc func(args ...string) (Query, error)

// RegisterNewMethodQueryFunc registers a factory function for a method query.
// This allows dynamic registration of new method types at runtime.
func RegisterNewMethodQueryFunc(method string, fn NewMethodQueryFunc) {
	methodQueryMux.Lock()
	defer methodQueryMux.Unlock()

	newMethodQueryFuncs[method] = fn
}

// NewMethodQuery creates a method query for the specified method name.
// Returns an error if the method is not registered.
func NewMethodQuery(method string, args ...string) (Query, error) {
	methodQueryMux.Lock()
	defer methodQueryMux.Unlock()

	fn, ok := newMethodQueryFuncs[method]
	if !ok {
		return nil, fmt.Errorf("unknown method: %s", method)
	}
	return fn(args...)
}

// CountQuery returns the count of elements in arrays or maps.
// For other node types, returns 0.
type CountQuery struct{}

// NewCountQuery creates a new CountQuery instance.
// Arguments are ignored for this query type.
func NewCountQuery(args ...string) (Query, error) {
	return &CountQuery{}, nil
}

// Exec returns the count of elements in the node.
// For arrays and maps, returns their length. For other types, returns 0.
func (q *CountQuery) Exec(n Node) ([]Node, error) {
	switch n.Type() {
	case TypeArray:
		return ToNodeValues(len(n.Array())), nil
	case TypeMap:
		return ToNodeValues(len(n.Map())), nil
	}
	return ToNodeValues(0), nil
}

func (q *CountQuery) String() string {
	return "count()"
}

// KeysQuery returns the keys of arrays (as indices) or maps.
// For arrays, returns numeric indices. For maps, returns string keys.
type KeysQuery struct{}

// NewKeysQuery creates a new KeysQuery instance.
// Arguments are ignored for this query type.
func NewKeysQuery(args ...string) (Query, error) {
	return &KeysQuery{}, nil
}

// Exec returns the keys of the node as an array.
// For arrays, returns numeric indices. For maps, returns string keys.
// Returns nil for other node types.
func (q *KeysQuery) Exec(n Node) ([]Node, error) {
	switch n.Type() {
	case TypeArray:
		a := n.Array()
		keys := make(Array, len(a))
		for i := 0; i < len(a); i++ {
			keys[i] = NumberValue(i)
		}
		return []Node{keys}, nil
	case TypeMap:
		strKeys := n.Map().Keys()
		keys := make(Array, len(strKeys))
		for i := range strKeys {
			keys[i] = StringValue(strKeys[i])
		}
		return []Node{keys}, nil
	}
	return nil, nil
}

func (q *KeysQuery) String() string {
	return "keys()"
}

// ValuesQuery returns the values of arrays or maps as an array.
// For arrays, returns the array itself. For maps, returns values in key order.
type ValuesQuery struct{}

// NewValuesQuery creates a new ValuesQuery instance.
// Arguments are ignored for this query type.
func NewValuesQuery(args ...string) (Query, error) {
	return &ValuesQuery{}, nil
}

// Exec returns the values of the node as an array.
// For arrays, returns the array itself. For maps, returns values in key order.
// Returns nil for other node types.
func (q *ValuesQuery) Exec(n Node) ([]Node, error) {
	switch n.Type() {
	case TypeArray:
		return []Node{n.Array()}, nil
	case TypeMap:
		m := n.Map()
		keys := m.Keys()
		values := make(Array, len(keys))
		for i, key := range keys {
			values[i] = m[key]
		}
		return []Node{values}, nil
	}
	return nil, nil
}

func (q *ValuesQuery) String() string {
	return "values()"
}

// init automatically registers the built-in method queries.
func init() {
	RegisterNewMethodQueryFunc("count", NewCountQuery)
	RegisterNewMethodQueryFunc("keys", NewKeysQuery)
	RegisterNewMethodQueryFunc("values", NewValuesQuery)
}
