package tree

import (
	"fmt"
	"strconv"
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

// EmptyQuery checks if the node is empty.
// Returns true for empty arrays, empty maps, null values, and empty strings.
type EmptyQuery struct{}

// NewEmptyQuery creates a new EmptyQuery instance.
// Arguments are ignored for this query type.
func NewEmptyQuery(args ...string) (Query, error) {
	return &EmptyQuery{}, nil
}

// Exec returns whether the node is empty.
// Returns true for empty arrays, empty maps, null values, and empty strings.
func (q *EmptyQuery) Exec(n Node) ([]Node, error) {
	switch n.Type() {
	case TypeArray:
		return ToNodeValues(len(n.Array()) == 0), nil
	case TypeMap:
		return ToNodeValues(len(n.Map()) == 0), nil
	case TypeStringValue:
		return ToNodeValues(n.Value().String() == ""), nil
	case TypeNilValue:
		return ToNodeValues(true), nil
	}
	return ToNodeValues(false), nil
}

func (q *EmptyQuery) String() string {
	return "empty()"
}

// TypeQuery returns the type name of the node.
// Returns "array", "object", "string", "number", "boolean", or "null".
type TypeQuery struct{}

// NewTypeQuery creates a new TypeQuery instance.
// Arguments are ignored for this query type.
func NewTypeQuery(args ...string) (Query, error) {
	return &TypeQuery{}, nil
}

// Exec returns the type name of the node as a string.
// Returns "array", "object", "string", "number", "boolean", or "null".
func (q *TypeQuery) Exec(n Node) ([]Node, error) {
	var typeName string
	switch n.Type() {
	case TypeArray:
		typeName = "array"
	case TypeMap:
		typeName = "object"
	case TypeStringValue:
		typeName = "string"
	case TypeNumberValue:
		typeName = "number"
	case TypeBoolValue:
		typeName = "boolean"
	case TypeNilValue:
		typeName = "null"
	default:
		typeName = "unknown"
	}
	return ToNodeValues(typeName), nil
}

func (q *TypeQuery) String() string {
	return "type()"
}

// HasQuery checks if the node has the specified key.
// Works with both arrays (numeric keys) and maps (string keys).
type HasQuery struct {
	Key string
}

// NewHasQuery creates a new HasQuery instance.
// Requires exactly one argument specifying the key to check.
func NewHasQuery(args ...string) (Query, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("has() requires exactly 1 argument, got %d", len(args))
	}
	return &HasQuery{Key: args[0]}, nil
}

// Exec returns whether the node has the specified key.
// For arrays, checks numeric indices. For maps, checks string keys.
func (q *HasQuery) Exec(n Node) ([]Node, error) {
	switch n.Type() {
	case TypeArray:
		if idx, err := strconv.Atoi(q.Key); err == nil {
			return ToNodeValues(n.Has(idx)), nil
		}
		return ToNodeValues(false), nil
	case TypeMap:
		return ToNodeValues(n.Has(q.Key)), nil
	}
	return ToNodeValues(false), nil
}

func (q *HasQuery) String() string {
	return fmt.Sprintf("has(%q)", q.Key)
}

// FirstQuery returns the first element of an array.
// Returns null for empty arrays or non-array types.
type FirstQuery struct{}

// NewFirstQuery creates a new FirstQuery instance.
// Arguments are ignored for this query type.
func NewFirstQuery(args ...string) (Query, error) {
	return &FirstQuery{}, nil
}

// Exec returns the first element of an array.
// Returns null for empty arrays or non-array types.
func (q *FirstQuery) Exec(n Node) ([]Node, error) {
	if a := n.Array(); a != nil {
		if len(a) > 0 {
			return []Node{a[0]}, nil
		}
	}
	return []Node{Nil}, nil
}

func (q *FirstQuery) String() string {
	return "first()"
}

// LastQuery returns the last element of an array.
// Returns null for empty arrays or non-array types.
type LastQuery struct{}

// NewLastQuery creates a new LastQuery instance.
// Arguments are ignored for this query type.
func NewLastQuery(args ...string) (Query, error) {
	return &LastQuery{}, nil
}

// Exec returns the last element of an array.
// Returns null for empty arrays or non-array types.
func (q *LastQuery) Exec(n Node) ([]Node, error) {
	if a := n.Array(); a != nil {
		if len(a) > 0 {
			return []Node{a[len(a)-1]}, nil
		}
	}
	return []Node{Nil}, nil
}

func (q *LastQuery) String() string {
	return "last()"
}

// FlattenQuery flattens nested arrays into a single array.
// Only flattens one level deep by default.
type FlattenQuery struct{}

// NewFlattenQuery creates a new FlattenQuery instance.
// Arguments are ignored for this query type.
func NewFlattenQuery(args ...string) (Query, error) {
	return &FlattenQuery{}, nil
}

// Exec flattens nested arrays into a single array.
// Only flattens one level deep. Returns the original node for non-arrays.
func (q *FlattenQuery) Exec(n Node) ([]Node, error) {
	if a := n.Array(); a != nil {
		var flattened Array
		for _, item := range a {
			if subArray := item.Array(); subArray != nil {
				flattened = append(flattened, subArray...)
			} else {
				flattened = append(flattened, item)
			}
		}
		return []Node{flattened}, nil
	}
	return []Node{n}, nil
}

func (q *FlattenQuery) String() string {
	return "flatten()"
}

// init automatically registers the built-in method queries.
func init() {
	RegisterNewMethodQueryFunc("count", NewCountQuery)
	RegisterNewMethodQueryFunc("keys", NewKeysQuery)
	RegisterNewMethodQueryFunc("values", NewValuesQuery)
	RegisterNewMethodQueryFunc("empty", NewEmptyQuery)
	RegisterNewMethodQueryFunc("type", NewTypeQuery)
	RegisterNewMethodQueryFunc("has", NewHasQuery)
	RegisterNewMethodQueryFunc("first", NewFirstQuery)
	RegisterNewMethodQueryFunc("last", NewLastQuery)
	RegisterNewMethodQueryFunc("flatten", NewFlattenQuery)
}
