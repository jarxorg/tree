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
  red
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
	isHelp       = false
	inputFormat  = format("json")
	outputFormat = format("json")
	isExpand     = false
)

func init() {
	flag.BoolVar(&isHelp, "h", false, "help for "+cmd)
	flag.Var(&inputFormat, "i", `input format (json or yaml) (default "json")`)
	flag.Var(&outputFormat, "o", `output format (json or yaml) (default "json")`)
	flag.BoolVar(&isExpand, "x", false, "expand results")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\nUsage:\n  %s\n\n", desc, usage)
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n%s", examplesText)
	}
}

func main() {
	flag.Parse()
	if isHelp {
		flag.Usage()
		return
	}
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
	dec := json.NewDecoder(os.Stdin)
	for dec.More() {
		n, err := tree.DecodeJSON(dec)
		if err != nil {
			return err
		}
		if err := evaluate(n); err != nil {
			return err
		}
	}
	return nil
}

func runYAML() error {
	dec := yaml.NewDecoder(os.Stdin)
	for {
		n, err := tree.DecodeYAML(dec)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if err := evaluate(n); err != nil {
			return err
		}
	}
	return nil
}

func evaluate(node tree.Node) error {
	node, err := tree.Find(node, flag.Arg(0))
	if err != nil {
		return err
	}
	if node == nil {
		return nil
	}
	if isExpand {
		return node.Each(func(_ interface{}, v tree.Node) error {
			return output(v)
		})
	}
	return output(node)
}

func output(node tree.Node) error {
	if node.Type().IsValue() {
		fmt.Println(node.Value().String())
		return nil
	}
	switch outputFormat {
	case "yaml":
		return outputYAML(node)
	}
	return outputJSON(node)
}

var outputYAMLCalled = 0

func outputYAML(node tree.Node) error {
	if outputYAMLCalled > 0 {
		fmt.Println("---")
	}
	out, err := tree.MarshalYAML(node)
	if err != nil {
		return err
	}
	fmt.Print(string(out))
	outputYAMLCalled++
	return nil
}

func outputJSON(node tree.Node) error {
	out, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
