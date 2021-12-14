package tree

import "testing"

func TestValues(t *testing.T) {
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
	for _, test := range tests {
		v := test.value
		vv := v.Value()
		if tt := v.Type(); tt&TypeValue == 0 {
			t.Errorf(`Error Type returns %v; want TypeValue`, tt)
		}
		if a := v.Array(); a != nil {
			t.Errorf(`Error Array returns %v; want nil`, a)
		}
		if m := v.Map(); m != nil {
			t.Errorf(`Error Map returns %v; want nil`, m)
		}
		if vv.(Node) != v {
			t.Errorf(`Error Value returns %v; want %v`, vv, v)
		}
		if b := vv.Bool(); b != test.b {
			t.Errorf(`Error Bool returns %v; want %v`, b, test.b)
		}
		if i := vv.Int(); i != test.i {
			t.Errorf(`Error Int returns %v; want %v`, i, test.i)
		}
		if i64 := vv.Int64(); i64 != test.i64 {
			t.Errorf(`Error Int64 returns %v; want %v`, i64, test.i64)
		}
		if f64 := vv.Float64(); f64 != test.f64 {
			t.Errorf(`Error Float64 returns %v; want %v`, f64, test.f64)
		}
		if s := vv.String(); s != test.s {
			t.Errorf(`Error String returns %v; want %v`, s, test.s)
		}
	}
}
