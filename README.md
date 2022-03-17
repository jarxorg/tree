# Tree

[![PkgGoDev](https://pkg.go.dev/badge/github.com/jarxorg/tree)](https://pkg.go.dev/github.com/jarxorg/tree)
[![Report Card](https://goreportcard.com/badge/github.com/jarxorg/tree)](https://goreportcard.com/report/github.com/jarxorg/tree)

Tree is a simple structure for dealing with dynamic or unknown JSON/YAML structures in Go.

## Syntax

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
	fmt.Println(string(j))
	fmt.Println()

	y, err := yaml.Marshal(group)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(y))
	fmt.Println()

	var n tree.Map
	if err := json.Unmarshal(j, &n); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", n)
	fmt.Println()

	r, err := tree.Find(n, ".Colors[1]")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", r)
	fmt.Println()

	// Output:
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
	//
	// "Red"
}
```

### Using other parsers

Tree may works on other parsers those has compatible with "encoding/json" or "gopkg.in/yaml.v2". See [examples](examples) directory.

## Query

| Query | Description | Results |
| - | - | - |
| .store.book[0] | The first book | {"category": "reference", "author": "Nigel Rees", "title": "Sayings of the Century", "price": 8.95} |
| .store.book[0].price | The price of the first book | 8.95 |
| .store.book.0.price | The price of the first book (using dot) | 8.95 |
| .store.book[:2].price | All prices of books[0:2] (index 2 is exclusive) | 8.95, 12.99 |
| .store.book[].author | All authors of all books | "Nigel Rees", "Evelyn Waugh", "Herman Melville", "J. R. R. Tolkien" |
| ..author | All authors |  "Nigel Rees", "Evelyn Waugh", "Herman Melville", "J. R. R. Tolkien" |
| ..author \| [0] | The first author | "Nigel Rees" |
| .store.book[(.category == "fiction" or .category == "reference") and .price < 10].title | All titles of books these are categoried into "fiction", "reference" and price < 10 | "Sayings of the Century", "Moby Dick" |

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

### Installation

```sh
go install github.com/jarxorg/tree/cmd/tq@latest
```

### Usage

```
tq is a portable command-line JSON/YAML processor.

Usage:
  tq [flags] [query] ([file...])

Flags:
  -e value
    	edit expression
  -h	help for tq
  -i value
    	input format (json or yaml) (default json)
  -o value
    	output format (json or yaml) (default json)
  -r	output raw strings
  -s	slurp all results into an array
  -t string
    	golang text/template string (ignore -o flag)
  -x	expand results

Examples:
  % echo '{"colors": ["red", "green", "blue"]}' | tq '.colors[0]'
  "red"

  % echo '{"users":[{"id":1,"name":"one"},{"id":2,"name":"two"}]}' | tq -x -t '{{.id}}: {{.name}}' '.users'
  1: one
  2: two

  % echo '{}' | tq -e '.colors = ["Red", "Green"]' -e '.colors += "Blue"' .
  {
    "colors": [
      "Red",
      "Green",
      "Blue"
    ]
  }
```

### for jq user

| tq | jq |
| - | - |
| tq '.store.book[0]' | jq '.store.book[0]' |
| tq '.store.book[]' | jq '.store.book[]' |
| tq '.store.book[:2].price' | jq '.store.book[:2][] \| .price' |
| tq '.store.book[.category == "fiction" and .price < 10].title' | jq '.store.book[] \| select(.category == "fiction" and .price < 10) \| .title' |
