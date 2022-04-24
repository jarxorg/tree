package tree

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func Test_MarshalJSON(t *testing.T) {
	want := `{"a":["1",2,true,null,null]}`
	n := Map{
		"a": Array{
			StringValue("1"),
			NumberValue(2),
			BoolValue(true),
			Nil,
			nil,
		},
	}
	got, err := MarshalJSON(n)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Errorf(`Error %#v marshaled %s; want %s`, n, string(got), want)
	}
}

func Test_Map_MarshalJSON(t *testing.T) {
	want := `{"a":["1",2,true]}`
	n := Map{
		"a": Array{
			StringValue("1"),
			NumberValue(2),
			BoolValue(true),
		},
	}
	got, err := json.Marshal(n)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Errorf(`Error %#v marshaled %s; want %s`, n, string(got), want)
	}
}

func Test_Array_MarshalJSON(t *testing.T) {
	want := `["1",2,true]`
	n := Array{
		StringValue("1"),
		NumberValue(2),
		BoolValue(true),
	}
	got, err := json.Marshal(n)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Errorf(`Error %#v marshaled %s; want %s`, n, string(got), want)
	}
}

func Test_DecodeJSON_Errors(t *testing.T) {
	tests := []struct {
		data   []byte
		errstr string
	}{
		{
			data:   []byte(`"`),
			errstr: "unexpected EOF",
		}, {
			data:   []byte(`}`),
			errstr: "invalid character '}' looking for beginning of value",
		}, {
			data:   []byte("{\n1"),
			errstr: `invalid character '1'`,
		},
	}
	for i, test := range tests {
		dec := json.NewDecoder(bytes.NewReader(test.data))
		got, err := DecodeJSON(dec)
		if got != nil {
			t.Errorf(`Error tests[%d] returns not nil %#v`, i, got)
		}
		if err == nil {
			t.Fatalf(`Error tests[%d] returns no error`, i)
		}
		if err.Error() != test.errstr {
			t.Errorf(`Error tests[%d] returns error %s; want %s`, i, err.Error(), test.errstr)
		}
	}
}

func Test_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		want Node
		data string
	}{
		{
			data: `{"a":1,"b":true,"c":null,"d":["1",2,true],"e":{"x":"x"}}`,
			want: Map{
				"a": NumberValue(1),
				"b": BoolValue(true),
				"c": nil,
				"d": Array{
					StringValue("1"),
					NumberValue(2),
					BoolValue(true),
				},
				"e": Map{
					"x": StringValue("x"),
				},
			},
		}, {
			data: `["1",2,true,null,{"a":1,"b":true,"c":null},["x"]]`,
			want: Array{
				StringValue("1"),
				NumberValue(2),
				BoolValue(true),
				nil,
				Map{
					"a": NumberValue(1),
					"b": BoolValue(true),
					"c": nil,
				},
				Array{
					StringValue("x"),
				},
			},
		}, {
			data: `1`,
			want: NumberValue(1),
		}, {
			data: `"str"`,
			want: StringValue("str"),
		}, {
			data: `true`,
			want: BoolValue(true),
		}, {
			data: `null`,
			want: nil,
		},
	}
	for _, test := range tests {
		got, err := UnmarshalJSON([]byte(test.data))
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf(`Error unmarshaled %#v; want %#v`, got, test.want)
		}
	}
}

func Test_Map_UnmarshalJSON(t *testing.T) {
	want := Map{
		"a": NumberValue(1),
		"b": BoolValue(true),
		"c": nil,
	}
	data := []byte(`{"a":1,"b":true,"c":null}`)
	var got Map
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf(`Error unmarshaled %#v; want %#v`, got, want)
	}
}

func Test_Array_UnmarshalJSON(t *testing.T) {
	want := Array{
		StringValue("1"),
		NumberValue(2),
		BoolValue(true),
	}
	data := []byte(`["1",2,true]`)
	got := Array{}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf(`Error unmarshaled %#v; want %#v`, got, want)
	}
}
