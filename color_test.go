package tree

import (
	"bytes"
	"testing"
)

func TestOutputColorJSON(t *testing.T) {
	tests := []struct {
		n    Node
		want string
	}{
		{
			n: Map{
				"num":  ToValue(1),
				"str":  ToValue("2"),
				"bool": ToValue(true),
				"null": Nil,
			},
			want: "{\n  \x1b[1;34m\"bool\"\x1b[0m: true,\n  \x1b[1;34m\"null\"\x1b[0m: \x1b[1;30mnull\x1b[0m,\n  \x1b[1;34m\"num\"\x1b[0m: 1,\n  \x1b[1;34m\"str\"\x1b[0m: \x1b[0;32m\"2\"\x1b[0m\n}\n",
		},
	}
	for i, test := range tests {
		out := new(bytes.Buffer)
		err := OutputColorJSON(out, test.n)
		if err != nil {
			t.Fatalf("tests[%d] %v", i, err)
		}
		if got := out.String(); got != test.want {
			t.Errorf("tests[%d] got %q; want %q\n%s", i, got, test.want, test.want)
		}
	}
}

func TestEncodeJSON(t *testing.T) {
	tests := []struct {
		e    *ColorEncoder
		n    Node
		want string
	}{
		{
			e: &ColorEncoder{IndentSize: 4, NoColor: true},
			n: Map{
				"a": ToValue(1),
				"b": Array{
					ToValue("2"),
					ToValue(true),
				},
				"c": Nil,
				"d": nil,
			},
			want: `{
    "a": 1,
    "b": [
        "2",
        true
    ],
    "c": null,
    "d": null
}
`,
		}, {
			e:    &ColorEncoder{IndentSize: 2, NoColor: true},
			n:    ToValue("\"\n\r\t"),
			want: "\"\\\"\\n\\r\\t\"\n",
		},
	}
	for i, test := range tests {
		out := new(bytes.Buffer)
		test.e.Out = out
		err := test.e.EncodeJSON(test.n)
		if err != nil {
			t.Fatalf("tests[%d] %v", i, err)
		}
		if got := out.String(); got != test.want {
			t.Errorf("tests[%d] got %q; want %q\n%s", i, got, test.want, test.want)
		}
	}
}

func TestOutputColorYAML(t *testing.T) {
	tests := []struct {
		n    Node
		want string
	}{
		{
			n: Map{
				"num":  ToValue(1),
				"str":  ToValue("2"),
				"bool": ToValue(true),
				"null": Nil,
			},
			want: "\x1b[1;34mbool\x1b[0m: true\n\x1b[1;34mnull\x1b[0m: \x1b[1;30mnull\x1b[0m\n\x1b[1;34mnum\x1b[0m: 1\n\x1b[1;34mstr\x1b[0m: \x1b[0;32m\"2\"\x1b[0m\n",
		},
	}
	for i, test := range tests {
		out := new(bytes.Buffer)
		err := OutputColorYAML(out, test.n)
		if err != nil {
			t.Fatalf("tests[%d] %v", i, err)
		}
		if got := out.String(); got != test.want {
			t.Errorf("tests[%d] got %q; want %q\n%s", i, got, test.want, test.want)
		}
	}
}

func TestEncodeYAML(t *testing.T) {
	tests := []struct {
		e    *ColorEncoder
		n    Node
		want string
	}{
		{
			e: &ColorEncoder{IndentSize: 2, NoColor: true},
			n: Map{
				"a": ToValue(1),
				"b": Array{
					ToValue("2"),
					ToValue(true),
				},
				"c": Nil,
				"d": nil,
			},
			want: `a: 1
b:
  - "2"
  - true
c: null
d: null
`,
		}, {
			e: &ColorEncoder{IndentSize: 2, NoColor: true},
			n: Map{
				"a": ToValue("line1\nline2\n"),
			},
			want: `a: |
  line1
  line2
`,
		}, {
			e: &ColorEncoder{IndentSize: 2, NoColor: true},
			n: Map{
				"a": ToValue("line1\nline2"),
			},
			want: `a: |-
  line1
  line2
`,
		}, {
			e: &ColorEncoder{IndentSize: 2, NoColor: true},
			n: Array{
				ToValue(1),
				Map{
					"a": ToValue(2),
					"b": ToValue(true),
				},
				Array{
					ToValue("c"),
					Nil,
				},
			},
			want: `- 1
- a: 2
  b: true
-
  - c
  - null
`,
		},
	}
	for i, test := range tests {
		out := new(bytes.Buffer)
		test.e.Out = out
		err := test.e.EncodeYAML(test.n)
		if err != nil {
			t.Fatalf("tests[%d] %v", i, err)
		}
		if got := out.String(); got != test.want {
			t.Errorf("tests[%d] got %q; want %q\n%s", i, got, test.want, test.want)
		}
	}
}
