package tree

import (
	"reflect"
	"testing"
)

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
