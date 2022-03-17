package tree

import (
	"fmt"
	"sort"
	"strconv"
)

// Type represents the Node type.
type Type int

// These variables are the Node types.
const (
	TypeArray Type = 1 << (32 - 1 - iota)
	TypeMap
	TypeValue
	TypeStringValue = TypeValue | iota
	TypeBoolValue
	TypeNumberValue
)

// IsArray returns t == TypeArray.
func (t Type) IsArray() bool {
	return t == TypeArray
}

// IsMap returns t == TypeMap.
func (t Type) IsMap() bool {
	return t == TypeMap
}

// IsValue returns true if t is TypeStringValue or TypeBoolValue or TypeNumberValue.
func (t Type) IsValue() bool {
	return t&TypeValue != 0
}

// IsStringValue returns t == TypeStringValue.
func (t Type) IsStringValue() bool {
	return t == TypeStringValue
}

// IsBoolValue returns t == TypeBoolValue.
func (t Type) IsBoolValue() bool {
	return t == TypeBoolValue
}

// IsNumberValue returns t == TypeNumberValue.
func (t Type) IsNumberValue() bool {
	return t == TypeNumberValue
}

// A Node is an element on the tree.
type Node interface {
	// Type returns this node type.
	Type() Type
	// Array returns this node as an Array.
	Array() Array
	// Map returns this node as a Map.
	Map() Map
	// Value returns this node as a Value.
	Value() Value
	// Has checks this node has key.
	Has(key interface{}) bool
	// Get returns array/map value that matched by the specified key.
	// The key type allows int or string.
	Get(key interface{}) Node
	// Each calls the callback function for each Array|Map values.
	// If the node type is not Array|Map then the callback called once with nil key and self as value.
	Each(cb func(key interface{}, v Node) error) error
	// Find finds a node using the query expression.
	Find(expr string) ([]Node, error)
}

// EditorNode is an interface that defines the methods to edit this node.
type EditorNode interface {
	Node
	Append(v Node) error
	Set(key interface{}, v Node) error
	Delete(key interface{}) error
}

// Array represents an array of Node.
type Array []Node

var _ EditorNode = (*Array)(nil)

// Type returns TypeArray.
func (n Array) Type() Type {
	return TypeArray
}

// Array returns this node as an Array.
func (n Array) Array() Array {
	return n
}

// Map returns nil.
func (n Array) Map() Map {
	return nil
}

// Value returns nil.
func (n Array) Value() Value {
	return nil
}

func (n Array) toIndex(key interface{}) (int, bool) {
	switch key.(type) {
	case int:
		if k := key.(int); k >= 0 {
			return k, k < len(n)
		}
	case string:
		k, err := strconv.Atoi(key.(string))
		if err == nil && k >= 0 {
			return k, k < len(n)
		}
	}
	return -1, false
}

// Has checks this node has key.
func (n Array) Has(key interface{}) bool {
	_, ok := n.toIndex(key)
	return ok
}

// Get returns an array value as Node.
func (n Array) Get(key interface{}) Node {
	if i, ok := n.toIndex(key); ok {
		return n[i]
	}
	return nil
}

// Each calls the callback function for each Array values.
func (n Array) Each(cb func(key interface{}, n Node) error) error {
	for i, v := range n {
		if err := cb(i, v); err != nil {
			return err
		}
	}
	return nil
}

// Find finds a node using the query expression.
func (n Array) Find(expr string) ([]Node, error) {
	return Find(n, expr)
}

// Append appends v to *n.
func (n *Array) Append(v Node) error {
	*n = append(*n, v)
	return nil
}

// Set sets v to n[key].
func (n *Array) Set(key interface{}, v Node) error {
	i, ok := n.toIndex(key)
	if i == -1 {
		return fmt.Errorf("Cannot index array with %v", key)
	}
	if !ok {
		a := make([]Node, i+1)
		copy(a, *n)
		*n = a
	}
	(*n)[i] = v
	return nil
}

// Delete deletes n[key].
func (n *Array) Delete(key interface{}) error {
	i, ok := n.toIndex(key)
	if i == -1 {
		return fmt.Errorf("Cannot index array with %v", key)
	}
	if ok {
		a := *n
		*n = append(a[0:i], a[i+1:]...)
	}
	return nil
}

// Map represents a map of Node.
type Map map[string]Node

var _ EditorNode = (Map)(nil)

// Type returns TypeMap.
func (n Map) Type() Type {
	return TypeMap
}

// Array returns nil.
func (n Map) Array() Array {
	return nil
}

// Map returns this node as a Map.
func (n Map) Map() Map {
	return n
}

// Value returns nil.
func (n Map) Value() Value {
	return nil
}

func (n Map) toKey(key interface{}) (string, bool) {
	switch key.(type) {
	case int:
		k := strconv.Itoa(key.(int))
		_, ok := n[k]
		return k, ok
	case string:
		k := key.(string)
		_, ok := n[k]
		return k, ok
	}
	return "", false
}

// Has checks this node has key.
func (n Map) Has(key interface{}) bool {
	_, ok := n.toKey(key)
	return ok
}

// Get returns an array value as Node.
func (n Map) Get(key interface{}) Node {
	if k, ok := n.toKey(key); ok {
		return n[k]
	}
	return nil
}

// Keys returns sorted keys of the map.
func (n Map) Keys() []string {
	keys := make([]string, len(n))
	i := 0
	for k := range n {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

// Values returns values of the map.
func (n Map) Values() []Node {
	values := make([]Node, len(n))
	for i, k := range n.Keys() {
		values[i] = n[k]
	}
	return values
}

// Append appends v to *n.
func (n Map) Append(v Node) error {
	return fmt.Errorf("Cannot append to map")
}

// Set sets v to n[key].
func (n Map) Set(key interface{}, v Node) error {
	switch key.(type) {
	case int:
		n[strconv.Itoa(key.(int))] = v
		return nil
	case string:
		n[key.(string)] = v
		return nil
	}
	return fmt.Errorf("Cannot index array with %v", key)
}

// Delete deletes n[key].
func (n Map) Delete(key interface{}) error {
	switch key.(type) {
	case int:
		delete(n, strconv.Itoa(key.(int)))
		return nil
	case string:
		delete(n, key.(string))
		return nil
	}
	return fmt.Errorf("Cannot index array with %v", key)
}

// Each calls the callback function for each Map values.
func (n Map) Each(cb func(key interface{}, n Node) error) error {
	for _, k := range n.Keys() {
		if err := cb(k, n[k]); err != nil {
			return err
		}
	}
	return nil
}

// Find finds a node using the query expression.
func (n Map) Find(expr string) ([]Node, error) {
	return Find(n, expr)
}
