# Tree

[![PkgGoDev](https://pkg.go.dev/badge/github.com/jarxorg/tree)](https://pkg.go.dev/github.com/jarxorg/tree)
[![Report Card](https://goreportcard.com/badge/github.com/jarxorg/tree)](https://goreportcard.com/report/github.com/jarxorg/tree)

Tree is a simple structure for dealing with dynamic or unknown JSON/YAML in Go.

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
	// NOTE: Get chain
	fmt.Println(group.Get("Colors").Get(1))
	fmt.Println()

	// NOTE: Output JSON
	j, err := json.Marshal(group)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(j))
	fmt.Println()

	// NOTE: Output YAML
	y, err := yaml.Marshal(group)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(y))
	fmt.Println()

	// NOTE: Unmarshal JSON
	var n tree.Map
	if err := json.Unmarshal(j, &n); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", n)
	fmt.Println()

	// NOTE: Find
	r, err := tree.Find(n, ".Colors[1]")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", r)
	fmt.Println()

	// Output:
	// Red
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
	//
	// "Red"
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
| .store.book[.title ~= "^S"].title | Titles beginning with "S" | "Sayings of the Century", "Sword of Honour" |

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
VERSION=0.7.0 GOOS=Darwin GOARCH=arm64; curl -fsSL "https://github.com/jarxorg/tree/releases/download/v${VERSION}/tree_${VERSION}_${GOOS}_${GOARCH}.tar.gz" | tar xz tq && mv tq /usr/local/bin
```

### Usage

```sh
tq is a command-line JSON/YAML processor.

Usage:
  tq [flags] [query] ([file...])

Flags:
  -e, --edit stringArray       edit expression
  -x, --expand                 expand results
  -h, --help                   help for tq
  -U, --inplace                update files, inplace
  -i, --input-format string    input format (json or yaml)
  -O, --output string          output file
  -o, --output-format string   output format (json or yaml, default json)
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
