package tree

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_Type(t *testing.T) {
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

func Test_Node(t *testing.T) {
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
	for i, test := range tests {
		n := test.n
		if tt := n.Type(); tt != test.t {
			t.Errorf(`Error tests[%d] Type returns %v; want %v`, i, tt, test.t)
		}
		if aa := n.Array(); !reflect.DeepEqual(aa, test.a) {
			t.Errorf(`Error tests[%d] Array returns %v; want %v`, i, aa, test.a)
		}
		if mm := n.Map(); !reflect.DeepEqual(mm, test.m) {
			t.Errorf(`Error tests[%d] Map returns %v; want %v`, i, mm, test.m)
		}
		if vv := n.Value(); !reflect.DeepEqual(vv, test.v) {
			t.Errorf(`Error tests[%d] Value returns %v; want %v`, i, vv, test.v)
		}
	}
}

func Test_Node_Get(t *testing.T) {
	tests := []struct {
		n    Node
		key  interface{}
		has  bool
		want Node
	}{
		{
			n:    Array{StringValue("a"), StringValue("b")},
			key:  1,
			has:  true,
			want: StringValue("b"),
		}, {
			n:    Array{StringValue("a"), StringValue("b")},
			key:  "1",
			has:  true,
			want: StringValue("b"),
		}, {
			n:   Array{StringValue("a"), StringValue("b")},
			key: 1.0,
		}, {
			n:   Array{StringValue("a"), StringValue("b")},
			key: 2,
		}, {
			n:   Array{StringValue("a"), nil},
			key: 1,
			has: true,
		}, {
			n:    Map{"1": NumberValue(10), "2": NumberValue(20)},
			key:  "1",
			has:  true,
			want: NumberValue(10),
		}, {
			n:    Map{"1": NumberValue(10), "2": NumberValue(20)},
			key:  1,
			has:  true,
			want: NumberValue(10),
		}, {
			n:   Map{"1": NumberValue(10), "2": NumberValue(20)},
			key: 1.0,
		}, {
			n:   Map{"1": NumberValue(10), "2": NumberValue(20)},
			key: "3",
		}, {
			n:   Map{"1": NumberValue(10), "2": nil},
			key: "2",
			has: true,
		}, {
			n: StringValue("str"),
		}, {
			n: BoolValue(true),
		}, {
			n: NumberValue(1),
		},
	}
	for i, test := range tests {
		if test.n.Has(test.key) != test.has {
			t.Errorf("Error tests[%d] has %v; want %v", i, !test.has, test.has)
		}
		got := test.n.Get(test.key)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Error tests[%d] got %q; want %q", i, got, test.want)
		}
	}
}

func Test_Node_Each(t *testing.T) {
	tests := []struct {
		n    Node
		want map[interface{}]Node
	}{
		{
			n:    Array{StringValue("a"), StringValue("b")},
			want: map[interface{}]Node{0: StringValue("a"), 1: StringValue("b")},
		}, {
			n:    Map{"a": NumberValue(0), "b": NumberValue(1)},
			want: map[interface{}]Node{"a": NumberValue(0), "b": NumberValue(1)},
		}, {
			n:    StringValue("str"),
			want: map[interface{}]Node{nil: StringValue("str")},
		}, {
			n:    BoolValue(true),
			want: map[interface{}]Node{nil: BoolValue(true)},
		}, {
			n:    NumberValue(1),
			want: map[interface{}]Node{nil: NumberValue(1)},
		},
	}
	for i, test := range tests {
		got := map[interface{}]Node{}
		err := test.n.Each(func(key interface{}, v Node) error {
			got[key] = v
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf(`Error tests[%d] got %v; want %v`, i, got, test.want)
		}
		wantErr := fmt.Errorf("test%d", i)
		gotErr := test.n.Each(func(key interface{}, v Node) error {
			return wantErr
		})
		if wantErr != gotErr {
			t.Errorf(`Error tests[%d] got error %v; want %v`, i, gotErr, wantErr)
		}
	}
}

func Test_Node_Find(t *testing.T) {
	tests := []struct {
		n    Node
		expr string
		want []Node
	}{
		{
			n:    Array{StringValue("a"), StringValue("b")},
			expr: ".[0]",
			want: []Node{StringValue("a")},
		}, {
			n:    Map{"1": NumberValue(10), "2": NumberValue(20)},
			expr: ".1",
			want: []Node{NumberValue(10)},
		},
	}
	for i, test := range tests {
		got, err := test.n.Find(test.expr)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Error tests[%d] returns %#v; want %#v", i, got, test.want)
		}
	}
}

