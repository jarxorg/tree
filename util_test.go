package tree

import (
	"reflect"
	"testing"
)

func TestToValue(t *testing.T) {
	tests := []struct {
		v    interface{}
		want Node
	}{
		{
			v:    nil,
			want: nil,
		}, {
			v:    "string",
			want: StringValue("string"),
		}, {
			v:    true,
			want: BoolValue(true),
		}, {
			v:    1,
			want: NumberValue(1),
		}, {
			v:    int64(2),
			want: NumberValue(2),
		}, {
			v:    int32(3),
			want: NumberValue(3),
		}, {
			v:    float64(4.4),
			want: NumberValue(4.4),
		}, {
			v:    float32(5.5),
			want: NumberValue(5.5),
		}, {
			v:    uint64(6),
			want: NumberValue(uint64(6)),
		}, {
			v:    uint32(7),
			want: NumberValue(uint32(7)),
		}, {
			v:    BoolValue(true),
			want: BoolValue(true),
		}, {
			v:    struct{}{},
			want: StringValue("struct {}{}"),
		},
	}
	for _, test := range tests {
		got := ToValue(test.v)
		if got != test.want {
			t.Errorf(`Error ToValue(%#v) returns %#v; want %#v`, test.v, got, test.want)
		}
	}
}

func TestToNode(t *testing.T) {
	tests := []struct {
		v    interface{}
		want Node
	}{
		{
			v:    nil,
			want: nil,
		}, {
			v:    StringValue("a"),
			want: StringValue("a"),
		}, {
			v:    map[string]interface{}{"a": 1, "b": true},
			want: Map{"a": NumberValue(1), "b": BoolValue(true)},
		}, {
			v:    []interface{}{"a", true, 1},
			want: Array{StringValue("a"), BoolValue(true), NumberValue(1)},
		},
	}
	for _, test := range tests {
		got := ToNode(test.v)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf(`Error %v ToNode %v; want %v`, test.v, got, test.want)
		}
	}
}
