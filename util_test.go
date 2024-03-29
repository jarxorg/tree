package tree

import (
	"reflect"
	"testing"
)

func Test_ToValue(t *testing.T) {
	tests := []struct {
		v    interface{}
		want Node
	}{
		{
			v:    nil,
			want: Nil,
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
	for i, test := range tests {
		got := ToValue(test.v)
		if got != test.want {
			t.Errorf("tests[%d] for %v; got %#v; want %#v", i, test.v, got, test.want)
		}
	}
}

func Test_ToNode(t *testing.T) {
	tests := []struct {
		v    interface{}
		want Node
	}{
		{
			v:    nil,
			want: Nil,
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
	for i, test := range tests {
		got := ToNode(test.v)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("tests[%d] for %v; got %v; want %v", i, test.v, got, test.want)
		}
	}
}

func Test_Walk(t *testing.T) {
	root := Array{
		Map{"ID": ToValue(1)},
		Map{"ID": ToValue(2), "Sub": Array{Map{"ID": ToValue(20)}}},
		Map{"ID": ToValue(3), "Sub": Array{Map{"ID": ToValue(30)}}},
	}

	tests := []struct {
		n    Node
		keys []interface{}
		skip bool
	}{
		{
			n:    root,
			keys: []interface{}{},
		}, {
			n:    root.Get(0),
			keys: []interface{}{0},
		}, {
			n:    root.Get(0).Get("ID"),
			keys: []interface{}{0, "ID"},
		}, {
			n:    root.Get(1),
			keys: []interface{}{1},
			skip: true,
		}, {
			n:    root.Get(2),
			keys: []interface{}{2},
		}, {
			n:    root.Get(2).Get("ID"),
			keys: []interface{}{2, "ID"},
		}, {
			n:    root.Get(2).Get("Sub"),
			keys: []interface{}{2, "Sub"},
		}, {
			n:    root.Get(2).Get("Sub").Get(0),
			keys: []interface{}{2, "Sub", 0},
		}, {
			n:    root.Get(2).Get("Sub").Get(0).Get("ID"),
			keys: []interface{}{2, "Sub", 0, "ID"},
		},
	}

	i := 0
	err := Walk(root, func(n Node, keys []interface{}) error {
		if i >= len(tests) {
			t.Fatalf("fn is called too many times %d", i)
			return nil
		}
		test := tests[i]
		i++

		if !reflect.DeepEqual(n, test.n) {
			t.Errorf("walk[%d] got %#v; want %#v", i, n, test.n)
		}
		if !reflect.DeepEqual(keys, test.keys) {
			t.Errorf("walk[%d] got %#v; want %#v", i, keys, test.n)
		}
		if test.skip {
			return SkipWalk
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(tests) != i {
		t.Errorf("fn is called %d times; want %d", i, len(tests))
	}
}

func Test_regexpMatchString(t *testing.T) {
	tests := []struct {
		expr   string
		value  string
		want   bool
		errstr string
	}{
		{
			expr:  `a`,
			value: "abc",
			want:  true,
		}, {
			expr:  `^[a-z]+$`,
			value: "abc",
			want:  true,
		}, {
			expr:  `x`,
			value: "abc",
			want:  false,
		}, {
			expr:   `(`,
			value:  "abc",
			errstr: "error parsing regexp: missing closing ): `(`",
		},
	}
	for i, test := range tests {
		got, err := regexpMatchString(test.expr, test.value)
		if test.errstr != "" {
			if err == nil {
				t.Fatalf("tests[%d] for %v; no error", i, test.expr)
			}
			if err.Error() != test.errstr {
				t.Errorf(`tests[%d] for %v; got %v want %v`, i, test.expr, err.Error(), test.errstr)
			}
			continue
		}
		if err != nil {
			t.Fatalf("tests[%d] for %v; %+v", i, test.expr, err)
		}
		if got != test.want {
			t.Errorf("tests[%d] for %v; got %v; want %v", i, test.expr, got, test.want)
		}
	}
}

func TestClone(t *testing.T) {
	tests := []struct {
		n    Node
		want Node
	}{
		{
			n:    ToValue(1),
			want: ToValue(1),
		}, {
			n:    Map{"a": ToValue(1), "b": ToValue(2)},
			want: Map{"a": ToValue(1), "b": ToValue(2)},
		},
	}
	for i, test := range tests {
		got := Clone(test.n)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf(`tests[%d]: unexpected %v; want %v`, i, got, test.want)
		}
	}
}

func TestCloneDeep(t *testing.T) {
	tests := []struct {
		n      Node
		want   Node
		update func(n Node)
	}{
		{
			n:      ToValue(1),
			want:   ToValue(1),
			update: func(n Node) {},
		}, {
			n:    Map{"a": ToArrayValues(1, 2), "b": ToArrayValues(3, 4)},
			want: Map{"a": ToArrayValues(1, 2), "b": ToArrayValues(3, 4)},
			update: func(n Node) {
				n.Map().Get("a").Array()[0] = ToValue(5)
			},
		},
	}
	for i, test := range tests {
		got := CloneDeep(test.n)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf(`tests[%d]: unexpected %v; want %v`, i, got, test.want)
		}
		test.update(test.n)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf(`tests[%d]: unexpected %v; want %v`, i, got, test.want)
		}
	}
}
