package tree

// Type represents the Node type.
type Type int

// These variables are the Node types.
const (
	TypeArray = iota
	TypeMap
	TypeValue
)

// Node provides
type Node interface {
	Type() Type
	Array() Array
	Map() Map
	Value() Value
}

// Array represents an array of Node.
type Array []Node

var _ Node = (Array)(nil)

func (n Array) Type() Type {
	return TypeArray
}

func (n Array) Array() Array {
	return n
}

func (n Array) Map() Map {
	return nil
}

func (n Array) Value() Value {
	return nil
}

// Map represents a map of Node.
type Map map[string]Node

var _ Node = (Map)(nil)

func (n Map) Type() Type {
	return TypeMap
}

func (n Map) Array() Array {
	return nil
}

func (n Map) Map() Map {
	return n
}

func (n Map) Value() Value {
	return nil
}
