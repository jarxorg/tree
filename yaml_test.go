package tree

import (
	"bytes"
	"log"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

func Test_MarshalYAML(t *testing.T) {
	want := `a:
- "1"
- 2
- true
- null
- null
`
	n := Map{
		"a": Array{
			StringValue("1"),
			NumberValue(2),
			BoolValue(true),
			Nil,
			nil,
		},
	}
	got, err := MarshalYAML(n)
	if err != nil {
		log.Fatal(err)
	}
	if string(got) != want {
		t.Errorf(`Error %#v marshaled %s; want %s`, n, string(got), want)
	}
}

func Test_Map_MarshalYAML(t *testing.T) {
	want := `a:
- "1"
- 2
- true
`
	n := Map{
		"a": Array{
			StringValue("1"),
			NumberValue(2),
			BoolValue(true),
		},
	}
	got, err := yaml.Marshal(n)
	if err != nil {
		log.Fatal(err)
	}
	if string(got) != want {
		t.Errorf(`Error %#v marshaled %s; want %s`, n, string(got), want)
	}
}

func Test_Array_MarshalYAML(t *testing.T) {
	want := `- "1"
- 2
- true
`
	n := Array{
		StringValue("1"),
		NumberValue(2),
		BoolValue(true),
	}
	got, err := yaml.Marshal(n)
	if err != nil {
		log.Fatal(err)
	}
	if string(got) != want {
		t.Errorf(`Error %#v marshaled %s; want %s`, n, string(got), want)
	}
}

func Test_DecodeYAML_Errors(t *testing.T) {
	tests := []struct {
		data   []byte
		errstr string
	}{
		{
			data:   []byte(`"`),
			errstr: "yaml: found unexpected end of stream",
		}, {
			data:   []byte(`}`),
			errstr: "yaml: did not find expected node content",
		}, {
			data:   []byte("{\n1"),
			errstr: `yaml: line 2: did not find expected ',' or '}'`,
		},
	}
	for i, test := range tests {
		dec := yaml.NewDecoder(bytes.NewReader(test.data))
		got, err := DecodeYAML(dec)
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

func Test_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		want Node
		data []byte
	}{
		{
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
			data: []byte(`a: 1
b: true
c: null
d: ["1",2,true]
e: {"x":"x"}
`),
		}, {
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
			data: []byte(`- "1"
- 2
- true
- null
- {"a":1,"b":true,"c":null}
- ["x"]
`),
		},
	}
	for _, test := range tests {
		got, err := UnmarshalYAML(test.data)
		if err != nil {
			log.Fatal(err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf(`Error unmarshaled %#v; want %#v`, got, test.want)
		}
	}
}

func Test_Map_UnmarshalYAML(t *testing.T) {
	want := Map{
		"a": NumberValue(1),
		"b": BoolValue(true),
		"c": nil,
	}
	data := []byte(`a: 1
b: true
c: null
`)
	var got Map
	if err := yaml.Unmarshal(data, &got); err != nil {
		log.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf(`Error unmarshaled %#v; want %#v`, got, want)
	}
}

func Test_Array_UnmarshalYAML(t *testing.T) {
	want := Array{
		StringValue("1"),
		NumberValue(2),
		BoolValue(true),
	}
	data := []byte(`- "1"
- 2
- true
`)
	got := Array{}
	if err := yaml.Unmarshal(data, &got); err != nil {
		log.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf(`Error unmarshaled %#v; want %#v`, got, want)
	}
}
