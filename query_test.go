package tree

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

// printToken prints token tree for debug.
func printToken(w io.Writer, t *token, depth int) {
	indent := strings.Repeat("\t", depth)
	fmt.Fprintf(w, "%s{%s} %s\n", indent, t.cmd, t.word)
	if len(t.children) > 0 {
		depth++
		for _, c := range t.children {
			printToken(w, c, depth)
		}
	}
}

func Test_Query(t *testing.T) {
	tests := []struct {
		q      Query
		n      Node
		want   Node
		errstr string
	}{
		{
			q:    NopQuery,
			n:    Array{},
			want: Array{},
		}, {
			q:    MapQuery("key"),
			n:    Map{"key": ToValue("value")},
			want: ToValue("value"),
		}, {
			q: MapQuery("key"),
		}, {
			q:      MapQuery("key"),
			n:      ToValue("not map"),
			errstr: `Cannot index array with string "key"`,
		}, {
			q:    ArrayQuery(0),
			n:    Array{ToValue(1)},
			want: ToValue(1),
		}, {
			q: ArrayQuery(0),
		}, {
			q:      ArrayQuery(0),
			n:      ToValue("not array"),
			errstr: `Cannot index array with index 0`,
		}, {
			q:    ArrayRangeQuery{0, 1},
			n:    Array{ToValue(0), ToValue(1), ToValue(2)},
			want: Array{ToValue(0), ToValue(1)},
		}, {
			q: ArrayRangeQuery{0, 1},
		}, {
			q:      ArrayRangeQuery{0, 1, 2},
			n:      Array{},
			errstr: `Invalid array range [0 1 2]`,
		}, {
			q:      ArrayRangeQuery{0, 1},
			n:      Map{},
			errstr: `Cannot index array with range 0:1`,
		}, {
			q:    FilterQuery{MapQuery("key"), ArrayQuery(0)},
			n:    Map{"key": Array{ToValue(1)}},
			want: ToValue(1),
		}, {
			q:      FilterQuery{MapQuery("key"), ArrayQuery(0)},
			n:      Map{"key": ToValue(1)},
			errstr: `Cannot index array with index 0`,
		}, {
			q: SelectQuery{Selectors: []Selector{
				Comparator{MapQuery("key"), EQ, ValueQuery{ToValue(1)}},
			}},
			n:    Array{Map{"key": ToValue(1)}, Map{"key": ToValue(2)}},
			want: Array{Map{"key": ToValue(1)}},
		}, {
			q:    SelectQuery{},
			n:    Array{Map{"key": ToValue(1)}, Map{"key": ToValue(2)}},
			want: Array{Map{"key": ToValue(1)}, Map{"key": ToValue(2)}},
		}, {
			q: SelectQuery{Selectors: []Selector{
				Comparator{MapQuery("key"), EQ, ValueQuery{ToValue(1)}},
			}},
			n: Map{},
		}, {
			q: SelectQuery{},
		}, {
			q: SelectQuery{Selectors: []Selector{
				Comparator{ArrayQuery(0), EQ, ValueQuery{ToValue(1)}},
			}},
			n:      Array{Map{"key": ToValue(1)}},
			errstr: `Cannot index array with index 0`,
		}, {
			q: SelectQuery{Selectors: []Selector{
				Comparator{ValueQuery{ToValue(1)}, EQ, ArrayQuery(0)},
			}},
			n:      Array{Map{"key": ToValue(1)}},
			errstr: `Cannot index array with index 0`,
		},
	}
	for i, test := range tests {
		got, err := test.q.Exec(test.n)
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
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf(`Error tests[%d] returns %v; want %v`, i, got, test.want)
		}
	}
}

