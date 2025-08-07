package tree

import (
	"reflect"
	"testing"
)

func Test_MethodQuery(t *testing.T) {
	testCases := []struct {
		caseName string
		q        Query
		n        Node
		want     []Node
		errstr   string
	}{
		{
			caseName: "count map",
			q:        &CountQuery{},
			n:        Map{"key1": ToValue(1), "key2": ToValue("a")},
			want:     []Node{ToValue(2)},
		}, {
			caseName: "count array",
			q:        &CountQuery{},
			n:        ToArrayValues(1, 2, 3),
			want:     []Node{ToValue(3)},
		}, {
			caseName: "keys map",
			q:        &KeysQuery{},
			n:        Map{"key1": ToValue(1), "key2": ToValue("a")},
			want:     []Node{ToArrayValues("key1", "key2")},
		}, {
			caseName: "keys array",
			q:        &KeysQuery{},
			n:        ToArrayValues(1, 2, 3),
			want:     []Node{ToArrayValues(0, 1, 2)},
		}, {
			caseName: "values map",
			q:        &ValuesQuery{},
			n:        Map{"key1": ToValue(1), "key2": ToValue("a")},
			want:     []Node{ToArrayValues(1, "a")},
		}, {
			caseName: "values array",
			q:        &ValuesQuery{},
			n:        ToArrayValues(1, 2, 3),
			want:     []Node{ToArrayValues(1, 2, 3)},
		}, {
			caseName: "empty array true",
			q:        &EmptyQuery{},
			n:        Array{},
			want:     []Node{BoolValue(true)},
		}, {
			caseName: "empty array false",
			q:        &EmptyQuery{},
			n:        ToArrayValues(1, 2),
			want:     []Node{BoolValue(false)},
		}, {
			caseName: "empty map true",
			q:        &EmptyQuery{},
			n:        Map{},
			want:     []Node{BoolValue(true)},
		}, {
			caseName: "empty string true",
			q:        &EmptyQuery{},
			n:        StringValue(""),
			want:     []Node{BoolValue(true)},
		}, {
			caseName: "empty string false",
			q:        &EmptyQuery{},
			n:        StringValue("test"),
			want:     []Node{BoolValue(false)},
		}, {
			caseName: "empty nil",
			q:        &EmptyQuery{},
			n:        Nil,
			want:     []Node{BoolValue(true)},
		}, {
			caseName: "type array",
			q:        &TypeQuery{},
			n:        ToArrayValues(1, 2),
			want:     []Node{StringValue("array")},
		}, {
			caseName: "type object",
			q:        &TypeQuery{},
			n:        Map{"key": ToValue("value")},
			want:     []Node{StringValue("object")},
		}, {
			caseName: "type string",
			q:        &TypeQuery{},
			n:        StringValue("test"),
			want:     []Node{StringValue("string")},
		}, {
			caseName: "type number",
			q:        &TypeQuery{},
			n:        NumberValue(42),
			want:     []Node{StringValue("number")},
		}, {
			caseName: "type boolean",
			q:        &TypeQuery{},
			n:        BoolValue(true),
			want:     []Node{StringValue("boolean")},
		}, {
			caseName: "type null",
			q:        &TypeQuery{},
			n:        Nil,
			want:     []Node{StringValue("null")},
		}, {
			caseName: "has key in map true",
			q:        &HasQuery{Key: "name"},
			n:        Map{"name": StringValue("test"), "age": NumberValue(30)},
			want:     []Node{BoolValue(true)},
		}, {
			caseName: "has key in map false",
			q:        &HasQuery{Key: "missing"},
			n:        Map{"name": StringValue("test")},
			want:     []Node{BoolValue(false)},
		}, {
			caseName: "has index in array true",
			q:        &HasQuery{Key: "1"},
			n:        ToArrayValues("a", "b", "c"),
			want:     []Node{BoolValue(true)},
		}, {
			caseName: "has index in array false",
			q:        &HasQuery{Key: "5"},
			n:        ToArrayValues("a", "b"),
			want:     []Node{BoolValue(false)},
		}, {
			caseName: "first element",
			q:        &FirstQuery{},
			n:        ToArrayValues("a", "b", "c"),
			want:     []Node{StringValue("a")},
		}, {
			caseName: "first empty array",
			q:        &FirstQuery{},
			n:        Array{},
			want:     []Node{Nil},
		}, {
			caseName: "first non-array",
			q:        &FirstQuery{},
			n:        StringValue("test"),
			want:     []Node{Nil},
		}, {
			caseName: "last element",
			q:        &LastQuery{},
			n:        ToArrayValues("a", "b", "c"),
			want:     []Node{StringValue("c")},
		}, {
			caseName: "last empty array",
			q:        &LastQuery{},
			n:        Array{},
			want:     []Node{Nil},
		}, {
			caseName: "flatten nested arrays",
			q:        &FlattenQuery{},
			n:        Array{ToArrayValues(1, 2), ToArrayValues(3, 4), NumberValue(5)},
			want:     []Node{ToArrayValues(1, 2, 3, 4, 5)},
		}, {
			caseName: "flatten non-array",
			q:        &FlattenQuery{},
			n:        StringValue("test"),
			want:     []Node{StringValue("test")},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.caseName, func(t *testing.T) {
			got, err := tc.q.Exec(tc.n)
			if tc.errstr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tc.errstr)
				}
				if err.Error() != tc.errstr {
					t.Errorf("got error %q; want %q", err.Error(), tc.errstr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v; want %v", got, tc.want)
			}
		})
	}
}
