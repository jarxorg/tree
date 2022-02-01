package tree

import "testing"

func Test_Value(t *testing.T) {
	tests := []struct {
		value Node
		b     bool
		i     int
		i64   int64
		f64   float64
		s     string
	}{
		{
			value: StringValue("test"),
			s:     "test",
		}, {
			value: BoolValue(true),
			b:     true,
			s:     "true",
		}, {
			value: NumberValue(1),
			i:     1,
			i64:   int64(1),
			f64:   float64(1),
			s:     "1",
		}, {
			value: NumberValue(2.3),
			i:     2,
			i64:   int64(2),
			f64:   float64(2.3),
			s:     "2.3",
		},
	}
	for i, test := range tests {
		v := test.value
		vv := v.Value()
		if tt := v.Type(); tt&TypeValue == 0 {
			t.Errorf(`Error tests[%d] Type returns %v; want TypeValue`, i, tt)
		}
		if a := v.Array(); a != nil {
			t.Errorf(`Error tests[%d] Array returns %v; want nil`, i, a)
		}
		if m := v.Map(); m != nil {
			t.Errorf(`Error tests[%d] Map returns %v; want nil`, i, m)
		}
		if vv.(Node) != v {
			t.Errorf(`Error tests[%d] Value returns %v; want %v`, i, vv, v)
		}
		if b := vv.Bool(); b != test.b {
			t.Errorf(`Error tests[%d] Bool returns %v; want %v`, i, b, test.b)
		}
		if ii := vv.Int(); ii != test.i {
			t.Errorf(`Error tests[%d] Int returns %v; want %v`, i, ii, test.i)
		}
		if i64 := vv.Int64(); i64 != test.i64 {
			t.Errorf(`Error tests[%d] Int64 returns %v; want %v`, i, i64, test.i64)
		}
		if f64 := vv.Float64(); f64 != test.f64 {
			t.Errorf(`Error tests[%d] Float64 returns %v; want %v`, i, f64, test.f64)
		}
		if s := vv.String(); s != test.s {
			t.Errorf(`Error tests[%d] String returns %v; want %v`, i, s, test.s)
		}
	}
}

func Test_Value_Compare(t *testing.T) {
	tests := []struct {
		n    Value
		op   Operator
		v    Value
		want bool
	}{
		{StringValue("x"), EQ, nil, false},
		{StringValue("x"), EQ, StringValue("x"), true},
		{StringValue("x"), EQ, StringValue("y"), false},
		{StringValue("1"), EQ, NumberValue(1), false},
		{StringValue("x"), GT, StringValue("a"), true},
		{StringValue("x"), GT, StringValue("x"), false},
		{StringValue("x"), GT, StringValue("y"), false},
		{StringValue("x"), GE, StringValue("a"), true},
		{StringValue("x"), GE, StringValue("x"), true},
		{StringValue("x"), GE, StringValue("y"), false},
		{StringValue("x"), LT, StringValue("a"), false},
		{StringValue("x"), LT, StringValue("x"), false},
		{StringValue("x"), LT, StringValue("y"), true},
		{StringValue("x"), LE, StringValue("a"), false},
		{StringValue("x"), LE, StringValue("x"), true},
		{StringValue("x"), LE, StringValue("y"), true},
		{StringValue("x"), Operator("unknown"), StringValue("x"), false},
		{NumberValue(1), EQ, nil, false},
		{NumberValue(1), EQ, NumberValue(1), true},
		{NumberValue(1), EQ, NumberValue(0), false},
		{NumberValue(1), EQ, NumberValue(1.0), true},
		{NumberValue(1), EQ, StringValue("1"), false},
		{NumberValue(1), GT, NumberValue(0), true},
		{NumberValue(1), GT, NumberValue(1), false},
		{NumberValue(1), GT, NumberValue(2), false},
		{NumberValue(1), GE, NumberValue(0), true},
		{NumberValue(1), GE, NumberValue(1), true},
		{NumberValue(1), GE, NumberValue(2), false},
		{NumberValue(1), LT, NumberValue(0), false},
		{NumberValue(1), LT, NumberValue(1), false},
		{NumberValue(1), LT, NumberValue(2), true},
		{NumberValue(1), LE, NumberValue(0), false},
		{NumberValue(1), LE, NumberValue(1), true},
		{NumberValue(1), LE, NumberValue(2), true},
		{NumberValue(1), Operator("unknown"), NumberValue(1), false},
		{BoolValue(true), EQ, BoolValue(true), true},
		{BoolValue(true), EQ, BoolValue(false), false},
		{BoolValue(true), EQ, StringValue("true"), false},
		{BoolValue(true), LT, BoolValue(true), false},
		{BoolValue(true), GT, BoolValue(true), false},
	}
	for i, test := range tests {
		got := test.n.Compare(test.op, test.v)
		if got != test.want {
			t.Errorf(`Error tests[%d] returns %v; want %v`, i, got, test.want)
		}
	}
}
