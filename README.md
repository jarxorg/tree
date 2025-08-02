# Tree

[![PkgGoDev](https://pkg.go.dev/badge/github.com/jarxorg/tree)](https://pkg.go.dev/github.com/jarxorg/tree)
[![Report Card](https://goreportcard.com/badge/github.com/jarxorg/tree)](https://goreportcard.com/report/github.com/jarxorg/tree)

Tree is a simple structure for dealing with dynamic or unknown JSON/YAML in Go.

## Features

- Parses json/yaml of unknown structure to get to nodes with fluent interface.
- Syntax similar to Go standard and map and slice.
- Find function can be specified the [Query](#query) expression.
- Edit function can be specified the [Edit](#edit) expression.
- Bundled 'tq' that is a portable command-line JSON/YAML processor.

## Road to 1.0

- Placeholders in query.
- Merging nodes.

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
func ExampleMarshalJSON() {
	group := tree.Map{
		"ID":     tree.ToValue(1),
		"Name":   tree.ToValue("Reds"),
		"Colors": tree.ToArrayValues("Crimson", "Red", "Ruby", "Maroon"),
	}
	b, err := json.Marshal(group)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

	// Output:
	// {"Colors":["Crimson","Red","Ruby","Maroon"],"ID":1,"Name":"Reds"}
}
```

```go
func ExampleUnmarshalJSON() {
	data := []byte(`[
  {"Name": "Platypus", "Order": "Monotremata"},
  {"Name": "Quoll",    "Order": "Dasyuromorphia"}
]`)

	var animals tree.Array
	err := json.Unmarshal(data, &animals)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", animals)

	// Output:
	// [map[Name:Platypus Order:Monotremata] map[Name:Quoll Order:Dasyuromorphia]]
}
```

### Using other parsers

Tree may works on other parsers those has compatible with "encoding/json" or "gopkg.in/yaml.v2". See [examples](examples) directory.

### Alternate json.RawMessage

For example, [Dynamic JSON in Go](https://eagain.net/articles/go-dynamic-json/) shows an example of using json.RawMessage.

It may be simpler to use tree.Map instead of json.RawMessage.

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jarxorg/tree"
)

const input = `
{
	"type": "sound",
	"msg": {
		"description": "dynamite",
		"authority": "the Bruce Dickinson"
	}
}
`

type Envelope struct {
	Type string
	Msg  tree.Map
}

func main() {
	env := Envelope{}
	if err := json.Unmarshal([]byte(input), &env); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", env)
	fmt.Printf("%#v\n", env.Msg.Get("description"))

	// Output:
	// main.Envelope{Type:"sound", Msg:tree.Map{"authority":"the Bruce Dickinson", "description":"dynamite"}}
	// "dynamite"
}
```

## Get

```go
func ExampleGet() {
	group := tree.Map{
		"ID":     tree.ToValue(1),
		"Name":   tree.ToValue("Reds"),
		"Colors": tree.ToArrayValues("Crimson", "Red", "Ruby", "Maroon"),
		"Nil":    nil,
	}
	fmt.Println(group.Get("Colors").Get(1))
	fmt.Println(group.Get("Colors", 2))
	fmt.Println(group.Get("Colors").Get(5).IsNil())
	fmt.Println(group.Get("Nil").IsNil())

	// Output:
	// Red
	// Ruby
	// true
	// true
}
```

## Find

```go
func ExampleFind() {
	group := tree.Map{
		"ID":     tree.ToValue(1),
		"Name":   tree.ToValue("Reds"),
		"Colors": tree.ToArrayValues("Crimson", "Red", "Ruby", "Maroon"),
	}

	rs, err := group.Find(".Colors[1:3]")
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range rs {
		fmt.Println(r)
	}

	// Output:
	// Red
	// Ruby
}
```

### Query

| Query | Description | Results |
| - | - | - |
| .store.book[0] | The first book | {"category": "reference", "author": "Nigel Rees", "title": "Sayings of the Century", "price": 8.95, "tags": [...]} |
| .store.book[0].price | The price of the first book | 8.95 |
| .store.book.0.price | The price of the first book (using dot) | 8.95 |
| .store.book[:2].price | All prices of books[0:2] (index 2 is exclusive) | 8.95, 12.99 |
| .store.book[].author | All authors of all books | "Nigel Rees", "Evelyn Waugh", "Herman Melville", "J. R. R. Tolkien" |
| ..author | All authors |  "Nigel Rees", "Evelyn Waugh", "Herman Melville", "J. R. R. Tolkien" |
| ..author \| [0] | The first author | "Nigel Rees" |
| .store.book[.tags[.name == "genre" and .value == "fiction"]].title | All titles of books tagged "fiction" | "Sword of Honour", "Moby Dick" |
| .store.book[(.category == "fiction" or .category == "reference") and .price < 10].title | All titles of books these are categoried into "fiction", "reference" and price < 10 | "Sayings of the Century", "Moby Dick" |
| .store.book[.title ~= "^S"].title | Titles beginning with "S" | "Sayings of the Century", "Sword of Honour" |
| .store.book.count() | Count books | 4 |
| .store.book[0].keys() | Sorted keys of the first book | ["author", "category", "price", "title"] |
| .store.book[0].values() | Values of the first book | ["Nigel Rees", "reference", 8.95, "Sayings of the Century"] |

#### Illustrative Object

```json
{
  "store": {
    "bicycle": {
      "color": "red",
      "price": 19.95
    },
    "book": [
      {
        "author": "Nigel Rees",
        "category": "reference",
        "price": 8.95,
        "title": "Sayings of the Century",
        "tags": [
          { "name": "genre", "value": "reference" },
          { "name": "era", "value": "20th century" },
          { "name": "theme", "value": "quotations" }
        ]
      },
      {
        "author": "Evelyn Waugh",
        "category": "fiction",
        "price": 12.99,
        "title": "Sword of Honour",
        "tags": [
          { "name": "genre", "value": "fiction" },
          { "name": "era", "value": "20th century" },
          { "name": "theme", "value": "WWII" }
        ]
      },
      {
        "author": "Herman Melville",
        "category": "fiction",
        "isbn": "0-553-21311-3",
        "price": 8.99,
        "title": "Moby Dick",
        "tags": [
          { "name": "genre", "value": "fiction" },
          { "name": "era", "value": "19th century" },
          { "name": "theme", "value": "whale hunting" }
        ]
      },
      {
        "author": "J. R. R. Tolkien",
        "category": "fiction",
        "isbn": "0-395-19395-8",
        "price": 22.99,
        "title": "The Lord of the Rings",
        "tags": [
          { "name": "genre", "value": "fantasy" },
          { "name": "era", "value": "20th century" },
          { "name": "theme", "value": "good vs evil" }
        ]
      }
    ]
  }
}
```

## Edit

```go
func ExampleEdit() {
	var group tree.Node = tree.Map{
		"ID":     tree.ToValue(1),
		"Name":   tree.ToValue("Reds"),
		"Colors": tree.ToArrayValues("Crimson", "Red", "Ruby", "Maroon"),
	}

	if err := tree.Edit(&group, ".Colors += \"Pink\""); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Append Pink to Colors:\n  %+v\n", group)

	if err := tree.Edit(&group, ".Name = \"Blue\""); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Set Blue to Name:\n  %+v\n", group)

	if err := tree.Edit(&group, ".Colors ^?"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Delete Colors:\n  %+v\n", group)

	// Output:
	// Append Pink to Colors:
	//   map[Colors:[Crimson Red Ruby Maroon Pink] ID:1 Name:Reds]
	// Set Blue to Name:
	//   map[Colors:[Crimson Red Ruby Maroon Pink] ID:1 Name:Blue]
	// Delete Colors:
	//   map[ID:1 Name:Blue]
}
```

## tq

tq is a portable command-line JSON/YAML processor.

### Installation

```sh
go install github.com/jarxorg/tree/cmd/tq@latest
```

Using Homebrew

```sh
brew tap jarxorg/tree
brew install jarxorg/tree/tq
```

Download binary

```sh
VERSION=0.8.2 GOOS=Darwin GOARCH=arm64; curl -fsSL "https://github.com/jarxorg/tree/releases/download/v${VERSION}/tree_${VERSION}_${GOOS}_${GOARCH}.tar.gz" | tar xz tq && mv tq /usr/local/bin
```

### Usage

```sh
tq is a command-line JSON/YAML processor.

Usage:
  tq [flags] [query] ([file...])

Flags:
  -c, --color                  output with colors
  -e, --edit stringArray       edit expression
  -x, --expand                 expand results
  -h, --help                   help for tq
  -U, --inplace                update files, inplace
  -i, --input-format string    input format (json or yaml)
  -j, --input-json             alias --input-format json
  -y, --input-yaml             alias --input-format yaml
  -O, --output string          output file
  -o, --output-format string   output format (json or yaml, default json)
  -J, --output-json            alias --output-format json
  -Y, --output-yaml            alias --output-format yaml
  -r, --raw                    output raw strings
  -s, --slurp                  slurp all results into an array
  -t, --template string        golang text/template string
  -v, --version                print version

Examples:
  % echo '{"colors": ["red", "green", "blue"]}' | tq '.colors[0]'
  "red"

  % echo '{"users":[{"id":1,"name":"one"},{"id":2,"name":"two"}]}' | tq -x -t '{{.id}}: {{.name}}' '.users'
  1: one
  2: two

  % echo '{}' | tq -e '.colors = ["red", "green"]' -e '.colors += "blue"' .
  {
    "colors": [
      "red",
      "green",
      "blue"
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


## Third-party library licenses

- [spf13/pflag](https://github.com/spf13/pflag/blob/master/LICENSE)
