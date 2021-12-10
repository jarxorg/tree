package tree

import (
	"strconv"
)

type Value interface {
	String() string
	Bool() bool
	Int() int
	Int64() int64
	Float64() float64
}

type StringValue string

var _ Node = StringValue("")

func (n StringValue) Type() Type {
	return TypeValue
}

func (n StringValue) Array() Array {
	return nil
}

func (n StringValue) Map() Map {
	return nil
}

func (n StringValue) Value() Value {
	return n
}

func (n StringValue) Bool() bool {
	return false
}

func (n StringValue) Int() int {
	return 0
}

func (n StringValue) Int64() int64 {
	return 0
}

func (n StringValue) Float64() float64 {
	return 0
}

func (n StringValue) String() string {
	return string(n)
}

type BoolValue bool

var _ Node = BoolValue(false)

func (n BoolValue) Type() Type {
	return TypeValue
}

func (n BoolValue) Array() Array {
	return nil
}

func (n BoolValue) Map() Map {
	return nil
}

func (n BoolValue) Value() Value {
	return n
}

func (n BoolValue) Bool() bool {
	return bool(n)
}

func (n BoolValue) Int() int {
	return 0
}

func (n BoolValue) Int64() int64 {
	return 0
}

func (n BoolValue) Float64() float64 {
	return 0
}

func (n BoolValue) String() string {
	return strconv.FormatBool(bool(n))
}

type Int64Value int64

var _ Node = Int64Value(0)

func (n Int64Value) Type() Type {
	return TypeValue
}

func (n Int64Value) Array() Array {
	return nil
}

func (n Int64Value) Map() Map {
	return nil
}

func (n Int64Value) Value() Value {
	return n
}

func (n Int64Value) Bool() bool {
	return false
}

func (n Int64Value) Int() int {
	return int(n)
}

func (n Int64Value) Int64() int64 {
	return int64(n)
}

func (n Int64Value) Float64() float64 {
	return float64(n)
}

func (n Int64Value) String() string {
	return strconv.FormatInt(int64(n), 64)
}

type Float64Value float64

var _ Node = Float64Value(0)

func (n Float64Value) Type() Type {
	return TypeValue
}

func (n Float64Value) Array() Array {
	return nil
}

func (n Float64Value) Map() Map {
	return nil
}

func (n Float64Value) Value() Value {
	return n
}

func (n Float64Value) Bool() bool {
	return false
}

func (n Float64Value) Int() int {
	return int(n)
}

func (n Float64Value) Int64() int64 {
	return int64(n)
}

func (n Float64Value) Float64() float64 {
	return float64(n)
}

func (n Float64Value) String() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}
