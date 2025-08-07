package tree

import (
	"reflect"
	"testing"
)

func Test_MethodQuery(t *testing.T) {
	testCases := []struct {
		caseName string
		q        Query
		n        Node
		want     []Node
		errstr   string
	}{
		{
			caseName: "count map",
			q:        &CountQuery{},
			n:        Map{"key1": ToValue(1), "key2": ToValue("a")},
			want:     []Node{ToValue(2)},
		}, {
			caseName: "count array",
			q:        &CountQuery{},
			n:        ToArrayValues(1, 2, 3),
			want:     []Node{ToValue(3)},
		}, {
			caseName: "keys map",
			q:        &KeysQuery{},
			n:        Map{"key1": ToValue(1), "key2": ToValue("a")},
			want:     []Node{ToArrayValues("key1", "key2")},
		}, {
			caseName: "keys array",
			q:        &KeysQuery{},
			n:        ToArrayValues(1, 2, 3),
			want:     []Node{ToArrayValues(0, 1, 2)},
		}, {
			caseName: "values map",
			q:        &ValuesQuery{},
			n:        Map{"key1": ToValue(1), "key2": ToValue("a")},
			want:     []Node{ToArrayValues(1, "a")},
		}, {
			caseName: "values array",
			q:        &ValuesQuery{},
			n:        ToArrayValues(1, 2, 3),
			want:     []Node{ToArrayValues(1, 2, 3)},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.caseName, func(t *testing.T) {
			got, err := tc.q.Exec(tc.n)
			if tc.errstr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tc.errstr)
				}
				if err.Error() != tc.errstr {
					t.Errorf("got error %q; want %q", err.Error(), tc.errstr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v; want %v", got, tc.want)
			}
		})
	}
}