func Test_ParseQuery(t *testing.T) {
	tests := []struct {
		expr string
		want Query
	}{
		{
			expr: `.`,
			want: NopQuery,
		}, {
			expr: `[]`,
			want: SelectQuery{},
		}, {
			expr: `.store.book[0]`,
			want: FilterQuery{
				MapQuery("store"),
				MapQuery("book"),
				ArrayQuery(0),
			},
		}, {
			expr: `."store"."book"[0]`,
			want: FilterQuery{
				MapQuery("store"),
				MapQuery("book"),
				ArrayQuery(0),
			},
		}, {
			expr: `.store.book[0:1]`,
			want: FilterQuery{
				MapQuery("store"),
				MapQuery("book"),
				ArrayRangeQuery{0, 1},
			},
		}, {
			expr: `.store.book[.category=="fiction" and .price < 10].title`,
			want: FilterQuery{
				MapQuery("store"),
				MapQuery("book"),
				SelectQuery{
					Selectors: []Selector{
						Comparator{MapQuery("category"), EQ, ValueQuery{ToValue("fiction")}},
						Comparator{MapQuery("price"), LT, ValueQuery{ToValue("10")}},
					},
				},
				MapQuery("title"),
			},
		}, {
			expr: `.store.book[.authors[0] == "Nigel Rees"]`,
			want: FilterQuery{
				MapQuery("store"),
				MapQuery("book"),
				SelectQuery{
					Selectors: []Selector{
						Comparator{FilterQuery{MapQuery("authors"), ArrayQuery(0)}, EQ, ValueQuery{ToValue("Nigel Rees")}},
					},
				},
			},
		},
	}

	for i, test := range tests {
		got, err := ParseQuery(test.expr)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf(`Error tests[%d]: "%s" returns %#v; want %#v`, i, test.expr, got, test.want)
		}
	}
}

func Test_ParseQuery_Errors(t *testing.T) {
	tests := []struct {
		expr   string
		errstr string
	}{
		{
			expr:   `<`,
			errstr: `Syntax error: invalid token <: "<"`,
		}, {
			expr:   `[`,
			errstr: `Syntax error: no right brackets: "["`,
		}, {
			expr:   `]`,
			errstr: `Syntax error: no left bracket '[': "]"`,
		}, {
			expr:   `[a]`,
			errstr: `Syntax error: invalid index: "[a]"`,
		}, {
			expr:   `[a:b]`,
			errstr: `Syntax error: invalid range: "[a:b]"`,
		}, {
			expr:   `[0:a]`,
			errstr: `Syntax error: invalid range: "[0:a]"`,
		}, {
			expr:   `[[l] == .r]`,
			errstr: `Syntax error: invalid index: "[[l] == .r]"`,
		}, {
			expr:   `[.l == [r]]`,
			errstr: `Syntax error: invalid index: "[.l == [r]]"`,
		}, {
			expr:   `.a[a]`,
			errstr: `Syntax error: invalid index: ".a[a]"`,
		},
	}
	for i, test := range tests {
		got, err := ParseQuery(test.expr)
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

// NOTE: Copy from https://github.com/stedolan/jq/wiki/For-JSONPath-users#illustrative-object
var testStoreJSON = `{
  "store": {
    "book": [{
        "category": "reference",
				"author": "Nigel Rees",
				"authors": ["Nigel Rees"],
        "title": "Sayings of the Century",
        "price": 8.95
      },
      {
        "category": "fiction",
        "author": "Evelyn Waugh",
        "title": "Sword of Honour",
        "price": 12.99
      },
      {
        "category": "fiction",
        "author": "Herman Melville",
        "title": "Moby Dick",
        "isbn": "0-553-21311-3",
        "price": 8.99
      },
      {
        "category": "fiction",
        "author": "J. R. R. Tolkien",
        "title": "The Lord of the Rings",
        "isbn": "0-395-19395-8",
        "price": 22.99
      }
    ],
    "bicycle": {
      "color": "red",
      "price": 19.95
    }
  }
}
`

func Test_Find(t *testing.T) {
	n, err := UnmarshalJSON([]byte(testStoreJSON))
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		expr string
		want Node
	}{
		{
			expr: `.store.book[0]`,
			want: n.Get("store").Get("book").Get(0),
		}, {
			expr: `.store.book[0].price`,
			want: n.Get("store").Get("book").Get(0).Get("price"),
		}, {
			expr: `.store.book[0:2]`,
			want: Array{
				n.Get("store").Get("book").Get(0),
				n.Get("store").Get("book").Get(1),
				n.Get("store").Get("book").Get(2),
			},
		}, {
			expr: `.store.book[1:2].price`,
			want: ToArrayValues(12.99, 8.99),
		}, {
			expr: `.store.book[].author`,
			want: ToArrayValues("Nigel Rees", "Evelyn Waugh", "Herman Melville", "J. R. R. Tolkien"),
		}, {
			expr: `.store.book[.category == "fiction" and .price < 10].title`,
			want: ToArrayValues("Moby Dick"),
		}, {
			expr: `.store.book[.authors[0] == "Nigel Rees"].title`,
			want: ToArrayValues("Sayings of the Century"),
		},
	}
	for i, test := range tests {
		got, err := Find(n, test.expr)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Error tests[%d] returns %#v; want %#v", i, got, test.want)
		}
	}
}
