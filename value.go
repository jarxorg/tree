package tree

import (
	"strconv"
)

// Operator represents an operator.
type Operator string

var (
	// EQ is `==`.
	EQ Operator = "=="
	// GT is `>`.
	GT Operator = ">"
	// GE is `>=`.
	GE Operator = ">="
	// LT is `<`.
	LT Operator = "<"
	// LE is `<=`.
	LE Operator = "<="
)

// Value provides the accessor of primitive value.
type Value interface {
	Node
	String() string
	Bool() bool
	Int() int
	Int64() int64
	Float64() float64
	Compare(op Operator, v Value) bool
}

// A StringValue represents a string value.
type StringValue string

var _ Node = StringValue("")

// Type returns TypeValue.
func (n StringValue) Type() Type {
	return TypeStringValue
}

// Array returns nil.
func (n StringValue) Array() Array {
	return nil
}

// Map returns nil.
func (n StringValue) Map() Map {
	return nil
}

// Value returns this.
func (n StringValue) Value() Value {
	return n
}

// Has returns false.
func (n StringValue) Has(key interface{}) bool {
	return false
}

// Get returns nil.
func (n StringValue) Get(key interface{}) Node {
	return nil
}

// Each calls cb(nil, n).
func (n StringValue) Each(cb func(key interface{}, n Node) error) error {
	return cb(nil, n)
}

// Find finds a node using the query expression.
func (n StringValue) Find(expr string) ([]Node, error) {
	return Find(n, expr)
}

// Bool returns false.
func (n StringValue) Bool() bool {
	return false
}

// Int returns 0.
func (n StringValue) Int() int {
	return 0
}

// Int64 returns 0.
func (n StringValue) Int64() int64 {
	return 0
}

// Float64 returns 0.
func (n StringValue) Float64() float64 {
	return 0
}

// String returns this as string.
func (n StringValue) String() string {
	return string(n)
}

// Compare compares n and v.
func (n StringValue) Compare(op Operator, v Value) bool {
	if v == nil || !v.Type().IsStringValue() {
		return false
	}
	sn := n.String()
	sv := v.String()
	switch op {
	case EQ:
		return sn == sv
	case GT:
		return sn > sv
	case GE:
		return sn >= sv
	case LT:
		return sn < sv
	case LE:
		return sn <= sv
	}
	return false
}

// A BoolValue represents a bool value.
type BoolValue bool

var _ Node = BoolValue(false)

// Type returns TypeValue.
func (n BoolValue) Type() Type {
	return TypeBoolValue
}

// Array returns nil.
func (n BoolValue) Array() Array {
	return nil
}

// Map returns nil.
func (n BoolValue) Map() Map {
	return nil
}

// Value returns this.
func (n BoolValue) Value() Value {
	return n
}

// Has returns false.
func (n BoolValue) Has(key interface{}) bool {
	return false
}

// Get returns nil.
func (n BoolValue) Get(key interface{}) Node {
	return nil
}

// Each calls cb(nil, n).
func (n BoolValue) Each(cb func(key interface{}, n Node) error) error {
	return cb(nil, n)
}

// Find finds a node using the query expression.
func (n BoolValue) Find(expr string) ([]Node, error) {
	return Find(n, expr)
}

// Bool returns this.
func (n BoolValue) Bool() bool {
	return bool(n)
}

// Int returns 0.
func (n BoolValue) Int() int {
	return 0
}

// Int64 returns 0.
func (n BoolValue) Int64() int64 {
	return 0
}

// Float64 returns 0.
func (n BoolValue) Float64() float64 {
	return 0
}

// String returns this as string.
func (n BoolValue) String() string {
	return strconv.FormatBool(bool(n))
}

// Compare compares n and v.
func (n BoolValue) Compare(op Operator, v Value) bool {
	if v == nil || !v.Type().IsBoolValue() {
		return false
	}
	switch op {
	case EQ:
		return n.Bool() == v.Bool()
	}
	return false
}

// A NumberValue represents an number value.
type NumberValue float64

var _ Node = NumberValue(0)

// Type returns TypeValue.
func (n NumberValue) Type() Type {
	return TypeNumberValue
}

// Array returns nil.
func (n NumberValue) Array() Array {
	return nil
}

// Map returns nil.
func (n NumberValue) Map() Map {
	return nil
}

// Value returns this.
func (n NumberValue) Value() Value {
	return n
}

// Has returns false.
func (n NumberValue) Has(key interface{}) bool {
	return false
}

// Get returns nil.
func (n NumberValue) Get(key interface{}) Node {
	return nil
}

// Each calls cb(nil, n).
func (n NumberValue) Each(cb func(key interface{}, n Node) error) error {
	return cb(nil, n)
}

// Find finds a node using the query expression.
func (n NumberValue) Find(expr string) ([]Node, error) {
	return Find(n, expr)
}

// Bool returns false.
func (n NumberValue) Bool() bool {
	return false
}

// Int returns int(n).
func (n NumberValue) Int() int {
	return int(n)
}

// Int64 returns int64(n).
func (n NumberValue) Int64() int64 {
	return int64(n)
}

// Float64 returns float64(n).
func (n NumberValue) Float64() float64 {
	return float64(n)
}

// String returns this as string using strconv.FormatFloat(float64(n), 'f', -1, 64).
func (n NumberValue) String() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

// Compare compares n and v.
func (n NumberValue) Compare(op Operator, v Value) bool {
	if v == nil || !v.Type().IsNumberValue() {
		return false
	}
	nv := v.Float64()
	nn := n.Float64()
	switch op {
	case EQ:
		return nn == nv
	case GT:
		return nn > nv
	case GE:
		return nn >= nv
	case LT:
		return nn < nv
	case LE:
		return nn <= nv
	}
	return false
}
