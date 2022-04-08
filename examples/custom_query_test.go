package examples

import (
	"fmt"
	"log"

	"github.com/jarxorg/tree"
)

type evenIndexQuery struct{}

func (q *evenIndexQuery) Exec(n tree.Node) ([]tree.Node, error) {
	if !n.Type().IsArray() {
		return nil, nil
	}
	var rs []tree.Node
	for i, r := range n.Array() {
		if i%2 == 0 {
			rs = append(rs, r)
		}
	}
	return rs, nil
}

func (q *evenIndexQuery) String() string {
	return "even-index-query"
}

func ExampleCustomQuery() {
	group := tree.Array{
		tree.Map{
			"ID":     tree.ToValue(1),
			"Name":   tree.ToValue("Reds"),
			"Colors": tree.ToArrayValues("Crimson", "Red", "Ruby", "Maroon"),
		},
		tree.Map{
			"ID":     tree.ToValue(2),
			"Name":   tree.ToValue("Greens"),
			"Colors": tree.ToArrayValues("Green", "Lime", "Olive", "Teal"),
		},
		tree.Map{
			"ID":     tree.ToValue(3),
			"Name":   tree.ToValue("Blues"),
			"Colors": tree.ToArrayValues("Aqua", "Blue", "Cyan", "SkyBlue"),
		},
	}

	q := tree.FilterQuery{
		tree.SelectQuery{},
		tree.MapQuery("Colors"),
		&evenIndexQuery{},
	}

	rs, err := q.Exec(group)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", rs)

	// Output:
	// [Crimson Ruby Green Olive Aqua Cyan]
}

type primaryColorSelector struct {
}

func (s *primaryColorSelector) Matches(n tree.Node) (bool, error) {
	if !n.Type().IsStringValue() {
		return false, nil
	}
	switch n.Value().String() {
	case "Red", "Green", "Blue":
		return true, nil
	}
	return false, nil
}

func (s *primaryColorSelector) String() string {
	return "odd-index-selector"
}

func ExampleCustomSelector() {
	group := tree.Array{
		tree.Map{
			"ID":     tree.ToValue(1),
			"Name":   tree.ToValue("Reds"),
			"Colors": tree.ToArrayValues("Crimson", "Red", "Ruby", "Maroon"),
		},
		tree.Map{
			"ID":     tree.ToValue(2),
			"Name":   tree.ToValue("Greens"),
			"Colors": tree.ToArrayValues("Green", "Lime", "Olive", "Teal"),
		},
		tree.Map{
			"ID":     tree.ToValue(3),
			"Name":   tree.ToValue("Blues"),
			"Colors": tree.ToArrayValues("Aqua", "Blue", "Cyan", "SkyBlue"),
		},
	}

	q := tree.FilterQuery{
		tree.SelectQuery{},
		tree.MapQuery("Colors"),
		tree.SelectQuery{&primaryColorSelector{}},
	}

	rs, err := q.Exec(group)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", rs)

	// Output:
	// [Red Green Blue]
}