func Test_EditorNode_Append(t *testing.T) {
	tests := []struct {
		n      EditorNode
		values []Node
		want   EditorNode
		errstr string
	}{
		{
			n:      &Array{NumberValue(1)},
			values: []Node{StringValue("2"), BoolValue(true)},
			want:   &Array{NumberValue(1), StringValue("2"), BoolValue(true)},
		}, {
			n:      Map{},
			values: []Node{StringValue("2")},
			errstr: "Cannot append to map",
		},
	}
	for i, test := range tests {
		var err error
		for _, value := range test.values {
			err = test.n.Append(value)
			if err != nil {
				break
			}
		}
		if test.errstr != "" {
			if err == nil {
				t.Fatalf("Fatal tests[%d] returns no error", i)
			}
			if err.Error() != test.errstr {
				t.Errorf(`Error tests[%d] returns error %s; want %s`, i, err.Error(), test.errstr)
			}
			continue
		}
		if err != nil {
			t.Fatal(err, i)
		}
		got := test.n
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf(`Error tests[%d] returns %v; want %v`, i, got, test.want)
		}
	}
}

func Test_EditorNode_Set(t *testing.T) {
	tests := []struct {
		n       EditorNode
		entries map[interface{}]Node
		want    EditorNode
		errstr  string
	}{
		{
			n: &Array{NumberValue(0), StringValue("1")},
			entries: map[interface{}]Node{
				0:   NumberValue(1),
				"1": StringValue("2"),
				2:   BoolValue(true),
			},
			want: &Array{NumberValue(1), StringValue("2"), BoolValue(true)},
		}, {
			n:       &Array{},
			entries: map[interface{}]Node{-2: StringValue("value")},
			errstr:  "Cannot index array with -2",
		}, {
			n: Map{
				"1": NumberValue(1),
				"2": StringValue("2"),
				"3": BoolValue(true),
			},
			entries: map[interface{}]Node{
				"1": NumberValue(10),
				"4": StringValue("40"),
				5:   BoolValue(true),
			},
			want: Map{
				"1": NumberValue(10),
				"2": StringValue("2"),
				"3": BoolValue(true),
				"4": StringValue("40"),
				"5": BoolValue(true),
			},
		}, {
			n:       Map{},
			entries: map[interface{}]Node{true: StringValue("value")},
			errstr:  "Cannot index array with true",
		},
	}
	for i, test := range tests {
		var err error
		for key, value := range test.entries {
			err = test.n.Set(key, value)
			if err != nil {
				break
			}
		}
		if test.errstr != "" {
			if err == nil {
				t.Fatalf("Fatal tests[%d] returns no error", i)
			}
			if err.Error() != test.errstr {
				t.Errorf(`Error tests[%d] returns error %s; want %s`, i, err.Error(), test.errstr)
			}
			continue
		}
		if err != nil {
			t.Fatal(err, i)
		}
		got := test.n
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf(`Error tests[%d] returns %v; want %v`, i, got, test.want)
		}
	}
}

func Test_EditorNode_Delete(t *testing.T) {
	tests := []struct {
		n      EditorNode
		keys   []interface{}
		want   EditorNode
		errstr string
	}{
		{
			n:    &Array{NumberValue(1), StringValue("1"), BoolValue(true)},
			keys: []interface{}{1, "1"},
			want: &Array{NumberValue(1)},
		}, {
			n:      &Array{},
			keys:   []interface{}{-1},
			errstr: "Cannot index array with -1",
		}, {
			n: Map{
				"1": NumberValue(1),
				"2": StringValue("2"),
				"3": BoolValue(true),
				"4": StringValue("4"),
				"5": BoolValue(true),
			},
			keys: []interface{}{"2", "4", 5, 7},
			want: Map{
				"1": NumberValue(1),
				"3": BoolValue(true),
			},
		}, {
			n:      Map{},
			keys:   []interface{}{true},
			errstr: "Cannot index array with true",
		},
	}
	for i, test := range tests {
		var err error
		for _, key := range test.keys {
			err = test.n.Delete(key)
			if err != nil {
				break
			}
		}
		if test.errstr != "" {
			if err == nil {
				t.Fatalf("Fatal tests[%d] returns no error", i)
			}
			if err.Error() != test.errstr {
				t.Errorf(`Error tests[%d] returns error %s; want %s`, i, err.Error(), test.errstr)
			}
			continue
		}
		if err != nil {
			t.Fatal(err, i)
		}
		got := test.n
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf(`Error tests[%d] returns %v; want %v`, i, got, test.want)
		}
	}
}
