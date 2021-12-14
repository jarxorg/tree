package tree

import "testing"

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
