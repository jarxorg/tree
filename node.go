package tree

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

// Map returns nil
func (n Array) Map() Map {
	return nil
}

// Value returns nil
func (n Array) Value() Value {
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
