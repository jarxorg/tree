package tree

import (
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
	// Get returns array/map value that matched by the specified key.
	// The key type allows int or string.
	Get(key interface{}) Node
	// Each calls the callback function for each Array|Map values.
	// If the node type is not Array|Map then the callback called once with nil key and self as value.
	Each(cb func(key interface{}, v Node) error) error
}

// Array represents an array of Node.
type Array []Node

var _ Node = (Array)(nil)

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

// Get returns an array value as Node.
func (n Array) Get(key interface{}) Node {
	switch key.(type) {
	case int:
		if k := key.(int); k >= 0 && k < len(n) {
			return n[k]
		}
	case string:
		k, err := strconv.Atoi(key.(string))
		if err == nil && k >= 0 && k < len(n) {
			return n[k]
		}
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

// Map represents a map of Node.
type Map map[string]Node

var _ Node = (Map)(nil)

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

// Get returns an array value as Node.
func (n Map) Get(key interface{}) Node {
	switch key.(type) {
	case int:
		return n[strconv.Itoa(key.(int))]
	case string:
		return n[key.(string)]
	}
	return nil
}

// Each calls the callback function for each Map values.
func (n Map) Each(cb func(key interface{}, n Node) error) error {
	for k, v := range n {
		if err := cb(k, v); err != nil {
			return err
		}
	}
	return nil
}
