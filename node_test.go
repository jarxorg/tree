package tree

import (
	"reflect"
	"testing"
)

func TestType(t *testing.T) {
	tests := []struct {
		typ  Type
		is   func() bool
		want bool
	}{
		{typ: TypeArray, is: TypeArray.IsArray, want: true},
		{typ: TypeArray, is: TypeArray.IsMap, want: false},
		{typ: TypeArray, is: TypeArray.IsValue, want: false},
		{typ: TypeMap, is: TypeMap.IsArray, want: false},
		{typ: TypeMap, is: TypeMap.IsMap, want: true},
		{typ: TypeMap, is: TypeMap.IsValue, want: false},
		{typ: TypeValue, is: TypeValue.IsArray, want: false},
		{typ: TypeValue, is: TypeValue.IsMap, want: false},
		{typ: TypeValue, is: TypeValue.IsValue, want: true},
		{typ: TypeStringValue, is: TypeStringValue.IsArray, want: false},
		{typ: TypeStringValue, is: TypeStringValue.IsMap, want: false},
		{typ: TypeStringValue, is: TypeStringValue.IsValue, want: true},
		{typ: TypeStringValue, is: TypeStringValue.IsStringValue, want: true},
		{typ: TypeStringValue, is: TypeStringValue.IsBoolValue, want: false},
		{typ: TypeStringValue, is: TypeStringValue.IsNumberValue, want: false},
		{typ: TypeBoolValue, is: TypeBoolValue.IsArray, want: false},
		{typ: TypeBoolValue, is: TypeBoolValue.IsMap, want: false},
		{typ: TypeBoolValue, is: TypeBoolValue.IsValue, want: true},
		{typ: TypeBoolValue, is: TypeBoolValue.IsStringValue, want: false},
		{typ: TypeBoolValue, is: TypeBoolValue.IsBoolValue, want: true},
		{typ: TypeBoolValue, is: TypeBoolValue.IsNumberValue, want: false},
		{typ: TypeNumberValue, is: TypeNumberValue.IsArray, want: false},
		{typ: TypeNumberValue, is: TypeNumberValue.IsMap, want: false},
		{typ: TypeNumberValue, is: TypeNumberValue.IsValue, want: true},
		{typ: TypeNumberValue, is: TypeNumberValue.IsStringValue, want: false},
		{typ: TypeNumberValue, is: TypeNumberValue.IsBoolValue, want: false},
		{typ: TypeNumberValue, is: TypeNumberValue.IsNumberValue, want: true},
	}
	for i, test := range tests {
		if got := test.is(); got != test.want {
			t.Errorf(`Error tests[%d] is %v; want %v`, i, got, test.want)
		}
	}
}

func TestNode(t *testing.T) {
	a := Array{}
	m := Map{}
	tests := []struct {
		n Node
		t Type
		a Array
		m Map
		v Value
	}{
		{
			n: m,
			t: TypeMap,
			m: m,
		}, {
			n: a,
			t: TypeArray,
			a: a,
		},
	}
	for _, test := range tests {
		n := test.n
		if tt := n.Type(); tt != test.t {
			t.Errorf(`Error Type returns %v; want %v`, tt, test.t)
		}
		if aa := n.Array(); !reflect.DeepEqual(aa, test.a) {
			t.Errorf(`Error Array returns %v; want %v`, aa, test.a)
		}
		if mm := n.Map(); !reflect.DeepEqual(mm, test.m) {
			t.Errorf(`Error Map returns %v; want %v`, mm, test.m)
		}
		if vv := n.Value(); !reflect.DeepEqual(vv, test.v) {
			t.Errorf(`Error Value returns %v; want %v`, vv, test.v)
		}
	}
}
