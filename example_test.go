package tree_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jarxorg/tree"
)

func ExampleMarshalJSON() {
	group := tree.Map{
		"ID":     tree.ToValue(1),
		"Name":   tree.ToValue("Reds"),
		"Colors": tree.ToArray("Crimson", "Red", "Ruby", "Maroon"),
	}
	b, err := json.Marshal(group)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

	// Output:
	// {"Colors":["Crimson","Red","Ruby","Maroon"],"ID":1,"Name":"Reds"}
}

func ExampleMarshalJSON_combined() {
	type ColorGroup struct {
		ID     int
		Name   string
		Colors tree.Array
	}
	group := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: tree.ToArray("Crimson", "Red", "Ruby", "Maroon"),
	}
	b, err := json.Marshal(group)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

	// Output:
	// {"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}
}

func ExampleUnmarshalJSON() {
	var jsonBlob = []byte(`[
  {"Name": "Platypus", "Order": "Monotremata"},
  {"Name": "Quoll",    "Order": "Dasyuromorphia"}
]`)

	var animals tree.Array
	err := json.Unmarshal(jsonBlob, &animals)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", animals)

	// Output:
	// [map[Name:Platypus Order:Monotremata] map[Name:Quoll Order:Dasyuromorphia]]
}

func ExampleUnmarshalJSON_combined() {
	var jsonBlob = []byte(`[
  {"Name": "Platypus", "Order": "Monotremata"},
  {"Name": "Quoll",    "Order": "Dasyuromorphia"}
]`)
	type Animal struct {
		Name  string
		Order tree.StringValue
	}
	var animals []Animal
	err := json.Unmarshal(jsonBlob, &animals)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", animals)

	// Output:
	// [{Name:Platypus Order:Monotremata} {Name:Quoll Order:Dasyuromorphia}]
}
