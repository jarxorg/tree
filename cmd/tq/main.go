package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jarxorg/tree"
	"gopkg.in/yaml.v2"
)

const (
	cmd          = "tq"
	desc         = cmd + " is a portable command-line JSON/YAML processor."
	usage        = cmd + " [flags] [query]"
	examplesText = `Examples:
  % echo '{"colors": ["red", "green", "blue"]}' | tq '.colors[0]'
  "red"
`
)

type format string

func (f *format) String() string {
	return string(*f)
}

func (f *format) Set(value string) error {
	switch value {
	case "json":
		*f = "json"
		return nil
	case "yaml":
		*f = "yaml"
		return nil
	}
	return fmt.Errorf("unknown format")
}

var (
	inputFormat  = format("json")
	outputFormat = format("json")
)

func init() {
	flag.Bool("h", false, "help for "+cmd)
	flag.Var(&inputFormat, "i", `input format (json or yaml) (default "json")`)
	flag.Var(&outputFormat, "o", `output format (json or yaml) (default "json")`)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\nUsage:\n  %s\n\n", desc, usage)
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n%s", examplesText)
	}
}

func main() {
	flag.Parse()
	handleError(run())
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	switch inputFormat {
	case "yaml":
		return runYAML()
	}
	return runJSON()
}

func runJSON() error {
	i := 0
	dec := json.NewDecoder(os.Stdin)
	for dec.More() {
		n, err := tree.DecodeJSON(dec)
		if err != nil {
			return err
		}
		if err := evaluate(n, i); err != nil {
			return err
		}
		i++
	}
	return nil
}

func runYAML() error {
	i := 0
	dec := yaml.NewDecoder(os.Stdin)
	for {
		n, err := tree.DecodeYAML(dec)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if err := evaluate(n, i); err != nil {
			return err
		}
		i++
	}
	return nil
}

func evaluate(node tree.Node, i int) error {
	node, err := tree.Find(node, flag.Arg(0))
	if err != nil {
		return err
	}
	switch outputFormat {
	case "yaml":
		out, err := tree.MarshalYAML(node)
		if err != nil {
			return err
		}
		if i > 0 {
			fmt.Println("---")
		}
		fmt.Print(string(out))
		return nil
	}
	out, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
