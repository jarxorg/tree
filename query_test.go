package tree

import (
	"reflect"
	"testing"
)

func TestParseQuery(t *testing.T) {
	tests := []struct {
		expr string
		want Query
	}{
		{
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
						Comparator{Left: MapQuery("category"), Operator: EQ, Right: ValueQuery("fiction")},
						Comparator{Left: MapQuery("price"), Operator: LT, Right: ValueQuery("10")},
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
						Comparator{Left: FilterQuery{MapQuery("authors"), ArrayQuery(0)}, Operator: EQ, Right: ValueQuery("Nigel Rees")},
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

func TestQueryExec(t *testing.T) {
	n, err := UnmarshalJSON([]byte(testStoreJSON))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		q    Query
		want Node
	}{
		{
			// NOTE: .store.book[0]
			q: FilterQuery{
				MapQuery("store"),
				MapQuery("book"),
				ArrayQuery(0),
			},
			want: n.Get("store").Get("book").Get(0),
		}, {
			// NOTE: .store.book[0:2]
			q: FilterQuery{
				MapQuery("store"),
				MapQuery("book"),
				ArrayRangeQuery([]int{0, 2}),
			},
			want: Array{
				n.Get("store").Get("book").Get(0),
				n.Get("store").Get("book").Get(1),
				n.Get("store").Get("book").Get(2),
			},
		}, {
			// NOTE: .store.book[.category=="fiction" and .price < 10].title
			q: FilterQuery{
				MapQuery("store"),
				MapQuery("book"),
				SelectQuery{
					Selectors: []Selector{
						Comparator{Left: MapQuery("category"), Operator: EQ, Right: ValueQuery("fiction")},
						Comparator{Left: MapQuery("price"), Operator: LT, Right: ValueQuery(10)},
					},
				},
				MapQuery("title"),
			},
			want: ToArrayValues("Moby Dick"),
		}, {
			// NOTE: .store.book[.authors[0] == "Nigel Rees"].title
			q: FilterQuery{
				MapQuery("store"),
				MapQuery("book"),
				SelectQuery{
					Selectors: []Selector{
						Comparator{Left: FilterQuery{MapQuery("authors"), ArrayQuery(0)}, Operator: EQ, Right: ValueQuery("Nigel Rees")},
					},
				},
				MapQuery("title"),
			},
			want: ToArrayValues("Sayings of the Century"),
		},
	}

	for i, test := range tests {
		got, err := test.q.Exec(n)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Error tests[%d] returns %#v; want %#v", i, got, test.want)
		}
	}
}
