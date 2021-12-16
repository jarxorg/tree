# github.com/jarxorg/tree

[![PkgGoDev](https://pkg.go.dev/badge/github.com/jarxorg/tree)](https://pkg.go.dev/github.com/jarxorg/tree)
[![Report Card](https://goreportcard.com/badge/github.com/jarxorg/tree)](https://goreportcard.com/report/github.com/jarxorg/tree)

Tree is simple tree structure on go.
The node can marshal/unmarshal json and yaml that is instead of interface{}.

## Examples

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jarxorg/tree"
	"gopkg.in/yaml.v2"
)

func main() {
	group := tree.Map{
		"ID":     tree.ToValue(1),
		"Name":   tree.ToValue("Reds"),
		"Colors": tree.ToArray("Crimson", "Red", "Ruby", "Maroon"),
	}
	j, err := json.Marshal(group)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	fmt.Println(string(j))
	fmt.Println()

	y, err := yaml.Marshal(group)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(y))
	fmt.Println()

	var n tree.Map
	if err := json.Unmarshal(j, &n); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", n)

	// Output:
	//
	// {"Colors":["Crimson","Red","Ruby","Maroon"],"ID":1,"Name":"Reds"}
	//
	// Colors:
	// - Crimson
	// - Red
	// - Ruby
	// - Maroon
	// ID: 1
	// Name: Reds
	//
	// tree.Map{"Colors":tree.Array{"Crimson", "Red", "Ruby", "Maroon"}, "ID":1, "Name":"Reds"}
}
```
