# Tree

[![PkgGoDev](https://pkg.go.dev/badge/github.com/jarxorg/tree)](https://pkg.go.dev/github.com/jarxorg/tree)
[![Report Card](https://goreportcard.com/badge/github.com/jarxorg/tree)](https://goreportcard.com/report/github.com/jarxorg/tree)

Tree is a simple structure for dealing with dynamic or unknown JSON/YAML structures in Go.

## Formats

### Go

```go
tree.Map{
	"ID":     tree.ToValue(1),
	"Name":   tree.ToValue("Reds"),
	"Colors": tree.ToArrayValues("Crimson", "Red", "Ruby", "Maroon"),
}
```

### JSON

```json
{
	"ID": 1,
	"Name": "Reds",
	"Colors": ["Crimson", "Red", "Ruby", "Maroon"]
}
```

### YAML

```yaml
ID: 1
Name: Reds
Colors:
- Crimson
- Red
- Ruby
- Maroon
```

## Marshal and Unmarshal

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
		"Colors": tree.ToArrayValues("Crimson", "Red", "Ruby", "Maroon"),
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

## Query

| Tree Query | Results |
| - | - |
| .store.book[0] | {"category": "reference", "author": "Nigel Rees", "title": "Sayings of the Century", "price": 8.95} |
| .store.book[0].price | 8.95 |
| .store.book[1:2].price | [12.99, 8.99] |
| .store.book[].author | ["Nigel Rees", "Evelyn Waugh", "Herman Melville", "J. R. R. Tolkien"] |
| .store.book[.category=="fiction" and .price < 10].title | ["Moby Dick"] |
| .store.book[.authors[0] == "Nigel Rees"].title | ["Sayings of the Century"] |

### Illustrative Object

```json
{
  "store": {
    "book": [{
        "category": "reference",
        "author": "Nigel Rees",
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
```

## tq - Command line tool

### Install

```sh
go install github.com/jarxorg/tree/cmd/tq@latest
```

### Usage

```
tq is a portable command-line JSON/YAML processor.

Usage:
  tq [flags] [query]

Flags:
  -h	help for tq
  -i value
    	input format (json or yaml) (default "json") (default json)
  -o value
    	output format (json or yaml) (default "json") (default json)

Examples:
  % echo '{"colors": ["red", "green", "blue"]}' | tq '.colors[0]'
  "red"
```
